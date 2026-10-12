package controller

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/keymutex"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/ovsdb/ovnnb"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func Test_handleAddOrUpdateVpc_staticRoutes(t *testing.T) {
	t.Parallel()

	vpcName := "test-vpc"

	// Policy variables for taking pointers
	srcIPPolicy := ovnnb.LogicalRouterStaticRoutePolicySrcIP
	dstIPPolicy := ovnnb.LogicalRouterStaticRoutePolicyDstIP

	// Internal static route created directly in OVN with kube-ovn vendor
	internalStaticRoute := &ovnnb.LogicalRouterStaticRoute{
		UUID: "internal-static-route-uuid",
		ExternalIDs: map[string]string{
			"vendor": util.CniTypeName,
		},
		IPPrefix:   "10.0.0.0/24",
		Nexthop:    "1.2.3.4",
		Policy:     &srcIPPolicy,
		RouteTable: util.MainRouteTable,
	}

	// Static route that matches VPC spec
	managedStaticRoute := &ovnnb.LogicalRouterStaticRoute{
		UUID: "managed-static-route-uuid",
		ExternalIDs: map[string]string{
			"vendor": util.CniTypeName,
		},
		IPPrefix:   "192.168.0.0/24",
		Nexthop:    "10.0.0.1",
		Policy:     &dstIPPolicy,
		RouteTable: util.MainRouteTable,
	}

	t.Run("only try to manage static routes with vendor kube-ovn", func(t *testing.T) {
		fakeController := newFakeController(t)
		ctrl := fakeController.fakeController
		fakeinformers := fakeController.fakeInformers
		mockOvnClient := fakeController.mockOvnClient

		// Initialize mutexes
		ctrl.vpcKeyMutex = keymutex.NewHashed(500)

		vpc := &kubeovnv1.Vpc{
			ObjectMeta: metav1.ObjectMeta{
				Name: vpcName,
			},
			Spec: kubeovnv1.VpcSpec{
				StaticRoutes: []*kubeovnv1.StaticRoute{
					{
						CIDR:       "192.168.0.0/24",
						NextHopIP:  "10.0.0.1",
						Policy:     kubeovnv1.PolicyDst,
						RouteTable: util.MainRouteTable,
					},
				},
				EnableExternal: false,
				PolicyRoutes:   nil,
			},
			Status: kubeovnv1.VpcStatus{
				Subnets:        []string{},
				EnableExternal: false,
			},
		}

		_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpc, metav1.CreateOptions{})
		require.NoError(t, err)

		err = fakeinformers.vpcInformer.Informer().GetStore().Add(vpc)
		require.NoError(t, err)

		existingKubeOvnRoutes := []*ovnnb.LogicalRouterStaticRoute{
			internalStaticRoute,
		}

		externalIDs := map[string]string{"vendor": util.CniTypeName}

		mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
		mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(existingKubeOvnRoutes, nil)
		mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
			Name: vpcName,
			Nat:  []string{},
		}, nil)
		mockOvnClient.EXPECT().DeleteLogicalRouterStaticRoute(vpcName, gomock.Any(), gomock.Any(), "10.0.0.0/24", "1.2.3.4").Return(nil)
		mockOvnClient.EXPECT().AddLogicalRouterStaticRoute(
			vpcName,
			util.MainRouteTable,
			"dst-ip",
			"192.168.0.0/24",
			nil,
			externalIDs,
			"10.0.0.1",
		).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
		mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
		mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
		mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		err = ctrl.handleAddOrUpdateVpc(vpcName)
		require.NoError(t, err)
	})

	t.Run("delete orphaned kube-ovn routes", func(t *testing.T) {
		fakeController := newFakeController(t)
		ctrl := fakeController.fakeController
		fakeinformers := fakeController.fakeInformers
		mockOvnClient := fakeController.mockOvnClient

		ctrl.vpcKeyMutex = keymutex.NewHashed(500)

		vpc := &kubeovnv1.Vpc{
			ObjectMeta: metav1.ObjectMeta{
				Name: vpcName,
			},
			Spec: kubeovnv1.VpcSpec{
				StaticRoutes: []*kubeovnv1.StaticRoute{
					{
						CIDR:       "192.168.0.0/24",
						NextHopIP:  "10.0.0.1",
						Policy:     kubeovnv1.PolicyDst,
						RouteTable: util.MainRouteTable,
					},
				},
				EnableExternal: false,
				PolicyRoutes:   nil,
			},
			Status: kubeovnv1.VpcStatus{
				Subnets:        []string{},
				EnableExternal: false,
			},
		}

		_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpc, metav1.CreateOptions{})
		require.NoError(t, err)

		err = fakeinformers.vpcInformer.Informer().GetStore().Add(vpc)
		require.NoError(t, err)

		existingKubeOvnRoutes := []*ovnnb.LogicalRouterStaticRoute{
			internalStaticRoute,
			managedStaticRoute,
		}

		externalIDs := map[string]string{"vendor": util.CniTypeName}

		mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
		mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(existingKubeOvnRoutes, nil)
		mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
			Name: vpcName,
			Nat:  []string{},
		}, nil)
		mockOvnClient.EXPECT().DeleteLogicalRouterStaticRoute(vpcName, gomock.Any(), gomock.Any(), "10.0.0.0/24", "1.2.3.4").Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
		mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
		mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
		mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		err = ctrl.handleAddOrUpdateVpc(vpcName)
		require.NoError(t, err)
	})

	t.Run("handle empty VPC static routes", func(t *testing.T) {
		fakeController := newFakeController(t)
		ctrl := fakeController.fakeController
		fakeinformers := fakeController.fakeInformers
		mockOvnClient := fakeController.mockOvnClient

		ctrl.vpcKeyMutex = keymutex.NewHashed(500)

		vpcEmpty := &kubeovnv1.Vpc{
			ObjectMeta: metav1.ObjectMeta{
				Name: vpcName,
			},
			Spec: kubeovnv1.VpcSpec{
				StaticRoutes:   []*kubeovnv1.StaticRoute{},
				EnableExternal: false,
				PolicyRoutes:   nil,
			},
			Status: kubeovnv1.VpcStatus{
				Subnets:        []string{},
				EnableExternal: false,
			},
		}

		_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpcEmpty, metav1.CreateOptions{})
		require.NoError(t, err)

		err = fakeinformers.vpcInformer.Informer().GetStore().Add(vpcEmpty)
		require.NoError(t, err)

		existingKubeOvnRoutes := []*ovnnb.LogicalRouterStaticRoute{
			internalStaticRoute,
			managedStaticRoute,
		}

		externalIDs := map[string]string{"vendor": util.CniTypeName}

		mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
		mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(existingKubeOvnRoutes, nil)
		mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
			Name: vpcName,
			Nat:  []string{},
		}, nil)
		mockOvnClient.EXPECT().DeleteLogicalRouterStaticRoute(vpcName, gomock.Any(), gomock.Any(), "10.0.0.0/24", "1.2.3.4").Return(nil)
		mockOvnClient.EXPECT().DeleteLogicalRouterStaticRoute(vpcName, gomock.Any(), gomock.Any(), "192.168.0.0/24", "10.0.0.1").Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
		mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
		mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
		mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		err = ctrl.handleAddOrUpdateVpc(vpcName)
		require.NoError(t, err)
	})

	t.Run("add static routes from VPC spec when none exist", func(t *testing.T) {
		fakeController := newFakeController(t)
		ctrl := fakeController.fakeController
		fakeinformers := fakeController.fakeInformers
		mockOvnClient := fakeController.mockOvnClient

		ctrl.vpcKeyMutex = keymutex.NewHashed(500)

		vpc := &kubeovnv1.Vpc{
			ObjectMeta: metav1.ObjectMeta{
				Name: vpcName,
			},
			Spec: kubeovnv1.VpcSpec{
				StaticRoutes: []*kubeovnv1.StaticRoute{
					{
						CIDR:       "192.168.0.0/24",
						NextHopIP:  "10.0.0.1",
						Policy:     kubeovnv1.PolicyDst,
						RouteTable: util.MainRouteTable,
					},
				},
				EnableExternal: false,
				PolicyRoutes:   nil,
			},
			Status: kubeovnv1.VpcStatus{
				Subnets:        []string{},
				EnableExternal: false,
			},
		}

		_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpc, metav1.CreateOptions{})
		require.NoError(t, err)

		err = fakeinformers.vpcInformer.Informer().GetStore().Add(vpc)
		require.NoError(t, err)

		externalIDs := map[string]string{"vendor": util.CniTypeName}

		mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
		mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(nil, nil)
		mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
			Name: vpcName,
			Nat:  []string{},
		}, nil)
		mockOvnClient.EXPECT().AddLogicalRouterStaticRoute(
			vpcName,
			util.MainRouteTable,
			"dst-ip",
			"192.168.0.0/24",
			nil,
			externalIDs,
			"10.0.0.1",
		).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
		mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
		mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
		mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		err = ctrl.handleAddOrUpdateVpc(vpcName)
		require.NoError(t, err)
	})
}

