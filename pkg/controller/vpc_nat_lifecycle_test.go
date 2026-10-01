package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/keymutex"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	kubeovnfake "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned/fake"
	kubeovninformerfactory "github.com/kubeovn/kube-ovn/pkg/client/informers/externalversions"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// TestHandleAddSkipsTerminating covers keys that were queued while the object was still live and
// only reached the worker after deletion started. The enqueue-side routing cannot catch those, so
// the add handlers must return early instead of allocating IPs or writing rules for a dying object.
func TestHandleAddSkipsTerminating(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	final := []string{util.KubeOVNControllerFinalizer}

	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-eip", DeletionTimestamp: &now, Finalizers: final},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "gw", QoSPolicy: "missing-qos"},
	}
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-fip", DeletionTimestamp: &now, Finalizers: final},
		Spec:       kubeovnv1.IptablesFIPRuleSpec{EIP: "dying-eip", InternalIP: "10.0.0.5"},
	}
	dnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-dnat", DeletionTimestamp: &now, Finalizers: final},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "dying-eip", ExternalPort: "80", InternalPort: "8080",
			InternalIP: "10.0.0.5", Protocol: "tcp",
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-snat", DeletionTimestamp: &now, Finalizers: final},
		Spec:       kubeovnv1.IptablesSnatRuleSpec{EIP: "dying-eip", InternalCIDR: "10.0.0.0/24"},
	}
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-qos", DeletionTimestamp: &now, Finalizers: final},
		Spec:       kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{fip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{dnat},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{snat},
		QoSPolicies:       []*kubeovnv1.QoSPolicy{qos},
	})
	require.NoError(t, err)
	c := fc.fakeController

	// None of these may report an error: an error would requeue the key forever, and any work done
	// here would be work against an object that is already being torn down.
	require.NoError(t, c.handleAddIptablesEip("dying-eip"))
	require.NoError(t, c.handleAddIptablesFip("dying-fip"))
	require.NoError(t, c.handleAddIptablesDnatRule("dying-dnat"))
	require.NoError(t, c.handleAddIptablesSnatRule("dying-snat"))
	require.NoError(t, c.handleAddQoSPolicy("dying-qos"))

	// The EIP must not have been allocated an address or stamped with a QoS credential.
	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "dying-eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, stored.Status.IP)
	require.NotContains(t, stored.Labels, util.QoSPolicyUIDLabel)
}

func TestSyncVpcNatGatewayCRKeepsQoSLabels(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	qos := &kubeovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{Name: "qos", UID: "qos-uid"}}
	gw := fakeGw("gw")
	gw.Spec.QoSPolicy = "qos"

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		ConfigMaps: []*corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatGatewayConfig, Namespace: "kube-system"},
				Data:       map[string]string{"enable-vpc-nat-gw": "true"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatConfig, Namespace: "kube-system"},
				Data:       map[string]string{"image": "kubeovn/vpc-nat-gateway:test"},
			},
		},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.syncNatUIDLabels(t.Context()))
	// initResourceOnce runs this after the workers start; it must preserve the backfilled credential.
	require.NoError(t, c.syncVpcNatGatewayCR())

	stored, err := c.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "qos", stored.Labels[util.QoSLabel])
	require.Equal(t, "qos-uid", stored.Labels[util.QoSPolicyUIDLabel])
}

// TestSyncVpcNatGatewayCRDanglingQoSDoesNotBlockStartup covers a gateway left pointing at a policy
// that no longer exists. syncVpcNatGatewayCR runs from initResourceOnce, whose errors are fatal, so
// a dangling reference must not keep the controller from starting; the backfill skips the same case.
func TestSyncVpcNatGatewayCRDanglingQoSDoesNotBlockStartup(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	gw := fakeGw("gw")
	gw.Spec.QoSPolicy = "gone-qos"

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		ConfigMaps: []*corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatGatewayConfig, Namespace: "kube-system"},
				Data:       map[string]string{"enable-vpc-nat-gw": "true"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatConfig, Namespace: "kube-system"},
				Data:       map[string]string{"image": "kubeovn/vpc-nat-gateway:test"},
			},
		},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.syncNatUIDLabels(t.Context()), "backfill skips the dangling reference")
	require.NoError(t, c.syncVpcNatGatewayCR(), "startup must not fail on a dangling reference")
}

// TestQoSPolicyReleaseReadsThroughAPI mirrors the EIP side: the release decision must not run on an
// informer cache that has not observed a referrer written moments earlier, or the policy drops its
// finalizer while still in use.
func TestQoSPolicyReleaseReadsThroughAPI(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "qos",
			UID:               "qos-uid",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
		Spec: kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "eip",
			Labels: map[string]string{util.QoSLabel: "qos", util.QoSPolicyUIDLabel: "qos-uid"},
		},
		Spec: kubeovnv1.IptablesEIPSpec{QoSPolicy: "qos"},
	}
	client := kubeovnfake.NewSimpleClientset()
	_, err := client.KubeovnV1().QoSPolicies().Create(t.Context(), qos, metav1.CreateOptions{})
	require.NoError(t, err)
	_, err = client.KubeovnV1().IptablesEIPs().Create(t.Context(), eip, metav1.CreateOptions{})
	require.NoError(t, err)
	factory := kubeovninformerfactory.NewSharedInformerFactory(client, 0)
	qosInformer := factory.Kubeovn().V1().QoSPolicies()
	eipInformer := factory.Kubeovn().V1().IptablesEIPs()
	require.NoError(t, qosInformer.Informer().GetStore().Add(qos))
	c := &Controller{
		qosPoliciesLister:  qosInformer.Lister(),
		iptablesEipsLister: eipInformer.Lister(),
		vpcNatGwKeyMutex:   keymutex.NewHashed(0),
		config:             &Configuration{KubeOvnClient: client},
	}

	cached, err := c.iptablesEipsLister.List(labels.SelectorFromSet(labels.Set{util.QoSPolicyUIDLabel: "qos-uid"}))
	require.NoError(t, err)
	require.Empty(t, cached, "the informer cache has not observed the referrer yet")

	require.NoError(t, c.handleUpdateQoSPolicy("qos"))
	stored, err := c.config.KubeOvnClient.KubeovnV1().QoSPolicies().Get(t.Context(), "qos", metav1.GetOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, stored.Finalizers, "the policy is still referenced and must keep its finalizer")
}

