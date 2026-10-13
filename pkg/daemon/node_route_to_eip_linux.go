package daemon

// This file implements node-local EIP access for VPC NAT Gateway.
//
// Problem: When a VPC NAT Gateway pod runs on a node, the node itself cannot
// directly access the EIP addresses configured on the gateway because the
// EIP traffic goes through the external macvlan interface attached to the pod.
//
// Solution: Create a macvlan sub-interface on the node with the same master
// interface as the NAT Gateway's external network, and add host routes for
// each EIP pointing to this sub-interface. This allows the node to reach
// the EIP addresses via Layer 2 on the external network.
//
// Key design: One macvlan sub-interface per master interface (not per subnet).
// Multiple subnets may share the same NAD (Network Attachment Definition) and
// thus the same master interface. Creating one interface per master avoids
// ARP conflicts that would occur if the same IP were configured on multiple
// interfaces.
//
// Control flow:
//  1. Subnet event (subnet with NadMacvlanMasterAnnotation):
//     enqueueAddSubnet/enqueueUpdateSubnet/enqueueDeleteSubnet
//     → macvlanSubnetQueue → runMacvlanSubnetWorker()
//     → reconcileMacvlanSubnet() → createMacvlanSubInterface() / deleteMacvlanSubInterface()
//
//  2. IptablesEIP event:
//     - Add/Update: enqueueAddIptablesEip / enqueueUpdateIptablesEip → iptablesEipQueue
//       → runIptablesEipWorker() → syncIptablesEipRoute() → addEIPRoute() / deleteEIPRoute()
//     - Delete: enqueueDeleteIptablesEip → iptablesEipDeleteQueue
//       → runIptablesEipDeleteWorker() → handleDeleteIptablesEipRoute() → deleteEIPRoute()
//
//  3. NAT GW Pod event (pod with VpcNatGatewayLabel):
//     - Update: enqueueUpdatePod() → handleNatGwPodUpdate() → enqueueEipsForNatGw()
//     - Delete: enqueueDeletePod() → handleNatGwPodDelete() → enqueueEipsForNatGw()
//     → syncIptablesEipRoute() → addEIPRoute() / deleteEIPRoute()
//
// Prerequisites:
//   - Subnet controller must set NadMacvlanMasterAnnotation when provider NAD is macvlan type
//   - EnableNodeLocalAccessVpcNatGwEIP config flag must be true (default)

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// eipRouteInfo holds information needed to add or delete an EIP route.
// Supports both legacy macvlan mode and high-performance veth-pair mode.
type eipRouteInfo struct {
	eipName     string // EIP name, used for add operations to check current state
	v4ip        string // IPv4 address of the EIP
	macvlanName string // macvlan sub-interface name for the route
	mode        string // "macvlan" or "veth"
}

// parseEIPDestination parses an EIP address into a host route destination (/32 for IPv4, /128 for IPv6)
func parseEIPDestination(eip string) (*net.IPNet, error) {
	ip := net.ParseIP(eip)
	if ip == nil {
		err := fmt.Errorf("invalid EIP address: %s", eip)
		klog.Error(err)
		return nil, err
	}

	mask := net.CIDRMask(32, 32)
	if ip.To4() == nil {
		mask = net.CIDRMask(128, 128)
	}
	return &net.IPNet{IP: ip, Mask: mask}, nil
}

// ensureMacvlanSubInterfaceUp ensures the macvlan sub-interface is up
func ensureMacvlanSubInterfaceUp(link netlink.Link) error {
	if link.Attrs().OperState != netlink.OperUp {
		if err := netlink.LinkSetUp(link); err != nil {
			err = fmt.Errorf("failed to set link %s up: %w", link.Attrs().Name, err)
			klog.Error(err)
			return err
		}
	}
	return nil
}

// createMacvlanSubInterface creates a macvlan sub-interface for node local EIP access.
// The macvlan name is derived from the master interface name, ensuring only one
// sub-interface is created per master interface even if multiple subnets use it.
func (c *Controller) createMacvlanSubInterface(masterIface string) error {
	macvlanName, err := util.GenMacvlanIfaceName(masterIface)
	if err != nil {
		return fmt.Errorf("createMacvlanSubInterface: %w", err)
	}

	// Check if sub-interface already exists
	if link, err := netlink.LinkByName(macvlanName); err == nil {
		klog.V(3).Infof("macvlan sub-interface %s already exists", macvlanName)
		return ensureMacvlanSubInterfaceUp(link)
	}

	master, err := netlink.LinkByName(masterIface)
	if err != nil {
		err = fmt.Errorf("failed to get master interface %s: %w", masterIface, err)
		klog.Error(err)
		return err
	}

	macvlan := &netlink.Macvlan{
		LinkAttrs: netlink.LinkAttrs{
			Name:        macvlanName,
			ParentIndex: master.Attrs().Index,
		},
		Mode: netlink.MACVLAN_MODE_BRIDGE,
	}

	if err := netlink.LinkAdd(macvlan); err != nil {
		err = fmt.Errorf("failed to create macvlan sub-interface %s: %w", macvlanName, err)
		klog.Error(err)
		return err
	}

	if err := netlink.LinkSetUp(macvlan); err != nil {
		if delErr := netlink.LinkDel(macvlan); delErr != nil {
			klog.Warningf("failed to cleanup macvlan interface %s after setup failed: %v", macvlanName, delErr)
		}
		err = fmt.Errorf("failed to set macvlan sub-interface %s up: %w", macvlanName, err)
		klog.Error(err)
		return err
	}

	klog.Infof("created macvlan sub-interface %s with master %s", macvlanName, masterIface)
	return nil
}

