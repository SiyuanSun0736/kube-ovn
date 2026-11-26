# VPC NAT Gateway nftables 表与规则参考

Kube-OVN VPC NAT Gateway 在运行中使用两张 nftables 表：
1. `table ip kube_ovn_stateless_nat`：无状态 NAT 表（`spec.dataplaneMode: stateless-nft` 时启用），处理南北向 FIP、独占 DNAT、单后端及多后端公网 EIP Service（固定 `N = 127` 五元组 Maglev 一致性哈希），通过 `notrack` 绕过 conntrack。
2. `table ip kube-ovn`：有状态与限速表，处理有状态网关的 Share-DNAT、ClusterIP / lanIP 内网负载均衡、无状态模式下的多后端发卡（Hairpin）兜底、TCP SYN 限速（`syn-limit`）与 TCP 并发连接限制（`conn-limit`）。

---

## 1. 数据面分流

| 流量类型 | 命中表 / 链路 | conntrack 状态 | 地址与端口转换 |
| :--- | :--- | :--- | :--- |
| 南北向 FIP（1:1） | `kube_ovn_stateless_nat` | `untracked` (notrack) | 入向 `sl_fip_dnat` 改目的 IP，出向 `sl_fip_snat` 改源 IP |
| 南北向独占 DNAT | `kube_ovn_stateless_nat` | `untracked` (notrack) | 入向 `sl_dnat_in` 改目的 IP:Port，出向 `sl_dnat_out` 改源 IP:Port |
| 单后端公网 EIP Service | `kube_ovn_stateless_nat` | `untracked` (notrack) | 动态写入 `sl_dnat_in` 与 `sl_dnat_out` 直通 |
| 多后端公网 EIP Service | `kube_ovn_stateless_nat` | `untracked` (notrack) | 入向五元组 `jhash mod 127` Maglev 一致性哈希分发至 `_ep_<bkhash>` 子链，出向汇聚至共享 `sl_dout_*` |
| ClusterIP / lanIP / 多后端发卡 / ClientIP 亲和性多后端 | `kube-ovn` | `tracked` (有状态) | 未经过公网卡，或需按客户端 IP 会话保持时保留在 `service-ips` vmap 走有状态路径 |
| 单后端发卡流量 (Hairpin) | `iptables` (HAIRPIN 链) | `tracked` (有状态) | 未命中出向公网卡 FIB 校验，回退有状态链路 |
| 多对一 SNAT 网段 | `iptables` (SNAT 链) | `tracked` (有状态) | 无状态回程无法反解多对一端口，走有状态路径 |

---

## 2. 无状态 NAT 表：`kube_ovn_stateless_nat`

### 2.1 基础链与 Hook 优先级

```nft
add table ip kube_ovn_stateless_nat
add chain ip kube_ovn_stateless_nat prerouting  { type filter hook prerouting priority -250; policy accept; }
add chain ip kube_ovn_stateless_nat forward     { type filter hook forward priority mangle; policy accept; }
add chain ip kube_ovn_stateless_nat postrouting { type filter hook postrouting priority 99; policy accept; }
```

#### prerouting 链（优先级 -250）

Netfilter 中 `NF_IP_PRI_CONNTRACK` 为 -200。prerouting 设为 -250，在连接跟踪之前执行，使报文进入 conntrack 前就能打上 `notrack` 并改写目的地址，不占用 conntrack 表项。

```nft
# 入向 FIP 查表跳转
add rule ip kube_ovn_stateless_nat prerouting iifname "$EXTERNAL_INTERFACE" ip daddr vmap @sl_fip_dnat

# 入向 DNAT 查表跳转
add rule ip kube_ovn_stateless_nat prerouting iifname "$EXTERNAL_INTERFACE" meta l4proto . th dport . ip daddr vmap @sl_dnat_in

# 出向 FIP：经 FIB 路由确认出接口为公网卡时打 notrack
add rule ip kube_ovn_stateless_nat prerouting iifname "$VPC_INTERFACE" fib daddr oifname "$EXTERNAL_INTERFACE" ip saddr @sl_notrack_ips notrack

# 出向 DNAT：经 FIB 路由确认出接口为公网卡时打 notrack
add rule ip kube_ovn_stateless_nat prerouting iifname "$VPC_INTERFACE" fib daddr oifname "$EXTERNAL_INTERFACE" meta l4proto . th sport . ip saddr @sl_notrack_tuples notrack
```