// TestFipRebindSwapsClaimBetweenPodOperations pins the claim ordering on the rebind path. The EIP
// in-use check counts the UID label, so the old EIP has to stay claimed until its rule is gone and
// the new one has to be claimed before its rule exists; either edge lets a concurrent release drop
// a finalizer with a live rule behind it.
func TestFipRebindSwapsClaimBetweenPodOperations(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	newFip := func(statusNatGwDp string) *kubeovnv1.IptablesFIPRule {
		return &kubeovnv1.IptablesFIPRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "fip",
				Labels: map[string]string{util.VpcNatGatewayNameLabel: "gw", util.EipV4IpLabel: "1.1.1.1", util.EipUIDLabel: "old-uid"},
			},
			Spec: kubeovnv1.IptablesFIPRuleSpec{EIP: "new-eip", InternalIP: "10.0.0.1"},
			Status: kubeovnv1.IptablesFIPRuleStatus{
				V4ip: "1.1.1.1", NatGwDp: statusNatGwDp, InternalIP: "10.0.0.1", Ready: true,
			},
		}
	}
	setup := func(t *testing.T, fip *kubeovnv1.IptablesFIPRule) *Controller {
		t.Helper()
		eip := &kubeovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{Name: "new-eip", UID: "new-uid"},
			Spec:       kubeovnv1.IptablesEIPSpec{V4ip: "2.2.2.2", NatGwDp: "gw"},
			Status:     kubeovnv1.IptablesEIPStatus{Ready: true, IP: "2.2.2.2"},
		}
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
			IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
			IptablesFIPs:   []*kubeovnv1.IptablesFIPRule{fip},
		})
		require.NoError(t, err)
		return fc.fakeController
	}
	storedLabel := func(t *testing.T, c *Controller) string {
		t.Helper()
		stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
		require.NoError(t, err)
		return stored.Labels[util.EipUIDLabel]
	}

	t.Run("gateway without a running instance keeps no old rule", func(t *testing.T) {
		// Status points at the live gateway, but no instance is running: the rules only lived in
		// the container's writable layer, so the old claim can be swapped right away. The new rule
		// itself cannot be programmed without an instance, which is the error returned here.
		c := setup(t, newFip("gw"))
		require.Error(t, c.handleUpdateIptablesFip("fip"))
		require.Equal(t, "new-uid", storedLabel(t, c), "no instance holds the old rule, so nothing is orphaned by the swap")
	})

	t.Run("new claim is written before the new rule", func(t *testing.T) {
		// Status points at a gateway that is already gone, so the removal is a no-op and the run
		// gets as far as creating the new rule, which fails on the missing pod.
		c := setup(t, newFip("retired-gw"))
		require.Error(t, c.handleUpdateIptablesFip("fip"))
		require.Equal(t, "new-uid", storedLabel(t, c), "the claim must land before the rule it covers")
	})
}

// TestCreateOrUpdateEipCRClearsQoSLabels pins the update branch to the create branch: writing the
// QoS credential only when the reference is non-empty leaves the previous UID behind, and the
// in-use check counts that label.
func TestCreateOrUpdateEipCRClearsQoSLabels(t *testing.T) {
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "eip",
			Labels: map[string]string{util.QoSLabel: "gone-qos", util.QoSPolicyUIDLabel: "gone-uid"},
		},
		Spec: kubeovnv1.IptablesEIPSpec{V4ip: "1.1.1.1", NatGwDp: "gw"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.createOrUpdateEipCR("eip", "1.1.1.1", "", "", "gw", "", "external", ""))

	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, stored.Labels[util.QoSLabel])
	require.Empty(t, stored.Labels[util.QoSPolicyUIDLabel], "a dropped reference must not keep counting")
}

func TestNatUIDBackfillKeepsGenerationMismatch(t *testing.T) {
	qos := &kubeovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{Name: "qos", UID: "new-qos-uid"}}
	eip := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip", UID: "new-eip-uid"}}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:  []*kubeovnv1.QoSPolicy{qos},
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
	})
	require.NoError(t, err)
	c := fc.fakeController

	_, ok, err := c.backfillEipUIDLabel("fip", util.FipUsingEip, "eip", true, map[string]string{util.EipUIDLabel: "old-eip-uid"})
	require.NoError(t, err)
	require.False(t, ok)
	_, _, ok, err = c.qosUIDLabels("eip", "eip", "qos", true, map[string]string{util.QoSLabel: "qos", util.QoSPolicyUIDLabel: "old-qos-uid"})
	require.NoError(t, err)
	require.False(t, ok)
}

// TestWaitNatLabelClaimsSyncedNamesStragglers covers the diagnostics on a path whose caller aborts
// startup, where a bare deadline error says nothing about what failed to converge.
func TestWaitNatLabelClaimsSyncedNamesStragglers(t *testing.T) {
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesFIPs: []*kubeovnv1.IptablesFIPRule{
			{ObjectMeta: metav1.ObjectMeta{Name: "fip"}},
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), 300*time.Millisecond)
	defer cancel()
	err = fc.fakeController.waitNatLabelClaimsSynced(ctx, []natLabelClaim{
		{natType: util.FipUsingEip, name: "fip", labelKey: util.EipUIDLabel, value: "never-written"},
	})
	require.ErrorContains(t, err, "fip")
	require.ErrorContains(t, err, "1 of 1 claims still unsynced")
}

// TestPatchIptableInfoRejectsUnknownType keeps an unroutable type from reporting success, which
// would leave the startup backfill waiting for a label nobody wrote.
func TestPatchIptableInfoRejectsUnknownType(t *testing.T) {
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{})
	require.NoError(t, err)
	require.ErrorContains(t, fc.fakeController.patchIptableInfo("name", "bogus", "[]"), "unknown nat type")
}

