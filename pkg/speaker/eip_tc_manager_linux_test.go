//go:build linux

package speaker

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/cilium/ebpf"
	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
)

// 本集成测试使用真实内核 Map 与私有测试接口。必须在隔离且具备特权的 Network Namespace 中运行，
// 严禁直接作用于宿主机生产节点已有的 TC 程序。
func requireBPFIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("KUBE_OVN_BPF_TEST") != "1" {
		t.Skip("requires an isolated privileged Linux network namespace")
	}
}

// newTestTCManager 构建测试专用的 linuxTCManager 实例并注入预设大小的 BPF Hash Map。
func newTestTCManager(t *testing.T) (*linuxTCManager, *ebpf.Map, *ebpf.Map) {
	t.Helper()
	requireBPFIntegration(t)
	targetSize := uint32(unsafe.Sizeof(bpfEipTarget{}))
	v4, err := ebpf.NewMap(&ebpf.MapSpec{Type: ebpf.Hash, KeySize: 4, ValueSize: targetSize, MaxEntries: 64})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, v4.Close()) })
	v6, err := ebpf.NewMap(&ebpf.MapSpec{Type: ebpf.Hash, KeySize: 16, ValueSize: targetSize, MaxEntries: 64})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, v6.Close()) })
	m := newTCManager(&Configuration{EnableTcBpfRedirect: true}).(*linuxTCManager)
	m.v4Map, m.v6Map = v4, v6
	return m, v4, v6
}

func TestTCIPv4PacketKey(t *testing.T) {
	m, v4, _ := newTestTCManager(t)
	ip := net.ParseIP("111.62.156.190")
	require.NoError(t, m.AddEIP(ip, TCRedirectTarget{TargetIfindex: 10, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}))
	var value bpfEipTarget
	require.NoError(t, v4.Lookup([]byte{111, 62, 156, 190}, &value), "map key must match packet bytes")
	require.Equal(t, uint32(10), value.Ifindex)
	require.NoError(t, m.DeleteEIP(ip))
	require.ErrorIs(t, v4.Lookup([]byte{111, 62, 156, 190}, &value), ebpf.ErrKeyNotExist)
}

func TestTCSyncRemovesStaleEntries(t *testing.T) {
	m, v4, v6 := newTestTCManager(t)
	old := bpfEipTarget{Ifindex: 10}
	require.NoError(t, v4.Put([]byte{1, 2, 3, 4}, old))
	require.NoError(t, v6.Put(net.ParseIP("2001:db8::1").To16(), old))
	require.NoError(t, m.SyncEIPs(map[string]TCRedirectTarget{}))
	var value bpfEipTarget
	require.ErrorIs(t, v4.Lookup([]byte{1, 2, 3, 4}, &value), ebpf.ErrKeyNotExist)
	require.ErrorIs(t, v6.Lookup(net.ParseIP("2001:db8::1").To16(), &value), ebpf.ErrKeyNotExist)
}

func TestTCRejectsInvalidTarget(t *testing.T) {
	m, _, _ := newTestTCManager(t)
	for _, target := range []TCRedirectTarget{{}, {TargetIfindex: 10}, {TargetIfindex: -1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}} {
		require.Error(t, m.AddEIP(net.ParseIP("2001:db8::1"), target))
	}
	require.Error(t, m.SyncEIPs(map[string]TCRedirectTarget{"invalid": {TargetIfindex: 10}}))
}

func TestTCAttachMissingObject(t *testing.T) {
	requireBPFIntegration(t)
	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "eip-test0"}}
	require.NoError(t, netlink.LinkAdd(link))
	t.Cleanup(func() { require.NoError(t, netlink.LinkDel(link)) })
	m := NewTCManager(&Configuration{EnableTcBpfRedirect: true, TcBpfProgPath: "/does-not-exist.o"})
	require.Error(t, m.Attach("eip-test0"))
	require.NoError(t, m.Close())
}

