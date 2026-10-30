//go:build linux

package speaker

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/unix"

	"github.com/cilium/ebpf"
	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
)

// These tests use real kernel maps and private interfaces. Run in an isolated,
// privileged network namespace, never against a node's existing TC programs.
func requireBPFIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("KUBE_OVN_BPF_TEST") != "1" {
		t.Skip("requires an isolated privileged Linux network namespace")
	}
}

func newTestTCMaps(t *testing.T) (*linuxTCManager, *ebpf.Map, *ebpf.Map) {
	t.Helper()
	requireBPFIntegration(t)
	v4, err := ebpf.NewMap(&ebpf.MapSpec{Type: ebpf.Hash, KeySize: 4, ValueSize: 32, MaxEntries: 64})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, v4.Close()) })
	v6, err := ebpf.NewMap(&ebpf.MapSpec{Type: ebpf.Hash, KeySize: 16, ValueSize: 32, MaxEntries: 64})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, v6.Close()) })
	m := newPlatformTCManager(&Configuration{EnableTcBpfRedirect: true}).(*linuxTCManager)
	m.v4Map, m.v6Map = v4, v6
	return m, v4, v6
}

func TestTCIPv4PacketKey(t *testing.T) {
	m, v4, _ := newTestTCMaps(t)
	ip := net.ParseIP("111.62.156.190")
	require.NoError(t, m.AddEIP(ip, TargetInfo{TargetIfindex: 10, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}))
	var value bpfEipTarget
	require.NoError(t, v4.Lookup([]byte{111, 62, 156, 190}, &value), "map key must match packet bytes")
	require.Equal(t, uint32(10), value.Ifindex)
	require.NoError(t, m.DeleteEIP(ip))
	require.ErrorIs(t, v4.Lookup([]byte{111, 62, 156, 190}, &value), ebpf.ErrKeyNotExist)
}

func TestTCSyncRemovesStaleEntries(t *testing.T) {
	m, v4, v6 := newTestTCMaps(t)
	old := bpfEipTarget{Ifindex: 10}
	require.NoError(t, v4.Put([]byte{1, 2, 3, 4}, old))
	require.NoError(t, v6.Put(net.ParseIP("2001:db8::1").To16(), old))
	require.NoError(t, m.SyncEIPs(map[string]TargetInfo{}))
	var value bpfEipTarget
	require.ErrorIs(t, v4.Lookup([]byte{1, 2, 3, 4}, &value), ebpf.ErrKeyNotExist)
	require.ErrorIs(t, v6.Lookup(net.ParseIP("2001:db8::1").To16(), &value), ebpf.ErrKeyNotExist)
}