// TestSyncNatUIDLabelsSkipsTerminating keeps the startup backfill from writing credentials the
// normal writers would refuse. Stamping a terminating policy or EIP revives a reference the
// release path was about to stop counting, which is the deadlock this series removes.
func TestSyncNatUIDLabelsSkipsTerminating(t *testing.T) {
	now := metav1.Now()
	fin := []string{util.KubeOVNControllerFinalizer}

	dyingQoS := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-qos", UID: "dying-qos-uid", DeletionTimestamp: &now, Finalizers: fin},
	}
	liveQoS := &kubeovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{Name: "live-qos", UID: "live-qos-uid"}}
	dyingEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-eip", UID: "dying-eip-uid", DeletionTimestamp: &now, Finalizers: fin},
	}
	// Live referrer pointing at a terminating policy: the policy is free to go, so nothing may
	// hand it a fresh claim.
	liveEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "live-eip", UID: "live-eip-uid"},
		Spec:       kubeovnv1.IptablesEIPSpec{QoSPolicy: "dying-qos"},
	}
	// Terminating referrer pointing at a live policy: it establishes no new binding.
	dyingFip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-fip", DeletionTimestamp: &now, Finalizers: fin},
		Spec:       kubeovnv1.IptablesFIPRuleSpec{EIP: "live-eip"},
	}
	// Live referrer pointing at a terminating EIP.
	liveSnat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "live-snat"},
		Spec:       kubeovnv1.IptablesSnatRuleSpec{EIP: "dying-eip"},
	}
	dyingGw := fakeGw("dying-gw")
	dyingGw.DeletionTimestamp = &now
	dyingGw.Finalizers = fin
	dyingGw.Spec.QoSPolicy = "live-qos"

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:       []*kubeovnv1.QoSPolicy{dyingQoS, liveQoS},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{dyingEip, liveEip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{dyingFip},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{liveSnat},
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{dyingGw},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.syncNatUIDLabels(t.Context()))

	kc := c.config.KubeOvnClient.KubeovnV1()
	eip, err := kc.IptablesEIPs().Get(t.Context(), "live-eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, eip.Labels[util.QoSPolicyUIDLabel], "a terminating policy must not gain a referrer")

	fip, err := kc.IptablesFIPRules().Get(t.Context(), "dying-fip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, fip.Labels[util.EipUIDLabel], "a terminating referrer establishes no binding")

	snat, err := kc.IptablesSnatRules().Get(t.Context(), "live-snat", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, snat.Labels[util.EipUIDLabel], "a terminating eip must not gain a referrer")

	gw, err := kc.VpcNatGateways().Get(t.Context(), "dying-gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, gw.Labels[util.QoSPolicyUIDLabel], "a terminating gateway establishes no binding")
}

// TestEipCleanupSurvivesMissingSubnet keeps a gone external subnet from holding an EIP's
// finalizer. The cleanup used to resolve the subnet before it even looked at DeletionTimestamp,
// so an EIP outliving its subnet could never finish deleting.
func TestEipCleanupSurvivesMissingSubnet(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "eip",
			UID:               "eip-uid",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
		// The subnet this points at is deliberately absent from the fixture.
		Spec:   kubeovnv1.IptablesEIPSpec{ExternalSubnet: "gone-subnet", NatGwDp: "gw"},
		Status: kubeovnv1.IptablesEIPStatus{IP: "1.1.1.1"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.handleUpdateIptablesEip("eip"))

	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, stored.Finalizers, "an unreferenced eip must release even without its subnet")
}

// TestPatchEipLabelClearsDroppedQoS pins the other half of the exempt-path rule: the recorded UID
// is preserved only while the EIP still points at the policy. Dropping the reference must clear
// both labels, or the policy keeps counting a referrer that no longer refers to it.
func TestPatchEipLabelClearsDroppedQoS(t *testing.T) {
	now := metav1.Now()
	dying := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "qos",
			UID:               "qos-uid",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "eip",
			Labels: map[string]string{util.QoSLabel: "qos", util.QoSPolicyUIDLabel: "qos-uid"},
		},
		// The reference is gone from the spec while the policy it named is still terminating.
		Spec: kubeovnv1.IptablesEIPSpec{QoSPolicy: "", NatGwDp: "gw"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:  []*kubeovnv1.QoSPolicy{dying},
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.patchEipLabel("eip"))

	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, stored.Labels[util.QoSLabel])
	require.Empty(t, stored.Labels[util.QoSPolicyUIDLabel], "a dropped reference must stop counting")
}

func TestNatRulesBecomeNotReadyWhenEipIsUnavailable(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", UID: "eip-uid"},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1"},
	}
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{Name: "fip", Labels: map[string]string{util.EipUIDLabel: "eip-uid"}},
		Spec:       kubeovnv1.IptablesFIPRuleSpec{EIP: "eip", InternalIP: "10.0.0.1"},
		Status: kubeovnv1.IptablesFIPRuleStatus{
			Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", InternalIP: "10.0.0.1",
		},
	}
	dnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "dnat", Labels: map[string]string{util.EipUIDLabel: "eip-uid"}},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip", Protocol: "tcp", ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
		Status: kubeovnv1.IptablesDnatRuleStatus{
			Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", Protocol: "tcp",
			ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "snat", Labels: map[string]string{util.EipUIDLabel: "eip-uid"}},
		Spec:       kubeovnv1.IptablesSnatRuleSpec{EIP: "eip", InternalCIDR: "10.0.0.0/24"},
		Status: kubeovnv1.IptablesSnatRuleStatus{
			Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", InternalCIDR: "10.0.0.0/24",
		},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{fip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{dnat},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{snat},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.ErrorContains(t, c.handleUpdateIptablesFip("fip"), "eip eip is not ready")
	require.ErrorContains(t, c.handleUpdateIptablesDnatRule("dnat"), "eip eip is not ready")
	require.ErrorContains(t, c.handleUpdateIptablesSnatRule("snat"), "eip eip is not ready")

	kc := c.config.KubeOvnClient.KubeovnV1()
	storedFip, err := kc.IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedFip.Status.Ready)
	require.Equal(t, "1.1.1.1", storedFip.Status.V4ip)
	storedDnat, err := kc.IptablesDnatRules().Get(t.Context(), "dnat", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedDnat.Status.Ready)
	require.Equal(t, "1.1.1.1", storedDnat.Status.V4ip)
	storedSnat, err := kc.IptablesSnatRules().Get(t.Context(), "snat", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedSnat.Status.Ready)
	require.Equal(t, "1.1.1.1", storedSnat.Status.V4ip)
}

