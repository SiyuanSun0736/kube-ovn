package speaker

// This file implements EIP BGP announcement via node route mode.
//
// Network plane: IptablesEIP (EIP plane)
//   EIPs are floating IPs assigned to vpc-nat-gw pods via iptables DNAT/SNAT rules.
//   They live in the NAT gateway pod's network namespace, NOT on a node interface.
//   In node-route-eip-mode the speaker runs on the HOST network of the node where the
//   vpc-nat-gw pod is scheduled, then announces the EIP prefix so upstream routers can
//   reach the NAT gateway pod via the node's underlay IP.
//
// Orthogonal plane: Service VIP (kube-proxy / kube-ipvs0)
//   When --enable-lb-svc-announce is also set, the speaker
//   additionally announces LoadBalancer ingress IPs from Service.status.loadBalancer.ingress.
//   These VIPs are programmed on the kube-ipvs0 virtual interface on every node by
//   kube-proxy/IPVS — a completely separate network plane from IptablesEIP.
//   The two planes do NOT conflict and CAN be enabled simultaneously.
//
// Mode: --node-route-eip-mode (mutually exclusive with --nat-gw-mode)
//   --nat-gw-mode:         speaker runs INSIDE the NAT gateway pod (overlay).
//   --node-route-eip-mode: speaker runs on the NODE host network (underlay).
//   Both announce IptablesEIP resources; running both would create conflicting routes.
//
// Control flow:
//  1. Speaker starts in NodeRouteEIPMode with host network enabled
//  2. On startup: enqueue all existing Ready EIPs for route recovery
//  3. On EIP add/update event: check if the associated vpc-nat-gw pod runs on this node
//  4. If yes and EIP is ready: announce the EIP via BGP
//  5. On EIP delete: withdraw the EIP from BGP announcement
//  6. Periodic reconcile (syncNodeRouteEIPs) also merges Service VIP prefixes
//     (expectedBgpLbServiceEip) if --enable-lb-svc-announce is enabled.
//
// Prerequisites:
//   - Speaker must run in host network mode with NodeRouteEIPMode enabled
//   - NodeName must be configured (via --node-name or NODE_NAME env)
//   - vpc-nat-gw pods must have proper labels for identification

import (
	"errors"
	"fmt"
	"net"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// initNodeRouteEIPMode initializes the EIP informer handlers for node route mode.
// This should be called during controller initialization when NodeRouteEIPMode is enabled.
func (c *Controller) initNodeRouteEIPMode() {
	if c.eipQueue == nil {
		c.eipQueue = workqueue.NewTypedRateLimitingQueueWithConfig(
			workqueue.DefaultTypedControllerRateLimiter[string](),
			workqueue.TypedRateLimitingQueueConfig[string]{Name: "NodeRouteEIP"},
		)
	}

	// Register EIP event handlers for node route mode
	eipInformer := c.kubeovnInformerFactory.Kubeovn().V1().IptablesEIPs().Informer()
	_, _ = eipInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.enqueueAddNodeRouteEIP,
		UpdateFunc: c.enqueueUpdateNodeRouteEIP,
		DeleteFunc: c.enqueueDeleteNodeRouteEIP,
	})
}

// enqueueAddNodeRouteEIP handles add events for IptablesEIP in node route mode.
// This is primarily for speaker restart recovery: when the speaker restarts,
// the informer triggers Add events for all existing resources. Ready EIPs
// need to be re-processed to restore their BGP announcements.
func (c *Controller) enqueueAddNodeRouteEIP(obj any) {
	eip, ok := obj.(*kubeovnv1.IptablesEIP)
	if !ok {
		klog.Errorf("expected IptablesEIP but got %T", obj)
		return
	}

	// Only enqueue ready EIPs - non-ready EIPs will be processed when they become ready
	if !eip.Status.Ready {
		return
	}

	klog.V(3).Infof("enqueue add iptables-eip %s for BGP route recovery", eip.Name)
	c.eipQueue.Add(eip.Name)
}

