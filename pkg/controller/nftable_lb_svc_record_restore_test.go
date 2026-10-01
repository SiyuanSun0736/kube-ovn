package controller

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"

	v1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/set"

	k8stesting "k8s.io/client-go/testing"

	kubeovnfake "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned/fake"
)

// Test_enqueueDelIptablesDnatRule_wakesOwnerService pins the record-loss wake-up: deleting a
// Service accounting record out from under its owner re-enqueues that Service so its reconcile
// restores the claim, while non-accounting share records and exclusive rules keep their existing
// behavior.
func Test_enqueueDelIptablesDnatRule_wakesOwnerService(t *testing.T) {
	t.Parallel()

	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(queue.ShutDown)
	delQueue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(delQueue.ShutDown)
	c := &Controller{
		config:                         &Configuration{EnableGwNftableLbSvc: true, EnableGwNftableSvcClusterIP: true},
		addOrUpdateGwNftableLbSvcQueue: queue,
		delIptablesDnatRuleQueue:       delQueue,
	}

	c.enqueueDelIptablesDnatRule(&kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "svc-record",
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: "default", util.NftableLbSvcNameLabel: "web", util.NftableLbSvcRecordLabel: "true",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{Type: kubeovnv1.DnatRuleTypeShare, EIP: "eip0"},
	})
	require.Equal(t, 1, queue.Len())
	key, _ := queue.Get()
	require.Equal(t, "default/web", key)
	queue.Done(key)

	// a hand-managed share record (no Service owner) only re-checks the EIP
	c.enqueueDelIptablesDnatRule(&kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "manual"},
		Spec:       kubeovnv1.IptablesDnatRuleSpec{Type: kubeovnv1.DnatRuleTypeShare, EIP: "eip0"},
	})
	// an exclusive rule goes to its own delete queue, never to Services
	c.enqueueDelIptablesDnatRule(&kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "exclusive"},
		Spec:       kubeovnv1.IptablesDnatRuleSpec{Type: kubeovnv1.DnatRuleTypeExclusive, EIP: "eip0"},
	})
	require.Equal(t, 0, queue.Len(), "only Service accounting records wake their owner Service")
}

// Test_handleAddOrUpdateGwNftableLbService_restoresDeletedRecords pins the recovery half of the
// accounting rule: after every accounting record of a Service is deleted out from under it, the
// next reconcile re-claims the records before it touches anything else, so claims never go
// missing for longer than one reconcile.
func Test_handleAddOrUpdateGwNftableLbService_restoresDeletedRecords(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{f.slice},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{f.eip},
		Pods:           gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.addOrUpdateGwNftableLbSvcQueue = workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(c.addOrUpdateGwNftableLbSvcQueue.ShutDown)
	c.resetIptablesEipQueue = workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(c.resetIptablesEipQueue.ShutDown)
	c.execRulesInPod = func(_ *v1.Pod, _ string, _ []string) error { return nil }
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	owner := f.namespace + "/" + f.svc.Name
	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(owner))
	claimed, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, claimed.Items, "the service claims its records on reconcile")

	// Delete every account record of the service, as an accident or a manual cleanup would, and
	// wait until the informer observes the loss.
	for _, item := range claimed.Items {
		require.NoError(t, c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Delete(context.Background(), item.Name, metav1.DeleteOptions{}))
	}
	require.Eventually(t, func() bool {
		rules, err := c.iptablesDnatRulesLister.List(labels.Everything())
		return err == nil && len(rules) == 0
	}, 5*time.Second, 10*time.Millisecond, "the informer must observe the record deletion")

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(owner))
	restored, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)
	require.Len(t, restored.Items, len(claimed.Items), "the reconcile restores the deleted claims")
	seen := set.New[string]()
	for _, item := range restored.Items {
		seen.Insert(item.Name)
	}
	for _, item := range claimed.Items {
		require.True(t, seen.Has(item.Name), "record %s must be re-claimed", item.Name)
	}
}

