package controller

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/ptr"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

type orderedVpcNatGatewayLister struct {
	gateways []*kubeovnv1.VpcNatGateway
}

func (l orderedVpcNatGatewayLister) List(labels.Selector) ([]*kubeovnv1.VpcNatGateway, error) {
	return l.gateways, nil
}

func (l orderedVpcNatGatewayLister) Get(name string) (*kubeovnv1.VpcNatGateway, error) {
	for _, gateway := range l.gateways {
		if gateway.Name == name {
			return gateway, nil
		}
	}
	return nil, nil
}

// fakeNatGwPod builds the running NAT gateway pod of a gateway, optionally already initialized.
func fakeNatGwPod(gwName string, initialized bool) *corev1.Pod {
	annotations := map[string]string{util.VpcNatGatewayAnnotation: gwName}
	if initialized {
		annotations[util.VpcNatGatewayInitAnnotation] = "true"
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GenNatGwName(gwName) + "-0",
			Namespace: "kube-system",
			Labels: map[string]string{
				"app":                   util.GenNatGwName(gwName),
				util.VpcNatGatewayLabel: "true",
			},
			Annotations: annotations,
		},
		Status: corev1.PodStatus{Phase: corev1.PodRunning},
	}
}

func TestVpcNatGwScriptConstants(t *testing.T) {
	// Verify the constants are set correctly for backward compatibility
	assert.Equal(t, "/kube-ovn", vpcNatGwScriptMountPath, "Script mount path should be /kube-ovn")
	assert.Equal(t, "nat-gw-script", vpcNatGwScriptVolumeName, "Volume name should be nat-gw-script")
	assert.Equal(t, "nat-gateway.sh", vpcNatGwScriptName, "Script name should be nat-gateway.sh")
	assert.Equal(t, "/kube-ovn/nat-gateway.sh", vpcNatGwScriptPath, "Script path should be /kube-ovn/nat-gateway.sh")
	assert.Equal(t, "vpc-nat-gw", vpcNatGwContainerName, "Container name should be vpc-nat-gw")
	assert.Equal(t, "vpc-nat-gw", vpcNatGwServiceAccountName, "ServiceAccount name should be vpc-nat-gw")
}

func TestEnqueueVpcNatGwRoutesTerminatingToDelete(t *testing.T) {
	t.Parallel()

	updateQueue := newTypedRateLimitingQueue[string]("VpcNatGateway", nil)
	deleteQueue := newTypedRateLimitingQueue[string]("DeleteVpcNatGateway", nil)
	t.Cleanup(updateQueue.ShutDown)
	t.Cleanup(deleteQueue.ShutDown)
	c := &Controller{
		addOrUpdateVpcNatGatewayQueue: updateQueue,
		delVpcNatGatewayQueue:         deleteQueue,
		config:                        &Configuration{PodNamespace: "kube-system"},
	}
	now := metav1.Now()
	liveGw := fakeGw("live-gw")
	dyingAddGw := fakeGw("dying-add-gw")
	dyingAddGw.DeletionTimestamp = &now
	dyingAddGw.Finalizers = []string{util.KubeOVNControllerFinalizer}
	dyingUpdateGw := fakeGw("dying-update-gw")
	dyingUpdateGw.DeletionTimestamp = &now
	dyingUpdateGw.Finalizers = []string{util.KubeOVNControllerFinalizer}

	c.enqueueAddVpcNatGw(liveGw)
	require.Equal(t, 1, updateQueue.Len())
	c.enqueueAddVpcNatGw(dyingAddGw)
	require.Equal(t, 1, updateQueue.Len())
	require.Equal(t, 1, deleteQueue.Len())
	c.enqueueUpdateVpcNatGw(liveGw, dyingUpdateGw)
	require.Equal(t, 1, updateQueue.Len())
	require.Equal(t, 2, deleteQueue.Len())
}

func TestHandleAddOrUpdateVpcNatGwSkipsTerminating(t *testing.T) {
	t.Parallel()

	now := metav1.Now()
	gw := fakeGw("dying-gw")
	gw.DeletionTimestamp = &now
	gw.Finalizers = []string{util.KubeOVNControllerFinalizer}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)
	require.NoError(t, fc.fakeController.handleAddOrUpdateVpcNatGw("dying-gw"))
}

func TestVpcNatGatewayRejectsStaleQoSGenerationAfterReferenceIsDropped(t *testing.T) {
	qos := &kubeovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{
		Name: "qos", UID: types.UID("new-uid"), Finalizers: []string{util.KubeOVNControllerFinalizer},
	}}
	gw := fakeGw("gw")
	gw.Labels = map[string]string{util.QoSPolicyUIDLabel: "old-uid"}
	gw.Status.QoSPolicy = "qos"
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)

	err = fc.fakeController.handleAddOrUpdateVpcNatGw("gw")
	require.ErrorContains(t, err, "previous generation of qos policy qos")
}

func TestVpcNatGatewayRejectsQoSWithoutControllerFinalizer(t *testing.T) {
	qos := &kubeovnv1.QoSPolicy{ObjectMeta: metav1.ObjectMeta{Name: "qos", UID: types.UID("qos-uid")}}
	gw := fakeGw("gw")
	gw.Labels = map[string]string{util.QoSPolicyUIDLabel: "qos-uid"}
	gw.Status.QoSPolicy = "qos"
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)

	err = fc.fakeController.handleAddOrUpdateVpcNatGw("gw")
	require.ErrorContains(t, err, "waiting for qos policy qos controller reconcile")
}

// TestHandleInitVpcNatGwSkipsTerminating pins the same guard on the init path, which pod update
// events keep feeding while the gateway pod is still terminating.
func TestHandleInitVpcNatGwSkipsTerminating(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	now := metav1.Now()
	gw := fakeGw("dying-gw")
	gw.DeletionTimestamp = &now
	gw.Finalizers = []string{util.KubeOVNControllerFinalizer}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)
	// Without the guard this reaches getNatGwPod and fails on the already gone pod.
	require.NoError(t, fc.fakeController.handleInitVpcNatGw("dying-gw"))
}

