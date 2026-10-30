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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestNatGwTargetUsesExternalNetworkAndPodMAC(t *testing.T) {
	subnet := &kubeovnv1.Subnet{ObjectMeta: metav1.ObjectMeta{Name: "external", Annotations: map[string]string{util.NadMacvlanMasterAnnotation: "bond2"}}, Spec: kubeovnv1.SubnetSpec{Provider: "wan.kube-system"}}
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "vpc-nat-gw-test-0", Namespace: "kube-system", Annotations: map[string]string{nadv1.NetworkStatusAnnot: `[{"name":"kube-system/lan","mac":"02:00:00:00:00:01"},{"name":"kube-system/wan","mac":"02:00:00:00:00:02"}]`}}}
	name, err := util.GenMacvlanIfaceName("bond2")
	require.NoError(t, err)
	lookup := func(actual string) (netlink.Link, error) {
		require.Equal(t, name, actual)
		return &netlink.Macvlan{LinkAttrs: netlink.LinkAttrs{Index: 12, ParentIndex: 8, HardwareAddr: net.HardwareAddr{2, 0, 0, 0, 0, 3}}, Mode: netlink.MACVLAN_MODE_BRIDGE}, nil
	}
	info, err := natGwTargetInfo(pod, subnet, "", lookup)
	require.NoError(t, err)
	require.Equal(t, 12, info.TargetIfindex)
	require.Equal(t, net.HardwareAddr{2, 0, 0, 0, 0, 2}, info.TargetMAC)
	for _, status := range []string{"", `[]`, `[{"name":"kube-system/wan","mac":""}]`, `[{"name":"kube-system/wan","mac":"00:00:00:00:00:00"}]`, `[{"name":"kube-system/lan","mac":"02:00:00:00:00:01"}]`} {
		pod.Annotations[nadv1.NetworkStatusAnnot] = status
		_, err := natGwTargetInfo(pod, subnet, "", lookup)
		require.Error(t, err, "status %q", status)
	}
	_, err = natGwTargetInfo(pod, subnet, "", func(string) (netlink.Link, error) { return nil, errors.New("link missing") })
	require.Error(t, err)
}
