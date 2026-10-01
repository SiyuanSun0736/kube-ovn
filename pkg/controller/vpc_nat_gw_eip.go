package controller

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/ovs"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func (c *Controller) enqueueAddIptablesEip(obj any) {
	eip := obj.(*kubeovnv1.IptablesEIP)
	key := cache.MetaObjectToName(eip).String()
	// A terminating object reconciles via the update queue for cleanup (handleAdd returns early).
	if enqueueUpdateIfTerminating(c.updateIptablesEipQueue, key, "iptables eip", eip.DeletionTimestamp) {
		if err := c.enqueueIptablesEipReferrers(eip, false); err != nil {
			klog.Errorf("failed to enqueue referrers of terminating eip %s during add replay: %v", key, err)
		}
		return
	}
	klog.Infof("enqueue add iptables eip %s", key)
	c.addIptablesEipQueue.Add(key)
	// A gateway-mode Service that references this EIP may have been waiting for it to be created.
	c.enqueueGwNftableLbServicesForEIP(eip.Name)
}

// enqueueIptablesEipReferrers wakes NAT rules that may have been waiting for this EIP to become ready.
func (c *Controller) enqueueIptablesEipReferrers(eip *kubeovnv1.IptablesEIP, usable bool) error {
	var errs []error
	eipUID := string(eip.UID)
	fips, err := c.iptablesFipsLister.List(labels.Everything())
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list fips referencing eip %s: %w", eip.Name, err))
	} else {
		for _, fip := range fips {
			boundUID := fip.Labels[util.EipUIDLabel]
			if fip.DeletionTimestamp.IsZero() && (fip.Spec.EIP == eip.Name || (!usable && boundUID == eipUID)) {
				if !usable && boundUID != "" && boundUID != eipUID {
					continue
				}
				switch {
				case fip.Status.V4ip == "" || fip.Status.NatGwDp == "" || fip.Status.InternalIP == "":
					c.addIptablesFipQueue.Add(fip.Name)
				case !usable || fip.Status.V4ip != eip.Status.IP || fip.Status.NatGwDp != eip.Spec.NatGwDp ||
					fip.Status.InternalIP != fip.Spec.InternalIP || boundUID != eipUID:
					c.updateIptablesFipQueue.Add(fip.Name)
				case !fip.Status.Ready:
					c.addIptablesFipQueue.Add(fip.Name)
				}
			}
		}
	}
	dnats, err := c.iptablesDnatRulesLister.List(labels.Everything())
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list dnats referencing eip %s: %w", eip.Name, err))
	} else {
		for _, dnat := range dnats {
			boundUID := dnat.Labels[util.EipUIDLabel]
			if dnat.DeletionTimestamp.IsZero() && (dnat.Spec.EIP == eip.Name || (!usable && boundUID == eipUID)) {
				if !usable && boundUID != "" && boundUID != eipUID {
					continue
				}
				switch {
				case dnat.Status.V4ip == "" || dnat.Status.NatGwDp == "" || dnat.Status.Protocol == "" ||
					dnat.Status.ExternalPort == "" || dnat.Status.InternalIP == "" || dnat.Status.InternalPort == "":
					c.addIptablesDnatRuleQueue.Add(dnat.Name)
				case !usable || dnat.Status.V4ip != eip.Status.IP || dnat.Status.NatGwDp != eip.Spec.NatGwDp ||
					dnat.Status.Protocol != dnat.Spec.Protocol || dnat.Status.ExternalPort != dnat.Spec.ExternalPort ||
					dnat.Status.InternalIP != dnat.Spec.InternalIP || dnat.Status.InternalPort != dnat.Spec.InternalPort || boundUID != eipUID:
					c.updateIptablesDnatRuleQueue.Add(dnat.Name)
				case !dnat.Status.Ready:
					c.addIptablesDnatRuleQueue.Add(dnat.Name)
				}
			}
		}
	}
	snats, err := c.iptablesSnatRulesLister.List(labels.Everything())
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list snats referencing eip %s: %w", eip.Name, err))
	} else {
		for _, snat := range snats {
			boundUID := snat.Labels[util.EipUIDLabel]
			if snat.DeletionTimestamp.IsZero() && (snat.Spec.EIP == eip.Name || (!usable && boundUID == eipUID)) {
				if !usable && boundUID != "" && boundUID != eipUID {
					continue
				}
				statusV4Cidr, _ := util.SplitStringIP(snat.Status.InternalCIDR)
				specV4Cidr, _ := util.SplitStringIP(snat.Spec.InternalCIDR)
				switch {
				case snat.Status.V4ip == "" || snat.Status.NatGwDp == "" || snat.Status.InternalCIDR == "":
					c.addIptablesSnatRuleQueue.Add(snat.Name)
				case !usable || snat.Status.V4ip != eip.Status.IP || snat.Status.NatGwDp != eip.Spec.NatGwDp ||
					statusV4Cidr != specV4Cidr || boundUID != eipUID:
					c.updateIptablesSnatRuleQueue.Add(snat.Name)
				case !snat.Status.Ready:
					c.addIptablesSnatRuleQueue.Add(snat.Name)
				}
			}
		}
	}
	return errors.Join(errs...)
}

func (c *Controller) enqueueUpdateIptablesEip(oldObj, newObj any) {
	oldEip := oldObj.(*kubeovnv1.IptablesEIP)
	newEip := newObj.(*kubeovnv1.IptablesEIP)
	if !newEip.DeletionTimestamp.IsZero() ||
		oldEip.Status.Redo != newEip.Status.Redo ||
		oldEip.Spec.QoSPolicy != newEip.Spec.QoSPolicy {
		key := cache.MetaObjectToName(newEip).String()
		klog.Infof("enqueue update iptables eip %s", key)
		c.updateIptablesEipQueue.Add(key)
	}

	// When the QoSLabel is cleared or switched, re-enqueue the previous QoS policy so it can drop
	// its finalizer once unused (the queue key is the policy name).
	c.enqueueQoSPolicyRelease(oldEip.Labels, newEip.Labels)
	if oldEip.Status.Ready != newEip.Status.Ready || oldEip.Status.IP != newEip.Status.IP ||
		(oldEip.DeletionTimestamp.IsZero() && !newEip.DeletionTimestamp.IsZero()) {
		usable := newEip.DeletionTimestamp.IsZero() && newEip.Status.Ready && newEip.Status.IP != ""
		if err := c.enqueueIptablesEipReferrers(newEip, usable); err != nil {
			klog.Errorf("failed to enqueue referrers of eip %s: %v", newEip.Name, err)
		}
		// The share records of a gateway-mode Service carry the EIP address, so they have to be
		// rewritten (or released) when it appears, disappears or becomes unusable.
		c.enqueueGwNftableLbServicesForEIP(newEip.Name)
	}
}