// deleteMacvlanSubInterface deletes the macvlan sub-interface for the given master interface.
func deleteMacvlanSubInterface(masterIface string) error {
	macvlanName, err := util.GenMacvlanIfaceName(masterIface)
	if err != nil {
		return fmt.Errorf("deleteMacvlanSubInterface: %w", err)
	}

	link, err := netlink.LinkByName(macvlanName)
	if err != nil {
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			return nil
		}
		err = fmt.Errorf("failed to get macvlan sub-interface %s: %w", macvlanName, err)
		klog.Error(err)
		return err
	}

	if err := netlink.LinkDel(link); err != nil {
		err = fmt.Errorf("failed to delete macvlan sub-interface %s: %w", macvlanName, err)
		klog.Error(err)
		return err
	}

	klog.Infof("deleted macvlan sub-interface %s", macvlanName)
	return nil
}

// addEIPRoute adds a route for EIP via the specified macvlan sub-interface.
func addEIPRoute(eip, macvlanSubIfName string) error {
	link, err := netlink.LinkByName(macvlanSubIfName)
	if err != nil {
		err = fmt.Errorf("failed to get macvlan sub-interface %s for EIP %s: %w", macvlanSubIfName, eip, err)
		klog.Error(err)
		return err
	}

	dst, err := parseEIPDestination(eip)
	if err != nil {
		klog.Error(err)
		return err
	}

	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Scope:     netlink.SCOPE_LINK,
	}

	// Use RouteReplace (RTM_NEWROUTE with NLM_F_CREATE|NLM_F_REPLACE) rather than del+add:
	// the kernel updates the FIB entry atomically, so an existing route is never removed
	// mid-flight (no traffic gap) and a missing one is created. This makes it safe to
	// re-apply idempotently on every keepalive cycle regardless of the current state.
	if err := netlink.RouteReplace(route); err != nil {
		err = fmt.Errorf("failed to add route for EIP %s via %s: %w", eip, macvlanSubIfName, err)
		klog.Error(err)
		return err
	}

	// V(3): RouteReplace is idempotent and re-applied on every keepalive cycle.
	klog.V(3).Infof("added route for EIP %s via macvlan sub-interface %s", eip, macvlanSubIfName)
	return nil
}

// deleteEIPRoute deletes the route for EIP from the specified macvlan sub-interface.
func deleteEIPRoute(eip, macvlanSubIfName string) error {
	link, err := netlink.LinkByName(macvlanSubIfName)
	if err != nil {
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			// Interface gone means route is already gone
			return nil
		}
		err = fmt.Errorf("failed to get macvlan interface %s for deleting EIP %s route: %w", macvlanSubIfName, eip, err)
		klog.Error(err)
		return err
	}

	dst, err := parseEIPDestination(eip)
	if err != nil {
		klog.Error(err)
		return err
	}

	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Scope:     netlink.SCOPE_LINK,
	}
	if err := netlink.RouteDel(route); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			klog.V(3).Infof("route for EIP %s not found on %s", eip, macvlanSubIfName)
			return nil
		}
		err = fmt.Errorf("failed to delete route for EIP %s from %s: %w", eip, macvlanSubIfName, err)
		klog.Error(err)
		return err
	}
	klog.Infof("deleted route for EIP %s from %s", eip, macvlanSubIfName)
	return nil
}

// addVethEIPRoute adds a route for EIP via the specified veth next-hop IP.
func addVethEIPRoute(eip, nextHop string) error {
	dst, err := parseEIPDestination(eip)
	if err != nil {
		return err
	}
	gw := net.ParseIP(nextHop)
	if gw == nil {
		return fmt.Errorf("invalid veth nextHop IP: %s", nextHop)
	}

	route := &netlink.Route{
		Dst:   dst,
		Gw:    gw,
		Scope: netlink.SCOPE_UNIVERSE,
	}

	if err := netlink.RouteReplace(route); err != nil {
		return fmt.Errorf("failed to add veth route for EIP %s via %s: %w", eip, nextHop, err)
	}
	klog.V(3).Infof("added veth route for EIP %s via %s", eip, nextHop)
	return nil
}

