//go:build linux

package speaker

import (
	"net"
	"testing"
	"unsafe"
)

// TestTCManagerTargetInfoLayout 严格断言内核 struct eip_target 与 Go 结构体 bpfEipTarget 在 64 位体系下的内存对齐
func TestTCManagerTargetInfoLayout(t *testing.T) {
	var target bpfEipTarget
	// 4 (Ifindex) + 6 (TargetMAC) + 2 (Pad) + 4 (Pad2) + 8 (Packets) + 8 (Bytes) = 32 字节
	if size := unsafe.Sizeof(target); size != 32 {
		t.Fatalf("bpfEipTarget 结构体内存大小异常: 期望 32 字节 (64-bit 严格对齐)，实际为 %d 字节", size)
	}

	// 验证 64 位 Packets 字段严格从 16 字节偏移量开始
	if offset := unsafe.Offsetof(target.Packets); offset != 16 {
		t.Fatalf("Packets 字段偏移量异常: 期望从 16 字节开始，实际为 %d", offset)
	}
}

// TestTCManagerLifecycle 验证 TCManager 初始化与基础生命周期调用安全
func TestTCManagerLifecycle(t *testing.T) {
	cfg := &Configuration{
		EnableTcBpfRedirect: false,
		TcBpfInterface:      "bond0",
		TcBpfProgPath:       "/etc/kube-ovn/tc_eip.o",
	}

	mgr := NewTCManager(cfg)
	if mgr == nil {
		t.Fatal("期望创建非空的 TCManager 实例")
	}

	eip := net.ParseIP("111.62.156.190")
	info := TargetInfo{
		TargetIfindex: 10,
		TargetMAC:     net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
	}

	// 验证在模拟测试环境下的接口增删调用不发生 panic 崩溃
	if err := mgr.AddEIP(eip, info); err != nil {
		t.Errorf("AddEIP: %v", err)
	}

	if err := mgr.DeleteEIP(eip); err != nil {
		t.Errorf("DeleteEIP: %v", err)
	}
}