func Test_handleAddOrUpdateVpc_policyRoutes_ecmpNextHops(t *testing.T) {
	t.Parallel()

	vpcName := "test-vpc-policy"

	t.Run("empty NextHopIP with non-reroute action should pass nil next-hops", func(t *testing.T) {
		fakeController := newFakeController(t)
		ctrl := fakeController.fakeController
		fakeinformers := fakeController.fakeInformers
		mockOvnClient := fakeController.mockOvnClient

		ctrl.vpcKeyMutex = keymutex.NewHashed(500)

		vpc := &kubeovnv1.Vpc{
			ObjectMeta: metav1.ObjectMeta{
				Name: vpcName,
			},
			Spec: kubeovnv1.VpcSpec{
				StaticRoutes:   []*kubeovnv1.StaticRoute{},
				EnableExternal: false,
				PolicyRoutes: []*kubeovnv1.PolicyRoute{
					{
						Priority:  200,
						Match:     "ip4.src == 10.0.0.0/8",
						Action:    kubeovnv1.PolicyRouteActionDrop,
						NextHopIP: "",
					},
				},
			},
			Status: kubeovnv1.VpcStatus{
				Subnets:        []string{},
				EnableExternal: false,
			},
		}

		_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpc, metav1.CreateOptions{})
		require.NoError(t, err)

		err = fakeinformers.vpcInformer.Informer().GetStore().Add(vpc)
		require.NoError(t, err)

		externalIDs := map[string]string{"vendor": util.CniTypeName}

		mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
		mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(nil, nil)
		mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
			Name: vpcName,
			Nat:  []string{},
		}, nil)
		mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
		// The key assertion: empty NextHopIP must produce nil next-hops, not [""]
		mockOvnClient.EXPECT().AddLogicalRouterPolicy(
			vpcName,
			200,
			"ip4.src == 10.0.0.0/8",
			string(kubeovnv1.PolicyRouteActionDrop),
			([]string)(nil),
			([]string)(nil),
			externalIDs,
		).Return(nil)
		mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
		mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
		mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		err = ctrl.handleAddOrUpdateVpc(vpcName)
		require.NoError(t, err)
	})

	t.Run("ECMP next-hops are split correctly for custom VPC policy routes", func(t *testing.T) {
		fakeController := newFakeController(t)
		ctrl := fakeController.fakeController
		fakeinformers := fakeController.fakeInformers
		mockOvnClient := fakeController.mockOvnClient

		ctrl.vpcKeyMutex = keymutex.NewHashed(500)

		vpc := &kubeovnv1.Vpc{
			ObjectMeta: metav1.ObjectMeta{
				Name: vpcName,
			},
			Spec: kubeovnv1.VpcSpec{
				StaticRoutes:   []*kubeovnv1.StaticRoute{},
				EnableExternal: false,
				PolicyRoutes: []*kubeovnv1.PolicyRoute{
					{
						Priority:  100,
						Match:     "ip4.dst == 10.1.0.0/16",
						Action:    kubeovnv1.PolicyRouteActionReroute,
						NextHopIP: "192.168.1.1,192.168.1.2",
					},
				},
			},
			Status: kubeovnv1.VpcStatus{
				Subnets:        []string{},
				EnableExternal: false,
			},
		}

		_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpc, metav1.CreateOptions{})
		require.NoError(t, err)

		err = fakeinformers.vpcInformer.Informer().GetStore().Add(vpc)
		require.NoError(t, err)

		externalIDs := map[string]string{"vendor": util.CniTypeName}

		mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
		mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
		mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(nil, nil)
		mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
			Name: vpcName,
			Nat:  []string{},
		}, nil)
		// No existing policies in OVN for this custom VPC
		mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
		// The key assertion: next-hops must be split into a slice, not wrapped as one element
		mockOvnClient.EXPECT().AddLogicalRouterPolicy(
			vpcName,
			100,
			"ip4.dst == 10.1.0.0/16",
			string(kubeovnv1.PolicyRouteActionReroute),
			[]string{"192.168.1.1", "192.168.1.2"},
			([]string)(nil),
			externalIDs,
		).Return(nil)
		mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
		mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
		mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
		err = ctrl.handleAddOrUpdateVpc(vpcName)
		require.NoError(t, err)
	})
}