// honorFinalizersOnDnatDelete makes the fake clientset honor finalizers on delete like the API
// server does: an object with finalizers transitions to terminating instead of disappearing. It
// goes through the tracker directly because the fake's own mutex is not reentrant.
func honorFinalizersOnDnatDelete(client *kubeovnfake.Clientset) {
	tracker := client.Tracker()
	gvr := kubeovnv1.SchemeGroupVersion.WithResource("iptables-dnat-rules")
	client.PrependReactor("delete", "iptables-dnat-rules", func(action k8stesting.Action) (bool, runtime.Object, error) {
		name := action.(k8stesting.DeleteAction).GetName()
		obj, err := tracker.Get(gvr, "", name)
		if err != nil {
			return true, nil, err
		}
		rule, ok := obj.(*kubeovnv1.IptablesDnatRule)
		if !ok || len(rule.Finalizers) == 0 {
			// no finalizers: fall through to the default tracker delete, which drops the object
			return false, nil, nil
		}
		updated := rule.DeepCopy()
		now := metav1.Now()
		updated.DeletionTimestamp = &now
		return true, updated, tracker.Update(gvr, updated, "")
	})
}

// Test_claimNftableLbRecords_migratesTerminatingRecord pins the stuck-Terminating recovery for a
// record that re-entered desired: the claim keeps its identity but migrates the legacy controller
// finalizers off so the record finishes terminating and a later pass can recreate it.
func Test_claimNftableLbRecords_migratesTerminatingRecord(t *testing.T) {
	t.Parallel()

	// Seed through the typed client: the fake tracker's key for objects passed to
	// NewSimpleClientset differs from the one the typed client uses, which would make every
	// follow-up call a silent no-op.
	client := kubeovnfake.NewSimpleClientset()
	honorFinalizersOnDnatDelete(client)
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "legacy-record",
			Finalizers: []string{util.DepreciatedFinalizerName, util.KubeOVNControllerFinalizer},
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: "default", util.NftableLbSvcNameLabel: "web", util.NftableLbSvcRecordLabel: "true",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.0.5", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	if _, err := client.KubeovnV1().IptablesDnatRules().Create(context.Background(), record, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := client.KubeovnV1().IptablesDnatRules().Delete(context.Background(), record.Name, metav1.DeleteOptions{}); err != nil {
		t.Fatal(err)
	}
	terminating, err := client.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, terminating.DeletionTimestamp.IsZero(), "the finalizer-honoring delete keeps the record terminating")
	client.ClearActions()

	c := &Controller{config: &Configuration{KubeOvnClient: client}}
	svc := &v1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}}
	want := terminating.DeepCopy()
	want.DeletionTimestamp = nil
	want.Finalizers = nil
	want.ResourceVersion = ""

	live, retiring, err := c.claimNftableLbRecords(svc, map[string]*kubeovnv1.IptablesDnatRule{record.Name: want},
		[]*kubeovnv1.IptablesDnatRule{terminating})
	require.NoError(t, err)
	require.Empty(t, live, "a terminating record is not claimed this pass")
	require.Empty(t, retiring, "a terminating record is not snapshotted")

	actions := client.Actions()
	require.True(t, testHasAction(actions, "patch"),
		"the legacy finalizers were migrated off so the record can finish terminating")
	require.False(t, testHasAction(actions, "create"), "the stuck record's name is not recreated while it exists")
}

// Test_settleNftableLbRecords_freesLegacyFinalizedRecords pins the stale-retirement closing path:
// deleting a record that still carries the pre-upgrade controller finalizer does not leave it
// terminating forever - the Service is its only owner and clears the finalizer after the identity
// cleanup.
func Test_settleNftableLbRecords_freesLegacyFinalizedRecords(t *testing.T) {
	t.Parallel()

	client := kubeovnfake.NewSimpleClientset()
	honorFinalizersOnDnatDelete(client)
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "legacy-stale",
			Finalizers: []string{util.KubeOVNControllerFinalizer},
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: "default", util.NftableLbSvcNameLabel: "web", util.NftableLbSvcRecordLabel: "true",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.0.5", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	if _, err := client.KubeovnV1().IptablesDnatRules().Create(context.Background(), record, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}
	created, err := client.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	client.ClearActions()

	c := &Controller{config: &Configuration{KubeOvnClient: client}}
	require.NoError(t, c.settleNftableLbRecords("203.0.113.10", "gw0", nil,
		[]*kubeovnv1.IptablesDnatRule{created}, nil, nil))

	actions := client.Actions()
	require.True(t, testHasAction(actions, "delete"), "the stale record got its delete request")
	require.True(t, testHasAction(actions, "patch"), "the stuck record's legacy finalizer was migrated off after it")
	after, err := client.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, after.DeletionTimestamp.IsZero(), "the record is terminating")
	require.Empty(t, after.Finalizers, "the stuck record is freed once its identities were cleaned")
}