// deleteVethEIPRoute deletes the route for EIP from host routing table.
func deleteVethEIPRoute(eip string) error {
	dst, err := parseEIPDestination(eip)
	if err != nil {
		return err
	}
	route := &netlink.Route{
		Dst: dst,
	}
	if err := netlink.RouteDel(route); err != nil && !errors.Is(err, syscall.ESRCH) {
		return fmt.Errorf("failed to delete veth route for EIP %s: %w", eip, err)
	}
	klog.Infof("deleted veth route for EIP %s", eip)
	return nil
}

// cleanHostLocalTableRouteForEIP deletes hijacked local table routes (table 255) for EIP.
func cleanHostLocalTableRouteForEIP(eip string) error {
	dst, err := parseEIPDestination(eip)
	if err != nil {
		return err
	}
	filter := &netlink.Route{
		Table: syscall.RT_TABLE_LOCAL,
		Dst:   dst,
	}
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_V4, filter, netlink.RT_FILTER_TABLE|netlink.RT_FILTER_DST)
	if err != nil {
		return nil
	}
	for _, r := range routes {
		if r.Type == syscall.RTN_LOCAL {
			if delErr := netlink.RouteDel(&r); delErr != nil {
				klog.Warningf("failed to delete local table route for EIP %s: %v", eip, delErr)
			} else {
				klog.Infof("successfully deleted hijacked local table route for EIP %s", eip)
			}
		}
	}
	return nil
}

// cleanHostStaleLinkRoutesForEIP deletes any stale scope-link routes for EIP (e.g., from macbond1).
func cleanHostStaleLinkRoutesForEIP(eip string) {
	dst, err := parseEIPDestination(eip)
	if err != nil {
		return
	}
	filter := &netlink.Route{
		Dst:   dst,
		Scope: netlink.SCOPE_LINK,
	}
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_V4, filter, netlink.RT_FILTER_DST|netlink.RT_FILTER_SCOPE)
	if err == nil {
		for _, r := range routes {
			if delErr := netlink.RouteDel(&r); delErr == nil {
				klog.Infof("successfully deleted stale scope-link route for EIP %s (linkIdx=%d)", eip, r.LinkIndex)
			}
		}
	}
}

// ensureHostForwarding ensures host IPv4 forwarding is enabled.
func ensureHostForwarding() {
	_ = os.WriteFile("/proc/sys/net/ipv4/ip_forward", []byte("1\n"), 0644)
}

// podNetworkStatusItem parses Multus network-status annotation.
type podNetworkStatusItem struct {
	Name      string   `json:"name"`
	Interface string   `json:"interface"`
	IPs       []string `json:"ips"`
	Default   bool     `json:"default"`
}

// getPodPtpIP extracts the point-to-point IP from gateway pod network-status annotation.
func getPodPtpIP(pod *corev1.Pod) string {
	const fallbackIP = "169.254.100.2"
	if pod == nil || pod.Annotations == nil {
		return fallbackIP
	}
	rawStatus := pod.Annotations["k8s.v1.cni.cncf.io/network-status"]
	if rawStatus == "" {
		return fallbackIP
	}
	var statuses []podNetworkStatusItem
	if err := json.Unmarshal([]byte(rawStatus), &statuses); err == nil {
		for _, st := range statuses {
			if !st.Default || st.Interface != "eth0" {
				for _, ip := range st.IPs {
					if strings.HasPrefix(ip, "169.254.") {
						return ip
					}
				}
			}
		}
	}
	return fallbackIP
}

