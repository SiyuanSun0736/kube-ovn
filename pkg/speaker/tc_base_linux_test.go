//go:build linux

package speaker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func TestTCIsMatchingTCFilter(t *testing.T) {
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

type nilAttrsLink struct{ netlink.Link }

func (nilAttrsLink) Attrs() *netlink.LinkAttrs { return nil }

func TestTCAttachBPFFilterNilLink(t *testing.T) {
	err := attachBPFFilter(tcFilterConfig{Link: nil})
	require.ErrorContains(t, err, "link is nil")

	// 验证 Link.Attrs() 为 nil 的防御分支
	err = attachBPFFilter(tcFilterConfig{
		Link:   nilAttrsLink{},
		ProgFD: 1,
	})
	require.ErrorContains(t, err, "link attrs is nil")
}

func TestTCAttachBPFFilterInvalidFD(t *testing.T) {
	err := attachBPFFilter(tcFilterConfig{
		Link:   &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
		ProgFD: 0,
	})
	require.ErrorContains(t, err, "bpf program fd is invalid")
}

func TestTCAttachBPFFilterEmptyAttachName(t *testing.T) {
	err := attachBPFFilter(tcFilterConfig{
		Link:   &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
		ProgFD: 1,
	})
	require.ErrorContains(t, err, "AttachName must not be empty")
}

type mockTCNetlinkOps struct {
	qdiscAddErr       error
	filterListFilters []netlink.Filter
	filterListErr     error
	filterReplaceErr  error
	filterDelErr      error

	replacedFilters []netlink.Filter
	deletedFilters  []netlink.Filter

	linkByNameFunc  func(string) (netlink.Link, error)
	linkByIndexFunc func(int) (netlink.Link, error)
	linkListFunc    func() ([]netlink.Link, error)
}

func (m *mockTCNetlinkOps) QdiscAdd(netlink.Qdisc) error {
	return m.qdiscAddErr
}

func (m *mockTCNetlinkOps) FilterList(netlink.Link, uint32) ([]netlink.Filter, error) {
	return m.filterListFilters, m.filterListErr
}

func (m *mockTCNetlinkOps) FilterReplace(f netlink.Filter) error {
	m.replacedFilters = append(m.replacedFilters, f)
	return m.filterReplaceErr
}

func (m *mockTCNetlinkOps) FilterDel(f netlink.Filter) error {
	m.deletedFilters = append(m.deletedFilters, f)
	return m.filterDelErr
}

func (m *mockTCNetlinkOps) LinkByName(name string) (netlink.Link, error) {
	if m.linkByNameFunc != nil {
		return m.linkByNameFunc(name)
	}
	return &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: 10, Name: name}}, nil
}

func (m *mockTCNetlinkOps) LinkByIndex(index int) (netlink.Link, error) {
	if m.linkByIndexFunc != nil {
		return m.linkByIndexFunc(index)
	}
	return &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: index, Name: fmt.Sprintf("mock%d", index)}}, nil
}

func (m *mockTCNetlinkOps) LinkList() ([]netlink.Link, error) {
	if m.linkListFunc != nil {
		return m.linkListFunc()
	}
	return nil, nil
}

