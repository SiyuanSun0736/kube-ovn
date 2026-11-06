//go:build linux

package speaker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"k8s.io/klog/v2"
)

// ensureClsactQdisc 确保指定网卡已附加 clsact qdisc（若已存在则忽略 EEXIST）。
func ensureClsactQdisc(link netlink.Link) error {
	qdisc := &netlink.GenericQdisc{QdiscAttrs: netlink.QdiscAttrs{
		LinkIndex: link.Attrs().Index, Handle: netlink.MakeHandle(0xffff, 0), Parent: netlink.HANDLE_CLSACT,
	}, QdiscType: "clsact"}
	if err := netlink.QdiscAdd(qdisc); err != nil && !errors.Is(err, unix.EEXIST) {
		return err
	}
	return nil
}

// tcFilterConfig 封装通用的 TC BPF Ingress 过滤器挂载参数
type tcFilterConfig struct {
	Link         netlink.Link
	ProgFD       int
	ProgID       uint32
	ProgTag      string
	FilterName   string
	Priority     uint16
	DirectAction bool
}

// attachBPFFilter 在指定网卡的 Ingress 队列上原子挂载或就地替换（Atomic Swap）eBPF TC 过滤器。
// 1. 自动确保网卡具有 clsact qdisc。
// 2. 检查现有过滤器，若已有相同 ID 与 DirectAction 的规则已就绪，直接幂等返回。
// 3. 自动探测并清理本程序历史遗留在其他 priority 上的旧规则。
// 4. 使用 netlink.FilterReplace 执行零中断无锁 RCU 热更。
func attachBPFFilter(cfg tcFilterConfig) error {
	if cfg.Link == nil {
		return errors.New("link is nil")
	}
	if err := ensureClsactQdisc(cfg.Link); err != nil {
		return err
	}
	filters, err := netlink.FilterList(cfg.Link, netlink.HANDLE_MIN_INGRESS)
	if err != nil {
		return err
	}
	linkName := cfg.Link.Attrs().Name
	for _, filter := range filters {
		attrs := filter.Attrs()
		bpf, ok := filter.(*netlink.BpfFilter)
		if attrs.Priority != cfg.Priority {
			// 若发现本程序历史遗留在其他 priority 上的旧规则（如升级修改了 priority 配置），主动剥离避免旧规则拦截或双规则共存
			if ok && isMatchingTCFilter(bpf, cfg.ProgID, cfg.ProgTag, cfg.FilterName) {
				if err := netlink.FilterDel(filter); err != nil && !errors.Is(err, unix.ENOENT) {
					klog.Warningf("failed to remove obsolete TC filter at priority %d on %s: %v", attrs.Priority, linkName, err)
				}
			}
			continue
		}
		if !ok || !isMatchingTCFilter(bpf, cfg.ProgID, cfg.ProgTag, cfg.FilterName) {
			return fmt.Errorf("TC priority %d on %s is already owned by another filter", cfg.Priority, linkName)
		}
		if bpf.Id == int(cfg.ProgID) && bpf.DirectAction == cfg.DirectAction {
			return nil
		}
	}
	return netlink.FilterReplace(&netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: cfg.Link.Attrs().Index, Parent: netlink.HANDLE_MIN_INGRESS,
			Handle: 1, Priority: cfg.Priority, Protocol: unix.ETH_P_ALL,
		},
		Fd: cfg.ProgFD, Name: cfg.FilterName, DirectAction: cfg.DirectAction,
	})
}

// isMatchingTCFilter 基于 Ingress parent、Handle=1、Protocol=ALL、DirectAction
// 以及 BPF Program ID / BPF 指令 Tag / 专属名称等多维特征进行精准交叉核验，
// 绝不误删与其他 CNI (如 Cilium/Calico/QoS) 共享的规则。
func isMatchingTCFilter(filter *netlink.BpfFilter, progID uint32, progTag string, filterNames ...string) bool {
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
	if progID > 0 && filter.Id == int(progID) {
		return true
	}
	// 强指纹 2：若 BPF 字节码指令 SHA Tag 匹配当前编译产物
	if progTag != "" && filter.Tag != "" && strings.EqualFold(filter.Tag, progTag) {
		return true
	}
	// 强指纹 3：专属 Filter 唯一名称标识
	for _, name := range filterNames {
		if name != "" && filter.Name == name {
			return true
		}
	}
	return false
}

// detachBPFFiltersOnLinks 安全遍历指定网卡（为空则遍历所有网卡）并按 matchFunc 卸载符合条件的 TC Filter。
func detachBPFFiltersOnLinks(linkName string, matchFunc func(*netlink.BpfFilter) bool) error {
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
			if bpf, ok := filter.(*netlink.BpfFilter); ok && matchFunc(bpf) {
				if err := netlink.FilterDel(filter); err != nil && !errors.Is(err, unix.ENOENT) {
					errs = append(errs, err)
				}
			}
		}
	}
	return errors.Join(errs...)
}

// pruneBPFHashMap 安全清理 BPF Hash Map 中的过期条目。
// 必须在删除前先完整收集过期 Key：直接在遍历游标处删除可能导致 BPF Hash Map 遍历重置，当保留部分条目时可能发生死循环。
func pruneBPFHashMap(targetMap *ebpf.Map, wanted map[string]bool) error {
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

// removeBPFPin 安全清理 Pinned Map 文件。
func removeBPFPin(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// loadCollectionWithPinnedMaps 加载 BPF CollectionSpec，并优先复用 bpffs 中既有的 Pinned Map，
// 实现进程重启/热更新期间数据面 0 丢包；若复用失败则自动回退并清理旧 Map。
func loadCollectionWithPinnedMaps(spec *ebpf.CollectionSpec, pinMapPaths map[string]string) (*ebpf.Collection, error) {
	for _, ms := range spec.Maps {
		ms.Pinning = ebpf.PinNone
	}
	if err := rlimit.RemoveMemlock(); err != nil {
		return nil, fmt.Errorf("remove BPF memlock limit: %w", err)
	}

	replacements := make(map[string]*ebpf.Map)
	for name, p := range pinMapPaths {
		if pinned, err := ebpf.LoadPinnedMap(p, nil); err == nil {
			replacements[name] = pinned
		}
	}

	coll, err := ebpf.NewCollectionWithOptions(spec, ebpf.CollectionOptions{MapReplacements: replacements})
	if err != nil && len(replacements) > 0 {
		klog.Warningf("failed to load BPF collection with existing pinned maps, falling back to new maps: %v", err)
		for _, r := range replacements {
			_ = r.Close()
		}
		replacements = make(map[string]*ebpf.Map)
		for _, p := range pinMapPaths {
			_ = removeBPFPin(p)
		}
		coll, err = ebpf.NewCollection(spec)
	}
	if err != nil {
		for _, r := range replacements {
			_ = r.Close()
		}
		return nil, fmt.Errorf("load BPF collection: %w", err)
	}

	// 仅对新创建（未复用）的 Map 挂载至 bpffs
	for name, path := range pinMapPaths {
		if _, ok := replacements[name]; ok {
			continue
		}
		if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			coll.Close()
			return nil, fmt.Errorf("create map dir %s: %w", filepath.Dir(path), err)
		}
		if err = coll.Maps[name].Pin(path); err != nil {
			coll.Close()
			return nil, fmt.Errorf("pin BPF map %s: %w", path, err)
		}
	}
	return coll, nil
}
