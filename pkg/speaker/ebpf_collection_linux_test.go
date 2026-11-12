//go:build linux

package speaker

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

type mockBPFPinOps struct {
	defaultBPFPinOps
	failMkdirAll  bool
	failPinPrefix string
	failRename    bool
	failRemove    bool
}

func (m *mockBPFPinOps) MkdirAll(path string, perm os.FileMode) error {
	if m.failMkdirAll {
		return errors.New("mock mkdirall error")
	}
	return m.defaultBPFPinOps.MkdirAll(path, perm)
}

func (m *mockBPFPinOps) Pin(target *ebpf.Map, path string) error {
	if m.failPinPrefix != "" && strings.Contains(path, m.failPinPrefix) {
		return errors.New("mock pin error")
	}
	return m.defaultBPFPinOps.Pin(target, path)
}

func (m *mockBPFPinOps) Rename(oldpath, newpath string) error {
	if m.failRename {
		return errors.New("mock rename error")
	}
	return m.defaultBPFPinOps.Rename(oldpath, newpath)
}

func (m *mockBPFPinOps) Remove(path string) error {
	if m.failRemove {
		return errors.New("mock remove error")
	}
	return m.defaultBPFPinOps.Remove(path)
}

func TestTCRemoveBPFPin(t *testing.T) {
	// 1. 删除不存在的文件应安全返回 nil
	require.NoError(t, removeBPFPin("/tmp/nonexistent-pin-file-12345"))

	// 2. 正常删除已有文件
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test-pin")
	require.NoError(t, os.WriteFile(filePath, []byte("test"), 0o600))
	require.FileExists(t, filePath)
	require.NoError(t, removeBPFPin(filePath))
	require.NoFileExists(t, filePath)

	// 3. 底层 Remove 返回非 ErrNotExist 错误
	origOps := pinOps
	pinOps = &mockBPFPinOps{failRemove: true}
	t.Cleanup(func() { pinOps = origOps })
	require.ErrorContains(t, removeBPFPin("/tmp/any"), "mock remove error")
}

