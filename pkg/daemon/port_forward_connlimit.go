package daemon

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

type portForwardKey struct {
	ip   string
	port int32
}

type portForwardConnLimitRule struct {
	ip        string
	port      int32
	connLimit int
}

func (r portForwardConnLimitRule) key() portForwardKey {
	return portForwardKey{ip: r.ip, port: r.port}
}

type iptablesRuleEntry struct {
	comment string
	limit   int
	raw     string
}

func iptablesRuleComment(ip string, port int32) string {
	return fmt.Sprintf("ovn-pf-%s-%d", ip, port)
}

func execIptables(protocol string, args ...string) ([]byte, error) {
	binary := "iptables"
	if protocol == kubeovnv1.ProtocolIPv6 {
		binary = "ip6tables"
	}
	cmd := exec.Command(binary, args...) // #nosec G204
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	return cmd.CombinedOutput()
}

func buildIptablesRuleSpec(rule portForwardConnLimitRule) []string {
	cidr := "/32"
	if util.CheckProtocol(rule.ip) == kubeovnv1.ProtocolIPv6 {
		cidr = "/128"
	}
	comment := iptablesRuleComment(rule.ip, rule.port)
	return []string{
		"-d", fmt.Sprintf("%s%s", rule.ip, cidr),
		"-p", "tcp",
		"-m", "tcp",
		"--dport", strconv.Itoa(int(rule.port)),
		"--tcp-flags", "FIN,SYN,RST,ACK", "SYN",
		"-m", "connlimit",
		"--connlimit-above", strconv.Itoa(rule.connLimit),
		"--connlimit-mask", "0",
		"-m", "comment",
		"--comment", comment,
		"-j", "DROP",
	}
}

func (c *Controller) enabledProtocols() []string {
	if c == nil || c.protocol == "" || c.protocol == kubeovnv1.ProtocolDual {
		return []string{kubeovnv1.ProtocolIPv4, kubeovnv1.ProtocolIPv6}
	}
	if c.protocol == kubeovnv1.ProtocolIPv6 {
		return []string{kubeovnv1.ProtocolIPv6}
	}
	return []string{kubeovnv1.ProtocolIPv4}
}

func (c *Controller) deleteIptablesRule(protocol string, rule portForwardConnLimitRule) error {
	spec := buildIptablesRuleSpec(rule)
	args := append([]string{"-w", "5", "-t", "mangle", "-D", util.OvnPortForwardIptablesChain}, spec...)
	out, err := execIptables(protocol, args...)
	if err != nil {
		outStr := string(out)
		if strings.Contains(outStr, "does a matching rule exist") || strings.Contains(outStr, "No chain/target/match by that name") {
			return nil
		}
		klog.Errorf("failed to delete port-forward connlimit rule %s:%d (protocol %s): %v, out: %s",
			rule.ip, rule.port, protocol, err, outStr)
		return err
	}
	return nil
}

func (c *Controller) addIptablesRule(protocol string, rule portForwardConnLimitRule) error {
	spec := buildIptablesRuleSpec(rule)
	args := append([]string{"-w", "5", "-t", "mangle", "-A", util.OvnPortForwardIptablesChain}, spec...)
	out, err := execIptables(protocol, args...)
	if err != nil {
		klog.Errorf("failed to add port-forward connlimit rule %s:%d (protocol %s): %v, out: %s",
			rule.ip, rule.port, protocol, err, string(out))
		return err
	}
	return nil
}

