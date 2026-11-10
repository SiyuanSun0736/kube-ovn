package controller

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/ptr"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/ipam"
	"github.com/kubeovn/kube-ovn/pkg/ovs"
	"github.com/kubeovn/kube-ovn/pkg/ovsdb/ovnnb"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// readCounter reads a metric value without pulling the prometheus testutil
// package (and its extra dependencies) into the test dependency tree.
func readCounter(t *testing.T, c prometheus.Counter) float64 {
	t.Helper()
	var m dto.Metric
	require.NoError(t, c.Write(&m))
	return m.GetCounter().GetValue()
}

func TestCheckIsPodVpcNatGw(t *testing.T) {
	tests := []struct {
		name                string
		pod                 *corev1.Pod
		networkAttachments  []*nadv1.NetworkAttachmentDefinition
		subnets             []*kubeovnv1.Subnet
		enableNonPrimaryCNI bool
		expectedIsVpcNatGw  bool
		expectedVpcGwName   string
		description         string
	}{
		{
			name: "Pod with default provider VPC NAT gateway annotation",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						util.VpcNatGatewayAnnotation: "test-nat-gw",
					},
				},
			},
			networkAttachments:  []*nadv1.NetworkAttachmentDefinition{},
			subnets:             []*kubeovnv1.Subnet{},
			enableNonPrimaryCNI: false,
			expectedIsVpcNatGw:  true,
			expectedVpcGwName:   "test-nat-gw",
			description:         "Should detect VPC NAT gateway with default provider",
		},
		{
			name: "Pod with custom provider VPC NAT gateway annotation in non-primary CNI mode",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						// Network attachment annotation to indicate this pod uses net1
						nadv1.NetworkAttachmentAnnot: `[{"name": "net1"}]`,
						// Custom provider VPC NAT gateway annotation
						util.VpcNatGatewayAnnotation: "test-nat-gw",
						// Kube-OVN annotations for net1 provider
						fmt.Sprintf(util.LogicalSwitchAnnotationTemplate, "net1.default.ovn"): "net1-subnet",
						fmt.Sprintf(util.LogicalRouterAnnotationTemplate, "net1.default.ovn"): "net1-vpc",
						fmt.Sprintf(util.IPAddressAnnotationTemplate, "net1.default.ovn"):     "192.168.1.10",
					},
				},
			},
			networkAttachments: []*nadv1.NetworkAttachmentDefinition{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "net1",
						Namespace: "default",
					},
					Spec: nadv1.NetworkAttachmentDefinitionSpec{
						Config: `{
							"cniVersion": "0.3.1",
							"name": "net1",
							"type": "kube-ovn",
							"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
							"provider": "net1.default.ovn"
						}`,
					},
				},
			},
			subnets: []*kubeovnv1.Subnet{
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
			enableNonPrimaryCNI: true,
			expectedIsVpcNatGw:  true,
			expectedVpcGwName:   "test-nat-gw",
			description:         "Should detect VPC NAT gateway with custom provider in non-primary CNI mode",
		},
		{
			name: "Pod without VPC NAT gateway annotation or with empty name",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						"other.annotation": "value",
					},
				},
			},
			networkAttachments:  []*nadv1.NetworkAttachmentDefinition{},
			subnets:             []*kubeovnv1.Subnet{},
			enableNonPrimaryCNI: false,
			expectedIsVpcNatGw:  false,
			expectedVpcGwName:   "",
			description:         "Should not detect VPC NAT gateway when annotation is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create controller with proper setup
			fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				NetworkAttachments: tt.networkAttachments,
				Subnets:            tt.subnets,
				Pods:               []*corev1.Pod{tt.pod},
			})
			require.NoError(t, err, "Failed to create fake controller")
			controller := fakeController.fakeController
			// Set the non-primary CNI mode
			controller.config.EnableNonPrimaryCNI = tt.enableNonPrimaryCNI

			// Call the method under test
			isVpcNatGw, vpcGwName := controller.checkIsPodVpcNatGw(tt.pod)

			// Verify results
			assert.Equal(t, tt.expectedIsVpcNatGw, isVpcNatGw, "IsVpcNatGw mismatch: %s", tt.description)
			assert.Equal(t, tt.expectedVpcGwName, vpcGwName, "VpcGwName mismatch: %s", tt.description)
		})
	}

	// Test additional edge cases in a single sub-test for efficiency
	t.Run("Edge cases", func(t *testing.T) {
		fakeController, err := newFakeControllerWithOptions(t, nil)
		require.NoError(t, err)
		controller := fakeController.fakeController
		// Test nil pod
		isVpcNatGw, vpcGwName := controller.checkIsPodVpcNatGw(nil)
		assert.False(t, isVpcNatGw, "Nil pod should not be VPC NAT gateway")
		assert.Equal(t, "", vpcGwName, "Nil pod should have empty gateway name")

		// Test pod with empty VPC NAT gateway name
		podWithEmptyGw := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test-pod",
				Namespace:   "default",
				Annotations: map[string]string{util.VpcNatGatewayAnnotation: ""},
			},
		}
		isVpcNatGw, vpcGwName = controller.checkIsPodVpcNatGw(podWithEmptyGw)
		assert.False(t, isVpcNatGw, "Pod with empty gateway name should not be VPC NAT gateway")
		assert.Equal(t, "", vpcGwName, "Pod with empty gateway name should return empty")

		// Test pod with no annotations
		podNoAnnotations := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test-pod",
				Namespace:   "default",
				Annotations: nil,
			},
		}
		isVpcNatGw, vpcGwName = controller.checkIsPodVpcNatGw(podNoAnnotations)
		assert.False(t, isVpcNatGw, "Pod with no annotations should not be VPC NAT gateway")
		assert.Equal(t, "", vpcGwName, "Pod with no annotations should return empty")
	})
}

func TestBackfillVpcNatGwLanIPFromPod(t *testing.T) {
	const (
		gwName    = "test-nat-gw"
		subnet    = "nat-subnet"
		provider  = "net1.default.ovn"
		lanIP     = "10.244.0.10"
		namespace = "default"
	)

	tests := []struct {
		name                   string
		gwSpecLanIP            string
		subnetProtocol         string
		givenGwName            string
		podOwnerName           string
		podNamespace           string
		controllerPodNamespace string
		podAnnotation          map[string]string
		expectedLanIP          string
	}{
		{
			name:                   "backfill lanIP from pod annotation",
			gwSpecLanIP:            "",
			subnetProtocol:         kubeovnv1.ProtocolIPv4,
			givenGwName:            gwName,
			podOwnerName:           util.GenNatGwName(gwName),
			podNamespace:           namespace,
			controllerPodNamespace: namespace,
			podAnnotation: map[string]string{
				fmt.Sprintf(util.IPAddressAnnotationTemplate, provider): lanIP,
			},
			expectedLanIP: lanIP,
		},
		{
			name:                   "derive gateway name from owner reference",
			gwSpecLanIP:            "",
			subnetProtocol:         kubeovnv1.ProtocolIPv4,
			givenGwName:            "",
			podOwnerName:           util.GenNatGwName(gwName),
			podNamespace:           namespace,
			controllerPodNamespace: namespace,
			podAnnotation: map[string]string{
				fmt.Sprintf(util.IPAddressAnnotationTemplate, provider): lanIP,
			},
			expectedLanIP: lanIP,
		},
		{
			name:                   "skip when spec lanIP already set",
			gwSpecLanIP:            "10.244.0.99",
			subnetProtocol:         kubeovnv1.ProtocolIPv4,
			givenGwName:            gwName,
			podOwnerName:           util.GenNatGwName(gwName),
			podNamespace:           namespace,
			controllerPodNamespace: namespace,
			podAnnotation: map[string]string{
				fmt.Sprintf(util.IPAddressAnnotationTemplate, provider): lanIP,
			},
			expectedLanIP: "10.244.0.99",
		},
		{
			name:                   "backfill lanIP from pod in custom namespace",
			gwSpecLanIP:            "",
			subnetProtocol:         kubeovnv1.ProtocolIPv4,
			givenGwName:            gwName,
			podOwnerName:           util.GenNatGwName(gwName),
			podNamespace:           "other-ns",
			controllerPodNamespace: namespace,
			podAnnotation: map[string]string{
				fmt.Sprintf(util.IPAddressAnnotationTemplate, provider): lanIP,
			},
			expectedLanIP: lanIP,
		},
		{
			name:                   "skip when lanIP annotation is invalid",
			gwSpecLanIP:            "",
			subnetProtocol:         kubeovnv1.ProtocolIPv4,
			givenGwName:            gwName,
			podOwnerName:           util.GenNatGwName(gwName),
			podNamespace:           namespace,
			controllerPodNamespace: namespace,
			podAnnotation: map[string]string{
				fmt.Sprintf(util.IPAddressAnnotationTemplate, provider): "not-an-ip",
			},
			expectedLanIP: "",
		},
		{
			name:                   "prefer IPv6 address for IPv6 subnet",
			gwSpecLanIP:            "",
			subnetProtocol:         kubeovnv1.ProtocolIPv6,
			givenGwName:            gwName,
			podOwnerName:           util.GenNatGwName(gwName),
			podNamespace:           namespace,
			controllerPodNamespace: namespace,
			podAnnotation: map[string]string{
				fmt.Sprintf(util.IPAddressAnnotationTemplate, provider): "10.244.0.10,fd00:10:16::10",
			},
			expectedLanIP: "fd00:10:16::10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gw := &kubeovnv1.VpcNatGateway{
				ObjectMeta: metav1.ObjectMeta{
					Name: gwName,
				},
				Spec: kubeovnv1.VpcNatGatewaySpec{
					Vpc:    "vpc-a",
					Subnet: subnet,
					LanIP:  tt.gwSpecLanIP,
				},
			}
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        util.GenNatGwPodName(gwName),
					Namespace:   tt.podNamespace,
					Annotations: tt.podAnnotation,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: appsv1.SchemeGroupVersion.String(),
							Kind:       util.KindStatefulSet,
							Name:       tt.podOwnerName,
						},
					},
				},
			}

			fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				Subnets: []*kubeovnv1.Subnet{
					{
						ObjectMeta: metav1.ObjectMeta{Name: subnet},
						Spec: kubeovnv1.SubnetSpec{
							Provider: provider,
							Protocol: tt.subnetProtocol,
						},
					},
				},
				VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
			})
			require.NoError(t, err)

			controller := fakeController.fakeController
			controller.config.PodNamespace = tt.controllerPodNamespace
			err = controller.backfillVpcNatGwLanIPFromPod(pod, tt.givenGwName)
			require.NoError(t, err)

			gotGw, err := controller.config.KubeOvnClient.KubeovnV1().VpcNatGateways().Get(
				context.Background(), gwName, metav1.GetOptions{},
			)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedLanIP, gotGw.Spec.LanIP)
		})
	}
}

