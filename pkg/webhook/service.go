package webhook

import (
	"context"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

func (v *ValidatingHook) ServiceUpdateHook(_ context.Context, req admission.Request) admission.Response {
	oldSvc := v1.Service{}
	if err := v.decoder.DecodeRaw(req.OldObject, &oldSvc); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}
	newSvc := v1.Service{}
	if err := v.decoder.DecodeRaw(req.Object, &newSvc); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	oldGateway := oldSvc.Annotations[util.VpcNatGatewayAnnotation]
	newGateway := newSvc.Annotations[util.VpcNatGatewayAnnotation]
	if oldGateway == "" {
		return ctrlwebhook.Allowed("bypass")
	}
	if oldGateway != newGateway {
		return ctrlwebhook.Errored(http.StatusBadRequest,
			fmt.Errorf("service %s/%s vpc nat gateway cannot change; delete and recreate the service", newSvc.Namespace, newSvc.Name))
	}
	// Only a LoadBalancer Service consumes the EIP annotation (it names the EIP its ingress IP
	// comes from); for a ClusterIP Service the annotation is inert, so editing it is harmless and
	// must not force a delete/recreate. Removing the annotation is allowed: the reconcile then
	// releases only the EIP identity and keeps the ClusterIP ones. Binding the Service to another
	// EIP still asks for a fresh Service, so an identity is never silently re-homed.
	oldEip, newEip := oldSvc.Annotations[util.EipAnnotation], newSvc.Annotations[util.EipAnnotation]
	if newSvc.Spec.Type == v1.ServiceTypeLoadBalancer && oldEip != newEip && newEip != "" {
		return ctrlwebhook.Errored(http.StatusBadRequest,
			fmt.Errorf("service %s/%s eip cannot change; remove the annotation to release only the EIP identity, or delete and recreate the service", newSvc.Namespace, newSvc.Name))
	}
	return ctrlwebhook.Allowed("bypass")
}