// reconcileMacvlanSubnet is the main handler for macvlan subnet events.
// It processes subnet add/update/delete events and manages macvlan sub-interfaces accordingly.
// This function is called by the macvlan subnet worker with independent retry mechanism.
//
// Key design: One macvlan sub-interface per master interface (not per subnet).
// Multiple subnets may share the same master interface via the same NAD.
// The macvlan interface is created when the first subnet with that master appears,
// and deleted only when no more subnets use that master.
//
// Handles all transitions of NadMacvlanMasterAnnotation:
//   - Annotation addition: create macvlan sub-interface (if not exists)
//   - Annotation removal: delete macvlan sub-interface (if no other subnet uses it)
//   - Annotation value change: handle as removal + addition
func (c *Controller) reconcileMacvlanSubnet(event *subnetEvent) error {
	var oldSubnet, newSubnet *kubeovnv1.Subnet
	if event.oldObj != nil {
		oldSubnet, _ = event.oldObj.(*kubeovnv1.Subnet)
	}
	if event.newObj != nil {
		newSubnet, _ = event.newObj.(*kubeovnv1.Subnet)
	}

	oldMaster := ""
	if oldSubnet != nil {
		oldMaster = oldSubnet.Annotations[util.NadMacvlanMasterAnnotation]
	}
	newMaster := ""
	if newSubnet != nil {
		newMaster = newSubnet.Annotations[util.NadMacvlanMasterAnnotation]
	}

	// No change in master annotation, nothing to do
	if oldMaster == newMaster {
		return nil
	}

	// Master annotation removed or changed: check if we should delete old macvlan
	if oldMaster != "" {
		if c.shouldDeleteMacvlanForMaster(oldMaster, oldSubnet.Name) {
			if err := deleteMacvlanSubInterface(oldMaster); err != nil {
				klog.Errorf("failed to cleanup macvlan for master %s (subnet %s): %v", oldMaster, oldSubnet.Name, err)
				return err
			}
		} else {
			klog.V(3).Infof("macvlan for master %s still in use by other subnets, skipping delete", oldMaster)
		}
	}

	// Master annotation added or changed: create new macvlan interface
	if newMaster != "" {
		klog.V(3).Infof("creating macvlan sub-interface for subnet %s (master=%s)", newSubnet.Name, newMaster)
		if err := c.createMacvlanSubInterface(newMaster); err != nil {
			klog.Errorf("failed to create macvlan for subnet %s (master=%s): %v", newSubnet.Name, newMaster, err)
			return err
		}
	}

	return nil
}

// shouldDeleteMacvlanForMaster checks if the macvlan interface for the given master
// should be deleted. Returns true only if no other subnet uses the same master.
// excludeSubnet is the subnet being deleted/changed, which should be excluded from the check.
func (c *Controller) shouldDeleteMacvlanForMaster(master, excludeSubnet string) bool {
	// Use label selector to only list macvlan-type subnets for efficiency
	subnets, err := c.subnetsLister.List(labels.SelectorFromSet(labels.Set{
		util.NadMacvlanTypeLabel: "true",
	}))
	if err != nil {
		klog.Errorf("failed to list macvlan subnets: %v", err)
		return false // Don't delete if we can't check
	}

	for _, subnet := range subnets {
		if subnet.Name == excludeSubnet {
			continue
		}
		if subnet.Annotations[util.NadMacvlanMasterAnnotation] == master {
			return false // Another subnet uses this master
		}
	}
	return true
}

// loopEnsureMacvlanSubInterfaces is a keepalive loop that periodically ensures the
// macvlan sub-interfaces for all macvlan-type subnets exist and are up on the node, and
// that their EIP routes are present.
//
// The sub-interfaces are normally created once in response to subnet events. If an
// interface is later lost or brought down (e.g. the master interface flaps, or the
// sub-interface is deleted/set down by an external actor), there is no event to trigger
// recovery, which previously required a manual daemon pod restart. This loop recreates
// the interface (createMacvlanSubInterface also brings it back up) automatically.
//
// The kernel flushes an interface's SCOPE_LINK routes when it is deleted or set down, so
// after ensuring the interfaces we unconditionally re-apply the EIP routes. RouteReplace
// is idempotent, so refreshing every cycle is cheap and, unlike detecting the exact
// down/up transition, has no polling blind spot: no matter how the interface flaps
// between two polls, the next cycle restores its routes.
//
// TOCTOU note: this loop and the subnet delete worker both act on the shared informer
// cache from separate goroutines. If a subnet's List() snapshot here is taken just
// before the subnet deletion propagates, the delete worker may remove the interface in
// between, and this loop then recreates it as an orphan (no subnet references its
// master anymore, so nothing deletes it afterwards). The impact is bounded and benign:
// a macvlan NAD/subnet is effectively a cluster singleton that is meant to persist, so
// deletion (and thus this window) is extremely rare; an orphaned sub-interface carries
// no IP and no routes, so it causes no network issues and only occurs when the feature
// is being torn down anyway.
func (c *Controller) loopEnsureMacvlanSubInterfaces() {
	subnets, err := c.subnetsLister.List(labels.SelectorFromSet(labels.Set{
		util.NadMacvlanTypeLabel: "true",
	}))
	if err != nil {
		klog.Errorf("failed to list macvlan subnets: %v", err)
		return
	}

	// Collect the set of unique master interfaces in use by macvlan subnets.
	masters := make(map[string]struct{})
	for _, subnet := range subnets {
		if master := subnet.Annotations[util.NadMacvlanMasterAnnotation]; master != "" {
			masters[master] = struct{}{}
		}
	}

	if len(masters) == 0 {
		return
	}

	for master := range masters {
		if err := c.createMacvlanSubInterface(master); err != nil {
			klog.Errorf("failed to ensure macvlan sub-interface for master %s: %v", master, err)
		}
	}

	// Idempotently re-apply EIP routes every cycle so they survive any interface flap.
	c.resyncLocalIptablesEipRoutes()
}

