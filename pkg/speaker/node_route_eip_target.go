//go:build linux

package speaker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/vishvananda/netlink"
	corev1 "k8s.io/api/core/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func (c *Controller) getNatGwTargetInfo(pod *corev1.Pod, eip *kubeovnv1.IptablesEIP) (TargetInfo, error) {
	attachmentSubnet := eip.Spec.ExternalSubnet
	if c.natgatewayLister != nil && eip.Spec.NatGwDp != "" {
		if gw, err := c.natgatewayLister.Get(eip.Spec.NatGwDp); err == nil && len(gw.Spec.ExternalSubnets) > 0 {
			attachmentSubnet = util.GetNatGwExternalNetwork(gw.Spec.ExternalSubnets)
		}
	}
	subnet, err := c.subnetsLister.Get(attachmentSubnet)
	if err != nil {
		subnet, err = c.subnetsLister.Get(eip.Spec.ExternalSubnet)
		if err != nil {
			return TargetInfo{}, err
		}
	}
	return natGwTargetInfo(pod, subnet, c.config.TcBpfInterface, netlink.LinkByName)
}

// TC Ingress 接收网卡与外网子网的 macvlan 父设备可以不是同一个设备。
// 通过后者进行直通路由并精准重定向到 Pod 的真实硬件 MAC 地址。
func natGwTargetInfo(pod *corev1.Pod, subnet *kubeovnv1.Subnet, defaultParent string, lookup func(string) (netlink.Link, error)) (TargetInfo, error) {
	if pod == nil {
		return TargetInfo{}, errors.New("NAT gateway pod is missing")
	}
	nadName, nadNamespace, ok := util.GetNadBySubnetProvider(subnet.Spec.Provider)
	if !ok {
		nadName, nadNamespace = subnet.Name, pod.Namespace
	}
	var statuses []nadv1.NetworkStatus
	if err := json.Unmarshal([]byte(pod.Annotations[nadv1.NetworkStatusAnnot]), &statuses); err != nil {
		return TargetInfo{}, fmt.Errorf("parse network status for Pod %s/%s: %w", pod.Namespace, pod.Name, err)
	}

	master := subnet.Annotations[util.NadMacvlanMasterAnnotation]
	if master != "" {
		name, err := util.GenMacvlanIfaceName(master)
		if err != nil {
			return TargetInfo{}, err
		}
		link, err := lookup(name)
		if err != nil {
			return TargetInfo{}, err
		}
		macvlan, ok := link.(*netlink.Macvlan)
		if !ok || macvlan.Mode != netlink.MACVLAN_MODE_BRIDGE {
			return TargetInfo{}, fmt.Errorf("interface %s is not a bridge-mode macvlan", name)
		}
		for _, status := range statuses {
			networkName := status.Name
			if !strings.Contains(networkName, "/") {
				networkName = pod.Namespace + "/" + networkName
			}
			if networkName != nadNamespace+"/"+nadName {
				continue
			}
			mac, err := net.ParseMAC(status.Mac)
			if err != nil {
				return TargetInfo{}, fmt.Errorf("invalid MAC for network %s: %w", networkName, err)
			}
			info := TargetInfo{TargetIfindex: link.Attrs().Index, ParentIfindex: link.Attrs().ParentIndex, TargetMAC: mac}
			if !validTarget(info) {
				return TargetInfo{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return info, nil
		}
	} else {
		// 非 macvlan 模式（例如 ptp / veth-pair 模式）：
		// 通过查询 Pod 外部 IP 的宿主机路由反查目标宿主机端 veth 设备。
		for _, status := range statuses {
			networkName := status.Name
			if !strings.Contains(networkName, "/") {
				networkName = pod.Namespace + "/" + networkName
			}
			if networkName != nadNamespace+"/"+nadName {
				continue
			}
			if len(status.IPs) == 0 {
				continue
			}
			podIP := net.ParseIP(status.IPs[0])
			if podIP == nil {
				continue
			}
			routes, err := netlink.RouteGet(podIP)
			if err != nil || len(routes) == 0 {
				return TargetInfo{}, fmt.Errorf("failed to get route for pod IP %s: %w", podIP, err)
			}
			// 拒绝经由外部网关的路由（例如默认路由或间接路由）。
			// 本地 Pod 的直连 veth 路由必须是直连路由，不能携带下一跳网关。
			if len(routes[0].Gw) > 0 && !routes[0].Gw.IsUnspecified() {
				return TargetInfo{}, fmt.Errorf("route for pod IP %s is not a direct route (has gateway %s)", podIP, routes[0].Gw)
			}
			targetLinkIndex := routes[0].LinkIndex
			mac, err := net.ParseMAC(status.Mac)
			if err != nil {
				return TargetInfo{}, fmt.Errorf("invalid MAC for network %s: %w", networkName, err)
			}
			parentIndex := 0
			if defaultParent != "" {
				if plink, err := lookup(defaultParent); err == nil {
					parentIndex = plink.Attrs().Index
				}
			}
			if parentIndex > 0 && targetLinkIndex == parentIndex {
				return TargetInfo{}, fmt.Errorf("route for pod IP %s points to ingress/parent interface %d, not a local pod veth", podIP, parentIndex)
			}
			info := TargetInfo{TargetIfindex: targetLinkIndex, ParentIfindex: parentIndex, TargetMAC: mac}
			if !validTarget(info) {
				return TargetInfo{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return info, nil
		}
	}
	return TargetInfo{}, fmt.Errorf("pod %s/%s has no network status for %s/%s", pod.Namespace, pod.Name, nadNamespace, nadName)
}