func parseIptablesChainRules(output string) []iptablesRuleEntry {
	var entries []iptablesRuleEntry
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ":") || strings.HasPrefix(line, "-N ") {
			continue
		}

		if !strings.HasPrefix(line, "-A "+util.OvnPortForwardIptablesChain) {
			continue
		}

		var comment string
		if idx := strings.Index(line, "--comment "); idx != -1 {
			rem := strings.TrimSpace(line[idx+len("--comment "):])
			if len(rem) > 0 {
				if rem[0] == '"' || rem[0] == '\'' {
					quote := rem[0]
					if end := strings.IndexByte(rem[1:], quote); end != -1 {
						comment = rem[1 : end+1]
					}
				} else {
					fields := strings.Fields(rem)
					if len(fields) > 0 {
						comment = fields[0]
					}
				}
			}
		}

		var limit int
		if idx := strings.Index(line, "--connlimit-above "); idx != -1 {
			rem := strings.TrimSpace(line[idx+len("--connlimit-above "):])
			fields := strings.Fields(rem)
			if len(fields) > 0 {
				if l, err := strconv.Atoi(fields[0]); err == nil {
					limit = l
				}
			}
		}

		entries = append(entries, iptablesRuleEntry{
			comment: comment,
			limit:   limit,
			raw:     line,
		})
	}
	return entries
}

func findIptablesRulesByComment(entries []iptablesRuleEntry, targetComment string) []iptablesRuleEntry {
	var res []iptablesRuleEntry
	for _, e := range entries {
		if e.comment == targetComment {
			res = append(res, e)
		}
	}
	return res
}

func (c *Controller) listIptablesRules(protocol string) ([]iptablesRuleEntry, error) {
	out, err := execIptables(protocol, "-w", "5", "-t", "mangle", "-S", util.OvnPortForwardIptablesChain)
	if err != nil {
		outStr := string(out)
		if strings.Contains(outStr, "No chain/target/match by that name") || strings.Contains(outStr, "does a matching rule exist") {
			return nil, fmt.Errorf("chain does not exist: %w", err)
		}
		return nil, fmt.Errorf("iptables -S failed: %w, out: %s", err, outStr)
	}
	return parseIptablesChainRules(string(out)), nil
}

func (c *Controller) ensureIptablesChain(protocol string) error {
	// Create chain if not exists (ignore error if already exists)
	_, _ = execIptables(protocol, "-w", "5", "-t", "mangle", "-N", util.OvnPortForwardIptablesChain)

	// Ensure jump in both PREROUTING and OUTPUT hooks
	for _, parent := range []string{"PREROUTING", "OUTPUT"} {
		if _, err := execIptables(protocol, "-w", "5", "-t", "mangle", "-C", parent, "-j", util.OvnPortForwardIptablesChain); err != nil {
			if out, insErr := execIptables(protocol, "-w", "5", "-t", "mangle", "-I", parent, "1", "-j", util.OvnPortForwardIptablesChain); insErr != nil {
				klog.Errorf("failed to insert %s jump in mangle %s (protocol %s): %v, out: %s",
					util.OvnPortForwardIptablesChain, parent, protocol, insErr, string(out))
				return insErr
			}
			klog.Infof("ensured %s jump in mangle %s for protocol %s", util.OvnPortForwardIptablesChain, parent, protocol)
		}
	}
	return nil
}

func (c *Controller) parsePortForwardConnLimitRules(svc *corev1.Service) []portForwardConnLimitRule {
	if svc == nil || svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return nil
	}

	limitStr := svc.Annotations[util.PortForwardConnLimitAnnotation]
	if limitStr == "" {
		return nil
	}

	parsed, err := strconv.Atoi(strings.TrimSpace(limitStr))
	if err != nil {
		klog.Warningf("invalid port-forward connlimit annotation %q on service %s/%s: %v", limitStr, svc.Namespace, svc.Name, err)
		return nil
	}
	if parsed <= 0 {
		return nil
	}
	connLimit := parsed

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
				connLimit: connLimit,
			})
		}
	}
	return rules
}

func extractServicePublicIPs(svc *corev1.Service) []string {
	ipSet := make(map[string]struct{})
	for _, ingress := range svc.Status.LoadBalancer.Ingress {
		if ingress.IP == "" {
			continue
		}
		if parsed := net.ParseIP(strings.TrimSpace(ingress.IP)); parsed != nil {
			ipSet[parsed.String()] = struct{}{}
		}
	}

	result := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		result = append(result, ip)
	}
	sort.Strings(result)
	return result
}

