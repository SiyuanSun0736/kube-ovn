package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func (c *Controller) enqueueAddIptablesFip(obj any) {
	fip := obj.(*kubeovnv1.IptablesFIPRule)
	key := cache.MetaObjectToName(fip).String()
	// A terminating object reconciles via the update queue for cleanup (handleAdd returns early).
	if enqueueUpdateIfTerminating(c.updateIptablesFipQueue, key, "fip", fip.DeletionTimestamp) {
		return
	}
	klog.V(3).Infof("enqueue add iptables fip %s", key)
	c.addIptablesFipQueue.Add(key)
}

func (c *Controller) enqueueUpdateIptablesFip(oldObj, newObj any) {
	oldFip := oldObj.(*kubeovnv1.IptablesFIPRule)
	newFip := newObj.(*kubeovnv1.IptablesFIPRule)
	key := cache.MetaObjectToName(newFip).String()
	if !newFip.DeletionTimestamp.IsZero() {
		klog.V(3).Infof("enqueue update to clean fip %s", key)
		c.updateIptablesFipQueue.Add(key)
		return
	}
	if newFip.Spec.EIP == "" || newFip.Spec.InternalIP == "" {
		klog.Warningf("skip enqueue fip %s: incomplete spec (eip=%q, internalIP=%q)", key, newFip.Spec.EIP, newFip.Spec.InternalIP)
		return
	}
	if oldFip.Status.V4ip != newFip.Status.V4ip ||
		oldFip.Spec.EIP != newFip.Spec.EIP ||
		oldFip.Status.Redo != newFip.Status.Redo ||
		oldFip.Spec.InternalIP != newFip.Spec.InternalIP {
		klog.V(3).Infof("enqueue update fip %s", key)
		c.updateIptablesFipQueue.Add(key)
		return
	}
}

func (c *Controller) enqueueDelIptablesFip(obj any) {
	var fip *kubeovnv1.IptablesFIPRule
	switch t := obj.(type) {
	case *kubeovnv1.IptablesFIPRule:
		fip = t
	case cache.DeletedFinalStateUnknown:
		f, ok := t.Obj.(*kubeovnv1.IptablesFIPRule)
		if !ok {
			klog.Warningf("unexpected object type: %T", t.Obj)
			return
		}
		fip = f
	default:
		klog.Warningf("unexpected type: %T", obj)
		return
	}

	key := cache.MetaObjectToName(fip).String()
	klog.V(3).Infof("enqueue delete iptables fip %s", key)
	c.delIptablesFipQueue.Add(key)
}

// enqueueIptablesEipRecheck re-runs an EIP's reconcile so it can finish deleting once the rules
// that referenced it are gone. Share records are written and deleted by the Service reconcile
// (never by the DNAT worker), so the EIP has to be woken explicitly when one stops referencing it,
// otherwise its finalizer waits for rules that are already gone.
func (c *Controller) enqueueIptablesEipRecheck(eipName string) {
	if eipName == "" || c.updateIptablesEipQueue == nil {
		return
	}
	klog.V(3).Infof("re-check iptables eip %s after its referencing rule changed", eipName)
	c.updateIptablesEipQueue.Add(eipName)
}

func (c *Controller) enqueueAddIptablesDnatRule(obj any) {
	dnat := obj.(*kubeovnv1.IptablesDnatRule)
	if dnat.Spec.Type == kubeovnv1.DnatRuleTypeShare {
		return
	}
	key := cache.MetaObjectToName(dnat).String()
	// A terminating object reconciles via the update queue for cleanup (handleAdd returns early).
	if enqueueUpdateIfTerminating(c.updateIptablesDnatRuleQueue, key, "dnat", dnat.DeletionTimestamp) {
		return
	}
	klog.V(3).Infof("enqueue add iptables dnat %s", key)
	c.addIptablesDnatRuleQueue.Add(key)
}

func (c *Controller) enqueueUpdateIptablesDnatRule(oldObj, newObj any) {
	oldDnat := oldObj.(*kubeovnv1.IptablesDnatRule)
	newDnat := newObj.(*kubeovnv1.IptablesDnatRule)
	if oldDnat.Spec.Type == kubeovnv1.DnatRuleTypeShare || newDnat.Spec.Type == kubeovnv1.DnatRuleTypeShare {
		// A share record is the Service's accounting object: the DNAT worker never touches it, but
		// the EIPs it referenced may become free (or newly referenced) with this change.
		if oldDnat.Spec.EIP != newDnat.Spec.EIP {
			c.enqueueIptablesEipRecheck(oldDnat.Spec.EIP)
			c.enqueueIptablesEipRecheck(newDnat.Spec.EIP)
		}
		return
	}
	key := cache.MetaObjectToName(newDnat).String()
	if !newDnat.DeletionTimestamp.IsZero() {
		klog.V(3).Infof("enqueue update to clean dnat %s", key)
		c.updateIptablesDnatRuleQueue.Add(key)
		return
	}
	if newDnat.Spec.Protocol != util.ProtocolTCP && newDnat.Spec.Protocol != util.ProtocolUDP {
		klog.Warningf("enqueue invalid dnat %s for protocol validation: %q", key, newDnat.Spec.Protocol)
		c.updateIptablesDnatRuleQueue.Add(key)
		return
	}
	// The rule's identity is its address, and a rule that serves a ClusterIP carries no EIP: only a
	// rule with neither address has nothing to reconcile. This enqueue is the redo's only entry
	// point (redoDnat patches the status and nothing else), so rejecting an EIP-less rule here would
	// leave the ClusterIP identity, its hairpin rule and its lo address unprogrammed on a gateway
	// instance that replaces the one they were programmed on.
	if (newDnat.Spec.EIP == "" && newDnat.Spec.ClusterIP == "") || newDnat.Spec.ExternalPort == "" ||
		newDnat.Spec.InternalIP == "" || newDnat.Spec.InternalPort == "" {
		klog.Warningf("skip enqueue dnat %s: incomplete spec (eip=%q, clusterIP=%q, externalPort=%q, protocol=%q, internalIP=%q, internalPort=%q)",
			key, newDnat.Spec.EIP, newDnat.Spec.ClusterIP, newDnat.Spec.ExternalPort, newDnat.Spec.Protocol, newDnat.Spec.InternalIP, newDnat.Spec.InternalPort)
		return
	}
	if oldDnat.Labels[util.VpcNatGatewayNameLabel] != newDnat.Labels[util.VpcNatGatewayNameLabel] ||
		oldDnat.Status.V4ip != newDnat.Status.V4ip ||
		oldDnat.Spec.EIP != newDnat.Spec.EIP ||
		oldDnat.Status.Redo != newDnat.Status.Redo ||
		oldDnat.Spec.Protocol != newDnat.Spec.Protocol ||
		oldDnat.Spec.InternalIP != newDnat.Spec.InternalIP ||
		oldDnat.Spec.ExternalPort != newDnat.Spec.ExternalPort ||
		oldDnat.Spec.InternalPort != newDnat.Spec.InternalPort {
		klog.V(3).Infof("enqueue update dnat %s", key)
		c.updateIptablesDnatRuleQueue.Add(key)
		return
	}
}

func (c *Controller) enqueueDelIptablesDnatRule(obj any) {
	var dnat *kubeovnv1.IptablesDnatRule
	switch t := obj.(type) {
	case *kubeovnv1.IptablesDnatRule:
		dnat = t
	case cache.DeletedFinalStateUnknown:
		d, ok := t.Obj.(*kubeovnv1.IptablesDnatRule)
		if !ok {
			klog.Warningf("unexpected object type: %T", t.Obj)
			return
		}
		dnat = d
	default:
		klog.Warningf("unexpected type: %T", obj)
		return
	}

	if dnat.Spec.Type == kubeovnv1.DnatRuleTypeShare {
		// The Service released this record; the EIP it referenced has to be re-checked so its
		// finalizer can clear now that no rule uses it any more.
		c.enqueueIptablesEipRecheck(dnat.Spec.EIP)
		// A Service accounting record carries no finalizer, so nothing stops it from being
		// deleted out from under its Service while the gateway keeps wiring the identity. The
		// Service reconcile is the only writer that restores the claim, and the EIP release
		// holds until it has settled, so wake the Service here.
		if util.IsNftableLbSvcRecord(dnat.Labels) && dnat.Labels[util.NftableLbSvcNameLabel] != "" {
			c.enqueueGwNftableLbService(dnat.Labels[util.NftableLbSvcNsLabel] + "/" + dnat.Labels[util.NftableLbSvcNameLabel])
		}
		return
	}
	key := cache.MetaObjectToName(dnat).String()
	klog.V(3).Infof("enqueue delete iptables dnat %s", key)
	c.delIptablesDnatRuleQueue.Add(key)
}

func (c *Controller) enqueueAddIptablesSnatRule(obj any) {
	snat := obj.(*kubeovnv1.IptablesSnatRule)
	key := cache.MetaObjectToName(snat).String()
	// A terminating object reconciles via the update queue for cleanup (handleAdd returns early).
	if enqueueUpdateIfTerminating(c.updateIptablesSnatRuleQueue, key, "snat", snat.DeletionTimestamp) {
		return
	}
	klog.V(3).Infof("enqueue add iptables snat %s", key)
	c.addIptablesSnatRuleQueue.Add(key)
}

func (c *Controller) enqueueUpdateIptablesSnatRule(oldObj, newObj any) {
	oldSnat := oldObj.(*kubeovnv1.IptablesSnatRule)
	newSnat := newObj.(*kubeovnv1.IptablesSnatRule)
	key := cache.MetaObjectToName(newSnat).String()
	if !newSnat.DeletionTimestamp.IsZero() {
		klog.V(3).Infof("enqueue update to clean snat %s", key)
		c.updateIptablesSnatRuleQueue.Add(key)
		return
	}
	if newSnat.Spec.EIP == "" || newSnat.Spec.InternalCIDR == "" {
		klog.Warningf("skip enqueue snat %s: incomplete spec (eip=%q, internalCIDR=%q)", key, newSnat.Spec.EIP, newSnat.Spec.InternalCIDR)
		return
	}
	if oldSnat.Status.V4ip != newSnat.Status.V4ip ||
		oldSnat.Spec.EIP != newSnat.Spec.EIP ||
		oldSnat.Status.Redo != newSnat.Status.Redo ||
		oldSnat.Spec.InternalCIDR != newSnat.Spec.InternalCIDR {
		klog.V(3).Infof("enqueue update snat %s", key)
		c.updateIptablesSnatRuleQueue.Add(key)
		return
	}
}

func (c *Controller) enqueueDelIptablesSnatRule(obj any) {
	var snat *kubeovnv1.IptablesSnatRule
	switch t := obj.(type) {
	case *kubeovnv1.IptablesSnatRule:
		snat = t
	case cache.DeletedFinalStateUnknown:
		s, ok := t.Obj.(*kubeovnv1.IptablesSnatRule)
		if !ok {
			klog.Warningf("unexpected object type: %T", t.Obj)
			return
		}
		snat = s
	default:
		klog.Warningf("unexpected type: %T", obj)
		return
	}

	key := cache.MetaObjectToName(snat).String()
	klog.V(3).Infof("enqueue delete iptables snat %s", key)
	c.delIptablesSnatRuleQueue.Add(key)
}

