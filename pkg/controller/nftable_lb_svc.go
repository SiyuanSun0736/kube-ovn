package controller

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// The nftable LB service feature makes a vpc-nat-gw act like kube-proxy. Service and
// EndpointSlice are the source of truth and the only trigger: this controller derives each
// VIP:port:protocol identity and its complete backend set, then writes the gateway nft map,
// hairpin rule, loopback VIP and VPC route directly.
//
// It owns exactly the identities its flags name (see validateServiceFeatureGates for the
// partition): a LoadBalancer Service's EIP identities with --enable-gw-nftable-lb-svc, and
// ClusterIP identities with --enable-gw-nftable-svc-cluster-ip. The classic OVN load balancer
// and the per-Service LB Pod implementations keep managing the parts they own; in particular an
// --enable-lb cluster keeps serving a gateway Service's ClusterIP from the OVN switch LBs.
//
// Binding model:
//   - A Service is handled when it names its gateway in util.VpcNatGatewayAnnotation.
//     A LoadBalancer Service also names the EIP that supplies its ingress address; a
//     ClusterIP Service uses its ClusterIP as the only VIP.
//   - One IptablesDnatRule record is kept per Service port/backend so operators can query the
//     programmed forwarding in one layer and EIP accounting continues to work. A type=share
//     object is only this ledger: its add/update/delete events never write or remove NAT state.
//     The write order is claim, rule, Ready: a record is persisted before the gateway rule it
//     accounts for, flips Ready only after the data plane converged, and a stale record (or the
//     retired-generation snapshot of an identity-changing update) is retired only after its rule
//     is gone, so cleanup always finds a claim for what exists.
//   - The Kube-OVN controller finalizer claims the whole data plane. Cleanup removes nft identities, hairpin,
//     loopback VIPs and routes, then deletes the records and finally releases the Service.
//   - EIP, Pod and ordinary gateway events do not trigger this controller. Their state is read
//     when a Service or EndpointSlice event reconciles the Service. Gateway instance replacement
//     is the recovery exception: it only re-enqueues bound Services, and records never redo.
//
// Traffic policy scope (by design):
//   This feature intentionally aligns with kube-proxy's *Cluster* traffic policy
//   only. It always load balances across all Ready endpoints of the Service and
//   does NOT implement ExternalTrafficPolicy=Local / InternalTrafficPolicy=Local,
//   topology-aware routing, or terminating-endpoint fallback. Those node-local
//   semantics do not map cleanly onto a centralized vpc-nat-gw (which is not
//   per-node), and Cluster policy is the intended, sufficient behavior here.
//   Consequently client source IP is not preserved (traffic is DNAT'd through the
//   gateway), matching what Cluster policy already implies.
//

// nftableLbSvcCandidate reports whether a Service can be considered by the gateway controller.
func nftableLbSvcCandidate(svc *v1.Service) bool {
	if svc == nil || svc.Annotations[util.VpcNatGatewayAnnotation] == "" {
		return false
	}
	return svc.Spec.Type == v1.ServiceTypeLoadBalancer || svc.Spec.Type == v1.ServiceTypeClusterIP
}

// nftableLbSvcQualifies reports whether at least one enabled gateway identity can serve the
// Service, and which ones. The two identities are gated independently: EnableGwNftableLbSvc
// serves the EIP a LoadBalancer Service publishes, EnableGwNftableSvcClusterIP serves the
// ClusterIP. One Service can have both, and they share a single record and data plane.
func (c *Controller) nftableLbSvcQualifies(svc *v1.Service) (serveEIP, serveClusterIP bool) {
	if !nftableLbSvcCandidate(svc) {
		return false, false
	}
	serveEIP = svc.Spec.Type == v1.ServiceTypeLoadBalancer && c.config.EnableGwNftableLbSvc &&
		svc.Annotations[util.EipAnnotation] != ""
	serveClusterIP = (svc.Spec.Type == v1.ServiceTypeLoadBalancer || svc.Spec.Type == v1.ServiceTypeClusterIP) &&
		c.config.EnableGwNftableSvcClusterIP
	return serveEIP, serveClusterIP
}

// nftableLbSvcGateway returns the gateway that serves the Service, or "" when none is named.
func nftableLbSvcGateway(svc *v1.Service) string {
	return svc.Annotations[util.VpcNatGatewayAnnotation]
}

// nftableLbSvcTeardownHoldsEip reports whether a Service that declares this EIP may still owe
// teardown work: a terminating EIP must keep its object (and thereby its address) until no
// referencing Service can still have the EIP identity wired, because the Service's cleanup
// reads the identity witness from the EIP object itself - the accounting records may be gone
// and the ingress may never have been published, so the object is the last source of the
// address. Two signals count as still holding the identity, and both terminate: the Service
// still qualifies for the EIP leg (its reconcile restores or tears down the identity, then
// drops the claim), or one of its records still carries the leg (a live claim or a trim
// tombstone, both retired by a successful reconcile). The bare finalizer must NOT hold: a
// Service serving only its ClusterIP keeps it forever while owing this EIP nothing.
func (c *Controller) nftableLbSvcTeardownHoldsEip(eip *kubeovnv1.IptablesEIP) bool {
	if !c.gwNftableLbSvcEnabled() || c.svcIndexer == nil {
		return false
	}
	svcObjs, err := c.svcIndexer.ByIndex(IndexGwNftableLbServiceByEip, eip.Name)
	if err != nil {
		// Bias towards holding: releasing on an unknown answer strands the teardown witness.
		klog.Errorf("failed to list nftable lb services of eip %s: %v", eip.Name, err)
		return true
	}
	for _, svcObj := range svcObjs {
		svc, ok := svcObj.(*v1.Service)
		if !ok || !slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer) {
			continue
		}
		serving, _ := c.nftableLbSvcQualifies(svc)
		if !serving {
			// The EIP leg is disabled or deconfigured; only still-unretired records (live claims
			// or trim tombstones) mean the identity can still be wired.
			rules, err := c.iptablesDnatRulesLister.List(labels.SelectorFromSet(labels.Set{
				util.NftableLbSvcNsLabel:   svc.Namespace,
				util.NftableLbSvcNameLabel: svc.Name,
			}))
			if err != nil {
				klog.Errorf("failed to list records of nftable lb service %s/%s: %v", svc.Namespace, svc.Name, err)
				return true
			}
			if !slices.ContainsFunc(rules, func(rule *kubeovnv1.IptablesDnatRule) bool {
				return rule.Spec.EIP == eip.Name || rule.Labels[util.EipV4IpLabel] != ""
			}) {
				continue
			}
		}
		klog.Infof("eip %s: service %s/%s still holds its identity, holding the release", eip.Name, svc.Namespace, svc.Name)
		return true
	}
	return false
}

func (c *Controller) enqueueGwNftableLbService(key string) {
	if c.addOrUpdateGwNftableLbSvcQueue == nil || key == "" {
		return
	}
	klog.V(3).Infof("enqueue add/update gateway nftable lb service %s", key)
	c.addOrUpdateGwNftableLbSvcQueue.Add(key)
}

func gwNftableLbSvcChanged(oldSvc, newSvc *v1.Service) bool {
	if !oldSvc.DeletionTimestamp.Equal(newSvc.DeletionTimestamp) ||
		oldSvc.Annotations[util.VpcNatGatewayAnnotation] != newSvc.Annotations[util.VpcNatGatewayAnnotation] ||
		oldSvc.Annotations[util.EipAnnotation] != newSvc.Annotations[util.EipAnnotation] ||
		oldSvc.Spec.Type != newSvc.Spec.Type ||
		!slices.Equal(util.ServiceClusterIPs(*oldSvc), util.ServiceClusterIPs(*newSvc)) ||
		!reflect.DeepEqual(oldSvc.Spec.Ports, newSvc.Spec.Ports) ||
		oldSvc.Spec.SessionAffinity != newSvc.Spec.SessionAffinity ||
		!reflect.DeepEqual(oldSvc.Spec.SessionAffinityConfig, newSvc.Spec.SessionAffinityConfig) {
		return true
	}
	return false
}

func (c *Controller) enqueueGwNftableLbServicesForNatGw(natGwName string) {
	if c.config == nil || (!c.config.EnableGwNftableLbSvc && !c.config.EnableGwNftableSvcClusterIP) || natGwName == "" || c.svcIndexer == nil {
		return
	}
	svcs, err := c.svcIndexer.ByIndex(IndexGwNftableLbServiceByGateway, natGwName)
	if err != nil {
		klog.Errorf("failed to list nftable lb services of nat gw %s: %v", natGwName, err)
		return
	}
	for _, obj := range svcs {
		if svc, ok := obj.(*v1.Service); ok {
			c.enqueueGwNftableLbService(svc.Namespace + "/" + svc.Name)
		}
	}
}

