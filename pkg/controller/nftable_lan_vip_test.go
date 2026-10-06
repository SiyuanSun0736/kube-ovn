package controller

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// lanVipFixture builds one gateway (lanIP 10.0.7.254 in vpc1) with three bound Services:
// svc-a (LoadBalancer, ports 80/tcp ClientIP-300s and 53/udp), svc-b (ClusterIP, port 80/tcp
// without affinity), and svc-c (terminating, port 80/tcp, must never contribute).
type lanVipFixture struct {
	vpc    *kubeovnv1.Vpc
	subnet *kubeovnv1.Subnet
	gw     *kubeovnv1.VpcNatGateway
	svcs   []*v1.Service
	slices []*discoveryv1.EndpointSlice
}

func newLanVipFixture() *lanVipFixture {
	const (
		gwName     = "lan-gw"
		subnetName = "lan-subnet"
		vpcName    = "vpc1"
		lanIP      = "10.0.7.254"
	)
	namespace := metav1.NamespaceDefault
	terminating := metav1.Now()

	svc := func(name, clusterIP string, svcType v1.ServiceType, ports []v1.ServicePort) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace, Name: name,
				Annotations: map[string]string{util.VpcNatGatewayAnnotation: gwName},
			},
			Spec: v1.ServiceSpec{Type: svcType, ClusterIP: clusterIP, Ports: ports},
		}
	}
	timeout := int32(300)
	svcA := svc("svc-a", "10.96.1.10", v1.ServiceTypeLoadBalancer, []v1.ServicePort{
		// shared with svc-b below: the 80/tcp identity must merge both backend sets
		// and resolve the affinity disagreement to ClientIP with svc-a's timeout
		{Name: "http", Port: 80, Protocol: v1.ProtocolTCP},
		{Name: "dns", Port: 53, Protocol: v1.ProtocolUDP},
	})
	svcA.Spec.SessionAffinity = v1.ServiceAffinityClientIP
	svcA.Spec.SessionAffinityConfig = &v1.SessionAffinityConfig{ClientIP: &v1.ClientIPConfig{TimeoutSeconds: &timeout}}

	svcB := svc("svc-b", "10.96.2.10", v1.ServiceTypeClusterIP, []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}})

	svcC := svc("svc-c", "10.96.3.10", v1.ServiceTypeClusterIP, []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}})
	svcC.Finalizers = []string{util.KubeOVNControllerFinalizer}
	svcC.DeletionTimestamp = &terminating

	// An unbound Service must never join the partition, even on a colliding port.
	svcD := svc("svc-d", "10.96.4.10", v1.ServiceTypeClusterIP, []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}})
	delete(svcD.Annotations, util.VpcNatGatewayAnnotation)

	// EndpointSlice addresses are shared by every port in the slice, so each Service port gets
	// its own slice here (which is also how the endpointslice controller lays them out).
	slice := func(svcName, portName string, port int32, protocol v1.Protocol, addresses ...string) *discoveryv1.EndpointSlice {
		eps := make([]discoveryv1.Endpoint, 0, len(addresses))
		for _, addr := range addresses {
			ready := true
			if addr == "10.0.7.22" {
				ready = false
			}
			eps = append(eps, discoveryv1.Endpoint{Addresses: []string{addr}, Conditions: discoveryv1.EndpointConditions{Ready: &ready}})
		}
		return &discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace, Name: svcName + "-" + portName,
				Labels: map[string]string{discoveryv1.LabelServiceName: svcName},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Ports:       []discoveryv1.EndpointPort{{Name: new(portName), Port: new(port), Protocol: new(protocol)}},
			Endpoints:   eps,
		}
	}
	return &lanVipFixture{
		vpc:    &kubeovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: vpcName}},
		subnet: &kubeovnv1.Subnet{ObjectMeta: metav1.ObjectMeta{Name: subnetName}, Spec: kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Protocol: kubeovnv1.ProtocolIPv4, Vpc: vpcName, CIDRBlock: "10.0.7.0/24", Gateway: "10.0.7.1"}},
		gw:     &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: gwName, UID: "lan-gw-uid"}, Spec: kubeovnv1.VpcNatGatewaySpec{Vpc: vpcName, Subnet: subnetName, LanIP: lanIP}},
		svcs:   []*v1.Service{svcA, svcB, svcC, svcD},
		slices: []*discoveryv1.EndpointSlice{
			slice("svc-a", "http", 8080, v1.ProtocolTCP, "10.0.7.11", "10.0.7.12"),
			slice("svc-a", "dns", 53, v1.ProtocolUDP, "10.0.7.13"),
			slice("svc-b", "http", 8081, v1.ProtocolTCP, "10.0.7.21", "10.0.7.22"), // second endpoint not ready
		},
	}
}

