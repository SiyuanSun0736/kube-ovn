//go:build linux

package speaker

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	nadutils "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/utils"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	corev1 "k8s.io/api/core/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func defaultFDBLookup(mac net.HardwareAddr) (int, error) {
	neighs, err := netlink.NeighList(0, unix.AF_BRIDGE)
	if err != nil {
		return 0, fmt.Errorf("failed to list bridge FDB entries: %w", err)
	}
	for _, neigh := range neighs {
		if bytes.Equal(neigh.HardwareAddr, mac) {
			return neigh.LinkIndex, nil
		}
	}
	return 0, fmt.Errorf("no bridge FDB entry found for MAC %s", mac)
}

// getNatGwTCRedirectTarget 获取 EIP 对应的 NAT 网关 Pod 的 TC 重定向直通目标。
// 直接以 EIP 本身指定的 ExternalSubnet 为准精准定位外网网卡。
func (c *Controller) getNatGwTCRedirectTarget(pod *corev1.Pod, eip *kubeovnv1.IptablesEIP) (TCRedirectTarget, error) {
	return c.getNatGwTCRedirectTargetInternal(pod, eip, netlink.LinkByName)
}

func (c *Controller) getNatGwTCRedirectTargetInternal(pod *corev1.Pod, eip *kubeovnv1.IptablesEIP, lookup func(string) (netlink.Link, error), fdbLookup ...func(net.HardwareAddr) (int, error)) (TCRedirectTarget, error) {
	if eip.Spec.ExternalSubnet == "" {
		return TCRedirectTarget{}, fmt.Errorf("iptables-eip %s has empty externalSubnet", eip.Name)
	}
	subnet, err := c.subnetsLister.Get(eip.Spec.ExternalSubnet)
	if err != nil {
		return TCRedirectTarget{}, fmt.Errorf("failed to get external subnet %s for EIP %s: %w", eip.Spec.ExternalSubnet, eip.Name, err)
	}
	return natGwTCRedirectTarget(pod, subnet, c.config.TcBpfInterface, lookup, fdbLookup...)
}

// TC Ingress 接收网卡与外网子网的 macvlan 父设备可以不是同一个设备。
// 通过后者进行直通路由并精准重定向到 Pod 的真实硬件 MAC 地址。
func natGwTCRedirectTarget(pod *corev1.Pod, subnet *kubeovnv1.Subnet, defaultParent string, lookup func(string) (netlink.Link, error), fdbLookup ...func(net.HardwareAddr) (int, error)) (TCRedirectTarget, error) {
	if pod == nil {
		return TCRedirectTarget{}, errors.New("NAT gateway pod is missing")
	}
	nadName, nadNamespace, ok := util.GetNadBySubnetProvider(subnet.Spec.Provider)
	if !ok {
		return TCRedirectTarget{}, fmt.Errorf("subnet %s provider %q is not a valid NAD format (<name>.<namespace>)", subnet.Name, subnet.Spec.Provider)
	}
	expectedNetName := nadNamespace + "/" + nadName

	statuses, err := nadutils.GetNetworkStatus(pod)
	if err != nil {
		return TCRedirectTarget{}, fmt.Errorf("failed to get network status for Pod %s/%s: %w", pod.Namespace, pod.Name, err)
	}

	master := subnet.Annotations[util.NadMacvlanMasterAnnotation]
	if master != "" {
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
			mac, err := net.ParseMAC(status.Mac)
			if err != nil {
				return TCRedirectTarget{}, fmt.Errorf("invalid MAC for network %s: %w", status.Name, err)
			}
			target := TCRedirectTarget{TargetIfindex: link.Attrs().Index, ParentIfindex: link.Attrs().ParentIndex, TargetMAC: mac}
			if !validTCRedirectTarget(target) {
				return TCRedirectTarget{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return target, nil
		}
	} else {
		// 非 macvlan 模式（veth-pair 模式）：
		// 纯二层无 IP 依赖，通过网桥 FDB 转发表根据 Pod MAC 反查宿主机端端口
		fdbFind := defaultFDBLookup
		if len(fdbLookup) > 0 && fdbLookup[0] != nil {
			fdbFind = fdbLookup[0]
		}
		for _, status := range statuses {
			if status.Name != expectedNetName {
				continue
			}
			mac, err := net.ParseMAC(status.Mac)
			if err != nil {
				return TCRedirectTarget{}, fmt.Errorf("invalid MAC for network %s: %w", status.Name, err)
			}
			targetLinkIndex, err := fdbFind(mac)
			if err != nil {
				return TCRedirectTarget{}, fmt.Errorf("failed to resolve host veth device for pod %s/%s via bridge FDB: %w", pod.Namespace, pod.Name, err)
			}
			if defaultParent != "" {
				if plink, err := lookup(defaultParent); err == nil && targetLinkIndex == plink.Attrs().Index {
					return TCRedirectTarget{}, fmt.Errorf("bridge FDB target device points to ingress interface %s (%d), not a local pod veth", defaultParent, plink.Attrs().Index)
				}
			}
			target := TCRedirectTarget{TargetIfindex: targetLinkIndex, ParentIfindex: 0, TargetMAC: mac}
			if !validTCRedirectTarget(target) {
				return TCRedirectTarget{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return target, nil
		}
	}
	return TCRedirectTarget{}, fmt.Errorf("pod %s/%s has no network status for %s", pod.Namespace, pod.Name, expectedNetName)
}
