package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestBuildPortForwardConnLimitRuleSpec(t *testing.T) {
	spec := buildPortForwardConnLimitRuleSpec("1.2.3.4", 30080, 20000)
	expected := []string{
		"-d", "1.2.3.4/32",
		"-p", "tcp",
		"-m", "tcp", "--dport", "30080",
		"--tcp-flags", "FIN,SYN,RST,ACK", "SYN",
		"-m", "connlimit",
		"--connlimit-above", "20000",
		"--connlimit-mask", "0",
		"-j", "DROP",
	}
	require.Equal(t, expected, spec)
}

func TestGetExternalInterface(t *testing.T) {
	c := &Controller{}
	require.Empty(t, c.getExternalInterface())

	c.config = &Configuration{
		Iface: "eth0",
	}
	require.Equal(t, "eth0", c.getExternalInterface())

	c.config.SetTunnelIface("bond0")
	require.Equal(t, "bond0", c.getExternalInterface())
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
					util.PortForwardConnLimitAnnotation:   "10000",
					util.MetalLBLoadBalancerIPsAnnotation: "198.51.100.1",
				},
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeLoadBalancer,
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

		rules := c.parsePortForwardConnLimitRules(svc)
		require.Len(t, rules, 1)
		require.Equal(t, "198.51.100.1", rules[0].ip)
		require.Equal(t, int32(8080), rules[0].port)
		require.Equal(t, "tcp", rules[0].protocol)
		require.Equal(t, 10000, rules[0].connLimit)
	})

	t.Run("fallback to tcp connlimit annotation", func(t *testing.T) {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					util.TCPConnLimitAnnotation:         "20000",
					util.MetalLBAllowSharedIPAnnotation: "198.51.100.2",
				},
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeLoadBalancer,
				Ports: []corev1.ServicePort{
					{
						Protocol: corev1.ProtocolTCP,
						Port:     9090,
					},
				},
			},
		}
		rules := c.parsePortForwardConnLimitRules(svc)
		require.Len(t, rules, 1)
		require.Equal(t, "198.51.100.2", rules[0].ip)
		require.Equal(t, int32(9090), rules[0].port)
		require.Equal(t, 20000, rules[0].connLimit)
	})

	t.Run("multiple ips and multiple ports", func(t *testing.T) {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					util.PortForwardConnLimitAnnotation:   "5000",
					util.MetalLBLoadBalancerIPsAnnotation: "198.51.100.10, 198.51.100.11",
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
		}

		rules := c.parsePortForwardConnLimitRules(svc)
		require.Len(t, rules, 4) // 2 IPs * 2 Ports

		keys := make(map[string]bool)
		for _, r := range rules {
			require.Equal(t, 5000, r.connLimit)
			keys[r.key()] = true
		}
		require.True(t, keys["198.51.100.10:30001/tcp"])
		require.True(t, keys["198.51.100.10:30002/tcp"])
		require.True(t, keys["198.51.100.11:30001/tcp"])
		require.True(t, keys["198.51.100.11:30002/tcp"])
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
		protocol:  "tcp",
		connLimit: 100,
	}
	require.Equal(t, "1.1.1.1:80/tcp", r.key())
}
