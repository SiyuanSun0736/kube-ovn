package speaker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	kubeovnlister "github.com/kubeovn/kube-ovn/pkg/client/listers/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

type failingTCManager struct {
	TCManager
	attachErr error
	syncErr   error
}

func (m *failingTCManager) Attach(string) error                  { return m.attachErr }
func (m *failingTCManager) SyncEIPs(map[string]TargetInfo) error { return m.syncErr }

type trackingServicesLister struct {
	listerv1.ServiceLister
	listed bool
}

func (l *trackingServicesLister) List(labels.Selector) ([]*corev1.Service, error) {
	l.listed = true
	return nil, nil
}

func TestTCFailureDoesNotBlockServiceReconciliation(t *testing.T) {
	tcErr := errors.New("TC unavailable")
	for _, mgr := range []*failingTCManager{{attachErr: tcErr}, {syncErr: tcErr}} {
		nodes := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
		require.NoError(t, nodes.Add(&corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node1", Labels: map[string]string{util.BgpSpeakLbVipLabel: "true"}}}))
		services := &trackingServicesLister{}
		c := &Controller{
			config:      &Configuration{NodeName: "node1", NodeRouteEIPMode: true, EnableTcBpfRedirect: true, EnableLbSvcAnnounce: true},
			eipLister:   kubeovnlister.NewIptablesEIPLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
			nodesLister: listerv1.NewNodeLister(nodes), servicesLister: services, tcManager: mgr,
		}
		require.ErrorIs(t, c.syncNodeRouteEIPs(), tcErr)
		require.True(t, services.listed, "TC errors must not prevent the independent Service VIP reconciliation")
	}
}