func (c *Controller) enqueueDelIptablesEip(obj any) {
	var eip *kubeovnv1.IptablesEIP
	switch t := obj.(type) {
	case *kubeovnv1.IptablesEIP:
		eip = t
	case cache.DeletedFinalStateUnknown:
		e, ok := t.Obj.(*kubeovnv1.IptablesEIP)
		if !ok {
			klog.Warningf("unexpected object type: %T", t.Obj)
			return
		}
		eip = e
	default:
		klog.Warningf("unexpected type: %T", obj)
		return
	}

	key := cache.MetaObjectToName(eip).String()
	klog.Infof("enqueue del iptables eip %s", key)
	c.delIptablesEipQueue.Add(eip)
	if err := c.enqueueIptablesEipReferrers(eip, false); err != nil {
		klog.Errorf("failed to enqueue referrers of deleted eip %s: %v", key, err)
	}
	// The Service owns the share records that reference this EIP and has to release them, or the
	// EIP's finalizer waits for rules nothing else deletes.
	c.enqueueGwNftableLbServicesForEIP(eip.Name)

	// Re-trigger QoS reconcile so it can drop its finalizer once unused. DeleteFunc runs after
	// the informer cache dropped this EIP; the queue key is the policy name.
	c.enqueueQoSPolicyRelease(eip.Labels, nil)
}

// natEipNamespace returns the namespace where the NAT gateway pod for the given EIP resides.
// Uses eip.Spec.Namespace when set; falls back to the controller's own PodNamespace.
func (c *Controller) natEipNamespace(eip *kubeovnv1.IptablesEIP) string {
	if eip.Spec.Namespace != "" {
		return eip.Spec.Namespace
	}
	// Derive the namespace from the referenced VpcNatGateway so that EIPs without an
	// explicit spec.namespace always locate the Pod in the correct namespace.
	// natGwNamespaceByName already falls back to c.config.PodNamespace when the GW is not found.
	return c.natGwNamespaceByName(eip.Spec.NatGwDp)
}

