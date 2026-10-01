package controller

import (
	"context"
	"errors"
	"maps"
	"slices"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	corelisters "k8s.io/client-go/listers/core/v1"
	discoverylisters "k8s.io/client-go/listers/discovery/v1"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	kubeovnfake "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned/fake"
	kubeovnlister "github.com/kubeovn/kube-ovn/pkg/client/listers/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/ovsdb/ovnnb"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func Test_nftableLbSvcQualifies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		svcType         v1.ServiceType
		annotations     map[string]string
		enableEIP       bool
		enableClusterIP bool
		wantEIP         bool
		wantClusterIP   bool
	}{
		{
			name: "loadbalancer with both identities enabled", svcType: v1.ServiceTypeLoadBalancer,
			annotations: map[string]string{util.VpcNatGatewayAnnotation: "gw0", util.EipAnnotation: "eip0"},
			enableEIP:   true, enableClusterIP: true, wantEIP: true, wantClusterIP: true,
		},
		{
			name: "loadbalancer with only eip identity enabled", svcType: v1.ServiceTypeLoadBalancer,
			annotations: map[string]string{util.VpcNatGatewayAnnotation: "gw0", util.EipAnnotation: "eip0"},
			enableEIP:   true, wantEIP: true,
		},
		{
			name: "loadbalancer with only cluster ip identity enabled", svcType: v1.ServiceTypeLoadBalancer,
			annotations:     map[string]string{util.VpcNatGatewayAnnotation: "gw0", util.EipAnnotation: "eip0"},
			enableClusterIP: true, wantClusterIP: true,
		},
		{
			name: "loadbalancer with both identities disabled", svcType: v1.ServiceTypeLoadBalancer,
			annotations: map[string]string{util.VpcNatGatewayAnnotation: "gw0", util.EipAnnotation: "eip0"},
		},
		{
			name: "loadbalancer without eip can still serve cluster ip", svcType: v1.ServiceTypeLoadBalancer,
			annotations: map[string]string{util.VpcNatGatewayAnnotation: "gw0"},
			enableEIP:   true, enableClusterIP: true, wantClusterIP: true,
		},
		{
			name: "cluster ip service enabled", svcType: v1.ServiceTypeClusterIP,
			annotations:     map[string]string{util.VpcNatGatewayAnnotation: "gw0"},
			enableClusterIP: true, wantClusterIP: true,
		},
		{
			name: "service without gateway", svcType: v1.ServiceTypeLoadBalancer,
			annotations: map[string]string{util.EipAnnotation: "eip0"},
			enableEIP:   true, enableClusterIP: true,
		},
		{
			name: "nodeport is not handled", svcType: v1.ServiceTypeNodePort,
			annotations: map[string]string{util.VpcNatGatewayAnnotation: "gw0", util.EipAnnotation: "eip0"},
			enableEIP:   true, enableClusterIP: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Controller{config: &Configuration{
				EnableGwNftableLbSvc: tt.enableEIP, EnableGwNftableSvcClusterIP: tt.enableClusterIP,
			}}
			svc := &v1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: tt.annotations}, Spec: v1.ServiceSpec{Type: tt.svcType}}
			gotEIP, gotClusterIP := c.nftableLbSvcQualifies(svc)
			require.Equal(t, tt.wantEIP, gotEIP)
			require.Equal(t, tt.wantClusterIP, gotClusterIP)
		})
	}
}

// TestNftableLbServiceIndexers pins what the two Service indexes hold: the EIP index holds only
// LoadBalancer Services that name an EIP (the conflict resolver can only contest an EIP identity),
// while the gateway index holds every candidate, including a LoadBalancer Service that has no EIP
// yet because the ClusterIP identity can still serve it.
func TestNftableLbServiceIndexers(t *testing.T) {
	t.Parallel()

	svc := func(name string, svcType v1.ServiceType, gateway, eip string) *v1.Service {
		annotations := map[string]string{}
		if gateway != "" {
			annotations[util.VpcNatGatewayAnnotation] = gateway
		}
		if eip != "" {
			annotations[util.EipAnnotation] = eip
		}
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: name, Annotations: annotations},
			Spec:       v1.ServiceSpec{Type: svcType},
		}
	}

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		IndexGwNftableLbServiceByEip:     indexGwNftableLbServiceByEip,
		IndexGwNftableLbServiceByGateway: indexGwNftableLbServiceByGateway,
	})
	for _, s := range []*v1.Service{
		svc("a", v1.ServiceTypeLoadBalancer, "gw0", "eip0"),
		svc("b", v1.ServiceTypeLoadBalancer, "gw0", ""), // no EIP annotation: ClusterIP identity only
		svc("c", v1.ServiceTypeLoadBalancer, "gw1", "eip1"),
		svc("d", v1.ServiceTypeClusterIP, "gw0", ""),    // a ClusterIP Service has no EIP
		svc("e", v1.ServiceTypeClusterIP, "", ""),       // no gateway: not a candidate
		svc("f", v1.ServiceTypeNodePort, "gw0", "eip0"), // NodePort is not handled
	} {
		require.NoError(t, indexer.Add(s))
	}

	names := func(objs []any) []string {
		out := make([]string, 0, len(objs))
		for _, obj := range objs {
			out = append(out, obj.(*v1.Service).Name)
		}
		slices.Sort(out)
		return out
	}

	byEIP, err := indexer.ByIndex(IndexGwNftableLbServiceByEip, "eip0")
	require.NoError(t, err)
	require.Equal(t, []string{"a"}, names(byEIP))

	byGateway, err := indexer.ByIndex(IndexGwNftableLbServiceByGateway, "gw0")
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b", "d"}, names(byGateway),
		"only gateway candidates are indexed, and a LoadBalancer Service without an EIP still is one")
}

func TestEnqueueGwNftableLbServiceWithoutOvnLb(t *testing.T) {
	t.Parallel()

	// Gateway mode runs with the OVN mode disabled, so Service events must feed its own queue.
	queue := newTypedRateLimitingQueue[string]("gw-nftable-lb-no-ovn-lb", nil)
	t.Cleanup(queue.ShutDown)
	c := &Controller{
		config:                         &Configuration{EnableOvnLB: false, EnableGwNftableLbSvc: true},
		addOrUpdateGwNftableLbSvcQueue: queue,
	}

	c.enqueueGwNftableLbService("default/web")
	require.Equal(t, 1, c.addOrUpdateGwNftableLbSvcQueue.Len())
	item, _ := c.addOrUpdateGwNftableLbSvcQueue.Get()
	require.Equal(t, "default/web", item)
	c.addOrUpdateGwNftableLbSvcQueue.Done(item)
}

func TestPodAndEipEventsDoNotEnqueueGwNftableLbService(t *testing.T) {
	t.Parallel()

	serviceQueue := newTypedRateLimitingQueue[string]("GwNftableLbService", nil)
	addEipQueue := newTypedRateLimitingQueue[string]("AddIptablesEip", nil)
	updateEipQueue := newTypedRateLimitingQueue[string]("UpdateIptablesEip", nil)
	delEipQueue := newTypedRateLimitingQueue[*kubeovnv1.IptablesEIP]("DelIptablesEip", nil)
	defer serviceQueue.ShutDown()
	defer addEipQueue.ShutDown()
	defer updateEipQueue.ShutDown()
	defer delEipQueue.ShutDown()

	endpointIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	serviceIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	emptyIndexer := func() cache.Indexer {
		return cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	}
	c := &Controller{
		config:                         &Configuration{EnableGwNftableLbSvc: true},
		addOrUpdateGwNftableLbSvcQueue: serviceQueue,
		addIptablesEipQueue:            addEipQueue,
		updateIptablesEipQueue:         updateEipQueue,
		delIptablesEipQueue:            delEipQueue,
		addIptablesFipQueue:            newTypedRateLimitingQueue[string]("AddIptablesFip", nil),
		updateIptablesFipQueue:         newTypedRateLimitingQueue[string]("UpdateIptablesFip", nil),
		addIptablesDnatRuleQueue:       newTypedRateLimitingQueue[string]("AddIptablesDnat", nil),
		updateIptablesDnatRuleQueue:    newTypedRateLimitingQueue[string]("UpdateIptablesDnat", nil),
		addIptablesSnatRuleQueue:       newTypedRateLimitingQueue[string]("AddIptablesSnat", nil),
		updateIptablesSnatRuleQueue:    newTypedRateLimitingQueue[string]("UpdateIptablesSnat", nil),
		iptablesFipsLister:             kubeovnlister.NewIptablesFIPRuleLister(emptyIndexer()),
		iptablesDnatRulesLister:        kubeovnlister.NewIptablesDnatRuleLister(emptyIndexer()),
		iptablesSnatRulesLister:        kubeovnlister.NewIptablesSnatRuleLister(emptyIndexer()),
		endpointSlicesLister:           discoverylisters.NewEndpointSliceLister(endpointIndexer),
		servicesLister:                 corelisters.NewServiceLister(serviceIndexer),
	}

	oldEip := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip0"}}
	newEip := oldEip.DeepCopy()
	newEip.ResourceVersion = "2"
	c.enqueueAddIptablesEip(oldEip)
	c.enqueueUpdateIptablesEip(oldEip, newEip)
	c.enqueueDelIptablesEip(oldEip)

	oldPod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod", Namespace: "default"}, Spec: v1.PodSpec{HostNetwork: true}}
	newPod := oldPod.DeepCopy()
	newPod.ResourceVersion = "2"
	c.enqueueAddPod(oldPod)
	c.enqueueUpdatePod(oldPod, newPod)
	c.enqueueDeletePod(oldPod)

	require.Zero(t, serviceQueue.Len())
}

func Test_nftableLbDnatRuleName(t *testing.T) {
	t.Parallel()

	// deterministic: same inputs produce the same name
	a := nftableLbDnatRuleName("ns", "svc", "tcp", "80", "10.0.0.1", "8080")
	b := nftableLbDnatRuleName("ns", "svc", "tcp", "80", "10.0.0.1", "8080")
	require.Equal(t, a, b)

	// different backend produces a different name
	c := nftableLbDnatRuleName("ns", "svc", "tcp", "80", "10.0.0.2", "8080")
	require.NotEqual(t, a, c)

	// name is DNS-1123 compliant and within the 63 char limit even for long svc names
	longName := "this-is-a-very-long-service-name-that-exceeds-the-kubernetes-limit-for-sure"
	got := nftableLbDnatRuleName("ns", longName, "tcp", "80", "10.0.0.1", "8080")
	require.LessOrEqual(t, len(got), 63)
	require.Empty(t, validation.IsDNS1123Subdomain(got))
}

func Test_nftableLbBackendResolver_missingTargetPod(t *testing.T) {
	t.Parallel()

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	c := &Controller{podsLister: corelisters.NewPodLister(indexer)}
	resolver := c.nftableLbBackendResolver(&v1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: "ns"}}, "gw-vpc")

	// EndpointSlice entry still references a Pod that no longer exists: the backend must be
	// skipped instead of falling back to the stale recorded IP (manual/selectorless slices
	// are not automatically pruned by the k8s endpointslice controller).
	ep := discoveryv1.Endpoint{
		Addresses: []string{"10.0.0.1"},
		TargetRef: &v1.ObjectReference{Kind: "Pod", Namespace: "ns", Name: "gone"},
	}
	ip, ok := resolver(ep)
	require.False(t, ok)
	require.Empty(t, ip)

	// Endpoints without a Pod target remain best-effort (e.g. externally maintained slices).
	ip, ok = resolver(discoveryv1.Endpoint{Addresses: []string{"10.0.0.2"}})
	require.True(t, ok)
	require.Equal(t, "10.0.0.2", ip)
}

// testNftableLbBackendIP resolves an endpoint to its first IPv4 address, matching the
// single-NIC (gateway-VPC-agnostic) behavior these table tests exercise.
func testNftableLbBackendIP(ep discoveryv1.Endpoint) (string, bool) {
	ip := firstIPv4(ep.Addresses)
	return ip, ip != ""
}

func Test_selectNftableLbBackendIPv4(t *testing.T) {
	t.Parallel()

	// single-NIC backend already in the gateway VPC: matches on its only NIC
	ip, ok := selectNftableLbBackendIPv4([]nftableLbNicCandidate{{ipv4: "10.0.0.1", vpc: "vpc1"}}, "vpc1")
	require.True(t, ok)
	require.Equal(t, "10.0.0.1", ip)

	// dual-NIC backend: primary in default VPC, secondary in the gateway VPC -> pick secondary
	ip, ok = selectNftableLbBackendIPv4([]nftableLbNicCandidate{
		{ipv4: "10.16.0.5", vpc: "ovn-cluster"},
		{ipv4: "192.168.0.5", vpc: "vpc1"},
	}, "vpc1")
	require.True(t, ok)
	require.Equal(t, "192.168.0.5", ip)

	// several NICs in the gateway VPC: deterministic lowest IPv4 (dual-NIC is the 99% case)
	ip, ok = selectNftableLbBackendIPv4([]nftableLbNicCandidate{
		{ipv4: "192.168.0.9", vpc: "vpc1"},
		{ipv4: "192.168.0.3", vpc: "vpc1"},
	}, "vpc1")
	require.True(t, ok)
	require.Equal(t, "192.168.0.3", ip)

	// has kube-ovn NICs but none in the gateway VPC -> no match (caller skips the backend)
	_, ok = selectNftableLbBackendIPv4([]nftableLbNicCandidate{{ipv4: "10.16.0.5", vpc: "ovn-cluster"}}, "vpc1")
	require.False(t, ok)

	// a NIC in the gateway VPC without an IPv4 is not selectable
	_, ok = selectNftableLbBackendIPv4([]nftableLbNicCandidate{{ipv4: "", vpc: "vpc1"}}, "vpc1")
	require.False(t, ok)

	// no candidates -> no match
	_, ok = selectNftableLbBackendIPv4(nil, "vpc1")
	require.False(t, ok)
}

