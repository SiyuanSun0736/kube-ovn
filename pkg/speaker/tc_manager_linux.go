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
	tcEIPFilterName       = "kube-ovn-eip"
	tcEIPv4MapName        = "kube_ovn_eip_v4_map"
	tcEIPv6MapName        = "kube_ovn_eip_v6_map"
	defaultTcEIPv4MapPath = "/sys/fs/bpf/kube_ovn_eip_v4_map"
	defaultTcEIPv6MapPath = "/sys/fs/bpf/kube_ovn_eip_v6_map"
)

// 结构体内存布局必须与 tc_eip.c 中的 struct eip_target 保持完全一致（总计 16 字节，严格 32 位对齐）
type bpfEipTarget struct {
	Ifindex       uint32
	TargetMAC     [6]byte
	Pad           [2]byte
	ParentIfindex uint32
}

type linuxTCManager struct {
	cfg         *Configuration
	mu          sync.Mutex
	collection  *ebpf.Collection
	v4Map       *ebpf.Map
	v6Map       *ebpf.Map
	progID      uint32
	progTag     string
	attachedDev string
	closed      bool
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

// mapPaths 返回 IPv4 和 IPv6 Pinned Map 路径。
func (m *linuxTCManager) mapPaths() (string, string) {
	if m.cfg.TcBpfMapPath == "" || m.cfg.TcBpfMapPath == defaultTcEIPv4MapPath {
		return defaultTcEIPv4MapPath, defaultTcEIPv6MapPath
	}
	v4 := m.cfg.TcBpfMapPath
	v6 := strings.TrimSuffix(v4, "_v4_map") + "_v6_map"
	return v4, v6
}

// ensureAttached 自动确保指定的物理网卡已挂载 TC BPF 程序，并通过内存缓存短路避免高频 Netlink 系统调用。
func (m *linuxTCManager) ensureAttached(linkName string) error {
	if m.closed {
		return errors.New("TC manager is closed")
	}
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}
	// 若处于测试注入独立 Map 模式（未加载完整 BPF Collection），跳过底层网卡挂载
	if m.collection == nil && (m.v4Map != nil || m.v6Map != nil) {
		return nil
	}
	if linkName == "" {
		linkName = m.cfg.TcBpfInterface
	}
	if linkName == "" {
		klog.Errorf("TC BPF interface must be explicitly specified (e.g. macbond1); refusing to attach to protect production from misrouting")
		return errors.New("tc-bpf-interface must be explicitly specified")
	}
	if m.attachedDev == linkName {
		return nil
	}
	if err := m.attach(linkName); err != nil {
		return err
	}
	m.attachedDev = linkName
	return nil
}

func (m *linuxTCManager) attach(linkName string) error {
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
	// 使用 FilterReplace 在指定网卡的 Ingress 队列上原子挂载或就地替换（Atomic Swap）eBPF TC 过滤器：
	// 1. LinkIndex/Parent: 挂载到宿主机物理网卡的 clsact Ingress 入向通道，在报文进入内核 Netfilter/conntrack 之前提前拦截；
	// 2. Handle/Priority/Protocol: 独占配置的优先级（默认 5），匹配所有二层以太网协议（ETH_P_ALL，含 IPv4/IPv6/VLAN）；
	// 3. DirectAction: 开启 direct-action (da) 快速通道，允许 BPF 程序直接执行 bpf_redirect 进行微秒级线速转发，跳过传统 TC Action 管道；
	// 4. 原子切替：内核依靠 RCU 机制实现热更无锁平滑升级，过程全程零中断、不漏包。
	return netlink.FilterReplace(&netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: link.Attrs().Index, Parent: netlink.HANDLE_MIN_INGRESS,
			Handle: 1, Priority: prio, Protocol: unix.ETH_P_ALL,
		},
		Fd: prog.FD(), Name: tcEIPFilterName, DirectAction: true,
	})
}

