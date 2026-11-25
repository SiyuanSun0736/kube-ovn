package controller

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	listers "github.com/kubeovn/kube-ovn/pkg/client/listers/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestNftDnatMapAddRule(t *testing.T) {
	t.Parallel()

	t.Run("stateless rule sorts and dedups backends", func(t *testing.T) {
		rule, err := nftDnatMapAddRule("tcp", "10.0.0.1", "80",
			[]string{"10.0.0.6:8080", "10.0.0.5:8080", "10.0.0.6:8080"}, "", 0)
		require.NoError(t, err)
		require.Equal(t, "10.0.0.1,80,tcp,none,0,10.0.0.5:8080@10.0.0.6:8080", rule)
	})

	t.Run("client ip affinity carries its timeout", func(t *testing.T) {
		rule, err := nftDnatMapAddRule("udp", "10.0.0.2", "53",
			[]string{"10.0.0.5:5353"}, kubeovnv1.DnatSessionAffinityClientIP, 30)
		require.NoError(t, err)
		require.Equal(t, "10.0.0.2,53,udp,clientip,30,10.0.0.5:5353", rule)
	})

	t.Run("client ip affinity falls back to the default timeout", func(t *testing.T) {
		rule, err := nftDnatMapAddRule("tcp", "10.0.0.3", "443",
			[]string{"10.0.0.5:8443"}, kubeovnv1.DnatSessionAffinityClientIP, 0)
		require.NoError(t, err)
		require.Equal(t, "10.0.0.3,443,tcp,clientip,"+
			"10800,10.0.0.5:8443", rule)
	})

	t.Run("empty vip is rejected", func(t *testing.T) {
		_, err := nftDnatMapAddRule("tcp", "", "80", []string{"10.0.0.5:8080"}, "", 0)
		require.ErrorContains(t, err, "empty IPv4 EIP")
	})

	t.Run("backend less rule is rejected", func(t *testing.T) {
		_, err := nftDnatMapAddRule("tcp", "10.0.0.1", "80", nil, "", 0)
		require.ErrorContains(t, err, "no backends")
	})
}

func TestNftDnatMapDelRule(t *testing.T) {
	t.Parallel()
	require.Equal(t, "10.0.0.1,80,tcp", nftDnatMapDelRule("tcp", "10.0.0.1", "80"))
}

// TestProgramNftableLbServiceIdentitiesReprogramsUnchangedRecords pins the instance-recovery
// contract: the accounting records describe the Service's desired state, not the state of the
// current gateway instance. When a gateway Pod is replaced the records still match the desired
// state, and the Service is re-enqueued precisely so the fresh, empty instance is programmed
// again; a record-based skip would leave it without maps or hairpin rules.
func TestProgramNftableLbServiceIdentitiesReprogramsUnchangedRecords(t *testing.T) {
	t.Parallel()

	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "svc.80.tcp.10.0.0.5.8080",
			Labels: map[string]string{
				util.VpcNatGatewayNameLabel: "gw",
				util.VpcDnatEPortLabel:      "80",
			},
		},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			ClusterIP:    "10.96.0.10",
			ExternalPort: "80",
			InternalIP:   "10.0.0.5",
			InternalPort: "8080",
			Protocol:     "tcp",
			Type:         kubeovnv1.DnatRuleTypeShare,
		},
	}

	programmed := map[string][]string{}
	c := &Controller{
		config: &Configuration{},
		execRulesInPod: func(_ *v1.Pod, operation string, rules []string) error {
			programmed[operation] = append(programmed[operation], rules...)
			return nil
		},
	}

	pods := []*v1.Pod{{ObjectMeta: metav1.ObjectMeta{Name: "gw-0", Namespace: metav1.NamespaceSystem}}}
	identities := sortedNftableLbIdentities(map[string]*kubeovnv1.IptablesDnatRule{record.Name: record}, "")
	// The records already match the desired state, as they do right after a gateway instance is
	// replaced (the reconcile that would refresh them has not run yet).
	require.NoError(t, c.programNftableLbServiceIdentities("", "", pods, identities, []*kubeovnv1.IptablesDnatRule{record}, nil))

	require.Equal(t, []string{"10.96.0.10,80,tcp,none,0,10.0.0.5:8080"}, programmed[natGwNftDnatMapAdd],
		"the identity must be programmed again for the new gateway instance")
	require.Empty(t, programmed[natGwNftDnatMapDel], "an identity the Service still wants must not be deleted")
}