func TestGetPodKubeovnNetsNonPrimaryCNI(t *testing.T) {
	tests := []struct {
		name                string
		pod                 *corev1.Pod
		networkAttachments  []*nadv1.NetworkAttachmentDefinition
		subnets             []*kubeovnv1.Subnet
		enableNonPrimaryCNI bool
		expectedNetCount    int
		expectError         bool
		description         string
	}{
		{
			name: "Non-primary CNI mode with network attachments",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						nadv1.NetworkAttachmentAnnot: `[{"name": "net1"}]`,
						// Kube-OVN annotations for net1 provider
						fmt.Sprintf(util.LogicalSwitchAnnotationTemplate, "net1.default.ovn"): "net1-subnet",
						fmt.Sprintf(util.LogicalRouterAnnotationTemplate, "net1.default.ovn"): "net1-vpc",
						fmt.Sprintf(util.IPAddressAnnotationTemplate, "net1.default.ovn"):     "192.168.1.10",
					},
				},
			},
			networkAttachments: []*nadv1.NetworkAttachmentDefinition{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "net1",
						Namespace: "default",
					},
					Spec: nadv1.NetworkAttachmentDefinitionSpec{
						Config: `{
							"cniVersion": "0.3.1",
							"name": "net1",
							"type": "kube-ovn",
							"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
							"provider": "net1.default.ovn"
						}`,
					},
				},
			},
			subnets: []*kubeovnv1.Subnet{
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
			enableNonPrimaryCNI: true,
			expectedNetCount:    1,
			expectError:         false,
			description:         "Should return only network attachment definitions in non-primary CNI mode",
		},
		{
			name: "Primary CNI mode vs Non-primary CNI behavior",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						nadv1.NetworkAttachmentAnnot: `[{"name": "net1"}]`,
						// Both custom and default provider annotations
						fmt.Sprintf(util.LogicalSwitchAnnotationTemplate, "net1.default.ovn"): "net1-subnet",
						fmt.Sprintf(util.LogicalSwitchAnnotationTemplate, util.OvnProvider):   "ovn-default",
						fmt.Sprintf(util.IPAddressAnnotationTemplate, "net1.default.ovn"):     "192.168.1.10",
						fmt.Sprintf(util.IPAddressAnnotationTemplate, util.OvnProvider):       "10.244.0.5",
					},
				},
			},
			networkAttachments: []*nadv1.NetworkAttachmentDefinition{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "net1",
						Namespace: "default",
					},
					Spec: nadv1.NetworkAttachmentDefinitionSpec{
						Config: `{
							"cniVersion": "0.3.1",
							"name": "net1",
							"type": "kube-ovn",
							"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
							"provider": "net1.default.ovn"
						}`,
					},
				},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "net1-subnet",
					},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "192.168.1.0/24",
						Provider:  "net1.default.ovn",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "ovn-default",
					},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.244.0.0/24",
						Provider:  util.OvnProvider,
						Default:   true,
					},
				},
			},
			enableNonPrimaryCNI: false, // This test will verify both modes
			expectedNetCount:    2,     // Both networks in primary mode
			expectError:         false,
			description:         "Should handle both network attachments and default network differently in primary vs non-primary modes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create controller with proper setup
			fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				NetworkAttachments: tt.networkAttachments,
				Subnets:            tt.subnets,
				Pods:               []*corev1.Pod{tt.pod},
			})
			require.NoError(t, err, "Failed to create fake controller")
			controller := fakeController.fakeController

			// Set the non-primary CNI mode
			controller.config.EnableNonPrimaryCNI = tt.enableNonPrimaryCNI

			// Call the method under test
			nets, err := controller.getPodKubeovnNets(tt.pod)

			// Check for errors
			if tt.expectError {
				assert.Error(t, err, "Expected an error but got none: %s", tt.description)
				return
			}
			require.NoError(t, err, "Unexpected error: %s", tt.description)

			// Verify network count
			assert.Equal(t, tt.expectedNetCount, len(nets), "Network count mismatch: %s", tt.description)

			// For the comparison test, also test non-primary mode
			if tt.name == "Primary CNI mode vs Non-primary CNI behavior" {
				controller.config.EnableNonPrimaryCNI = true
				netsNonPrimary, err := controller.getPodKubeovnNets(tt.pod)
				require.NoError(t, err, "Unexpected error in non-primary mode")
				assert.Equal(t, 1, len(netsNonPrimary), "Non-primary mode should return only network attachments")
			}
		})
	}
}

