package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/utils/keymutex"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	kubeovnfake "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned/fake"
	kubeovninformerfactory "github.com/kubeovn/kube-ovn/pkg/client/informers/externalversions"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func Test_getVipIps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		svc      *v1.Service
		expected []string
	}{
		{
			name: "annotation with single IPv4",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						util.SwitchLBRuleVipsAnnotation: "10.0.0.1",
					},
				},
			},
			expected: []string{"10.0.0.1"},
		},
		{
			name: "annotation with dual-stack IPs",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						util.SwitchLBRuleVipsAnnotation: "10.0.0.1,fd00::1",
					},
				},
			},
			expected: []string{"10.0.0.1", "fd00::1"},
		},
		{
			name: "annotation with empty value should return no IPs",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						util.SwitchLBRuleVipsAnnotation: "",
					},
				},
			},
			expected: nil,
		},
		{
			name: "annotation with trailing comma should filter empty elements",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						util.SwitchLBRuleVipsAnnotation: "10.0.0.1,",
					},
				},
			},
			expected: []string{"10.0.0.1"},
		},
		{
			name: "annotation with leading comma should filter empty elements",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						util.SwitchLBRuleVipsAnnotation: ",10.0.0.1",
					},
				},
			},
			expected: []string{"10.0.0.1"},
		},
		{
			name: "no annotation falls back to ClusterIPs",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
				},
				Spec: v1.ServiceSpec{
					ClusterIP:  "10.96.0.1",
					ClusterIPs: []string{"10.96.0.1"},
				},
			},
			expected: []string{"10.96.0.1"},
		},
		{
			name: "no annotation with external IP from subnet",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						util.ServiceExternalIPFromSubnetAnnotation: "external-subnet",
					},
				},
				Spec: v1.ServiceSpec{
					ClusterIP:  "10.96.0.1",
					ClusterIPs: []string{"10.96.0.1"},
				},
				Status: v1.ServiceStatus{
					LoadBalancer: v1.LoadBalancerStatus{
						Ingress: []v1.LoadBalancerIngress{
							{IP: "192.168.1.1"},
						},
					},
				},
			},
			expected: []string{"10.96.0.1", "192.168.1.1"},
		},
		{
			name: "no annotation with empty ingress IP should be filtered",
			svc: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						util.ServiceExternalIPFromSubnetAnnotation: "external-subnet",
					},
				},
				Spec: v1.ServiceSpec{
					ClusterIP:  "10.96.0.1",
					ClusterIPs: []string{"10.96.0.1"},
				},
				Status: v1.ServiceStatus{
					LoadBalancer: v1.LoadBalancerStatus{
						Ingress: []v1.LoadBalancerIngress{
							{IP: "192.168.1.1"},
							{IP: ""},
							{Hostname: "lb.example.com"},
						},
					},
				},
			},
			expected: []string{"10.96.0.1", "192.168.1.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getVipIps(tt.svc)
			require.Equal(t, tt.expected, got)
		})
	}
}