func TestProgramNftableLbServiceIdentities_Stateless(t *testing.T) {
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "dnat-record",
			Labels: map[string]string{
				util.NftableLbSvcNsLabel:   "default",
				util.NftableLbSvcNameLabel: "web",
			},
		},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP:          "eip1",
			ExternalPort: "80",
			Protocol:     "tcp",
			InternalIP:   "10.0.0.5",
			InternalPort: "8080",
		},
		Status: kubeovnv1.IptablesDnatRuleStatus{
			V4ip: "1.1.1.1",
		},
	}

	programmed := map[string][]string{}
	gwIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	c := &Controller{
		config:              &Configuration{},
		vpcNatGatewayLister: listers.NewVpcNatGatewayLister(gwIndexer),
		execRulesInPod: func(_ *v1.Pod, operation string, rules []string) error {
			programmed[operation] = append(programmed[operation], rules...)
			return nil
		},
	}
	gw := &kubeovnv1.VpcNatGateway{
		ObjectMeta: metav1.ObjectMeta{Name: "gw"},
		Spec:       kubeovnv1.VpcNatGatewaySpec{DataplaneMode: util.NatDataplaneModeStatelessNft},
	}
	require.NoError(t, gwIndexer.Add(gw))

	pods := []*v1.Pod{{ObjectMeta: metav1.ObjectMeta{Name: "gw-0", Namespace: metav1.NamespaceSystem}}}
	identities := sortedNftableLbIdentities(map[string]*kubeovnv1.IptablesDnatRule{record.Name: record}, "1.1.1.1")

	// 1. Single backend on EIP on stateless gateway produces stateless-dnat-add
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "1.1.1.1", pods, identities, nil, nil))
	require.Equal(t, []string{"1.1.1.1,80,tcp,10.0.0.5,8080"}, programmed[natGwStatelessDnatAdd])

	// 2. Removing an identity uses stateless-dnat-del and stateful nft-dnat-map-del
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "1.1.1.1", pods, nil, []*kubeovnv1.IptablesDnatRule{record}, nil))
	require.Contains(t, programmed[natGwStatelessDnatDel], "1.1.1.1,80,tcp")

	// 3. Multi-backend identity programs stateless-dnat-add for external ingress and keeps stateful nft-dnat-map-add for Hairpin
	programmed = map[string][]string{}
	multiBackendIdentities := []*nftableLbIdentity{
		{
			vip:          "1.1.1.1",
			externalPort: "80",
			protocol:     "tcp",
			backends:     []string{"10.0.0.5:8080", "10.0.0.6:8080"},
		},
	}
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "1.1.1.1", pods, multiBackendIdentities, nil, nil))
	require.Equal(t, []string{"1.1.1.1,80,tcp,10.0.0.5:8080@10.0.0.6:8080"}, programmed[natGwStatelessDnatAdd])
	require.Equal(t, []string{"1.1.1.1,80,tcp,none,0,10.0.0.5:8080@10.0.0.6:8080"}, programmed[natGwNftDnatMapAdd])
	require.Empty(t, programmed[natGwStatelessDnatDel])

	// 4. ClusterIP identity stays stateful on a stateless gateway
	programmed = map[string][]string{}
	clusterIPRecord := record.DeepCopy()
	clusterIPRecord.Spec.ClusterIP = "10.96.0.10"
	clusterIdentities := sortedNftableLbIdentities(map[string]*kubeovnv1.IptablesDnatRule{clusterIPRecord.Name: clusterIPRecord}, "1.1.1.1")
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "1.1.1.1", pods, clusterIdentities, nil, nil))
	require.Equal(t, []string{"1.1.1.1,80,tcp,10.0.0.5,8080"}, programmed[natGwStatelessDnatAdd])
	require.Equal(t, []string{"10.96.0.10,80,tcp,none,0,10.0.0.5:8080"}, programmed[natGwNftDnatMapAdd])

	// 5. Zero-backend identity does not fail and triggers rule cleanup
	programmed = map[string][]string{}
	zeroBackendIdentities := []*nftableLbIdentity{
		{
			vip:          "1.1.1.1",
			externalPort: "80",
			protocol:     "tcp",
			backends:     []string{},
		},
	}
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "1.1.1.1", pods, zeroBackendIdentities, nil, nil))
	require.Empty(t, programmed[natGwStatelessDnatAdd])
	require.Contains(t, programmed[natGwStatelessDnatDel], "1.1.1.1,80,tcp")
	require.Contains(t, programmed[natGwNftDnatMapDel], "1.1.1.1,80,tcp")

	// 6. Multi-backend identity with ClientIP affinity stays stateful and cleans up stateless rule
	programmed = map[string][]string{}
	affinityIdentities := []*nftableLbIdentity{
		{
			vip:             "1.1.1.1",
			externalPort:    "80",
			protocol:        "tcp",
			backends:        []string{"10.0.0.5:8080", "10.0.0.6:8080"},
			affinity:        string(v1.ServiceAffinityClientIP),
			affinityTimeout: 10800,
		},
	}
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "1.1.1.1", pods, affinityIdentities, nil, nil))
	require.Empty(t, programmed[natGwStatelessDnatAdd])
	require.Contains(t, programmed[natGwStatelessDnatDel], "1.1.1.1,80,tcp")
	require.Equal(t, []string{"1.1.1.1,80,tcp,clientip,10800,10.0.0.5:8080@10.0.0.6:8080"}, programmed[natGwNftDnatMapAdd])

	// 7. Dropping EIP while retaining ClusterIP still cleans up the retired EIP stateless rule even when current eipIP is empty
	programmed = map[string][]string{}
	clusterOnlyIdentities := []*nftableLbIdentity{
		{
			vip:          "10.96.0.10",
			externalPort: "80",
			protocol:     "tcp",
			backends:     []string{"10.0.0.5:8080", "10.0.0.6:8080"},
		},
	}
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "", pods, clusterOnlyIdentities, []*kubeovnv1.IptablesDnatRule{clusterIPRecord}, nil))
	require.Contains(t, programmed[natGwStatelessDnatDel], "1.1.1.1,80,tcp")
	require.Contains(t, programmed[natGwNftDnatMapDel], "1.1.1.1,80,tcp")

	// 8. More than 127 backends falls back to stateful LB and cleans up stateless rule
	programmed = map[string][]string{}
	largeBackends := make([]string, 0, maxStatelessMaglevBackends+1)
	for i := range maxStatelessMaglevBackends + 1 {
		largeBackends = append(largeBackends, fmt.Sprintf("10.0.%d.%d:8080", i/250, (i%250)+1))
	}
	largeIdentities := []*nftableLbIdentity{
		{
			vip:          "1.1.1.1",
			externalPort: "80",
			protocol:     "tcp",
			backends:     largeBackends,
		},
	}
	require.NoError(t, c.programNftableLbServiceIdentities("gw", "1.1.1.1", pods, largeIdentities, nil, nil))
	require.Empty(t, programmed[natGwStatelessDnatAdd])
	require.Contains(t, programmed[natGwStatelessDnatDel], "1.1.1.1,80,tcp")
	require.NotEmpty(t, programmed[natGwNftDnatMapAdd])
}

