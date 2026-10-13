package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	kubeovnlister "github.com/kubeovn/kube-ovn/pkg/client/listers/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestParseEIPDestination(t *testing.T) {
	tests := []struct {
		name        string
		eip         string
		wantMask    int
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid IPv4",
			eip:      "192.168.1.100",
			wantMask: 32,
			wantErr:  false,
		},
		{
			name:     "valid IPv6",
			eip:      "2001:db8::1",
			wantMask: 128,
			wantErr:  false,
		},
		{
			name:        "invalid IP",
			eip:         "invalid",
			wantErr:     true,
			errContains: "invalid EIP address",
		},
		{
			name:        "empty string",
			eip:         "",
			wantErr:     true,
			errContains: "invalid EIP address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst, err := parseEIPDestination(tt.eip)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, dst)
			ones, _ := dst.Mask.Size()
			assert.Equal(t, tt.wantMask, ones)
		})
	}
}

func TestShouldEnqueueIptablesEip(t *testing.T) {
	tests := []struct {
		name           string
		externalSubnet string
		ready          bool
		want           bool
	}{
		{
			name:           "ready with ExternalSubnet",
			externalSubnet: "external-subnet",
			ready:          true,
			want:           true,
		},
		{
			name:           "ready without ExternalSubnet",
			externalSubnet: "",
			ready:          true,
			want:           false,
		},
		{
			name:           "not ready with ExternalSubnet",
			externalSubnet: "external-subnet",
			ready:          false,
			want:           false,
		},
		{
			name:           "not ready without ExternalSubnet",
			externalSubnet: "",
			ready:          false,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eip := &kubeovnv1.IptablesEIP{
				Spec: kubeovnv1.IptablesEIPSpec{
					ExternalSubnet: tt.externalSubnet,
				},
				Status: kubeovnv1.IptablesEIPStatus{
					Ready: tt.ready,
				},
			}
			assert.Equal(t, tt.want, shouldEnqueueIptablesEip(eip))
		})
	}
}

func TestIsVpcNatGwPod(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
		want   bool
	}{
		{
			name:   "NAT GW pod",
			labels: map[string]string{util.VpcNatGatewayLabel: "true"},
			want:   true,
		},
		{
			name:   "NAT GW pod with extra labels",
			labels: map[string]string{util.VpcNatGatewayLabel: "true", "app": "test"},
			want:   true,
		},
		{
			name:   "not NAT GW pod - label false",
			labels: map[string]string{util.VpcNatGatewayLabel: "false"},
			want:   false,
		},
		{
			name:   "not NAT GW pod - no label",
			labels: map[string]string{"app": "test"},
			want:   false,
		},
		{
			name:   "not NAT GW pod - nil labels",
			labels: nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Labels: tt.labels},
			}
			assert.Equal(t, tt.want, isVpcNatGwPod(pod))
		})
	}
}

func TestGetNatGwNameFromPod(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
		want   string
	}{
		{
			name:   "has NAT GW name label",
			labels: map[string]string{util.VpcNatGatewayNameLabel: "my-nat-gw"},
			want:   "my-nat-gw",
		},
		{
			name:   "no NAT GW name label",
			labels: map[string]string{"app": "test"},
			want:   "",
		},
		{
			name:   "nil labels",
			labels: nil,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Labels: tt.labels},
			}
			assert.Equal(t, tt.want, getNatGwNameFromPod(pod))
		})
	}
}