func TestNatRulesReleaseClaimWhenEipIsDeleting(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	deletingEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", UID: "eip-uid", DeletionTimestamp: &now, Finalizers: []string{util.KubeOVNControllerFinalizer}},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "gone-gw"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1"},
	}
	goneEip := deletingEip.DeepCopy()
	goneEip.Name = "gone-eip"
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "fip",
			Labels:      eipClaimLabels("gone-gw", "1.1.1.1", "", "eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "eip"},
		},
		Spec:   kubeovnv1.IptablesFIPRuleSpec{EIP: "eip", InternalIP: "10.0.0.1"},
		Status: kubeovnv1.IptablesFIPRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "gone-gw", InternalIP: "10.0.0.1"},
	}
	dnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "dnat",
			Labels:      eipClaimLabels("gone-gw", "1.1.1.1", "80", "eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "eip"},
		},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip", Protocol: "tcp", ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
		Status: kubeovnv1.IptablesDnatRuleStatus{
			Ready: true, V4ip: "1.1.1.1", NatGwDp: "gone-gw", Protocol: "tcp",
			ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "snat",
			Labels:      eipClaimLabels("gone-gw", "1.1.1.1", "", "eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "eip"},
		},
		Spec:   kubeovnv1.IptablesSnatRuleSpec{EIP: "eip", InternalCIDR: "10.0.0.0/24"},
		Status: kubeovnv1.IptablesSnatRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "gone-gw", InternalCIDR: "10.0.0.0/24"},
	}
	goneFip := fip.DeepCopy()
	goneFip.Name = "gone-fip"
	goneFip.Spec.EIP = "gone-eip"
	goneFip.Annotations[util.VpcEipAnnotation] = "gone-eip"
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{deletingEip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{fip, goneFip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{dnat},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{snat},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("UpdateIptablesEip", nil)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	require.NoError(t, c.handleUpdateIptablesFip("fip"))
	require.NoError(t, c.handleUpdateIptablesFip("gone-fip"))
	require.NoError(t, c.handleUpdateIptablesDnatRule("dnat"))
	require.NoError(t, c.handleUpdateIptablesSnatRule("snat"))
	require.Equal(t, 2, c.updateIptablesEipQueue.Len())
	item, shutdown := c.updateIptablesEipQueue.Get()
	require.False(t, shutdown)
	require.Equal(t, "eip", item)
	c.updateIptablesEipQueue.Done(item)
	c.updateIptablesEipQueue.Forget(item)
	item, shutdown = c.updateIptablesEipQueue.Get()
	require.False(t, shutdown)
	require.Equal(t, "gone-eip", item)
	c.updateIptablesEipQueue.Done(item)
	c.updateIptablesEipQueue.Forget(item)

	kc := c.config.KubeOvnClient.KubeovnV1()
	storedFip, err := kc.IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedFip.Status.Ready, storedFip.Labels, storedFip.Annotations)
	storedGoneFip, err := kc.IptablesFIPRules().Get(t.Context(), "gone-fip", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedGoneFip.Status.Ready, storedGoneFip.Labels, storedGoneFip.Annotations)
	storedDnat, err := kc.IptablesDnatRules().Get(t.Context(), "dnat", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedDnat.Status.Ready, storedDnat.Labels, storedDnat.Annotations)
	storedSnat, err := kc.IptablesSnatRules().Get(t.Context(), "snat", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedSnat.Status.Ready, storedSnat.Labels, storedSnat.Annotations)
}

func TestNatRuleReleaseUsesBoundEipWhenSpecReboundToMissingEip(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	oldEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "old-eip", UID: "old-eip-uid", DeletionTimestamp: &now, Finalizers: []string{util.KubeOVNControllerFinalizer}},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "gone-gw"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1"},
	}
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "fip",
			Labels:      eipClaimLabels("gone-gw", "1.1.1.1", "", "old-eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "old-eip"},
		},
		Spec:   kubeovnv1.IptablesFIPRuleSpec{EIP: "missing-new-eip", InternalIP: "10.0.0.1"},
		Status: kubeovnv1.IptablesFIPRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "gone-gw", InternalIP: "10.0.0.1"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs: []*kubeovnv1.IptablesEIP{oldEip},
		IptablesFIPs: []*kubeovnv1.IptablesFIPRule{fip},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("UpdateIptablesEip", nil)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	require.NoError(t, c.handleUpdateIptablesFip("fip"))
	require.Equal(t, 1, c.updateIptablesEipQueue.Len())
	item, shutdown := c.updateIptablesEipQueue.Get()
	require.False(t, shutdown)
	require.Equal(t, "old-eip", item)
	c.updateIptablesEipQueue.Done(item)
	c.updateIptablesEipQueue.Forget(item)

	storedFip, err := c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedFip.Status.Ready, storedFip.Labels, storedFip.Annotations)
}

func TestNatRulesReleaseClaimWhenNatDisabled(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "false"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", UID: "eip-uid", DeletionTimestamp: &now, Finalizers: []string{util.KubeOVNControllerFinalizer}},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1"},
	}
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "fip",
			Labels:      eipClaimLabels("gw", "1.1.1.1", "", "eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "eip"},
		},
		Spec:   kubeovnv1.IptablesFIPRuleSpec{EIP: "eip", InternalIP: "10.0.0.1"},
		Status: kubeovnv1.IptablesFIPRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", InternalIP: "10.0.0.1"},
	}
	dnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "dnat",
			Labels:      eipClaimLabels("gw", "1.1.1.1", "80", "eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "eip"},
		},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip", Protocol: "tcp", ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
		Status: kubeovnv1.IptablesDnatRuleStatus{
			Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", Protocol: "tcp",
			ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "snat",
			Labels:      eipClaimLabels("gw", "1.1.1.1", "", "eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "eip"},
		},
		Spec:   kubeovnv1.IptablesSnatRuleSpec{EIP: "eip", InternalCIDR: "10.0.0.0/24"},
		Status: kubeovnv1.IptablesSnatRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", InternalCIDR: "10.0.0.0/24"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{fip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{dnat},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{snat},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("UpdateIptablesEip", nil)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	require.NoError(t, c.handleUpdateIptablesFip("fip"))
	require.NoError(t, c.handleUpdateIptablesDnatRule("dnat"))
	require.NoError(t, c.handleUpdateIptablesSnatRule("snat"))
	require.Equal(t, 1, c.updateIptablesEipQueue.Len())

	kc := c.config.KubeOvnClient.KubeovnV1()
	storedFip, err := kc.IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedFip.Status.Ready, storedFip.Labels, storedFip.Annotations)
	storedDnat, err := kc.IptablesDnatRules().Get(t.Context(), "dnat", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedDnat.Status.Ready, storedDnat.Labels, storedDnat.Annotations)
	storedSnat, err := kc.IptablesSnatRules().Get(t.Context(), "snat", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimReleased(t, storedSnat.Status.Ready, storedSnat.Labels, storedSnat.Annotations)
}