func newLanVipController(t *testing.T, f *lanVipFixture, enabled bool) *fakeController {
	t.Helper()
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       f.svcs,
		EndpointSlices: f.slices,
		Pods:           gatewayPods(f.gw.Name, f.gw.Spec.LanIP),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLanipVip = enabled
	c.natGwLanVipSyncQueue = newTypedRateLimitingQueue[string]("test-lanvip-sync", nil)
	t.Cleanup(c.natGwLanVipSyncQueue.ShutDown)
	return fc
}

func TestDesiredNatGwLanVipRulesMergesSharedPorts(t *testing.T) {
	f := newLanVipFixture()
	fc := newLanVipController(t, f, true)

	rules, err := fc.fakeController.desiredNatGwLanVipRules(f.gw)
	require.NoError(t, err)
	require.Equal(t, []string{
		// Service-level affinity (like kube-proxy): both identities of svc-a carry its
		// ClientIP setting; the shared 80/tcp port additionally merges svc-b's backends,
		// and the terminating svc-c contributes nothing
		"10.0.7.254,53,udp,clientip,300,10.0.7.13:53",
		"10.0.7.254,80,tcp,clientip,300,10.0.7.11:8080@10.0.7.12:8080@10.0.7.21:8081",
	}, rules)
}

func TestDesiredNatGwLanVipRulesSkipsBackendlessPorts(t *testing.T) {
	f := newLanVipFixture()
	// nobody ready serves port 80 any more, port 53 of svc-a remains
	f.slices[0].Endpoints = nil
	f.slices[2].Endpoints = nil
	fc := newLanVipController(t, f, true)

	rules, err := fc.fakeController.desiredNatGwLanVipRules(f.gw)
	require.NoError(t, err)
	require.Equal(t, []string{"10.0.7.254,53,udp,clientip,300,10.0.7.13:53"}, rules,
		"a port with zero ready backends produces no rule (the gateway-side sync GCs it)")
}

func TestHandleSyncNatGwLanVipProgramsPartition(t *testing.T) {
	f := newLanVipFixture()
	fc := newLanVipController(t, f, true)
	c := fc.fakeController

	type execCall struct {
		op    string
		rules []string
	}
	var calls []execCall
	c.execRulesInPod = func(_ *v1.Pod, operation string, rules []string) error {
		calls = append(calls, execCall{operation, rules})
		return nil
	}

	require.NoError(t, c.handleSyncNatGwLanVip(f.gw.Name))
	require.Len(t, calls, 1, "one batched exec for the single gateway instance")
	require.Equal(t, natGwNftLanVipSync, calls[0].op)
	require.Len(t, calls[0].rules, 2)

	// the partition state is informer-derived; re-syncing an unchanged instance+rule set would
	// program byte-identical nft state, so the redundant gateway exec is skipped
	require.NoError(t, c.handleSyncNatGwLanVip(f.gw.Name))
	require.Len(t, calls, 1, "unchanged partition skips the redundant gateway exec")

	// a recreated gateway instance (new pod UID) invalidates the cache and re-programs
	// (the gateway pods are read through the kube client, so the fake one is updated)
	pod := gatewayPods(f.gw.Name, f.gw.Spec.LanIP)[0]
	pod.UID = "recreated-uid"
	_, err := c.config.KubeClient.CoreV1().Pods(pod.Namespace).Update(context.Background(), pod, metav1.UpdateOptions{})
	require.NoError(t, err)
	require.NoError(t, c.handleSyncNatGwLanVip(f.gw.Name))
	require.Len(t, calls, 2, "new gateway instance forces re-program")

	// intent changes re-program too (deleting svc-b shrinks the shared 80/tcp identity to
	// svc-a's backends), again only once per actual change
	require.NoError(t, fc.fakeInformers.serviceInformer.Informer().GetIndexer().Delete(f.svcs[1]))
	require.NoError(t, c.handleSyncNatGwLanVip(f.gw.Name))
	require.Len(t, calls, 3, "intent change forces re-program")
	require.NotEqual(t, calls[1].rules, calls[2].rules)
	require.NoError(t, c.handleSyncNatGwLanVip(f.gw.Name))
	require.Len(t, calls, 3, "and the identical follow-up sync is skipped again")
}

func TestHandleSyncNatGwLanVipFlagOffWipesPartition(t *testing.T) {
	f := newLanVipFixture()
	fc := newLanVipController(t, f, false)
	c := fc.fakeController

	var rules []string
	c.execRulesInPod = func(_ *v1.Pod, _ string, r []string) error {
		rules = r
		return nil
	}
	require.NoError(t, c.handleSyncNatGwLanVip(f.gw.Name))
	require.Empty(t, rules, "with the feature disabled the sync wipes the partition")
}

