// +build ignore
// Copyright 2026 Kube-OVN Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
 * Kube-OVN BGP 动态联动 eBPF TC 旁路路由加速与宿主机 conntrack 抗穿透过滤器
 *
 * 【背景与痛点】
 * 在机房 BGP 架构中，外部公网对租户 EIP 的流量首先进入物理宿主机网卡（如 bond0）。
 * 由于物理网卡位于宿主机根网络空间（init_net），默认开启 Netfilter conntrack 跟踪。
 * 当遭遇公网高并发或 SYN 洪峰（如常山2故障暴露的 400W+ 僵尸连接）时，宿主机全局 conntrack
 * 表瞬间耗尽，导致同物理机上的其它 Pod、K8s 管理网络以及 SSH 运维彻底失联瘫痪。
 *
 * 【本程序核心设计与机制】
 * 1. 本 eBPF TC 分类器挂载在物理宿主机网卡的 clsact ingress 队列上，设置最高优先级 priority 1，
 *    以 direct-action (da) 模式运行。
 * 2. 外部报文在驱动收包最前端即被本程序捕获：
 *    - 提取目的 IP 地址并在内核 BPF Hash Map (kube_ovn_eip_v4_map / kube_ovn_eip_v6_map) 中检索；
 *    - 若命中 EIP：表明该流量目的为本地承载的 vpc-nat-gw 网关 Pod，直接调用内核短路原语
 *      bpf_redirect(target_ifindex, 0)，以 TC_ACT_REDIRECT 动作直接推入网关 Pod 网卡；
 *      【核心价值】：数据包 100% 绕过宿主机整个 Netfilter 与 conntrack 模块，宿主机表项增量为 0！
 *    - 若未命中 EIP：表明为集群节点间管理通信、NodePort、宿主机自身 SSH 或普通 Pod 流量，
 *      返回 TC_ACT_OK，报文原样穿透给后面的 Cilium TC 探针及原生 Linux 内核网络协议栈，
 *      与 Cilium NetworkPolicy、Hubble 监控完全共生且互不影响。
 */

#include <linux/bpf.h>
#include <linux/pkt_cls.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/ipv6.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

// BPF Map 最大支持的 EIP 记录条目数
#define MAX_EIP_ENTRIES 65536

// 重定向目标信息结构体（内存严格按 64 位对齐，总大小为 32 字节）
struct eip_target {
    __u32 ifindex;          // 本地网关 Pod 对应的宿主机接口设备索引（macvlan 或 veth）
    __u8  target_mac[6];    // 目标网卡硬件 MAC 地址（若非 0 则执行二层 MAC 重写）
    __u8  pad[2];           // 2 字节填充，对齐到 4 字节
    __u32 pad2;             // 4 字节显式填充，使后方的 64 位整型严格按 8 字节边界对齐
    __u64 packets;          // 命中直通的累计数据包计数
    __u64 bytes;            // 命中直通的累计字节统计
};

// IPv4 EIP ➔ 本地网关网卡索引与统计映射表
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_EIP_ENTRIES);
    __type(key, __u32);                 // IPv4 目的地址（网络字节序）
    __type(value, struct eip_target);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} kube_ovn_eip_v4_map SEC(".maps");

// IPv6 EIP ➔ 本地网关网卡索引与统计映射表
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_EIP_ENTRIES);
    __type(key, struct in6_addr);       // IPv6 目的地址
    __type(value, struct eip_target);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} kube_ovn_eip_v6_map SEC(".maps");

// 处理 IPv4 报文逻辑
static __always_inline int process_ipv4(struct __sk_buff *skb, void *data, void *data_end) {
    struct iphdr *iph = data + sizeof(struct ethhdr);
    // 边界检查，防止越界访问引发内核校验器拒绝
    if ((void *)(iph + 1) > data_end) {
        return TC_ACT_OK;
    }

    __u32 daddr = iph->daddr;
    // 检索 BPF Map
    struct eip_target *target = bpf_map_lookup_elem(&kube_ovn_eip_v4_map, &daddr);
    if (!target) {
        // 未命中本节点承载的 EIP：原样放行给后续的 Cilium 探针及宿主机内核协议栈
        return TC_ACT_OK;
    }

    // 原子递增转发统计计数
    __sync_fetch_and_add(&target->packets, 1);
    __sync_fetch_and_add(&target->bytes, skb->len);

    // 若配置了目标容器网卡的 MAC 地址，则原子改写以太网目的 MAC
    if (target->target_mac[0] | target->target_mac[1] | target->target_mac[2] |
        target->target_mac[3] | target->target_mac[4] | target->target_mac[5]) {
        if (bpf_skb_store_bytes(skb, offsetof(struct ethhdr, h_dest), target->target_mac, 6, 0) < 0) {
            return TC_ACT_SHOT;
        }
    }

    // 内核短路直通转发：直接推入网关 Pod 网卡的 Ingress 队列（flags=0）
    // 彻底避开宿主机根空间 Netfilter 与 conntrack 表
    return bpf_redirect(target->ifindex, 0);
}

// 处理 IPv6 报文逻辑
static __always_inline int process_ipv6(struct __sk_buff *skb, void *data, void *data_end) {
    struct ipv6hdr *ip6h = data + sizeof(struct ethhdr);
    if ((void *)(ip6h + 1) > data_end) {
        return TC_ACT_OK;
    }

    struct in6_addr daddr = ip6h->daddr;
    struct eip_target *target = bpf_map_lookup_elem(&kube_ovn_eip_v6_map, &daddr);
    if (!target) {
        return TC_ACT_OK;
    }

    __sync_fetch_and_add(&target->packets, 1);
    __sync_fetch_and_add(&target->bytes, skb->len);

    if (target->target_mac[0] | target->target_mac[1] | target->target_mac[2] |
        target->target_mac[3] | target->target_mac[4] | target->target_mac[5]) {
        if (bpf_skb_store_bytes(skb, offsetof(struct ethhdr, h_dest), target->target_mac, 6, 0) < 0) {
            return TC_ACT_SHOT;
        }
    }

    return bpf_redirect(target->ifindex, 0);
}

// TC 分类器入口钩子函数
SEC("classifier")
int tc_eip_ingress(struct __sk_buff *skb) {
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;

    // 解析以太网头
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) {
        return TC_ACT_OK;
    }

    __u16 h_proto = bpf_ntohs(eth->h_proto);

    // 遇到 VLAN Tag 流量时放行，保持与底层网络封装兼容
    if (h_proto == ETH_P_8021Q || h_proto == ETH_P_8021AD) {
        return TC_ACT_OK;
    }

    if (h_proto == ETH_P_IP) {
        return process_ipv4(skb, data, data_end);
    } else if (h_proto == ETH_P_IPV6) {
        return process_ipv6(skb, data, data_end);
    }

    // 其它未知协议报文默认放行
    return TC_ACT_OK;
}

char _license[] SEC("license") = "Apache-2.0";
