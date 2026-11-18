package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestBuildNftAddRuleCmd(t *testing.T) {
	t.Run("ipv4 rule command", func(t *testing.T) {
		r := portForwardConnLimitRule{
			ip:        "1.2.3.4",
			port:      80,
			connLimit: 1000,
		}
		expected := `add rule inet kube_ovn_pf prerouting ip daddr 1.2.3.4 tcp dport 80 tcp flags & (fin|syn|rst|ack) == syn ct count over 1000 drop comment "ovn-pf-1.2.3.4-80"`
		require.Equal(t, expected, buildNftAddRuleCmd(r))
	})

	t.Run("ipv6 rule command", func(t *testing.T) {
		r := portForwardConnLimitRule{
			ip:        "fd00::100",
			port:      443,
			connLimit: 5000,
		}
		expected := `add rule inet kube_ovn_pf prerouting ip6 daddr fd00::100 tcp dport 443 tcp flags & (fin|syn|rst|ack) == syn ct count over 5000 drop comment "ovn-pf-fd00::100-443"`
		require.Equal(t, expected, buildNftAddRuleCmd(r))
	})
}

func TestParseNftRuleHandles(t *testing.T) {
	output := `table inet kube_ovn_pf {
	chain prerouting { # handle 1
		type filter hook prerouting priority -150; policy accept;
		ip daddr 1.2.3.4 tcp dport 80 tcp flags syn / fin,syn,rst,ack ct count over 100 drop comment "ovn-pf-1.2.3.4-80" # handle 2
		ip daddr 1.2.3.4 tcp dport 8080 tcp flags syn / fin,syn,rst,ack ct count over 500 drop comment "ovn-pf-1.2.3.4-8080" # handle 3
		ip6 daddr fd00::100 tcp dport 443 tcp flags syn / fin,syn,rst,ack ct count over 200 drop comment "ovn-pf-fd00::100-443" # handle 4
	}
}`

	t.Run("match specific handle", func(t *testing.T) {
		handles := parseNftRuleHandles(output, "ovn-pf-1.2.3.4-80")
		require.Equal(t, []string{"2"}, handles)

		handles = parseNftRuleHandles(output, "ovn-pf-1.2.3.4-8080")
		require.Equal(t, []string{"3"}, handles)

		handles = parseNftRuleHandles(output, "ovn-pf-fd00::100-443")
		require.Equal(t, []string{"4"}, handles)
	})

	t.Run("no match", func(t *testing.T) {
		handles := parseNftRuleHandles(output, "ovn-pf-1.2.3.4-9090")
		require.Empty(t, handles)
	})
}

func TestParseNftRuleLimit(t *testing.T) {
	output := `table inet kube_ovn_pf {
	chain prerouting { # handle 1
		type filter hook prerouting priority -150; policy accept;
		ip daddr 1.2.3.4 tcp dport 80 tcp flags syn / fin,syn,rst,ack ct count over 100 drop comment "ovn-pf-1.2.3.4-80" # handle 2
		ip6 daddr fd00::100 tcp dport 443 tcp flags syn / fin,syn,rst,ack ct count over 5000 drop comment "ovn-pf-fd00::100-443" # handle 3
	}
}`

	t.Run("found limit", func(t *testing.T) {
		limit, found := parseNftRuleLimit(output, "ovn-pf-1.2.3.4-80")
		require.True(t, found)
		require.Equal(t, 100, limit)

		limit, found = parseNftRuleLimit(output, "ovn-pf-fd00::100-443")
		require.True(t, found)
		require.Equal(t, 5000, limit)
	})

	t.Run("not found", func(t *testing.T) {
		limit, found := parseNftRuleLimit(output, "ovn-pf-nonexistent-80")
		require.False(t, found)
		require.Equal(t, 0, limit)
	})
}

func TestParsePortForwardConnLimitRules(t *testing.T) {
	c := &Controller{}

	t.Run("nil service", func(t *testing.T) {
		rules := c.parsePortForwardConnLimitRules(nil)
		require.Nil(t, rules)
	})

	t.Run("non-loadbalancer service", func(t *testing.T) {
		svc := &corev1.Service{
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeClusterIP,
			},
		}
		rules := c.parsePortForwardConnLimitRules(svc)
		require.Nil(t, rules)
	})

	t.Run("missing connlimit annotation", func(t *testing.T) {
		svc := &corev1.Service{
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeLoadBalancer,
			},
		}
		rules := c.parsePortForwardConnLimitRules(svc)
		require.Nil(t, rules)
	})

	t.Run("invalid or zero connlimit annotation", func(t *testing.T) {
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