func TestTCLoadCollectionWithPinnedMaps(t *testing.T) {
	newSpec := func(valueSize, maxEntries uint32) *ebpf.CollectionSpec {
		return &ebpf.CollectionSpec{Maps: map[string]*ebpf.MapSpec{
			"test_map": {Name: "test_map", Type: ebpf.Hash, KeySize: 4, ValueSize: valueSize, MaxEntries: maxEntries},
		}}
	}

	t.Run("fresh", func(t *testing.T) {
		path := testTCPinPath(t)
		spec := newSpec(4, 1)
		coll, err := loadCollectionWithPinnedMaps(spec, map[string]string{"test_map": path})
		require.NoError(t, err)
		t.Cleanup(coll.Close)

		require.NoError(t, coll.Maps["test_map"].Put(uint32(1), uint32(42)))
		pinned, err := ebpf.LoadPinnedMap(path, nil)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, pinned.Close()) })

		var value uint32
		require.NoError(t, pinned.Lookup(uint32(1), &value))
		require.Equal(t, uint32(42), value)

		entries, err := os.ReadDir(filepath.Dir(path))
		require.NoError(t, err)
		require.Len(t, entries, 1, "temporary pins must not remain after commit")
	})

	t.Run("incompatible_fallback", func(t *testing.T) {
		path := testTCPinPath(t)
		// 预设不兼容 Spec 的旧 Pin
		old, err := ebpf.NewMap(&ebpf.MapSpec{Type: ebpf.Hash, KeySize: 4, ValueSize: 8, MaxEntries: 1})
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, old.Close()) })
		require.NoError(t, old.Pin(path))

		// 触发 fallback 重建
		spec := newSpec(4, 1)
		coll, err := loadCollectionWithPinnedMaps(spec, map[string]string{"test_map": path})
		require.NoError(t, err)
		t.Cleanup(coll.Close)

		require.NoError(t, coll.Maps["test_map"].Put(uint32(1), uint32(42)))
		pinned, err := ebpf.LoadPinnedMap(path, nil)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, pinned.Close()) })

		var value uint32
		require.NoError(t, pinned.Lookup(uint32(1), &value))
		require.Equal(t, uint32(42), value)

		entries, err := os.ReadDir(filepath.Dir(path))
		require.NoError(t, err)
		require.Len(t, entries, 1, "temporary pins must not remain after commit")
	})

	t.Run("reused_hot_reload_preserves_entries_without_tmp", func(t *testing.T) {
		path := testTCPinPath(t)
		spec := newSpec(4, 10)

		// 第一次调用：初始化并写入多个测试 Key
		coll1, err := loadCollectionWithPinnedMaps(spec, map[string]string{"test_map": path})
		require.NoError(t, err)
		require.NoError(t, coll1.Maps["test_map"].Put(uint32(1), uint32(101)))
		require.NoError(t, coll1.Maps["test_map"].Put(uint32(2), uint32(202)))
		// 模拟服务退出 / 重启：关闭 Collection 句柄，但 Pinned Map 依然保留在 bpffs
		coll1.Close()

		// 第二次调用：同 Spec 与路径再次加载，验证复用成功路径
		coll2, err := loadCollectionWithPinnedMaps(spec, map[string]string{"test_map": path})
		require.NoError(t, err)
		t.Cleanup(coll2.Close)

		// 断言已编程条目未丢失（实现热更新 0 丢包）
		var v1, v2 uint32
		require.NoError(t, coll2.Maps["test_map"].Lookup(uint32(1), &v1))
		require.Equal(t, uint32(101), v1)
		require.NoError(t, coll2.Maps["test_map"].Lookup(uint32(2), &v2))
		require.Equal(t, uint32(202), v2)

		// 断言复用路径跳过 Phase 1/2，未产生任何 _tmp 残留
		require.NoFileExists(t, path+"_tmp", "reused path must skip Phase 1 and not create _tmp pin")
		entries, err := os.ReadDir(filepath.Dir(path))
		require.NoError(t, err)
		require.Len(t, entries, 1, "only the pinned map file must exist in bpffs directory")
	})

	t.Run("reused_pinned_map_handle_closed_no_fd_leak", func(t *testing.T) {
		path := testTCPinPath(t)
		spec := newSpec(4, 1)

		// 1. 初始化 Pin 文件
		coll1, err := loadCollectionWithPinnedMaps(spec, map[string]string{"test_map": path})
		require.NoError(t, err)
		coll1.Close()

		countFDs := func() int {
			entries, err := os.ReadDir("/proc/self/fd")
			if err != nil {
				t.Skipf("cannot read /proc/self/fd: %v", err)
			}
			return len(entries)
		}

		before := countFDs()

		// 2. 模拟热重载加载既有 Pin
		coll2, err := loadCollectionWithPinnedMaps(spec, map[string]string{"test_map": path})
		require.NoError(t, err)

		// 3. 关闭 coll2 后，打开的 FD 数量必须精准恢复，断言 LoadPinnedMap 打开的临时 FD 已被 defer 显式 Close
		coll2.Close()
		after := countFDs()

		require.Equal(t, before, after, "replacements handle must be closed to prevent FD leaks")
	})
}

