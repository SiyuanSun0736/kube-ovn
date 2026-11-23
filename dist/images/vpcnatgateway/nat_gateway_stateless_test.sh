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
grep -F 'add element ip kube_ovn_stateless_nat sl_fip_snat { 10.0.0.2 : jump sl_fout_10_0_0_2 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_notrack_ips { 10.0.0.2 }' "$tmp_dir/batches"
grep -F '@th,160,32 set 0x0a000002' "$tmp_dir/batches"
grep -F '@th,192,32 set 0xc0000201' "$tmp_dir/batches"
grep -F 'fip 192.0.2.1,10.0.0.2' "$tmp_dir/hairpin"

stateless_fip_del '192.0.2.1'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_fip_dnat { 192.0.2.1 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_fin_192_0_2_1' "$tmp_dir/calls"

stateless_dnat_add '192.0.2.11,80,tcp,10.0.0.2,8080'
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_in { tcp . 80 . 192.0.2.11 : jump sl_din_tcp_192_0_2_11_80 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.2 : jump sl_dout_tcp_10_0_0_2_8080 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_notrack_tuples { tcp . 8080 . 10.0.0.2 }' "$tmp_dir/batches"
grep -F '@th,136,8 6 @th,224,16 80 notrack ip daddr set 10.0.0.2 @th,160,32 set 0x0a000002 @th,224,16 set 8080' "$tmp_dir/batches"
grep -F '@th,136,8 6 @th,240,16 8080 notrack ip saddr set 192.0.2.11 @th,192,32 set 0xc000020b @th,240,16 set 80' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_dnat_icmp_out { 10.0.0.2 : jump sl_docmp_10_0_0_2 }' "$tmp_dir/batches"

stateless_dnat_del '192.0.2.11,80,tcp'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_in { tcp . 80 . 192.0.2.11 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_din_tcp_192_0_2_11_80' "$tmp_dir/calls"

stateless_dnat_add '192.0.2.11,80,tcp,10.0.0.2,8080'
stateless_dnat_del '192.0.2.11,80,tcp,10.0.0.2,8080'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_dnat_out { tcp . 8080 . 10.0.0.2 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_dout_tcp_10_0_0_2_8080' "$tmp_dir/calls"

stateless_snat_add '192.0.2.2,10.0.0.4/32'
grep -F 'add element ip kube_ovn_stateless_nat sl_snat_in { 192.0.2.2 : jump sl_sin_192_0_2_2 }' "$tmp_dir/batches"
grep -F 'add element ip kube_ovn_stateless_nat sl_snat_out { 10.0.0.4 : jump sl_sout_10_0_0_4 }' "$tmp_dir/batches"
grep -F '@th,192,32 set 0xc0000202' "$tmp_dir/batches"
grep -F 'tcp flags & (syn | rst | ack) == syn drop' "$tmp_dir/batches"
grep -F 'icmp type echo-request drop' "$tmp_dir/batches"

stateless_snat_del '192.0.2.2,10.0.0.4'
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_snat_out { 10.0.0.4 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_sout_10_0_0_4' "$tmp_dir/calls"
grep -Fx 'delete element ip kube_ovn_stateless_nat sl_snat_in { 192.0.2.2 }' "$tmp_dir/calls"
grep -Fx 'delete chain ip kube_ovn_stateless_nat sl_sin_192_0_2_2' "$tmp_dir/calls"

# Invalid deletion inputs must succeed as no-op and never block cleanup
stateless_snat_del '192.0.2.2,10.0.0.0/24'
stateless_dnat_del '192.0.2.11,080,tcp'
for invalid in '192.0.2.2,10.0.0.0/24' '192.0.2.2,::1' '999.1.1.1,10.0.0.1'; do
    if stateless_snat_add "$invalid"; then echo "accepted invalid mapping $invalid" >&2; exit 1; fi
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

# nft_find_map_key_by_target 单行多元素测试
nft() {
    echo 'table ip kube_ovn_stateless_nat {
	map sl_fip_snat {
		type ipv4_addr : verdict
		elements = { 10.0.0.2 : jump sl_fout_10_0_0_2, 10.0.0.3 : jump sl_fout_10_0_0_3 }
	}
}'
}
found_key=$(nft_find_map_key_by_target sl_fip_snat "sl_fout_10_0_0_3")
if [[ "$found_key" != "10.0.0.3" ]]; then echo "nft_find_map_key_by_target expected 10.0.0.3, got $found_key" >&2; exit 1; fi

# nft_del_dnat_icmp_rules 严格 Token 匹配测试（防 80 误匹配 1280 或 8080）
nft() {
    if [[ "$1" == "-a" ]]; then
        cat << 'EOF'
table ip kube_ovn_stateless_nat {
	chain sl_dicmp_192_0_2_1 {
		add rule ip kube_ovn_stateless_nat sl_dicmp_192_0_2_1 @th,136,8 6 @th,224,16 80 notrack # handle 3
		add rule ip kube_ovn_stateless_nat sl_dicmp_192_0_2_1 @th,136,8 6 @th,224,16 0x50 notrack # handle 4
		add rule ip kube_ovn_stateless_nat sl_dicmp_192_0_2_1 @th,136,8 6 @th,224,16 8080 notrack # handle 5
		add rule ip kube_ovn_stateless_nat sl_dicmp_192_0_2_1 @th,136,8 6 @th,224,16 1280 notrack # handle 6
		add rule ip kube_ovn_stateless_nat sl_dicmp_192_0_2_1 @th,136,8 6 @th,224,16 0x500 notrack # handle 7
	}
}
EOF
    elif [[ "$1" == delete* ]]; then
        printf '%s\n' "$*" >> "$tmp_dir/del_calls"
    fi
}
nft_del_dnat_icmp_rules "sl_dicmp_192_0_2_1" 6 80 224
if ! grep -q "handle 3" "$tmp_dir/del_calls" || ! grep -q "handle 4" "$tmp_dir/del_calls"; then
    echo "nft_del_dnat_icmp_rules failed to delete handles 3 and 4: $(cat "$tmp_dir/del_calls")" >&2
    exit 1
fi
if grep -qE "handle (5|6|7)" "$tmp_dir/del_calls"; then
    echo "nft_del_dnat_icmp_rules mistakenly deleted handles 5/6/7: $(cat "$tmp_dir/del_calls")" >&2
    exit 1
fi

echo 'stateless shell tests passed'