出向报文只有在 FIB 查表确认下一跳走 `$EXTERNAL_INTERFACE` 时才打 `notrack`。若内网 Pod 访问同一个网关绑定的 EIP（发卡流量），出接口仍是 `$VPC_INTERFACE`，规则不命中，报文继续走 iptables 发卡规则与 conntrack。

出向 DNAT 仅对 `@sl_notrack_tuples` 里记录的 `(proto, sport, saddr)` 打 `notrack`，同一 Pod 其他端口或协议的流量不受影响。

#### forward 链（优先级 mangle）

```nft
add rule ip kube_ovn_stateless_nat forward ct state untracked tcp flags & (syn | rst) == syn tcp option maxseg size set rt mtu
```

对 `ct state untracked` 的 TCP SYN 报文，将 MSS 钳制为出接口路由 MTU，避免 Geneve/STT/VXLAN 等 VPC 隧道封装导致报文超过 PMTU 而分片或丢包。

#### postrouting 链（优先级 99）

iptables POSTROUTING SNAT 优先级为 100。无状态 postrouting 取 99，在 iptables SNAT 之前改写源地址。

```nft
# VPC 出向 FIP 改写源地址
add rule ip kube_ovn_stateless_nat postrouting iifname "$VPC_INTERFACE" oifname "$EXTERNAL_INTERFACE" ct state untracked ip saddr vmap @sl_fip_snat

# VPC 出向 DNAT 反向改写源地址与源端口
add rule ip kube_ovn_stateless_nat postrouting iifname "$VPC_INTERFACE" oifname "$EXTERNAL_INTERFACE" ct state untracked meta l4proto . th sport . ip saddr vmap @sl_dnat_out
```

---

### 2.2 集合与字典定义

```nft
# 出向判定 notrack 的集合
add set ip kube_ovn_stateless_nat sl_notrack_ips     { type ipv4_addr; }
add set ip kube_ovn_stateless_nat sl_notrack_tuples  { type inet_proto . inet_service . ipv4_addr; }

# 查表分发至子链的字典
add map ip kube_ovn_stateless_nat sl_fip_dnat        { type ipv4_addr : verdict; }
add map ip kube_ovn_stateless_nat sl_fip_snat        { type ipv4_addr : verdict; }
add map ip kube_ovn_stateless_nat sl_dnat_in         { type inet_proto . inet_service . ipv4_addr : verdict; }
add map ip kube_ovn_stateless_nat sl_dnat_out        { type inet_proto . inet_service . ipv4_addr : verdict; }
```

---

### 2.3 映射子链

前述基础链、默认规则与全局 Set/Map 在网关初始化函数 `stateless_init` 中一次性建立；下述映射子链及 map 元素由 `stateless_fip_add / stateless_dnat_add` 动态添加（删除时对应 `stateless_fip_del / stateless_dnat_del`）。

#### FIP 映射子链

```nft
# 入向子链
add chain ip kube_ovn_stateless_nat sl_fin_192_0_2_10
add rule ip kube_ovn_stateless_nat sl_fin_192_0_2_10 notrack ip daddr set 10.0.1.20

# 出向子链
add chain ip kube_ovn_stateless_nat sl_fout_192_0_2_10
add rule ip kube_ovn_stateless_nat sl_fout_192_0_2_10 ip saddr set 192.0.2.10

# 关联元素
add element ip kube_ovn_stateless_nat sl_fip_dnat    { 192.0.2.10 : jump sl_fin_192_0_2_10 }
add element ip kube_ovn_stateless_nat sl_fip_snat    { 10.0.1.20 : jump sl_fout_192_0_2_10 }
add element ip kube_ovn_stateless_nat sl_notrack_ips { 10.0.1.20 }
```

