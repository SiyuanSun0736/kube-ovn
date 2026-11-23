package controller

import (
	"errors"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	listers "github.com/kubeovn/kube-ovn/pkg/client/listers/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestNatGatewayOperation(t *testing.T) {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	c := &Controller{vpcNatGatewayLister: listers.NewVpcNatGatewayLister(indexer)}
	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw"}}
	require.NoError(t, indexer.Add(gw))
	for _, op := range []string{natGwSubnetFipAdd, natGwSubnetFipDel, natGwDnatAdd, natGwDnatDel, natGwSnatAdd, natGwSnatDel} {
		got, err := c.natGatewayOperation("gw", op)
		require.NoError(t, err)
		require.Equal(t, op, got)
	}
	gw = gw.DeepCopy()
	gw.Spec.DataplaneMode = util.NatDataplaneModeStatelessNft
	require.NoError(t, indexer.Update(gw))
	for op, want := range map[string]string{
		natGwSubnetFipAdd: "stateless-fip-add", natGwSubnetFipDel: "stateless-fip-del",
		natGwDnatAdd: "stateless-dnat-add", natGwDnatDel: "stateless-dnat-del",
		natGwSnatAdd: natGwSnatAdd, natGwSnatDel: natGwSnatDel,
	} {
		got, err := c.natGatewayOperation("gw", op)
		require.NoError(t, err)
		require.Equal(t, want, got)
	}
	require.Equal(t, util.NatDataplaneModeStatelessNft, c.getNatGwDataplaneMode(gw))
	_, err := c.natGatewayOperation("missing", natGwDnatAdd)
	require.Error(t, err)
}

func TestStatelessInitFailureLeavesPodUninitialized(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })
	gw := fakeGw("gw")
	gw.Spec.DataplaneMode = util.NatDataplaneModeStatelessNft
	pod := fakeNatGwPod("gw", false)
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw}, Pods: []*corev1.Pod{pod}})
	require.NoError(t, err)
	failure := errors.New("nft transaction failed")
	fc.fakeController.execRulesInPod = func(_ *corev1.Pod, operation string, _ []string) error {
		if operation == "stateless-init" {
			return failure
		}
		return nil
	}
	require.ErrorIs(t, fc.fakeController.handleInitVpcNatGw("gw"), failure)
	current, err := fc.fakeController.config.KubeClient.CoreV1().Pods(pod.Namespace).Get(t.Context(), pod.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, current.Annotations[util.VpcNatGatewayInitAnnotation])
}

func TestInitVpcNatGwFailsOnInvalidDataplaneMode(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })
	gw := fakeGw("gw")
	gw.Spec.DataplaneMode = "unsupported-mode"
	pod := fakeNatGwPod("gw", false)
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw}, Pods: []*corev1.Pod{pod}})
	require.NoError(t, err)
	err = fc.fakeController.handleInitVpcNatGw("gw")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported dataplaneMode")
}

