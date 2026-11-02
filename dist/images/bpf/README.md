# EIP TC 旁路直通加速与抗 conntrack 击穿防护

Kube-OVN 主 Dockerfile 在独立构建阶段编译 `tc_eip.c`，并将 `/etc/kube-ovn/tc_eip.o` 同时打包进 amd64 与 arm64 镜像中。
在安装有 clang、libbpf 头文件及 Linux UAPI 头文件的 Linux 环境下，可直接执行 `make -C dist/images/bpf` 进行独立编译。

通过启动参数 `--node-route-eip-mode --enable-tc-bpf-redirect` 启用 Speaker 的 TC BPF 直通功能。
Speaker 必须运行在宿主机网络空间（hostNetwork），具备加载 BPF 程序、管理 TC 过滤器的权限，
并能够向挂载在 `/sys/fs/bpf` 的 bpffs 写入文件（通常为挂载宿主机 bpffs 的 privileged 容器）。

`--tc-bpf-interface` 指定物理 Ingress 接收网卡；若未指定，则自动探测默认 IPv4 出向网卡。
网关 Pod 的外网子网通过其 `ovn.kubernetes.io/nad-macvlan-master` 注解独立指定 macvlan 主设备（Master）。
宿主机 macvlan 必须以 bridge 模式运行。网关 Pod 的 Multus network-status 必须包含该外网 NAD 及其真实硬件 MAC 地址。
未解析或非法的目标会被拒绝，并在全量对账过程中从旁路 Map 中剔除并按需重试。
带有内联 VLAN Tag 的报文将直接放行，不走旁路直通。

`--tc-bpf-map-path` 指定 IPv4 Pinned Map 挂载路径，默认为 `/sys/fs/bpf/kube_ovn_eip_v4_map`。
IPv6 沿用该路径，将其中的 `_v4_map` 替换为 `_v6_map`（若无此后缀则直接追加 `_v6_map`）。
Speaker 仅加载并固化（Pin）当前附加程序实际使用的 Map，不会按名称随意探测其它已存在的 Pinned Map。

TC 过滤器挂载在 Ingress 方向，设置优先级 priority 5（可通过 `--tc-bpf-priority` 自定义，默认 5；预留 1-4 给更高优先级防护）、handle 1，名称为 `kube-ovn-eip`。
网络分工与架构设计：
- **Kube-OVN TC**：只占用外网出口网卡（如 `macbond1`），专职负责公网 EIP 的 eBPF 短路直通与防 conntrack 击穿；
- **Cilium TC**：只占用内网出口网卡（`bond0`，负责集群管理通信）以及各个 Pod 虚拟网卡 `*_h`（容器侧安全策略）。
两者在网络接口上严格物理隔离、各司其职，互不干扰。
若当前 priority 已被非本程序过滤器占用，系统将报错拒绝挂载以避免覆盖冲突。
非 EIP 报文返回 `TC_ACT_UNSPEC` (-1)，正常放行进入宿主机 Linux 网络协议栈。
在 direct-action 模式下，`TC_ACT_OK` (0) 专用于 macvlan bridge 本地交付链路（报文已改写为 Pod MAC 再次进入 Ingress 时），直接终止后续 TC 过滤并交由驱动将报文交付给目标 Pod RX 处理，彻底避开重定向死循环与宿主机 Netfilter / conntrack 跟踪表。
命中 EIP 的公网报文会通过宿主机 macvlan 直通转发，且目的 MAC 被改写为网关 Pod 的实际 MAC。
在 macvlan bridge 本地交付模式下，报文会重新进入 Master 的 Ingress；在 Master 上已寻址到 Pod MAC 的报文将直接返回 `TC_ACT_OK` 进入 macvlan RX 处理，避免重定向死循环，且彻底规避穿越宿主机 Netfilter 与 conntrack 跟踪表。

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