#### DNAT 映射子链

```nft
# 入向子链
add chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80
add rule ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80 notrack ip daddr set 10.0.1.20 tcp dport set 8080

# 出向子链
add chain ip kube_ovn_stateless_nat sl_dout_tcp_192_0_2_10_80
add rule ip kube_ovn_stateless_nat sl_dout_tcp_192_0_2_10_80 ip saddr set 192.0.2.10 tcp sport set 80

# 关联元素
add element ip kube_ovn_stateless_nat sl_dnat_in         { tcp . 80 . 192.0.2.10 : jump sl_din_tcp_192_0_2_10_80 }
add element ip kube_ovn_stateless_nat sl_dnat_out        { tcp . 8080 . 10.0.1.20 : jump sl_dout_tcp_192_0_2_10_80 }
add element ip kube_ovn_stateless_nat sl_notrack_tuples  { tcp . 8080 . 10.0.1.20 }
```

#### 多后端 Service DNAT（五元组 Maglev 一致性哈希，`N = 127`）

```nft
# 各后端独立入向子链（按 ip:port 16 位 md5 短哈希命名，与下标解耦）
add chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80_ep_0efcb8350f3a691c
add rule ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80_ep_0efcb8350f3a691c notrack ip daddr set 10.0.1.10 tcp dport set 8080

add chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80_ep_0018b634e9d8d64f
add rule ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80_ep_0018b634e9d8d64f notrack ip daddr set 10.0.1.11 tcp dport set 8080

# 入向主链：五元组 jhash（固定 seed）mod 127 查 Maglev vmap
add chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80
add rule ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_10_80 jhash ip saddr . ip daddr . meta l4proto . th sport . th dport mod 127 seed 0x1a2b3c4d vmap { 0 : goto sl_din_tcp_192_0_2_10_80_ep_0efcb8350f3a691c, 1 : goto sl_din_tcp_192_0_2_10_80_ep_0018b634e9d8d64f, ..., 126 : goto sl_din_tcp_192_0_2_10_80_ep_0efcb8350f3a691c }

# 共享出向子链与多后端三元组汇聚
add chain ip kube_ovn_stateless_nat sl_dout_tcp_192_0_2_10_80
add rule ip kube_ovn_stateless_nat sl_dout_tcp_192_0_2_10_80 ip saddr set 192.0.2.10 tcp sport set 80

add element ip kube_ovn_stateless_nat sl_dnat_in         { tcp . 80 . 192.0.2.10 : jump sl_din_tcp_192_0_2_10_80 }
add element ip kube_ovn_stateless_nat sl_dnat_out        { tcp . 8080 . 10.0.1.10 : jump sl_dout_tcp_192_0_2_10_80, tcp . 8080 . 10.0.1.11 : jump sl_dout_tcp_192_0_2_10_80 }
add element ip kube_ovn_stateless_nat sl_notrack_tuples  { tcp . 8080 . 10.0.1.10, tcp . 8080 . 10.0.1.11 }
```

---

## 3. 共享与限速表：`kube-ovn`

### 3.1 Service 负载均衡（Share-DNAT）

```nft
add table ip kube-ovn
add chain ip kube-ovn prerouting { type nat hook prerouting priority -150; }
add map ip kube-ovn service-ips { type ipv4_addr . inet_proto . inet_service : verdict; }
add rule ip kube-ovn prerouting ip daddr . meta l4proto . th dport vmap @service-ips
```

针对 `dnat-<idhash>` 子链（`idhash = md5(eip:dport:proto)[0:12]`）：

