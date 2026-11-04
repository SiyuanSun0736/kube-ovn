//go:build linux

package speaker

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"k8s.io/klog/v2"
)

const (
	tcEIPFilterName = "kube-ovn-eip"
	tcEIPv4MapName  = "kube_ovn_eip_v4_map"
	tcEIPv6MapName  = "kube_ovn_eip_v6_map"
)

// 结构体内存布局必须与 tc_eip.c 中的 struct eip_target 保持完全一致（总计 32 字节，严格 64 位对齐）
type bpfEipTarget struct {
	Ifindex       uint32
	TargetMAC     [6]byte
	Pad           [2]byte
	ParentIfindex uint32
	Packets       uint64
	Bytes         uint64
}

type linuxTCManager struct {
	cfg        *Configuration
	mu         sync.Mutex
	collection *ebpf.Collection
	v4Map      *ebpf.Map
	v6Map      *ebpf.Map
	closed     bool
}

func newTCManager(cfg *Configuration) TCManager { return &linuxTCManager{cfg: cfg} }

func ensureClsactQdisc(link netlink.Link) error {
	qdisc := &netlink.GenericQdisc{QdiscAttrs: netlink.QdiscAttrs{
		LinkIndex: link.Attrs().Index, Handle: netlink.MakeHandle(0xffff, 0), Parent: netlink.HANDLE_CLSACT,
	}, QdiscType: "clsact"}
	if err := netlink.QdiscAdd(qdisc); err != nil && !errors.Is(err, unix.EEXIST) {
		return err
	}
	return nil
}

// mapPaths 根据配置的 IPv4 Pinned Map 路径推导，IPv6 Map 为其同级同名映射文件
func (m *linuxTCManager) mapPaths() (string, string) {
	v4 := m.cfg.TcBpfMapPath
	if v4 == "" {
		v4 = "/sys/fs/bpf/kube_ovn_eip_v4_map"
	}
	v6 := strings.TrimSuffix(v4, "_v4_map") + "_v6_map"
	return v4, v6
}

// Attach 具备幂等性，支持网卡设备或过滤器重建后的自动恢复。
// 失败的尝试不会留下成功挂载标记，可安全重试。
func (m *linuxTCManager) Attach(linkName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return errors.New("TC manager is closed")
	}
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}
	if linkName == "" {
		linkName = m.cfg.TcBpfInterface
	}
	if linkName == "" {
		klog.Errorf("TC BPF interface must be explicitly specified (e.g. macbond1); refusing to attach to protect production from misrouting")
		return errors.New("tc-bpf-interface must be explicitly specified")
	}
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return err
	}
	if err = ensureClsactQdisc(link); err != nil {
		return err
	}
	if m.collection == nil {
		if err = m.load(); err != nil {
			return err
		}
	}
	prog := m.collection.Programs["tc_eip_ingress"]
	filters, err := netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	if err != nil {
		return err
	}
	info, err := prog.Info()
	if err != nil {
		return err
	}
	id, _ := info.ID()
	prio := m.priority()
	for _, filter := range filters {
		attrs := filter.Attrs()
		bpf, ok := filter.(*netlink.BpfFilter)
		if attrs.Priority != prio {
			// 若发现本程序历史遗留在其他 priority 上的旧规则（如升级修改了 priority 配置），主动剥离避免旧规则拦截或双规则共存
			if ok && m.isEIPFilter(bpf) {
				if err := netlink.FilterDel(filter); err != nil && !errors.Is(err, unix.ENOENT) {
					klog.Warningf("failed to remove obsolete TC filter at priority %d on %s: %v", attrs.Priority, linkName, err)
				}
			}
			continue
		}
		if !ok || !m.isEIPFilter(bpf) {
			return fmt.Errorf("TC priority %d on %s is already owned by another filter", prio, linkName)
		}
		if bpf.Id == int(id) && bpf.DirectAction {
			return nil
		}
	}
	return netlink.FilterReplace(&netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: link.Attrs().Index, Parent: netlink.HANDLE_MIN_INGRESS,
			Handle: 1, Priority: prio, Protocol: unix.ETH_P_ALL,
		},
		Fd: prog.FD(), Name: tcEIPFilterName, DirectAction: true,
	})
}

