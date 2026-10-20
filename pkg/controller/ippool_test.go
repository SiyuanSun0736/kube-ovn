package controller

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/set"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
	"k8s.io/utils/keymutex"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/internal"
	"github.com/kubeovn/kube-ovn/pkg/ipam"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestExpandIPPoolAddresses(t *testing.T) {
	addresses, err := util.ExpandIPPoolAddresses([]string{
		"10.0.0.1",
		"2001:db8::1",
		"192.168.1.0/24",
		"10.0.0.1", // duplicate should be removed
		" 2001:db8::1 ",
	})
	require.NoError(t, err)
	require.Equal(t, []string{
		"10.0.0.1/32",
		"192.168.1.0/24",
		"2001:db8::1/128",
	}, addresses)
}

func TestExpandIPPoolAddressesRange(t *testing.T) {
	addresses, err := util.ExpandIPPoolAddresses([]string{"10.0.0.0..10.0.0.3"})
	require.NoError(t, err)
	require.Equal(t, []string{"10.0.0.0/30"}, addresses)

	addresses, err = util.ExpandIPPoolAddresses([]string{"10.0.0.1..10.0.0.5"})
	require.NoError(t, err)
	require.Equal(t, []string{
		"10.0.0.1/32",
		"10.0.0.2/31",
		"10.0.0.4/31",
	}, addresses)

	addresses, err = util.ExpandIPPoolAddresses([]string{"2001:db8::1..2001:db8::4"})
	require.NoError(t, err)
	require.Equal(t, []string{
		"2001:db8::1/128",
		"2001:db8::2/127",
		"2001:db8::4/128",
	}, addresses)
}

func TestExpandIPPoolAddressesInvalid(t *testing.T) {
	_, err := util.ExpandIPPoolAddresses([]string{"10.0.0.1..2001:db8::1"})
	require.Error(t, err)

	_, err = util.ExpandIPPoolAddresses([]string{"foo"})
	require.Error(t, err)
}

func TestIPPoolAddressSetName(t *testing.T) {
	require.Equal(t, "foo.bar", util.IPPoolAddressSetName("foo-bar"))
	require.Equal(t, "123pool", util.IPPoolAddressSetName("123pool"))
}

// Additional comprehensive tests for IPPool utilities

func TestExpandIPPoolAddressesEdgeCases(t *testing.T) {
	t.Run("Empty and whitespace entries", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{
			"",
			"   ",
			"\t",
			"10.0.0.1",
		})
		require.NoError(t, err)
		require.Len(t, addresses, 1)
		require.Equal(t, "10.0.0.1/32", addresses[0])
	})

	t.Run("CIDR normalization", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{
			"192.168.1.5/24", // non-canonical, should normalize to 192.168.1.0/24
		})
		require.NoError(t, err)
		require.Equal(t, []string{"192.168.1.0/24"}, addresses)
	})

	t.Run("Large range", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{
			"10.0.0.0..10.0.0.255",
		})
		require.NoError(t, err)
		require.Equal(t, []string{"10.0.0.0/24"}, addresses)
	})

	t.Run("Complex mixed input", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{
			"10.0.0.1",
			"10.0.0.5..10.0.0.10",
			"192.168.1.0/24",
			"2001:db8::1",
			"2001:db8::100..2001:db8::103",
		})
		require.NoError(t, err)
		require.Greater(t, len(addresses), 5) // Should have multiple CIDRs
		require.Contains(t, addresses, "10.0.0.1/32")
		require.Contains(t, addresses, "192.168.1.0/24")
		require.Contains(t, addresses, "2001:db8::1/128")
	})
}

func TestExpandIPPoolAddressesRangeEdgeCases(t *testing.T) {
	t.Run("Unaligned IPv4 range", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{"10.0.0.1..10.0.0.5"})
		require.NoError(t, err)
		require.Equal(t, []string{
			"10.0.0.1/32",
			"10.0.0.2/31",
			"10.0.0.4/31",
		}, addresses)
	})

	t.Run("Power of 2 range", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{"10.0.0.0..10.0.0.7"})
		require.NoError(t, err)
		require.Equal(t, []string{"10.0.0.0/29"}, addresses)
	})

	t.Run("Single IP range", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{"10.0.0.1..10.0.0.1"})
		require.NoError(t, err)
		require.Equal(t, []string{"10.0.0.1/32"}, addresses)
	})

	t.Run("IPv6 unaligned range", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{"2001:db8::1..2001:db8::4"})
		require.NoError(t, err)
		require.Equal(t, []string{
			"2001:db8::1/128",
			"2001:db8::2/127",
			"2001:db8::4/128",
		}, addresses)
	})
}

func TestExpandIPPoolAddressesErrorConditions(t *testing.T) {
	t.Run("Completely invalid input", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddresses([]string{"totally-invalid"})
		require.Error(t, err)
	})

	t.Run("Invalid range - missing end", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddresses([]string{"10.0.0.1.."})
		require.Error(t, err)
	})

	t.Run("Invalid range - missing start", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddresses([]string{"..10.0.0.1"})
		require.Error(t, err)
	})

	t.Run("Invalid CIDR - out of range prefix", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddresses([]string{"10.0.0.0/33"})
		require.Error(t, err)
	})

	t.Run("Invalid CIDR - IPv6 out of range", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddresses([]string{"2001:db8::/129"})
		require.Error(t, err)
	})

	t.Run("Mixed IP families in range", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddresses([]string{"10.0.0.1..2001:db8::1"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "mixes IPv4 and IPv6")
	})

	t.Run("Range with end less than start", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddresses([]string{"10.0.0.100..10.0.0.1"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "start is greater than end")
	})
}