// handleAddIptablesFip creates a FIP rule from scratch.
//
// Responsibility:
//   - Bring the FIP from an empty Status to a fully consistent state across all 4 dimensions:
//     1. iptables rule in NAT GW Pod  2. FIP Status  3. FIP Labels  4. EIP Status
//   - On success, Status MUST be fully populated (V4ip, InternalIP, NatGwDp, Ready=true).
//     This is the contract that handleUpdateIptablesFip relies on: a complete Status means
//     the old values are reliable and can be used for spec-change detection and rule deletion.
//   - On failure at any step, returns error to retry. Partial state (e.g., iptables rule
//     created but Status not yet patched) may exist; the next retry uses current Spec
//     and must converge to the desired state.
//
// Error state: if this handler never completes (Status.V4ip stays empty), the resource is
// in an error state. The update handler MUST NOT attempt NAT operations in this state —
// it lacks reliable old values for safe rule replacement.
func (c *Controller) handleAddIptablesFip(key string) error {
	fip, err := c.iptablesFipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	// The key may have been queued while the object was still live; the update queue owns cleanup.
	if !fip.DeletionTimestamp.IsZero() {
		return nil
	}

	if vpcNatEnabled != "true" {
		return errors.New("iptables nat gw not enable")
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle add iptables fip %s", key)

	if fip.Status.V4ip != "" && fip.Status.NatGwDp != "" && fip.Status.InternalIP != "" {
		eip, getErr := c.getBindableEip(fip.Spec.EIP)
		if getErr != nil || fip.Status.V4ip != eip.Status.IP || fip.Status.NatGwDp != eip.Spec.NatGwDp ||
			fip.Status.InternalIP != fip.Spec.InternalIP || fip.Labels[util.EipUIDLabel] != string(eip.UID) {
			c.updateIptablesFipQueue.Add(key)
			return nil
		}
		if fip.Status.Ready {
			return nil
		}
	}
	klog.V(3).Infof("handle add fip %s", key)

	if err := c.validateFipRule(fip); err != nil {
		return err
	}

	eip, err := c.getBindableEip(fip.Spec.EIP)
	if err != nil {
		klog.Errorf("failed to get eip, %v", err)
		return err
	}

	if err = c.fipTryUseEip(key, eip); err != nil {
		err = fmt.Errorf("failed to create fip %s, %w", key, err)
		klog.Error(err)
		return err
	}

	// we add the finalizer **before** we run "createFipInPod". This is because if we
	// added the finalizer after, then it is possible that the FIP is deleted after
	// we run createFipInPod but before the finalizer is created, and
	// then we can be left with IPtables rules in the VPC Nat
	// Gateway pod which are unmanaged.
	if err = c.handleAddIptablesFipFinalizer(key); err != nil {
		klog.Errorf("failed to handle add finalizer for fip, %v", err)
		return err
	}

	// Claim the EIP before touching the gateway pod: the EIP in-use check counts this label, so a
	// claim written only after the rules exist can be missed by a concurrent EIP release.
	if err = c.patchFipLabel(key, eip); err != nil {
		klog.Errorf("failed to update label for fip %s, %v", key, err)
		return err
	}

	if err = c.createFipInPod(eip.Spec.NatGwDp, eip.Status.IP, fip.Spec.InternalIP); err != nil {
		klog.Errorf("failed to create fip, %v", err)
		return err
	}
	if err = c.patchFipStatus(key, eip.Status.IP, eip.Spec.V6ip, eip.Spec.NatGwDp, "", true); err != nil {
		klog.Errorf("failed to patch status for fip %s, %v", key, err)
		return err
	}
	if err = c.patchEipStatus(fip.Spec.EIP, "", "", "", true); err != nil {
		// refresh eip nats
		klog.Errorf("failed to patch fip use eip %s, %v", key, err)
		return err
	}
	// patchFipLabel updated the FIP's EipUIDLabel via the API, but the informer cache
	// may not have synced yet when patchEipStatus called getIptablesEipNat above, causing
	// it to miss the FIP and leave EIP.Status.Nat stale. Schedule a delayed reset.
	c.resetIptablesEipQueue.AddAfter(fip.Spec.EIP, 3*time.Second)
	return nil
}

func (c *Controller) fipTryUseEip(fipName string, eip *kubeovnv1.IptablesEIP) error {
	// check if has another fip using this eip already
	selector := labels.SelectorFromSet(labels.Set{util.EipUIDLabel: string(eip.UID)})
	usingFips, err := c.iptablesFipsLister.List(selector)
	if err != nil {
		klog.Errorf("failed to get fips, %v", err)
		return err
	}
	for _, uf := range usingFips {
		if uf.Name != fipName {
			err = fmt.Errorf("%s is using by the other fip %s", eip.Status.IP, uf.Name)
			klog.Error(err)
			return err
		}
	}
	return nil
}

// handleUpdateIptablesFip handles FIP deletion, spec changes, and redo.
//
// FIP involves 4 dimensions of data that must be kept consistent:
//
//	Dimension        Storage                Description
//	──────────────── ────────────────────── ──────────────────────────────────────────
//	1. iptables rule NAT GW Pod             v4ip <-> internalIP DNAT/SNAT pair
//	2. FIP Status    CR .status             V4ip, InternalIP, NatGwDp, Ready
//	3. FIP Label     CR .labels/.annotations EipV4IpLabel, VpcNatGatewayNameLabel, VpcEipAnnotation
//	4. EIP Status    EIP CR .status         Nat field records which NAT rules use this EIP
//
// Precondition for spec-change path:
//   - Status.V4ip != "" (Status is complete, populated by handleAddIptablesFip).
//     A complete Status means the old values are reliable and reflect what was actually
//     created in the Pod. This is the ONLY safe basis for deleting old iptables rules.
//   - If Status.V4ip == "", the resource is in an error state (handleAdd never completed).
//     The update handler MUST NOT attempt any NAT operations — it has no reliable old
//     values to delete and creating new rules could leave stale data.
//     It should log the error and return, leaving recovery to the add handler (or future updates that fix the spec).
//     IMPORTANT: If a change fails, we must wait for retry until success. Continuing to change spec
//     on top of a failed state (dirty status) can introduce more zombie rules that are hard to track.
//
// Paths:
//
//  1. Delete path (DeletionTimestamp set):
//     Uses Status (with Spec fallback for best-effort cleanup) to locate the old iptables rule.
//     Removes finalizer, then async-resets EIP to refresh dimension 4.
//
//  2. Spec change path (Status.V4ip != "" AND Status differs from Spec+EIP):
//     Old values come from Status; new values come from Spec + EIP CR.
//     Steps strictly ordered to maintain all 4 dimensions:
//     a. patchFipStatus(ready=false)  (mark dimension 2 dirty; crash leaves a visible not-ready signal)
//     b. finalDeleteFipInPod  (clean dimension 1 with old values from Status)
//     c. patchFipLabel   (swap dimension 3 once the old rule is gone, before the new one exists)
//     d. createFipInPod  (create dimension 1 with new values from Spec+EIP)
//     e. patchFipStatus  (update dimension 2 to match new iptables rule, mark ready=true)
//     f. patchEipStatus  (update dimension 4 on new EIP)
//     g. resetOldEip     (async clean dimension 4 on old EIP, if EIP changed)
//     If any step fails, returns error to retry. Iptables operations are idempotent.
//
//  3. Redo path (gateway pod restarted):
//     Re-creates the iptables rule using Status.V4ip (the known-good external IP)
//     and Status.InternalIP (the destination IP recorded when the rule was originally created).
//     Only touches dimension 1; dimensions 2-4 are already correct.
func (c *Controller) handleUpdateIptablesFip(key string) error {
	cachedFip, err := c.iptablesFipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle update iptables fip %s", key)

	// should delete
	if !cachedFip.DeletionTimestamp.IsZero() {
		if vpcNatEnabled == "true" {
			if err = c.finalDeleteFipInPod(key, cachedFip); err != nil {
				return err
			}
		}
		if err = c.handleDelIptablesFipFinalizer(key); err != nil {
			klog.Errorf("failed to handle del finalizer for fip, %v", err)
			return err
		}
		//  reset eip
		c.resetIptablesEipQueue.AddAfter(cachedFip.Spec.EIP, 3*time.Second)
		return nil
	}
	klog.V(3).Infof("handle update fip %s", key)
	// add or update should make sure vpc nat enabled
	if vpcNatEnabled != "true" {
		if released, releaseErr := c.releaseDeletedEipRef(
			key, util.FipUsingEip, cachedFip.Spec.EIP, cachedFip.Labels, cachedFip.Annotations,
			func() error { return c.finalDeleteFipInPod(key, cachedFip) },
			func() error { return c.patchFipStatus(key, "", "", "", "", false) },
		); released || releaseErr != nil {
			return releaseErr
		}
		return errors.New("iptables nat gw not enable")
	}

	if err := c.validateFipRule(cachedFip); err != nil {
		return err
	}

	eip, err := c.getBindableEip(cachedFip.Spec.EIP)
	if err != nil {
		klog.Errorf("failed to get eip, %v", err)
		if released, releaseErr := c.releaseDeletedEipRef(
			key, util.FipUsingEip, cachedFip.Spec.EIP, cachedFip.Labels, cachedFip.Annotations,
			func() error { return c.finalDeleteFipInPod(key, cachedFip) },
			func() error { return c.patchFipStatus(key, "", "", "", "", false) },
		); released || releaseErr != nil {
			return releaseErr
		}
		if cachedFip.Status.Ready {
			if patchErr := c.patchFipStatus(key, "", "", "", "", false); patchErr != nil {
				return fmt.Errorf("failed to mark fip %s not ready after its eip became unavailable: %w", key, patchErr)
			}
		}
		return err
	}

	if err = c.fipTryUseEip(key, eip); err != nil {
		err = fmt.Errorf("failed to update fip %s, %w", key, err)
		klog.Error(err)
		return err
	}

	if eip.Spec.NatGwDp == "" {
		klog.Errorf("fip %s: eip %s has empty NatGwDp, skip binding", key, cachedFip.Spec.EIP)
		return nil
	}

	// Error state: Status incomplete means handleAdd never completed.
	// Do not attempt NAT operations — no reliable old values for safe rule replacement.
	// All spec-corresponding Status fields must be populated for the update handler to proceed.
	if cachedFip.Status.V4ip == "" || cachedFip.Status.NatGwDp == "" || cachedFip.Status.InternalIP == "" {
		klog.Errorf("fip %s has incomplete status (V4ip=%q, NatGwDp=%q, InternalIP=%q), skipping NAT operations; waiting for add handler to complete",
			key, cachedFip.Status.V4ip, cachedFip.Status.NatGwDp, cachedFip.Status.InternalIP)
		return nil
	}

	// spec change: compare Status (old, what's in Pod) vs Spec+EIP (new, desired)
	oldV4ip := cachedFip.Status.V4ip
	newV4ip := eip.Status.IP
	newInternalIP := cachedFip.Spec.InternalIP

	// Warn if we are modifying a resource that might be in a dirty state from a previous failed update.
	if !cachedFip.Status.Ready {
		// TODO: consider using a webhook to reject spec changes when the resource is not ready,
		// to prevent users from modifying the spec when the previous update has not yet been fully reconciled.
		klog.Warningf("fip %s is being updated while not Ready (previous update likely failed). This may lead to stale iptables rules.", key)
	}

	// Verify new parameters are valid before modifying any state.
	// eip.Status.IP can be empty if EIP itself is in error state.
	if newV4ip == "" || newInternalIP == "" {
		klog.Errorf("skipping fip %s update: incomplete new parameters (v4ip=%q, internalIP=%q)", key, newV4ip, newInternalIP)
		return nil
	}

	if oldV4ip != newV4ip || cachedFip.Status.NatGwDp != eip.Spec.NatGwDp ||
		cachedFip.Status.InternalIP != newInternalIP || cachedFip.Labels[util.EipUIDLabel] != string(eip.UID) {
		// Mark FIP as not ready before starting the update.
		// This ensures that if the controller crashes or the update fails midway,
		// the resource will be left in a non-ready state, indicating a potential inconsistency.
		if cachedFip.Status.Ready {
			klog.V(3).Infof("fip %s spec changed, marking as not ready before update", key)
			if err = c.patchFipStatus(key, oldV4ip, cachedFip.Status.V6ip, cachedFip.Status.NatGwDp, "", false); err != nil {
				klog.Errorf("failed to mark fip %s as not ready, %v", key, err)
				return err
			}
		}
		// delete old rule; finalDeleteFipInPod resolves (natGwDp, v4ip) from Status
		if err = c.finalDeleteFipInPod(key, cachedFip); err != nil {
			return err
		}
		// Swap the claim between the two pod operations: the old EIP stays claimed until its rule is
		// gone, and the new one is claimed before its rule exists. The in-use check counts this label,
		// so either edge would let a concurrent release drop a finalizer with a live rule behind it.
		if err = c.patchFipLabel(key, eip); err != nil {
			klog.Errorf("failed to update label for fip %s, %v", key, err)
			return err
		}
		c.enqueueDeletingOldIptablesEip(cachedFip.Annotations[util.VpcEipAnnotation], cachedFip.Spec.EIP)
		if err = c.createFipInPod(eip.Spec.NatGwDp, newV4ip, newInternalIP); err != nil {
			klog.Errorf("failed to create fip %s, %v", key, err)
			return err
		}
		if err = c.patchFipStatus(key, newV4ip, eip.Spec.V6ip, eip.Spec.NatGwDp, "", true); err != nil {
			klog.Errorf("failed to patch status for fip %s, %v", key, err)
			return err
		}
		if err = c.patchEipStatus(cachedFip.Spec.EIP, "", "", "", true); err != nil {
			klog.Errorf("failed to patch fip use eip %s, %v", key, err)
			return err
		}
		// Same refresh the add path schedules: patchEipStatus read the new EIP's nat off the
		// informer, which may not carry the label patched moments ago.
		c.resetIptablesEipQueue.AddAfter(cachedFip.Spec.EIP, 3*time.Second)
		// Reset old EIP after all 4 dimensions are updated.
		// This must NOT be in enqueue — async reset there could race with handler's
		// patchEipStatus, causing the old EIP's nat label to be cleared before the
		// new EIP is fully bound, leaving a window where the old EIP appears free.
		if oldEipName := cachedFip.Annotations[util.VpcEipAnnotation]; oldEipName != "" && oldEipName != cachedFip.Spec.EIP {
			c.resetIptablesEipQueue.AddAfter(oldEipName, 3*time.Second)
		}
		if err = c.handleAddIptablesFipFinalizer(key); err != nil {
			klog.Errorf("failed to handle add finalizer for fip %s, %v", key, err)
			return err
		}
		return nil
	}

	// redo
	if !cachedFip.Status.Ready &&
		cachedFip.Status.Redo != "" &&
		cachedFip.Status.V4ip != "" &&
		cachedFip.DeletionTimestamp.IsZero() {
		klog.V(3).Infof("reapply fip '%s' in pod", key)
		gwPod, err := c.getNatGwPod(cachedFip.Status.NatGwDp, c.natGwNamespaceByName(cachedFip.Status.NatGwDp))
		if err != nil {
			klog.Error(err)
			return err
		}
		// If Pod started before the redo timestamp, it has not restarted since
		// the redo was marked — iptables rules are still intact, skip re-creation.
		fipRedo, _ := time.ParseInLocation("2006-01-02T15:04:05", cachedFip.Status.Redo, time.Local)
		if len(gwPod.Status.ContainerStatuses) == 0 || gwPod.Status.ContainerStatuses[0].State.Running == nil {
			return fmt.Errorf("fip %s: gateway pod container not running, will retry redo", key)
		}
		if gwPod.Status.ContainerStatuses[0].State.Running.StartedAt.Before(&metav1.Time{Time: fipRedo}) {
			klog.V(3).Infof("fip %s: pod started before redo mark, rules intact, skip", key)
			return nil
		}
		if err = c.createFipInPod(cachedFip.Status.NatGwDp, cachedFip.Status.V4ip, cachedFip.Status.InternalIP); err != nil {
			klog.Errorf("failed to create fip, %v", err)
			return err
		}
		if err = c.patchFipStatus(key, "", "", "", "", true); err != nil {
			klog.Errorf("failed to patch status for fip %s, %v", key, err)
			return err
		}
	}
	if err = c.handleAddIptablesFipFinalizer(key); err != nil {
		klog.Errorf("failed to handle add finalizer for fip %s, %v", key, err)
		return err
	}
	return nil
}

func (c *Controller) handleDelIptablesFip(key string) error {
	klog.V(3).Infof("deleted iptables fip %s", key)
	return nil
}

// handleAddIptablesDnatRule creates a DNAT rule from scratch.
//
// Responsibility:
//   - Bring the DNAT from an empty Status to a fully consistent state across all 4 dimensions:
//     1. iptables rule in NAT GW Pod  2. DNAT Status  3. DNAT Labels  4. EIP Status
//   - On success, Status MUST be fully populated (V4ip, Protocol, ExternalPort, InternalIP,
//     InternalPort, NatGwDp, Ready=true). This is the contract that handleUpdateIptablesDnatRule
//     relies on: a complete Status provides reliable old values for spec-change detection
//     and rule deletion.
//   - On failure at any step, returns error to retry. Partial state may exist; the next
//     retry uses current Spec and must converge to the desired state.
//
// Error state: if this handler never completes (Status.V4ip stays empty), the resource is
// in an error state. The update handler MUST NOT attempt NAT operations in this state.
func (c *Controller) handleAddIptablesDnatRule(key string) error {
	dnat, err := c.iptablesDnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	// The key may have been queued while the object was still live; the update queue owns cleanup.
	if !dnat.DeletionTimestamp.IsZero() {
		return nil
	}

	if dnat.Spec.Type == kubeovnv1.DnatRuleTypeShare {
		return nil
	}

	if vpcNatEnabled != "true" {
		return errors.New("iptables nat gw not enable")
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle add iptables dnat rule %s", key)

	if dnat.Status.V4ip != "" && dnat.Status.NatGwDp != "" && dnat.Status.Protocol != "" &&
		dnat.Status.ExternalPort != "" && dnat.Status.InternalIP != "" && dnat.Status.InternalPort != "" {
		stale := false
		if dnatUsesEip(&dnat.Spec) {
			eip, getErr := c.getBindableEip(dnat.Spec.EIP)
			stale = getErr != nil || dnat.Status.V4ip != eip.Status.IP || dnat.Status.NatGwDp != eip.Spec.NatGwDp ||
				dnat.Status.Protocol != dnat.Spec.Protocol || dnat.Status.ExternalPort != dnat.Spec.ExternalPort ||
				dnat.Status.InternalIP != dnat.Spec.InternalIP || dnat.Status.InternalPort != dnat.Spec.InternalPort ||
				dnat.Labels[util.EipUIDLabel] != string(eip.UID)
		} else {
			stale = dnat.Status.Protocol != dnat.Spec.Protocol || dnat.Status.ExternalPort != dnat.Spec.ExternalPort ||
				dnat.Status.InternalIP != dnat.Spec.InternalIP || dnat.Status.InternalPort != dnat.Spec.InternalPort ||
				dnat.Status.V4ip != dnat.Spec.ClusterIP
		}
		if stale {
			c.updateIptablesDnatRuleQueue.Add(key)
			return nil
		}
		if dnat.Status.Ready {
			return nil
		}
	}
	klog.V(3).Infof("handle add iptables dnat %s", key)

	if err := c.validateDnatRule(dnat); err != nil {
		return err
	}

	gwName, v4ip, v6ip, err := c.resolveDnatAddress(dnat)
	if err != nil {
		klog.Error(err)
		return err
	}
	if dup, duplicateErr := c.isDnatDuplicated(gwName, dnat); dup || duplicateErr != nil {
		if dup {
			c.recorder.Event(dnat, corev1.EventTypeWarning, "DnatIdentityConflict", duplicateErr.Error())
		}
		klog.Error(duplicateErr)
		return duplicateErr
	}
	// Add the finalizer **before** creating rules in Pod. If we added it after,
	// the DNAT could be deleted after createDnatInPod but before the finalizer,
	// leaving unmanaged iptables rules in the gateway pod.
	if err = c.handleAddIptablesDnatFinalizer(key); err != nil {
		klog.Errorf("failed to handle add finalizer for dnat, %v", err)
		return err
	}

	// Claim the EIP before touching the gateway pod: the EIP in-use check counts this label, so a
	// claim written only after the rules exist can be missed by a concurrent EIP release.
	if err = c.patchDnatLabel(key, dnat); err != nil {
		klog.Errorf("failed to patch label for dnat %s, %v", key, err)
		return err
	}

	// Only exclusive DNAT reaches this worker. Share objects are Service accounting records.
	if err = c.createDnatInPod(gwName, dnat.Spec.Protocol,
		v4ip, dnat.Spec.InternalIP,
		dnat.Spec.ExternalPort, dnat.Spec.InternalPort); err != nil {
		klog.Errorf("failed to create dnat, %v", err)
		return err
	}
	if err = c.patchDnatStatus(key, v4ip, v6ip, gwName, "", true); err != nil {
		klog.Errorf("failed to patch status for dnat %s, %v", key, err)
		return err
	}
	if dnat.Spec.EIP == "" {
		return nil
	}
	if err = c.patchEipStatus(dnat.Spec.EIP, "", "", "", true); err != nil {
		// refresh eip nats
		klog.Errorf("failed to patch dnat use eip %s, %v", key, err)
		return err
	}
	// patchDnatLabel updated the DNAT's EipUIDLabel via the API, but the informer cache
	// may not have synced yet when patchEipStatus called getIptablesEipNat above, causing
	// it to miss the DNAT and leave EIP.Status.Nat stale. Schedule a delayed reset.
	c.resetIptablesEipQueue.AddAfter(dnat.Spec.EIP, 3*time.Second)
	return nil
}

// dnatNeedsSpecCleanup reports whether the rule is in the state a crashed spec change leaves
// behind: its status still points at the identity the data plane was programmed with, while the
// spec already describes another one. Both identities then have to be cleaned up, which is what the
// caller does below. Every EIP rule can reach this state (a hand-managed one whose EIP or port was
// edited, or a Service-driven one whose EIP changed).
func dnatNeedsSpecCleanup(dnat *kubeovnv1.IptablesDnatRule) bool {
	// The divergent identity below is resolved through the EIP, so a rule without one has nothing to
	// look up (a ClusterIP rule records no second identity in its status): guarding here keeps the
	// caller from calling GetEip("") and logging a misleading "eip not found" on every deletion.
	return dnatUsesEip(&dnat.Spec) && !dnat.Status.Ready && dnat.Status.V4ip != ""
}

// resolveDnatAddress returns what the data plane needs from the rule's address: the gateway that
// serves it, the IPv4 address it programs and the IPv6 address to record in its status (only an EIP
// can carry one). A rule that serves a ClusterIP has no EIP, so its address is its own.
func (c *Controller) resolveDnatAddress(dnat *kubeovnv1.IptablesDnatRule) (gwName, v4ip, v6ip string, err error) {
	if !dnatUsesEip(&dnat.Spec) && dnatServesClusterIP(&dnat.Spec) {
		return dnat.Labels[util.VpcNatGatewayNameLabel], dnat.Spec.ClusterIP, "", nil
	}
	eip, err := c.getBindableEip(dnat.Spec.EIP)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get eip %s: %w", dnat.Spec.EIP, err)
	}
	return eip.Spec.NatGwDp, eip.Status.IP, eip.Spec.V6ip, nil
}

// handleUpdateIptablesDnatRule handles DNAT rule deletion, spec changes, and redo.
//
// DNAT involves 4 dimensions of data consistency (see handleUpdateIptablesFip for general pattern):
//  1. iptables rule  - (protocol, v4ip, externalPort) -> (internalIP, internalPort)
//  2. DNAT Status    - V4ip, Protocol, InternalIP, ExternalPort, InternalPort, NatGwDp, Ready
//  3. DNAT Label     - EipV4IpLabel, VpcNatGatewayNameLabel, VpcDnatEPortLabel, VpcEipAnnotation
//  4. EIP Status     - Nat field
//
// Key difference from FIP: DNAT identity = (eip, externalPort, protocol).
// del_dnat in nat-gateway.sh matches by identity only (lenient deletion),
// so even if controller passes stale internalIP/internalPort, deletion still works correctly.
//
// Precondition for spec-change path:
//   - Status.V4ip != "" (Status is complete, populated by handleAddIptablesDnatRule).
//     Status provides reliable old values for detecting what changed and deleting old rules.
//   - If Status.V4ip == "", the resource is in an error state (handleAdd never completed).
//     The update handler MUST NOT attempt any NAT operations — it should log the error
//     and return, leaving recovery to the add handler.
//     IMPORTANT: If a change fails, we must wait for retry until success. Continuing to change spec
//     on top of a failed state (dirty status) can introduce more zombie rules that are hard to track.
func (c *Controller) handleUpdateIptablesDnatRule(key string) error {
	cachedDnat, err := c.iptablesDnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}

	if cachedDnat.Spec.Type == kubeovnv1.DnatRuleTypeShare {
		return nil
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle update iptables dnat %s", key)

	// should delete
	if !cachedDnat.DeletionTimestamp.IsZero() {
		if vpcNatEnabled == "true" {
			if err = c.finalDeleteDnatInPod(key, cachedDnat); err != nil {
				return err
			}
		}
		if err = c.handleDelIptablesDnatFinalizer(key); err != nil {
			klog.Errorf("failed to handle del finalizer for dnat %s, %v", key, err)
			return err
		}
		//  reset eip
		if cachedDnat.Spec.EIP != "" {
			c.resetIptablesEipQueue.AddAfter(cachedDnat.Spec.EIP, 3*time.Second)
		}
		return nil
	}
	klog.V(3).Infof("handle update dnat %s", key)

	// add or update should make sure vpc nat enabled. Checked before validateDnatRule and the
	// isDnatDuplicated lister scan (matching handleAddIptablesDnatRule) so we do not scan when
	// the NAT GW is disabled.
	if vpcNatEnabled != "true" {
		if released, releaseErr := c.releaseDeletedEipRef(
			key, util.DnatUsingEip, cachedDnat.Spec.EIP, cachedDnat.Labels, cachedDnat.Annotations,
			func() error { return c.finalDeleteDnatInPod(key, cachedDnat) },
			func() error { return c.patchDnatStatus(key, "", "", "", "", false) },
		); released || releaseErr != nil {
			return releaseErr
		}
		return errors.New("iptables nat gw not enable")
	}

	if err := c.validateDnatRule(cachedDnat); err != nil {
		return err
	}

	gwName := cachedDnat.Labels[util.VpcNatGatewayNameLabel]
	var eip *kubeovnv1.IptablesEIP
	if dnatUsesEip(&cachedDnat.Spec) {
		eip, err = c.getBindableEip(cachedDnat.Spec.EIP)
		if err != nil {
			klog.Errorf("failed to get eip, %v", err)
			if released, releaseErr := c.releaseDeletedEipRef(
				key, util.DnatUsingEip, cachedDnat.Spec.EIP, cachedDnat.Labels, cachedDnat.Annotations,
				func() error { return c.finalDeleteDnatInPod(key, cachedDnat) },
				func() error { return c.patchDnatStatus(key, "", "", "", "", false) },
			); released || releaseErr != nil {
				return releaseErr
			}
			if cachedDnat.Status.Ready {
				if patchErr := c.patchDnatStatus(key, "", "", "", "", false); patchErr != nil {
					return fmt.Errorf("failed to mark dnat %s not ready after its eip became unavailable: %w", key, patchErr)
				}
			}
			return err
		}
		gwName = eip.Spec.NatGwDp
	}
	if dup, duplicateErr := c.isDnatDuplicated(gwName, cachedDnat); dup || duplicateErr != nil {
		if dup && cachedDnat.Status.Ready {
			if err = c.patchDnatStatus(key, "", "", "", "", false); err != nil {
				return err
			}
		}
		if dup {
			c.recorder.Event(cachedDnat, corev1.EventTypeWarning, "DnatIdentityConflict", duplicateErr.Error())
		}
		klog.Errorf("failed to update dnat, %v", duplicateErr)
		return duplicateErr
	}

	if !dnatUsesEip(&cachedDnat.Spec) {
		// A rule that serves a ClusterIP is programmed by the Service reconcile only.
		return nil
	}
	if eip.Spec.NatGwDp == "" {
		klog.Errorf("dnat %s: eip %s has empty NatGwDp, skip binding", key, cachedDnat.Spec.EIP)
		return nil
	}

	// Error state: Status incomplete means handleAdd never completed.
	// Do not attempt NAT operations — no reliable old values for safe rule replacement.
	// All spec-corresponding Status fields must be populated for the update handler to proceed.
	if cachedDnat.Status.V4ip == "" || cachedDnat.Status.NatGwDp == "" || cachedDnat.Status.Protocol == "" ||
		cachedDnat.Status.ExternalPort == "" || cachedDnat.Status.InternalIP == "" || cachedDnat.Status.InternalPort == "" {
		klog.Errorf("dnat %s has incomplete status (V4ip=%q, NatGwDp=%q, Protocol=%q, ExternalPort=%q, InternalIP=%q, InternalPort=%q), skipping NAT operations; waiting for add handler to complete",
			key, cachedDnat.Status.V4ip, cachedDnat.Status.NatGwDp, cachedDnat.Status.Protocol,
			cachedDnat.Status.ExternalPort, cachedDnat.Status.InternalIP, cachedDnat.Status.InternalPort)
		return nil
	}

	// spec change: compare Status (old, what's in Pod) vs Spec+EIP (new, desired)
	oldV4ip := cachedDnat.Status.V4ip
	oldProtocol := cachedDnat.Status.Protocol
	oldExternalPort := cachedDnat.Status.ExternalPort
	newV4ip := eip.Status.IP
	newProtocol := cachedDnat.Spec.Protocol
	newInternalIP := cachedDnat.Spec.InternalIP
	newExternalPort := cachedDnat.Spec.ExternalPort
	newInternalPort := cachedDnat.Spec.InternalPort

	// Warn if we are modifying a resource that might be in a dirty state from a previous failed update.
	if !cachedDnat.Status.Ready {
		// TODO: consider using a webhook to reject spec changes when the resource is not ready,
		// to prevent users from modifying the spec when the previous update has not yet been fully reconciled.
		klog.Warningf("dnat %s is being updated while not Ready (previous update likely failed). This may lead to stale iptables rules.", key)
	}

	// Verify new parameters are valid before modifying any state.
	// eip.Status.IP can be empty if EIP itself is in error state.
	if newV4ip == "" || newInternalIP == "" || newExternalPort == "" || newProtocol == "" || newInternalPort == "" {
		klog.Errorf("skipping dnat %s update: incomplete new parameters (v4ip=%q, internalIP=%q, externalPort=%q, protocol=%q, internalPort=%q)",
			key, newV4ip, newInternalIP, newExternalPort, newProtocol, newInternalPort)
		return nil
	}

	if oldV4ip != newV4ip || cachedDnat.Status.NatGwDp != eip.Spec.NatGwDp ||
		oldProtocol != newProtocol || oldExternalPort != newExternalPort || cachedDnat.Status.InternalIP != newInternalIP ||
		cachedDnat.Status.InternalPort != newInternalPort || cachedDnat.Labels[util.EipUIDLabel] != string(eip.UID) {
		// Mark DNAT as not ready before starting the update.
		// This ensures that if the controller crashes or the update fails midway,
		// the resource will be left in a non-ready state, indicating a potential inconsistency.
		if cachedDnat.Status.Ready {
			klog.V(3).Infof("dnat %s spec changed, marking as not ready before update", key)
			if err = c.patchDnatStatus(key, oldV4ip, cachedDnat.Status.V6ip, cachedDnat.Status.NatGwDp, "", false); err != nil {
				klog.Errorf("failed to mark dnat %s as not ready, %v", key, err)
				return err
			}
		}
		// delete old rule; finalDeleteDnatInPod resolves identity from Status
		if err = c.finalDeleteDnatInPod(key, cachedDnat); err != nil {
			return err
		}
		// Swap the claim between the two pod operations: the old EIP stays claimed until its rule is
		// gone, and the new one is claimed before its rule exists. The in-use check counts this label,
		// so either edge would let a concurrent release drop a finalizer with a live rule behind it.
		if err = c.patchDnatLabel(key, cachedDnat); err != nil {
			klog.Errorf("failed to patch label for dnat %s, %v", key, err)
			return err
		}
		c.enqueueDeletingOldIptablesEip(cachedDnat.Annotations[util.VpcEipAnnotation], cachedDnat.Spec.EIP)

		// Only exclusive DNAT reaches this worker.
		if err = c.createDnatInPod(eip.Spec.NatGwDp, newProtocol,
			newV4ip, newInternalIP,
			newExternalPort, newInternalPort); err != nil {
			klog.Errorf("failed to create dnat %s, %v", key, err)
			return err
		}
		if err = c.patchDnatStatus(key, newV4ip, eip.Spec.V6ip, eip.Spec.NatGwDp, "", true); err != nil {
			klog.Errorf("failed to patch status for dnat %s, %v", key, err)
			return err
		}
		if err = c.patchEipStatus(cachedDnat.Spec.EIP, "", "", "", true); err != nil {
			klog.Errorf("failed to patch dnat use eip %s, %v", key, err)
			return err
		}
		// Same refresh the add path schedules: patchEipStatus read the new EIP's nat off the
		// informer, which may not carry the label patched moments ago.
		c.resetIptablesEipQueue.AddAfter(cachedDnat.Spec.EIP, 3*time.Second)
		// Reset old EIP after all 4 dimensions are updated.
		// This must NOT be in enqueue — async reset there could race with handler's
		// patchEipStatus, causing the old EIP's nat label to be cleared before the
		// new EIP is fully bound, leaving a window where the old EIP appears free.
		if oldEipName := cachedDnat.Annotations[util.VpcEipAnnotation]; oldEipName != "" && oldEipName != cachedDnat.Spec.EIP {
			c.resetIptablesEipQueue.AddAfter(oldEipName, 3*time.Second)
		}
		if err = c.handleAddIptablesDnatFinalizer(key); err != nil {
			klog.Errorf("failed to handle add finalizer for dnat %s, %v", key, err)
			return err
		}
		return nil
	}

	// redo
	if !cachedDnat.Status.Ready &&
		cachedDnat.Status.Redo != "" &&
		cachedDnat.Status.V4ip != "" &&
		cachedDnat.DeletionTimestamp.IsZero() {
		klog.V(3).Infof("reapply dnat in pod for %s", key)
		gwPod, err := c.getNatGwPod(cachedDnat.Status.NatGwDp, c.natGwNamespaceByName(cachedDnat.Status.NatGwDp))
		if err != nil {
			klog.Error(err)
			return err
		}
		// If Pod started before the redo timestamp, it has not restarted since
		// the redo was marked — iptables rules are still intact, skip re-creation.
		dnatRedo, _ := time.ParseInLocation("2006-01-02T15:04:05", cachedDnat.Status.Redo, time.Local)
		if len(gwPod.Status.ContainerStatuses) == 0 || gwPod.Status.ContainerStatuses[0].State.Running == nil {
			return fmt.Errorf("dnat %s: gateway pod container not running, will retry redo", key)
		}
		if gwPod.Status.ContainerStatuses[0].State.Running.StartedAt.Before(&metav1.Time{Time: dnatRedo}) {
			klog.V(3).Infof("dnat %s: pod started before redo mark, rules intact, skip", key)
			return nil
		}

		if err = c.createDnatInPod(cachedDnat.Status.NatGwDp, cachedDnat.Status.Protocol,
			cachedDnat.Status.V4ip, cachedDnat.Status.InternalIP,
			cachedDnat.Status.ExternalPort, cachedDnat.Status.InternalPort); err != nil {
			klog.Errorf("failed to create dnat %s, %v", key, err)
			return err
		}
		if err = c.patchDnatStatus(key, "", "", "", "", true); err != nil {
			klog.Errorf("failed to patch status for dnat %s, %v", key, err)
			return err
		}
	}
	if err = c.handleAddIptablesDnatFinalizer(key); err != nil {
		klog.Errorf("failed to handle add finalizer for dnat %s, %v", key, err)
		return err
	}
	return nil
}

func (c *Controller) handleDelIptablesDnatRule(key string) error {
	klog.V(3).Infof("deleted iptables dnat %s", key)
	return nil
}

// handleAddIptablesSnatRule creates a SNAT rule from scratch.
//
// Responsibility:
//   - Bring the SNAT from an empty Status to a fully consistent state across all 4 dimensions:
//     1. iptables rule in NAT GW Pod  2. SNAT Status  3. SNAT Labels  4. EIP Status
//   - On success, Status MUST be fully populated (V4ip, InternalCIDR, NatGwDp, Ready=true).
//     This is the contract that handleUpdateIptablesSnatRule relies on: a complete Status
//     provides reliable old values for spec-change detection and rule deletion.
//   - On failure at any step, returns error to retry. Partial state may exist; the next
//     retry uses current Spec and must converge to the desired state.
//
// Error state: if this handler never completes (Status.V4ip stays empty), the resource is
// in an error state. The update handler MUST NOT attempt NAT operations in this state.
func (c *Controller) handleAddIptablesSnatRule(key string) error {
	snat, err := c.iptablesSnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	// The key may have been queued while the object was still live; the update queue owns cleanup.
	if !snat.DeletionTimestamp.IsZero() {
		return nil
	}

	if vpcNatEnabled != "true" {
		return errors.New("iptables nat gw not enable")
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle add iptables snat rule %s", key)

	if snat.Status.V4ip != "" && snat.Status.NatGwDp != "" && snat.Status.InternalCIDR != "" {
		eip, getErr := c.getBindableEip(snat.Spec.EIP)
		statusV4Cidr, _ := util.SplitStringIP(snat.Status.InternalCIDR)
		specV4Cidr, _ := util.SplitStringIP(snat.Spec.InternalCIDR)
		if getErr != nil || snat.Status.V4ip != eip.Status.IP || snat.Status.NatGwDp != eip.Spec.NatGwDp ||
			statusV4Cidr != specV4Cidr || snat.Labels[util.EipUIDLabel] != string(eip.UID) {
			c.updateIptablesSnatRuleQueue.Add(key)
			return nil
		}
		if snat.Status.Ready {
			return nil
		}
	}
	klog.V(3).Infof("handle add iptables snat %s", key)

	if err := c.validateSnatRule(snat); err != nil {
		return err
	}

	eip, err := c.getBindableEip(snat.Spec.EIP)
	if err != nil {
		klog.Errorf("failed to get eip, %v", err)
		return err
	}
	// create snat
	v4Cidr, _ := util.SplitStringIP(snat.Spec.InternalCIDR)
	if v4Cidr == "" {
		// only support IPv4 snat
		err = fmt.Errorf("failed to get snat v4 internal cidr, original cidr is %s", snat.Spec.InternalCIDR)
		return err
	}
	// Add the finalizer **before** creating rules in Pod. If we added it after,
	// the SNAT could be deleted after createSnatInPod but before the finalizer,
	// leaving unmanaged iptables rules in the gateway pod.
	if err = c.handleAddIptablesSnatFinalizer(key); err != nil {
		klog.Errorf("failed to handle add finalizer for snat, %v", err)
		return err
	}

	// Claim the EIP before touching the gateway pod: the EIP in-use check counts this label, so a
	// claim written only after the rules exist can be missed by a concurrent EIP release.
	if err = c.patchSnatLabel(key, eip); err != nil {
		klog.Errorf("failed to patch label for snat %s, %v", key, err)
		return err
	}

	memberID := resolveSnatMemberID(eip, snat)
	if err = c.createSnatInPodWithMember(eip.Spec.NatGwDp, eip.Status.IP, v4Cidr, memberID); err != nil {
		klog.Errorf("failed to create snat, %v", err)
		return err
	}
	if err = c.patchSnatStatus(key, eip.Status.IP, eip.Spec.V6ip, eip.Spec.NatGwDp, "", true); err != nil {
		klog.Errorf("failed to update status for snat %s, %v", key, err)
		return err
	}
	if err = c.patchEipStatus(snat.Spec.EIP, "", "", "", true); err != nil {
		// refresh eip nats
		klog.Errorf("failed to patch snat use eip %s, %v", key, err)
		return err
	}
	// patchSnatLabel updated the SNAT's EipUIDLabel via the API, but the informer cache
	// may not have synced yet when patchEipStatus called getIptablesEipNat above, causing
	// it to silently miss the SNAT and leave EIP.Status.Nat stale. Schedule a delayed reset
	// so that after the informer syncs the label, the EIP nat status is corrected.
	c.resetIptablesEipQueue.AddAfter(snat.Spec.EIP, 3*time.Second)
	return nil
}

// handleUpdateIptablesSnatRule handles SNAT rule deletion, spec changes, and redo.
//
// SNAT involves 4 dimensions of data consistency (see handleUpdateIptablesFip for general pattern):
//  1. iptables rule  - v4ip SNAT for internalCIDR
//  2. SNAT Status    - V4ip, InternalCIDR, NatGwDp, Ready
//  3. SNAT Label     - EipV4IpLabel, VpcNatGatewayNameLabel, VpcEipAnnotation
//  4. EIP Status     - Nat field
//
// Key difference from FIP/DNAT: SNAT is a 1:N model. One EIP can serve multiple CIDRs,
// and one CIDR can have multiple EIPs (for port exhaustion mitigation via --random-fully).
// SNAT identity = (eip, internalCIDR). del_snat in nat-gateway.sh matches by identity
// to locate the exact rule.
//
// Precondition for spec-change path:
//   - Status.V4ip != "" (Status is complete, populated by handleAddIptablesSnatRule).
//     Status provides reliable old values for detecting what changed and deleting old rules.
//   - If Status.V4ip == "", the resource is in an error state (handleAdd never completed).
//     The update handler MUST NOT attempt any NAT operations — it should log the error
//     and return, leaving recovery to the add handler.
//     IMPORTANT: If a change fails, we must wait for retry until success. Continuing to change spec
//     on top of a failed state (dirty status) can introduce more zombie rules that are hard to track.
func (c *Controller) handleUpdateIptablesSnatRule(key string) error {
	cachedSnat, err := c.iptablesSnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}

	c.vpcNatGwKeyMutex.LockKey(key)
	defer func() { _ = c.vpcNatGwKeyMutex.UnlockKey(key) }()
	klog.Infof("handle update iptables snat rule %s", key)

	// should delete
	if !cachedSnat.DeletionTimestamp.IsZero() {
		if vpcNatEnabled == "true" {
			if err = c.finalDeleteSnatInPod(key, cachedSnat); err != nil {
				return err
			}
		}
		if err = c.handleDelIptablesSnatFinalizer(key); err != nil {
			klog.Errorf("failed to handle del finalizer for snat %s, %v", key, err)
			return err
		}
		c.resetIptablesEipQueue.AddAfter(cachedSnat.Spec.EIP, 3*time.Second)
		return nil
	}
	klog.V(3).Infof("handle update snat %s", key)

	if err := c.validateSnatRule(cachedSnat); err != nil {
		return err
	}

	eip, err := c.getBindableEip(cachedSnat.Spec.EIP)
	if err != nil {
		klog.Errorf("failed to get eip, %v", err)
		if released, releaseErr := c.releaseDeletedEipRef(
			key, util.SnatUsingEip, cachedSnat.Spec.EIP, cachedSnat.Labels, cachedSnat.Annotations,
			func() error { return c.finalDeleteSnatInPod(key, cachedSnat) },
			func() error { return c.patchSnatStatus(key, "", "", "", "", false) },
		); released || releaseErr != nil {
			return releaseErr
		}
		if cachedSnat.Status.Ready {
			if patchErr := c.patchSnatStatus(key, "", "", "", "", false); patchErr != nil {
				return fmt.Errorf("failed to mark snat %s not ready after its eip became unavailable: %w", key, patchErr)
			}
		}
		return err
	}

	// add or update should make sure vpc nat enabled
	if vpcNatEnabled != "true" {
		return errors.New("iptables nat gw not enable")
	}

	if eip.Spec.NatGwDp == "" {
		klog.Errorf("snat %s: eip %s has empty NatGwDp, skip binding", key, cachedSnat.Spec.EIP)
		return nil
	}

	// Error state: Status incomplete means handleAdd never completed.
	// Do not attempt NAT operations — no reliable old values for safe rule replacement.
	// All spec-corresponding Status fields must be populated for the update handler to proceed.
	if cachedSnat.Status.V4ip == "" || cachedSnat.Status.NatGwDp == "" || cachedSnat.Status.InternalCIDR == "" {
		klog.Errorf("snat %s has incomplete status (V4ip=%q, NatGwDp=%q, InternalCIDR=%q), skipping NAT operations; waiting for add handler to complete",
			key, cachedSnat.Status.V4ip, cachedSnat.Status.NatGwDp, cachedSnat.Status.InternalCIDR)
		return nil
	}

	// spec change: compare Status (old, what's in Pod) vs Spec+EIP (new, desired)
	// SNAT identity = (v4ip, internalCIDR), all fields are identity — no non-identity fields.
	oldV4ip := cachedSnat.Status.V4ip
	oldV4Cidr, _ := util.SplitStringIP(cachedSnat.Status.InternalCIDR)
	newV4ip := eip.Status.IP
	newV4Cidr, _ := util.SplitStringIP(cachedSnat.Spec.InternalCIDR)

	// Warn if we are modifying a resource that might be in a dirty state from a previous failed update.
	if !cachedSnat.Status.Ready {
		// TODO: consider using a webhook to reject spec changes when the resource is not ready,
		// to prevent users from modifying the spec when the previous update has not yet been fully reconciled.
		klog.Warningf("snat %s is being updated while not Ready (previous update likely failed). This may lead to stale iptables rules.", key)
	}

	// Verify new parameters are valid before modifying any state.
	// eip.Status.IP can be empty if EIP itself is in error state;
	// SplitStringIP can return empty v4 part for IPv6-only CIDR.
	if newV4ip == "" || newV4Cidr == "" {
		klog.Errorf("skipping snat %s update: incomplete new parameters (v4ip=%q, v4Cidr=%q)", key, newV4ip, newV4Cidr)
		return nil
	}

	desiredMember := resolveSnatMemberID(eip, cachedSnat)
	currentMember := cachedSnat.Labels[util.NatGatewayMemberLabel]

	if oldV4ip != newV4ip || cachedSnat.Status.NatGwDp != eip.Spec.NatGwDp || oldV4Cidr != newV4Cidr ||
		cachedSnat.Labels[util.EipUIDLabel] != string(eip.UID) || currentMember != desiredMember {
		// Mark SNAT as not ready before starting the update.
		// This ensures that if the controller crashes or the update fails midway,
		// the resource will be left in a non-ready state, indicating a potential inconsistency.
		if cachedSnat.Status.Ready {
			klog.V(3).Infof("snat %s spec changed, marking as not ready before update", key)
			if err = c.patchSnatStatus(key, oldV4ip, cachedSnat.Status.V6ip, cachedSnat.Status.NatGwDp, "", false); err != nil {
				klog.Errorf("failed to mark snat %s as not ready, %v", key, err)
				return err
			}
		}
		// delete old rule; finalDeleteSnatInPod resolves identity from Status
		if err = c.finalDeleteSnatInPod(key, cachedSnat); err != nil {
			return err
		}
		// Swap the claim between the two pod operations: the old EIP stays claimed until its rule is
		// gone, and the new one is claimed before its rule exists. The in-use check counts this label,
		// so either edge would let a concurrent release drop a finalizer with a live rule behind it.
		if err = c.patchSnatLabel(key, eip); err != nil {
			klog.Errorf("failed to patch label for snat %s, %v", key, err)
			return err
		}
		c.enqueueDeletingOldIptablesEip(cachedSnat.Annotations[util.VpcEipAnnotation], cachedSnat.Spec.EIP)
		memberID := resolveSnatMemberID(eip, cachedSnat)
		if err = c.createSnatInPodWithMember(eip.Spec.NatGwDp, newV4ip, newV4Cidr, memberID); err != nil {
			klog.Errorf("failed to create snat %s, %v", key, err)
			return err
		}
		if err = c.patchSnatStatus(key, newV4ip, eip.Spec.V6ip, eip.Spec.NatGwDp, "", true); err != nil {
			klog.Errorf("failed to patch status for snat %s, %v", key, err)
			return err
		}
		if err = c.patchEipStatus(cachedSnat.Spec.EIP, "", "", "", true); err != nil {
			klog.Errorf("failed to patch snat use eip %s, %v", key, err)
			return err
		}
		// Same refresh the add path schedules: patchEipStatus read the new EIP's nat off the
		// informer, which may not carry the label patched moments ago.
		c.resetIptablesEipQueue.AddAfter(cachedSnat.Spec.EIP, 3*time.Second)
		// Reset old EIP after all 4 dimensions are updated.
		// This must NOT be in enqueue — async reset there could race with handler's
		// patchEipStatus, causing the old EIP's nat label to be cleared before the
		// new EIP is fully bound, leaving a window where the old EIP appears free.
		if oldEipName := cachedSnat.Annotations[util.VpcEipAnnotation]; oldEipName != "" && oldEipName != cachedSnat.Spec.EIP {
			c.resetIptablesEipQueue.AddAfter(oldEipName, 3*time.Second)
		}
		if err = c.handleAddIptablesSnatFinalizer(key); err != nil {
			klog.Errorf("failed to handle add finalizer for snat %s, %v", key, err)
			return err
		}
		return nil
	}

	// redo
	if !cachedSnat.Status.Ready &&
		cachedSnat.Status.Redo != "" &&
		cachedSnat.Status.V4ip != "" &&
		cachedSnat.DeletionTimestamp.IsZero() {
		gwPod, err := c.getNatGwPod(cachedSnat.Status.NatGwDp, c.natGwNamespaceByName(cachedSnat.Status.NatGwDp))
		if err != nil {
			klog.Error(err)
			return err
		}
		// If Pod started before the redo timestamp, it has not restarted since
		// the redo was marked — iptables rules are still intact, skip re-creation.
		snatRedo, _ := time.ParseInLocation("2006-01-02T15:04:05", cachedSnat.Status.Redo, time.Local)
		if len(gwPod.Status.ContainerStatuses) == 0 || gwPod.Status.ContainerStatuses[0].State.Running == nil {
			return fmt.Errorf("snat %s: gateway pod container not running, will retry redo", key)
		}
		if gwPod.Status.ContainerStatuses[0].State.Running.StartedAt.Before(&metav1.Time{Time: snatRedo}) {
			klog.V(3).Infof("snat %s: pod started before redo mark, rules intact, skip", key)
			return nil
		}
		var eip *kubeovnv1.IptablesEIP
		if cachedSnat.Spec.EIP != "" {
			eip, _ = c.iptablesEipsLister.Get(cachedSnat.Spec.EIP)
		}
		memberID := resolveRecordedSnatMemberID(cachedSnat, eip)
		if err = c.createSnatInPodWithMember(cachedSnat.Status.NatGwDp, cachedSnat.Status.V4ip, cachedSnat.Status.InternalCIDR, memberID); err != nil {
			klog.Errorf("failed to create new snat, %v", err)
			return err
		}
		if err = c.patchSnatStatus(key, "", "", "", "", true); err != nil {
			klog.Errorf("failed to patch status for snat %s, %v", key, err)
			return err
		}
	}
	if err = c.handleAddIptablesSnatFinalizer(key); err != nil {
		klog.Errorf("failed to handle add finalizer for snat %s, %v", key, err)
		return err
	}
	return nil
}

func (c *Controller) handleDelIptablesSnatRule(key string) error {
	klog.V(3).Infof("deleted iptables snat %s", key)
	return nil
}

func (c *Controller) syncIptablesFipFinalizer(cl client.Client) error {
	// migrate depreciated finalizer to new finalizer
	rules := &kubeovnv1.IptablesFIPRuleList{}
	return migrateFinalizers(cl, rules, func(i int) (client.Object, client.Object) {
		if i < 0 || i >= len(rules.Items) {
			return nil, nil
		}
		return rules.Items[i].DeepCopy(), rules.Items[i].DeepCopy()
	})
}

func (c *Controller) handleAddIptablesFipFinalizer(key string) error {
	cachedIptablesFip, err := c.iptablesFipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if !cachedIptablesFip.DeletionTimestamp.IsZero() || controllerutil.ContainsFinalizer(cachedIptablesFip, util.KubeOVNControllerFinalizer) {
		return nil
	}
	newIptablesFip := cachedIptablesFip.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesFip, util.DepreciatedFinalizerName)
	controllerutil.AddFinalizer(newIptablesFip, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesFip, newIptablesFip)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables fip '%s', %v", cachedIptablesFip.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Patch(context.Background(), cachedIptablesFip.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to add finalizer for iptables fip '%s', %v", cachedIptablesFip.Name, err)
		return err
	}
	return nil
}