// testHasAction reports whether the fake clientset recorded an action with the given verb.
func testHasAction(actions []k8stesting.Action, verb string) bool {
	for _, action := range actions {
		if action.GetVerb() == verb {
			return true
		}
	}
	return false
}

// Test_cleanupNftableLbService_migratesLegacyFinalizers pins the closing path for a Service that
// never gained the annotations of the new design: its leftover share records from before the
// upgrade keep the controller finalizer the share DNAT workers no longer clear, and cleanup must
// migrate it off after removing the identities the records account for.
func Test_cleanupNftableLbService_migratesLegacyFinalizers(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	terminating := metav1.Now()
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}
	f.svc.DeletionTimestamp = &terminating
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "upgraded-record",
			Finalizers: []string{util.KubeOVNControllerFinalizer},
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: f.namespace, util.NftableLbSvcNameLabel: f.svc.Name,
				util.NftableLbSvcUIDLabel: string(f.svc.UID), util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel: f.gw.Name,
				util.EipV4IpLabel:           "172.20.0.5", util.EipUIDLabel: "owned-eip-uid",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "owned-eip", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.7.2", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:              []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:           []*kubeovnv1.Subnet{f.subnet},
		Services:          []*v1.Service{f.svc},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{f.eip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{record},
		Pods:              gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.execRulesInPod = func(_ *v1.Pod, _ string, _ []string) error { return nil }
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()
	client, ok := c.config.KubeOvnClient.(*kubeovnfake.Clientset)
	require.True(t, ok)
	honorFinalizersOnDnatDelete(client)

	require.NoError(t, c.cleanupNftableLbService(f.svc, f.namespace, f.svc.Name))

	kept, err := client.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, kept.DeletionTimestamp.IsZero(), "the record is terminating")
	require.Empty(t, kept.Finalizers, "the legacy finalizer is migrated off after the identity cleanup, freeing the record")
	cleaned, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, cleaned.Finalizers, "the Service itself is released")
}

// Test_claimNftableLbRecords_snapshotsRetiringIdentity pins the claim lifecycle for an in-place
// update that drops an identity: the closing ClusterIP gate strips Spec.ClusterIP from the
// record, but the old nft map and hairpin are removed only by the program phase that follows. A
// failed pass would retry with a ledger that no longer knows the old identity (Status.V4ip
// records only the EIP when both legs share one record). The claim must therefore keep a
// retiring snapshot of the old content until the data plane converged, and a retried claim must
// not pile up further snapshots.
func Test_claimNftableLbRecords_snapshotsRetiringIdentity(t *testing.T) {
	t.Parallel()

	client := kubeovnfake.NewSimpleClientset()
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gate-flip-record",
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: "default", util.NftableLbSvcNameLabel: "web", util.NftableLbSvcRecordLabel: "true",
				util.EipV4IpLabel: "203.0.113.10", util.EipUIDLabel: "eip0-uid",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.0.5", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	current, err := client.KubeovnV1().IptablesDnatRules().Create(context.Background(), record, metav1.CreateOptions{})
	require.NoError(t, err)

	want := record.DeepCopy()
	want.Spec.ClusterIP = "" // the ClusterIP gate closed; the EIP leg stays desired
	want.ResourceVersion = ""

	c := &Controller{config: &Configuration{KubeOvnClient: client}}
	svc := &v1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}}
	live, retiring, err := c.claimNftableLbRecords(svc, map[string]*kubeovnv1.IptablesDnatRule{record.Name: want},
		[]*kubeovnv1.IptablesDnatRule{current})
	require.NoError(t, err)
	require.Contains(t, live, record.Name)
	require.Len(t, retiring, 1, "an identity-dropping update snapshots the retired content")

	snapshot, err := client.KubeovnV1().IptablesDnatRules().Get(context.Background(), retiring[0], metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "10.96.1.5", snapshot.Spec.ClusterIP, "the snapshot keeps the retired identity")
	require.Equal(t, "eip0-uid", snapshot.Labels[util.EipUIDLabel], "the snapshot keeps the old EIP claim until cleanup")
	require.Empty(t, snapshot.Finalizers, "the snapshot needs no finalizer migration")

	updated := live[record.Name]
	require.Empty(t, updated.Spec.ClusterIP, "the accounting record itself converges to desired")

	// The retry view after the failed program phase: the updated record alone has lost the old
	// ClusterIP, the snapshot restores it for the cleanup derivation.
	identities := nftableLbExistingIdentities([]*kubeovnv1.IptablesDnatRule{updated, snapshot})
	require.Contains(t, identities, "10.96.1.5/80/tcp", "a retried pass must still find the retiring ClusterIP identity")

	// A reconciled claim (specs equal) creates no further snapshot.
	_, retiring, err = c.claimNftableLbRecords(svc, map[string]*kubeovnv1.IptablesDnatRule{record.Name: want},
		[]*kubeovnv1.IptablesDnatRule{updated})
	require.NoError(t, err)
	require.Empty(t, retiring, "a converged record is not snapshotted again")
}