func TestCanonicalizeIPPoolEntries(t *testing.T) {
	t.Run("Basic functionality", func(t *testing.T) {
		result, err := util.CanonicalizeIPPoolEntries([]string{
			"10.0.0.1",
			"192.168.1.0/24",
		})
		require.NoError(t, err)
		require.True(t, result["10.0.0.1/32"])
		require.True(t, result["192.168.1.0/24"])
		require.False(t, result["10.0.0.2/32"])
	})

	t.Run("Empty input", func(t *testing.T) {
		result, err := util.CanonicalizeIPPoolEntries([]string{})
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("Duplicate detection", func(t *testing.T) {
		result, err := util.CanonicalizeIPPoolEntries([]string{
			"10.0.0.1",
			"10.0.0.1",
		})
		require.NoError(t, err)
		require.Len(t, result, 1)
		require.True(t, result["10.0.0.1/32"])
	})
}

func TestNormalizeAddressSetEntries(t *testing.T) {
	t.Run("Standard OVN format", func(t *testing.T) {
		result := util.NormalizeAddressSetEntries(`"10.0.0.1/32" "10.0.0.2/32"`)
		require.Equal(t, set.New("10.0.0.1/32", "10.0.0.2/32"), result)
	})

	t.Run("Extra whitespace", func(t *testing.T) {
		result := util.NormalizeAddressSetEntries(`  "10.0.0.1/32"   "10.0.0.2/32"  `)
		require.Equal(t, set.New("10.0.0.1/32", "10.0.0.2/32"), result)
	})

	t.Run("Empty input", func(t *testing.T) {
		result := util.NormalizeAddressSetEntries("")
		require.Empty(t, result)
	})

	t.Run("No quotes", func(t *testing.T) {
		result := util.NormalizeAddressSetEntries("10.0.0.1/32 10.0.0.2/32")
		require.Equal(t, set.New("10.0.0.1/32", "10.0.0.2/32"), result)
	})

	t.Run("Mixed formats", func(t *testing.T) {
		result := util.NormalizeAddressSetEntries(`"192.168.1.0/24" "2001:db8::1/128"`)
		require.Equal(t, set.New("192.168.1.0/24", "2001:db8::1/128"), result)
	})
}

func TestNormalizeIP(t *testing.T) {
	t.Run("Standard IPv4", func(t *testing.T) {
		ip, err := util.NormalizeIP("192.168.1.1")
		require.NoError(t, err)
		require.Equal(t, "192.168.1.1", ip.String())
		require.NotNil(t, ip.To4())
	})

	t.Run("Standard IPv6", func(t *testing.T) {
		ip, err := util.NormalizeIP("2001:db8::1")
		require.NoError(t, err)
		require.Equal(t, "2001:db8::1", ip.String())
		require.Nil(t, ip.To4())
	})

	t.Run("IPv6 full form", func(t *testing.T) {
		ip, err := util.NormalizeIP("2001:0db8:0000:0000:0000:0000:0000:0001")
		require.NoError(t, err)
		require.Equal(t, "2001:db8::1", ip.String())
	})

	t.Run("Whitespace handling", func(t *testing.T) {
		ip, err := util.NormalizeIP("  10.0.0.1  ")
		require.NoError(t, err)
		require.Equal(t, "10.0.0.1", ip.String())
	})

	t.Run("Invalid IP", func(t *testing.T) {
		_, err := util.NormalizeIP("999.999.999.999")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid IP address")
	})

	t.Run("Not an IP", func(t *testing.T) {
		_, err := util.NormalizeIP("hostname")
		require.Error(t, err)
	})

	t.Run("IPv4 loopback", func(t *testing.T) {
		ip, err := util.NormalizeIP("127.0.0.1")
		require.NoError(t, err)
		require.Equal(t, "127.0.0.1", ip.String())
	})

	t.Run("IPv6 loopback", func(t *testing.T) {
		ip, err := util.NormalizeIP("::1")
		require.NoError(t, err)
		require.Equal(t, "::1", ip.String())
	})
}

func TestExpandIPPoolAddressesForOVNIntegration(t *testing.T) {
	// These tests use ExpandIPPoolAddressesForOVN which enforces OVN address set limitation
	t.Run("Mixed IPv4 and IPv6 - should fail", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddressesForOVN([]string{
			"10.0.0.1",
			"2001:db8::1",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "mixed IPv4 and IPv6 addresses are not supported")
	})

	t.Run("Pure IPv4 pool", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{
			"192.168.1.0/30",
			"10.0.0.1..10.0.0.5",
		})
		require.NoError(t, err)
		require.NotEmpty(t, addresses)
		// Verify all are IPv4
		for _, addr := range addresses {
			require.NotContains(t, addr, ":")
		}
		// CIDR should be preserved
		require.Contains(t, addresses, "192.168.1.0/30")
	})

	t.Run("Pure IPv6 pool", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{
			"2001:db8::/126",
			"fd00::1..fd00::3",
		})
		require.NoError(t, err)
		require.NotEmpty(t, addresses)
		// Verify all are IPv6
		for _, addr := range addresses {
			require.Contains(t, addr, ":")
		}
		// CIDR should be preserved
		require.Contains(t, addresses, "2001:db8::/126")
	})

	t.Run("Single IPs are simplified", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{
			"10.0.0.1",
			"10.0.0.2",
		})
		require.NoError(t, err)
		require.Len(t, addresses, 2)
		// Should not have /32 suffix
		require.Contains(t, addresses, "10.0.0.1")
		require.Contains(t, addresses, "10.0.0.2")
		require.NotContains(t, addresses, "10.0.0.1/32")
		require.NotContains(t, addresses, "10.0.0.2/32")
	})

	t.Run("Mixed in range notation", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddressesForOVN([]string{
			"10.0.0.1..10.0.0.5",
			"fd00::1",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "mixed IPv4 and IPv6")
	})

	t.Run("Mixed in CIDR notation", func(t *testing.T) {
		_, err := util.ExpandIPPoolAddressesForOVN([]string{
			"192.168.1.0/24",
			"2001:db8::/64",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "mixed IPv4 and IPv6")
	})

	t.Run("Range expansion with simplification", func(t *testing.T) {
		// Range that expands to /32 should be simplified
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{
			"10.0.0.1..10.0.0.1", // Single IP range
		})
		require.NoError(t, err)
		require.Len(t, addresses, 1)
		require.Equal(t, "10.0.0.1", addresses[0])
		require.NotContains(t, addresses[0], "/32")
	})

	t.Run("Range with multiple CIDRs - some simplified", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{
			"10.0.0.1..10.0.0.5",
		})
		require.NoError(t, err)
		require.NotEmpty(t, addresses)
		// Should contain simplified /32 and non-simplified /31
		require.Contains(t, addresses, "10.0.0.1")    // Simplified from /32
		require.Contains(t, addresses, "10.0.0.2/31") // Not /32, keep as-is
		require.Contains(t, addresses, "10.0.0.4/31") // Not /32, keep as-is
	})

	t.Run("Empty input", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{})
		require.NoError(t, err)
		require.Empty(t, addresses)
	})

	t.Run("Only whitespace entries", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{"", "  ", "\t"})
		require.NoError(t, err)
		require.Empty(t, addresses)
	})

	t.Run("IPv6 single IP simplified", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{"2001:db8::1"})
		require.NoError(t, err)
		require.Len(t, addresses, 1)
		require.Equal(t, "2001:db8::1", addresses[0])
		require.NotContains(t, addresses[0], "/128")
	})

	t.Run("IPv6 range with /128 simplified", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{"fd00::1..fd00::1"})
		require.NoError(t, err)
		require.Len(t, addresses, 1)
		require.Equal(t, "fd00::1", addresses[0])
	})

	t.Run("Duplicate IPs deduplicated and simplified", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{
			"10.0.0.1",
			"10.0.0.1",
			" 10.0.0.1 ",
		})
		require.NoError(t, err)
		require.Len(t, addresses, 1)
		require.Equal(t, "10.0.0.1", addresses[0])
	})

	t.Run("CIDR normalization preserved", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddressesForOVN([]string{"192.168.1.5/24"})
		require.NoError(t, err)
		require.Len(t, addresses, 1)
		require.Equal(t, "192.168.1.0/24", addresses[0])
	})
}