func TestTCLoadCollectionTwoPhaseCommitFailures(t *testing.T) {
	newSpec := func() *ebpf.CollectionSpec {
		return &ebpf.CollectionSpec{Maps: map[string]*ebpf.MapSpec{
			"map1": {Name: "map1", Type: ebpf.Hash, KeySize: 4, ValueSize: 4, MaxEntries: 1},
			"map2": {Name: "map2", Type: ebpf.Hash, KeySize: 4, ValueSize: 4, MaxEntries: 1},
		}}
	}

	createDummyPinnedMap := func(path string, val uint32) *ebpf.Map {
		m, err := ebpf.NewMap(&ebpf.MapSpec{Type: ebpf.Hash, KeySize: 4, ValueSize: 4, MaxEntries: 1})
		require.NoError(t, err)
		t.Cleanup(func() { _ = m.Close() })
		require.NoError(t, m.Pin(path))
		require.NoError(t, m.Put(uint32(1), val))
		return m
	}

	assertPinnedValue := func(path string, expected uint32) {
		pinned, err := ebpf.LoadPinnedMap(path, nil)
		require.NoError(t, err)
		defer pinned.Close()
		var val uint32
		require.NoError(t, pinned.Lookup(uint32(1), &val))
		require.Equal(t, expected, val)
	}

	origOps := pinOps
	t.Cleanup(func() { pinOps = origOps })

	t.Run("mkdir_failure_preserves_old_and_cleans_tmp", func(t *testing.T) {
		p1 := testTCPinPath(t)
		p2 := filepath.Join(filepath.Dir(p1), "map2")
		// 预先写入 dummy 旧 Pin
		createDummyPinnedMap(p1, 991)

		pinOps = &mockBPFPinOps{failMkdirAll: true}
		coll, err := loadCollectionWithPinnedMaps(newSpec(), map[string]string{"map1": p1, "map2": p2})
		require.ErrorContains(t, err, "create map dir")
		require.Nil(t, coll)

		// 断言旧 Pin 未受破坏且数据完整
		assertPinnedValue(p1, 991)
	})

	t.Run("phase1_pin_failure_cleans_tmp_and_preserves_old", func(t *testing.T) {
		p1 := testTCPinPath(t)
		p2 := filepath.Join(filepath.Dir(p1), "map2")
		// 预先写入两个旧 Pin
		createDummyPinnedMap(p1, 991)
		createDummyPinnedMap(p2, 992)

		// 模拟 map2 Pin 失败（因为 spec 与已有 map 兼容会触发 reuse，为了迫使进入 Phase 1，我们用不兼容 spec 触发 fallback）
		fallbackSpec := &ebpf.CollectionSpec{Maps: map[string]*ebpf.MapSpec{
			"map1": {Name: "map1", Type: ebpf.Hash, KeySize: 4, ValueSize: 8, MaxEntries: 1},
			"map2": {Name: "map2", Type: ebpf.Hash, KeySize: 4, ValueSize: 8, MaxEntries: 1},
		}}
		pinOps = &mockBPFPinOps{failPinPrefix: "map2"}
		coll, err := loadCollectionWithPinnedMaps(fallbackSpec, map[string]string{"map1": p1, "map2": p2})
		require.ErrorContains(t, err, "pin BPF map")
		require.Nil(t, coll)

		// 验证 Phase 1 回滚：map1 的 _tmp 必须已被 cleanupTmp() 清理
		require.NoFileExists(t, p1+"_tmp", "temporary pin for map1 must be cleaned up on failure")
		require.NoFileExists(t, p2+"_tmp")

		// 验证旧 Pin 绝未被破坏或提前删除
		assertPinnedValue(p1, 991)
		assertPinnedValue(p2, 992)
	})

	t.Run("phase2_rename_failure_cleans_tmp", func(t *testing.T) {
		p1 := testTCPinPath(t)
		pinOps = &mockBPFPinOps{failRename: true}
		coll, err := loadCollectionWithPinnedMaps(newSpec(), map[string]string{"map1": p1})
		require.ErrorContains(t, err, "atomic swap pinned map")
		require.Nil(t, coll)

		// 验证临时 Pin 已被清理
		require.NoFileExists(t, p1+"_tmp")
	})
}

func TestTCPruneBPFHashMapNil(t *testing.T) {
	require.ErrorContains(t, pruneBPFHashMap(nil, map[string]bool{}), "target map is nil")
}

func TestTCLoadCollectionWithPinnedMapsMissingMap(t *testing.T) {
	if os.Getenv("KUBE_OVN_BPF_TEST") != "1" {
		t.Skip("requires privileged BPF environment")
	}
	if err := unix.Mount("bpf", "/sys/fs/bpf", "bpf", 0, ""); err != nil {
		require.ErrorIs(t, err, unix.EBUSY)
	}
	dir, err := os.MkdirTemp("/sys/fs/bpf", "bpf-missing-")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.RemoveAll(dir)) })

	spec := &ebpf.CollectionSpec{
		Maps: map[string]*ebpf.MapSpec{
			"valid_map": {
				Type:       ebpf.Hash,
				KeySize:    4,
				ValueSize:  4,
				MaxEntries: 10,
			},
		},
	}
	pinPaths := map[string]string{
		"missing_map": filepath.Join(dir, "missing_map"),
	}
	coll, err := loadCollectionWithPinnedMaps(spec, pinPaths)
	require.ErrorContains(t, err, "missing_map not found in collection")
	require.Nil(t, coll)
}
