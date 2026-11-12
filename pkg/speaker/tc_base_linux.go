//go:build linux

package speaker

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"k8s.io/klog/v2"
)

type tcNetlinkOps interface {
	QdiscAdd(netlink.Qdisc) error
	FilterList(netlink.Link, uint32) ([]netlink.Filter, error)
	FilterReplace(netlink.Filter) error
	FilterDel(netlink.Filter) error
	LinkByName(string) (netlink.Link, error)
	LinkByIndex(int) (netlink.Link, error)
	LinkList() ([]netlink.Link, error)
}

type defaultTCNetlinkOps struct{}

func (defaultTCNetlinkOps) QdiscAdd(q netlink.Qdisc) error {
	return netlink.QdiscAdd(q)
}

func (defaultTCNetlinkOps) FilterList(link netlink.Link, parent uint32) ([]netlink.Filter, error) {
	return netlink.FilterList(link, parent)
}

func (defaultTCNetlinkOps) FilterReplace(f netlink.Filter) error {
	return netlink.FilterReplace(f)
}

func (defaultTCNetlinkOps) FilterDel(f netlink.Filter) error {
	return netlink.FilterDel(f)
}

func (defaultTCNetlinkOps) LinkByName(name string) (netlink.Link, error) {
	return netlink.LinkByName(name)
}

func (defaultTCNetlinkOps) LinkByIndex(index int) (netlink.Link, error) {
	return netlink.LinkByIndex(index)
}

func (defaultTCNetlinkOps) LinkList() ([]netlink.Link, error) {
	return netlink.LinkList()
}

var netlinkOps tcNetlinkOps = defaultTCNetlinkOps{}

// ensureClsactQdisc 确保指定网卡已附加 clsact qdisc（若已存在则忽略 EEXIST）。
func ensureClsactQdisc(link netlink.Link) error {
	qdisc := &netlink.GenericQdisc{QdiscAttrs: netlink.QdiscAttrs{
		LinkIndex: link.Attrs().Index, Handle: netlink.MakeHandle(0xffff, 0), Parent: netlink.HANDLE_CLSACT,
	}, QdiscType: "clsact"}
	if err := netlinkOps.QdiscAdd(qdisc); err != nil && !errors.Is(err, unix.EEXIST) {
		return err
	}
	return nil
}

// tcFilterConfig 封装通用的 TC BPF Ingress 过滤器挂载参数。
type tcFilterConfig struct {
	Link        netlink.Link
	ProgFD      int
	ProgID      uint32
	ProgTag     string
	AttachName  string
	LegacyNames []string
	Priority    uint16
}

// attachBPFFilter 在指定网卡的 Ingress 队列上原子挂载或就地替换（Atomic Swap）eBPF TC 过滤器。
func attachBPFFilter(cfg tcFilterConfig) error {
	if cfg.Link == nil {
		return errors.New("link is nil")
	}
	if cfg.Link.Attrs() == nil {
		return errors.New("link attrs is nil")
	}
	if cfg.ProgFD <= 0 {
		return errors.New("bpf program fd is invalid")
	}
	if cfg.AttachName == "" {
		return errors.New("AttachName must not be empty")
	}
	if err := ensureClsactQdisc(cfg.Link); err != nil {
		return err
	}
	filters, err := netlinkOps.FilterList(cfg.Link, netlink.HANDLE_MIN_INGRESS)
	if err != nil {
		return err
	}
	allNames := append([]string{cfg.AttachName}, cfg.LegacyNames...)
	linkName := cfg.Link.Attrs().Name
	var matchedCurrent bool
	var obsolete []netlink.Filter
	for _, filter := range filters {
		attrs := filter.Attrs()
		if attrs == nil {
			continue
		}
		bpf, ok := filter.(*netlink.BpfFilter)
		if attrs.Priority != cfg.Priority {
			// 先挂后删（Make-Before-Break）：先收集旧规则，挂载新规则成功后再清理
			if ok && isMatchingTCFilter(bpf, cfg.ProgID, cfg.ProgTag, allNames...) {
				obsolete = append(obsolete, filter)
			}
			continue
		}
		if !ok || !isMatchingTCFilter(bpf, cfg.ProgID, cfg.ProgTag, allNames...) {
			return fmt.Errorf("TC priority %d on %s is already owned by another filter", cfg.Priority, linkName)
		}
		// 严禁在此直接 return nil：否则当同 priority 规则排在前面时，排在后面的异 priority 废弃规则会被跳过清理。
		if isCurrentTCFilter(bpf, cfg.ProgID, cfg.AttachName) {
			matchedCurrent = true
		}
	}
	if !matchedCurrent {
		if err := netlinkOps.FilterReplace(&netlink.BpfFilter{
			FilterAttrs: netlink.FilterAttrs{
				LinkIndex: cfg.Link.Attrs().Index, Parent: netlink.HANDLE_MIN_INGRESS,
				Handle: 1, Priority: cfg.Priority, Protocol: unix.ETH_P_ALL,
			},
			Fd: cfg.ProgFD, Name: cfg.AttachName, DirectAction: true,
		}); err != nil {
			return err
		}
	}
	for _, filter := range obsolete {
		if err := netlinkOps.FilterDel(filter); err != nil && !errors.Is(err, unix.ENOENT) && !errors.Is(err, unix.ENODEV) && !errors.Is(err, unix.ESRCH) {
			klog.Warningf("failed to remove obsolete TC filter at priority %d on %s: %v", filter.Attrs().Priority, linkName, err)
		}
	}
	return nil
}