func TestInitVpcNatGwContinuesAfterInitializedGateway(t *testing.T) {
	old := vpcNatEnabled
	vpcNatEnabled = "true"
	t.Cleanup(func() { vpcNatEnabled = old })

	initializedGw := fakeGw("initialized-gw")
	pendingGw := fakeGw("pending-gw")
	pods := []*corev1.Pod{
		fakeNatGwPod(initializedGw.Name, true),
		fakeNatGwPod(pendingGw.Name, false),
	}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{Pods: pods})
	require.NoError(t, err)
	fc.fakeController.vpcNatGatewayLister = orderedVpcNatGatewayLister{
		gateways: []*kubeovnv1.VpcNatGateway{initializedGw, pendingGw},
	}
	queue := newTypedRateLimitingQueue[string]("InitVpcNatGateway", nil)
	t.Cleanup(queue.ShutDown)
	fc.fakeController.initVpcNatGatewayQueue = queue

	require.NoError(t, fc.fakeController.initVpcNatGw())
	require.Equal(t, 1, queue.Len())
	item, shutdown := queue.Get()
	require.False(t, shutdown)
	queue.Done(item)
	require.Equal(t, pendingGw.Name, item)
}

// TestResyncVpcNatGwConfigDrivesInitScan covers the startup race between the enable transition and
// the init scan: vpcNatEnabled is still "unknown" while initResourceOnce runs, so the scan must be
// driven by the resync that flips the flag, otherwise the restart silently skips every gateway.
func TestResyncVpcNatGwConfigDrivesInitScan(t *testing.T) {
	oldEnabled, oldVersion := vpcNatEnabled, VpcNatCmVersion
	vpcNatEnabled = "unknown"
	VpcNatCmVersion = ""
	t.Cleanup(func() {
		vpcNatEnabled = oldEnabled
		VpcNatCmVersion = oldVersion
	})

	initializedGw := fakeGw("initialized-gw")
	pendingGw := fakeGw("pending-gw")
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{initializedGw, pendingGw},
		Pods: []*corev1.Pod{
			fakeNatGwPod(initializedGw.Name, true),
			fakeNatGwPod(pendingGw.Name, false),
		},
		ConfigMaps: []*corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{Name: util.VpcNatGatewayConfig, Namespace: "kube-system"},
				Data:       map[string]string{"enable-vpc-nat-gw": "true"},
			},
		},
	})
	require.NoError(t, err)

	addOrUpdateQueue := newTypedRateLimitingQueue[string]("AddOrUpdateVpcNatGateway", nil)
	initQueue := newTypedRateLimitingQueue[string]("InitVpcNatGateway", nil)
	t.Cleanup(addOrUpdateQueue.ShutDown)
	t.Cleanup(initQueue.ShutDown)
	fc.fakeController.addOrUpdateVpcNatGatewayQueue = addOrUpdateQueue
	fc.fakeController.initVpcNatGatewayQueue = initQueue

	fc.fakeController.resyncVpcNatGwConfig()

	require.Equal(t, "true", vpcNatEnabled, "the resync owns the enable transition")
	require.Equal(t, 2, addOrUpdateQueue.Len())
	require.Equal(t, 1, initQueue.Len(), "only the uninitialized gateway needs a scan enqueue")
	item, shutdown := initQueue.Get()
	require.False(t, shutdown)
	initQueue.Done(item)
	require.Equal(t, pendingGw.Name, item)
}

func TestExecNatGwQoSWaitsForQoSPolicyStatus(t *testing.T) {
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "pending-qos"},
		Spec: kubeovnv1.QoSPolicySpec{
			Shared:      true,
			BindingType: kubeovnv1.QoSBindingTypeNatGw,
		},
	}
	gw := fakeGw("gw")
	gw.Spec.QoSPolicy = "pending-qos"
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies:    []*kubeovnv1.QoSPolicy{qos},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)

	err = fc.fakeController.execNatGwQoS(gw, "pending-qos", QoSAdd)
	require.ErrorContains(t, err, "status to match the spec")
}

func TestExecNatGwQoSAllowsCleanupOfTerminatingPolicy(t *testing.T) {
	now := metav1.Now()
	qos := &kubeovnv1.QoSPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "dying-qos", DeletionTimestamp: &now},
		Spec:       kubeovnv1.QoSPolicySpec{Shared: true, BindingType: kubeovnv1.QoSBindingTypeNatGw},
		Status:     kubeovnv1.QoSPolicyStatus{Shared: true, BindingType: kubeovnv1.QoSBindingTypeNatGw},
	}
	gw := fakeGw("gw")
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		QoSPolicies: []*kubeovnv1.QoSPolicy{qos},
	})
	require.NoError(t, err)

	// No bandwidth rules means no pod lookup is needed; this isolates the cleanup guard.
	require.NoError(t, fc.fakeController.execNatGwQoS(gw, "dying-qos", QoSDel))
}

