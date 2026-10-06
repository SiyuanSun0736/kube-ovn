package controller

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// The lanIP-as-Service-VIP feature gives every Service bound to a vpc-nat-gw (annotation
// ovn.kubernetes.io/vpc_nat_gw, LoadBalancer or ClusterIP type) one more frontend: the
// gateway's lanIP plus the Service port. It is a pure add-on next to the existing paths -
// the OVN load balancer, the EIP and the ClusterIP identities keep working unchanged - and
// it needs no routes and no loopback addresses because the lanIP is directly reachable from
// inside the VPC subnet.
//
// Data plane: nftables only. The identities live in the same kube-ovn table's service-ips
// map as the share DNAT of the EIP/ClusterIP legs (keyed ip daddr . proto . port, so they
// can never collide), plus a feature-owned postrouting SNAT chain that pins the reply to
// the gateway instance that performed the DNAT. The gateway script reconciles the whole
// partition to a complete identity set per invocation (nft-lanvip-sync), so the controller
// keeps one write path: compute the desired set from Service/EndpointSlice intent, exec
// once per gateway Pod.
//
// Port sharing: the lanIP port space is per-gateway. Two Services declaring the same
// port+protocol on one gateway MERGE their backends into a single identity (a shared
// pool), not a conflict; session-affinity disagreements resolve deterministically to
// ClientIP with the maximum timeout, and get a Warning event.
//
// Lifecycle: no CRD records and no finalizers are involved. The partition reconciliation
// is idempotent and complete-set, so it doubles as the cleanup path: with
// --enable-gw-nftable-lanip-vip=false the desired set is empty and the sync wipes whatever
// the feature had programmed. Startup informer replay and gateway-redo wake-ups provide
// the triggers; Service and EndpointSlice events keep the union current afterwards.

// natGwNftLanVipSync is the gateway-script operation that reconciles the lanIP identity
// partition of one gateway Pod to a complete set.
const natGwNftLanVipSync = "nft-lanvip-sync"

// enqueueNatGwLanVipSync queues a partition reconciliation of the named gateway.
func (c *Controller) enqueueNatGwLanVipSync(gwName string) {
	if c.natGwLanVipSyncQueue == nil || gwName == "" {
		return
	}
	klog.V(3).Infof("enqueue lanIP vip sync of nat gw %s", gwName)
	c.natGwLanVipSyncQueue.Add(gwName)
}

// enqueueNatGwLanVipForService maps a Service key to its gateway and wakes the gateway's
// partition reconciliation. Used by the EndpointSlice handlers, where only the key exists.
func (c *Controller) enqueueNatGwLanVipForService(namespace, name string) {
	if c.natGwLanVipSyncQueue == nil {
		return
	}
	svc, err := c.servicesLister.Services(namespace).Get(name)
	if err != nil {
		// A deleted Service reaches the partition sync through its own delete event (its
		// gateway annotation is read from the tombstone object there).
		return
	}
	c.enqueueNatGwLanVipSync(svc.Annotations[util.VpcNatGatewayAnnotation])
}

// lanVipIdentity accumulates the merged desired state of one (port, protocol) identity on
// the gateway's lanIP: the union of all contributing Services' backends and the resolved
// session affinity.
type lanVipIdentity struct {
	port     string
	protocol string
	backends map[string]struct{}
	clientIP bool
	timeout  int32
	// affinitySource records what the affinity was merged from, so a disagreement between
	// the contributing Services can be reported once with the actual values.
	affinitySource string
}

// lanVipAffinityMerge records one Service whose session affinity merged into a shared lanIP
// identity. The merge set is reported once per change, not on every reconcile (see
// emitNatGwLanVipAffinityMerges).
type lanVipAffinityMerge struct {
	svc     *v1.Service
	key     string // identity key "port/protocol"
	summary string // this Service's "clientIP=<bool> timeout=<s>"
	source  string // the identity's seed contributor (deterministic: Services are sorted)
	lanIP   string
}

