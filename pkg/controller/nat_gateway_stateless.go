package controller

import (
	"fmt"
	"strconv"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// activeNatGatewayDataplaneMode 返回网关权威生效的 dataplaneMode。
// 若 Status 已记录模式则以 Status 为准；未同步或初始化前直接以 Spec 期望模式为准。
func (c *Controller) activeNatGatewayDataplaneMode(gw *kubeovnv1.VpcNatGateway) (string, error) {
	if gw == nil {
		return "", nil
	}
	if gw.Status.DataplaneMode != "" {
		return gw.Status.DataplaneMode, nil
	}
	return util.NatGatewayDataplaneMode(gw)
}

// natGatewayOperation 根据指定 VPC NAT 网关的 dataplaneMode 动态解析底层操作指令，查询失败时重新排队重试且严禁回退有状态。
func (c *Controller) natGatewayOperation(dp, operation string) (string, error) {
	if dp == "" {
		return operation, nil
	}
	gw, err := c.vpcNatGatewayLister.Get(dp)
	if err != nil {
		return "", fmt.Errorf("get gateway %s for %s: %w", dp, operation, err)
	}
	mode, err := c.activeNatGatewayDataplaneMode(gw)
	if err != nil {
		return "", err
	}
	if mode != util.NatDataplaneModeStatelessNft {
		return operation, nil
	}
	switch operation {
	case natGwSubnetFipAdd:
		return natGwStatelessFipAdd, nil
	case natGwSubnetFipDel:
		return natGwStatelessFipDel, nil
	default:
		return operation, nil
	}
}

// isStatelessNatGateway 判断指定 VPC NAT 网关是否启用了 stateless-nft 无状态数据面模式。
func (c *Controller) isStatelessNatGateway(name string) (bool, error) {
	if name == "" {
		return false, nil
	}
	gw, err := c.vpcNatGatewayLister.Get(name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	mode, err := c.activeNatGatewayDataplaneMode(gw)
	if err != nil {
		return false, err
	}
	return mode == util.NatDataplaneModeStatelessNft, nil
}

func fipMappingEstablished(fip *kubeovnv1.IptablesFIPRule) bool {
	if fip == nil {
		return false
	}
	appliedEip := fip.Annotations[util.VpcEipAnnotation]
	if appliedEip == "" {
		appliedEip = fip.Spec.EIP
	}
	return (fip.Status.Ready || fip.Status.V4ip != "" || fip.Labels[util.EipUIDLabel] != "") &&
		fip.Spec.InternalIP == fip.Status.InternalIP &&
		fip.Spec.EIP == appliedEip
}

func dnatMappingEstablished(dnat *kubeovnv1.IptablesDnatRule) bool {
	if dnat == nil {
		return false
	}
	appliedEip := dnat.Annotations[util.VpcEipAnnotation]
	if appliedEip == "" {
		appliedEip = dnat.Spec.EIP
	}
	return (dnat.Status.Ready || dnat.Status.V4ip != "" || dnat.Labels[util.EipUIDLabel] != "") &&
		dnat.Spec.InternalIP == dnat.Status.InternalIP &&
		dnat.Spec.InternalPort == dnat.Status.InternalPort &&
		dnat.Spec.ExternalPort == dnat.Status.ExternalPort &&
		strings.EqualFold(dnat.Spec.Protocol, dnat.Status.Protocol) &&
		dnat.Spec.EIP == appliedEip
}

func dnatMatchesDesired(existing, desired *kubeovnv1.IptablesDnatRule) bool {
	if existing == nil || desired == nil {
		return false
	}
	return (existing.Status.Ready || existing.Status.V4ip != "" || existing.Labels[util.EipUIDLabel] != "") &&
		existing.Spec.InternalIP == desired.Spec.InternalIP &&
		existing.Spec.InternalPort == desired.Spec.InternalPort &&
		existing.Spec.ExternalPort == desired.Spec.ExternalPort &&
		strings.EqualFold(existing.Spec.Protocol, desired.Spec.Protocol) &&
		existing.Spec.EIP == desired.Spec.EIP
}

// statelessMapping 描述无状态规则的公私网选择器，无状态反向改写无法区分指向同一后端的多个公网身份。
type statelessMapping struct {
	kind, name, eip, host, protocol, port, externalPort, owner string
	isShare                                                    bool
	isEstablished                                              bool
}

// validateStatelessMapping 校验无状态映射的 EIP 独占性与后端反向改写唯一性，防止规则冲突。
func (c *Controller) validateStatelessMapping(gwName string, candidate statelessMapping) error {
	if gwName == "" {
		return nil
	}
	enabled, err := c.isStatelessNatGateway(gwName)
	if err != nil || !enabled {
		return err
	}

	host, err := util.StatelessNATIPv4Host(candidate.host)
	if err != nil {
		return fmt.Errorf("stateless gateway %s %s %s: %w", gwName, candidate.kind, candidate.name, err)
	}
	candidate.host = host
	candidate.protocol = strings.ToLower(candidate.protocol)

	gatewayEips, err := c.getGatewayEIPs(gwName)
	if err != nil {
		return err
	}

	var mappings []statelessMapping
	switch candidate.kind {
	case "FIP":
		mappings, err = c.getGatewayFipRules(gwName, gatewayEips)
	case "DNAT":
		mappings, err = c.getGatewayDnatRules(gwName, gatewayEips)
	default:
		return nil
	}
	if err != nil {
		return err
	}

	for _, existing := range mappings {
		if err := checkStatelessMappingConflict(gwName, candidate, existing, gatewayEips); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) getGatewayEIPs(gwName string) (map[string]bool, error) {
	eips, err := c.iptablesEipsLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	gatewayEips := make(map[string]bool, len(eips))
	for _, item := range eips {
		if item.Spec.NatGwDp == gwName || item.Labels[util.VpcNatGatewayNameLabel] == gwName {
			gatewayEips[item.Name] = true
		}
	}
	return gatewayEips, nil
}

func (c *Controller) getGatewayFipRules(gwName string, gatewayEips map[string]bool) ([]statelessMapping, error) {
	fips, err := c.iptablesFipsLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var mappings []statelessMapping
	for _, item := range fips {
		if item.DeletionTimestamp != nil || item.Spec.InternalIP == "" {
			continue
		}
		if gatewayEips[item.Spec.EIP] || item.Labels[util.VpcNatGatewayNameLabel] == gwName || item.Status.NatGwDp == gwName {
			mappings = append(mappings, statelessMapping{
				kind:          "FIP",
				name:          item.Name,
				eip:           item.Spec.EIP,
				host:          item.Spec.InternalIP,
				isEstablished: fipMappingEstablished(item),
			})
		}
	}
	return mappings, nil
}

func (c *Controller) getGatewayDnatRules(gwName string, gatewayEips map[string]bool) ([]statelessMapping, error) {
	dnats, err := c.iptablesDnatRulesLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var mappings []statelessMapping
	for _, item := range dnats {
		if item.DeletionTimestamp != nil || item.Spec.InternalIP == "" {
			continue
		}
		if gatewayEips[item.Spec.EIP] || item.Labels[util.VpcNatGatewayNameLabel] == gwName || item.Status.NatGwDp == gwName {
			owner := ""
			if item.Labels[util.NftableLbSvcNsLabel] != "" && item.Labels[util.NftableLbSvcNameLabel] != "" {
				owner = item.Labels[util.NftableLbSvcNsLabel] + "/" + item.Labels[util.NftableLbSvcNameLabel]
			}
			mappings = append(mappings, statelessMapping{
				kind:          "DNAT",
				name:          item.Name,
				eip:           item.Spec.EIP,
				host:          item.Spec.InternalIP,
				protocol:      strings.ToLower(item.Spec.Protocol),
				port:          item.Spec.InternalPort,
				externalPort:  item.Spec.ExternalPort,
				owner:         owner,
				isShare:       item.Spec.Type == kubeovnv1.DnatRuleTypeShare,
				isEstablished: dnatMappingEstablished(item),
			})
		}
	}
	return mappings, nil
}

func checkStatelessMappingConflict(gwName string, candidate, existing statelessMapping, gatewayEips map[string]bool) error {
	if (existing.eip != "" && !gatewayEips[existing.eip] && existing.eip != candidate.eip) || existing.name == candidate.name {
		return nil
	}
	// 已建立生效的映射在重放恢复时，不受未生效的新建冲突请求阻断
	if candidate.isEstablished && !existing.isEstablished {
		return nil
	}

	// 1. 同 EIP 独占性与端口冲突检查
	if existing.eip != "" && existing.eip == candidate.eip {
		if candidate.kind == "FIP" {
			return fmt.Errorf("stateless gateway %s: EIP %s is already used by FIP %s", gwName, candidate.eip, existing.name)
		}
		// DNAT: 独占模式与共享模式互斥
		if existing.isShare != candidate.isShare {
			return fmt.Errorf("stateless gateway %s: EIP %s is already used by DNAT %s", gwName, candidate.eip, existing.name)
		}
		existingExtPort, _ := strconv.Atoi(existing.externalPort)
		candidateExtPort, _ := strconv.Atoi(candidate.externalPort)
		sameExtPort := existing.externalPort == candidate.externalPort || (existingExtPort != 0 && existingExtPort == candidateExtPort)
		if existing.protocol == candidate.protocol && sameExtPort {
			if !existing.isShare || existing.owner == "" || candidate.owner == "" || existing.owner != candidate.owner {
				return fmt.Errorf("stateless gateway %s: EIP %s port %s/%s is already used by DNAT %s", gwName, candidate.eip, candidate.protocol, candidate.externalPort, existing.name)
			}
		}
	}

	// 2. 后端反向唯一性检查（同一个后端不可指向不同公网身份）
	existingHost, err := util.StatelessNATIPv4Host(existing.host)
	if err != nil {
		return nil
	}
	if candidate.kind == "FIP" {
		if existingHost == candidate.host && existing.eip != candidate.eip {
			return fmt.Errorf("stateless gateway %s: backend %s conflicts with FIP %s; reverse mapping must be unique", gwName, candidate.host, existing.name)
		}
		return nil
	}

	// DNAT: 协议 + 端口 + 内网 IP 三元组唯一
	existingPort, _ := strconv.Atoi(existing.port)
	candidatePort, _ := strconv.Atoi(candidate.port)
	samePort := existing.port == candidate.port || (existingPort != 0 && existingPort == candidatePort)
	if existingHost == candidate.host && existing.protocol == candidate.protocol && samePort {
		if existing.eip != candidate.eip || existing.externalPort != candidate.externalPort {
			return fmt.Errorf("stateless gateway %s: backend %s conflicts with DNAT %s; reverse mapping must be unique", gwName, candidate.host, existing.name)
		}
	}
	return nil
}