1. 未开启会话保持（`affinity=none`）：
```nft
add rule ip kube-ovn dnat-5a1b2c3d4e5f meta l4proto tcp dnat ip addr . port to numgen random mod 2 map { 0 : 10.0.1.10 . 8080, 1 : 10.0.1.11 . 8080 }
```
通过 `numgen random mod N` 查 map 轮询选后端。

2. 开启客户端 IP 会话保持（`affinity=clientip`）：
```nft
# 调度链
add rule ip kube-ovn dnat-5a1b2c3d4e5f ip saddr @aff-5a1b2c3d4e5f-bkhash1 goto ep-5a1b2c3d4e5f-bkhash1
add rule ip kube-ovn dnat-5a1b2c3d4e5f ip saddr @aff-5a1b2c3d4e5f-bkhash2 goto ep-5a1b2c3d4e5f-bkhash2
add rule ip kube-ovn dnat-5a1b2c3d4e5f numgen random mod 2 vmap { 0 : goto ep-5a1b2c3d4e5f-bkhash1, 1 : goto ep-5a1b2c3d4e5f-bkhash2 }

# 后端 Endpoint 链与动态超时 Set
add set ip kube-ovn aff-5a1b2c3d4e5f-bkhash1 { type ipv4_addr; flags dynamic, timeout; timeout 10800s; }
add rule ip kube-ovn ep-5a1b2c3d4e5f-bkhash1 update @aff-5a1b2c3d4e5f-bkhash1 { ip saddr }
add rule ip kube-ovn ep-5a1b2c3d4e5f-bkhash1 meta l4proto tcp dnat to 10.0.1.10:8080
```
已知客户端由 `@aff-*` 集合直接跳至原后端子链；新客户端由 `numgen` 分配后写入集合，记录客户端 IP 并维护超时时间。

---

### 3.2 入口 SYN 限速（`syn-limit`）

```nft
add chain ip kube-ovn syn-limit { type filter hook prerouting priority -300; policy accept; }
add set ip kube-ovn syn_src_meter { type ipv4_addr; flags dynamic, timeout; size 262144; timeout 10s; }
add set ip kube-ovn syn_dst_meter { type ipv4_addr; flags dynamic, timeout; size 2048; timeout 10s; }

add rule ip kube-ovn syn-limit iifname != "$EXTERNAL_INTERFACE" return
add rule ip kube-ovn syn-limit tcp flags & (fin|syn|rst|ack) == syn update @syn_src_meter { ip saddr limit rate over 1000/second burst 2000 packets } drop
add rule ip kube-ovn syn-limit tcp flags & (fin|syn|rst|ack) == syn update @syn_dst_meter { ip daddr limit rate over 20000/second burst 40000 packets } drop
```

优先级 -300 等同 raw，在连接跟踪前执行。`syn_src_meter` 限制单源 IP 的 SYN 速率，`syn_dst_meter` 限制单个 EIP 的总 SYN 速率。

---

### 3.3 TCP 并发连接限制（`conn-limit`）

```nft
add chain ip kube-ovn conn-limit { type filter hook prerouting priority -150; policy accept; }
add rule ip kube-ovn conn-limit iifname != "$EXTERNAL_INTERFACE" return
add rule ip kube-ovn conn-limit tcp flags & (fin|syn|rst|ack) == syn ct count over 50000 drop
```

在 `stateless-nft` 模式下，流量已由 notrack 绕过 conntrack，此链无统计意义，控制器会自动清理该链规则。

---

## 4. 约束与边界

### 4.1 反向映射唯一性约束

无状态模式出向根据 `(protocol, internal_port, internal_ip)` 查 `sl_dnat_out`。同一个 Pod 的 `(protocol, port)` 在该网关内只能由一个公网身份发布。

若把同一个 Pod:port 同时绑定到两个 EIP（例如一个独占 DNAT 和一个 Service），Pod 回包会被 `sl_dnat_out` 改写为其中一个 EIP，导致发往另一个 EIP 的连接收到 RST。控制器在 `validateStatelessMapping` 时拦截并拒绝该配置。

### 4.2 多后端 Service 无状态与 Hairpin 分流