func TestIsVpcNatGwChanged(t *testing.T) {
	tests := []struct {
		name     string
		gw       *kubeovnv1.VpcNatGateway
		expected bool
	}{
		{
			name: "no changes returns false",
			gw: &kubeovnv1.VpcNatGateway{
				Spec: kubeovnv1.VpcNatGatewaySpec{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
				},
				Status: kubeovnv1.VpcNatGatewayStatus{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
				},
			},
			expected: false,
		},
		{
			name: "ExternalSubnets changed returns true",
			gw: &kubeovnv1.VpcNatGateway{
				Spec: kubeovnv1.VpcNatGatewaySpec{
					ExternalSubnets: []string{"subnet2"},
					Selector:        []string{"node=worker1"},
				},
				Status: kubeovnv1.VpcNatGatewayStatus{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
				},
			},
			expected: true,
		},
		{
			name: "Selector changed returns true",
			gw: &kubeovnv1.VpcNatGateway{
				Spec: kubeovnv1.VpcNatGatewaySpec{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker2"},
				},
				Status: kubeovnv1.VpcNatGatewayStatus{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
				},
			},
			expected: true,
		},
		{
			name: "Tolerations changed returns true",
			gw: &kubeovnv1.VpcNatGateway{
				Spec: kubeovnv1.VpcNatGatewaySpec{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
					Tolerations:     []corev1.Toleration{{Key: "new-key"}},
				},
				Status: kubeovnv1.VpcNatGatewayStatus{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
				},
			},
			expected: true,
		},
		{
			name: "Affinity changed returns true",
			gw: &kubeovnv1.VpcNatGateway{
				Spec: kubeovnv1.VpcNatGatewaySpec{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
					Affinity: corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{},
					},
				},
				Status: kubeovnv1.VpcNatGatewayStatus{
					ExternalSubnets: []string{"subnet1"},
					Selector:        []string{"node=worker1"},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVpcNatGwChanged(tt.gw)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSubnetProvider(t *testing.T) {
	tests := []struct {
		name             string
		subnetName       string
		subnets          []*kubeovnv1.Subnet
		expectedProvider string
		expectError      bool
		description      string
	}{
		{
			name:       "Valid OVN subnets with different providers",
			subnetName: "ovn-default",
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "ovn-default",
					},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.244.0.0/24",
						Provider:  util.OvnProvider,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "net1-subnet",
					},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "192.168.1.0/24",
						Provider:  "net1.default.ovn",
					},
				},
			},
			expectedProvider: util.OvnProvider,
			expectError:      false,
			description:      "Should return correct provider for valid OVN subnet among multiple subnets",
		},
		{
			name:       "Non-existent subnet",
			subnetName: "non-existent",
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "ovn-default",
					},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.244.0.0/24",
						Provider:  util.OvnProvider,
					},
				},
			},
			expectedProvider: "",
			expectError:      true,
			description:      "Should return error for non-existent subnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create controller with subnets
			fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				Subnets: tt.subnets,
			})
			require.NoError(t, err, "Failed to create fake controller")
			controller := fakeController.fakeController
			// Call the method under test
			provider, err := controller.GetSubnetProvider(tt.subnetName)

			// Check for errors
			if tt.expectError {
				assert.Error(t, err, "Expected an error but got none: %s", tt.description)
				return
			}
			require.NoError(t, err, "Unexpected error: %s", tt.description)

			// Verify provider
			assert.Equal(t, tt.expectedProvider, provider, "Provider mismatch: %s", tt.description)
		})
	}

	// Test multiple provider scenarios in a single comprehensive test
	t.Run("Multiple provider scenarios", func(t *testing.T) {
		subnets := []*kubeovnv1.Subnet{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "default-subnet"},
				Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "custom-subnet"},
				Spec:       kubeovnv1.SubnetSpec{Provider: "custom.provider.ovn"},
			},
		}

		fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			Subnets: subnets,
		})
		require.NoError(t, err)
		controller := fakeController.fakeController
		// Test default provider
		provider, err := controller.GetSubnetProvider("default-subnet")
		require.NoError(t, err)
		assert.Equal(t, util.OvnProvider, provider)

		// Test custom provider
		provider, err = controller.GetSubnetProvider("custom-subnet")
		require.NoError(t, err)
		assert.Equal(t, "custom.provider.ovn", provider)

		// Test non-existent subnet
		_, err = controller.GetSubnetProvider("missing-subnet")
		assert.Error(t, err, "Should error for missing subnet")
	})
}

func TestGetExternalSubnetNad(t *testing.T) {
	tests := []struct {
		name              string
		gw                *kubeovnv1.VpcNatGateway
		subnets           []*kubeovnv1.Subnet
		podNamespace      string
		expectedNamespace string
		expectedName      string
		expectError       bool
	}{
		{
			name: "provider with 3 parts (name.namespace.ovn)",
			gw: &kubeovnv1.VpcNatGateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gw"},
				Spec:       kubeovnv1.VpcNatGatewaySpec{ExternalSubnets: []string{"external-subnet"}},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "external-subnet"},
					Spec:       kubeovnv1.SubnetSpec{Provider: "real-eip.kube-system.ovn"},
				},
			},
			podNamespace:      "kube-system",
			expectedNamespace: "kube-system",
			expectedName:      "real-eip",
			expectError:       false,
		},
		{
			name: "provider with 2 parts (name.namespace)",
			gw: &kubeovnv1.VpcNatGateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gw"},
				Spec:       kubeovnv1.VpcNatGatewaySpec{ExternalSubnets: []string{"external-subnet"}},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "external-subnet"},
					Spec:       kubeovnv1.SubnetSpec{Provider: "my-nad.default"},
				},
			},
			podNamespace:      "kube-system",
			expectedNamespace: "default",
			expectedName:      "my-nad",
			expectError:       false,
		},
		{
			name: "provider is ovn (fallback to subnet name)",
			gw: &kubeovnv1.VpcNatGateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gw"},
				Spec:       kubeovnv1.VpcNatGatewaySpec{ExternalSubnets: []string{"ovn-vpc-external-network"}},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "ovn-vpc-external-network"},
					Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider},
				},
			},
			podNamespace:      "kube-system",
			expectedNamespace: "kube-system",
			expectedName:      "ovn-vpc-external-network",
			expectError:       false,
		},
		{
			name: "empty provider (fallback to subnet name)",
			gw: &kubeovnv1.VpcNatGateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gw"},
				Spec:       kubeovnv1.VpcNatGatewaySpec{ExternalSubnets: []string{"my-external-subnet"}},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "my-external-subnet"},
					Spec:       kubeovnv1.SubnetSpec{Provider: ""},
				},
			},
			podNamespace:      "kube-system",
			expectedNamespace: "kube-system",
			expectedName:      "my-external-subnet",
			expectError:       false,
		},
		{
			name: "empty ExternalSubnets (use default)",
			gw: &kubeovnv1.VpcNatGateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gw"},
				Spec:       kubeovnv1.VpcNatGatewaySpec{ExternalSubnets: []string{}},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "ovn-vpc-external-network"},
					Spec:       kubeovnv1.SubnetSpec{Provider: "external.default.ovn"},
				},
			},
			podNamespace:      "kube-system",
			expectedNamespace: "default",
			expectedName:      "external",
			expectError:       false,
		},
		{
			name: "subnet not found",
			gw: &kubeovnv1.VpcNatGateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gw"},
				Spec:       kubeovnv1.VpcNatGatewaySpec{ExternalSubnets: []string{"non-existent-subnet"}},
			},
			subnets:      []*kubeovnv1.Subnet{},
			podNamespace: "kube-system",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				Subnets: tt.subnets,
			})
			require.NoError(t, err)
			controller := fakeController.fakeController
			controller.config.PodNamespace = tt.podNamespace

			namespace, name, err := controller.getExternalSubnetNad(tt.gw)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedNamespace, namespace, "namespace mismatch")
			assert.Equal(t, tt.expectedName, name, "name mismatch")
		})
	}
}