func (c *Controller) reconcilePortForwardConnLimit(oldService, newService *corev1.Service) error {
	if c != nil {
		c.pfConnLimitMutex.Lock()
		defer c.pfConnLimitMutex.Unlock()
	}

	// Re-read current Service from Lister to avoid applying stale event snapshot on retry
	if c != nil && c.servicesLister != nil {
		svcToLookup := newService
		if svcToLookup == nil {
			svcToLookup = oldService
		}
		if svcToLookup != nil {
			currentSvc, err := c.servicesLister.Services(svcToLookup.Namespace).Get(svcToLookup.Name)
			if err == nil && currentSvc != nil {
				newService = currentSvc
			} else if newService != nil {
				// Service was deleted from K8s while stale event was queued
				newService = nil
			}
		}
	}

	oldRules := c.parsePortForwardConnLimitRules(oldService)
	newRules := c.parsePortForwardConnLimitRules(newService)

	if len(oldRules) == 0 && len(newRules) == 0 {
		return nil
	}

	// Group by protocol (IPv4 vs IPv6)
	rulesByProtoOld := map[string][]portForwardConnLimitRule{
		kubeovnv1.ProtocolIPv4: nil,
		kubeovnv1.ProtocolIPv6: nil,
	}
	for _, r := range oldRules {
		proto := kubeovnv1.ProtocolIPv4
		if util.CheckProtocol(r.ip) == kubeovnv1.ProtocolIPv6 {
			proto = kubeovnv1.ProtocolIPv6
		}
		rulesByProtoOld[proto] = append(rulesByProtoOld[proto], r)
	}

	rulesByProtoNew := map[string][]portForwardConnLimitRule{
		kubeovnv1.ProtocolIPv4: nil,
		kubeovnv1.ProtocolIPv6: nil,
	}
	for _, r := range newRules {
		proto := kubeovnv1.ProtocolIPv4
		if util.CheckProtocol(r.ip) == kubeovnv1.ProtocolIPv6 {
			proto = kubeovnv1.ProtocolIPv6
		}
		rulesByProtoNew[proto] = append(rulesByProtoNew[proto], r)
	}

	svcName := ""
	if newService != nil {
		svcName = newService.Namespace + "/" + newService.Name
	} else if oldService != nil {
		svcName = oldService.Namespace + "/" + oldService.Name
	}

	for _, proto := range c.enabledProtocols() {
		pOld := rulesByProtoOld[proto]
		pNew := rulesByProtoNew[proto]
		if len(pOld) == 0 && len(pNew) == 0 {
			continue
		}

		if err := c.reconcileProtocolRules(proto, pOld, pNew, svcName); err != nil {
			klog.Errorf("failed to reconcile port-forward connlimit rules for protocol %s: %v", proto, err)
			return err
		}
	}

	return nil
}

func (c *Controller) getOtherServicesClaimedPortForwards(excludeSvc string) (map[portForwardKey]struct{}, error) {
	if c == nil || c.servicesLister == nil {
		return nil, nil
	}
	svcs, err := c.servicesLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("failed to list services to check claimed port-forwards: %v", err)
		return nil, err
	}
	claimed := make(map[portForwardKey]struct{})
	for _, svc := range svcs {
		fullName := svc.Namespace + "/" + svc.Name
		if fullName == excludeSvc {
			continue
		}
		for _, r := range c.parsePortForwardConnLimitRules(svc) {
			claimed[r.key()] = struct{}{}
		}
	}
	return claimed, nil
}