func Test_firstIPv4(t *testing.T) {
	t.Parallel()

	require.Equal(t, "10.0.0.1", firstIPv4([]string{"fd00::1", "10.0.0.1"}))
	require.Equal(t, "10.0.0.1", firstIPv4([]string{"10.0.0.1"}))
	require.Empty(t, firstIPv4([]string{"fd00::1"}))
	require.Empty(t, firstIPv4(nil))
}

func Test_buildDesiredNftableLbDnatRules(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}, Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{
				{Name: "http", Port: 80, Protocol: v1.ProtocolTCP},
				{Name: "sctp", Port: 90, Protocol: v1.ProtocolSCTP},
			},
		},
	}

	endpointSlices := []*discoveryv1.EndpointSlice{
		{
			Ports: []discoveryv1.EndpointPort{
				{Name: new("http"), Port: new(int32(8080))},
				{Name: new("sctp"), Port: new(int32(9090))},
			},
			Endpoints: []discoveryv1.Endpoint{
				{Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
				{Addresses: []string{"10.0.0.2"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
				// not ready -> skipped
				{Addresses: []string{"10.0.0.3"}, Conditions: discoveryv1.EndpointConditions{Ready: new(false)}},
				// IPv6 -> skipped (share DNAT is IPv4 only)
				{Addresses: []string{"fd00::1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
			},
		},
	}

	desired := buildDesiredNftableLbDnatRules(svc, "eip0", "gw0", endpointSlices, testNftableLbBackendIP)

	// 2 ready IPv4 backends on the tcp port, and the sctp port is skipped: share DNAT only supports
	// tcp/udp for now (see the TODO on util.ValidateProtocol)
	tcpRules := make(map[string]*kubeovnv1.IptablesDnatRule)
	for _, rule := range desired {
		require.Equal(t, "tcp", rule.Spec.Protocol)
		tcpRules[rule.Spec.InternalIP] = rule
	}
	require.Len(t, desired, 2)
	require.Len(t, tcpRules, 2)
	require.Equal(t, "80", tcpRules["10.0.0.1"].Spec.ExternalPort)
	require.Equal(t, "8080", tcpRules["10.0.0.1"].Spec.InternalPort)

	backends := make(map[string]*kubeovnv1.IptablesDnatRule)
	for _, rule := range tcpRules {
		require.Equal(t, "eip0", rule.Spec.EIP)
		require.Equal(t, "80", rule.Spec.ExternalPort)
		require.Equal(t, "8080", rule.Spec.InternalPort)
		require.Equal(t, "tcp", rule.Spec.Protocol)
		require.Equal(t, kubeovnv1.DnatRuleTypeShare, rule.Spec.Type)
		require.Equal(t, "default", rule.Labels[util.NftableLbSvcNsLabel])
		require.Equal(t, "web", rule.Labels[util.NftableLbSvcNameLabel])
		backends[rule.Spec.InternalIP] = rule
	}
	require.Contains(t, backends, "10.0.0.1")
	require.Contains(t, backends, "10.0.0.2")
	require.NotContains(t, backends, "10.0.0.3")
	require.NotContains(t, backends, "fd00::1")
}

func Test_buildDesiredNftableLbDnatRules_unnamedPort(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}, Spec: v1.ServiceSpec{
			Type:  v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{{Port: 80, Protocol: v1.ProtocolTCP}},
		},
	}
	endpointSlices := []*discoveryv1.EndpointSlice{
		{
			Ports:     []discoveryv1.EndpointPort{{Port: new(int32(8080))}},
			Endpoints: []discoveryv1.Endpoint{{Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}}},
		},
	}

	require.Len(t, buildDesiredNftableLbDnatRules(svc, "eip0", "gw0", endpointSlices, testNftableLbBackendIP), 1)
}

func Test_buildDesiredNftableLbDnatRules_noMatchingPort(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}, Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{
				{Name: "http", Port: 80, Protocol: v1.ProtocolTCP},
			},
		},
	}
	endpointSlices := []*discoveryv1.EndpointSlice{
		{
			Ports: []discoveryv1.EndpointPort{
				{Name: new("other"), Port: new(int32(8080))},
			},
			Endpoints: []discoveryv1.Endpoint{
				{Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
			},
		},
	}

	desired := buildDesiredNftableLbDnatRules(svc, "eip0", "gw0", endpointSlices, testNftableLbBackendIP)
	require.Empty(t, desired)
}

func Test_buildDesiredNftableLbDnatRules_sessionAffinity(t *testing.T) {
	t.Parallel()

	endpointSlices := []*discoveryv1.EndpointSlice{
		{
			Ports: []discoveryv1.EndpointPort{
				{Name: new("http"), Port: new(int32(8080))},
			},
			Endpoints: []discoveryv1.Endpoint{
				{Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
			},
		},
	}
	ports := []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}}

	// ClientIP affinity with explicit timeout is propagated to the generated rules
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}, Spec: v1.ServiceSpec{
			Type:            v1.ServiceTypeLoadBalancer,
			Ports:           ports,
			SessionAffinity: v1.ServiceAffinityClientIP,
			SessionAffinityConfig: &v1.SessionAffinityConfig{
				ClientIP: &v1.ClientIPConfig{TimeoutSeconds: new(int32(600))},
			},
		},
	}
	desired := buildDesiredNftableLbDnatRules(svc, "eip0", "gw0", endpointSlices, testNftableLbBackendIP)
	require.Len(t, desired, 1)
	for _, rule := range desired {
		require.Equal(t, kubeovnv1.DnatSessionAffinityClientIP, rule.Spec.SessionAffinity)
		require.Equal(t, int32(600), rule.Spec.SessionAffinityTimeoutSeconds)
	}

	// No affinity: fields stay empty (stateless numgen-random balancing)
	svcNone := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}, Spec: v1.ServiceSpec{
			Type:            v1.ServiceTypeLoadBalancer,
			Ports:           ports,
			SessionAffinity: v1.ServiceAffinityNone,
		},
	}
	desiredNone := buildDesiredNftableLbDnatRules(svcNone, "eip0", "gw0", endpointSlices, testNftableLbBackendIP)
	require.Len(t, desiredNone, 1)
	for _, rule := range desiredNone {
		require.Equal(t, kubeovnv1.DnatSessionAffinityNone, rule.Spec.SessionAffinity)
		require.Zero(t, rule.Spec.SessionAffinityTimeoutSeconds)
	}
}

func Test_nftableLbSvcClusterIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		svc      v1.Service
		expected string
	}{
		{
			name:     "single stack",
			svc:      v1.Service{Spec: v1.ServiceSpec{ClusterIP: "10.96.1.5", ClusterIPs: []string{"10.96.1.5"}}},
			expected: "10.96.1.5",
		},
		{
			name:     "dual stack prefers IPv4",
			svc:      v1.Service{Spec: v1.ServiceSpec{ClusterIP: "10.96.1.5", ClusterIPs: []string{"10.96.1.5", "fd00:10:16::1"}}},
			expected: "10.96.1.5",
		},
		{
			name:     "legacy ClusterIP field only",
			svc:      v1.Service{Spec: v1.ServiceSpec{ClusterIP: "10.96.1.5"}},
			expected: "10.96.1.5",
		},
		{
			name: "headless",
			svc:  v1.Service{Spec: v1.ServiceSpec{ClusterIP: v1.ClusterIPNone}},
		},
		{
			name: "not allocated yet",
			svc:  v1.Service{Spec: v1.ServiceSpec{}},
		},
		{
			name: "IPv6 only (share DNAT is IPv4 only)",
			svc:  v1.Service{Spec: v1.ServiceSpec{ClusterIPs: []string{"fd00:10:16::1"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, nftableLbSvcClusterIP(&tt.svc))
		})
	}
}

func Test_buildDesiredNftableLbDnatRules_clusterIP(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}, Spec: v1.ServiceSpec{
			Type:       v1.ServiceTypeLoadBalancer,
			ClusterIP:  "10.96.1.5",
			ClusterIPs: []string{"10.96.1.5"},
			Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
		},
	}
	endpointSlices := []*discoveryv1.EndpointSlice{
		{
			Ports: []discoveryv1.EndpointPort{{Name: new("http"), Port: new(int32(8080))}},
			Endpoints: []discoveryv1.Endpoint{
				{Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
				{Addresses: []string{"10.0.0.2"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
			},
		},
	}

	// the internal VIP is part of every generated rule: it is the durable record of the ClusterIP
	// identity the gateway programs next to the EIP, and it is what a ClusterIP Service has
	desired := buildDesiredNftableLbDnatRules(svc, "eip0", "gw0", endpointSlices, testNftableLbBackendIP)
	require.Len(t, desired, 2)
	for _, rule := range desired {
		require.Equal(t, "10.96.1.5", rule.Spec.ClusterIP)
	}

	// a headless Service has no internal VIP: the EIP identity is programmed alone
	svcHeadless := svc.DeepCopy()
	svcHeadless.Spec.ClusterIP = v1.ClusterIPNone
	svcHeadless.Spec.ClusterIPs = nil
	desired = buildDesiredNftableLbDnatRules(svcHeadless, "eip0", "gw0", endpointSlices, testNftableLbBackendIP)
	require.Len(t, desired, 2)
	for _, rule := range desired {
		require.Empty(t, rule.Spec.ClusterIP)
	}
}

func Test_nftableLbDnatIdentity(t *testing.T) {
	t.Parallel()

	// protocol is normalized to lower case so TCP and tcp map to the same identity
	require.Equal(t, "eip0/80/tcp", nftableLbDnatIdentity("eip0", "80", "TCP"))
	require.Equal(t, nftableLbDnatIdentity("eip0", "80", "tcp"), nftableLbDnatIdentity("eip0", "80", "TCP"))
	require.NotEqual(t, nftableLbDnatIdentity("eip0", "80", "tcp"), nftableLbDnatIdentity("eip0", "443", "tcp"))
}

func Test_nftableLbSvcIdentities(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{Port: 80, Protocol: v1.ProtocolTCP},
				{Port: 53, Protocol: v1.ProtocolUDP},
				{Port: 90, Protocol: v1.ProtocolSCTP},
			},
		},
	}
	ids := nftableLbSvcIdentities(svc, "eip0")
	// identities must match what buildDesiredNftableLbDnatRules would program, including skipping
	// the ports it does not program: share DNAT only supports tcp/udp for now
	require.ElementsMatch(t, []string{"eip0/80/tcp", "eip0/53/udp"}, ids)
}

func controllerOwnerReference(apiVersion, kind, name string, uid types.UID) metav1.OwnerReference {
	controller := true
	return metav1.OwnerReference{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		UID:        uid,
		Controller: &controller,
	}
}

func TestResolveNftableLbConflictsWithNoReadyBackends(t *testing.T) {
	t.Parallel()

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		IndexGwNftableLbServiceByEip: indexGwNftableLbServiceByEip,
	})
	port := int32(80)
	newService := func(namespace, name string) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
				Annotations: map[string]string{
					util.EipAnnotation:           "eip0",
					util.VpcNatGatewayAnnotation: "gw0",
				},
			}, Spec: v1.ServiceSpec{
				Type:  v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{{Port: port, Protocol: v1.ProtocolTCP}},
			},
		}
	}
	winner := newService("ns", "a-winner")
	loser := newService("ns", "z-loser")
	require.NoError(t, indexer.Add(winner))
	require.NoError(t, indexer.Add(loser))

	ruleIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	queue := newTypedRateLimitingQueue[string]("gw-nftable-lb-conflict-test", nil)
	t.Cleanup(queue.ShutDown)
	controller := &Controller{
		svcIndexer:                     indexer,
		iptablesDnatRulesLister:        kubeovnlister.NewIptablesDnatRuleLister(ruleIndexer),
		recorder:                       record.NewFakeRecorder(2),
		addOrUpdateGwNftableLbSvcQueue: queue,
	}

	yielded, err := controller.resolveNftableLbConflicts(loser, "ns/z-loser", "203.0.113.10", map[string]*kubeovnv1.IptablesDnatRule{})
	require.NoError(t, err)
	require.Equal(t, map[string]struct{}{"203.0.113.10/80/tcp": {}}, yielded,
		"a loser must be detected from Service ports even without ready backends")

	desired := map[string]*kubeovnv1.IptablesDnatRule{
		"backend": {
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{util.EipV4IpLabel: "203.0.113.10", util.EipUIDLabel: "eip-uid"},
			},
			Spec: kubeovnv1.IptablesDnatRuleSpec{
				EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
			},
		},
	}
	yielded, err = controller.resolveNftableLbConflicts(loser, "ns/z-loser", "203.0.113.10", desired)
	require.NoError(t, err)
	require.Len(t, yielded, 1)
	require.Empty(t, desired["backend"].Spec.EIP)
	require.Equal(t, "10.96.1.5", desired["backend"].Spec.ClusterIP)
	require.NotContains(t, desired["backend"].Labels, util.EipV4IpLabel)
	require.NotContains(t, desired["backend"].Labels, util.EipUIDLabel)
}

func TestResolveNftableLbConflictsIgnoresTerminatingServices(t *testing.T) {
	t.Parallel()

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		IndexGwNftableLbServiceByEip: indexGwNftableLbServiceByEip,
	})
	newService := func(name string, deleting bool) *v1.Service {
		svc := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns",
				Name:      name,
				Annotations: map[string]string{
					util.EipAnnotation:           "eip0",
					util.VpcNatGatewayAnnotation: "gw0",
				},
			}, Spec: v1.ServiceSpec{
				Type:  v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{{Port: 80, Protocol: v1.ProtocolTCP}},
			},
		}
		if deleting {
			now := metav1.Now()
			svc.DeletionTimestamp = &now
		}
		return svc
	}
	// The terminating Service sorts first, so it would own the identity and starve the live
	// Service until the informer cache dropped it.
	terminating := newService("a-terminating", true)
	loser := newService("z-loser", false)
	require.NoError(t, indexer.Add(terminating))
	require.NoError(t, indexer.Add(loser))

	ruleIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	queue := newTypedRateLimitingQueue[string]("gw-nftable-lb-terminating-test", nil)
	t.Cleanup(queue.ShutDown)
	controller := &Controller{
		svcIndexer:                     indexer,
		iptablesDnatRulesLister:        kubeovnlister.NewIptablesDnatRuleLister(ruleIndexer),
		recorder:                       record.NewFakeRecorder(1),
		addOrUpdateGwNftableLbSvcQueue: queue,
	}

	yielded, err := controller.resolveNftableLbConflicts(loser, "ns/z-loser", "203.0.113.10", map[string]*kubeovnv1.IptablesDnatRule{})
	require.NoError(t, err)
	require.Empty(t, yielded, "a terminating Service must release the identity to its successor")
}