// enqueueUpdateNodeRouteEIP handles update events for IptablesEIP in node route mode.
// This processes EIPs that become ready, become non-ready, or have their NAT gateway changed.
func (c *Controller) enqueueUpdateNodeRouteEIP(_, newObj any) {
	eip, ok := newObj.(*kubeovnv1.IptablesEIP)
	if !ok {
		klog.Errorf("expected IptablesEIP but got %T", newObj)
		return
	}

	// Skip EIPs that are being deleted
	if eip.DeletionTimestamp != nil {
		return
	}

	// Always enqueue on update to handle:
	// - EIP becomes ready: announce route
	// - EIP becomes non-ready: withdraw route
	// - BGP annotation changed: add/remove route
	// - NatGwDp changed: update route announcement
	klog.V(3).Infof("enqueue update iptables-eip %s", eip.Name)
	c.eipQueue.Add(eip.Name)
}

// enqueueDeleteNodeRouteEIP handles delete events for IptablesEIP in node route mode.
// When an EIP is deleted, we withdraw its BGP announcement immediately (not queued).
func (c *Controller) enqueueDeleteNodeRouteEIP(obj any) {
	var eip *kubeovnv1.IptablesEIP
	switch t := obj.(type) {
	case *kubeovnv1.IptablesEIP:
		eip = t
	case cache.DeletedFinalStateUnknown:
		e, ok := t.Obj.(*kubeovnv1.IptablesEIP)
		if !ok {
			klog.Warningf("unexpected object type in DeletedFinalStateUnknown: %T", t.Obj)
			return
		}
		eip = e
	default:
		klog.Warningf("unexpected object type: %T", obj)
		return
	}

	klog.V(3).Infof("withdrawing routes for deleted iptables-eip %s", eip.Name)
	c.withdrawEIPRoutes(eip)
}

// withdrawEIPRoutes withdraws BGP routes for an EIP.
// This is called when an EIP is deleted, becomes non-ready, loses BGP annotation,
// or when the NAT gateway pod moves to another node.
func (c *Controller) withdrawEIPRoutes(eip *kubeovnv1.IptablesEIP) {
	var errs []error
	var withdrawn []string
	for _, ip := range []string{eip.Spec.V4ip, eip.Spec.V6ip} {
		if ip == "" {
			continue
		}
		if c.config.EnableTcBpfRedirect {
			if err := c.tcManager.DeleteEIP(net.ParseIP(ip)); err != nil {
				errs = append(errs, err)
			}
		}
		if !c.isRouteAnnounced(ip) {
			klog.V(3).Infof("BGP route for EIP %s not announced, skipping withdraw", ip)
			continue
		}
		if err := c.delRoute(ip); err != nil {
			klog.Errorf("failed to withdraw BGP route for EIP %s: %v", ip, err)
			errs = append(errs, err)
		} else {
			withdrawn = append(withdrawn, ip)
		}
	}

	if len(withdrawn) > 0 {
		klog.Infof("withdrawn BGP routes for iptables-eip %s: %v", eip.Name, withdrawn)
	}
	if len(errs) > 0 {
		klog.Errorf("errors withdrawing BGP routes for EIP %s: %v", eip.Name, errors.Join(errs...))
	}
}

// runNodeRouteEIPWorker runs the worker for processing EIP events in node route mode.
func (c *Controller) runNodeRouteEIPWorker() {
	for c.processNextNodeRouteEIPItem() {
	}
}

// processNextNodeRouteEIPItem processes the next EIP work item from the queue.
func (c *Controller) processNextNodeRouteEIPItem() bool {
	key, shutdown := c.eipQueue.Get()
	if shutdown {
		return false
	}

	err := func(key string) error {
		defer c.eipQueue.Done(key)
		if err := c.handleAddOrUpdateNodeRouteEIP(key); err != nil {
			c.eipQueue.AddRateLimited(key)
			return fmt.Errorf("error processing EIP %q: %w, requeuing", key, err)
		}
		c.eipQueue.Forget(key)
		return nil
	}(key)
	if err != nil {
		klog.Error(err)
	}
	return true
}

