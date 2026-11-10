package controller

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
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
		fakeController.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey(vpcName).Return(7, nil)
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

		// The router tunnel key is persisted together with Standby.
		persisted, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Get(context.Background(), vpcName, metav1.GetOptions{})
		require.NoError(t, err)
		require.Equal(t, 7, persisted.Status.TunnelKey)
		require.True(t, persisted.Status.Standby)
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
		fakeController.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey(vpcName).Return(7, nil)
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
		fakeController.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey(vpcName).Return(7, nil)
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
		fakeController.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey(vpcName).Return(7, nil)
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
		fakeController.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey(vpcName).Return(7, nil)
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
		fakeController.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey(vpcName).Return(7, nil)
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
		NextHopIP: "172.31.255.253,172.31.255.254",
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

func Test_reconcileVpcTunnelKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		persisted int
		sbKey     int
		sbErr     error
		callSB    bool
		want      int
		wantErr   string
	}{
		{name: "syncs missing key from SB", sbKey: 7, callSB: true, want: 7},
		{name: "keeps valid persisted key without SB call", persisted: 5, want: 5},
		{name: "resyncs out-of-range persisted key", persisted: util.MaxTunnelKey + 1, sbKey: 7, callSB: true, want: 7},
		{name: "SB error requeues", sbErr: errors.New("datapath binding not found"), callSB: true, wantErr: "datapath binding not found"},
		{name: "invalid SB key requeues", sbKey: 0, callSB: true, wantErr: "invalid tunnel key 0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fc := newFakeController(t)
			if tt.callSB {
				fc.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey("vpc-a").Return(tt.sbKey, tt.sbErr)
			}
			vpc := &kubeovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: tt.persisted}}
			err := fc.fakeController.reconcileVpcTunnelKey(vpc)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				require.Equal(t, tt.persisted, vpc.Status.TunnelKey, "an invalid key must never be recorded")
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, vpc.Status.TunnelKey)
		})
	}
}

// Subnets wait for Standby, so a VPC whose router key cannot be read must not become Standby.
func Test_handleAddOrUpdateVpc_tunnelKeyErrorKeepsStandbyFalse(t *testing.T) {
	t.Parallel()

	vpcName := "test-vpc"
	fakeController := newFakeController(t)
	ctrl := fakeController.fakeController
	mockOvnClient := fakeController.mockOvnClient
	ctrl.vpcKeyMutex = keymutex.NewHashed(500)

	vpc := &kubeovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: vpcName}}
	_, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Create(context.Background(), vpc, metav1.CreateOptions{})
	require.NoError(t, err)
	require.NoError(t, fakeController.fakeInformers.vpcInformer.Informer().GetStore().Add(vpc))

	mockOvnClient.EXPECT().CreateLogicalRouter(vpcName).Return(nil)
	mockOvnClient.EXPECT().UpdateLogicalRouter(gomock.Any(), gomock.Any()).Return(nil)
	mockOvnClient.EXPECT().GetLogicalRouter(vpcName, false).Return(&ovnnb.LogicalRouter{Name: vpcName}, nil)
	mockOvnClient.EXPECT().ListLogicalRouterStaticRoutes(vpcName, nil, nil, "", gomock.Any()).Return(nil, nil)
	mockOvnClient.EXPECT().ListLogicalRouterPolicies(vpcName, -1, nil, true).Return(nil, nil)
	fakeController.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey(vpcName).Return(0, errors.New("datapath binding not found"))

	require.ErrorContains(t, ctrl.handleAddOrUpdateVpc(vpcName), "datapath binding not found")

	persisted, err := ctrl.config.KubeOvnClient.KubeovnV1().Vpcs().Get(context.Background(), vpcName, metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, persisted.Status.Standby)
	require.Zero(t, persisted.Status.TunnelKey)
}