func TestExpandIPPoolAddressesGeneralUse(t *testing.T) {
	// ExpandIPPoolAddresses (without ForOVN suffix) allows mixed IP families
	t.Run("Mixed IPv4 and IPv6 - allowed for general use", func(t *testing.T) {
		addresses, err := util.ExpandIPPoolAddresses([]string{
			"10.0.0.1",
			"2001:db8::1",
		})
		require.NoError(t, err)
		require.NotEmpty(t, addresses)
		require.Contains(t, addresses, "10.0.0.1/32")
		require.Contains(t, addresses, "2001:db8::1/128")
	})
}

func TestEnqueueAddIPPoolRoutesTerminatingToDeleteQueue(t *testing.T) {
	c := &Controller{
		addOrUpdateIPPoolQueue: newTypedRateLimitingQueue[string]("AddIPPool", nil),
		deleteIPPoolQueue:      newTypedRateLimitingQueue[*kubeovnv1.IPPool]("DeleteIPPool", nil),
	}
	t.Cleanup(func() {
		c.addOrUpdateIPPoolQueue.ShutDown()
		c.deleteIPPoolQueue.ShutDown()
	})

	// After a controller restart the informer (re)list surfaces a terminating
	// ippool only via the add handler; it must go to the delete queue or the
	// finalizer is never removed.
	now := metav1.Now()
	terminating := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "stuck",
			DeletionTimestamp: &now,
			Finalizers:        []string{util.KubeOVNControllerFinalizer},
		},
	}
	c.enqueueAddIPPool(terminating)
	require.Equal(t, 0, c.addOrUpdateIPPoolQueue.Len())
	require.Equal(t, 1, c.deleteIPPoolQueue.Len())
	item, _ := c.deleteIPPoolQueue.Get()
	require.Same(t, terminating, item)
	c.deleteIPPoolQueue.Done(item)

	// A live ippool keeps going to the add/update queue.
	live := &kubeovnv1.IPPool{ObjectMeta: metav1.ObjectMeta{Name: "live"}}
	c.enqueueAddIPPool(live)
	require.Equal(t, 1, c.addOrUpdateIPPoolQueue.Len())
	require.Equal(t, 0, c.deleteIPPoolQueue.Len())
}