// Test_settleNftableLbRecords_retiresSnapshot pins the release order: the retiring snapshot's
// delete only happens in settle, after the gateway cleanup of the identities it preserved, so
// its old EIP UID claim outlives the rule it covered. Settle tolerates an already-deleted
// snapshot (the delete of a crashed pass is replayed).
func Test_settleNftableLbRecords_retiresSnapshot(t *testing.T) {
	t.Parallel()

	client := kubeovnfake.NewSimpleClientset()
	snapshot := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gate-flip-record-r12345678",
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: "default", util.NftableLbSvcNameLabel: "web", util.NftableLbSvcRecordLabel: "true",
				util.EipV4IpLabel: "203.0.113.10", util.EipUIDLabel: "eip0-uid",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.0.5", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	if _, err := client.KubeovnV1().IptablesDnatRules().Create(context.Background(), snapshot, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}

	c := &Controller{config: &Configuration{KubeOvnClient: client}}
	require.NoError(t, c.settleNftableLbRecords("203.0.113.10", "gw0", nil, nil, nil,
		[]string{snapshot.Name, "already-gone-snapshot"}))
	_, err := client.KubeovnV1().IptablesDnatRules().Get(context.Background(), snapshot.Name, metav1.GetOptions{})
	require.True(t, k8serrors.IsNotFound(err), "the retiring snapshot is deleted once the data plane converged")
}

// Test_handleUpdateIptablesEip_releasesOnceNoRecordsClaim pins the settled release path: with a
// terminating EIP whose records are gone (the owning Service already ran its cleanup, which
// tears down the data plane even without records), the EIP reconcile must finish the release -
// no hold may spin on the Service's leftover annotation. The Service itself is kept out of the
// way: it still declares the EIP, which is exactly the input a wrong hold would spin on.
func Test_handleUpdateIptablesEip_releasesOnceNoRecordsClaim(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	terminating := metav1.Now()
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name: "released-eip", UID: "released-eip-uid",
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
			DeletionTimestamp: &terminating,
		},
		Spec:   kubeovnv1.IptablesEIPSpec{NatGwDp: f.gw.Name, V4ip: "172.20.0.5"},
		Status: kubeovnv1.IptablesEIPStatus{IP: "172.20.0.5", Ready: true},
	}
	// The annotated Service stays live and unchanged, as after its completed cleanup.
	f.svc.Annotations[util.EipAnnotation] = eip.Name

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.updateIptablesEipQueue = workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)

	require.NoError(t, c.handleUpdateIptablesEip(eip.Name))
	released, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(context.Background(), eip.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, released.Finalizers, "with no records left the terminating EIP finishes its release")
}