func TestDiffPolicyRouteWithLogical_NormalizesNextHopIPs(t *testing.T) {
	t.Parallel()

	target := []*kubeovnv1.PolicyRoute{{
		Priority:  32000,
		Match:     "ip4.src == 172.16.8.149/32",
		Action:    kubeovnv1.PolicyRouteActionReroute,
		NextHopIP: " 172.31.255.253 , 172.31.255.254 ",
	}}

	existing := []*ovnnb.LogicalRouterPolicy{{
		Priority: 32000,
		Match:    "ip4.src == 172.16.8.149/32",
		Action:   string(kubeovnv1.PolicyRouteActionReroute),
		Nexthops: []string{"172.31.255.254", "172.31.255.253"},
	}}

	dels, adds := diffPolicyRouteWithLogical(existing, target)
	require.Empty(t, dels)
	require.Empty(t, adds)
}

func TestSplitNextHops(t *testing.T) {
	t.Parallel()

	require.Nil(t, splitNextHops(""))
	require.Equal(t, []string{"172.31.255.253", "172.31.255.254"}, splitNextHops("172.31.255.253, 172.31.255.254"))
	require.Equal(t, []string{"172.31.255.253", "172.31.255.254"}, splitNextHops("  172.31.255.253 ,  172.31.255.254  "))
}

