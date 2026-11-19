package daemon

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

const nftGateComment = "ovn-pf-gate"

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

type nftRuleEntry struct {
	handle  string
	comment string
	limit   int
}

func nftRuleComment(ip string, port int32) string {
	return fmt.Sprintf("ovn-pf-%s-%d", ip, port)
}

func (c *Controller) getExternalInterface() string {
	if c == nil || c.config == nil {
		return ""
	}
	if iface := c.config.GetTunnelIface(); iface != "" {
		return iface
	}
	return c.config.Iface
}

func buildNftGateRule(extIface string) string {
	if extIface != "" {
		return fmt.Sprintf("add rule inet %s %s iifname != %q return comment %q",
			util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, extIface, nftGateComment)
	}
	return fmt.Sprintf("add rule inet %s %s iifname %q return comment %q",
		util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, "lo", nftGateComment)
}

func buildNftAddRuleCmd(rule portForwardConnLimitRule) string {
	ipFamily := "ip"
	if util.CheckProtocol(rule.ip) == kubeovnv1.ProtocolIPv6 {
		ipFamily = "ip6"
	}
	comment := nftRuleComment(rule.ip, rule.port)
	return fmt.Sprintf("add rule inet %s %s %s daddr %s tcp dport %d tcp flags & (fin|syn|rst|ack) == syn ct count over %d drop comment %q",
		util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, ipFamily, rule.ip, rule.port, rule.connLimit, comment)
}

func parseNftChainRules(output string) []nftRuleEntry {
	var entries []nftRuleEntry
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "table ") || strings.HasPrefix(line, "chain ") || line == "}" {
			continue
		}

		idxHandle := strings.LastIndex(line, "# handle ")
		if idxHandle == -1 {
			continue
		}
		handle := strings.TrimSpace(line[idxHandle+len("# handle "):])
		if handle == "" {
			continue
		}

		var comment string
		if _, rem, found := strings.Cut(line, "comment "); found {
			rem = strings.TrimSpace(rem)
			if len(rem) > 0 && (rem[0] == '"' || rem[0] == '\'') {
				quote := rem[0]
				if end := strings.IndexByte(rem[1:], quote); end != -1 {
					comment = rem[1 : end+1]
				}
			}
		}

		var limit int
		if _, rem, found := strings.Cut(line, "ct count over "); found {
			fields := strings.Fields(rem)
			if len(fields) > 0 {
				if l, err := strconv.Atoi(fields[0]); err == nil {
					limit = l
				}
			}
		}

		entries = append(entries, nftRuleEntry{
			handle:  handle,
			comment: comment,
			limit:   limit,
		})
	}
	return entries
}

func findNftRulesByComment(entries []nftRuleEntry, targetComment string) []nftRuleEntry {
	var res []nftRuleEntry
	for _, e := range entries {
		if e.comment == targetComment {
			res = append(res, e)
		}
	}
	return res
}

func (c *Controller) listNftRules() ([]nftRuleEntry, error) {
	out, err := exec.Command("nft", "-a", "list", "chain", "inet", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain).CombinedOutput() // #nosec G204
	if err != nil {
		outStr := string(out)
		if strings.Contains(outStr, "No such file or directory") || strings.Contains(outStr, "does not exist") {
			return nil, fmt.Errorf("table or chain does not exist: %w", err)
		}
		return nil, fmt.Errorf("nft command failed: %w, out: %s", err, outStr)
	}
	return parseNftChainRules(string(out)), nil
}

func isNftTableOrChainNotExist(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "No such file or directory") || strings.Contains(msg, "does not exist")
}

func (c *Controller) ensureNftTableAndChain() error {
	out, err := exec.Command("nft", "-a", "list", "chain", "inet", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain).CombinedOutput() // #nosec G204
	if err == nil {
		entries := parseNftChainRules(string(out))
		for _, e := range entries {
			if e.comment == nftGateComment {
				return nil
			}
		}
		gateCmd := buildNftGateRule(c.getExternalInterface())
		cmd := exec.Command("nft", "-f", "-") // #nosec G204 // #nosec G204
		cmd.Stdin = strings.NewReader(gateCmd + "\n")
		if gateOut, gateErr := cmd.CombinedOutput(); gateErr != nil {
			klog.Errorf("failed to insert nft gate rule: %v, out: %s", gateErr, string(gateOut))
			return gateErr
		}
		return nil
	}

	gateRule := buildNftGateRule(c.getExternalInterface())
	script := fmt.Sprintf(`add table inet %[1]s
add chain inet %[1]s %[2]s { type filter hook prerouting priority -150; policy accept; }
%[3]s
`, util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, gateRule)

	cmd := exec.Command("nft", "-f", "-") // #nosec G204 // #nosec G204
	cmd.Stdin = strings.NewReader(script)
	if initOut, initErr := cmd.CombinedOutput(); initErr != nil {
		klog.Errorf("failed to init nft table %s and chain %s: %v, out: %s", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, initErr, string(initOut))
		return initErr
	}
	klog.Infof("initialized nft table %s and chain %s with gate rule", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain)
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
	return result
}

