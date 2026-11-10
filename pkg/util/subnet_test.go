package util

import (
	"testing"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

func TestIsValidTunnelKey(t *testing.T) {
	tests := []struct {
		key  int
		want bool
	}{
		{key: -1, want: false},
		{key: 0, want: false},
		{key: 1, want: true},
		{key: MaxTunnelKey, want: true},
		{key: MaxTunnelKey + 1, want: false},
	}
	for _, tt := range tests {
		if got := IsValidTunnelKey(tt.key); got != tt.want {
			t.Errorf("IsValidTunnelKey(%d) = %v, want %v", tt.key, got, tt.want)
		}
	}
}

func TestIsOvnVpcSubnet(t *testing.T) {
	tests := []struct {
		name   string
		subnet *kubeovnv1.Subnet
		want   bool
	}{
		{name: "nil", subnet: nil, want: false},
		{
			name: "implicit default VPC non-vlan subnet",
			subnet: &kubeovnv1.Subnet{
				Spec: kubeovnv1.SubnetSpec{Provider: ""},
			},
			want: true,
		},
		{
			name: "default cluster router VPC non-vlan subnet",
			subnet: &kubeovnv1.Subnet{
				Spec: kubeovnv1.SubnetSpec{Provider: OvnProvider, Vpc: DefaultVpc},
			},
			want: true,
		},
		{
			name: "custom VPC non-vlan subnet",
			subnet: &kubeovnv1.Subnet{
				Spec: kubeovnv1.SubnetSpec{Provider: OvnProvider, Vpc: "custom-vpc"},
			},
			want: true,
		},
		{
			name: "default VPC vlan underlay subnet",
			subnet: &kubeovnv1.Subnet{
				Spec: kubeovnv1.SubnetSpec{Provider: OvnProvider, Vpc: DefaultVpc, Vlan: "vlan-a"},
			},
			want: false,
		},
		{
			name: "custom VPC vlan underlay subnet",
			subnet: &kubeovnv1.Subnet{
				Spec: kubeovnv1.SubnetSpec{Provider: OvnProvider, Vpc: "custom-vpc", Vlan: "vlan-a"},
			},
			want: false,
		},
		{
			name: "non-OVN provider",
			subnet: &kubeovnv1.Subnet{
				Spec: kubeovnv1.SubnetSpec{Provider: "external.provider"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOvnVpcSubnet(tt.subnet); got != tt.want {
				t.Errorf("IsOvnVpcSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTunnelKeyAnnotationValidForSubnet(t *testing.T) {
	// The subnet's own logical switch key must no longer be accepted: pods carry the VPC key.
	vpcSubnet := &kubeovnv1.Subnet{
		Spec:   kubeovnv1.SubnetSpec{Provider: OvnProvider},
		Status: kubeovnv1.SubnetStatus{TunnelKey: 1234},
	}
	vlanSubnet := &kubeovnv1.Subnet{Spec: kubeovnv1.SubnetSpec{Provider: OvnProvider, Vlan: "vlan-a"}}
	tests := []struct {
		name         string
		annotations  map[string]string
		provider     string
		subnet       *kubeovnv1.Subnet
		vpcTunnelKey int
		want         bool
	}{
		{name: "matching default provider VPC key", annotations: map[string]string{TunnelKeyAnnotation: "5678"}, provider: OvnProvider, subnet: vpcSubnet, vpcTunnelKey: 5678, want: true},
		{name: "matching custom provider VPC key", annotations: map[string]string{"net1.default.ovn.kubernetes.io/tunnel_key": "5678"}, provider: "net1.default.ovn", subnet: vpcSubnet, vpcTunnelKey: 5678, want: true},
		{name: "subnet key is stale", annotations: map[string]string{TunnelKeyAnnotation: "1234"}, provider: OvnProvider, subnet: vpcSubnet, vpcTunnelKey: 5678},
		{name: "missing key", annotations: map[string]string{}, provider: OvnProvider, subnet: vpcSubnet, vpcTunnelKey: 5678},
		{name: "invalid key", annotations: map[string]string{TunnelKeyAnnotation: "invalid"}, provider: OvnProvider, subnet: vpcSubnet, vpcTunnelKey: 5678},
		{name: "VPC key not synced", annotations: map[string]string{TunnelKeyAnnotation: "0"}, provider: OvnProvider, subnet: vpcSubnet},
		{name: "vlan requires no key", annotations: map[string]string{}, provider: OvnProvider, subnet: vlanSubnet, want: true},
		{name: "non-OVN or unidentified subnet requires no key", annotations: map[string]string{}, provider: "net1.default", subnet: nil, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTunnelKeyAnnotationValidForSubnet(tt.annotations, tt.provider, tt.subnet, tt.vpcTunnelKey); got != tt.want {
				t.Errorf("IsTunnelKeyAnnotationValidForSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOvnProvider(t *testing.T) {
	testCases := []struct {
		name     string
		provider string
		expected bool
	}{
		{
			name:     "empty provider",
			provider: "",
			expected: true,
		},
		{
			name:     "ovn provider",
			provider: OvnProvider,
			expected: true,
		},
		{
			name:     "ovn provider with namespace",
			provider: "namespace.cluster.ovn",
			expected: true,
		},
		{
			name:     "non ovn provider",
			provider: "other-provider",
			expected: false,
		},
		{
			name:     "invalid provider format",
			provider: "invalid.format",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsOvnProvider(tc.provider)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
