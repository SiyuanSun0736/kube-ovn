//go:build linux

package speaker

import (
	"errors"
	"net"
	"testing"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

// fakeSubnetLister 模拟测试用的 SubnetLister
type fakeSubnetLister struct {
	subnets map[string]*kubeovnv1.Subnet
}

func (f *fakeSubnetLister) List(_ labels.Selector) ([]*kubeovnv1.Subnet, error) {
	var list []*kubeovnv1.Subnet
	for _, s := range f.subnets {
		list = append(list, s)
	}
	return list, nil
}

func (f *fakeSubnetLister) Get(name string) (*kubeovnv1.Subnet, error) {
	if s, ok := f.subnets[name]; ok {
		return s, nil
	}
	return nil, apierrors.NewNotFound(kubeovnv1.Resource("subnet"), name)
}

// fakeNatGatewayLister 模拟测试用的 VpcNatGatewayLister
type fakeNatGatewayLister struct {
	gateways map[string]*kubeovnv1.VpcNatGateway
}

func (f *fakeNatGatewayLister) List(_ labels.Selector) ([]*kubeovnv1.VpcNatGateway, error) {
	var list []*kubeovnv1.VpcNatGateway
	for _, gw := range f.gateways {
		list = append(list, gw)
	}
	return list, nil
}

func (f *fakeNatGatewayLister) Get(name string) (*kubeovnv1.VpcNatGateway, error) {
	if gw, ok := f.gateways[name]; ok {
		return gw, nil
	}
	return nil, apierrors.NewNotFound(kubeovnv1.Resource("vpcnatgateway"), name)
}

// TestGetNatGwTCRedirectTarget 验证 Controller.getNatGwTCRedirectTarget 对 EIP ExternalSubnet 的校验与定位
func TestGetNatGwTCRedirectTarget(t *testing.T) {
	t.Run("empty external subnet", func(t *testing.T) {
		c := &Controller{
			config:        &Configuration{TcBpfInterface: "bond0"},
			subnetsLister: &fakeSubnetLister{subnets: map[string]*kubeovnv1.Subnet{}},
		}
		eip := &kubeovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{Name: "eip-empty-subnet"},
			Spec:       kubeovnv1.IptablesEIPSpec{ExternalSubnet: ""},
		}
		pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "nat-gw-pod", Namespace: "kube-system"}}
		_, err := c.getNatGwTCRedirectTarget(pod, eip)
		require.ErrorContains(t, err, "has empty externalSubnet")
	})

	t.Run("external subnet not found should return error without guessing", func(t *testing.T) {
		c := &Controller{
			config: &Configuration{TcBpfInterface: "bond0"},
			subnetsLister: &fakeSubnetLister{subnets: map[string]*kubeovnv1.Subnet{
				"gw-ext-subnet": {
					ObjectMeta: metav1.ObjectMeta{Name: "gw-ext-subnet"},
				},
			}},
		}
		eip := &kubeovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{Name: "eip-with-diff-subnet"},
			Spec: kubeovnv1.IptablesEIPSpec{
				NatGwDp:        "gw-1",
				ExternalSubnet: "las-derived-subnet-not-found",
			},
		}
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vpc-nat-gw-gw-1-0",
				Namespace: "kube-system",
			},
		}

		lookup := func(actual string) (netlink.Link, error) {
			return nil, errors.New("should not be called")
		}

		// 验证不搞越界猜测，明确 fail-fast 返回错误
		_, err := c.getNatGwTCRedirectTargetInternal(pod, eip, lookup)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get external subnet las-derived-subnet-not-found")
	})
}