func (c *Controller) handleDelIptablesFipFinalizer(key string) error {
	cachedIptablesFip, err := c.iptablesFipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if len(cachedIptablesFip.GetFinalizers()) == 0 {
		return nil
	}
	newIptablesFip := cachedIptablesFip.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesFip, util.DepreciatedFinalizerName)
	controllerutil.RemoveFinalizer(newIptablesFip, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesFip, newIptablesFip)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables fip '%s', %v", cachedIptablesFip.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Patch(context.Background(), cachedIptablesFip.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to remove finalizer from iptables fip '%s', %v", cachedIptablesFip.Name, err)
		return err
	}
	return nil
}

func (c *Controller) syncIptablesDnatFinalizer(cl client.Client) error {
	// migrate depreciated finalizer to new finalizer
	rules := &kubeovnv1.IptablesDnatRuleList{}
	return migrateFinalizers(cl, rules, func(i int) (client.Object, client.Object) {
		if i < 0 || i >= len(rules.Items) {
			return nil, nil
		}
		return rules.Items[i].DeepCopy(), rules.Items[i].DeepCopy()
	})
}

func (c *Controller) handleAddIptablesDnatFinalizer(key string) error {
	cachedIptablesDnat, err := c.iptablesDnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if !cachedIptablesDnat.DeletionTimestamp.IsZero() || controllerutil.ContainsFinalizer(cachedIptablesDnat, util.KubeOVNControllerFinalizer) {
		return nil
	}
	newIptablesDnat := cachedIptablesDnat.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesDnat, util.DepreciatedFinalizerName)
	controllerutil.AddFinalizer(newIptablesDnat, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesDnat, newIptablesDnat)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables dnat '%s', %v", cachedIptablesDnat.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Patch(context.Background(), cachedIptablesDnat.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to add finalizer for iptables dnat '%s', %v", cachedIptablesDnat.Name, err)
		return err
	}
	return nil
}