// Hairpin SNAT is hot-applied: on a Pod whose state init already recorded, a spec change is
// pushed and recorded, and the recorded state then stops the gateway from reprogramming the same
// state on every resync.
func TestSyncNatGwHairpinSNAT(t *testing.T) {
	gw := fakeGw("gw")
	gw.Status.HairpinSNAT = hairpinSNATEnable // what handleInitVpcNatGw records on a fresh Pod
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)
	ctrl := fc.fakeController

	var calls [][]string
	ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
		calls = append(calls, append([]string{operation}, rules...))
		return nil
	}

	// converged: the spec is unset, which is the default the status already records
	require.NoError(t, ctrl.syncNatGwHairpinSNAT(gw, fakeNatGwPod("gw", true)))
	require.Empty(t, calls)

	// opting out programs it and records it
	gw.Spec.EnableLanIPAsHairpinSNAT = ptr.To(false)
	require.NoError(t, ctrl.syncNatGwHairpinSNAT(gw, fakeNatGwPod("gw", true)))
	require.Equal(t, [][]string{{natGwHairpinSNATSet, "disable"}}, calls)
	patched, err := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "disable", patched.Status.HairpinSNAT)

	// already programmed: nothing is pushed again (the stored object carries the spec the API has,
	// so state it explicitly rather than relying on what the fake client kept)
	patched.Spec.EnableLanIPAsHairpinSNAT = ptr.To(false)
	require.NoError(t, ctrl.syncNatGwHairpinSNAT(patched, fakeNatGwPod("gw", true)))
	require.Len(t, calls, 1)

	// and switching it back on is pushed too
	patched.Spec.EnableLanIPAsHairpinSNAT = ptr.To(true)
	require.NoError(t, ctrl.syncNatGwHairpinSNAT(patched, fakeNatGwPod("gw", true)))
	require.Equal(t, []string{natGwHairpinSNATSet, "enable"}, calls[1])

	// clearing the field goes back to the default, which is on
	patched.Spec.EnableLanIPAsHairpinSNAT = nil
	patched.Status.HairpinSNAT = hairpinSNATDisable
	require.NoError(t, ctrl.syncNatGwHairpinSNAT(patched, fakeNatGwPod("gw", true)))
	require.Equal(t, []string{natGwHairpinSNATSet, "enable"}, calls[2])
}

// The fallback is on unless the gateway opts out.
// The default is on, matching kube-proxy's Cluster traffic policy: a backend that programs no
// route of its own still gets its replies back through the gateway. Only an explicit false opts
// out, which is why the field is a pointer - a plain bool could not tell "unset" from "off".
func TestHairpinSNATModeDefaultsToOn(t *testing.T) {
	t.Parallel()

	// unset: the CRD default fills this in as true, and a cluster whose CRD predates the field
	// must land on the same behavior rather than silently inverting it
	require.Equal(t, "enable", hairpinSNATMode(fakeGw("unset")))

	gw := fakeGw("on")
	gw.Spec.EnableLanIPAsHairpinSNAT = ptr.To(true)
	require.Equal(t, "enable", hairpinSNATMode(gw))

	off := fakeGw("off")
	off.Spec.EnableLanIPAsHairpinSNAT = ptr.To(false)
	require.Equal(t, "disable", hairpinSNATMode(off))
}

// The init path is the only place a replaced Pod gets its hairpin fallback back: it programs the
// state the spec asks for and records it, so the status is never a claim the data plane does not
// back.
func TestHandleInitVpcNatGwAppliesHairpinSNAT(t *testing.T) {
	for _, tt := range []struct {
		name     string
		enabled  bool
		wantOps  [][]string
		wantStat string
	}{
		{
			name: "default gateway records the state it programmed",
			wantOps: [][]string{
				{natGwInit},
				{natGwHairpinSNATSet, hairpinSNATDisable},
				{natGwTCPConnLimitSet, "50000", "stateful"},
				{natGwTCPSynLimitSet, util.DefaultNatGwTCPSynSrcRate, util.DefaultNatGwTCPSynSrcBurst, util.DefaultNatGwTCPSynDstRate, util.DefaultNatGwTCPSynDstBurst},
			},
			wantStat: hairpinSNATDisable,
		},
		{
			name:    "fallback requested is pushed after init",
			enabled: true,
			wantOps: [][]string{
				{natGwInit},
				{natGwHairpinSNATSet, hairpinSNATEnable},
				{natGwTCPConnLimitSet, "50000", "stateful"},
				{natGwTCPSynLimitSet, util.DefaultNatGwTCPSynSrcRate, util.DefaultNatGwTCPSynSrcBurst, util.DefaultNatGwTCPSynDstRate, util.DefaultNatGwTCPSynDstBurst},
			},
			wantStat: hairpinSNATEnable,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			old := vpcNatEnabled
			vpcNatEnabled = "true"
			t.Cleanup(func() { vpcNatEnabled = old })

			gw := fakeGw("gw")
			gw.Spec.EnableLanIPAsHairpinSNAT = ptr.To(tt.enabled)
			pod := fakeNatGwPod("gw", false)
			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
				Pods:           []*corev1.Pod{pod},
			})
			require.NoError(t, err)
			ctrl := fc.fakeController
			for _, q := range []*workqueue.TypedRateLimitingInterface[string]{
				&ctrl.updateVpcFloatingIPQueue, &ctrl.updateVpcDnatQueue, &ctrl.updateVpcSnatQueue,
				&ctrl.updateVpcSubnetQueue, &ctrl.updateVpcEipQueue,
			} {
				*q = newTypedRateLimitingQueue[string]("test", nil)
				t.Cleanup((*q).ShutDown)
			}

			var calls [][]string
			ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
				calls = append(calls, append([]string{operation}, rules...))
				return nil
			}

			require.NoError(t, ctrl.handleInitVpcNatGw("gw"))
			require.Equal(t, tt.wantOps, calls)

			// the programmed fallback is recorded, so the update path does not push it again
			patched, err := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "gw", metav1.GetOptions{})
			require.NoError(t, err)
			require.Equal(t, tt.wantStat, patched.Status.HairpinSNAT)
			require.NoError(t, ctrl.syncNatGwHairpinSNAT(patched, fakeNatGwPod("gw", true)))
			require.Len(t, calls, len(tt.wantOps))

			// and the Pod is marked initialized only after all of it succeeded
			initialized, err := ctrl.config.KubeClient.CoreV1().Pods(pod.Namespace).Get(t.Context(), pod.Name, metav1.GetOptions{})
			require.NoError(t, err)
			require.Equal(t, "true", initialized.Annotations[util.VpcNatGatewayInitAnnotation])
		})
	}
}

