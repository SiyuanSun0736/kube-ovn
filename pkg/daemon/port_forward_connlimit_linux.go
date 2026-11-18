package daemon

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

type portForwardConnLimitRule struct {
	ip        string
	port      int32
	protocol  string
	connLimit int
}

func (r portForwardConnLimitRule) key() string {
	return fmt.Sprintf("%s:%d/%s", r.ip, r.port, r.protocol)
}

func buildPortForwardConnLimitRuleSpec(ip string, port int32, limit int) []string {
	return []string{
		"-d", ip + "/32",
		"-p", "tcp",
		"-m", "tcp", "--dport", strconv.Itoa(int(port)),
		"--tcp-flags", "FIN,SYN,RST,ACK", "SYN",
		"-m", "connlimit",
		"--connlimit-above", strconv.Itoa(limit),
		"--connlimit-mask", "0",
		"-j", "DROP",
	}
}

func (c *Controller) getExternalInterface() string {
	if c.config == nil {
		return ""
	}
	if iface := c.config.GetTunnelIface(); iface != "" {
		return iface
	}
	return c.config.Iface
}

func (c *Controller) ensurePortForwardConnLimitChain(protocol string) error {
	ipt := c.iptables[protocol]
	if ipt == nil {
		return nil
	}

	chainExists, err := ipt.ChainExists(util.Mangle, util.OvnPortForwardConnLimit)
	if err != nil {
		klog.Errorf("failed to check existence of chain %s in table %s: %v", util.OvnPortForwardConnLimit, util.Mangle, err)
		return err
	}
	if !chainExists {
		if err := ipt.NewChain(util.Mangle, util.OvnPortForwardConnLimit); err != nil {
			klog.Errorf("failed to create chain %s in table %s: %v", util.OvnPortForwardConnLimit, util.Mangle, err)
			return err
		}
		klog.Infof("created chain %s in table %s", util.OvnPortForwardConnLimit, util.Mangle)
	}

	// Ensure jump from PREROUTING to OVN-PF-CONNLIMIT at top
	jumpRule := []string{"-j", util.OvnPortForwardConnLimit}
	exists, err := ipt.Exists(util.Mangle, util.Prerouting, jumpRule...)
	if err != nil {
		klog.Errorf("failed to check jump to %s in %s PREROUTING: %v", util.OvnPortForwardConnLimit, util.Mangle, err)
		return err
	}
	if !exists {
		if err := ipt.Insert(util.Mangle, util.Prerouting, 1, jumpRule...); err != nil {
			klog.Errorf("failed to insert jump to %s in %s PREROUTING: %v", util.OvnPortForwardConnLimit, util.Mangle, err)
			return err
		}
		klog.Infof("inserted jump to %s in %s PREROUTING", util.OvnPortForwardConnLimit, util.Mangle)
	}

	// Gate rule: packets not coming from external interface return immediately (PR #143 style)
	extIface := c.getExternalInterface()
	var gateRule []string
	if extIface != "" {
		gateRule = []string{"!", "-i", extIface, "-j", "RETURN"}
	} else {
		gateRule = []string{"-i", "lo", "-j", "RETURN"}
	}

	ruleExists, err := ipt.Exists(util.Mangle, util.OvnPortForwardConnLimit, gateRule...)
	if err != nil {
		klog.Errorf("failed to check gate rule %+v in %s: %v", gateRule, util.OvnPortForwardConnLimit, err)
		return err
	}
	if !ruleExists {
		if err := ipt.Insert(util.Mangle, util.OvnPortForwardConnLimit, 1, gateRule...); err != nil {
			klog.Errorf("failed to insert gate rule %+v in %s: %v", gateRule, util.OvnPortForwardConnLimit, err)
			return err
		}
		klog.Infof("inserted gate rule %+v in %s", gateRule, util.OvnPortForwardConnLimit)
	}

	return nil
}

