package util

import (
	"fmt"
	"net/netip"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

// NatGatewayDataplaneMode resolves the configured API field.
func NatGatewayDataplaneMode(gw *kubeovnv1.VpcNatGateway) (string, error) {
	mode := gw.Spec.DataplaneMode
	if mode != "" && mode != NatDataplaneModeStatelessNft {
		return "", fmt.Errorf("gateway %s: unsupported dataplaneMode %q", gw.Name, mode)
	}
	return mode, nil
}

// StatelessNATIPv4Host rejects many-to-one SNAT: without conntrack there is no
// way to reconstruct the private destination from a shared public address.
func StatelessNATIPv4Host(value string) (string, error) {
	if ip, err := netip.ParseAddr(value); err == nil && ip.Is4() {
		return ip.String(), nil
	}
	if prefix, err := netip.ParsePrefix(value); err == nil && prefix.Addr().Is4() && prefix.Bits() == 32 {
		return prefix.Addr().String(), nil
	}
	return "", fmt.Errorf("stateless NAT requires a single IPv4 address or /32, got %q", value)
}
