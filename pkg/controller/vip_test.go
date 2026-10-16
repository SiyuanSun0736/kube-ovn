package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestEnqueueAddVirtualIPRoutesTerminatingToUpdateQueue(t *testing.T) {
	c := &Controller{
		addVirtualIPQueue:    newTypedRateLimitingQueue[string]("AddVirtualIP", nil),
		updateVirtualIPQueue: newTypedRateLimitingQueue[string]("UpdateVirtualIP", nil),
	}
	t.Cleanup(func() {
		c.addVirtualIPQueue.ShutDown()
		c.updateVirtualIPQueue.ShutDown()
	})

	// After a controller restart the informer (re)list surfaces a terminating vip
	// only via the add handler; it must go to the update queue or the finalizer is
	// never removed.
	now := metav1.Now()
	terminating := &kubeovnv1.Vip{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "stuck",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
	}
	c.enqueueAddVirtualIP(terminating)
	require.Equal(t, 0, c.addVirtualIPQueue.Len())
	require.Equal(t, 1, c.updateVirtualIPQueue.Len())
	item, _ := c.updateVirtualIPQueue.Get()
	c.updateVirtualIPQueue.Done(item)

	// A live vip keeps going to the add queue.
	live := &kubeovnv1.Vip{ObjectMeta: metav1.ObjectMeta{Name: "live"}}
	c.enqueueAddVirtualIP(live)
	require.Equal(t, 1, c.addVirtualIPQueue.Len())
	require.Equal(t, 0, c.updateVirtualIPQueue.Len())
}