func Test_nftableLbDnatSpecEqual(t *testing.T) {
	t.Parallel()

	base := kubeovnv1.IptablesDnatRuleSpec{
		EIP: "eip0", ExternalPort: "80", Protocol: "tcp",
		InternalIP: "10.0.0.1", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
	}
	same := base
	require.True(t, nftableLbDnatSpecEqual(&base, &same))

	// session-affinity drift is detected (drives rule recreate)
	affinity := base
	affinity.SessionAffinity = kubeovnv1.DnatSessionAffinityClientIP
	affinity.SessionAffinityTimeoutSeconds = 600
	require.False(t, nftableLbDnatSpecEqual(&base, &affinity))

	// EIP drift (annotation changed) is detected
	eip := base
	eip.EIP = "eip1"
	require.False(t, nftableLbDnatSpecEqual(&base, &eip))
}

func Test_chooseNftableLbOwner(t *testing.T) {
	t.Parallel()

	// The lexicographically smallest Service key wins.
	require.Equal(t, "ns/a", chooseNftableLbOwner(map[string]struct{}{"ns/a": {}, "ns/b": {}}))
	require.Equal(t, "ns/a", chooseNftableLbOwner(map[string]struct{}{"ns/b": {}, "ns/a": {}}))

	// single owner wins trivially
	require.Equal(t, "ns/only", chooseNftableLbOwner(map[string]struct{}{"ns/only": {}}))
}

func Test_desiredNatGwVipState(t *testing.T) {
	t.Parallel()

	const gwName = "gw0"
	ruleIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	controller := &Controller{
		iptablesDnatRulesLister: kubeovnlister.NewIptablesDnatRuleLister(ruleIndexer),
		iptablesEipsLister:      kubeovnlister.NewIptablesEIPLister(eipIndexer),
	}

	owned := func(name, eip, clusterIP string) *kubeovnv1.IptablesDnatRule {
		rule := &kubeovnv1.IptablesDnatRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					util.NftableLbSvcNsLabel:     "default",
					util.NftableLbSvcNameLabel:   "web",
					util.NftableLbSvcRecordLabel: "true",
					util.VpcNatGatewayNameLabel:  gwName,
				},
			}, Spec: kubeovnv1.IptablesDnatRuleSpec{
				EIP: eip, ExternalPort: "80", Protocol: "tcp",
				InternalIP: "10.0.0.1", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
			},
		}
		rule.Spec.ClusterIP = clusterIP
		return rule
	}

	require.NoError(t, eipIndexer.Add(&kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip0"}, Status: kubeovnv1.IptablesEIPStatus{IP: "203.0.113.10"}}))
	require.NoError(t, eipIndexer.Add(&kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "not-ready"}}))

	// both VIPs of a generated rule are routed, once per VIP even with two backends
	require.NoError(t, ruleIndexer.Add(owned("a", "eip0", "10.96.1.5")))
	require.NoError(t, ruleIndexer.Add(owned("b", "eip0", "10.96.1.5")))
	// a terminating rule must not keep a route alive
	terminating := owned("terminating", "eip0", "10.96.1.9")
	now := metav1.Now()
	terminating.DeletionTimestamp = &now
	terminating.Finalizers = []string{"keep"}
	require.NoError(t, ruleIndexer.Add(terminating))
	// an EIP that is not ready yet has no address to route. Its ClusterIP is still routed: the
	// internal VIP does not depend on the EIP.
	require.NoError(t, ruleIndexer.Add(owned("not-ready", "not-ready", "10.96.1.7")))

	vips, clusterIPs, err := controller.desiredNatGwVipState(gwName, nil)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"ip4.dst == 203.0.113.10": "203.0.113.10",
		"ip4.dst == 10.96.1.5":    "10.96.1.5",
		"ip4.dst == 10.96.1.7":    "10.96.1.7",
	}, vips)
	// only the ClusterIPs are held on lo: the EIP already lives on the external interface, and a
	// terminating rule's ClusterIP is released
	require.Equal(t, []string{"10.96.1.5", "10.96.1.7"}, clusterIPs)

	// an IPv6-only ClusterIP has no IPv4 route (share DNAT is IPv4 only)
	require.NoError(t, ruleIndexer.Add(owned("v6", "eip0", "fd00::1")))
	vips, clusterIPs, err = controller.desiredNatGwVipState(gwName, nil)
	require.NoError(t, err)
	require.NotContains(t, vips, "ip4.dst == fd00::1")
	require.Len(t, vips, 3)
	require.NotContains(t, clusterIPs, "fd00::1")

	// A Service-driven rule carries its ClusterIP in the spec (the field replaces the annotation):
	// both of its addresses belong to the gateway state, the ingress IP and the internal VIP. The
	// rule is listed on its own so the counts above stay independent of this case.
	viaSpec := ownedShareRule("viaspec", "eip0", "")
	viaSpec.Spec.ClusterIP = "10.96.1.8"
	require.NoError(t, ruleIndexer.Add(viaSpec))
	vips, clusterIPs, err = controller.desiredNatGwVipState(gwName, nil)
	require.NoError(t, err)
	require.Equal(t, "10.96.1.8", vips["ip4.dst == 10.96.1.8"], "the internal VIP of the rule is routed")
	require.Contains(t, clusterIPs, "10.96.1.8", "and held on lo")
	require.Equal(t, "203.0.113.10", vips["ip4.dst == 203.0.113.10"], "next to the ingress IP it aligns with")
}

// Test_desiredNatGwVipState_keepsCacheLaggingServices pins the informer-lag window of the
// complete-set sync: a Service that already claimed and programmed its identities is not yet in
// the DNAT record cache, and the sync deriving the gateway's whole VIP set from that cache alone
// would delete its lo address and VIP route with no event to repair them. The intent of live
// Services bound to the gateway floors the set.
func Test_desiredNatGwVipState_keepsCacheLaggingServices(t *testing.T) {
	t.Parallel()

	const gwName = "gw0"
	svcIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		IndexGwNftableLbServiceByGateway: indexGwNftableLbServiceByGateway,
	})
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	controller := &Controller{
		config:                  &Configuration{EnableGwNftableLbSvc: true, EnableGwNftableSvcClusterIP: true},
		svcIndexer:              svcIndexer,
		iptablesDnatRulesLister: kubeovnlister.NewIptablesDnatRuleLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
		iptablesEipsLister:      kubeovnlister.NewIptablesEIPLister(eipIndexer),
	}
	newLbSvc := func(namespace, name, clusterIP string) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace, Name: name,
				Annotations: map[string]string{
					util.EipAnnotation: "eip0", util.VpcNatGatewayAnnotation: gwName,
				},
			}, Spec: v1.ServiceSpec{
				Type:       v1.ServiceTypeLoadBalancer,
				ClusterIP:  clusterIP,
				ClusterIPs: []string{clusterIP},
				Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
			},
		}
	}
	// live Service whose fresh records have not reached the cache yet
	require.NoError(t, svcIndexer.Add(newLbSvc("ns", "lagging", "10.96.9.9")))
	// terminating Services are being cleaned up and contribute no floor
	dying := newLbSvc("ns", "dying", "10.96.9.10")
	dying.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	require.NoError(t, svcIndexer.Add(dying))
	// the owner's floor comes from its desired records below, never from its own intent
	require.NoError(t, svcIndexer.Add(newLbSvc("ns", "owner", "10.96.9.11")))
	require.NoError(t, eipIndexer.Add(&kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip0"}, Status: kubeovnv1.IptablesEIPStatus{IP: "203.0.113.10"},
	}))

	desired := map[string]*kubeovnv1.IptablesDnatRule{
		"owner-record": {
			ObjectMeta: metav1.ObjectMeta{Name: "owner-record"},
			Spec: kubeovnv1.IptablesDnatRuleSpec{
				ClusterIP: "10.96.9.20", ExternalPort: "80", Protocol: "tcp", Type: kubeovnv1.DnatRuleTypeShare,
			},
		},
	}
	vips, clusterIPs, err := controller.desiredNatGwVipStateForService(gwName, nil, "ns/owner", desired, "")
	require.NoError(t, err)

	require.Equal(t, "10.96.9.9", vips["ip4.dst == 10.96.9.9"], "the lagging Service's internal VIP must survive")
	require.Equal(t, "203.0.113.10", vips["ip4.dst == 203.0.113.10"], "the lagging Service's public VIP must survive")
	require.Equal(t, "10.96.9.20", vips["ip4.dst == 10.96.9.20"], "the owner's desired records contribute")
	require.NotContains(t, vips, "ip4.dst == 10.96.9.10", "a terminating Service contributes no floor")
	require.NotContains(t, vips, "ip4.dst == 10.96.9.11", "the owner's own intent is covered by its desired records, not the floor")
	require.Equal(t, []string{"10.96.9.20", "10.96.9.9"}, clusterIPs)
}

func TestDesiredNatGwVipStateConvergesAfterLabelUpdate(t *testing.T) {
	t.Parallel()

	const gwName = "gw0"
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	controller := &Controller{
		iptablesDnatRulesLister: kubeovnlister.NewIptablesDnatRuleLister(indexer),
		iptablesEipsLister:      kubeovnlister.NewIptablesEIPLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
	}
	first := clusterIPServedRule("first", gwName, "10.96.1.5")
	second := clusterIPServedRule("second", gwName, "10.96.1.6")

	// The second create reconcile can see only itself while the first rule's gateway label is still
	// absent from the informer cache, so its full-set sync temporarily narrows the desired set.
	delete(first.Labels, util.VpcNatGatewayNameLabel)
	require.NoError(t, indexer.Add(first))
	require.NoError(t, indexer.Add(second))
	_, clusterIPs, err := controller.desiredNatGwVipState(gwName, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"10.96.1.6"}, clusterIPs)

	// The label patch emits the update event that calls syncNatGwVipStateFromCache. Once visible,
	// deriving the set from cache restores both VIPs.
	first.Labels[util.VpcNatGatewayNameLabel] = gwName
	require.NoError(t, indexer.Update(first))
	vips, clusterIPs, err := controller.desiredNatGwVipState(gwName, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"10.96.1.5", "10.96.1.6"}, clusterIPs)
	require.Equal(t, map[string]string{
		"ip4.dst == 10.96.1.5": "10.96.1.5",
		"ip4.dst == 10.96.1.6": "10.96.1.6",
	}, vips)
}

// Test_natGwVipStateSerializedPerGateway pins that both halves of a gateway's shared VIP state take
// the gateway lock. Both are a read-modify-write of state derived from the gateway's rules, while
// the DNAT handlers lock by rule name and run add and update in separate workers, so two rules of
// one gateway reconcile concurrently: without the lock the call that lands last would win and drop
// an address or route the other one had just made desired.
//
// The lock is asserted by holding it: the call must not finish while another holder has the key,
// and must finish once it is released. No instance is passed, so nothing is executed in a Pod.
func Test_desiredNatGwVipState_applying(t *testing.T) {
	t.Parallel()

	const gwName = "gw0"
	ruleIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	controller := &Controller{
		iptablesDnatRulesLister: kubeovnlister.NewIptablesDnatRuleLister(ruleIndexer),
		iptablesEipsLister:      kubeovnlister.NewIptablesEIPLister(eipIndexer),
	}
	require.NoError(t, eipIndexer.Add(&kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip0"}, Status: kubeovnv1.IptablesEIPStatus{IP: "203.0.113.10"}}))

	// The rule as the handler holds it: created with the owner labels, the gateway label patched
	// through the API but not observed here, so the lister is empty.
	applying := ownedShareRule("rule-a", "eip0", "10.96.1.5")
	delete(applying.Labels, util.VpcNatGatewayNameLabel)

	vips, clusterIPs, err := controller.desiredNatGwVipState(gwName, nil)
	require.NoError(t, err)
	require.Empty(t, vips, "the cache has not observed the rule yet")
	require.Empty(t, clusterIPs)

	vips, clusterIPs, err = controller.desiredNatGwVipState(gwName, applying)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"ip4.dst == 203.0.113.10": "203.0.113.10",
		"ip4.dst == 10.96.1.5":    "10.96.1.5",
	}, vips, "the rule being programmed must contribute its VIPs without waiting for the cache")
	require.Equal(t, []string{"10.96.1.5"}, clusterIPs)

	// The teardown direction: the cache still shows the rule as live, but the caller knows it is
	// terminating and its state has to go.
	live := ownedShareRule("rule-a", "eip0", "10.96.1.5")
	require.NoError(t, ruleIndexer.Add(live))
	terminating := live.DeepCopy()
	now := metav1.Now()
	terminating.DeletionTimestamp = &now
	terminating.Finalizers = []string{"keep"}

	vips, clusterIPs, err = controller.desiredNatGwVipState(gwName, terminating)
	require.NoError(t, err)
	require.Empty(t, vips, "a rule the caller is tearing down must not keep its route alive")
	require.Empty(t, clusterIPs)

	// It replaces that one rule only: a sibling of the same Service keeps the shared VIPs.
	require.NoError(t, ruleIndexer.Add(ownedShareRule("rule-b", "eip0", "10.96.1.5")))
	vips, clusterIPs, err = controller.desiredNatGwVipState(gwName, terminating)
	require.NoError(t, err)
	require.Len(t, vips, 2)
	require.Equal(t, []string{"10.96.1.5"}, clusterIPs)
}