func (c *Controller) handleAddIptablesEip(key string) error {
	cachedEip, err := c.iptablesEipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	// The key may have been queued while the object was still live; the update queue owns cleanup.
	if !cachedEip.DeletionTimestamp.IsZero() {
		return nil
	}

	if vpcNatEnabled != "true" {
		return errors.New("iptables nat gw not enable")
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle add iptables eip %s", key)

	if cachedEip.Status.Ready && cachedEip.Status.IP != "" {
		if cachedEip.Spec.QoSPolicy == "" {
			return nil
		}
		qos, getErr := c.getAvailableQoSPolicy(cachedEip.Spec.QoSPolicy)
		if getErr != nil {
			c.updateIptablesEipQueue.Add(key)
			return nil
		}
		if !controllerutil.ContainsFinalizer(qos, util.KubeOVNControllerFinalizer) {
			c.updateIptablesEipQueue.Add(key)
			return nil
		}
		if cachedEip.Labels[util.QoSPolicyUIDLabel] == string(qos.UID) {
			return nil
		}
		c.updateIptablesEipQueue.Add(key)
		return nil
	}

	if _, err = c.getBindableQoSPolicy(cachedEip.Spec.QoSPolicy); err != nil {
		return err
	}
	if err = c.checkNatGwNotTerminating(cachedEip.Spec.NatGwDp); err != nil {
		return err
	}

	subnetName := util.GetExternalNetwork(cachedEip.Spec.ExternalSubnet)
	subnet, err := c.subnetsLister.Get(subnetName)
	if err != nil {
		klog.Errorf("failed to get subnet %s: %v", subnetName, err)
		return err
	}

	// v6 ip address can not use upper case
	if util.ContainsUppercase(cachedEip.Spec.V6ip) {
		err := fmt.Errorf("eip %s v6 ip address %s can not contain upper case", cachedEip.Name, cachedEip.Spec.V6ip)
		klog.Error(err)
		return err
	}

	// make sure vpc nat gw pod is ready before eip allocation
	if _, err := c.getNatGwPod(cachedEip.Spec.NatGwDp, c.natEipNamespace(cachedEip)); err != nil {
		klog.Error(err)
		return err
	}

	var v4ip, v6ip, mac string
	portName := ovs.PodNameToPortName(cachedEip.Name, cachedEip.Namespace, subnet.Spec.Provider)
	if cachedEip.Spec.V4ip != "" {
		if v4ip, v6ip, mac, err = c.acquireStaticEip(cachedEip.Name, cachedEip.Namespace, portName, cachedEip.Spec.V4ip, subnet.Name); err != nil {
			klog.Errorf("failed to acquire static eip, err: %v", err)
			return err
		}
	} else {
		// Random allocate
		if v4ip, v6ip, mac, err = c.acquireEip(cachedEip.Name, cachedEip.Namespace, portName, subnet.Name); err != nil {
			klog.Errorf("failed to allocate eip, err: %v", err)
			return err
		}
	}
	eipV4Cidr, _ := util.SplitStringIP(subnet.Spec.CIDRBlock)
	if v4ip == "" || eipV4Cidr == "" {
		err = fmt.Errorf("subnet %s does not support ipv4", subnet.Name)
		klog.Error(err)
		return err
	}
	addrV4, err := util.GetIPAddrWithMask(v4ip, eipV4Cidr)
	if err != nil {
		err = fmt.Errorf("failed to get eip %s with mask by cidr %s: %w", v4ip, eipV4Cidr, err)
		klog.Error(err)
		return err
	}

	if err = c.createEipInPod(cachedEip.Spec.NatGwDp, addrV4, c.natEipNamespace(cachedEip)); err != nil {
		klog.Errorf("failed to create eip '%s' in pod, %v", key, err)
		return err
	}

	if cachedEip.Spec.QoSPolicy != "" {
		if err = c.addEipQoS(cachedEip, v4ip); err != nil {
			klog.Errorf("failed to add qos '%s' in pod, %v", key, err)
			return err
		}
	}
	// Resolve the gateway namespace for spec.namespace backfill inside createOrUpdateEipCR.
	gwNamespace := c.natGwNamespaceByName(cachedEip.Spec.NatGwDp)
	if err = c.createOrUpdateEipCR(key, v4ip, v6ip, mac, cachedEip.Spec.NatGwDp, cachedEip.Spec.QoSPolicy, subnet.Name, gwNamespace); err != nil {
		klog.Errorf("failed to update eip %s, %v", key, err)
		return err
	}

	// Trigger subnet status update after all operations complete
	// At this point: IPAM allocated, IptablesEIP CR created with labels+status+finalizer
	c.updateSubnetStatusQueue.Add(subnet.Name)
	return nil
}

func (c *Controller) handleResetIptablesEip(key string) error {
	if _, err := c.iptablesEipsLister.Get(key); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	klog.Infof("handle reset eip %s", key)
	if err := c.patchEipLabel(key); err != nil {
		klog.Errorf("failed to patch label for eip %s, %v", key, err)
		return err
	}
	if err := c.patchEipStatus(key, "", "", "", true); err != nil {
		klog.Errorf("failed to reset nat for eip %s, %v", key, err)
		return err
	}
	return nil
}

func (c *Controller) handleUpdateIptablesEip(key string) error {
	cachedEip, err := c.iptablesEipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle update iptables eip %s", key)

	// The cleanup path runs before the external subnet is resolved: a missing or malformed subnet
	// must not hold the finalizer of an EIP that is already on its way out.
	if !cachedEip.DeletionTimestamp.IsZero() {
		klog.Infof("clean eip %q in pod", key)

		// Check if EIP is still being used by any NAT rules (FIP/DNAT/SNAT)
		// Only remove finalizer when no NAT rules are using it
		// Read through the API server: a rule that claimed this EIP moments ago may not be in the
		// informer cache yet, and releasing on that stale view would orphan its gateway pod rules.
		nat, err := c.getIptablesEipNatFromAPI(cachedEip)
		if err != nil {
			klog.Errorf("failed to get eip %s nat rules, %v", key, err)
			return err
		}
		if nat != "" {
			klog.Infof("eip %s is still being used by NAT rules: %s, waiting for them to be deleted", key, nat)
			return nil
		}
		if c.nftableLbSvcTeardownHoldsEip(cachedEip) {
			// No NAT rule claims the EIP any more, but a referencing Service has not finished its
			// teardown: its cleanup reads the identity witness (the address) from this very
			// object, because its records may be gone and its ingress never published. Releasing
			// now would strand the programmed identities with no witness left at all.
			klog.Infof("eip %s: no nat rules left, but a nftable lb service has not settled its teardown, holding deletion", key)
			c.updateIptablesEipQueue.AddAfter(key, 5*time.Second)
			return nil
		}

		if vpcNatEnabled == "true" {
			var v4Cidr string
			subnet, err := c.subnetsLister.Get(util.GetExternalNetwork(cachedEip.Spec.ExternalSubnet))
			switch {
			case err == nil:
				if v4Cidr, _ = util.SplitStringIP(subnet.Spec.CIDRBlock); v4Cidr == "" {
					klog.Warningf("external subnet of eip %s has no ipv4 cidr, skip cleaning its address in pod", key)
				}
			case k8serrors.IsNotFound(err):
				// Only reachable if the subnet's finalizer was forced off: it holds while any
				// address is allocated, and this EIP is one of them.
				klog.Warningf("external subnet of eip %s is gone, skip cleaning its address in pod", key)
			default:
				klog.Errorf("failed to get external subnet of eip %s: %v", key, err)
				return err
			}
			if v4Cidr != "" {
				v4ipCidr, err := util.GetIPAddrWithMask(cachedEip.Status.IP, v4Cidr)
				if err != nil {
					err = fmt.Errorf("failed to get eip %s with mask by cidr %s: %w", cachedEip.Status.IP, v4Cidr, err)
					klog.Error(err)
					return err
				}
				if err = c.deleteEipInPod(cachedEip.Spec.NatGwDp, v4ipCidr, c.natEipNamespace(cachedEip)); err != nil {
					klog.Errorf("failed to clean eip '%s' in pod, %v", key, err)
					return err
				}
			}
		}
		if cachedEip.Status.QoSPolicy != "" {
			if err = c.delEipQoS(cachedEip, cachedEip.Status.IP); err != nil {
				klog.Errorf("failed to del qos '%s' in pod, %v", key, err)
				return err
			}
		}
		// Release IP from IPAM before removing finalizer
		c.ipam.ReleaseAddressByPod(key, cachedEip.Spec.ExternalSubnet)

		// Now remove finalizer, which will trigger subnet status update.
		// QoS reconcile is re-triggered from the EIP DeleteFunc after the cache drops this EIP.
		if err = c.handleDelIptablesEipFinalizer(key); err != nil {
			klog.Errorf("failed to handle del finalizer for eip %s, %v", key, err)
			return err
		}

		return nil
	}

	subnetName := util.GetExternalNetwork(cachedEip.Spec.ExternalSubnet)
	subnet, err := c.subnetsLister.Get(subnetName)
	if err != nil {
		klog.Errorf("failed to get subnet %s: %v", subnetName, err)
		return err
	}

	v4Cidr, _ := util.SplitStringIP(subnet.Spec.CIDRBlock)
	if v4Cidr == "" {
		err = fmt.Errorf("subnet %s does not support ipv4", subnet.Name)
		klog.Error(err)
		return err
	}

	klog.Infof("handle update eip %s", key)
	// v6 ip address can not use upper case
	if util.ContainsUppercase(cachedEip.Spec.V6ip) {
		err := fmt.Errorf("eip %s v6 ip address %s can not contain upper case", cachedEip.Name, cachedEip.Spec.V6ip)
		klog.Error(err)
		return err
	}
	// eip change ip
	if c.eipChangeIP(cachedEip) {
		err := fmt.Errorf("not support eip change ip, old ip '%s', new ip '%s'", cachedEip.Status.IP, cachedEip.Spec.V4ip)
		klog.Error(err)
		return err
	}
	// make sure vpc nat enabled
	if vpcNatEnabled != "true" {
		err := errors.New("iptables nat gw not enable")
		klog.Error(err)
		return err
	}

	// update qos
	var desiredQoS *kubeovnv1.QoSPolicy
	if cachedEip.Status.QoSPolicy == cachedEip.Spec.QoSPolicy && cachedEip.Spec.QoSPolicy != "" {
		if desiredQoS, err = c.getAvailableQoSPolicy(cachedEip.Spec.QoSPolicy); err != nil {
			if cachedEip.Status.Ready {
				if patchErr := c.patchEipStatus(key, "", "", "", false); patchErr != nil {
					return fmt.Errorf("failed to mark eip %s not ready after its qos policy became unavailable: %w", key, patchErr)
				}
			}
			return err
		}
		if !controllerutil.ContainsFinalizer(desiredQoS, util.KubeOVNControllerFinalizer) {
			if cachedEip.Status.Ready {
				if patchErr := c.patchEipStatus(key, "", "", "", false); patchErr != nil {
					return fmt.Errorf("failed to mark eip %s not ready after its qos policy lost the controller finalizer: %w", key, patchErr)
				}
			}
			return fmt.Errorf("qos policy %s is not ready; wait for its first controller reconcile before referencing it", cachedEip.Spec.QoSPolicy)
		}
		uidMatches := cachedEip.Labels[util.QoSPolicyUIDLabel] == string(desiredQoS.UID)
		if uidMatches && !cachedEip.Status.Ready && cachedEip.Status.Redo == "" && cachedEip.Status.IP != "" {
			if _, err = c.getBindableQoSPolicy(cachedEip.Spec.QoSPolicy); err != nil {
				return err
			}
			if err = c.patchEipStatus(key, "", "", "", true); err != nil {
				return fmt.Errorf("failed to mark eip %s ready after its qos policy recovered: %w", key, err)
			}
			return nil
		}
	}
	qosUIDMismatch := desiredQoS != nil && cachedEip.Labels[util.QoSPolicyUIDLabel] != string(desiredQoS.UID)
	if cachedEip.Status.QoSPolicy != cachedEip.Spec.QoSPolicy || qosUIDMismatch {
		if _, err = c.getBindableQoSPolicy(cachedEip.Spec.QoSPolicy); err != nil {
			return err
		}
		if err = c.checkNatGwNotTerminating(cachedEip.Spec.NatGwDp); err != nil {
			return err
		}
		if qosUIDMismatch && cachedEip.Status.QoSPolicy == cachedEip.Spec.QoSPolicy {
			if cachedEip.Status.Ready {
				if err = c.patchEipStatus(key, "", "", "", false); err != nil {
					return fmt.Errorf("failed to mark eip %s not ready before qos policy generation rebind: %w", key, err)
				}
			}
			for _, direction := range []kubeovnv1.QoSPolicyRuleDirection{kubeovnv1.QoSDirectionIngress, kubeovnv1.QoSDirectionEgress} {
				if err = c.delEipQoSInPod(cachedEip.Spec.NatGwDp, cachedEip.Status.IP, c.natEipNamespace(cachedEip), direction); err != nil {
					return err
				}
			}
		} else if cachedEip.Status.QoSPolicy != "" {
			if err = c.delEipQoS(cachedEip, cachedEip.Status.IP); err != nil {
				klog.Errorf("failed to del qos '%s' in pod, %v", key, err)
				return err
			}
		}
		if cachedEip.Spec.QoSPolicy != "" {
			if err = c.addEipQoS(cachedEip, cachedEip.Status.IP); err != nil {
				klog.Errorf("failed to add qos '%s' in pod, %v", key, err)
				return err
			}
		}

		// The FIP/DNAT/SNAT rebind swaps its claim between the two pod operations; here it can only
		// come after both, because delEipQoS reads the rules to remove off the old policy and
		// releasing it early would let it be collected with its bandwidth rules still programmed.
		if err = c.patchEipLabel(key); err != nil {
			klog.Errorf("failed to label qos in eip, %v", err)
			return err
		}

		ready := cachedEip.Status.Ready || (cachedEip.Status.Redo == "" && cachedEip.Status.IP != "")
		if err = c.patchEipQoSStatus(key, cachedEip.Spec.QoSPolicy, ready); err != nil {
			klog.Errorf("failed to patch status for eip %s, %v", key, err)
			return err
		}
	}

	// redo
	if !cachedEip.Status.Ready &&
		cachedEip.Status.Redo != "" &&
		cachedEip.Status.IP != "" &&
		cachedEip.DeletionTimestamp.IsZero() {
		gwPod, err := c.getNatGwPod(cachedEip.Spec.NatGwDp, c.natEipNamespace(cachedEip))
		if err != nil {
			klog.Error(err)
			return err
		}
		// compare gw pod started time with eip redo time. if redo time before gw pod started. redo again
		eipRedo, _ := time.ParseInLocation("2006-01-02T15:04:05", cachedEip.Status.Redo, time.Local)
		if cachedEip.Status.Ready && cachedEip.Status.IP != "" && gwPod.Status.ContainerStatuses[0].State.Running.StartedAt.Before(&metav1.Time{Time: eipRedo}) {
			// already ok
			klog.V(3).Infof("eip %s already ok", key)
			return nil
		}
		addrV4, err := util.GetIPAddrWithMask(cachedEip.Status.IP, v4Cidr)
		if err != nil {
			err = fmt.Errorf("failed to get eip %s with mask by cidr %s: %w", cachedEip.Status.IP, v4Cidr, err)
			klog.Error(err)
			return err
		}
		if err = c.createEipInPod(cachedEip.Spec.NatGwDp, addrV4, c.natEipNamespace(cachedEip)); err != nil {
			klog.Errorf("failed to create eip, %v", err)
			return err
		}

		if cachedEip.Spec.QoSPolicy != "" {
			if err = c.addEipQoS(cachedEip, cachedEip.Status.IP); err != nil {
				klog.Errorf("failed to add qos '%s' in pod, %v", key, err)
				return err
			}
		}

		if err = c.patchEipStatus(key, "", "", cachedEip.Spec.QoSPolicy, true); err != nil {
			klog.Errorf("failed to patch status for eip %s, %v", key, err)
			return err
		}
	}
	if err = c.handleAddIptablesEipFinalizer(key); err != nil {
		klog.Errorf("failed to handle add finalizer for eip, %v", err)
		return err
	}
	return nil
}

func (c *Controller) handleDelIptablesEip(eip *kubeovnv1.IptablesEIP) error {
	klog.Infof("handle delete iptables eip %s", eip.Name)

	// For IptablesEIPs deleted without finalizer (race condition or direct deletion),
	// we need to ensure subnet status is updated as a safety net.
	externalNetwork := util.GetExternalNetwork(eip.Spec.ExternalSubnet)
	if externalNetwork != "" {
		c.updateSubnetStatusQueue.Add(externalNetwork)
	}

	return nil
}

func (c *Controller) GetEip(eipName string) (*kubeovnv1.IptablesEIP, error) {
	cachedEip, err := c.iptablesEipsLister.Get(eipName)
	if err != nil {
		klog.Errorf("failed to get eip %s, %v", eipName, err)
		return nil, err
	}
	if cachedEip.Status.IP == "" {
		return nil, fmt.Errorf("eip '%s' is not ready, has no v4ip", eipName)
	}
	eip := cachedEip.DeepCopy()
	return eip, nil
}

func (c *Controller) createEipInPod(dp, addrV4, ns string) error {
	gwPod, err := c.getNatGwPod(dp, ns)
	if err != nil {
		klog.Error(err)
		return err
	}
	return c.execNatGwRules(gwPod, natGwEipAdd, []string{addrV4})
}

// natGwDataPlaneGone reports whether the gateway has no data plane left to clean up, so an
// in-pod cleanup can be skipped instead of retried forever. That is the case when the
// VpcNatGateway no longer exists, when it is being deleted, and when none of its instances is
// running: the gateway data plane (iptables rules, nft maps, addresses) lives in the container's
// writable layer, so an instance that is not running holds no state, and a replacement instance
// starts empty and is programmed from the CRs that are still live at that point.
// Other errors are returned as-is for the caller to handle.
func (c *Controller) natGwDataPlaneGone(dp string) (bool, error) {
	gw, err := c.vpcNatGatewayLister.Get(dp)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}
	if !gw.DeletionTimestamp.IsZero() {
		return true, nil
	}
	pods, err := c.listNatGwPods(gw)
	if err != nil {
		return false, err
	}
	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodRunning && pod.DeletionTimestamp == nil {
			return false, nil
		}
	}
	klog.Infof("nat gw %s has no running instance, skipping data plane cleanup", dp)
	return true, nil
}