func TestTCRejectsInvalidTarget(t *testing.T) {
	m, _, _ := newTestTCMaps(t)
	for _, info := range []TargetInfo{{}, {TargetIfindex: 10}, {TargetIfindex: -1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}} {
		require.Error(t, m.AddEIP(net.ParseIP("2001:db8::1"), info))
	}
	require.Error(t, m.SyncEIPs(map[string]TargetInfo{"invalid": {TargetIfindex: 10}}))
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
	// bpffs is local to the disposable container's mount namespace.
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
	cfg := &Configuration{EnableTcBpfRedirect: true, TcBpfProgPath: "/missing.o", TcBpfMapPath: pin}
	m := NewTCManager(cfg)
	require.Error(t, m.Attach(link.Name))
	cfg.TcBpfProgPath = os.Getenv("KUBE_OVN_BPF_OBJECT")
	require.NoError(t, m.Attach(link.Name), "retry must load the object after it becomes available")
	require.NoError(t, m.Attach(link.Name), "attachment must be idempotent")
	ip := net.ParseIP("2001:db8::1")
	require.NoError(t, m.AddEIP(ip, TargetInfo{TargetIfindex: link.Index, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}))
	// Simulate a process restart without clean shutdown. The old program remains
	// attached; a new instance must discard its stale maps.
	restarted := NewTCManager(cfg)
	require.NoError(t, restarted.Attach(link.Name))
	require.NoError(t, restarted.SyncEIPs(map[string]TargetInfo{}))
	v6, err := ebpf.LoadPinnedMap(strings.TrimSuffix(pin, "_v4_map")+"_v6_map", nil)
	require.NoError(t, err)
	var value bpfEipTarget
	require.ErrorIs(t, v6.Lookup(ip.To16(), &value), ebpf.ErrKeyNotExist)
	require.NoError(t, v6.Close())
	// Keep an unrelated filter at another priority while disabling EIP bypass.
	foreign := &netlink.MatchAll{FilterAttrs: netlink.FilterAttrs{LinkIndex: link.Index, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 2, Priority: 10, Protocol: unix.ETH_P_ALL}}
	require.NoError(t, netlink.FilterAdd(foreign))
	disabled := NewTCManager(&Configuration{TcBpfMapPath: pin})
	require.NoError(t, disabled.Close())
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

func TestTCSyncPreservesExpectedAndCounters(t *testing.T) {
	m, v4, _ := newTestTCMaps(t)
	ip := net.ParseIP("192.0.2.10")
	info := TargetInfo{TargetIfindex: 10, ParentIfindex: 1, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0, 1}}
	require.NoError(t, m.AddEIP(ip, info))
	var value bpfEipTarget
	require.NoError(t, v4.Lookup(ip.To4(), &value))
	value.Packets = 42
	require.NoError(t, v4.Put(ip.To4(), value))
	require.NoError(t, v4.Put([]byte{1, 2, 3, 4}, value))
	require.NoError(t, m.SyncEIPs(map[string]TargetInfo{ip.String(): info}))
	require.NoError(t, v4.Lookup(ip.To4(), &value))
	require.Equal(t, uint64(42), value.Packets)
	require.ErrorIs(t, v4.Lookup([]byte{1, 2, 3, 4}, &value), ebpf.ErrKeyNotExist)
	// An invalid replacement must not retain the old forwarding entry.
	require.Error(t, m.SyncEIPs(map[string]TargetInfo{ip.String(): {}}))
	require.ErrorIs(t, v4.Lookup(ip.To4(), &value), ebpf.ErrKeyNotExist)
}

func TestTCMacvlanDelivery(t *testing.T) {
	pin := testTCPinPath(t)
	run := func(command string) {
		t.Helper()
		output, err := exec.Command("sh", "-ec", command).CombinedOutput()
		require.NoError(t, err, "%s: %s", command, output)
	}
	// Both the outside sender and gateway are in separate namespaces. The host
	// deliberately has no route to the EIP; successful ping requires TC bypass.
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
	m := NewTCManager(&Configuration{EnableTcBpfRedirect: true, TcBpfProgPath: os.Getenv("KUBE_OVN_BPF_OBJECT"), TcBpfMapPath: pin})
	t.Cleanup(func() { require.NoError(t, m.Close()) })
	require.NoError(t, m.Attach("eip-uplink"))
	host, err := netlink.LinkByName("eip-host")
	require.NoError(t, err)
	require.NoError(t, m.AddEIP(net.ParseIP("203.0.113.10"), TargetInfo{TargetIfindex: host.Attrs().Index, ParentIfindex: host.Attrs().ParentIndex, TargetMAC: net.HardwareAddr{2, 0, 0, 0, 0x10, 2}}))
	output, pingErr := exec.Command("ip", "netns", "exec", "eip-outside", "ping", "-c", "2", "-W", "2", "203.0.113.10").CombinedOutput()
	var counters bpfEipTarget
	require.NoError(t, m.(*linuxTCManager).v4Map.Lookup([]byte{203, 0, 113, 10}, &counters))
	require.Equal(t, uint64(2), counters.Packets, "local macvlan delivery must not recirculate through TC")
	require.NoError(t, pingErr, "%s", output)
}