func TestAcquireAddressWithSpecifiedSubnet(t *testing.T) {
	tests := []struct {
		name           string
		pod            *corev1.Pod
		namespaces     []*corev1.Namespace
		subnets        []*kubeovnv1.Subnet
		setupIPAM      func(*Controller)
		expectError    bool
		expectedSubnet string
		description    string
	}{
		{
			name: "User specifies subnet - should succeed",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						util.LogicalSwitchAnnotation: "subnet1",
						util.IPAddressAnnotation:     "10.0.1.10",
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
						Annotations: map[string]string{
							util.LogicalSwitchAnnotation: "subnet1,subnet2",
						},
					},
				},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "subnet1"},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.0.1.0/24",
						Protocol:  kubeovnv1.ProtocolIPv4,
						Provider:  util.OvnProvider,
					},
					Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "subnet2"},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.0.1.0/24",
						Protocol:  kubeovnv1.ProtocolIPv4,
						Provider:  util.OvnProvider,
					},
					Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
				},
			},
			expectError:    false,
			expectedSubnet: "subnet1",
			description:    "Should allocate from specified subnet",
		},
		{
			name: "User specifies subnet but IP occupied - should NOT fallback",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						util.LogicalSwitchAnnotation: "subnet1",
						util.IPAddressAnnotation:     "10.0.1.10",
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
						Annotations: map[string]string{
							util.LogicalSwitchAnnotation: "subnet1,subnet2",
						},
					},
				},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "subnet1"},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.0.1.0/24",
						Protocol:  kubeovnv1.ProtocolIPv4,
						Provider:  util.OvnProvider,
					},
					Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "subnet2"},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.0.1.0/24",
						Protocol:  kubeovnv1.ProtocolIPv4,
						Provider:  util.OvnProvider,
					},
					Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
				},
			},
			setupIPAM: func(c *Controller) {
				_, _, _, _ = c.ipam.GetStaticAddress("other-pod.default", "other-pod.default", "10.0.1.10", nil, "subnet1", true)
			},
			expectError: true,
			description: "Should NOT fallback to subnet2 when IP is occupied in specified subnet1",
		},
		{
			name: "No subnet specified - should try all namespace subnets",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						util.IPAddressAnnotation: "10.0.2.10",
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
						Annotations: map[string]string{
							util.LogicalSwitchAnnotation: "subnet1,subnet2",
						},
					},
				},
			},
			subnets: []*kubeovnv1.Subnet{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "subnet1"},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.0.1.0/24",
						Protocol:  kubeovnv1.ProtocolIPv4,
						Provider:  util.OvnProvider,
					},
					Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "subnet2"},
					Spec: kubeovnv1.SubnetSpec{
						CIDRBlock: "10.0.2.0/24",
						Protocol:  kubeovnv1.ProtocolIPv4,
						Provider:  util.OvnProvider,
					},
					Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
				},
			},
			expectError:    false,
			expectedSubnet: "subnet2",
			description:    "Should try all subnets and find matching one when no subnet specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				Namespaces: tt.namespaces,
				Subnets:    tt.subnets,
				Pods:       []*corev1.Pod{tt.pod},
			})
			require.NoError(t, err)
			controller := fakeController.fakeController
			controller.ipam = newIPAMForTest(tt.subnets)

			if tt.setupIPAM != nil {
				tt.setupIPAM(controller)
			}

			podNets, err := controller.getPodKubeovnNets(tt.pod)
			require.NoError(t, err)
			require.Greater(t, len(podNets), 0)

			_, _, _, subnet, err := controller.acquireAddress(tt.pod, podNets[0])

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				require.NoError(t, err, tt.description)
				assert.Equal(t, tt.expectedSubnet, subnet.Name, tt.description)
			}
		})
	}
}

func newIPAMForTest(subnets []*kubeovnv1.Subnet) *ipam.IPAM {
	ipamInstance := ipam.NewIPAM()
	for _, subnet := range subnets {
		excludeIPs := subnet.Spec.ExcludeIps
		if len(excludeIPs) == 0 {
			excludeIPs = []string{}
		}
		s, err := ipam.NewSubnet(subnet.Name, subnet.Spec.CIDRBlock, excludeIPs)
		if err != nil {
			panic(err)
		}
		ipamInstance.Subnets[subnet.Name] = s
	}
	return ipamInstance
}

// TestReconcileAllocateSubnets_gatedOnTunnelKey verifies the core behavior of this change: a pod
// is not persisted with an IP from an OVN subnet whose VPC tunnel key has not been synced from
// OVN SB yet, even when the subnet's own key is synced; the allocation errors so the pod requeues.
func TestReconcileAllocateSubnets_gatedOnTunnelKey(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				util.LogicalSwitchAnnotation: "ovn-subnet",
			},
		},
	}
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ovn-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.1.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  util.OvnProvider,
			Vpc:       "vpc-a",
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100, TunnelKey: 42},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Namespaces: []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
		Subnets:    []*kubeovnv1.Subnet{subnet},
		Vpcs:       []*kubeovnv1.Vpc{{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}}},
		Pods:       []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.ipam = newIPAMForTest([]*kubeovnv1.Subnet{subnet})

	podNets, err := c.getPodKubeovnNets(pod)
	require.NoError(t, err)
	require.Greater(t, len(podNets), 0)

	_, err = c.reconcileAllocateSubnets(pod, podNets)
	require.Error(t, err)
	require.Contains(t, err.Error(), "vpc vpc-a tunnel key of subnet ovn-subnet not observed")

	persisted, err := c.config.KubeClient.CoreV1().Pods("default").Get(context.Background(), pod.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.NotEqual(t, "true", persisted.Annotations[util.AllocatedAnnotation])
	require.Empty(t, persisted.Annotations[util.IPAddressAnnotation])
	require.Empty(t, persisted.Annotations[util.TunnelKeyAnnotation])
}

// TestReconcileAllocateSubnets_allowedWhenTunnelKeySynced is the happy-path counterpart of the gate
// test: once the VPC tunnel key is synced, allocation must proceed and record the VPC key (not the
// subnet's own key) as the pod tunnel_key annotation. It guards against over-blocking allocation.
func TestReconcileAllocateSubnets_allowedWhenTunnelKeySynced(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				util.LogicalSwitchAnnotation: "ovn-subnet",
			},
		},
	}
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ovn-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.1.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  util.OvnProvider,
			Vpc:       "vpc-a",
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100, TunnelKey: 42},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Namespaces: []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
		Subnets:    []*kubeovnv1.Subnet{subnet},
		Vpcs:       []*kubeovnv1.Vpc{{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: 1234}}},
		Pods:       []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.ipam = newIPAMForTest([]*kubeovnv1.Subnet{subnet})

	// The only OVN NB call on the happy path is creating the logical switch port.
	fc.mockOvnClient.EXPECT().CreateLogicalSwitchPort(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)

	podNets, err := c.getPodKubeovnNets(pod)
	require.NoError(t, err)
	require.Greater(t, len(podNets), 0)

	allocated, err := c.reconcileAllocateSubnets(pod, podNets)
	require.NoError(t, err)
	require.NotNil(t, allocated)
	tunnelKeyKey := fmt.Sprintf(util.TunnelKeyAnnotationTemplate, podNets[0].ProviderName)
	require.Equal(t, "1234", allocated.Annotations[tunnelKeyKey])
	require.Equal(t, "true", allocated.Annotations[util.AllocatedAnnotation])
	require.NotEmpty(t, allocated.Annotations[util.IPAddressAnnotation])
	require.Equal(t, subnet.Name, allocated.Annotations[util.LogicalSwitchAnnotation])

	persisted, err := c.config.KubeClient.CoreV1().Pods("default").Get(context.Background(), pod.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, allocated.Annotations[util.IPAddressAnnotation], persisted.Annotations[util.IPAddressAnnotation])
	require.Equal(t, "1234", persisted.Annotations[util.TunnelKeyAnnotation])
	require.Equal(t, "true", persisted.Annotations[util.AllocatedAnnotation])
}

func TestReconcileAllocateSubnetsUsesResolvedSubnetForTunnelKey(t *testing.T) {
	initialSubnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ipam-only-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.3.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  "net1.default",
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
	}
	resolvedSubnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ovn-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.4.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  util.OvnProvider,
			Vpc:       "vpc-b",
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100, TunnelKey: 42},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				util.IPAddressAnnotation: "10.0.4.5",
			},
		},
	}
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				util.LogicalSwitchAnnotation: initialSubnet.Name + "," + resolvedSubnet.Name,
			},
		},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Namespaces: []*corev1.Namespace{namespace},
		Subnets:    []*kubeovnv1.Subnet{initialSubnet, resolvedSubnet},
		Vpcs:       []*kubeovnv1.Vpc{{ObjectMeta: metav1.ObjectMeta{Name: "vpc-b"}, Status: kubeovnv1.VpcStatus{TunnelKey: 2345}}},
		Pods:       []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.ipam = newIPAMForTest([]*kubeovnv1.Subnet{initialSubnet, resolvedSubnet})

	fc.mockOvnClient.EXPECT().CreateLogicalSwitchPort(
		resolvedSubnet.Name, gomock.Any(), "10.0.4.5", gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)

	podNets, err := c.getPodKubeovnNets(pod)
	require.NoError(t, err)
	require.Len(t, podNets, 1)
	require.Equal(t, initialSubnet.Name, podNets[0].Subnet.Name)

	allocated, err := c.reconcileAllocateSubnets(pod, podNets)
	require.NoError(t, err)
	require.Equal(t, resolvedSubnet.Name, allocated.Annotations[util.LogicalSwitchAnnotation])
	require.Equal(t, "2345", allocated.Annotations[util.TunnelKeyAnnotation])
}