func TestStatelessMappingConflicts(t *testing.T) {
	gwIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	fipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	dnatIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	c := &Controller{
		vpcNatGatewayLister:     listers.NewVpcNatGatewayLister(gwIndexer),
		iptablesEipsLister:      listers.NewIptablesEIPLister(eipIndexer),
		iptablesFipsLister:      listers.NewIptablesFIPRuleLister(fipIndexer),
		iptablesDnatRulesLister: listers.NewIptablesDnatRuleLister(dnatIndexer),
	}
	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw"}, Spec: kubeovnv1.VpcNatGatewaySpec{DataplaneMode: util.NatDataplaneModeStatelessNft}}
	require.NoError(t, gwIndexer.Add(gw))
	eip := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip"}, Spec: kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"}}
	require.NoError(t, eipIndexer.Add(eip))
	other := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "other"}, Spec: kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"}}
	require.NoError(t, eipIndexer.Add(other))

	existing := &kubeovnv1.IptablesFIPRule{ObjectMeta: metav1.ObjectMeta{Name: "existing"}, Spec: kubeovnv1.IptablesFIPRuleSpec{EIP: "eip", InternalIP: "10.0.0.2"}}
	require.NoError(t, fipIndexer.Add(existing))

	for _, tt := range []struct {
		name, eip, host string
		invalid         bool
	}{
		{"new", "eip", "10.0.0.3", true},
		{"new", "other", "10.0.0.2", true},
		{"new", "other", "10.0.0.3", false},
		{"existing", "eip", "10.0.0.2", false},
	} {
		rule := &kubeovnv1.IptablesFIPRule{ObjectMeta: metav1.ObjectMeta{Name: tt.name}, Spec: kubeovnv1.IptablesFIPRuleSpec{EIP: tt.eip, InternalIP: tt.host}}
		err := c.validateFipRule(rule)
		if tt.invalid {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestStatelessAndShareDnatConflicts(t *testing.T) {
	gwIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	dnatIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	c := &Controller{
		vpcNatGatewayLister:     listers.NewVpcNatGatewayLister(gwIndexer),
		iptablesEipsLister:      listers.NewIptablesEIPLister(eipIndexer),
		iptablesDnatRulesLister: listers.NewIptablesDnatRuleLister(dnatIndexer),
		iptablesFipsLister:      listers.NewIptablesFIPRuleLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
		iptablesSnatRulesLister: listers.NewIptablesSnatRuleLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
	}
	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw"}, Spec: kubeovnv1.VpcNatGatewaySpec{DataplaneMode: util.NatDataplaneModeStatelessNft}}
	require.NoError(t, gwIndexer.Add(gw))
	eip1 := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip1"}, Spec: kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"}}
	eip2 := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip2"}, Spec: kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"}}
	require.NoError(t, eipIndexer.Add(eip1))
	require.NoError(t, eipIndexer.Add(eip2))

	// Existing share rule (e.g. from an nftable LoadBalancer Service)
	shareExisting := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "share-existing"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip1", InternalIP: "10.0.0.2", Protocol: "tcp", InternalPort: "8080", ExternalPort: "80", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	require.NoError(t, dnatIndexer.Add(shareExisting))

	// 1. Two stateful share rules with same backend endpoint are permitted (conntrack tracks individual sessions)
	shareAnotherOk := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "share-another-ok"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip2", InternalIP: "10.0.0.2", Protocol: "tcp", InternalPort: "8080", ExternalPort: "8088", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	require.NoError(t, c.validateDnatRule(shareAnotherOk))

	// 1.1 Another share rule with different backend port is permitted
	shareAnotherPortOk := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "share-another-port-ok"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip2", InternalIP: "10.0.0.2", Protocol: "tcp", InternalPort: "8081", ExternalPort: "8088", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	require.NoError(t, c.validateDnatRule(shareAnotherPortOk))

	// 2. A stateless DNAT rule CANNOT share the same backend with a share rule
	statelessConflictBackend := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "sl-conflict-backend"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip2", InternalIP: "10.0.0.2", Protocol: "tcp", InternalPort: "8080", ExternalPort: "8088",
		},
	}
	require.Error(t, c.validateDnatRule(statelessConflictBackend))

	// 3. A stateless DNAT rule CANNOT share the EIP of a stateful share rule
	statelessConflictEip := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "sl-conflict-eip"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip1", InternalIP: "10.0.0.3", Protocol: "tcp", InternalPort: "8081", ExternalPort: "8088",
		},
	}
	require.Error(t, c.validateDnatRule(statelessConflictEip))

	// 4. A stateless DNAT rule with different EIP and different backend is allowed
	statelessOk := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "sl-ok"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip2", InternalIP: "10.0.0.3", Protocol: "tcp", InternalPort: "8081", ExternalPort: "8088",
		},
	}
	require.NoError(t, c.validateDnatRule(statelessOk))

	// 5. Exclusive DNAT rules on the same EIP CANNOT collide on externalPort
	require.NoError(t, dnatIndexer.Add(statelessOk))
	exclusivePortConflict := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "sl-port-conflict"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip2", InternalIP: "10.0.0.4", Protocol: "tcp", InternalPort: "8082", ExternalPort: "8088",
		},
	}
	require.Error(t, c.validateDnatRule(exclusivePortConflict))

	// 5.1 Exclusive DNAT rules on the same EIP with distinct externalPort are allowed
	exclusivePortOk := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "sl-port-ok"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip2", InternalIP: "10.0.0.4", Protocol: "tcp", InternalPort: "8082", ExternalPort: "8089",
		},
	}
	require.NoError(t, c.validateDnatRule(exclusivePortOk))

	// 5.2 Case-insensitive protocol matching detects conflict even if protocol is uppercase
	candidateUpperProtocol := statelessMapping{
		kind:         "DNAT",
		name:         "sl-port-upper-conflict",
		eip:          "eip2",
		host:         "10.0.0.4",
		protocol:     "TCP",
		port:         "8082",
		externalPort: "8088",
	}
	require.Error(t, c.validateStatelessMapping("gw", candidateUpperProtocol))
}

