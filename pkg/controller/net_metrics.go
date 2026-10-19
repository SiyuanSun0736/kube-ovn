package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	metricSubnetAvailableIPs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subnet_available_ip_count",
			Help: "The available num of ip address in subnet.",
		},
		[]string{
			"subnet_name",
			"protocol",
			"subnet_cidr",
		})

	metricSubnetUsedIPs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subnet_used_ip_count",
			Help: "The used num of ip address in subnet.",
		},
		[]string{
			"subnet_name",
			"protocol",
			"subnet_cidr",
		})

	metricCentralSubnetInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "centralized_subnet_info",
			Help: "Provide information for centralized subnet.",
		},
		[]string{
			"subnet_name",
			"enable_ecmp",
			"gateway_node",
			"active_gateway",
			"match",
			"nexthops",
		})

	metricSubnetIPAMInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subnet_ipam_info",
			Help: "Provide information for subnet ip address management.",
		},
		[]string{
			"subnet_name",
			"cidr",
			"free",
			"reserved",
			"available",
			"using",
		})

	metricSubnetIPAssignedInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subnet_ip_assign_info",
			Help: "Provide information for subnet ip address assigned info.",
		},
		[]string{
			"subnet_name",
			"ip",
			"pod_name",
		})

	metricNatGwReadyMembers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_ovn_nat_gw_ready_members",
			Help: "Number of ready and healthy HA gateway member pods for a VpcNatGateway.",
		},
		[]string{
			"gateway",
			"vpc",
		})

	metricPolicyRouteResolveErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kube_ovn_policy_route_resolve_errors_total",
			Help: "Cumulative number of errors when resolving dynamic next-hop LAN IPs for gateway policy routes.",
		},
		[]string{
			"vpc",
			"gateway",
		})
)

func registerMetrics() {
	metrics.Registry.MustRegister(metricSubnetAvailableIPs)
	metrics.Registry.MustRegister(metricSubnetUsedIPs)
	metrics.Registry.MustRegister(metricCentralSubnetInfo)
	metrics.Registry.MustRegister(metricSubnetIPAMInfo)
	metrics.Registry.MustRegister(metricSubnetIPAssignedInfo)
	metrics.Registry.MustRegister(metricNatGwReadyMembers)
	metrics.Registry.MustRegister(metricPolicyRouteResolveErrors)
	registerTunnelKeyMetrics()
}