func TestReconcileAllocateSubnetsOVNVlanKeepsNicAnnotationsWithoutTunnelKey(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				util.LogicalSwitchAnnotation: "vlan-subnet",
				util.TunnelKeyAnnotation:     "1234",
			},
		},
	}
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "vlan-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.5.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  util.OvnProvider,
			Vlan:      "vlan-a",
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
	}
	vlan := &kubeovnv1.Vlan{
		ObjectMeta: metav1.ObjectMeta{Name: "vlan-a"},
		Spec:       kubeovnv1.VlanSpec{ID: 100, Provider: "provider-a"},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Namespaces: []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
		Subnets:    []*kubeovnv1.Subnet{subnet},
		Vlans:      []*kubeovnv1.Vlan{vlan},
		Pods:       []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.config.PodNicType = util.VethType
	c.ipam = newIPAMForTest([]*kubeovnv1.Subnet{subnet})

	fc.mockOvnClient.EXPECT().CreateLogicalSwitchPort(
		subnet.Name, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)

	allocated, err := c.reconcileAllocateSubnets(pod, []*kubeovnNet{{
		Type:         providerTypeOriginal,
		ProviderName: util.OvnProvider,
		Subnet:       subnet,
	}})
	require.NoError(t, err)
	require.Equal(t, subnet.Name, allocated.Annotations[util.LogicalSwitchAnnotation])
	require.Equal(t, util.VethType, allocated.Annotations[fmt.Sprintf(util.PodNicAnnotationTemplate, util.OvnProvider)])
	_, ok := allocated.Annotations[util.TunnelKeyAnnotation]
	require.False(t, ok, "OVN vlan/underlay subnet must keep the default tunnel key")
}

func TestReconcileAllocateSubnetsClearsStaleTunnelKeyForNonOVNProvider(t *testing.T) {
	const provider = "net1.default"
	logicalSwitchKey := fmt.Sprintf(util.LogicalSwitchAnnotationTemplate, provider)
	podNicKey := fmt.Sprintf(util.PodNicAnnotationTemplate, provider)
	tunnelKeyKey := fmt.Sprintf(util.TunnelKeyAnnotationTemplate, provider)
	allocatedKey := fmt.Sprintf(util.AllocatedAnnotationTemplate, provider)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				logicalSwitchKey: "old-ovn-subnet",
				podNicKey:        util.VethType,
				tunnelKeyKey:     "1234",
			},
		},
	}
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ipam-only-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.2.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  provider,
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Namespaces: []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
		Subnets:    []*kubeovnv1.Subnet{subnet},
		Pods:       []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.ipam = newIPAMForTest([]*kubeovnv1.Subnet{subnet})

	allocated, err := c.reconcileAllocateSubnets(pod, []*kubeovnNet{{
		Type:         providerTypeIPAM,
		ProviderName: provider,
		Subnet:       subnet,
	}})
	require.NoError(t, err)
	require.NotNil(t, allocated)
	require.Equal(t, "true", allocated.Annotations[allocatedKey])
	_, ok := allocated.Annotations[logicalSwitchKey]
	require.False(t, ok, "non-OVN provider must not keep a stale logical_switch")
	_, ok = allocated.Annotations[podNicKey]
	require.False(t, ok, "non-OVN provider must not keep a stale pod_nic")
	_, ok = allocated.Annotations[tunnelKeyKey]
	require.False(t, ok, "non-OVN provider must not keep a stale tunnel_key")
}

func TestGetNamedPortByNsReturnsCopy(t *testing.T) {
	np := NewNamedPort()
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: 80},
					},
				},
			},
		},
	}

	np.AddNamedPortByPod(pod)

	result := np.GetNamedPortByNs("test-ns")
	require.NotNil(t, result)
	assert.Contains(t, result, "http")

	// Mutating the returned map should not affect internal state
	delete(result, "http")

	result2 := np.GetNamedPortByNs("test-ns")
	require.NotNil(t, result2)
	assert.Contains(t, result2, "http", "internal map should not be affected by mutation of returned copy")
}

func TestDeleteNamedPortByPodWithRestartableInitContainers(t *testing.T) {
	restartAlways := corev1.ContainerRestartPolicyAlways
	np := NewNamedPort()
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-pod",
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{
				{
					Name:          "sidecar",
					RestartPolicy: &restartAlways,
					Ports: []corev1.ContainerPort{
						{Name: "metrics", ContainerPort: 9090},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: 80},
					},
				},
			},
		},
	}

	np.AddNamedPortByPod(pod)
	result := np.GetNamedPortByNs("test-ns")
	require.NotNil(t, result)
	assert.Contains(t, result, "http")
	assert.Contains(t, result, "metrics")

	np.DeleteNamedPortByPod(pod)
	result = np.GetNamedPortByNs("test-ns")
	assert.Empty(t, result, "both regular and sidecar init container named ports should be deleted")
}

func TestVpcTunnelKey(t *testing.T) {
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Vpcs: []*kubeovnv1.Vpc{
			{ObjectMeta: metav1.ObjectMeta{Name: "synced"}, Status: kubeovnv1.VpcStatus{TunnelKey: 1234}},
			{ObjectMeta: metav1.ObjectMeta{Name: "pending"}},
		},
	})
	require.NoError(t, err)
	c := fc.fakeController

	subnet := func(vpc string) *kubeovnv1.Subnet {
		// The subnet's own key must never leak into the VPC key.
		return &kubeovnv1.Subnet{Spec: kubeovnv1.SubnetSpec{Vpc: vpc}, Status: kubeovnv1.SubnetStatus{TunnelKey: 42}}
	}
	assert.Equal(t, 1234, c.vpcTunnelKey(subnet("synced")))
	assert.Equal(t, 0, c.vpcTunnelKey(subnet("pending")), "unsynced VPC key is not ready")
	assert.Equal(t, 0, c.vpcTunnelKey(subnet("missing")), "unknown VPC is not ready")
	assert.Equal(t, 0, c.vpcTunnelKey(subnet("")), "unformatted subnet is not ready")
}

