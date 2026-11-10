package ovs

import (
	"context"
	"fmt"

	"github.com/kubeovn/kube-ovn/pkg/ovsdb/ovnsb"
)

// GetLogicalSwitchTunnelKey retrieves the tunnel_key for a logical switch from OVN SB Datapath_Binding.
func (c *OVNSbClient) GetLogicalSwitchTunnelKey(lsName string) (int, error) {
	return c.getDatapathTunnelKey("logical switch", "logical-switch", lsName)
}

// GetLogicalRouterTunnelKey retrieves the tunnel_key for a logical router (VPC) from OVN SB
// Datapath_Binding. northd allocates switch and router keys from one namespace.
func (c *OVNSbClient) GetLogicalRouterTunnelKey(lrName string) (int, error) {
	return c.getDatapathTunnelKey("logical router", "logical-router", lrName)
}

// getDatapathTunnelKey matches Datapath_Binding by external_ids:name and the external ID northd
// sets for the datapath kind (logical-switch/logical-router = NB UUID). Matching the kind keeps a
// same-named datapath of the other kind from being accepted while the wanted one is not
// published yet.
func (c *OVNSbClient) getDatapathTunnelKey(kind, kindID, name string) (int, error) {
	if name == "" {
		return 0, fmt.Errorf("%s name is empty", kind)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	dpList := make([]ovnsb.DatapathBinding, 0)
	if err := c.ovsDbClient.WhereCache(func(dp *ovnsb.DatapathBinding) bool {
		if dp.ExternalIDs == nil {
			return false
		}
		return dp.ExternalIDs["name"] == name && dp.ExternalIDs[kindID] != ""
	}).List(ctx, &dpList); err != nil {
		return 0, fmt.Errorf("failed to list datapath binding for %s %s: %w", kind, name, err)
	}

	if len(dpList) == 0 {
		return 0, fmt.Errorf("datapath binding not found for %s %s", kind, name)
	}

	if len(dpList) > 1 {
		// Each datapath has exactly one binding; more means an NB/SB inconsistency.
		return 0, fmt.Errorf("found %d datapath bindings for %s %s, expected exactly one: OVN SB database may be inconsistent", len(dpList), kind, name)
	}

	return dpList[0].TunnelKey, nil
}