func TestTCNonEIPContinuesClassification(t *testing.T) {
	requireBPFIntegration(t)
	path := os.Getenv("KUBE_OVN_BPF_OBJECT")
	require.NotEmpty(t, path)
	spec, err := ebpf.LoadCollectionSpec(path)
	require.NoError(t, err)
	for _, m := range spec.Maps {
		m.Pinning = ebpf.PinNone
	}
	coll, err := ebpf.NewCollection(spec)
	require.NoError(t, err)
	defer coll.Close()
	for _, protocol := range []uint16{unix.ETH_P_IP, unix.ETH_P_IPV6, unix.ETH_P_ARP, unix.ETH_P_8021Q} {
		packet := make([]byte, 64)
		packet[12], packet[13], packet[14] = byte(protocol>>8), byte(protocol), 0x45
		result, _, err := coll.Programs["tc_eip_ingress"].Test(packet)
		require.NoError(t, err)
		require.Equal(t, uint32(0xffffffff), result, "non-EIP traffic must continue to subsequent TC filters")
	}
	for _, ipv6 := range []bool{false, true} {
		packet := make([]byte, 80)
		packet[12], packet[13], packet[14] = 8, 0, 0x45
		key := []byte{192, 0, 2, 10}
		copy(packet[30:34], key)
		m := coll.Maps[tcEIPv4MapName]
		if ipv6 {
			packet[12], packet[13], packet[14] = 0x86, 0xdd, 0x60
			key = net.ParseIP("2001:db8::10").To16()
			copy(packet[38:54], key)
			m = coll.Maps[tcEIPv6MapName]
		}
		target := bpfEipTarget{Ifindex: 10, ParentIfindex: 1, TargetMAC: [6]byte{2, 0, 0, 0, 0, 1}}
		require.NoError(t, m.Put(key, target))
		result, out, err := coll.Programs["tc_eip_ingress"].Test(packet)
		require.NoError(t, err)
		require.Equal(t, uint32(7), result, "EIP must redirect")
		require.Equal(t, target.TargetMAC[:], out[:6])
	}
}

func testTCPinPath(t *testing.T) string {
	t.Helper()
	requireBPFIntegration(t)
	// bpffs 作用域限定在此一次性测试容器的挂载命名空间内
	if err := unix.Mount("bpf", "/sys/fs/bpf", "bpf", 0, ""); err != nil {
		require.ErrorIs(t, err, unix.EBUSY)
	}
	dir, err := os.MkdirTemp("/sys/fs/bpf", "eip-test-")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.RemoveAll(dir)) })
	return filepath.Join(dir, "eip_v4_map")
}

func TestTCRetryRestartAndDisable(t *testing.T) {
	pin := testTCPinPath(t)
	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "eip-retry0"}}
	require.NoError(t, netlink.LinkAdd(link))
	t.Cleanup(func() { require.NoError(t, netlink.LinkDel(link)) })
	cfg := &Configuration{
		EnableTcBpfRedirect: true,
		TcBpfInterface:      link.Name,
		TcBpfProgPath:       "/missing.o",
		TcBpfMapPath:        pin,
	}
	m := NewTCManager(cfg)
	require.Error(t, m.Attach(link.Name))
	cfg.TcBpfProgPath = os.Getenv("KUBE_OVN_BPF_OBJECT")
	require.NoError(t, m.Attach(link.Name), "retry must load the object after it becomes available")
	require.NoError(t, m.Attach(link.Name), "attachment must be idempotent")
	ip := net.ParseIP("2001:db8::1")
	require.NoError(t, m.AddEIP(ip, TCRedirectTarget{TargetIfindex: link.Index, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}))
	// 模拟进程未优雅退出即重启。旧 TC 分类器程序保持挂载，新实例必须丢弃旧的陈旧 Map
	restarted := NewTCManager(cfg)
	require.NoError(t, restarted.Attach(link.Name))
	require.NoError(t, restarted.SyncEIPs(map[string]TCRedirectTarget{}))
	v6, err := ebpf.LoadPinnedMap(strings.TrimSuffix(pin, "_v4_map")+"_v6_map", nil)
	require.NoError(t, err)
	var value bpfEipTarget
	require.ErrorIs(t, v6.Lookup(ip.To16(), &value), ebpf.ErrKeyNotExist)
	require.NoError(t, v6.Close())
	// 停用 EIP 旁路时，必须保留挂载在其他 priority 上的无关过滤器
	foreign := &netlink.MatchAll{FilterAttrs: netlink.FilterAttrs{LinkIndex: link.Index, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 2, Priority: 10, Protocol: unix.ETH_P_ALL}}
	require.NoError(t, netlink.FilterAdd(foreign))
	disabled := NewTCManager(&Configuration{TcBpfMapPath: pin})
	require.NoError(t, disabled.Clean())
	filters, err := netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	require.NoError(t, err)
	require.Len(t, filters, 1)
	require.Equal(t, uint16(10), filters[0].Attrs().Priority)
	_, err = os.Stat(pin)
	require.ErrorIs(t, err, os.ErrNotExist)
	require.NoError(t, restarted.Close())
	require.NoError(t, m.Close())
	require.Error(t, m.Attach(link.Name), "a stopped worker cannot resurrect an attachment")
}