// newBgpLbVipController creates a minimal controller wired for handleAddBgpLbVipService tests.
func newBgpLbVipController(t *testing.T, vip *kubeovnv1.Vip, svc *v1.Service) *Controller {
	t.Helper()

	fc, err := newFakeControllerWithOptions(t, nil)
	require.NoError(t, err)
	ctrl := fc.fakeController

	// Wire virtualIpsLister from a dedicated informer factory.
	kubeovnClient := kubeovnfake.NewSimpleClientset()
	if vip != nil {
		_, err = kubeovnClient.KubeovnV1().Vips().Create(context.Background(), vip, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	stopCh := make(chan struct{})
	t.Cleanup(func() { close(stopCh) })
	vipInformerFactory := kubeovninformerfactory.NewSharedInformerFactoryWithOptions(kubeovnClient, 0)
	vipInformer := vipInformerFactory.Kubeovn().V1().Vips()
	vipInformerFactory.Start(stopCh)
	vipInformerFactory.WaitForCacheSync(stopCh)

	ctrl.virtualIpsLister = vipInformer.Lister()
	ctrl.svcKeyMutex = keymutex.NewHashed(0)

	// Register bgpVip indexer — must mirror the production indexer in controller.go:
	// index by both BgpVipAnnotation and MetalLBAllowSharedIPAnnotation so that
	// cleanupBgpLbVipServiceBindingByVip can look up Services via either annotation.
	err = fc.fakeInformers.serviceInformer.Informer().AddIndexers(cache.Indexers{
		bgpVipIndexName: func(obj any) ([]string, error) {
			s, ok := obj.(*v1.Service)
			if !ok {
				return nil, nil
			}
			var keys []string
			if v := s.Annotations[util.MetalLBAllowSharedIPAnnotation]; v != "" {
				keys = append(keys, v)
			} else if v := s.Annotations[util.BgpVipAnnotation]; v != "" {
				keys = append(keys, v)
			}
			return keys, nil
		},
	})
	require.NoError(t, err)
	ctrl.svcByBgpVipIndexer = fc.fakeInformers.serviceInformer.Informer().GetIndexer()

	// Seed service into the services lister cache.
	if svc != nil {
		err = fc.fakeInformers.serviceInformer.Informer().GetIndexer().Add(svc)
		require.NoError(t, err)
	}

	// Also register VIP in the indexer so the lister can find it immediately.
	if vip != nil {
		err = vipInformer.Informer().GetIndexer().Add(vip)
		require.NoError(t, err)
	}

	ctrl.config.KubeClient = fc.fakeController.config.KubeClient
	if svc != nil {
		_, err = ctrl.config.KubeClient.CoreV1().Services(svc.Namespace).Create(
			context.Background(), svc, metav1.CreateOptions{})
		if err != nil {
			// already exists from fake construction — ignore
			_ = err
		}
	}

	return ctrl
}

func TestHandleAddBgpLbVipService(t *testing.T) {
	t.Parallel()

	const (
		ns      = metav1.NamespaceDefault
		svcName = "my-lb-svc"
		vipName = "my-lb-vip"
		vipIP   = "203.0.113.5"
	)

	readyVIP := func() *kubeovnv1.Vip {
		return &kubeovnv1.Vip{
			ObjectMeta: metav1.ObjectMeta{Name: vipName},
			Spec: kubeovnv1.VipSpec{
				Type:   util.BgpLbVip,
				Subnet: "external",
			},
			Status: kubeovnv1.VipStatus{V4ip: vipIP},
		}
	}

	lbSvc := func(annotation string) *v1.Service {
		svc := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      svcName,
				Namespace: ns,
			},
			Spec: v1.ServiceSpec{Type: v1.ServiceTypeLoadBalancer},
		}
		if annotation != "" {
			svc.Annotations = map[string]string{util.BgpVipAnnotation: annotation}
		}
		return svc
	}

	key := cache.MetaObjectToName(lbSvc(vipName)).String()

	t.Run("no annotation: noop", func(t *testing.T) {
		t.Parallel()
		ctrl := newBgpLbVipController(t, readyVIP(), lbSvc(""))
		require.NoError(t, ctrl.handleAddBgpLbVipService(key))
	})

	t.Run("non-LB service type: noop", func(t *testing.T) {
		t.Parallel()
		svc := lbSvc(vipName)
		svc.Spec.Type = v1.ServiceTypeClusterIP
		ctrl := newBgpLbVipController(t, readyVIP(), svc)
		require.NoError(t, ctrl.handleAddBgpLbVipService(key))
	})

	t.Run("vip not found: error returned", func(t *testing.T) {
		t.Parallel()
		ctrl := newBgpLbVipController(t, nil, lbSvc("nonexistent-vip"))
		err := ctrl.handleAddBgpLbVipService(key)
		require.Error(t, err)
	})

	t.Run("vip wrong type: error returned", func(t *testing.T) {
		t.Parallel()
		vip := readyVIP()
		vip.Spec.Type = util.SwitchLBRuleVip // wrong type
		ctrl := newBgpLbVipController(t, vip, lbSvc(vipName))
		err := ctrl.handleAddBgpLbVipService(key)
		require.Error(t, err)
	})

	t.Run("vip has no IP yet: skip silently", func(t *testing.T) {
		t.Parallel()
		vip := readyVIP()
		vip.Status.V4ip = ""
		ctrl := newBgpLbVipController(t, vip, lbSvc(vipName))
		err := ctrl.handleAddBgpLbVipService(key)
		require.NoError(t, err)
	})

	t.Run("happy path: ingress set, bgp=true no longer written by controller", func(t *testing.T) {
		t.Parallel()
		ctrl := newBgpLbVipController(t, readyVIP(), lbSvc(vipName))
		require.NoError(t, ctrl.handleAddBgpLbVipService(key))

		updated, err := ctrl.config.KubeClient.CoreV1().Services(ns).Get(
			context.Background(), svcName, metav1.GetOptions{})
		require.NoError(t, err)
		require.Equal(t, []v1.LoadBalancerIngress{{IP: vipIP}}, updated.Status.LoadBalancer.Ingress)
		// The speaker now gates on bgp-vip / allow-shared-ip directly;
		// the controller no longer writes ovn.kubernetes.io/bgp=true.
		require.Empty(t, updated.Annotations[util.BgpAnnotation])
	})

	t.Run("idempotent: second call is a noop", func(t *testing.T) {
		t.Parallel()
		ctrl := newBgpLbVipController(t, readyVIP(), lbSvc(vipName))
		require.NoError(t, ctrl.handleAddBgpLbVipService(key))
		require.NoError(t, ctrl.handleAddBgpLbVipService(key))
	})
}