// Attach 具备幂等性，支持网卡设备或过滤器重建后的自动恢复。
// 失败的尝试不会留下成功挂载标记，可安全重试。
func (m *linuxTCManager) Attach(linkName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ensureAttached(linkName)
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
	// 优先尝试复用 bpffs 中既有的 Pinned Map，实现进程重启/热更新期间数据面 0 丢包。
	for _, ms := range spec.Maps {
		ms.Pinning = ebpf.PinNone
	}
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("remove BPF memlock limit: %w", err)
	}

	v4Path, v6Path := m.mapPaths()
	mapPaths := map[string]string{tcEIPv4MapName: v4Path, tcEIPv6MapName: v6Path}
	replacements := make(map[string]*ebpf.Map)
	for name, p := range mapPaths {
		if pinned, err := ebpf.LoadPinnedMap(p, nil); err == nil {
			replacements[name] = pinned
		}
	}

	coll, err := ebpf.NewCollectionWithOptions(spec, ebpf.CollectionOptions{MapReplacements: replacements})
	if err != nil && len(replacements) > 0 {
		klog.Warningf("failed to load TC program with existing pinned maps, falling back to new maps: %v", err)
		for _, r := range replacements {
			_ = r.Close()
		}
		replacements = make(map[string]*ebpf.Map)
		for _, p := range mapPaths {
			_ = removePin(p)
		}
		coll, err = ebpf.NewCollection(spec)
	}
	if err != nil {
		for _, r := range replacements {
			_ = r.Close()
		}
		return fmt.Errorf("load TC program: %w", err)
	}

	// 仅对新创建（未复用）的 Map 挂载至 bpffs
	for name, path := range mapPaths {
		if _, ok := replacements[name]; ok {
			continue
		}
		if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			coll.Close()
			return fmt.Errorf("create map dir %s: %w", filepath.Dir(path), err)
		}
		if err = coll.Maps[name].Pin(path); err != nil {
			coll.Close()
			return fmt.Errorf("pin EIP map %s: %w", path, err)
		}
	}
	m.collection = coll
	m.v4Map, m.v6Map = coll.Maps[tcEIPv4MapName], coll.Maps[tcEIPv6MapName]
	prog := coll.Programs["tc_eip_ingress"]
	if info, err := prog.Info(); err == nil {
		if id, ok := info.ID(); ok {
			m.progID = uint32(id)
		}
		m.progTag = info.Tag
	}
	return nil
}

const defaultTcBpfPriority uint16 = 5

func (m *linuxTCManager) priority() uint16 {
	if m != nil && m.cfg != nil && m.cfg.TcBpfPriority > 0 {
		return m.cfg.TcBpfPriority
	}
	return defaultTcBpfPriority
}