// resyncLocalIptablesEipRoutes re-enqueues all ready IptablesEIPs so the route worker can
// converge each one on this node: syncIptablesEipRoute adds the route when the NAT GW pod
// is local and deletes it otherwise.
//
// It deliberately processes every ready EIP, not just local ones. The pod informer here is
// node-scoped and only has an update handler, so when a NAT GW pod is rescheduled to
// another node the old node reliably gets neither a usable update (Case 1 of
// handleNatGwPodUpdate cannot fire under a node-scoped informer) nor a delete event. This
// periodic full pass is what removes the now-stale route on the node the pod moved away
// from, and also self-heals after a daemon restart. Skipping non-local EIPs would leak
// those routes.
func (c *Controller) resyncLocalIptablesEipRoutes() {
	eips, err := c.iptablesEipsLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("failed to list iptables-eips for route resync: %v", err)
		return
	}

	for _, eip := range eips {
		if !shouldEnqueueIptablesEip(eip) {
			continue
		}
		info := c.buildEipRouteInfo(eip)
		if info == nil {
			continue
		}
		// V(4): this runs every keepalive cycle, so keep it quiet on the steady state.
		klog.V(4).Infof("re-enqueue iptables-eip %s to sync its route", eip.Name)
		c.iptablesEipQueue.Add(*info)
	}
}

// findLocalNatGwPod finds the NAT GW pod for the given NatGwDp on the local node.
// Supports both single-instance gateways and dual-member HA gateways (e.g. phase2-gw with phase2-gw-m1/m2).
func (c *Controller) findLocalNatGwPod(natGwDp string) *corev1.Pod {
	if natGwDp == "" {
		return nil
	}

	// 1. Direct match with deterministic single-instance name
	podName := util.GenNatGwPodName(natGwDp)
	if pod, err := c.podsLister.Pods(c.config.PodNamespace).Get(podName); err == nil {
		if pod.Spec.NodeName == c.config.NodeName {
			return pod
		}
	}

	// 2. Dual-member / HA match: find any local pod belonging to this logical gateway
	pods, err := c.podsLister.Pods(c.config.PodNamespace).List(labels.Everything())
	if err != nil {
		return nil
	}
	prefix := "vpc-nat-gw-" + natGwDp
	for _, pod := range pods {
		if pod.Spec.NodeName != c.config.NodeName {
			continue
		}
		// Match by pod name prefix (e.g. vpc-nat-gw-phase2-gw-m1-0 matches phase2-gw)
		if strings.HasPrefix(pod.Name, prefix) {
			return pod
		}
		// Match by labels
		if pod.Labels != nil {
			if gw := pod.Labels[util.VpcNatGatewayNameLabel]; gw == natGwDp || strings.HasPrefix(gw, natGwDp) {
				return pod
			}
			if app := pod.Labels["app"]; app == prefix || strings.HasPrefix(app, prefix) {
				return pod
			}
		}
	}
	return nil
}

// hasNatGwPodOnLocalNode checks if the NAT GW pod for the given NatGwDp is scheduled on the local node.
func (c *Controller) hasNatGwPodOnLocalNode(natGwDp string) bool {
	return c.findLocalNatGwPod(natGwDp) != nil
}

// shouldEnqueueIptablesEip checks if an EIP should be enqueued for route processing.
// Returns true if EIP is ready and has ExternalSubnet configured.
func shouldEnqueueIptablesEip(eip *kubeovnv1.IptablesEIP) bool {
	return eip.Spec.ExternalSubnet != "" && eip.Status.Ready
}

// buildEipRouteInfo builds eipRouteInfo from an EIP object.
// Returns nil if the EIP cannot be processed.
func (c *Controller) buildEipRouteInfo(eip *kubeovnv1.IptablesEIP) *eipRouteInfo {
	if eip.Spec.V4ip == "" {
		return nil
	}

	subnet, err := c.subnetsLister.Get(eip.Spec.ExternalSubnet)
	if err != nil {
		klog.Errorf("failed to get subnet %s for EIP %s: %v", eip.Spec.ExternalSubnet, eip.Name, err)
		return nil
	}

	if util.IsMacvlanSubnet(subnet) {
		if util.IsHAOrBgpEIP(eip) {
			klog.Errorf("HA/BGP IptablesEIP %s cannot use macvlan external subnet %s; HA mode only supports veth-ptp external subnets", eip.Name, subnet.Name)
			return nil
		}
		master := subnet.Annotations[util.NadMacvlanMasterAnnotation]
		macvlanName, err := util.GenMacvlanIfaceName(master)
		if err != nil {
			klog.Errorf("failed to generate macvlan name for EIP %s (master=%s): %v", eip.Name, master, err)
			return nil
		}

		return &eipRouteInfo{
			eipName:     eip.Name,
			v4ip:        eip.Spec.V4ip,
			macvlanName: macvlanName,
			mode:        "macvlan",
		}
	}

	// Non-macvlan subnet (e.g., ptp veth-pair subnet for HA / BGP Anycast)
	return &eipRouteInfo{
		eipName: eip.Name,
		v4ip:    eip.Spec.V4ip,
		mode:    "veth",
	}
}

