package controller

import (
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// natGatewayOperation 根据指定 VPC NAT 网关的数据面模式（dataplaneMode），动态解析并路由底层操作指令。
//
// 设计与容错原则：
// 1. 当 gateway 的 dataplaneMode 为 "stateless-nft" 时，将标准的有状态 NAT 操作指令转换为无状态 nft 指令：
//   - natGwSubnetFipAdd -> "stateless-fip-add"
//   - natGwSubnetFipDel -> "stateless-fip-del"
//   - natGwDnatAdd      -> "stateless-dnat-add"
//   - natGwDnatDel      -> "stateless-dnat-del"
//   - natGwSnatAdd      -> "stateless-snat-add"
//   - natGwSnatDel      -> "stateless-snat-del"
//  2. 严禁静默降级（Never Fall Back）：若在查询网关 CR 或校验 dataplaneMode 时发生任何错误，
//     必须立即将错误向上层抛出并重新排入限速工作队列重试，绝对不能静默降级为传统有状态 iptables NAT 模式。
//     若静默回退，会导致本应走无状态的流量被有状态 conntrack 捕获，产生不可逆的脏连接状态与网络黑洞。
func (c *Controller) natGatewayOperation(dp, operation string) (string, error) {
	gw, err := c.vpcNatGatewayLister.Get(dp)
	if err != nil {
		return "", fmt.Errorf("get gateway %s for %s: %w", dp, operation, err)
	}
	mode, err := util.NatGatewayDataplaneMode(gw)
	if err != nil {
		return "", err
	}
	if mode != util.NatDataplaneModeStatelessNft {
		return operation, nil
	}
	switch operation {
	case natGwSubnetFipAdd:
		return "stateless-fip-add", nil
	case natGwSubnetFipDel:
		return "stateless-fip-del", nil
	case natGwDnatAdd:
		return "stateless-dnat-add", nil
	case natGwDnatDel:
		return "stateless-dnat-del", nil
	case natGwSnatAdd:
		return "stateless-snat-add", nil
	case natGwSnatDel:
		return "stateless-snat-del", nil
	default:
		return operation, nil
	}
}

// isStatelessNatGateway 查询并判断指定名称的 VPC NAT 网关是否启用了 "stateless-nft" 无状态数据面模式。
// 若 gateway 名称为空则安全返回 (false, nil)；若查询或模式校验失败则向上返回错误。
func (c *Controller) isStatelessNatGateway(name string) (bool, error) {
	if name == "" {
		return false, nil
	}
	gw, err := c.vpcNatGatewayLister.Get(name)
	if err != nil {
		return false, err
	}
	mode, err := util.NatGatewayDataplaneMode(gw)
	return mode == util.NatDataplaneModeStatelessNft, err
}

// statelessMapping 结构体用于描述一条 NAT 规则在无状态校验体系下的公网与私网选择器元数据。
//
// 关键约束背景：
// 在无状态（notrack）网络架构中，内核不会维护连接表。对于私网 Pod 出向主动发包，网关必须依靠
// 报文的源 IP（及端口）以“单射（Injective/One-to-One）”的方式精确反向改写为对应的公网 EIP。
// 如果同一个私网 IP 被映射到了多个不同的公网身份，网关在出向将无法判定应该使用哪个 EIP 进行改写。
type statelessMapping struct {
	kind     string // 规则类别："FIP"、"SNAT" 或 "DNAT"
	name     string // K8s 规则资源名称
	eip      string // 关联的公网 EIP 资源名称
	host     string // 私网主机 IPv4 地址（已标准化为单 IP 格式）
	protocol string // 协议类型（tcp / udp，仅 DNAT 使用）
	port     string // 私网目标端口（仅 DNAT 使用）
	isShare  bool   // 是否为共享型 LoadBalancer Service DNAT 规则
}

// validateStatelessMapping 对候选的无状态映射执行严格的排他性、互斥性与反向二义性冲突校验。
//
// 校验规则详解：
//  1. 作用域判定：若 EIP 尚未绑定到网关，或网关不是 "stateless-nft" 模式，则直接跳过无状态冲突校验。
//  2. 私网地址规范化：校验 candidate.host 必须为合法的单主机 IPv4 或 /32（由 StatelessNATIPv4Host 规整）。
//  3. 收集网关全量规则：从 Lister 中提取同一网关下的所有 FIP、SNAT 与 DNAT 映射。
//  4. 冲突判定矩阵：
//     a. 公网 EIP 独占性检查：
//     - FIP 和主机级 SNAT 独占整个公网 EIP；
//     - 只有非共享的独占端口型 DNAT 规则之间才允许共享同一个公网 EIP；
//     - 无状态 EIP 严禁与任何有状态的共享型 Service DNAT 规则（isShare=true）混用。
//     b. 私网后端单射唯一性检查（Anti-Ambiguity）：
//     - 在同一个无状态网关下，私网后端 IP（及协议+端口）必须具备全局唯一性，
//     防止私网 Pod 发起出向请求时发生反向改写歧义。
func (c *Controller) validateStatelessMapping(eipName string, candidate statelessMapping) error {
	eip, err := c.iptablesEipsLister.Get(eipName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// EIP not found or not bound yet; let EIP resolution handle it
			return nil
		}
		return err
	}
	if eip.Spec.NatGwDp == "" {
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

	eips, err := c.iptablesEipsLister.List(labels.Everything())
	if err != nil {
		return err
	}
	gatewayEips := map[string]bool{}
	for _, item := range eips {
		if item.Spec.NatGwDp == eip.Spec.NatGwDp {
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
		mappings = append(mappings, statelessMapping{kind: "FIP", name: item.Name, eip: item.Spec.EIP, host: item.Spec.InternalIP})
	}
	for _, item := range snats {
		mappings = append(mappings, statelessMapping{kind: "SNAT", name: item.Name, eip: item.Spec.EIP, host: item.Spec.InternalCIDR})
	}
	for _, item := range dnats {
		mappings = append(mappings, statelessMapping{
			kind:     "DNAT",
			name:     item.Name,
			eip:      item.Spec.EIP,
			host:     item.Spec.InternalIP,
			protocol: item.Spec.Protocol,
			port:     item.Spec.InternalPort,
			isShare:  item.Spec.Type == kubeovnv1.DnatRuleTypeShare,
		})
	}

	for _, existing := range mappings {
		if !gatewayEips[existing.eip] || (existing.kind == candidate.kind && existing.name == candidate.name) {
			continue
		}
		// Stateful Share records (LoadBalancer Services) are permitted to share backends
		// with other Share records and reuse identities in retire snapshots.
		if existing.isShare && candidate.isShare {
			continue
		}
		// FIP and host SNAT consume the entire public address; exclusive port DNAT may share
		// it only with other exclusive port DNAT rules. A stateless EIP cannot be shared with
		// a stateful Share-DNAT rule.
		if existing.eip == candidate.eip && (existing.isShare || candidate.isShare || existing.kind != "DNAT" || candidate.kind != "DNAT") {
			return fmt.Errorf("stateless gateway %s: EIP %s is already used by %s %s", eip.Spec.NatGwDp, candidate.eip, existing.kind, existing.name)
		}
		existingHost, err := util.StatelessNATIPv4Host(existing.host)
		if err != nil {
			return fmt.Errorf("stateless gateway %s has invalid %s %s: %w", eip.Spec.NatGwDp, existing.kind, existing.name, err)
		}
		// Stateless backends must have unambiguous reverse routing and cannot overlap with
		// stateful Service backends.
		if existingHost == candidate.host && (existing.kind != "DNAT" || candidate.kind != "DNAT" || (existing.protocol == candidate.protocol && existing.port == candidate.port)) {
			return fmt.Errorf("stateless gateway %s: backend %s conflicts with %s %s; reverse mapping must be unique", eip.Spec.NatGwDp, candidate.host, existing.kind, existing.name)
		}
	}
	return nil
}
