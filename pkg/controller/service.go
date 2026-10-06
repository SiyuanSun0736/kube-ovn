package controller

import (
	"context"
	"fmt"
	"maps"
	"net"
	"reflect"
	"slices"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

type vpcService struct {
	Vips     []string
	Vpc      string
	Protocol v1.Protocol
	Svc      *v1.Service
}

type updateSvcObject struct {
	key      string
	oldPorts []v1.ServicePort
	newPorts []v1.ServicePort
	// hadBgpVipAnnotation records whether the previous Service object carried either
	// the bgp-vip or allow-shared-ip annotation. Used by handleUpdateService to
	// distinguish "annotation was removed from a controller-owned Service" (clear
	// stale ingress) from "unrelated LoadBalancer Service with no bgp annotation"
	// (leave ingress untouched).
	hadBgpVipAnnotation bool
}

func (c *Controller) enqueueAddService(obj any) {
	svc := obj.(*v1.Service)
	key := cache.MetaObjectToName(svc).String()
	c.enqueueGwNftableLbService(key)
	c.enqueueNatGwLanVipSync(svc.Annotations[util.VpcNatGatewayAnnotation])
	klog.V(3).Infof("enqueue add endpoint %s", key)
	c.addOrUpdateEndpointSliceQueue.Add(key)

	if c.config.EnableNP {
		netpols, err := c.svcMatchNetworkPolicies(svc)
		if err != nil {
			utilruntime.HandleError(err)
			return
		}
		for _, np := range netpols {
			c.updateNpQueue.Add(np)
		}
	}

	if c.config.EnablePodLbSvc || c.config.EnableBgpLbVip {
		klog.V(3).Infof("enqueue add service %s for lb processing", key)
		c.addServiceQueue.Add(key)
	}
}

func (c *Controller) enqueueDeleteService(obj any) {
	var svc *v1.Service
	switch t := obj.(type) {
	case *v1.Service:
		svc = t
	case cache.DeletedFinalStateUnknown:
		s, ok := t.Obj.(*v1.Service)
		if !ok {
			klog.Warningf("unexpected object type: %T", t.Obj)
			return
		}
		svc = s
	default:
		klog.Warningf("unexpected type: %T", obj)
		return
	}

	// The gateway nftable LB mode owns only the EIP data plane of such a Service; its ClusterIP
	// VIPs keep belonging to the classic implementation, so its cleanup always runs here.
	klog.Infof("enqueue delete service %s/%s", svc.Namespace, svc.Name)
	c.enqueueGwNftableLbService(svc.Namespace + "/" + svc.Name)
	// The gateway annotation survives on the tombstone object: its lanIP partition loses this
	// Service's backends now, and the identity is dropped once nobody else claims the port.
	c.enqueueNatGwLanVipSync(svc.Annotations[util.VpcNatGatewayAnnotation])

	if svc.Annotations[util.SwitchLBRuleVipsAnnotation] != "" ||
		(svc.Spec.ClusterIP != v1.ClusterIPNone && svc.Spec.ClusterIP != "") ||
		svc.Annotations[util.ServiceExternalIPFromSubnetAnnotation] != "" {
		if c.config.EnableNP {
			netpols, err := c.svcMatchNetworkPolicies(svc)
			if err != nil {
				utilruntime.HandleError(err)
				return
			}
			for _, np := range netpols {
				c.updateNpQueue.Add(np)
			}
		}

		ips := util.ServiceClusterIPs(*svc)
		if vip := svc.Annotations[util.SwitchLBRuleVipsAnnotation]; vip != "" {
			ips = strings.Split(vip, ",")
		}

		if svc.Annotations[util.ServiceExternalIPFromSubnetAnnotation] != "" {
			for _, ingress := range svc.Status.LoadBalancer.Ingress {
				ips = append(ips, ingress.IP)
			}
		}

		for _, port := range svc.Spec.Ports {
			vpcSvc := &vpcService{
				Protocol: port.Protocol,
				Vpc:      svc.Annotations[util.VpcAnnotation],
				Svc:      svc,
			}
			for _, ip := range ips {
				vpcSvc.Vips = append(vpcSvc.Vips, util.JoinHostPort(ip, port.Port))
			}
			klog.V(3).Infof("delete vpc service: %v", vpcSvc)
			c.deleteServiceQueue.Add(vpcSvc)
		}
	}
}

func (c *Controller) enqueueUpdateService(oldObj, newObj any) {
	oldSvc := oldObj.(*v1.Service)
	newSvc := newObj.(*v1.Service)
	if oldSvc.ResourceVersion == newSvc.ResourceVersion {
		return
	}

	oldClusterIps := getVipIps(oldSvc)
	newClusterIps := getVipIps(newSvc)
	var ipsToDel []string
	for _, oldClusterIP := range oldClusterIps {
		if !slices.Contains(newClusterIps, oldClusterIP) {
			ipsToDel = append(ipsToDel, oldClusterIP)
		}
	}

	key := cache.MetaObjectToName(newSvc).String()
	// The gateway mode's reconcile only depends on the Service fields gwNftableLbSvcChanged
	// compares, so an unrelated status or annotation churn does not wake it. The lanIP
	// partition consumes the same fields (gateway annotation, ports, affinity), so it shares
	// the filter; both gateways matter because the annotation may have moved the Service.
	if gwNftableLbSvcChanged(oldSvc, newSvc) {
		c.enqueueGwNftableLbService(key)
		c.enqueueNatGwLanVipSync(oldSvc.Annotations[util.VpcNatGatewayAnnotation])
		c.enqueueNatGwLanVipSync(newSvc.Annotations[util.VpcNatGatewayAnnotation])
	}
	klog.V(3).Infof("enqueue update service %s", key)
	if len(ipsToDel) != 0 {
		ipsToDelStr := strings.Join(ipsToDel, ",")
		key = strings.Join([]string{key, ipsToDelStr}, "#")
	}

	updateSvc := &updateSvcObject{
		key:                 key,
		oldPorts:            oldSvc.Spec.Ports,
		newPorts:            newSvc.Spec.Ports,
		hadBgpVipAnnotation: oldSvc.Annotations[util.MetalLBAllowSharedIPAnnotation] != "" || oldSvc.Annotations[util.BgpVipAnnotation] != "",
	}
	c.updateServiceQueue.Add(updateSvc)
}

// handleDeleteService cleans up resources associated with a deleted Service.
//
// When EnableBgpLbVip=true and EnableOvnLB=false, this function has no interaction
// with OVN whatsoever. The only cleanup that runs is cleanBgpLbVipService, which
// is a pure IPAM-level operation (BGP withdrawal happens automatically when the
// speaker no longer sees the bgp annotation on the Service object).
func (c *Controller) handleDeleteService(service *vpcService) error {
	key := cache.MetaObjectToName(service.Svc).String()

	c.svcKeyMutex.LockKey(key)
	defer func() { _ = c.svcKeyMutex.UnlockKey(key) }()
	klog.Infof("handle delete service %s", key)

	// OVN LB VIP cleanup is only relevant when the classic OVN LB mode is active;
	// see deleteServiceOvnLBVips for details.
	if c.config.EnableOvnLB {
		if err := c.deleteServiceOvnLBVips(service); err != nil {
			return err
		}
	}

	if c.config.EnablePodLbSvc && service.Svc.Spec.Type == v1.ServiceTypeLoadBalancer {
		if err := c.deleteLbSvc(service.Svc); err != nil {
			klog.Errorf("failed to delete service %s, %v", service.Svc.Name, err)
			return err
		}
	}

	if c.config.EnableBgpLbVip && service.Svc.Spec.Type == v1.ServiceTypeLoadBalancer {
		if err := c.cleanBgpLbVipService(service.Svc); err != nil {
			klog.Errorf("failed to clean bgp-lb-vip for service %s: %v", service.Svc.Name, err)
			return err
		}
	}

	return nil
}

// deleteServiceOvnLBVips removes ClusterIP VIPs from OVN load balancer tables for a
// deleted Service. Only called when EnableOvnLB=true (classic OVN LB mode).
func (c *Controller) deleteServiceOvnLBVips(service *vpcService) error {
	svcs, err := c.servicesLister.Services(v1.NamespaceAll).List(labels.Everything())
	if err != nil {
		klog.Errorf("failed to list svc, %v", err)
		return err
	}

	var (
		vpcLB             [2]string
		vpcLbConfig       = c.GenVpcLoadBalancer(service.Vpc)
		ignoreHealthCheck = true
	)

	switch service.Protocol {
	case v1.ProtocolTCP:
		vpcLB = [2]string{vpcLbConfig.TCPLoadBalancer, vpcLbConfig.TCPSessLoadBalancer}
	case v1.ProtocolUDP:
		vpcLB = [2]string{vpcLbConfig.UDPLoadBalancer, vpcLbConfig.UDPSessLoadBalancer}
	case v1.ProtocolSCTP:
		vpcLB = [2]string{vpcLbConfig.SctpLoadBalancer, vpcLbConfig.SctpSessLoadBalancer}
	}

	for _, vip := range service.Vips {
		var (
			ip    string
			found bool
		)
		ip = parseVipAddr(vip)

		for _, svc := range svcs {
			if slices.Contains(util.ServiceClusterIPs(*svc), ip) {
				found = true
				break
			}
		}
		if found {
			continue
		}

		for _, lb := range vpcLB {
			if c.config.EnableOVNLBPreferLocal {
				if err = c.OVNNbClient.LoadBalancerDeleteIPPortMapping(lb, vip); err != nil {
					klog.Errorf("failed to delete ip port mapping for vip %s from LB %s: %v", vip, lb, err)
					return err
				}
			}

			if err = c.OVNNbClient.LoadBalancerDeleteVip(lb, vip, ignoreHealthCheck); err != nil {
				klog.Errorf("failed to delete vip %s from LB %s: %v", vip, lb, err)
				return err
			}
		}
	}

	return nil
}

// handleUpdateService reconciles a Service update.
//
// When EnableBgpLbVip=true: annotation changes (adding, removing, or switching the
// bgp-vip / allow-shared-ip annotation) must be reconciled immediately so that
// status.loadBalancer.ingress tracks the current annotation. reconcileBgpLbVipServiceLocked
// is idempotent and only calls UpdateStatus when the ingress actually differs.
func (c *Controller) handleUpdateService(svcObject *updateSvcObject) error {
	key := svcObject.key
	keys := strings.Split(key, "#")
	key = keys[0]
	var ipsToDel []string
	if len(keys) == 2 {
		ipsToDelStr := keys[1]
		ipsToDel = strings.Split(ipsToDelStr, ",")
	}

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Error(err)
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	c.svcKeyMutex.LockKey(key)
	defer func() { _ = c.svcKeyMutex.UnlockKey(key) }()
	klog.Infof("handle update service %s", key)

	svc, err := c.servicesLister.Services(namespace).Get(name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}

	ips := getVipIps(svc)

	// OVN ClusterIP VIP management is only relevant in classic OVN LB mode;
	// see syncServiceOvnLBVips for details.
	if c.config.EnableOvnLB {
		if err := c.syncServiceOvnLBVips(key, svc, ips, ipsToDel); err != nil {
			return err
		}
	}

	if c.config.EnablePodLbSvc && svc.Spec.Type == v1.ServiceTypeLoadBalancer {
		changed, err := c.checkLbSvcDeployAnnotationChanged(svc)
		if err != nil {
			klog.Errorf("failed to check annotation change for lb svc %s: %v", key, err)
			return err
		}

		// only process svc.spec.ports update
		if !changed {
			klog.Infof("update loadbalancer service %s", key)
			pod, err := c.getLbSvcPod(name, namespace)
			if err != nil {
				klog.Errorf("failed to get pod for lb svc %s: %v", key, err)
				if strings.Contains(err.Error(), "not found") {
					return nil
				}
				return err
			}

			toDel := diffSvcPorts(svcObject.oldPorts, svcObject.newPorts)
			if err := c.delDnatRules(pod, toDel, svc); err != nil {
				klog.Errorf("failed to delete dnat rules, err: %v", err)
				return err
			}
			if err = c.updatePodAttachNets(pod, svc); err != nil {
				klog.Errorf("failed to update pod attachment network for lb svc %s: %v", key, err)
				return err
			}
		}
	}

	if c.config.EnableBgpLbVip {
		// Re-reconcile so that annotation changes (add/remove/switch bgp-vip or
		// allow-shared-ip) are immediately reflected in status.loadBalancer.ingress.
		// The mutex is already held here; call the locked variant directly.
		//
		// Ownership guard: only clear stale ingress when this controller previously
		// owned the Service (old annotation was set). Services managed by other LB
		// controllers are left completely untouched.
		vipName := svc.Annotations[util.MetalLBAllowSharedIPAnnotation]
		if vipName == "" {
			vipName = svc.Annotations[util.BgpVipAnnotation]
		}
		if vipName == "" && svcObject.hadBgpVipAnnotation && len(svc.Status.LoadBalancer.Ingress) > 0 {
			// Annotation was removed from a controller-owned Service; clear the stale
			// ingress so kube-proxy and the BGP speaker see a clean state.
			// EnableBgpLbVip and EnablePodLbSvc are mutually exclusive (config validate),
			// so no other LB controller owns the ingress field when this path runs.
			targetSvc := svc.DeepCopy()
			targetSvc.Status.LoadBalancer.Ingress = nil
			if _, err := c.config.KubeClient.CoreV1().Services(svc.Namespace).UpdateStatus(
				context.Background(), targetSvc, metav1.UpdateOptions{},
			); err != nil {
				klog.Errorf("failed to clear stale ingress for service %s/%s: %v", svc.Namespace, svc.Name, err)
				return err
			}
			klog.Infof("bgp-lb-vip service %s/%s: annotation removed, cleared stale ingress", svc.Namespace, svc.Name)
			return nil
		}
		if vipName != "" {
			return c.reconcileBgpLbVipServiceLocked(key, svc)
		}
		return nil
	}

	return nil
}

// syncServiceOvnLBVips reconciles ClusterIP VIPs in OVN load balancer tables on Service
// update: adds missing VIPs, removes stale VIPs, and enqueues endpoint/SLR re-syncs as
// needed. Only called when EnableOvnLB=true (classic OVN LB mode).
func (c *Controller) syncServiceOvnLBVips(key string, svc *v1.Service, ips, ipsToDel []string) error {
	// Re-derive ingress IPs from a fresh subnet check instead of trusting the
	// ServiceExternalIPFromSubnetAnnotation annotation.  That annotation is
	// user-writable; reading it here would allow a user who can PATCH Service
	// annotations to inject arbitrary ingress IPs into the OVN load balancer.
	// The SwitchLBRuleVipsAnnotation path is kept as-is (explicit VIPs).
	if _, hasExplicitVips := svc.Annotations[util.SwitchLBRuleVipsAnnotation]; !hasExplicitVips {
		freshIPs := util.ServiceClusterIPs(*svc)
		subnets, err := c.subnetsLister.List(labels.Everything())
		if err != nil {
			klog.Errorf("failed to list subnets for service %s: %v", key, err)
			return err
		}
	ingressLoop:
		for _, ingress := range svc.Status.LoadBalancer.Ingress {
			if ingress.IP == "" {
				continue
			}
			for _, subnet := range subnets {
				if util.CIDRContainIP(subnet.Spec.CIDRBlock, ingress.IP) {
					freshIPs = append(freshIPs, ingress.IP)
					continue ingressLoop
				}
			}
		}
		ips = freshIPs
	}
	vpcName := svc.Annotations[util.VpcAnnotation]
	if vpcName == "" {
		vpcName = c.config.ClusterRouter
	}
	vpc, err := c.vpcsLister.Get(vpcName)
	if err != nil {
		klog.Errorf("failed to get vpc %s of lb, %v", vpcName, err)
		return err
	}

	tcpLb, udpLb, sctpLb := vpc.Status.TCPLoadBalancer, vpc.Status.UDPLoadBalancer, vpc.Status.SctpLoadBalancer
	oTCPLb, oUDPLb, oSctpLb := vpc.Status.TCPSessionLoadBalancer, vpc.Status.UDPSessionLoadBalancer, vpc.Status.SctpSessionLoadBalancer
	if svc.Spec.SessionAffinity == v1.ServiceAffinityClientIP {
		tcpLb, udpLb, sctpLb, oTCPLb, oUDPLb, oSctpLb = oTCPLb, oUDPLb, oSctpLb, tcpLb, udpLb, sctpLb
	}

	var tcpVips, udpVips, sctpVips []string
	for _, port := range svc.Spec.Ports {
		for _, ip := range ips {
			switch port.Protocol {
			case v1.ProtocolTCP:
				tcpVips = append(tcpVips, util.JoinHostPort(ip, port.Port))
			case v1.ProtocolUDP:
				udpVips = append(udpVips, util.JoinHostPort(ip, port.Port))
			case v1.ProtocolSCTP:
				sctpVips = append(sctpVips, util.JoinHostPort(ip, port.Port))
			}
		}
	}

	var (
		needUpdateEndpointQueue = false
		ignoreHealthCheck       = true
	)

	// for service update
	updateVip := func(lbName, oLbName string, svcVips []string) error {
		if len(lbName) == 0 {
			return nil
		}

		lb, err := c.OVNNbClient.GetLoadBalancer(lbName, false)
		if err != nil {
			klog.Errorf("failed to get LB %s: %v", lbName, err)
			return err
		}
		lbVIPs := maps.Clone(lb.Vips)
		klog.V(3).Infof("existing vips of LB %s: %v", lbName, lbVIPs)
		for _, vip := range svcVips {
			if err := c.OVNNbClient.LoadBalancerDeleteVip(oLbName, vip, ignoreHealthCheck); err != nil {
				klog.Errorf("failed to delete vip %s from LB %s: %v", vip, oLbName, err)
				return err
			}

			if _, ok := lbVIPs[vip]; !ok {
				klog.Infof("add vip %s to LB %s", vip, lbName)
				needUpdateEndpointQueue = true
			}
		}
		for vip := range lbVIPs {
			if ip := parseVipAddr(vip); (slices.Contains(ips, ip) && !slices.Contains(svcVips, vip)) || slices.Contains(ipsToDel, ip) {
				klog.Infof("remove stale vip %s from LB %s", vip, lbName)
				if err := c.OVNNbClient.LoadBalancerDeleteVip(lbName, vip, ignoreHealthCheck); err != nil {
					klog.Errorf("failed to delete vip %s from LB %s: %v", vip, lbName, err)
					return err
				}
			}
		}

		if len(oLbName) == 0 {
			return nil
		}

		oLb, err := c.OVNNbClient.GetLoadBalancer(oLbName, false)
		if err != nil {
			klog.Errorf("failed to get LB %s: %v", oLbName, err)
			return err
		}
		oLbVIPs := maps.Clone(oLb.Vips)
		klog.V(3).Infof("existing vips of LB %s: %v", oLbName, oLbVIPs)
		for vip := range oLbVIPs {
			if ip := parseVipAddr(vip); slices.Contains(ips, ip) || slices.Contains(ipsToDel, ip) {
				klog.Infof("remove stale vip %s from LB %s", vip, oLbName)
				if err = c.OVNNbClient.LoadBalancerDeleteVip(oLbName, vip, ignoreHealthCheck); err != nil {
					klog.Errorf("failed to delete vip %s from LB %s: %v", vip, oLbName, err)
					return err
				}
			}
		}
		return nil
	}

	if err = updateVip(tcpLb, oTCPLb, tcpVips); err != nil {
		klog.Error(err)
		return err
	}
	if err = updateVip(udpLb, oUDPLb, udpVips); err != nil {
		klog.Error(err)
		return err
	}
	if err = updateVip(sctpLb, oSctpLb, sctpVips); err != nil {
		klog.Error(err)
		return err
	}

	if err := c.checkServiceLBIPBelongToSubnet(svc); err != nil {
		klog.Error(err)
		return err
	}

	if needUpdateEndpointQueue {
		c.addOrUpdateEndpointSliceQueue.Add(key)
	}
	// add the svc key which has the same vip
	vip, ok := svc.Annotations[util.SwitchLBRuleVipsAnnotation]
	if ok && vip != "" {
		allSlrs, err := c.switchLBRuleLister.List(labels.Everything())
		if err != nil {
			klog.Error(err)
			return err
		}
		for _, slr := range allSlrs {
			if slr.Spec.Vip == vip {
				slrKey := fmt.Sprintf("%s/slr-%s", slr.Spec.Namespace, slr.Name)
				c.addOrUpdateEndpointSliceQueue.Add(slrKey)
			}
		}
	}
	return nil
}

// Parse key of map, [fd00:10:96::11c9]:10665 for example
func parseVipAddr(vip string) string {
	host, _, err := net.SplitHostPort(vip)
	if err != nil {
		klog.Errorf("failed to parse vip %q: %v", vip, err)
		return ""
	}
	return host
}

// handleAddService is the entry point for LoadBalancer Service creation.
//
// Mode matrix:
//   - EnableBgpLbVip=true : pure IPAM path — binds a VIP CR to the Service's
//     externalIPs and lets the BGP speaker announce it. No OVN LB objects are
//     created or modified; OVN is not involved at all.
//   - EnablePodLbSvc=true  : Pod-based path — creates a dedicated lb-svc Pod that
//     provides external connectivity via iptables NAT inside OVN.
//   - EnableOvnLB=true    : classic OVN LB path — handled via the update/delete
//     workers (add just enqueues endpoint sync).
func (c *Controller) handleAddService(key string) error {
	if c.config.EnableBgpLbVip {
		// BGP LB EIP mode: only IPAM and BGP announcement are involved.
		// enable-bgp-lb-vip and enable-lb-svc are mutually exclusive; the
		// lb-svc Pod-based flow is intentionally skipped.
		klog.Infof("dispatch add service %s to bgp-lb-vip handler", key)
		return c.handleAddBgpLbVipService(key)
	}
	if !c.config.EnablePodLbSvc {
		return nil
	}

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Error(err)
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	c.svcKeyMutex.LockKey(key)
	defer func() { _ = c.svcKeyMutex.UnlockKey(key) }()
	klog.Infof("handle add service %s", key)

	svc, err := c.servicesLister.Services(namespace).Get(name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if svc.Spec.Type != v1.ServiceTypeLoadBalancer {
		return nil
	}
	// Skip non kube-ovn lb-svc.
	if _, ok := svc.Annotations[util.AttachmentProvider]; !ok {
		return nil
	}

	klog.Infof("handle add loadbalancer service %s", key)

	if err = c.validateSvc(svc); err != nil {
		c.recorder.Event(svc, v1.EventTypeWarning, "ValidateSvcFailed", err.Error())
		klog.Errorf("failed to validate lb svc %s: %v", key, err)
		return err
	}

	nad, err := c.getAttachNetworkForService(svc)
	if err != nil {
		c.recorder.Event(svc, v1.EventTypeWarning, "GetNADFailed", err.Error())
		klog.Errorf("failed to check attachment network of lb svc %s: %v", key, err)
		return err
	}

	if err = c.createLbSvcPod(svc, nad); err != nil {
		klog.Errorf("failed to create lb svc pod for %s: %v", key, err)
		return err
	}

	var pod *v1.Pod
	for {
		pod, err = c.getLbSvcPod(name, namespace)
		if err != nil {
			klog.Warningf("pod for lb svc %s is not running: %v", key, err)
			time.Sleep(time.Second)
		}
		if pod != nil {
			break
		}

		// It's important here to check existing of svc, used to break the loop.
		_, err = c.servicesLister.Services(namespace).Get(name)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return nil
			}
			klog.Error(err)
			return err
		}
	}

	loadBalancerIP, err := c.getPodAttachIP(pod, svc)
	if err != nil {
		klog.Errorf("failed to get loadBalancerIP: %v", err)
		return err
	}

	svc, err = c.servicesLister.Services(namespace).Get(name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	targetSvc := svc.DeepCopy()
	if err = c.updatePodAttachNets(pod, targetSvc); err != nil {
		klog.Errorf("failed to update pod attachment network for service %s/%s: %v", namespace, name, err)
		return err
	}

	// compatible with IPv4 and IPv6 dual stack subnet.
	var ingress []v1.LoadBalancerIngress
	for ip := range strings.SplitSeq(loadBalancerIP, ",") {
		if ip != "" && net.ParseIP(ip) != nil {
			ingress = append(ingress, v1.LoadBalancerIngress{IP: ip})
		}
	}
	targetSvc.Status.LoadBalancer.Ingress = ingress
	if !equality.Semantic.DeepEqual(svc.Status, targetSvc.Status) {
		if _, err = c.config.KubeClient.CoreV1().Services(namespace).
			UpdateStatus(context.Background(), targetSvc, metav1.UpdateOptions{}); err != nil {
			klog.Errorf("failed to update status of service %s/%s: %v", namespace, name, err)
			return err
		}
	}

	return nil
}

func getVipIps(svc *v1.Service) []string {
	var ips []string
	if vip, ok := svc.Annotations[util.SwitchLBRuleVipsAnnotation]; ok {
		for _, ip := range strings.Split(vip, ",") {
			if ip != "" {
				ips = append(ips, ip)
			}
		}
	} else {
		ips = util.ServiceClusterIPs(*svc)
		if svc.Annotations[util.ServiceExternalIPFromSubnetAnnotation] != "" {
			for _, ingress := range svc.Status.LoadBalancer.Ingress {
				if ingress.IP != "" {
					ips = append(ips, ingress.IP)
				}
			}
		}
	}
	return ips
}

func diffSvcPorts(oldPorts, newPorts []v1.ServicePort) (toDel []v1.ServicePort) {
	for _, oldPort := range oldPorts {
		found := false
		for _, newPort := range newPorts {
			if reflect.DeepEqual(oldPort, newPort) {
				found = true
				break
			}
		}
		if !found {
			toDel = append(toDel, oldPort)
		}
	}

	return toDel
}

func (c *Controller) checkServiceLBIPBelongToSubnet(svc *v1.Service) error {
	if svc.Spec.Type != v1.ServiceTypeLoadBalancer {
		return nil
	}
	// Note: do NOT short-circuit when len(ingress)==0. An empty ingress list means
	// any previously set ServiceExternalIPFromSubnetAnnotation must be removed.
	// Returning early here would leave the annotation stale after ingress is cleared.

	svc = svc.DeepCopy()
	if svc.Annotations == nil {
		svc.Annotations = map[string]string{}
	}

	origAnnotation := svc.Annotations[util.ServiceExternalIPFromSubnetAnnotation]

	subnets, err := c.subnetsLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("failed to list subnets: %v", err)
		return err
	}

	isServiceExternalIPFromSubnet := false
outer:
	for _, subnet := range subnets {
		for _, ingress := range svc.Status.LoadBalancer.Ingress {
			if util.CIDRContainIP(subnet.Spec.CIDRBlock, ingress.IP) {
				svc.Annotations[util.ServiceExternalIPFromSubnetAnnotation] = subnet.Name
				isServiceExternalIPFromSubnet = true
				break outer
			}
		}
	}

	if !isServiceExternalIPFromSubnet {
		delete(svc.Annotations, util.ServiceExternalIPFromSubnetAnnotation)
	}

	newAnnotation := svc.Annotations[util.ServiceExternalIPFromSubnetAnnotation]
	if newAnnotation == origAnnotation {
		return nil
	}

	klog.Infof("Service %s/%s external IP belongs to subnet: %v", svc.Namespace, svc.Name, isServiceExternalIPFromSubnet)
	if _, err = c.config.KubeClient.CoreV1().Services(svc.Namespace).Update(context.Background(), svc, metav1.UpdateOptions{}); err != nil {
		klog.Errorf("failed to update service %s/%s: %v", svc.Namespace, svc.Name, err)
		return err
	}

	return nil
}

