package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestParseIptablesChainRules(t *testing.T) {
	rawOutput := `-N KUBE-OVN-PF
-A KUBE-OVN-PF -d 1.2.3.4/32 -p tcp -m tcp --dport 80 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 100 --connlimit-mask 0 --connlimit-saddr -m comment --comment "ovn-pf-1.2.3.4-80" -j DROP
-A KUBE-OVN-PF -d 2001:db8::1/128 -p tcp -m tcp --dport 443 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 200 --connlimit-mask 0 --connlimit-saddr -m comment --comment 'ovn-pf-2001:db8::1-443' -j DROP
-A KUBE-OVN-PF -p tcp -m tcp --dport 22 -j ACCEPT
`

	entries := parseIptablesChainRules(rawOutput)
	require.Len(t, entries, 3)

	// IPv4 rule
	require.Equal(t, "ovn-pf-1.2.3.4-80", entries[0].comment)
	require.Equal(t, 100, entries[0].limit)
	require.Contains(t, entries[0].raw, "--connlimit-above 100")

	// IPv6 rule with single quote
	require.Equal(t, "ovn-pf-2001:db8::1-443", entries[1].comment)
	require.Equal(t, 200, entries[1].limit)

	// Rule without comment
	require.Equal(t, "", entries[2].comment)
	require.Equal(t, 0, entries[2].limit)

	// Find by comment
	matched := findIptablesRulesByComment(entries, "ovn-pf-1.2.3.4-80")
	require.Len(t, matched, 1)
	require.Equal(t, 100, matched[0].limit)

	matchedNone := findIptablesRulesByComment(entries, "non-existent")
	require.Empty(t, matchedNone)
}

func TestParseIptablesChainRules_DuplicateOrGhost(t *testing.T) {
	rawOutput := `-N KUBE-OVN-PF
-A KUBE-OVN-PF -d 1.2.3.4/32 -p tcp -m tcp --dport 80 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 100 --connlimit-mask 0 -m comment --comment "ovn-pf-1.2.3.4-80" -j DROP
-A KUBE-OVN-PF -d 1.2.3.4/32 -p tcp -m tcp --dport 80 --tcp-flags FIN,SYN,RST,ACK SYN -m connlimit --connlimit-above 200 --connlimit-mask 0 -m comment --comment "ovn-pf-1.2.3.4-80" -j DROP
`
	entries := parseIptablesChainRules(rawOutput)
	require.Len(t, entries, 2)
	require.Equal(t, 100, entries[0].limit)
	require.Equal(t, 200, entries[1].limit)
	matched := findIptablesRulesByComment(entries, "ovn-pf-1.2.3.4-80")
	require.Len(t, matched, 2)
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

func TestBuildIptablesRuleSpec(t *testing.T) {
	ruleV4 := portForwardConnLimitRule{
		ip:        "1.2.3.4",
		port:      80,
		connLimit: 100,
	}
	specV4 := buildIptablesRuleSpec(ruleV4)
	expectedV4 := []string{
		"-d", "1.2.3.4/32",
		"-p", "tcp",
		"-m", "tcp",
		"--dport", "80",
		"--tcp-flags", "FIN,SYN,RST,ACK", "SYN",
		"-m", "connlimit",
		"--connlimit-above", "100",
		"--connlimit-mask", "0",
		"-m", "comment",
		"--comment", "ovn-pf-1.2.3.4-80",
		"-j", "DROP",
	}
	require.Equal(t, expectedV4, specV4)

	ruleV6 := portForwardConnLimitRule{
		ip:        "2001:db8::1",
		port:      443,
		connLimit: 500,
	}
	specV6 := buildIptablesRuleSpec(ruleV6)
	expectedV6 := []string{
		"-d", "2001:db8::1/128",
		"-p", "tcp",
		"-m", "tcp",
		"--dport", "443",
		"--tcp-flags", "FIN,SYN,RST,ACK", "SYN",
		"-m", "connlimit",
		"--connlimit-above", "500",
		"--connlimit-mask", "0",
		"-m", "comment",
		"--comment", "ovn-pf-2001:db8::1-443",
		"-j", "DROP",
	}
	require.Equal(t, expectedV6, specV6)
}

func TestEnabledProtocols(t *testing.T) {
	c := &Controller{}

	c.protocol = kubeovnv1.ProtocolDual
	require.Equal(t, []string{kubeovnv1.ProtocolIPv4, kubeovnv1.ProtocolIPv6}, c.enabledProtocols())

	c.protocol = kubeovnv1.ProtocolIPv4
	require.Equal(t, []string{kubeovnv1.ProtocolIPv4}, c.enabledProtocols())

	c.protocol = kubeovnv1.ProtocolIPv6
	require.Equal(t, []string{kubeovnv1.ProtocolIPv6}, c.enabledProtocols())

	var nilCtrl *Controller
	require.Equal(t, []string{kubeovnv1.ProtocolIPv4, kubeovnv1.ProtocolIPv6}, nilCtrl.enabledProtocols())
}

func TestGetOtherServicesClaimedPortForwards(t *testing.T) {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	c := &Controller{
		servicesLister: corelisters.NewServiceLister(indexer),
	}

	svc1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "svc1",
			Namespace:   "default",
			Annotations: map[string]string{util.PortForwardConnLimitAnnotation: "100"},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{Protocol: corev1.ProtocolTCP, Port: 80},
			},
		},
		Status: corev1.ServiceStatus{
			LoadBalancer: corev1.LoadBalancerStatus{
				Ingress: []corev1.LoadBalancerIngress{{IP: "1.2.3.4"}},
			},
		},
	}
	svc2 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "svc2",
			Namespace:   "default",
			Annotations: map[string]string{util.PortForwardConnLimitAnnotation: "200"},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{Protocol: corev1.ProtocolTCP, Port: 443},
			},
		},
		Status: corev1.ServiceStatus{
			LoadBalancer: corev1.LoadBalancerStatus{
				Ingress: []corev1.LoadBalancerIngress{{IP: "1.2.3.4"}},
			},
		},
	}

	require.NoError(t, indexer.Add(svc1))
	require.NoError(t, indexer.Add(svc2))

	claimed, err := c.getOtherServicesClaimedPortForwards("default/svc1")
	require.NoError(t, err)
	require.Len(t, claimed, 1)
	_, exists := claimed[portForwardKey{ip: "1.2.3.4", port: 443}]
	require.True(t, exists)
	_, exists = claimed[portForwardKey{ip: "1.2.3.4", port: 80}]
	require.False(t, exists)
}