func (c *Controller) parsePortForwardConnLimitRules(svc *corev1.Service) []portForwardConnLimitRule {
	if svc == nil || svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return nil
	}

	limitStr := svc.Annotations[util.PortForwardConnLimitAnnotation]
	if limitStr == "" {
		limitStr = svc.Annotations[util.TCPConnLimitAnnotation]
	}

	var connLimit int
	if limitStr != "" {
		parsed, err := strconv.Atoi(strings.TrimSpace(limitStr))
		if err != nil {
			klog.Warningf("invalid port-forward connlimit annotation %q on service %s/%s: %v", limitStr, svc.Namespace, svc.Name, err)
			return nil
		}
		if parsed <= 0 {
			// Limit <= 0 explicitly disables connlimit
			return nil
		}
		connLimit = parsed
	} else {
		// No connlimit annotation specified on service
		return nil
	}

	ips := extractServicePublicIPs(svc)
	if len(ips) == 0 {
		return nil
	}

	var tcpPorts []int32
	for _, port := range svc.Spec.Ports {
		if port.Protocol == corev1.ProtocolTCP && port.Port > 0 {
			tcpPorts = append(tcpPorts, port.Port)
		}
	}
	if len(tcpPorts) == 0 {
		return nil
	}

	rules := make([]portForwardConnLimitRule, 0, len(ips)*len(tcpPorts))
	for _, ip := range ips {
		for _, port := range tcpPorts {
			rules = append(rules, portForwardConnLimitRule{
				ip:        ip,
				port:      port,
				protocol:  "tcp",
				connLimit: connLimit,
			})
		}
	}
	return rules
}

func extractServicePublicIPs(svc *corev1.Service) []string {
	ipSet := make(map[string]struct{})

	addIP := func(raw string) {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			return
		}
		if parsed := net.ParseIP(trimmed); parsed != nil {
			ipSet[parsed.String()] = struct{}{}
		}
	}

	// 1. metallb loadBalancerIPs annotation (supports comma-separated list)
	if lbIPs := svc.Annotations[util.MetalLBLoadBalancerIPsAnnotation]; lbIPs != "" {
		for ip := range strings.SplitSeq(lbIPs, ",") {
			addIP(ip)
		}
	}

	// 2. metallb allow-shared-ip annotation
	if sharedIP := svc.Annotations[util.MetalLBAllowSharedIPAnnotation]; sharedIP != "" {
		addIP(sharedIP)
	}

	// 3. spec.loadBalancerIP
	if svc.Spec.LoadBalancerIP != "" {
		addIP(svc.Spec.LoadBalancerIP)
	}

	// 4. status.loadBalancer.ingress
	for _, ingress := range svc.Status.LoadBalancer.Ingress {
		if ingress.IP != "" {
			addIP(ingress.IP)
		}
	}

	result := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		result = append(result, ip)
	}
	return result
}

func (c *Controller) reconcilePortForwardConnLimit(oldService, newService *corev1.Service) error {
	oldRules := c.parsePortForwardConnLimitRules(oldService)
	newRules := c.parsePortForwardConnLimitRules(newService)

	if len(oldRules) == 0 && len(newRules) == 0 {
		return nil
	}

	oldMap := make(map[string]portForwardConnLimitRule, len(oldRules))
	for _, r := range oldRules {
		oldMap[r.key()] = r
	}

	newMap := make(map[string]portForwardConnLimitRule, len(newRules))
	for _, r := range newRules {
		newMap[r.key()] = r
	}

	// Rules to delete: in old but not in new, or limit changed
	for key, oldRule := range oldMap {
		newRule, exists := newMap[key]
		if !exists || newRule.connLimit != oldRule.connLimit {
			if err := c.deletePortForwardConnLimitRule(oldRule); err != nil {
				klog.Errorf("failed to delete port-forward connlimit rule for %s: %v", key, err)
			}
		}
	}

	// Rules to add: in new but not in old, or limit changed
	for key, newRule := range newMap {
		oldRule, exists := oldMap[key]
		if !exists || oldRule.connLimit != newRule.connLimit {
			if err := c.addPortForwardConnLimitRule(newRule); err != nil {
				klog.Errorf("failed to add port-forward connlimit rule for %s: %v", key, err)
				return err
			}
		}
	}

	return nil
}