// Test_natGwVipRouteNextHops_pendingInit pins that an instance which still needs initialization is
// not routed to. The vpc-nat-gw container has no readiness probe, so the kubelet reports it ready
// as soon as it runs, before the controller programmed its chains and identities.
func Test_natGwVipRouteNextHops_pendingInit(t *testing.T) {
	t.Parallel()

	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw0"}, Spec: kubeovnv1.VpcNatGatewaySpec{Vpc: "vpc0", Subnet: "subnet0"}}
	pods := gatewayPodsOnNodes("gw0", map[string]string{"node0": "10.0.7.254", "node1": "10.0.7.253"})

	nextHops, err := natGwVipRouteNextHops(gw, pods)
	require.NoError(t, err)
	require.Equal(t, []string{"10.0.7.253", "10.0.7.254"}, nextHops)

	// node0 sorts first, so pods[0] is its instance: it is running and ready, but the init mark of
	// its container instance is missing.
	delete(pods[0].Annotations, util.VpcNatGatewayInitAnnotation)
	nextHops, err = natGwVipRouteNextHops(gw, pods)
	require.NoError(t, err)
	require.Equal(t, []string{"10.0.7.253"}, nextHops,
		"an instance that is not initialized yet has no share DNAT identity to serve the VIP")
}

// gatewayPods returns a ready, initialized NAT gateway instance as the VIP routes' next hop. The
// init annotation is what tells the controller the instance's chains and identities are programmed;
// without it the instance is not a next hop (see natGwVipRouteNextHops).
func gatewayPods(gwName, ip string) []*v1.Pod {
	return []*v1.Pod{{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   metav1.NamespaceSystem,
			Name:        util.GenNatGwPodName(gwName),
			Annotations: map[string]string{util.VpcNatGatewayInitAnnotation: "true"},
			Labels: map[string]string{
				"app":                       util.GenNatGwName(gwName),
				util.VpcNatGatewayLabel:     "true",
				util.VpcNatGatewayNameLabel: gwName,
			},
		}, Spec: v1.PodSpec{NodeName: "node0"},
		Status: v1.PodStatus{
			Phase:      v1.PodRunning,
			PodIPs:     []v1.PodIP{{IP: ip}},
			PodIP:      ip,
			Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}},
		},
	}}
}

// gatewayPodsOnNodes returns one ready NAT gateway instance per node/LAN IP pair, so the VIP route
// next hop set covers more than one replica.
func gatewayPodsOnNodes(gwName string, nodeIPs map[string]string) []*v1.Pod {
	pods := make([]*v1.Pod, 0, len(nodeIPs))
	// Sorted, so a test that singles out one instance always picks the same one.
	for _, node := range slices.Sorted(maps.Keys(nodeIPs)) {
		ip := nodeIPs[node]
		pods = append(pods, &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:   metav1.NamespaceSystem,
				Name:        util.GenNatGwPodName(gwName) + "-" + node,
				Annotations: map[string]string{util.VpcNatGatewayInitAnnotation: "true"},
				Labels: map[string]string{
					"app":                       util.GenNatGwName(gwName),
					util.VpcNatGatewayLabel:     "true",
					util.VpcNatGatewayNameLabel: gwName,
				},
			}, Spec: v1.PodSpec{NodeName: node},
			Status: v1.PodStatus{
				Phase:      v1.PodRunning,
				PodIPs:     []v1.PodIP{{IP: ip}},
				PodIP:      ip,
				Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}},
			},
		})
	}
	return pods
}

// ownedShareRule is a share DNAT rule generated by the nftable LB service feature for a Service
// whose ClusterIP is recorded on the rule.
func ownedShareRule(name, eip, clusterIP string) *kubeovnv1.IptablesDnatRule {
	rule := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				util.NftableLbSvcNsLabel:     "default",
				util.NftableLbSvcNameLabel:   "web",
				util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel:  "gw0",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: eip, ClusterIP: clusterIP, ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.7.2", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
	return rule
}

func Test_syncNatGwVipState(t *testing.T) {
	t.Parallel()

	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw0"}, Spec: kubeovnv1.VpcNatGatewaySpec{Vpc: "vpc0"}}
	eip := &kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip0"}, Status: kubeovnv1.IptablesEIPStatus{IP: "172.20.0.5"}}
	externalIDs := natGwVipRouteExternalIDs("gw0")

	t.Run("creates a route per VIP", func(t *testing.T) {
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways:    []*kubeovnv1.VpcNatGateway{gw},
			Pods:              gatewayPods("gw0", "10.0.7.254"),
			IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
			IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownedShareRule("rule-a", "eip0", "10.96.1.5")},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableGwNftableLbSvc = true

		fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority, externalIDs, false).Return(nil, nil)
		added := map[string][]string{}
		fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy("vpc0", util.NatGatewayVipPolicyPriority, gomock.Any(), string(kubeovnv1.PolicyRouteActionReroute), gomock.Any(), nil, externalIDs).
			DoAndReturn(func(_ string, _ int, match, _ string, nextHops, _ []string, _ map[string]string) error {
				added[match] = nextHops
				return nil
			}).Times(2)

		require.NoError(t, c.syncNatGwVipState("gw0", nil))
		require.Equal(t, map[string][]string{
			"ip4.dst == 172.20.0.5": {"10.0.7.254"},
			"ip4.dst == 10.96.1.5":  {"10.0.7.254"},
		}, added)
	})

	t.Run("removes routes that no rule references and refreshes the next hops", func(t *testing.T) {
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways:    []*kubeovnv1.VpcNatGateway{gw},
			Pods:              gatewayPods("gw0", "10.0.7.254"),
			IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
			IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownedShareRule("rule-a", "eip0", "10.96.1.5")},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableGwNftableLbSvc = true

		stale := &ovnnb.LogicalRouterPolicy{UUID: "stale-uuid", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.9"}
		// Keep the expected next hop but drift the action: comparing nexthops alone would leave
		// this policy as "allow" forever and the VIP would bypass the gateway.
		drifted := &ovnnb.LogicalRouterPolicy{
			UUID: "drifted-uuid", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.5",
			Action: "allow", Nexthops: []string{"10.0.7.254"},
		}
		fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority, externalIDs, false).
			Return([]*ovnnb.LogicalRouterPolicy{stale, drifted}, nil)
		fc.mockOvnClient.EXPECT().DeleteLogicalRouterPolicyByUUID("vpc0", "stale-uuid").Return(nil)
		updated := map[string][]string{}
		fc.mockOvnClient.EXPECT().UpdateLogicalRouterPolicy(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(policy *ovnnb.LogicalRouterPolicy, _ ...any) error {
				require.Equal(t, string(kubeovnv1.PolicyRouteActionReroute), policy.Action)
				updated[policy.Match] = policy.Nexthops
				return nil
			}).Times(1)
		// the EIP's route did not exist yet
		fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy("vpc0", util.NatGatewayVipPolicyPriority, "ip4.dst == 172.20.0.5",
			string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, externalIDs).Return(nil)

		require.NoError(t, c.syncNatGwVipState("gw0", nil))
		require.Equal(t, 29210, util.NatGatewayVipPolicyPriority)
		require.Equal(t, map[string][]string{"ip4.dst == 10.96.1.5": {"10.0.7.254"}}, updated)
	})

	// This is what the gateway-move path in handleUpdateIptablesDnatRule relies on: once the rule's
	// label points at its new gateway, reconciling the one it left drops the routes it was the last
	// user of, instead of leaving them steering its addresses to a gateway that no longer serves it.
	t.Run("a rule that moved to another gateway leaves no route behind", func(t *testing.T) {
		moved := ownedShareRule("rule-a", "eip0", "10.96.1.5")
		moved.Labels[util.VpcNatGatewayNameLabel] = "gw1"
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways:    []*kubeovnv1.VpcNatGateway{gw},
			Pods:              gatewayPods("gw0", "10.0.7.254"),
			IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
			IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{moved},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableGwNftableLbSvc = true

		left := []*ovnnb.LogicalRouterPolicy{
			{UUID: "clusterip-uuid", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.5"},
			{UUID: "eip-uuid", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 172.20.0.5"},
		}
		fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority, externalIDs, false).Return(left, nil)
		var deleted []string
		fc.mockOvnClient.EXPECT().DeleteLogicalRouterPolicyByUUID("vpc0", gomock.Any()).
			DoAndReturn(func(_, uuid string) error {
				deleted = append(deleted, uuid)
				return nil
			}).Times(2)

		require.NoError(t, c.syncNatGwVipState("gw0", nil))
		require.ElementsMatch(t, []string{"clusterip-uuid", "eip-uuid"}, deleted)
	})

	t.Run("routes a VIP to every ready instance (ECMP)", func(t *testing.T) {
		// An HA gateway serves the VIP from any replica, so the route must carry all of them.
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
			Pods: gatewayPodsOnNodes("gw0", map[string]string{
				"node0": "10.0.7.254",
				"node1": "10.0.7.253",
			}),
			IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
			IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownedShareRule("rule-a", "eip0", "10.96.1.5")},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableGwNftableLbSvc = true

		fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority, externalIDs, false).Return(nil, nil)
		added := map[string][]string{}
		fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy("vpc0", util.NatGatewayVipPolicyPriority, gomock.Any(), string(kubeovnv1.PolicyRouteActionReroute), gomock.Any(), nil, externalIDs).
			DoAndReturn(func(_ string, _ int, match, _ string, nextHops, _ []string, _ map[string]string) error {
				added[match] = nextHops
				return nil
			}).Times(2)

		require.NoError(t, c.syncNatGwVipState("gw0", nil))
		require.Equal(t, map[string][]string{
			"ip4.dst == 172.20.0.5": {"10.0.7.253", "10.0.7.254"},
			"ip4.dst == 10.96.1.5":  {"10.0.7.253", "10.0.7.254"},
		}, added)
	})

	t.Run("routes a VIP only to the ready instances", func(t *testing.T) {
		pods := gatewayPodsOnNodes("gw0", map[string]string{
			"node0": "10.0.7.254",
			"node1": "10.0.7.253",
		})
		pods[1].Status.Conditions = []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionFalse}}
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways:    []*kubeovnv1.VpcNatGateway{gw},
			Pods:              pods,
			IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
			IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownedShareRule("rule-a", "eip0", "10.96.1.5")},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableGwNftableLbSvc = true

		fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority, externalIDs, false).Return(nil, nil)
		added := map[string][]string{}
		fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy("vpc0", util.NatGatewayVipPolicyPriority, gomock.Any(), string(kubeovnv1.PolicyRouteActionReroute), gomock.Any(), nil, externalIDs).
			DoAndReturn(func(_ string, _ int, match, _ string, nextHops, _ []string, _ map[string]string) error {
				added[match] = nextHops
				return nil
			}).Times(2)

		require.NoError(t, c.syncNatGwVipState("gw0", nil))
		ready := added["ip4.dst == 10.96.1.5"]
		require.Equal(t, []string{"10.0.7.254"}, ready, "a not ready replica must not be a next hop")
	})

	t.Run("drops every route while no instance has an IPv4 address", func(t *testing.T) {
		// a gateway whose VPC-side address is IPv6 only has no data plane for the IPv4 share DNAT,
		// so an ip4.dst route must not be left pointing at its IPv6 address
		v6Only := gatewayPods("gw0", "fd00::10")
		v6Only[0].Status.PodIPs = []v1.PodIP{{IP: "fd00::10"}}
		v6Only[0].Status.PodIP = "fd00::10"
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways:    []*kubeovnv1.VpcNatGateway{gw},
			Pods:              v6Only,
			IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
			IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownedShareRule("rule-a", "eip0", "10.96.1.5")},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableGwNftableLbSvc = true

		existing := &ovnnb.LogicalRouterPolicy{UUID: "u1", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.5"}
		fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority, externalIDs, false).
			Return([]*ovnnb.LogicalRouterPolicy{existing}, nil)
		fc.mockOvnClient.EXPECT().DeleteLogicalRouterPolicyByUUID("vpc0", "u1").Return(nil)

		require.NoError(t, c.syncNatGwVipState("gw0", nil))
	})

	t.Run("drops every route while no instance is reachable", func(t *testing.T) {
		notReady := gatewayPods("gw0", "10.0.7.254")
		notReady[0].Status.Conditions = []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionFalse}}
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			VpcNatGateways:    []*kubeovnv1.VpcNatGateway{gw},
			Pods:              notReady,
			IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
			IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownedShareRule("rule-a", "eip0", "10.96.1.5")},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableGwNftableLbSvc = true

		existing := &ovnnb.LogicalRouterPolicy{UUID: "u1", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.5"}
		fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority, externalIDs, false).
			Return([]*ovnnb.LogicalRouterPolicy{existing}, nil)
		fc.mockOvnClient.EXPECT().DeleteLogicalRouterPolicyByUUID("vpc0", "u1").Return(nil)

		require.NoError(t, c.syncNatGwVipState("gw0", nil))
	})
}

// Test_syncNatGwVipStateDisabled pins the feature gate: with --enable-gw-nftable-lb-svc off the
// feature owns no VIP, so the sync must not touch OVN at all (the mock client fails any call).
// Test_syncNatGwVipStateDisabled pins the feature-off semantics: the sync is a teardown path
// only. Unclaimed routes are removed (Services being deleted must be able to release their VIP
// routes), while claimed routes are left as-is and nothing is created or rewritten.
func Test_syncNatGwVipStateDisabled(t *testing.T) {
	t.Parallel()

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{{ObjectMeta: metav1.ObjectMeta{Name: "gw0"}, Spec: kubeovnv1.VpcNatGatewaySpec{Vpc: "vpc0", Subnet: "subnet0"}}},
		Pods:              gatewayPods("gw0", "10.0.7.254"),
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownedShareRule("rule-a", "eip0", "10.96.1.5")},
	})
	require.NoError(t, err)
	fc.fakeController.config.EnableGwNftableLbSvc = false
	fc.fakeController.config.EnableGwNftableSvcClusterIP = false

	stale := &ovnnb.LogicalRouterPolicy{UUID: "stale-uuid", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.9"}
	// Claimed by rule-a, but its action drifted: with the feature off it is left alone.
	drifted := &ovnnb.LogicalRouterPolicy{
		UUID: "drifted-uuid", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.5",
		Action: "allow", Nexthops: []string{"10.0.7.254"},
	}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies("vpc0", util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs("gw0"), false).Return([]*ovnnb.LogicalRouterPolicy{stale, drifted}, nil)
	fc.mockOvnClient.EXPECT().DeleteLogicalRouterPolicyByUUID("vpc0", "stale-uuid").Return(nil)
	// The strict mock fails the test on any update or add: nothing may be programmed once the
	// feature is disabled.

	require.NoError(t, fc.fakeController.syncNatGwVipState("gw0", nil))
}

