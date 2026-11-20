package controller

import (
	"fmt"
	"strconv"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// activeNatGatewayDataplaneMode 返回网关权威生效的 dataplaneMode。
// 一旦网关已通过 handleInitVpcNatGw 完成数据面初始化，其模式具有绝对权威，屏蔽任何未生效的 Spec 篡改；
// 若网关仍在初始化前（新建网关），则以 Spec 期望的模式为准。
func (c *Controller) activeNatGatewayDataplaneMode(gw *kubeovnv1.VpcNatGateway) (string, error) {
	if gw == nil {
		return "", nil
	}
	if gw.Status.DataplaneMode != "" {
		return gw.Status.DataplaneMode, nil
	}
	if c != nil && c.podsLister != nil {
		pod, err := c.getNatGwPod(gw.Name, c.natGwNamespace(gw))
		if err == nil && pod != nil && pod.Annotations[util.VpcNatGatewayInitAnnotation] != "" {
			// Pod 已经完成了初始化流程，且 Status.DataplaneMode 为空，表明其为已初始化的默认 stateful 网关
			return "", nil
		}
	}
	return util.NatGatewayDataplaneMode(gw)
}

// natGatewayOperation 根据指定 VPC NAT 网关的 dataplaneMode 动态解析底层操作指令，查询失败时重新排队重试且严禁回退有状态。
func (c *Controller) natGatewayOperation(dp, operation string) (string, error) {
	if c == nil || c.vpcNatGatewayLister == nil {
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
	case natGwDnatAdd:
		return natGwStatelessDnatAdd, nil
	case natGwDnatDel:
		return natGwStatelessDnatDel, nil
	case natGwSnatAdd:
		return natGwStatelessSnatAdd, nil
	case natGwSnatDel:
		return natGwStatelessSnatDel, nil
	default:
		return operation, nil
	}
}

// isStatelessNatGateway 判断指定 VPC NAT 网关是否启用了 stateless-nft 无状态数据面模式。
func (c *Controller) isStatelessNatGateway(name string) (bool, error) {
	if name == "" {
		return false, nil
	}
	if c == nil || c.vpcNatGatewayLister == nil {
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

// statelessMapping 描述无状态规则的公私网选择器，无状态反向改写无法区分指向同一后端的多个公网身份。
type statelessMapping struct {
	kind, name, eip, host, protocol, port, externalPort string
	isShare                                             bool
	isEstablished                                       bool
}

// validateStatelessMapping 校验无状态映射的 EIP 独占性与后端反向改写唯一性，防止规则冲突。
func (c *Controller) validateStatelessMapping(eipName string, candidate statelessMapping) error {
	if c == nil || c.iptablesEipsLister == nil || c.vpcNatGatewayLister == nil {
		return nil
	}
	eip, err := c.iptablesEipsLister.Get(eipName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// EIP not found or not bound yet; let EIP resolution handle it
			return nil
		}
		return err
	}
	if eip.Spec.NatGwDp == "" || !eip.DeletionTimestamp.IsZero() {
		return nil
	}
	enabled, err := c.isStatelessNatGateway(eip.Spec.NatGwDp)
	if err != nil || !enabled {
		return err
	}

	host, err := util.StatelessNATIPv4Host(candidate.host)
	if err != nil {
		return fmt.Errorf("stateless gateway %s %s %s: %w", eip.Spec.NatGwDp, candidate.kind, candidate.name, err)
	}
	candidate.host = host

	gwName := eip.Spec.NatGwDp

	eips, err := c.iptablesEipsLister.List(labels.Everything())
	if err != nil {
		return err
	}
	gatewayEips := map[string]bool{}
	for _, item := range eips {
		if item.Spec.NatGwDp == gwName || item.Labels[util.VpcNatGatewayNameLabel] == gwName {
			gatewayEips[item.Name] = true
		}
	}

	fips, err := c.iptablesFipsLister.List(labels.Everything())
	if err != nil {
		return err
	}
	snats, err := c.iptablesSnatRulesLister.List(labels.Everything())
	if err != nil {
		return err
	}
	dnats, err := c.iptablesDnatRulesLister.List(labels.Everything())
	if err != nil {
		return err
	}

	var mappings []statelessMapping
	for _, item := range fips {
		appliedEip := item.Annotations[util.VpcEipAnnotation]
		if appliedEip == "" {
			appliedEip = item.Spec.EIP
		}
		appliedGw := item.Status.NatGwDp
		if appliedGw == "" {
			appliedGw = item.Labels[util.VpcNatGatewayNameLabel]
		}
		specEstablished := (item.Status.Ready || item.Status.V4ip != "" || item.Labels[util.EipUIDLabel] != "") &&
			item.Spec.InternalIP == item.Status.InternalIP && item.Spec.EIP == appliedEip
		if gatewayEips[item.Spec.EIP] || item.Labels[util.VpcNatGatewayNameLabel] == gwName {
			mappings = append(mappings, statelessMapping{kind: "FIP", name: item.Name, eip: item.Spec.EIP, host: item.Spec.InternalIP, isEstablished: specEstablished})
		}
		if item.Status.InternalIP != "" && (appliedGw == gwName || gatewayEips[appliedEip]) {
			if !specEstablished {
				mappings = append(mappings, statelessMapping{kind: "FIP", name: item.Name, eip: appliedEip, host: item.Status.InternalIP, isEstablished: true})
			}
		}
	}
	for _, item := range snats {
		appliedEip := item.Annotations[util.VpcEipAnnotation]
		if appliedEip == "" {
			appliedEip = item.Spec.EIP
		}
		appliedGw := item.Status.NatGwDp
		if appliedGw == "" {
			appliedGw = item.Labels[util.VpcNatGatewayNameLabel]
		}
		specEstablished := (item.Status.Ready || item.Status.V4ip != "" || item.Labels[util.EipUIDLabel] != "") &&
			item.Spec.InternalCIDR == item.Status.InternalCIDR && item.Spec.EIP == appliedEip
		if gatewayEips[item.Spec.EIP] || item.Labels[util.VpcNatGatewayNameLabel] == gwName {
			mappings = append(mappings, statelessMapping{kind: "SNAT", name: item.Name, eip: item.Spec.EIP, host: item.Spec.InternalCIDR, isEstablished: specEstablished})
		}
		if item.Status.InternalCIDR != "" && (appliedGw == gwName || gatewayEips[appliedEip]) {
			if !specEstablished {
				mappings = append(mappings, statelessMapping{kind: "SNAT", name: item.Name, eip: appliedEip, host: item.Status.InternalCIDR, isEstablished: true})
			}
		}
	}
	for _, item := range dnats {
		appliedEip := item.Annotations[util.VpcEipAnnotation]
		if appliedEip == "" {
			appliedEip = item.Spec.EIP
		}
		appliedGw := item.Status.NatGwDp
		if appliedGw == "" {
			appliedGw = item.Labels[util.VpcNatGatewayNameLabel]
		}
		specEstablished := (item.Status.Ready || item.Status.V4ip != "" || item.Labels[util.EipUIDLabel] != "") &&
			item.Spec.InternalIP == item.Status.InternalIP &&
			item.Spec.InternalPort == item.Status.InternalPort &&
			item.Spec.ExternalPort == item.Status.ExternalPort &&
			item.Spec.Protocol == item.Status.Protocol &&
			item.Spec.EIP == appliedEip
		if gatewayEips[item.Spec.EIP] || item.Labels[util.VpcNatGatewayNameLabel] == gwName {
			mappings = append(mappings, statelessMapping{
				kind:          "DNAT",
				name:          item.Name,
				eip:           item.Spec.EIP,
				host:          item.Spec.InternalIP,
				protocol:      item.Spec.Protocol,
				port:          item.Spec.InternalPort,
				externalPort:  item.Spec.ExternalPort,
				isShare:       item.Spec.Type == kubeovnv1.DnatRuleTypeShare,
				isEstablished: specEstablished,
			})
		}
		if item.Status.InternalIP != "" && (appliedGw == gwName || gatewayEips[appliedEip]) {
			if !specEstablished {
				mappings = append(mappings, statelessMapping{
					kind:          "DNAT",
					name:          item.Name,
					eip:           appliedEip,
					host:          item.Status.InternalIP,
					protocol:      item.Status.Protocol,
					port:          item.Status.InternalPort,
					externalPort:  item.Status.ExternalPort,
					isShare:       item.Spec.Type == kubeovnv1.DnatRuleTypeShare,
					isEstablished: true,
				})
			}
		}
	}

	for _, existing := range mappings {
		if (!gatewayEips[existing.eip] && existing.eip != candidate.eip) || (existing.kind == candidate.kind && existing.name == candidate.name) {
			continue
		}
		// 已建立生效的映射在重放恢复时，不受未生效的新建冲突请求阻断
		if candidate.isEstablished && !existing.isEstablished {
			continue
		}
		// Stateful Share records (LoadBalancer Services) are handled by conntrack and nftable LB;
		// they do not participate in stateless reverse routing, so two stateful share rules may
		// reuse the same backend without conflict.
		if existing.isShare && candidate.isShare {
			continue
		}
		// FIP and host SNAT consume the entire public address; exclusive port DNAT may share
		// it only with other exclusive port DNAT rules on distinct external ports. A stateless EIP cannot be
		// shared with a stateful Share-DNAT rule.
		if existing.eip == candidate.eip {
			if existing.isShare || candidate.isShare || existing.kind != "DNAT" || candidate.kind != "DNAT" {
				return fmt.Errorf("stateless gateway %s: EIP %s is already used by %s %s", eip.Spec.NatGwDp, candidate.eip, existing.kind, existing.name)
			}
			existingExtPort, _ := strconv.Atoi(existing.externalPort)
			candidateExtPort, _ := strconv.Atoi(candidate.externalPort)
			if existing.protocol == candidate.protocol && (existing.externalPort == candidate.externalPort || (existingExtPort != 0 && existingExtPort == candidateExtPort)) {
				return fmt.Errorf("stateless gateway %s: EIP %s port %s/%s is already used by %s %s", eip.Spec.NatGwDp, candidate.eip, candidate.protocol, candidate.externalPort, existing.kind, existing.name)
			}
		}
		existingHost, err := util.StatelessNATIPv4Host(existing.host)
		if err != nil {
			continue
		}
		existingPort, _ := strconv.Atoi(existing.port)
		candidatePort, _ := strconv.Atoi(candidate.port)
		samePort := existing.port == candidate.port || (existingPort != 0 && existingPort == candidatePort)
		// Stateless backends must have unambiguous reverse routing and cannot overlap with
		// stateful Service backends.
		if existingHost == candidate.host && (existing.kind != "DNAT" || candidate.kind != "DNAT" || (existing.protocol == candidate.protocol && samePort)) {
			return fmt.Errorf("stateless gateway %s: backend %s conflicts with %s %s; reverse mapping must be unique", eip.Spec.NatGwDp, candidate.host, existing.kind, existing.name)
		}
	}
	return nil
}
