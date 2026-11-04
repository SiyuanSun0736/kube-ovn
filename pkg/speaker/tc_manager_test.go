//go:build linux

package speaker

import (
	"net"
	"testing"
	"unsafe"

	"github.com/vishvananda/netlink"
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

// TestTCManagerPriority 验证 TC 过滤器优先级默认值及自定义配置
func TestTCManagerPriority(t *testing.T) {
	// 默认优先级应为 5，预留 1-4 给更高优先级防护探针
	mgrDefault := NewTCManager(&Configuration{})
	if lm, ok := mgrDefault.(*linuxTCManager); ok {
		if p := lm.priority(); p != 5 {
			t.Fatalf("期望默认 priority 5，实际为 %d", p)
		}
	}

	// 自定义优先级覆盖生效
	mgrCustom := NewTCManager(&Configuration{TcBpfPriority: 7})
	if lm, ok := mgrCustom.(*linuxTCManager); ok {
		if p := lm.priority(); p != 7 {
			t.Fatalf("期望自定义 priority 7，实际为 %d", p)
		}
	}
}

// TestTCManagerIsEIPFilter 验证过滤器精准匹配逻辑（仅匹配当前配置的 priority，不匹配历史残留 priority 1）
func TestTCManagerIsEIPFilter(t *testing.T) {
	mgr := NewTCManager(&Configuration{TcBpfPriority: 5}).(*linuxTCManager)

	// 匹配配置 priority 5
	fMatch := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			Handle:   1,
			Priority: 5,
		},
		Name: tcEIPFilterName,
	}
	if !mgr.isEIPFilter(fMatch) {
		t.Errorf("期望匹配 priority 5 的本程序过滤器，但匹配失败")
	}

	// 历史 priority 1 不再匹配
	fOldPriority := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			Handle:   1,
			Priority: 1,
		},
		Name: tcEIPFilterName,
	}
	if mgr.isEIPFilter(fOldPriority) {
		t.Errorf("priority 1 过滤器不应被匹配")
	}

	// handle 不为 1 不匹配
	fDiffHandle := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			Handle:   2,
			Priority: 5,
		},
		Name: tcEIPFilterName,
	}
	if mgr.isEIPFilter(fDiffHandle) {
		t.Errorf("handle!=1 过滤器不应被匹配")
	}

	// 非本程序名字不匹配
	fOtherName := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			Handle:   1,
			Priority: 5,
		},
		Name: "cilium-netdev",
	}
	if mgr.isEIPFilter(fOtherName) {
		t.Errorf("非 kube-ovn 过滤器不应被匹配")
	}
}

// TestTCManagerAttachRequiresInterface 验证未指定接口时 Attach 明确报错，禁止误挂载到默认路由网卡
func TestTCManagerAttachRequiresInterface(t *testing.T) {
	mgr := NewTCManager(&Configuration{
		EnableTcBpfRedirect: true,
		TcBpfInterface:      "",
	})
	if err := mgr.Attach(""); err == nil {
		t.Fatal("期望未指定接口时 Attach 报错，但返回了 nil")
	}
}