func (c *Controller) reconcileProtocolRules(protocol string, oldRules, newRules []portForwardConnLimitRule, currentSvcName string) error {
	oldMap := make(map[portForwardKey]portForwardConnLimitRule, len(oldRules))
	for _, r := range oldRules {
		oldMap[r.key()] = r
	}
	newMap := make(map[portForwardKey]portForwardConnLimitRule, len(newRules))
	for _, r := range newRules {
		newMap[r.key()] = r
	}

	// Determine deleted and added/updated rules
	toDeleteRules := make(map[portForwardKey]portForwardConnLimitRule)
	for key, oldRule := range oldMap {
		if _, exists := newMap[key]; !exists {
			toDeleteRules[key] = oldRule
		}
	}

	toAddOrUpdate := make(map[portForwardKey]portForwardConnLimitRule)
	for key, newRule := range newMap {
		oldRule, exists := oldMap[key]
		if !exists || oldRule.connLimit != newRule.connLimit {
			toAddOrUpdate[key] = newRule
		}
	}

	if len(toDeleteRules) == 0 && len(toAddOrUpdate) == 0 {
		return nil
	}

	if err := c.ensureIptablesChain(protocol); err != nil {
		return err
	}

	entries, err := c.listIptablesRules(protocol)
	if err != nil && !strings.Contains(err.Error(), "chain does not exist") {
		klog.Errorf("failed to list iptables rules for protocol %s: %v", protocol, err)
		return err
	}

	claimedByOthers, err := c.getOtherServicesClaimedPortForwards(currentSvcName)
	if err != nil {
		klog.Errorf("failed to get claimed port-forwards for service %s: %v", currentSvcName, err)
		return err
	}

	// 1. Delete removed rules (incremental delete, preserving other services' rules and accounting)
	for _, r := range toDeleteRules {
		if _, claimed := claimedByOthers[r.key()]; claimed {
			klog.Infof("port-forward %s:%d is still claimed by another service, skip deleting rule", r.ip, r.port)
			continue
		}
		targetComment := iptablesRuleComment(r.ip, r.port)
		for _, e := range entries {
			if e.comment == targetComment {
				delRule := portForwardConnLimitRule{ip: r.ip, port: r.port, connLimit: e.limit}
				if delErr := c.deleteIptablesRule(protocol, delRule); delErr != nil {
					return delErr
				}
			}
		}
	}

	// 2. Add or update rules (add new rule first, then delete old rule to prevent unmetered window)
	for _, r := range toAddOrUpdate {
		targetComment := iptablesRuleComment(r.ip, r.port)
		matchedRuleFound := false
		var oldRulesToDelete []portForwardConnLimitRule
		for _, e := range entries {
			if e.comment == targetComment {
				if e.limit == r.connLimit && !matchedRuleFound {
					// Rule with exact same limit already exists in kernel; preserve it to avoid resetting connlimit accounting
					matchedRuleFound = true
				} else {
					oldRulesToDelete = append(oldRulesToDelete, portForwardConnLimitRule{ip: r.ip, port: r.port, connLimit: e.limit})
				}
			}
		}
		if !matchedRuleFound {
			if addErr := c.addIptablesRule(protocol, r); addErr != nil {
				klog.Errorf("failed to add port-forward connlimit rule %s:%d (limit=%d, protocol %s): %v",
					r.ip, r.port, r.connLimit, protocol, addErr)
				return addErr
			}
		}
		for _, delRule := range oldRulesToDelete {
			if delErr := c.deleteIptablesRule(protocol, delRule); delErr != nil {
				klog.Errorf("failed to delete replaced port-forward connlimit rule %s:%d (old limit=%d, protocol %s): %v",
					delRule.ip, delRule.port, delRule.connLimit, protocol, delErr)
				return delErr
			}
		}
	}

	klog.Infof("successfully reconciled iptables port-forward connlimit rules (protocol %s, adds=%d, dels=%d)",
		protocol, len(toAddOrUpdate), len(toDeleteRules))
	return nil
}