// enqueueGwNftableLbServicesForEIP wakes the gateway-mode Services that reference this EIP. The
// Service reconcile is the only writer of the share records that point at the EIP, so it has to
// learn about an EIP that appeared, became usable, became unusable, or is being deleted: an EIP
// whose referencing records are never released cannot finish deleting.
func (c *Controller) enqueueGwNftableLbServicesForEIP(eipName string) {
	if c.config == nil || (!c.config.EnableGwNftableLbSvc && !c.config.EnableGwNftableSvcClusterIP) || eipName == "" || c.svcIndexer == nil {
		return
	}
	svcs, err := c.svcIndexer.ByIndex(IndexGwNftableLbServiceByEip, eipName)
	if err != nil {
		klog.Errorf("failed to list nftable lb services of eip %s: %v", eipName, err)
		return
	}
	for _, obj := range svcs {
		if svc, ok := obj.(*v1.Service); ok {
			c.enqueueGwNftableLbService(svc.Namespace + "/" + svc.Name)
		}
	}
}

func (c *Controller) handleAddOrUpdateGwNftableLbService(key string) error {
	// With both feature gates off still run: nftableLbSvcQualifies qualifies nothing then, so
	// the reconcile is a teardown-only path for Services the feature had managed, and their
	// controller finalizer would never be released if this handler refused to work.
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	klog.Infof("handle add/update gateway nftable lb service %s", key)

	cachedSvc, err := c.servicesLister.Services(namespace).Get(name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// service deleted: remove all rules owned by it
			return c.cleanupNftableLbService(nil, namespace, name)
		}
		klog.Error(err)
		return err
	}

	serveEIP, serveClusterIP := c.nftableLbSvcQualifies(cachedSvc)
	// service being deleted or no longer qualifying: clean up owned rules (and clear the
	// ingress IP we published, if any)
	if !cachedSvc.DeletionTimestamp.IsZero() || (!serveEIP && !serveClusterIP) {
		return c.cleanupNftableLbService(cachedSvc, namespace, name)
	}

	// The Service names its gateway; a LoadBalancer Service additionally takes its ingress IP
	// from an EIP, which must belong to that gateway so every enabled identity of the Service
	// port is served by one data plane.
	gwName := nftableLbSvcGateway(cachedSvc)
	var eip *kubeovnv1.IptablesEIP
	var eipName, eipIP string
	if serveEIP {
		eipName = cachedSvc.Annotations[util.EipAnnotation]
		eip, err = c.GetEip(eipName)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return c.cleanupNftableLbService(cachedSvc, namespace, name)
			}
			// EIP not ready yet: requeue and retry once it has an IPv4 address.
			klog.Errorf("nftable lb service %s references eip %s which is not ready: %v", key, eipName, err)
			return err
		}
		// share DNAT is implemented with `ip daddr`/`ip saddr` and only supports IPv4
		if util.CheckProtocol(eip.Status.IP) != kubeovnv1.ProtocolIPv4 {
			klog.Errorf("nftable lb service %s references eip %s without an IPv4 address, skipping (share DNAT is IPv4 only)", key, eipName)
			return c.cleanupNftableLbService(cachedSvc, namespace, name)
		}
		if eip.Spec.NatGwDp != "" && eip.Spec.NatGwDp != gwName {
			klog.Errorf("nftable lb service %s: eip %s belongs to nat gw %s, not to the annotated gw %s", key, eipName, eip.Spec.NatGwDp, gwName)
			return c.cleanupNftableLbService(cachedSvc, namespace, name)
		}
		if !eip.DeletionTimestamp.IsZero() {
			// The EIP is going away. The Service is the only writer of the share records that
			// reference it, so it has to release them here: the EIP's finalizer waits for those
			// rules to disappear, and a record nobody else deletes would wait forever.
			klog.Infof("nftable lb service %s: eip %s is terminating, releasing its share records", key, eipName)
			return c.cleanupNftableLbService(cachedSvc, namespace, name)
		}
		if !eip.Status.Ready {
			// The DNAT worker gate (getBindableEip) used to hold a rule until its EIP was
			// programmed on the gateway; the Service writes share DNAT directly now and must wait
			// itself: an unready EIP is not live on the gateway (e.g. while a gateway replacement
			// is being redone), so claiming records with Ready status and publishing the ingress
			// IP would announce an address that drops traffic. Wait without touching an already
			// established binding; the EIP's usable transition re-enqueues this Service, and the
			// delayed retry is the fallback.
			klog.Infof("nftable lb service %s: eip %s is not ready yet, retrying later", key, eipName)
			if c.addOrUpdateGwNftableLbSvcQueue != nil {
				c.addOrUpdateGwNftableLbSvcQueue.AddAfter(key, 2*time.Second)
			}
			return nil
		}
		eipIP = eip.Status.IP
	}

	// The gateway lives in its VpcNatGateway's VPC and can only DNAT to backends reachable
	// there, so backend IPs are resolved against that VPC (see nftableLbBackendResolver).
	natGw, err := c.vpcNatGatewayLister.Get(gwName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return c.cleanupNftableLbService(cachedSvc, namespace, name)
		}
		klog.Errorf("nftable lb service %s: failed to get nat gateway %s: %v", key, gwName, err)
		return err
	}
	if !natGw.DeletionTimestamp.IsZero() {
		return c.cleanupNftableLbService(cachedSvc, namespace, name)
	}

	// The internal VIP is what a ClusterIP Service is served through, and what a LoadBalancer
	// Service is served through from inside its VPC. An IPv6-only ClusterIP has no share DNAT
	// data plane (it is IPv4 only), so the Service keeps being load balanced by OVN.
	if nftableLbSvcClusterIP(cachedSvc) == "" && eipIP == "" {
		klog.Errorf("nftable lb service %s has no IPv4 clusterIP and no eip address, skipping (share DNAT is IPv4 only)", key)
		return c.cleanupNftableLbService(cachedSvc, namespace, name)
	}

	// Claim the Service before touching the gateway. Stop this pass after the write so subsequent
	// status and record updates use the fresh ResourceVersion delivered by the informer.
	added, err := c.ensureServiceControllerFinalizer(cachedSvc)
	if err != nil {
		return err
	}
	if added {
		if c.addOrUpdateGwNftableLbSvcQueue != nil {
			c.addOrUpdateGwNftableLbSvcQueue.AddAfter(key, time.Second)
		}
		return nil
	}

	// Do not publish accounting records or ingress status before a gateway instance exists. A
	// replacement instance starts empty; the gateway redo path enqueues this Service again, and
	// this delayed retry also covers a temporarily unavailable instance.
	gone, err := c.natGwDataPlaneGone(gwName)
	if err != nil {
		return err
	}
	if gone {
		if c.addOrUpdateGwNftableLbSvcQueue != nil {
			c.addOrUpdateGwNftableLbSvcQueue.AddAfter(key, 2*time.Second)
		}
		return nil
	}

	endpointSlices, err := c.endpointSlicesLister.EndpointSlices(namespace).List(labels.Set{discoveryv1.LabelServiceName: name}.AsSelector())
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return c.cleanupNftableLbService(cachedSvc, namespace, name)
		}
		klog.Error(err)
		return err
	}

	desired := buildDesiredNftableLbDnatRulesForIdentities(cachedSvc, eipName, gwName, serveEIP, serveClusterIP, endpointSlices, c.nftableLbBackendResolver(cachedSvc, natGw.Spec.Vpc))
	for _, record := range desired {
		record.Labels[util.VpcDnatEPortLabel] = record.Spec.ExternalPort
		if eip != nil {
			record.Labels[util.EipV4IpLabel] = eipIP
			record.Labels[util.EipUIDLabel] = string(eip.UID)
		}
	}

	// A share DNAT identity (EIP, external port and protocol) has one Service writer. When
	// several Services declare it, the deterministic Service winner keeps it. The loser must
	// also never delete what it yielded: the winner owns that nft map and hairpin now (its add
	// rewrites the complete backend set), and no event would wake it to repair the damage.
	var yielded map[string]struct{}
	if serveEIP {
		yielded, err = c.resolveNftableLbConflicts(cachedSvc, key, eipIP, desired)
		if err != nil {
			return err
		}
	}
	// The Service is the authoritative writer; record events never trigger NAT changes. Write the
	// claim before the rule it covers (CODE_STYLE.md): a rejected or half-finished record sync
	// must never leave programmed share DNAT state that no record accounts for, because cleanup
	// would find no ledger for it and release the Service finalizer anyway. A failed claim aborts
	// the pass before the gateway is touched; a failed program keeps the claims for the retry.
	existing, err := c.existingNftableLbSvcRules(namespace, name)
	if err != nil {
		return err
	}
	// Records deleted out of band before this pass leave no ledger for an identity the Service
	// trimmed since (its EIP annotation is gone, or a gate closed): the restored records would
	// only describe the identities still served, and the old nft map and hairpin would be
	// deleted by nobody. Recover the evidence first so any later failure replays it.
	evidence, clusterEvidence, err := c.ensureNftableLbTrimmedIdentityRecords(cachedSvc, existing, serveEIP, serveClusterIP)
	if err != nil {
		return err
	}
	ledger := append(slices.Clone(existing), evidence...)
	ledger = append(ledger, clusterEvidence...)
	if !serveEIP && slices.ContainsFunc(ledger, func(rule *kubeovnv1.IptablesDnatRule) bool {
		return rule.Spec.EIP != "" || rule.Labels[util.EipV4IpLabel] != ""
	}) {
		// The EIP identity is leaving: withdraw the address from the Service first, while the
		// records still document the ownership. Clearing after the claim stripped the EIP would
		// be unrepetable: a failed status write is retried with records that no longer carry the
		// evidence, and the stale ingress IP would stay forever.
		if err := c.clearNftableLbSvcIngressIP(cachedSvc); err != nil {
			return err
		}
	}
	live, retiring, err := c.claimNftableLbRecords(cachedSvc, desired, existing)
	if err != nil {
		return err
	}
	if err = c.programNftableLbServiceDirect(cachedSvc, natGw.Name, eipIP, desired, ledger); err != nil {
		return err
	}
	if err = c.settleNftableLbRecords(eipIP, gwName, desired, append(slices.Clone(existing), evidence...), live, retiring); err != nil {
		return err
	}

	// The Service already programmed the gateway, so record status is not a handoff condition.
	if len(yielded) > 0 {
		return c.clearNftableLbSvcIngressIP(cachedSvc)
	}
	if !serveEIP {
		// The EIP identity is gone from the ledger; the ingress was withdrawn before the claim.
		return nil
	}
	if err = c.ensureNftableLbSvcIngressIP(cachedSvc, eipIP); err != nil {
		klog.Errorf("failed to set ingress ip for nftable lb service %s: %v", key, err)
		return err
	}
	if eipName != "" {
		c.resetIptablesEipQueue.AddAfter(eipName, 3*time.Second)
	}
	return nil
}