func (m *linuxTCManager) load() error {
	path := m.cfg.TcBpfProgPath
	if path == "" {
		path = "/etc/kube-ovn/tc_eip.o"
	}
	spec, err := ebpf.LoadCollectionSpec(path)
	if err != nil {
		return fmt.Errorf("load TC object: %w", err)
	}
	if spec.Programs["tc_eip_ingress"] == nil || spec.Maps[tcEIPv4MapName] == nil || spec.Maps[tcEIPv6MapName] == nil {
		return errors.New("TC object is missing its EIP classifier or maps")
	}
	// 启动时使用全新代际的 Map。Speaker 重启后旧的 Pin 绝不能带入过期的 ifindex；
	// 周期性对账逻辑会自动重新填充当前有效的目标网关信息。
	for _, ms := range spec.Maps {
		ms.Pinning = ebpf.PinNone
	}
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("remove BPF memlock limit: %w", err)
	}
	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		return fmt.Errorf("load TC program: %w", err)
	}
	if err = m.cleanupFilters(""); err != nil {
		coll.Close()
		return err
	}
	v4, v6 := m.mapPaths()
	for name, path := range map[string]string{tcEIPv4MapName: v4, tcEIPv6MapName: v6} {
		if err = os.MkdirAll(filepath.Dir(path), 0o755); err == nil {
			err = removePin(path)
		}
		if err == nil {
			err = coll.Maps[name].Pin(path)
		}
		if err != nil {
			coll.Close()
			return fmt.Errorf("pin EIP map %s: %w", path, err)
		}
	}
	m.collection = coll
	m.v4Map, m.v6Map = coll.Maps[tcEIPv4MapName], coll.Maps[tcEIPv6MapName]
	return nil
}

const defaultTcBpfPriority uint16 = 5

func (m *linuxTCManager) priority() uint16 {
	if m != nil && m.cfg != nil && m.cfg.TcBpfPriority > 0 {
		return m.cfg.TcBpfPriority
	}
	return defaultTcBpfPriority
}

// 严格仅匹配本程序的 TC 过滤器（包括原始 tc loader 生成的名字）。
// 绝不删除与其他 CNI 共享的整个 priority 规则或 clsact qdisc 队列。
// 基于 Handle=1 与唯一 Filter 名称识别，解耦 priority 以支持跨 priority 升级迁移与统一清理。
func (m *linuxTCManager) isEIPFilter(filter *netlink.BpfFilter) bool {
	return filter != nil && filter.Handle == 1 &&
		(filter.Name == tcEIPFilterName || filter.Name == "tc_eip.o:[classifier]")
}

func (m *linuxTCManager) cleanupFilters(linkName string) error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}
	var errs []error
	for _, link := range links {
		if linkName != "" && link.Attrs().Name != linkName {
			continue
		}
		filters, err := netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, filter := range filters {
			if bpf, ok := filter.(*netlink.BpfFilter); ok && m.isEIPFilter(bpf) {
				if err := netlink.FilterDel(filter); err != nil && !errors.Is(err, unix.ENOENT) {
					errs = append(errs, err)
				}
			}
		}
	}
	return errors.Join(errs...)
}

func (m *linuxTCManager) Detach(linkName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cleanupFilters(linkName)
}

func eipMapKey(ip net.IP) ([]byte, error) {
	if v4 := ip.To4(); v4 != nil {
		return []byte(v4), nil
	}
	if v6 := ip.To16(); v6 != nil {
		return []byte(v6), nil
	}
	return nil, fmt.Errorf("invalid EIP %q", ip)
}

func validTarget(info TargetInfo) bool {
	return info.ParentIfindex >= 0 && uint64(info.ParentIfindex) <= uint64(^uint32(0)) && info.TargetIfindex > 0 && uint64(info.TargetIfindex) <= uint64(^uint32(0)) &&
		len(info.TargetMAC) == 6 && info.TargetMAC[0]&1 == 0 && !bytes.Equal(info.TargetMAC, make([]byte, 6))
}