// desiredNatGwLanVipRules computes the complete lanIP identity set of one gateway in the
// nft-dnat-map-add rule format, one rule per (servicePort, protocol): the backend union of
// every candidate Service bound to the gateway, in a deterministic order. Identities with
// no ready backends are omitted (the data plane has no empty maps; the partition GC on the
// gateway removes whatever they left behind).
func (c *Controller) desiredNatGwLanVipRules(gw *kubeovnv1.VpcNatGateway) ([]string, error) {
	lanIP := gw.Spec.LanIP

	svcObjs, err := c.svcIndexer.ByIndex(IndexGwNftableLbServiceByGateway, gw.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to list services bound to nat gw %s: %w", gw.Name, err)
	}

	// The informer index returns Services in arbitrary order, but merge order decides which
	// Service seeds an identity's affinitySource (and thus which Service the merge warning
	// names). Sort by namespace/name so attribution is deterministic, like the identity keys
	// and backends already are further down.
	svcs := make([]*v1.Service, 0, len(svcObjs))
	for _, svcObj := range svcObjs {
		if svc, ok := svcObj.(*v1.Service); ok {
			svcs = append(svcs, svc)
		}
	}
	slices.SortFunc(svcs, func(a, b *v1.Service) int {
		if c := strings.Compare(a.Namespace, b.Namespace); c != 0 {
			return c
		}
		return strings.Compare(a.Name, b.Name)
	})

	identities := make(map[string]*lanVipIdentity)
	var merges []lanVipAffinityMerge
	for _, svc := range svcs {
		if !svc.DeletionTimestamp.IsZero() {
			// A terminating Service stops contributing; its cleanup path re-syncs the partition.
			continue
		}

		endpointSlices, err := c.endpointSlicesLister.EndpointSlices(svc.Namespace).List(
			labels.Set{discoveryv1.LabelServiceName: svc.Name}.AsSelector(),
		)
		if err != nil {
			return nil, err
		}
		resolve := c.nftableLbBackendResolver(svc, gw.Spec.Vpc)

		// Session affinity is a Service-level property (like kube-proxy): it applies to every
		// port of the Service, so each identity the Service contributes to carries it; the
		// identity key below merges contributions per (port, protocol).
		svcClientIP := svc.Spec.SessionAffinity == v1.ServiceAffinityClientIP
		var svcTimeout int32
		if svcClientIP {
			if cfg := svc.Spec.SessionAffinityConfig; cfg != nil && cfg.ClientIP != nil && cfg.ClientIP.TimeoutSeconds != nil {
				svcTimeout = *cfg.ClientIP.TimeoutSeconds
			}
		}

		for _, port := range svc.Spec.Ports {
			protocol := strings.ToLower(string(port.Protocol))
			if protocol != "tcp" && protocol != "udp" {
				// Matches what the share DNAT generation programs (tcp/udp only).
				continue
			}
			key := strconv.Itoa(int(port.Port)) + "/" + protocol
			id, known := identities[key]
			if !known {
				id = &lanVipIdentity{
					port:     strconv.Itoa(int(port.Port)),
					protocol: protocol,
					backends: make(map[string]struct{}),
				}
				identities[key] = id
			}

			// Merge the session affinity deterministically: ClientIP wins, the timeout is the
			// maximum. A disagreement between contributing Services is legal (port sharing is
			// by design) but surprising, so it is reported.
			id.clientIP = id.clientIP || svcClientIP
			if svcTimeout > id.timeout {
				id.timeout = svcTimeout
			}
			summary := fmt.Sprintf("clientIP=%t timeout=%d", svcClientIP, svcTimeout)
			switch {
			case id.affinitySource == "":
				id.affinitySource = fmt.Sprintf("%s/%s (%s)", svc.Namespace, svc.Name, summary)
			case id.affinitySource != "" && !strings.HasSuffix(id.affinitySource, "("+summary+")"):
				// The log/event emission is deferred to the end of the reconcile and gated by
				// signature, so a steady disagreement does not spam every pass.
				merges = append(merges, lanVipAffinityMerge{
					svc: svc, key: key, summary: summary, source: id.affinitySource, lanIP: lanIP,
				})
			}

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
					// Cluster traffic policy only: every Ready endpoint contributes.
					if !endpointReady(endpoint) {
						continue
					}
					address, ok := resolve(endpoint)
					if !ok {
						continue
					}
					id.backends[fmt.Sprintf("%s:%d", address, targetPort)] = struct{}{}
				}
			}
		}
	}

	keys := make([]string, 0, len(identities))
	for key := range identities {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	rules := make([]string, 0, len(identities))
	for _, key := range keys {
		id := identities[key]
		backends := make([]string, 0, len(id.backends))
		for backend := range id.backends {
			backends = append(backends, backend)
		}
		if len(backends) == 0 {
			// A port with no ready backend is not programmed (nft maps cannot be empty): the
			// gateway-side complete-set sync removes the identity if it existed.
			klog.V(2).Infof("nat gw %s lanIP identity %s:%s/%s has no ready backends, skipping", gw.Name, lanIP, id.port, id.protocol)
			continue
		}
		affinity := kubeovnv1.DnatSessionAffinityNone
		if id.clientIP {
			affinity = kubeovnv1.DnatSessionAffinityClientIP
		}
		rule, err := nftDnatMapAddRule(id.protocol, lanIP, id.port, backends, affinity, id.timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to build lanIP vip identity %s:%s/%s of nat gw %s: %w", lanIP, id.port, id.protocol, gw.Name, err)
		}
		rules = append(rules, rule)
	}
	c.emitNatGwLanVipAffinityMerges(gw.Name, merges)
	return rules, nil
}