func Test_enqueueUpdateVpc_wakesPodsWhenTunnelKeySynced(t *testing.T) {
	t.Parallel()

	subnets := []*kubeovnv1.Subnet{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "ovn-subnet"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "10.0.1.0/24", Provider: util.OvnProvider, Vpc: "vpc-a", Default: true},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "net1-subnet"},
			Spec:       kubeovnv1.SubnetSpec{CIDRBlock: "10.0.2.0/24", Provider: "net1.default.ovn", Vpc: "vpc-b"},
		},
	}
	nad := &nadv1.NetworkAttachmentDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "net1", Namespace: "default"},
		Spec: nadv1.NetworkAttachmentDefinitionSpec{
			Config: `{"cniVersion": "0.3.1", "name": "net1", "type": "kube-ovn", "provider": "net1.default.ovn"}`,
		},
	}
	pod := func(name string, annotations map[string]string, hostNetwork bool) *corev1.Pod {
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default", Annotations: annotations},
			Spec:       corev1.PodSpec{HostNetwork: hostNetwork},
			Status:     corev1.PodStatus{Phase: corev1.PodRunning},
		}
	}
	allocated := func(key string) map[string]string {
		a := map[string]string{util.AllocatedAnnotation: "true", util.LogicalSwitchAnnotation: "ovn-subnet"}
		if key != "" {
			a[util.TunnelKeyAnnotation] = key
		}
		return a
	}
	// Primary NIC converged; a hot-plugged attachment still waits for vpc-b's key.
	hotplug := allocated("4")
	hotplug[nadv1.NetworkAttachmentAnnot] = "default/net1"
	pods := []*corev1.Pod{
		pod("pending", nil, false),
		pod("attachment-pending", hotplug, false),
		pod("needs-repair", allocated(""), false),
		pod("converged", allocated("4"), false),
		pod("host", nil, true),
	}
	synced := &kubeovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: "vpc-b"}, Status: kubeovnv1.VpcStatus{TunnelKey: 5}}
	unsynced := synced.DeepCopy()
	unsynced.Status.TunnelKey = 0

	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Namespaces:         []*corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}},
		Subnets:            subnets,
		NetworkAttachments: []*nadv1.NetworkAttachmentDefinition{nad},
		Vpcs: []*kubeovnv1.Vpc{
			{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: 4}},
			synced,
		},
		Pods: pods,
	})
	require.NoError(t, err)
	c := fc.fakeController
	c.addOrUpdatePodQueue = newTypedRateLimitingQueue[string]("AddOrUpdatePod", nil)

	drain := func(q workqueue.TypedRateLimitingInterface[string]) []string {
		var keys []string
		for q.Len() != 0 {
			key, _ := q.Get()
			q.Done(key)
			keys = append(keys, key)
		}
		slices.Sort(keys)
		return keys
	}

	c.enqueueUpdateVpc(synced, synced.DeepCopy())
	require.Empty(t, c.vpcTunnelKeySyncedCh, "no transition, no wake")

	// The callback only signals; a burst of transitions coalesces into one sweep.
	c.enqueueUpdateVpc(unsynced, synced)
	c.enqueueUpdateVpc(unsynced, synced)
	require.Len(t, c.vpcTunnelKeySyncedCh, 1)
	require.Zero(t, c.addOrUpdatePodQueue.Len(), "the sweep must not run in the informer callback")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.runVpcTunnelKeyWaker(ctx)
	require.Eventually(t, func() bool {
		return c.addOrUpdatePodQueue.Len() == 2 && c.repairTunnelKeyQueue.Len() == 1
	}, 5*time.Second, 10*time.Millisecond)
	require.Equal(t, []string{"default/attachment-pending", "default/pending"}, drain(c.addOrUpdatePodQueue))
	require.Equal(t, []string{"default/needs-repair"}, drain(c.repairTunnelKeyQueue))

	// The waker consumes the signal before it sweeps, so a key that syncs after a sweep started
	// queues another sweep instead of being absorbed by the running one.
	require.Empty(t, c.vpcTunnelKeySyncedCh)
	late := pod("late", nil, false)
	require.NoError(t, fc.fakeInformers.podInformer.Informer().GetStore().Add(late))
	c.enqueueUpdateVpc(unsynced, synced)
	require.Eventually(t, func() bool { return c.addOrUpdatePodQueue.Len() == 3 }, 5*time.Second, 10*time.Millisecond)
	cancel()
	require.Equal(t, []string{"default/attachment-pending", "default/late", "default/pending"}, drain(c.addOrUpdatePodQueue))
}