func (c *Controller) deleteEipInPod(dp, v4Cidr, ns string) error {
	// A gateway with no running instance holds no data plane to clean up: the rules live in the
	// container's writable layer, so a replacement instance starts empty and is programmed from
	// the live CRs. Only a gateway that is known to be running has to be reached.
	deleted, err := c.natGwDataPlaneGone(dp)
	if err != nil {
		klog.Error(err)
		return err
	}
	if deleted {
		return nil
	}
	gwPod, err := c.getNatGwPod(dp, ns)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Pod is temporarily absent (e.g. being recreated); retry quietly.
			klog.V(4).Infof("nat gw pod %s not found, will retry eip pod cleanup", dp)
		} else {
			klog.Error(err)
		}
		return err
	}
	var delRules []string
	rule := v4Cidr
	delRules = append(delRules, rule)
	if err = c.execNatGwRules(gwPod, natGwEipDel, delRules); err != nil {
		klog.Error(err)
		return err
	}
	return nil
}

func (c *Controller) addOrUpdateEIPBandwidthLimitRules(eip *kubeovnv1.IptablesEIP, v4ip string, rules kubeovnv1.QoSPolicyBandwidthLimitRules) error {
	var err error
	for _, rule := range rules {
		if err = c.addEipQoSInPod(eip.Spec.NatGwDp, v4ip, c.natEipNamespace(eip), rule.Direction, rule.Priority, rule.RateMax, rule.BurstMax); err != nil {
			klog.Errorf("failed to set %s eip '%s' qos in pod, %v", rule.Direction, eip.Name, err)
			return err
		}
	}
	return nil
}