func TestTCSyncPreservesExpectedAndCleansStale(t *testing.T) {
	m, v4, _ := newTestTCManager(t)
	ip := net.ParseIP("192.0.2.10")
	target := TCRedirectTarget{TargetIfindex: 10, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}
	require.NoError(t, m.AddEIP(ip, target))
	var value bpfEipTarget
	require.NoError(t, v4.Lookup(ip.To4(), &value))
	require.NoError(t, v4.Put([]byte{1, 2, 3, 4}, value))
	require.NoError(t, m.SyncEIPs(map[string]TCRedirectTarget{ip.String(): target}))
	require.NoError(t, v4.Lookup(ip.To4(), &value))
	require.ErrorIs(t, v4.Lookup([]byte{1, 2, 3, 4}, &value), ebpf.ErrKeyNotExist)
	// 非法的替换目标绝不能保留旧的转发条目
	require.Error(t, m.SyncEIPs(map[string]TCRedirectTarget{ip.String(): {}}))
	require.ErrorIs(t, v4.Lookup(ip.To4(), &value), ebpf.ErrKeyNotExist)
}

func TestTCSyncRetainedDeclarativeProtection(t *testing.T) {
	m, v4, _ := newTestTCManager(t)
	dummy := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: fmt.Sprintf("eip-ret-%d", time.Now().UnixNano()%100000)}}
	require.NoError(t, netlink.LinkAdd(dummy))
	t.Cleanup(func() { _ = netlink.LinkDel(dummy) })

	ipRetain := net.ParseIP("192.0.2.20")
	ipStale := net.ParseIP("192.0.2.30")

	// 1. 预先注入两个 EIP 条目
	require.NoError(t, m.AddEIP(ipRetain, TCRedirectTarget{TargetIfindex: dummy.Index, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}))
	require.NoError(t, m.AddEIP(ipStale, TCRedirectTarget{TargetIfindex: dummy.Index, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 2}}))

	// 2. 模拟全量 Reconcile：
	// expected 为空（假设两个 EIP 都未能解析出最新目标）
	// 但 ipRetain 属于本地活跃 Pod 暂时抖动（传入 retainedIPs）
	// ipStale 则已彻底删除（未传入 retainedIPs）
	require.NoError(t, m.SyncEIPs(map[string]TCRedirectTarget{}, ipRetain.String()))

	// 3. 验证声明保留的 ipRetain 依然稳稳留在 BPF Map 中
	var val bpfEipTarget
	require.NoError(t, v4.Lookup(ipRetain.To4(), &val), "活跃声明保留的 EIP 必须留在 Map 中不被误清空")

	// 4. 验证未声明保留的 ipStale 被干净清理
	require.ErrorIs(t, v4.Lookup(ipStale.To4(), &val), ebpf.ErrKeyNotExist, "不再期望存在的 EIP 必须被 prune 清理")
}