func TestReserveIPPoolIPs(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet1"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.0.0.0/28",
			Gateway:    "10.0.0.1",
			ExcludeIps: []string{"10.0.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool1"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 3},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	pool := ippool.DeepCopy()
	apply := func() []string {
		t.Helper()
		ips, err := ctrl.reserveIPPoolIPs(pool)
		require.NoError(t, err)
		// the reservation holds the addresses, activating them hands them out;
		// recording them in the status is the caller's job
		require.NoError(t, ctrl.ipam.AddOrUpdateIPPool(pool.Spec.Subnet, pool.Name, ips))
		pool.Status.IPs = ips
		return ips
	}

	// the reservation covers the picked addresses
	require.Equal(t, []string{"10.0.0.2", "10.0.0.3", "10.0.0.4"}, apply())

	// growing keeps the existing addresses and appends newly picked ones
	pool.Spec.IPCount = 5
	require.Equal(t, []string{"10.0.0.2", "10.0.0.3", "10.0.0.4", "10.0.0.5", "10.0.0.6"}, apply())

	// shrinking releases the most recently added addresses first
	pool.Spec.IPCount = 2
	require.Equal(t, []string{"10.0.0.2", "10.0.0.3"}, apply())

	// not enough free addresses in the subnet, which only holds 13 of them
	pool.Spec.IPCount = 20
	_, err = ctrl.reserveIPPoolIPs(pool)
	require.ErrorContains(t, err, "does not have enough free addresses")

	// a failed reservation is not recorded and did not touch the committed one
	require.Equal(t, []string{"10.0.0.2", "10.0.0.3"}, pool.Status.IPs)
	pool.Spec.IPCount = 2
	require.Equal(t, []string{"10.0.0.2", "10.0.0.3"}, apply())

	// spec.ips is reserved as it is
	pool.Spec.IPCount = 0
	pool.Spec.IPs = []string{"10.0.0.10"}
	require.Equal(t, []string{"10.0.0.10"}, apply())
}

func TestHandleAddOrUpdateIPPoolWithIPCount(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet2"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.1.0.0/24",
			Gateway:    "10.1.0.1",
			ExcludeIps: []string{"10.1.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool2"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	fakeController.mockOvnClient.EXPECT().DeleteAddressSet(gomock.Any()).Return(nil).AnyTimes()
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))

	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []string{"10.1.0.2", "10.1.0.3"}, got.Status.IPs)
	require.True(t, got.Status.IsReady())

	// the reconcile left the membership allocatable
	ip, _, _, err := ctrl.ipam.GetRandomAddress("pod1.default", "nic1", nil, subnet.Name, ippool.Name, nil, true)
	require.NoError(t, err)
	require.Contains(t, got.Status.IPs, ip)
}

func TestReserveIPPoolIPsDualStack(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet3"},
		Spec: kubeovnv1.SubnetSpec{
			CIDRBlock:  "10.2.0.0/24,fd00::/120",
			Gateway:    "10.2.0.1,fd00::1",
			ExcludeIps: []string{"10.2.0.1", "fd00::1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool3"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	// ipCount applies to each IP family
	pool := ippool.DeepCopy()
	ips, err := ctrl.reserveIPPoolIPs(pool)
	require.NoError(t, err)
	require.Equal(t, []string{"10.2.0.2", "10.2.0.3", "fd00::2", "fd00::3"}, ips)
	pool.Status.IPs = ips

	pool.Spec.IPCount = 1
	ips, err = ctrl.reserveIPPoolIPs(pool)
	require.NoError(t, err)
	require.Equal(t, []string{"10.2.0.2", "fd00::2"}, ips)
}

// TestPatchIPPoolStatusKeepsIPs guards the regression where a statistics-only status
// update, built from an informer cache that has not caught up with the last reconcile
// yet, rolled status.ips back to its previous value.
func TestPatchIPPoolStatusKeepsIPs(t *testing.T) {
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool4"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: "subnet4", IPCount: 2},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{IPPools: []*kubeovnv1.IPPool{ippool}})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ippools := ctrl.config.KubeOvnClient.KubeovnV1().IPPools()

	ips := []string{"10.3.0.2", "10.3.0.3"}
	require.NoError(t, ctrl.patchIPPoolIPs(ippool, ips))
	require.Equal(t, ips, ippool.Status.IPs)
	got, err := ippools.Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, ips, got.Status.IPs)

	// a stale copy, as returned by the lister before the patch above is observed: it
	// still carries the previous selection, which an unguarded status patch would write
	// back - a nil would not, json omitempty drops it from the merge patch
	stale := ippool.DeepCopy()
	stale.Status.IPs = []string{"10.3.0.2"}
	stale.Status.V4UsingIPs = internal.BigInt{Int: *big.NewInt(1)}
	require.NoError(t, ctrl.patchIPPoolStatus(stale))

	got, err = ippools.Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, ips, got.Status.IPs)
	require.Equal(t, int64(1), got.Status.V4UsingIPs.Int64())
}

// TestPatchIPPoolIPsIgnoresStaleCache covers a reconcile whose informer view carries the
// new spec but not the status the previous reconcile persisted: skipping the write
// because the cached list happens to match would leave the persisted membership ahead
// of the one in IPAM, with no event left to correct it.
func TestPatchIPPoolIPsIgnoresStaleCache(t *testing.T) {
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool14"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: "subnet14", IPCount: 3},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{IPPools: []*kubeovnv1.IPPool{ippool}})
	require.NoError(t, err)
	ctrl := fakeController.fakeController

	// the previous reconcile grew the pool to three addresses
	require.NoError(t, ctrl.patchIPPoolIPs(ippool.DeepCopy(), []string{"10.14.0.2", "10.14.0.3", "10.14.0.4"}))

	// this reconcile shrinks it back, working on an object version that carries the new
	// spec but the status from before the growth
	stale := ippool.DeepCopy()
	stale.Spec.IPCount = 2
	stale.Status.IPs = []string{"10.14.0.2", "10.14.0.3"}
	require.NoError(t, ctrl.patchIPPoolIPs(stale, []string{"10.14.0.2", "10.14.0.3"}))

	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []string{"10.14.0.2", "10.14.0.3"}, got.Status.IPs)
}