// enqueueAddIptablesEip handles add events for IptablesEIP.
// This is primarily for daemon restart recovery: when the daemon restarts, the informer
// triggers Add events for all existing resources. Ready EIPs need to be re-processed
// to restore their routes on the node.
// Note: For newly created EIPs, they start with Ready=false and become Ready via an
// Update event, which is handled by enqueueUpdateIptablesEip.
func (c *Controller) enqueueAddIptablesEip(obj any) {
	eip := obj.(*kubeovnv1.IptablesEIP)

	// Clear deleted mark if EIP exists (handles edge case of stale mark)
	c.deletedEIPs.Delete(eip.Name)

	if !shouldEnqueueIptablesEip(eip) {
		return
	}

	info := c.buildEipRouteInfo(eip)
	if info == nil {
		return
	}

	klog.V(3).Infof("enqueue add iptables-eip %s for route recovery", eip.Name)
	c.iptablesEipQueue.Add(*info)
}

// enqueueUpdateIptablesEip handles update events for IptablesEIP.
// This is for normal runtime: when an EIP transitions to Ready state (after NAT Gateway
// configures iptables rules), add its route to the node.
func (c *Controller) enqueueUpdateIptablesEip(_, newObj any) {
	eip := newObj.(*kubeovnv1.IptablesEIP)

	// Skip EIPs that are being deleted
	if eip.DeletionTimestamp != nil {
		return
	}

	// Clear deleted mark if EIP was recreated with the same name
	c.deletedEIPs.Delete(eip.Name)

	if !shouldEnqueueIptablesEip(eip) {
		return
	}

	info := c.buildEipRouteInfo(eip)
	if info == nil {
		return
	}

	klog.V(3).Infof("enqueue update iptables-eip %s", eip.Name)
	c.iptablesEipQueue.Add(*info)
}

// enqueueDeleteIptablesEip handles delete events for IptablesEIP.
// It marks the EIP as deleted to prevent race conditions with queued add events,
// then enqueues the V4ip for route deletion with retry support.
func (c *Controller) enqueueDeleteIptablesEip(obj any) {
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

	// Mark EIP as deleted to prevent race conditions with queued add events
	c.deletedEIPs.Store(eip.Name, true)

	// Build route info - subnet must exist (protected by finalizer)
	info := c.buildEipRouteInfo(eip)
	if info == nil {
		return
	}

	klog.V(3).Infof("enqueue delete iptables-eip %s", eip.Name)
	c.iptablesEipDeleteQueue.Add(*info)
}

// runIptablesEipDeleteWorker runs the worker for deleting IptablesEIP routes.
func (c *Controller) runIptablesEipDeleteWorker() {
	for c.processNextIptablesEipDeleteItem() {
	}
}

// processNextIptablesEipDeleteItem processes the next delete work item.
func (c *Controller) processNextIptablesEipDeleteItem() bool {
	item, shutdown := c.iptablesEipDeleteQueue.Get()
	if shutdown {
		return false
	}

	err := func(info eipRouteInfo) error {
		defer c.iptablesEipDeleteQueue.Done(info)
		if err := c.handleDeleteIptablesEipRoute(info); err != nil {
			c.iptablesEipDeleteQueue.AddRateLimited(info)
			return fmt.Errorf("error deleting EIP route for %q: %w, requeuing", info.v4ip, err)
		}
		c.iptablesEipDeleteQueue.Forget(info)
		// Clean up the deleted mark after successful route deletion to prevent memory leak.
		// This is safe because:
		// 1. The route has been deleted successfully
		// 2. Any new EIP with the same name will trigger Add/Update events which clear the mark
		c.deletedEIPs.Delete(info.eipName)
		return nil
	}(item)
	if err != nil {
		klog.Error(err)
		return true
	}
	return true
}