// getBindableEip returns the EIP a NAT rule (fip/dnat/snat) wants to bind to, refusing the
// ones which are being deleted. Bindings established against a terminating EIP would keep its
// finalizer alive forever, since the EIP waits for all NAT rules referencing it to go away.
// This check and the caller's credential write are not atomic, see the binding rules in
// CODE_STYLE.md for the window that remains and why it is bounded rather than eliminated.
func (c *Controller) getBindableEip(eipName string) (*kubeovnv1.IptablesEIP, error) {
	eip, err := c.GetEip(eipName)
	if err != nil {
		return nil, err
	}
	if !eip.DeletionTimestamp.IsZero() {
		return nil, fmt.Errorf("eip %s is terminating, retry later", eipName)
	}
	if !eip.Status.Ready || eip.Status.IP == "" {
		return nil, fmt.Errorf("eip %s is not ready, retry later", eipName)
	}
	// The rules land in the gateway pod, so binding to a gateway on its way out only programs a
	// pod that is about to disappear.
	if err := c.checkNatGwNotTerminating(eip.Spec.NatGwDp); err != nil {
		return nil, err
	}
	return eip, nil
}

func (c *Controller) checkNatGwNotTerminating(gwName string) error {
	gw, err := c.vpcNatGatewayLister.Get(gwName)
	if err != nil {
		klog.Errorf("failed to get vpc nat gw %s, %v", gwName, err)
		return err
	}
	if !gw.DeletionTimestamp.IsZero() {
		return fmt.Errorf("vpc nat gw %s is terminating, retry later", gwName)
	}
	return nil
}

func (c *Controller) getAvailableQoSPolicy(qosPolicyName string) (*kubeovnv1.QoSPolicy, error) {
	if qosPolicyName == "" {
		return nil, nil
	}
	qosPolicy, err := c.qosPoliciesLister.Get(qosPolicyName)
	if err != nil {
		// A referenced policy that does not exist must not be reported as bindable: the caller
		// would stamp an empty UID credential the in-use check can never match.
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("qos policy %s does not exist; create it before referencing it: %w", qosPolicyName, err)
		}
		klog.Errorf("failed to get qos policy %s, %v", qosPolicyName, err)
		return nil, err
	}
	if !qosPolicy.DeletionTimestamp.IsZero() {
		return nil, fmt.Errorf("qos policy %s is terminating; wait for its deletion to complete before referencing it", qosPolicyName)
	}
	return qosPolicy, nil
}

func (c *Controller) getBindableQoSPolicy(qosPolicyName string) (*kubeovnv1.QoSPolicy, error) {
	qosPolicy, err := c.getAvailableQoSPolicy(qosPolicyName)
	if err != nil || qosPolicy == nil {
		return qosPolicy, err
	}
	if !qosPolicyStatusMatchesSpec(qosPolicy) {
		return nil, fmt.Errorf("qos policy %s is not ready; wait for its status to match the spec before referencing it", qosPolicyName)
	}
	if !controllerutil.ContainsFinalizer(qosPolicy, util.KubeOVNControllerFinalizer) {
		return nil, fmt.Errorf("qos policy %s is not ready; wait for its first controller reconcile before referencing it", qosPolicyName)
	}
	return qosPolicy, nil
}

// add tc rule for eip in nat gw pod
func (c *Controller) addEipQoS(eip *kubeovnv1.IptablesEIP, v4ip string) error {
	var err error
	// Reporting success for a missing policy would leave the EIP ready while its QoS rules are absent.
	qosPolicy, err := c.qosPoliciesLister.Get(eip.Spec.QoSPolicy)
	if err != nil {
		klog.Errorf("failed to get qos policy %s, %v", eip.Spec.QoSPolicy, err)
		return err
	}
	if !qosPolicy.Status.Shared {
		eips, err := c.iptablesEipsLister.List(
			labels.SelectorFromSet(labels.Set{util.QoSPolicyUIDLabel: string(qosPolicy.UID)}),
		)
		if err != nil {
			klog.Errorf("failed to get eip list, %v", err)
			return err
		}
		if len(eips) != 0 {
			if eips[0].Name != eip.Name {
				err := fmt.Errorf("not support unshared qos policy %s to related to multiple eip", eip.Spec.QoSPolicy)
				klog.Error(err)
				return err
			}
		}
	}
	return c.addOrUpdateEIPBandwidthLimitRules(eip, v4ip, qosPolicy.Status.BandwidthLimitRules)
}