func TestReconcileBgpLbVipServiceLocked(t *testing.T) {
	t.Parallel()

	const (
		ns      = metav1.NamespaceDefault
		svcName = "my-lb-svc"
		vipName = "my-lb-vip"
		vipIP   = "203.0.113.5"
	)

	vip := &kubeovnv1.Vip{
		ObjectMeta: metav1.ObjectMeta{Name: vipName},
		Spec: kubeovnv1.VipSpec{
			Type:   util.BgpLbVip,
			Subnet: "external",
		},
		Status: kubeovnv1.VipStatus{V4ip: vipIP},
	}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Annotations: map[string]string{
				util.BgpVipAnnotation: vipName,
			},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
		},
	}

	ctrl := newBgpLbVipController(t, vip, svc)
	key := cache.MetaObjectToName(svc).String()

	require.NoError(t, ctrl.reconcileBgpLbVipServiceLocked(key, svc))

	updated, err := ctrl.config.KubeClient.CoreV1().Services(ns).Get(
		context.Background(), svcName, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, updated.Spec.ExternalIPs)
	require.Equal(t, []v1.LoadBalancerIngress{{IP: vipIP}}, updated.Status.LoadBalancer.Ingress)
	// bgp=true is no longer written by the controller; the speaker gates directly on bgp-vip.
	require.Empty(t, updated.Annotations[util.BgpAnnotation])
}

