//go:build linux

package speaker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func TestIsMatchingTCFilter(t *testing.T) {
	baseFilter := func() *netlink.BpfFilter {
		return &netlink.BpfFilter{
			FilterAttrs: netlink.FilterAttrs{
				Handle:   1,
				Parent:   netlink.HANDLE_MIN_INGRESS,
				Protocol: unix.ETH_P_ALL,
				Priority: 5,
			},
			DirectAction: true,
			Name:         "test-tc-filter",
		}
	}

	// 1. nil filter
	require.False(t, isMatchingTCFilter(nil, 0, ""))

	// 2. 名字完全匹配
	f := baseFilter()
	require.True(t, isMatchingTCFilter(f, 0, "", "test-tc-filter"))

	// 3. 名字匹配多候选之一
	require.True(t, isMatchingTCFilter(f, 0, "", "other-filter", "test-tc-filter"))

	// 4. 强指纹：ProgID 匹配
	fID := baseFilter()
	fID.Name = ""
	fID.Id = 99
	require.True(t, isMatchingTCFilter(fID, 99, ""))

	// 5. 强指纹：ProgTag 匹配（忽略大小写）
	fTag := baseFilter()
	fTag.Name = ""
	fTag.Tag = "08218004CAF6A189"
	require.True(t, isMatchingTCFilter(fTag, 0, "08218004caf6a189"))

	// 6. Parent 不为 Ingress 拒绝
	fEgress := baseFilter()
	fEgress.Parent = netlink.HANDLE_MIN_EGRESS
	require.False(t, isMatchingTCFilter(fEgress, 0, "", "test-tc-filter"))

	// 7. DirectAction 为 false 拒绝
	fNoDA := baseFilter()
	fNoDA.DirectAction = false
	require.False(t, isMatchingTCFilter(fNoDA, 0, "", "test-tc-filter"))

	// 8. Handle != 1 拒绝
	fDiffHandle := baseFilter()
	fDiffHandle.Handle = 2
	require.False(t, isMatchingTCFilter(fDiffHandle, 0, "", "test-tc-filter"))

	// 9. 不匹配任何特征拒绝
	fAlien := baseFilter()
	fAlien.Name = "cilium-netdev"
	require.False(t, isMatchingTCFilter(fAlien, 0, "", "test-tc-filter"))
}

func TestRemoveBPFPin(t *testing.T) {
	// 1. 删除不存在的文件应安全返回 nil
	require.NoError(t, removeBPFPin("/tmp/nonexistent-pin-file-12345"))

	// 2. 正常删除已有文件
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test-pin")
	require.NoError(t, os.WriteFile(filePath, []byte("test"), 0o600))
	require.FileExists(t, filePath)
	require.NoError(t, removeBPFPin(filePath))
	require.NoFileExists(t, filePath)
}

func TestAttachBPFFilterNilLink(t *testing.T) {
	err := attachBPFFilter(tcFilterConfig{Link: nil})
	require.ErrorContains(t, err, "link is nil")
}