// nftableLbIdentity is one share DNAT map: one VIP, external port and protocol, with all
// backends that the Service wants behind it. It is an internal desired-state value, not a CRD.
type nftableLbIdentity struct {
	vip             string
	externalPort    string
	protocol        string
	backends        []string
	affinity        string
	affinityTimeout int32
}

// nftableLbProgram is the executable desired state of one identity. Keeping construction pure
// makes the Service writer testable without a Kubernetes pod-exec server.
type nftableLbProgram struct {
	identity    *nftableLbIdentity
	hairpinRule string
}

// buildNftableLbIdentities groups Service accounting records into the identities the Service
// controller will eventually program. A LoadBalancer record with both EIP and ClusterIP produces
// two identities with the same backend set; a ClusterIP-only record produces one.
func buildNftableLbIdentities(records map[string]*kubeovnv1.IptablesDnatRule, eipIP string) map[string]*nftableLbIdentity {
	identities := make(map[string]*nftableLbIdentity)
	add := func(vip string, rule *kubeovnv1.IptablesDnatRule) {
		if vip == "" {
			return
		}
		key := vip + "/" + rule.Spec.ExternalPort + "/" + strings.ToLower(rule.Spec.Protocol)
		id := identities[key]
		if id == nil {
			id = &nftableLbIdentity{
				vip:             vip,
				externalPort:    rule.Spec.ExternalPort,
				protocol:        strings.ToLower(rule.Spec.Protocol),
				affinity:        rule.Spec.SessionAffinity,
				affinityTimeout: rule.Spec.SessionAffinityTimeoutSeconds,
			}
			identities[key] = id
		}
		id.backends = append(id.backends, fmt.Sprintf("%s:%s", rule.Spec.InternalIP, rule.Spec.InternalPort))
	}
	for _, rule := range records {
		if rule.Spec.EIP != "" {
			add(eipIP, rule)
		}
		add(rule.Spec.ClusterIP, rule)
	}
	for _, id := range identities {
		id.backends = dedupSortedBackends(id.backends)
	}
	return identities
}

// programNftableLbServiceIdentities programs every desired identity on every Pod and removes the
// identities the Service no longer wants, with one exec per operation and Pod.
//
// It takes no view of the accounting records on purpose: the records describe the Service's
// desired state, not the state of the current gateway instance. Gateway-instance replacement
// re-enqueues the Service because the new instance starts empty, and an EIP reallocation changes
// the programmed VIP while the records still carry the old one, so a record-based skip would leave
// the fresh instance without rules and blackhole the Service.
// protected holds the identities another Service currently claims: they are stale for this
// Service but owned elsewhere, so their deletion is suppressed.
func (c *Controller) programNftableLbServiceIdentities(pods []*v1.Pod, programs []nftableLbProgram,
	existing []*kubeovnv1.IptablesDnatRule, protected map[string]struct{},
) error {
	wanted := make(map[string]struct{}, len(programs))
	addRules := make([]string, 0, len(programs))
	hairpinAddRules := make([]string, 0, len(programs))
	for _, program := range programs {
		identity := program.identity
		wanted[identity.vip+"/"+identity.externalPort+"/"+identity.protocol] = struct{}{}
		rule, err := nftDnatMapAddRule(identity.protocol, identity.vip, identity.externalPort,
			identity.backends, identity.affinity, identity.affinityTimeout)
		if err != nil {
			return fmt.Errorf("failed to build share dnat identity %s:%s/%s: %w",
				identity.vip, identity.externalPort, identity.protocol, err)
		}
		addRules = append(addRules, rule)
		hairpinAddRules = append(hairpinAddRules, program.hairpinRule)
	}

	// One exec per Pod for the whole Service: the gateway script accepts several rules per
	// invocation, so a multi-port Service no longer costs one pod-exec per identity and field.
	if len(addRules) != 0 {
		if err := c.execNatGwRulesInPods(pods, natGwNftDnatMapAdd, addRules); err != nil {
			return err
		}
	}
	if len(hairpinAddRules) != 0 {
		if err := c.execNatGwRulesInPods(pods, natGwVipHairpinAdd, hairpinAddRules); err != nil {
			return err
		}
	}

	delRules := make([]string, 0, len(existing))
	hairpinDelRules := make([]string, 0, len(existing))
	for key, identity := range nftableLbExistingIdentities(existing) {
		if _, ok := wanted[key]; ok {
			continue
		}
		if _, ok := protected[key]; ok {
			// The Service lost this EIP identity to another Service winner, whose claim proves it
			// owns the nft map and hairpin now: deleting here would tear down what the winner just
			// programmed, and no event tells the winner to rebuild it. The loser's records already
			// dropped the identity, so nothing else references it from this Service.
			continue
		}
		delRules = append(delRules, nftDnatMapDelRule(identity.protocol, identity.vip, identity.externalPort))
		hairpinDelRules = append(hairpinDelRules, fmt.Sprintf("%s,%s,%s", identity.vip, identity.externalPort, identity.protocol))
	}
	if len(delRules) != 0 {
		if err := c.execNatGwRulesInPods(pods, natGwNftDnatMapDel, delRules); err != nil {
			return err
		}
	}
	if len(hairpinDelRules) != 0 {
		if err := c.execNatGwRulesInPods(pods, natGwVipHairpinDel, hairpinDelRules); err != nil {
			return err
		}
	}
	return nil
}

// programNftableLbServiceDirect writes the complete desired share-DNAT identities for one
// Service. The Service and EndpointSlices are the source of truth; the accounting records are
// claimed before this runs. ledger holds this Service's pre-claim records plus any recovered
// evidence of trimmed identities; identities the ledger no longer wants are deleted - except
// the ones another Service claims: a claim passes before its data plane does, so a live record
// can mean programmed state whose owner gets no repair event from this pass. An identity
// claimed by nobody (e.g. yielded to a backend-less winner) is this Service's stale wiring to
// tear down itself.
func (c *Controller) programNftableLbServiceDirect(svc *v1.Service, gateway, eipIP string,
	desired map[string]*kubeovnv1.IptablesDnatRule, ledger []*kubeovnv1.IptablesDnatRule,
) error {
	pods, err := c.getNatGwPods(gateway, c.natGwNamespaceByName(gateway), false)
	if err != nil {
		return err
	}

	protected, err := c.nftableLbClaimedIdentities(gateway, svc.Namespace+"/"+svc.Name)
	if err != nil {
		return err
	}

	programs := buildNftableLbPrograms(desired, eipIP)
	if err = c.programNftableLbServiceIdentities(pods, programs, ledger, protected); err != nil {
		return fmt.Errorf("failed to program share dnat identities of service %s/%s: %w", svc.Namespace, svc.Name, err)
	}

	owner := svc.Namespace + "/" + svc.Name
	desiredVipRoutes, clusterIPs, err := c.desiredNatGwVipStateForService(gateway, nil, owner, desired, eipIP)
	if err != nil {
		return err
	}
	if err = c.execNatGwRulesInPods(pods, natGwVipAddrSync, clusterIPs); err != nil {
		return fmt.Errorf("failed to sync cluster VIPs of service %s/%s: %w", svc.Namespace, svc.Name, err)
	}
	// The state was computed above for the address sync; reconcile the routes with it instead of
	// scanning the DNAT lister a second time.
	if err = c.applyNatGwVipStateForService(gateway, desiredVipRoutes); err != nil {
		return fmt.Errorf("failed to sync VIP routes of service %s/%s: %w", svc.Namespace, svc.Name, err)
	}
	return nil
}

func (c *Controller) existingNftableLbSvcRules(namespace, name string) ([]*kubeovnv1.IptablesDnatRule, error) {
	return c.iptablesDnatRulesLister.List(labels.SelectorFromSet(labels.Set{
		util.NftableLbSvcNsLabel: namespace, util.NftableLbSvcNameLabel: name,
	}))
}

