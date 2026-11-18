package daemon

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
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

func nftRuleComment(ip string, port int32) string {
	return fmt.Sprintf("ovn-pf-%s-%d", ip, port)
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

func ensureNftTableAndChain() error {
	if err := exec.Command("nft", "list", "chain", "inet", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain).Run(); err == nil {
		return nil
	}

	script := fmt.Sprintf(`add table inet %[1]s
add chain inet %[1]s %[2]s { type filter hook prerouting priority -150; policy accept; }
`, util.OvnPortForwardNftTable, util.OvnPortForwardNftChain)

	cmd := exec.Command("nft", "-f", "-")
	cmd.Stdin = strings.NewReader(script)
	if out, err := cmd.CombinedOutput(); err != nil {
		klog.Errorf("failed to init nft table %s and chain %s: %v, out: %s", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, err, string(out))
		return err
	}
	klog.Infof("initialized nft table %s and chain %s", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain)
	return nil
}

func parseNftRuleHandles(output, targetComment string) []string {
	var handles []string
	expectedComment := fmt.Sprintf("comment %q", targetComment)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, expectedComment) {
			idx := strings.LastIndex(line, "# handle ")
			if idx != -1 {
				handle := strings.TrimSpace(line[idx+len("# handle "):])
				if handle != "" {
					handles = append(handles, handle)
				}
			}
		}
	}
	return handles
}

func parseNftRuleLimit(output, targetComment string) (int, bool) {
	expectedComment := fmt.Sprintf("comment %q", targetComment)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, expectedComment) {
			idx := strings.Index(line, "ct count over ")
			if idx != -1 {
				fields := strings.Fields(line[idx+len("ct count over "):])
				if len(fields) > 0 {
					if limit, err := strconv.Atoi(fields[0]); err == nil {
						return limit, true
					}
				}
			}
		}
	}
	return 0, false
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
	if err := ensureNftTableAndChain(); err != nil {
		return err
	}

	comment := nftRuleComment(rule.ip, rule.port)
	out, err := exec.Command("nft", "-a", "list", "chain", "inet", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain).CombinedOutput()
	if err == nil {
		if limit, found := parseNftRuleLimit(string(out), comment); found && limit == rule.connLimit {
			klog.V(4).Infof("nft rule for %s:%d with limit %d already exists, preserving accounting match", rule.ip, rule.port, rule.connLimit)
			return nil
		}
		if handles := parseNftRuleHandles(string(out), comment); len(handles) > 0 {
			if err := c.deletePortForwardConnLimitRule(rule); err != nil {
				klog.Errorf("failed to delete existing nft rule before update for %s:%d: %v", rule.ip, rule.port, err)
			}
		}
	}

	cmdStr := buildNftAddRuleCmd(rule)
	cmd := exec.Command("nft", "-f", "-")
	cmd.Stdin = strings.NewReader(cmdStr + "\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		klog.Errorf("failed to add nft connlimit rule for %s:%d: %v, out: %s", rule.ip, rule.port, err, string(out))
		return err
	}
	klog.Infof("added nft port-forward connlimit rule for %s:%d limit %d", rule.ip, rule.port, rule.connLimit)
	return nil
}

func (c *Controller) deletePortForwardConnLimitRule(rule portForwardConnLimitRule) error {
	comment := nftRuleComment(rule.ip, rule.port)
	out, err := exec.Command("nft", "-a", "list", "chain", "inet", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain).CombinedOutput()
	if err != nil {
		return nil
	}

	handles := parseNftRuleHandles(string(out), comment)
	if len(handles) == 0 {
		return nil
	}

	var cmds []string
	for _, handle := range handles {
		cmds = append(cmds, fmt.Sprintf("delete rule inet %s %s handle %s", util.OvnPortForwardNftTable, util.OvnPortForwardNftChain, handle))
	}

	cmd := exec.Command("nft", "-f", "-")
	cmd.Stdin = strings.NewReader(strings.Join(cmds, "\n") + "\n")
	if delOut, delErr := cmd.CombinedOutput(); delErr != nil {
		klog.Errorf("failed to delete nft rule for %s:%d (handles %v): %v, out: %s", rule.ip, rule.port, handles, delErr, string(delOut))
		return delErr
	}
	klog.Infof("deleted nft rule for %s:%d (handles %v)", rule.ip, rule.port, handles)
	return nil
}