func TestDiffPolicyRouteWithLogical_HandlesLegacyNextHopField(t *testing.T) {
	t.Parallel()

	nextHop := "172.31.255.253"
	target := []*kubeovnv1.PolicyRoute{{
		Priority:  32000,
		Match:     "ip4.src == 172.16.8.149/32",
		Action:    kubeovnv1.PolicyRouteActionReroute,
		NextHopIP: nextHop,
	}}

	existing := []*ovnnb.LogicalRouterPolicy{{
		Priority: 32000,
		Match:    "ip4.src == 172.16.8.149/32",
		Action:   string(kubeovnv1.PolicyRouteActionReroute),
		Nexthop:  &nextHop,
	}}

	dels, adds := diffPolicyRouteWithLogical(existing, target)
	require.Empty(t, dels)
	require.Empty(t, adds)
}

func TestDiffPolicyRouteWithLogical_HandlesBFDSessions(t *testing.T) {
	t.Parallel()

	t.Run("matching bfd sessions normalized", func(t *testing.T) {
		target := []*kubeovnv1.PolicyRoute{{
			Priority:  32000,
			Match:     "ip4.src == 172.16.8.149/32",
			Action:    kubeovnv1.PolicyRouteActionReroute,
			NextHopIP: "172.31.255.253,172.31.255.254",
			BfdID:     "bfd-uuid-2, bfd-uuid-1",
		}}

		existing := []*ovnnb.LogicalRouterPolicy{{
			Priority:    32000,
			Match:       "ip4.src == 172.16.8.149/32",
			Action:      string(kubeovnv1.PolicyRouteActionReroute),
			Nexthops:    []string{"172.31.255.253", "172.31.255.254"},
			BFDSessions: []string{"bfd-uuid-1", "bfd-uuid-2"},
		}}

		dels, adds := diffPolicyRouteWithLogical(existing, target)
		require.Empty(t, dels)
		require.Empty(t, adds)
	})

	t.Run("differing bfd sessions triggers update", func(t *testing.T) {
		target := []*kubeovnv1.PolicyRoute{{
			Priority:  32000,
			Match:     "ip4.src == 172.16.8.149/32",
			Action:    kubeovnv1.PolicyRouteActionReroute,
			NextHopIP: "172.31.255.253",
			BfdID:     "bfd-uuid-new",
		}}

		existing := []*ovnnb.LogicalRouterPolicy{{
			Priority:    32000,
			Match:       "ip4.src == 172.16.8.149/32",
			Action:      string(kubeovnv1.PolicyRouteActionReroute),
			Nexthops:    []string{"172.31.255.253"},
			BFDSessions: []string{"bfd-uuid-old"},
		}}

		dels, adds := diffPolicyRouteWithLogical(existing, target)
		require.Len(t, dels, 1)
		require.Len(t, adds, 1)
		require.Equal(t, "bfd-uuid-old", dels[0].BfdID)
		require.Equal(t, "bfd-uuid-new", adds[0].BfdID)
	})
}