func nftableLbExistingIdentities(records []*kubeovnv1.IptablesDnatRule) map[string]*nftableLbIdentity {
	identities := make(map[string]*nftableLbIdentity)
	for _, record := range records {
		protocol := strings.ToLower(record.Spec.Protocol)
		add := func(vip string) {
			if vip == "" {
				return
			}
			key := vip + "/" + record.Spec.ExternalPort + "/" + protocol
			identities[key] = &nftableLbIdentity{vip: vip, externalPort: record.Spec.ExternalPort, protocol: protocol}
		}
		vip := record.Status.V4ip
		if vip == "" {
			vip = record.Labels[util.EipV4IpLabel]
		}
		add(vip)
		if record.Spec.ClusterIP != vip {
			add(record.Spec.ClusterIP)
		}
	}
	return identities
}

func buildNftableLbPrograms(records map[string]*kubeovnv1.IptablesDnatRule, eipIP string) []nftableLbProgram {
	identities := buildNftableLbIdentities(records, eipIP)
	keys := make([]string, 0, len(identities))
	for key := range identities {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	programs := make([]nftableLbProgram, 0, len(keys))
	for _, key := range keys {
		id := identities[key]
		programs = append(programs, nftableLbProgram{
			identity:    id,
			hairpinRule: fmt.Sprintf("%s,%s,%s", id.vip, id.externalPort, id.protocol),
		})
	}
	return programs
}

// claimNftableLbRecords persists the accounting records of the desired identities before the
// gateway sees the rules they cover: a share DNAT rule must never exist without a record that
// accounts for it (and carries the EIP's UID claim). Records are written without status; Ready
// is published by settleNftableLbRecords once the data plane converged. It returns the live
// record object for each desired name a status can be published for, plus the names of the
// retiring snapshots this pass created.
//
// An in-place update that drops an identity (the EIP annotation removed, the ClusterIP gate
// closed, a conflicting EIP identity yielded) would erase the only ledger the gateway cleanup
// reads and release the old EIP's UID claim before the old nft map and hairpin are actually
// removed: a failed gateway pass would then retry with records that no longer know what to
// delete. Before such an update the record is therefore snapshotted under a derived name
// (acquire the new claim, keep the old one): the snapshot preserves the old identities and UID
// claim for retried passes, and settleNftableLbRecords retires it only after the data plane
// converged - the same claim-outlives-rule order the stale-record path already follows.
func (c *Controller) claimNftableLbRecords(svc *v1.Service,
	desired map[string]*kubeovnv1.IptablesDnatRule, existing []*kubeovnv1.IptablesDnatRule,
) (map[string]*kubeovnv1.IptablesDnatRule, []string, error) {
	client := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules()
	existingByName := make(map[string]*kubeovnv1.IptablesDnatRule, len(existing))
	for _, current := range existing {
		existingByName[current.Name] = current
	}
	live := make(map[string]*kubeovnv1.IptablesDnatRule, len(desired))
	var retiring []string
	for name, want := range desired {
		current := existingByName[name]
		switch {
		case current == nil:
			created, err := client.Create(context.Background(), want, metav1.CreateOptions{})
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create nftable LB record %s for Service %s/%s: %w", name, svc.Namespace, svc.Name, err)
			}
			live[name] = created
		case !current.DeletionTimestamp.IsZero():
			// The user is removing an accounting object. It has no data-plane semantics; keep its
			// identity programmed, do not recreate the record in this pass. A legacy record also
			// carries the pre-upgrade controller finalizer that nobody clears any more: migrate it
			// off so the record can finish terminating instead of blocking its recreate forever.
			if err := c.releaseNftableLbRecordFinalizers(current); err != nil {
				return nil, nil, err
			}
			continue
		default:
			specChanged := !nftableLbDnatSpecEqual(&current.Spec, &want.Spec)
			claimChanged := current.Labels[util.EipV4IpLabel] != want.Labels[util.EipV4IpLabel] ||
				current.Labels[util.EipUIDLabel] != want.Labels[util.EipUIDLabel]
			if specChanged || claimChanged || !maps.Equal(current.Labels, want.Labels) || len(current.Finalizers) != 0 {
				if specChanged || claimChanged {
					// The update drops identity evidence or an EIP claim from the ledger: keep a
					// snapshot of the retired content until the data plane caught up.
					snapshotName, err := c.snapshotNftableLbRecord(svc, current)
					if err != nil {
						return nil, nil, err
					}
					retiring = append(retiring, snapshotName)
				}
				updated := current.DeepCopy()
				updated.Spec = want.Spec
				updated.Labels = maps.Clone(want.Labels)
				updated.Finalizers = nil
				var err error
				current, err = client.Update(context.Background(), updated, metav1.UpdateOptions{})
				if err != nil {
					return nil, nil, fmt.Errorf("failed to update nftable LB record %s: %w", current.Name, err)
				}
			}
			live[name] = current
		}
	}
	return live, retiring, nil
}

// nftableLbRetireRuleName derives the deterministic name of a record's retiring snapshot. The
// content hash distinguishes successive retired generations of one record (its UID never
// changes), and makes the create idempotent when a pass crashed between snapshot and update.
func nftableLbRetireRuleName(record *kubeovnv1.IptablesDnatRule) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s/%s/%#v/%s/%s", record.Name, record.UID,
		record.Spec, record.Labels[util.EipV4IpLabel], record.Labels[util.EipUIDLabel])))
	suffix := hex.EncodeToString(sum[:])[:8]
	base := record.Name
	if len(base) > 63-len("-r")-len(suffix) {
		base = base[:63-len("-r")-len(suffix)]
	}
	return base + "-r" + suffix
}

// snapshotNftableLbRecord persists a retiring snapshot of a record whose update drops identity
// or claim evidence. The snapshot is a plain record: no finalizers to migrate, and no status
// (the identity derivation of a snapshot falls back to its stored spec and EIP labels). It
// carries the Service ownership labels, so it joins every later existing/ledger listing of this
// Service until settle deletes it. An AlreadyExists create means an earlier pass crashed after
// snapshotting the same content, which deterministically maps to the same name.
func (c *Controller) snapshotNftableLbRecord(svc *v1.Service, current *kubeovnv1.IptablesDnatRule) (string, error) {
	snapshot := &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:   nftableLbRetireRuleName(current),
			Labels: maps.Clone(current.Labels),
		},
		Spec: *current.Spec.DeepCopy(),
	}
	_, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Create(context.Background(), snapshot, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", fmt.Errorf("failed to snapshot retiring nftable LB record %s for Service %s/%s: %w", current.Name, svc.Namespace, svc.Name, err)
	}
	return snapshot.Name, nil
}

// releaseNftableLbRecordFinalizers removes the controller finalizers from a Service accounting
// record. Pre-upgrade share rules carry them, and every share DNAT worker that used to clear
// them early-returns for these records now: once the Service ran the data-plane cleanup for the
// record's identities, nobody else will free the record, and its stuck claim would hold its EIP
// forever.
func (c *Controller) releaseNftableLbRecordFinalizers(record *kubeovnv1.IptablesDnatRule) error {
	if !slices.Contains(record.Finalizers, util.DepreciatedFinalizerName) &&
		!slices.Contains(record.Finalizers, util.KubeOVNControllerFinalizer) {
		return nil
	}
	updated := record.DeepCopy()
	updated.Finalizers = slices.DeleteFunc(updated.Finalizers, func(finalizer string) bool {
		return finalizer == util.DepreciatedFinalizerName || finalizer == util.KubeOVNControllerFinalizer
	})
	patch, err := util.GenerateMergePatchPayload(record, updated)
	if err != nil {
		return err
	}
	_, err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Patch(context.Background(), record.Name,
		types.MergePatchType, patch, metav1.PatchOptions{}, "")
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to migrate legacy finalizers of nftable LB record %s: %w", record.Name, err)
	}
	return nil
}

// settleNftableLbRecords runs once the gateway holds the desired data plane: it publishes Ready
// status for the claimed records, then retires the records the Service no longer wants and the
// retiring snapshots claimNftableLbRecords took of identity-changing updates. A stale record or
// snapshot is deleted only here so its claim outlives the rule it covered.
func (c *Controller) settleNftableLbRecords(eipIP, gateway string,
	desired map[string]*kubeovnv1.IptablesDnatRule, existing []*kubeovnv1.IptablesDnatRule, live map[string]*kubeovnv1.IptablesDnatRule,
	retiring []string,
) error {
	client := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules()
	for _, current := range existing {
		if _, ok := desired[current.Name]; ok {
			continue
		}
		if err := client.Delete(context.Background(), current.Name, metav1.DeleteOptions{}); err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete stale nftable LB record %s: %w", current.Name, err)
		}
		// The record's identities were removed in the program phase above, so the legacy
		// controller finalizer is migrated off: nobody else clears it, and the stuck record
		// would hold its EIP claim forever.
		if err := c.releaseNftableLbRecordFinalizers(current); err != nil {
			return err
		}
	}
	for _, name := range retiring {
		// Snapshots carry no finalizers: their old EIP UID claim is released by this delete,
		// which runs only after the data-plane cleanup of the identities they preserved.
		if err := client.Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete retiring nftable LB record snapshot %s: %w", name, err)
		}
	}
	for name := range desired {
		record, ok := live[name]
		if !ok {
			// terminating record: its status is published by a later pass once it was recreated.
			continue
		}
		if err := c.setNftableLbRecordStatus(record, eipIP, gateway); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) setNftableLbRecordStatus(record *kubeovnv1.IptablesDnatRule, eipIP, gateway string) error {
	vip := record.Spec.ClusterIP
	if record.Spec.EIP != "" {
		vip = eipIP
	}
	updated := record.DeepCopy()
	updated.Status.Ready = true
	updated.Status.V4ip = vip
	updated.Status.NatGwDp = gateway
	updated.Status.Protocol = record.Spec.Protocol
	updated.Status.ExternalPort = record.Spec.ExternalPort
	updated.Status.InternalIP = record.Spec.InternalIP
	updated.Status.InternalPort = record.Spec.InternalPort
	_, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().UpdateStatus(context.Background(), updated, metav1.UpdateOptions{})
	if k8serrors.IsNotFound(err) {
		return nil
	}
	return err
}