// Test_handleAddOrUpdateGwNftableLbService_restoresTrimmedIdentityEvidence pins the trim path
// under record loss: the Service's records were deleted out of band and its EIP annotation was
// removed before the restore pass ran. The reconcile still serves the ClusterIP, but the old
// EIP identity (witnessed by the published ingress IP) must be torn down in the same pass, the
// ingress withdrawn, and - if the gateway delete fails - a retried pass must still find the old
// identity, because the evidence was persisted before the claim narrowed the ledger.
func Test_handleAddOrUpdateGwNftableLbService_restoresTrimmedIdentityEvidence(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	// Records are gone (never seeded), the annotation was removed in the gap, and the ingress
	// still publishes the EIP address the Service once served.
	delete(f.svc.Annotations, util.EipAnnotation)
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}
	f.svc.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{{IP: "172.20.0.5"}}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{f.slice},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{f.eip},
		Pods:           gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true

	var deleteCalls [][]string
	failFirstDelete := true
	c.execRulesInPod = func(_ *v1.Pod, operation string, rules []string) error {
		if operation == natGwNftDnatMapDel {
			deleteCalls = append(deleteCalls, append([]string(nil), rules...))
			if failFirstDelete {
				failFirstDelete = false
				return errors.New("injected gateway failure")
			}
		}
		return nil
	}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	key := f.namespace + "/" + f.svc.Name
	require.Error(t, c.handleAddOrUpdateGwNftableLbService(key), "the injected delete failure aborts the first pass")
	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(key), "the retry replays the trimmed identity from the persisted evidence")

	require.Len(t, deleteCalls, 2, "both passes attempt the old EIP identity delete")
	for _, rules := range deleteCalls {
		require.Contains(t, rules, "172.20.0.5,80,tcp", "the trimmed EIP identity is torn down even after the first pass failed")
	}

	svc, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, svc.Status.LoadBalancer.Ingress, "the stale ingress address is withdrawn")

	rules, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)
	require.Len(t, rules.Items, 1, "only the still-served ClusterIP record remains; the evidence was retired in settle")
	restored := rules.Items[0]
	require.Equal(t, "10.96.1.5", restored.Spec.ClusterIP)
	require.Empty(t, restored.Spec.EIP)
	require.Empty(t, restored.Labels[util.EipUIDLabel], "the released EIP claim is not resurrected")
}

// Test_handleUpdateIptablesEip_holdsWhileServiceTeardownPending pins the EIP-vs-Service release
// order: with the accounting records gone and the ingress never published, the terminating EIP
// object is the only remaining witness of the address a referencing Service still has to tear
// down. The EIP finalizer must hold while the Service carries the feature finalizer; once the
// Service cleanup ran (records, data plane and finally its finalizer), the EIP release proceeds
// in a later pass.
func Test_handleUpdateIptablesEip_holdsWhileServiceTeardownPending(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}
	terminating := metav1.Now()
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name: f.eip.Name, UID: f.eip.UID,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
			DeletionTimestamp: &terminating,
		},
		Spec:   f.eip.Spec,
		Status: f.eip.Status,
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{eip},
		Pods:           gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.updateIptablesEipQueue = workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)
	c.execRulesInPod = func(_ *v1.Pod, _ string, _ []string) error { return nil }
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()

	// No records, no published ingress: the EIP object itself is the teardown witness.
	require.NoError(t, c.handleUpdateIptablesEip(eip.Name))
	held, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(context.Background(), eip.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Contains(t, held.Finalizers, util.KubeOVNControllerFinalizer,
		"the EIP release holds while the Service's teardown is pending")

	// The Service reconcile runs the terminating-EIP cleanup: the held object still supplies the
	// address, so the teardown ledger removes the EIP identity before the finalizer drops.
	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+f.svc.Name))
	require.Eventually(t, func() bool {
		svc, err := c.servicesLister.Services(f.namespace).Get(f.svc.Name)
		return err == nil && !slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer)
	}, 5*time.Second, 10*time.Millisecond, "the informer observes the Service finalizer release")

	require.NoError(t, c.handleUpdateIptablesEip(eip.Name))
	released, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(context.Background(), eip.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, released.Finalizers, "the EIP finishes its release once the Service settled")
}