func TestTCAttachBPFFilterOpsBranches(t *testing.T) {
	origOps := netlinkOps
	defer func() { netlinkOps = origOps }()

	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: 10, Name: "dummy0"}}
	const (
		progFD  = 42
		progID  = 1001
		progTag = "08218004caf6a189"
	)

	// 1. 幂等短路分支：同 priority 上已存在匹配本程序有效 ProgID 的规则
	mock := &mockTCNetlinkOps{
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Id:           1001,
				Tag:          "08218004caf6a189",
				Name:         "kube-ovn-eip",
			},
		},
	}
	netlinkOps = mock
	err := attachBPFFilter(tcFilterConfig{
		Link:        link,
		ProgFD:      progFD,
		ProgID:      progID,
		ProgTag:     progTag,
		AttachName:  "kube-ovn-eip",
		LegacyNames: []string{"tc_eip.o:[classifier]"},
		Priority:    5,
	})
	require.NoError(t, err)
	require.Empty(t, mock.replacedFilters)
	require.Empty(t, mock.deletedFilters)

	// 2. 异 priority 剥离分支：旧 priority 7 存在旧名称规则，新挂载在 priority 5
	mock = &mockTCNetlinkOps{
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 7, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Name:         "tc_eip.o:[classifier]",
			},
		},
	}
	netlinkOps = mock
	err = attachBPFFilter(tcFilterConfig{
		Link:        link,
		ProgFD:      progFD,
		ProgID:      progID,
		ProgTag:     progTag,
		AttachName:  "kube-ovn-eip",
		LegacyNames: []string{"tc_eip.o:[classifier]"},
		Priority:    5,
	})
	require.NoError(t, err)
	require.Len(t, mock.deletedFilters, 1)
	require.Equal(t, uint16(7), mock.deletedFilters[0].Attrs().Priority)
	require.Len(t, mock.replacedFilters, 1)
	require.Equal(t, uint16(5), mock.replacedFilters[0].Attrs().Priority)

	// 3. 同 priority 冲突分支：目标 priority 5 已被第三方规则占用（名字不匹配）
	mock = &mockTCNetlinkOps{
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Name:         "cilium-netdev",
			},
		},
	}
	netlinkOps = mock
	err = attachBPFFilter(tcFilterConfig{
		Link:        link,
		ProgFD:      progFD,
		ProgID:      progID,
		ProgTag:     progTag,
		AttachName:  "kube-ovn-eip",
		LegacyNames: []string{"tc_eip.o:[classifier]"},
		Priority:    5,
	})
	require.ErrorContains(t, err, "TC priority 5 on dummy0 is already owned by another filter")
	require.Empty(t, mock.replacedFilters)
	require.Empty(t, mock.deletedFilters)

	// 4. 消除 Ordering Hazard 分支：同 priority 5 已存在匹配本程序的规则，且排在后面的 priority 7 存在旧名称规则
	// 必须能够短路（不触发 FilterReplace），同时不能提前 return，必须成功剥离 priority 7 的旧规则
	mock = &mockTCNetlinkOps{
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Id:           1001,
				Tag:          "08218004caf6a189",
				Name:         "kube-ovn-eip",
			},
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 7, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Name:         "tc_eip.o:[classifier]",
			},
		},
	}
	netlinkOps = mock
	err = attachBPFFilter(tcFilterConfig{
		Link:        link,
		ProgFD:      progFD,
		ProgID:      progID,
		ProgTag:     progTag,
		AttachName:  "kube-ovn-eip",
		LegacyNames: []string{"tc_eip.o:[classifier]"},
		Priority:    5,
	})
	require.NoError(t, err)
	require.Empty(t, mock.replacedFilters)
	require.Len(t, mock.deletedFilters, 1)
	require.Equal(t, uint16(7), mock.deletedFilters[0].Attrs().Priority)

	// 5. ProgID 不匹配触发就地替换分支：同 priority 上存在属于本组件的旧规则（Name 匹配），但 ProgID 不一致（如升级或重启后重新加载），必须执行 FilterReplace
	mock = &mockTCNetlinkOps{
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Id:           999, // 旧程序 ID
				Name:         "kube-ovn-eip",
			},
		},
	}
	netlinkOps = mock
	err = attachBPFFilter(tcFilterConfig{
		Link:        link,
		ProgFD:      progFD,
		ProgID:      progID,
		ProgTag:     progTag,
		AttachName:  "kube-ovn-eip",
		LegacyNames: []string{"tc_eip.o:[classifier]"},
		Priority:    5,
	})
	require.NoError(t, err)
	require.Len(t, mock.replacedFilters, 1)
	require.Equal(t, uint16(5), mock.replacedFilters[0].Attrs().Priority)
	require.Empty(t, mock.deletedFilters)

	// 6. progID == 0 防御分支：若 program 无有效内核 ID（progID == 0），即使现有规则 bpf.Id == 0 且 Name 匹配，绝不可误判定为 matchedCurrent 发生零值误判短路，必须执行 FilterReplace
	mock = &mockTCNetlinkOps{
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Id:           0,
				Name:         "kube-ovn-eip",
			},
		},
	}
	netlinkOps = mock
	err = attachBPFFilter(tcFilterConfig{
		Link:        link,
		ProgFD:      42,
		ProgID:      0,
		ProgTag:     "",
		AttachName:  "kube-ovn-eip",
		LegacyNames: []string{"tc_eip.o:[classifier]"},
		Priority:    5,
	})
	require.NoError(t, err)
	require.Len(t, mock.replacedFilters, 1)
	require.Empty(t, mock.deletedFilters)

	// 7. 历史名称规则就地规范化分支：同 priority 上存在属于本组件的旧规则（ProgID 一致但 Name 为历史旧名），必须执行 FilterReplace 规范化为 AttachName
	mock = &mockTCNetlinkOps{
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Id:           1001,
				Tag:          "08218004caf6a189",
				Name:         "tc_eip.o:[classifier]",
			},
		},
	}
	netlinkOps = mock
	err = attachBPFFilter(tcFilterConfig{
		Link:        link,
		ProgFD:      progFD,
		ProgID:      progID,
		ProgTag:     progTag,
		AttachName:  "kube-ovn-eip",
		LegacyNames: []string{"tc_eip.o:[classifier]"},
		Priority:    5,
	})
	require.NoError(t, err)
	require.Len(t, mock.replacedFilters, 1)
	require.Equal(t, uint16(5), mock.replacedFilters[0].Attrs().Priority)
	require.Equal(t, "kube-ovn-eip", mock.replacedFilters[0].(*netlink.BpfFilter).Name)
	require.Empty(t, mock.deletedFilters)
}