// handleAddOrUpdateNodeRouteEIP processes an EIP add or update event.
// It checks if the associated vpc-nat-gw pod runs on the local node and
// announces the EIP via BGP if so, or withdraws it if conditions are not met.
func (c *Controller) handleAddOrUpdateNodeRouteEIP(eipName string) error {
	klog.V(3).Infof("handling add/update iptables-eip %s in node route mode", eipName)

	eip, err := c.eipLister.Get(eipName)
	if err != nil {
		// EIP was deleted, skip processing
		klog.V(3).Infof("iptables-eip %s not found, may have been deleted", eipName)
		return nil
	}

	// Skip non-ready EIPs - they will be processed when they become ready
	// Periodic reconcile (syncNodeRouteEIPs) will clean up any stale routes
	if !eip.Status.Ready {
		return nil
	}

	// Check if BGP annotation is enabled for this EIP
	if eip.Annotations[util.BgpAnnotation] != "true" {
		klog.V(3).Infof("iptables-eip %s does not have BGP annotation, skipping", eipName)
		// Withdraw any existing routes for this EIP (in case annotation was removed)
		c.withdrawEIPRoutes(eip)
		return nil
	}

	// Check if the NAT gateway pod is running on the local node
	if !c.hasNatGwPodOnLocalNode(eip) {
		klog.V(3).Infof("NAT GW pod for iptables-eip %s not on local node %s, withdrawing routes",
			eipName, c.config.NodeName)
		// Withdraw any existing routes for this EIP (in case pod moved to another node)
		c.withdrawEIPRoutes(eip)
		return nil
	}

	// Announce routes only if not already announced (idempotent)
	var errs []error
	var redirectTarget *TCRedirectTarget
	if c.config.EnableTcBpfRedirect {
		var err error
		redirectTarget, _, err = c.resolveEipRedirectTarget(eip)
		if err != nil {
			errs = append(errs, err)
		}
	}

	var announced []string
	for _, ip := range []string{eip.Spec.V4ip, eip.Spec.V6ip} {
		if ip == "" {
			continue
		}
		if c.config.EnableTcBpfRedirect {
			if redirectTarget == nil || !validTCRedirectTarget(*redirectTarget) {
				klog.V(3).Infof("deferring BGP route announcement for EIP %s until TC redirect target is ready", ip)
				continue
			}
			if err := c.tcManager.AddEIP(net.ParseIP(ip), *redirectTarget); err != nil {
				klog.Errorf("failed to add EIP %s to TC BPF map: %v", ip, err)
				errs = append(errs, err)
				continue
			}
		}
		if c.isRouteAnnounced(ip) {
			klog.V(3).Infof("BGP route for EIP %s already announced, skipping", ip)
			continue
		}
		if err := c.addRoute(ip); err != nil {
			klog.Errorf("failed to announce BGP route for EIP %s: %v", ip, err)
			errs = append(errs, err)
		} else {
			announced = append(announced, ip)
		}
	}

	if len(announced) > 0 {
		klog.Infof("announced BGP routes for iptables-eip %s: %v", eipName, announced)
	}

	return errors.Join(errs...)
}

// hasNatGwPodOnLocalNode checks if the NAT gateway pod for this EIP is running on the local node.
// NAT gateway pod name follows the pattern: vpc-nat-gw-{natGwDp}-0
// Returns true only if the pod exists, is Running, and is scheduled on the local node.
func (c *Controller) hasNatGwPodOnLocalNode(eip *kubeovnv1.IptablesEIP) bool {
	if eip.Spec.NatGwDp == "" {
		klog.Errorf("iptables-eip %s has empty NatGwDp field", eip.Name)
		return false
	}

	podName := util.GenNatGwPodName(eip.Spec.NatGwDp)
	// Use gwPodsLister which watches pods in VpcNatGwNamespace
	pod, err := c.gwPodsLister.Pods(c.config.VpcNatGwNamespace).Get(podName)
	if err != nil {
		klog.V(3).Infof("failed to get NAT GW pod %s/%s: %v", c.config.VpcNatGwNamespace, podName, err)
		return false
	}

	// Only announce routes for Running pods on local node
	return pod.Spec.NodeName == c.config.NodeName && isPodRunning(pod)
}

