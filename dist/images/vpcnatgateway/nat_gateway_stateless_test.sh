#!/usr/bin/env bash
set -Eeuo pipefail
script="$(cd "$(dirname "$0")" && pwd)/nat-gateway.sh"
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT
export NAT_GW_ENV_FILE="$tmp_dir/env"
eval "$(sed '/^opt=\$1$/,$d' "$script")"
# Capture external operations; no privileged network changes in unit tests.
# Mock nft binary to record invocations, simulate table listing, and inspect command batches.
nft() {
    printf '%s\n' "$*" >> "$tmp_dir/calls"
    if [[ "$*" == 'list tables' ]]; then
        printf 'table ip %s\n' ko_sl_fip_192_0_2_1 ko_sl_dnat_192_0_2_11_tcp_80 ko_sl_snat_192_0_2_2_10_0_0_4
    fi
    if [[ "$*" == '-f -' ]]; then cat >> "$tmp_dir/batches"; fi
}

# Mock hairpin iptables functions to verify they are invoked alongside nftables transactions.
add_floating_ip() { echo "fip $*" >> "$tmp_dir/hairpin"; }
del_floating_ip() { echo "del-fip $*" >> "$tmp_dir/hairpin"; }
add_dnat() { echo "dnat $*" >> "$tmp_dir/hairpin"; }
del_dnat() { echo "del-dnat $*" >> "$tmp_dir/hairpin"; }
check_inited() { :; }
stateless_fip_add '192.0.2.1,10.0.0.2'
grep -F 'ip daddr 192.0.2.1 notrack ip daddr set 10.0.0.2' "$tmp_dir/batches"
grep -F 'fib daddr oifname "net1" notrack' "$tmp_dir/batches"
grep -F 'ct state untracked ip saddr 10.0.0.2 ip saddr set 192.0.2.1' "$tmp_dir/batches"
grep -F 'fip 192.0.2.1,10.0.0.2' "$tmp_dir/hairpin"
# Replaying the mapping flushes only its own table, without duplicate rules.
stateless_fip_add '192.0.2.1,10.0.0.2'
[[ $(grep -c 'flush table ip ko_sl_fip_192_0_2_1' "$tmp_dir/batches") == 2 ]]
stateless_fip_del '192.0.2.1'
grep -Fx 'delete table ip ko_sl_fip_192_0_2_1' "$tmp_dir/calls"
stateless_dnat_add '192.0.2.11,80,tcp,10.0.0.2,8080'
grep -F 'tcp dport 80 notrack ip daddr set 10.0.0.2 tcp dport set 8080' "$tmp_dir/batches"
grep -F 'tcp sport 8080 ip saddr set 192.0.2.11 tcp sport set 80' "$tmp_dir/batches"
stateless_dnat_del '192.0.2.11,80,tcp'
grep -Fx 'delete table ip ko_sl_dnat_192_0_2_11_tcp_80' "$tmp_dir/calls"
stateless_snat_add '192.0.2.2,10.0.0.4/32'
stateless_snat_del '192.0.2.2,10.0.0.4'
grep -Fx 'delete table ip ko_sl_snat_192_0_2_2_10_0_0_4' "$tmp_dir/calls"
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
echo 'stateless shell tests passed'
