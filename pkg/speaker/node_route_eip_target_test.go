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

	t.Run("external subnet not found", func(t *testing.T) {
		c := &Controller{
			config:        &Configuration{TcBpfInterface: "bond0"},
			subnetsLister: &fakeSubnetLister{subnets: map[string]*kubeovnv1.Subnet{}},
		}
		eip := &kubeovnv1.IptablesEIP{
			ObjectMeta: metav1.ObjectMeta{Name: "eip-missing-subnet"},
			Spec:       kubeovnv1.IptablesEIPSpec{ExternalSubnet: "non-exist-subnet"},
		}
		pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "nat-gw-pod", Namespace: "kube-system"}}
		_, err := c.getNatGwTCRedirectTarget(pod, eip)
		require.ErrorContains(t, err, "failed to get external subnet non-exist-subnet")
	})
}

// TestNatGwTCRedirectTarget 验证 natGwTCRedirectTarget 在各类网络模式与边界输入下的直通目标解析
func TestNatGwTCRedirectTarget(t *testing.T) {
	t.Run("nil pod", func(t *testing.T) {
		subnet := &kubeovnv1.Subnet{ObjectMeta: metav1.ObjectMeta{Name: "ext-subnet"}}
		_, err := natGwTCRedirectTarget(nil, subnet, "bond0", nil)
		require.ErrorContains(t, err, "NAT gateway pod is missing")
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

	t.Run("veth ptp mode missing status", func(t *testing.T) {
		// 非 macvlan 模式下未包含外网网卡状态时报错
		subnet := &kubeovnv1.Subnet{
			ObjectMeta: metav1.ObjectMeta{Name: "ext-ptp"},
			Spec:       kubeovnv1.SubnetSpec{Provider: "wan.kube-system"},
		}
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vpc-nat-gw-ptp-0",
				Namespace: "kube-system",
				Annotations: map[string]string{
					nadv1.NetworkStatusAnnot: `[{"name":"kube-system/lan","mac":"02:00:00:00:00:01"}]`,
				},
			},
		}
		_, err := natGwTCRedirectTarget(pod, subnet, "bond0", nil)
		require.ErrorContains(t, err, "has no network status for kube-system/wan")
	})
}