// resolveEipRedirectTarget 解析当前 EIP 在本地网关 Pod 的重定向直通目标。
// 返回值：
//   - target: 解析成功的直通目标（解析失败或不适用时为 nil）
//   - retainStale: 当活跃 Pod 仅因元数据短暂就绪延迟（如 Multus NetworkStatus 写入中）时为 true，
//     指示声明式保留存量 TC BPF 映射，避免流量被瞬时中断
//   - err: 解析异常错误
func (c *Controller) resolveEipRedirectTarget(eip *kubeovnv1.IptablesEIP) (*TCRedirectTarget, bool, error) {
	podName := util.GenNatGwPodName(eip.Spec.NatGwDp)
	pod, err := c.gwPodsLister.Pods(c.config.VpcNatGwNamespace).Get(podName)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			// 非 NotFound 的瞬时缓存或 API 抖动：声明式保留存量映射以防网络中断
			return nil, true, fmt.Errorf("failed to get NAT GW pod %s: %w", podName, err)
		}
		// Pod 已不存在，清理对应的 UID 缓存
		c.gwTargetCacheMu.Lock()
		delete(c.gwTargetCache, eip.Spec.NatGwDp)
		c.gwTargetCacheMu.Unlock()
		return nil, false, err
	}

	// 若 Pod 处于删除中或非 Running 阶段，属于确定性失效：绝不保留旧映射，并清理 UID 缓存
	if isPodDeleting(pod) || !isPodRunning(pod) {
		klog.Warningf("NAT GW pod %s for EIP %s is terminating or not running, dropping redirect target", podName, eip.Name)
		c.gwTargetCacheMu.Lock()
		delete(c.gwTargetCache, eip.Spec.NatGwDp)
		c.gwTargetCacheMu.Unlock()
		return nil, false, nil
	}

	// 1. 尝试从长效 UID 缓存读取（只要 Pod.UID 未变，网卡 ifindex 和 MAC 恒定，无需重复查询 Netlink）
	c.gwTargetCacheMu.RLock()
	cached, ok := c.gwTargetCache[eip.Spec.NatGwDp]
	c.gwTargetCacheMu.RUnlock()
	if ok && cached.podUID == pod.UID {
		return &cached.target, false, nil
	}

	// 2. 缓存未命中或 Pod 重建（UID 变更），实时解析网卡目标
	target, err := c.getNatGwTCRedirectTarget(pod, eip)
	if err != nil {
		// Pod 正常活跃，但网卡/子网元数据暂未就绪（如 NetworkStatus 写入延迟），保留存量映射
		klog.Warningf("failed to resolve redirect target for EIP %s while pod %s is active: %v", eip.Name, podName, err)
		return nil, true, err
	}

	// 3. 解析成功，存入长效 UID 缓存
	c.gwTargetCacheMu.Lock()
	c.gwTargetCache[eip.Spec.NatGwDp] = gwTargetCacheEntry{
		podUID: pod.UID,
		target: target,
	}
	c.gwTargetCacheMu.Unlock()

	return &target, false, nil
}