func TestEnqueueUpdateIPPool(t *testing.T) {
	ctrl := newFakeController(t).fakeController
	ctrl.addOrUpdateIPPoolQueue = newTypedRateLimitingQueue[string]("TestAddIPPool", nil)
	defer ctrl.addOrUpdateIPPoolQueue.ShutDown()

	old := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool5"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: "subnet5", IPCount: 2},
	}

	// a status-only change must not be enqueued, otherwise recording status.ips would
	// trigger an endless reconcile loop
	unchanged := old.DeepCopy()
	unchanged.Status.IPs = []string{"10.4.0.2", "10.4.0.3"}
	ctrl.enqueueUpdateIPPool(old, unchanged)
	require.Equal(t, 0, ctrl.addOrUpdateIPPoolQueue.Len())

	// changing the requested number of addresses must trigger a reconcile
	resized := old.DeepCopy()
	resized.Spec.IPCount = 3
	ctrl.enqueueUpdateIPPool(old, resized)
	require.Equal(t, 1, ctrl.addOrUpdateIPPoolQueue.Len())
	key, _ := ctrl.addOrUpdateIPPoolQueue.Get()
	require.Equal(t, old.Name, key)
}

func TestHandleAddOrUpdateIPPoolResolveFailure(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet6"},
		Spec:       kubeovnv1.SubnetSpec{Protocol: kubeovnv1.ProtocolIPv4, CIDRBlock: "10.5.0.0/24"},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool6"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2},
		Status:     kubeovnv1.IPPoolStatus{IPs: []string{"10.5.0.2", "10.5.0.3"}},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)

	// the subnet is not registered in IPAM yet, so no address can be reserved
	err = ctrl.handleAddOrUpdateIPPool(ippool.Name)
	require.ErrorContains(t, err, "does not exist in IPAM")

	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, ippool.Status.IPs, got.Status.IPs, "a failed reconcile must not drop the recorded addresses")
	require.False(t, got.Status.IsReady())
}

func TestHandleAddOrUpdateIPPoolAddressSetWithIPCount(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet7"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.6.0.0/24",
			Gateway:    "10.6.0.1",
			ExcludeIps: []string{"10.6.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool7"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2, EnableAddressSet: true},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	asName := util.IPPoolAddressSetName(ippool.Name)
	mockOvnClient := fakeController.mockOvnClient
	mockOvnClient.EXPECT().CreateAddressSet(asName, map[string]string{ippoolKey: ippool.Name}).Return(nil)
	// the address set must hold the addresses picked for spec.ipCount, not spec.ips
	mockOvnClient.EXPECT().AddressSetUpdateAddress(asName, "10.6.0.2", "10.6.0.3").Return(nil)

	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))

	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []string{"10.6.0.2", "10.6.0.3"}, got.Status.IPs)
}

// TestHandleAddOrUpdateIPPoolAddressSetFailureKeepsMembershipReserved pins the exposure
// invariant: while the address set cannot be reconciled, the new membership must not
// become allocatable - neither through the pool nor, now that it is reserved, through
// the subnet, so that the published address set can never describe a foreign pod.
func TestHandleAddOrUpdateIPPoolAddressSetFailureKeepsMembershipReserved(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet12"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.12.0.0/24",
			Gateway:    "10.12.0.1",
			ExcludeIps: []string{"10.12.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool12"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2, EnableAddressSet: true},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	asName := util.IPPoolAddressSetName(ippool.Name)
	fakeController.mockOvnClient.EXPECT().CreateAddressSet(asName, gomock.Any()).Return(errors.New("ovn is down"))

	err = ctrl.handleAddOrUpdateIPPool(ippool.Name)
	require.ErrorContains(t, err, "ovn is down")

	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []string{"10.12.0.2", "10.12.0.3"}, got.Status.IPs,
		"the reservation is persisted for the retry and for the restore after a restart")
	require.False(t, got.Status.IsReady())

	// the pool hands out nothing, and the subnet hands out neither of the two addresses
	_, _, _, err = ctrl.ipam.GetRandomAddress("pod1.default", "nic1", nil, subnet.Name, ippool.Name, nil, true)
	require.ErrorIs(t, err, ipam.ErrNoAvailable)
	ip, _, _, err := ctrl.ipam.GetRandomAddress("pod2.default", "nic2", nil, subnet.Name, "", nil, true)
	require.NoError(t, err)
	require.Equal(t, "10.12.0.4", ip, "the reserved addresses stay out of reach while the address set is missing")
}

// TestInitIPAMRestoresIPCountPoolAsReservation covers the startup restore: a membership
// recorded in status.ips whose data plane may never have been published must come back
// as a reservation, so a pod cannot be handed an address the address set does not
// describe. The reconcile that follows activates it.
func TestInitIPAMRestoresIPCountPoolAsReservation(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet15"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.15.0.0/24",
			Gateway:    "10.15.0.1",
			ExcludeIps: []string{"10.15.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool15"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2},
		Status:     kubeovnv1.IPPoolStatus{IPs: []string{"10.15.0.2", "10.15.0.3"}},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	fakeController.mockOvnClient.EXPECT().DeleteAddressSet(gomock.Any()).Return(nil).AnyTimes()
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	require.NoError(t, ctrl.initIPPoolIPAM(ippool))
	_, _, _, err = ctrl.ipam.GetRandomAddress("pod1.default", "nic1", nil, subnet.Name, ippool.Name, nil, true)
	require.ErrorIs(t, err, ipam.ErrNoAvailable)
	// the addresses are reserved all the same, the subnet does not hand them out either
	ip, _, _, err := ctrl.ipam.GetRandomAddress("pod2.default", "nic2", nil, subnet.Name, "", nil, true)
	require.NoError(t, err)
	require.Equal(t, "10.15.0.4", ip)

	// the reconcile republishes the data plane and activates the same addresses
	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))
	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, ippool.Status.IPs, got.Status.IPs)
	ip, _, _, err = ctrl.ipam.GetRandomAddress("pod1.default", "nic1", nil, subnet.Name, ippool.Name, nil, true)
	require.NoError(t, err)
	require.Contains(t, ippool.Status.IPs, ip)
}