// emitNatGwLanVipAffinityMerges reports Services whose session affinity merged on a shared
// lanIP identity. desiredNatGwLanVipRules runs on every Service/EndpointSlice event of the
// gateway, so the log line and the Warning event are emitted only when the merge set
// changes; otherwise EndpointSlice churn would re-fire the same disagreement on every pass
// and spam the event recorder.
func (c *Controller) emitNatGwLanVipAffinityMerges(gwName string, merges []lanVipAffinityMerge) {
	if len(merges) == 0 {
		c.natGwLanVipMergeNotes.Delete(gwName)
		return
	}
	parts := make([]string, 0, len(merges))
	for _, m := range merges {
		parts = append(parts, fmt.Sprintf("%s|%s/%s|%s", m.key, m.svc.Namespace, m.svc.Name, m.summary))
	}
	sig := strings.Join(parts, "\n")
	if prev, ok := c.natGwLanVipMergeNotes.Load(gwName); ok && prev == sig {
		return
	}
	c.natGwLanVipMergeNotes.Store(gwName, sig)
	for _, m := range merges {
		klog.Warningf("nat gw %s lanIP identity %s: session affinity %s of service %s/%s merges into %s (ClientIP with the max timeout wins)",
			gwName, m.key, m.summary, m.svc.Namespace, m.svc.Name, m.source)
		c.recorder.Eventf(m.svc, v1.EventTypeWarning, "NatGwLanVipAffinityMerged",
			"session affinity (%s) of this port merges with %s on gateway %s lanIP %s: identity uses ClientIP with the maximum timeout",
			m.summary, m.source, gwName, m.lanIP)
	}
}

// handleSyncNatGwLanVip reconciles one gateway's lanIP identity partition. It runs whether
// or not the feature is enabled: with the flag off the desired set is empty and the sync
// wipes the leftover partition, which is what lets a flipped flag strand nothing.
func (c *Controller) handleSyncNatGwLanVip(gwName string) error {
	klog.Infof("sync lanIP vip identities of nat gw %s", gwName)
	gw, err := c.vpcNatGatewayLister.Get(gwName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// The gateway is gone and its Pods (and their rules) cascade away with it.
			c.natGwLanVipMergeNotes.Delete(gwName)
			c.natGwLanVipLastSync.Delete(gwName)
			return nil
		}
		return err
	}
	if !gw.DeletionTimestamp.IsZero() {
		return nil
	}

	var rules []string
	if c.config.EnableGwNftableLanipVip {
		// Share DNAT is IPv4 only. An empty/IPv6 lanIP means the gateway is not set up for
		// this yet (the spec is backfilled by the gateway reconcile); a gateway update
		// re-enqueues this sync, and there is no partition to wipe for an address the
		// feature never had.
		if util.CheckProtocol(gw.Spec.LanIP) != kubeovnv1.ProtocolIPv4 {
			klog.Infof("nat gw %s has no IPv4 lanIP yet, skipping lanIP vip sync", gwName)
			return nil
		}
		rules, err = c.desiredNatGwLanVipRules(gw)
		if err != nil {
			return err
		}
	}

	// Do not program into a gateway whose data plane is not there (yet). The gateway redo
	// path re-enqueues this sync once an instance is up; the delayed retry is the fallback.
	gone, err := c.natGwDataPlaneGone(gwName)
	if err != nil {
		return err
	}
	if gone {
		if c.natGwLanVipSyncQueue != nil {
			c.natGwLanVipSyncQueue.AddAfter(gwName, 2*time.Second)
		}
		return nil
	}

	pods, err := c.getNatGwPods(gwName, c.natGwNamespaceByName(gwName), false)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		if c.natGwLanVipSyncQueue != nil {
			c.natGwLanVipSyncQueue.AddAfter(gwName, 2*time.Second)
		}
		return nil
	}

	// The desired set is computed purely from informer state and the script reconciles the
	// partition to exactly that set, so while the gateway instance set (pod UIDs) and the rule
	// set are both unchanged, a re-exec would program byte-identical nft state. Skip it:
	// EndpointSlice churn drives this path heavily, and execing per event per gateway is the
	// dominant cost of the feature at scale. A restarted instance (new UID) or any intent
	// change invalidates the signature and re-programs. No periodic self-heal is attempted
	// for out-of-band nft edits inside the gateway Pod; the instance redo (UID change) or the
	// next intent change re-syncs.
	sig := natGwLanVipSyncSignature(pods, rules)
	if prev, ok := c.natGwLanVipLastSync.Load(gwName); ok && prev == sig {
		klog.V(3).Infof("nat gw %s: lanIP vip partition unchanged, skipping exec", gwName)
		return nil
	}
	if err = c.execNatGwRulesInPods(pods, natGwNftLanVipSync, rules); err != nil {
		return fmt.Errorf("failed to sync lanIP vip identities of nat gw %s: %w", gwName, err)
	}
	c.natGwLanVipLastSync.Store(gwName, sig)
	klog.Infof("nat gw %s: synced %d lanIP vip identities to %d pod(s)", gwName, len(rules), len(pods))
	return nil
}

// natGwLanVipSyncSignature identifies the exact gateway state a sync would program: which
// instances would receive it (their UIDs, so a recreated Pod can never be mistaken for the
// old one) plus the full desired rule set (already deterministically ordered).
func natGwLanVipSyncSignature(pods []*v1.Pod, rules []string) string {
	uids := make([]string, 0, len(pods))
	for _, pod := range pods {
		uids = append(uids, string(pod.UID))
	}
	slices.Sort(uids)
	return strings.Join(uids, " ") + "\x00" + strings.Join(rules, "\n")
}