func TestNatRulesMoveClaimWhenOldEipIsDeletingDuringRebind(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	oldEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "old-eip", UID: "old-eip-uid", DeletionTimestamp: &now, Finalizers: []string{util.KubeOVNControllerFinalizer}},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "old-gw"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1"},
	}
	newEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "new-eip", UID: "new-eip-uid"},
		Spec:       kubeovnv1.IptablesEIPSpec{V4ip: "2.2.2.2", NatGwDp: "new-gw"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: true, IP: "2.2.2.2"},
	}
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "fip",
			Labels:      eipClaimLabels("old-gw", "1.1.1.1", "", "old-eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "old-eip"},
		},
		Spec:   kubeovnv1.IptablesFIPRuleSpec{EIP: "new-eip", InternalIP: "10.0.0.1"},
		Status: kubeovnv1.IptablesFIPRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "old-gw", InternalIP: "10.0.0.1"},
	}
	dnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "dnat",
			Labels:      eipClaimLabels("old-gw", "1.1.1.1", "80", "old-eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "old-eip"},
		},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "new-eip", Protocol: "tcp", ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
		Status: kubeovnv1.IptablesDnatRuleStatus{
			Ready: true, V4ip: "1.1.1.1", NatGwDp: "old-gw", Protocol: "tcp",
			ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "snat",
			Labels:      eipClaimLabels("old-gw", "1.1.1.1", "", "old-eip-uid"),
			Annotations: map[string]string{util.VpcEipAnnotation: "old-eip"},
		},
		Spec:   kubeovnv1.IptablesSnatRuleSpec{EIP: "new-eip", InternalCIDR: "10.0.0.0/24"},
		Status: kubeovnv1.IptablesSnatRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "old-gw", InternalCIDR: "10.0.0.0/24"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{fakeGw("new-gw")},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{oldEip, newEip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{fip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{dnat},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{snat},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("UpdateIptablesEip", nil)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	require.Error(t, c.handleUpdateIptablesFip("fip"))
	require.Error(t, c.handleUpdateIptablesDnatRule("dnat"))
	require.Error(t, c.handleUpdateIptablesSnatRule("snat"))
	require.Equal(t, 1, c.updateIptablesEipQueue.Len())
	item, shutdown := c.updateIptablesEipQueue.Get()
	require.False(t, shutdown)
	require.Equal(t, "old-eip", item)
	c.updateIptablesEipQueue.Done(item)
	c.updateIptablesEipQueue.Forget(item)

	kc := c.config.KubeOvnClient.KubeovnV1()
	storedFip, err := kc.IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimMoved(t, storedFip.Labels, storedFip.Annotations)
	storedDnat, err := kc.IptablesDnatRules().Get(t.Context(), "dnat", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimMoved(t, storedDnat.Labels, storedDnat.Annotations)
	storedSnat, err := kc.IptablesSnatRules().Get(t.Context(), "snat", metav1.GetOptions{})
	require.NoError(t, err)
	assertEipClaimMoved(t, storedSnat.Labels, storedSnat.Annotations)
}

func eipClaimLabels(gateway, v4ip, externalPort, uid string) map[string]string {
	labels := map[string]string{
		util.VpcNatGatewayNameLabel: gateway,
		util.EipV4IpLabel:           v4ip,
		util.EipUIDLabel:            uid,
	}
	if externalPort != "" {
		labels[util.VpcDnatEPortLabel] = externalPort
	}
	return labels
}

func assertEipClaimReleased(t *testing.T, ready bool, labels, annotations map[string]string) {
	t.Helper()
	require.False(t, ready)
	require.NotContains(t, labels, util.VpcNatGatewayNameLabel)
	require.NotContains(t, labels, util.VpcDnatEPortLabel)
	require.NotContains(t, labels, util.EipV4IpLabel)
	require.NotContains(t, labels, util.EipUIDLabel)
	require.NotContains(t, annotations, util.VpcEipAnnotation)
}

func assertEipClaimMoved(t *testing.T, labels, annotations map[string]string) {
	t.Helper()
	require.Equal(t, "new-eip-uid", labels[util.EipUIDLabel])
	require.Equal(t, "2.2.2.2", labels[util.EipV4IpLabel])
	require.Equal(t, "new-gw", labels[util.VpcNatGatewayNameLabel])
	require.Equal(t, "new-eip", annotations[util.VpcEipAnnotation])
}

func TestNatRuleAddWaitsForEipReady(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", UID: "eip-uid"},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1"},
	}
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{Name: "fip"},
		Spec:       kubeovnv1.IptablesFIPRuleSpec{EIP: "eip", InternalIP: "10.0.0.1"},
	}
	dnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "dnat"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip", Protocol: "tcp", ExternalPort: "80", InternalIP: "10.0.0.2", InternalPort: "8080",
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "snat"},
		Spec:       kubeovnv1.IptablesSnatRuleSpec{EIP: "eip", InternalCIDR: "10.0.0.0/24"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{fip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{dnat},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{snat},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.ErrorContains(t, c.handleAddIptablesFip("fip"), "eip eip is not ready")
	require.ErrorContains(t, c.handleAddIptablesDnatRule("dnat"), "eip eip is not ready")
	require.ErrorContains(t, c.handleAddIptablesSnatRule("snat"), "eip eip is not ready")

	kc := c.config.KubeOvnClient.KubeovnV1()
	storedFip, err := kc.IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedFip.Status.Ready)
	require.Empty(t, storedFip.Labels[util.EipUIDLabel])
	storedDnat, err := kc.IptablesDnatRules().Get(t.Context(), "dnat", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedDnat.Status.Ready)
	require.Empty(t, storedDnat.Labels[util.EipUIDLabel])
	storedSnat, err := kc.IptablesSnatRules().Get(t.Context(), "snat", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedSnat.Status.Ready)
	require.Empty(t, storedSnat.Labels[util.EipUIDLabel])
}

func TestEipBecomesNotReadyWhenQoSPolicyIsTerminating(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "qos",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
		Spec:   kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
		Status: kubeovnv1.QoSPolicyStatus{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip"},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip: "1.1.1.1", QoSPolicy: "qos", NatGwDp: "gw", ExternalSubnet: "external",
		},
		Status: kubeovnv1.IptablesEIPStatus{Ready: true, IP: "1.1.1.1", QoSPolicy: "qos"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: "external"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.ErrorContains(t, c.handleUpdateIptablesEip("eip"), "qos policy qos is terminating")
	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, stored.Status.Ready)
	require.Equal(t, "1.1.1.1", stored.Status.IP, "dependency invalidation must preserve cleanup identity")
	require.Equal(t, "qos", stored.Status.QoSPolicy, "the last successfully applied policy remains authoritative")
}

