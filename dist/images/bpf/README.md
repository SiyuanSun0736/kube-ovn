# EIP TC 旁路直通加速与抗 conntrack 击穿防护

Kube-OVN 主 Dockerfile 在独立构建阶段编译 `tc_eip.c`，并将 `/etc/kube-ovn/tc_eip.o` 打包进镜像。
在具备 clang 与 Linux UAPI 头文件环境下，可执行 `make -C dist/images/bpf` 独立编译。

## 编译与内核兼容性

`tc_eip.c` 仅依赖标准网络头（`ethhdr`, `iphdr`）、稳定 UAPI 结构 `struct __sk_buff` 以及基础 Helper（`bpf_map_lookup_elem`, `bpf_skb_store_bytes`, `bpf_redirect`）。
编译产物 `tc_eip.o` 为通用 ELF 目标文件，免去节点现场编译和 kernel-headers 依赖，直接适配 Linux 5.4+ 的所有 5.x 及 6.x 内核。

## 运行配置与网络架构

通过启动参数 `--node-route-eip-mode --enable-tc-bpf-redirect` 启用 Speaker 的 TC BPF 直通功能。
Speaker 必须运行在宿主机网络空间（hostNetwork），具备加载 BPF 程序、管理 TC 过滤器的权限，
并能够向挂载在 `/sys/fs/bpf` 的 bpffs 写入文件（通常为挂载宿主机 bpffs 的 privileged 容器）。

`--tc-bpf-interface` 指定物理 Ingress 接收网卡；若未指定，则自动探测默认 IPv4 出向网卡。
网关 Pod 的外网子网通过其 `ovn.kubernetes.io/nad-macvlan-master` 注解独立指定 macvlan 主设备（Master）。
宿主机 macvlan 必须以 bridge 模式运行。网关 Pod 的 Multus network-status 必须包含该外网 NAD 及其真实硬件 MAC 地址。
未解析或非法的目标会被拒绝，并在全量对账过程中从旁路 Map 中剔除并按需重试。
支持解析单层（802.1Q/802.1ad）及双层 QinQ VLAN Tag 报文，改写最外层 MAC 并保留 VLAN 封装直通网关 Pod。

`--tc-bpf-map-path` 指定 IPv4 Pinned Map 挂载路径，默认为 `/sys/fs/bpf/kube_ovn_eip_v4_map`。
IPv6 沿用该路径，将其中的 `_v4_map` 替换为 `_v6_map`（若无此后缀则直接追加 `_v6_map`）。
Speaker 仅加载并固化（Pin）当前附加程序实际使用的 Map，不会按名称随意探测其它已存在的 Pinned Map。

TC 过滤器挂载在 Ingress 方向，设置优先级 priority 5（可通过 `--tc-bpf-priority` 自定义，默认 5；预留 1-4 给更高优先级防护）、handle 1，名称为 `kube-ovn-eip`。
网络分工与架构设计：
- **Kube-OVN TC**：只占用外网出口网卡（如 `macbond1`），专职负责公网 EIP 的 eBPF 短路直通与防 conntrack 击穿；
- **Cilium TC**：只占用内网出口网卡（`bond0`，负责集群管理通信）以及各个 Pod 虚拟网卡 `*_h`（容器侧安全策略）。
两者在网络接口上严格物理隔离、各司其职，互不干扰。
若当前 priority 已被非本程序过滤器占用，系统将报错拒绝挂载以避免覆盖冲突。

### 数据面收发流向与 macvlan bridge 防死循环机制

程序工作在 `direct-action` (da) 模式下，针对不同数据流的处理流程如下：

1. **非 EIP 报文（宿主机自身流量 / 其它 Pod 流量）**：
   - 返回 `TC_ACT_UNSPEC` (-1)；
   - 保持数据包原样，正常放行进入宿主机 Linux 网络协议栈（路由选择、Netfilter / conntrack 等常规路径）。

2. **外网命中 EIP 报文（首次到达 Ingress）**：
   - 报文源 MAC 为上游交换机 MAC，目的 MAC 为宿主机物理网卡（如 `macbond1`）MAC；
   - BPF 查表命中 EIP，判定目标 MAC 非全零，调用 `bpf_skb_store_bytes` 将目的 MAC 原子改写为目标 Pod 的硬件 MAC；
   - 调用 `bpf_redirect(target_ifindex, 0)` 将报文重定向到目标网络设备（如宿主机上的 macvlan bridge 接口），跳过宿主机 conntrack。