func TestReserveIPPoolIPsStaleAddresses(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet8"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.7.0.0/24",
			Gateway:    "10.7.0.1",
			ExcludeIps: []string{"10.7.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool8"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 3},
		// 10.8.0.2 is left over from the previous CIDR of the subnet
		Status: kubeovnv1.IPPoolStatus{IPs: []string{"10.7.0.5", "10.8.0.2", "10.7.0.6"}},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	// addresses outside the current CIDR are dropped, the missing one is picked again
	pool := ippool.DeepCopy()
	ips, err := ctrl.reserveIPPoolIPs(pool)
	require.NoError(t, err)
	require.Equal(t, []string{"10.7.0.5", "10.7.0.6", "10.7.0.2"}, ips)

	// a pool whose subnet does not exist yet cannot be resolved
	pool.Spec.Subnet = "noSuchSubnet"
	_, err = ctrl.reserveIPPoolIPs(pool)
	require.Error(t, err)
}

// TestReserveIPPoolIPsExceedingMaxCount covers pools created before the CRD gained its
// upper bound: such a pool is rejected instead of carving a large chunk out of the subnet.
func TestReserveIPPoolIPsExceedingMaxCount(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet9"},
		Spec:       kubeovnv1.SubnetSpec{Protocol: kubeovnv1.ProtocolIPv4, CIDRBlock: "10.9.0.0/26", Gateway: "10.9.0.1"},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool9"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: maxIPPoolIPCount + 1},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, nil))

	_, err = ctrl.reserveIPPoolIPs(ippool.DeepCopy())
	require.ErrorContains(t, err, "more than the maximum")

	// the limit itself is allowed
	pool := ippool.DeepCopy()
	pool.Spec.IPCount = maxIPPoolIPCount
	ips, err := ctrl.reserveIPPoolIPs(pool)
	require.NoError(t, err)
	require.Len(t, ips, maxIPPoolIPCount)
}

// TestReserveIPPoolIPsRecoversHeldSelection covers a retry whose informer view of the
// pool lags behind the last reconcile: status.ips is not updated yet, so the reserved
// addresses must be recovered from what IPAM holds instead of being reported as
// unavailable.
func TestReserveIPPoolIPsRecoversHeldSelection(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet10"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.10.0.0/30",
			Gateway:    "10.10.0.1",
			ExcludeIps: []string{"10.10.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool10"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 1},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	// the subnet holds a single usable address; the pool reserved and activated it
	ips, err := ctrl.reserveIPPoolIPs(ippool.DeepCopy())
	require.NoError(t, err)
	require.Equal(t, []string{"10.10.0.2"}, ips)
	require.NoError(t, ctrl.ipam.AddOrUpdateIPPool(subnet.Name, ippool.Name, ips))

	// the status update did not reach the informer yet, so the retry reserves with
	// status.ips still empty: the committed address is recovered from what IPAM
	// holds, not requested as a new pick
	ips, err = ctrl.reserveIPPoolIPs(ippool.DeepCopy())
	require.NoError(t, err)
	require.Equal(t, []string{"10.10.0.2"}, ips)
}

// TestHandleAddOrUpdateIPPoolSwitchSelectionMode covers switching between spec.ipCount
// and spec.ips: leaving count mode clears the recorded addresses, and re-entering it
// keeps the pool's current addresses instead of resurrecting the stale selection.
func TestHandleAddOrUpdateIPPoolSwitchSelectionMode(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet11"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.11.0.0/28",
			Gateway:    "10.11.0.1",
			ExcludeIps: []string{"10.11.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool11"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	fakeController.mockOvnClient.EXPECT().DeleteAddressSet(gomock.Any()).Return(nil).AnyTimes()
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))
	ippools := ctrl.config.KubeOvnClient.KubeovnV1().IPPools()

	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))
	got, err := ippools.Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []string{"10.11.0.2", "10.11.0.3"}, got.Status.IPs)

	// switch to spec.ips: the former addresses are released and the recorded ones
	// are cleared, they must not be mistaken for the pool's members later
	got.Spec.IPCount = 0
	got.Spec.IPs = []string{"10.11.0.10"}
	_, err = ippools.Update(context.Background(), got, metav1.UpdateOptions{})
	require.NoError(t, err)
	waitIPPoolCached(t, ctrl, ippool.Name, func(cached *kubeovnv1.IPPool) bool {
		return cached.Spec.IPCount == 0 && len(cached.Spec.IPs) == 1
	})
	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))
	got, err = ippools.Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, got.Status.IPs)

	// another pool claims one of the former addresses
	other := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool12"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPs: []string{"10.11.0.3"}},
	}
	_, err = ippools.Create(context.Background(), other, metav1.CreateOptions{})
	require.NoError(t, err)
	waitIPPoolCached(t, ctrl, other.Name, func(*kubeovnv1.IPPool) bool { return true })
	require.NoError(t, ctrl.handleAddOrUpdateIPPool(other.Name))

	// switch back to spec.ipCount: the pool keeps its current address and must not
	// collide with the other pool's claim on the stale selection
	got.Spec.IPCount = 2
	got.Spec.IPs = nil
	_, err = ippools.Update(context.Background(), got, metav1.UpdateOptions{})
	require.NoError(t, err)
	waitIPPoolCached(t, ctrl, ippool.Name, func(cached *kubeovnv1.IPPool) bool {
		return cached.Spec.IPCount == 2 && len(cached.Spec.IPs) == 0
	})
	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))
	got, err = ippools.Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	// keeps its own current address, skips the address claimed by the other pool
	require.Equal(t, []string{"10.11.0.10", "10.11.0.2"}, got.Status.IPs)
}

