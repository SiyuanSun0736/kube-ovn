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
 * 功能：Kube-OVN 公网 EIP eBPF TC Ingress 旁路直通分类器。
 * 作用：在物理网卡驱动层命中 EIP 时直接短路重定向推入本地网关 Pod，绕过宿主机 Netfilter conntrack 规避连接表被击穿；非 EIP 流量正常放行。
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
    __u32 parent_ifindex;             // macvlan master ifindex; also aligns the following counters
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

static __always_inline int mac_eq(const __u8 *a, const __u8 *b) {
    #pragma unroll
    for (int i = 0; i < ETH_ALEN; i++) {
        if (a[i] != b[i]) {
            return 0;
        }
    }
    return 1;
}

// macvlan bridge egress reinjects local traffic through the master's ingress.
// Once the destination MAC is rewritten, let its RX handler deliver to the Pod
// instead of redirecting the same packet back through the macvlan indefinitely.
static __always_inline int is_local_delivery(struct __sk_buff *skb, void *data,
                                             struct eip_target *target) {
    struct ethhdr *eth = data;
    return skb->ifindex == target->parent_ifindex &&
           mac_eq(eth->h_dest, target->target_mac);
}

// In direct-action mode TC_ACT_OK ends classification; it does not continue
// to a lower-priority Cilium filter. Unhandled packets must return UNSPEC.
// OK is reserved below for matched EIPs already on the local delivery path.
// 处理 IPv4 报文逻辑
static __always_inline int process_ipv4(struct __sk_buff *skb, void *data, void *data_end) {
    struct iphdr *iph = data + sizeof(struct ethhdr);
    // 边界检查，防止越界访问引发内核校验器拒绝
    if ((void *)(iph + 1) > data_end) {
        return TC_ACT_UNSPEC;
    }

    __u32 daddr = iph->daddr;
    // 检索 BPF Map
    struct eip_target *target = bpf_map_lookup_elem(&kube_ovn_eip_v4_map, &daddr);
    if (!target || !target->ifindex) {
        // 未命中本节点承载的 EIP：原样放行给后续的 Cilium 探针及宿主机内核协议栈
        return TC_ACT_UNSPEC;
    }

    if (is_local_delivery(skb, data, target)) {
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

    // 内核短路直通转发：从宿主机 macvlan 的 Egress 发出（flags=0），由目标 Pod MAC 选择接收端
    // 彻底避开宿主机根空间 Netfilter 与 conntrack 表
    return bpf_redirect(target->ifindex, 0);
}

// 处理 IPv6 报文逻辑
static __always_inline int process_ipv6(struct __sk_buff *skb, void *data, void *data_end) {
    struct ipv6hdr *ip6h = data + sizeof(struct ethhdr);
    if ((void *)(ip6h + 1) > data_end) {
        return TC_ACT_UNSPEC;
    }

    struct in6_addr daddr = ip6h->daddr;
    struct eip_target *target = bpf_map_lookup_elem(&kube_ovn_eip_v6_map, &daddr);
    if (!target || !target->ifindex) {
        return TC_ACT_UNSPEC;
    }

    if (is_local_delivery(skb, data, target)) {
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
        return TC_ACT_UNSPEC;
    }

    __u16 h_proto = bpf_ntohs(eth->h_proto);

    // 遇到 VLAN Tag 流量时放行，保持与底层网络封装兼容
    if (h_proto == ETH_P_8021Q || h_proto == ETH_P_8021AD) {
        return TC_ACT_UNSPEC;
    }

    if (h_proto == ETH_P_IP) {
        return process_ipv4(skb, data, data_end);
    } else if (h_proto == ETH_P_IPV6) {
        return process_ipv6(skb, data, data_end);
    }

    // 其它未知协议报文默认放行
    return TC_ACT_UNSPEC;
}

char _license[] SEC("license") = "Apache-2.0";
