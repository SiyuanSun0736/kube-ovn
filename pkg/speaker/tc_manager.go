package speaker

import (
	"net"
)

// TCRedirectTarget 封装 eBPF TC 层重定向到本地网关 Pod 所需的目标网卡与 MAC 信息
type TCRedirectTarget struct {
	// TargetIfindex 是目标网关 Pod 在物理宿主机根网络空间所对应的接口设备索引（如 macvlan 子接口或 veth 宿主机端）
	TargetIfindex int
	// ParentIfindex 标识 macvlan 父设备接口索引（Master ifindex），本地重定向的数据包会在此重新进入 Ingress 接收队列
	ParentIfindex int
	// TargetMAC 是目标网关 Pod 内部外部网卡的真实硬件 MAC 地址
	TargetMAC net.HardwareAddr
}

// TCManager 接口定义了宿主机物理网卡 eBPF TC 分类器生命周期管理与 BPF Map 路由状态同步能力
type TCManager interface {
	// Attach 在指定的物理网卡上附加 clsact qdisc 并在 Ingress 钩子上挂载配置优先级的 eBPF TC 分类器
	Attach(linkName string) error

	// AddEIP 向内核 BPF Map (kube_ovn_eip_v4_map / kube_ovn_eip_v6_map) 中原子写入或更新 EIP ➔ 目标网卡重定向条目
	AddEIP(eip net.IP, target TCRedirectTarget) error

	// DeleteEIP 从内核 BPF Map 中原子剔除已销毁或漂移走的 EIP 重定向条目
	DeleteEIP(eip net.IP) error

	// SyncEIPs 根据当前期望的本地活跃 EIP 集合全量同步 BPF Map，执行对账并清理可能残留的孤儿条目。
	// retainedIPs 可选指定当前属于本节点活跃范围、但因 Pod 元数据短暂就绪延迟而需声明式保留存量映射的 EIP 列表。
	SyncEIPs(expected map[string]TCRedirectTarget, retainedIPs ...string) error

	// Close 仅释放用户态内存对象与 Map/Prog 的本地 FD 句柄，保留内核态已挂载的 Filter 与 Pinned Map。
	// 专为 Speaker 进程重启 / 热更新设计，实现数据面 0 丢包。
	// See also: Clean()
	Close() error

	// Clean 彻底销毁网卡 TC 过滤器并移除 Pinned Map 文件（仅在节点下线或功能彻底禁用时调用）。
	// 警告：此操作具有破坏性，会导致数据面直通中断，不可用于平滑热重启。
	// See also: Close()
	Clean() error
}

// NewTCManager 根据当前 Speaker 配置构造 TCManager 实例
func NewTCManager(cfg *Configuration) TCManager {
	return newTCManager(cfg)
}