func (c *Controller) handleDelIptablesDnatFinalizer(key string) error {
	cachedIptablesDnat, err := c.iptablesDnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if len(cachedIptablesDnat.GetFinalizers()) == 0 {
		return nil
	}
	newIptablesDnat := cachedIptablesDnat.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesDnat, util.DepreciatedFinalizerName)
	controllerutil.RemoveFinalizer(newIptablesDnat, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesDnat, newIptablesDnat)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables dnat '%s', %v", cachedIptablesDnat.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Patch(context.Background(), cachedIptablesDnat.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to remove finalizer from iptables dnat '%s', %v", cachedIptablesDnat.Name, err)
		return err
	}
	return nil
}

func (c *Controller) patchFipLabel(key string, eip *kubeovnv1.IptablesEIP) error {
	oriFip, err := c.iptablesFipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	fip := oriFip.DeepCopy()
	var needUpdateLabel, needUpdateAnno bool
	var op string
	if len(fip.Labels) == 0 {
		op = "add"
		fip.Labels = map[string]string{
			util.VpcNatGatewayNameLabel: eip.Spec.NatGwDp,
			util.EipV4IpLabel:           eip.Spec.V4ip,
			util.EipUIDLabel:            string(eip.UID),
		}
		needUpdateLabel = true
	} else if fip.Labels[util.VpcNatGatewayNameLabel] != eip.Spec.NatGwDp ||
		fip.Labels[util.EipV4IpLabel] != eip.Spec.V4ip ||
		fip.Labels[util.EipUIDLabel] != string(eip.UID) {
		op = "replace"
		fip.Labels[util.VpcNatGatewayNameLabel] = eip.Spec.NatGwDp
		fip.Labels[util.EipV4IpLabel] = eip.Spec.V4ip
		fip.Labels[util.EipUIDLabel] = string(eip.UID)
		needUpdateLabel = true
	}
	if needUpdateLabel {
		if err := c.updateIptableLabels(fip.Name, op, util.FipUsingEip, fip.Labels); err != nil {
			klog.Error(err)
			return err
		}
	}

	if len(fip.Annotations) == 0 {
		op = "add"
		needUpdateAnno = true
		fip.Annotations = map[string]string{
			util.VpcEipAnnotation: eip.Name,
		}
	} else if fip.Annotations[util.VpcEipAnnotation] != eip.Name {
		op = "replace"
		needUpdateAnno = true
		fip.Annotations[util.VpcEipAnnotation] = eip.Name
	}
	if needUpdateAnno {
		if err := c.updateIptableAnnotations(fip.Name, op, util.FipUsingEip, fip.Annotations); err != nil {
			klog.Error(err)
			return err
		}
	}
	return nil
}