// TestLanVipEnqueueWiring pins the event glue of the feature: Service add/update/delete and
// EndpointSlice events must wake the partition reconciliation of the gateways involved (both
// gateways when the annotation moves), or the merged identities would go stale.
func TestLanVipEnqueueWiring(t *testing.T) {
	f := newLanVipFixture()
	fc := newLanVipController(t, f, true)
	c := fc.fakeController
	// the enqueue handlers fan out to several queues; create the ones their paths touch
	c.addOrUpdateEndpointSliceQueue = newTypedRateLimitingQueue[string]("test-eps", nil)
	c.deleteServiceQueue = newTypedRateLimitingQueue[*vpcService]("test-del-svc", nil)
	c.updateServiceQueue = newTypedRateLimitingQueue[*updateSvcObject]("test-upd-svc", nil)
	t.Cleanup(func() {
		c.addOrUpdateEndpointSliceQueue.ShutDown()
		c.deleteServiceQueue.ShutDown()
		c.updateServiceQueue.ShutDown()
	})

	// the workqueue dedups repeated keys, so drain after every step and compare the set
	drain := func() []string {
		var got []string
		for c.natGwLanVipSyncQueue.Len() > 0 {
			item, _ := c.natGwLanVipSyncQueue.Get()
			got = append(got, item)
			c.natGwLanVipSyncQueue.Done(item)
			c.natGwLanVipSyncQueue.Forget(item)
		}
		slices.Sort(got)
		return got
	}

	svcA := f.svcs[0]
	c.enqueueAddService(svcA)
	require.Equal(t, []string{"lan-gw"}, drain(), "service add wakes its gateway")

	moved := svcA.DeepCopy()
	moved.ResourceVersion = "2"
	moved.Annotations[util.VpcNatGatewayAnnotation] = "other-gw"
	moved.Spec.Ports = append(slices.Clone(svcA.Spec.Ports), v1.ServicePort{Name: "extra", Port: 8080, Protocol: v1.ProtocolTCP})
	c.enqueueUpdateService(svcA, moved)
	require.Equal(t, []string{"lan-gw", "other-gw"}, drain(), "update wakes the old and the new gateway")

	c.enqueueDeleteService(svcA)
	require.Equal(t, []string{"lan-gw"}, drain(), "delete (tombstone annotation) wakes the gateway")

	c.enqueueAddEndpointSlice(f.slices[0])
	require.Equal(t, []string{"lan-gw"}, drain(), "endpoint slice add wakes the owning gateway")

	c.enqueueDeleteEndpointSlice(f.slices[0])
	require.Equal(t, []string{"lan-gw"}, drain(), "endpoint slice delete wakes the owning gateway")

	// an unbound Service touches no queue
	c.enqueueDeleteService(f.svcs[3])
	require.Empty(t, drain())
}

func TestHandleSyncNatGwLanVipSkipsWhenGatewayUnavailable(t *testing.T) {
	f := newLanVipFixture()

	// gateway object gone: nothing to do, no error
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{})
	require.NoError(t, err)
	require.NoError(t, fc.fakeController.handleSyncNatGwLanVip("ghost-gw"))

	// no lanIP yet (backfill pending): skipped quietly, a gateway update re-enqueues the sync
	fNoIP := newLanVipFixture()
	fNoIP.gw.Spec.LanIP = ""
	fc = newLanVipController(t, fNoIP, true)
	executed := false
	fc.fakeController.execRulesInPod = func(_ *v1.Pod, _ string, _ []string) error {
		executed = true
		return nil
	}
	require.NoError(t, fc.fakeController.handleSyncNatGwLanVip(f.gw.Name))
	require.False(t, executed)

	// no running gateway instance: retried later instead of programming into thin air
	fNoPods := newLanVipFixture()
	op := &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{fNoPods.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fNoPods.gw},
		Subnets:        []*kubeovnv1.Subnet{fNoPods.subnet},
		Services:       fNoPods.svcs,
		EndpointSlices: fNoPods.slices,
	}
	fc, err = newFakeControllerWithOptions(t, op)
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLanipVip = true
	c.natGwLanVipSyncQueue = newTypedRateLimitingQueue[string]("test-lanvip-sync-empty", nil)
	t.Cleanup(c.natGwLanVipSyncQueue.ShutDown)
	c.execRulesInPod = func(_ *v1.Pod, _ string, _ []string) error {
		executed = true
		return nil
	}
	require.NoError(t, c.handleSyncNatGwLanVip(f.gw.Name))
	require.False(t, executed, "with no running gateway instance nothing is programmed")
}

func drainEvents(recorder record.EventRecorder) []string {
	faker, ok := recorder.(*record.FakeRecorder)
	if !ok {
		return nil
	}
	var events []string
	for {
		select {
		case e := <-faker.Events:
			events = append(events, e)
		default:
			return events
		}
	}
}

