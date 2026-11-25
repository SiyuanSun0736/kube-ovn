#!/usr/bin/env bash
set -Eeuo pipefail
script="$(cd "$(dirname "$0")" && pwd)/nat-gateway.sh"
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT
export NAT_GW_ENV_FILE="$tmp_dir/env"
eval "$(sed '/^opt=\$1$/,$d' "$script")"
# 拦截外部操作：单元测试中不触发特权网络命名空间改动
nft() {
    printf '%s\n' "$*" >> "$tmp_dir/calls"
    if [[ "$*" == '-f -' ]]; then cat >> "$tmp_dir/batches"; fi
}

# Mock 发卡 iptables 函数
add_floating_ip() { echo "fip $*" >> "$tmp_dir/hairpin"; }
del_floating_ip() { echo "del-fip $*" >> "$tmp_dir/hairpin"; }
add_dnat() { echo "dnat $*" >> "$tmp_dir/hairpin"; }
del_dnat() { echo "del-dnat $*" >> "$tmp_dir/hairpin"; }
check_inited() { :; }

stateless_init
grep -F 'add map ip kube_ovn_stateless_nat sl_fip_dnat' "$tmp_dir/batches"
grep -F 'add map ip kube_ovn_stateless_nat sl_dnat_in' "$tmp_dir/batches"
grep -F 'add set ip kube_ovn_stateless_nat sl_notrack_ips' "$tmp_dir/batches"
grep -F 'add set ip kube_ovn_stateless_nat sl_notrack_tuples' "$tmp_dir/batches"
grep -F 'ip saddr @sl_notrack_ips notrack' "$tmp_dir/batches"

stateless_fip_add '192.0.2.1,10.0.0.2'
grep -F 'add element ip kube_ovn_stateless_nat sl_fip_dnat { 192.0.2.1 : jump sl_fin_192_0_2_1 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_fip_snat { 10.0.0.2 : jump sl_fout_192_0_2_1 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_notrack_ips { 10.0.0.2 }' "$tmp_dir/batches"
grep -F 'fip 192.0.2.1,10.0.0.2' "$tmp_dir/hairpin"

stateless_fip_del '192.0.2.1'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_fip_dnat { 192.0.2.1 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_fin_192_0_2_1' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_fout_192_0_2_1' "$tmp_dir/calls"

stateless_fip_add '192.0.2.1,10.0.0.2'
stateless_fip_del '192.0.2.1,10.0.0.2'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_fip_snat { 10.0.0.2 }' "$tmp_dir/calls"
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_notrack_ips { 10.0.0.2 }' "$tmp_dir/calls"

stateless_dnat_add '192.0.2.11,80,tcp,10.0.0.2,8080'
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_in { tcp . 80 . 192.0.2.11 : jump sl_din_tcp_192_0_2_11_80 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.2 : jump sl_dout_tcp_192_0_2_11_80 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_notrack_tuples { tcp . 8080 . 10.0.0.2 }' "$tmp_dir/batches"

stateless_dnat_del '192.0.2.11,80,tcp'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_in { tcp . 80 . 192.0.2.11 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_11_80' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_dout_tcp_192_0_2_11_80' "$tmp_dir/calls"

stateless_dnat_add '192.0.2.11,80,tcp,10.0.0.2,8080'
stateless_dnat_del '192.0.2.11,80,tcp,10.0.0.2,8080'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.2 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_dout_tcp_192_0_2_11_80' "$tmp_dir/calls"

# 大写协议标准化处理断言
stateless_dnat_add '192.0.2.12,80,TCP,10.0.0.3,8080'
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_in { tcp . 80 . 192.0.2.12 : jump sl_din_tcp_192_0_2_12_80 }' "$tmp_dir/batches"
stateless_dnat_del '192.0.2.12,80,TCP'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_in { tcp . 80 . 192.0.2.12 }' "$tmp_dir/calls"

# 多后端五元组一致性哈希（N=127 Maglev）与级联清理测试
: > "$tmp_dir/batches"
: > "$tmp_dir/calls"
stateless_dnat_add '192.0.2.15,80,tcp,10.0.0.2:8080@10.0.0.3:8080'
bk1=$(nft_backend_hash 10.0.0.2 8080)
bk2=$(nft_backend_hash 10.0.0.3 8080)
grep -F "add chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk1}" "$tmp_dir/batches"
grep -F "add chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk2}" "$tmp_dir/batches"
expected_seed="0x$(echo -n 'sl_din_tcp_192_0_2_15_80' | md5sum | cut -c1-8)"
grep -F "jhash ip saddr . ip daddr . meta l4proto . th sport . th dport mod 127 seed ${expected_seed} vmap {" "$tmp_dir/batches"
grep -F '0 : goto sl_din_tcp_192_0_2_15_80_ep_' "$tmp_dir/batches"
grep -F '126 : goto sl_din_tcp_192_0_2_15_80_ep_' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.2 : jump sl_dout_tcp_192_0_2_15_80 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.3 : jump sl_dout_tcp_192_0_2_15_80 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_notrack_tuples { tcp . 8080 . 10.0.0.2 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_notrack_tuples { tcp . 8080 . 10.0.0.3 }' "$tmp_dir/batches"

# 模拟 nft -nn list chain 返回多后端主链与子链，验证缩容到单后端及全删时的精确清理
nft() {
    printf '%s\n' "$*" >> "$tmp_dir/calls"
    if [[ "$*" == '-f -' ]]; then
        cat >> "$tmp_dir/batches"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80" ]]; then
        echo "jhash ip saddr . ip daddr . meta l4proto . th sport . th dport mod 127 vmap { 0 : goto sl_din_tcp_192_0_2_15_80_ep_${bk1}, 1 : goto sl_din_tcp_192_0_2_15_80_ep_${bk2} }"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk1}" ]]; then
        echo "notrack ip daddr set 10.0.0.2 tcp dport set 8080"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk2}" ]]; then
        echo "notrack ip daddr set 10.0.0.3 tcp dport set 8080"
        return 0
    fi
}
: > "$tmp_dir/batches"
: > "$tmp_dir/calls"
: > "$tmp_dir/hairpin"
stateless_dnat_add '192.0.2.15,80,tcp,10.0.0.2:8080'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.3 }' "$tmp_dir/calls"
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_notrack_tuples { tcp . 8080 . 10.0.0.3 }' "$tmp_dir/calls"
grep -Fx "delete chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk1}" "$tmp_dir/calls"
grep -Fx "delete chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk2}" "$tmp_dir/calls"
grep -Fx 'del-dnat 192.0.2.15,80,tcp' "$tmp_dir/hairpin"
grep -Fx 'dnat 192.0.2.15,80,tcp,10.0.0.2,8080' "$tmp_dir/hairpin"

# 验证从单后端扩容到多后端时清理旧单后端 Hairpin DNAT 规则
nft() {
    printf '%s\n' "$*" >> "$tmp_dir/calls"
    if [[ "$*" == '-f -' ]]; then
        cat >> "$tmp_dir/batches"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80" ]]; then
        echo "notrack ip daddr set 10.0.0.2 tcp dport set 8080"
        return 0
    fi
}
: > "$tmp_dir/hairpin"
stateless_dnat_add '192.0.2.15,80,tcp,10.0.0.2:8080@10.0.0.3:8080'
grep -Fx 'del-dnat 192.0.2.15,80,tcp' "$tmp_dir/hairpin"
if grep -q '^dnat ' "$tmp_dir/hairpin"; then
    echo "multi-backend should not add single-backend hairpin dnat" >&2
    exit 1
fi

# 恢复多后端链 mock，验证全量删除级联清理
nft() {
    printf '%s\n' "$*" >> "$tmp_dir/calls"
    if [[ "$*" == '-f -' ]]; then
        cat >> "$tmp_dir/batches"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80" ]]; then
        echo "jhash ip saddr . ip daddr . meta l4proto . th sport . th dport mod 127 seed 0x6f907617 vmap { 0 : goto sl_din_tcp_192_0_2_15_80_ep_${bk1}, 1 : goto sl_din_tcp_192_0_2_15_80_ep_${bk2} }"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk1}" ]]; then
        echo "notrack ip daddr set 10.0.0.2 tcp dport set 8080"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk2}" ]]; then
        echo "notrack ip daddr set 10.0.0.3 tcp dport set 8080"
        return 0
    fi
}

