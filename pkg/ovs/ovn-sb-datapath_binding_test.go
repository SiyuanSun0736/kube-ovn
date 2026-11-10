package ovs

import (
	"testing"

	"github.com/stretchr/testify/require"

	ovsclient "github.com/kubeovn/kube-ovn/pkg/ovsdb/client"
	"github.com/kubeovn/kube-ovn/pkg/ovsdb/ovnsb"
)

// newDatapathBinding mimics northd: kindID is "logical-switch" or "logical-router" and holds the
// NB UUID of the datapath.
func newDatapathBinding(kindID, name string, tunnelKey int) *ovnsb.DatapathBinding {
	return &ovnsb.DatapathBinding{
		UUID:        ovsclient.NamedUUID(),
		ExternalIDs: map[string]string{"name": name, kindID: "nb-uuid-" + name},
		TunnelKey:   tunnelKey,
	}
}

func (suite *OvnClientTestSuite) testGetLogicalSwitchTunnelKey() {
	t := suite.T()

	sbClient := suite.ovnSBClient
	nbClient := suite.ovnNBClient
	lsName := "test-ls-tunnel-key"

	// Create logical switch in NB to trigger datapath binding creation in SB
	// Note: In real OVN, northd syncs NB to SB. In test, we simulate by creating datapath binding directly.
	err := nbClient.CreateBareLogicalSwitch(lsName)
	require.NoError(t, err)

	// In mock test environment, we need to create datapath binding manually
	// since there's no northd to sync NB to SB
	tunnelKey := 100
	dp := newDatapathBinding("logical-switch", lsName, tunnelKey)
	ops, err := sbClient.Create(dp)
	require.NoError(t, err)
	err = sbClient.Transact("datapath-binding-add", ops)
	require.NoError(t, err)

	t.Run("should return tunnel key for existing ls", func(t *testing.T) {
		key, err := sbClient.GetLogicalSwitchTunnelKey(lsName)
		require.NoError(t, err)
		require.Equal(t, tunnelKey, key)
	})

	t.Run("should fail with empty ls name", func(t *testing.T) {
		key, err := sbClient.GetLogicalSwitchTunnelKey("")
		require.Equal(t, 0, key)
		require.ErrorContains(t, err, "logical switch name is empty")
	})

	t.Run("should fail for non-existent ls", func(t *testing.T) {
		key, err := sbClient.GetLogicalSwitchTunnelKey("non-existent-ls")
		require.Equal(t, 0, key)
		require.ErrorContains(t, err, "datapath binding not found")
	})
}

func (suite *OvnClientTestSuite) testGetLogicalRouterTunnelKey() {
	t := suite.T()

	sbClient := suite.ovnSBClient
	addBinding := func(kindID, name string, tunnelKey int) {
		ops, err := sbClient.Create(newDatapathBinding(kindID, name, tunnelKey))
		require.NoError(t, err)
		require.NoError(t, sbClient.Transact("datapath-binding-add", ops))
	}
	// No northd in the test environment, so bindings are created directly.
	addBinding("logical-router", "test-lr-tunnel-key", 200)
	addBinding("logical-router", "test-lr-duplicate", 201)
	addBinding("logical-router", "test-lr-duplicate", 202)
	// Only the same-named switch is published so far: its key must not be taken for the router.
	addBinding("logical-switch", "test-lr-not-published", 203)
	// Both published under one name: each kind resolves to its own key.
	addBinding("logical-switch", "test-shared-name", 204)
	addBinding("logical-router", "test-shared-name", 205)

	t.Run("should return tunnel key for existing lr", func(t *testing.T) {
		key, err := sbClient.GetLogicalRouterTunnelKey("test-lr-tunnel-key")
		require.NoError(t, err)
		require.Equal(t, 200, key)
	})

	t.Run("should fail with empty lr name", func(t *testing.T) {
		key, err := sbClient.GetLogicalRouterTunnelKey("")
		require.Equal(t, 0, key)
		require.ErrorContains(t, err, "logical router name is empty")
	})

	t.Run("should fail for non-existent lr", func(t *testing.T) {
		key, err := sbClient.GetLogicalRouterTunnelKey("non-existent-lr")
		require.Equal(t, 0, key)
		require.ErrorContains(t, err, "datapath binding not found")
	})

	t.Run("should fail closed on duplicate router bindings", func(t *testing.T) {
		key, err := sbClient.GetLogicalRouterTunnelKey("test-lr-duplicate")
		require.Equal(t, 0, key)
		require.ErrorContains(t, err, "expected exactly one")
	})

	t.Run("should not accept a same-named switch while the router is unpublished", func(t *testing.T) {
		key, err := sbClient.GetLogicalRouterTunnelKey("test-lr-not-published")
		require.Equal(t, 0, key)
		require.ErrorContains(t, err, "datapath binding not found")
	})

	t.Run("should resolve each kind of a shared name to its own key", func(t *testing.T) {
		key, err := sbClient.GetLogicalRouterTunnelKey("test-shared-name")
		require.NoError(t, err)
		require.Equal(t, 205, key)
		key, err = sbClient.GetLogicalSwitchTunnelKey("test-shared-name")
		require.NoError(t, err)
		require.Equal(t, 204, key)
	})
}