func (m *linuxTCManager) mapFor(key []byte) (*ebpf.Map, error) {
	if m.closed {
		return nil, errors.New("TC manager is closed")
	}
	target := m.v6Map
	if len(key) == 4 {
		target = m.v4Map
	}
	if target == nil {
		return nil, errors.New("EIP BPF maps are not loaded")
	}
	return target, nil
}

func (m *linuxTCManager) addEIP(ip net.IP, info TargetInfo) error {
	if !validTarget(info) {
		return fmt.Errorf("invalid TC target for %s: ifindex=%d mac=%s", ip, info.TargetIfindex, info.TargetMAC)
	}
	key, err := eipMapKey(ip)
	if err != nil {
		return err
	}
	targetMap, err := m.mapFor(key)
	if err != nil {
		return err
	}
	// #nosec G115 -- validTarget 已校验 ifindex 均为非负数且符合 uint32 范围
	target := bpfEipTarget{Ifindex: uint32(info.TargetIfindex), ParentIfindex: uint32(info.ParentIfindex)}
	copy(target.TargetMAC[:], info.TargetMAC)
	var old bpfEipTarget
	if err = targetMap.Lookup(key, &old); err == nil {
		// 周期性全量对账过程中，若目标参数未发生变化，则不重置运行时的报文和字节统计计数器
		if old.Ifindex == target.Ifindex && old.ParentIfindex == target.ParentIfindex && old.TargetMAC == target.TargetMAC {
			return nil
		}
	} else if !errors.Is(err, ebpf.ErrKeyNotExist) {
		return err
	}
	return targetMap.Update(key, &target, ebpf.UpdateAny)
}

func (m *linuxTCManager) AddEIP(ip net.IP, info TargetInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}
	return m.addEIP(ip, info)
}

func (m *linuxTCManager) DeleteEIP(ip net.IP) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}
	key, err := eipMapKey(ip)
	if err != nil {
		return err
	}
	targetMap, err := m.mapFor(key)
	if err != nil {
		return err
	}
	if err = targetMap.Delete(key); errors.Is(err, ebpf.ErrKeyNotExist) {
		return nil
	}
	return err
}

func (m *linuxTCManager) SyncEIPs(expected map[string]TargetInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}
	wanted := make(map[string]bool, len(expected))
	var errs []error
	for address, info := range expected {
		ip := net.ParseIP(address)
		if err := m.addEIP(ip, info); err != nil {
			errs = append(errs, err)
			continue
		}
		key, err := eipMapKey(ip)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		wanted[string(key)] = true
	}
	for _, targetMap := range []*ebpf.Map{m.v4Map, m.v6Map} {
		if targetMap == nil {
			errs = append(errs, errors.New("EIP BPF maps are not loaded"))
			continue
		}
		if err := pruneEIPMap(targetMap, wanted); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func pruneEIPMap(targetMap *ebpf.Map, wanted map[string]bool) error {
	// 必须在删除前先完整收集过期 Key：直接在遍历游标处删除可能导致 BPF Hash Map 遍历重置，当保留部分条目时可能发生死循环。
	var stale [][]byte
	var key any
	for {
		next, err := targetMap.NextKeyBytes(key)
		if err != nil {
			return err
		}
		if next == nil {
			break
		}
		if !wanted[string(next)] {
			stale = append(stale, next)
		}
		key = next
	}
	var errs []error
	for _, key := range stale {
		if err := targetMap.Delete(key); err != nil && !errors.Is(err, ebpf.ErrKeyNotExist) {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func removePin(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// Close 会在功能停用时尽力而为（Best-Effort）清理遗留的 TC 挂载、BPF Collection 与 Pinned Map，并聚合返回所有错误。
func (m *linuxTCManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var errs []error
	if err := m.cleanupFilters(""); err != nil {
		errs = append(errs, err)
	}
	if m.collection != nil {
		m.collection.Close()
		m.collection = nil
		m.v4Map = nil
		m.v6Map = nil
	}
	v4, v6 := m.mapPaths()
	if err := removePin(v4); err != nil {
		errs = append(errs, err)
	}
	if err := removePin(v6); err != nil {
		errs = append(errs, err)
	}
	if len(errs) == 0 {
		m.closed = true
	}
	return errors.Join(errs...)
}