func TestStatelessGatewayRejectsTCPConnLimit(t *testing.T) {
	limit := int32(1000)
	gw := fakeGw("gw")
	gw.Spec.DataplaneMode = util.NatDataplaneModeStatelessNft
	gw.Spec.TCPConnLimit = &limit
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw}})
	require.NoError(t, err)
	err = fc.fakeController.handleAddOrUpdateVpcNatGw("gw")
	require.ErrorContains(t, err, "stateless dataplane mode does not support tcpConnLimit")
}

func TestStatelessConflictWithLanIPService(t *testing.T) {
	gwIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	dnatIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	c := &Controller{
		vpcNatGatewayLister:     listers.NewVpcNatGatewayLister(gwIndexer),
		iptablesEipsLister:      listers.NewIptablesEIPLister(eipIndexer),
		iptablesDnatRulesLister: listers.NewIptablesDnatRuleLister(dnatIndexer),
		iptablesFipsLister:      listers.NewIptablesFIPRuleLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
		iptablesSnatRulesLister: listers.NewIptablesSnatRuleLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
	}
	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw"}, Spec: kubeovnv1.VpcNatGatewaySpec{DataplaneMode: util.NatDataplaneModeStatelessNft}}
	require.NoError(t, gwIndexer.Add(gw))
	eip := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip"}, Spec: kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"}}
	require.NoError(t, eipIndexer.Add(eip))

	// Existing exclusive stateless DNAT rule
	exclusiveDnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "sl-dnat"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip", InternalIP: "10.0.0.2", Protocol: "tcp", InternalPort: "8080", ExternalPort: "80",
		},
	}
	require.NoError(t, dnatIndexer.Add(exclusiveDnat))

	// LanIP Service VIP (EIP is empty) pointing to same backend should be rejected
	candidate := statelessMapping{
		kind:         "DNAT",
		name:         "lanip-svc",
		eip:          "",
		host:         "10.0.0.2",
		protocol:     "tcp",
		port:         "8080",
		externalPort: "80",
		isShare:      true,
	}
	err := c.validateStatelessMapping("gw", candidate)
	require.Error(t, err)
	require.ErrorContains(t, err, "reverse mapping must be unique")
}

func TestStatelessRetryPreservesMode(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	gw := fakeGw("gw")
	gw.Spec.DataplaneMode = util.NatDataplaneModeStatelessNft
	gw.Status.HairpinSNAT = "enable" // simulate partial init from a previous attempt
	pod := fakeNatGwPod("gw", false)
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		Pods:           []*corev1.Pod{pod},
	})
	require.NoError(t, err)

	for _, q := range []*workqueue.TypedRateLimitingInterface[string]{
		&fc.fakeController.updateVpcFloatingIPQueue, &fc.fakeController.updateVpcDnatQueue, &fc.fakeController.updateVpcSnatQueue,
		&fc.fakeController.updateVpcSubnetQueue, &fc.fakeController.updateVpcEipQueue,
	} {
		*q = newTypedRateLimitingQueue[string]("test", nil)
		t.Cleanup((*q).ShutDown)
	}

	var calls [][]string
	fc.fakeController.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
		calls = append(calls, append([]string{operation}, rules...))
		return nil
	}

	require.NoError(t, fc.fakeController.handleInitVpcNatGw("gw"))
	// Must still execute stateless-init and not degrade to stateful
	foundStatelessInit := false
	for _, call := range calls {
		if len(call) > 0 && call[0] == "stateless-init" {
			foundStatelessInit = true
			break
		}
	}
	require.True(t, foundStatelessInit, "stateless-init must be called even when HairpinSNAT is already populated")

	patched, err := fc.fakeController.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, util.NatDataplaneModeStatelessNft, patched.Status.DataplaneMode)
}