func (c *Controller) ensureServiceControllerFinalizer(svc *v1.Service) (bool, error) {
	if slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer) {
		return false, nil
	}
	updated := svc.DeepCopy()
	updated.Finalizers = append(updated.Finalizers, util.KubeOVNControllerFinalizer)
	if _, err := c.config.KubeClient.CoreV1().Services(svc.Namespace).Update(context.Background(), updated, metav1.UpdateOptions{}); err != nil {
		return false, err
	}
	return true, nil
}

// nftableLbSvcClusterIP returns the IPv4 ClusterIP that the gateway programs as a second share
// DNAT identity next to the Service's EIP, or "" when the Service has none: a headless Service,
// a ClusterIP that is not allocated yet, or an IPv6-only Service (share DNAT is IPv4 only).
func nftableLbSvcClusterIP(svc *v1.Service) string {
	return firstIPv4(util.ServiceClusterIPs(*svc))
}

// nftableLbDnatSpecEqual compares the controller-managed fields of two share DNAT specs.
// It intentionally ignores server-defaulted or unrelated fields so only meaningful drift
// (identity, backend, or session-affinity changes) triggers a rule recreate.
func nftableLbDnatSpecEqual(a, b *kubeovnv1.IptablesDnatRuleSpec) bool {
	return a.EIP == b.EIP &&
		a.ClusterIP == b.ClusterIP &&
		a.ExternalPort == b.ExternalPort &&
		a.Protocol == b.Protocol &&
		a.InternalIP == b.InternalIP &&
		a.InternalPort == b.InternalPort &&
		a.Type == b.Type &&
		a.SessionAffinity == b.SessionAffinity &&
		a.SessionAffinityTimeoutSeconds == b.SessionAffinityTimeoutSeconds
}

// ensureNftableLbSvcIngressIP sets the Service's status.loadBalancer.ingress to the EIP's
// IPv4 address so the LoadBalancer reports the externally reachable IP instead of staying
// <pending>. It is a no-op when the ingress list already advertises exactly that IP.
func (c *Controller) ensureNftableLbSvcIngressIP(svc *v1.Service, ip string) error {
	if len(svc.Status.LoadBalancer.Ingress) != 1 || svc.Status.LoadBalancer.Ingress[0].IP != ip {
		updated := svc.DeepCopy()
		updated.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{{IP: ip}}
		var err error
		if svc, err = c.config.KubeClient.CoreV1().Services(svc.Namespace).UpdateStatus(context.Background(), updated, metav1.UpdateOptions{}); err != nil {
			klog.Errorf("failed to update status of nftable lb service %s/%s: %v", svc.Namespace, svc.Name, err)
			return err
		}
		klog.Infof("set nftable lb service %s/%s ingress ip to %s", svc.Namespace, svc.Name, ip)
	}
	return nil
}

// clearNftableLbSvcIngressIP removes a LoadBalancer ingress IP that this mode previously
// published (so EXTERNAL-IP returns to <pending> once the Service leaves nftable-lb-svc
// mode). It is a no-op when there is nothing to clear or the Service is being deleted.
func (c *Controller) clearNftableLbSvcIngressIP(svc *v1.Service) error {
	if svc == nil || !svc.DeletionTimestamp.IsZero() {
		return nil
	}
	if !slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer) || len(svc.Status.LoadBalancer.Ingress) == 0 {
		return nil
	}
	updated := svc.DeepCopy()
	updated.Status.LoadBalancer.Ingress = nil
	_, err := c.config.KubeClient.CoreV1().Services(svc.Namespace).UpdateStatus(context.Background(), updated, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("failed to clear ingress ip of nftable lb service %s/%s: %v", svc.Namespace, svc.Name, err)
		return err
	}
	klog.Infof("cleared nftable lb service %s/%s ingress ip", svc.Namespace, svc.Name)
	return nil
}

// nftableLbClaimedIdentities returns the share DNAT identities of one gateway that the
// accounting records of Services other than selfKey claim. The write order (claim, rule, Ready)
// makes a live record the exact handover witness: the successor's claim precedes its data plane,
// so an identity another Service's record carries may already be programmed for that successor.
// Read through the API server: a record claimed moments ago may not be in the informer cache
// yet, and deleting on that stale view would tear down the successor's identity (the same
// release-side rule getIptablesEipNatFromAPI follows for an EIP's finalizer).
func (c *Controller) nftableLbClaimedIdentities(gwName, selfKey string) (map[string]struct{}, error) {
	rules, err := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{util.VpcNatGatewayNameLabel: gwName}).String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list share records claiming identities of nat gw %s: %w", gwName, err)
	}
	var others []*kubeovnv1.IptablesDnatRule
	for i := range rules.Items {
		rule := &rules.Items[i]
		if rule.Spec.Type != kubeovnv1.DnatRuleTypeShare || !util.IsNftableLbSvcRecord(rule.Labels) ||
			util.NftableLbSvcOwnerKey(rule.Labels) == selfKey {
			continue
		}
		others = append(others, rule)
	}
	claimed := make(map[string]struct{})
	for key := range nftableLbExistingIdentities(others) {
		claimed[key] = struct{}{}
	}
	return claimed, nil
}

// nftableLbSvcTeardownLedger synthesizes the ledger entries a Service's own intent still
// witnesses when its accounting records are gone: one per served port, carrying the EIP address
// it published and its ClusterIP. Only identity fields are filled - deletion is keyed by
// vip/port/protocol, so backends are irrelevant. The EIP address comes from the EIP object, or
// from the published ingress IP: this controller only ever sets it to the EIP's address, which
// makes it a durable witness once the EIP itself is gone.
func (c *Controller) nftableLbSvcTeardownLedger(svc *v1.Service) []*kubeovnv1.IptablesDnatRule {
	eipIP := ""
	if eipName := svc.Annotations[util.EipAnnotation]; eipName != "" {
		if eip, err := c.iptablesEipsLister.Get(eipName); err != nil {
			klog.Warningf("nftable lb service %s/%s teardown: failed to resolve eip %s: %v", svc.Namespace, svc.Name, eipName, err)
		} else {
			eipIP = eip.Status.IP
		}
	}
	if eipIP == "" && len(svc.Status.LoadBalancer.Ingress) == 1 {
		eipIP = svc.Status.LoadBalancer.Ingress[0].IP
	}
	clusterIP := nftableLbSvcClusterIP(svc)
	var ledger []*kubeovnv1.IptablesDnatRule
	for _, port := range svc.Spec.Ports {
		protocol := strings.ToLower(string(port.Protocol))
		if protocol != "tcp" && protocol != "udp" {
			// Matches what the generation programs (share DNAT supports tcp/udp only).
			continue
		}
		rule := &kubeovnv1.IptablesDnatRule{
			Spec: kubeovnv1.IptablesDnatRuleSpec{
				ClusterIP:    clusterIP,
				ExternalPort: strconv.Itoa(int(port.Port)),
				Protocol:     protocol,
			},
		}
		if eipIP != "" {
			rule.Labels = map[string]string{util.EipV4IpLabel: eipIP}
		}
		ledger = append(ledger, rule)
	}
	return ledger
}

