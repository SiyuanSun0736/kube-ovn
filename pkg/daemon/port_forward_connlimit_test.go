package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestParseIptablesChainRules(t *testing.T) {
	rawOutput := `-N KUBE-OVN-PF
-A KUBE-OVN-PF -i lo -m comment --comment ovn-pf-gate -j RETURN
-A KUBE-OVN-PF ! -i eth0 -m comment --comment ovn-pf-gate -j RETURN
-A KUBE-OVN-PF -d 1.2.3.4/32 -p tcp -m tcp --dport 80 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 100 --connlimit-mask 0 --connlimit-saddr -m comment --comment "ovn-pf-1.2.3.4-80" -j DROP
-A KUBE-OVN-PF -d 2001:db8::1/128 -p tcp -m tcp --dport 443 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 200 --connlimit-mask 0 --connlimit-saddr -m comment --comment 'ovn-pf-2001:db8::1-443' -j DROP
-A KUBE-OVN-PF -p tcp -m tcp --dport 22 -j ACCEPT
`

	entries := parseIptablesChainRules(rawOutput)
	require.Len(t, entries, 5)

	// Gate rule 1
	require.Equal(t, "ovn-pf-gate", entries[0].comment)
	require.Equal(t, 0, entries[0].limit)
	require.Contains(t, entries[0].raw, "-i lo")

	// Gate rule 2
	require.Equal(t, "ovn-pf-gate", entries[1].comment)
	require.Equal(t, 0, entries[1].limit)
	require.Contains(t, entries[1].raw, "! -i eth0")

	// IPv4 rule
	require.Equal(t, "ovn-pf-1.2.3.4-80", entries[2].comment)
	require.Equal(t, 100, entries[2].limit)
	require.Contains(t, entries[2].raw, "--connlimit-above 100")

	// IPv6 rule with single quote
	require.Equal(t, "ovn-pf-2001:db8::1-443", entries[3].comment)
	require.Equal(t, 200, entries[3].limit)

	// Rule without comment
	require.Equal(t, "", entries[4].comment)
	require.Equal(t, 0, entries[4].limit)

	// Find by comment
	matched := findIptablesRulesByComment(entries, "ovn-pf-1.2.3.4-80")
	require.Len(t, matched, 1)
	require.Equal(t, 100, matched[0].limit)

	matchedGate := findIptablesRulesByComment(entries, "ovn-pf-gate")
	require.Len(t, matchedGate, 2)

	matchedNone := findIptablesRulesByComment(entries, "non-existent")
	require.Empty(t, matchedNone)
}

func TestBuildIptablesGateRules(t *testing.T) {
	rulesV4 := buildIptablesGateRules(kubeovnv1.ProtocolIPv4)
	require.Len(t, rulesV4, 1)
	require.Equal(t, `-A KUBE-OVN-PF -s 127.0.0.0/8 -m comment --comment "ovn-pf-gate" -j RETURN`, rulesV4[0])

	rulesV6 := buildIptablesGateRules(kubeovnv1.ProtocolIPv6)
	require.Len(t, rulesV6, 1)
	require.Equal(t, `-A KUBE-OVN-PF -s ::1/128 -m comment --comment "ovn-pf-gate" -j RETURN`, rulesV6[0])
}

func TestBuildIptablesAddRuleCmd(t *testing.T) {
	ruleV4 := portForwardConnLimitRule{
		ip:        "1.1.1.1",
		port:      80,
		connLimit: 100,
	}
	expectedV4 := `-A KUBE-OVN-PF -d 1.1.1.1/32 -p tcp -m tcp --dport 80 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 100 --connlimit-mask 0 -m comment --comment "ovn-pf-1.1.1.1-80" -j DROP`
	require.Equal(t, expectedV4, buildIptablesAddRuleCmd(ruleV4))

	ruleV6 := portForwardConnLimitRule{
		ip:        "2001:db8::1",
		port:      443,
		connLimit: 500,
	}
	expectedV6 := `-A KUBE-OVN-PF -d 2001:db8::1/128 -p tcp -m tcp --dport 443 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 500 --connlimit-mask 0 -m comment --comment "ovn-pf-2001:db8::1-443" -j DROP`
	require.Equal(t, expectedV6, buildIptablesAddRuleCmd(ruleV6))
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

func TestExtractServicePublicIPs(t *testing.T) {
	svc := &corev1.Service{
		Status: corev1.ServiceStatus{
			LoadBalancer: corev1.LoadBalancerStatus{
				Ingress: []corev1.LoadBalancerIngress{
					{IP: "198.51.100.20"},
					{IP: "198.51.100.5"},
					{IP: "198.51.100.10"},
					{IP: "invalid-ip"},
					{IP: ""},
				},
			},
		},
	}
	ips := extractServicePublicIPs(svc)
	require.Equal(t, []string{"198.51.100.10", "198.51.100.20", "198.51.100.5"}, ips)
}