// gcPortForwardConnLimitRules scans current iptables rules, cleans up orphan rules,
// and restores missing active Service rules without resetting connlimit accounting.
func (c *Controller) gcPortForwardConnLimitRules() error {
	if c == nil || c.servicesLister == nil {
		return nil
	}

	c.pfConnLimitMutex.Lock()
	defer c.pfConnLimitMutex.Unlock()

	svcs, err := c.servicesLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("failed to list services for port-forward connlimit gc: %v", err)
		return err
	}

	activeRulesByProto := map[string]map[portForwardKey]portForwardConnLimitRule{
		kubeovnv1.ProtocolIPv4: make(map[portForwardKey]portForwardConnLimitRule),
		kubeovnv1.ProtocolIPv6: make(map[portForwardKey]portForwardConnLimitRule),
	}
	activeComments := make(map[string]bool)
	for _, svc := range svcs {
		for _, r := range c.parsePortForwardConnLimitRules(svc) {
			proto := kubeovnv1.ProtocolIPv4
			if util.CheckProtocol(r.ip) == kubeovnv1.ProtocolIPv6 {
				proto = kubeovnv1.ProtocolIPv6
			}
			activeRulesByProto[proto][r.key()] = r
			activeComments[iptablesRuleComment(r.ip, r.port)] = true
		}
	}

	for _, proto := range c.enabledProtocols() {
		if err := c.ensureIptablesChain(proto); err != nil {
			klog.Errorf("failed to ensure iptables chain for protocol %s during gc: %v", proto, err)
			continue
		}

		entries, err := c.listIptablesRules(proto)
		if err != nil {
			if strings.Contains(err.Error(), "chain does not exist") {
				continue
			}
			klog.Errorf("failed to list iptables rules for protocol %s during gc: %v", proto, err)
			continue
		}

		// 1. Clean up orphan rules
		for _, e := range entries {
			if strings.HasPrefix(e.comment, "ovn-pf-") && !activeComments[e.comment] {
				commentBody := strings.TrimPrefix(e.comment, "ovn-pf-")
				lastHyphen := strings.LastIndex(commentBody, "-")
				if lastHyphen <= 0 {
					continue
				}
				ip := commentBody[:lastHyphen]
				portStr := commentBody[lastHyphen+1:]
				port, parseErr := strconv.ParseInt(portStr, 10, 32)
				if parseErr != nil {
					continue
				}
				orphanRule := portForwardConnLimitRule{ip: ip, port: int32(port), connLimit: e.limit}
				klog.Infof("found orphan port-forward connlimit rule %q in protocol %s, deleting", e.comment, proto)
				if delErr := c.deleteIptablesRule(proto, orphanRule); delErr != nil {
					klog.Errorf("failed to delete orphan rule %q (protocol %s): %v", e.comment, proto, delErr)
				}
			}
		}

		// 2. Self-heal missing active rules
		activeRules := activeRulesByProto[proto]
		for _, r := range activeRules {
			targetComment := iptablesRuleComment(r.ip, r.port)
			matchedRuleFound := false
			var oldRulesToDelete []portForwardConnLimitRule
			for _, e := range entries {
				if e.comment == targetComment {
					if e.limit == r.connLimit && !matchedRuleFound {
						matchedRuleFound = true
					} else {
						oldRulesToDelete = append(oldRulesToDelete, portForwardConnLimitRule{ip: r.ip, port: r.port, connLimit: e.limit})
					}
				}
			}
			if !matchedRuleFound {
				klog.Infof("restoring missing port-forward connlimit rule %s:%d (limit=%d, protocol %s)",
					r.ip, r.port, r.connLimit, proto)
				if addErr := c.addIptablesRule(proto, r); addErr != nil {
					klog.Errorf("failed to add missing rule %s:%d (protocol %s): %v", r.ip, r.port, proto, addErr)
					continue
				}
			}
			for _, delRule := range oldRulesToDelete {
				if delErr := c.deleteIptablesRule(proto, delRule); delErr != nil {
					klog.Errorf("failed to delete replaced rule %s:%d (protocol %s): %v", delRule.ip, delRule.port, proto, delErr)
				}
			}
		}
	}
	return nil
}

func (c *Controller) gcPortForwardConnLimitRulesWorker() {
	if err := c.gcPortForwardConnLimitRules(); err != nil {
		klog.Errorf("failed to garbage collect orphan port-forward connlimit rules: %v", err)
	}
}
