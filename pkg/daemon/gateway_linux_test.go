package daemon

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestGetTunnelIptablesRules(t *testing.T) {
	tunnelMatchSet := "ovn40tunnel-nodes"

	t.Run("geneve with resolved tunnel interface and match-set", func(t *testing.T) {
		cfg := &Configuration{
			NetworkType:         util.NetworkTypeGeneve,
			Iface:               "en.*,bond0", // Regex selector preserved in Iface
			EnableTunnelNotrack: true,
		}
		cfg.SetTunnelIface("bond0") // Resolved device stored in tunnelIface

		c := &Controller{config: cfg}
		active, obsolete := c.getTunnelIptablesRules(tunnelMatchSet)

		require.Len(t, active, 2)
		require.Equal(t, util.Raw, active[0].Table)
		require.Equal(t, util.Prerouting, active[0].Chain)
		require.Equal(t, strings.Fields("-i bond0 -p udp -m udp --dport 6081 -m set --match-set ovn40tunnel-nodes src -j NOTRACK"), active[0].Rule)

		require.Equal(t, util.Raw, active[1].Table)
		require.Equal(t, util.Output, active[1].Chain)
		require.Equal(t, strings.Fields("-o bond0 -p udp -m udp --dport 6081 -m set --match-set ovn40tunnel-nodes dst -j NOTRACK"), active[1].Rule)

		// Obsolete rules should clean inactive port rules
		var hasMatchSet4789Pre bool
		for _, r := range obsolete {
			ruleStr := strings.Join(r.Rule, " ")
			if ruleStr == "-i bond0 -p udp -m udp --dport 4789 -m set --match-set ovn40tunnel-nodes src -j NOTRACK" && r.Chain == util.Prerouting {
				hasMatchSet4789Pre = true
			}
		}
		require.True(t, hasMatchSet4789Pre, "inactive 4789 match-set rule must be obsolete")
	})

	t.Run("disabled tunnel notrack flag emits no active rules and cleans all rules", func(t *testing.T) {
		cfg := &Configuration{
			NetworkType:         util.NetworkTypeGeneve,
			EnableTunnelNotrack: false, // Flag disabled
		}
		cfg.SetTunnelIface("bond0")

		c := &Controller{config: cfg}
		active, obsolete := c.getTunnelIptablesRules(tunnelMatchSet)

		require.Empty(t, active)
		require.NotEmpty(t, obsolete)

		var hasMatchSet6081Pre bool
		for _, r := range obsolete {
			ruleStr := strings.Join(r.Rule, " ")
			if ruleStr == "-i bond0 -p udp -m udp --dport 6081 -m set --match-set ovn40tunnel-nodes src -j NOTRACK" && r.Chain == util.Prerouting {
				hasMatchSet6081Pre = true
			}
		}
		require.True(t, hasMatchSet6081Pre, "active port match-set rule must be cleaned up when feature is disabled")
	})

	t.Run("vxlan without iface emits no active rules", func(t *testing.T) {
		cfg := &Configuration{
			NetworkType:         util.NetworkTypeVxlan,
			EnableTunnelNotrack: true,
		}
		cfg.SetTunnelIface("")

		c := &Controller{config: cfg}
		active, obsolete := c.getTunnelIptablesRules(tunnelMatchSet)

		require.Empty(t, active)
		require.Empty(t, obsolete)
	})

	t.Run("vlan network type emits no active rules and cleans all ports", func(t *testing.T) {
		cfg := &Configuration{
			NetworkType:         util.NetworkTypeVlan,
			EnableTunnelNotrack: true,
		}
		cfg.SetTunnelIface("bond0")

		c := &Controller{config: cfg}
		active, obsolete := c.getTunnelIptablesRules(tunnelMatchSet)

		require.Empty(t, active)
		require.NotEmpty(t, obsolete)
	})
}
