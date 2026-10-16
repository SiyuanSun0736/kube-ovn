package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// assertEnqueueAddRoutingWithFinalizer checks the OVN add-path routing: a live object goes to the
// add queue, a terminating object with a finalizer goes to the update queue (deletion cleanup), and
// a terminating object whose finalizer is already gone is skipped (cleanup done, awaiting deletion).
func assertEnqueueAddRoutingWithFinalizer(
	t *testing.T,
	addQueue, updateQueue workqueue.TypedRateLimitingInterface[string],
	enqueue func(any),
	live, terminatingWithFinalizer, terminatingNoFinalizer any,
) {
	t.Helper()
	enqueue(live)
	require.Equal(t, 1, addQueue.Len(), "live object should go to the add queue")
	require.Equal(t, 0, updateQueue.Len(), "live object must not go to the update queue")

	enqueue(terminatingWithFinalizer)
	require.Equal(t, 1, addQueue.Len(), "terminating object must not go to the add queue")
	require.Equal(t, 1, updateQueue.Len(), "terminating object with finalizer should go to the update queue")

	enqueue(terminatingNoFinalizer)
	require.Equal(t, 1, addQueue.Len(), "terminating object without finalizer must not go to the add queue")
	require.Equal(t, 1, updateQueue.Len(), "terminating object without finalizer must not be re-enqueued")
}

func TestEnqueueAddOvnEip(t *testing.T) {
	t.Parallel()
	c := &Controller{
		addOvnEipQueue:    newTypedRateLimitingQueue[string]("AddOvnEip", nil),
		updateOvnEipQueue: newTypedRateLimitingQueue[string]("UpdateOvnEip", nil),
	}
	t.Cleanup(c.addOvnEipQueue.ShutDown)
	t.Cleanup(c.updateOvnEipQueue.ShutDown)
	now := metav1.Now()
	fin := []string{util.KubeOVNControllerFinalizer}
	assertEnqueueAddRoutingWithFinalizer(t, c.addOvnEipQueue, c.updateOvnEipQueue, c.enqueueAddOvnEip,
		&kubeovnv1.OvnEip{ObjectMeta: metav1.ObjectMeta{Name: "live-eip"}},
		&kubeovnv1.OvnEip{ObjectMeta: metav1.ObjectMeta{Name: "terminating-eip", DeletionTimestamp: &now, Finalizers: fin}},
		&kubeovnv1.OvnEip{ObjectMeta: metav1.ObjectMeta{Name: "terminating-eip-no-finalizer", DeletionTimestamp: &now}},
	)
}

func TestEnqueueAddOvnFip(t *testing.T) {
	t.Parallel()
	c := &Controller{
		addOvnFipQueue:    newTypedRateLimitingQueue[string]("AddOvnFip", nil),
		updateOvnFipQueue: newTypedRateLimitingQueue[string]("UpdateOvnFip", nil),
	}
	t.Cleanup(c.addOvnFipQueue.ShutDown)
	t.Cleanup(c.updateOvnFipQueue.ShutDown)
	now := metav1.Now()
	fin := []string{util.KubeOVNControllerFinalizer}
	assertEnqueueAddRoutingWithFinalizer(t, c.addOvnFipQueue, c.updateOvnFipQueue, c.enqueueAddOvnFip,
		&kubeovnv1.OvnFip{ObjectMeta: metav1.ObjectMeta{Name: "live-fip"}},
		&kubeovnv1.OvnFip{ObjectMeta: metav1.ObjectMeta{Name: "terminating-fip", DeletionTimestamp: &now, Finalizers: fin}},
		&kubeovnv1.OvnFip{ObjectMeta: metav1.ObjectMeta{Name: "terminating-fip-no-finalizer", DeletionTimestamp: &now}},
	)
}

func TestEnqueueAddOvnDnatRule(t *testing.T) {
	t.Parallel()
	c := &Controller{
		addOvnDnatRuleQueue:    newTypedRateLimitingQueue[string]("AddOvnDnat", nil),
		updateOvnDnatRuleQueue: newTypedRateLimitingQueue[string]("UpdateOvnDnat", nil),
	}
	t.Cleanup(c.addOvnDnatRuleQueue.ShutDown)
	t.Cleanup(c.updateOvnDnatRuleQueue.ShutDown)
	now := metav1.Now()
	fin := []string{util.KubeOVNControllerFinalizer}
	assertEnqueueAddRoutingWithFinalizer(t, c.addOvnDnatRuleQueue, c.updateOvnDnatRuleQueue, c.enqueueAddOvnDnatRule,
		&kubeovnv1.OvnDnatRule{ObjectMeta: metav1.ObjectMeta{Name: "live-dnat"}},
		&kubeovnv1.OvnDnatRule{ObjectMeta: metav1.ObjectMeta{Name: "terminating-dnat", DeletionTimestamp: &now, Finalizers: fin}},
		&kubeovnv1.OvnDnatRule{ObjectMeta: metav1.ObjectMeta{Name: "terminating-dnat-no-finalizer", DeletionTimestamp: &now}},
	)
}

func TestEnqueueAddOvnSnatRule(t *testing.T) {
	t.Parallel()
	c := &Controller{
		addOvnSnatRuleQueue:    newTypedRateLimitingQueue[string]("AddOvnSnat", nil),
		updateOvnSnatRuleQueue: newTypedRateLimitingQueue[string]("UpdateOvnSnat", nil),
	}
	t.Cleanup(c.addOvnSnatRuleQueue.ShutDown)
	t.Cleanup(c.updateOvnSnatRuleQueue.ShutDown)
	now := metav1.Now()
	fin := []string{util.KubeOVNControllerFinalizer}
	assertEnqueueAddRoutingWithFinalizer(t, c.addOvnSnatRuleQueue, c.updateOvnSnatRuleQueue, c.enqueueAddOvnSnatRule,
		&kubeovnv1.OvnSnatRule{ObjectMeta: metav1.ObjectMeta{Name: "live-snat"}},
		&kubeovnv1.OvnSnatRule{ObjectMeta: metav1.ObjectMeta{Name: "terminating-snat", DeletionTimestamp: &now, Finalizers: fin}},
		&kubeovnv1.OvnSnatRule{ObjectMeta: metav1.ObjectMeta{Name: "terminating-snat-no-finalizer", DeletionTimestamp: &now}},
	)
}
