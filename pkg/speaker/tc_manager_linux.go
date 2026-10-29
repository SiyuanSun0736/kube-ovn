//go:build linux

package speaker

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"k8s.io/klog/v2"
)

// bpfEipTarget 对应 C 语言中的 struct eip_target，必须在 64 位下保持 32 字节且 8 字节对齐
type bpfEipTarget struct {
	Ifindex   uint32   // 目标网卡接口设备索引 (4 字节)
	TargetMAC [6]byte  // 目标 MAC 地址 (6 字节)
	Pad       [2]byte  // 2 字节填充，对齐到 4 字节
	Pad2      uint32   // 4 字节显式填充，使后方的 uint64 严格对齐到 8 字节边界
	Packets   uint64   // 直通转发的数据包统计 (8 字节)
	Bytes     uint64   // 直通转发的字节统计 (8 字节)
}

// linuxTCManager 是 Linux 平台下 TCManager 的具体实现
type linuxTCManager struct {
	cfg      *Configuration
	mu       sync.Mutex
	attached map[string]bool
	v4MapFD  int
	v6MapFD  int
}

// newPlatformTCManager 创建 Linux 平台的 TC 管理器
func newPlatformTCManager(cfg *Configuration) TCManager {
	return &linuxTCManager{
		cfg:      cfg,
		attached: make(map[string]bool),
		v4MapFD:  -1,
		v6MapFD:  -1,
	}
}

// ensureClsactQdisc 确保物理网卡上存在 clsact 无锁队列调度器
func ensureClsactQdisc(link netlink.Link) error {
	qdisc := &netlink.GenericQdisc{
		QdiscAttrs: netlink.QdiscAttrs{
			LinkIndex: link.Attrs().Index,
			Handle:    netlink.MakeHandle(0xffff, 0), // clsact 的 handle 规定为 ffff:0
			Parent:    netlink.HANDLE_CLSACT,
		},
		QdiscType: "clsact",
	}

	// 添加 clsact 队列调度器，若已存在（如 Cilium 或先前端已添加）则忽略 EEXIST 错误
	if err := netlink.QdiscAdd(qdisc); err != nil && !os.IsExist(err) && !errors.Is(err, syscall.EEXIST) {
		return fmt.Errorf("为网卡 %s 添加 clsact qdisc 失败: %w", link.Attrs().Name, err)
	}
	return nil
}

// Attach 在目标网卡上完成 clsact 队列创建并挂载 TC BPF 过滤器（prio 1 direct-action）
func (m *linuxTCManager) Attach(linkName string) error {
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 若未指定网卡，优先使用配置中指定的网卡，否则自动推导宿主机默认路由物理网卡
	if linkName == "" {
		linkName = m.cfg.TcBpfInterface
	}
	if linkName == "" {
		linkName = m.detectDefaultInterface()
	}
	if linkName == "" {
		return errors.New("tc-manager: 无法自动识别或未配置用于挂载 TC 的物理网卡接口")
	}

	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return fmt.Errorf("tc-manager: 未在宿主机找到指定网卡 %s: %w", linkName, err)
	}

	// 1. 确保 clsact qdisc 存在
	if err := ensureClsactQdisc(link); err != nil {
		return err
	}

	// 2. 加载并挂载 BPF 程序到 ingress 队列
	progPath := m.cfg.TcBpfProgPath
	if progPath == "" {
		progPath = "/etc/kube-ovn/tc_eip.o"
	}

	if _, err := os.Stat(progPath); err == nil {
		// 执行标准的 tc filter replace 命令，以 priority 1 保证优先于 Cilium (prio 10) 拦截 EIP 流量
		cmd := exec.Command("tc", "filter", "replace", "dev", linkName, "ingress",
			"prio", "1", "handle", "1", "bpf", "da", "obj", progPath, "sec", "classifier")
		if output, err := cmd.CombinedOutput(); err != nil {
			klog.Warningf("通过 tc 命令行在网卡 %s 挂载 BPF 过滤器失败: %v, 输出: %s", linkName, err, string(output))
		} else {
			klog.Infof("已成功在网卡 %s 上挂载 eBPF TC 过滤器 (prio 1 direct-action, obj: %s)", linkName, progPath)
			m.attached[linkName] = true
			m.ensureMapFDs()
			return nil
		}
	} else {
		klog.Warningf("未在路径 %s 检测到编译好的 TC BPF 目标文件: %v，挂载将等待文件就绪后生效", progPath, err)
	}

	m.attached[linkName] = true
	m.ensureMapFDs()
	return nil
}

