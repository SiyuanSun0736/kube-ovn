package util

import (
	"fmt"
	"net/netip"
	"strings"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

// NatGatewayDataplaneMode 校验并提取 VpcNatGateway 的 DataplaneMode 配置。
func NatGatewayDataplaneMode(gw *kubeovnv1.VpcNatGateway) (string, error) {
	mode := strings.TrimSpace(gw.Spec.DataplaneMode)
	if mode != "" && mode != NatDataplaneModeStatelessNft {
		return "", fmt.Errorf("gateway %s: unsupported dataplaneMode %q", gw.Name, mode)
	}
	return mode, nil
}

// StatelessNATIPv4Host 解析并标准化无状态 NAT 的 IPv4 主机地址或 /32 前缀，拒绝多对一网段。
func StatelessNATIPv4Host(value string) (string, error) {
	if ip, err := netip.ParseAddr(value); err == nil && ip.Is4() {
		return ip.String(), nil
	}
	if prefix, err := netip.ParsePrefix(value); err == nil && prefix.Addr().Is4() && prefix.Bits() == 32 {
		return prefix.Addr().String(), nil
	}
	return "", fmt.Errorf("stateless NAT requires a single IPv4 address or /32, got %q", value)
}