func (c *Controller) releaseDeletedEipRef(
	key, natType, eipName string,
	currentLabels, currentAnnotations map[string]string,
	deleteInPod func() error,
	markNotReady func() error,
) (bool, error) {
	boundEipName := boundIptablesEipName(eipName, currentAnnotations)
	deleted, err := c.eipDeletedOrDeleting(boundEipName)
	if err != nil || !deleted {
		return false, err
	}
	if vpcNatEnabled == "true" {
		if err := deleteInPod(); err != nil {
			return true, err
		}
	}
	if err := markNotReady(); err != nil {
		return true, err
	}
	if err := c.releaseIptablesEipClaim(key, natType, currentLabels, currentAnnotations); err != nil {
		return true, err
	}
	if c.updateIptablesEipQueue != nil {
		c.updateIptablesEipQueue.Add(boundEipName)
	}
	return true, nil
}

func boundIptablesEipName(specEipName string, annotations map[string]string) string {
	if eipName := annotations[util.VpcEipAnnotation]; eipName != "" {
		return eipName
	}
	return specEipName
}

func (c *Controller) eipDeletedOrDeleting(eipName string) (bool, error) {
	eip, err := c.iptablesEipsLister.Get(eipName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}
	return !eip.DeletionTimestamp.IsZero(), nil
}

func (c *Controller) releaseIptablesEipClaim(key, natType string, currentLabels, currentAnnotations map[string]string) error {
	labels := copyLabels(currentLabels)
	delete(labels, util.VpcNatGatewayNameLabel)
	delete(labels, util.VpcDnatEPortLabel)
	delete(labels, util.EipV4IpLabel)
	delete(labels, util.EipUIDLabel)
	if len(labels) != len(currentLabels) {
		if err := c.updateIptableLabels(key, patchLabelsOp(currentLabels), natType, labels); err != nil {
			return err
		}
	}
	annotations := copyLabels(currentAnnotations)
	delete(annotations, util.VpcEipAnnotation)
	if len(annotations) != len(currentAnnotations) {
		if err := c.updateIptableAnnotations(key, patchLabelsOp(currentAnnotations), natType, annotations); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) enqueueDeletingOldIptablesEip(oldEipName, currentEipName string) {
	if oldEipName == "" || oldEipName == currentEipName || c.updateIptablesEipQueue == nil {
		return
	}
	deleted, err := c.eipDeletedOrDeleting(oldEipName)
	if err != nil {
		klog.Errorf("failed to check old eip %s deletion state: %v", oldEipName, err)
		return
	}
	if deleted {
		c.updateIptablesEipQueue.Add(oldEipName)
	}
}

func (c *Controller) syncIptablesSnatFinalizer(cl client.Client) error {
	rules := &kubeovnv1.IptablesSnatRuleList{}
	return migrateFinalizers(cl, rules, func(i int) (client.Object, client.Object) {
		if i < 0 || i >= len(rules.Items) {
			return nil, nil
		}
		return rules.Items[i].DeepCopy(), rules.Items[i].DeepCopy()
	})
}

func (c *Controller) handleAddIptablesSnatFinalizer(key string) error {
	cachedIptablesSnat, err := c.iptablesSnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if !cachedIptablesSnat.DeletionTimestamp.IsZero() || controllerutil.ContainsFinalizer(cachedIptablesSnat, util.KubeOVNControllerFinalizer) {
		return nil
	}
	newIptablesSnat := cachedIptablesSnat.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesSnat, util.DepreciatedFinalizerName)
	controllerutil.AddFinalizer(newIptablesSnat, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesSnat, newIptablesSnat)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables snat '%s', %v", cachedIptablesSnat.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesSnatRules().Patch(context.Background(), cachedIptablesSnat.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to add finalizer for iptables snat '%s', %v", cachedIptablesSnat.Name, err)
		return err
	}
	return nil
}

func (c *Controller) handleDelIptablesSnatFinalizer(key string) error {
	cachedIptablesSnat, err := c.iptablesSnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if len(cachedIptablesSnat.GetFinalizers()) == 0 {
		return nil
	}
	newIptablesSnat := cachedIptablesSnat.DeepCopy()
	controllerutil.RemoveFinalizer(newIptablesSnat, util.DepreciatedFinalizerName)
	controllerutil.RemoveFinalizer(newIptablesSnat, util.KubeOVNControllerFinalizer)
	patch, err := util.GenerateMergePatchPayload(cachedIptablesSnat, newIptablesSnat)
	if err != nil {
		klog.Errorf("failed to generate patch payload for iptables snat '%s', %v", cachedIptablesSnat.Name, err)
		return err
	}
	if _, err := c.config.KubeOvnClient.KubeovnV1().IptablesSnatRules().Patch(context.Background(), cachedIptablesSnat.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, ""); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to remove finalizer from iptables snat '%s', %v", cachedIptablesSnat.Name, err)
		return err
	}
	return nil
}

