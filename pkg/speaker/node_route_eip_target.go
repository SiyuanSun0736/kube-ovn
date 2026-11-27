//go:build linux

package speaker

import (
	"errors"
	"fmt"
	"net"

	nadutils "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/utils"
	"github.com/vishvananda/netlink"
	corev1 "k8s.io/api/core/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// getNatGwTCRedirectTarget 获取 EIP 对应的 NAT 网关 Pod 的 TC 重定向直通目标。
// 直接以 EIP 本身指定的 ExternalSubnet 为准精准定位外网网卡。
func (c *Controller) getNatGwTCRedirectTarget(pod *corev1.Pod, eip *kubeovnv1.IptablesEIP) (TCRedirectTarget, error) {
	return c.getNatGwTCRedirectTargetInternal(pod, eip, netlink.LinkByName, netlink.LinkByAlias)
}

func (c *Controller) getNatGwTCRedirectTargetInternal(pod *corev1.Pod, eip *kubeovnv1.IptablesEIP, lookup, lookupByAlias func(string) (netlink.Link, error)) (TCRedirectTarget, error) {
	if eip.Spec.ExternalSubnet == "" {
		return TCRedirectTarget{}, fmt.Errorf("iptables-eip %s has empty externalSubnet", eip.Name)
	}
	subnet, err := c.subnetsLister.Get(eip.Spec.ExternalSubnet)
	if err != nil {
		return TCRedirectTarget{}, fmt.Errorf("failed to get external subnet %s for EIP %s: %w", eip.Spec.ExternalSubnet, eip.Name, err)
	}
	return natGwTCRedirectTarget(pod, subnet, c.config.TcBpfInterface, lookup, lookupByAlias)
}

// TC Ingress 接收网卡与外网子网的 macvlan 父设备可以不是同一个设备。
// 通过后者进行直通路由并精准重定向到 Pod 的真实硬件 MAC 地址。
func natGwTCRedirectTarget(pod *corev1.Pod, subnet *kubeovnv1.Subnet, defaultParent string, lookup, lookupByAlias func(string) (netlink.Link, error)) (TCRedirectTarget, error) {
	if pod == nil {
		return TCRedirectTarget{}, errors.New("NAT gateway pod is missing")
	}
	nadName, nadNamespace, ok := util.GetNadBySubnetProvider(subnet.Spec.Provider)
	if !ok {
		return TCRedirectTarget{}, fmt.Errorf("subnet %s provider %q is not a valid NAD format (<name>.<namespace>)", subnet.Name, subnet.Spec.Provider)
	}
	expectedNetName := nadNamespace + "/" + nadName

	statuses, statusErr := nadutils.GetNetworkStatus(pod)

	master := subnet.Annotations[util.NadMacvlanMasterAnnotation]
	if master != "" {
		if statusErr != nil {
			return TCRedirectTarget{}, fmt.Errorf("failed to get network status for Pod %s/%s: %w", pod.Namespace, pod.Name, statusErr)
		}
		name, err := util.GenMacvlanIfaceName(master)
		if err != nil {
			return TCRedirectTarget{}, err
		}
		link, err := lookup(name)
		if err != nil {
			return TCRedirectTarget{}, err
		}
		macvlan, ok := link.(*netlink.Macvlan)
		if !ok || macvlan.Mode != netlink.MACVLAN_MODE_BRIDGE {
			return TCRedirectTarget{}, fmt.Errorf("interface %s is not a bridge-mode macvlan", name)
		}
		for _, status := range statuses {
			if status.Name != expectedNetName {
				continue
			}
			gwNet1Mac, err := net.ParseMAC(status.Mac)
			if err != nil {
				return TCRedirectTarget{}, fmt.Errorf("invalid MAC for network %s: %w", status.Name, err)
			}
			gwRedirectTarget := TCRedirectTarget{TargetIfindex: link.Attrs().Index, ParentIfindex: link.Attrs().ParentIndex, TargetMAC: gwNet1Mac}
			if !validTCRedirectTarget(gwRedirectTarget) {
				return TCRedirectTarget{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return gwRedirectTarget, nil
		}
	} else {
		// 非 macvlan 模式（例如 OVS Underlay veth-pair 或 ptp 模式）：
		// 1. 优先通过 CNI 写入宿主端 veth 的确定性 ifalias 直接定位（不依赖宿主机路由，也不依赖 Multus network-status / IPAM）；
		// 2. 未命中 alias 时回退至通过 Pod 外部 IP 的宿主机路由反查宿主机端 veth 设备。
		var alias string
		var aliasErr error
		if lookupByAlias != nil {
			ifName := "net1"
			macStr := pod.Annotations[fmt.Sprintf(util.MacAddressAnnotationTemplate, subnet.Spec.Provider)]
			if statusErr == nil {
				for _, status := range statuses {
					if status.Name != expectedNetName {
						continue
					}
					if status.Interface != "" {
						ifName = status.Interface
					}
					if status.Mac != "" {
						macStr = status.Mac
					}
					break
				}
			}
			if macStr != "" {
				gwNet1Mac, err := net.ParseMAC(macStr)
				if err != nil {
					return TCRedirectTarget{}, fmt.Errorf("invalid MAC for network %s: %w", expectedNetName, err)
				}
				alias = util.GenNatGwHostVethAlias(pod.Namespace, pod.Name, string(pod.UID), ifName)
				var hostLink netlink.Link
				hostLink, aliasErr = lookupByAlias(alias)
				if aliasErr == nil && hostLink != nil {
					gwRedirectTarget := TCRedirectTarget{TargetIfindex: hostLink.Attrs().Index, ParentIfindex: 0, TargetMAC: gwNet1Mac}
					if !validTCRedirectTarget(gwRedirectTarget) {
						return TCRedirectTarget{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
					}
					return gwRedirectTarget, nil
				}
			}
		}
		if statusErr != nil {
			if aliasErr != nil {
				return TCRedirectTarget{}, fmt.Errorf("failed to locate host veth by alias %q for pod %s/%s (%s): %w (network status error: %w)", alias, pod.Namespace, pod.Name, expectedNetName, aliasErr, statusErr)
			}
			return TCRedirectTarget{}, fmt.Errorf("failed to get network status for Pod %s/%s: %w", pod.Namespace, pod.Name, statusErr)
		}
		for _, status := range statuses {
			if status.Name != expectedNetName {
				continue
			}
			if len(status.IPs) == 0 {
				continue
			}
			gwNet1Mac, err := net.ParseMAC(status.Mac)
			if err != nil {
				return TCRedirectTarget{}, fmt.Errorf("invalid MAC for network %s: %w", status.Name, err)
			}
			podIP := net.ParseIP(status.IPs[0])
			if podIP == nil {
				continue
			}
			routes, err := netlink.RouteGet(podIP)
			if err != nil || len(routes) == 0 {
				return TCRedirectTarget{}, fmt.Errorf("failed to get route for pod IP %s: %w", podIP, err)
			}
			// 拒绝经由外部网关的路由（例如默认路由或间接路由）。
			// 本地 Pod 的直连 veth 路由必须是直连路由，不能携带下一跳网关。
			if len(routes[0].Gw) > 0 && !routes[0].Gw.IsUnspecified() {
				return TCRedirectTarget{}, fmt.Errorf("route for pod IP %s is not a direct route (has gateway %s)", podIP, routes[0].Gw)
			}
			gwHostIfindex := routes[0].LinkIndex
			if defaultParent != "" {
				if plink, err := lookup(defaultParent); err == nil && gwHostIfindex == plink.Attrs().Index {
					return TCRedirectTarget{}, fmt.Errorf("route for pod IP %s points to ingress interface %s (%d), not a local pod veth", podIP, defaultParent, plink.Attrs().Index)
				}
			}
			// veth 设备无父设备（ParentIfindex 为 0）
			gwRedirectTarget := TCRedirectTarget{TargetIfindex: gwHostIfindex, ParentIfindex: 0, TargetMAC: gwNet1Mac}
			if !validTCRedirectTarget(gwRedirectTarget) {
				return TCRedirectTarget{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return gwRedirectTarget, nil
		}
		if aliasErr != nil {
			return TCRedirectTarget{}, fmt.Errorf("failed to locate host veth by alias %q for pod %s/%s (%s): %w", alias, pod.Namespace, pod.Name, expectedNetName, aliasErr)
		}
	}
	return TCRedirectTarget{}, fmt.Errorf("pod %s/%s has no network status for %s", pod.Namespace, pod.Name, expectedNetName)
}