// 严格仅匹配本程序的 TC 过滤器。
// 基于 Filter 类型、Ingress parent、Handle=1、Protocol=ALL、DirectAction 快速转发通道
// 以及 BPF Program ID / BPF 指令 Tag / 专属名称等多维特征进行精准交叉核验，
// 绝不误删与其他 CNI (如 Cilium/Calico/QoS) 共享的规则。
func (m *linuxTCManager) isEIPFilter(filter *netlink.BpfFilter) bool {
	if filter == nil {
		return false
	}
	attrs := filter.Attrs()
	if attrs == nil || attrs.Parent != netlink.HANDLE_MIN_INGRESS || attrs.Handle != 1 || attrs.Protocol != unix.ETH_P_ALL {
		return false
	}
	if !filter.DirectAction {
		return false
	}
	// 强指纹 1：若当前加载了 BPF 实例，且程序运行时内核 ID 相同
	if m.progID > 0 && filter.Id == int(m.progID) {
		return true
	}
	// 强指纹 2：若 BPF 字节码指令 SHA Tag 匹配当前编译产物
	if m.progTag != "" && filter.Tag != "" && strings.EqualFold(filter.Tag, m.progTag) {
		return true
	}
	// 强指纹 3：专属 Filter 唯一名称标识
	return filter.Name == tcEIPFilterName || filter.Name == "tc_eip.o:[classifier]"
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
	m.attachedDev = ""
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

func validTCRedirectTarget(target TCRedirectTarget) bool {
	return target.ParentIfindex >= 0 && uint64(target.ParentIfindex) <= uint64(^uint32(0)) && target.TargetIfindex > 0 && uint64(target.TargetIfindex) <= uint64(^uint32(0)) &&
		len(target.TargetMAC) == 6 && target.TargetMAC[0]&1 == 0 && !bytes.Equal(target.TargetMAC, make([]byte, 6))
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

func (m *linuxTCManager) addEIP(ip net.IP, target TCRedirectTarget) error {
	if !validTCRedirectTarget(target) {
		return fmt.Errorf("invalid TC target for %s: ifindex=%d mac=%s", ip, target.TargetIfindex, target.TargetMAC)
	}
	key, err := eipMapKey(ip)
	if err != nil {
		return err
	}
	targetMap, err := m.mapFor(key)
	if err != nil {
		return err
	}
	// #nosec G115 -- validTCRedirectTarget 已校验 ifindex 均为非负数且符合 uint32 范围
	bpfTarget := bpfEipTarget{Ifindex: uint32(target.TargetIfindex), ParentIfindex: uint32(target.ParentIfindex)}
	copy(bpfTarget.TargetMAC[:], target.TargetMAC)
	var old bpfEipTarget
	if err = targetMap.Lookup(key, &old); err == nil {
		if old == bpfTarget {
			return nil
		}
	} else if !errors.Is(err, ebpf.ErrKeyNotExist) {
		return err
	}
	return targetMap.Update(key, &bpfTarget, ebpf.UpdateAny)
}

func (m *linuxTCManager) AddEIP(ip net.IP, target TCRedirectTarget) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}
	if err := m.ensureAttached(m.cfg.TcBpfInterface); err != nil {
		return err
	}
	return m.addEIP(ip, target)
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

func (m *linuxTCManager) SyncEIPs(expected map[string]TCRedirectTarget, retainedIPs ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}
	if len(expected) > 0 || len(retainedIPs) > 0 {
		if err := m.ensureAttached(m.cfg.TcBpfInterface); err != nil {
			return err
		}
	}
	wanted := make(map[string]bool, len(expected)+len(retainedIPs))
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

	// 声明式保留逻辑：处理属于当前节点活跃范围、仅因 Pod 元数据短暂就绪延迟而需保留存量映射的 EIP
	for _, address := range retainedIPs {
		ip := net.ParseIP(address)
		if ip == nil {
			continue
		}
		key, err := eipMapKey(ip)
		if err != nil {
			continue
		}
		if wanted[string(key)] {
			continue
		}
		targetMap, err := m.mapFor(key)
		if err != nil {
			continue
		}
		var old bpfEipTarget
		if err := targetMap.Lookup(key, &old); err == nil && old.Ifindex > 0 {
			// 严格核验宿主机上存量目标网卡设备的物理存活性，仅在设备存活时保留旧条目
			if link, err := netlink.LinkByIndex(int(old.Ifindex)); err == nil && link != nil {
				wanted[string(key)] = true
			}
		}
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

// Close 关闭并释放已打开的用户态 BPF Collection 与文件描述符句柄，
// 但保留内核中的 TC Filter 规则与 Pinned Map 数据面，支持进程重启与热更新期间数据面 0 丢包。
func (m *linuxTCManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attachedDev = ""
	if m.collection != nil {
		m.collection.Close()
		m.collection = nil
		m.v4Map = nil
		m.v6Map = nil
		m.progID = 0
		m.progTag = ""
	}
	m.closed = true
	return nil
}

// Clean 会在功能显式停用时尽力而为（Best-Effort）彻底清理遗留的 TC 挂载、BPF Collection 与 Pinned Map，并聚合返回所有错误。
func (m *linuxTCManager) Clean() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attachedDev = ""
	var errs []error
	if err := m.cleanupFilters(""); err != nil {
		errs = append(errs, err)
	}
	if m.collection != nil {
		m.collection.Close()
		m.collection = nil
		m.v4Map = nil
		m.v6Map = nil
		m.progID = 0
		m.progTag = ""
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