在 `stateless-nft` 模式下，公网 EIP Service 无论 1 个还是多个 Ready 后端均下发至 `kube_ovn_stateless_nat`（多后端采用固定 `N = 127` 五元组 Maglev 一致性哈希），0 后端时自动清理无状态规则。其中 `N = 127` 槽位按单个 `(EIP, ExternalPort, Protocol)` 对外端口独立分配，不同 Service 或同一 Service 的不同端口互不共享，单端口支持的最大 Ready 后端数为 127 个；超过 127 个后端时控制器会回退至 `table ip kube-ovn` 有状态路径，避免槽位溢出导致流量黑洞（后端数在 127↔128 临界点跨越切换时，因 `notrack` 与 `tracked` 状态变化，已建立的外部长连接会被中断）。

同时，`jhash` 使用由入向链名推导的固定 `seed`，以保证网关多副本及规则重放时流映射一致；其可接受的权衡是哈希值对外部可预测（若恶意构造特定源端口，存在将流量集中到单个后端的风险）。

若多后端 Service 配置了 `sessionAffinity: ClientIP`，由于五元组无状态哈希包含客户端源端口无法保证按源 IP 会话保持，控制器会跳过无状态路径（并清理可能残留的无状态规则），使其保留在 `table ip kube-ovn` 走有状态 `aff-*` 会话保持。

同时，多后端场景会保留 `table ip kube-ovn` 中的 Share-DNAT 表项以兜底 VPC 内部发卡（Hairpin）流量；ClusterIP 与 `lanIP` 流量不经过公网卡，严格保留在有状态路径，防止与公网 EIP 产生出向三元组反向映射冲突。

### 4.3 Maglev 槽位生成算法（`stateless_maglev_vmap_entries`）

[nat-gateway.sh](nat-gateway.sh) 用纯 POSIX / BusyBox `awk` 生成固定 127 个槽位的 Maglev 一致性哈希表。

因为 POSIX / BusyBox `awk` 没有内置 `ord()` 函数，`BEGIN` 块先构建了一张 ASCII `32..126` 的字符编码查表字典 `ord[]`，用来逐字符读取排序后的 `<ip>:<port>` 字符串，再计算两组独立哈希：

```awk
h1 = int((h1 * 33 + ch) % 2147483647)
h2 = int((h2 * 131 + ch + 17) % 2147483647)
```

`h1`（初值 `5381`，乘数 `33`）和 `h2`（初值 `0`，乘数 `131`、偏移 `17`）各自独立散列，减少不同后端出现相同探测序列的概率。`awk` 数值底层是 IEEE 754 双精度浮点数（精确整数上限为 $2^{53}$），每步对梅森质数 `2147483647`（$2^{31}-1$）取模后，中间值最大约 $2.81 \times 10^{11}$，不会出现精度丢失。

接着用这两组哈希算出每个后端的偏好序列参数：起始偏移 `offset = h1 % 127`（`0..126`），探测步长 `skip = (h2 % 126) + 1`（`1..126`）。由于 `127` 是质数，任意 `skip` 都与 `127` 互质，因此序列 `(offset + k * skip) mod 127`（$k = 0..126$）能够不重复地遍历全部 127 个槽位。

所有后端按字典序排序后，轮流沿各自的偏好序列占领下一个空槽（`slot[c] == -1`），直到 127 个槽位全部填满（`vmap` 不允许留空槽）。扩缩容时存活后端的 `offset` 和 `skip` 保持不变，只有约 $1/N$ 的槽位归属会发生变化。

### 4.4 分片与 ICMP

notrack 绕过了 conntrack defrag，UDP 大包分片中除首片外的报文不带传输层端口，无法匹配 DNAT 规则。

入向 ICMP unreachable / frag-needed 报文的外层 IP 会被改写，但内嵌载荷里的原始 IP 不做翻译。TCP 由 forward 链的 MSS 钳制避免分片，UDP 需由应用层控制包长不超过 MTU。