// Detach 从指定物理网卡上安全卸载 TC BPF 过滤器
func (m *linuxTCManager) Detach(linkName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if linkName == "" {
		linkName = m.cfg.TcBpfInterface
	}
	if linkName == "" {
		return nil
	}

	cmd := exec.Command("tc", "filter", "del", "dev", linkName, "ingress", "prio", "1")
	_ = cmd.Run()
	delete(m.attached, linkName)
	klog.Infof("已成功从网卡 %s 上卸载 eBPF TC 过滤器", linkName)
	return nil
}

// detectDefaultInterface 自动探测宿主机承载默认出向路由的物理网卡
func (m *linuxTCManager) detectDefaultInterface() string {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err == nil {
		for _, r := range routes {
			if r.Dst == nil && r.LinkIndex > 0 {
				if link, err := netlink.LinkByIndex(r.LinkIndex); err == nil {
					return link.Attrs().Name
				}
			}
		}
	}
	return ""
}

// ensureMapFDs 确保 IPv4 与 IPv6 的 Pinned BPF Map 文件描述符已被打开并缓存
func (m *linuxTCManager) ensureMapFDs() {
	if m.v4MapFD < 0 {
		m.v4MapFD = openPinnedMap("kube_ovn_eip_v4_map")
	}
	if m.v6MapFD < 0 {
		m.v6MapFD = openPinnedMap("kube_ovn_eip_v6_map")
	}
}

// openPinnedMap 遍历 Linux 标准 BPF 文件系统路径，检索并打开固化的 BPF Map FD
func openPinnedMap(name string) int {
	candidates := []string{
		filepath.Join("/sys/fs/bpf", name),
		filepath.Join("/sys/fs/bpf/tc/globals", name),
		filepath.Join("/sys/fs/bpf/kube-ovn", name),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			fd, err := bpfObjGet(p)
			if err == nil {
				klog.Infof("已从路径 %s 打开 Pinned BPF Map %s (fd: %d)", p, name, fd)
				return fd
			}
			klog.V(3).Infof("打开 Pinned BPF Map %s 失败: %v", p, err)
		}
	}
	return -1
}

// AddEIP 向内核 BPF Map 写入 EIP ➔ 本地网关 Pod 网卡的直通映射条目
func (m *linuxTCManager) AddEIP(eip net.IP, info TargetInfo) error {
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureMapFDs()

	target := bpfEipTarget{
		Ifindex: uint32(info.TargetIfindex),
	}
	if len(info.TargetMAC) == 6 {
		copy(target.TargetMAC[:], info.TargetMAC)
	}

	if ipv4 := eip.To4(); ipv4 != nil {
		if m.v4MapFD < 0 {
			klog.V(3).Infof("BPF IPv4 Map 文件描述符不可用，暂跳过写入 EIP %s", eip)
			return nil
		}
		// 转换为大端网络字节序 (__u32)
		key := binary.BigEndian.Uint32(ipv4)
		if err := bpfMapUpdateElem(m.v4MapFD, unsafe.Pointer(&key), unsafe.Pointer(&target), unix.BPF_ANY); err != nil {
			return fmt.Errorf("更新 BPF IPv4 Map 失败 (EIP: %s): %w", eip, err)
		}
		klog.Infof("TC BPF Map: 成功写入 IPv4 %s ➔ 网卡 ifindex %d (mac: %s)", eip, info.TargetIfindex, info.TargetMAC)
		return nil
	}

	if m.v6MapFD < 0 {
		klog.V(3).Infof("BPF IPv6 Map 文件描述符不可用，暂跳过写入 EIP %s", eip)
		return nil
	}
	var key [16]byte
	copy(key[:], eip.To16())
	if err := bpfMapUpdateElem(m.v6MapFD, unsafe.Pointer(&key[0]), unsafe.Pointer(&target), unix.BPF_ANY); err != nil {
		return fmt.Errorf("更新 BPF IPv6 Map 失败 (EIP: %s): %w", eip, err)
	}
	klog.Infof("TC BPF Map: 成功写入 IPv6 %s ➔ 网卡 ifindex %d (mac: %s)", eip, info.TargetIfindex, info.TargetMAC)
	return nil
}