// The hairpin status writer must not carry the rest of the status with it: its lister snapshot
// can predate what another path has already written, and a whole-status merge patch would put the
// stale values back.
func TestPatchNatGwHairpinSNATStatusTouchesOnlyItsField(t *testing.T) {
	t.Parallel()

	gw := fakeGw("gw")
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
	})
	require.NoError(t, err)
	ctrl := fc.fakeController

	// another path moves the QoS status forward; the lister still holds the snapshot without it
	gwClient := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways()
	_, err = gwClient.Patch(t.Context(), "gw", types.MergePatchType,
		[]byte(`{"status":{"qosPolicy":"qos-set-elsewhere"}}`), metav1.PatchOptions{}, "status")
	require.NoError(t, err)
	stale, err := ctrl.vpcNatGatewayLister.Get("gw")
	require.NoError(t, err)
	require.Empty(t, stale.Status.QoSPolicy, "the lister must still be stale for this test to mean anything")

	require.NoError(t, ctrl.patchNatGwHairpinSNATStatus("gw", hairpinSNATEnable))

	got, err := gwClient.Get(t.Context(), "gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, hairpinSNATEnable, got.Status.HairpinSNAT)
	require.Equal(t, "qos-set-elsewhere", got.Status.QoSPolicy)
}

// ... and the other way round: a whole-status writer must not clear the hairpin field it never
// observed, which is what the omitempty on the status field buys. Both writers serialize the
// whole struct from their own lister snapshot, and handleInitVpcNatGw runs them a few lines
// apart in one reconcile, so the lister cannot have caught up in between: each of them has to
// leave the field alone.
func TestNatGwStatusWritersKeepHairpinSNAT(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		write func(ctrl *Controller) error
	}{
		{
			name:  "patchNatGwStatus",
			write: func(ctrl *Controller) error { return ctrl.patchNatGwStatus("gw") },
		},
		{
			name:  "patchNatGwQoSStatus",
			write: func(ctrl *Controller) error { return ctrl.patchNatGwQoSStatus("gw", "qos") },
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gw := fakeGw("gw")
			gw.Spec.ExternalSubnets = []string{"ext"}
			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
			})
			require.NoError(t, err)
			ctrl := fc.fakeController
			gwClient := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways()

			require.NoError(t, ctrl.patchNatGwHairpinSNATStatus("gw", hairpinSNATEnable))
			require.NoError(t, ctrl.patchNatGwTCPConnLimitStatus("gw", 50000))
			// the lister snapshot these writers use has neither hairpinSnat nor tcpConnLimit
			require.NoError(t, tt.write(ctrl))

			got, err := gwClient.Get(t.Context(), "gw", metav1.GetOptions{})
			require.NoError(t, err)
			require.Equal(t, hairpinSNATEnable, got.Status.HairpinSNAT)
			require.Equal(t, ptr.To(int32(50000)), got.Status.TCPConnLimit)
		})
	}
}

// A gateway Pod that still runs a script without hairpin-snat-set (controller upgraded ahead of
// the gateway image) must not be sent the command: it fails on every retry and stops the whole
// gateway reconcile from converging. The empty status is what marks such a Pod.
func TestSyncNatGwHairpinSNATSkipsUnprogrammedPod(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		status  string
		enabled bool
	}{
		{name: "never programmed, default gateway"},
		{name: "never programmed, gateway asking for the fallback", enabled: true},
		{name: "recorded unsupported, default gateway", status: hairpinSNATUnsupported},
		{name: "recorded unsupported, gateway asking for the fallback", status: hairpinSNATUnsupported, enabled: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gw := fakeGw("gw")
			gw.Status.HairpinSNAT = tt.status
			gw.Spec.EnableLanIPAsHairpinSNAT = ptr.To(tt.enabled)
			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
			})
			require.NoError(t, err)
			ctrl := fc.fakeController

			var calls [][]string
			ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
				calls = append(calls, append([]string{operation}, rules...))
				return errors.New("Unknown command: hairpin-snat-set")
			}

			require.NoError(t, ctrl.syncNatGwHairpinSNAT(gw, fakeNatGwPod("gw", true)))
			require.Empty(t, calls)
		})
	}
}