// ensureNftableLbTrimmedIdentityRecords recreates the ledger of identities the Service trimmed
// while its accounting records were missing. Records deleted out of band erase exactly the
// evidence a trim reconcile needs to remove the old nft map and hairpin: the restored records
// only describe the identities still served, and the ingress-clearing guard above would never
// fire either. It returns the persisted EIP-leg tombstones (retired by settle once the data
// plane converged) and the in-memory ClusterIP-leg entries, which join only the program view.
// A leg another Service took over is protected by the claim check in the program phase.
func (c *Controller) ensureNftableLbTrimmedIdentityRecords(svc *v1.Service,
	existing []*kubeovnv1.IptablesDnatRule, serveEIP, serveClusterIP bool,
) (persisted, inMemory []*kubeovnv1.IptablesDnatRule, err error) {
	if !slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer) {
		// The finalizer claims everything the feature programmed; without it nothing exists.
		return nil, nil, nil
	}
	// Only recover identities a witness proves were wired, or every gate-flipped Service would
	// churn create/delete tombstones per reconcile (their settle deletion wakes the owner again
	// - a hot loop). The EIP leg's witness is the published ingress: this controller only ever
	// puts the EIP's address there, and the trim clears it, so a finished trim stops qualifying.
	// It is also consumed mid-pass (the ingress is withdrawn before the gateway is touched), so
	// the evidence must be persisted to survive a failed pass.
	trimmedEip := !serveEIP && len(svc.Status.LoadBalancer.Ingress) != 0 && !slices.ContainsFunc(existing, func(rule *kubeovnv1.IptablesDnatRule) bool {
		return rule.Spec.EIP != "" || rule.Labels[util.EipV4IpLabel] != ""
	})
	// The ClusterIP leg needs evidence only when the whole ledger was lost: its address is the
	// immutable spec field, so a stripped record set or the in-flight pass above keeps it
	// readable otherwise. The witness is never consumed, so the evidence stays in memory: no
	// records to churn, and a failed pass rederives it identically.
	lostClusterLedger := !serveClusterIP && len(existing) == 0 && nftableLbSvcClusterIP(svc) != ""
	if !trimmedEip && !lostClusterLedger {
		return nil, nil, nil
	}

	client := c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules()
	for _, intent := range c.nftableLbSvcTeardownLedger(svc) {
		eipIP := intent.Labels[util.EipV4IpLabel]
		if trimmedEip && eipIP != "" {
			record := &kubeovnv1.IptablesDnatRule{
				ObjectMeta: metav1.ObjectMeta{
					Name: nftableLbDnatRuleName(svc.Namespace, svc.Name, intent.Spec.Protocol, intent.Spec.ExternalPort, "trimmed", "evidence"),
					Labels: map[string]string{
						util.NftableLbSvcNsLabel:     svc.Namespace,
						util.NftableLbSvcNameLabel:   svc.Name,
						util.NftableLbSvcUIDLabel:    string(svc.UID),
						util.NftableLbSvcRecordLabel: "true",
						util.VpcNatGatewayNameLabel:  nftableLbSvcGateway(svc),
						util.EipV4IpLabel:            eipIP,
					},
				}, Spec: kubeovnv1.IptablesDnatRuleSpec{
					ClusterIP:    intent.Spec.ClusterIP,
					ExternalPort: intent.Spec.ExternalPort,
					Protocol:     intent.Spec.Protocol,
					Type:         kubeovnv1.DnatRuleTypeShare,
					// A tombstone: only the identity (vip, port, protocol) is read back; the backend
					// fields exist just to satisfy admission.
					InternalIP: "0.0.0.0", InternalPort: "1",
				},
			}
			if _, err := client.Create(context.Background(), record, metav1.CreateOptions{}); err != nil && !k8serrors.IsAlreadyExists(err) {
				return nil, nil, fmt.Errorf("failed to persist trimmed identity evidence %s for Service %s/%s: %w", record.Name, svc.Namespace, svc.Name, err)
			}
			persisted = append(persisted, record)
		}
		if lostClusterLedger && intent.Spec.ClusterIP != "" {
			inMemory = append(inMemory, &kubeovnv1.IptablesDnatRule{
				Spec: kubeovnv1.IptablesDnatRuleSpec{
					ClusterIP:    intent.Spec.ClusterIP,
					ExternalPort: intent.Spec.ExternalPort,
					Protocol:     intent.Spec.Protocol,
					Type:         kubeovnv1.DnatRuleTypeShare,
				},
			})
		}
	}
	return persisted, inMemory, nil
}

// cleanupNftableLbService removes all share DNAT rules generated for the given Service and clears
// the LoadBalancer ingress IP this controller marked as managed. Existing rules are also evidence
// of ownership when cleanup resumes after a controller restart.
// svc may be nil (the Service was already deleted), in which case only the rules are removed.
func (c *Controller) cleanupNftableLbService(svc *v1.Service, namespace, name string) error {
	rules, err := c.iptablesDnatRulesLister.List(labels.SelectorFromSet(labels.Set{
		util.NftableLbSvcNsLabel:   namespace,
		util.NftableLbSvcNameLabel: name,
	}))
	if err != nil {
		klog.Errorf("failed to list nftable lb dnat records for service %s/%s: %v", namespace, name, err)
		return err
	}

	// Records are the ledger of the identities this Service programmed. Remove those identities
	// directly; deleting a record never drives the gateway.
	byGateway := make(map[string][]*kubeovnv1.IptablesDnatRule)
	for _, rule := range rules {
		gw := rule.Labels[util.VpcNatGatewayNameLabel]
		if gw == "" {
			gw = rule.Status.NatGwDp
		}
		if gw != "" {
			byGateway[gw] = append(byGateway[gw], rule)
		}
	}
	if svc != nil && slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer) {
		// The records may be gone while their rules are not: claim-before-rule proves a record
		// precedes its rule, not that it outlives it, and a Service whose records were deleted
		// can enter teardown before the restore pass ran. The Service's own intent (its ports
		// and addresses) is the remaining witness of which identities could exist; the deletes
		// below are idempotent and the claim check still protects any identity a successor
		// took over, so widening the ledger with intent tears down exactly the leftovers.
		if gw := nftableLbSvcGateway(svc); gw != "" {
			byGateway[gw] = append(byGateway[gw], c.nftableLbSvcTeardownLedger(svc)...)
		}
	}
	for gw, ledger := range byGateway {
		gone, err := c.natGwDataPlaneGone(gw)
		if err != nil {
			return err
		}
		if gone {
			continue
		}
		pods, err := c.getNatGwPods(gw, c.natGwNamespaceByName(gw), false)
		if err != nil {
			return err
		}
		claimed, err := c.nftableLbClaimedIdentities(gw, namespace+"/"+name)
		if err != nil {
			return err
		}
		delRules := make([]string, 0, len(ledger))
		hairpinDelRules := make([]string, 0, len(ledger))
		for key, identity := range nftableLbExistingIdentities(ledger) {
			if _, ok := claimed[key]; ok {
				// Another Service's record claims this identity: it owns the nft map and hairpin
				// now (its claim precedes its programming, and its add rewrote the complete
				// backend set). Deleting here would tear down the wiring of a successor no event
				// wakes to repair - the same handover the program phase honors via yielded.
				continue
			}
			delRules = append(delRules, nftDnatMapDelRule(identity.protocol, identity.vip, identity.externalPort))
			hairpinDelRules = append(hairpinDelRules, fmt.Sprintf("%s,%s,%s", identity.vip, identity.externalPort, identity.protocol))
		}
		if len(delRules) != 0 {
			if err = c.execNatGwRulesInPods(pods, natGwNftDnatMapDel, delRules); err != nil {
				return fmt.Errorf("failed to remove identities of service %s/%s: %w", namespace, name, err)
			}
		}
		if len(hairpinDelRules) != 0 {
			if err = c.execNatGwRulesInPods(pods, natGwVipHairpinDel, hairpinDelRules); err != nil {
				return fmt.Errorf("failed to remove hairpins of service %s/%s: %w", namespace, name, err)
			}
		}
		owner := namespace + "/" + name
		_, clusterIPs, stateErr := c.desiredNatGwVipStateForService(gw, nil, owner, nil, "")
		if stateErr != nil {
			return stateErr
		}
		if err = c.execNatGwRulesInPods(pods, natGwVipAddrSync, clusterIPs); err != nil {
			return fmt.Errorf("failed to release cluster VIPs of service %s/%s: %w", namespace, name, err)
		}
		if err = c.syncNatGwVipStateForService(gw, nil, owner, nil, ""); err != nil {
			return fmt.Errorf("failed to release VIP routes of service %s/%s: %w", namespace, name, err)
		}
	}

	if svc != nil && slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer) {
		if err = c.clearNftableLbSvcIngressIP(svc); err != nil {
			return err
		}
	}

	for _, rule := range rules {
		if err = c.config.KubeOvnClient.KubeovnV1().IptablesDnatRules().Delete(context.Background(), rule.Name, metav1.DeleteOptions{}); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
			return fmt.Errorf("failed to delete nftable lb dnat record %s for service %s/%s: %w", rule.Name, namespace, name, err)
		}
		// The identities the record accounted for are removed above; migrate the legacy
		// controller finalizer off, or the record stays terminating forever.
		if err = c.releaseNftableLbRecordFinalizers(rule); err != nil {
			return err
		}
	}
	if svc == nil || !slices.Contains(svc.Finalizers, util.KubeOVNControllerFinalizer) {
		return nil
	}
	// clearNftableLbSvcIngressIP may have updated status and metadata. Read the latest object before
	// removing the finalizer so the cleanup does not fail on a stale ResourceVersion.
	latest, err := c.config.KubeClient.CoreV1().Services(svc.Namespace).Get(context.Background(), svc.Name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	latest.Finalizers = slices.DeleteFunc(slices.Clone(latest.Finalizers), func(finalizer string) bool {
		return finalizer == util.KubeOVNControllerFinalizer
	})
	_, err = c.config.KubeClient.CoreV1().Services(svc.Namespace).Update(context.Background(), latest, metav1.UpdateOptions{})
	return err
}