// isCurrentTCFilter 校验规则是否已是当前 Program 在目标 Priority 上的正式生效规则。
func isCurrentTCFilter(bpf *netlink.BpfFilter, progID uint32, attachName string) bool {
	return bpf != nil && bpf.Name == attachName && progID > 0 && bpf.Id == int(progID) && bpf.DirectAction
}

// isMatchingTCFilter 基于 Ingress parent、Handle=1、Protocol=ALL、DirectAction
// 以及 BPF Program ID / BPF 指令 Tag / 专属名称等多维特征进行精准交叉核验。
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
	if progID > 0 && filter.Id == int(progID) {
		return true
	}
	if progTag != "" && filter.Tag != "" && strings.EqualFold(filter.Tag, progTag) {
		return true
	}
	for _, name := range filterNames {
		if name != "" && filter.Name == name {
			return true
		}
	}
	return false
}

// detachBPFFiltersOnLinks 安全遍历指定网卡（为空则遍历所有网卡）并按 matchFunc 卸载符合条件的 TC Filter。
func detachBPFFiltersOnLinks(linkName string, matchFunc func(*netlink.BpfFilter) bool) error {
	if matchFunc == nil {
		return errors.New("matchFunc must not be nil")
	}
	var links []netlink.Link
	if linkName != "" {
		link, err := netlinkOps.LinkByName(linkName)
		if err != nil {
			var linkNotFound netlink.LinkNotFoundError
			if errors.As(err, &linkNotFound) || errors.Is(err, unix.ENODEV) || errors.Is(err, unix.ENOENT) {
				return nil
			}
			return err
		}
		links = []netlink.Link{link}
	} else {
		var err error
		links, err = netlinkOps.LinkList()
		if err != nil {
			return err
		}
	}
	var errs []error
	for _, link := range links {
		if link == nil || link.Attrs() == nil {
			continue
		}
		filters, err := netlinkOps.FilterList(link, netlink.HANDLE_MIN_INGRESS)
		if err != nil {
			if errors.Is(err, unix.ENOENT) || errors.Is(err, unix.ENODEV) {
				klog.V(4).Infof("interface %s has no clsact qdisc or is gone, skipping filter detach: %v", link.Attrs().Name, err)
				continue
			}
			errs = append(errs, err)
			continue
		}
		for _, filter := range filters {
			if bpf, ok := filter.(*netlink.BpfFilter); ok && matchFunc(bpf) {
				if err := netlinkOps.FilterDel(filter); err != nil && !errors.Is(err, unix.ENOENT) && !errors.Is(err, unix.ENODEV) && !errors.Is(err, unix.ESRCH) {
					errs = append(errs, err)
				}
			}
		}
	}
	return errors.Join(errs...)
}