func TestEipRecoversAfterDroppingUnavailableQoSPolicy(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "qos",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
		Spec:   kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
		Status: kubeovnv1.QoSPolicyStatus{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "eip",
			Labels: map[string]string{util.QoSLabel: "qos", util.QoSPolicyUIDLabel: "qos-uid"},
		},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip: "1.1.1.1", NatGwDp: "gw", ExternalSubnet: "external",
		},
		Status: kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1", QoSPolicy: "qos"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: "external"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.handleUpdateIptablesEip("eip"))
	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.True(t, stored.Status.Ready)
	require.Empty(t, stored.Status.QoSPolicy)
	require.Empty(t, stored.Labels[util.QoSPolicyUIDLabel])

	redoEip := eip.DeepCopy()
	redoEip.Name = "redo-eip"
	redoEip.Status.Redo = time.Now().Format("2006-01-02T15:04:05")
	redoController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{redoEip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: "external"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	require.Error(t, redoController.fakeController.handleUpdateIptablesEip("redo-eip"), "redo still waits for the gateway pod")
	stored, err = redoController.fakeController.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "redo-eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, stored.Status.Ready, "qos recovery must not complete a pending redo")
	require.Empty(t, stored.Status.QoSPolicy)
}

func TestEipRecoversWhenBoundQoSPolicyBecomesUsable(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "qos", UID: "new-uid", Finalizers: []string{util.KubeOVNControllerFinalizer}},
		Spec:       kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
		Status:     kubeovnv1.QoSPolicyStatus{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", Labels: map[string]string{util.QoSPolicyUIDLabel: "new-uid"}},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip: "1.1.1.1", QoSPolicy: "qos", NatGwDp: "gw", ExternalSubnet: "external",
		},
		Status: kubeovnv1.IptablesEIPStatus{Ready: false, IP: "1.1.1.1", QoSPolicy: "qos"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: "external"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.handleUpdateIptablesEip("eip"))
	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.True(t, stored.Status.Ready)
	require.Equal(t, "qos", stored.Status.QoSPolicy)

	staleEip := eip.DeepCopy()
	staleEip.Name = "stale-eip"
	staleEip.Labels[util.QoSPolicyUIDLabel] = "old-uid"
	staleEip.Status.Ready = true
	staleController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{staleEip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: "external"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	// No gateway instance is running, so the rebind has no data plane to reconcile and is
	// deferred until one appears (the replacement instance starts empty and is programmed from
	// the live CRs): the call itself succeeds and the stale claim is left for that reconcile.
	require.NoError(t, staleController.fakeController.handleUpdateIptablesEip("stale-eip"))
}

func TestReadyEipAndFipAddReplayRouteGenerationMismatchToUpdate(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "qos", UID: "new-qos-uid", Finalizers: []string{util.KubeOVNControllerFinalizer}},
		Spec:       kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
		Status:     kubeovnv1.QoSPolicyStatus{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", UID: "new-eip-uid", Labels: map[string]string{util.QoSPolicyUIDLabel: "old-qos-uid"}},
		Spec:       kubeovnv1.IptablesEIPSpec{QoSPolicy: "qos"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: true, IP: "1.1.1.1", QoSPolicy: "qos"},
	}
	fip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{Name: "fip", Labels: map[string]string{util.EipUIDLabel: "old-eip-uid"}},
		Spec:       kubeovnv1.IptablesFIPRuleSpec{EIP: "eip", InternalIP: "10.0.0.1"},
		Status:     kubeovnv1.IptablesFIPRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", InternalIP: "10.0.0.1"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:  []*kubeovnv1.QoSPolicy{qos},
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
		IptablesFIPs: []*kubeovnv1.IptablesFIPRule{fip},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("ReplayEipUpdate", nil)
	c.updateIptablesFipQueue = newTypedRateLimitingQueue[string]("ReplayFipUpdate", nil)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)
	t.Cleanup(c.updateIptablesFipQueue.ShutDown)

	require.NoError(t, c.handleAddIptablesEip("eip"))
	require.Equal(t, 1, c.updateIptablesEipQueue.Len())
	require.NoError(t, c.handleAddIptablesFip("fip"))
	require.Equal(t, 1, c.updateIptablesFipQueue.Len())
	require.Error(t, c.handleUpdateIptablesFip("fip"), "the old rule cannot be removed without its gateway pod")
	storedFip, err := c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Get(t.Context(), "fip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedFip.Status.Ready, "a failed generation rebind must not leave the fip ready")
}

func TestReadyEipAddReplayMarksUnavailableQoSNotReady(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", Labels: map[string]string{util.QoSPolicyUIDLabel: "old-uid"}},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip: "1.1.1.1", QoSPolicy: "missing-qos", NatGwDp: "gw", ExternalSubnet: "external",
		},
		Status: kubeovnv1.IptablesEIPStatus{Ready: true, IP: "1.1.1.1", QoSPolicy: "missing-qos"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: "external"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("UnavailableQoSEipUpdate", nil)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	require.NoError(t, c.handleAddIptablesEip("eip"))
	require.Equal(t, 1, c.updateIptablesEipQueue.Len())
	require.Error(t, c.handleUpdateIptablesEip("eip"))
	storedEip, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedEip.Status.Ready)
}