func TestCleanupBgpLbVipServiceBindingByVip(t *testing.T) {
	t.Parallel()

	fc, err := newFakeControllerWithOptions(t, nil)
	require.NoError(t, err)
	ctrl := fc.fakeController

	const ns = metav1.NamespaceDefault
	bindingSvc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "svc-bound",
			Namespace: ns,
			Annotations: map[string]string{
				util.BgpVipAnnotation: "vip-a",
			},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
		},
		Status: v1.ServiceStatus{
			LoadBalancer: v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{{IP: "203.0.113.10"}},
			},
		},
	}
	unrelatedSvc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "svc-other",
			Namespace: ns,
			Annotations: map[string]string{
				util.BgpVipAnnotation: "vip-b",
			},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
		},
		Status: v1.ServiceStatus{
			LoadBalancer: v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{{IP: "203.0.113.11"}},
			},
		},
	}

	err = fc.fakeInformers.serviceInformer.Informer().AddIndexers(cache.Indexers{
		bgpVipIndexName: func(obj any) ([]string, error) {
			s, ok := obj.(*v1.Service)
			if !ok {
				return nil, nil
			}
			var keys []string
			if v := s.Annotations[util.MetalLBAllowSharedIPAnnotation]; v != "" {
				keys = append(keys, v)
			} else if v := s.Annotations[util.BgpVipAnnotation]; v != "" {
				keys = append(keys, v)
			}
			return keys, nil
		},
	})
	require.NoError(t, err)
	ctrl.svcByBgpVipIndexer = fc.fakeInformers.serviceInformer.Informer().GetIndexer()

	for _, svc := range []*v1.Service{bindingSvc, unrelatedSvc} {
		err = fc.fakeInformers.serviceInformer.Informer().GetIndexer().Add(svc)
		require.NoError(t, err)
		_, err = ctrl.config.KubeClient.CoreV1().Services(ns).Create(context.Background(), svc, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	require.NoError(t, ctrl.cleanupBgpLbVipServiceBindingByVip("vip-a"))

	cleaned, err := ctrl.config.KubeClient.CoreV1().Services(ns).Get(context.Background(), bindingSvc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Empty(t, cleaned.Status.LoadBalancer.Ingress)
	// BgpVipAnnotation is a user-managed field; cleanupBgpLbVipServiceBinding must
	// NOT remove it so the Service can auto-rebind if the VIP is re-created.
	require.Equal(t, "vip-a", cleaned.Annotations[util.BgpVipAnnotation])

	stillBound, err := ctrl.config.KubeClient.CoreV1().Services(ns).Get(context.Background(), unrelatedSvc.Name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Equal(t, []v1.LoadBalancerIngress{{IP: "203.0.113.11"}}, stillBound.Status.LoadBalancer.Ingress)
	require.Equal(t, "vip-b", stillBound.Annotations[util.BgpVipAnnotation])
}

func TestHandleAddServicePrefersBgpLbVipWhenBothEnabled(t *testing.T) {
	t.Parallel()

	newService := func() *v1.Service {
		return &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lb-svc",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					util.AttachmentProvider: "default/provider",
				},
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{{
					Name:     "tcp",
					Protocol: v1.ProtocolTCP,
					Port:     80,
				}},
			},
		}
	}

	t.Run("both enabled uses bgp path and returns noop without bgp-vip annotation", func(t *testing.T) {
		fc, err := newFakeControllerWithOptions(t, nil)
		require.NoError(t, err)

		ctrl := fc.fakeController
		ctrl.svcKeyMutex = keymutex.NewHashed(0)
		ctrl.config.EnableBgpLbVip = true
		ctrl.config.EnablePodLbSvc = true

		svc := newService()
		require.NoError(t, fc.fakeInformers.serviceInformer.Informer().GetIndexer().Add(svc))

		key := cache.MetaObjectToName(svc).String()
		require.NoError(t, ctrl.handleAddService(key))
	})

	t.Run("only lb-svc enabled enters lb-svc path and errors without NAD", func(t *testing.T) {
		fc, err := newFakeControllerWithOptions(t, nil)
		require.NoError(t, err)

		ctrl := fc.fakeController
		ctrl.svcKeyMutex = keymutex.NewHashed(0)
		ctrl.config.EnableBgpLbVip = false
		ctrl.config.EnablePodLbSvc = true

		svc := newService()
		require.NoError(t, fc.fakeInformers.serviceInformer.Informer().GetIndexer().Add(svc))

		key := cache.MetaObjectToName(svc).String()
		require.Error(t, ctrl.handleAddService(key))
	})
}
