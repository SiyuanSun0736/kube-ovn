package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestValidateDnat(t *testing.T) {
	c := &Controller{}

	tests := []struct {
		name    string
		dnat    *kubeovnv1.IptablesDnatRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid dnat rule",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: false,
		},
		{
			name: "neither eip nor clusterIP",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "one of eip and clusterIP must be set",
		},
		{
			name: "empty externalPort",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "invalid externalPort",
		},
		{
			name: "invalid externalPort not a number",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "abc",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "invalid externalPort",
		},
		{
			name: "invalid externalPort out of range",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "70000",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "must be between 1 and 65535",
		},
		{
			name: "invalid externalPort zero",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "0",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "must be between 1 and 65535",
		},
		{
			name: "empty internalPort",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "invalid internalPort",
		},
		{
			name: "invalid internalPort not a number",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "xyz",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "invalid internalPort",
		},
		{
			name: "empty internalIP",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "internalIP cannot be empty",
		},
		{
			name: "invalid internalIP",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "not-an-ip",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "invalid internalIP",
		},
		{
			name: "empty protocol",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "",
				},
			},
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name: "invalid protocol",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "icmp",
				},
			},
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name: "mixed-case protocol",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "Tcp",
				},
			},
			wantErr: true,
			errMsg:  "lowercase tcp or udp",
		},
		{
			name: "uppercase protocol",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "80",
					InternalPort: "8080",
					InternalIP:   "10.0.0.1",
					Protocol:     "TCP",
				},
			},
			wantErr: true,
			errMsg:  "lowercase tcp or udp",
		},
		{
			name: "invalid IPv6 internalIP - not supported",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "443",
					InternalPort: "8443",
					InternalIP:   "fd00::1",
					Protocol:     "tcp",
				},
			},
			wantErr: true,
			errMsg:  "must be IPv4",
		},
		{
			name: "max valid port",
			dnat: &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-dnat"},
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					EIP:          "test-eip",
					ExternalPort: "65535",
					InternalPort: "65535",
					InternalIP:   "10.0.0.1",
					Protocol:     "tcp",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateDnatRule(tt.dnat)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFip(t *testing.T) {
	c := &Controller{}

	tests := []struct {
		name    string
		fip     *kubeovnv1.IptablesFIPRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid fip rule",
			fip: &kubeovnv1.IptablesFIPRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-fip"},
				Spec: kubeovnv1.IptablesFIPRuleSpec{
					EIP:        "test-eip",
					InternalIP: "10.0.0.1",
				},
			},
			wantErr: false,
		},
		{
			name: "empty eip",
			fip: &kubeovnv1.IptablesFIPRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-fip"},
				Spec: kubeovnv1.IptablesFIPRuleSpec{
					EIP:        "",
					InternalIP: "10.0.0.1",
				},
			},
			wantErr: true,
			errMsg:  "eip cannot be empty",
		},
		{
			name: "empty internalIP",
			fip: &kubeovnv1.IptablesFIPRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-fip"},
				Spec: kubeovnv1.IptablesFIPRuleSpec{
					EIP:        "test-eip",
					InternalIP: "",
				},
			},
			wantErr: true,
			errMsg:  "internalIP cannot be empty",
		},
		{
			name: "invalid internalIP",
			fip: &kubeovnv1.IptablesFIPRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-fip"},
				Spec: kubeovnv1.IptablesFIPRuleSpec{
					EIP:        "test-eip",
					InternalIP: "invalid-ip",
				},
			},
			wantErr: true,
			errMsg:  "invalid internalIP",
		},
		{
			name: "invalid IPv6 internalIP - not supported",
			fip: &kubeovnv1.IptablesFIPRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-fip"},
				Spec: kubeovnv1.IptablesFIPRuleSpec{
					EIP:        "test-eip",
					InternalIP: "2001:db8::1",
				},
			},
			wantErr: true,
			errMsg:  "must be IPv4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateFipRule(tt.fip)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSnat(t *testing.T) {
	c := &Controller{}

	tests := []struct {
		name    string
		snat    *kubeovnv1.IptablesSnatRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid snat rule",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "10.0.0.0/24",
				},
			},
			wantErr: false,
		},
		{
			name: "empty eip",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "",
					InternalCIDR: "10.0.0.0/24",
				},
			},
			wantErr: true,
			errMsg:  "eip cannot be empty",
		},
		{
			name: "empty internalCIDR",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "",
				},
			},
			wantErr: true,
			errMsg:  "internalCIDR cannot be empty",
		},
		{
			name: "valid single IP address",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "10.0.0.1",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid single IPv6 address - not supported",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "fd00::1",
				},
			},
			wantErr: true,
			errMsg:  "must be IPv4",
		},
		{
			name: "invalid internalCIDR - malformed IP",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "10.0.0.256",
				},
			},
			wantErr: true,
			errMsg:  "invalid internalCIDR",
		},
		{
			name: "invalid internalCIDR - invalid format",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "invalid-cidr",
				},
			},
			wantErr: true,
			errMsg:  "invalid internalCIDR",
		},
		{
			name: "invalid IPv6 internalCIDR - not supported",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "fd00::/64",
				},
			},
			wantErr: true,
			errMsg:  "must be IPv4",
		},
		{
			name: "invalid multiple CIDRs - not supported",
			snat: &kubeovnv1.IptablesSnatRule{
				ObjectMeta: metav1.ObjectMeta{Name: "test-snat"},
				Spec: kubeovnv1.IptablesSnatRuleSpec{
					EIP:          "test-eip",
					InternalCIDR: "10.0.0.0/24,192.168.1.0/24",
				},
			},
			wantErr: true,
			errMsg:  "contains multiple CIDRs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateSnatRule(tt.snat)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDeleteFipInPod_NatGwGone verifies that deleteFipInPod returns nil (skips
// cleanup) when the VpcNatGateway CRD no longer exists.
func TestDeleteFipInPod_NatGwGone(t *testing.T) {
	t.Parallel()
	fc, err := newFakeControllerWithOptions(t, nil)
	require.NoError(t, err)
	err = fc.fakeController.deleteFipInPod("missing-gw", "10.0.0.1")
	require.NoError(t, err, "should skip cleanup when gateway CRD is gone")
}

// TestDeleteFipInPod_NatGwExistsPodMissing verifies that deleteFipInPod returns
// an error to trigger a retry when the gateway CRD exists but the pod is absent.
func TestDeleteFipInPod_NatGwExistsPodMissing(t *testing.T) {
	t.Parallel()
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("test-gw")},
	})
	require.NoError(t, err)
	err = fc.fakeController.deleteFipInPod("test-gw", "10.0.0.1")
	require.NoError(t, err, "a gateway without a running instance holds no data plane to clean up")
}