// TestResolveEipRedirectTargetCache 验证 UID 缓存机制：只要 Pod UID 相同直接复用缓存，Pod 重建或失效则失效更新
func TestResolveEipRedirectTargetCache(t *testing.T) {
	gwSubnet := &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gw-ext",
		},
		Spec: kubeovnv1.SubnetSpec{Provider: "wan.kube-system"},
	}
	podIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	podLister := listerv1.NewPodLister(podIndexer)

	c := &Controller{
		config: &Configuration{
			TcBpfInterface:    "bond0",
			VpcNatGwNamespace: "kube-system",
		},
		gwPodsLister:  podLister,
		subnetsLister: &fakeSubnetLister{subnets: map[string]*kubeovnv1.Subnet{"gw-ext": gwSubnet}},
		gwTargetCache: make(map[string]gwTargetCacheEntry),
	}

	podName := util.GenNatGwPodName("gw-test")
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: "kube-system",
			UID:       types.UID("pod-uid-1111"),
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	err := podIndexer.Add(pod)
	require.NoError(t, err)

	eip := &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{Name: "eip-test"},
		Spec: kubeovnv1.IptablesEIPSpec{
			NatGwDp:        "gw-test",
			ExternalSubnet: "gw-ext",
		},
	}

	// 1. 预先在缓存中填入 uid-1111 的已解析 Target
	expectedTarget := TCRedirectTarget{
		TargetIfindex: 10,
		ParentIfindex: 2,
		TargetMAC:     net.HardwareAddr{0x02, 0x00, 0x00, 0x00, 0x00, 0x01},
	}
	c.gwTargetCache["gw-test"] = gwTargetCacheEntry{
		podUID: types.UID("pod-uid-1111"),
		target: expectedTarget,
	}

	// 2. 调用 resolveEipRedirectTarget，应直接命中 UID 缓存并返回，完全不调用 getNatGwTCRedirectTarget
	target, retainStale, err := c.resolveEipRedirectTarget(eip)
	require.NoError(t, err)
	require.False(t, retainStale)
	require.NotNil(t, target)
	require.Equal(t, expectedTarget, *target)

	// 3. Pod 被删除重建，UID 变为 pod-uid-2222
	podRebuilt := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: "kube-system",
			UID:       types.UID("pod-uid-2222"),
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	err = podIndexer.Update(podRebuilt)
	require.NoError(t, err)

	// 调用时因 UID 不匹配，缓存失效，尝试重新解析（因为未 mock netlink 会报错，证明没有走旧缓存）
	_, _, err = c.resolveEipRedirectTarget(eip)
	require.Error(t, err)

	// 4. Pod 处于 Terminating 状态（带有 DeletionTimestamp），清理缓存并返回 nil
	now := metav1.Now()
	podDeleting := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              podName,
			Namespace:         "kube-system",
			UID:               types.UID("pod-uid-1111"),
			DeletionTimestamp: &now,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	c.gwTargetCache["gw-test"] = gwTargetCacheEntry{
		podUID: types.UID("pod-uid-1111"),
		target: expectedTarget,
	}
	err = podIndexer.Update(podDeleting)
	require.NoError(t, err)

	target, retainStale, err = c.resolveEipRedirectTarget(eip)
	require.NoError(t, err)
	require.False(t, retainStale)
	require.Nil(t, target)
	// 验证缓存已被清理
	_, exists := c.gwTargetCache["gw-test"]
	require.False(t, exists)
}