func (c *Controller) delEIPBandwidthLimitRules(eip *kubeovnv1.IptablesEIP, v4ip string, rules kubeovnv1.QoSPolicyBandwidthLimitRules) error {
	var err error
	for _, rule := range rules {
		if err = c.delEipQoSInPod(eip.Spec.NatGwDp, v4ip, c.natEipNamespace(eip), rule.Direction); err != nil {
			klog.Errorf("failed to del %s eip '%s' qos in pod, %v", rule.Direction, eip.Name, err)
			return err
		}
	}
	return nil
}

// del tc rule for eip in nat gw pod
func (c *Controller) delEipQoS(eip *kubeovnv1.IptablesEIP, v4ip string) error {
	var err error
	qosPolicy, err := c.qosPoliciesLister.Get(eip.Status.QoSPolicy)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to get qos policy %s, %v", eip.Status.QoSPolicy, err)
		return err
	}

	return c.delEIPBandwidthLimitRules(eip, v4ip, qosPolicy.Status.BandwidthLimitRules)
}

func (c *Controller) addEipQoSInPod(
	dp, v4ip, ns string, direction kubeovnv1.QoSPolicyRuleDirection, priority int, rate string,
	burst string,
) error {
	if v4ip == "" {
		klog.Infof("v4ip is empty for nat gateway %s, skipping QoS rule addition", dp)
		return nil
	}
	var operation string
	gwPod, err := c.getNatGwPod(dp, ns)
	if err != nil {
		klog.Error(err)
		return err
	}
	var addRules []string
	rule := fmt.Sprintf("%s,%d,%s,%s", v4ip, priority, rate, burst)
	addRules = append(addRules, rule)

	switch direction {
	case kubeovnv1.QoSDirectionIngress:
		operation = natGwEipIngressQoSAdd
	case kubeovnv1.QoSDirectionEgress:
		operation = natGwEipEgressQoSAdd
	}

	return c.execNatGwRules(gwPod, operation, addRules)
}

func (c *Controller) delEipQoSInPod(dp, v4ip, ns string, direction kubeovnv1.QoSPolicyRuleDirection) error {
	var operation string
	// Same CRD / pod sentinel logic as deleteEipInPod: skip when the gateway is
	// gone, and a gateway with no running instance has no data plane left to clean up.
	deleted, err := c.natGwDataPlaneGone(dp)
	if err != nil {
		klog.Error(err)
		return err
	}
	if deleted {
		return nil
	}
	gwPod, err := c.getNatGwPod(dp, ns)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			klog.V(4).Infof("nat gw pod %s not found, will retry eip qos cleanup", dp)
		} else {
			klog.Error(err)
		}
		return err
	}
	var delRules []string
	delRules = append(delRules, v4ip)

	switch direction {
	case kubeovnv1.QoSDirectionIngress:
		operation = natGwEipIngressQoSDel
	case kubeovnv1.QoSDirectionEgress:
		operation = natGwEipEgressQoSDel
	}

	return c.execNatGwRules(gwPod, operation, delRules)
}

func (c *Controller) acquireStaticEip(name, _, nicName, ip, externalSubnet string) (string, string, string, error) {
	checkConflict := true
	var v4ip, v6ip, mac string
	var err error
	for ipStr := range strings.SplitSeq(ip, ",") {
		if net.ParseIP(ipStr) == nil {
			return "", "", "", fmt.Errorf("failed to parse eip ip %s", ipStr)
		}
	}

	if v4ip, v6ip, mac, err = c.ipam.GetStaticAddress(name, nicName, ip, nil, externalSubnet, checkConflict); err != nil {
		klog.Errorf("failed to get static ip %v, mac %v, subnet %v, err %v", ip, mac, externalSubnet, err)
		return "", "", "", err
	}
	return v4ip, v6ip, mac, nil
}

func (c *Controller) acquireEip(name, _, nicName, externalSubnet string) (string, string, string, error) {
	var skippedAddrs []string
	for {
		ipv4, ipv6, mac, err := c.ipam.GetRandomAddress(name, nicName, nil, externalSubnet, "", skippedAddrs, true)
		if err != nil {
			klog.Error(err)
			return "", "", "", err
		}

		ipv4OK, ipv6OK, err := c.validatePodIP(name, externalSubnet, ipv4, ipv6)
		if err != nil {
			klog.Error(err)
			return "", "", "", err
		}
		if ipv4OK && ipv6OK {
			return ipv4, ipv6, mac, nil
		}
		if !ipv4OK {
			skippedAddrs = append(skippedAddrs, ipv4)
		}
		if !ipv6OK {
			skippedAddrs = append(skippedAddrs, ipv6)
		}
	}
}

func (c *Controller) eipChangeIP(eip *kubeovnv1.IptablesEIP) bool {
	if eip.Status.IP == "" {
		// eip created but not ready
		return false
	}
	if eip.Status.IP != eip.Spec.V4ip {
		return true
	}
	return false
}

func (c *Controller) GetGwBySubnet(name string) (string, string, error) {
	subnet, err := c.subnetsLister.Get(name)
	if err != nil {
		err = fmt.Errorf("faile to get subnet %q: %w", name, err)
		klog.Error(err)
		return "", "", err
	}
	v4, v6 := util.SplitStringIP(subnet.Spec.Gateway)
	return v4, v6, nil
}