// TestEnqueueIptablesEipWakesGwNftableLbServices pins the EIP-event wake-up: the Service reconcile
// is the only writer of the share records that reference an EIP, so an EIP appearing or going away
// must enqueue the Services bound to it (otherwise a referencing record is never released and the
// EIP's finalizer waits forever).
func TestEnqueueIptablesEipWakesGwNftableLbServices(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web",
			Namespace: metav1.NamespaceDefault,
			Annotations: map[string]string{
				util.VpcNatGatewayAnnotation: "gw",
				util.EipAnnotation:           "eip0",
			},
		},
		Spec: v1.ServiceSpec{
			Type:  v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{{Protocol: v1.ProtocolTCP, Port: 80}},
		},
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{Services: []*v1.Service{svc}})
	require.NoError(t, err)
	ctrl := fc.fakeController
	ctrl.config.EnableGwNftableLbSvc = true
	ctrl.addOrUpdateGwNftableLbSvcQueue = newTypedRateLimitingQueue[string]("AddOrUpdateGwNftableLbSvc", nil)
	ctrl.addIptablesEipQueue = newTypedRateLimitingQueue[string]("AddIptablesEip", nil)
	ctrl.delIptablesEipQueue = newTypedRateLimitingQueue[*kubeovnv1.IptablesEIP]("DelIptablesEip", nil)
	t.Cleanup(func() {
		ctrl.addOrUpdateGwNftableLbSvcQueue.ShutDown()
		ctrl.addIptablesEipQueue.ShutDown()
		ctrl.delIptablesEipQueue.ShutDown()
	})

	eip := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip0"}}
	ctrl.enqueueAddIptablesEip(eip)
	require.Equal(t, 1, ctrl.addOrUpdateGwNftableLbSvcQueue.Len(),
		"creating the EIP must wake the Service that references it")
	item, _ := ctrl.addOrUpdateGwNftableLbSvcQueue.Get()
	require.Equal(t, "default/web", item)
	ctrl.addOrUpdateGwNftableLbSvcQueue.Done(item)

	// The workqueue de-duplicates a key that is already queued, so drain it first: the delete event
	// has to enqueue the Service again on its own.
	ctrl.enqueueDelIptablesEip(eip)
	require.Equal(t, 1, ctrl.addOrUpdateGwNftableLbSvcQueue.Len(),
		"deleting the EIP must wake the Service so it releases its records")
}
