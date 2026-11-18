//go:build !linux

package daemon

import (
	corev1 "k8s.io/api/core/v1"
)

type portForwardConnLimitRule struct {
	ip        string
	port      int32
	protocol  string
	connLimit int
}

func (c *Controller) parsePortForwardConnLimitRules(svc *corev1.Service) []portForwardConnLimitRule {
	return nil
}

func (c *Controller) reconcilePortForwardConnLimit(oldService, newService *corev1.Service) error {
	return nil
}