// Test_enqueueGwNftableLbService_gatesOffStillEnqueues pins the exit path's trigger: with both
// feature gates off, a Service that still carries the feature's finalizer and records must be
// enqueued anyway, because this queue is the only owner that can release either.
func Test_enqueueGwNftableLbService_gatesOffStillEnqueues(t *testing.T) {
	t.Parallel()

	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	t.Cleanup(queue.ShutDown)
	c := &Controller{
		config:                         &Configuration{EnableGwNftableLbSvc: false, EnableGwNftableSvcClusterIP: false},
		addOrUpdateGwNftableLbSvcQueue: queue,
	}
	c.enqueueGwNftableLbService("default/left-behind")
	require.Equal(t, 1, queue.Len())
}

// Test_handleAddOrUpdateGwNftableLbService_gatesOffStillCleansUp reproduces the shutdown sequence:
// the feature was on, the gates were flipped to false and the controller restarted, and the user
// now deletes a Service the feature manages. Cleanup must still run: remove the identities from
// the gateway, retire the records, withdraw the VIP routes, and release the Service finalizer,
// without programming anything new.
func Test_handleAddOrUpdateGwNftableLbService_gatesOffStillCleansUp(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	terminating := metav1.Now()
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}
	f.svc.DeletionTimestamp = &terminating
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "leftover-record",
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: f.namespace, util.NftableLbSvcNameLabel: f.svc.Name,
				util.NftableLbSvcUIDLabel: string(f.svc.UID), util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel: f.gw.Name,
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
	// both gates off: the feature is disabled
	c.config.EnableGwNftableLbSvc = false
	c.config.EnableGwNftableSvcClusterIP = false

	var execs []string
	c.execRulesInPod = func(_ *v1.Pod, operation string, _ []string) error {
		execs = append(execs, operation)
		return nil
	}
	staleRoute := &ovnnb.LogicalRouterPolicy{UUID: "leftover-uuid", Priority: util.NatGatewayVipPolicyPriority, Match: "ip4.dst == 10.96.1.5"}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return([]*ovnnb.LogicalRouterPolicy{staleRoute}, nil)
	fc.mockOvnClient.EXPECT().DeleteLogicalRouterPolicyByUUID(util.DefaultVpc, "leftover-uuid").Return(nil)
	// The strict mock fails the test on any route add or update: disabled means teardown only.

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+f.svc.Name))

	require.Contains(t, execs, natGwNftDnatMapDel, "the leftover identities are removed from the gateway")
	require.Contains(t, execs, natGwVipHairpinDel)
	require.Contains(t, execs, natGwVipAddrSync, "the lo VIPs are released")
	_, err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.True(t, k8serrors.IsNotFound(err), "the leftover accounting record is retired")
	svc, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, svc.Finalizers, "the Service finalizer is released so the deletion can finish")
}

// TestHandleAddOrUpdateVpcNatGwSyncsVipRoutes pins the gateway side convergence of the share DNAT
// VIP routes: an HA scale down removes an instance without creating a new one, so no DNAT rule is
// rewritten and the only chance to drop the gone instance from the routes is the gateway reconcile.
func TestHandleAddOrUpdateVpcNatGwSyncsVipRoutes(t *testing.T) {
	const (
		gwName     = "sync-gw"
		subnetName = "nat-subnet"
		eipName    = "sync-eip"
		vpcName    = util.DefaultVpc
		liveIP     = "10.20.0.10"
		goneIP     = "10.20.0.11"
		clusterIP  = "10.96.1.5"
	)
	namespace := metav1.NamespaceSystem
	stsName := util.GenNatGwName(gwName)
	gwLabels := util.GenNatGwLabels(gwName)
	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: gwName, UID: "gw-uid"}, Spec: kubeovnv1.VpcNatGatewaySpec{Vpc: vpcName, Subnet: subnetName}}
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: stsName, Namespace: namespace, UID: "sts-uid"}, Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: gwLabels},
			Template: v1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: gwLabels, Annotations: map[string]string{util.VpcNatGatewayAnnotation: gwName}}},
		},
	}
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: stsName + "-0", Namespace: namespace, Labels: gwLabels,
			Annotations: map[string]string{
				util.IPAddressAnnotation:         liveIP,
				util.VpcNatGatewayInitAnnotation: "true",
				util.VpcNatGatewayAnnotation:     gwName,
			},
			OwnerReferences: []metav1.OwnerReference{controllerOwnerReference(appsv1.SchemeGroupVersion.String(), util.KindStatefulSet, stsName, sts.UID)},
		}, Spec: v1.PodSpec{NodeName: "node-1"},
		Status: v1.PodStatus{
			Phase:      v1.PodRunning,
			PodIPs:     []v1.PodIP{{IP: liveIP}},
			Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}},
			ContainerStatuses: []v1.ContainerStatus{{
				Name: "vpc-nat-gw", ContainerID: "containerd://sync-instance",
				State: v1.ContainerState{Running: &v1.ContainerStateRunning{}},
			}},
		},
	}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default", Name: "sync-service",
			Annotations: map[string]string{util.VpcNatGatewayAnnotation: gwName, util.EipAnnotation: eipName},
		}, Spec: v1.ServiceSpec{Type: v1.ServiceTypeLoadBalancer},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: eipName}, Spec: kubeovnv1.IptablesEIPSpec{NatGwDp: gwName, V4ip: "172.20.0.5"},
		Status: kubeovnv1.IptablesEIPStatus{IP: "172.20.0.5"},
	}
	rule := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sync-rule",
			Labels: map[string]string{
				util.NftableLbSvcNsLabel:    "default",
				util.NftableLbSvcNameLabel:  "web",
				util.VpcNatGatewayNameLabel: gwName,
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: eipName, ClusterIP: clusterIP,
			ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.7.2", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}

	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{{ObjectMeta: metav1.ObjectMeta{Name: vpcName}}},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		Subnets: []*kubeovnv1.Subnet{
			{ObjectMeta: metav1.ObjectMeta{Name: subnetName}, Spec: kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Protocol: kubeovnv1.ProtocolIPv4, Vpc: vpcName, CIDRBlock: "10.20.0.0/24", Gateway: "10.20.0.1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "ovn-vpc-external-network"}, Spec: kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Protocol: kubeovnv1.ProtocolIPv4, Gateway: "192.168.0.1"}},
		},
		StatefulSets:      []*appsv1.StatefulSet{sts},
		Pods:              []*v1.Pod{pod},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{eip},
		Services:          []*v1.Service{svc},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{rule},
	})
	require.NoError(t, err)
	controller := fakeController.fakeController
	controller.config.EnableGwNftableLbSvc = true
	controller.config.EnableOvnLB = true
	controller.addOrUpdateGwNftableLbSvcQueue = newTypedRateLimitingQueue[string]("gw-nftable-lb-svc", nil)
	t.Cleanup(controller.addOrUpdateGwNftableLbSvcQueue.ShutDown)

	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = "unknown" })

	// Gateway reconcile only re-enqueues the bound Service. The Service is the sole writer of
	// share-DNAT VIP routes, so no OVN mock is expected here.
	require.NoError(t, controller.handleUpdateVpcDnat(gwName))
	require.Equal(t, 1, controller.addOrUpdateGwNftableLbSvcQueue.Len())
}

// TestHandleUpdateEndpointSliceSkipsOwnedVips pins the OVN load balancer side of the ownership rule:
// for a Service the gateway serves, the handler releases the VIP instead of programming it.
// nftableLbSvcOwnershipFixture is the shared setup of the ownership tests: one Service handled by
// the nftable LB service feature, its EIP, the gateway serving it, and the VPC load balancer that
// must not program its VIPs while the gateway owns them.
type nftableLbSvcOwnershipFixture struct {
	svc       *v1.Service
	slice     *discoveryv1.EndpointSlice
	eip       *kubeovnv1.IptablesEIP
	gw        *kubeovnv1.VpcNatGateway
	vpc       *kubeovnv1.Vpc
	subnet    *kubeovnv1.Subnet
	namespace string
}

func newNftableLbSvcOwnershipFixture() *nftableLbSvcOwnershipFixture {
	const (
		gwName     = "owned-gw"
		subnetName = "owned-subnet"
		eipName    = "owned-eip"
		vpcName    = util.DefaultVpc
		svcName    = "owned-web"
		clusterIP  = "10.96.1.5"
	)
	namespace := metav1.NamespaceDefault
	return &nftableLbSvcOwnershipFixture{
		namespace: namespace,
		gw:        &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: gwName, UID: "gw-uid"}, Spec: kubeovnv1.VpcNatGatewaySpec{Vpc: vpcName, Subnet: subnetName}},
		eip: &kubeovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{Name: eipName, UID: "owned-eip-uid"},
			Spec:       kubeovnv1.IptablesEIPSpec{NatGwDp: gwName, V4ip: "172.20.0.5"},
			Status:     kubeovnv1.IptablesEIPStatus{IP: "172.20.0.5", Ready: true},
		},
		vpc:    &kubeovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: vpcName}, Status: kubeovnv1.VpcStatus{TCPLoadBalancer: "vpc-tcp-load", TCPSessionLoadBalancer: "vpc-tcp-sess"}},
		subnet: &kubeovnv1.Subnet{ObjectMeta: metav1.ObjectMeta{Name: subnetName}, Spec: kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Protocol: kubeovnv1.ProtocolIPv4, Vpc: vpcName, CIDRBlock: "10.0.7.0/24", Gateway: "10.0.7.1"}},
		svc: &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      svcName,
				UID:       "owned-web-uid",
				Annotations: map[string]string{
					util.EipAnnotation:           eipName,
					util.VpcNatGatewayAnnotation: gwName,
				},
			}, Spec: v1.ServiceSpec{
				Type:       v1.ServiceTypeLoadBalancer,
				ClusterIP:  clusterIP,
				ClusterIPs: []string{clusterIP},
				Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
			},
		},
		slice: &discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace, Name: svcName + "-abc",
				Labels: map[string]string{discoveryv1.LabelServiceName: svcName},
			}, AddressType: discoveryv1.AddressTypeIPv4,
			Ports:     []discoveryv1.EndpointPort{{Name: new("http"), Port: new(int32(8080)), Protocol: new(v1.ProtocolTCP)}},
			Endpoints: []discoveryv1.Endpoint{{Addresses: []string{"10.0.7.2"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}}},
		},
	}
}

func TestNftableLbServiceClaimsFinalizerBeforeProgramming(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{f.slice},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{f.eip},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+f.svc.Name))
	updated, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, updated.Status.LoadBalancer.Ingress)
	require.Contains(t, updated.Finalizers, util.KubeOVNControllerFinalizer)
}

func clusterIPServedRule(name, gateway, clusterIP string) *kubeovnv1.IptablesDnatRule {
	return &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				util.NftableLbSvcNsLabel:     "default",
				util.NftableLbSvcNameLabel:   "web",
				util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel:  gateway,
				util.VpcDnatEPortLabel:       "80",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			ClusterIP:    clusterIP,
			ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.7.2", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		},
	}
}

// Test_desiredNatGwVipState_clusterIPServedRules pins that a rule serving a ClusterIP owns its
// address on its own: without an owner label it still contributes the VIP route and the address to
// hold on lo, otherwise the rule would program a share DNAT identity nothing routes to.
func Test_desiredNatGwVipState_clusterIPServedRules(t *testing.T) {
	t.Parallel()

	const gwName = "gw0"
	ruleIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	require.NoError(t, eipIndexer.Add(&kubeovnv1.IptablesEIP{ObjectMeta: metav1.ObjectMeta{Name: "eip0"}, Status: kubeovnv1.IptablesEIPStatus{IP: "203.0.113.10"}}))
	controller := &Controller{
		iptablesDnatRulesLister: kubeovnlister.NewIptablesDnatRuleLister(ruleIndexer),
		iptablesEipsLister:      kubeovnlister.NewIptablesEIPLister(eipIndexer),
	}
	require.NoError(t, ruleIndexer.Add(clusterIPServedRule("c1", gwName, "10.96.1.5")))
	require.NoError(t, ruleIndexer.Add(clusterIPServedRule("c2", gwName, "10.96.1.6")))

	vips, clusterIPs, err := controller.desiredNatGwVipState(gwName, nil)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"ip4.dst == 10.96.1.5": "10.96.1.5",
		"ip4.dst == 10.96.1.6": "10.96.1.6",
	}, vips)
	require.Equal(t, []string{"10.96.1.5", "10.96.1.6"}, clusterIPs)

	// A terminating rule releases its address again.
	terminating := clusterIPServedRule("c1", gwName, "10.96.1.5")
	now := metav1.Now()
	terminating.DeletionTimestamp = &now
	terminating.Finalizers = []string{"keep"}
	vips, clusterIPs, err = controller.desiredNatGwVipState(gwName, terminating)
	require.NoError(t, err)
	require.NotContains(t, vips, "ip4.dst == 10.96.1.5")
	require.Equal(t, []string{"10.96.1.6"}, clusterIPs)
}