func TestTCDetachBPFFiltersOnLinksWithOps(t *testing.T) {
	origOps := netlinkOps
	defer func() { netlinkOps = origOps }()

	// 1. linkName 为空时，通过 ops.LinkList 枚举网卡并卸载匹配规则，且支持跳过 nil link
	mock := &mockTCNetlinkOps{
		linkListFunc: func() ([]netlink.Link, error) {
			return []netlink.Link{
				nil,
				nilAttrsLink{},
				&netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: 1, Name: "dummy1"}},
			}, nil
		},
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 1, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Name:         "kube-ovn-eip",
			},
		},
	}
	netlinkOps = mock

	err := detachBPFFiltersOnLinks("", func(b *netlink.BpfFilter) bool {
		return b.Name == "kube-ovn-eip"
	})
	require.NoError(t, err)
	require.Len(t, mock.deletedFilters, 1)

	// 2. linkName 非空时，仅通过 ops.LinkByName 查询指定网卡
	mock = &mockTCNetlinkOps{
		linkByNameFunc: func(name string) (netlink.Link, error) {
			require.Equal(t, "eth0", name)
			return &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: 2, Name: "eth0"}}, nil
		},
	}
	netlinkOps = mock
	err = detachBPFFiltersOnLinks("eth0", func(_ *netlink.BpfFilter) bool { return false })
	require.NoError(t, err)

	// 3. LinkByName 返回 unix.ENOENT 或 unix.ENODEV 时视为已清理返回 nil
	for _, notFoundErr := range []error{unix.ENOENT, unix.ENODEV, netlink.LinkNotFoundError{}} {
		mock = &mockTCNetlinkOps{
			linkByNameFunc: func(_ string) (netlink.Link, error) {
				return nil, notFoundErr
			},
		}
		netlinkOps = mock
		require.NoError(t, detachBPFFiltersOnLinks("gone0", func(_ *netlink.BpfFilter) bool { return true }))
	}

	// 4. FilterDel 返回 unix.ESRCH 时静默忽略
	mock = &mockTCNetlinkOps{
		linkByNameFunc: func(name string) (netlink.Link, error) {
			return &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: 3, Name: name}}, nil
		},
		filterListFilters: []netlink.Filter{
			&netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 3, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				DirectAction: true,
				Name:         "kube-ovn-eip",
			},
		},
		filterDelErr: unix.ESRCH,
	}
	netlinkOps = mock
	require.NoError(t, detachBPFFiltersOnLinks("eth1", func(_ *netlink.BpfFilter) bool { return true }))
}