func (c *Controller) patchFipStatus(key, v4ip, v6ip, natGwDp, redo string, ready bool) error {
	oriFip, err := c.iptablesFipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	fip := oriFip.DeepCopy()
	var changed bool
	if fip.Status.Ready != ready {
		fip.Status.Ready = ready
		changed = true
	}
	if redo != "" && fip.Status.Redo != redo {
		fip.Status.Redo = redo
		changed = true
	}

	if ready && v4ip != "" && fip.Status.V4ip != v4ip {
		fip.Status.V4ip = v4ip
		fip.Status.V6ip = v6ip
		fip.Status.NatGwDp = natGwDp
		changed = true
	}
	if ready && fip.Spec.InternalIP != "" && fip.Status.InternalIP != fip.Spec.InternalIP {
		fip.Status.InternalIP = fip.Spec.InternalIP
		changed = true
	}

	if changed {
		bytes, err := fip.Status.Bytes()
		if err != nil {
			klog.Error(err)
			return err
		}
		if _, err = c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Patch(context.Background(), fip.Name,
			types.MergePatchType, bytes, metav1.PatchOptions{}, "status"); err != nil {
			if k8serrors.IsNotFound(err) {
				return nil
			}
			klog.Errorf("failed to patch fip %s, %v", fip.Name, err)
			return err
		}
	}
	return nil
}

func (c *Controller) redoFip(key, redo string, eipReady bool) error {
	fip, err := c.iptablesFipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to get fip %s, %v", key, err)
		return err
	}
	if redo != "" && redo != fip.Status.Redo {
		if !eipReady {
			if err = c.patchEipLabel(fip.Spec.EIP); err != nil {
				err = fmt.Errorf("failed to patch eip %s, %w", fip.Spec.EIP, err)
				klog.Error(err)
				return err
			}
			if err = c.patchEipStatus(fip.Spec.EIP, "", redo, "", false); err != nil {
				err = fmt.Errorf("failed to patch eip %s, %w", fip.Spec.EIP, err)
				klog.Error(err)
				return err
			}
		}
		if err = c.patchFipStatus(key, "", "", "", redo, false); err != nil {
			err = fmt.Errorf("failed to patch fip %s, %w", fip.Name, err)
			klog.Error(err)
			return err
		}
	}
	return err
}

// patchDnatLabel records who owns the rule: for an EIP rule the gateway, port, EIP address and EIP
// UID, all derived from the EIP; for a ClusterIP rule the serving gateway comes from its label. The
// labels are what the rule consumers (share backend aggregation, redo, VIP state) select on, so
// they are reconciled on every pass rather than only at creation.
func (c *Controller) patchDnatLabel(key string, rule *kubeovnv1.IptablesDnatRule) error {
	oriDnat, err := c.iptablesDnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	dnat := oriDnat.DeepCopy()
	var needUpdateLabel, needUpdateAnno bool
	var op string
	if dnat.Spec.EIP == "" {
		gateway := rule.Labels[util.VpcNatGatewayNameLabel]
		if dnat.Labels[util.VpcNatGatewayNameLabel] == gateway &&
			dnat.Labels[util.VpcDnatEPortLabel] == rule.Spec.ExternalPort {
			return nil
		}
		if _, ok := dnat.Labels[util.VpcNatGatewayNameLabel]; !ok || len(dnat.Labels) == 0 {
			op = "add"
		} else {
			op = "replace"
		}
		if dnat.Labels == nil {
			dnat.Labels = map[string]string{}
		}
		dnat.Labels[util.VpcNatGatewayNameLabel] = gateway
		dnat.Labels[util.VpcDnatEPortLabel] = rule.Spec.ExternalPort
		if err := c.updateIptableLabels(dnat.Name, op, util.DnatUsingEip, dnat.Labels); err != nil {
			klog.Error(err)
			return err
		}
		return nil
	}
	eip, err := c.iptablesEipsLister.Get(rule.Spec.EIP)
	if err != nil {
		return fmt.Errorf("failed to get eip %s of dnat %s: %w", rule.Spec.EIP, key, err)
	}
	if len(dnat.Labels) == 0 {
		op = "add"
		dnat.Labels = map[string]string{
			util.VpcNatGatewayNameLabel: eip.Spec.NatGwDp,
			util.VpcDnatEPortLabel:      dnat.Spec.ExternalPort,
			util.EipV4IpLabel:           eip.Spec.V4ip,
			util.EipUIDLabel:            string(eip.UID),
		}
		needUpdateLabel = true
	} else if dnat.Labels[util.VpcNatGatewayNameLabel] != eip.Spec.NatGwDp ||
		dnat.Labels[util.VpcDnatEPortLabel] != dnat.Spec.ExternalPort ||
		dnat.Labels[util.EipV4IpLabel] != eip.Spec.V4ip ||
		dnat.Labels[util.EipUIDLabel] != string(eip.UID) {
		op = "replace"
		dnat.Labels[util.VpcNatGatewayNameLabel] = eip.Spec.NatGwDp
		dnat.Labels[util.VpcDnatEPortLabel] = dnat.Spec.ExternalPort
		dnat.Labels[util.EipV4IpLabel] = eip.Spec.V4ip
		dnat.Labels[util.EipUIDLabel] = string(eip.UID)
		needUpdateLabel = true
	}
	if needUpdateLabel {
		if err := c.updateIptableLabels(dnat.Name, op, util.DnatUsingEip, dnat.Labels); err != nil {
			klog.Error(err)
			return err
		}
	}

	if len(dnat.Annotations) == 0 {
		op = "add"
		needUpdateAnno = true
		dnat.Annotations = map[string]string{
			util.VpcEipAnnotation: eip.Name,
		}
	} else if dnat.Annotations[util.VpcEipAnnotation] != eip.Name {
		op = "replace"
		needUpdateAnno = true
		dnat.Annotations[util.VpcEipAnnotation] = eip.Name
	}
	if needUpdateAnno {
		if err := c.updateIptableAnnotations(dnat.Name, op, util.DnatUsingEip, dnat.Annotations); err != nil {
			klog.Error(err)
			return err
		}
	}
	return nil
}

func (c *Controller) patchDnatStatus(key, v4ip, v6ip, natGwDp, redo string, ready bool) error {
	oriDnat, err := c.iptablesDnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	var changed bool
	dnat := oriDnat.DeepCopy()
	if dnat.Status.Ready != ready {
		dnat.Status.Ready = ready
		changed = true
	}
	if redo != "" && dnat.Status.Redo != redo {
		dnat.Status.Redo = redo
		changed = true
	}
	if ready && v4ip != "" && dnat.Status.V4ip != v4ip {
		dnat.Status.V4ip = v4ip
		dnat.Status.V6ip = v6ip
		dnat.Status.NatGwDp = natGwDp
		changed = true
	}
	if ready && dnat.Spec.Protocol != "" && dnat.Status.Protocol != dnat.Spec.Protocol {
		dnat.Status.Protocol = dnat.Spec.Protocol
		changed = true
	}
	if ready && dnat.Spec.InternalIP != "" && dnat.Status.InternalIP != dnat.Spec.InternalIP {
		dnat.Status.InternalIP = dnat.Spec.InternalIP
		changed = true
	}
	if ready && dnat.Spec.InternalPort != "" && dnat.Status.InternalPort != dnat.Spec.InternalPort {
		dnat.Status.InternalPort = dnat.Spec.InternalPort
		changed = true
	}
	if ready && dnat.Spec.ExternalPort != "" && dnat.Status.ExternalPort != dnat.Spec.ExternalPort {
		dnat.Status.ExternalPort = dnat.Spec.ExternalPort
		changed = true
	}

	if changed {
		bytes, err := dnat.Status.Bytes()
		if err != nil {
			klog.Error(err)
			return err
		}
		if _, err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Patch(context.Background(), dnat.Name,
			types.MergePatchType, bytes, metav1.PatchOptions{}, "status"); err != nil {
			if k8serrors.IsNotFound(err) {
				return nil
			}
			klog.Errorf("failed to patch dnat %s, %v", dnat.Name, err)
			return err
		}
	}
	return nil
}

func (c *Controller) redoDnat(key, redo string, eipReady bool) error {
	dnat, err := c.iptablesDnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to get dnat %s, %v", key, err)
		return err
	}
	if redo != "" && redo != dnat.Status.Redo {
		if !eipReady && dnat.Spec.EIP != "" {
			if err = c.patchEipStatus(dnat.Spec.EIP, "", redo, "", false); err != nil {
				err = fmt.Errorf("failed to patch eip %s, %w", dnat.Spec.EIP, err)
				klog.Error(err)
				return err
			}
		}
		if err = c.patchDnatStatus(key, "", "", "", redo, false); err != nil {
			err = fmt.Errorf("failed to patch dnat %s, %w", key, err)
			klog.Error(err)
			return err
		}
	}
	return nil
}

func (c *Controller) patchSnatLabel(key string, eip *kubeovnv1.IptablesEIP) error {
	oriSnat, err := c.iptablesSnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	snat := oriSnat.DeepCopy()
	var needUpdateLabel, needUpdateAnno bool
	var op string
	eipMember := getMemberIDFromMeta(eip.Labels, eip.Annotations)
	currentMember := snat.Labels[util.NatGatewayMemberLabel]

	if len(snat.Labels) == 0 {
		op = "add"
		snat.Labels = map[string]string{
			util.VpcNatGatewayNameLabel: eip.Spec.NatGwDp,
			util.EipV4IpLabel:           eip.Spec.V4ip,
			util.EipUIDLabel:            string(eip.UID),
		}
		if eipMember != "" {
			snat.Labels[util.NatGatewayMemberLabel] = eipMember
		}
		needUpdateLabel = true
	} else if snat.Labels[util.VpcNatGatewayNameLabel] != eip.Spec.NatGwDp ||
		snat.Labels[util.EipV4IpLabel] != eip.Spec.V4ip ||
		snat.Labels[util.EipUIDLabel] != string(eip.UID) ||
		currentMember != eipMember {
		op = "replace"
		snat.Labels[util.VpcNatGatewayNameLabel] = eip.Spec.NatGwDp
		snat.Labels[util.EipV4IpLabel] = eip.Spec.V4ip
		snat.Labels[util.EipUIDLabel] = string(eip.UID)
		if eipMember != "" {
			snat.Labels[util.NatGatewayMemberLabel] = eipMember
		} else {
			delete(snat.Labels, util.NatGatewayMemberLabel)
		}
		needUpdateLabel = true
	}
	if needUpdateLabel {
		if err := c.updateIptableLabels(snat.Name, op, util.SnatUsingEip, snat.Labels); err != nil {
			klog.Error(err)
			return err
		}
	}

	if len(snat.Annotations) == 0 {
		op = "add"
		needUpdateAnno = true
		snat.Annotations = map[string]string{
			util.VpcEipAnnotation: eip.Name,
		}
	} else if snat.Annotations[util.VpcEipAnnotation] != eip.Name {
		op = "replace"
		needUpdateAnno = true
		snat.Annotations[util.VpcEipAnnotation] = eip.Name
	}
	if needUpdateAnno {
		if err := c.updateIptableAnnotations(snat.Name, op, util.SnatUsingEip, snat.Annotations); err != nil {
			klog.Error(err)
			return err
		}
	}
	return nil
}

func (c *Controller) patchSnatStatus(key, v4ip, v6ip, natGwDp, redo string, ready bool) error {
	oriSnat, err := c.iptablesSnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	snat := oriSnat.DeepCopy()
	var changed bool
	if snat.Status.Ready != ready {
		snat.Status.Ready = ready
		changed = true
	}
	if redo != "" && snat.Status.Redo != redo {
		snat.Status.Redo = redo
		changed = true
	}
	if ready && v4ip != "" && snat.Status.V4ip != v4ip {
		snat.Status.V4ip = v4ip
		snat.Status.V6ip = v6ip
		snat.Status.NatGwDp = natGwDp
		changed = true
	}
	if ready && snat.Spec.InternalCIDR != "" {
		v4CidrSpec, _ := util.SplitStringIP(snat.Spec.InternalCIDR)
		if v4CidrSpec != "" {
			v4Cidr, _ := util.SplitStringIP(snat.Status.InternalCIDR)
			if v4Cidr != v4CidrSpec {
				snat.Status.InternalCIDR = v4CidrSpec
				changed = true
			}
		}
	}

	if changed {
		bytes, err := snat.Status.Bytes()
		if err != nil {
			klog.Error(err)
			return err
		}
		if _, err = c.config.KubeOvnClient.KubeovnV1().IptablesSnatRules().Patch(context.Background(), snat.Name,
			types.MergePatchType, bytes, metav1.PatchOptions{}, "status"); err != nil {
			if k8serrors.IsNotFound(err) {
				return nil
			}
			klog.Errorf("failed to patch snat %s, %v", snat.Name, err)
			return err
		}
	}
	return nil
}

func (c *Controller) redoSnat(key, redo string, eipReady bool) error {
	snat, err := c.iptablesSnatRulesLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to get snat %s, %v", key, err)
		return err
	}
	if redo != "" && redo != snat.Status.Redo {
		if !eipReady {
			if err = c.patchEipStatus(snat.Spec.EIP, "", redo, "", false); err != nil {
				err = fmt.Errorf("failed to patch eip %s, %w", snat.Spec.EIP, err)
				klog.Error(err)
				return err
			}
		}
		if err = c.patchSnatStatus(key, "", "", "", redo, false); err != nil {
			err = fmt.Errorf("failed to patch snat %s, %w", key, err)
			klog.Error(err)
			return err
		}
	}
	return nil
}