func TestHandleRepairTunnelKey(t *testing.T) {
	// The subnet's own key (42) differs from every VPC key so tests prove pods get the VPC key.
	ovnSubnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ovn-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.1.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  util.OvnProvider,
			Vpc:       "vpc-a",
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100, TunnelKey: 42},
	}

	ovnPod := func(tunnelKey *string) *corev1.Pod {
		annotations := map[string]string{
			util.LogicalSwitchAnnotation: "ovn-subnet",
			util.AllocatedAnnotation:     "true",
		}
		if tunnelKey != nil {
			annotations[util.TunnelKeyAnnotation] = *tunnelKey
		}
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test-pod",
				Namespace:   "default",
				Annotations: annotations,
			},
		}
	}

	deletingPod := ovnPod(nil)
	deletingPod.DeletionTimestamp = ptr.To(metav1.Now())
	hostNetworkPod := ovnPod(nil)
	hostNetworkPod.Spec.HostNetwork = true

	tests := []struct {
		name          string
		pod           *corev1.Pod
		subnet        *kubeovnv1.Subnet
		vpcTunnelKey  int
		wantErr       bool
		wantAnnotated bool
		wantValue     string
	}{
		{
			name:   "patches missing tunnel_key from VPC status",
			pod:    ovnPod(nil),
			subnet: ovnSubnet, vpcTunnelKey: 1234,
			wantAnnotated: true,
			wantValue:     "1234",
		},
		{
			name:   "leaves matching valid annotation untouched",
			pod:    ovnPod(ptr.To("1234")),
			subnet: ovnSubnet, vpcTunnelKey: 1234,
			wantAnnotated: true,
			wantValue:     "1234",
		},
		{
			name:   "repairs stale but valid annotation",
			pod:    ovnPod(ptr.To("999")),
			subnet: ovnSubnet, vpcTunnelKey: 1234,
			wantAnnotated: true,
			wantValue:     "1234",
		},
		{
			name:   "migrates subnet key to VPC key",
			pod:    ovnPod(ptr.To("42")),
			subnet: ovnSubnet, vpcTunnelKey: 1234,
			wantAnnotated: true,
			wantValue:     "1234",
		},
		{
			name:   "repairs zero annotation",
			pod:    ovnPod(ptr.To("0")),
			subnet: ovnSubnet, vpcTunnelKey: 1234,
			wantAnnotated: true,
			wantValue:     "1234",
		},
		{
			name:   "repairs non-numeric annotation",
			pod:    ovnPod(ptr.To("invalid")),
			subnet: ovnSubnet, vpcTunnelKey: 1234,
			wantAnnotated: true,
			wantValue:     "1234",
		},
		{
			name:   "repairs out-of-range annotation",
			pod:    ovnPod(ptr.To(strconv.Itoa(util.MaxTunnelKey + 1))),
			subnet: ovnSubnet, vpcTunnelKey: 1234,
			wantAnnotated: true,
			wantValue:     "1234",
		},
		{
			name:   "requeues when VPC tunnel key not synced yet",
			pod:    ovnPod(nil),
			subnet: ovnSubnet, vpcTunnelKey: 0,
			wantErr: true,
		},
		{
			name:   "requeues when VPC tunnel key is out of range",
			pod:    ovnPod(nil),
			subnet: ovnSubnet, vpcTunnelKey: util.MaxTunnelKey + 1,
			wantErr: true,
		},
		{
			name: "removes stale key from OVN vlan underlay subnet",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						util.LogicalSwitchAnnotation: "vlan-subnet",
						util.AllocatedAnnotation:     "true",
						util.TunnelKeyAnnotation:     "1234",
					},
				},
			},
			subnet: &kubeovnv1.Subnet{
				ObjectMeta: metav1.ObjectMeta{Name: "vlan-subnet"},
				Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vlan: "vlan-a"},
				Status:     kubeovnv1.SubnetStatus{TunnelKey: 0},
			},
		},
		{
			name: "removes stale key without logical switch",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						util.AllocatedAnnotation: "true",
						util.TunnelKeyAnnotation: "7",
					},
				},
			},
			subnet: &kubeovnv1.Subnet{
				ObjectMeta: metav1.ObjectMeta{Name: "unused-subnet"},
				Spec:       kubeovnv1.SubnetSpec{Provider: "net1.default"},
			},
		},
		{
			name: "removes stale key from non-ovn subnet",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						util.LogicalSwitchAnnotation: "underlay-subnet",
						util.AllocatedAnnotation:     "true",
						util.TunnelKeyAnnotation:     "7",
					},
				},
			},
			subnet: &kubeovnv1.Subnet{
				ObjectMeta: metav1.ObjectMeta{Name: "underlay-subnet"},
				Spec:       kubeovnv1.SubnetSpec{Provider: "underlay.default"},
				Status:     kubeovnv1.SubnetStatus{TunnelKey: 7},
			},
		},
		{name: "skips deleting pod", pod: deletingPod, subnet: ovnSubnet, vpcTunnelKey: 1234},
		{name: "skips hostNetwork pod", pod: hostNetworkPod, subnet: ovnSubnet, vpcTunnelKey: 1234},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				Namespaces: []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
				Subnets:    []*kubeovnv1.Subnet{tt.subnet},
				Vpcs:       []*kubeovnv1.Vpc{{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: tt.vpcTunnelKey}}},
				Pods:       []*corev1.Pod{tt.pod},
			})
			require.NoError(t, err)
			c := fc.fakeController

			err = c.handleRepairTunnelKey("default/test-pod")
			if tt.wantErr {
				require.Error(t, err)
				updated, getErr := c.config.KubeClient.CoreV1().Pods("default").Get(context.Background(), "test-pod", metav1.GetOptions{})
				require.NoError(t, getErr)
				_, ok := updated.Annotations[util.TunnelKeyAnnotation]
				require.False(t, ok, "invalid VPC tunnel key must never be patched")
				return
			}
			require.NoError(t, err)

			updated, err := c.config.KubeClient.CoreV1().Pods("default").Get(context.Background(), "test-pod", metav1.GetOptions{})
			require.NoError(t, err)
			if tt.wantAnnotated {
				require.Equal(t, tt.wantValue, updated.Annotations[util.TunnelKeyAnnotation])
			} else {
				_, ok := updated.Annotations[util.TunnelKeyAnnotation]
				require.False(t, ok, "tunnel_key annotation must not be added")
			}
		})
	}

	t.Run("pod not found is a no-op", func(t *testing.T) {
		fc, err := newFakeControllerWithOptions(t, nil)
		require.NoError(t, err)
		require.NoError(t, fc.fakeController.handleRepairTunnelKey("default/does-not-exist"))
	})
}

func TestHandleRepairTunnelKeyMultiNIC(t *testing.T) {
	subnets := []*kubeovnv1.Subnet{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "ovn-a"},
			Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "vpc-a"},
			Status:     kubeovnv1.SubnetStatus{TunnelKey: 42},
		},
		{
			// The subnet key is synced but its VPC key is not: the NIC must still wait.
			ObjectMeta: metav1.ObjectMeta{Name: "ovn-b"},
			Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "vpc-b"},
			Status:     kubeovnv1.SubnetStatus{TunnelKey: 43},
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "multi-nic",
			Namespace: "default",
			Annotations: map[string]string{
				nadv1.NetworkAttachmentAnnot: "default/net1, default/net2, default/net3, default/net4",
				// NIC 1: allocated on ovn-a, missing tunnel_key -> patched
				util.AllocatedAnnotation:     "true",
				util.LogicalSwitchAnnotation: "ovn-a",
				// NIC 2: allocated on ovn-b whose VPC key is not synced -> requeued, not patched
				"net1.ovn.kubernetes.io/allocated":      "true",
				"net1.ovn.kubernetes.io/logical_switch": "ovn-b",
				// NIC 3: allocated non-OVN NIC with no logical_switch -> silent skip, not counted
				"net2.kubernetes.io/allocated": "true",
				// NIC 4: logical_switch points to a subnet that does not exist -> skipped and counted
				"net3.ovn.kubernetes.io/allocated":      "true",
				"net3.ovn.kubernetes.io/logical_switch": "ovn-c",
				// NIC 5: has logical_switch but is NOT allocated -> skipped
				"net4.ovn.kubernetes.io/logical_switch": "ovn-a",
			},
		},
	}

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Namespaces: []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
		Subnets:    subnets,
		Vpcs: []*kubeovnv1.Vpc{
			{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: 1234}},
			{ObjectMeta: metav1.ObjectMeta{Name: "vpc-b"}},
		},
		Pods: []*corev1.Pod{pod},
	})
	require.NoError(t, err)
	c := fc.fakeController

	skippedBefore := readCounter(t, metricPodTunnelKeySkipped)
	err = c.handleRepairTunnelKey("default/multi-nic")
	require.Error(t, err, "must requeue while any allocated VPC key is not synced yet")

	updated, err := c.config.KubeClient.CoreV1().Pods("default").Get(context.Background(), "multi-nic", metav1.GetOptions{})
	require.NoError(t, err)
	// NIC 1 must be patched even though NIC 2 is not ready (partial progress lands).
	require.Equal(t, "1234", updated.Annotations[util.TunnelKeyAnnotation])
	// NIC 2 must not be patched with a stale zero key.
	_, ok := updated.Annotations["net1.ovn.kubernetes.io/tunnel_key"]
	require.False(t, ok)
	// NIC 5 (not allocated) must not be patched.
	_, ok = updated.Annotations["net4.ovn.kubernetes.io/tunnel_key"]
	require.False(t, ok)
	// Only the unresolvable logical_switch (NIC 4) counts as skipped; the
	// non-OVN NIC 3 with no logical_switch must be silent.
	require.Equal(t, float64(1), readCounter(t, metricPodTunnelKeySkipped)-skippedBefore)
}