func TestTCMacvlanDelivery(t *testing.T) {
	pin := testTCPinPath(t)
	run := func(command string) {
		t.Helper()
		output, err := exec.Command("sh", "-ec", command).CombinedOutput()
		require.NoError(t, err, "%s: %s", command, output)
	}
	// 外部发送方与网关 Pod 分别处于独立的网络命名空间。
	// 宿主机网络故意不配置到 EIP 的路由；ping 连通必须依赖 TC 旁路直通转发生效。
	t.Cleanup(func() {
		run("ip netns del eip-outside 2>/dev/null || true; ip netns del eip-gateway 2>/dev/null || true; ip link del eip-uplink 2>/dev/null || true")
	})
	run(`ip netns add eip-outside
ip netns add eip-gateway
ip link add eip-uplink type veth peer name eip-external
ip link set eip-external netns eip-outside
ip addr add 198.18.0.1/24 dev eip-uplink
ip link set eip-uplink up
ip link add eip-host link eip-uplink type macvlan mode bridge
ip link set eip-host up
ip link add eip-pod link eip-uplink address 02:00:00:00:10:02 type macvlan mode bridge
ip link set eip-pod netns eip-gateway
ip -n eip-gateway link set lo up
ip -n eip-gateway link set eip-pod up
ip -n eip-gateway addr add 198.18.0.2/24 dev eip-pod
ip -n eip-gateway addr add 203.0.113.10/32 dev eip-pod
ip -n eip-gateway route add default via 198.18.0.3
ip -n eip-outside link set lo up
ip -n eip-outside link set eip-external up
ip -n eip-outside addr add 198.18.0.3/24 dev eip-external
ip -n eip-outside route add 203.0.113.10/32 via 198.18.0.1`)
	m := NewTCManager(&Configuration{EnableTcBpfRedirect: true, TcBpfInterface: "eip-uplink", TcBpfProgPath: os.Getenv("KUBE_OVN_BPF_OBJECT"), TcBpfMapPath: pin})
	t.Cleanup(func() { require.NoError(t, m.Close()) })
	require.NoError(t, m.Attach("eip-uplink"))
	host, err := netlink.LinkByName("eip-host")
	require.NoError(t, err)
	require.NoError(t, m.AddEIP(net.ParseIP("203.0.113.10"), TCRedirectTarget{TargetIfindex: host.Attrs().Index, ParentIfindex: host.Attrs().ParentIndex, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0x10, 2}}))
	output, pingErr := exec.Command("ip", "netns", "exec", "eip-outside", "ping", "-c", "2", "-W", "2", "203.0.113.10").CombinedOutput()
	var counters bpfEipTarget
	require.NoError(t, m.(*linuxTCManager).v4Map.Lookup([]byte{203, 0, 113, 10}, &counters))
	require.Equal(t, uint32(host.Attrs().Index), counters.Ifindex, "target ifindex must match")
	require.NoError(t, pingErr, "%s", output)
}