// TestNatGwTCRedirectTarget 验证 natGwTCRedirectTarget 在各类网络模式与边界输入下的直通目标解析
func TestNatGwTCRedirectTarget(t *testing.T) {
	t.Run("nil pod", func(t *testing.T) {
		subnet := &kubeovnv1.Subnet{ObjectMeta: metav1.ObjectMeta{Name: "ext-subnet"}}
		_, err := natGwTCRedirectTarget(nil, subnet, "", nil)
		require.ErrorContains(t, err, "NAT gateway pod is missing")
	})

	t.Run("invalid provider format", func(t *testing.T) {
		subnet := &kubeovnv1.Subnet{
			ObjectMeta: metav1.ObjectMeta{Name: "ext-subnet"},
			Spec:       kubeovnv1.SubnetSpec{Provider: "invalid-provider-without-dot"},
		}
		pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "nat-gw", Namespace: "kube-system"}}
		_, err := natGwTCRedirectTarget(pod, subnet, "", nil)
		require.ErrorContains(t, err, "is not a valid NAD format")
	})

	t.Run("macvlan mode target resolution", func(t *testing.T) {
		subnet := &kubeovnv1.Subnet{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "external",
				Annotations: map[string]string{util.NadMacvlanMasterAnnotation: "bond2"},
			},
			Spec: kubeovnv1.SubnetSpec{Provider: "wan.kube-system"},
		}
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vpc-nat-gw-test-0",
				Namespace: "kube-system",
				Annotations: map[string]string{
					nadv1.NetworkStatusAnnot: `[{"name":"kube-system/lan","mac":"02:00:00:00:00:01"},{"name":"kube-system/wan","mac":"02:00:00:00:00:02"}]`,
				},
			},
		}
		name, err := util.GenMacvlanIfaceName("bond2")
		require.NoError(t, err)

		lookup := func(actual string) (netlink.Link, error) {
			require.Equal(t, name, actual)
			return &netlink.Macvlan{
				LinkAttrs: netlink.LinkAttrs{Index: 12, ParentIndex: 8, HardwareAddr: net.HardwareAddr{2, 0, 0, 0, 0, 3}},
				Mode:      netlink.MACVLAN_MODE_BRIDGE,
			}, nil
		}

		// 1. 成功解析直通网卡索引与 MAC
		target, err := natGwTCRedirectTarget(pod, subnet, "", lookup)
		require.NoError(t, err)
		require.Equal(t, 12, target.TargetIfindex)
		require.Equal(t, net.HardwareAddr{2, 0, 0, 0, 0, 2}, target.TargetMAC)

		// 2. 异常 network-status 处理（包括空串、空数组、空 MAC、全零 MAC 等）
		invalidStatuses := []string{
			"",
			`[]`,
			`[{"name":"kube-system/wan","mac":""}]`,
			`[{"name":"kube-system/wan","mac":"00:00:00:00:00:00"}]`,
			`[{"name":"kube-system/lan","mac":"02:00:00:00:00:01"}]`,
		}
		for _, status := range invalidStatuses {
			pod.Annotations[nadv1.NetworkStatusAnnot] = status
			_, err := natGwTCRedirectTarget(pod, subnet, "", lookup)
			require.Error(t, err, "status %q", status)
		}

		// 3. 设备不存在时报错
		pod.Annotations[nadv1.NetworkStatusAnnot] = `[{"name":"kube-system/wan","mac":"02:00:00:00:00:02"}]`
		_, err = natGwTCRedirectTarget(pod, subnet, "", func(string) (netlink.Link, error) {
			return nil, errors.New("link missing")
		})
		require.Error(t, err)
	})

	t.Run("veth bridge fdb mode target resolution", func(t *testing.T) {
		subnet := &kubeovnv1.Subnet{
			ObjectMeta: metav1.ObjectMeta{Name: "ext-veth"},
			Spec:       kubeovnv1.SubnetSpec{Provider: "wan.kube-system"},
		}
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vpc-nat-gw-veth-0",
				Namespace: "kube-system",
				Annotations: map[string]string{
					nadv1.NetworkStatusAnnot: `[{"name":"kube-system/wan","mac":"02:00:00:00:00:02"}]`,
				},
			},
		}

		fdbLookup := func(mac net.HardwareAddr) (int, error) {
			if mac.String() == "02:00:00:00:00:02" {
				return 25, nil
			}
			return 0, errors.New("mac not found in fdb")
		}

		// 1. 成功解析 veth 目标设备索引与 MAC（ParentIfindex 为 0）
		target, err := natGwTCRedirectTarget(pod, subnet, "bond0", func(string) (netlink.Link, error) {
			return &netlink.Device{LinkAttrs: netlink.LinkAttrs{Index: 10}}, nil
		}, fdbLookup)
		require.NoError(t, err)
		require.Equal(t, 25, target.TargetIfindex)
		require.Equal(t, 0, target.ParentIfindex)
		require.Equal(t, net.HardwareAddr{2, 0, 0, 0, 0, 2}, target.TargetMAC)

		// 2. FDB 中未找到 MAC 记录时报错
		target, err = natGwTCRedirectTarget(pod, subnet, "bond0", nil, func(net.HardwareAddr) (int, error) {
			return 0, errors.New("fdb entry not found")
		})
		require.ErrorContains(t, err, "failed to resolve host veth device for pod")

		// 3. 校验重定向目标不能等于 Ingress 父设备（防环回）
		_, err = natGwTCRedirectTarget(pod, subnet, "bond0", func(string) (netlink.Link, error) {
			return &netlink.Device{LinkAttrs: netlink.LinkAttrs{Index: 25}}, nil
		}, fdbLookup)
		require.ErrorContains(t, err, "bridge FDB target device points to ingress interface")

		// 4. network-status 缺少对应网络时报错
		pod.Annotations[nadv1.NetworkStatusAnnot] = `[{"name":"kube-system/lan","mac":"02:00:00:00:00:01"}]`
		_, err = natGwTCRedirectTarget(pod, subnet, "bond0", nil, fdbLookup)
		require.ErrorContains(t, err, "has no network status for kube-system/wan")
	})

	t.Run("fdbLookupInternal device name filtering", func(t *testing.T) {
		mac := net.HardwareAddr{0x02, 0x00, 0x00, 0x00, 0x00, 0x02}
		neighs := []netlink.Neigh{
			{LinkIndex: 1, HardwareAddr: mac}, // bond0 (物理口/bond口，应被跳过)
			{LinkIndex: 2, HardwareAddr: mac}, // 119ec988_net1_h (Pod veth 口，应被选中)
		}

		links := map[int]string{
			1: "bond0",
			2: "119ec988_net1_h",
		}

		neighListMock := func(int, int) ([]netlink.Neigh, error) {
			return neighs, nil
		}
		linkByIndexMock := func(idx int) (netlink.Link, error) {
			name, ok := links[idx]
			if !ok {
				return nil, errors.New("link not found")
			}
			return &netlink.Device{LinkAttrs: netlink.LinkAttrs{Index: idx, Name: name}}, nil
		}

		// 成功过滤跳过 bond0，精准命中 119ec988_net1_h
		idx, err := fdbLookupInternal(mac, neighListMock, linkByIndexMock)
		require.NoError(t, err)
		require.Equal(t, 2, idx)

		// 测试全部为非 Pod 网卡时报错
		links[2] = "ens7f0"
		_, err = fdbLookupInternal(mac, neighListMock, linkByIndexMock)
		require.ErrorContains(t, err, "no bridge FDB entry found on pod veth interfaces")

		// 测试 isPodVethName
		require.True(t, isPodVethName("119ec988_net1_h"))
		require.True(t, isPodVethName("a1a3a0b2527d_h"))
		require.True(t, isPodVethName("vethab07e886"))
		require.False(t, isPodVethName("bond0"))
		require.False(t, isPodVethName("ens7f0"))
		require.False(t, isPodVethName("docker0"))
	})
}
