package daemon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestParseNftChainRules(t *testing.T) {
	rawOutput := `table inet kube_ovn_pf {
	chain prerouting {
		type filter hook prerouting priority -150; policy accept;
		iifname != "eth0" return comment "ovn-pf-gate" # handle 2
		ip daddr 1.2.3.4 tcp dport 80 tcp flags & (fin|syn|rst|ack) == syn ct count over 100 drop comment "ovn-pf-1.2.3.4-80" # handle 3
		ip6 daddr 2001:db8::1 tcp dport 443 tcp flags & (fin|syn|rst|ack) == syn ct count over 200 drop comment 'ovn-pf-2001:db8::1-443' # handle 4
		tcp dport 22 accept # handle 5
	}
}`

	entries := parseNftChainRules(rawOutput)
	require.Len(t, entries, 4)

	// Gate rule
	require.Equal(t, "2", entries[0].handle)
	require.Equal(t, "ovn-pf-gate", entries[0].comment)
	require.Equal(t, 0, entries[0].limit)

	// IPv4 rule
	require.Equal(t, "3", entries[1].handle)
	require.Equal(t, "ovn-pf-1.2.3.4-80", entries[1].comment)
	require.Equal(t, 100, entries[1].limit)

	// IPv6 rule with single quote
	require.Equal(t, "4", entries[2].handle)
	require.Equal(t, "ovn-pf-2001:db8::1-443", entries[2].comment)
	require.Equal(t, 200, entries[2].limit)

	// Rule without comment
	require.Equal(t, "5", entries[3].handle)
	require.Equal(t, "", entries[3].comment)
	require.Equal(t, 0, entries[3].limit)

	// Find by comment
	matched := findNftRulesByComment(entries, "ovn-pf-1.2.3.4-80")
	require.Len(t, matched, 1)
	require.Equal(t, "3", matched[0].handle)
	require.Equal(t, 100, matched[0].limit)

	matchedGate := findNftRulesByComment(entries, "ovn-pf-gate")
	require.Len(t, matchedGate, 1)
	require.Equal(t, "2", matchedGate[0].handle)

	matchedNone := findNftRulesByComment(entries, "non-existent")
	require.Empty(t, matchedNone)
}

func TestBuildNftGateRule(t *testing.T) {
	ruleWithNic := buildNftGateRule("eth0")
	require.Equal(t, `add rule inet kube_ovn_pf prerouting iifname != "eth0" return comment "ovn-pf-gate"`, ruleWithNic)

	ruleWithoutNic := buildNftGateRule("")
	require.Equal(t, `add rule inet kube_ovn_pf prerouting iifname "lo" return comment "ovn-pf-gate"`, ruleWithoutNic)
}

func TestBuildNftAddRuleCmd(t *testing.T) {
	ruleV4 := portForwardConnLimitRule{
		ip:        "1.1.1.1",
		port:      80,
		connLimit: 100,
	}
	expectedV4 := `add rule inet kube_ovn_pf prerouting ip daddr 1.1.1.1 tcp dport 80 tcp flags & (fin|syn|rst|ack) == syn ct count over 100 drop comment "ovn-pf-1.1.1.1-80"`
	require.Equal(t, expectedV4, buildNftAddRuleCmd(ruleV4))

	ruleV6 := portForwardConnLimitRule{
		ip:        "2001:db8::1",
		port:      443,
		connLimit: 500,
	}
	expectedV6 := `add rule inet kube_ovn_pf prerouting ip6 daddr 2001:db8::1 tcp dport 443 tcp flags & (fin|syn|rst|ack) == syn ct count over 500 drop comment "ovn-pf-2001:db8::1-443"`
	require.Equal(t, expectedV6, buildNftAddRuleCmd(ruleV6))
}

func TestIsNftTableOrChainNotExist(t *testing.T) {
	require.False(t, isNftTableOrChainNotExist(nil))
	require.False(t, isNftTableOrChainNotExist(errors.New("permission denied")))
	require.True(t, isNftTableOrChainNotExist(errors.New("Error: No such file or directory")))
	require.True(t, isNftTableOrChainNotExist(errors.New("table does not exist")))
}