// waitIPPoolCached waits until the informer cache has observed an update made through
// the client: the handler resolves pools from the lister, which lags behind the
// tracker.
func waitIPPoolCached(t *testing.T, ctrl *Controller, name string, condition func(*kubeovnv1.IPPool) bool) {
	t.Helper()
	require.Eventually(t, func() bool {
		cached, err := ctrl.ippoolLister.Get(name)
		return err == nil && condition(cached)
	}, 5*time.Second, 10*time.Millisecond)
}

// TestHandleAddOrUpdateIPPoolSpecIPsUnaffected pins the scope of the reservation: a pool
// listing spec.ips picks nothing, so it is not reserved and keeps the behaviour it had
// before spec.ipCount existed - its addresses are committed by the activation, and a
// failed address set leaves the membership exactly as the previous reconcile left it.
func TestHandleAddOrUpdateIPPoolSpecIPsUnaffected(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet16"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.16.0.0/24",
			Gateway:    "10.16.0.1",
			ExcludeIps: []string{"10.16.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool16"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPs: []string{"10.16.0.2"}, EnableAddressSet: true},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	asName := util.IPPoolAddressSetName(ippool.Name)
	fakeController.mockOvnClient.EXPECT().CreateAddressSet(asName, gomock.Any()).Return(errors.New("ovn is down"))
	require.ErrorContains(t, ctrl.handleAddOrUpdateIPPool(ippool.Name), "ovn is down")

	// nothing was committed, exactly as before: the pool does not exist in IPAM yet
	require.Nil(t, ctrl.ipam.Subnets[subnet.Name].IPPools[ippool.Name])
	require.Empty(t, ctrl.ipam.Subnets[subnet.Name].IPPools[""].V4Using.String())

	// and the retry commits the listing
	fakeController.mockOvnClient.EXPECT().CreateAddressSet(asName, gomock.Any()).Return(nil)
	fakeController.mockOvnClient.EXPECT().AddressSetUpdateAddress(asName, "10.16.0.2").Return(nil)
	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))
	ip, _, _, err := ctrl.ipam.GetRandomAddress("pod1.default", "nic1", nil, subnet.Name, ippool.Name, nil, true)
	require.NoError(t, err)
	require.Equal(t, "10.16.0.2", ip)
}

// TestHandleAddOrUpdateIPPoolAddressSetFailureReportsNoAvailable covers what the status
// advertises while a membership waits for its data plane: the pod controller picks an
// ippool of a namespace by status.v4AvailableIPs alone, so a reserved pool must report
// none, otherwise pods of that namespace pick it and fail instead of falling back.
func TestHandleAddOrUpdateIPPoolAddressSetFailureReportsNoAvailable(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet17"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv4,
			CIDRBlock:  "10.17.0.0/24",
			Gateway:    "10.17.0.1",
			ExcludeIps: []string{"10.17.0.1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool17"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2, EnableAddressSet: true},
		// what a previous pass published while the membership was still active
		Status: kubeovnv1.IPPoolStatus{
			V4AvailableIPs:     internal.BigInt{Int: *big.NewInt(2)},
			V4AvailableIPRange: "10.17.0.2-10.17.0.3",
		},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))
	ippools := ctrl.config.KubeOvnClient.KubeovnV1().IPPools()

	asName := util.IPPoolAddressSetName(ippool.Name)
	fakeController.mockOvnClient.EXPECT().CreateAddressSet(asName, gomock.Any()).Return(errors.New("ovn is down"))
	require.ErrorContains(t, ctrl.handleAddOrUpdateIPPool(ippool.Name), "ovn is down")

	got, err := ippools.Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, got.Status.IsReady())
	require.Zero(t, got.Status.V4AvailableIPs.Int64(), "a pool that hands out nothing must advertise nothing")
	require.Empty(t, got.Status.V4AvailableIPRange)

	// the retry succeeds and the pool advertises its addresses again
	fakeController.mockOvnClient.EXPECT().CreateAddressSet(asName, gomock.Any()).Return(nil)
	fakeController.mockOvnClient.EXPECT().AddressSetUpdateAddress(asName, "10.17.0.2", "10.17.0.3").Return(nil)
	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))
	got, err = ippools.Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.True(t, got.Status.IsReady())
	require.Equal(t, int64(2), got.Status.V4AvailableIPs.Int64())
}

