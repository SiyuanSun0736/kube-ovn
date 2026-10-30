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
	subnet, err := c.subnetsLister.Get(eip.Spec.ExternalSubnet)
	if err != nil {
		return TargetInfo{}, err
	}
	return natGwTargetInfo(pod, subnet, c.config.TcBpfInterface, netlink.LinkByName)
}

// The TC ingress device and the external subnet's macvlan master need not be
// the same device. Route through the latter and address the Pod's actual MAC.
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
		// Non-macvlan mode (e.g. ptp / veth-pair mode):
		// Lookup the target veth device via the Pod's external IP route.
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
				return TargetInfo{}, fmt.Errorf("failed to get route for pod IP %s: %v", podIP, err)
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
			info := TargetInfo{TargetIfindex: targetLinkIndex, ParentIfindex: parentIndex, TargetMAC: mac}
			if !validTarget(info) {
				return TargetInfo{}, fmt.Errorf("invalid target for Pod %s/%s", pod.Namespace, pod.Name)
			}
			return info, nil
		}
	}
	return TargetInfo{}, fmt.Errorf("pod %s/%s has no network status for %s/%s", pod.Namespace, pod.Name, nadNamespace, nadName)
}