func TestParsePortForwardConnLimitRules(t *testing.T) {
	c := &Controller{}

	t.Run("nil service or non-loadbalancer", func(t *testing.T) {
		require.Nil(t, c.parsePortForwardConnLimitRules(nil))

		svc := &corev1.Service{
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeClusterIP,
			},
		}
		require.Nil(t, c.parsePortForwardConnLimitRules(svc))
	})

	t.Run("invalid annotation", func(t *testing.T) {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					util.PortForwardConnLimitAnnotation: "not-a-number",
				},
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeLoadBalancer,
			},
		}
		require.Nil(t, c.parsePortForwardConnLimitRules(svc))

		svc.Annotations[util.PortForwardConnLimitAnnotation] = "0"
		require.Nil(t, c.parsePortForwardConnLimitRules(svc))

		svc.Annotations[util.PortForwardConnLimitAnnotation] = "-10"
		require.Nil(t, c.parsePortForwardConnLimitRules(svc))
	})

	t.Run("valid standard annotation and single ip", func(t *testing.T) {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					util.PortForwardConnLimitAnnotation: "10000",
				},
			},
			Spec: corev1.ServiceSpec{
				Type:           corev1.ServiceTypeLoadBalancer,
				LoadBalancerIP: "198.51.100.1",
				Ports: []corev1.ServicePort{
					{
						Protocol: corev1.ProtocolTCP,
						Port:     8080,
					},
					{
						Protocol: corev1.ProtocolUDP,
						Port:     8081,
					},
				},
			},
		}

		// Spec.LoadBalancerIP alone should not trigger rules until allocated in status
		require.Nil(t, c.parsePortForwardConnLimitRules(svc))

		// Once provisioned in status.loadBalancer.ingress, rules should be generated
		svc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{
			{IP: "198.51.100.1"},
		}

		rules := c.parsePortForwardConnLimitRules(svc)
		require.Len(t, rules, 1)
		require.Equal(t, "198.51.100.1", rules[0].ip)
		require.Equal(t, int32(8080), rules[0].port)
		require.Equal(t, 10000, rules[0].connLimit)
	})

	t.Run("multiple ips from status ingress and multiple ports", func(t *testing.T) {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					util.PortForwardConnLimitAnnotation: "5000",
				},
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeLoadBalancer,
				Ports: []corev1.ServicePort{
					{
						Protocol: corev1.ProtocolTCP,
						Port:     30001,
					},
					{
						Protocol: corev1.ProtocolTCP,
						Port:     30002,
					},
				},
			},
			Status: corev1.ServiceStatus{
				LoadBalancer: corev1.LoadBalancerStatus{
					Ingress: []corev1.LoadBalancerIngress{
						{IP: "198.51.100.10"},
						{IP: "198.51.100.11"},
					},
				},
			},
		}

		rules := c.parsePortForwardConnLimitRules(svc)
		require.Len(t, rules, 4) // 2 IPs * 2 Ports

		keys := make(map[portForwardKey]bool)
		for _, r := range rules {
			require.Equal(t, 5000, r.connLimit)
			keys[r.key()] = true
		}
		require.True(t, keys[portForwardKey{ip: "198.51.100.10", port: 30001}])
		require.True(t, keys[portForwardKey{ip: "198.51.100.10", port: 30002}])
		require.True(t, keys[portForwardKey{ip: "198.51.100.11", port: 30001}])
		require.True(t, keys[portForwardKey{ip: "198.51.100.11", port: 30002}])
	})

	t.Run("status ingress ip extraction", func(t *testing.T) {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					util.PortForwardConnLimitAnnotation: "15000",
				},
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeLoadBalancer,
				Ports: []corev1.ServicePort{
					{
						Protocol: corev1.ProtocolTCP,
						Port:     443,
					},
				},
			},
			Status: corev1.ServiceStatus{
				LoadBalancer: corev1.LoadBalancerStatus{
					Ingress: []corev1.LoadBalancerIngress{
						{IP: "203.0.113.50"},
					},
				},
			},
		}

		rules := c.parsePortForwardConnLimitRules(svc)
		require.Len(t, rules, 1)
		require.Equal(t, "203.0.113.50", rules[0].ip)
		require.Equal(t, int32(443), rules[0].port)
		require.Equal(t, 15000, rules[0].connLimit)
	})
}

func TestPortForwardConnLimitRuleKey(t *testing.T) {
	r := portForwardConnLimitRule{
		ip:        "1.1.1.1",
		port:      80,
		connLimit: 100,
	}
	expected := portForwardKey{
		ip:   "1.1.1.1",
		port: 80,
	}
	require.Equal(t, expected, r.key())
}