// resolveNftableLbConflicts resolves conflicts between Services that declare the same
// EIP:externalPort:protocol identity. Manual share DNAT is unsupported and is not considered:
// share CRs are accounting records of Services, not independent forwarding intent. It returns
// the vip-keyed identities this Service yields to another winner: the caller hands them to the
// data-plane sync, which must not delete them (the winner owns them already).
func (c *Controller) resolveNftableLbConflicts(svc *v1.Service, key, eipIP string, desired map[string]*kubeovnv1.IptablesDnatRule) (map[string]struct{}, error) {
	selfKey := svc.Namespace + "/" + svc.Name
	eipName := svc.Annotations[util.EipAnnotation]
	if eipName == "" {
		// Only an EIP:port identity can be contested: several Services can point at the same EIP,
		// while a ClusterIP is unique to its Service and its nft map is shared by nobody else.
		return nil, nil
	}

	// Identities come from Service intent, not desired backends. A Service can have no
	// ready endpoints temporarily, but it must still lose a contested identity and must
	// not publish an EIP that routes to another Service's backends.
	wanted := make(map[string]struct{})
	for _, id := range nftableLbSvcIdentities(svc, eipName) {
		wanted[id] = struct{}{}
	}

	// Competing owners are derived from the Service objects that reference the same EIP and
	// expose the same port+protocol (stable intent), not from existing rules, so the winner
	// is deterministic regardless of reconcile ordering (a loser must never program a rule
	// even briefly before the winner's rules exist). Self is always a candidate.
	owners := make(map[string]map[string]struct{}, len(wanted))
	for id := range wanted {
		owners[id] = map[string]struct{}{selfKey: {}}
	}

	// Look up only the services referencing this EIP via the informer index (O(matched)),
	// instead of scanning every Service in the cluster.
	svcObjs, err := c.svcIndexer.ByIndex(IndexGwNftableLbServiceByEip, eipName)
	if err != nil {
		klog.Errorf("failed to query services by eip for nftable lb conflict check %s: %v", key, err)
		return nil, err
	}
	for _, obj := range svcObjs {
		s, ok := obj.(*v1.Service)
		// The index only holds qualifying Services, so skip just self and Services that are
		// going away: a terminating owner must not block its successor.
		if !ok || (s.Namespace == svc.Namespace && s.Name == svc.Name) || !s.DeletionTimestamp.IsZero() {
			continue
		}
		for _, id := range nftableLbSvcIdentities(s, eipName) {
			if _, ok := owners[id]; ok {
				owners[id][s.Namespace+"/"+s.Name] = struct{}{}
			}
		}
	}

	// Resolve every Service identity even when there are no desired backend rules.
	droppedIdentities := make(map[string]string)
	for id := range wanted {
		winner := chooseNftableLbOwner(owners[id])
		if winner == selfKey {
			continue
		}
		droppedIdentities[id] = nftableLbOwnerDesc(winner)
	}

	// Drop available backend rules for identities this service does not win.
	for name, rule := range desired {
		id := nftableLbDnatIdentity(rule.Spec.EIP, rule.Spec.ExternalPort, rule.Spec.Protocol)
		if _, conflicted := droppedIdentities[id]; conflicted {
			// Only the EIP identity is contested: a record that also carries the ClusterIP
			// stays, minus the EIP fields, so the internal VIP keeps being served.
			if rule.Spec.ClusterIP == "" {
				delete(desired, name)
				continue
			}
			rule.Spec.EIP = ""
			delete(rule.Labels, util.EipV4IpLabel)
			delete(rule.Labels, util.EipUIDLabel)
		}
	}

	for id, winnerDesc := range droppedIdentities {
		klog.Warningf("nftable lb service %s yields share DNAT identity %s to %s; a given EIP:port can back only one owner", key, id, winnerDesc)
		c.recorder.Eventf(svc, v1.EventTypeWarning, "NftableLbSvcConflict",
			"share DNAT identity %s is owned by %s; this service will not program it (a given EIP:port can back only one LoadBalancer service)", id, winnerDesc)
	}

	// retry so this service can take over once the current owner releases the identity
	if len(droppedIdentities) > 0 {
		c.addOrUpdateGwNftableLbSvcQueue.AddAfter(key, 10*time.Second)
	}

	// The data plane keys identities by the EIP's address rather than its name.
	yielded := make(map[string]struct{}, len(droppedIdentities))
	for id := range droppedIdentities {
		yielded[eipIP+strings.TrimPrefix(id, eipName)] = struct{}{}
	}
	return yielded, nil
}

// nftableLbSvcIdentities returns the EIP-anchored share DNAT identities (eip/externalPort/protocol)
// a LoadBalancer Service would program, one per servicePort. Only an EIP identity can be contested
// between Services, which is why the conflict resolver uses it and nothing else.
func nftableLbSvcIdentities(svc *v1.Service, eipName string) []string {
	ids := make([]string, 0, len(svc.Spec.Ports))
	for _, port := range svc.Spec.Ports {
		protocol := strings.ToLower(string(port.Protocol))
		if protocol != "tcp" && protocol != "udp" {
			// Matches the ports the generation programs (see buildDesiredNftableLbDnatRules).
			continue
		}
		ids = append(ids, nftableLbDnatIdentity(eipName, strconv.Itoa(int(port.Port)), protocol))
	}
	return ids
}

// nftableLbDnatIdentity returns the share DNAT identity key (eip/externalPort/protocol)
// that determines which backends are aggregated into a single nft map.
func nftableLbDnatIdentity(eip, externalPort, protocol string) string {
	return eip + "/" + externalPort + "/" + strings.ToLower(protocol)
}

// chooseNftableLbOwner deterministically selects the lexicographically smallest Service key.
func chooseNftableLbOwner(owners map[string]struct{}) string {
	best := ""
	for owner := range owners {
		if best == "" || owner < best {
			best = owner
		}
	}
	return best
}

func nftableLbOwnerDesc(owner string) string {
	return "service " + owner
}

func endpointPortMatchesServicePort(port discoveryv1.EndpointPort, servicePort v1.ServicePort) bool {
	if port.Name == nil {
		if servicePort.Name != "" {
			return false
		}
	} else if *port.Name != servicePort.Name {
		return false
	}
	return port.Protocol == nil || *port.Protocol == servicePort.Protocol
}

// buildDesiredNftableLbDnatRules builds one rule per (servicePort, ready backend). backendIP resolves each ready
// endpoint to the single IPv4 the gateway must DNAT to (the NIC in the gateway's VPC; see
// nftableLbBackendResolver); an endpoint it cannot resolve is skipped. The map is keyed by
// rule name; the name deterministically encodes the full identity so that an unchanged
// backend always maps to the same rule (idempotent reconcile).
func buildDesiredNftableLbDnatRules(svc *v1.Service, eipName, gateway string, endpointSlices []*discoveryv1.EndpointSlice, backendIP func(discoveryv1.Endpoint) (string, bool)) map[string]*kubeovnv1.IptablesDnatRule {
	return buildDesiredNftableLbDnatRulesForIdentities(svc, eipName, gateway, true, true, endpointSlices, backendIP)
}

