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

const pfGateComment = "ovn-pf-gate"

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

func execIptablesRestore(protocol, script string) ([]byte, error) {
	binary := "iptables-restore"
	if protocol == kubeovnv1.ProtocolIPv6 {
		binary = "ip6tables-restore"
	}
	cmd := exec.Command(binary, "--noflush") // #nosec G204
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	cmd.Stdin = strings.NewReader(script)
	return cmd.CombinedOutput()
}

func buildIptablesGateRules(protocol string) []string {
	if protocol == kubeovnv1.ProtocolIPv6 {
		return []string{
			fmt.Sprintf("-A %s -s ::1/128 -m comment --comment %q -j RETURN", util.OvnPortForwardIptablesChain, pfGateComment),
		}
	}
	return []string{
		fmt.Sprintf("-A %s -s 127.0.0.0/8 -m comment --comment %q -j RETURN", util.OvnPortForwardIptablesChain, pfGateComment),
	}
}

func buildIptablesAddRuleCmd(rule portForwardConnLimitRule) string {
	cidr := "/32"
	if util.CheckProtocol(rule.ip) == kubeovnv1.ProtocolIPv6 {
		cidr = "/128"
	}
	comment := iptablesRuleComment(rule.ip, rule.port)
	return fmt.Sprintf("-A %s -d %s%s -p tcp -m tcp --dport %d --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above %d --connlimit-mask 0 -m comment --comment %q -j DROP",
		util.OvnPortForwardIptablesChain, rule.ip, cidr, rule.port, rule.connLimit, comment)
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
	out, err := execIptables(protocol, "-t", "mangle", "-S", util.OvnPortForwardIptablesChain)
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
	_, _ = execIptables(protocol, "-t", "mangle", "-N", util.OvnPortForwardIptablesChain)

	// Ensure jump in both PREROUTING and OUTPUT hooks
	for _, parent := range []string{"PREROUTING", "OUTPUT"} {
		if _, err := execIptables(protocol, "-t", "mangle", "-C", parent, "-j", util.OvnPortForwardIptablesChain); err != nil {
			if out, insErr := execIptables(protocol, "-t", "mangle", "-I", parent, "1", "-j", util.OvnPortForwardIptablesChain); insErr != nil {
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

	for _, proto := range []string{kubeovnv1.ProtocolIPv4, kubeovnv1.ProtocolIPv6} {
		pOld := rulesByProtoOld[proto]
		pNew := rulesByProtoNew[proto]
		if len(pOld) == 0 && len(pNew) == 0 {
			continue
		}

		if err := c.reconcileProtocolRules(proto, pOld, pNew); err != nil {
			klog.Errorf("failed to reconcile port-forward connlimit rules for protocol %s: %v", proto, err)
			return err
		}
	}

	return nil
}

func (c *Controller) reconcileProtocolRules(protocol string, oldRules, newRules []portForwardConnLimitRule) error {
	oldMap := make(map[portForwardKey]portForwardConnLimitRule, len(oldRules))
	for _, r := range oldRules {
		oldMap[r.key()] = r
	}
	newMap := make(map[portForwardKey]portForwardConnLimitRule, len(newRules))
	for _, r := range newRules {
		newMap[r.key()] = r
	}

	// Determine deleted and added/updated comments
	toDeleteComments := make(map[string]bool)
	for key, oldRule := range oldMap {
		if _, exists := newMap[key]; !exists {
			toDeleteComments[iptablesRuleComment(oldRule.ip, oldRule.port)] = true
		}
	}

	toAddOrUpdate := make(map[portForwardKey]portForwardConnLimitRule)
	for key, newRule := range newMap {
		oldRule, exists := oldMap[key]
		if !exists || oldRule.connLimit != newRule.connLimit {
			toAddOrUpdate[key] = newRule
		}
	}

	if len(toDeleteComments) == 0 && len(toAddOrUpdate) == 0 {
		return nil
	}

	if err := c.ensureIptablesChain(protocol); err != nil {
		return err
	}

	entries, _ := c.listIptablesRules(protocol)

	// Build new chain contents in atomic buffer
	var lines []string
	lines = append(lines, "*mangle")
	lines = append(lines, fmt.Sprintf(":%s - [0:0]", util.OvnPortForwardIptablesChain))
	lines = append(lines, fmt.Sprintf("-F %s", util.OvnPortForwardIptablesChain))

	// Re-insert gate rules at chain head
	lines = append(lines, buildIptablesGateRules(protocol)...)

	// Keep existing rules not belonging to toDeleteComments or toAddOrUpdate
	for _, e := range entries {
		if e.comment == "" || e.comment == pfGateComment {
			continue
		}
		if toDeleteComments[e.comment] {
			continue
		}
		// If rule is being updated, don't keep old version
		isUpdating := false
		for _, r := range toAddOrUpdate {
			if e.comment == iptablesRuleComment(r.ip, r.port) {
				isUpdating = true
				break
			}
		}
		if !isUpdating {
			lines = append(lines, e.raw)
		}
	}

	// Append added/updated rules
	for _, r := range toAddOrUpdate {
		lines = append(lines, buildIptablesAddRuleCmd(r))
	}

	lines = append(lines, "COMMIT\n")

	script := strings.Join(lines, "\n")
	if out, err := execIptablesRestore(protocol, script); err != nil {
		klog.Errorf("failed to apply iptables port-forward connlimit restore script (protocol %s): %v, out: %s", protocol, err, string(out))
		return err
	}

	klog.Infof("successfully applied iptables port-forward connlimit batch (protocol %s, adds=%d, dels=%d)",
		protocol, len(toAddOrUpdate), len(toDeleteComments))
	return nil
}

// gcPortForwardConnLimitRules scans current iptables rules and cleans up orphan rules
// whose comment starts with "ovn-pf-" but is no longer associated with any active Service.
func (c *Controller) gcPortForwardConnLimitRules() {
	if c == nil || c.servicesLister == nil {
		return
	}

	svcs, err := c.servicesLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("failed to list services for port-forward connlimit gc: %v", err)
		return
	}

	activeComments := make(map[string]bool)
	for _, svc := range svcs {
		for _, r := range c.parsePortForwardConnLimitRules(svc) {
			activeComments[iptablesRuleComment(r.ip, r.port)] = true
		}
	}

	for _, proto := range []string{kubeovnv1.ProtocolIPv4, kubeovnv1.ProtocolIPv6} {
		entries, err := c.listIptablesRules(proto)
		if err != nil || len(entries) == 0 {
			continue
		}

		hasOrphan := false
		for _, e := range entries {
			if strings.HasPrefix(e.comment, "ovn-pf-") && e.comment != pfGateComment {
				if !activeComments[e.comment] {
					hasOrphan = true
					klog.Infof("found orphan port-forward connlimit rule %q in protocol %s, queueing for gc", e.comment, proto)
				}
			}
		}

		if !hasOrphan {
			continue
		}

		// Rebuild chain keeping only active rules
		var lines []string
		lines = append(lines, "*mangle")
		lines = append(lines, fmt.Sprintf(":%s - [0:0]", util.OvnPortForwardIptablesChain))
		lines = append(lines, fmt.Sprintf("-F %s", util.OvnPortForwardIptablesChain))
		lines = append(lines, buildIptablesGateRules(proto)...)

		for _, e := range entries {
			if strings.HasPrefix(e.comment, "ovn-pf-") && e.comment != pfGateComment {
				if activeComments[e.comment] {
					lines = append(lines, e.raw)
				}
			}
		}
		lines = append(lines, "COMMIT\n")

		if out, err := execIptablesRestore(proto, strings.Join(lines, "\n")); err != nil {
			klog.Errorf("failed to garbage collect orphan port-forward connlimit rules (protocol %s): %v, out: %s", proto, err, string(out))
			continue
		}
		klog.Infof("successfully garbage collected orphan port-forward connlimit rules for protocol %s", proto)
	}
}