func (c *Controller) createOrUpdateEipCR(key, v4ip, v6ip, mac, natGwDp, qos, externalNet, gwNamespace string) error {
	qosPolicy, err := c.getBindableQoSPolicy(qos)
	if err != nil {
		return err
	}
	qosPolicyUID := ""
	if qosPolicy != nil {
		qosPolicyUID = string(qosPolicy.UID)
	}
	needCreate := false
	cachedEip, err := c.iptablesEipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			needCreate = true
		} else {
			klog.Errorf("failed to get eip %s, %v", key, err)
			return err
		}
	}
	if needCreate {
		klog.V(3).Infof("create eip cr %s", key)
		// Create CR with finalizer, labels and status all at once
		_, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Create(context.Background(), &kubeovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{
				Name:       key,
				Finalizers: []string{util.KubeOVNControllerFinalizer},
				Labels: map[string]string{
					util.SubnetNameLabel:        externalNet,
					util.EipV4IpLabel:           v4ip,
					util.VpcNatGatewayNameLabel: natGwDp,
					util.QoSLabel:               qos,
					util.QoSPolicyUIDLabel:      qosPolicyUID,
				},
			},
			Spec: kubeovnv1.IptablesEIPSpec{
				V4ip:           v4ip,
				V6ip:           v6ip,
				MacAddress:     mac,
				NatGwDp:        natGwDp,
				QoSPolicy:      qos,
				ExternalSubnet: externalNet,
				Namespace:      gwNamespace,
			},
			Status: kubeovnv1.IptablesEIPStatus{
				IP:        v4ip,
				Ready:     true,
				QoSPolicy: qos,
				Nat:       "",
				Redo:      "",
			},
		}, metav1.CreateOptions{})
		if err != nil {
			errMsg := fmt.Errorf("failed to create eip crd %s, %w", key, err)
			klog.Error(errMsg)
			return errMsg
		}
	} else {
		eip := cachedEip.DeepCopy()

		// Ensure labels are set correctly before any update
		if eip.Labels == nil {
			eip.Labels = make(map[string]string)
		}
		eip.Labels[util.SubnetNameLabel] = externalNet
		eip.Labels[util.VpcNatGatewayNameLabel] = natGwDp
		eip.Labels[util.EipV4IpLabel] = v4ip
		eip.Labels[util.QoSLabel] = qos
		eip.Labels[util.QoSPolicyUIDLabel] = qosPolicyUID
		if v4ip != "" {
			klog.V(3).Infof("update eip cr %s", key)
			eip.Spec.V4ip = v4ip
			eip.Spec.V6ip = v6ip
			eip.Spec.NatGwDp = natGwDp
			eip.Spec.MacAddress = mac
			eip.Spec.ExternalSubnet = externalNet
			// Auto-populate spec.namespace from VpcNatGateway for NAMESPACE column visibility.
			if eip.Spec.Namespace == "" && gwNamespace != "" {
				eip.Spec.Namespace = gwNamespace
			}
			// Update with labels and spec in one call
			if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Update(context.Background(), eip, metav1.UpdateOptions{}); err != nil {
				errMsg := fmt.Errorf("failed to update eip crd %s, %w", key, err)
				klog.Error(errMsg)
				return errMsg
			}
			if eip.Status.IP == "" {
				// eip is ip holder, not support change ip
				eip.Status.IP = v4ip
				// TODO:// ipv6
			}
			eip.Status.Ready = true
			eip.Status.QoSPolicy = qos
			bytes, err := eip.Status.Bytes()
			if err != nil {
				klog.Error(err)
				return err
			}
			if _, err = c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Patch(context.Background(), key, types.MergePatchType,
				bytes, metav1.PatchOptions{}, "status"); err != nil {
				if k8serrors.IsNotFound(err) {
					return nil
				}
				klog.Errorf("failed to patch eip %s, %v", eip.Name, err)
				return err
			}
		}

		if err = c.handleAddIptablesEipFinalizer(key); err != nil {
			klog.Errorf("failed to handle add or update finalizer for eip, %v", err)
			return err
		}
	}
	// Trigger subnet status update after all operations complete
	c.updateSubnetStatusQueue.AddAfter(externalNet, 300*time.Millisecond)
	return nil
}

func (c *Controller) syncIptablesEipFinalizer(cl client.Client) error {
	// migrate depreciated finalizer to new finalizer
	eips := &kubeovnv1.IptablesEIPList{}
	return migrateFinalizers(cl, eips, func(i int) (client.Object, client.Object) {
		if i < 0 || i >= len(eips.Items) {
			return nil, nil
		}
		return eips.Items[i].DeepCopy(), eips.Items[i].DeepCopy()
	})
}

func (c *Controller) handleAddIptablesEipFinalizer(key string) error {
	cachedIptablesEip, err := c.iptablesEipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if !cachedIptablesEip.DeletionTimestamp.IsZero() {
		return nil
	}
	newIptablesEip := cachedIptablesEip.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesEip, util.DepreciatedFinalizerName)
	controllerutil.AddFinalizer(newIptablesEip, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesEip, newIptablesEip)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables eip '%s', %v", cachedIptablesEip.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Patch(context.Background(), cachedIptablesEip.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to add finalizer for iptables eip '%s', %v", cachedIptablesEip.Name, err)
		return err
	}

	// Trigger subnet status update after finalizer is processed as a fallback
	// This handles cases where finalizer was not added during creation
	// AddFinalizer is idempotent, so this is safe even if finalizer already exists
	externalNetwork := util.GetExternalNetwork(cachedIptablesEip.Spec.ExternalSubnet)
	c.updateSubnetStatusQueue.Add(externalNetwork)
	return nil
}

func (c *Controller) handleDelIptablesEipFinalizer(key string) error {
	cachedIptablesEip, err := c.iptablesEipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if len(cachedIptablesEip.GetFinalizers()) == 0 {
		return nil
	}
	newIptablesEip := cachedIptablesEip.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesEip, util.DepreciatedFinalizerName)
	controllerutil.RemoveFinalizer(newIptablesEip, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesEip, newIptablesEip)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables eip '%s', %v", cachedIptablesEip.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Patch(context.Background(), cachedIptablesEip.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to remove finalizer from iptables eip '%s', %v", cachedIptablesEip.Name, err)
		return err
	}

	// Trigger subnet status update after finalizer is removed
	// This ensures subnet status reflects the IP release
	// Add delay to ensure API server completes the finalizer removal
	externalNetwork := util.GetExternalNetwork(cachedIptablesEip.Spec.ExternalSubnet)
	c.updateSubnetStatusQueue.AddAfter(externalNetwork, 300*time.Millisecond)
	return nil
}

func (c *Controller) patchEipQoSStatus(key, qos string, ready bool) error {
	var changed bool
	oriEip, err := c.iptablesEipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	eip := oriEip.DeepCopy()
	if eip.Status.Ready != ready {
		eip.Status.Ready = ready
		changed = true
	}

	// update status.qosPolicy
	if eip.Status.QoSPolicy != qos {
		eip.Status.QoSPolicy = qos
		changed = true
	}

	if changed {
		bytes, err := eip.Status.Bytes()
		if err != nil {
			klog.Error(err)
			return err
		}
		if _, err = c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Patch(context.Background(), key, types.MergePatchType,
			bytes, metav1.PatchOptions{}, "status"); err != nil {
			if k8serrors.IsNotFound(err) {
				return nil
			}
			klog.Errorf("failed to patch eip %s, %v", eip.Name, err)
			return err
		}
	}
	return nil
}

func (c *Controller) getIptablesEipNat(eip *kubeovnv1.IptablesEIP) (string, error) {
	nats := make([]string, 0, 3)
	selector := labels.SelectorFromSet(labels.Set{util.EipUIDLabel: string(eip.UID)})
	dnats, err := c.iptablesDnatRulesLister.List(selector)
	if err != nil {
		klog.Errorf("failed to get dnats, %v", err)
		return "", err
	}
	if len(dnats) != 0 {
		nats = append(nats, util.DnatUsingEip)
	}
	fips, err := c.iptablesFipsLister.List(selector)
	if err != nil {
		klog.Errorf("failed to get fips, %v", err)
		return "", err
	}
	if len(fips) != 0 {
		nats = append(nats, util.FipUsingEip)
	}
	snats, err := c.iptablesSnatRulesLister.List(selector)
	if err != nil {
		klog.Errorf("failed to get snats, %v", err)
		return "", err
	}
	if len(snats) != 0 {
		nats = append(nats, util.SnatUsingEip)
	}
	nat := strings.Join(nats, ",")
	return nat, nil
}

