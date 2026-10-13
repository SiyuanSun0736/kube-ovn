package v1

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VpcNatGatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []VpcNatGateway `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced
// +resourceName=vpc-nat-gateways
//
// VpcNatGateway represents a NAT gateway for a VPC, implemented as a StatefulSet Pod.
//
// Architecture note:
// The NAT gateway Pod does NOT support hot updates. Any changes to Spec fields (ExternalSubnets,
// Selector, Tolerations, Affinity, etc.) will trigger a StatefulSet template update,
// which causes the Pod to be recreated via RollingUpdate strategy. This is by design because:
//  1. Network configuration (routes, iptables rules) is initialized at Pod startup
//  2. Runtime state (vpc_cidrs, init status) is managed by separate handlers and will be
//     automatically restored after Pod recreation through the normal reconciliation flow
//
// The exceptions are QoSPolicy and EnableLanIPAsHairpinSNAT, which are applied without Pod restart.
type VpcNatGateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   VpcNatGatewaySpec   `json:"spec"`
	Status VpcNatGatewayStatus `json:"status"`
}

type VpcNatGatewaySpec struct {
	// Namespace where the NAT gateway StatefulSet/Pod will be created.
	// If empty, defaults to the kube-ovn controller's own namespace (typically kube-system).
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
	// VPC name for the NAT gateway. This field is immutable after creation.
	Vpc string `json:"vpc"`
	// Subnet name for the NAT gateway. This field is immutable after creation.
	Subnet string `json:"subnet"`
	// External subnets accessible through the NAT gateway
	ExternalSubnets []string `json:"externalSubnets"`
	// LAN IP address for the NAT gateway. This field is immutable after creation.
	LanIP string `json:"lanIp"`
	// Pod selector for the NAT gateway
	Selector    []string            `json:"selector"`
	Tolerations []corev1.Toleration `json:"tolerations"`
	Affinity    corev1.Affinity     `json:"affinity"`
	// QoS policy name to apply to the NAT gateway
	QoSPolicy string `json:"qosPolicy"`
	// BGP speaker configuration
	BgpSpeaker VpcBgpSpeaker `json:"bgpSpeaker"`
	// Static routes for the NAT gateway
	Routes []Route `json:"routes"`
	// Disable default EIP assignment
	NoDefaultEIP bool `json:"noDefaultEIP"`
	// EnableLanIPAsHairpinSNAT turns on the fallback hairpin SNAT, off by default: VPC traffic
	// that this gateway DNAT'd back towards the VPC then leaves with the gateway lanIP as source,
	// which guarantees the reply returns to the instance holding the conntrack entry. Leave it off
	// when the backend preserves the client address and routes its replies back through this
	// gateway itself (an LB/envoy pod with a policy route to the lanIP).
	// It covers every VIP that is not a floating IP (Service ClusterIPs, the gateway lanIP,
	// share-DNAT VIPs); traffic dialed at a FIP is always SNAT'd to that FIP, independently of
	// this field. Applied without Pod recreation.
	// +kubebuilder:validation:Optional
	EnableLanIPAsHairpinSNAT bool `json:"enableLanIpAsHairpinSnat,omitempty"`
	// User-defined annotations for the StatefulSet NAT gateway Pod template.
	// Only effective at creation time; updates to this field are not detected.
	Annotations map[string]string `json:"annotations,omitempty"`
}

type VpcBgpSpeaker struct {
	Enabled               bool            `json:"enabled"`
	ASN                   uint32          `json:"asn"`
	RemoteASN             uint32          `json:"remoteAsn"`
	Neighbors             []string        `json:"neighbors"`
	HoldTime              metav1.Duration `json:"holdTime"`
	RouterID              string          `json:"routerId"`
	Password              string          `json:"password"` // #nosec G117
	EnableGracefulRestart bool            `json:"enableGracefulRestart"`
	ExtraArgs             []string        `json:"extraArgs"`
}

// TODO: Consider removing redundant Status fields since statefulset template changes always trigger Pod recreation.
type VpcNatGatewayStatus struct {
	QoSPolicy string `json:"qosPolicy" patchStrategy:"merge"`
	// HairpinSNAT is the state of the lanIP fallback currently programmed on the gateway Pod,
	// "enable" or "disable" (see Spec.EnableLanIPAsHairpinSNAT), or "unsupported" when that Pod
	// runs a gateway script that cannot switch it, in which case the Spec field only takes
	// effect once the Pod is replaced by a gateway image that can. omitempty keeps the other
	// status writers, which all serialize the whole struct from their own lister snapshot, from
	// patching a value they never read back over the one this field's own writer just set.
	HairpinSNAT     string              `json:"hairpinSnat,omitempty" patchStrategy:"merge"`
	ExternalSubnets []string            `json:"externalSubnets" patchStrategy:"merge"`
	Selector        []string            `json:"selector" patchStrategy:"merge"`
	Tolerations     []corev1.Toleration `json:"tolerations" patchStrategy:"merge"`
	Affinity        corev1.Affinity     `json:"affinity" patchStrategy:"merge"`
}

type Route struct {
	CIDR      string `json:"cidr"`
	NextHopIP string `json:"nextHopIP"`
}

func (s *VpcNatGatewayStatus) Bytes() ([]byte, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	newStr := fmt.Sprintf(`{"status": %s}`, string(bytes))
	klog.V(5).Info("status body", newStr)
	return []byte(newStr), nil
}