func (c *Controller) reconcilePortForwardConnLimit(oldService, newService *corev1.Service) error {
	oldRules := c.parsePortForwardConnLimitRules(oldService)
	newRules := c.parsePortForwardConnLimitRules(newService)

	if len(oldRules) == 0 && len(newRules) == 0 {
		return nil
	}

	oldMap := make(map[portForwardKey]portForwardConnLimitRule, len(oldRules))
	for _, r := range oldRules {
		oldMap[r.key()] = r
	}

	newMap := make(map[portForwardKey]portForwardConnLimitRule, len(newRules))
	for _, r := range newRules {
		newMap[r.key()] = r
	}

	// Rules to delete: in old but not in new, or limit changed
	for key, oldRule := range oldMap {
		newRule, exists := newMap[key]
		if !exists || newRule.connLimit != oldRule.connLimit {
			if err := c.deletePortForwardConnLimitRule(oldRule); err != nil {
				klog.Errorf("failed to delete port-forward connlimit rule for %s:%d: %v", key.ip, key.port, err)
			}
		}
	}

	// Rules to add: in new but not in old, or limit changed
	for key, newRule := range newMap {
		oldRule, exists := oldMap[key]
		if !exists || oldRule.connLimit != newRule.connLimit {
			if err := c.addPortForwardConnLimitRule(newRule); err != nil {
				klog.Errorf("failed to add port-forward connlimit rule for %s:%d: %v", key.ip, key.port, err)
				return err
			}
		}
	}

	return nil
}

func (c *Controller) addPortForwardConnLimitRule(rule portForwardConnLimitRule) error {
	if err := c.ensureNftTableAndChain(); err != nil {
		return err
	}

	comment := nftRuleComment(rule.ip, rule.port)
	entries, err := c.listNftRules()
	if err != nil {
		klog.Errorf("failed to list nft rules before adding rule for %s:%d: %v", rule.ip, rule.port, err)
		return err
	}

	matched := findNftRulesByComment(entries, comment)
	if len(matched) == 1 && matched[0].limit == rule.connLimit {
		klog.V(4).Infof("nft rule for %s:%d with limit %d already exists, preserving accounting match", rule.ip, rule.port, rule.connLimit)
		return nil
	}

	// Atomically delete existing matching rules and add new rule in one batch
	var batch []string
	for _, m := range matched {
		batch = append(batch, fmt.Sprintf("delete rule inet %s %s handle %s", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, m.handle))
	}
	batch = append(batch, buildNftAddRuleCmd(rule))

	cmd := exec.Command("nft", "-f", "-") // #nosec G204
	cmd.Stdin = strings.NewReader(strings.Join(batch, "\n") + "\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		klog.Errorf("failed to apply nft connlimit rule for %s:%d: %v, out: %s", rule.ip, rule.port, err, string(out))
		return err
	}
	klog.Infof("applied nft port-forward connlimit rule for %s:%d limit %d (deleted %d old handles)", rule.ip, rule.port, rule.connLimit, len(matched))
	return nil
}

func (c *Controller) deletePortForwardConnLimitRule(rule portForwardConnLimitRule) error {
	comment := nftRuleComment(rule.ip, rule.port)
	entries, err := c.listNftRules()
	if err != nil {
		if isNftTableOrChainNotExist(err) {
			return nil
		}
		klog.Errorf("failed to list nft rules for deleting %s:%d: %v", rule.ip, rule.port, err)
		return err
	}

	matched := findNftRulesByComment(entries, comment)
	if len(matched) == 0 {
		return nil
	}

	var cmds []string
	for _, m := range matched {
		cmds = append(cmds, fmt.Sprintf("delete rule inet %s %s handle %s", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, m.handle))
	}

	cmd := exec.Command("nft", "-f", "-") // #nosec G204
	cmd.Stdin = strings.NewReader(strings.Join(cmds, "\n") + "\n")
	if delOut, delErr := cmd.CombinedOutput(); delErr != nil {
		if isNftTableOrChainNotExist(delErr) {
			return nil
		}
		klog.Errorf("failed to delete nft rule for %s:%d (handles %v): %v, out: %s", rule.ip, rule.port, matched, delErr, string(delOut))
		return delErr
	}
	klog.Infof("deleted nft rule for %s:%d (handles %v)", rule.ip, rule.port, matched)
	return nil
}

// gcPortForwardConnLimitRules scans current nftables rules and cleans up orphan rules
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
			activeComments[nftRuleComment(r.ip, r.port)] = true
		}
	}

	entries, err := c.listNftRules()
	if err != nil {
		return
	}

	var toDelete []string
	for _, e := range entries {
		if strings.HasPrefix(e.comment, "ovn-pf-") && e.comment != nftGateComment {
			if !activeComments[e.comment] {
				toDelete = append(toDelete, e.handle)
				klog.Infof("found orphan port-forward connlimit rule %q (handle %s), queueing for deletion", e.comment, e.handle)
			}
		}
	}

	if len(toDelete) == 0 {
		return
	}

	var cmds []string
	for _, h := range toDelete {
		cmds = append(cmds, fmt.Sprintf("delete rule inet %s %s handle %s", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, h))
	}

	cmd := exec.Command("nft", "-f", "-") // #nosec G204
	cmd.Stdin = strings.NewReader(strings.Join(cmds, "\n") + "\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		klog.Errorf("failed to garbage collect orphan port-forward connlimit rules (handles %v): %v, out: %s", toDelete, err, string(out))
		return
	}
	klog.Infof("successfully garbage collected %d orphan port-forward connlimit rules (handles %v)", len(toDelete), toDelete)
}
