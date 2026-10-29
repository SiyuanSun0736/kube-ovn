//go:build !linux

package speaker

import (
	"net"

	"k8s.io/klog/v2"
)

// noopTCManager 为非 Linux 操作系统（如开发调试用的 macOS / Windows）提供空操作实现
type noopTCManager struct{}

// newPlatformTCManager 在非 Linux 平台下返回空操作管理器，保障跨平台编译无障碍
func newPlatformTCManager(_ *Configuration) TCManager {
	return &noopTCManager{}
}

func (m *noopTCManager) Attach(linkName string) error {
	klog.V(4).Infof("TCManager: 非 Linux 环境忽略网卡 %s 的 TC 挂载请求", linkName)
	return nil
}

func (m *noopTCManager) Detach(linkName string) error {
	klog.V(4).Infof("TCManager: 非 Linux 环境忽略网卡 %s 的 TC 卸载请求", linkName)
	return nil
}

func (m *noopTCManager) AddEIP(eip net.IP, info TargetInfo) error {
	klog.V(4).Infof("TCManager: 非 Linux 环境忽略 AddEIP 操作 (EIP: %s, 目标 ifindex: %d)", eip, info.TargetIfindex)
	return nil
}

func (m *noopTCManager) DeleteEIP(eip net.IP) error {
	klog.V(4).Infof("TCManager: 非 Linux 环境忽略 DeleteEIP 操作 (EIP: %s)", eip)
	return nil
}

func (m *noopTCManager) SyncEIPs(expected map[string]TargetInfo) error {
	return nil
}

func (m *noopTCManager) Close() error {
	return nil
}
