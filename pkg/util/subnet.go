package util

import (
	"fmt"
	"strconv"
	"strings"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

// MaxTunnelKey is OVN's maximum 24-bit Datapath_Binding.tunnel_key value.
const MaxTunnelKey = 1<<24 - 1

func IsValidTunnelKey(key int) bool {
	return key > 0 && key <= MaxTunnelKey
}

// IsOvnVpcSubnet reports whether a subnet participates in the VPC tunnel-key
// guarantee. It includes default and custom VPCs; vlan is the only
// VPC/underlay discriminator after the OVN-provider check.
func IsOvnVpcSubnet(subnet *kubeovnv1.Subnet) bool {
	return subnet != nil && IsOvnProvider(subnet.Spec.Provider) && subnet.Spec.Vlan == ""
}

// IsTunnelKeyAnnotationValidForSubnet reports whether a provider annotation
// carries the exact valid tunnel key of the VPC (logical router) owning a
// non-vlan OVN VPC subnet. Other subnet types do not require the annotation
// and therefore return true.
func IsTunnelKeyAnnotationValidForSubnet(annotations map[string]string, provider string, subnet *kubeovnv1.Subnet, vpcTunnelKey int) bool {
	if !IsOvnVpcSubnet(subnet) {
		return true
	}
	tunnelKey, err := strconv.Atoi(annotations[fmt.Sprintf(TunnelKeyAnnotationTemplate, provider)])
	return err == nil && IsValidTunnelKey(tunnelKey) && tunnelKey == vpcTunnelKey
}

func IsOvnProvider(provider string) bool {
	if provider == "" || provider == OvnProvider {
		return true
	}
	if fields := strings.Split(provider, "."); len(fields) == 3 && fields[2] == OvnProvider {
		return true
	}
	return false
}

func GetNadBySubnetProvider(provider string) (nadName, nadNamespace string, existNad bool) {
	fields := strings.Split(provider, ".")
	switch {
	case len(fields) == 3 && fields[2] == OvnProvider:
		return fields[0], fields[1], true
	case len(fields) == 2:
		return fields[0], fields[1], true
	}
	return "", "", false
}