func (c *Controller) createFipInPod(dp, v4ip, internalIP string) error {
	gwPods, err := c.getNatGwPods(dp, c.natGwNamespaceByName(dp), false)
	if err != nil {
		klog.Error(err)
		return err
	}
	var addRules []string
	rule := fmt.Sprintf("%s,%s", v4ip, internalIP)
	addRules = append(addRules, rule)

	stateless := c.isStatelessDpMode(dp)
	op := natGwSubnetFipAdd
	if stateless {
		op = natGwStatelessFipAdd
	}

	var firstErr error
	for _, gwPod := range gwPods {
		if err = c.execNatGwRules(gwPod, op, addRules); err != nil {
			klog.Errorf("failed to create fip in pod %s/%s, err: %v", gwPod.Namespace, gwPod.Name, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// finalDeleteFipInPod resolves (natGwDp, v4ip) from the FIP CR's Status,
// with best-effort fallback to the EIP resource when Status is incomplete.
// Then delegates to deleteFipInPod to execute the actual shell deletion.
//
// Used by both the delete path (Status may be empty) and the spec-change path
// (Status guaranteed non-empty by the V4ip guard).
func (c *Controller) finalDeleteFipInPod(key string, cachedFip *kubeovnv1.IptablesFIPRule) error {
	klog.V(3).Infof("final delete fip '%s' in pod", key)
	var firstErr error
	statusV4ip := cachedFip.Status.V4ip
	statusNatGwDp := cachedFip.Status.NatGwDp
	if statusV4ip == "" {
		klog.Warningf("fip %s has empty Status.V4ip, fallback to eip %s", key, cachedFip.Spec.EIP)
		eip, err := c.GetEip(cachedFip.Spec.EIP)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				klog.Errorf("fip %s: eip %s not found, skip pod cleanup", key, cachedFip.Spec.EIP)
				return nil
			}
			klog.Errorf("failed to get eip %s for fip %s, %v", cachedFip.Spec.EIP, key, err)
			return err
		}
		statusV4ip = eip.Status.IP
		if statusNatGwDp == "" {
			statusNatGwDp = eip.Spec.NatGwDp
		}
	}
	if statusV4ip == "" || statusNatGwDp == "" {
		klog.Warningf("fip %s: skip status-based cleanup due to incomplete identity (v4ip=%q, natGwDp=%q)", key, statusV4ip, statusNatGwDp)
	} else if err := c.deleteFipInPod(statusNatGwDp, statusV4ip); err != nil {
		klog.Errorf("failed to delete fip %s, %v", key, err)
		firstErr = err
	}

	// Spec-change crash: Status has old IP (V4ip != "") but Ready=false means a spec
	// change crashed midway. Pod may have a new-IP rule while Status still points to the old IP.
	if !cachedFip.Status.Ready && cachedFip.Status.V4ip != "" {
		eip, err := c.GetEip(cachedFip.Spec.EIP)
		if err != nil {
			if !k8serrors.IsNotFound(err) {
				return err
			}
			klog.Warningf("fip %s not ready: eip %s not found, skip spec-based cleanup", key, cachedFip.Spec.EIP)
			return firstErr
		}
		specV4ip := eip.Status.IP
		specNatGwDp := eip.Spec.NatGwDp
		if specV4ip == "" || specNatGwDp == "" {
			klog.Warningf("fip %s not ready: skip spec-based cleanup due to incomplete spec identity (v4ip=%q, natGwDp=%q)", key, specV4ip, specNatGwDp)
			return firstErr
		}
		if specV4ip != statusV4ip || specNatGwDp != statusNatGwDp {
			if err = c.deleteFipInPod(specNatGwDp, specV4ip); err != nil {
				klog.Errorf("failed spec-based cleanup for fip %s, %v", key, err)
				if firstErr == nil {
					firstErr = err
				}
			}
		}
	}
	return firstErr
}

// finalDeleteDnatInPod resolves (natGwDp, protocol, v4ip, externalPort) from the DNAT CR's
// Status, with best-effort fallback to EIP/Spec when Status is incomplete.
// Then delegates to deleteDnatInPod to execute the actual shell deletion.
func (c *Controller) finalDeleteDnatInPod(key string, cachedDnat *kubeovnv1.IptablesDnatRule) error {
	klog.V(3).Infof("final delete dnat '%s' in pod", key)
	var firstErr error
	statusV4ip := cachedDnat.Status.V4ip
	statusNatGwDp := cachedDnat.Status.NatGwDp
	if statusV4ip == "" && !dnatUsesEip(&cachedDnat.Spec) && dnatServesClusterIP(&cachedDnat.Spec) {
		// A rule serving a ClusterIP carries its identity itself, so Status is only needed for
		// the fields the data plane was actually programmed with.
		klog.Warningf("dnat %s has empty Status.V4ip, fallback to clusterIP %s", key, cachedDnat.Spec.ClusterIP)
		statusV4ip = cachedDnat.Spec.ClusterIP
		if statusNatGwDp == "" {
			statusNatGwDp = cachedDnat.Labels[util.VpcNatGatewayNameLabel]
		}
	} else if statusV4ip == "" {
		klog.Warningf("dnat %s has empty Status.V4ip, fallback to eip %s", key, cachedDnat.Spec.EIP)
		eip, err := c.GetEip(cachedDnat.Spec.EIP)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				klog.Errorf("dnat %s: eip %s not found, skip pod cleanup", key, cachedDnat.Spec.EIP)
				return nil
			}
			return err
		}
		statusV4ip = eip.Status.IP
		if statusNatGwDp == "" {
			statusNatGwDp = eip.Spec.NatGwDp
		}
	}
	statusProtocol := cachedDnat.Status.Protocol
	if statusProtocol == "" {
		klog.Warningf("dnat %s has empty Status.Protocol, fallback to Spec", key)
		statusProtocol = cachedDnat.Spec.Protocol
	}
	if statusProtocol == "" {
		klog.Errorf("dnat %s has v4ip %s but protocol is empty in both Status and Spec, skip pod cleanup", key, statusV4ip)
		return nil
	}
	statusExternalPort := cachedDnat.Status.ExternalPort
	if statusExternalPort == "" {
		klog.Warningf("dnat %s has empty Status.ExternalPort, fallback to Spec", key)
		statusExternalPort = cachedDnat.Spec.ExternalPort
	}
	if statusExternalPort == "" {
		klog.Errorf("dnat %s has v4ip %s but externalPort is empty in both Status and Spec, skip pod cleanup", key, statusV4ip)
		return nil
	}
	statusInternalIP := cachedDnat.Status.InternalIP
	if statusInternalIP == "" {
		statusInternalIP = cachedDnat.Spec.InternalIP
	}
	statusInternalPort := cachedDnat.Status.InternalPort
	if statusInternalPort == "" {
		statusInternalPort = cachedDnat.Spec.InternalPort
	}
	if statusV4ip == "" || statusNatGwDp == "" {
		klog.Warningf("dnat %s: skip status-based cleanup due to incomplete identity (v4ip=%q, natGwDp=%q)", key, statusV4ip, statusNatGwDp)
	} else if err := c.deleteDnatInPodWithInternal(statusNatGwDp, statusProtocol,
		statusV4ip, statusExternalPort, statusInternalIP, statusInternalPort); err != nil {
		klog.Errorf("failed to delete dnat %s, %v", key, err)
		firstErr = err
	}

	// Spec-change crash: Status has old IP (V4ip != "") but Ready=false means a spec
	// change crashed midway. Pod may have a new-IP rule while Status still points to the old IP.
	if dnatNeedsSpecCleanup(cachedDnat) {
		eip, err := c.GetEip(cachedDnat.Spec.EIP)
		if err != nil {
			if !k8serrors.IsNotFound(err) {
				return err
			}
			klog.Warningf("dnat %s not ready: eip %s not found, skip spec-based cleanup", key, cachedDnat.Spec.EIP)
			return firstErr
		}
		specV4ip := eip.Status.IP
		specNatGwDp := eip.Spec.NatGwDp
		specProtocol := cachedDnat.Spec.Protocol
		specExternalPort := cachedDnat.Spec.ExternalPort
		if specV4ip == "" || specNatGwDp == "" || specProtocol == "" || specExternalPort == "" {
			klog.Warningf("dnat %s not ready: skip spec-based cleanup due to incomplete spec identity (v4ip=%q, natGwDp=%q, protocol=%q, externalPort=%q)",
				key, specV4ip, specNatGwDp, specProtocol, specExternalPort)
			return firstErr
		}
		if specV4ip != statusV4ip || specNatGwDp != statusNatGwDp || specProtocol != statusProtocol || specExternalPort != statusExternalPort {
			if err = c.deleteDnatInPodWithInternal(specNatGwDp, specProtocol, specV4ip, specExternalPort, cachedDnat.Spec.InternalIP, cachedDnat.Spec.InternalPort); err != nil {
				klog.Errorf("failed spec-based cleanup for dnat %s, %v", key, err)
				if firstErr == nil {
					firstErr = err
				}
			}
		}
	}
	return firstErr
}

func dnatCleanupEipName(dnat *kubeovnv1.IptablesDnatRule) string {
	if eipName := dnat.Annotations[util.VpcEipAnnotation]; eipName != "" {
		return eipName
	}
	return dnat.Spec.EIP
}

// finalDeleteSnatInPod resolves (natGwDp, v4ip, v4Cidr) from the SNAT CR's Status,
// with best-effort fallback to EIP/Spec when Status is incomplete.
// Then delegates to deleteSnatInPod to execute the actual shell deletion.
func (c *Controller) finalDeleteSnatInPod(key string, cachedSnat *kubeovnv1.IptablesSnatRule) error {
	klog.V(3).Infof("final delete snat '%s' in pod", key)
	var firstErr error
	statusV4ip := cachedSnat.Status.V4ip
	statusNatGwDp := cachedSnat.Status.NatGwDp
	if statusV4ip == "" {
		klog.Warningf("snat %s has empty Status.V4ip, fallback to eip %s", key, cachedSnat.Spec.EIP)
		eip, err := c.GetEip(cachedSnat.Spec.EIP)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				klog.Errorf("snat %s: eip %s not found, skip pod cleanup", key, cachedSnat.Spec.EIP)
				return nil
			}
			return err
		}
		statusV4ip = eip.Status.IP
		if statusNatGwDp == "" {
			statusNatGwDp = eip.Spec.NatGwDp
		}
	}
	statusV4Cidr, _ := util.SplitStringIP(cachedSnat.Status.InternalCIDR)
	if statusV4Cidr == "" {
		klog.Warningf("snat %s has empty Status.InternalCIDR, fallback to Spec", key)
		statusV4Cidr, _ = util.SplitStringIP(cachedSnat.Spec.InternalCIDR)
	}
	if statusV4Cidr == "" {
		klog.Errorf("snat %s has v4ip %s but v4Cidr is empty in both Status and Spec, skip pod cleanup", key, statusV4ip)
		return nil
	}
	var eip *kubeovnv1.IptablesEIP
	if cachedSnat.Spec.EIP != "" {
		eip, _ = c.iptablesEipsLister.Get(cachedSnat.Spec.EIP)
	}
	memberID := resolveRecordedSnatMemberID(cachedSnat, eip)
	if statusV4ip == "" || statusNatGwDp == "" {
		klog.Warningf("snat %s: skip status-based cleanup due to incomplete identity (v4ip=%q, natGwDp=%q)", key, statusV4ip, statusNatGwDp)
	} else if err := c.deleteSnatInPodWithMember(statusNatGwDp, statusV4ip, statusV4Cidr, memberID); err != nil {
		klog.Errorf("failed to delete snat %s, %v", key, err)
		firstErr = err
	}

	// Spec-change crash: Status has old IP (V4ip != "") but Ready=false means a spec
	// change crashed midway. Pod may have a new-IP rule while Status still points to the old IP.
	if !cachedSnat.Status.Ready && cachedSnat.Status.V4ip != "" {
		eip, err := c.GetEip(cachedSnat.Spec.EIP)
		if err != nil {
			if !k8serrors.IsNotFound(err) {
				return err
			}
			klog.Warningf("snat %s not ready: eip %s not found, skip spec-based cleanup", key, cachedSnat.Spec.EIP)
			return firstErr
		}
		specV4ip := eip.Status.IP
		specNatGwDp := eip.Spec.NatGwDp
		specV4Cidr, _ := util.SplitStringIP(cachedSnat.Spec.InternalCIDR)
		if specV4ip == "" || specNatGwDp == "" || specV4Cidr == "" {
			klog.Warningf("snat %s not ready: skip spec-based cleanup due to incomplete spec identity (v4ip=%q, natGwDp=%q, v4Cidr=%q)",
				key, specV4ip, specNatGwDp, specV4Cidr)
			return firstErr
		}
		if specV4ip != statusV4ip || specNatGwDp != statusNatGwDp || specV4Cidr != statusV4Cidr {
			specMemberID := resolveSnatMemberID(eip, cachedSnat)
			if err = c.deleteSnatInPodWithMember(specNatGwDp, specV4ip, specV4Cidr, specMemberID); err != nil {
				klog.Errorf("failed spec-based cleanup for snat %s, %v", key, err)
				if firstErr == nil {
					firstErr = err
				}
			}
		}
	}
	return firstErr
}

