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
)

const (
	tcEIPFilterName = "kube-ovn-eip"
	tcEIPv4MapName  = "kube_ovn_eip_v4_map"
	tcEIPv6MapName  = "kube_ovn_eip_v6_map"
)

// Keep this layout in sync with struct eip_target in tc_eip.c.
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

func newPlatformTCManager(cfg *Configuration) TCManager { return &linuxTCManager{cfg: cfg} }

func ensureClsactQdisc(link netlink.Link) error {
	qdisc := &netlink.GenericQdisc{QdiscAttrs: netlink.QdiscAttrs{
		LinkIndex: link.Attrs().Index, Handle: netlink.MakeHandle(0xffff, 0), Parent: netlink.HANDLE_CLSACT,
	}, QdiscType: "clsact"}
	if err := netlink.QdiscAdd(qdisc); err != nil && !errors.Is(err, unix.EEXIST) {
		return err
	}
	return nil
}

// mapPaths honors the configured IPv4 pin path; the IPv6 pin is its sibling.
func (m *linuxTCManager) mapPaths() (string, string) {
	v4 := m.cfg.TcBpfMapPath
	if v4 == "" {
		v4 = "/sys/fs/bpf/kube_ovn_eip_v4_map"
	}
	v6 := strings.TrimSuffix(v4, "_v4_map") + "_v6_map"
	return v4, v6
}

// Attach is idempotent, including recovery after a link or filter was recreated.
// A failed attempt leaves no successful attachment marker and can be retried.
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
		var err error
		linkName, err = detectDefaultInterface()
		if err != nil {
			return err
		}
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
	for _, filter := range filters {
		attrs := filter.Attrs()
		if attrs.Priority != 1 {
			continue
		}
		bpf, ok := filter.(*netlink.BpfFilter)
		if !ok || !isEIPFilter(bpf) {
			return fmt.Errorf("TC priority 1 on %s is already owned by another filter", linkName)
		}
		if bpf.Id == int(id) && bpf.DirectAction {
			return nil
		}
	}
	return netlink.FilterReplace(&netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: link.Attrs().Index, Parent: netlink.HANDLE_MIN_INGRESS,
			Handle: 1, Priority: 1, Protocol: unix.ETH_P_ALL,
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
	// Start from an empty generation. Old pins must not resurrect stale ifindices
	// after a speaker restart; periodic reconciliation repopulates current targets.
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

func detectDefaultInterface() (string, error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		return "", err
	}
	for _, r := range routes {
		if r.Dst == nil && r.LinkIndex > 0 {
			link, err := netlink.LinkByIndex(r.LinkIndex)
			if err != nil {
				return "", err
			}
			return link.Attrs().Name, nil
		}
	}
	return "", errors.New("no default IPv4 interface found; specify --tc-bpf-interface")
}

// Match only our filter, including the name emitted by the original tc loader.
// Never delete an entire priority or clsact qdisc shared with another CNI.
func isEIPFilter(filter *netlink.BpfFilter) bool {
	return filter.Handle == 1 && filter.Priority == 1 &&
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
			if bpf, ok := filter.(*netlink.BpfFilter); ok && isEIPFilter(bpf) {
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
	return info.ParentIfindex > 0 && uint64(info.ParentIfindex) <= uint64(^uint32(0)) && info.TargetIfindex > 0 && uint64(info.TargetIfindex) <= uint64(^uint32(0)) &&
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
	// #nosec G115 -- validTarget checks both indices are positive and fit uint32.
	target := bpfEipTarget{Ifindex: uint32(info.TargetIfindex), ParentIfindex: uint32(info.ParentIfindex)}
	copy(target.TargetMAC[:], info.TargetMAC)
	var old bpfEipTarget
	if err = targetMap.Lookup(key, &old); err == nil {
		// Do not reset live packet counters during periodic reconciliation.
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
	// Collect keys before deleting: deleting the iteration cursor restarts a BPF
	// hash-map walk, potentially looping forever when some entries are retained.
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

// Close also cleans up a previous process's attachments when the feature is
// disabled. Closing the userspace FD alone never detaches a TC program.
func (m *linuxTCManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	if err := m.cleanupFilters(""); err != nil {
		return err
	}
	if m.collection != nil {
		m.collection.Close()
		m.collection = nil
		m.v4Map = nil
		m.v6Map = nil
	}
	v4, v6 := m.mapPaths()
	return errors.Join(removePin(v4), removePin(v6))
}