// getIptablesEipNatFromAPI reports which NAT rules reference the EIP, reading from the API server
// rather than the informer cache. The finalizer release decision must not run on a cache that has
// not observed a freshly written reference yet, or the EIP is dropped while rules still claim it.
func (c *Controller) getIptablesEipNatFromAPI(eip *kubeovnv1.IptablesEIP) (string, error) {
	// Do not set a Limit here: with a label selector the server may return zero items alongside a
	// continue token, and this result decides whether the EIP's finalizer is released.
	opts := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{util.EipUIDLabel: string(eip.UID)}).String(),
	}
	client := c.config.KubeOvnClient.KubeovnV1()
	nats := make([]string, 0, 3)
	dnats, err := client.IptablesDnatRules().List(context.Background(), opts)
	if err != nil {
		klog.Errorf("failed to list dnats, %v", err)
		return "", err
	}
	if len(dnats.Items) != 0 {
		nats = append(nats, util.DnatUsingEip)
	}
	fips, err := client.IptablesFIPRules().List(context.Background(), opts)
	if err != nil {
		klog.Errorf("failed to list fips, %v", err)
		return "", err
	}
	if len(fips.Items) != 0 {
		nats = append(nats, util.FipUsingEip)
	}
	snats, err := client.IptablesSnatRules().List(context.Background(), opts)
	if err != nil {
		klog.Errorf("failed to list snats, %v", err)
		return "", err
	}
	if len(snats.Items) != 0 {
		nats = append(nats, util.SnatUsingEip)
	}
	return strings.Join(nats, ","), nil
}

func (c *Controller) patchEipStatus(key, v4ip, redo, qos string, ready bool) error {
	oriEip, err := c.iptablesEipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	eip := oriEip.DeepCopy()
	var changed bool
	if eip.Status.Ready != ready {
		eip.Status.Ready = ready
		changed = true
	}

	if redo != "" && eip.Status.Redo != redo {
		eip.Status.Redo = redo
		changed = true
	}

	if ready && v4ip != "" && eip.Status.IP != v4ip {
		eip.Status.IP = v4ip
		changed = true
	}

	nat, err := c.getIptablesEipNat(oriEip)
	if err != nil {
		err = fmt.Errorf("failed to get eip nat: %w", err)
		klog.Error(err)
		return err
	}
	// nat record all kinds of nat rules using this eip
	klog.V(3).Infof("nat of eip %s is %s", eip.Name, nat)
	if eip.Status.Nat != nat {
		eip.Status.Nat = nat
		changed = true
	}

	if qos != "" && eip.Status.QoSPolicy != qos {
		eip.Status.QoSPolicy = qos
		changed = true
	}

	if changed {
		bytes, err := eip.Status.Bytes()
		if err != nil {
			klog.Error(err)
			return err
		}
		if _, err = c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Patch(context.Background(), key, types.MergePatchType,
			bytes, metav1.PatchOptions{}, "status"); err != nil {
			if k8serrors.IsNotFound(err) {
				return nil
			}
			klog.Errorf("failed to patch eip %s, %v", eip.Name, err)
			return err
		}
	}
	return nil
}

func (c *Controller) patchEipLabel(eipName string) error {
	oriEip, err := c.iptablesEipsLister.Get(eipName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	externalNetwork := util.GetExternalNetwork(oriEip.Spec.ExternalSubnet)

	// patchEipLabel writes util.QoSLabel and util.QoSPolicyUIDLabel on the update/redo path, and the
	// QoS policy in-use check is keyed on the UID label. Refuse to point a live EIP at a terminating
	// QoS policy: handleResetIptablesEip runs from a 3s delayed queue keyed by the EIP name, so it
	// can otherwise attach the tombstone of a previous instance to a freshly recreated EIP of the
	// same name.
	// Only a reference that would actually change is rejected. Rewriting the label the EIP already
	// carries adds nothing to the in-use count, and failing it would abort the whole gateway's redo
	// loop in handleUpdateVpcFloatingIP, which gives up on the first error and would leave unrelated
	// FIPs unapplied after a gateway pod restart.
	// A terminating EIP is exempt too: it establishes no new binding, its own deletion is unblocked
	// by the status patch handleResetIptablesEip issues right after this call, and the QoS policy is
	// released anyway once the EIP DeleteFunc fires.
	qosPolicy, qosErr := c.getBindableQoSPolicy(oriEip.Spec.QoSPolicy)
	if qosErr != nil && oriEip.DeletionTimestamp.IsZero() && oriEip.Labels[util.QoSLabel] != oriEip.Spec.QoSPolicy {
		return qosErr
	}
	// On an exempt path the policy is unreadable, so keep the recorded UID: clearing it would drop
	// this reference from the in-use count while the EIP still points at the policy.
	qosPolicyUID := oriEip.Labels[util.QoSPolicyUIDLabel]
	if qosErr == nil {
		qosPolicyUID = ""
		if qosPolicy != nil {
			qosPolicyUID = string(qosPolicy.UID)
		}
	}

	eip := oriEip.DeepCopy()
	var needUpdateLabel bool
	var op string
	if len(eip.Labels) == 0 {
		op = "add"
		needUpdateLabel = true
		eip.Labels = map[string]string{
			util.SubnetNameLabel:        externalNetwork,
			util.VpcNatGatewayNameLabel: eip.Spec.NatGwDp,
			util.QoSLabel:               eip.Spec.QoSPolicy,
			util.QoSPolicyUIDLabel:      qosPolicyUID,
			util.EipV4IpLabel:           eip.Spec.V4ip,
		}
	} else if eip.Labels[util.VpcNatGatewayNameLabel] != eip.Spec.NatGwDp || eip.Labels[util.QoSLabel] != eip.Spec.QoSPolicy ||
		eip.Labels[util.QoSPolicyUIDLabel] != qosPolicyUID || eip.Labels[util.SubnetNameLabel] != externalNetwork ||
		eip.Labels[util.EipV4IpLabel] != eip.Spec.V4ip {
		op = "replace"
		needUpdateLabel = true
		eip.Labels[util.SubnetNameLabel] = externalNetwork
		eip.Labels[util.VpcNatGatewayNameLabel] = eip.Spec.NatGwDp
		eip.Labels[util.QoSLabel] = eip.Spec.QoSPolicy
		eip.Labels[util.QoSPolicyUIDLabel] = qosPolicyUID
		eip.Labels[util.EipV4IpLabel] = eip.Spec.V4ip
	}
	if needUpdateLabel {
		if err := c.updateIptableLabels(eip.Name, op, "eip", eip.Labels); err != nil {
			klog.Errorf("failed to update label of eip %s, %v", eip.Name, err)
			return err
		}
	}
	return nil
}