// Test_handleUpdateIptablesEip_releasedWhenServiceServesOnlyClusterIP pins the hold's scope: a
// Service whose EIP leg is disabled keeps its controller finalizer for the ClusterIP leg
// forever, but owes a terminating EIP nothing once the trim retired the EIP identity from its
// records. Repeated ClusterIP-only reconciles must not wedge the EIP release.
func Test_handleUpdateIptablesEip_releasedWhenServiceServesOnlyClusterIP(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{f.slice},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{f.eip},
		Pods:           gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	// The EIP leg gate is off: the Service keeps serving its ClusterIP and never again enters
	// the EIP path, so its finalizer stays.
	c.config.EnableGwNftableLbSvc = false
	c.config.EnableGwNftableSvcClusterIP = true
	c.updateIptablesEipQueue = workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(c.updateIptablesEipQueue.ShutDown)
	c.execRulesInPod = func(_ *v1.Pod, _ string, _ []string) error { return nil }
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	key := f.namespace + "/" + f.svc.Name
	for i := 0; i < 3; i++ {
		require.NoError(t, c.handleAddOrUpdateGwNftableLbService(key))
		require.Eventually(t, func() bool {
			rules, err := c.iptablesDnatRulesLister.List(labels.Everything())
			if err != nil || len(rules) != 1 {
				return false
			}
			// Settled: only the ClusterIP claim survives, the trim tombstone is retired.
			return rules[0].Spec.EIP == "" && rules[0].Labels[util.EipV4IpLabel] == ""
		}, 5*time.Second, 10*time.Millisecond, "the informer observes the settled ledger between passes")
	}

	terminating := metav1.Now()
	eip, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(context.Background(), f.eip.Name, metav1.GetOptions{})
	require.NoError(t, err)
	eip.Finalizers = []string{util.KubeOVNControllerFinalizer}
	eip.DeletionTimestamp = &terminating
	if _, err = c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Update(context.Background(), eip, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}
	require.Eventually(t, func() bool {
		cached, err := c.iptablesEipsLister.Get(eip.Name)
		return err == nil && !cached.DeletionTimestamp.IsZero()
	}, 5*time.Second, 10*time.Millisecond, "the informer observes the terminating EIP")

	require.NoError(t, c.handleUpdateIptablesEip(eip.Name))
	released, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Get(context.Background(), eip.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, released.Finalizers, "a Service owing no EIP leg must not hold the release, finalizer or not")
}

// Test_handleAddOrUpdateGwNftableLbService_noTombstoneChurnAfterTrim pins the steady state after
// a trim: a Service whose EIP leg is disabled and whose ingress is already withdrawn must not
// recreate evidence records on every reconcile - a tombstone's settle deletion wakes the owner,
// which would recreate it forever (and transiently hold a terminating EIP via its EIP label).
// The pass touches only the still-served ClusterIP claim.
func Test_handleAddOrUpdateGwNftableLbService_noTombstoneChurnAfterTrim(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	delete(f.svc.Annotations, util.EipAnnotation)
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}
	// No ingress published: either the EIP leg never served (gate always off) or its trim
	// already completed. Neither can owe the gateway an EIP identity.

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{f.slice},
		Pods:           gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = false
	c.config.EnableGwNftableSvcClusterIP = true
	c.execRulesInPod = func(_ *v1.Pod, _ string, _ []string) error { return nil }
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+f.svc.Name))
	require.Eventually(t, func() bool {
		rules, err := c.iptablesDnatRulesLister.List(labels.Everything())
		return err == nil && len(rules) == 1
	}, 5*time.Second, 10*time.Millisecond, "the informer observes the claimed record")

	client, ok := c.config.KubeOvnClient.(*kubeovnfake.Clientset)
	require.True(t, ok)
	client.ClearActions()
	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+f.svc.Name))
	for _, action := range client.Actions() {
		require.NotEqual(t, "create", action.GetVerb(), "a settled ledger must not recreate evidence records")
		require.NotEqual(t, "delete", action.GetVerb(), "a settled ledger has no tombstones or stale records to delete")
	}
}