// TestHandleAddOrUpdateIPPoolSpecIPsStatusUnchanged pins the scope of the counter refresh:
// a pool listing spec.ips is never reserved, so what it reports on a failed reconcile must
// stay what it was before spec.ipCount existed - the numbers of the last pass.
func TestHandleAddOrUpdateIPPoolSpecIPsStatusUnchanged(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet18"},
		Spec:       kubeovnv1.SubnetSpec{Protocol: kubeovnv1.ProtocolIPv4, CIDRBlock: "10.18.0.0/24", Gateway: "10.18.0.1"},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool18"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPs: []string{"10.18.0.2"}, EnableAddressSet: true},
		// what the previous pass published
		Status: kubeovnv1.IPPoolStatus{
			V4AvailableIPs:     internal.BigInt{Int: *big.NewInt(1)},
			V4AvailableIPRange: "10.18.0.2",
		},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, nil))
	fakeController.mockOvnClient.EXPECT().CreateAddressSet(gomock.Any(), gomock.Any()).Return(errors.New("ovn is down"))

	require.ErrorContains(t, ctrl.handleAddOrUpdateIPPool(ippool.Name), "ovn is down")

	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, int64(1), got.Status.V4AvailableIPs.Int64(), "the counters of a spec.ips pool are left as they were")
	require.Equal(t, "10.18.0.2", got.Status.V4AvailableIPRange)
	require.False(t, got.Status.IsReady())
}

// TestHandleAddOrUpdateIPPoolActivationFailureReason pins the condition reason a failed
// IPAM commit reports: it predates spec.ipCount and whatever watches ippool events keys
// on it, so the new reservation step must not rename it.
func TestHandleAddOrUpdateIPPoolActivationFailureReason(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet19"},
		Spec:       kubeovnv1.SubnetSpec{Protocol: kubeovnv1.ProtocolIPv4, CIDRBlock: "10.19.0.0/24", Gateway: "10.19.0.1"},
	}
	held := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool19a"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPs: []string{"10.19.0.2"}},
	}
	// claims the address the other pool already holds, so the activation fails
	conflicting := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool19b"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPs: []string{"10.19.0.2"}},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{held, conflicting},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	fakeController.mockOvnClient.EXPECT().DeleteAddressSet(gomock.Any()).Return(nil).AnyTimes()
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, nil))

	require.NoError(t, ctrl.handleAddOrUpdateIPPool(held.Name))
	require.ErrorContains(t, ctrl.handleAddOrUpdateIPPool(conflicting.Name), "conflict IPs")

	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), conflicting.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, got.Status.IsReady())
	require.Equal(t, "UpdateIPAMFailed", got.Status.GetCondition(kubeovnv1.Ready).Reason)
}

// TestReserveIPPoolIPsIPv6Only covers spec.ipCount on a single stack IPv6 subnet: the
// per-family counts must ask for IPv6 addresses only, as asking for IPv4 ones the subnet
// cannot have would report it as exhausted.
func TestReserveIPPoolIPsIPv6Only(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet20"},
		Spec: kubeovnv1.SubnetSpec{
			Protocol:   kubeovnv1.ProtocolIPv6,
			CIDRBlock:  "fd20::/120",
			Gateway:    "fd20::1",
			ExcludeIps: []string{"fd20::1"},
		},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool20"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPCount: 2},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	ctrl.ippoolKeyMutex = keymutex.NewHashed(97)
	fakeController.mockOvnClient.EXPECT().DeleteAddressSet(gomock.Any()).Return(nil).AnyTimes()
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, subnet.Spec.ExcludeIps))

	ips, err := ctrl.reserveIPPoolIPs(ippool.DeepCopy())
	require.NoError(t, err)
	require.Equal(t, []string{"fd20::2", "fd20::3"}, ips)

	require.NoError(t, ctrl.handleAddOrUpdateIPPool(ippool.Name))
	got, err := ctrl.config.KubeOvnClient.KubeovnV1().IPPools().Get(context.Background(), ippool.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []string{"fd20::2", "fd20::3"}, got.Status.IPs)
	require.True(t, got.Status.IsReady())
	require.Equal(t, int64(2), got.Status.V6AvailableIPs.Int64())
	require.Zero(t, got.Status.V4AvailableIPs.Int64())

	_, ip, _, err := ctrl.ipam.GetRandomAddress("pod1.default", "nic1", nil, subnet.Name, ippool.Name, nil, true)
	require.NoError(t, err)
	require.Contains(t, got.Status.IPs, ip)
}

// TestInitIPPoolIPAMSpecIPs covers the other half of the startup restore: a pool listing
// spec.ips is restored active, as it was before spec.ipCount existed - it has no
// membership that could be waiting for a data plane.
func TestInitIPPoolIPAMSpecIPs(t *testing.T) {
	subnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "subnet21"},
		Spec:       kubeovnv1.SubnetSpec{Protocol: kubeovnv1.ProtocolIPv4, CIDRBlock: "10.21.0.0/24", Gateway: "10.21.0.1"},
	}
	ippool := &kubeovnv1.IPPool{
		ObjectMeta: metav1.ObjectMeta{Name: "pool21"},
		Spec:       kubeovnv1.IPPoolSpec{Subnet: subnet.Name, IPs: []string{"10.21.0.2"}},
	}
	fakeController, err := newFakeControllerWithOptions(t, &FakeControllerOptions{
		Subnets: []*kubeovnv1.Subnet{subnet},
		IPPools: []*kubeovnv1.IPPool{ippool},
	})
	require.NoError(t, err)
	ctrl := fakeController.fakeController
	require.NoError(t, ctrl.ipam.AddOrUpdateSubnet(subnet.Name, subnet.Spec.CIDRBlock, subnet.Spec.Gateway, nil))

	require.NoError(t, ctrl.initIPPoolIPAM(ippool))
	ip, _, _, err := ctrl.ipam.GetRandomAddress("pod1.default", "nic1", nil, subnet.Name, ippool.Name, nil, true)
	require.NoError(t, err, "a restored spec.ips pool hands out its addresses right away")
	require.Equal(t, "10.21.0.2", ip)
}