func Test_handleAddOrUpdateVpc_policyRoutes_withBfd(t *testing.T) {
	t.Parallel()

	vpcName := "test-vpc-bfd-policy"
	fakeController := newFakeController(t)
	ctrl := fakeController.fakeController
	fakeinformers := fakeController.fakeInformers
	mockOvnClient := fakeController.mockOvnClient

	ctrl.vpcKeyMutex = keymutex.NewHashed(500)

	vpc := &kubeovnv1.Vpc{
		ObjectMeta: metav1.ObjectMeta{
			Name: vpcName,
		},
		Spec: kubeovnv1.VpcSpec{
			PolicyRoutes: []*kubeovnv1.PolicyRoute{
				{
					Priority:  100,
					Match:     "ip4.dst == 10.1.0.0/16",
					Action:    kubeovnv1.PolicyRouteActionReroute,
					NextHopIP: "192.168.1.1,192.168.1.2",
					BfdID:     "bfd-session-1,bfd-session-2",
				},
			},
		},
	}

	_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpc, metav1.CreateOptions{})
	require.NoError(t, err)

	err = fakeinformers.vpcInformer.Informer().GetStore().Add(vpc)
	require.NoError(t, err)

	externalIDs := map[string]string{"vendor": util.CniTypeName}

	mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
	mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
	mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(nil, nil)
	mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
		Name: vpcName,
		Nat:  []string{},
	}, nil)
	mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
	mockOvnClient.EXPECT().AddLogicalRouterPolicy(
		vpcName,
		100,
		"ip4.dst == 10.1.0.0/16",
		string(kubeovnv1.PolicyRouteActionReroute),
		[]string{"192.168.1.1", "192.168.1.2"},
		[]string{"bfd-session-1", "bfd-session-2"},
		externalIDs,
	).Return(nil)
	mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
	mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
	mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
	mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)

	err = ctrl.handleAddOrUpdateVpc(vpcName)
	require.NoError(t, err)
}

func Test_handleAddOrUpdateVpc_policyRoutes_with_Gateway(t *testing.T) {
	t.Parallel()

	vpcName := "test-vpc-gw-policy"
	gwName := "test-gw"

	gw := &kubeovnv1.VpcNatGateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: gwName,
		},
		Spec: kubeovnv1.VpcNatGatewaySpec{
			Vpc:    vpcName,
			Subnet: "net1",
			LanIP:  "192.168.1.100",
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-gw-0",
			Namespace: "kube-system",
			Labels: map[string]string{
				"app":                    util.GenNatGwName(gwName),
				util.VpcNatGatewayLabel: "true",
			},
		},
		Spec: corev1.PodSpec{
			NodeName: "node1",
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			},
			PodIPs: []corev1.PodIP{
				{IP: "192.168.1.10"},
			},
		},
	}

	vpc := &kubeovnv1.Vpc{
		ObjectMeta: metav1.ObjectMeta{
			Name: vpcName,
		},
		Spec: kubeovnv1.VpcSpec{
			PolicyRoutes: []*kubeovnv1.PolicyRoute{
				{
					Priority: 100,
					Match:    "ip4.dst == 10.1.0.0/16",
					Action:   kubeovnv1.PolicyRouteActionReroute,
					Gateway:  gwName,
				},
			},
		},
	}

	opts := &FakeControllerOptions{
		Vpcs:           []*kubeovnv1.Vpc{vpc},
		VpcNatGateways: []*kubeovnv1.VpcNatGateway{gw},
		Pods:           []*corev1.Pod{pod},
	}
	fakeController, err := newFakeControllerWithOptions(t, opts)
	require.NoError(t, err)

	ctrl := fakeController.fakeController
	mockOvnClient := fakeController.mockOvnClient
	ctrl.vpcKeyMutex = keymutex.NewHashed(500)

	externalIDs := map[string]string{"vendor": util.CniTypeName}

	mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
	mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
	mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", externalIDs).Return(nil, nil)
	mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{
		Name: vpcName,
		Nat:  []string{},
	}, nil)
	mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
	mockOvnClient.EXPECT().AddLogicalRouterPolicy(
		vpcName,
		100,
		"ip4.dst == 10.1.0.0/16",
		string(kubeovnv1.PolicyRouteActionReroute),
		[]string{"192.168.1.10"},
		nil,
		externalIDs,
	).Return(nil)
	mockOvnClient.EXPECT().ListLogicalSwitch(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalSwitch{}, nil).AnyTimes()
	mockOvnClient.EXPECT().ListLogicalRouter(gomock.Any(), gomock.Any()).Return([]ovnnb.LogicalRouter{}, nil).AnyTimes()
	mockOvnClient.EXPECT().DeleteLogicalRouterPort(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)
	mockOvnClient.EXPECT().DeleteHAChassisGroup(fmt.Sprintf("bfd@%s", vpcName)).Return(nil)

	err = ctrl.handleAddOrUpdateVpc(vpcName)
	require.NoError(t, err)
}