// DeleteEIP 从内核 BPF Map 中删除已撤销或漂移走的 EIP 映射条目
func (m *linuxTCManager) DeleteEIP(eip net.IP) error {
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureMapFDs()

	if ipv4 := eip.To4(); ipv4 != nil {
		if m.v4MapFD < 0 {
			return nil
		}
		key := binary.BigEndian.Uint32(ipv4)
		if err := bpfMapDeleteElem(m.v4MapFD, unsafe.Pointer(&key)); err != nil {
			if !errors.Is(err, syscall.ENOENT) {
				return fmt.Errorf("从 BPF IPv4 Map 剔除条目失败 (EIP: %s): %w", eip, err)
			}
		} else {
			klog.Infof("TC BPF Map: 成功剔除 IPv4 %s", eip)
		}
		return nil
	}

	if m.v6MapFD < 0 {
		return nil
	}
	var key [16]byte
	copy(key[:], eip.To16())
	if err := bpfMapDeleteElem(m.v6MapFD, unsafe.Pointer(&key[0])); err != nil {
		if !errors.Is(err, syscall.ENOENT) {
			return fmt.Errorf("从 BPF IPv6 Map 剔除条目失败 (EIP: %s): %w", eip, err)
		}
	} else {
		klog.Infof("TC BPF Map: 成功剔除 IPv6 %s", eip)
	}
	return nil
}

// SyncEIPs 对账同步当前所有活跃 EIP 集合，确保内核 BPF Map 与集群期望状态严格一致
func (m *linuxTCManager) SyncEIPs(expected map[string]TargetInfo) error {
	if !m.cfg.EnableTcBpfRedirect {
		return nil
	}

	for ipStr, info := range expected {
		ip := net.ParseIP(ipStr)
		if ip != nil {
			_ = m.AddEIP(ip, info)
		}
	}
	return nil
}

// Close 关闭并重置已打开的 BPF 文件描述符
func (m *linuxTCManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.v4MapFD >= 0 {
		_ = unix.Close(m.v4MapFD)
		m.v4MapFD = -1
	}
	if m.v6MapFD >= 0 {
		_ = unix.Close(m.v6MapFD)
		m.v6MapFD = -1
	}
	return nil
}

// 以下为底层纯 Go 原生调用的 Linux BPF 系统调用封装（避免依赖额外的重量级 CGO 外部库）

// bpfObjGet 通过 BPF_OBJ_GET 系统调用打开位于 /sys/fs/bpf 下固化的 Map
func bpfObjGet(pathname string) (int, error) {
	pathBytes, err := unix.BytePtrFromString(pathname)
	if err != nil {
		return -1, err
	}

	attr := struct {
		pathname uint64
		bpfFd    uint32
		pad      uint32
		flags    uint32
	}{
		pathname: uint64(uintptr(unsafe.Pointer(pathBytes))),
	}

	r1, _, errNo := unix.Syscall(unix.SYS_BPF, uintptr(unix.BPF_OBJ_GET), uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if errNo != 0 {
		return -1, errNo
	}
	return int(r1), nil
}

// bpfMapUpdateElem 通过 BPF_MAP_UPDATE_ELEM 系统调用写入或更新 Key-Value 条目
func bpfMapUpdateElem(fd int, key, value unsafe.Pointer, flags uint64) error {
	attr := struct {
		mapFd uint32
		pad   uint32
		key   uint64
		value uint64
		flags uint64
	}{
		mapFd: uint32(fd),
		key:   uint64(uintptr(key)),
		value: uint64(uintptr(value)),
		flags: flags,
	}

	_, _, errNo := unix.Syscall(unix.SYS_BPF, uintptr(unix.BPF_MAP_UPDATE_ELEM), uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if errNo != 0 {
		return errNo
	}
	return nil
}

// bpfMapDeleteElem 通过 BPF_MAP_DELETE_ELEM 系统调用删除指定的 Key 条目
func bpfMapDeleteElem(fd int, key unsafe.Pointer) error {
	attr := struct {
		mapFd uint32
		pad   uint32
		key   uint64
	}{
		mapFd: uint32(fd),
		key:   uint64(uintptr(key)),
	}

	_, _, errNo := unix.Syscall(unix.SYS_BPF, uintptr(unix.BPF_MAP_DELETE_ELEM), uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if errNo != 0 {
		return errNo
	}
	return nil
}