// handleAddBgpLbVipService binds a pre-allocated VIP (type=bgp_lb_vip) to a
// LoadBalancer Service. Control flow:
//
//	Vip(spec/status) -> service controller -> Service annotation/status -> speaker -> BGP announce.
//
// The VIP is identified by the ovn.kubernetes.io/bgp-vip annotation.
// Once bound:
//   - svc.status.loadBalancer.ingress is set so the BGP speaker announces the IP and
//     kube-proxy installs DNAT rules (for type=LoadBalancer, kube-proxy processes ingress
//     IPs automatically — spec.externalIPs is NOT set to avoid duplicate kubectl output).
//   - svc.annotations[ovn.kubernetes.io/bgp] is set to "true" so the speaker picks it up.
func (c *Controller) handleAddBgpLbVipService(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Error(err)
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	c.svcKeyMutex.LockKey(key)
	defer func() { _ = c.svcKeyMutex.UnlockKey(key) }()

	svc, err := c.servicesLister.Services(namespace).Get(name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	if svc.Spec.Type != v1.ServiceTypeLoadBalancer {
		return nil
	}

	return c.reconcileBgpLbVipServiceLocked(key, svc)
}

func (c *Controller) reconcileBgpLbVipServiceLocked(key string, svc *v1.Service) error {
	if svc.Spec.Type != v1.ServiceTypeLoadBalancer {
		return nil
	}

	// Resolve VIP CR name: MetalLB compat annotation takes priority over bgp-vip.
	// metallb.universe.tf/allow-shared-ip value is the VIP CR name (same semantics
	// as ovn.kubernetes.io/bgp-vip); naming it after the IP is a naming convention only.
	vipName := svc.Annotations[util.MetalLBAllowSharedIPAnnotation]
	if vipName == "" {
		vipName = svc.Annotations[util.BgpVipAnnotation]
	}
	if vipName == "" {
		// Service does not request a BGP LB VIP; nothing to do.
		// Stale-ingress cleanup after annotation removal is handled in handleUpdateService,
		// which has ownership context via updateSvcObject.hadBgpVipAnnotation.
		return nil
	}

	namespace, name := svc.Namespace, svc.Name

	klog.Infof("reconcile bgp-lb-vip service %s: vip=%s ingress=%v", key, vipName, svc.Status.LoadBalancer.Ingress)

	vip, err := c.virtualIpsLister.Get(vipName)
	if err != nil {
		klog.Errorf("failed to get vip %s for service %s: %v", vipName, key, err)
		return err
	}
	// Only bgp_lb_vip is supported here because this path expects an IPAM-only VIP
	// that is written into Service ingress for speaker announcement.
	if vip.Spec.Type != util.BgpLbVip {
		return fmt.Errorf("vip %s has type %q, expected %q", vipName, vip.Spec.Type, util.BgpLbVip)
	}
	if vip.Status.V4ip == "" && vip.Status.V6ip == "" {
		// IP not yet allocated. enqueueUpdateVirtualIP in vip.go will re-enqueue this
		// Service once the VIP receives its IP, so return nil to avoid retry noise.
		klog.Infof("bgp-lb-vip service %s: vip %s has no IP yet, skip", key, vipName)
		return nil
	}

	// Build ingress list from all allocated IPs (dual-stack and IPv6-only supported).
	var ingress []v1.LoadBalancerIngress
	if vip.Status.V4ip != "" {
		ingress = append(ingress, v1.LoadBalancerIngress{IP: vip.Status.V4ip})
	}
	if vip.Status.V6ip != "" {
		ingress = append(ingress, v1.LoadBalancerIngress{IP: vip.Status.V6ip})
	}
	klog.Infof("bgp-lb-vip service %s resolved vip %s to ips %v", key, vipName, ingress)

	targetSvc := svc.DeepCopy()

	// Set status.loadBalancer.ingress so the BGP speaker discovers the IP(s).
	// The speaker uses ovn.kubernetes.io/bgp-vip or metallb.universe.tf/allow-shared-ip
	// as the announcement gate — no extra controller-written annotation is needed.
	// kube-proxy also uses ingress IPs for DNAT rules on type=LoadBalancer services.
	// We only set Ingress here (not spec.externalIPs) — consistent with MetalLB.
	if !equality.Semantic.DeepEqual(targetSvc.Status.LoadBalancer.Ingress, ingress) {
		klog.Infof("bgp-lb-vip service %s updating status ingress: %v -> %v", key, targetSvc.Status.LoadBalancer.Ingress, ingress)
		targetSvc.Status.LoadBalancer.Ingress = ingress
		if _, err = c.config.KubeClient.CoreV1().Services(namespace).UpdateStatus(
			context.Background(), targetSvc, metav1.UpdateOptions{},
		); err != nil {
			klog.Errorf("failed to update status for service %s/%s: %v", namespace, name, err)
			return err
		}
		klog.Infof("bgp-lb-vip service %s status updated successfully", key)
	}

	klog.Infof("bgp-lb-vip service %s bound to vip %s ingress=%v", key, vipName, ingress)
	return nil
}

func (c *Controller) cleanupBgpLbVipServiceBindingByVip(vipName string) error {
	objs, err := c.svcByBgpVipIndexer.ByIndex(bgpVipIndexName, vipName)
	if err != nil {
		klog.Errorf("failed to index services for bgp-lb-vip %s cleanup: %v", vipName, err)
		return err
	}

	for _, obj := range objs {
		svc, ok := obj.(*v1.Service)
		if !ok {
			continue
		}
		// Apply the same VIP priority rule as reconcile: allow-shared-ip takes precedence
		// over bgp-vip.  If the Service's effective VIP is different from the one being
		// deleted (e.g. bgp-vip=X, allow-shared-ip=Y and we are deleting X), the Service
		// is still actively bound via Y — skip cleanup to avoid incorrectly clearing ingress.
		effectiveVip := svc.Annotations[util.MetalLBAllowSharedIPAnnotation]
		if effectiveVip == "" {
			effectiveVip = svc.Annotations[util.BgpVipAnnotation]
		}
		if effectiveVip != vipName {
			continue
		}
		if err := c.cleanupBgpLbVipServiceBinding(svc); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) cleanupBgpLbVipServiceBinding(svc *v1.Service) error {
	if len(svc.Status.LoadBalancer.Ingress) == 0 {
		return nil
	}
	targetSvc := svc.DeepCopy()
	// Clear LoadBalancerStatus so kube-proxy removes DNAT rules and the BGP speaker
	// (which gates on bgp-vip / allow-shared-ip) withdraws the route on next reconcile.
	// Do NOT remove BgpVipAnnotation — it is user-managed; if the VIP CR is re-created
	// the Service must be able to auto-rebind.
	targetSvc.Status.LoadBalancer = v1.LoadBalancerStatus{}
	if _, err := c.config.KubeClient.CoreV1().Services(targetSvc.Namespace).UpdateStatus(
		context.Background(), targetSvc, metav1.UpdateOptions{},
	); err != nil {
		klog.Errorf("failed to cleanup service %s/%s status for bgp-lb-vip: %v", svc.Namespace, svc.Name, err)
		return err
	}
	return nil
}

// cleanBgpLbVipService is called on Service deletion.
// The VIP CR lifecycle is managed independently by the user; no action needed here.
// (The Service object is being deleted so kube-proxy will remove its forwarding rules
// automatically. The VIP CR remains reusable by another Service.)
func (c *Controller) cleanBgpLbVipService(svc *v1.Service) error {
	if len(svc.Status.LoadBalancer.Ingress) == 0 {
		return nil
	}
	klog.Infof("clean bgp-lb-vip for deleted service %s/%s", svc.Namespace, svc.Name)
	// The Service object is being deleted, so kube-proxy will remove Service-related
	// forwarding state from the deleted Service spec/status. The VIP CR is managed
	// separately and remains reusable by another Service.
	return nil
}
