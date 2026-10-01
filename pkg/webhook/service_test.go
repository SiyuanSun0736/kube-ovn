package webhook

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

func updateRequest(t *testing.T, oldObj, newObj any) admission.Request {
	t.Helper()
	oldRaw, err := json.Marshal(oldObj)
	require.NoError(t, err)
	newRaw, err := json.Marshal(newObj)
	require.NoError(t, err)
	return admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{
			Operation: admissionv1.Update,
			Object:    runtime.RawExtension{Raw: newRaw},
			OldObject: runtime.RawExtension{Raw: oldRaw},
		},
	}
}

func TestServiceUpdateHook(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, v1.AddToScheme(scheme))
	hook := &ValidatingHook{decoder: admission.NewDecoder(scheme)}

	service := func(serviceType v1.ServiceType, gateway, eip string) *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "web",
				Annotations: map[string]string{
					util.VpcNatGatewayAnnotation: gateway,
					util.EipAnnotation:           eip,
				},
			},
			Spec: v1.ServiceSpec{Type: serviceType},
		}
	}
	loadBalancer := func(gateway, eip string) *v1.Service {
		return service(v1.ServiceTypeLoadBalancer, gateway, eip)
	}
	clusterIP := func(gateway, eip string) *v1.Service {
		return service(v1.ServiceTypeClusterIP, gateway, eip)
	}

	t.Run("allows initial binding", func(t *testing.T) {
		resp := hook.ServiceUpdateHook(context.Background(), updateRequest(t, clusterIP("", ""), clusterIP("gw-a", "eip-a")))
		require.True(t, resp.Allowed)
	})

	t.Run("rejects changing or removing gateway", func(t *testing.T) {
		resp := hook.ServiceUpdateHook(context.Background(), updateRequest(t, loadBalancer("gw-a", "eip-a"), loadBalancer("gw-b", "eip-a")))
		require.False(t, resp.Allowed)
		require.Contains(t, resp.Result.Message, "vpc nat gateway cannot change")
		resp = hook.ServiceUpdateHook(context.Background(), updateRequest(t, loadBalancer("gw-a", "eip-a"), loadBalancer("", "")))
		require.False(t, resp.Allowed)
		require.Contains(t, resp.Result.Message, "vpc nat gateway cannot change")
	})

	t.Run("rejects changing the eip of a LoadBalancer service", func(t *testing.T) {
		resp := hook.ServiceUpdateHook(context.Background(), updateRequest(t, loadBalancer("gw-a", "eip-a"), loadBalancer("gw-a", "eip-b")))
		require.False(t, resp.Allowed)
		require.Contains(t, resp.Result.Message, "eip cannot change")
	})

	// The documented identity trimming flow: dropping the EIP annotation releases only the EIP
	// identity while the ClusterIP identities keep being served.
	t.Run("allows removing the eip of a LoadBalancer service", func(t *testing.T) {
		resp := hook.ServiceUpdateHook(context.Background(), updateRequest(t, loadBalancer("gw-a", "eip-a"), loadBalancer("gw-a", "")))
		require.True(t, resp.Allowed)
	})

	t.Run("rejects late-binding an eip to a gateway-bound LoadBalancer service", func(t *testing.T) {
		resp := hook.ServiceUpdateHook(context.Background(), updateRequest(t, loadBalancer("gw-a", ""), loadBalancer("gw-a", "eip-a")))
		require.False(t, resp.Allowed)
		require.Contains(t, resp.Result.Message, "eip cannot change")
	})

	// A ClusterIP Service never reads the eip annotation, so editing it must not force a
	// delete/recreate of the Service.
	t.Run("allows changing the eip annotation of a ClusterIP service", func(t *testing.T) {
		resp := hook.ServiceUpdateHook(context.Background(), updateRequest(t, clusterIP("gw-a", "eip-a"), clusterIP("gw-a", "eip-b")))
		require.True(t, resp.Allowed)
	})
}