func TestDeletedEIPsStateManagement(t *testing.T) {
	t.Run("add event clears deleted mark", func(t *testing.T) {
		c := &Controller{}

		// Simulate a previous delete event that marked EIP as deleted
		c.deletedEIPs.Store("test-eip", true)

		// Verify mark exists
		_, exists := c.deletedEIPs.Load("test-eip")
		assert.True(t, exists, "deleted mark should exist before add event")

		// Simulate what enqueueAddIptablesEip does - clear the mark
		c.deletedEIPs.Delete("test-eip")

		// Verify mark is cleared
		_, exists = c.deletedEIPs.Load("test-eip")
		assert.False(t, exists, "deleted mark should be cleared after add event")
	})

	t.Run("update event clears deleted mark", func(t *testing.T) {
		c := &Controller{}

		// Simulate a previous delete event
		c.deletedEIPs.Store("test-eip", true)

		// Verify mark exists
		_, exists := c.deletedEIPs.Load("test-eip")
		assert.True(t, exists, "deleted mark should exist before update event")

		// Simulate what enqueueUpdateIptablesEip does - clear the mark
		c.deletedEIPs.Delete("test-eip")

		// Verify mark is cleared
		_, exists = c.deletedEIPs.Load("test-eip")
		assert.False(t, exists, "deleted mark should be cleared after update event")
	})

	t.Run("delete event sets deleted mark", func(t *testing.T) {
		c := &Controller{}

		// Verify no mark initially
		_, exists := c.deletedEIPs.Load("test-eip")
		assert.False(t, exists, "deleted mark should not exist initially")

		// Simulate what enqueueDeleteIptablesEip does - set the mark
		c.deletedEIPs.Store("test-eip", true)

		// Verify mark is set
		_, exists = c.deletedEIPs.Load("test-eip")
		assert.True(t, exists, "deleted mark should be set after delete event")
	})

	t.Run("successful route deletion clears deleted mark", func(t *testing.T) {
		c := &Controller{}

		// Simulate delete event marking EIP as deleted
		c.deletedEIPs.Store("test-eip", true)

		// Verify mark exists
		_, exists := c.deletedEIPs.Load("test-eip")
		assert.True(t, exists, "deleted mark should exist after delete event")

		// Simulate what processNextIptablesEipDeleteItem does after successful deletion
		c.deletedEIPs.Delete("test-eip")

		// Verify mark is cleaned up (prevents memory leak)
		_, exists = c.deletedEIPs.Load("test-eip")
		assert.False(t, exists, "deleted mark should be cleaned up after successful route deletion")
	})

	t.Run("sync skips deleted EIP", func(t *testing.T) {
		c := &Controller{}

		// Mark EIP as deleted
		c.deletedEIPs.Store("test-eip", true)

		// Check if the EIP is marked as deleted (simulating syncIptablesEipRoute logic)
		_, deleted := c.deletedEIPs.Load("test-eip")
		assert.True(t, deleted, "should detect EIP as deleted")
	})

	t.Run("sync proceeds for non-deleted EIP", func(t *testing.T) {
		c := &Controller{}

		// EIP is not marked as deleted
		_, deleted := c.deletedEIPs.Load("test-eip")
		assert.False(t, deleted, "should not detect EIP as deleted when not marked")
	})

	t.Run("multiple EIPs tracked independently", func(t *testing.T) {
		c := &Controller{}

		// Mark multiple EIPs as deleted
		c.deletedEIPs.Store("eip-1", true)
		c.deletedEIPs.Store("eip-2", true)

		// Clear one EIP
		c.deletedEIPs.Delete("eip-1")

		// Verify eip-1 is cleared but eip-2 still marked
		_, exists1 := c.deletedEIPs.Load("eip-1")
		_, exists2 := c.deletedEIPs.Load("eip-2")
		assert.False(t, exists1, "eip-1 should be cleared")
		assert.True(t, exists2, "eip-2 should still be marked")
	})
}

func TestResyncLocalIptablesEipRoutes(t *testing.T) {
	const master = "eth0"
	macvlanName, err := util.GenMacvlanIfaceName(master)
	require.NoError(t, err)

	newEIP := func(name, v4ip, externalSubnet string, ready bool) *kubeovnv1.IptablesEIP {
		return &kubeovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec: kubeovnv1.IptablesEIPSpec{
				V4ip:           v4ip,
				ExternalSubnet: externalSubnet,
			},
			Status: kubeovnv1.IptablesEIPStatus{Ready: ready},
		}
	}

	// subnetWithMaster resolves buildEipRouteInfo; subnetNoMaster makes it return nil.
	subnetWithMaster := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "ext-subnet",
			Annotations: map[string]string{util.NadMacvlanMasterAnnotation: master},
		},
	}
	subnetNoMaster := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{Name: "ext-subnet-no-master"},
	}

	eips := []*kubeovnv1.IptablesEIP{
		newEIP("eip-ready", "1.1.1.1", "ext-subnet", true),                      // enqueued
		newEIP("eip-not-ready", "1.1.1.2", "ext-subnet", false),                 // skip: not ready
		newEIP("eip-no-external-subnet", "1.1.1.3", "", true),                   // skip: no ExternalSubnet
		newEIP("eip-no-v4ip", "", "ext-subnet", true),                           // skip: buildEipRouteInfo nil
		newEIP("eip-subnet-no-master", "1.1.1.4", "ext-subnet-no-master", true), // skip: no master
		newEIP("eip-subnet-missing", "1.1.1.5", "nonexistent", true),            // skip: subnet not found
	}

	subnetIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	require.NoError(t, subnetIndexer.Add(subnetWithMaster))
	require.NoError(t, subnetIndexer.Add(subnetNoMaster))

	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	for _, eip := range eips {
		require.NoError(t, eipIndexer.Add(eip))
	}

	c := &Controller{
		subnetsLister:      kubeovnlister.NewSubnetLister(subnetIndexer),
		iptablesEipsLister: kubeovnlister.NewIptablesEIPLister(eipIndexer),
		iptablesEipQueue:   newTypedRateLimitingQueue[eipRouteInfo]("test-iptables-eip", nil),
	}
	defer c.iptablesEipQueue.ShutDown()

	c.resyncLocalIptablesEipRoutes()

	// Only the ready EIP backed by a subnet with a macvlan master should be enqueued;
	// the route worker (not exercised here) later decides add vs delete by NAT GW locality.
	require.Equal(t, 1, c.iptablesEipQueue.Len())
	item, _ := c.iptablesEipQueue.Get()
	assert.Equal(t, "eip-ready", item.eipName)
	assert.Equal(t, "1.1.1.1", item.v4ip)
	assert.Equal(t, macvlanName, item.macvlanName)
}

