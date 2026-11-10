package daemon

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	kubeovnlister "github.com/kubeovn/kube-ovn/pkg/client/listers/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func TestPodTunnelKeyReady(t *testing.T) {
	subnetIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	for _, subnet := range []*kubeovnv1.Subnet{
		{
			// Status.TunnelKey is the subnet's own logical switch key; pods must carry the VPC key.
			ObjectMeta: metav1.ObjectMeta{Name: "vpc-subnet"},
			Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "vpc-a"},
			Status:     kubeovnv1.SubnetStatus{TunnelKey: 99},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "pending-vpc-subnet"},
			Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "pending-vpc"},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "orphan-subnet"},
			Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vpc: "missing-vpc"},
			Status:     kubeovnv1.SubnetStatus{TunnelKey: 77},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "vlan-subnet"},
			Spec:       kubeovnv1.SubnetSpec{Provider: util.OvnProvider, Vlan: "vlan-a", Vpc: "missing-vpc"},
		},
	} {
		if err := subnetIndexer.Add(subnet); err != nil {
			t.Fatal(err)
		}
	}
	vpcIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	for _, vpc := range []*kubeovnv1.Vpc{
		{ObjectMeta: metav1.ObjectMeta{Name: "vpc-a"}, Status: kubeovnv1.VpcStatus{TunnelKey: 1234}},
		{ObjectMeta: metav1.ObjectMeta{Name: "pending-vpc"}},
	} {
		if err := vpcIndexer.Add(vpc); err != nil {
			t.Fatal(err)
		}
	}
	csh := cniServerHandler{Controller: &Controller{
		subnetsLister: kubeovnlister.NewSubnetLister(subnetIndexer),
		vpcsLister:    kubeovnlister.NewVpcLister(vpcIndexer),
	}}
	pod := func(provider, subnet, tunnelKey string) *corev1.Pod {
		annotations := map[string]string{}
		if subnet != "" {
			annotations[fmt.Sprintf(util.LogicalSwitchAnnotationTemplate, provider)] = subnet
		}
		if tunnelKey != "" {
			annotations[fmt.Sprintf(util.TunnelKeyAnnotationTemplate, provider)] = tunnelKey
		}
		return &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Annotations: annotations}}
	}

	tests := []struct {
		name     string
		pod      *corev1.Pod
		provider string
		want     bool
		wantErr  bool
	}{
		{name: "OVN provider without logical switch needs no VPC key", pod: pod(util.OvnProvider, "", ""), provider: util.OvnProvider, want: true},
		{name: "non-OVN provider without logical switch", pod: pod("net1.default", "", ""), provider: "net1.default", want: true},
		{name: "key matches VPC status", pod: pod(util.OvnProvider, "vpc-subnet", "1234"), provider: util.OvnProvider, want: true},
		{name: "custom OVN provider key matches", pod: pod("net1.default.ovn", "vpc-subnet", "1234"), provider: "net1.default.ovn", want: true},
		{name: "subnet key from a pre-upgrade controller is accepted", pod: pod(util.OvnProvider, "vpc-subnet", "99"), provider: util.OvnProvider, want: true},
		{name: "subnet key accepted without a resolvable VPC", pod: pod(util.OvnProvider, "orphan-subnet", "77"), provider: util.OvnProvider, want: true},
		{name: "VPC key missing", pod: pod(util.OvnProvider, "vpc-subnet", ""), provider: util.OvnProvider, want: false},
		{name: "VPC key zero", pod: pod(util.OvnProvider, "vpc-subnet", "0"), provider: util.OvnProvider, want: false},
		{name: "VPC key non-numeric", pod: pod(util.OvnProvider, "vpc-subnet", "invalid"), provider: util.OvnProvider, want: false},
		{name: "VPC key differs from VPC status", pod: pod(util.OvnProvider, "vpc-subnet", "999"), provider: util.OvnProvider, want: false},
		{name: "VPC status not ready", pod: pod(util.OvnProvider, "pending-vpc-subnet", "1234"), provider: util.OvnProvider, want: false},
		{name: "missing VPC", pod: pod(util.OvnProvider, "orphan-subnet", "1234"), provider: util.OvnProvider, want: false, wantErr: true},
		{name: "vlan subnet needs no key or VPC", pod: pod(util.OvnProvider, "vlan-subnet", ""), provider: util.OvnProvider, want: true},
		{name: "vlan stale key does not block CNI", pod: pod(util.OvnProvider, "vlan-subnet", "1234"), provider: util.OvnProvider, want: true},
		{name: "missing subnet", pod: pod(util.OvnProvider, "missing-subnet", "1234"), provider: util.OvnProvider, want: false, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := csh.podTunnelKeyReady(tt.pod, tt.provider)
			if (err != nil) != tt.wantErr {
				t.Fatalf("podTunnelKeyReady() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("podTunnelKeyReady() = %v, want %v", got, tt.want)
			}
		})
	}
}