// handleDeleteIptablesEipRoute deletes the route for an IptablesEIP.
func (c *Controller) handleDeleteIptablesEipRoute(info eipRouteInfo) error {
	klog.Infof("deleting iptables-eip route for %s (v4ip=%s, mode=%s, macvlan=%s)", info.eipName, info.v4ip, info.mode, info.macvlanName)
	if info.mode == "veth" {
		return deleteVethEIPRoute(info.v4ip)
	}
	return deleteEIPRoute(info.v4ip, info.macvlanName)
}

// runIptablesEipWorker runs the worker for syncing IptablesEIP routes.
func (c *Controller) runIptablesEipWorker() {
	for c.processNextIptablesEipItem() {
	}
}

// processNextIptablesEipItem processes the next work item.
func (c *Controller) processNextIptablesEipItem() bool {
	item, shutdown := c.iptablesEipQueue.Get()
	if shutdown {
		return false
	}

	err := func(info eipRouteInfo) error {
		defer c.iptablesEipQueue.Done(info)
		if err := c.syncIptablesEipRoute(info); err != nil {
			c.iptablesEipQueue.AddRateLimited(info)
			return fmt.Errorf("error syncing EIP route for %q: %w, requeuing", info.eipName, err)
		}
		c.iptablesEipQueue.Forget(info)
		return nil
	}(item)
	if err != nil {
		klog.Error(err)
		return true
	}
	return true
}

// syncIptablesEipRoute syncs the route for an IptablesEIP.
// It adds the route if NAT GW pod is on this node, or deletes the route if not.
func (c *Controller) syncIptablesEipRoute(info eipRouteInfo) error {
	klog.V(3).Infof("syncing iptables-eip route for %s (mode=%s)", info.eipName, info.mode)

	// Check if EIP was deleted - skip if it was to prevent race conditions
	if _, deleted := c.deletedEIPs.Load(info.eipName); deleted {
		klog.V(3).Infof("iptables-eip %s was deleted, skipping add", info.eipName)
		return nil
	}

	eip, err := c.iptablesEipsLister.Get(info.eipName)
	if err != nil {
		klog.V(3).Infof("iptables-eip %s not found: %v", info.eipName, err)
		return nil
	}

	// Only add routes for ready EIPs. An EIP becomes ready after NAT Gateway
	// pod successfully configures iptables rules. Before that, adding routes
	// would cause traffic to be blackholed.
	if !eip.Status.Ready {
		klog.V(3).Infof("iptables-eip %s not ready, skipping route", info.eipName)
		return nil
	}

	// If NAT GW pod is not on this node, delete routes (if any) and return
	if !c.hasNatGwPodOnLocalNode(eip.Spec.NatGwDp) {
		klog.V(3).Infof("NAT GW pod for iptables-eip %s not on local node, deleting routes if exist", info.eipName)
		if info.mode == "veth" {
			_ = deleteVethEIPRoute(info.v4ip)
		} else {
			if err := deleteEIPRoute(info.v4ip, info.macvlanName); err != nil {
				klog.V(3).Infof("failed to delete route for EIP %s (may not exist): %v", info.eipName, err)
			}
		}
		return nil
	}

	// 1. High-performance veth-pair mode (L3 Point-to-Point / BGP Anycast)
	if info.mode == "veth" {
		// Clean any hijacked routes in table 255 (local)
		_ = cleanHostLocalTableRouteForEIP(info.v4ip)
		// Clean any stale scope-link routes (e.g. dev macbond1 scope link)
		cleanHostStaleLinkRoutesForEIP(info.v4ip)
		// Ensure host IPv4 forwarding is active
		ensureHostForwarding()

		// Dynamically obtain pod's real ptp IP from network-status
		localPod := c.findLocalNatGwPod(eip.Spec.NatGwDp)
		if localPod == nil {
			klog.Errorf("failed to get local NAT GW pod for gateway %s", eip.Spec.NatGwDp)
			return fmt.Errorf("local NAT GW pod for %s not found", eip.Spec.NatGwDp)
		}
		ptpIP := getPodPtpIP(localPod)

		// Install penetration route via point-to-point interface
		if err := addVethEIPRoute(info.v4ip, ptpIP); err != nil {
			klog.Errorf("failed to add veth route for iptables-eip %s (V4ip=%s, ptpIP=%s): %v", info.eipName, info.v4ip, ptpIP, err)
			return err
		}

		// Also clean any legacy macvlan route if existed
		if info.macvlanName != "" {
			_ = deleteEIPRoute(info.v4ip, info.macvlanName)
		}
		return nil
	}

	// 2. Legacy macvlan mode: Add IPv4 route via macvlan sub-interface
	if err := addEIPRoute(info.v4ip, info.macvlanName); err != nil {
		klog.Errorf("failed to add IPv4 route for iptables-eip %s (V4ip=%s, macvlan=%s): %v", info.eipName, info.v4ip, info.macvlanName, err)
		return err
	}

	return nil
}