3. **macvlan bridge 本地交付与防死循环链路（二次进入 Ingress）**：
   - **二次注入成因**：在 macvlan bridge 模式下，重定向到宿主机端 macvlan 接口后，macvlan 驱动会将报文的接收端借道父物理设备进行二层本地交换分发（`macvlan_handle_frame` / `netif_rx` 重新入队）。这会导致报文在宿主机软中断（SoftIRQ）上下文中**再次触发父网卡（`macbond1`）的 TC Ingress Hook**。
   - **死循环隐患**：此时报文的目的 IP 依然是该 EIP，若再次无脑重定向，报文将在物理网卡与 macvlan 间发生无限 Ping-Pong 死循环，直至软中断打满丢包。
   - **哨兵判定与终结**：BPF 程序入口处通过 `is_local_delivery` 校验：
     ```c
     skb->ifindex == target->parent_ifindex && mac_eq(eth->h_dest, target->target_mac)
     ```
     当前网卡为父网卡且**目的 MAC 已经被改写为 Pod MAC**，精准识别出这是“已完成旁路、正在本地交付中的报文”；
   - **直接放行与驱动截胡（彻底绕过宿主机 Netfilter）**：
     - 立即返回 `TC_ACT_OK` (0)，直接终止后续 TC 过滤；
     - 内核协议栈紧接着执行二层接收钩子 `rx_handler`（即 `macvlan_handle_frame`）。由于目的 MAC 已经是 Pod MAC，macvlan 驱动将报文判定归属目标 Pod，直接将其交由 Pod 接收队列（RX），并向宿主机内核返回 `RX_HANDLER_CONSUMED`；
     - 宿主机协议栈在此处直接中断处理流程（`goto out`），**根本不会继续向下流转至宿主机的 `ip_rcv` 与 Netfilter / conntrack 模块**。实机核验表明宿主机 conntrack 表项始终为 0，彻底切断重定向死循环并保持全局免跟踪。

### 数据流向拓扑链对比

- **架构 A（Veth-Pair 模式）**：
  ```text
  外部交换机 ──► macbond1 (TC Ingress) ──► bpf_redirect(vethXXXX, 0) ──► Pod 内部接收
  ```
  *(跨 Netns 独立设备直接投递，物理接口仅触发 1 次，天然不回流)*

- **架构 B（Macvlan Bridge 模式）**：
  ```text
  外部交换机 ──► macbond1 (首次 TC Ingress: 改写为 Pod MAC) ──► bpf_redirect(macvlan, 0) ──► macvlan bridge 驱动 (借道 netif_rx 回灌) ──► macbond1 (二次 TC Ingress: is_local_delivery 命中放行 TC_ACT_OK) ──► rx_handler (macvlan_handle_frame 消费) ──► Pod 内部接收
  ```
  *(借道 SoftIRQ 重新入队触发二次 Ingress，由哨兵精准放行并经驱动层截胡，全链路 0 宿主机 conntrack)*

### 热更新与生命周期管理

参考 Cilium Agent 的平滑热更新机制，Speaker 重启或更新期间**内核 BPF 程序永不 detach**：
- **热更新与复用**：新实例启动时若程序未变直接复用，若字节码升级则通过 `netlink.FilterReplace` 在内核原子替换（Atomic Swap），全程零中断；
- **主动卸载**：仅在显式配置 `--enable-tc-bpf-redirect=false` 停用功能或迁移 priority 时，才会主动 detach 清理。

启动时会用全新代际的空 Map 替换旧的 EIP 挂载，并通过最佳努力（Best-effort）机制在宣布 BGP 路由的同时异步建立与核对 TC BPF 直通映射；若 TC 挂载或映射同步发生失败，系统将自动回退到宿主机内核协议栈转发以保障 BGP 路由与服务的连续性。
周期性对账逻辑会自动重试加载或挂载失败的操作、更新变动的目标网卡信息，并清理已漂移到其它节点的陈旧条目。
正常退出时会清理自己拥有的过滤器和固化的 Pinned Map。
若需停用本特性，请保持 node-route 模式开启并以 `--enable-tc-bpf-redirect=false` 重启；下一个 Speaker 实例也会主动清理非正常退出留下的挂载痕迹。
严禁在同一物理节点上同时运行多个 node-route speaker。

## 内核回归测试

此类测试必须在具备 bpffs、`ip` 和 `ping` 工具的专用隔离特权 Linux 容器或独立 Network Namespace 中运行。
测试会创建私有虚拟接口和网络命名空间，真实执行 BPF Map 增删与程序校验，并验证在宿主机故意不配置 EIP 路由的情况下通过 TC 旁路实现 macvlan 连通性。

```sh
make -C dist/images/bpf
KUBE_OVN_BPF_TEST=1 \
KUBE_OVN_BPF_OBJECT="$PWD/dist/images/bpf/tc_eip.o" \
go test ./pkg/speaker -run 'Test(TC|NatGwTarget)' -count=1 -v
```

若未设置 `KUBE_OVN_BPF_TEST=1`，需要内核特权的测试用例将被自动跳过。
