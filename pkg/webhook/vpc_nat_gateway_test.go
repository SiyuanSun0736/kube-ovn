package webhook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	ovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestValidateIptablesDnatProtocol(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		protocol string
		wantErr  bool
	}{
		{protocol: "tcp"},
		{protocol: "udp"},
		{protocol: "TCP", wantErr: true},
		{protocol: "UDP", wantErr: true},
		{protocol: "Tcp", wantErr: true},
		{protocol: "icmp", wantErr: true},
	} {
		t.Run(test.protocol, func(t *testing.T) {
			err := validateIptablesDnatProtocol(test.protocol)
			if test.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "must be lowercase")
				return
			}
			require.NoError(t, err)
		})
	}
}

// TestEipUIDSelectorIsolatesGenerations pins admission to the same credential the controller's
// in-use check counts. Two EIPs can carry the same address, so selecting by address made the
// webhook block a deletion because of a rule belonging to a different EIP.
func TestEipUIDSelectorIsolatesGenerations(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	require.NoError(t, ovnv1.AddToScheme(scheme))

	mine := &ovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "mine", UID: "mine-uid"},
		Status:     ovnv1.IptablesEIPStatus{IP: "1.1.1.1"},
	}
	other := &ovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "other", UID: "other-uid"},
		Status:     ovnv1.IptablesEIPStatus{IP: "1.1.1.1"},
	}
	// Belongs to "other" but shares the address, which is exactly what used to be miscounted.
	fip := &ovnv1.IptablesFIPRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "other-fip",
			Labels: map[string]string{
				util.EipV4IpLabel: "1.1.1.1",
				util.EipUIDLabel:  "other-uid",
			},
		},
		Spec: ovnv1.IptablesFIPRuleSpec{EIP: "other"},
	}
	reader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(mine, other, fip).Build()

	list := &ovnv1.IptablesFIPRuleList{}
	require.NoError(t, reader.List(t.Context(), list, eipUIDSelector(mine)))
	require.Empty(t, list.Items, "a rule owned by another EIP must not count")

	list = &ovnv1.IptablesFIPRuleList{}
	require.NoError(t, reader.List(t.Context(), list, eipUIDSelector(other)))
	require.Len(t, list.Items, 1)
}

// TestValidateQoSPolicyRef covers the admission guard shared by IptablesEIP and VpcNatGateway.
// The controller keeps a referrer pointing at a missing or terminating policy out of Ready, so
// admission has to reject it up front instead of leaving the user with a silent retry loop.
func TestValidateQoSPolicyRef(t *testing.T) {
	t.Parallel()

	now := metav1.Now()
	scheme := runtime.NewScheme()
	require.NoError(t, ovnv1.AddToScheme(scheme))
	reader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(
		&ovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{
			Name:       "live-qos",
			Finalizers: []string{util.KubeOVNControllerFinalizer},
		}},
		&ovnv1.QoSPolicy{
			ObjectMeta: metav1.ObjectMeta{Name: "pending-qos", Finalizers: []string{util.KubeOVNControllerFinalizer}},
			Spec:       ovnv1.QoSPolicySpec{BindingType: ovnv1.QoSBindingTypeEIP},
		},
		&ovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{Name: "new-qos"}},
		&ovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{
			Name:              "dying-qos",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		}},
	).Build()

	t.Run("empty reference is allowed", func(t *testing.T) {
		require.NoError(t, validateQoSPolicyRef(t.Context(), reader, ""))
	})

	t.Run("existing policy is allowed", func(t *testing.T) {
		require.NoError(t, validateQoSPolicyRef(t.Context(), reader, "live-qos"))
	})

	t.Run("policy without controller reconcile is rejected", func(t *testing.T) {
		require.ErrorContains(t, validateQoSPolicyRef(t.Context(), reader, "new-qos"), "not ready")
	})

	t.Run("policy with stale status is rejected", func(t *testing.T) {
		require.ErrorContains(t, validateQoSPolicyRef(t.Context(), reader, "pending-qos"), "not ready")
	})

	t.Run("missing policy is rejected", func(t *testing.T) {
		err := validateQoSPolicyRef(t.Context(), reader, "missing-qos")
		require.Error(t, err)
		require.True(t, k8serrors.IsNotFound(err))
		require.ErrorContains(t, err, "create it before referencing it")
	})

	t.Run("terminating policy is rejected", func(t *testing.T) {
		err := validateQoSPolicyRef(t.Context(), reader, "dying-qos")
		require.ErrorContains(t, err, "terminating")
		require.ErrorContains(t, err, "wait for its deletion to complete")
	})
}