// At init the command is sent, and a script that does not support it is tolerated only when the
// gateway does not ask for the fallback: the status then stays empty instead of claiming a state
// the data plane does not have. A gateway that does ask for it fails loudly, because the running
// image cannot deliver it.
func TestHandleInitVpcNatGwToleratesLegacyScript(t *testing.T) {
	for _, tt := range []struct {
		name      string
		enabled   bool
		expectErr bool
	}{
		{name: "default gateway initializes anyway"},
		{name: "requested fallback fails loudly", enabled: true, expectErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			old := vpcNatEnabled
			vpcNatEnabled = "true"
			t.Cleanup(func() { vpcNatEnabled = old })

			gw := fakeGw("gw")
			gw.Spec.EnableLanIPAsHairpinSNAT = ptr.To(tt.enabled)
			pod := fakeNatGwPod("gw", false)
			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
				Pods:           []*corev1.Pod{pod},
			})
			require.NoError(t, err)
			ctrl := fc.fakeController
			for _, q := range []*workqueue.TypedRateLimitingInterface[string]{
				&ctrl.updateVpcFloatingIPQueue, &ctrl.updateVpcDnatQueue, &ctrl.updateVpcSnatQueue,
				&ctrl.updateVpcSubnetQueue, &ctrl.updateVpcEipQueue,
			} {
				*q = newTypedRateLimitingQueue[string]("test", nil)
				t.Cleanup((*q).ShutDown)
			}
			ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, _ []string) error {
				if operation == natGwHairpinSNATSet {
					return errors.New("Unknown command: hairpin-snat-set")
				}
				return nil
			}

			err = ctrl.handleInitVpcNatGw("gw")
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			patched, getErr := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "gw", metav1.GetOptions{})
			require.NoError(t, getErr)
			require.Equal(t, hairpinSNATUnsupported, patched.Status.HairpinSNAT,
				"a script that cannot switch the fallback must be recorded as such, not as a programmed state")
			// ... and the Pod is still marked initialized, i.e. the gateway is usable
			initialized, podErr := ctrl.config.KubeClient.CoreV1().Pods(pod.Namespace).Get(t.Context(), pod.Name, metav1.GetOptions{})
			require.NoError(t, podErr)
			require.Equal(t, "true", initialized.Annotations[util.VpcNatGatewayInitAnnotation])
		})
	}
}

func TestNatGwDesiredTCPConnLimit(t *testing.T) {
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{})
	require.NoError(t, err)
	ctrl := fc.fakeController

	// 1. Default when unset: 50,000
	gw := fakeGw("gw")
	assert.Equal(t, int32(util.DefaultNatGwTCPConnLimit), ctrl.getNatGwDesiredTCPConnLimit(gw))

	// 2. Custom annotation tcp-connlimit
	gw.Annotations = map[string]string{util.VpcNatTCPConnLimitAnnotation: "15000"}
	assert.Equal(t, int32(15000), ctrl.getNatGwDesiredTCPConnLimit(gw))

	// 3. Spec field overrides annotation
	gw.Spec.TCPConnLimit = ptr.To(int32(30000))
	assert.Equal(t, int32(30000), ctrl.getNatGwDesiredTCPConnLimit(gw))

	// 4. Spec field 0 disables connlimit
	gw.Spec.TCPConnLimit = ptr.To(int32(0))
	assert.Equal(t, int32(0), ctrl.getNatGwDesiredTCPConnLimit(gw))

	// 5. Negative spec field falls back to default directly (authoritative, does not fall through to annotation)
	gw.Spec.TCPConnLimit = ptr.To(int32(-1))
	gw.Annotations = map[string]string{util.VpcNatTCPConnLimitAnnotation: "15000"}
	assert.Equal(t, int32(util.DefaultNatGwTCPConnLimit), ctrl.getNatGwDesiredTCPConnLimit(gw))

	// 6. Removing spec field restores annotation or default
	gw.Spec.TCPConnLimit = nil
	assert.Equal(t, int32(15000), ctrl.getNatGwDesiredTCPConnLimit(gw))
	delete(gw.Annotations, util.VpcNatTCPConnLimitAnnotation)
	assert.Equal(t, int32(util.DefaultNatGwTCPConnLimit), ctrl.getNatGwDesiredTCPConnLimit(gw))
}

func TestSyncNatGwTCPConnLimit(t *testing.T) {
	gw := fakeGw("gw")
	pod := fakeNatGwPod("gw", true)
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		Pods:           []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	ctrl := fc.fakeController

	var calls [][]string
	ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
		calls = append(calls, append([]string{operation}, rules...))
		return nil
	}

	// 1. Initial sync programs default (50000, stateful)
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, pod, true))
	require.Equal(t, [][]string{{natGwTCPConnLimitSet, "50000", "stateful"}}, calls)

	patched, err := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), "gw", metav1.GetOptions{})
	require.NoError(t, err)
	require.NotNil(t, patched.Status.TCPConnLimit)
	require.Equal(t, int32(50000), *patched.Status.TCPConnLimit)

	// 2. Converged: status matches desired, no re-exec
	calls = nil
	gw.Status.TCPConnLimit = ptr.To(int32(50000))
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[util.VpcNatTCPConnLimitStateAnnotation] = "50000"
	pod.Annotations[util.VpcNatGatewayDataplaneModeAnnotation] = "stateful"
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, pod, false))
	require.Empty(t, calls)

	// 3. Dynamic modification: spec changed to 20000
	gw.Spec.TCPConnLimit = ptr.To(int32(20000))
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, pod, false))
	require.Equal(t, [][]string{{natGwTCPConnLimitSet, "20000", "stateful"}}, calls)

	// 4. Disable: spec set to 0
	calls = nil
	gw.Status.TCPConnLimit = ptr.To(int32(20000))
	gw.Spec.TCPConnLimit = ptr.To(int32(0))
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, pod, false))
	require.Equal(t, [][]string{{natGwTCPConnLimitSet, "0", "stateful"}}, calls)

	// 5. Stateless mode: pass dataplane mode to script and publish effective limit 0
	calls = nil
	gw.Status.TCPConnLimit = ptr.To(int32(0))
	gw.Spec.TCPConnLimit = ptr.To(int32(15000))
	gw.Annotations = map[string]string{util.VpcNatGatewayDataplaneModeAnnotation: "stateless"}
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, pod, false))
	require.Equal(t, [][]string{{natGwTCPConnLimitSet, "15000", "stateless"}}, calls)
	patchedGw, getErr := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), gw.Name, metav1.GetOptions{})
	require.NoError(t, getErr)
	require.Equal(t, ptr.To(int32(0)), patchedGw.Status.TCPConnLimit)

	// 6. Negative spec value safely falls back to default 50000
	calls = nil
	delete(gw.Annotations, util.VpcNatGatewayDataplaneModeAnnotation)
	gw.Spec.TCPConnLimit = ptr.To(int32(-100))
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, pod, false))
	require.Equal(t, [][]string{{natGwTCPConnLimitSet, "50000", "stateful"}}, calls)
}