func TestTCPriorityMigrationAndCleanup(t *testing.T) {
	requireBPFIntegration(t)
	pin := testTCPinPath(t)
	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: fmt.Sprintf("eip-prio-%d", time.Now().UnixNano()%100000)}}
	require.NoError(t, netlink.LinkAdd(link))
	t.Cleanup(func() { _ = netlink.LinkDel(link) })

	// 1. 模拟旧实例：使用 Priority 5 挂载
	cfg1 := &Configuration{
		EnableTcBpfRedirect: true,
		TcBpfInterface:      link.Name,
		TcBpfProgPath:       os.Getenv("KUBE_OVN_BPF_OBJECT"),
		TcBpfMapPath:        pin,
		TcBpfPriority:       5,
	}
	m1 := NewTCManager(cfg1)
	require.NoError(t, m1.Attach(link.Name))

	// 验证 Priority 5 过滤器已就绪
	filters, err := netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	require.NoError(t, err)
	require.Len(t, filters, 1)
	require.Equal(t, uint16(5), filters[0].Attrs().Priority)

	// 2. 模拟运维调整 priority 为 7 并重启新实例
	cfg2 := &Configuration{
		EnableTcBpfRedirect: true,
		TcBpfInterface:      link.Name,
		TcBpfProgPath:       os.Getenv("KUBE_OVN_BPF_OBJECT"),
		TcBpfMapPath:        pin,
		TcBpfPriority:       7,
	}
	m2 := NewTCManager(cfg2)
	require.NoError(t, m2.Attach(link.Name))

	// 挂载第三方（如 Cilium/QoS）的过滤器，验证不会被误删
	foreign := &netlink.MatchAll{FilterAttrs: netlink.FilterAttrs{LinkIndex: link.Index, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 2, Priority: 10, Protocol: unix.ETH_P_ALL}}
	require.NoError(t, netlink.FilterAdd(foreign))

	// 3. 验证旧 Priority 5 过滤器已被自动剥离，当前仅留有 Priority 7 和第三方 Priority 10
	filters, err = netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	require.NoError(t, err)
	require.Len(t, filters, 2)
	var foundPrio7, foundPrio10 bool
	for _, f := range filters {
		if f.Attrs().Priority == 7 {
			foundPrio7 = true
		}
		if f.Attrs().Priority == 10 {
			foundPrio10 = true
		}
	}
	require.True(t, foundPrio7, "必须成功挂载新 Priority 7 规则")
	require.True(t, foundPrio10, "必须保留第三方 Priority 10 规则")

	// 4. 模拟停用功能执行 Clean()，验证本程序所有 priority 规则彻底清理，且不伤及第三方规则
	require.NoError(t, m2.Clean())
	filters, err = netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	require.NoError(t, err)
	require.Len(t, filters, 1)
	require.Equal(t, uint16(10), filters[0].Attrs().Priority)
}

func TestTCHotReloadPreservesMapAndFilters(t *testing.T) {
	requireBPFIntegration(t)
	pin := testTCPinPath(t)
	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: fmt.Sprintf("eip-reload-%d", time.Now().UnixNano()%100000)}}
	require.NoError(t, netlink.LinkAdd(link))
	t.Cleanup(func() { _ = netlink.LinkDel(link) })

	cfg := &Configuration{
		EnableTcBpfRedirect: true,
		TcBpfInterface:      link.Name,
		TcBpfProgPath:       os.Getenv("KUBE_OVN_BPF_OBJECT"),
		TcBpfMapPath:        pin,
	}

	// 1. 启动第一代实例并写入条目
	m1 := NewTCManager(cfg)
	require.NoError(t, m1.Attach(link.Name))
	ip := net.ParseIP("183.240.180.136")
	target := TCRedirectTarget{TargetIfindex: 100, ParentIfindex: 10, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0x18, 1}}
	require.NoError(t, m1.AddEIP(ip, target))

	// 2. 模拟进程优雅退出/热更：执行 Close()，此时内核 Filter 与 Pinned Map 必须完好保留
	require.NoError(t, m1.Close())
	filters, err := netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	require.NoError(t, err)
	require.Len(t, filters, 1, "TC filter must remain active after graceful Close")

	// 3. 启动第二代实例：应该成功复用既有 Map
	m2 := NewTCManager(cfg)
	require.NoError(t, m2.Attach(link.Name))

	// 验证第二代实例直接继承了第一代的 EIP 条目，业务在进程重启期间完全无中断
	var old bpfEipTarget
	require.NoError(t, m2.(*linuxTCManager).v4Map.Lookup(ip.To4(), &old), "reloaded manager must inherit existing entries")
	require.Equal(t, uint32(100), old.Ifindex)

	// 4. 显式 Clean 之后，规则与 Map 彻底被清理
	require.NoError(t, m2.Clean())
	filters, err = netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	require.NoError(t, err)
	require.Empty(t, filters, "TC filter must be removed after Clean")
	_, err = os.Stat(pin)
	require.True(t, os.IsNotExist(err), "pinned map must be removed after Clean")
}
