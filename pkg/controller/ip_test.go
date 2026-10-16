package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func Test_handleUpdateIP_deletedSubnet(t *testing.T) {
	t.Parallel()

	now := metav1.Now()
	ip := &kubeovnv1.IP{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-ip",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
		Spec: kubeovnv1.IPSpec{
			Subnet:    "deleted-subnet",
			Namespace: "default",
			PodName:   "test-pod",
		},
	}

	fakeCtrl, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		IPs: []*kubeovnv1.IP{ip},
	})
	require.NoError(t, err)

	ctrl := fakeCtrl.fakeController

	// Shut down work queues to avoid goroutine leaks
	t.Cleanup(func() {
		ctrl.updateSubnetStatusQueue.ShutDown()
		ctrl.syncVirtualPortsQueue.ShutDown()
	})

	// The subnet "deleted-subnet" does not exist in the fake client.
	// This must not panic (previously caused NPE in isOvnSubnet).
	err = ctrl.handleUpdateIP("test-ip")
	require.NoError(t, err)

	// Verify the subnet status update was enqueued
	require.Equal(t, 1, ctrl.updateSubnetStatusQueue.Len())
}

func TestEnqueueAddIPRoutesTerminatingToUpdateQueue(t *testing.T) {
	c := &Controller{
		addIPQueue:              newTypedRateLimitingQueue[string]("AddIP", nil),
		updateIPQueue:           newTypedRateLimitingQueue[string]("UpdateIP", nil),
		updateSubnetStatusQueue: newTypedRateLimitingQueue[string]("UpdateSubnetStatus", nil),
	}
	t.Cleanup(func() {
		c.addIPQueue.ShutDown()
		c.updateIPQueue.ShutDown()
		c.updateSubnetStatusQueue.ShutDown()
	})

	// After a controller restart the informer (re)list surfaces a terminating ip
	// only via the add handler; it must go to the update queue or the finalizer is
	// never removed.
	now := metav1.Now()
	terminating := &kubeovnv1.IP{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "pod.default.10.0.0.1",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
		Spec: kubeovnv1.IPSpec{Subnet: "subnet"},
	}
	c.enqueueAddIP(terminating)
	require.Equal(t, 0, c.addIPQueue.Len())
	require.Equal(t, 0, c.updateSubnetStatusQueue.Len())
	require.Equal(t, 1, c.updateIPQueue.Len())
	item, _ := c.updateIPQueue.Get()
	c.updateIPQueue.Done(item)

	// A live ip keeps going to the add queue.
	live := &kubeovnv1.IP{
		ObjectMeta: metav1.ObjectMeta{Name: "pod.default.10.0.0.2"},
		Spec:       kubeovnv1.IPSpec{Subnet: "subnet"},
	}
	c.enqueueAddIP(live)
	require.Equal(t, 1, c.addIPQueue.Len())
	require.Equal(t, 0, c.updateIPQueue.Len())
}
