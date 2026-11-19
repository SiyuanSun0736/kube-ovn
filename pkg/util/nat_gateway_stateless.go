package util

import (
	"fmt"
	"net/netip"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

// NatGatewayDataplaneMode 校验并提取 VpcNatGateway 的 DataplaneMode 配置。
// 若配置了未知或不支持的模式则返回错误。
func NatGatewayDataplaneMode(gw *kubeovnv1.VpcNatGateway) (string, error) {
	mode := gw.Spec.DataplaneMode
	if mode != "" && mode != NatDataplaneModeStatelessNft {
		return "", fmt.Errorf("gateway %s: unsupported dataplaneMode %q", gw.Name, mode)
	}
	return mode, nil
}

// StatelessNATIPv4Host 解析并标准化无状态 NAT 的 IPv4 主机地址或 /32 前缀。
// 拒绝非 IPv4、大于 /32 的子网段以及非法字符串，因为无状态 NAT 不经过 conntrack，
// 必须严格要求公网与私网具备确定性的单主机映射。
func StatelessNATIPv4Host(value string) (string, error) {
	if ip, err := netip.ParseAddr(value); err == nil && ip.Is4() {
		return ip.String(), nil
	}
	if prefix, err := netip.ParsePrefix(value); err == nil && prefix.Addr().Is4() && prefix.Bits() == 32 {
		return prefix.Addr().String(), nil
	}
	return "", fmt.Errorf("stateless NAT requires a single IPv4 address or /32, got %q", value)
}
