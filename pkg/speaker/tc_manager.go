package speaker

import (
	"net"
)

// TargetInfo 封装 eBPF TC 层重定向到本地网关 Pod 所需的目标网卡信息
type TargetInfo struct {
	// TargetIfindex 是目标网关 Pod 在物理宿主机根网络空间所对应的接口设备索引（如 macvlan 子接口或 veth 宿主机端）
	TargetIfindex int
	// ParentIfindex identifies the macvlan master, where locally redirected packets re-enter ingress.
	ParentIfindex int
	// TargetMAC is the external interface MAC inside the gateway Pod.
	TargetMAC net.HardwareAddr
}

// TCManager 接口定义了宿主机物理网卡 eBPF TC 分类器生命周期管理与 BPF Map 路由状态同步能力
type TCManager interface {
	// Attach 在指定的物理网卡上附加 clsact qdisc 并在 Ingress 钩子上挂载优先级为 1 的 eBPF TC 分类器
	Attach(linkName string) error

	// Detach 从指定网卡上卸载 eBPF TC 分类器
	Detach(linkName string) error

	// AddEIP 向内核 BPF Map (kube_ovn_eip_v4_map / kube_ovn_eip_v6_map) 中原子写入或更新 EIP ➔ 目标网卡重定向条目
	AddEIP(eip net.IP, info TargetInfo) error

	// DeleteEIP 从内核 BPF Map 中原子剔除已销毁或漂移走的 EIP 重定向条目
	DeleteEIP(eip net.IP) error

	// SyncEIPs 根据当前期望的本地活跃 EIP 集合全量同步 BPF Map，执行对账并清理可能残留的孤儿条目
	SyncEIPs(expected map[string]TargetInfo) error

	// Close 关闭并释放已打开的 BPF 文件描述符与相关句柄
	Close() error
}

// NewTCManager 根据当前 Speaker 配置构造对应平台的 TCManager 实例
// 在 Linux 平台下返回具备实际 Netlink 挂载与 BPF Syscall 能力的实例；在非 Linux 平台下返回空实现供跨平台编译
func NewTCManager(cfg *Configuration) TCManager {
	return newPlatformTCManager(cfg)
}