func TestSyncNatGwTCPConnLimit_ReplacedPodAndSkew(t *testing.T) {
	gw := fakeGw("gw")
	gw.Status.TCPConnLimit = ptr.To(int32(10000))
	gw.Spec.TCPConnLimit = ptr.To(int32(10000))

	newPod := fakeNatGwPod("gw", true)
	skewPod := fakeNatGwPod("gw-skew", true)

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		Pods:           []*corev1.Pod{newPod, skewPod},
	})
	require.NoError(t, err)
	ctrl := fc.fakeController

	var calls [][]string
	ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
		calls = append(calls, append([]string{operation}, rules...))
		return nil
	}

	// Case 1: Replaced pod during init has force=true, must program even though status.tcpConnLimit == desired
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, newPod, true))
	require.Equal(t, [][]string{{natGwTCPConnLimitSet, "10000", "stateful"}}, calls)

	// Case 2: Dataplane mode change without threshold change triggers convergence
	calls = nil
	gw.Annotations = map[string]string{util.VpcNatGatewayDataplaneModeAnnotation: "stateless"}
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, newPod, false))
	require.Equal(t, [][]string{{natGwTCPConnLimitSet, "10000", "stateless"}}, calls)

	// Case 3: Unsupported older image returns unknown command error -> handled gracefully
	calls = nil
	ctrl.execRulesInPod = func(_ *corev1.Pod, _ string, _ []string) error {
		return errors.New("exit status 1: Unknown command: tcp-connlimit-set")
	}
	require.NoError(t, ctrl.syncNatGwTCPConnLimit(gw, skewPod, false))
	patchedGw, getErr := ctrl.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(t.Context(), gw.Name, metav1.GetOptions{})
	require.NoError(t, getErr)
	require.Nil(t, patchedGw.Status.TCPConnLimit, "status.tcpConnLimit must be cleared when pod image does not support the command")
}

func TestInitVpcNatGateway_SecurityRuleFailurePropagation(t *testing.T) {
	for _, tc := range []struct {
		name       string
		setupGw    func(*kubeovnv1.VpcNatGateway)
		failOp     string
		failErr    error
		expectErr  bool
		expectAnno string
		expectVal  string
	}{
		{
			name:       "connlimit failure tolerated during init (fail-open)",
			failOp:     natGwTCPConnLimitSet,
			failErr:    errors.New("iptables rule installation failed"),
			expectAnno: util.VpcNatTCPConnLimitStateAnnotation,
			expectVal:  "",
		},
		{
			name: "custom connlimit failure blocks init",
			setupGw: func(g *kubeovnv1.VpcNatGateway) {
				g.Spec.TCPConnLimit = ptr.To(int32(20000))
			},
			failOp:    natGwTCPConnLimitSet,
			failErr:   errors.New("iptables rule installation failed"),
			expectErr: true,
		},
		{
			name:       "syn limit failure tolerated during init (fail-open)",
			failOp:     natGwTCPSynLimitSet,
			failErr:    errors.New("nftables rule installation failed"),
			expectAnno: util.VpcNatTCPSynLimitStateAnnotation,
			expectVal:  "",
		},
		{
			name: "custom syn limit failure blocks init",
			setupGw: func(g *kubeovnv1.VpcNatGateway) {
				g.Annotations = map[string]string{util.VpcNatTCPSynDstRateLimitAnnotation: "5000/s"}
			},
			failOp:    natGwTCPSynLimitSet,
			failErr:   errors.New("nftables rule installation failed"),
			expectErr: true,
		},
		{
			name:       "unsupported connlimit command tolerated during init",
			failOp:     natGwTCPConnLimitSet,
			failErr:    errors.New("exit status 1: stderr=Unknown command: tcp-connlimit-set stdout="),
			expectAnno: util.VpcNatTCPConnLimitStateAnnotation,
			expectVal:  tcpConnLimitUnsupported,
		},
		{
			name:       "unsupported syn limit command tolerated during init",
			failOp:     natGwTCPSynLimitSet,
			failErr:    errors.New("exit status 1: stderr=Unknown command: tcp-syn-limit-set stdout="),
			expectAnno: util.VpcNatTCPSynLimitStateAnnotation,
			expectVal:  tcpSynLimitUnsupported,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			old := vpcNatEnabled
			vpcNatEnabled = "true"
			t.Cleanup(func() { vpcNatEnabled = old })

			gw := fakeGw("gw-sec")
			if tc.setupGw != nil {
				tc.setupGw(gw)
			}
			pod := fakeNatGwPod("gw-sec", false)
			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
				Pods:           []*corev1.Pod{pod},
			})
			require.NoError(t, err)
			ctrl := fc.fakeController
			for _, q := range []*workqueue.TypedRateLimitingInterface[string]{
				&ctrl.updateVpcFloatingIPQueue, &ctrl.updateVpcDnatQueue, &ctrl.updateVpcSnatQueue,
				&ctrl.updateVpcSubnetQueue, &ctrl.updateVpcEipQueue,
			} {
				*q = newTypedRateLimitingQueue[string]("test", nil)
				t.Cleanup((*q).ShutDown)
			}
			ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, _ []string) error {
				if operation == tc.failOp {
					return tc.failErr
				}
				return nil
			}

			err = ctrl.handleInitVpcNatGw("gw-sec")
			if tc.expectErr {
				require.Error(t, err, "custom limiter failure must fail init to alert operator")
				return
			}
			require.NoError(t, err, "limiter failure must not gate data plane initialization")
			initializedPod, getErr := ctrl.config.KubeClient.CoreV1().Pods(pod.Namespace).Get(t.Context(), pod.Name, metav1.GetOptions{})
			require.NoError(t, getErr)
			require.Equal(t, "true", initializedPod.Annotations[util.VpcNatGatewayInitAnnotation], "gateway must be marked initialized")
			require.Equal(t, tc.expectVal, initializedPod.Annotations[tc.expectAnno], "failed limiter must be marked unsupported on pod")
		})
	}
}