// TestIptablesDnatUpdateServiceRecord pins that the mutable-accounting escape hatch for Service
// records honors the same identity semantics everywhere: an empty controller identity disables
// identity checking, so it must not funnel record updates into the immutable branches, while the
// record's content is still validated.
// hookCache is a cache.Cache limited to reads on the objects the test seeded; the informer
// methods of the embedded interface are never exercised by the admission paths under test.
type hookCache struct {
	cache.Cache
	reader client.Reader
}

func (f *hookCache) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	return f.reader.Get(ctx, key, obj, opts...)
}

func (f *hookCache) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return f.reader.List(ctx, list, opts...)
}

func TestIptablesDnatUpdateServiceRecord(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	require.NoError(t, ovnv1.AddToScheme(scheme))
	gw := &ovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw0"}}
	// hookCache serves only reads: the embedded nil cache.Cache satisfies the rest of the
	// interface the hook never calls in these paths.
	reader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(gw).Build()
	cacheClient := &hookCache{reader: reader}

	newRecord := func(affinity string) *ovnv1.IptablesDnatRule {
		return &ovnv1.IptablesDnatRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "web.80.tcp.10.0.0.5.8080",
				Labels: map[string]string{
					util.NftableLbSvcNsLabel:     "default",
					util.NftableLbSvcNameLabel:   "web",
					util.NftableLbSvcUIDLabel:    "svc-uid",
					util.NftableLbSvcRecordLabel: "true",
					util.VpcNatGatewayNameLabel:  "gw0",
					util.VpcDnatEPortLabel:       "80",
				},
			}, Spec: ovnv1.IptablesDnatRuleSpec{
				Type: ovnv1.DnatRuleTypeShare, ClusterIP: "10.96.1.5",
				ExternalPort: "80", Protocol: "tcp", InternalIP: "10.0.0.5", InternalPort: "8080",
				SessionAffinity: affinity, SessionAffinityTimeoutSeconds: func() int32 {
					if affinity == ovnv1.DnatSessionAffinityClientIP {
						return 300
					}
					return 0
				}(),
			},
		}
	}
	dnatUpdate := func(t *testing.T, username string, oldObj, newObj any) admission.Request {
		t.Helper()
		req := updateRequest(t, oldObj, newObj)
		req.UserInfo.Username = username
		return req
	}
	newHook := func(identity string) *ValidatingHook {
		return &ValidatingHook{decoder: admission.NewDecoder(scheme), cache: cacheClient, controllerIdentity: identity}
	}

	old := newRecord(ovnv1.DnatSessionAffinityNone)
	updated := newRecord(ovnv1.DnatSessionAffinityClientIP)

	t.Run("empty identity still lets the accounting update through", func(t *testing.T) {
		resp := newHook("").iptablesDnatUpdateHook(t.Context(), dnatUpdate(t, "some-controller-user", old, updated))
		require.True(t, resp.Allowed, "with identity checking disabled the record update must not hit the immutable branches")
	})
	t.Run("empty identity still validates the record content", func(t *testing.T) {
		broken := updated.DeepCopy()
		delete(broken.Labels, util.VpcNatGatewayNameLabel)
		resp := newHook("").iptablesDnatUpdateHook(t.Context(), dnatUpdate(t, "some-controller-user", old, broken))
		require.False(t, resp.Allowed)
		require.Contains(t, resp.Result.Message, "gateway label is required")
	})
	t.Run("matching identity keeps the existing behavior", func(t *testing.T) {
		resp := newHook("system:serviceaccount:kube-system:kube-ovn-controller").iptablesDnatUpdateHook(
			t.Context(), dnatUpdate(t, "system:serviceaccount:kube-system:kube-ovn-controller", old, updated),
		)
		require.True(t, resp.Allowed)
	})
	t.Run("a stranger is still held by the immutable branches", func(t *testing.T) {
		resp := newHook("system:serviceaccount:kube-system:kube-ovn-controller").iptablesDnatUpdateHook(
			t.Context(), dnatUpdate(t, "system:serviceaccount:other:intruder", old, updated),
		)
		require.False(t, resp.Allowed)
		require.Contains(t, resp.Result.Message, "immutable")
	})
}

