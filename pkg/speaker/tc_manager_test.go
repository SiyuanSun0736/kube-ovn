//go:build linux

package speaker

import (
	"net"
	"testing"
	"unsafe"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// TestTCManagerRedirectTargetLayout 严格断言内核 struct eip_target 与 Go 结构体 bpfEipTarget 在 64 位体系下的内存对齐
func TestTCManagerRedirectTargetLayout(t *testing.T) {
	var target bpfEipTarget
	// 4 (Ifindex) + 6 (TargetMAC) + 2 (Pad) + 4 (ParentIfindex) = 16 字节
	if size := unsafe.Sizeof(target); size != 16 {
		t.Fatalf("bpfEipTarget struct size mismatch: expected 16 bytes, got %d", size)
	}

	// 验证 ParentIfindex 字段严格从 12 字节偏移量开始
	if offset := unsafe.Offsetof(target.ParentIfindex); offset != 12 {
		t.Fatalf("ParentIfindex field offset mismatch: expected 12 bytes, got %d", offset)
	}
}

// TestNewTCManager 验证 TCManager 初始化与基础生命周期调用安全
func TestNewTCManager(t *testing.T) {
	cfg := &Configuration{
		EnableTcBpfRedirect: false,
		TcBpfInterface:      "bond0",
		TcBpfProgPath:       "/etc/kube-ovn/tc_eip.o",
	}

	mgr := NewTCManager(cfg)
	if mgr == nil {
		t.Fatal("expected non-nil TCManager instance")
	}

	eip := net.ParseIP("111.62.156.190")
	target := TCRedirectTarget{
		TargetIfindex: 10,
		TargetMAC:     net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
	}

	// 验证在模拟测试环境下的接口增删调用不发生 panic 崩溃
	if err := mgr.AddEIP(eip, target); err != nil {
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
			t.Fatalf("expected default priority 5, got %d", p)
		}
	}

	// 自定义优先级覆盖生效
	mgrCustom := NewTCManager(&Configuration{TcBpfPriority: 7})
	if lm, ok := mgrCustom.(*linuxTCManager); ok {
		if p := lm.priority(); p != 7 {
			t.Fatalf("expected custom priority 7, got %d", p)
		}
	}
}

// TestTCManagerIsEIPFilter 验证过滤器多维交叉核验逻辑（多重属性 + BPF Tag / ID 校验）
func TestTCManagerIsEIPFilter(t *testing.T) {
	mgr := NewTCManager(&Configuration{TcBpfPriority: 5}).(*linuxTCManager)
	mgr.progID = 42
	mgr.progTag = "08218004caf6a189"

	baseFilter := func() *netlink.BpfFilter {
		return &netlink.BpfFilter{
			FilterAttrs: netlink.FilterAttrs{
				Handle:   1,
				Parent:   netlink.HANDLE_MIN_INGRESS,
				Protocol: unix.ETH_P_ALL,
				Priority: 5,
			},
			DirectAction: true,
			Name:         tcEIPFilterName,
		}
	}

	// 1. 标准名字匹配
	fMatch := baseFilter()
	if !mgr.isEIPFilter(fMatch) {
		t.Errorf("standard kube-ovn filter should match")
	}

	// 2. 跨 priority 识别
	fOtherPriority := baseFilter()
	fOtherPriority.Priority = 7
	if !mgr.isEIPFilter(fOtherPriority) {
		t.Errorf("filter with different priority should match")
	}

	// 3. 强指纹：通过 BPF Program ID 匹配
	fIDMatch := baseFilter()
	fIDMatch.Name = ""
	fIDMatch.Id = 42
	if !mgr.isEIPFilter(fIDMatch) {
		t.Errorf("filter matching program ID should match")
	}

	// 4. 强指纹：通过 BPF Tag 匹配
	fTagMatch := baseFilter()
	fTagMatch.Name = ""
	fTagMatch.Tag = "08218004caf6a189"
	if !mgr.isEIPFilter(fTagMatch) {
		t.Errorf("filter matching BPF tag should match")
	}

	// 5. Parent 不为 Ingress 拒绝
	fEgress := baseFilter()
	fEgress.Parent = netlink.HANDLE_MIN_EGRESS
	if mgr.isEIPFilter(fEgress) {
		t.Errorf("egress filter should not match")
	}

	// 6. DirectAction 为 false 拒绝
	fNoDA := baseFilter()
	fNoDA.DirectAction = false
	if mgr.isEIPFilter(fNoDA) {
		t.Errorf("filter without DirectAction should not match")
	}

	// 7. Handle!=1 拒绝
	fDiffHandle := baseFilter()
	fDiffHandle.Handle = 2
	if mgr.isEIPFilter(fDiffHandle) {
		t.Errorf("filter with Handle!=1 should not match")
	}

	// 8. 非本程序名字拒绝（如 Cilium 且无匹配 Tag/ID）
	fCilium := baseFilter()
	fCilium.Name = "cilium-netdev"
	if mgr.isEIPFilter(fCilium) {
		t.Errorf("non-kube-ovn filter without matching fingerprints should not match")
	}
}

// TestTCManagerAttachRequiresInterface 验证未指定接口时 Attach 明确报错，禁止误挂载到默认路由网卡
func TestTCManagerAttachRequiresInterface(t *testing.T) {
	mgr := NewTCManager(&Configuration{
		EnableTcBpfRedirect: true,
		TcBpfInterface:      "",
	})
	if err := mgr.Attach(""); err == nil {
		t.Fatal("expected error when attaching without interface, got nil")
	}
}