func TestSyncNatGwTCPSynRateLimit(t *testing.T) {
	gw := fakeGw("gw")
	pod := fakeNatGwPod("gw", true)
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		Pods:           []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	ctrl := fc.fakeController

	var calls [][]string
	ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
		calls = append(calls, append([]string{operation}, rules...))
		return nil
	}

	// 1. Default rate limits (must not mutate caller's cached pod object in-place)
	require.NoError(t, ctrl.syncNatGwTCPSynRateLimit(gw, pod, false))
	require.Empty(t, pod.Annotations[util.VpcNatTCPSynLimitStateAnnotation], "caller pod object must not be mutated in-place")
	require.Equal(t, [][]string{{
		natGwTCPSynLimitSet,
		util.DefaultNatGwTCPSynSrcRate,
		util.DefaultNatGwTCPSynSrcBurst,
		util.DefaultNatGwTCPSynDstRate,
		util.DefaultNatGwTCPSynDstBurst,
	}}, calls)

	// Repeated reconcile without changes (using updated pod from apiserver) must skip reprogramming to preserve buckets
	pod, err = ctrl.config.KubeClient.CoreV1().Pods(pod.Namespace).Get(t.Context(), pod.Name, metav1.GetOptions{})
	require.NoError(t, err)
	calls = nil
	require.NoError(t, ctrl.syncNatGwTCPSynRateLimit(gw, pod, false))
	require.Empty(t, calls, "unchanged syn rate limit must not reprogram pod")

	// 2. Custom dst and src rate limits via annotations
	calls = nil
	gw.Annotations = map[string]string{
		util.VpcNatTCPSynSrcRateLimitAnnotation:  "5000/s",
		util.VpcNatTCPSynSrcBurstLimitAnnotation: "10000",
		util.VpcNatTCPSynDstRateLimitAnnotation:  "8000/s",
		util.VpcNatTCPSynDstBurstLimitAnnotation: "16000",
	}
	require.NoError(t, ctrl.syncNatGwTCPSynRateLimit(gw, pod, false))
	require.Equal(t, [][]string{{
		natGwTCPSynLimitSet,
		"5000/s",
		"10000",
		"8000/s",
		"16000",
	}}, calls)

	// 3. User input containing 'unknown command' in value must NOT be misclassified as unsupported command
	ctrl.execRulesInPod = func(_ *corev1.Pod, _ string, _ []string) error {
		return errors.New("exit status 1: Error: invalid SYN dst_rate 'unknown command'")
	}
	gw.Annotations[util.VpcNatTCPSynDstRateLimitAnnotation] = "unknown command"
	err = ctrl.syncNatGwTCPSynRateLimit(gw, pod, false)
	require.Error(t, err)
	require.NotEqual(t, tcpSynLimitUnsupported, pod.Annotations[util.VpcNatTCPSynLimitStateAnnotation])

	// 4. Disable destination SYN rate limit via "0"
	calls = nil
	ctrl.execRulesInPod = func(_ *corev1.Pod, operation string, rules []string) error {
		calls = append(calls, append([]string{operation}, rules...))
		return nil
	}
	gw.Annotations = map[string]string{
		util.VpcNatTCPSynDstRateLimitAnnotation: "0",
	}
	require.NoError(t, ctrl.syncNatGwTCPSynRateLimit(gw, pod, false))
	require.Equal(t, [][]string{{
		natGwTCPSynLimitSet,
		util.DefaultNatGwTCPSynSrcRate,
		util.DefaultNatGwTCPSynSrcBurst,
		"0",
		util.DefaultNatGwTCPSynDstBurst,
	}}, calls)

	// 5. Dataplane mode trims whitespace and falls back to stateful when empty
	gw.Annotations = map[string]string{util.VpcNatGatewayDataplaneModeAnnotation: "  "}
	require.Equal(t, "stateful", ctrl.getNatGwDataplaneMode(gw))
	gw.Annotations = map[string]string{util.VpcNatGatewayDataplaneModeAnnotation: " stateless "}
	require.Equal(t, "stateless", ctrl.getNatGwDataplaneMode(gw))
}

func TestIsUnsupportedCommandError(t *testing.T) {
	require.True(t, isUnsupportedCommandError(errors.New("exit status 1: Unknown command: tcp-connlimit-set"), natGwTCPConnLimitSet))
	require.True(t, isUnsupportedCommandError(errors.New("exit status 1: stderr=Unknown command: tcp-connlimit-set stdout="), natGwTCPConnLimitSet))
	require.True(t, isUnsupportedCommandError(errors.New("exit status 1: stderr=UNSUPPORTED_COMMAND: tcp-connlimit-set stdout="), natGwTCPConnLimitSet))
	require.True(t, isUnsupportedCommandError(errors.New("exit status 127: stderr=UNSUPPORTED_COMMAND: tcp-connlimit-set\nUnknown command: tcp-connlimit-set stdout="), natGwTCPConnLimitSet))
	require.True(t, isUnsupportedCommandError(errors.New("Unknown command: tcp-syn-limit-set"), natGwTCPSynLimitSet))
	require.False(t, isUnsupportedCommandError(errors.New("Error: invalid SYN src_rate 'unknown command'"), natGwTCPSynLimitSet))
	require.False(t, isUnsupportedCommandError(errors.New("exit status 1: Unknown command: other-op"), natGwTCPConnLimitSet))

	// False positive prevention: user input reflecting sentinel or unknown command inside validation error MUST NOT match
	require.False(t, isUnsupportedCommandError(errors.New("exit status 1: stderr=Error: invalid SYN src_rate 'UNSUPPORTED_COMMAND: tcp-syn-limit-set' stdout="), natGwTCPSynLimitSet))
	require.False(t, isUnsupportedCommandError(errors.New("exit status 1: stderr=Error: invalid SYN src_rate 'Unknown command: tcp-syn-limit-set' stdout="), natGwTCPSynLimitSet))
	require.False(t, isUnsupportedCommandError(errors.New("exit status 1: stderr=Error: invalid tcp connlimit value 'UNSUPPORTED_COMMAND: tcp-connlimit-set': must be an integer stdout="), natGwTCPConnLimitSet))

	// False positive prevention: "unknown command: tcp-connlimit-set" inside stdout must NOT match when stderr is present
	require.False(t, isUnsupportedCommandError(errors.New("exit status 1: stderr=iptables error stdout=unknown command: tcp-connlimit-set"), natGwTCPConnLimitSet))
	require.False(t, isUnsupportedCommandError(nil, natGwTCPConnLimitSet))
}
