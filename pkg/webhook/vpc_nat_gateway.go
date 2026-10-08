package webhook

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	ovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

var (
	vpcNatGatewayGVK = metav1.GroupVersionKind{Group: ovnv1.SchemeGroupVersion.Group, Version: ovnv1.SchemeGroupVersion.Version, Kind: "VpcNatGateway"}
	iptablesEIPGVK   = metav1.GroupVersionKind{Group: ovnv1.SchemeGroupVersion.Group, Version: ovnv1.SchemeGroupVersion.Version, Kind: "IptablesEIP"}
	iptablesDnatRule = metav1.GroupVersionKind{Group: ovnv1.SchemeGroupVersion.Group, Version: ovnv1.SchemeGroupVersion.Version, Kind: "IptablesDnatRule"}
	iptablesSnatRule = metav1.GroupVersionKind{Group: ovnv1.SchemeGroupVersion.Group, Version: ovnv1.SchemeGroupVersion.Version, Kind: "IptablesSnatRule"}
	iptablesFIPRule  = metav1.GroupVersionKind{Group: ovnv1.SchemeGroupVersion.Group, Version: ovnv1.SchemeGroupVersion.Version, Kind: "IptablesFIPRule"}
)

func (v *ValidatingHook) VpcNatGwCreateOrUpdateHook(ctx context.Context, req admission.Request) admission.Response {
	gw := ovnv1.VpcNatGateway{}
	if err := v.decoder.DecodeRaw(req.Object, &gw); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	// On update: enforce spec.namespace immutability.
	// Changing the namespace would create a new StatefulSet in the new namespace while
	// leaving the old one orphaned; there is no migration path for the running workload.
	if len(req.OldObject.Raw) > 0 {
		gwOld := ovnv1.VpcNatGateway{}
		if err := v.decoder.DecodeRaw(req.OldObject, &gwOld); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if gwOld.Spec.Namespace != gw.Spec.Namespace {
			err := fmt.Errorf("VpcNatGateway %q: spec.namespace is immutable (old: %q, new: %q)",
				gw.Name, gwOld.Spec.Namespace, gw.Spec.Namespace)
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		// Changing the VPC would leave routes, NAT state and interfaces in the old VPC; moving a
		// running gateway is not a supported migration. Delete and recreate it instead.
		if gwOld.Spec.Vpc != gw.Spec.Vpc {
			err := fmt.Errorf("VpcNatGateway %q: spec.vpc is immutable (old: %q, new: %q); delete and recreate the gateway to move it",
				gw.Name, gwOld.Spec.Vpc, gw.Spec.Vpc)
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}

		oldMode := ""
		if gwOld.Annotations != nil {
			oldMode = gwOld.Annotations[util.NatGatewayDataplaneModeAnnotation]
		}
		newMode := ""
		if gw.Annotations != nil {
			newMode = gw.Annotations[util.NatGatewayDataplaneModeAnnotation]
		}
		if oldMode != newMode {
			err := fmt.Errorf("VpcNatGateway %q: annotation %q is immutable once initialized (old: %q, new: %q); delete and recreate gateway to switch dataplane mode",
				gw.Name, util.NatGatewayDataplaneModeAnnotation, oldMode, newMode)
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
	}

	if err := v.ValidateVpcNatConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatGW(ctx, &gw); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) VpcNatGwDeleteHook(ctx context.Context, req admission.Request) admission.Response {
	eipList := ovnv1.IptablesEIPList{}
	if err := v.client.List(ctx, &eipList, cli.MatchingLabels{util.VpcNatGatewayNameLabel: req.Name}); err != nil {
		return ctrlwebhook.Errored(http.StatusInternalServerError, err)
	}
	if len(eipList.Items) != 0 {
		err := fmt.Errorf("vpc-nat-gateway \"%s\" cannot be deleted if any eip is in use", req.Name)
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	// A rule serving a ClusterIP has no EIP, so the EIP check above cannot see it. Without this
	// the gateway could be deleted while such rules still program its data plane, leaving rules
	// (and the VPC routes of their VIPs) behind with no gateway left to clean them up. Service
	// accounting records carry the gateway label when created.
	dnatList := ovnv1.IptablesDnatRuleList{}
	if err := v.client.List(ctx, &dnatList); err != nil {
		return ctrlwebhook.Errored(http.StatusInternalServerError, err)
	}
	for i := range dnatList.Items {
		dnat := &dnatList.Items[i]
		// A terminating rule is already on its way out and must not block the gateway forever.
		if dnat.Spec.EIP != "" || dnat.Labels[util.VpcNatGatewayNameLabel] != req.Name || !dnat.DeletionTimestamp.IsZero() {
			continue
		}
		err := fmt.Errorf("vpc-nat-gateway %q cannot be deleted while dnat rule %q serves clusterIP %s on it",
			req.Name, dnat.Name, dnat.Spec.ClusterIP)
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesEIPCreateHook(ctx context.Context, req admission.Request) admission.Response {
	eip := ovnv1.IptablesEIP{}
	if err := v.decoder.DecodeRaw(req.Object, &eip); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	// IPAM-only EIPs (lbSvc / bgp_lb_vip) have no NatGwDp; they only need an IP
	// allocation — skip NAT-gateway config checks, but still validate IP format
	// and subnet/CIDR fields so malformed EIPs are rejected at admission time.
	if eip.Spec.NatGwDp == "" {
		if err := v.validateIptablesEIPIPFields(ctx, &eip); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		return ctrlwebhook.Allowed("by pass")
	}

	if err := v.ValidateVpcNatConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateIptablesEIP(ctx, &eip); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesEIPUpdateHook(ctx context.Context, req admission.Request) admission.Response {
	eipNew := ovnv1.IptablesEIP{}
	if err := v.decoder.DecodeRaw(req.Object, &eipNew); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	eipOld := ovnv1.IptablesEIP{}
	if err := v.decoder.DecodeRaw(req.OldObject, &eipOld); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	// IPAM-only EIPs (lbSvc / bgp_lb_vip) have no NatGwDp. Their immutability
	// rule: once an IP is allocated (Spec.V4ip or Spec.V6ip set), it cannot change.
	// IP format and subnet/CIDR fields are still validated on every spec change.
	//
	// Guard: if the old object already had a NatGwDp, this update is attempting to
	// clear an immutable field (migrating a NAT-backed EIP to IPAM-only is not
	// supported). Reject before entering the IPAM-only branch so the NAT-gateway
	// immutability checks below cannot be bypassed.
	if eipNew.Spec.NatGwDp == "" {
		if eipOld.Spec.NatGwDp != "" {
			err := fmt.Errorf("IptablesEIP %q: NatGwDp is immutable once set (cannot be cleared, old: %s)",
				eipNew.Name, eipOld.Spec.NatGwDp)
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if eipOld.Spec.V4ip != "" && eipNew.Spec.V4ip != eipOld.Spec.V4ip {
			err := fmt.Errorf("IptablesEIP %q: V4ip is immutable once allocated (old: %s, new: %s)",
				eipNew.Name, eipOld.Spec.V4ip, eipNew.Spec.V4ip)
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if eipOld.Spec.V6ip != "" && eipNew.Spec.V6ip != eipOld.Spec.V6ip {
			err := fmt.Errorf("IptablesEIP %q: V6ip is immutable once allocated (old: %s, new: %s)",
				eipNew.Name, eipOld.Spec.V6ip, eipNew.Spec.V6ip)
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if eipOld.Spec != eipNew.Spec {
			if err := v.validateIptablesEIPIPFields(ctx, &eipNew); err != nil {
				return ctrlwebhook.Errored(http.StatusBadRequest, err)
			}
		}
		return ctrlwebhook.Allowed("by pass")
	}

	// IptablesEIP is an internal resource of a NatGwDp. Once created and Ready,
	// its Spec (including V4ip address) is immutable — the Ready check below blocks
	// any Spec change. NatGwDp is additionally immutable once set (even before Ready),
	// because migrating an EIP across gateways is not supported.
	//
	// This immutability is a key invariant for NAT rule controllers: they rely on
	// EIP.Status.IP being stable for the lifetime of the EIP resource.

	// NatGwDp is immutable once set: changing the gateway would require migrating
	// all associated NAT rules (FIP/DNAT/SNAT) to a different gateway pod,
	// which is not supported by the update handlers.
	if eipOld.Spec.NatGwDp != "" && eipNew.Spec.NatGwDp != eipOld.Spec.NatGwDp {
		err := fmt.Errorf("IptablesEIP %q: NatGwDp is immutable once set (old: %s, new: %s)",
			eipNew.Name, eipOld.Spec.NatGwDp, eipNew.Spec.NatGwDp)
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if eipOld.Spec != eipNew.Spec {
		if eipOld.Status.Ready && eipNew.Status.Redo == eipOld.Status.Redo {
			err := fmt.Errorf("IptablesEIP \"%s\" is ready,not support change", eipNew.Name)
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if err := v.ValidateVpcNatConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}

		if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}

		if err := v.ValidateIptablesEIP(ctx, &eipNew); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
	}
	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesEIPDeleteHook(ctx context.Context, req admission.Request) admission.Response {
	eip := ovnv1.IptablesEIP{}
	if err := v.decoder.DecodeRaw(req.OldObject, &eip); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	// IPAM-only EIPs (lbSvc / bgp_lb_vip) have no NatGwDp and are managed
	// without a VPC NAT gateway. Skip gateway-config validation so these EIPs
	// remain deletable even when --enable-vpc-nat-gw is off or the NAT configmap
	// is absent — exactly the deployments this PR adds support for.
	if eip.Spec.NatGwDp == "" {
		return ctrlwebhook.Allowed("IPAM-only EIP, no NAT gateway config required")
	}

	if err := v.ValidateVpcNatConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if eip.Status.Ready {
		// Match the controller's in-use check: select the NAT rules that claimed this EIP
		// generation by UID. Driving the lookup off eip.Status.Nat instead would skip the check
		// whenever that status is stale, and selecting by address would count rules belonging to a
		// different EIP that happens to carry the same one.
		fipList := ovnv1.IptablesFIPRuleList{}
		snatList := ovnv1.IptablesSnatRuleList{}
		dnatList := ovnv1.IptablesDnatRuleList{}
		selector := eipUIDSelector(&eip)
		for _, list := range []cli.ObjectList{&fipList, &snatList, &dnatList} {
			if err := v.cache.List(ctx, list, selector); err != nil && !k8serrors.IsNotFound(err) {
				return ctrlwebhook.Errored(http.StatusInternalServerError, err)
			}
		}

		var inUse []string
		if len(fipList.Items) != 0 {
			inUse = append(inUse, util.FipUsingEip)
		}
		if len(snatList.Items) != 0 {
			inUse = append(inUse, util.SnatUsingEip)
		}
		if len(dnatList.Items) != 0 {
			inUse = append(inUse, util.DnatUsingEip)
		}
		if len(inUse) != 0 {
			err := fmt.Errorf("eip \"%s\" is still in use,you need to delete the %s of eip first", eip.Name, strings.Join(inUse, ","))
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesDnatCreateHook(ctx context.Context, req admission.Request) admission.Response {
	dnat := ovnv1.IptablesDnatRule{}
	if err := v.decoder.DecodeRaw(req.Object, &dnat); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.validateNftableLbOwnerLabels(req, nil, dnat.Labels); err != nil {
		return ctrlwebhook.Errored(http.StatusForbidden, err)
	}
	if dnat.Spec.Type == ovnv1.DnatRuleTypeShare && !util.IsNftableLbSvcRecord(dnat.Labels) {
		return ctrlwebhook.Errored(http.StatusBadRequest, fmt.Errorf(
			"dnat %q: type=share is reserved for nftable LB Service accounting records", dnat.Name,
		))
	}

	if err := v.ValidateVpcNatConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateIptablesDnat(ctx, &dnat); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesDnatUpdateHook(ctx context.Context, req admission.Request) admission.Response {
	dnatNew := ovnv1.IptablesDnatRule{}
	if err := v.decoder.DecodeRaw(req.Object, &dnatNew); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	dnatOld := ovnv1.IptablesDnatRule{}
	if err := v.decoder.DecodeRaw(req.OldObject, &dnatOld); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.validateNftableLbOwnerLabels(req, dnatOld.Labels, dnatNew.Labels); err != nil {
		return ctrlwebhook.Errored(http.StatusForbidden, err)
	}

	if dnatNew.Spec == dnatOld.Spec {
		// The owner labels decide which rules aggregate their backends into one nft map, so a
		// change to them has to pass the same identity validation as a spec change would.
		if util.NftableLbSvcOwnerID(dnatOld.Labels) != util.NftableLbSvcOwnerID(dnatNew.Labels) {
			if err := v.ValidateIptablesDnat(ctx, &dnatNew); err != nil {
				return ctrlwebhook.Errored(http.StatusBadRequest, err)
			}
		}
		return ctrlwebhook.Allowed("by pass")
	}

	// Service records are mutable accounting data, not forwarding intent. Only the authenticated
	// controller can change their spec; validate the new record, but do not apply ordinary DNAT
	// immutability rules that exist to protect a data-plane writer. An empty identity disables
	// identity checking everywhere (see validateNftableLbOwnerLabels), so it must not funnel
	// record updates into the immutable branches either.
	if (v.controllerIdentity == "" || req.UserInfo.Username == v.controllerIdentity) &&
		util.IsNftableLbSvcRecord(dnatOld.Labels) && util.IsNftableLbSvcRecord(dnatNew.Labels) {
		if err := v.ValidateIptablesDnat(ctx, &dnatNew); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		return ctrlwebhook.Allowed("by pass")
	}

	// A rule serving a ClusterIP is created and deleted together with the Service it comes
	// from, and never updated: its identity is its address, so any spec change would silently
	// reprogram a different VIP. Delete and recreate the rule instead.
	if dnatOld.Spec.ClusterIP != "" {
		return ctrlwebhook.Errored(http.StatusBadRequest, fmt.Errorf(
			"dnat %q serves clusterIP %s and is immutable after creation; delete and recreate it to change it",
			dnatNew.Name, dnatOld.Spec.ClusterIP,
		))
	}
	{
		// Type is immutable after creation
		oldType := dnatOld.Spec.Type
		if oldType == "" {
			oldType = ovnv1.DnatRuleTypeExclusive
		}
		newType := dnatNew.Spec.Type
		if newType == "" {
			newType = ovnv1.DnatRuleTypeExclusive
		}
		if oldType != newType {
			return ctrlwebhook.Errored(http.StatusBadRequest, fmt.Errorf("dnat type is immutable after creation: cannot change from %q to %q", oldType, newType))
		}
		if dnatOld.Spec.SessionAffinity != dnatNew.Spec.SessionAffinity ||
			dnatOld.Spec.SessionAffinityTimeoutSeconds != dnatNew.Spec.SessionAffinityTimeoutSeconds {
			return ctrlwebhook.Errored(http.StatusBadRequest, fmt.Errorf("dnat %q sessionAffinity is immutable after creation", dnatNew.Name))
		}

		if err := v.ValidateVpcNatConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if err := v.ValidateIptablesDnat(ctx, &dnatNew); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesSnatCreateHook(ctx context.Context, req admission.Request) admission.Response {
	snat := ovnv1.IptablesSnatRule{}
	if err := v.decoder.DecodeRaw(req.Object, &snat); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateIptablesSnat(ctx, &snat); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesSnatUpdateHook(ctx context.Context, req admission.Request) admission.Response {
	snatNew := ovnv1.IptablesSnatRule{}
	if err := v.decoder.DecodeRaw(req.Object, &snatNew); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	snatOld := ovnv1.IptablesSnatRule{}
	if err := v.decoder.DecodeRaw(req.OldObject, &snatOld); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if snatNew.Spec != snatOld.Spec {
		if err := v.ValidateVpcNatConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if err := v.ValidateIptablesSnat(ctx, &snatNew); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesFipCreateHook(ctx context.Context, req admission.Request) admission.Response {
	fip := ovnv1.IptablesFIPRule{}
	if err := v.decoder.DecodeRaw(req.Object, &fip); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := v.ValidateIptablesFip(ctx, &fip); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) iptablesFipUpdateHook(ctx context.Context, req admission.Request) admission.Response {
	fipNew := ovnv1.IptablesFIPRule{}
	if err := v.decoder.DecodeRaw(req.Object, &fipNew); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	fipOld := ovnv1.IptablesFIPRule{}
	if err := v.decoder.DecodeRaw(req.OldObject, &fipOld); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if fipNew.Spec != fipOld.Spec {
		if err := v.ValidateVpcNatConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if err := v.ValidateVpcNatGatewayConfig(ctx); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
		if err := v.ValidateIptablesFip(ctx, &fipNew); err != nil {
			return ctrlwebhook.Errored(http.StatusBadRequest, err)
		}
	}
	return ctrlwebhook.Allowed("bypass")
}

func (v *ValidatingHook) ValidateVpcNatGW(ctx context.Context, gw *ovnv1.VpcNatGateway) error {
	prefix, err := v.getNatGwNamePrefix(ctx)
	if err != nil {
		return err
	}
	if err := util.ValidateNatGwStatefulSetNameLength(prefix, gw.Name); err != nil {
		return err
	}

	if gw.Spec.Vpc == "" {
		return errors.New("parameter \"vpc\" cannot be empty")
	}
	vpc := &ovnv1.Vpc{}
	key := types.NamespacedName{Name: gw.Spec.Vpc}
	if err := v.cache.Get(ctx, key, vpc); err != nil {
		return err
	}

	if gw.Spec.Subnet == "" {
		return errors.New("parameter \"subnet\" cannot be empty")
	}

	subnet := &ovnv1.Subnet{}
	key = types.NamespacedName{Name: gw.Spec.Subnet}
	if err := v.cache.Get(ctx, key, subnet); err != nil {
		return err
	}

	if net.ParseIP(gw.Spec.LanIP) == nil {
		err := fmt.Errorf("lanIP %s is not a valid", gw.Spec.LanIP)
		return err
	}

	if !util.CIDRContainIP(subnet.Spec.CIDRBlock, gw.Spec.LanIP) {
		err := fmt.Errorf("lanIP %s is not in the range of subnet %s, cidr %v",
			gw.Spec.LanIP, subnet.Name, subnet.Spec.CIDRBlock)
		return err
	}

	for _, t := range gw.Spec.Tolerations {
		if t.Operator != corev1.TolerationOpExists &&
			t.Operator != corev1.TolerationOpEqual {
			err := fmt.Errorf("invalid taint operator: %s, supported params: \"Equal\", \"Exists\"", t.Operator)
			return err
		}

		if t.Effect != corev1.TaintEffectNoSchedule &&
			t.Effect != corev1.TaintEffectNoExecute &&
			t.Effect != corev1.TaintEffectPreferNoSchedule {
			err := fmt.Errorf("invalid taint effect: %s, supported params: \"NoSchedule\", \"PreferNoSchedule\", \"NoExecute\"", t.Effect)
			return err
		}
	}

	if err := validateQoSPolicyRef(ctx, v.cache, gw.Spec.QoSPolicy); err != nil {
		return err
	}

	return nil
}

// eipUIDSelector matches the NAT rules that claimed this EIP generation. It is the same credential
// the controller's in-use check counts, so admission and the finalizer never disagree.
func eipUIDSelector(eip *ovnv1.IptablesEIP) cli.MatchingLabels {
	return cli.MatchingLabels{util.EipUIDLabel: string(eip.UID)}
}

// validateQoSPolicyRef rejects a reference to a QoSPolicy that does not exist or is terminating.
// The controller keeps such a referrer out of Ready, so failing at admission hands the user the
// error directly. The read comes from an informer cache, so this is best effort and the controller
// stays the authority.
func validateQoSPolicyRef(ctx context.Context, reader cli.Reader, qosPolicy string) error {
	if qosPolicy == "" {
		return nil
	}
	qos := &ovnv1.QoSPolicy{}
	if err := reader.Get(ctx, types.NamespacedName{Name: qosPolicy}, qos); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("qos policy %s does not exist; create it before referencing it: %w", qosPolicy, err)
		}
		return err
	}
	if !qos.DeletionTimestamp.IsZero() {
		return fmt.Errorf("qos policy %s is terminating; wait for its deletion to complete before referencing it", qosPolicy)
	}
	if !controllerutil.ContainsFinalizer(qos, util.KubeOVNControllerFinalizer) || !qos.StatusMatchesSpec() {
		return fmt.Errorf("qos policy %s is not ready; wait for its controller status to match the spec before referencing it", qosPolicy)
	}
	return nil
}

func (v *ValidatingHook) getNatGwNamePrefix(ctx context.Context) (string, error) {
	cm := &corev1.ConfigMap{}
	cmKey := cli.ObjectKey{Namespace: metav1.NamespaceSystem, Name: util.VpcNatConfig}
	if err := v.cache.Get(ctx, cmKey, cm); err != nil {
		return "", err
	}

	prefix := strings.TrimSpace(cm.Data["natGwNamePrefix"])
	if prefix == "" {
		return util.VpcNatGwNameDefaultPrefix, nil
	}
	return prefix, nil
}

func (v *ValidatingHook) ValidateVpcNatGatewayConfig(ctx context.Context) error {
	cm := &corev1.ConfigMap{}
	cmKey := types.NamespacedName{Namespace: "kube-system", Name: util.VpcNatGatewayConfig}
	if err := v.cache.Get(ctx, cmKey, cm); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("configMap \"%s\" not configured", util.VpcNatGatewayConfig)
		}
		return err
	}

	if cm.Data["enable-vpc-nat-gw"] != "true" {
		err := fmt.Errorf("parameter \"enable-vpc-nat-gw\" in ConfigMap \"%s\" not true", util.VpcNatGatewayConfig)
		return err
	}

	return nil
}

func (v *ValidatingHook) ValidateIptablesEIP(ctx context.Context, eip *ovnv1.IptablesEIP) error {
	if eip.Spec.NatGwDp == "" {
		return errors.New("parameter \"natGwDp\" cannot be empty")
	}
	if err := validateQoSPolicyRef(ctx, v.cache, eip.Spec.QoSPolicy); err != nil {
		return err
	}
	return v.validateIptablesEIPIPFields(ctx, eip)
}

// validateIptablesEIPIPFields validates the IP format and subnet/CIDR fields of an
// IptablesEIP. It is intentionally separate from the NatGwDp (NAT-gateway) check so
// that IPAM-only EIPs (lbSvc / bgp_lb_vip, which have no NatGwDp) still receive full
// structural validation on create and update.
func (v *ValidatingHook) validateIptablesEIPIPFields(ctx context.Context, eip *ovnv1.IptablesEIP) error {
	subnet := &ovnv1.Subnet{}
	externalNetwork := util.GetExternalNetwork(eip.Spec.ExternalSubnet)
	key := types.NamespacedName{Name: externalNetwork}
	if err := v.cache.Get(ctx, key, subnet); err != nil {
		return err
	}

	if eip.Spec.V4ip != "" {
		if net.ParseIP(eip.Spec.V4ip) == nil {
			return fmt.Errorf("v4ip %s is not a valid", eip.Spec.V4ip)
		}

		if !util.CIDRContainIP(subnet.Spec.CIDRBlock, eip.Spec.V4ip) {
			err := fmt.Errorf("the vip %s is not in the range of subnet \"%s\", cidr %v",
				eip.Spec.V4ip, subnet.Name, subnet.Spec.CIDRBlock)
			return err
		}
	}

	if eip.Spec.V6ip != "" {
		// v6 ip address can not use upper case
		if util.ContainsUppercase(eip.Spec.V6ip) {
			err := fmt.Errorf("eip %s v6 ip address %s can not contain upper case", eip.Name, eip.Spec.V6ip)
			return err
		}
		if net.ParseIP(eip.Spec.V6ip) == nil {
			err := fmt.Errorf("v6ip %s is not a valid", eip.Spec.V6ip)
			return err
		}

		if !util.CIDRContainIP(subnet.Spec.CIDRBlock, eip.Spec.V6ip) {
			err := fmt.Errorf("the vip %s is not in the range of subnet \"%s\", cidr %v",
				eip.Spec.V6ip, subnet.Name, subnet.Spec.CIDRBlock)
			return err
		}
	}

	return nil
}

func validateIptablesDnatProtocol(protocol string) error {
	if protocol != util.ProtocolTCP && protocol != util.ProtocolUDP {
		return fmt.Errorf("invalid iptable protocol: %s,supported params: \"tcp\", \"udp\"; protocol must be lowercase", protocol)
	}
	return nil
}

func (v *ValidatingHook) ValidateIptablesDnat(ctx context.Context, dnat *ovnv1.IptablesDnatRule) error {
	// A rule addresses at least one VIP. An EIP is the public one (allocated by the EIP and bound
	// on the external interface), a ClusterIP is the internal one (held on lo). A Service handled
	// by the nftable LB service feature carries both on one rule, its ingress IP and its ClusterIP
	// aligned on the same identity; a ClusterIP Service only carries the ClusterIP.
	if dnat.Spec.EIP == "" && dnat.Spec.ClusterIP == "" {
		return errors.New(`parameter "eip" or "clusterIP" cannot be empty`)
	}

	var eip *ovnv1.IptablesEIP
	if dnat.Spec.EIP != "" {
		eip = &ovnv1.IptablesEIP{}
		key := types.NamespacedName{Name: dnat.Spec.EIP}
		if err := v.cache.Get(ctx, key, eip); err != nil {
			return err
		}
		gateway := dnat.Labels[util.VpcNatGatewayNameLabel]
		if gateway != "" && eip.Spec.NatGwDp != "" && gateway != eip.Spec.NatGwDp {
			return fmt.Errorf("dnat %q: gateway label %q does not serve eip %q (natGwDp %q)",
				dnat.Name, gateway, dnat.Spec.EIP, eip.Spec.NatGwDp)
		}
	} else {
		// Without an eip the gateway comes from the controller-owned label, and the internal VIP
		// flow only exists for Service share records.
		gateway := dnat.Labels[util.VpcNatGatewayNameLabel]
		if gateway == "" {
			return fmt.Errorf("dnat %q: gateway label is required with clusterIP when there is no eip", dnat.Name)
		}
		gw := &ovnv1.VpcNatGateway{}
		if err := v.cache.Get(ctx, types.NamespacedName{Name: gateway}, gw); err != nil {
			return fmt.Errorf("dnat %q: gateway %q: %w", dnat.Name, gateway, err)
		}
	}
	if dnat.Spec.ClusterIP != "" {
		if net.ParseIP(dnat.Spec.ClusterIP) == nil || util.CheckProtocol(dnat.Spec.ClusterIP) != ovnv1.ProtocolIPv4 {
			return fmt.Errorf("dnat %q: clusterIP %q must be an IPv4 address", dnat.Name, dnat.Spec.ClusterIP)
		}
		if dnat.Spec.Type != ovnv1.DnatRuleTypeShare {
			return fmt.Errorf("dnat %q: clusterIP requires type=%s", dnat.Name, ovnv1.DnatRuleTypeShare)
		}
	}

	if dnat.Spec.ExternalPort == "" {
		return errors.New("parameter \"externalPort\" cannot be empty")
	}

	if dnat.Spec.InternalPort == "" {
		return errors.New("parameter \"internalPort\" cannot be empty")
	}

	if err := util.ValidatePort(dnat.Spec.ExternalPort); err != nil {
		return fmt.Errorf("externalPort %s is not a valid port: %w", dnat.Spec.ExternalPort, err)
	}

	if err := util.ValidatePort(dnat.Spec.InternalPort); err != nil {
		return fmt.Errorf("internalPort %s is not a valid port: %w", dnat.Spec.InternalPort, err)
	}

	if net.ParseIP(dnat.Spec.InternalIP) == nil {
		err := fmt.Errorf("internalIP %s is not a valid ip", dnat.Spec.InternalIP)
		return err
	}

	if err := validateIptablesDnatProtocol(dnat.Spec.Protocol); err != nil {
		return err
	}

	// Default type to Exclusive if empty
	dnatType := dnat.Spec.Type
	if dnatType == "" {
		dnatType = ovnv1.DnatRuleTypeExclusive
	}

	// Share type DNAT is implemented with IPv4-only nft rules (ip daddr / ip saddr).
	// Reject it for a v6-only EIP. EIPs that are still being allocated (both addresses
	// empty) are allowed here; the controller requeues until the IPv4 address is ready.
	if eip != nil && dnatType == ovnv1.DnatRuleTypeShare && eip.Spec.V4ip == "" && eip.Spec.V6ip != "" {
		return fmt.Errorf("dnat %q cannot use type=share: EIP %q is IPv6-only, share dnat only supports IPv4", dnat.Name, dnat.Spec.EIP)
	}

	switch dnat.Spec.SessionAffinity {
	case ovnv1.DnatSessionAffinityNone:
		if dnat.Spec.SessionAffinityTimeoutSeconds != 0 {
			return fmt.Errorf("dnat %q sets sessionAffinityTimeoutSeconds without ClientIP affinity", dnat.Name)
		}
	case ovnv1.DnatSessionAffinityClientIP:
		if dnatType != ovnv1.DnatRuleTypeShare {
			return fmt.Errorf("dnat %q ClientIP affinity requires type=share", dnat.Name)
		}
		if dnat.Spec.SessionAffinityTimeoutSeconds < 0 || dnat.Spec.SessionAffinityTimeoutSeconds > 86400 {
			return fmt.Errorf("dnat %q sessionAffinityTimeoutSeconds is out of range", dnat.Name)
		}
	default:
		return fmt.Errorf("dnat %q has invalid sessionAffinity %q", dnat.Name, dnat.Spec.SessionAffinity)
	}

	// Check type conflict with existing DNAT rules sharing the same identity. The identity is the
	// VIP the rule addresses, whichever flow it belongs to.
	gwName := dnat.Labels[util.VpcNatGatewayNameLabel]
	if eip != nil {
		gwName = eip.Spec.NatGwDp
	}
	if gwName != "" {
		dnatList := &ovnv1.IptablesDnatRuleList{}
		if err := v.cache.List(ctx, dnatList, cli.MatchingLabels{
			util.VpcNatGatewayNameLabel: gwName,
			util.VpcDnatEPortLabel:      dnat.Spec.ExternalPort,
		}); err != nil {
			return fmt.Errorf("failed to list iptables DNAT rules: %w", err)
		}

		for _, existing := range dnatList.Items {
			// Skip self (for update hooks)
			if existing.Name == dnat.Name {
				continue
			}
			// Only check same identity (VIP + Protocol)
			if !sameDnatIdentitySpec(&existing.Spec, &dnat.Spec) || existing.Spec.Protocol != dnat.Spec.Protocol {
				continue
			}

			existingType := existing.Spec.Type
			if existingType == "" {
				existingType = ovnv1.DnatRuleTypeExclusive
			}

			// Both must be Share type to coexist
			if dnatType != ovnv1.DnatRuleTypeShare || existingType != ovnv1.DnatRuleTypeShare {
				return fmt.Errorf("dnat %q (type=%s) conflicts with existing dnat %q (type=%s) sharing the same identity (eip=%s, port=%s, protocol=%s): both must be type=share to coexist",
					dnat.Name, dnatType, existing.Name, existingType,
					dnat.Spec.EIP, dnat.Spec.ExternalPort, dnat.Spec.Protocol)
			}
		}
	}

	// Check FIP/DNAT exclusivity: FIP claims all traffic to the EIP (EXCLUSIVE_DNAT),
	// which shadows any port-specific DNAT rules (SHARED_DNAT) for the same EIP. A ClusterIP
	// rule has no EIP, so the FIP plane cannot shadow it.
	if eip != nil {
		fipList := &ovnv1.IptablesFIPRuleList{}
		if err := v.cache.List(ctx, fipList, eipUIDSelector(eip)); err != nil {
			return fmt.Errorf("failed to list iptables FIP rules: %w", err)
		}
		if len(fipList.Items) != 0 {
			return fmt.Errorf("EIP %q is already used by FIP rule %q; floating IP requires exclusive use of the EIP (FIP matches all traffic, shadowing port-specific DNAT rules)", dnat.Spec.EIP, fipList.Items[0].Name)
		}
	}

	return nil
}

// sameDnatIdentitySpec mirrors the controller's identity rule: an EIP is the identity whenever
// either side has one, otherwise the ClusterIP is. It is duplicated here because the webhook must
// reject a conflicting rule before it is created.
func sameDnatIdentitySpec(a, b *ovnv1.IptablesDnatRuleSpec) bool {
	if a.EIP != "" || b.EIP != "" {
		return a.EIP == b.EIP
	}
	return a.ClusterIP == b.ClusterIP
}

// nftableLbSvcOwnerLabels are controller-owned accounting metadata. They identify which Service
// produced a record; records themselves never aggregate or write data-plane state.
var nftableLbSvcOwnerLabels = []string{
	util.NftableLbSvcNsLabel,
	util.NftableLbSvcNameLabel,
	util.NftableLbSvcUIDLabel,
	util.NftableLbSvcRecordLabel,
}

func (v *ValidatingHook) validateNftableLbOwnerLabels(req admission.Request, oldLabels, newLabels map[string]string) error {
	// An empty identity disables the check: it is the escape hatch for a deployment whose
	// controller authenticates as a user this webhook cannot be told about.
	if v.controllerIdentity == "" || req.UserInfo.Username == v.controllerIdentity {
		return nil
	}
	for _, key := range nftableLbSvcOwnerLabels {
		if oldLabels[key] != newLabels[key] {
			return fmt.Errorf("label %q is owned by the kube-ovn controller and cannot be set by %q", key, req.UserInfo.Username)
		}
	}
	return nil
}

func (v *ValidatingHook) ValidateIptablesSnat(ctx context.Context, snat *ovnv1.IptablesSnatRule) error {
	if snat.Spec.EIP == "" {
		return errors.New("parameter \"eip\" cannot be empty")
	}
	eip := &ovnv1.IptablesEIP{}
	key := types.NamespacedName{Name: snat.Spec.EIP}
	if err := v.cache.Get(ctx, key, eip); err != nil {
		return err
	}

	if err := util.CheckCidrs(snat.Spec.InternalCIDR); err != nil {
		return fmt.Errorf("invalid cidr %s", snat.Spec.InternalCIDR)
	}

	return nil
}

func (v *ValidatingHook) ValidateIptablesFip(ctx context.Context, fip *ovnv1.IptablesFIPRule) error {
	if fip.Spec.EIP == "" {
		err := errors.New("parameter \"eip\" cannot be empty")
		return err
	}
	eip := &ovnv1.IptablesEIP{}
	key := types.NamespacedName{Name: fip.Spec.EIP}
	if err := v.cache.Get(ctx, key, eip); err != nil {
		return err
	}

	if net.ParseIP(fip.Spec.InternalIP) == nil {
		err := fmt.Errorf("internalIP %s is not a valid", fip.Spec.InternalIP)
		return err
	}

	// Check FIP/DNAT exclusivity: FIP claims all traffic to the EIP (EXCLUSIVE_DNAT),
	// which shadows any port-specific DNAT rules (SHARED_DNAT) for the same EIP.
	dnatList := &ovnv1.IptablesDnatRuleList{}
	if err := v.cache.List(ctx, dnatList, eipUIDSelector(eip)); err != nil {
		return fmt.Errorf("failed to list iptables DNAT rules: %w", err)
	}
	if len(dnatList.Items) != 0 {
		return fmt.Errorf("EIP %q is already used by DNAT rule %q; floating IP requires exclusive use of the EIP (FIP matches all traffic, shadowing port-specific DNAT rules)", fip.Spec.EIP, dnatList.Items[0].Name)
	}

	return nil
}
