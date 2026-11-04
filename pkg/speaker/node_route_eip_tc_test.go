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

func (m *failingTCManager) Attach(string) error                             { return m.attachErr }
func (m *failingTCManager) SyncEIPs(map[string]TargetInfo, ...string) error { return m.syncErr }

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

type trackingTCManager struct {
	TCManager
	closeErr error
	closed   bool
}

func (m *trackingTCManager) Close() error {
	m.closed = true
	return m.closeErr
}

func TestTCCloseFailureRetriesUntilSuccess(t *testing.T) {
	mgr := &trackingTCManager{closeErr: errors.New("permission denied")}
	services := &trackingServicesLister{}
	nodes := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	require.NoError(t, nodes.Add(&corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node1", Labels: map[string]string{util.BgpSpeakLbVipLabel: "true"}}}))
	c := &Controller{
		config:           &Configuration{NodeName: "node1", NodeRouteEIPMode: true, EnableTcBpfRedirect: false, EnableLbSvcAnnounce: true},
		tcManager:        mgr,
		tcCleanupPending: true,
		eipLister:        kubeovnlister.NewIptablesEIPLister(cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})),
		servicesLister:   services,
		nodesLister:      listerv1.NewNodeLister(nodes),
	}
	// 第一次 Reconcile：Close 失败，tcCleanupPending 必须保留为 true 以便下次重试，同时业务对账不被中断
	c.Reconcile()
	require.True(t, mgr.closed, "tcManager.Close must be called")
	require.True(t, c.tcCleanupPending, "tcCleanupPending must be retained on Close failure for retry")
	require.True(t, services.listed, "Reconciliation must proceed even if TC cleanup fails")

	// 第二次 Reconcile：故障恢复，Close 成功，tcCleanupPending 清除为 false
	mgr.closeErr = nil
	c.Reconcile()
	require.False(t, c.tcCleanupPending, "tcCleanupPending must be cleared when Close succeeds")
}
