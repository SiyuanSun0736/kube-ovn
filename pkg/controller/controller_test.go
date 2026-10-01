package controller

// Unified Fake Controller for Testing
//
// This file provides a unified approach to creating fake controllers for testing.
// The main function is newFakeControllerWithOptions() which accepts optional parameters
// for subnets, NADs (Network Attachment Definitions), pods, and namespaces.
//
// The fake controller properly initializes:
// - Kubernetes fake client with pods and namespaces
// - NAD fake client with network attachment definitions (populated via API)
// - KubeOVN fake client with subnets (populated via API)
// - All necessary informers with proper synchronization
// - Mock OVN client for OVN operations

import (
	"context"
	"testing"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	nadfake "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned/fake"
	nadinformers "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/informers/externalversions"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/keymutex"

	mockovs "github.com/kubeovn/kube-ovn/mocks/pkg/ovs"
	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	kubeovnfake "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned/fake"
	kubeovninformerfactory "github.com/kubeovn/kube-ovn/pkg/client/informers/externalversions"
	kubeovninformer "github.com/kubeovn/kube-ovn/pkg/client/informers/externalversions/kubeovn/v1"
	ovnipam "github.com/kubeovn/kube-ovn/pkg/ipam"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

type fakeControllerInformers struct {
	vpcInformer       kubeovninformer.VpcInformer
	vpcNatGwInformer  kubeovninformer.VpcNatGatewayInformer
	subnetInformer    kubeovninformer.SubnetInformer
	ipInformer        kubeovninformer.IPInformer
	vlanInformer      kubeovninformer.VlanInformer
	serviceInformer   coreinformers.ServiceInformer
	namespaceInformer coreinformers.NamespaceInformer
	podInformer       coreinformers.PodInformer
	nodeInformer      coreinformers.NodeInformer
}

type fakeController struct {
	fakeController  *Controller
	fakeInformers   *fakeControllerInformers
	mockOvnClient   *mockovs.MockNbClient
	mockOvnSbClient *mockovs.MockSbClient
}

func alwaysReady() bool { return true }

// FakeControllerOptions holds optional parameters for creating a fake controller
type FakeControllerOptions struct {
	Subnets            []*kubeovnv1.Subnet
	VpcNatGateways     []*kubeovnv1.VpcNatGateway
	QoSPolicies        []*kubeovnv1.QoSPolicy
	IptablesEIPs       []*kubeovnv1.IptablesEIP
	IptablesFIPs       []*kubeovnv1.IptablesFIPRule
	IptablesDnatRules  []*kubeovnv1.IptablesDnatRule
	IptablesSnatRules  []*kubeovnv1.IptablesSnatRule
	IPs                []*kubeovnv1.IP
	Vlans              []*kubeovnv1.Vlan
	NetworkAttachments []*nadv1.NetworkAttachmentDefinition
	Pods               []*corev1.Pod
	Namespaces         []*corev1.Namespace
	EndpointSlices     []*discoveryv1.EndpointSlice
	ConfigMaps         []*corev1.ConfigMap
	Vpcs               []*kubeovnv1.Vpc
	Services           []*corev1.Service
	StatefulSets       []*appsv1.StatefulSet
}

// newFakeControllerWithOptions creates a fake controller with optional pre-populated objects
func newFakeControllerWithOptions(t *testing.T, opts *FakeControllerOptions) (*fakeController, error) {
	if opts == nil {
		opts = &FakeControllerOptions{}
	}

	// Create default namespace if none provided
	defaultNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				util.LogicalSwitchAnnotation: util.DefaultSubnet,
			},
		},
	}

	namespaces := opts.Namespaces
	if len(namespaces) == 0 {
		namespaces = []*corev1.Namespace{defaultNamespace}
	}

	// Create fake Kubernetes client with namespaces and pods
	kubeObjects := make([]runtime.Object, 0, len(namespaces)+len(opts.Pods))
	for _, ns := range namespaces {
		kubeObjects = append(kubeObjects, ns)
	}
	for _, pod := range opts.Pods {
		kubeObjects = append(kubeObjects, pod)
	}
	for _, slice := range opts.EndpointSlices {
		kubeObjects = append(kubeObjects, slice)
	}
	for _, svc := range opts.Services {
		kubeObjects = append(kubeObjects, svc)
	}
	for _, sts := range opts.StatefulSets {
		kubeObjects = append(kubeObjects, sts)
	}
	for _, cm := range opts.ConfigMaps {
		kubeObjects = append(kubeObjects, cm)
	}
	kubeClient := fake.NewSimpleClientset(kubeObjects...)

	// Create fake NAD client
	nadClient := nadfake.NewSimpleClientset()
	for _, nad := range opts.NetworkAttachments {
		_, err := nadClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(nad.Namespace).Create(
			context.Background(), nad, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}

	// Create fake KubeOVN client
	kubeovnClient := kubeovnfake.NewSimpleClientset()
	for _, vpc := range opts.Vpcs {
		_, err := kubeovnClient.KubeovnV1().Vpcs().Create(
			context.Background(), vpc, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, subnet := range opts.Subnets {
		_, err := kubeovnClient.KubeovnV1().Subnets().Create(
			context.Background(), subnet, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, gw := range opts.VpcNatGateways {
		_, err := kubeovnClient.KubeovnV1().VpcNatGateways().Create(
			context.Background(), gw, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, qos := range opts.QoSPolicies {
		_, err := kubeovnClient.KubeovnV1().QoSPolicies().Create(
			context.Background(), qos, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, eip := range opts.IptablesEIPs {
		_, err := kubeovnClient.KubeovnV1().IptablesEIPs().Create(
			context.Background(), eip, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, fip := range opts.IptablesFIPs {
		_, err := kubeovnClient.KubeovnV1().IptablesFIPRules().Create(
			context.Background(), fip, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, dnat := range opts.IptablesDnatRules {
		_, err := kubeovnClient.KubeovnV1().IptablesDnatRules().Create(
			context.Background(), dnat, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, snat := range opts.IptablesSnatRules {
		_, err := kubeovnClient.KubeovnV1().IptablesSnatRules().Create(
			context.Background(), snat, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, ip := range opts.IPs {
		_, err := kubeovnClient.KubeovnV1().IPs().Create(
			context.Background(), ip, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	for _, vlan := range opts.Vlans {
		_, err := kubeovnClient.KubeovnV1().Vlans().Create(
			context.Background(), vlan, metav1.CreateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}

	// Create informer factories
	kubeInformerFactory := informers.NewSharedInformerFactory(kubeClient, 0)
	serviceInformer := kubeInformerFactory.Core().V1().Services()
	// The gateway nftable LB feature looks Services up by gateway and by EIP through informer
	// indexes, so the fake registers the same indexes the controller does.
	if err := serviceInformer.Informer().AddIndexers(cache.Indexers{
		IndexGwNftableLbServiceByEip:     indexGwNftableLbServiceByEip,
		IndexGwNftableLbServiceByGateway: indexGwNftableLbServiceByGateway,
	}); err != nil {
		return nil, err
	}
	namespaceInformer := kubeInformerFactory.Core().V1().Namespaces()
	podInformer := kubeInformerFactory.Core().V1().Pods()
	nodeInformer := kubeInformerFactory.Core().V1().Nodes()
	configMapInformer := kubeInformerFactory.Core().V1().ConfigMaps()
	endpointSliceInformer := kubeInformerFactory.Discovery().V1().EndpointSlices()

	nadInformerFactory := nadinformers.NewSharedInformerFactory(nadClient, 0)
	nadInformer := nadInformerFactory.K8sCniCncfIo().V1().NetworkAttachmentDefinitions()

	kubeovnInformerFactory := kubeovninformerfactory.NewSharedInformerFactory(kubeovnClient, 0)
	vpcInformer := kubeovnInformerFactory.Kubeovn().V1().Vpcs()
	subnetInformer := kubeovnInformerFactory.Kubeovn().V1().Subnets()
	ipInformer := kubeovnInformerFactory.Kubeovn().V1().IPs()
	vpcNatGwInformer := kubeovnInformerFactory.Kubeovn().V1().VpcNatGateways()
	qosPolicyInformer := kubeovnInformerFactory.Kubeovn().V1().QoSPolicies()
	iptablesEipInformer := kubeovnInformerFactory.Kubeovn().V1().IptablesEIPs()
	iptablesFipInformer := kubeovnInformerFactory.Kubeovn().V1().IptablesFIPRules()
	iptablesDnatRuleInformer := kubeovnInformerFactory.Kubeovn().V1().IptablesDnatRules()
	iptablesSnatRuleInformer := kubeovnInformerFactory.Kubeovn().V1().IptablesSnatRules()
	vlanInformer := kubeovnInformerFactory.Kubeovn().V1().Vlans()

	fakeInformers := &fakeControllerInformers{
		vpcInformer:       vpcInformer,
		vpcNatGwInformer:  vpcNatGwInformer,
		subnetInformer:    subnetInformer,
		ipInformer:        ipInformer,
		vlanInformer:      vlanInformer,
		serviceInformer:   serviceInformer,
		namespaceInformer: namespaceInformer,
		podInformer:       podInformer,
		nodeInformer:      nodeInformer,
	}

	// Create mock OVN clients
	mockCtrl := gomock.NewController(t)
	mockOvnClient := mockovs.NewMockNbClient(mockCtrl)
	mockOvnSbClient := mockovs.NewMockSbClient(mockCtrl)

	// Create controller with all informers
	ctrl := &Controller{
		servicesLister:          serviceInformer.Lister(),
		svcIndexer:              serviceInformer.Informer().GetIndexer(),
		namespacesLister:        namespaceInformer.Lister(),
		podsLister:              podInformer.Lister(),
		endpointSlicesLister:    endpointSliceInformer.Lister(),
		nodesLister:             nodeInformer.Lister(),
		configMapsLister:        configMapInformer.Lister(),
		vpcsLister:              vpcInformer.Lister(),
		vpcSynced:               alwaysReady,
		subnetsLister:           subnetInformer.Lister(),
		subnetSynced:            alwaysReady,
		ipsLister:               ipInformer.Lister(),
		ipSynced:                alwaysReady,
		vlansLister:             vlanInformer.Lister(),
		netAttachLister:         nadInformer.Lister(),
		netAttachSynced:         alwaysReady,
		vpcNatGatewayLister:     vpcNatGwInformer.Lister(),
		qosPoliciesLister:       qosPolicyInformer.Lister(),
		iptablesEipsLister:      iptablesEipInformer.Lister(),
		iptablesFipsLister:      iptablesFipInformer.Lister(),
		iptablesDnatRulesLister: iptablesDnatRuleInformer.Lister(),
		iptablesSnatRulesLister: iptablesSnatRuleInformer.Lister(),
		vpcNatGwKeyMutex:        keymutex.NewHashed(0),
		vpcNatGwExecKeyMutex:    keymutex.NewHashed(0),
		deletingPodObjMap:       xsync.NewMap[string, *corev1.Pod](),
		OVNNbClient:             mockOvnClient,
		OVNSbClient:             mockOvnSbClient,
		ipam:                    ovnipam.NewIPAM(),
		recorder:                record.NewFakeRecorder(100),
		subnetKeyMutex:          keymutex.NewHashed(0),
		podKeyMutex:             keymutex.NewHashed(0),
		addOrUpdateSubnetQueue:  newTypedRateLimitingQueue[string]("AddOrUpdateSubnet", nil),
		repairTunnelKeyQueue:    newTypedRateLimitingQueue[string]("RepairTunnelKey", nil),
		syncVirtualPortsQueue:   newTypedRateLimitingQueue[string]("SyncVirtualPort", nil),
		updateSubnetStatusQueue: newTypedRateLimitingQueue[string]("UpdateSubnetStatus", nil),
	}

	ctrl.config = &Configuration{
		ClusterRouter:        "ovn-cluster",
		DefaultLogicalSwitch: "ovn-default",
		NodeSwitch:           "join",
		KubeOvnClient:        kubeovnClient,
		KubeClient:           kubeClient,
		PodNamespace:         "kube-system",
		AttachNetClient:      nadClient,
	}

	// Start informers and wait for sync
	stopCh := make(chan struct{})
	t.Cleanup(func() { close(stopCh) })

	kubeInformerFactory.Start(stopCh)
	nadInformerFactory.Start(stopCh)
	kubeovnInformerFactory.Start(stopCh)

	kubeInformerFactory.WaitForCacheSync(stopCh)
	nadInformerFactory.WaitForCacheSync(stopCh)
	kubeovnInformerFactory.WaitForCacheSync(stopCh)

	return &fakeController{
		fakeController:  ctrl,
		fakeInformers:   fakeInformers,
		mockOvnClient:   mockOvnClient,
		mockOvnSbClient: mockOvnSbClient,
	}, nil
}

// newFakeController creates a basic fake controller
func newFakeController(t *testing.T) *fakeController {
	controller, err := newFakeControllerWithOptions(t, nil)
	require.NoError(t, err)
	return controller
}

func Test_allSubnetReady(t *testing.T) {
	fakeController, err := newFakeControllerWithOptions(t, nil)
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	mockOvnClient := fakeController.mockOvnClient

	subnets := []string{"ovn-default", "join"}

	t.Run("all subnet ready", func(t *testing.T) {
		mockOvnClient.EXPECT().LogicalSwitchExists(gomock.Any()).Return(true, nil).Times(2)

		ready, err := ctrl.allSubnetReady(subnets...)
		require.NoError(t, err)
		require.True(t, ready)
	})

	t.Run("some subnet are not ready", func(t *testing.T) {
		mockOvnClient.EXPECT().LogicalSwitchExists(subnets[0]).Return(true, nil)
		mockOvnClient.EXPECT().LogicalSwitchExists(subnets[1]).Return(false, nil)

		ready, err := ctrl.allSubnetReady(subnets...)
		require.NoError(t, err)
		require.False(t, ready)
	})
}

// TestFakeControllerWithOptions demonstrates usage of the unified fake controller
func TestFakeControllerWithOptions(t *testing.T) {
	// Example: creating a fake controller with NADs, subnets, and pods
	nads := []*nadv1.NetworkAttachmentDefinition{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "net1",
				Namespace: "default",
			},
			Spec: nadv1.NetworkAttachmentDefinitionSpec{
				Config: `{"cniVersion": "0.3.1", "name": "net1", "type": "kube-ovn"}`,
			},
		},
	}

	subnets := []*kubeovnv1.Subnet{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "net1-subnet"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "192.168.1.0/24"},
		},
	}

	pods := []*corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
				Annotations: map[string]string{
					nadv1.NetworkAttachmentAnnot: `[{"name": "net1"}]`,
				},
			},
		},
	}

	opts := &FakeControllerOptions{
		Subnets:            subnets,
		NetworkAttachments: nads,
		Pods:               pods,
	}

	fakeCtrl, err := newFakeControllerWithOptions(t, opts)
	require.NoError(t, err)
	ctrl := fakeCtrl.fakeController

	// Verify that the fake controller was created successfully
	require.NotNil(t, ctrl)
	require.NotNil(t, ctrl.config)
	require.NotNil(t, ctrl.config.AttachNetClient)
	require.NotNil(t, ctrl.config.KubeOvnClient)

	// Verify that NADs can be retrieved
	nadClient := ctrl.config.AttachNetClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions("default")
	retrievedNAD, err := nadClient.Get(context.Background(), "net1", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "net1", retrievedNAD.Name)

	// Verify that subnets can be retrieved
	subnetClient := ctrl.config.KubeOvnClient.KubeovnV1().Subnets()
	retrievedSubnet, err := subnetClient.Get(context.Background(), "net1-subnet", metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, "net1-subnet", retrievedSubnet.Name)
}