func (c *Controller) addPortForwardConnLimitRule(rule portForwardConnLimitRule) error {
	proto := getIPProtocol(rule.ip)
	ipt := c.iptables[proto]
	if ipt == nil {
		return nil
	}

	if err := c.ensurePortForwardConnLimitChain(proto); err != nil {
		klog.Errorf("failed to ensure chain %s for rule %s:%d: %v", util.OvnPortForwardConnLimit, rule.ip, rule.port, err)
		return err
	}

	ruleSpec := buildPortForwardConnLimitRuleSpec(rule.ip, rule.port, rule.connLimit)
	exists, err := ipt.Exists(util.Mangle, util.OvnPortForwardConnLimit, ruleSpec...)
	if err != nil {
		klog.Errorf("failed to check existence of connlimit rule for %s:%d: %v", rule.ip, rule.port, err)
		return err
	}
	if !exists {
		klog.Infof("adding port-forward connlimit rule: -d %s --dport %d limit %d", rule.ip, rule.port, rule.connLimit)
		if err := ipt.Append(util.Mangle, util.OvnPortForwardConnLimit, ruleSpec...); err != nil {
			klog.Errorf("failed to append connlimit rule for %s:%d: %v", rule.ip, rule.port, err)
			return err
		}
	}

	return nil
}

func (c *Controller) deletePortForwardConnLimitRule(rule portForwardConnLimitRule) error {
	proto := getIPProtocol(rule.ip)
	ipt := c.iptables[proto]
	if ipt == nil {
		return nil
	}

	chainExists, err := ipt.ChainExists(util.Mangle, util.OvnPortForwardConnLimit)
	if err != nil || !chainExists {
		return nil
	}

	// 1. If rule with exact limit exists, delete it directly
	exactRuleSpec := buildPortForwardConnLimitRuleSpec(rule.ip, rule.port, rule.connLimit)
	exists, err := ipt.Exists(util.Mangle, util.OvnPortForwardConnLimit, exactRuleSpec...)
	if err == nil && exists {
		klog.Infof("deleting exact port-forward connlimit rule: -d %s --dport %d limit %d", rule.ip, rule.port, rule.connLimit)
		if err := ipt.Delete(util.Mangle, util.OvnPortForwardConnLimit, exactRuleSpec...); err != nil {
			klog.Errorf("failed to delete rule for %s:%d: %v", rule.ip, rule.port, err)
		}
	}

	// 2. Also delete any matching IP + Port rules regardless of limit threshold (handles limit updates)
	return c.deletePortForwardConnLimitByIPPort(proto, rule.ip, rule.port)
}

func (c *Controller) deletePortForwardConnLimitByIPPort(proto, ip string, port int32) error {
	ipt := c.iptables[proto]
	if ipt == nil {
		return nil
	}

	rules, err := ipt.List(util.Mangle, util.OvnPortForwardConnLimit)
	if err != nil {
		klog.Errorf("failed to list rules in chain %s: %v", util.OvnPortForwardConnLimit, err)
		return nil
	}

	ipPattern := ip
	portPattern := strconv.Itoa(int(port))

	for _, rawRule := range rules {
		// Example: -A OVN-PF-CONNLIMIT -d 1.2.3.4/32 -p tcp -m tcp --dport 20022 --tcp-flags SYN,ACK SYN ...
		if strings.Contains(rawRule, ipPattern) && strings.Contains(rawRule, portPattern) && strings.Contains(rawRule, "connlimit") {
			tokens := strings.Fields(rawRule)
			if len(tokens) > 2 && tokens[0] == "-A" && tokens[1] == util.OvnPortForwardConnLimit {
				ruleSpec := tokens[2:]
				klog.Infof("deleting stale port-forward connlimit rule: %s", strings.Join(ruleSpec, " "))
				if err := ipt.Delete(util.Mangle, util.OvnPortForwardConnLimit, ruleSpec...); err != nil {
					klog.Errorf("failed to delete matched rule %q: %v", rawRule, err)
				}
			}
		}
	}
	return nil
}

func getIPProtocol(ip string) string {
	if util.CheckProtocol(ip) == kubeovnv1.ProtocolIPv6 {
		return kubeovnv1.ProtocolIPv6
	}
	return kubeovnv1.ProtocolIPv4
}
