package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

func TestNatGatewayDataplaneMode(t *testing.T) {
	for _, tt := range []struct {
		spec, want string
		invalid    bool
	}{
		{"", "", false},
		{"stateless-nft", "stateless-nft", false},
		{"bad", "", true},
	} {
		gw := &kubeovnv1.VpcNatGateway{ObjectMeta: metav1.ObjectMeta{Name: "gw"}, Spec: kubeovnv1.VpcNatGatewaySpec{DataplaneMode: tt.spec}}
		mode, err := NatGatewayDataplaneMode(gw)
		if tt.invalid {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.want, mode)
		}
	}
}

func TestStatelessNATIPv4Host(t *testing.T) {
	for _, value := range []string{"10.0.0.2", "10.0.0.2/32"} {
		ip, err := StatelessNATIPv4Host(value)
		require.NoError(t, err)
		require.Equal(t, "10.0.0.2", ip)
	}
	for _, value := range []string{"", "10.0.0.0/24", "::1", "::ffff:10.0.0.2", "10.0.0.2/31", "garbage"} {
		_, err := StatelessNATIPv4Host(value)
		require.Error(t, err, value)
	}
}