// Test_sameDnatIdentity pins the identity comparison both flows share: the two address kinds can
// never be confused with each other, since a rule carries exactly one of them.
func Test_sameDnatIdentity(t *testing.T) {
	t.Parallel()

	eipA := &kubeovnv1.IptablesDnatRuleSpec{EIP: "eip-a"}
	eipB := &kubeovnv1.IptablesDnatRuleSpec{EIP: "eip-b"}
	// A Service-driven rule aligns the ingress IP and the ClusterIP on one identity.
	eipASvc := &kubeovnv1.IptablesDnatRuleSpec{EIP: "eip-a", ClusterIP: "10.96.1.5"}
	cipA := &kubeovnv1.IptablesDnatRuleSpec{ClusterIP: "10.96.1.5"}
	cipA2 := &kubeovnv1.IptablesDnatRuleSpec{ClusterIP: "10.96.1.5"}
	cipB := &kubeovnv1.IptablesDnatRuleSpec{ClusterIP: "10.96.1.6"}

	require.True(t, sameDnatIdentity(eipA, eipA))
	require.False(t, sameDnatIdentity(eipA, eipB))
	require.True(t, sameDnatIdentity(eipA, eipASvc),
		"the EIP:port nft map is shared, so a rule that does not carry the ClusterIP still collides")
	require.False(t, sameDnatIdentity(eipA, cipA), "an EIP identity and a ClusterIP identity are distinct")
	require.True(t, sameDnatIdentity(cipA, cipA2), "a ClusterIP identity does not depend on the gateway")
	require.False(t, sameDnatIdentity(cipA, cipB))
}

// Test_buildDesiredNftableLbDnatRules_clusterIPService pins the rule shape of the two Service
// kinds: a ClusterIP Service is served through its internal VIP alone, a LoadBalancer Service
// through the EIP and that same VIP aligned on one rule.
func Test_buildDesiredNftableLbDnatRules_clusterIPService(t *testing.T) {
	t.Parallel()

	endpointSlices := []*discoveryv1.EndpointSlice{
		{
			Ports: []discoveryv1.EndpointPort{{Name: new("http"), Port: new(int32(8080))}},
			Endpoints: []discoveryv1.Endpoint{
				{Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}},
			},
		},
	}
	newSvc := func(svcType v1.ServiceType, eip string) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default", Name: "web",
				Annotations: map[string]string{
					util.VpcNatGatewayAnnotation: "gw0",
					util.EipAnnotation:           eip,
				},
			}, Spec: v1.ServiceSpec{
				Type:       svcType,
				ClusterIP:  "10.96.1.5",
				ClusterIPs: []string{"10.96.1.5"},
				Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
			},
		}
	}

	// ClusterIP Service: no EIP at all, so its record carries the internal VIP.
	desired := buildDesiredNftableLbDnatRules(newSvc(v1.ServiceTypeClusterIP, ""), "", "gw0", endpointSlices, testNftableLbBackendIP)
	require.Len(t, desired, 1)
	for _, rule := range desired {
		require.Empty(t, rule.Spec.EIP, "a ClusterIP Service has no public address")
		require.Equal(t, "10.96.1.5", rule.Spec.ClusterIP)
		require.Equal(t, kubeovnv1.DnatRuleTypeShare, rule.Spec.Type)
	}

	// LoadBalancer Service: both addresses on the same rule, so the two VIPs of one Service port
	// stay aligned and are programmed and released together.
	desired = buildDesiredNftableLbDnatRules(newSvc(v1.ServiceTypeLoadBalancer, "eip0"), "eip0", "gw0", endpointSlices, testNftableLbBackendIP)
	require.Len(t, desired, 1)
	for _, rule := range desired {
		require.Equal(t, "eip0", rule.Spec.EIP)
		require.Equal(t, "10.96.1.5", rule.Spec.ClusterIP)
	}
}

// Test_buildDesiredNftableLbDnatRules_featureGates pins the per-identity feature gates: each
// address is dropped by its own switch, and one Service port served by both keeps a single record
// that carries both fields.
func Test_buildDesiredNftableLbDnatRules_featureGates(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default", Name: "web",
			Annotations: map[string]string{util.VpcNatGatewayAnnotation: "gw0", util.EipAnnotation: "eip0"},
		}, Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer, ClusterIP: "10.96.1.5", ClusterIPs: []string{"10.96.1.5"},
			Ports: []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
		},
	}
	endpointSlices := []*discoveryv1.EndpointSlice{{
		Ports: []discoveryv1.EndpointPort{{Name: new("http"), Port: new(int32(8080))}},
		Endpoints: []discoveryv1.Endpoint{{
			Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)},
		}},
	}}
	tests := []struct {
		name                          string
		enableEIP, enableClusterIP    bool
		wantEIP, wantClusterIP        string
		wantIdentities, wantRuleCount int
	}{
		{name: "both disabled"},
		{name: "eip only", enableEIP: true, wantEIP: "eip0", wantIdentities: 1, wantRuleCount: 1},
		{name: "cluster ip only", enableClusterIP: true, wantClusterIP: "10.96.1.5", wantIdentities: 1, wantRuleCount: 1},
		{name: "both enabled", enableEIP: true, enableClusterIP: true, wantEIP: "eip0", wantClusterIP: "10.96.1.5", wantIdentities: 2, wantRuleCount: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desired := buildDesiredNftableLbDnatRulesForIdentities(svc, "eip0", "gw0", tt.enableEIP, tt.enableClusterIP, endpointSlices, testNftableLbBackendIP)
			require.Len(t, desired, tt.wantRuleCount)
			require.Len(t, buildNftableLbIdentities(desired, "203.0.113.10"), tt.wantIdentities)
			for _, rule := range desired {
				require.Equal(t, tt.wantEIP, rule.Spec.EIP)
				require.Equal(t, tt.wantClusterIP, rule.Spec.ClusterIP)
			}
		})
	}
}

// Test_handleAddOrUpdateGwNftableLbService_clusterIPService pins the ClusterIP Service reconcile:
// the gateway comes from the annotation, no EIP is looked up, and the generated rule is the one
// that serves the internal VIP.
func Test_handleAddOrUpdateGwNftableLbService_clusterIPService(t *testing.T) {
	const (
		gwName     = "cip-gw"
		subnetName = "cip-subnet"
		vpcName    = util.DefaultVpc
		svcName    = "cip-web"
		clusterIP  = "10.96.1.9"
	)
	namespace := metav1.NamespaceDefault
	gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: gwName, UID: "gw-uid"}, Spec: kubeovnv1.VpcNatGatewaySpec{Vpc: vpcName, Subnet: subnetName}}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace, Name: svcName,
			Annotations: map[string]string{util.VpcNatGatewayAnnotation: gwName},
		}, Spec: v1.ServiceSpec{
			Type:       v1.ServiceTypeClusterIP,
			ClusterIP:  clusterIP,
			ClusterIPs: []string{clusterIP},
			Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
		},
	}
	slice := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace, Name: svcName + "-abc",
			Labels: map[string]string{discoveryv1.LabelServiceName: svcName},
		}, AddressType: discoveryv1.AddressTypeIPv4,
		Ports: []discoveryv1.EndpointPort{{Name: new("http"), Port: new(int32(8080)), Protocol: new(v1.ProtocolTCP)}},
		Endpoints: []discoveryv1.Endpoint{{
			Addresses:  []string{"10.0.0.7"},
			Conditions: discoveryv1.EndpointConditions{Ready: new(true)},
		}},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{{ObjectMeta: metav1.ObjectMeta{Name: vpcName}, Status: kubeovnv1.VpcStatus{TCPLoadBalancer: "vpc-tcp-load"}}},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		Subnets:        []*kubeovnv1.Subnet{{ObjectMeta: metav1.ObjectMeta{Name: subnetName}, Spec: kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Protocol: kubeovnv1.ProtocolIPv4, Vpc: vpcName, CIDRBlock: "10.0.0.0/24", Gateway: "10.0.0.1"}}},
		Services:       []*v1.Service{svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{slice},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.addOrUpdateGwNftableLbSvcQueue = newTypedRateLimitingQueue[string]("test-gw-nftable-lb-svc", nil)
	t.Cleanup(c.addOrUpdateGwNftableLbSvcQueue.ShutDown)

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(namespace+"/"+svcName))
	require.Eventually(t, func() bool { return c.addOrUpdateGwNftableLbSvcQueue.Len() == 1 },
		2*time.Second, 50*time.Millisecond, "finalizer claim schedules the programming pass")

	updated, err := c.config.KubeClient.CoreV1().Services(namespace).Get(context.Background(), svcName, metav1.GetOptions{})
	require.NoError(t, err)
	require.Contains(t, updated.Finalizers, util.KubeOVNControllerFinalizer)

	// No accounting record exists before the gateway was programmed.
	rules, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(),
		metav1.ListOptions{LabelSelector: util.NftableLbSvcNameLabel + "=" + svcName})
	require.NoError(t, err)
	require.Empty(t, rules.Items)
}

// Test_handleAddOrUpdateGwNftableLbService_claimBeforeRule pins the write order that keeps
// programmed share DNAT state accountable: the claim is persisted before the gateway is touched,
// so a record create rejected by admission fails the pass with an untouched data plane instead of
// an orphaned rule cleanup can never find. Ready and the ingress IP land only after the data plane
// converged.
func Test_handleAddOrUpdateGwNftableLbService_claimBeforeRule(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	// skip the finalizer-claim pass: this test drives the programming pass directly
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{f.slice},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{f.eip},
		Pods:           gatewayPods("owned-gw", "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.resetIptablesEipQueue = newTypedRateLimitingQueue[string]("test-reset-eip", nil)
	t.Cleanup(c.resetIptablesEipQueue.ShutDown)

	var execs []string
	c.execRulesInPod = func(_ *v1.Pod, operation string, _ []string) error {
		execs = append(execs, operation)
		return nil
	}

	// Simulated admission rejection of the accounting record (the webhook conflict check). The
	// captured claim must not carry status: Ready is published only after the data plane converged.
	// The mock OVN client has no expectations yet, so any NB call fails the test on its own.
	rejectRecords := true
	var claimed *kubeovnv1.IptablesDnatRule
	c.config.KubeOvnClient.(*kubeovnfake.Clientset).PrependReactor("create", "iptables-dnat-rules",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			claimed = action.(k8stesting.CreateAction).GetObject().(*kubeovnv1.IptablesDnatRule)
			if rejectRecords {
				return true, nil, errors.New("dnat rejected: conflicting identity")
			}
			return false, nil, nil
		})

	key := f.namespace + "/" + f.svc.Name
	require.Error(t, c.handleAddOrUpdateGwNftableLbService(key))
	require.NotNil(t, claimed)
	require.False(t, claimed.Status.Ready)
	require.Empty(t, execs, "a rejected claim must stop the pass before the gateway is touched")

	externalIDs := natGwVipRouteExternalIDs("owned-gw")
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority, externalIDs, false).Return(nil, nil)
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, externalIDs).Return(nil).Times(2)

	rejectRecords = false
	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(key))
	require.Contains(t, execs, natGwNftDnatMapAdd)
	require.Contains(t, execs, natGwVipHairpinAdd)
	require.Contains(t, execs, natGwVipAddrSync)

	rules, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(),
		metav1.ListOptions{LabelSelector: util.NftableLbSvcNameLabel + "=" + f.svc.Name})
	require.NoError(t, err)
	require.Len(t, rules.Items, 1)
	record := rules.Items[0]
	require.True(t, record.Status.Ready)
	require.Equal(t, "172.20.0.5", record.Status.V4ip)
	require.Equal(t, "owned-eip-uid", record.Labels[util.EipUIDLabel], "the persisted claim carries the EIP UID")

	updated, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []v1.LoadBalancerIngress{{IP: "172.20.0.5"}}, updated.Status.LoadBalancer.Ingress)
}

// Test_handleAddOrUpdateGwNftableLbService_waitsForReadyEip pins the bindability gate the DNAT
// worker used to provide (getBindableEip): a LoadBalancer Service whose EIP is not ready yet
// (e.g. during a gateway replacement redo, IP still carrying over) must not claim records with
// Ready status or publish the ingress IP before the EIP is live on the gateway. The pass waits
// and schedules a retry; the EIP's ready transition re-enqueues the Service.
func Test_handleAddOrUpdateGwNftableLbService_waitsForReadyEip(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}
	f.eip.Status.Ready = false

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		EndpointSlices: []*discoveryv1.EndpointSlice{f.slice},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{f.eip},
		Pods:           gatewayPods("owned-gw", "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.addOrUpdateGwNftableLbSvcQueue = newTypedRateLimitingQueue[string]("test-gw-nftable-lb-svc", nil)
	t.Cleanup(c.addOrUpdateGwNftableLbSvcQueue.ShutDown)

	var execs []string
	c.execRulesInPod = func(_ *v1.Pod, operation string, _ []string) error {
		execs = append(execs, operation)
		return nil
	}

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+f.svc.Name))
	require.Empty(t, execs, "an unready eip must stop the pass before the gateway is touched")
	require.Eventually(t, func() bool { return c.addOrUpdateGwNftableLbSvcQueue.Len() == 1 },
		5*time.Second, 100*time.Millisecond, "a retry is scheduled while the eip is not ready")

	rules, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(),
		metav1.ListOptions{LabelSelector: util.NftableLbSvcNameLabel + "=" + f.svc.Name})
	require.NoError(t, err)
	require.Empty(t, rules.Items, "no Ready record may be published before the eip is live")

	updated, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, updated.Status.LoadBalancer.Ingress)
}