// TestDeleteDnatInPod_NatGwGone verifies that deleteDnatInPod returns nil when
// the VpcNatGateway CRD no longer exists.
func TestDeleteDnatInPod_NatGwGone(t *testing.T) {
	t.Parallel()
	fc, err := newFakeControllerWithOptions(t, nil)
	require.NoError(t, err)
	err = fc.fakeController.deleteDnatInPod("missing-gw", "tcp", "10.0.0.1", "80")
	require.NoError(t, err, "should skip cleanup when gateway CRD is gone")
}

// TestDeleteDnatInPod_NatGwExistsPodMissing verifies that deleteDnatInPod
// returns an error to trigger a retry when the pod is absent.
func TestDeleteDnatInPod_NatGwExistsPodMissing(t *testing.T) {
	t.Parallel()
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("test-gw")},
	})
	require.NoError(t, err)
	err = fc.fakeController.deleteDnatInPod("test-gw", "tcp", "10.0.0.1", "80")
	require.NoError(t, err, "a gateway without a running instance holds no data plane to clean up")
}

// TestDeleteSnatInPod_NatGwGone verifies that deleteSnatInPod returns nil when
// the VpcNatGateway CRD no longer exists.
func TestDeleteSnatInPod_NatGwGone(t *testing.T) {
	t.Parallel()
	fc, err := newFakeControllerWithOptions(t, nil)
	require.NoError(t, err)
	err = fc.fakeController.deleteSnatInPod("missing-gw", "10.0.0.1", "192.168.1.0/24")
	require.NoError(t, err, "should skip cleanup when gateway CRD is gone")
}

// TestDeleteSnatInPod_NatGwExistsPodMissing verifies that deleteSnatInPod
// returns an error to trigger a retry when the pod is absent.
func TestDeleteSnatInPod_NatGwExistsPodMissing(t *testing.T) {
	t.Parallel()
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{fakeGw("test-gw")},
	})
	require.NoError(t, err)
	err = fc.fakeController.deleteSnatInPod("test-gw", "10.0.0.1", "192.168.1.0/24")
	require.NoError(t, err, "a gateway without a running instance holds no data plane to clean up")
}