func TestNatRuleAddReplayRoutesCompleteStatusToUpdateWhenEipIsUnavailable(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", UID: "new-eip-uid"},
		Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: "gw"},
		Status:     kubeovnv1.IptablesEIPStatus{IP: "1.1.1.1"},
	}
	readyFip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{Name: "ready-fip", Labels: map[string]string{util.EipUIDLabel: "new-eip-uid"}},
		Spec:       kubeovnv1.IptablesFIPRuleSpec{EIP: "eip", InternalIP: "10.0.0.1"},
		Status:     kubeovnv1.IptablesFIPRuleStatus{Ready: true, V4ip: "1.1.1.1", NatGwDp: "gw", InternalIP: "10.0.0.1"},
	}
	staleFip := readyFip.DeepCopy()
	staleFip.Name = "stale-fip"
	staleFip.Status.Ready = false
	staleFip.Labels[util.EipUIDLabel] = "old-eip-uid"
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
		IptablesFIPs:   []*kubeovnv1.IptablesFIPRule{readyFip, staleFip},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesFipQueue = newTypedRateLimitingQueue[string]("UnavailableEipFipUpdate", nil)
	t.Cleanup(c.updateIptablesFipQueue.ShutDown)

	require.NoError(t, c.handleAddIptablesFip("ready-fip"))
	require.NoError(t, c.handleAddIptablesFip("stale-fip"))
	require.Equal(t, 2, c.updateIptablesFipQueue.Len())
	require.Error(t, c.handleUpdateIptablesFip("ready-fip"))
	require.Error(t, c.handleUpdateIptablesFip("stale-fip"))
	storedReadyFip, err := c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Get(t.Context(), "ready-fip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedReadyFip.Status.Ready)
	storedStaleFip, err := c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Get(t.Context(), "stale-fip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedStaleFip.Status.Ready)
	require.Equal(t, "old-eip-uid", storedStaleFip.Labels[util.EipUIDLabel], "an unavailable eip must not claim a new generation")
}

func TestReadyEipAddReplayRoutesFinalizerLossToUpdate(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "qos", UID: "qos-uid"},
		Spec:       kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
		Status:     kubeovnv1.QoSPolicyStatus{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip", Labels: map[string]string{util.QoSPolicyUIDLabel: "qos-uid"}},
		Spec:       kubeovnv1.IptablesEIPSpec{V4ip: "1.1.1.1", QoSPolicy: "qos"},
		Status:     kubeovnv1.IptablesEIPStatus{Ready: true, IP: "1.1.1.1", QoSPolicy: "qos"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:  []*kubeovnv1.QoSPolicy{qos},
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: util.GetExternalNetwork("")},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("FinalizerLossEipUpdate", nil)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	require.NoError(t, c.handleAddIptablesEip("eip"))
	require.Equal(t, 1, c.updateIptablesEipQueue.Len())
	require.ErrorContains(t, c.handleUpdateIptablesEip("eip"), "first controller reconcile")
	storedEip, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, storedEip.Status.Ready)
}

func TestTerminatingQoSPolicyReferenceReleaseClosesLifecycle(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "qos", UID: "qos-uid", DeletionTimestamp: &now,
			Finalizers: []string{util.KubeOVNControllerFinalizer},
		},
		Spec:   kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
		Status: kubeovnv1.QoSPolicyStatus{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "eip",
			Labels: map[string]string{util.QoSLabel: "qos", util.QoSPolicyUIDLabel: "qos-uid"},
		},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip: "1.1.1.1", QoSPolicy: "qos", NatGwDp: "gw", ExternalSubnet: "external",
		},
		Status: kubeovnv1.IptablesEIPStatus{Ready: true, IP: "1.1.1.1", QoSPolicy: "qos"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("gw")},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
		Subnets: []*kubeovnv1.Subnet{{
			ObjectMeta: metav1.ObjectMeta{Name: "external"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "1.1.1.0/24"},
		}},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.updateQoSPolicyQueue = newTypedRateLimitingQueue[string]("UpdateQoSPolicy", nil)
	c.updateIptablesEipQueue = newTypedRateLimitingQueue[string]("UpdateIptablesEip", nil)
	t.Cleanup(c.updateQoSPolicyQueue.ShutDown)
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	liveQos := qos.DeepCopy()
	liveQos.DeletionTimestamp = nil
	c.enqueueUpdateQoSPolicy(liveQos, qos)
	require.Equal(t, 1, c.updateIptablesEipQueue.Len())
	require.ErrorContains(t, c.handleUpdateIptablesEip("eip"), "qos policy qos is terminating")

	// Simulate the informer observation of the user's Spec update; the API object is updated too so
	// the final assertions inspect the same generation the handler reconciles.
	cachedEip, err := c.iptablesEipsLister.Get("eip")
	require.NoError(t, err)
	cachedEip.Spec.QoSPolicy = ""
	apiEip, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	apiEip.Spec.QoSPolicy = ""
	_, err = c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Update(t.Context(), apiEip, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.NoError(t, c.handleUpdateIptablesEip("eip"))
	storedEip, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.True(t, storedEip.Status.Ready)
	require.Empty(t, storedEip.Status.QoSPolicy)
	require.Empty(t, storedEip.Labels[util.QoSPolicyUIDLabel])

	require.NoError(t, c.handleUpdateQoSPolicy("qos"))
	storedQos, err := c.config.KubeOvnClient.KubeovnV1().QoSPolicies().Get(t.Context(), "qos", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, storedQos.Finalizers)
}

func TestDeletingEipReleasesQoSPolicyFinalizerAfterEipDisappears(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "false"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "qos",
			UID:               "qos-uid",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
		Spec:   kubeovnv1.QoSPolicySpec{BindingType: kubeovnv1.QoSBindingTypeEIP},
		Status: kubeovnv1.QoSPolicyStatus{BindingType: kubeovnv1.QoSBindingTypeEIP},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "eip",
			UID:               "eip-uid",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
			Labels:            map[string]string{util.QoSLabel: "qos", util.QoSPolicyUIDLabel: "qos-uid"},
		},
		Spec:   kubeovnv1.IptablesEIPSpec{QoSPolicy: "qos", NatGwDp: "gw"},
		Status: kubeovnv1.IptablesEIPStatus{Ready: true, IP: "1.1.1.1", QoSPolicy: "qos"},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:  []*kubeovnv1.QoSPolicy{qos},
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.handleUpdateIptablesEip("eip"))
	storedEip, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, storedEip.Finalizers)
	require.Equal(t, "qos", storedEip.Labels[util.QoSLabel], "deleting eip does not need an extra self-label patch")
	require.Equal(t, "qos", storedEip.Status.QoSPolicy, "deleting eip does not need an extra self-status patch")

	require.NoError(t, c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Delete(t.Context(), "eip", metav1.DeleteOptions{}))
	require.NoError(t, c.handleUpdateQoSPolicy("qos"))
	storedQos, err := c.config.KubeOvnClient.KubeovnV1().QoSPolicies().Get(t.Context(), "qos", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, storedQos.Finalizers)
}