// Test_handleAddOrUpdateGwNftableLbService_yieldKeepsWinnerIdentity pins the conflict handoff:
// when a new Service takes an EIP identity over, the previous owner's reconcile must not delete
// the nft map and hairpin the winner just programmed (its add already rewrote the complete
// backend set, and no event would wake the winner to repair a deletion).
func Test_handleAddOrUpdateGwNftableLbService_yieldKeepsWinnerIdentity(t *testing.T) {
	const (
		gwName    = "owned-gw"
		eipName   = "owned-eip"
		eipIP     = "172.20.0.5"
		loserName = "z-owner"
		loserIP   = "10.96.1.60"
		backend   = "10.0.7.21"
	)
	f := newNftableLbSvcOwnershipFixture()
	newLbSvc := func(name, clusterIP string) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: f.namespace, Name: name, UID: types.UID(name + "-uid"),
				Annotations: map[string]string{util.EipAnnotation: eipName, util.VpcNatGatewayAnnotation: gwName},
				Finalizers:  []string{util.KubeOVNControllerFinalizer},
			}, Spec: v1.ServiceSpec{
				Type:       v1.ServiceTypeLoadBalancer,
				ClusterIP:  clusterIP,
				ClusterIPs: []string{clusterIP},
				Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
			},
		}
	}
	// "a-winner" sorts before the previous owner and holds the identity's winner seat.
	winner, loser := newLbSvc("a-winner", "10.96.1.50"), newLbSvc(loserName, loserIP)
	slice := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: f.namespace, Name: loserName + "-abc",
			Labels: map[string]string{discoveryv1.LabelServiceName: loserName},
		}, AddressType: discoveryv1.AddressTypeIPv4,
		Ports:     []discoveryv1.EndpointPort{{Name: new("http"), Port: new(int32(8080)), Protocol: new(v1.ProtocolTCP)}},
		Endpoints: []discoveryv1.Endpoint{{Addresses: []string{backend}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}}},
	}
	// The previous owner's ledger from when it still held the identity.
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: nftableLbDnatRuleName(f.namespace, loserName, "tcp", "80", backend, "8080"),
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: f.namespace, util.NftableLbSvcNameLabel: loserName,
				util.NftableLbSvcUIDLabel: string(loser.UID), util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel: gwName,
				util.EipV4IpLabel:           eipIP, util.EipUIDLabel: "owned-eip-uid",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: eipName, ClusterIP: loserIP, ExternalPort: "80", Protocol: "tcp",
			InternalIP: backend, InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		}, Status: kubeovnv1.IptablesDnatRuleStatus{Ready: true, V4ip: eipIP, NatGwDp: gwName},
	}
	// The winner already claimed the identity: its record is the takeover witness (claim, rule,
	// Ready) that suppresses the loser's delete even before the winner's programming lands.
	winnerRecord := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: nftableLbDnatRuleName(f.namespace, "a-winner", "tcp", "80", "10.0.7.22", "8080"),
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: f.namespace, util.NftableLbSvcNameLabel: "a-winner",
				util.NftableLbSvcUIDLabel: string(winner.UID), util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel: gwName,
				util.EipV4IpLabel:           eipIP, util.EipUIDLabel: "owned-eip-uid",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: eipName, ClusterIP: "10.96.1.50", ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.7.22", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		}, Status: kubeovnv1.IptablesDnatRuleStatus{Ready: true, V4ip: eipIP, NatGwDp: gwName},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:              []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:           []*kubeovnv1.Subnet{f.subnet},
		Services:          []*v1.Service{winner, loser},
		EndpointSlices:    []*discoveryv1.EndpointSlice{slice},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{f.eip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{record, winnerRecord},
		Pods:              gatewayPods(gwName, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.addOrUpdateGwNftableLbSvcQueue = newTypedRateLimitingQueue[string]("test-gw-nftable-lb-svc", nil)
	t.Cleanup(c.addOrUpdateGwNftableLbSvcQueue.ShutDown)

	type execCall struct {
		op    string
		rules []string
	}
	var execs []execCall
	c.execRulesInPod = func(_ *v1.Pod, operation string, rules []string) error {
		execs = append(execs, execCall{operation, append([]string(nil), rules...)})
		return nil
	}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(gwName), false).Return(nil, nil)
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+loserName))

	for _, call := range execs {
		require.NotEqual(t, natGwNftDnatMapDel, call.op, "the loser must not delete the winner's nft map")
		require.NotEqual(t, natGwVipHairpinDel, call.op, "the loser must not delete the winner's hairpin")
	}
	added := false
	for _, call := range execs {
		if call.op == natGwNftDnatMapAdd {
			added = true
			for _, rule := range call.rules {
				require.NotContains(t, rule, eipIP, "the loser must not program the EIP identity it lost")
			}
		}
	}
	require.True(t, added, "the loser still programs its own ClusterIP identity")

	updated, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, updated.Spec.EIP, "the loser's claim drops the identity it yielded")
	require.Equal(t, loserIP, updated.Spec.ClusterIP)
}

// Test_handleAddOrUpdateGwNftableLbService_backendlessWinnerKeepsNoStaleMap pins the other half
// of the handoff: a winner with no ready endpoints claims nothing and programs nothing, so the
// yielded identity has no takeover to protect - the loser must tear its own stale map and
// hairpin down instead of leaving the address answered with the loser's backends while claimed
// by nobody. When the winner gains endpoints later, its endpoint event reconciles and programs
// the identity itself.
func Test_handleAddOrUpdateGwNftableLbService_backendlessWinnerKeepsNoStaleMap(t *testing.T) {
	const (
		gwName    = "owned-gw"
		eipName   = "owned-eip"
		eipIP     = "172.20.0.5"
		loserName = "z-owner"
		loserIP   = "10.96.1.60"
		backend   = "10.0.7.21"
	)
	f := newNftableLbSvcOwnershipFixture()
	newLbSvc := func(name, clusterIP string) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: f.namespace, Name: name, UID: types.UID(name + "-uid"),
				Annotations: map[string]string{util.EipAnnotation: eipName, util.VpcNatGatewayAnnotation: gwName},
				Finalizers:  []string{util.KubeOVNControllerFinalizer},
			}, Spec: v1.ServiceSpec{
				Type:       v1.ServiceTypeLoadBalancer,
				ClusterIP:  clusterIP,
				ClusterIPs: []string{clusterIP},
				Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
			},
		}
	}
	// "a-winner" sorts first but has no endpoints at all: no slice exists for it.
	winner, loser := newLbSvc("a-winner", "10.96.1.50"), newLbSvc(loserName, loserIP)
	loserSlice := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: f.namespace, Name: loserName + "-abc",
			Labels: map[string]string{discoveryv1.LabelServiceName: loserName},
		}, AddressType: discoveryv1.AddressTypeIPv4,
		Ports:     []discoveryv1.EndpointPort{{Name: new("http"), Port: new(int32(8080)), Protocol: new(v1.ProtocolTCP)}},
		Endpoints: []discoveryv1.Endpoint{{Addresses: []string{backend}, Conditions: discoveryv1.EndpointConditions{Ready: new(true)}}},
	}
	// The previous owner's ledger from when it still held the identity.
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: nftableLbDnatRuleName(f.namespace, loserName, "tcp", "80", backend, "8080"),
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: f.namespace, util.NftableLbSvcNameLabel: loserName,
				util.NftableLbSvcUIDLabel: string(loser.UID), util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel: gwName,
				util.EipV4IpLabel:           eipIP, util.EipUIDLabel: "owned-eip-uid",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: eipName, ClusterIP: loserIP, ExternalPort: "80", Protocol: "tcp",
			InternalIP: backend, InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		}, Status: kubeovnv1.IptablesDnatRuleStatus{Ready: true, V4ip: eipIP, NatGwDp: gwName},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:              []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:           []*kubeovnv1.Subnet{f.subnet},
		Services:          []*v1.Service{winner, loser},
		EndpointSlices:    []*discoveryv1.EndpointSlice{loserSlice},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{f.eip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{record},
		Pods:              gatewayPods(gwName, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true
	c.addOrUpdateGwNftableLbSvcQueue = newTypedRateLimitingQueue[string]("test-gw-nftable-lb-svc", nil)
	t.Cleanup(c.addOrUpdateGwNftableLbSvcQueue.ShutDown)
	c.resetIptablesEipQueue = newTypedRateLimitingQueue[string]("test-reset-eip", nil)
	t.Cleanup(c.resetIptablesEipQueue.ShutDown)

	type execCall struct {
		op    string
		rules []string
	}
	var execs []execCall
	c.execRulesInPod = func(_ *v1.Pod, operation string, rules []string) error {
		execs = append(execs, execCall{operation, append([]string(nil), rules...)})
		return nil
	}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(gwName), false).Return(nil, nil).AnyTimes()
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+loserName))

	var staleDeletes, clusterIPAdds int
	for _, call := range execs {
		for _, rule := range call.rules {
			if call.op == natGwNftDnatMapDel && strings.Contains(rule, eipIP) {
				staleDeletes++
			}
			if call.op == natGwNftDnatMapAdd && strings.Contains(rule, loserIP) {
				clusterIPAdds++
			}
		}
	}
	require.Positive(t, staleDeletes, "the loser tears down the stale EIP map the backend-less winner does not claim")
	require.Positive(t, clusterIPAdds, "the loser keeps programming its own ClusterIP identity")

	updated, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, updated.Spec.EIP, "the loser's claim drops the identity it yielded")

	// The winner's own pass: nothing to claim and nothing to program, but nothing to delete
	// either — the stale map is already gone.
	execs = nil
	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+winner.Name))
	for _, call := range execs {
		require.NotEqual(t, natGwNftDnatMapDel, call.op, "the winner has no stale wiring to remove")
		if call.op == natGwNftDnatMapAdd {
			for _, rule := range call.rules {
				require.NotContains(t, rule, eipIP, "a backend-less winner must not program the EIP identity")
			}
		}
	}
	published, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), winner.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []v1.LoadBalancerIngress{{IP: eipIP}}, published.Status.LoadBalancer.Ingress,
		"the winner reports its ingress address; with no backends the address now refuses traffic instead of hitting the loser")
}

// Test_cleanupNftableLbService_keepsSuccessorIdentity pins the identity handover on the
// teardown path: the conflict resolver lets a successor take over a terminating Service's EIP
// identity, and the dying owner's cleanup must not delete the nft map and hairpin the
// successor's claim (and possibly programming) covers - no event would wake the successor to
// repair that. The owner's unmatched identities (here its ClusterIP leg) are still removed.
func Test_cleanupNftableLbService_keepsSuccessorIdentity(t *testing.T) {
	const (
		gwName      = "owned-gw"
		eipName     = "owned-eip"
		eipIP       = "172.20.0.5"
		ownerName   = "z-owner"
		ownerIP     = "10.96.1.60"
		successorIP = "10.96.1.50"
		backend     = "10.0.7.21"
	)
	f := newNftableLbSvcOwnershipFixture()
	newLbSvc := func(name, clusterIP string, uid types.UID) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: f.namespace, Name: name, UID: uid,
				Annotations: map[string]string{util.EipAnnotation: eipName, util.VpcNatGatewayAnnotation: gwName},
				Finalizers:  []string{util.KubeOVNControllerFinalizer},
			}, Spec: v1.ServiceSpec{
				Type:       v1.ServiceTypeLoadBalancer,
				ClusterIP:  clusterIP,
				ClusterIPs: []string{clusterIP},
				Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
			},
		}
	}
	owner := newLbSvc(ownerName, ownerIP, types.UID(ownerName+"-uid"))
	now := metav1.Now()
	owner.DeletionTimestamp = &now
	successor := newLbSvc("a-winner", successorIP, types.UID("a-winner-uid"))
	newRecord := func(svc *v1.Service, clusterIP string) *kubeovnv1.IptablesDnatRule {
		return &kubeovnv1.IptablesDnatRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: nftableLbDnatRuleName(f.namespace, svc.Name, "tcp", "80", backend, "8080"),
				Labels: map[string]string{
					util.NftableLbSvcNsLabel: f.namespace, util.NftableLbSvcNameLabel: svc.Name,
					util.NftableLbSvcUIDLabel: string(svc.UID), util.NftableLbSvcRecordLabel: "true",
					util.VpcNatGatewayNameLabel: gwName,
					util.EipV4IpLabel:           eipIP, util.EipUIDLabel: "owned-eip-uid",
				},
			}, Spec: kubeovnv1.IptablesDnatRuleSpec{
				EIP: eipName, ClusterIP: clusterIP, ExternalPort: "80", Protocol: "tcp",
				InternalIP: backend, InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
			}, Status: kubeovnv1.IptablesDnatRuleStatus{Ready: true, V4ip: eipIP, NatGwDp: gwName},
		}
	}
	ownerRecord, successorRecord := newRecord(owner, ownerIP), newRecord(successor, successorIP)

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:              []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:           []*kubeovnv1.Subnet{f.subnet},
		Services:          []*v1.Service{owner, successor},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{f.eip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{ownerRecord, successorRecord},
		Pods:              gatewayPods(gwName, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true

	type execCall struct {
		op    string
		rules []string
	}
	var execs []execCall
	c.execRulesInPod = func(_ *v1.Pod, operation string, rules []string) error {
		execs = append(execs, execCall{operation, append([]string(nil), rules...)})
		return nil
	}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(gwName), false).Return(nil, nil).AnyTimes()
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	require.NoError(t, c.cleanupNftableLbService(owner, f.namespace, ownerName))

	var eipDeletes, clusterIPDeletes int
	for _, call := range execs {
		if call.op != natGwNftDnatMapDel && call.op != natGwVipHairpinDel {
			continue
		}
		for _, rule := range call.rules {
			if strings.Contains(rule, eipIP) {
				eipDeletes++
			}
			if call.op == natGwNftDnatMapDel && strings.Contains(rule, ownerIP) {
				clusterIPDeletes++
			}
		}
	}
	require.Zero(t, eipDeletes, "the terminating owner must not delete the identity the successor claimed")
	require.Positive(t, clusterIPDeletes, "the old owner's own ClusterIP identity is still removed")

	_, err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), ownerRecord.Name, metav1.GetOptions{})
	require.True(t, k8serrors.IsNotFound(err), "the old owner's record is retired")
	_, err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), successorRecord.Name, metav1.GetOptions{})
	require.NoError(t, err, "the successor's record is untouched")
}