func (c *Controller) deleteFipInPod(dp, v4ip string) error {
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
	gwPods, err := c.getNatGwPods(dp, c.natGwNamespaceByName(dp), false)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			klog.V(4).Infof("nat gw pod %s not found, will retry fip pod cleanup", dp)
		} else {
			klog.Error(err)
		}
		return err
	}
	// del_floating_ip matches by EIP only (FIP is 1:1, identity = EIP)
	stateless := c.isStatelessDpMode(dp)
	op := natGwSubnetFipDel
	rules := []string{v4ip}
	if stateless {
		op = natGwStatelessFipDel
		// In stateless mode, rule matches "eip,internalIp" or eip
		rules = []string{fmt.Sprintf("%s,", v4ip)}
	}
	var firstErr error
	for _, gwPod := range gwPods {
		if err = c.execNatGwRules(gwPod, op, rules); err != nil {
			klog.Errorf("failed to delete fip in pod %s/%s, err: %v", gwPod.Namespace, gwPod.Name, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (c *Controller) createDnatInPod(dp, protocol, v4ip, internalIP, externalPort, internalPort string) error {
	gwPods, err := c.getNatGwPods(dp, c.natGwNamespaceByName(dp), false)
	if err != nil {
		klog.Errorf("failed to get nat gw pods, %v", err)
		return err
	}
	var addRules []string
	rule := fmt.Sprintf("%s,%s,%s,%s,%s", v4ip, externalPort, protocol, internalIP, internalPort)
	addRules = append(addRules, rule)

	stateless := c.isStatelessDpMode(dp)
	op := natGwDnatAdd
	if stateless {
		op = natGwStatelessDnatAdd
	}

	var firstErr error
	for _, gwPod := range gwPods {
		if err = c.execNatGwRules(gwPod, op, addRules); err != nil {
			klog.Errorf("failed to create dnat in pod %s/%s, err: %v", gwPod.Namespace, gwPod.Name, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (c *Controller) deleteDnatInPod(dp, protocol, v4ip, externalPort string) error {
	return c.deleteDnatInPodWithInternal(dp, protocol, v4ip, externalPort, "", "")
}

func (c *Controller) deleteDnatInPodWithInternal(dp, protocol, v4ip, externalPort, internalIP, internalPort string) error {
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
	gwPods, err := c.getNatGwPods(dp, c.natGwNamespaceByName(dp), false)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			klog.V(4).Infof("nat gw pod %s not found, will retry dnat pod cleanup", dp)
		} else {
			klog.Error(err)
		}
		return err
	}

	stateless := c.isStatelessDpMode(dp)
	op := natGwDnatDel
	var rules []string
	if stateless {
		op = natGwStatelessDnatDel
		rules = []string{fmt.Sprintf("%s,%s,%s,%s,%s", v4ip, externalPort, protocol, internalIP, internalPort)}
	} else {
		// del_dnat matches by identity triplet (EIP, ExternalPort, Protocol) only
		rules = []string{fmt.Sprintf("%s,%s,%s", v4ip, externalPort, protocol)}
	}

	var firstErr error
	for _, gwPod := range gwPods {
		if err = c.execNatGwRules(gwPod, op, rules); err != nil {
			klog.Errorf("failed to delete dnat in pod %s/%s, err: %v", gwPod.Namespace, gwPod.Name, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (c *Controller) createSnatInPod(dp, v4ip, internalCIDR string) error {
	return c.createSnatInPodWithMember(dp, v4ip, internalCIDR, "")
}

func (c *Controller) createSnatInPodWithMember(dp, v4ip, internalCIDR, memberID string) error {
	internalCIDR = normalizeSnatInternalCIDR(internalCIDR)
	gwPods, err := c.getNatGwPods(dp, c.natGwNamespaceByName(dp), false)
	if err != nil {
		klog.Errorf("failed to get nat gw pods, %v", err)
		return err
	}
	gwPods, err = c.filterNatGwPodsByMember(gwPods, memberID)
	if err != nil {
		klog.Errorf("failed to filter nat gw pods for member %s: %v", memberID, err)
		return err
	}

	stateless := c.isStatelessDpMode(dp)
	op := natGwSnatAdd
	if stateless {
		op = natGwStatelessSnatAdd
	}

	var firstErr error
	for _, gwPod := range gwPods {
		var rules []string
		rule := fmt.Sprintf("%s,%s", v4ip, internalCIDR)

		if !stateless {
			version, err := c.getIptablesVersion(gwPod)
			if err != nil {
				version = "1.0.0"
				klog.Warningf("failed to checking iptables version, assuming version at least %s: %v", version, err)
			}
			if util.CompareVersion(version, "1.6.2") >= 1 {
				rule = fmt.Sprintf("%s,%s", rule, "--random-fully")
			}
		}

		rules = append(rules, rule)
		if err = c.execNatGwRules(gwPod, op, rules); err != nil {
			klog.Errorf("failed to exec nat gateway rule in pod %s/%s, err: %v", gwPod.Namespace, gwPod.Name, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (c *Controller) deleteSnatInPod(dp, v4ip, internalCIDR string) error {
	return c.deleteSnatInPodWithMember(dp, v4ip, internalCIDR, "")
}

func (c *Controller) deleteSnatInPodWithMember(dp, v4ip, internalCIDR, memberID string) error {
	internalCIDR = normalizeSnatInternalCIDR(internalCIDR)

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
	gwPods, err := c.getNatGwPods(dp, c.natGwNamespaceByName(dp), false)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			klog.V(4).Infof("nat gw pod %s not found, will retry snat pod cleanup", dp)
		} else {
			klog.Error(err)
		}
		return err
	}
	gwPods, err = c.filterNatGwPodsByMember(gwPods, memberID)
	if err != nil {
		klog.Errorf("failed to filter nat gw pods for member %s during delete: %v", memberID, err)
		return err
	}

	stateless := c.isStatelessDpMode(dp)
	op := natGwSnatDel
	if stateless {
		op = natGwStatelessSnatDel
	}

	// del nat
	var delRules []string
	rule := fmt.Sprintf("%s,%s", v4ip, internalCIDR)
	delRules = append(delRules, rule)
	var firstErr error
	for _, gwPod := range gwPods {
		if err = c.execNatGwRules(gwPod, op, delRules); err != nil {
			klog.Errorf("failed to delete snat in pod %s/%s, err: %v", gwPod.Namespace, gwPod.Name, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (c *Controller) updateIptableLabels(name, op, natType string, labels map[string]string) error {
	patchPayloadTemplate := `[{ "op": "%s", "path": "/metadata/labels", "value": %s }]`
	raw, _ := json.Marshal(labels)
	patchPayload := fmt.Sprintf(patchPayloadTemplate, op, raw)

	if err := c.patchIptableInfo(name, natType, patchPayload); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to patch label for %s %s, %v", natType, name, err)
		return err
	}
	return nil
}

func (c *Controller) updateIptableAnnotations(name, op, natType string, anno map[string]string) error {
	patchPayloadTemplate := `[{ "op": "%s", "path": "/metadata/annotations", "value": %s }]`
	raw, _ := json.Marshal(anno)
	patchPayload := fmt.Sprintf(patchPayloadTemplate, op, raw)

	if err := c.patchIptableInfo(name, natType, patchPayload); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("failed to patch annotations for %s %s, %v", natType, name, err)
		return err
	}
	return nil
}

func (c *Controller) patchIptableInfo(name, natType, patchPayload string) error {
	var err error
	switch natType {
	case util.FipUsingEip:
		_, err = c.config.KubeOvnClient.KubeovnV1().IptablesFIPRules().Patch(context.Background(), name,
			types.JSONPatchType, []byte(patchPayload), metav1.PatchOptions{})
	case util.SnatUsingEip:
		_, err = c.config.KubeOvnClient.KubeovnV1().IptablesSnatRules().Patch(context.Background(), name,
			types.JSONPatchType, []byte(patchPayload), metav1.PatchOptions{})
	case util.DnatUsingEip:
		_, err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Patch(context.Background(), name,
			types.JSONPatchType, []byte(patchPayload), metav1.PatchOptions{})
	case "eip":
		_, err = c.config.KubeOvnClient.KubeovnV1().IptablesEIPs().Patch(context.Background(), name,
			types.JSONPatchType, []byte(patchPayload), metav1.PatchOptions{})
	default:
		// Silently writing nothing would strand the startup backfill until its poll deadline.
		return fmt.Errorf("unknown nat type %s", natType)
	}
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	return nil
}

// validateDnatRule validates IptablesDnatRule fields to prevent malformed iptables commands.
func (c *Controller) validateDnatRule(dnat *kubeovnv1.IptablesDnatRule) error {
	var err error
	// A rule addresses at least one VIP: a public one through an EIP (bound on the external
	// interface) and/or the internal ClusterIP (held on lo). A Service handled by the nftable LB
	// service feature carries both on one rule. The webhook enforces this too; the check is
	// repeated here because the controller must not program a rule whose identity it cannot resolve.
	if !dnatUsesEip(&dnat.Spec) && !dnatServesClusterIP(&dnat.Spec) {
		err = fmt.Errorf("%s: one of eip and clusterIP must be set", dnat.Name)
		klog.Error(err)
		return err
	}
	if dnatServesClusterIP(&dnat.Spec) {
		if !dnatUsesEip(&dnat.Spec) && dnat.Labels[util.VpcNatGatewayNameLabel] == "" {
			err = fmt.Errorf("%s: gateway label is required with clusterIP when there is no eip", dnat.Name)
			klog.Error(err)
			return err
		}
		if util.CheckProtocol(dnat.Spec.ClusterIP) != kubeovnv1.ProtocolIPv4 {
			err = fmt.Errorf("%s: clusterIP %q must be IPv4, share dnat is IPv4 only", dnat.Name, dnat.Spec.ClusterIP)
			klog.Error(err)
			return err
		}
		if dnat.Spec.Type != kubeovnv1.DnatRuleTypeShare {
			err = fmt.Errorf("%s: clusterIP requires type=%s: the address is shared by all backends", dnat.Name, kubeovnv1.DnatRuleTypeShare)
			klog.Error(err)
			return err
		}
	}
	if err = util.ValidatePort(dnat.Spec.ExternalPort); err != nil {
		err = fmt.Errorf("%s: invalid externalPort %q: %w", dnat.Name, dnat.Spec.ExternalPort, err)
		klog.Error(err)
		return err
	}
	if err = util.ValidatePort(dnat.Spec.InternalPort); err != nil {
		err = fmt.Errorf("%s: invalid internalPort %q: %w", dnat.Name, dnat.Spec.InternalPort, err)
		klog.Error(err)
		return err
	}
	if dnat.Spec.InternalIP == "" {
		err = fmt.Errorf("%s: internalIP cannot be empty", dnat.Name)
		klog.Error(err)
		return err
	}
	if !util.IsValidIP(dnat.Spec.InternalIP) {
		err = fmt.Errorf("%s: invalid internalIP %q", dnat.Name, dnat.Spec.InternalIP)
		klog.Error(err)
		return err
	}
	// iptables NAT only supports IPv4
	if util.CheckProtocol(dnat.Spec.InternalIP) != kubeovnv1.ProtocolIPv4 {
		err = fmt.Errorf("%s: internalIP %q must be IPv4, IPv6 is not supported", dnat.Name, dnat.Spec.InternalIP)
		klog.Error(err)
		return err
	}
	if dnat.Spec.Protocol != util.ProtocolTCP && dnat.Spec.Protocol != util.ProtocolUDP {
		err = fmt.Errorf("%s: invalid protocol %q: protocol must be lowercase tcp or udp", dnat.Name, dnat.Spec.Protocol)
		klog.Error(err)
		return err
	}
	return nil
}

// validateFipRule validates IptablesFIPRule fields to prevent malformed iptables commands.
func (c *Controller) validateFipRule(fip *kubeovnv1.IptablesFIPRule) error {
	var err error
	if fip.Spec.EIP == "" {
		err = fmt.Errorf("%s: eip cannot be empty", fip.Name)
		klog.Error(err)
		return err
	}
	if fip.Spec.InternalIP == "" {
		err = fmt.Errorf("%s: internalIP cannot be empty", fip.Name)
		klog.Error(err)
		return err
	}
	if !util.IsValidIP(fip.Spec.InternalIP) {
		err = fmt.Errorf("%s: invalid internalIP %q", fip.Name, fip.Spec.InternalIP)
		klog.Error(err)
		return err
	}
	// iptables NAT only supports IPv4
	if util.CheckProtocol(fip.Spec.InternalIP) != kubeovnv1.ProtocolIPv4 {
		err = fmt.Errorf("%s: internalIP %q must be IPv4, IPv6 is not supported", fip.Name, fip.Spec.InternalIP)
		klog.Error(err)
		return err
	}
	return nil
}

// validateSnatRule validates IptablesSnatRule fields to prevent malformed iptables commands.
func (c *Controller) validateSnatRule(snat *kubeovnv1.IptablesSnatRule) error {
	var err error
	if snat.Spec.EIP == "" {
		err = fmt.Errorf("%s: eip cannot be empty", snat.Name)
		klog.Error(err)
		return err
	}
	internalCIDR := snat.Spec.InternalCIDR
	if internalCIDR == "" {
		err = fmt.Errorf("%s: internalCIDR cannot be empty", snat.Name)
		klog.Error(err)
		return err
	}
	// iptables NAT only supports single IPv4 CIDR or IP, ip6tables is not used here
	if strings.Count(internalCIDR, "/") > 1 {
		err = fmt.Errorf("%s: internalCIDR %q contains multiple CIDRs, only single CIDR or IP is allowed", snat.Name, internalCIDR)
		klog.Error(err)
		return err
	}
	if strings.Contains(internalCIDR, "/") {
		if err = util.CheckCidrs(internalCIDR); err != nil {
			err = fmt.Errorf("%s: invalid internalCIDR %q: %w", snat.Name, internalCIDR, err)
			klog.Error(err)
			return err
		}
	} else {
		if !util.IsValidIP(internalCIDR) {
			err = fmt.Errorf("%s: invalid internalCIDR %q", snat.Name, internalCIDR)
			klog.Error(err)
			return err
		}
	}
	if util.CheckProtocol(internalCIDR) != kubeovnv1.ProtocolIPv4 {
		err = fmt.Errorf("%s: internalCIDR %q must be IPv4, IPv6 is not supported", snat.Name, internalCIDR)
		klog.Error(err)
		return err
	}
	return nil
}

func (c *Controller) isStatelessDpMode(dp string) bool {
	gw, err := c.vpcNatGatewayLister.Get(dp)
	if err != nil {
		return false
	}
	if gw.Annotations != nil {
		if mode := gw.Annotations[util.NatGatewayDataplaneModeAnnotation]; mode == "stateless" || mode == "stateless-nft" {
			return true
		}
	}
	return false
}

func (c *Controller) filterNatGwPodsByMember(pods []*corev1.Pod, memberID string) ([]*corev1.Pod, error) {
	if memberID == "" {
		return pods, nil
	}
	var matched []*corev1.Pod
	for _, p := range pods {
		if p.Labels != nil {
			if m := p.Labels[util.NatGatewayMemberLabel]; m == memberID {
				matched = append(matched, p)
				continue
			}
			if m := p.Labels[util.NatGatewayMemberLegacyLabel]; m == memberID {
				matched = append(matched, p)
				continue
			}
		}
	}
	if len(matched) > 0 {
		return matched, nil
	}
	return nil, fmt.Errorf("nat gw member pod %q is not found or not ready", memberID)
}

// normalizeSnatInternalCIDR converts a bare IPv4 (e.g. "10.0.0.5") — a shape
// accepted by validateSnatRule — to its canonical "<ip>/32" form so the NAT
// gateway script can assume every SNAT rule carries an explicit prefix length.
// This keeps downstream logic (longest-prefix ordering, idempotency checks)
// free of bare-IP special cases.
func normalizeSnatInternalCIDR(cidr string) string {
	if cidr == "" || strings.Contains(cidr, "/") {
		return cidr
	}
	return cidr + "/32"
}

// getMemberIDFromMeta extracts the gateway member identifier from labels and annotations,
// honoring both standard and legacy label keys.
func getMemberIDFromMeta(labels, annotations map[string]string) string {
	if labels != nil {
		if m := labels[util.NatGatewayMemberLabel]; m != "" {
			return m
		}
		if m := labels[util.NatGatewayMemberLegacyLabel]; m != "" {
			return m
		}
	}
	if annotations != nil {
		if m := annotations[util.NatGatewayMemberLabel]; m != "" {
			return m
		}
		if m := annotations[util.NatGatewayMemberLegacyLabel]; m != "" {
			return m
		}
	}
	return ""
}

// resolveSnatMemberID resolves the assigned NAT gateway member identifier from the associated
// EIP or SNAT rule metadata. EIP ownership takes precedence, falling back to SNAT rule metadata
// only when no EIP is associated.
func resolveSnatMemberID(eip *kubeovnv1.IptablesEIP, snat *kubeovnv1.IptablesSnatRule) string {
	if eip != nil {
		return getMemberIDFromMeta(eip.Labels, eip.Annotations)
	}
	if snat != nil {
		return getMemberIDFromMeta(snat.Labels, snat.Annotations)
	}
	return ""
}

// resolveRecordedSnatMemberID resolves the member recorded on the SNAT rule itself.
// When deleting or cleaning up an existing recorded rule, this ensures that the pod member
// where the rule was actually deployed is targeted, even if the EIP has since been reassigned or unassigned.
func resolveRecordedSnatMemberID(snat *kubeovnv1.IptablesSnatRule, eip *kubeovnv1.IptablesEIP) string {
	if snat != nil {
		if m := getMemberIDFromMeta(snat.Labels, snat.Annotations); m != "" {
			return m
		}
	}
	if eip != nil {
		if m := getMemberIDFromMeta(eip.Labels, eip.Annotations); m != "" {
			return m
		}
	}
	return ""
}