func TestPodProvidersNeedingTunnelKeyRepair(t *testing.T) {
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "ovn-subnet"},
				Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "vpc-a"},
				Status:     kubeovnv1.SubnetStatus{TunnelKey: 42},
			},
			{ObjectMeta: metav1.ObjectMeta{Name: "vlan-subnet"}, Spec: kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vlan: "vlan-a"}},
		},
		Vpcs: []*kubeovnv1.Vpc{{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: 1234}}},
	})
	require.NoError(t, err)
	c := fc.fakeController

	tests := []struct {
		name        string
		annotations map[string]string
		wantRepair  bool
	}{
		{name: "missing VPC key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "ovn-subnet"}, wantRepair: true},
		{name: "matching VPC key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "ovn-subnet", util.TunnelKeyAnnotation: "1234"}},
		{name: "stale VPC key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "ovn-subnet", util.TunnelKeyAnnotation: "999"}, wantRepair: true},
		{name: "legacy subnet key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "ovn-subnet", util.TunnelKeyAnnotation: "42"}, wantRepair: true},
		{name: "invalid VPC key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "ovn-subnet", util.TunnelKeyAnnotation: "invalid"}, wantRepair: true},
		{name: "not allocated", annotations: map[string]string{util.LogicalSwitchAnnotation: "ovn-subnet"}},
		{name: "non-OVN without key", annotations: map[string]string{util.AllocatedAnnotation: "true"}},
		{name: "non-OVN stale key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.TunnelKeyAnnotation: "1234"}, wantRepair: true},
		{name: "vlan without key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "vlan-subnet"}},
		{name: "vlan stale key", annotations: map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "vlan-subnet", util.TunnelKeyAnnotation: "1234"}, wantRepair: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Annotations: tt.annotations}}
			require.Equal(t, tt.wantRepair, len(c.podProvidersNeedingTunnelKeyRepair(pod)) != 0)
		})
	}
}

func TestRepairPodTunnelKeyOnStartup(t *testing.T) {
	pod := func(name, subnet string, phase corev1.PodPhase, sts bool) *corev1.Pod {
		p := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
				Annotations: map[string]string{
					util.AllocatedAnnotation:     "true",
					util.LogicalSwitchAnnotation: subnet,
				},
			},
			Status: corev1.PodStatus{Phase: phase},
		}
		if sts {
			p.OwnerReferences = []metav1.OwnerReference{{Kind: util.KindStatefulSet, Name: "sts", APIVersion: appsv1.SchemeGroupVersion.String()}}
		}
		return p
	}

	tests := []struct {
		name       string
		pod        *corev1.Pod
		tunnelKey  int
		wantKey    string
		wantQueued bool
	}{
		{name: "repairs synchronously when key is ready", pod: pod("ready", "ready-subnet", corev1.PodRunning, false), tunnelKey: 1234, wantKey: "1234"},
		{name: "defers until VPC worker provides key", pod: pod("deferred", "deferred-subnet", corev1.PodRunning, false), wantQueued: true},
		{name: "repairs non-alive STS selected by InitIPAM", pod: pod("sts-0", "ready-subnet", corev1.PodSucceeded, true), tunnelKey: 1234, wantKey: "1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &kubeovnv1.Subnet{
				ObjectMeta: metav1.ObjectMeta{Name: tt.pod.Annotations[util.LogicalSwitchAnnotation]},
				Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "vpc-a"},
				Status:     kubeovnv1.SubnetStatus{TunnelKey: 42},
			}
			vpc := &kubeovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: tt.tunnelKey}}
			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{Subnets: []*kubeovnv1.Subnet{subnet}, Vpcs: []*kubeovnv1.Vpc{vpc}, Pods: []*corev1.Pod{tt.pod}})
			require.NoError(t, err)
			c := fc.fakeController

			providers := c.podProvidersNeedingTunnelKeyRepair(tt.pod)
			require.NotEmpty(t, providers)
			c.repairPodTunnelKeyOnStartup(tt.pod, providers)

			updated, err := c.config.KubeClient.CoreV1().Pods("default").Get(context.Background(), tt.pod.Name, metav1.GetOptions{})
			require.NoError(t, err)
			require.Equal(t, tt.wantKey, updated.Annotations[util.TunnelKeyAnnotation])
			if tt.wantQueued {
				require.Equal(t, 1, c.repairTunnelKeyQueue.Len())
			} else {
				require.Zero(t, c.repairTunnelKeyQueue.Len())
			}
		})
	}
}

func TestAliveSiblingVMPod(t *testing.T) {
	vmiOwner := func(vmName string) []metav1.OwnerReference {
		return []metav1.OwnerReference{
			{
				APIVersion: kubevirtv1.SchemeGroupVersion.String(),
				Kind:       util.KindVirtualMachineInstance,
				Name:       vmName,
			},
		}
	}
	vmPod := func(name, vmName string, phase corev1.PodPhase, deleted bool) *corev1.Pod {
		p := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:       "ns",
				Name:            name,
				OwnerReferences: vmiOwner(vmName),
			},
			Status: corev1.PodStatus{Phase: phase},
		}
		if deleted {
			now := metav1.Now()
			p.DeletionTimestamp = &now
			grace := int64(0)
			p.DeletionGracePeriodSeconds = &grace
		}
		return p
	}

	tests := []struct {
		name           string
		pods           []*corev1.Pod
		vmName         string
		excludePodName string
		wantName       string
	}{
		{
			name:           "no siblings",
			pods:           []*corev1.Pod{vmPod("virt-launcher-vm-aaa", "vm", corev1.PodRunning, false)},
			vmName:         "vm",
			excludePodName: "virt-launcher-vm-aaa",
			wantName:       "",
		},
		{
			name: "alive sibling exists",
			pods: []*corev1.Pod{
				vmPod("virt-launcher-vm-aaa", "vm", corev1.PodSucceeded, true),
				vmPod("virt-launcher-vm-bbb", "vm", corev1.PodRunning, false),
			},
			vmName:         "vm",
			excludePodName: "virt-launcher-vm-aaa",
			wantName:       "virt-launcher-vm-bbb",
		},
		{
			name: "only completed siblings",
			pods: []*corev1.Pod{
				vmPod("virt-launcher-vm-aaa", "vm", corev1.PodSucceeded, true),
				vmPod("virt-launcher-vm-bbb", "vm", corev1.PodSucceeded, false),
			},
			vmName:         "vm",
			excludePodName: "virt-launcher-vm-aaa",
			wantName:       "",
		},
		{
			name: "sibling belongs to different vm",
			pods: []*corev1.Pod{
				vmPod("virt-launcher-other-xxx", "other", corev1.PodRunning, false),
			},
			vmName:         "vm",
			excludePodName: "virt-launcher-vm-aaa",
			wantName:       "",
		},
		{
			name: "non-vm pod ignored",
			pods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: "plain-pod"},
					Status:     corev1.PodStatus{Phase: corev1.PodRunning},
				},
			},
			vmName:         "vm",
			excludePodName: "virt-launcher-vm-aaa",
			wantName:       "",
		},
		{
			name: "excluded pod ignored even when alive",
			pods: []*corev1.Pod{
				vmPod("virt-launcher-vm-aaa", "vm", corev1.PodRunning, false),
			},
			vmName:         "vm",
			excludePodName: "virt-launcher-vm-aaa",
			wantName:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aliveSiblingVMPod(tt.pods, tt.vmName, tt.excludePodName)
			if tt.wantName == "" {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tt.wantName, got.Name)
		})
	}
}

func TestClearRoutedPatch(t *testing.T) {
	ovnNet := func(provider string) *kubeovnNet {
		return &kubeovnNet{
			ProviderName: provider,
			Subnet:       &kubeovnv1.Subnet{Spec: kubeovnv1.SubnetSpec{Provider: util.OvnProvider}},
		}
	}
	attachNet := &kubeovnNet{
		ProviderName: "macvlan.default",
		Subnet:       &kubeovnv1.Subnet{Spec: kubeovnv1.SubnetSpec{Provider: "macvlan.default"}},
	}
	podWith := func(annotations map[string]string) *corev1.Pod {
		return &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: "p", Annotations: annotations}}
	}

	routedKey := fmt.Sprintf(util.RoutedAnnotationTemplate, util.OvnProvider)
	attachRoutedKey := fmt.Sprintf(util.RoutedAnnotationTemplate, attachNet.ProviderName)

	// A routed OVN network is reset so its reconcile runs again.
	patch := clearRoutedPatch(podWith(map[string]string{routedKey: "true"}), []*kubeovnNet{ovnNet(util.OvnProvider)})
	require.Len(t, patch, 1)
	value, ok := patch[routedKey]
	require.True(t, ok)
	assert.Nil(t, value, "the annotation must be removed, not set to a value")

	// A pod that has not been routed yet has a pending reconcile already.
	assert.Empty(t, clearRoutedPatch(podWith(nil), []*kubeovnNet{ovnNet(util.OvnProvider)}))

	// A non-OVN attachment has no port group of ours, so it is left alone.
	assert.Empty(t, clearRoutedPatch(podWith(map[string]string{attachRoutedKey: "true"}), []*kubeovnNet{attachNet}))
}

