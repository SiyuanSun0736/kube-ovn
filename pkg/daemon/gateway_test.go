package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestGetCidrByProtocol(t *testing.T) {
	cases := []struct {
		name     string
		cidr     string
		protocol string
		wantErr  bool
		expetced string
	}{{
		name:     "ipv4 only",
		cidr:     "1.1.1.0/24",
		protocol: kubeovnv1.ProtocolIPv4,
		expetced: "1.1.1.0/24",
	}, {
		name:     "ipv6 only",
		cidr:     "2001:db8::/120",
		protocol: kubeovnv1.ProtocolIPv6,
		expetced: "2001:db8::/120",
	}, {
		name:     "get ipv4 from ipv6",
		cidr:     "2001:db8::/120",
		protocol: kubeovnv1.ProtocolIPv4,
	}, {
		name:     "get ipv4 from dual stack",
		cidr:     "1.1.1.0/24,2001:db8::/120",
		protocol: kubeovnv1.ProtocolIPv4,
		expetced: "1.1.1.0/24",
	}, {
		name:     "get ipv6 from ipv4",
		cidr:     "1.1.1.0/24",
		protocol: kubeovnv1.ProtocolIPv6,
	}, {
		name:     "get ipv6 from dual stack",
		cidr:     "1.1.1.0/24,2001:db8::/120",
		protocol: kubeovnv1.ProtocolIPv6,
		expetced: "2001:db8::/120",
	}, {
		name:     "invalid cidr",
		cidr:     "foo bar",
		protocol: kubeovnv1.ProtocolIPv4,
		wantErr:  true,
	}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := getCidrByProtocol(c.cidr, c.protocol)
			if (err != nil) != c.wantErr {
				t.Errorf("getCidrByProtocol(%q, %q) error = %v, wantErr = %v", c.cidr, c.protocol, err, c.wantErr)
			}
			require.Equal(t, c.expetced, got)
		})
	}
}

func TestGetTunnelPort(t *testing.T) {
	cases := []struct {
		name        string
		networkType string
		expected    int
	}{
		{
			name:        "geneve network type",
			networkType: util.NetworkTypeGeneve,
			expected:    6081,
		},
		{
			name:        "vxlan network type",
			networkType: util.NetworkTypeVxlan,
			expected:    4789,
		},
		{
			name:        "vlan network type",
			networkType: util.NetworkTypeVlan,
			expected:    0,
		},
		{
			name:        "stt network type",
			networkType: util.NetworkTypeStt,
			expected:    0,
		},
		{
			name:        "empty network type defaults to 0",
			networkType: "",
			expected:    0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Controller{
				config: &Configuration{
					NetworkType: tc.networkType,
				},
			}
			require.Equal(t, tc.expected, c.getTunnelPort())
		})
	}

	t.Run("nil config returns 0", func(t *testing.T) {
		c := &Controller{}
		require.Equal(t, 0, c.getTunnelPort())
	})
}