// syncNodeRouteEIPs performs a full reconciliation of all EIPs in node route mode.
// This method finds all EIPs associated with local NAT gateway pods and announces them.
// It also withdraws any EIPs that should no longer be announced.
//
// When EnableLbSvcAnnounce is also enabled, Service VIP prefixes are
// merged into the expected set so that node-route-eip mode and Service VIP announcements
// can operate concurrently on the same speaker instance.
func (c *Controller) syncNodeRouteEIPs() error {
	expectedPrefixes := make(prefixMap)
	expectedEipTargets := make(map[string]TCRedirectTarget)
	var retainedEIPs []string
	var targetErrors []error

	// EIP plane: announce IptablesEIP addresses for local NAT gateway pods.
	eips, err := c.eipLister.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("failed to list EIPs: %w", err)
	}

	type targetResolution struct {
		target      *TCRedirectTarget
		retainStale bool
		err         error
	}
	targetCache := make(map[string]targetResolution)

	for _, eip := range eips {
		// Only process ready EIPs with BGP annotation
		if eip.Annotations[util.BgpAnnotation] != "true" || !eip.Status.Ready {
			continue
		}

		// Only announce EIPs for local NAT gateway pods
		if !c.hasNatGwPodOnLocalNode(eip) {
			continue
		}

		var redirectTarget *TCRedirectTarget
		var retainStale bool
		if c.config.EnableTcBpfRedirect {
			cacheKey := eip.Spec.NatGwDp
			if cached, ok := targetCache[cacheKey]; ok {
				redirectTarget, retainStale = cached.target, cached.retainStale
				if cached.err != nil && redirectTarget == nil {
					targetErrors = append(targetErrors, cached.err)
				}
			} else {
				var err error
				redirectTarget, retainStale, err = c.resolveEipRedirectTarget(eip)
				targetCache[cacheKey] = targetResolution{target: redirectTarget, retainStale: retainStale, err: err}
				if err != nil && redirectTarget == nil {
					targetErrors = append(targetErrors, err)
				}
			}
		}

		for _, ip := range []string{eip.Spec.V4ip, eip.Spec.V6ip} {
			if ip != "" {
				addExpectedPrefix(ip, expectedPrefixes)
				if redirectTarget != nil {
					expectedEipTargets[ip] = *redirectTarget
				} else if retainStale {
					retainedEIPs = append(retainedEIPs, ip)
				}
			}
		}
	}

	if c.config.EnableTcBpfRedirect {
		err := c.tcManager.Attach(c.config.TcBpfInterface)
		if err == nil {
			err = c.tcManager.SyncEIPs(expectedEipTargets, retainedEIPs...)
		}
		if err != nil {
			// TC eBPF 挂载或同步失败时不应撤销 BGP 路由；降级回退到宿主机内核协议栈以优先保障服务连续性
			klog.Errorf("failed to sync EIP redirection targets, falling back to kernel datapath: %v", err)
			targetErrors = append(targetErrors, err)
		}
	}

	// Service VIP plane: merge LoadBalancer ingress IPs when Service VIP announcement is
	// enabled alongside node-route-eip mode. The two planes are orthogonal:
	//   - EIP plane  → IptablesEIP resources, NAT gateway traffic path.
	//   - Service VIP plane → Service.status.loadBalancer.ingress on kube-ipvs0.
	// Merging both into expectedPrefixes ensures reconcileRoutes does not withdraw
	// Service VIP routes that are legitimately announced.
	if c.config.EnableLbSvcAnnounce {
		node := c.getLocalNode()
		if shouldAnnounceLbSvcIngressOnNode(node) {
			services, err := c.servicesLister.List(labels.Everything())
			if err != nil {
				return fmt.Errorf("failed to list services for bgp-lb-vip: %w", err)
			}
			expectedBgpLbServiceEip := make(prefixMap)
			collectSvcBgpPrefixes(services, node, expectedBgpLbServiceEip)
			mergePrefixMap(expectedBgpLbServiceEip, expectedPrefixes)
		}
	}

	return errors.Join(c.reconcileRoutes(expectedPrefixes), errors.Join(targetErrors...))
}

// startNodeRouteEIPWorkers starts the worker goroutines for processing EIP events.
// This should be called from the controller's Run method.
func (c *Controller) startNodeRouteEIPWorkers(stopCh <-chan struct{}, workers int) {
	klog.Infof("starting %d node route EIP workers", workers)
	for i := 0; i < workers; i++ {
		go wait.Until(c.runNodeRouteEIPWorker, time.Second, stopCh)
	}
}

// shutdownNodeRouteEIPWorkers shuts down the EIP work queue.
func (c *Controller) shutdownNodeRouteEIPWorkers() {
	if c.eipQueue != nil {
		c.eipQueue.ShutDown()
	}
}

// enqueueAllReadyEIPs enqueues all ready EIPs on speaker startup for route recovery.
// This ensures that after a speaker restart, all local EIPs are re-announced.
func (c *Controller) enqueueAllReadyEIPs() error {
	eips, err := c.eipLister.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("failed to list EIPs for startup recovery: %w", err)
	}

	count := 0
	for _, eip := range eips {
		// Only enqueue ready EIPs with BGP annotation
		if !eip.Status.Ready || eip.Annotations[util.BgpAnnotation] != "true" {
			continue
		}
		klog.V(3).Infof("enqueue ready iptables-eip %s on startup", eip.Name)
		c.eipQueue.Add(eip.Name)
		count++
	}

	klog.Infof("enqueued %d ready EIPs for startup recovery", count)
	return nil
}

// ReconcileNodeRouteEIPs performs periodic reconciliation in node route EIP mode.
// Called by the controller's Reconcile method to ensure BGP state is consistent.
func (c *Controller) ReconcileNodeRouteEIPs() {
	if err := c.syncNodeRouteEIPs(); err != nil {
		klog.Errorf("failed to reconcile node route EIPs: %v", err)
	}
}