: > "$tmp_dir/calls"
stateless_dnat_del '192.0.2.15,80,tcp'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_in { tcp . 80 . 192.0.2.15 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80' "$tmp_dir/calls"
grep -Fx "delete chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk1}" "$tmp_dir/calls"
grep -Fx "delete chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_15_80_ep_${bk2}" "$tmp_dir/calls"
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.2 }' "$tmp_dir/calls"
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.3 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_dout_tcp_192_0_2_15_80' "$tmp_dir/calls"

# 验证单后端变更与幂等行为：变更时先 del_dnat 再 add_dnat，未变更时不触发 del_dnat
nft() {
    printf '%s\n' "$*" >> "$tmp_dir/calls"
    if [[ "$*" == '-f -' ]]; then
        cat >> "$tmp_dir/batches"
        return 0
    fi
    if [[ "$*" == "-nn list chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_16_80" ]]; then
        echo "notrack ip daddr set 10.0.0.2 tcp dport set 8080"
        return 0
    fi
}
: > "$tmp_dir/hairpin"
stateless_dnat_add '192.0.2.16,80,tcp,10.0.0.2,8080'
if grep -q '^del-dnat ' "$tmp_dir/hairpin"; then
    echo "unchanged single backend should not trigger del-dnat" >&2
    exit 1
fi
grep -Fx 'dnat 192.0.2.16,80,tcp,10.0.0.2,8080' "$tmp_dir/hairpin"

: > "$tmp_dir/hairpin"
stateless_dnat_add '192.0.2.16,80,tcp,10.0.0.4,8080'
grep -Fx 'del-dnat 192.0.2.16,80,tcp' "$tmp_dir/hairpin"
grep -Fx 'dnat 192.0.2.16,80,tcp,10.0.0.4,8080' "$tmp_dir/hairpin"

# Invalid deletion inputs must succeed as no-op and never block cleanup
stateless_dnat_del '192.0.2.11,080,tcp'
for invalid in '192.0.2.2,::1' '999.1.1.1,10.0.0.1'; do
    if stateless_fip_add "$invalid"; then echo "accepted invalid mapping $invalid" >&2; exit 1; fi
done
if stateless_dnat_add '192.0.2.1,80,sctp,10.0.0.2,8080'; then exit 1; fi
# Installation failures must reach the controller; never report success after a failed batch.
nft() { return 1; }
if stateless_fip_add '192.0.2.1,10.0.0.2'; then exit 1; fi

# nft_del_safe 必须正确保留非零退出码
nft() { echo "Permission denied" >&2; return 13; }
set +e
nft_del_safe "delete chain ip kube_ovn_stateless_nat sl_test"
del_err=$?
set -e
if [[ $del_err -ne 13 ]]; then echo "nft_del_safe failed to capture exit code: $del_err" >&2; exit 1; fi


echo 'stateless shell tests passed'