func TestAlignVMIPNodeName(t *testing.T) {
	newIP := func(name, nodeName string, labels map[string]string) *kubeovnv1.IP {
		return &kubeovnv1.IP{
			ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels},
			Spec:       kubeovnv1.IPSpec{PodName: "vm", Namespace: "ns", NodeName: nodeName},
		}
	}
	ports := []ovnnb.LogicalSwitchPort{{Name: "vm.ns.ovn"}, {Name: "vm-attach.ns.ovn"}, {Name: "missing.ns.ovn"}}

	fakeCtrl, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IPs: []*kubeovnv1.IP{
			newIP("vm.ns.ovn", "old-node", map[string]string{util.NodeNameLabel: "old-node"}),
			newIP("vm-attach.ns.ovn", "target-node", map[string]string{util.NodeNameLabel: "target-node"}),
		},
	})
	require.NoError(t, err)
	ctrl := fakeCtrl.fakeController

	require.NoError(t, ctrl.alignVMIPNodeName(ports, "source-node"))

	for _, name := range []string{"vm.ns.ovn", "vm-attach.ns.ovn"} {
		ip, err := ctrl.config.KubeOvnClient.KubeovnV1().IPs().Get(context.Background(), name, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, "source-node", ip.Spec.NodeName, "%s must point at the node running the vm", name)
		assert.Equal(t, "source-node", ip.Labels[util.NodeNameLabel], "%s label must agree with the spec", name)
	}

	// An IP that is already aligned is left untouched, and a missing one is not an error.
	require.NoError(t, ctrl.alignVMIPNodeName(ports, "source-node"))
}

func TestAliveVMSiblingOf(t *testing.T) {
	vmiOwner := func(vmName string) []metav1.OwnerReference {
		return []metav1.OwnerReference{{
			APIVersion: kubevirtv1.SchemeGroupVersion.String(),
			Kind:       util.KindVirtualMachineInstance,
			Name:       vmName,
		}}
	}
	vmPod := func(name, vmName string, phase corev1.PodPhase, createdBy string) *corev1.Pod {
		labels := map[string]string{"vm.kubevirt.io/name": vmName}
		if createdBy != "" {
			labels[kubevirtv1.CreatedByLabel] = createdBy
		}
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:       "ns",
				Name:            name,
				Labels:          labels,
				OwnerReferences: vmiOwner(vmName),
			},
			Status: corev1.PodStatus{Phase: phase},
		}
	}

	const uid = "6c1d1b0a-0000-4000-8000-000000000001"

	t.Run("labelled pods are matched by the vmi label", func(t *testing.T) {
		dead := vmPod("vm-aaa", "vm", corev1.PodSucceeded, uid)
		alive := vmPod("vm-bbb", "vm", corev1.PodRunning, uid)
		// A pod of another namespace user carries no created-by label and must not be picked up.
		other := vmPod("other-ccc", "other", corev1.PodRunning, "")
		fakeCtrl, err := newFakeControllerWithOptions(t, &FakeControllerOptions{Pods: []*corev1.Pod{dead, alive, other}})
		require.NoError(t, err)

		sibling, err := fakeCtrl.fakeController.aliveVMSiblingOf(dead, "vm")
		require.NoError(t, err)
		require.NotNil(t, sibling)
		assert.Equal(t, "vm-bbb", sibling.Name)
	})

	t.Run("a pod without the label still finds its sibling", func(t *testing.T) {
		dead := vmPod("vm-aaa", "vm", corev1.PodSucceeded, "")
		alive := vmPod("vm-bbb", "vm", corev1.PodRunning, uid)
		fakeCtrl, err := newFakeControllerWithOptions(t, &FakeControllerOptions{Pods: []*corev1.Pod{dead, alive}})
		require.NoError(t, err)

		sibling, err := fakeCtrl.fakeController.aliveVMSiblingOf(dead, "vm")
		require.NoError(t, err)
		require.NotNil(t, sibling, "a missing label must fall back to the whole namespace, never to \"no sibling\"")
		assert.Equal(t, "vm-bbb", sibling.Name)
	})

	t.Run("no alive sibling", func(t *testing.T) {
		dead := vmPod("vm-aaa", "vm", corev1.PodSucceeded, uid)
		fakeCtrl, err := newFakeControllerWithOptions(t, &FakeControllerOptions{Pods: []*corev1.Pod{dead}})
		require.NoError(t, err)

		sibling, err := fakeCtrl.fakeController.aliveVMSiblingOf(dead, "vm")
		require.NoError(t, err)
		assert.Nil(t, sibling)
	})
}

func TestResetVMPodRouted(t *testing.T) {
	routedKey := fmt.Sprintf(util.RoutedAnnotationTemplate, util.OvnProvider)
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ovn-default"},
		Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "ovn-cluster"},
	}
	newPod := func() *corev1.Pod {
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns",
				Name:      "vm-bbb",
				Annotations: map[string]string{
					util.LogicalSwitchAnnotation:                                    "ovn-default",
					fmt.Sprintf(util.AllocatedAnnotationTemplate, util.OvnProvider): "true",
					routedKey: "true",
				},
			},
			Status: corev1.PodStatus{Phase: corev1.PodRunning},
		}
	}

	t.Run("the routed annotation is dropped", func(t *testing.T) {
		fakeCtrl, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			Subnets: []*kubeovnv1.Subnet{subnet},
			Pods:    []*corev1.Pod{newPod()},
		})
		require.NoError(t, err)
		ctrl := fakeCtrl.fakeController

		alive, err := ctrl.resetVMPodRouted(newPod())
		require.NoError(t, err)
		assert.True(t, alive)

		patched, err := ctrl.config.KubeClient.CoreV1().Pods("ns").Get(context.Background(), "vm-bbb", metav1.GetOptions{})
		require.NoError(t, err)
		assert.NotContains(t, patched.Annotations, routedKey, "the routed annotation must be removed, not blanked")
	})

	t.Run("a pod that is already gone is reported", func(t *testing.T) {
		fakeCtrl, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			Subnets: []*kubeovnv1.Subnet{subnet},
		})
		require.NoError(t, err)
		ctrl := fakeCtrl.fakeController

		// The pod exists only in the informer cache: the patch reports NotFound, which is how the
		// caller learns that the sibling it picked is gone.
		alive, err := ctrl.resetVMPodRouted(newPod())
		require.NoError(t, err)
		assert.False(t, alive, "a vanished sibling must be reported so its node is not used")
	})
}