// isVpcNatGwPod checks if a pod is a VPC NAT Gateway pod
func isVpcNatGwPod(pod *corev1.Pod) bool {
	return pod.Labels[util.VpcNatGatewayLabel] == "true"
}

// getNatGwNameFromPod extracts the NAT GW name from a NAT GW pod via label.
func getNatGwNameFromPod(pod *corev1.Pod) string {
	return pod.Labels[util.VpcNatGatewayNameLabel]
}

// enqueueEipsForNatGw enqueues all EIPs associated with the given NAT GW for route processing.
// This is called when a NAT GW pod node changes or phase changes.
// The EIP handler will decide whether to add or delete routes based on current state.
func (c *Controller) enqueueEipsForNatGw(natGwName string) {
	eips, err := c.iptablesEipsLister.List(labels.SelectorFromSet(labels.Set{
		util.VpcNatGatewayNameLabel: natGwName,
	}))
	if err != nil {
		klog.Errorf("failed to list EIPs for NAT GW %s: %v", natGwName, err)
		return
	}

	for _, eip := range eips {
		if !shouldEnqueueIptablesEip(eip) {
			continue
		}
		info := c.buildEipRouteInfo(eip)
		if info == nil {
			continue
		}
		klog.Infof("enqueue iptables-eip %s for NAT GW %s pod event", eip.Name, natGwName)
		c.iptablesEipQueue.Add(*info)
	}
}

// handleNatGwPodUpdate handles NAT GW pod update events.
// When a NAT GW pod moves to/from this node or phase changes, enqueue all its EIPs for route processing.
// The EIP handler will decide whether to add or delete routes based on current pod location.
func (c *Controller) handleNatGwPodUpdate(oldPod, newPod *corev1.Pod) {
	if !isVpcNatGwPod(newPod) {
		return
	}

	oldNodeName := oldPod.Spec.NodeName
	newNodeName := newPod.Spec.NodeName

	// Skip if pod is not related to this node (neither old nor new node is this node)
	if oldNodeName != c.config.NodeName && newNodeName != c.config.NodeName {
		return
	}

	natGwName := getNatGwNameFromPod(newPod)
	if natGwName == "" {
		return
	}

	// Case 1: Pod moved from this node to another node - enqueue to delete routes
	if oldNodeName == c.config.NodeName && newNodeName != c.config.NodeName {
		klog.Infof("NAT GW pod %s moved from this node to %s, enqueuing EIPs to delete routes", newPod.Name, newNodeName)
		c.enqueueEipsForNatGw(natGwName)
		return
	}

	// Case 2: Pod moved from another node to this node - enqueue to add routes
	if oldNodeName != c.config.NodeName && newNodeName == c.config.NodeName {
		if newPod.Status.Phase == corev1.PodRunning {
			klog.Infof("NAT GW pod %s moved to this node, enqueuing EIPs to add routes", newPod.Name)
			c.enqueueEipsForNatGw(natGwName)
		}
		return
	}

	// Case 3: Pod is on this node and phase changed to Running - enqueue to add routes
	if oldPod.Status.Phase != corev1.PodRunning && newPod.Status.Phase == corev1.PodRunning {
		klog.Infof("NAT GW pod %s became running on this node, enqueuing EIPs to add routes", newPod.Name)
		c.enqueueEipsForNatGw(natGwName)
		return
	}

	// Case 4: Pod is on this node and phase changed from Running to non-Running
	// (e.g., being deleted with DeletionTimestamp set, or terminated) - enqueue to delete routes
	if oldPod.Status.Phase == corev1.PodRunning && newPod.Status.Phase != corev1.PodRunning {
		klog.Infof("NAT GW pod %s no longer running on this node (phase: %s), enqueuing EIPs to delete routes", newPod.Name, newPod.Status.Phase)
		c.enqueueEipsForNatGw(natGwName)
	}
}

// handleNatGwPodDelete handles NAT GW pod delete events. The pod informer is node-scoped,
// so a delete event only reaches the node the pod was running on, which is exactly where
// its EIP routes become stale. Enqueue the gateway's EIPs so the route worker removes them
// promptly instead of waiting for the periodic resync (the route worker deletes the routes
// because the NAT GW pod is no longer local).
func (c *Controller) handleNatGwPodDelete(pod *corev1.Pod) {
	if !isVpcNatGwPod(pod) {
		return
	}

	natGwName := getNatGwNameFromPod(pod)
	if natGwName == "" {
		return
	}

	klog.Infof("NAT GW pod %s deleted from this node, enqueuing EIPs to delete routes", pod.Name)
	c.enqueueEipsForNatGw(natGwName)
}
