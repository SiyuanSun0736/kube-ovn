package util

import (
	"fmt"
	"net/netip"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

// NatGatewayDataplaneMode validates and extracts the DataplaneMode from a VpcNatGateway.
// It returns an error if an unknown or unsupported mode is specified.
func NatGatewayDataplaneMode(gw *kubeovnv1.VpcNatGateway) (string, error) {
	mode := gw.Spec.DataplaneMode
	if mode != "" && mode != NatDataplaneModeStatelessNft {
		return "", fmt.Errorf("gateway %s: unsupported dataplaneMode %q", gw.Name, mode)
	}
	return mode, nil
}

// StatelessNATIPv4Host parses and normalizes a candidate IPv4 host address or /32 prefix.
// It rejects non-IPv4 addresses, CIDRs broader than /32, and invalid strings, because
// stateless 1:1 NAT and SNAT require a deterministic single private IP without conntrack.
func StatelessNATIPv4Host(value string) (string, error) {
	if ip, err := netip.ParseAddr(value); err == nil && ip.Is4() {
		return ip.String(), nil
	}
	if prefix, err := netip.ParsePrefix(value); err == nil && prefix.Addr().Is4() && prefix.Bits() == 32 {
		return prefix.Addr().String(), nil
	}
	return "", fmt.Errorf("stateless NAT requires a single IPv4 address or /32, got %q", value)
}
