package controller

import (
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// natGatewayOperation translates generic NAT operations (e.g., natGwSubnetFipAdd)
// into their stateless-nft equivalents when the gateway dataplaneMode is stateless-nft.
// It never falls back to stateful NAT on a failed lookup: the caller returns the error
// to the rate-limited work queue before publishing Ready.
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

// isStatelessNatGateway reports whether the specified VPC NAT gateway is configured
// with the stateless-nft dataplane mode.
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

// statelessMapping describes the public and private selectors of one rule.
// Reverse rewriting cannot disambiguate two public identities for one backend.
type statelessMapping struct {
	kind, name, eip, host, protocol, port string
	isShare                               bool
}

// validateStatelessMapping enforces mutual exclusivity and disambiguation rules
// for stateless NAT mappings belonging to the same gateway. It prevents multiple rules
// from claiming the same private host or colliding on public EIP assignments.
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
