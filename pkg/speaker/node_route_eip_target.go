//go:build linux

package speaker

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	nadutils "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/utils"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	corev1 "k8s.io/api/core/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// isPodVethName 校验设备名称是否属于 Pod 在宿主机端的 host 端 veth 网卡设备：
// 1. 以 "_net1_h" 结尾：Multus 辅助网卡（如网关业务直通外网口）在宿主机端的配对端；
// 2. 以 "_h" 结尾：Kube-OVN/PTP 等 CNI 默认在宿主机端的 Pod 配对端；
// 3. 以 "veth" 开头：通用或旧版 CNI 创建的宿主机端 veth 设备。
// 用于在遍历 Linux 网桥 FDB 转发表时排除底层物理网卡（如 ens7f0）、Bond 聚合口（如 bond0）等非 Pod 设备。
func isPodVethName(name string) bool {
	return strings.HasSuffix(name, "_net1_h") ||
		strings.HasSuffix(name, "_h") ||
		strings.HasPrefix(name, "veth")
}

// fdbLookupInternal 通过 Linux 内核网桥 FDB 转发表（AF_BRIDGE），仅凭 Pod MAC 地址反查其在宿主机端挂载的 veth 接口 ifindex：
// 1. 调用 netlink.NeighList(0, unix.AF_BRIDGE) 查询宿主机所有网桥接口的 FDB 记录（等价于 `bridge fdb show`）；
// 2. 校验表项对应的设备名必须满足 isPodVethName（_net1_h / _h / veth...），直接排除 bond0、物理网卡等非 Pod 设备；
// 3. 在符合规则的 Pod veth 接口表项中比对 HardwareAddr 是否与 Pod 的 MAC 精确一致。
func fdbLookupInternal(mac net.HardwareAddr, neighList func(int, int) ([]netlink.Neigh, error), linkByIndex func(int) (netlink.Link, error)) (int, error) {
	neighs, err := neighList(0, unix.AF_BRIDGE)
	if err != nil {
		return 0, fmt.Errorf("failed to list bridge FDB entries: %w", err)
	}
	for _, neigh := range neighs {
		link, err := linkByIndex(neigh.LinkIndex)
		if err != nil {
			continue
		}
		// 1. 先判断网卡名称是否符合 Pod veth 规则
		if !isPodVethName(link.Attrs().Name) {
			continue
		}
		// 2. 规则符合后再比对 MAC 地址
		if bytes.Equal(neigh.HardwareAddr, mac) {
			return neigh.LinkIndex, nil
		}
	}
	return 0, fmt.Errorf("no bridge FDB entry found on pod veth interfaces for MAC %s", mac)
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
		// 纯二层无 IP 依赖，彻底消除对 Pod IP (status.IPs[0]) 和路由表 (netlink.RouteGet) 的依赖。
		// 1. 直接以 Pod 的 network-status 中的 status.Mac 为准；
		// 2. 通过宿主机网桥 FDB 转发表（AF_BRIDGE）反查接入该网桥的 host 端 veth 网卡设备；
		// 3. 严格限定匹配到的设备名属于 Pod host 端 veth（_net1_h / _h / veth...），防止误匹配到物理口/Bond 口。
		fdbFind := func(m net.HardwareAddr) (int, error) {
			return fdbLookupInternal(m, netlink.NeighList, netlink.LinkByIndex)
		}
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
			// 防环回安全校验：重定向目标网卡不能是 Ingress 接收父设备自身（例如 bond0）
			if defaultParent != "" {
				if plink, err := lookup(defaultParent); err == nil && targetLinkIndex == plink.Attrs().Index {
					return TCRedirectTarget{}, fmt.Errorf("bridge FDB target device points to ingress interface %s (%d), not a local pod veth", defaultParent, plink.Attrs().Index)
				}
			}
			// veth-pair 点对点模式直通 Pod，ParentIfindex 为 0
			target := TCRedirectTarget{TargetIfindex: targetLinkIndex, ParentIfindex: 0, TargetMAC: mac}
			if !validTCRedirectTarget(target) {
				return TCRedirectTarget{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return target, nil
		}
	}
	return TCRedirectTarget{}, fmt.Errorf("pod %s/%s has no network status for %s", pod.Namespace, pod.Name, expectedNetName)
}