func TestGatewayDropsTerminatingQoSPolicyAndReleasesFinalizer(t *testing.T) {
	now := metav1.Now()
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "qos", UID: "qos-uid", DeletionTimestamp: &now,
			Finalizers: []string{util.KubeOVNControllerFinalizer},
		},
		Spec:   kubeovnv1.QoSPolicySpec{Shared: true, BindingType: kubeovnv1.QoSBindingTypeNatGw},
		Status: kubeovnv1.QoSPolicyStatus{Shared: true, BindingType: kubeovnv1.QoSBindingTypeNatGw},
	}
	gw := fakeGw("gw")
	// Spec already reflects the user's dropped reference; Status and labels describe the last
	// successful binding that still needs ordered cleanup.
	gw.Status.QoSPolicy = "qos"
	gw.Labels = map[string]string{util.QoSLabel: "qos", util.QoSPolicyUIDLabel: "qos-uid"}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)
	c := fc.fakeController

	// Empty rules avoid a pod fixture while still proving that cleanup reads a terminating policy
	// through the unguarded data-plane path.
	require.NoError(t, c.execNatGwQoS(gw, "qos", QoSDel))
	require.NoError(t, c.updateCrdNatGwLabels("gw", ""))
	require.NoError(t, c.patchNatGwQoSStatus("gw", ""))

	storedGw, err := c.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, storedGw.Status.QoSPolicy)
	require.Empty(t, storedGw.Labels[util.QoSLabel])
	require.Empty(t, storedGw.Labels[util.QoSPolicyUIDLabel])

	require.NoError(t, c.handleUpdateQoSPolicy("qos"))
	storedQos, err := c.config.KubeOvnClient.KubeovnV1().QoSPolicies().Get(t.Context(), "qos", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, storedQos.Finalizers)
}

// TestSyncVpcNatGatewayCRSkipsTerminating closes the same hole as the backfill on the other
// startup path: updateCrdNatGwLabels resolves a live policy and stamps its UID, so a gateway
// already being deleted would hand the policy a referrer that outlives the reconcile.
func TestSyncVpcNatGatewayCRSkipsTerminating(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	gw := fakeGw("dying-gw")
	gw.DeletionTimestamp = &now
	gw.Finalizers = []string{util.KubeOVNControllerFinalizer}
	gw.Spec.QoSPolicy = "live-qos"

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{{ObjectMeta: metav1.ObjectMeta{Name: "live-qos", UID: "live-qos-uid"}}},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		ConfigMaps: []*corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatGatewayConfig, Namespace: "kube-system"},
				Data:       map[string]string{"enable-vpc-nat-gw": "true"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatConfig, Namespace: "kube-system"},
				Data:       map[string]string{"image": "kubeovn/vpc-nat-gateway:test"},
			},
		},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.syncVpcNatGatewayCR())

	stored, err := c.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "dying-gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, stored.Labels[util.QoSPolicyUIDLabel], "a terminating gateway establishes no binding")
}

// TestSyncNatUIDLabelsMigratesExistingBindings is the counterpart to the skip test: refusing to
// stamp a terminating object must not drop a binding that is already programmed. An older
// controller only wrote the address label, so on upgrade these rules carry no UID and the release
// check, which selects on UID alone, would see the EIP as unreferenced.
func TestSyncNatUIDLabelsMigratesExistingBindings(t *testing.T) {
	now := metav1.Now()
	fin := []string{util.KubeOVNControllerFinalizer}

	dyingEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-eip", UID: "dying-eip-uid", DeletionTimestamp: &now, Finalizers: fin},
	}
	dyingQoS := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-qos", UID: "dying-qos-uid", DeletionTimestamp: &now, Finalizers: fin},
	}
	// Terminating rule whose old-style credential proves the rule reached the gateway pod.
	dyingFip := &kubeovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "bound-fip",
			DeletionTimestamp: &now,
			Finalizers:        fin,
			Labels:            map[string]string{util.EipV4IpLabel: "1.1.1.1"},
		},
		Spec: kubeovnv1.IptablesFIPRuleSpec{EIP: "dying-eip"},
	}
	// Live rule bound through status rather than the old label.
	boundSnat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "bound-snat"},
		Spec:       kubeovnv1.IptablesSnatRuleSpec{EIP: "dying-eip"},
		Status:     kubeovnv1.IptablesSnatRuleStatus{V4ip: "1.1.1.1"},
	}
	// Live EIP already bound to a policy that is now terminating.
	boundEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "bound-eip",
			UID:    "bound-eip-uid",
			Labels: map[string]string{util.QoSLabel: "dying-qos"},
		},
		Spec: kubeovnv1.IptablesEIPSpec{QoSPolicy: "dying-qos"},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:       []*kubeovnv1.QoSPolicy{dyingQoS},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{dyingEip, boundEip},
		IptablesFIPs:      []*kubeovnv1.IptablesFIPRule{dyingFip},
		IptablesSnatRules: []*kubeovnv1.IptablesSnatRule{boundSnat},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.syncNatUIDLabels(t.Context()))

	kc := c.config.KubeOvnClient.KubeovnV1()
	fip, err := kc.IptablesFIPRules().Get(t.Context(), "bound-fip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "dying-eip-uid", fip.Labels[util.EipUIDLabel], "a programmed rule must keep counting")

	snat, err := kc.IptablesSnatRules().Get(t.Context(), "bound-snat", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "dying-eip-uid", snat.Labels[util.EipUIDLabel], "status proves the rule exists")

	eip, err := kc.IptablesEIPs().Get(t.Context(), "bound-eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "dying-qos-uid", eip.Labels[util.QoSPolicyUIDLabel], "an established qos binding must migrate")
}

// TestSyncNatUIDLabelsClearsDroppedQoS covers the credential an interrupted run can strand. The
// update handler only rewrites these labels when spec and status disagree, which they no longer
// do once the reference is gone, so nothing else would ever stop counting this referrer.
func TestSyncNatUIDLabelsClearsDroppedQoS(t *testing.T) {
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "eip",
			UID:    "eip-uid",
			Labels: map[string]string{util.QoSLabel: "gone-qos", util.QoSPolicyUIDLabel: "gone-qos-uid"},
		},
		Spec: kubeovnv1.IptablesEIPSpec{QoSPolicy: ""},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IptablesEIPs: []*kubeovnv1.IptablesEIP{eip},
	})
	require.NoError(t, err)
	c := fc.fakeController

	require.NoError(t, c.syncNatUIDLabels(t.Context()))

	stored, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(t.Context(), "eip", metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, stored.Labels[util.QoSLabel])
	require.Empty(t, stored.Labels[util.QoSPolicyUIDLabel], "a dropped reference must stop counting")
}