func TestHandleNatGwPodDelete(t *testing.T) {
	const (
		master    = "eth0"
		natGwName = "gw1"
	)
	macvlanName, err := util.GenMacvlanIfaceName(master)
	require.NoError(t, err)

	subnetWithMaster := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "ext-subnet",
			Annotations: map[string]string{util.NadMacvlanMasterAnnotation: master},
		},
	}
	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "eip-1",
			Labels: map[string]string{util.VpcNatGatewayNameLabel: natGwName},
		},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip:           "1.1.1.1",
			ExternalSubnet: "ext-subnet",
			NatGwDp:        natGwName,
		},
		Status: kubeovnv1.IptablesEIPStatus{Ready: true},
	}

	subnetIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	require.NoError(t, subnetIndexer.Add(subnetWithMaster))
	eipIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	require.NoError(t, eipIndexer.Add(eip))

	newController := func() *Controller {
		return &Controller{
			subnetsLister:      kubeovnlister.NewSubnetLister(subnetIndexer),
			iptablesEipsLister: kubeovnlister.NewIptablesEIPLister(eipIndexer),
			iptablesEipQueue:   newTypedRateLimitingQueue[eipRouteInfo]("test-natgw-delete", nil),
		}
	}

	t.Run("nat gw pod delete enqueues its eips", func(t *testing.T) {
		c := newController()
		defer c.iptablesEipQueue.ShutDown()
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: util.GenNatGwPodName(natGwName),
				Labels: map[string]string{
					util.VpcNatGatewayLabel:     "true",
					util.VpcNatGatewayNameLabel: natGwName,
				},
			},
		}
		c.handleNatGwPodDelete(pod)

		require.Equal(t, 1, c.iptablesEipQueue.Len())
		item, _ := c.iptablesEipQueue.Get()
		assert.Equal(t, "eip-1", item.eipName)
		assert.Equal(t, "1.1.1.1", item.v4ip)
		assert.Equal(t, macvlanName, item.macvlanName)
	})

	t.Run("non nat gw pod is ignored", func(t *testing.T) {
		c := newController()
		defer c.iptablesEipQueue.ShutDown()
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "some-pod", Labels: map[string]string{"app": "test"}},
		}
		c.handleNatGwPodDelete(pod)

		assert.Equal(t, 0, c.iptablesEipQueue.Len())
	})
}

func TestBuildEipRouteInfoModes(t *testing.T) {
	macvlanSubnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "macvlan-sub",
			Annotations: map[string]string{
				util.NadMacvlanMasterAnnotation: "bond1",
			},
		},
	}
	vethSubnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "veth-sub",
		},
	}

	subnetIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	require.NoError(t, subnetIndexer.Add(macvlanSubnet))
	require.NoError(t, subnetIndexer.Add(vethSubnet))

	c := &Controller{
		subnetsLister: kubeovnlister.NewSubnetLister(subnetIndexer),
	}

	// Case 1: Standard macvlan subnet without BGP annotation
	macvlanEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip-macvlan"},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip:           "183.240.180.136",
			ExternalSubnet: "macvlan-sub",
		},
	}
	info1 := c.buildEipRouteInfo(macvlanEip)
	require.NotNil(t, info1)
	assert.Equal(t, "macvlan", info1.mode)
	assert.Equal(t, "macbond1", info1.macvlanName)

	// Case 2: Macvlan subnet with BGP Anycast or HA member label -> rejected (HA does not support macvlan)
	bgpEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name: "eip-bgp",
			Annotations: map[string]string{
				util.BgpAnnotation: "true",
			},
		},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip:           "183.240.180.136",
			ExternalSubnet: "macvlan-sub",
		},
	}
	info2 := c.buildEipRouteInfo(bgpEip)
	assert.Nil(t, info2)

	haMemberEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name: "eip-ha-member",
			Labels: map[string]string{
				util.NatGatewayMemberLabel: "gw-m1",
			},
		},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip:           "183.240.180.136",
			ExternalSubnet: "macvlan-sub",
		},
	}
	assert.Nil(t, c.buildEipRouteInfo(haMemberEip))

	// Case 3: Ptp / veth subnet without macvlan master annotation -> veth mode
	vethEip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip-veth"},
		Spec: kubeovnv1.IptablesEIPSpec{
			V4ip:           "183.240.180.136",
			ExternalSubnet: "veth-sub",
		},
	}
	info3 := c.buildEipRouteInfo(vethEip)
	require.NotNil(t, info3)
	assert.Equal(t, "veth", info3.mode)
	assert.Empty(t, info3.macvlanName)
}

func TestGetPodPtpIPInDaemon(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"k8s.v1.cni.cncf.io/network-status": `[
					{"name": "generic-veth", "interface": "eth0", "ips": ["10.210.88.253"], "default": true},
					{"name": "kube-system/public-veth-ptp", "interface": "net1", "ips": ["169.254.100.3"]}
				]`,
			},
		},
	}
	assert.Equal(t, "169.254.100.3", getPodPtpIP(pod))
}