// shareDnat builds an IptablesDnatRule with the identity labels the share DNAT consumers
// select on. dnatType may be "" (defaults to exclusive behavior).
func shareDnat(name, gw, eip, eport, proto, intIP, intPort, dnatType string) *kubeovnv1.IptablesDnatRule {
	return &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				util.VpcNatGatewayNameLabel: gw,
				util.VpcDnatEPortLabel:      eport,
			},
		},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP:          eip,
			ExternalPort: eport,
			Protocol:     proto,
			InternalIP:   intIP,
			InternalPort: intPort,
			Type:         dnatType,
		},
	}
}

func TestDedupSortedBackends(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty entries dropped", in: []string{"", ""}, want: nil},
		{
			name: "dedup and sort",
			in:   []string{"10.0.0.2:80", "10.0.0.1:80", "10.0.0.2:80", "", "10.0.0.3:80"},
			want: []string{"10.0.0.1:80", "10.0.0.2:80", "10.0.0.3:80"},
		},
		{
			name: "already unique stays sorted",
			in:   []string{"10.0.0.1:80", "10.0.0.2:80"},
			want: []string{"10.0.0.1:80", "10.0.0.2:80"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, dedupSortedBackends(tt.in))
		})
	}
}

func TestDnatCleanupEipName(t *testing.T) {
	t.Parallel()

	dnat := shareDnat("dnat", "gw", "new-eip", "80", "tcp", "10.0.0.1", "8080", kubeovnv1.DnatRuleTypeShare)
	dnat.Annotations = map[string]string{util.VpcEipAnnotation: "old-eip"}
	assert.Equal(t, "old-eip", dnatCleanupEipName(dnat), "rebind cleanup must use the old EIP sibling set")

	dnat.Annotations = nil
	assert.Equal(t, "new-eip", dnatCleanupEipName(dnat), "legacy rules fall back to the current spec EIP")
}

// assertEnqueueAddRouting checks that the add handler routes a live object to the add queue and a
// terminating object to the update queue (the deletion-cleanup path relied on after a restart).
func assertEnqueueAddRouting(
	t *testing.T,
	addQueue, updateQueue workqueue.TypedRateLimitingInterface[string],
	enqueue func(any),
	live, terminating any,
) {
	t.Helper()
	enqueue(live)
	require.Equal(t, 1, addQueue.Len(), "live object should go to the add queue")
	require.Equal(t, 0, updateQueue.Len(), "live object must not go to the update queue")

	enqueue(terminating)
	require.Equal(t, 1, addQueue.Len(), "terminating object must not go to the add queue")
	require.Equal(t, 1, updateQueue.Len(), "terminating object should go to the update queue for cleanup")
}

func TestEnqueueAddIptablesFip(t *testing.T) {
	t.Parallel()
	c := &Controller{
		addIptablesFipQueue:    newTypedRateLimitingQueue[string]("AddIptablesFip", nil),
		updateIptablesFipQueue: newTypedRateLimitingQueue[string]("UpdateIptablesFip", nil),
	}
	t.Cleanup(c.addIptablesFipQueue.ShutDown)
	t.Cleanup(c.updateIptablesFipQueue.ShutDown)
	now := metav1.Now()
	assertEnqueueAddRouting(
		t, c.addIptablesFipQueue, c.updateIptablesFipQueue, c.enqueueAddIptablesFip,
		&kubeovnv1.IptablesFIPRule{ObjectMeta: metav1.ObjectMeta{Name: "live-fip"}},
		&kubeovnv1.IptablesFIPRule{ObjectMeta: metav1.ObjectMeta{Name: "terminating-fip", DeletionTimestamp: &now}},
	)
}

func TestEnqueueAddIptablesDnatRule(t *testing.T) {
	t.Parallel()
	c := &Controller{
		addIptablesDnatRuleQueue:    newTypedRateLimitingQueue[string]("AddIptablesDnat", nil),
		updateIptablesDnatRuleQueue: newTypedRateLimitingQueue[string]("UpdateIptablesDnat", nil),
	}
	t.Cleanup(c.addIptablesDnatRuleQueue.ShutDown)
	t.Cleanup(c.updateIptablesDnatRuleQueue.ShutDown)
	now := metav1.Now()
	assertEnqueueAddRouting(
		t, c.addIptablesDnatRuleQueue, c.updateIptablesDnatRuleQueue, c.enqueueAddIptablesDnatRule,
		&kubeovnv1.IptablesDnatRule{ObjectMeta: metav1.ObjectMeta{Name: "live-dnat"}},
		&kubeovnv1.IptablesDnatRule{ObjectMeta: metav1.ObjectMeta{Name: "terminating-dnat", DeletionTimestamp: &now}},
	)
}