// A shared-port affinity disagreement must merge deterministically: the seed contributor and
// the resulting identities cannot depend on the informer index iteration order (random per
// call), and a steady disagreement must not re-fire its Warning event on every pass.
func TestDesiredNatGwLanVipRulesDeterministicMergeAttribution(t *testing.T) {
	f := newLanVipFixture()
	fc := newLanVipController(t, f, true)
	c := fc.fakeController

	reference, err := c.desiredNatGwLanVipRules(f.gw)
	require.NoError(t, err)
	events := drainEvents(c.recorder)
	require.Len(t, events, 1, "exactly one affinity-merge event: svc-b merges into svc-a on 80/tcp")
	require.Contains(t, events[0], "NatGwLanVipAffinityMerged")
	// the message carries the merged Service's own summary (svc-b: no affinity) and names the
	// seed contributor (svc-a: sorted first) — both are order-independent now
	require.Contains(t, events[0], "clientIP=false timeout=0", "the event reports svc-b's summary")
	require.Contains(t, events[0], "default/svc-a", "the seed contributor is the sorted-first Service")

	for i := 0; i < 10; i++ {
		rules, err := c.desiredNatGwLanVipRules(f.gw)
		require.NoError(t, err)
		require.Equal(t, reference, rules, "iteration %d: merge output must be deterministic", i)
	}
	require.Empty(t, drainEvents(c.recorder), "a steady disagreement re-emits nothing")
}

// The merge warning fires on changes of the gateway's disagreement set only: an unchanged
// disagreement stays silent, a resolved disagreement resets the state, and re-introducing it
// reports exactly once more.
func TestNatGwLanVipAffinityMergeEventOnlyOnChange(t *testing.T) {
	f := newLanVipFixture()
	fc := newLanVipController(t, f, true)
	c := fc.fakeController
	updateSvc := func(svc *v1.Service) {
		t.Helper()
		require.NoError(t, fc.fakeInformers.serviceInformer.Informer().GetIndexer().Update(svc))
	}

	_, err := c.desiredNatGwLanVipRules(f.gw)
	require.NoError(t, err)
	require.Len(t, drainEvents(c.recorder), 1, "first emission for the fresh disagreement")

	for i := 0; i < 3; i++ {
		_, err := c.desiredNatGwLanVipRules(f.gw)
		require.NoError(t, err)
	}
	require.Empty(t, drainEvents(c.recorder), "unchanged disagreement stays silent")

	// Resolving the disagreement (svc-b adopts the identity's ClientIP/300s) resets the state.
	svcB := f.svcs[1].DeepCopy()
	svcB.Spec.SessionAffinity = v1.ServiceAffinityClientIP
	svcB.Spec.SessionAffinityConfig = f.svcs[0].Spec.SessionAffinityConfig.DeepCopy()
	updateSvc(svcB)
	_, err = c.desiredNatGwLanVipRules(f.gw)
	require.NoError(t, err)
	require.Empty(t, drainEvents(c.recorder), "resolved disagreement emits nothing")

	// Re-introducing the disagreement is a change again and reports exactly once more.
	svcB.Spec.SessionAffinity = v1.ServiceAffinityNone
	svcB.Spec.SessionAffinityConfig = nil
	updateSvc(svcB)
	_, err = c.desiredNatGwLanVipRules(f.gw)
	require.NoError(t, err)
	require.Len(t, drainEvents(c.recorder), 1, "a new disagreement fires exactly once")
	_, err = c.desiredNatGwLanVipRules(f.gw)
	require.NoError(t, err)
	require.Empty(t, drainEvents(c.recorder), "and it goes silent again afterwards")
}

// Deleting a gateway must also drop its lanVIP controller-side caches: nothing re-enqueues the
// partition sync afterwards (the Services leave the gateway index), so without the explicit
// cleanup the maps would grow by one entry per gateway ever created.
func TestHandleDelVpcNatGwClearsLanVipCaches(t *testing.T) {
	f := newLanVipFixture()
	fc := newLanVipController(t, f, true)
	c := fc.fakeController

	c.natGwLanVipMergeNotes.Store(f.gw.Name, "80/tcp|ns/svc|clientIP=true timeout=300")
	c.natGwLanVipLastSync.Store(f.gw.Name, "uid1\x00rules")

	require.NoError(t, c.handleDelVpcNatGw("kube-system/"+f.gw.Name))
	_, ok := c.natGwLanVipMergeNotes.Load(f.gw.Name)
	require.False(t, ok, "merge notes forgotten on gateway delete")
	_, ok = c.natGwLanVipLastSync.Load(f.gw.Name)
	require.False(t, ok, "last-sync signature forgotten on gateway delete")
}