// TestHandleDeletePodVMPortGroups covers the branch decision this feature lives in: a VM pod that
// is deleted while another virt-launcher pod of the same VMI is alive must not strip the port groups
// of the shared LSP, and must hand the memberships to that pod; without a live sibling it still
// removes them. gomock fails on the unexpected call, so "RemovePortFromPortGroups is not called" is
// a hard assertion here.
func TestHandleDeletePodVMPortGroups(t *testing.T) {
	const (
		namespace = "ns"
		vmName    = "vm"
		vmiUID    = "6c1d1b0a-0000-4000-8000-000000000001"
	)
	portName := ovs.PodNameToPortName(vmName, namespace, util.OvnProvider)
	routedKey := fmt.Sprintf(util.RoutedAnnotationTemplate, util.OvnProvider)

	vmPodOwnerRef := []metav1.OwnerReference{{
		APIVersion: kubevirtv1.SchemeGroupVersion.String(),
		Kind:       util.KindVirtualMachineInstance,
		Name:       vmName,
	}}
	vmPod := func(name, uid, node string, phase corev1.PodPhase, deleting bool) *corev1.Pod {
		p := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:       namespace,
				Name:            name,
				UID:             types.UID(uid),
				Labels:          map[string]string{kubevirtv1.CreatedByLabel: vmiUID},
				OwnerReferences: vmPodOwnerRef,
				Annotations: map[string]string{
					util.LogicalSwitchAnnotation:                         "ovn-default",
					util.AllocatedAnnotation:                             "true",
					fmt.Sprintf(util.IPAddressAnnotationTemplate, "ovn"): "10.0.0.2",
					routedKey: "true",
				},
			},
			Spec:   corev1.PodSpec{NodeName: node, RestartPolicy: corev1.RestartPolicyNever},
			Status: corev1.PodStatus{Phase: phase},
		}
		if deleting {
			now := metav1.Now()
			p.DeletionTimestamp = &now
		}
		return p
	}

	newController := func(t *testing.T, pods []*corev1.Pod) (*fakeController, *kubeovnv1.IP) {
		t.Helper()
		ip := &kubeovnv1.IP{
			// createOrUpdateIPCR writes the node into both, so the fixture carries both.
			ObjectMeta: metav1.ObjectMeta{Name: portName, Labels: map[string]string{util.NodeNameLabel: "node-a"}},
			Spec: kubeovnv1.IPSpec{
				PodName:   vmName,
				Namespace: namespace,
				Subnet:    "ovn-default",
				NodeName:  "node-a",
			},
		}
		fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
			Namespaces: []*corev1.Namespace{{
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{util.LogicalSwitchAnnotation: "ovn-default"},
				},
			}},
			Subnets: []*kubeovnv1.Subnet{{
				ObjectMeta: metav1.ObjectMeta{Name: "ovn-default"},
				Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "ovn-cluster"},
			}},
			Pods: pods,
			IPs:  []*kubeovnv1.IP{ip},
		})
		require.NoError(t, err)
		c := fc.fakeController
		c.config.EnableKeepVMIP = true

		kvCtrl := gomock.NewController(t)
		kvClient := kubecli.NewMockKubevirtClient(kvCtrl)
		vmiIface := kubecli.NewMockVirtualMachineInstanceInterface(kvCtrl)
		vmIface := kubecli.NewMockVirtualMachineInterface(kvCtrl)
		kvClient.EXPECT().VirtualMachineInstance(namespace).Return(vmiIface).AnyTimes()
		vmiIface.EXPECT().Get(gomock.Any(), vmName, gomock.Any()).
			Return(&kubevirtv1.VirtualMachineInstance{}, nil).AnyTimes()
		kvClient.EXPECT().VirtualMachine(namespace).Return(vmIface).AnyTimes()
		vmIface.EXPECT().Get(gomock.Any(), vmName, gomock.Any()).
			Return(&kubevirtv1.VirtualMachine{}, nil).AnyTimes()
		c.config.KubevirtClient = kvClient

		return fc, ip
	}

	t.Run("a live sibling keeps the port groups and its own node", func(t *testing.T) {
		deleting := vmPod("virt-launcher-vm-aaa", "uid-a", "node-a", corev1.PodSucceeded, true)
		sibling := vmPod("virt-launcher-vm-bbb", "uid-b", "node-b", corev1.PodRunning, false)
		fc, _ := newController(t, []*corev1.Pod{deleting, sibling})
		c := fc.fakeController

		fc.mockOvnClient.EXPECT().
			ListNormalLogicalSwitchPorts(gomock.Any(), map[string]string{"pod": namespace + "/" + vmName}).
			Return([]ovnnb.LogicalSwitchPort{{Name: portName}}, nil).Times(2)
		fc.mockOvnClient.EXPECT().CleanLogicalSwitchPortMigrateOptions(portName).Return(nil)

		key := namespace + "/" + deleting.Name
		c.deletingPodObjMap.Store(key, deleting)
		require.NoError(t, c.handleDeletePod(key))

		patched, err := c.config.KubeClient.CoreV1().Pods(namespace).Get(context.Background(), sibling.Name, metav1.GetOptions{})
		require.NoError(t, err)
		assert.NotContains(t, patched.Annotations, routedKey, "the surviving pod must be told to re-assert its memberships")

		ip, err := c.config.KubeOvnClient.KubeovnV1().IPs().Get(context.Background(), portName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, "node-b", ip.Spec.NodeName, "the IP must follow the node that runs the vm")
		assert.Equal(t, "node-b", ip.Labels[util.NodeNameLabel])
	})

	t.Run("a sibling that is not placed yet leaves the records alone", func(t *testing.T) {
		deleting := vmPod("virt-launcher-vm-aaa", "uid-a", "node-a", corev1.PodSucceeded, true)
		pending := vmPod("virt-launcher-vm-bbb", "uid-b", "", corev1.PodPending, false)
		fc, _ := newController(t, []*corev1.Pod{deleting, pending})
		c := fc.fakeController

		fc.mockOvnClient.EXPECT().
			ListNormalLogicalSwitchPorts(gomock.Any(), map[string]string{"pod": namespace + "/" + vmName}).
			Return([]ovnnb.LogicalSwitchPort{{Name: portName}}, nil).Times(2)
		fc.mockOvnClient.EXPECT().CleanLogicalSwitchPortMigrateOptions(portName).Return(nil)

		key := namespace + "/" + deleting.Name
		c.deletingPodObjMap.Store(key, deleting)
		require.NoError(t, c.handleDeletePod(key))

		ip, err := c.config.KubeOvnClient.KubeovnV1().IPs().Get(context.Background(), portName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, "node-a", ip.Spec.NodeName, "a pod without a node must not blank the record")
		assert.Equal(t, "node-a", ip.Labels[util.NodeNameLabel])
	})

	t.Run("without a live sibling the memberships are removed", func(t *testing.T) {
		deleting := vmPod("virt-launcher-vm-aaa", "uid-a", "node-a", corev1.PodSucceeded, true)
		fc, ip := newController(t, []*corev1.Pod{deleting})
		c := fc.fakeController

		fc.mockOvnClient.EXPECT().
			ListNormalLogicalSwitchPorts(gomock.Any(), map[string]string{"pod": namespace + "/" + vmName}).
			Return([]ovnnb.LogicalSwitchPort{{Name: portName}}, nil).Times(2)
		fc.mockOvnClient.EXPECT().CleanLogicalSwitchPortMigrateOptions(portName).Return(nil)
		fc.mockOvnClient.EXPECT().RemovePortFromPortGroups(portName).Return(nil)

		key := namespace + "/" + deleting.Name
		c.deletingPodObjMap.Store(key, deleting)
		require.NoError(t, c.handleDeletePod(key))

		assert.Equal(t, "node-a", ip.Spec.NodeName, "with no live pod the record is left alone")
	})
}

// newNatGwInitQueue wires the init queue the NAT gateway pod handlers enqueue into.
func newNatGwInitQueue(t *testing.T, c *Controller) workqueue.TypedRateLimitingInterface[string] {
	t.Helper()
	queue := newTypedRateLimitingQueue[string]("InitVpcNatGateway", nil)
	t.Cleanup(queue.ShutDown)
	c.initVpcNatGatewayQueue = queue
	return queue
}

func requireQueueItem(t *testing.T, queue workqueue.TypedRateLimitingInterface[string], want string) {
	t.Helper()
	require.Equal(t, 1, queue.Len())
	item, shutdown := queue.Get()
	require.False(t, shutdown)
	queue.Done(item)
	require.Equal(t, want, item)
}

// TestHandleAddOrUpdatePodQueuesNatGwInit covers the fast resync path: a NAT gateway pod that
// already owns its addresses never reaches reconcileAllocateSubnets, so the add/update handler
// itself must queue the missing initialization, and only for a pod that still needs it.
func TestHandleAddOrUpdatePodQueuesNatGwInit(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ovn-subnet"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock: "10.0.1.0/24",
			Protocol:  kubeovnv1.ProtocolIPv4,
			Provider:  util.OvnProvider,
		},
		Status: kubeovnv1.SubnetStatus{V4AvailableIPs: 100, TunnelKey: 1234},
	}

	tests := []struct {
		name        string
		phase       corev1.PodPhase
		initialized bool
		terminating bool
		want        int
	}{
		{name: "running pod without the init annotation is queued", phase: corev1.PodRunning, want: 1},
		{name: "initialized pod is skipped", phase: corev1.PodRunning, initialized: true},
		{name: "pod that is not running yet is skipped", phase: corev1.PodPending},
		{name: "terminating pod is skipped", phase: corev1.PodRunning, terminating: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := fakeNatGwPod("gw", tt.initialized)
			pod.Namespace = "default"
			pod.Spec.NodeName = "node-1"
			pod.Status.Phase = tt.phase
			if tt.terminating {
				now := metav1.Now()
				pod.DeletionTimestamp = &now
			}
			pod.Annotations[util.LogicalSwitchAnnotation] = subnet.Name
			pod.Annotations[util.AllocatedAnnotation] = "true"
			pod.Annotations[util.IPAddressAnnotation] = "10.0.1.10"
			pod.Annotations[util.MacAddressAnnotation] = "00:00:00:00:00:10"
			pod.Annotations[util.RoutedAnnotation] = "true"

			fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
				Namespaces: []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
				Subnets:    []*kubeovnv1.Subnet{subnet},
				Pods:       []*corev1.Pod{pod},
			})
			require.NoError(t, err)
			c := fc.fakeController
			queue := newNatGwInitQueue(t, c)
			fc.mockOvnClient.EXPECT().ListNormalLogicalSwitchPorts(gomock.Any(), gomock.Any()).Return(nil, nil)

			require.NoError(t, c.handleAddOrUpdatePod("default/"+pod.Name))
			if tt.want == 0 {
				require.Zero(t, queue.Len())
				return
			}
			requireQueueItem(t, queue, "gw")
		})
	}
}