// buildDesiredNftableLbDnatRulesForIdentities derives the records of the identities the enabled
// feature gates select. The ClusterIP identity is dropped when its gate is off, the EIP identity
// when its gate is off or the Service names no EIP, and both fields stay on one record when both
// are served: everything else (nft maps, hairpin, lo addresses, VIP routes) is derived from these
// records, so the gates are applied here once.
func buildDesiredNftableLbDnatRulesForIdentities(svc *v1.Service, eipName, gateway string, serveEIP, serveClusterIP bool, endpointSlices []*discoveryv1.EndpointSlice, backendIP func(discoveryv1.Endpoint) (string, bool)) map[string]*kubeovnv1.IptablesDnatRule {
	desired := make(map[string]*kubeovnv1.IptablesDnatRule)

	// The internal VIP the gateway programs with the same backends (see file header).
	clusterIP := ""
	if serveClusterIP {
		clusterIP = nftableLbSvcClusterIP(svc)
	}
	if !serveEIP {
		eipName = ""
	}
	if eipName == "" && clusterIP == "" {
		return desired
	}

	// Translate the Service's client-IP session affinity into the share DNAT fields. All
	// backends of one identity carry the same affinity settings, matching kube-proxy where
	// affinity is a per-ServicePort property.
	sessionAffinity := kubeovnv1.DnatSessionAffinityNone
	var affinityTimeout int32
	if svc.Spec.SessionAffinity == v1.ServiceAffinityClientIP {
		sessionAffinity = kubeovnv1.DnatSessionAffinityClientIP
		if cfg := svc.Spec.SessionAffinityConfig; cfg != nil && cfg.ClientIP != nil && cfg.ClientIP.TimeoutSeconds != nil {
			affinityTimeout = *cfg.ClientIP.TimeoutSeconds
		}
	}

	for _, port := range svc.Spec.Ports {
		protocol := strings.ToLower(string(port.Protocol))
		// share DNAT only supports tcp/udp for now; see the TODO on util.ValidateProtocol for what
		// extending it takes (the nft data plane itself is protocol agnostic).
		if protocol != "tcp" && protocol != "udp" {
			klog.Warningf("skipping service %s/%s port %d: nftable lb service only supports tcp/udp for now, got %s",
				svc.Namespace, svc.Name, port.Port, port.Protocol)
			continue
		}
		externalPort := strconv.Itoa(int(port.Port))

		for _, endpointSlice := range endpointSlices {
			var targetPort int32
			for _, p := range endpointSlice.Ports {
				if endpointPortMatchesServicePort(p, port) && p.Port != nil {
					targetPort = *p.Port
					break
				}
			}
			if targetPort == 0 {
				continue
			}

			for _, endpoint := range endpointSlice.Endpoints {
				// Cluster traffic policy only (see file header): every Ready endpoint of the
				// Service is a backend, regardless of node locality. ExternalTrafficPolicy/
				// InternalTrafficPolicy=Local, topology, and terminating-endpoint fallback are
				// intentionally not honored here.
				if !endpointReady(endpoint) {
					continue
				}
				// Resolve to the IPv4 the gateway can reach in its VPC; skip a backend the
				// gateway cannot reach (resolver emits a Warning event once per pod).
				address, ok := backendIP(endpoint)
				if !ok {
					continue
				}
				internalPort := strconv.Itoa(int(targetPort))
				name := nftableLbDnatRuleName(svc.Namespace, svc.Name, protocol, externalPort, address, internalPort)
				rule := &kubeovnv1.IptablesDnatRule{
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
						Labels: map[string]string{
							util.NftableLbSvcNsLabel:     svc.Namespace,
							util.NftableLbSvcNameLabel:   svc.Name,
							util.NftableLbSvcUIDLabel:    string(svc.UID),
							util.NftableLbSvcRecordLabel: "true",
							util.VpcNatGatewayNameLabel:  gateway,
						},
					},
					Spec: kubeovnv1.IptablesDnatRuleSpec{
						// The enabled addresses of the Service port share one record: ClusterIP is
						// the internal VIP the gateway holds on lo, while EIP is the public one a
						// LoadBalancer Service publishes. Either field can be independently disabled.
						EIP:                           eipName,
						ClusterIP:                     clusterIP,
						ExternalPort:                  externalPort,
						Protocol:                      protocol,
						InternalIP:                    address,
						InternalPort:                  internalPort,
						Type:                          kubeovnv1.DnatRuleTypeShare,
						SessionAffinity:               sessionAffinity,
						SessionAffinityTimeoutSeconds: affinityTimeout,
					},
				}
				desired[name] = rule
			}
		}
	}

	return desired
}

// nftableLbNicCandidate is a resolvable kube-ovn NIC of a backend pod: its IPv4 address and
// the VPC its subnet belongs to. Used to pick the backend IP reachable from the gateway.
type nftableLbNicCandidate struct {
	ipv4 string
	vpc  string
}

// selectNftableLbBackendIPv4 returns the IPv4 of the NIC that sits in gwVpc, i.e. the address
// the vpc-nat-gw (which lives in gwVpc) can actually DNAT to. A single-NIC backend in the
// gateway VPC matches on its only NIC. When a backend has several NICs in the gateway VPC the
// lowest IPv4 is chosen deterministically: a dual-NIC backend (one default-VPC NIC for
// kube-proxy + one gateway-VPC NIC for the gateway) covers ~99% of cases and has exactly one
// match, so per-NIC selection is intentionally not exposed.
func selectNftableLbBackendIPv4(candidates []nftableLbNicCandidate, gwVpc string) (string, bool) {
	var matches []string
	for _, candidate := range candidates {
		if candidate.ipv4 != "" && candidate.vpc == gwVpc {
			matches = append(matches, candidate.ipv4)
		}
	}
	if len(matches) == 0 {
		return "", false
	}
	slices.Sort(matches)
	return matches[0], true
}

// nftableLbBackendResolver returns a resolver mapping a ready endpoint to the single IPv4 the
// gateway must DNAT to. For a multi-NIC backend it selects the NIC in the gateway's VPC
// (gwVpc); a backend that has kube-ovn NICs but none in gwVpc is unreachable from the gateway
// and is skipped (with a one-shot Warning event) to avoid a black-hole share DNAT. An endpoint
// whose Pod target no longer exists is skipped for the same reason; only endpoints without a
// Pod target (e.g. externally maintained EndpointSlices) fall back to their primary IPv4.
func (c *Controller) nftableLbBackendResolver(svc *v1.Service, gwVpc string) func(discoveryv1.Endpoint) (string, bool) {
	type resolved struct {
		ip string
		ok bool
	}
	// Cache one result per backend pod: a pod appears once per Service port, so memoizing
	// both avoids repeated lister traversals and emits at most one skip event per pod.
	cache := make(map[string]resolved)
	return func(ep discoveryv1.Endpoint) (string, bool) {
		primary := firstIPv4(ep.Addresses)
		if gwVpc == "" {
			return primary, primary != ""
		}
		pod := c.nftableLbEndpointPod(ep, svc.Namespace)
		if pod == nil {
			// An endpoint that still references a Pod which no longer exists must not fall
			// back to its recorded address: the backend is gone and DNATing to it would black
			// hole traffic (kube-ovn removes automatic EndpointSlices entries, but manually
			// maintained EndpointSlices without a Service selector are not updated for us).
			// Endpoints without a Pod target still fall back to their primary IPv4.
			if ep.TargetRef != nil && ep.TargetRef.Kind == "Pod" && ep.TargetRef.Name != "" {
				return "", false
			}
			return primary, primary != ""
		}
		podKey := pod.Namespace + "/" + pod.Name
		if r, done := cache[podKey]; done {
			return r.ip, r.ok
		}

		// default to the primary endpoint IP (best effort for non-kube-ovn pods)
		ip, ok := primary, primary != ""
		candidates := c.nftableLbPodNicCandidates(pod)
		if selected, matched := selectNftableLbBackendIPv4(candidates, gwVpc); matched {
			ip, ok = selected, true
		} else if len(candidates) > 0 {
			// has kube-ovn NICs but none in the gateway VPC: unreachable from the gateway
			ip, ok = "", false
			c.recorder.Eventf(svc, v1.EventTypeWarning, "NftableLbSvcBackendSkipped",
				"backend pod %s has no NIC in gateway VPC %q; skipping it to avoid an unreachable share DNAT", podKey, gwVpc)
		}
		cache[podKey] = resolved{ip, ok}
		return ip, ok
	}
}

// nftableLbEndpointPod returns the backend Pod referenced by an endpoint, or nil when the
// endpoint has no Pod target or the pod is not in cache.
func (c *Controller) nftableLbEndpointPod(ep discoveryv1.Endpoint, defaultNamespace string) *v1.Pod {
	if ep.TargetRef == nil || ep.TargetRef.Kind != "Pod" || ep.TargetRef.Name == "" {
		return nil
	}
	namespace := ep.TargetRef.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	pod, err := c.podsLister.Pods(namespace).Get(ep.TargetRef.Name)
	if err != nil {
		return nil
	}
	return pod
}

// nftableLbPodNicCandidates returns one candidate per resolvable kube-ovn NIC of the pod
// (primary and attached), pairing the NIC's IPv4 with the VPC of its subnet.
func (c *Controller) nftableLbPodNicCandidates(pod *v1.Pod) []nftableLbNicCandidate {
	providers, err := c.getPodProviders(pod)
	if err != nil {
		klog.Warningf("failed to get providers for backend pod %s/%s: %v", pod.Namespace, pod.Name, err)
		return nil
	}
	var candidates []nftableLbNicCandidate
	for _, provider := range providers {
		subnetName, err := getSubnetByProvider(pod, provider)
		if err != nil {
			continue
		}
		subnet, err := c.subnetsLister.Get(subnetName)
		if err != nil {
			continue
		}
		ips := strings.Split(pod.Annotations[fmt.Sprintf(util.IPAddressAnnotationTemplate, provider)], ",")
		candidates = append(candidates, nftableLbNicCandidate{
			ipv4: firstIPv4(ips),
			vpc:  subnet.Spec.Vpc,
		})
	}
	return candidates
}

// firstIPv4 returns the first IPv4 address in the list, or "" when there is none.
func firstIPv4(addresses []string) string {
	for _, address := range addresses {
		if util.CheckProtocol(address) == kubeovnv1.ProtocolIPv4 {
			return address
		}
	}
	return ""
}

// nftableLbDnatRuleName builds a deterministic, DNS-1123 compliant name that uniquely
// encodes the rule identity. The name is "lb-<sanitized svc name>-<12 hex hash>"; the
// svc name portion is truncated so the total length never exceeds 63 characters.
func nftableLbDnatRuleName(namespace, name, protocol, externalPort, backendIP, internalPort string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{namespace, name, protocol, externalPort, backendIP, internalPort}, "/")))
	hash := hex.EncodeToString(sum[:])[:12]

	const prefix = "lb-"
	// reserve room for prefix, the "-" separator and the 12-char hash
	maxNameLen := 63 - len(prefix) - 1 - len(hash)
	svcPart := name
	if len(svcPart) > maxNameLen {
		svcPart = svcPart[:maxNameLen]
	}
	return fmt.Sprintf("%s%s-%s", prefix, svcPart, hash)
}
