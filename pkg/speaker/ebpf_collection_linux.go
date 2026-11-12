//go:build linux

package speaker

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"k8s.io/klog/v2"
)

// pruneBPFHashMap 安全清理 BPF Hash Map 中的过期条目。
// 必须在删除前先完整收集过期 Key：直接在遍历游标处删除可能导致 BPF Hash Map 遍历重置，当保留部分条目时可能发生死循环。
func pruneBPFHashMap(targetMap *ebpf.Map, wanted map[string]bool) error {
	if targetMap == nil {
		return errors.New("target map is nil")
	}
	var stale [][]byte
	var key any
	for {
		next, err := targetMap.NextKeyBytes(key)
		if err != nil {
			if errors.Is(err, ebpf.ErrKeyNotExist) {
				break
			}
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
	for _, key := range stale {
		if err := targetMap.Delete(key); err != nil && !errors.Is(err, ebpf.ErrKeyNotExist) {
			keyDesc := hex.EncodeToString(key)
			if len(key) == net.IPv4len || len(key) == net.IPv6len {
				keyDesc = net.IP(key).String()
			}
			klog.Errorf("failed to delete stale entry %s from BPF map %s: %v", keyDesc, targetMap, err)
		}
	}
	return nil
}

type bpfMapPinOps interface {
	MkdirAll(path string, perm os.FileMode) error
	Pin(m *ebpf.Map, path string) error
	Rename(oldpath, newpath string) error
	Remove(path string) error
}

type defaultBPFMapPinOps struct{}

func (defaultBPFMapPinOps) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (defaultBPFMapPinOps) Pin(m *ebpf.Map, path string) error {
	return m.Pin(path)
}

func (defaultBPFMapPinOps) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (defaultBPFMapPinOps) Remove(path string) error {
	return os.Remove(path)
}

var mapPinOps bpfMapPinOps = defaultBPFMapPinOps{}

// removeBPFPin 安全清理 Pinned Map 文件。
func removeBPFPin(path string) error {
	if err := mapPinOps.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// loadCollectionWithPinnedMaps 加载 BPF CollectionSpec，并优先复用 bpffs 中既有的 Pinned Map，
// 实现进程重启/热更新期间数据面 0 丢包；若复用失败则安全回退，
// 仅在新 Collection 创建成功后原子替换 Pin 路径，确保旧实例回滚能力不受损。
func loadCollectionWithPinnedMaps(spec *ebpf.CollectionSpec, pinMapPaths map[string]string) (*ebpf.Collection, error) {
	for _, ms := range spec.Maps {
		ms.Pinning = ebpf.PinNone
	}
	if err := rlimit.RemoveMemlock(); err != nil {
		return nil, fmt.Errorf("remove BPF memlock limit: %w", err)
	}

	replacements := make(map[string]*ebpf.Map)
	reused := make(map[string]bool)
	// 根据 cilium/ebpf 契约：CollectionOptions.MapReplacements 内部仅克隆 (Clone) 传入的 Map 句柄，
	// 调用方仍然拥有 LoadPinnedMap 打开的原始句柄的所有权，必须显式 Close()，否则每次加载均会导致 FD 泄漏。
	defer func() {
		for _, r := range replacements {
			_ = r.Close()
		}
	}()

	for name, p := range pinMapPaths {
		if pinned, err := ebpf.LoadPinnedMap(p, nil); err == nil {
			replacements[name] = pinned
			reused[name] = true
		}
	}

	coll, err := ebpf.NewCollectionWithOptions(spec, ebpf.CollectionOptions{MapReplacements: replacements})
	if err != nil && len(replacements) > 0 {
		klog.Warningf("failed to load BPF collection with existing pinned maps, falling back to new maps: %v", err)
		// 绝不在这里提前 removeBPFPin(p)，避免新 Collection 加载失败退出时旧 Pin 被破坏导致无法回滚；
		// 清空 reused 标记，确保回退创建的全新 Map 实例在后续两阶段提交中被完整 Pin 到 bpffs。
		reused = make(map[string]bool)
		coll, err = ebpf.NewCollection(spec)
	}
	if err != nil {
		return nil, fmt.Errorf("load BPF collection: %w", err)
	}

	type pendingPin struct {
		name    string
		tmpPath string
		dstPath string
	}
	var pending []pendingPin
	cleanupTmp := func() {
		for _, p := range pending {
			_ = removeBPFPin(p.tmpPath)
		}
	}

	names := make([]string, 0, len(pinMapPaths))
	for name := range pinMapPaths {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		path := pinMapPaths[name]
		if reused[name] {
			continue
		}
		if err = mapPinOps.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			cleanupTmp()
			coll.Close()
			return nil, fmt.Errorf("create map dir %s: %w", filepath.Dir(path), err)
		}
		targetMap, ok := coll.Maps[name]
		if !ok || targetMap == nil {
			cleanupTmp()
			coll.Close()
			return nil, fmt.Errorf("BPF map %s not found in collection", name)
		}
		// bpffs 保留含点的文件名，使用不含点的后缀避免 Pin 返回 EPERM。
		tmpPath := path + "_tmp"
		_ = removeBPFPin(tmpPath)
		if err = mapPinOps.Pin(targetMap, tmpPath); err != nil {
			_ = removeBPFPin(tmpPath)
			cleanupTmp()
			coll.Close()
			return nil, fmt.Errorf("pin BPF map %s to tmp %s: %w", name, tmpPath, err)
		}
		pending = append(pending, pendingPin{name: name, tmpPath: tmpPath, dstPath: path})
	}

	// Phase 2: 全部 Pin 成功，统一原子重命名覆盖目标路径。
	// 设计妥协说明：
	// 若 Phase 2 执行过程中出现极端罕见故障（如前一个 Map 重命名成功但后续 Map 失败），
	// 已成功覆盖的目标路径无法回滚至旧 Pin 文件。鉴于所有待重命名文件均位于同一 bpffs 挂载点内且权限目录均已验证就绪，
	// 内核 VFS rename(2) 操作在同文件系统内失败概率极低，此两阶段设计可在极小复杂度下杜绝绝大多数中间态。
	for _, p := range pending {
		if err = mapPinOps.Rename(p.tmpPath, p.dstPath); err != nil {
			cleanupTmp()
			coll.Close()
			return nil, fmt.Errorf("atomic swap pinned map %s: %w", p.dstPath, err)
		}
	}

	return coll, nil
}
