//go:build linux

package speaker

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/cilium/ebpf"
	"github.com/vishvananda/netlink"
	"k8s.io/klog/v2"
)

const (
	tcEIPFilterName = "kube-ovn-eip"
	tcEIPv4MapName  = "kube_ovn_eip_v4_map"
	tcEIPv6MapName  = "kube_ovn_eip_v6_map"
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

// mapPaths 根据配置的 IPv4 Pinned Map 路径推导，IPv6 Map 为其同级同名映射文件
func (m *linuxTCManager) mapPaths() (string, string) {
	v4 := m.cfg.TcBpfMapPath
	if v4 == "" {
		v4 = "/sys/fs/bpf/kube_ovn_eip_v4_map"
	}
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
		linkName = m.attachedDev
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
	if m.collection == nil {
		if err = m.load(); err != nil {
			return err
		}
	}
	prog := m.collection.Programs["tc_eip_ingress"]
	info, err := prog.Info()
	if err != nil {
		return err
	}
	id, _ := info.ID()
	return attachBPFFilter(tcFilterConfig{
		Link:         link,
		ProgFD:       prog.FD(),
		ProgID:       uint32(id),
		ProgTag:      info.Tag,
		FilterName:   tcEIPFilterName,
		Priority:     m.priority(),
		DirectAction: true,
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

	v4Path, v6Path := m.mapPaths()
	mapPaths := map[string]string{tcEIPv4MapName: v4Path, tcEIPv6MapName: v6Path}
	coll, err := loadCollectionWithPinnedMaps(spec, mapPaths)
	if err != nil {
		return fmt.Errorf("load TC program: %w", err)
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
func (m *linuxTCManager) isEIPFilter(filter *netlink.BpfFilter) bool {
	return isMatchingTCFilter(filter, m.progID, m.progTag, tcEIPFilterName, "tc_eip.o:[classifier]")
}

func (m *linuxTCManager) cleanupFilters(linkName string) error {
	return detachBPFFiltersOnLinks(linkName, m.isEIPFilter)
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
		if err := pruneBPFHashMap(targetMap, wanted); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
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
	if err := removeBPFPin(v4); err != nil {
		errs = append(errs, err)
	}
	if err := removeBPFPin(v6); err != nil {
		errs = append(errs, err)
	}
	if len(errs) == 0 {
		m.closed = true
	}
	return errors.Join(errs...)
}