// Test_cleanupNftableLbService_recordlessTeardownDerivesIntent pins the record-deletion window:
// a Service whose accounting records were deleted can enter teardown before the restore pass
// ran. The ledger list then comes back empty, but the gateway still holds whatever the Service
// programmed - the cleanup must derive the identities from the Service's own intent (ports,
// EIP address, ClusterIP), run the data-plane teardown, and only then release the finalizer.
func Test_cleanupNftableLbService_recordlessTeardownDerivesIntent(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	now := metav1.Now()
	f.svc.DeletionTimestamp = &now
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}

	// Every accounting record of the Service is gone (deleted out from under it; nothing was
	// restored). The gateway keeps running, so the rules the records accounted for still exist.
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:        []*kubeovnv1.Subnet{f.subnet},
		Services:       []*v1.Service{f.svc},
		IptablesEIPs:   []*kubeovnv1.IptablesEIP{f.eip},
		Pods:           gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true
	c.config.EnableGwNftableSvcClusterIP = true

	type execCall struct {
		op    string
		rules []string
	}
	var execs []execCall
	c.execRulesInPod = func(_ *v1.Pod, operation string, rules []string) error {
		execs = append(execs, execCall{operation, append([]string(nil), rules...)})
		return nil
	}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(f.namespace+"/"+f.svc.Name))

	var delVips []string
	addrSynced := false
	for _, call := range execs {
		switch call.op {
		case natGwNftDnatMapDel:
			for _, rule := range call.rules {
				delVips = append(delVips, strings.Split(rule, ",")[0])
			}
		case natGwVipAddrSync:
			addrSynced = true
		}
	}
	require.Contains(t, delVips, "172.20.0.5", "the EIP identity of the deleted-record service is torn down")
	require.Contains(t, delVips, "10.96.1.5", "the ClusterIP identity of the deleted-record service is torn down")
	require.True(t, addrSynced, "the lo address sync also runs without a surviving record")

	released, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, released.Finalizers, "the finalizer is released only after the data-plane teardown")

	rules, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)
	require.Empty(t, rules.Items, "teardown creates no accounting records")
}

func TestClearNftableLbSvcIngressIP(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  metav1.NamespaceDefault,
			Name:       "web",
			Finalizers: []string{util.KubeOVNControllerFinalizer},
		}, Status: v1.ServiceStatus{LoadBalancer: v1.LoadBalancerStatus{
			Ingress: []v1.LoadBalancerIngress{{IP: "203.0.113.10"}},
		}},
	}
	client := k8sfake.NewSimpleClientset(svc)
	c := &Controller{config: &Configuration{KubeClient: client}}

	require.NoError(t, c.clearNftableLbSvcIngressIP(svc))
	actions := client.Actions()
	require.Len(t, actions, 1)
	require.Equal(t, "status", actions[0].GetSubresource())
	updated, err := client.CoreV1().Services(svc.Namespace).Get(context.Background(), svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, updated.Status.LoadBalancer.Ingress)
	require.Contains(t, updated.Finalizers, util.KubeOVNControllerFinalizer)
}

// Test_shareClusterIPBackends pins the decision the teardown of one share rule makes for the internal
// VIP it serves: the siblings of the ClusterIP identity are the live rules of the gateway serving
// that same address, which is not the same set as the siblings of the EIP identity (sameDnatIdentity
// compares the EIP whenever either rule has one). A rule that is the only one serving its ClusterIP
// therefore gets an empty result, which is what makes the caller drop the address instead of
// rebuilding it with the backends of a rule that serves another one.
func TestCleanupNftableLbServiceIgnoresPredecessorRules(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	f.svc.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{{IP: "172.20.0.9"}}

	stale := buildDesiredNftableLbDnatRules(f.svc, f.eip.Name, f.gw.Name, []*discoveryv1.EndpointSlice{f.slice}, testNftableLbBackendIP)
	require.Len(t, stale, 1)
	staleRule := slices.Collect(maps.Values(stale))[0]
	staleRule.Labels[util.NftableLbSvcUIDLabel] = "a-previous-incarnation"

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:              []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:           []*kubeovnv1.Subnet{f.subnet},
		Services:          []*v1.Service{f.svc},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{f.eip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{staleRule},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.EnableGwNftableLbSvc = true

	require.NoError(t, c.cleanupNftableLbService(f.svc, f.namespace, f.svc.Name))

	updated, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []v1.LoadBalancerIngress{{IP: "172.20.0.9"}}, updated.Status.LoadBalancer.Ingress,
		"a predecessor's rule must not authorize clearing this Service's ingress IP")

	_, err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), staleRule.Name, metav1.GetOptions{})
	require.True(t, k8serrors.IsNotFound(err), "the predecessor's rule must still be deleted")
}

func TestBuildNftableLbIdentities(t *testing.T) {
	t.Parallel()

	records := map[string]*kubeovnv1.IptablesDnatRule{
		"a": {
			Spec: kubeovnv1.IptablesDnatRuleSpec{
				EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
				InternalIP: "10.0.0.2", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
			},
		},
		"b": {
			Spec: kubeovnv1.IptablesDnatRuleSpec{
				EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
				InternalIP: "10.0.0.3", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
			},
		},
	}

	got := buildNftableLbIdentities(records, "203.0.113.10")
	require.Len(t, got, 2)
	require.ElementsMatch(t, []string{"10.0.0.2:8080", "10.0.0.3:8080"}, got["203.0.113.10/80/tcp"].backends)
	require.ElementsMatch(t, []string{"10.0.0.2:8080", "10.0.0.3:8080"}, got["10.96.1.5/80/tcp"].backends)
}

func TestBuildNftableLbPrograms(t *testing.T) {
	t.Parallel()

	records := map[string]*kubeovnv1.IptablesDnatRule{
		"a": {Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "TCP",
			InternalIP: "10.0.0.2", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		}},
	}
	programs := buildNftableLbPrograms(records, "203.0.113.10")
	require.Len(t, programs, 2)
	require.Equal(t, "10.96.1.5,80,tcp", programs[0].hairpinRule)
	require.Equal(t, "203.0.113.10,80,tcp", programs[1].hairpinRule)
	require.Equal(t, []string{"10.0.0.2:8080"}, programs[0].identity.backends)
	require.Equal(t, []string{"10.0.0.2:8080"}, programs[1].identity.backends)
}

func TestNftableLbExistingIdentities(t *testing.T) {
	t.Parallel()

	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{util.EipV4IpLabel: "203.0.113.10"}}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip0", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "TCP",
		},
	}
	got := nftableLbExistingIdentities([]*kubeovnv1.IptablesDnatRule{record})
	require.Contains(t, got, "203.0.113.10/80/tcp")
	require.Contains(t, got, "10.96.1.5/80/tcp")
}

// Test_gatewayAnnotatedServiceKeepsItsClassicLifecycle pins the ownership partition between the
// two implementations: the gateway nftable LB mode owns only a Service's EIP identity, so the
// classic OVN paths keep managing its ClusterIP VIPs and backends. Without that, adding the
// gateway annotations would freeze the ClusterIP's backend set, and deleting the Service would
// leak its ClusterIP VIP in the OVN load balancer tables.
func Test_gatewayAnnotatedServiceKeepsItsClassicLifecycle(t *testing.T) {
	t.Parallel()

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default", Name: "web",
			Annotations: map[string]string{util.VpcNatGatewayAnnotation: "gw0", util.EipAnnotation: "eip0"},
		}, Spec: v1.ServiceSpec{
			Type:       v1.ServiceTypeLoadBalancer,
			ClusterIP:  "10.96.1.5",
			ClusterIPs: []string{"10.96.1.5"},
			Ports:      []v1.ServicePort{{Name: "http", Port: 80, Protocol: v1.ProtocolTCP}},
		},
	}

	t.Run("updates keep flowing to the OVN path", func(t *testing.T) {
		gwQueue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
		delQueue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[*vpcService]())
		updateQueue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[*updateSvcObject]())
		t.Cleanup(gwQueue.ShutDown)
		t.Cleanup(delQueue.ShutDown)
		t.Cleanup(updateQueue.ShutDown)
		c := &Controller{
			config:                         &Configuration{EnableOvnLB: true, EnableGwNftableLbSvc: true},
			addOrUpdateGwNftableLbSvcQueue: gwQueue,
			deleteServiceQueue:             delQueue,
			updateServiceQueue:             updateQueue,
		}

		oldSvc := svc.DeepCopy()
		oldSvc.ResourceVersion = "1"
		oldSvc.Annotations = nil
		newSvc := svc.DeepCopy()
		newSvc.ResourceVersion = "2"
		c.enqueueUpdateService(oldSvc, newSvc)

		require.Equal(t, 1, gwQueue.Len(), "gaining the annotations wakes the gateway reconcile")
		require.Equal(t, 0, delQueue.Len(), "nothing is torn down for an update")
		require.Equal(t, 1, updateQueue.Len(), "the OVN ClusterIP VIP resync keeps running for the gateway-managed service")
	})

	t.Run("deletes release the OVN ClusterIP VIPs", func(t *testing.T) {
		delQueue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[*vpcService]())
		t.Cleanup(delQueue.ShutDown)
		c := &Controller{
			config:             &Configuration{EnableOvnLB: true, EnableGwNftableLbSvc: true},
			deleteServiceQueue: delQueue,
		}

		c.enqueueDeleteService(svc.DeepCopy())
		require.Equal(t, 1, delQueue.Len(), "the ClusterIP VIP cleanup of a gateway-managed service runs")
		item, _ := delQueue.Get()
		require.Equal(t, []string{"10.96.1.5:80"}, item.Vips)
	})
}

// Test_handleAddOrUpdateGwNftableLbService_ingressClearIsRepeatable pins the retry semantics of
// withdrawing the ingress IP on the gate flip: the address is withdrawn before the records stop
// documenting the EIP ownership, so a failed status write is retried from intact evidence instead
// of leaving the stale EXTERNAL-IP behind forever.
func Test_handleAddOrUpdateGwNftableLbService_ingressClearIsRepeatable(t *testing.T) {
	f := newNftableLbSvcOwnershipFixture()
	f.svc.Finalizers = []string{util.KubeOVNControllerFinalizer}
	f.svc.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{{IP: "172.20.0.5"}}
	record := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: nftableLbDnatRuleName(f.namespace, f.svc.Name, "tcp", "80", "10.0.7.2", "8080"),
			Labels: map[string]string{
				util.NftableLbSvcNsLabel: f.namespace, util.NftableLbSvcNameLabel: f.svc.Name,
				util.NftableLbSvcUIDLabel: string(f.svc.UID), util.NftableLbSvcRecordLabel: "true",
				util.VpcNatGatewayNameLabel: f.gw.Name,
				util.EipV4IpLabel:           "172.20.0.5", util.EipUIDLabel: "owned-eip-uid",
			},
		}, Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "owned-eip", ClusterIP: "10.96.1.5", ExternalPort: "80", Protocol: "tcp",
			InternalIP: "10.0.7.2", InternalPort: "8080", Type: kubeovnv1.DnatRuleTypeShare,
		}, Status: kubeovnv1.IptablesDnatRuleStatus{Ready: true, V4ip: "172.20.0.5", NatGwDp: f.gw.Name},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs:              []*kubeovnv1.Vpc{f.vpc},
		VpcNatGateways:    []*kubeovnv1.VpcNatGateway{f.gw},
		Subnets:           []*kubeovnv1.Subnet{f.subnet},
		Services:          []*v1.Service{f.svc},
		EndpointSlices:    []*discoveryv1.EndpointSlice{f.slice},
		IptablesEIPs:      []*kubeovnv1.IptablesEIP{f.eip},
		IptablesDnatRules: []*kubeovnv1.IptablesDnatRule{record},
		Pods:              gatewayPods(f.gw.Name, "10.0.7.254"),
	})
	require.NoError(t, err)
	c := fc.fakeController
	// The EIP gate is off, the ClusterIP gate stays on: the EIP identity retires.
	c.config.EnableGwNftableLbSvc = false
	c.config.EnableGwNftableSvcClusterIP = true

	kubeClient, ok := c.config.KubeClient.(*k8sfake.Clientset)
	require.True(t, ok)
	injected := true
	kubeClient.PrependReactor("update", "services", func(action k8stesting.Action) (bool, runtime.Object, error) {
		if injected && action.(k8stesting.UpdateAction).GetSubresource() == "status" {
			injected = false
			return true, nil, errors.New("injected status write failure")
		}
		return false, nil, nil
	})

	var execs []string
	c.execRulesInPod = func(_ *v1.Pod, operation string, _ []string) error {
		execs = append(execs, operation)
		return nil
	}
	fc.mockOvnClient.EXPECT().ListLogicalRouterPolicies(util.DefaultVpc, util.NatGatewayVipPolicyPriority,
		natGwVipRouteExternalIDs(f.gw.Name), false).Return(nil, nil).AnyTimes()
	fc.mockOvnClient.EXPECT().AddLogicalRouterPolicy(util.DefaultVpc, util.NatGatewayVipPolicyPriority, gomock.Any(),
		string(kubeovnv1.PolicyRouteActionReroute), []string{"10.0.7.254"}, nil, gomock.Any()).Return(nil).AnyTimes()

	key := f.namespace + "/" + f.svc.Name
	require.Error(t, c.handleAddOrUpdateGwNftableLbService(key), "the failed withdrawal aborts the pass")

	stillService, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Len(t, stillService.Status.LoadBalancer.Ingress, 1, "the ingress stays published while the withdrawal could not complete")
	stillRecord, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "owned-eip", stillRecord.Spec.EIP, "the ledger still documents the EIP claim, so the retry sees the same ownership")
	require.Empty(t, execs, "nothing was reprogrammed before the address was withdrawn")

	require.NoError(t, c.handleAddOrUpdateGwNftableLbService(key), "the retry withdraws the address and retires the identity")
	finished, err := c.config.KubeClient.CoreV1().Services(f.namespace).Get(context.Background(), f.svc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, finished.Status.LoadBalancer.Ingress)
	stripped, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Get(context.Background(), record.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, stripped.Spec.EIP, "the ledger dropped the EIP only after the address was withdrawn")
	require.Contains(t, execs, natGwNftDnatMapDel, "the EIP identity was removed from the gateway")
	require.Contains(t, execs, natGwNftDnatMapAdd, "the ClusterIP identity kept being served")
}
