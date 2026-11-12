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
 * 拓扑架构与分工：
 *   - Kube-OVN TC：只占用 macbond1（对外发布公网 BGP EIP 路由的接口），专职负责公网 EIP 短路直通与防 conntrack 击穿；
 *   - Cilium TC：只占用 bond0（集群管理通信）+ 各个 Pod 虚拟网卡 *_h（容器侧安全策略）。
 *   两者在网络接口上严格物理隔离、各司其职，互不干扰。
 * 作用：在 macbond1 Ingress 命中 EIP 时直接短路重定向推入本地网关 Pod，绕过宿主机 Netfilter conntrack 规避连接表被击穿；非 EIP 流量正常放行。
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

// 重定向目标信息结构体（内存严格按 32 位对齐，总大小为 16 字节）
struct eip_target {
    __u32 ifindex;          // 本地网关 Pod 对应的宿主机接口设备索引（macvlan 或 veth）
    __u8  target_mac[6];    // 目标网卡硬件 MAC 地址（若非 0 则执行二层 MAC 重写）
    __u8  pad[2];           // 2 字节填充，对齐到 4 字节
    __u32 parent_ifindex;   // macvlan 父设备接口索引（Master ifindex）
};
_Static_assert(sizeof(struct eip_target) == 16, "struct eip_target size must be 16 bytes");

// VLAN 报文头结构体（IEEE 802.1Q / 802.1ad，共 4 字节）
struct vlan_hdr {
    __be16 h_vlan_TCI;                  // 优先权、CFI 与 VLAN 标识符 (VID)
    __be16 h_vlan_encapsulated_proto;   // 内层封装的上层协议（如 ETH_P_IP 或内层 802.1Q）
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

// macvlan bridge 模式下，Egress 本地流量可能会重新注入父设备的 Ingress 接收队列。
// 一旦目的 MAC 已经被改写为目标 Pod MAC，直接放行交由内核驱动 RX 处理并交付给 Pod，
// 避免相同数据包在 macvlan 与物理网卡之间发生无限死循环重定向。
static __always_inline int is_local_delivery(struct __sk_buff *skb, void *data,
                                             struct eip_target *target) {
    struct ethhdr *eth = data;
    return target->parent_ifindex &&
           skb->ifindex == target->parent_ifindex &&
           mac_eq(eth->h_dest, target->target_mac);
}

// 返回值与动作说明：
// - Kube-OVN TC 只挂载在 macbond1，Cilium TC 只挂载在 bond0 及各个 Pod 虚拟网卡 *_h；
// - TC_ACT_UNSPEC (-1)：未命中 EIP 的常规网络流量（如节点管理、BGP 协商）原样放行，交付宿主机 Linux 内核协议栈正常处理；
// - TC_ACT_OK (0)：专用于“已命中 EIP 且正处于 macvlan 本地交付链路”的报文，直接终止 TC 处理交由内核驱动 RX 处理，规避重定向死循环；
// - TC_ACT_SHOT (2)：仅在底层严重异常（如 bpf_skb_store_bytes 越界失败）时安全丢包。
// 处理 IPv4 报文逻辑
static __always_inline int process_ipv4(struct __sk_buff *skb, void *data, void *data_end, __u32 l3_offset) {
    struct iphdr *iph = data + l3_offset;
    // 边界检查，防止越界访问引发内核校验器拒绝
    if ((void *)(iph + 1) > data_end) {
        return TC_ACT_UNSPEC;
    }

    __u32 daddr = iph->daddr;
    // 检索 BPF Map
    struct eip_target *target = bpf_map_lookup_elem(&kube_ovn_eip_v4_map, &daddr);
    if (!target || !target->ifindex) {
        // 未命中本节点承载的 EIP：原样放行交付宿主机内核协议栈正常处理
        return TC_ACT_UNSPEC;
    }

    if (is_local_delivery(skb, data, target)) {
        return TC_ACT_OK;
    }

    // 若配置了目标容器网卡的 MAC 地址，则原子改写以太网目的 MAC（外层 MAC 偏移为 0）
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
static __always_inline int process_ipv6(struct __sk_buff *skb, void *data, void *data_end, __u32 l3_offset) {
    struct ipv6hdr *ip6h = data + l3_offset;
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
    __u32 l3_offset = sizeof(struct ethhdr);

    // 支持解析 802.1Q 与 802.1ad VLAN 标签（包括单层与双层 QinQ），兼容外网/宿主机带 Tag 的网络流量
    #pragma unroll
    for (int i = 0; i < 2; i++) {
        if (h_proto == ETH_P_8021Q || h_proto == ETH_P_8021AD) {
            struct vlan_hdr *vhdr = data + l3_offset;
            if ((void *)(vhdr + 1) > data_end) {
                return TC_ACT_UNSPEC;
            }
            h_proto = bpf_ntohs(vhdr->h_vlan_encapsulated_proto);
            l3_offset += sizeof(struct vlan_hdr);
        }
    }

    if (h_proto == ETH_P_IP) {
        return process_ipv4(skb, data, data_end, l3_offset);
    } else if (h_proto == ETH_P_IPV6) {
        return process_ipv6(skb, data, data_end, l3_offset);
    }

    // 其它未知协议报文默认放行
    return TC_ACT_UNSPEC;
}

char _license[] SEC("license") = "Apache-2.0";