func TestVpcNatGatewayDataplaneModeImmutable(t *testing.T) {
	t.Parallel()
	scheme := runtime.NewScheme()
	require.NoError(t, ovnv1.AddToScheme(scheme))

	hook := &ValidatingHook{
		decoder: admission.NewDecoder(scheme),
	}

	oldGw := &ovnv1.VpcNatGateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gw-test",
			Annotations: map[string]string{
				util.NatGatewayDataplaneModeAnnotation: "stateless",
			},
		},
		Spec: ovnv1.VpcNatGatewaySpec{
			Vpc:       "test-vpc",
			Subnet:    "test-subnet",
			LanIP:     "10.0.0.1",
			Namespace: "test-ns",
		},
	}

	newGw := oldGw.DeepCopy()
	newGw.Annotations[util.NatGatewayDataplaneModeAnnotation] = "stateful"

	req := updateRequest(t, oldGw, newGw)
	resp := hook.VpcNatGwCreateOrUpdateHook(context.Background(), req)
	require.False(t, resp.Allowed)
	require.Contains(t, resp.Result.Message, "is immutable once initialized")
}

func TestStatelessHAVpcNatGatewayRejectsMacvlanSubnet(t *testing.T) {
	t.Parallel()
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, ovnv1.AddToScheme(scheme))

	natCm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatConfig, Namespace: metav1.NamespaceSystem},
		Data:       map[string]string{"natGwNamePrefix": "vpc-nat-gw"},
	}
	vpc := &ovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: "test-vpc"}}
	lanSub := &ovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "lan-sub"},
		Spec:       ovnv1.SubnetSpec{CIDRBlock: "10.0.0.0/24"},
	}
	macvlanSub := &ovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "macvlan-ext",
			Annotations: map[string]string{util.NadMacvlanMasterAnnotation: "bond1"},
			Labels:      map[string]string{util.NadMacvlanTypeLabel: "true"},
		},
		Spec: ovnv1.SubnetSpec{CIDRBlock: "183.240.180.128/27"},
	}
	vethSub := &ovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "veth-ext"},
		Spec:       ovnv1.SubnetSpec{CIDRBlock: "183.240.180.128/27"},
	}
	statelessGw := &ovnv1.VpcNatGateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "ha-gw-m1",
			Annotations: map[string]string{util.NatGatewayDataplaneModeAnnotation: "stateless"},
		},
		Spec: ovnv1.VpcNatGatewaySpec{
			Vpc:             "test-vpc",
			Subnet:          "lan-sub",
			LanIP:           "10.0.0.254",
			ExternalSubnets: []string{"veth-ext"},
		},
	}

	reader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(natCm, vpc, lanSub, macvlanSub, vethSub, statelessGw).Build()
	hook := &ValidatingHook{
		decoder: admission.NewDecoder(scheme),
		cache:   &hookCache{reader: reader},
	}

	t.Run("stateless HA gateway with veth subnet is allowed", func(t *testing.T) {
		require.NoError(t, hook.ValidateVpcNatGW(t.Context(), statelessGw))
	})

	t.Run("stateless HA gateway with macvlan subnet is rejected", func(t *testing.T) {
		badGw := statelessGw.DeepCopy()
		badGw.Spec.ExternalSubnets = []string{"macvlan-ext"}
		err := hook.ValidateVpcNatGW(t.Context(), badGw)
		require.ErrorContains(t, err, "cannot use macvlan external subnet")
	})

	t.Run("BGP Anycast EIP on macvlan subnet is rejected", func(t *testing.T) {
		bgpEip := &ovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "bgp-eip",
				Annotations: map[string]string{util.BgpAnnotation: "true"},
			},
			Spec: ovnv1.IptablesEIPSpec{
				NatGwDp:        "ha-gw-m1",
				ExternalSubnet: "macvlan-ext",
				V4ip:           "183.240.180.136",
			},
		}
		err := hook.ValidateIptablesEIP(t.Context(), bgpEip)
		require.ErrorContains(t, err, "cannot use macvlan external subnet")
	})

	t.Run("EIP without BGP annotation but attached to stateless HA gateway on macvlan subnet is rejected", func(t *testing.T) {
		eip := &ovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{Name: "plain-eip-on-ha-gw"},
			Spec: ovnv1.IptablesEIPSpec{
				NatGwDp:        "ha-gw-m1",
				ExternalSubnet: "macvlan-ext",
				V4ip:           "183.240.180.136",
			},
		}
		err := hook.ValidateIptablesEIP(t.Context(), eip)
		require.ErrorContains(t, err, "cannot use macvlan external subnet")
	})
}

