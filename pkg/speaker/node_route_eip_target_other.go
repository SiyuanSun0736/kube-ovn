//go:build !linux

package speaker

import (
	"errors"

	corev1 "k8s.io/api/core/v1"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
)

func (c *Controller) getNatGwTargetInfo(_ *corev1.Pod, _ *kubeovnv1.IptablesEIP) (TargetInfo, error) {
	return TargetInfo{}, errors.New("TC EIP redirection requires Linux")
}