func TestTCIsMatchingTCFilterB1Regression(t *testing.T) {
	// 验证 B1 回归场景：
	// 旧版本可能保留了名字为 "tc_eip.o:[classifier]" 的 Filter 规则，
	// 且由于升级重编译，ProgTag 与内核 ProgID 均不相同。
	// 新版本必须仅通过 FilterNames 命中并识别该旧规则，以支持跨 Priority 安全剥离。
	legacyFilter := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			Handle:   1,
			Parent:   netlink.HANDLE_MIN_INGRESS,
			Protocol: unix.ETH_P_ALL,
			Priority: 5,
		},
		DirectAction: true,
		Id:           1001,                    // 旧程序在内核中的 ProgID
		Tag:          "aaaaaaaaaaaaaaaa",      // 旧程序字节码 Tag
		Name:         "tc_eip.o:[classifier]", // 历史版本规则名
	}

	newProgID := uint32(2002)
	newProgTag := "bbbbbbbbbbbbbbbb"
	filterNames := []string{tcEIPFilterName, "tc_eip.o:[classifier]"}

	// 1. 在 ProgID 与 ProgTag 完全不一致的情况下，必须通过历史 FilterName 成功识别旧规则（B1 核心诉求）
	require.True(t, isMatchingTCFilter(legacyFilter, newProgID, newProgTag, filterNames...))

	// 2. 若未配置历史 FilterName，则由于 ID/Tag 均不匹配而判定非本程序，无法识别
	require.False(t, isMatchingTCFilter(legacyFilter, newProgID, newProgTag, tcEIPFilterName))

	// 3. 第三方（如 Cilium/Calico）规则绝不误匹配
	alienFilter := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			Handle:   1,
			Parent:   netlink.HANDLE_MIN_INGRESS,
			Protocol: unix.ETH_P_ALL,
			Priority: 10,
		},
		DirectAction: true,
		Name:         "cilium-netdev",
	}
	require.False(t, isMatchingTCFilter(alienFilter, newProgID, newProgTag, filterNames...))
}

func TestTCAttachBPFFilterMigrationFailurePreservesOldFilter(t *testing.T) {
	for _, conflict := range []bool{true, false} {
		name := "replace failure"
		if conflict {
			name = "priority conflict"
		}
		t.Run(name, func(t *testing.T) {
			old := &netlink.BpfFilter{
				FilterAttrs: netlink.FilterAttrs{
					LinkIndex: 10, Parent: netlink.HANDLE_MIN_INGRESS, Handle: 1, Priority: 5, Protocol: unix.ETH_P_ALL,
				},
				Name: "kube-ovn-eip", DirectAction: true,
			}
			mock := &mockTCNetlinkOps{
				filterListFilters: []netlink.Filter{old},
				filterReplaceErr:  unix.EPERM,
			}
			if conflict {
				mock.filterListFilters = append(mock.filterListFilters, &netlink.MatchAll{
					FilterAttrs: netlink.FilterAttrs{Priority: 7},
				})
			}
			origOps := netlinkOps
			netlinkOps = mock
			t.Cleanup(func() { netlinkOps = origOps })

			err := attachBPFFilter(tcFilterConfig{
				Link:       &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Index: 10, Name: "dummy0"}},
				ProgFD:     42,
				AttachName: "kube-ovn-eip", Priority: 7,
			})
			if conflict {
				require.ErrorContains(t, err, "already owned by another filter")
				require.Empty(t, mock.replacedFilters)
			} else {
				require.ErrorIs(t, err, unix.EPERM)
				require.Len(t, mock.replacedFilters, 1)
			}
			require.Empty(t, mock.deletedFilters, "failed migration must preserve the working filter")
		})
	}
}

func TestTCDetachBPFFiltersOnLinksNilMatchFunc(t *testing.T) {
	err := detachBPFFiltersOnLinks("eth0", nil)
	require.ErrorContains(t, err, "matchFunc must not be nil")
}
