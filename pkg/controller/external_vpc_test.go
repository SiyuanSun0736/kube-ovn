package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/ovsdb/ovnnb"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// Imported VPCs never reach handleAddOrUpdateVpc, so the periodic sync must record their router
// key or pods on their subnets would wait at the tunnel_key gate forever.
func TestSyncExternalVpcRecordsTunnelKey(t *testing.T) {
	existing := &kubeovnv1.Vpc{ObjectMeta: metav1.ObjectMeta{Name: "ext-old", Labels: map[string]string{util.VpcExternalLabel: "true"}}}
	fc, err := newFakeControllerWithOptions(t, &FakeControllerOptions{Vpcs: []*kubeovnv1.Vpc{existing}})
	require.NoError(t, err)
	c := fc.fakeController

	routers := []ovnnb.LogicalRouter{{Name: "ext-old"}, {Name: "ext-new"}}
	fc.mockOvnClient.EXPECT().ListLogicalRouter(false, gomock.Any()).Return(routers, nil)
	fc.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey("ext-old").Return(9, nil)
	fc.mockOvnSbClient.EXPECT().GetLogicalRouterTunnelKey("ext-new").Return(10, nil)

	c.syncExternalVpc()

	for name, want := range map[string]int{"ext-old": 9, "ext-new": 10} {
		vpc, err := c.config.KubeOvnClient.KubeovnV1().Vpcs().Get(context.Background(), name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Equal(t, want, vpc.Status.TunnelKey, name)
	}
}