func TestEnqueueUpdateIptablesDnatRuleRejectsUppercaseProtocol(t *testing.T) {
	t.Parallel()
	queue := newTypedRateLimitingQueue[string]("UpdateIptablesDnat", nil)
	t.Cleanup(queue.ShutDown)
	c := &Controller{updateIptablesDnatRuleQueue: queue}
	oldDnat := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{Name: "dnat"},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP: "eip", ExternalPort: "80", InternalPort: "8080", InternalIP: "10.0.0.1", Protocol: "TCP",
		},
	}
	newDnat := oldDnat.DeepCopy()

	c.enqueueUpdateIptablesDnatRule(oldDnat, newDnat)
	item, shutdown := queue.Get()
	require.False(t, shutdown)
	require.Equal(t, "dnat", item)
	queue.Done(item)
}

func TestEnqueueAddIptablesSnatRule(t *testing.T) {
	t.Parallel()
	c := &Controller{
		addIptablesSnatRuleQueue:    newTypedRateLimitingQueue[string]("AddIptablesSnat", nil),
		updateIptablesSnatRuleQueue: newTypedRateLimitingQueue[string]("UpdateIptablesSnat", nil),
	}
	t.Cleanup(c.addIptablesSnatRuleQueue.ShutDown)
	t.Cleanup(c.updateIptablesSnatRuleQueue.ShutDown)
	now := metav1.Now()
	assertEnqueueAddRouting(
		t, c.addIptablesSnatRuleQueue, c.updateIptablesSnatRuleQueue, c.enqueueAddIptablesSnatRule,
		&kubeovnv1.IptablesSnatRule{ObjectMeta: metav1.ObjectMeta{Name: "live-snat"}},
		&kubeovnv1.IptablesSnatRule{ObjectMeta: metav1.ObjectMeta{Name: "terminating-snat", DeletionTimestamp: &now}},
	)
}

func TestNormalizeSnatInternalCIDR(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", normalizeSnatInternalCIDR(""))
	assert.Equal(t, "10.0.0.1/32", normalizeSnatInternalCIDR("10.0.0.1"))
	assert.Equal(t, "10.0.0.0/24", normalizeSnatInternalCIDR("10.0.0.0/24"))
}

func TestResolveSnatMemberID(t *testing.T) {
	t.Parallel()
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{util.NatGatewayMemberLabel: "member-a"},
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{util.NatGatewayMemberLabel: "member-b"},
		},
	}
	// EIP takes precedence
	assert.Equal(t, "member-a", resolveSnatMemberID(eip, snat))
	// Fallback to snat
	assert.Equal(t, "member-b", resolveSnatMemberID(nil, snat))
	// Unsharded EIP takes precedence over snat
	unshardedEip := &kubeovnv1.IptablesEIP{}
	assert.Equal(t, "", resolveSnatMemberID(unshardedEip, snat))
	// Legacy label
	legacySnat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{util.NatGatewayMemberLegacyLabel: "member-legacy"},
		},
	}
	assert.Equal(t, "member-legacy", resolveSnatMemberID(nil, legacySnat))
}

func TestFilterNatGwPodsByMember(t *testing.T) {
	t.Parallel()
	c := &Controller{}
	pods := []*corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "gw-0",
				Labels: map[string]string{util.NatGatewayMemberLabel: "member-a"},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "gw-1",
				Labels: map[string]string{util.NatGatewayMemberLabel: "member-b"},
			},
		},
	}

	// Empty memberID matches all
	matched, err := c.filterNatGwPodsByMember(pods, "")
	require.NoError(t, err)
	assert.Len(t, matched, 2)

	// Matching member-a
	matched, err = c.filterNatGwPodsByMember(pods, "member-a")
	require.NoError(t, err)
	require.Len(t, matched, 1)
	assert.Equal(t, "gw-0", matched[0].Name)

	// Non-existing member returns error
	_, err = c.filterNatGwPodsByMember(pods, "non-existent")
	assert.Error(t, err)
}

func TestResolveRecordedSnatMemberID(t *testing.T) {
	t.Parallel()
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{util.NatGatewayMemberLabel: "member-new"},
		},
	}
	snat := &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{util.NatGatewayMemberLabel: "member-old"},
		},
	}
	// For recorded state cleanup, SNAT recorded member takes precedence over current EIP member
	assert.Equal(t, "member-old", resolveRecordedSnatMemberID(snat, eip))
	// Fallback to EIP if SNAT has no member metadata
	assert.Equal(t, "member-new", resolveRecordedSnatMemberID(nil, eip))
	// Unsharded returns empty
	assert.Equal(t, "", resolveRecordedSnatMemberID(nil, nil))
}
