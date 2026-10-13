#!/usr/bin/env bash
set -Eeuo pipefail

script="$(cd "$(dirname "$0")" && pwd)/nat-gateway.sh"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

# Keep the script from reading an interface file left behind on the host.
export NAT_GW_ENV_FILE="$tmp_dir/nat-gateway.env"

# Load definitions without executing the command dispatcher at the end.
eval "$(sed '/^opt=\$1$/,$d' "$script")"

# A negative assertion cannot be written as "! cmd": POSIX has set -e ignore any command whose
# return value is inverted with !, so such a line never fails the test, however wrong the code is.
# refute runs the command in a subshell, which also contains the exit of a function under test.
refute() {
    if ( "$@" ) >/dev/null 2>&1; then
        echo "expected to fail, but it succeeded: $*" >&2
        exit 1
    fi
}

VPC_INTERFACE=eth0
EXTERNAL_INTERFACE=net1
vpc_addr=10.0.7.254
vpc_has_addr=true

iptables_cmd=iptables
iptables() {
    # -S is what current_hairpin_fallback_state reads; everything else goes to the chain stub.
    if [[ "$*" == "-t nat -S HAIRPIN_SNAT" ]]; then sed 's/^/-A HAIRPIN_SNAT /' "$hairpin_state"; return 0; fi
    ipt_stub "$@"
}
inited_calls=0
# The real check_inited greps the iptables-save output; stub it so the test output stays quiet
# while still proving the VIP commands refuse to run on an uninitialized gateway.
check_inited() { inited_calls=$((inited_calls+1)); }

lo_addrs="$tmp_dir/lo.addrs"
# Addresses that live on lo without this feature's label: the loopback address itself and whatever
# else the gateway or the operator put there. The sync must never touch them.
lo_foreign="$tmp_dir/lo.foreign"
ip_log="$tmp_dir/ip.log"
ipt_log="$tmp_dir/iptables.log"
hairpin_state="$tmp_dir/hairpin.rules"
hairpin_fip_state="$tmp_dir/hairpin_fip.rules"
snat_state="$tmp_dir/exclusive_snat.rules"
dnat_state="$tmp_dir/exclusive_dnat.rules"
: > "$lo_addrs"
printf '127.0.0.1/8\n10.99.0.1/32\n' > "$lo_foreign"
: > "$ip_log"
: > "$ipt_log"
: > "$hairpin_state"
: > "$hairpin_fip_state"
: > "$snat_state"
: > "$dnat_state"

ip() {
    printf 'ip %s\n' "$*" >> "$ip_log"
    case "$*" in
        "-4 addr show dev $VPC_INTERFACE")
            [[ "$vpc_has_addr" == true ]] && printf 'inet %s/24 scope global %s\n' "$vpc_addr" "$VPC_INTERFACE"
            ;;
        "-4 addr show dev lo label $VIP_ADDR_LABEL")
            while read -r addr; do printf 'inet %s scope host %s\n' "$addr" "$VIP_ADDR_LABEL"; done < "$lo_addrs"
            ;;
        "-4 addr show dev lo")
            # The unscoped view, so that dropping the label filter from the script is caught here
            # (the sync would then try to release addresses it does not own).
            while read -r addr; do printf 'inet %s scope host %s\n' "$addr" "$VIP_ADDR_LABEL"; done < "$lo_addrs"
            while read -r addr; do printf 'inet %s scope host lo\n' "$addr"; done < "$lo_foreign"
            ;;
        "addr add "*" dev lo label $VIP_ADDR_LABEL")
            local addr="${3}"
            grep -qxF "$addr" "$lo_addrs" || printf '%s\n' "$addr" >> "$lo_addrs"
            ;;
        "addr del "*" dev lo")
            local addr="${3}"
            if grep -qxF "$addr" "$lo_foreign"; then
                echo "refused to delete unlabeled lo address: $addr" >&2
                return 1
            fi
            grep -vxF "$addr" "$lo_addrs" > "$lo_addrs.new" || true
            mv "$lo_addrs.new" "$lo_addrs"
            ;;
        *)
            echo "unexpected ip invocation: $*" >&2
            return 1
            ;;
    esac
}

# Only the chains and forms the two ensure_* helpers use are emulated.
ipt_stub() {
    printf 'iptables %s\n' "$*" >> "$ipt_log"
    local full="$*" op chain rule state
    op="${3}"
    chain="${4}"
    rule="${full#-t nat $op $chain}"
    rule="${rule# }"
    case "$chain" in
        HAIRPIN_SNAT) state="$hairpin_state" ;;
        HAIRPIN_FIP) state="$hairpin_fip_state" ;;
        EXCLUSIVE_SNAT) state="$snat_state" ;;
        EXCLUSIVE_DNAT) state="$dnat_state" ;;
        *)
            echo "unexpected iptables chain: $*" >&2
            return 1
            ;;
    esac
    case "$op" in
        -N) ;;
        -F) : > "$state" ;;
        -C) grep -qxF -e "$rule" "$state" ;;
        -A) grep -qxF -e "$rule" "$state" || printf '%s\n' "$rule" >> "$state" ;;
        -D)
            grep -vxF -e "$rule" "$state" > "$state.new" || true
            mv "$state.new" "$state"
            ;;
        *)
            echo "unexpected iptables invocation: $*" >&2
            return 1
            ;;
    esac
}

# The migration helper locates the rules to rewrite in the iptables-save output.
# iptables-save prints a bare address as a /32 prefix; the DNAT lookups in add/del_floating_ip
# match on that form, so the stub normalizes it the same way.
iptables_save() {
    local rule
    while IFS= read -r rule; do [ -n "$rule" ] && printf -- '-A HAIRPIN_SNAT %s\n' "$rule"; done < "$hairpin_state"
    while IFS= read -r rule; do [ -n "$rule" ] && printf -- '-A HAIRPIN_FIP %s\n' "$rule"; done < "$hairpin_fip_state"
    while IFS= read -r rule; do [ -n "$rule" ] && printf -- '-A EXCLUSIVE_SNAT %s\n' "$rule"; done < "$snat_state"
    while IFS= read -r rule; do [ -n "$rule" ] && printf -- '-A EXCLUSIVE_DNAT %s\n' "$rule"; done < "$dnat_state" |
        sed -E 's#(-d [0-9.]+) #\1/32 #'
    return 0
}
iptables_save_cmd=iptables_save
conntrack_log="$tmp_dir/conntrack.log"
conntrack() { printf '%s\n' "$*" >> "$conntrack_log"; }

nft_log="$tmp_dir/nft.log"
: > "$nft_log"
nft() {
    if [[ "$*" == "-f -" ]]; then
        cat >> "$nft_log"
        return 0
    fi
    echo "unexpected nft invocation: $*" >&2
    return 1
}

# The ClusterIPs are held on lo as single /32s, under a label that scopes the set this feature
# owns, so the sync can release the ones the controller no longer asks for.
vip_addr_sync 10.96.1.5
grep -qxF '10.96.1.5/32' "$lo_addrs"
grep -qF 'ip addr add 10.96.1.5/32 dev lo label lo:ko-vip' "$ip_log"
# re-running must converge, not stack addresses or re-add what is already held
add_calls="$(grep -c 'addr add' "$ip_log")"
vip_addr_sync 10.96.1.5
[[ "$(wc -l < "$lo_addrs")" == 1 ]]
[[ "$(grep -c 'addr add' "$ip_log")" == "$add_calls" ]]
# a VIP that is not in the desired set is released, whoever left it behind: this is what makes the
# lo state converge without tracking which rule last referenced an address
vip_addr_sync 10.96.2.7
grep -qxF '10.96.2.7/32' "$lo_addrs"
refute grep -qxF '10.96.1.5/32' "$lo_addrs"
# an empty desired set releases everything this feature owns
vip_addr_sync
[[ ! -s "$lo_addrs" ]]
# ... and nothing else: the label is what scopes the set, so the loopback address and any other
# unlabeled address on lo survive (the ip stub fails the test if one of them is deleted)
[[ "$(wc -l < "$lo_foreign")" == 2 ]]
# ... and then does not call ip addr del again
del_calls="$(grep -c 'addr del' "$ip_log")"
vip_addr_sync
[[ "$(grep -c 'addr del' "$ip_log")" == "$del_calls" ]]
# the VIP set is only programmed on an initialized gateway
[[ "$inited_calls" == 5 ]]
refute vip_addr_sync 'not-an-ip'

# Hairpin SNAT: one wildcard rule, rebuilt from scratch, replaces whatever an older version left
# in the chain (per-EIP and per-identity rules).
printf '%s\n' \
    '-m mark --mark 0x1/0x1 -o eth0 -m conntrack --ctstate DNAT --ctorigdst 203.0.113.10 -j SNAT --to-source 203.0.113.10' \
    '-m mark --mark 0x1/0x1 -o eth0 -p tcp -m conntrack --ctstate DNAT --ctorigdst 10.96.1.5 --ctorigdstport 80 -j SNAT --to-source 10.0.7.254' \
    > "$hairpin_state"
ensure_hairpin_snat
# the lanIP fallback is off by default, so init leaves the per-FIP jump alone: traffic dialed at a
# FIP is SNAT'd to that FIP whatever the switch says
[[ "$(wc -l < "$hairpin_state")" == 1 ]]
[[ "$(head -1 "$hairpin_state")" == '-m mark --mark 0x1/0x1 -o eth0 -m conntrack --ctstate DNAT -j HAIRPIN_FIP' ]]
# the legacy nft hairpin chains of the lanVIP feature go with them
grep -qF 'delete chain ip kube-ovn postrouting' "$nft_log"
grep -qF 'delete chain ip kube-ovn lanvip-snat' "$nft_log"
# rebuilding converges instead of stacking
ensure_hairpin_snat
[[ "$(wc -l < "$hairpin_state")" == 1 ]]

# Switching the fallback on appends it after the per-FIP jump, so a FIP still wins over it.
: > "$ipt_log"
hairpin_snat_mode enable
[[ "$(wc -l < "$hairpin_state")" == 2 ]]
[[ "$(head -1 "$hairpin_state")" == '-m mark --mark 0x1/0x1 -o eth0 -m conntrack --ctstate DNAT -j HAIRPIN_FIP' ]]
[[ "$(tail -1 "$hairpin_state")" == '-m mark --mark 0x1/0x1 -o eth0 -m conntrack --ctstate DNAT -j MASQUERADE --random-fully' ]]
# the fallback rule is address-independent: no VPC address is read to build it
vpc_has_addr=false
hairpin_snat_mode enable
[[ "$(wc -l < "$hairpin_state")" == 2 ]]
vpc_has_addr=true
# ... and switching it back off keeps the per-FIP jump
hairpin_snat_mode disable
[[ "$(wc -l < "$hairpin_state")" == 1 ]]
grep -qxF -- '-m mark --mark 0x1/0x1 -o eth0 -m conntrack --ctstate DNAT -j HAIRPIN_FIP' "$hairpin_state"
# switching the fallback must not touch the per-FIP chain: rebuilding it would leave the FIP
# hairpin rules missing for the duration of an operation that has nothing to do with them
refute grep -q -- '-t nat -F HAIRPIN_FIP' "$ipt_log"
refute grep -q -- '-t nat -A HAIRPIN_FIP ' "$ipt_log"
# ... and an unknown mode is refused rather than silently rewriting the chain
refute hairpin_snat_mode bogus

# init converges the form of the chain, not the policy: an operator who hot-replaces the script
# and runs init by hand must not silently lose the fallback the gateway was running with, because
# the controller records that state and would see the gateway as converged.
hairpin_snat_mode enable
[[ "$(wc -l < "$hairpin_state")" == 2 ]]
[[ "$(current_hairpin_fallback_state)" == enable ]]
hairpin_snat_mode disable
[[ "$(current_hairpin_fallback_state)" == disable ]]
# ... and init must actually be the caller of that decision, otherwise the function above is
# dead code and a manual init silently resets the policy (the chain rebuild is form-only).
grep -qF 'ensure_hairpin_snat "$(current_hairpin_fallback_state)"' "$script"
# the state is only programmed on an initialized gateway (one check_inited per mode call)
[[ "$inited_calls" == 10 ]]

# The per-FIP chain is derived from the authoritative FIP map in EXCLUSIVE_DNAT, so a gateway
# whose script was hot-replaced gets it at init without the controller replaying every FIP.
printf '%s\n' \
    '-d 203.0.113.80/32 -j DNAT --to-destination 10.0.7.8' \
    '-d 203.0.113.81/32 -j DNAT --to-destination 10.0.7.9' \
    > "$dnat_state"
printf '%s\n' '-m conntrack --ctorigdst 203.0.113.99 -j SNAT --to-source 203.0.113.99 --random-fully' > "$hairpin_fip_state"
ensure_hairpin_fip
[[ "$(wc -l < "$hairpin_fip_state")" == 2 ]]
# the rule is keyed on the dialed EIP and SNATs to that same EIP, the form this gateway used
# before the hairpin rules were unified
grep -qxF -- '-m conntrack --ctorigdst 203.0.113.80 -j SNAT --to-source 203.0.113.80 --random-fully' "$hairpin_fip_state"
refute grep -qF -- '203.0.113.99' "$hairpin_fip_state"
: > "$dnat_state"
: > "$hairpin_fip_state"

# FIP egress SNAT must be scoped to the external interface, otherwise it also catches the
# VPC-bound hairpin traffic of its own client. Rules an older script wrote carry no -o and
# add_floating_ip never repairs them (it returns early once the DNAT rule exists), so init
# rewrites them.
printf '%s\n' \
    '-s 10.0.7.1/32 -j SNAT --to-source 203.0.113.20' \
    '-o net1 -s 10.0.7.2/32 -j SNAT --to-source 203.0.113.21' \
    > "$snat_state"
ensure_exclusive_snat_oif
[[ "$(wc -l < "$snat_state")" == 2 ]]
grep -qxF -- '-o net1 -s 10.0.7.1/32 -j SNAT --to-source 203.0.113.20' "$snat_state"
# a rule that is already scoped is left alone, down to its field order
grep -qxF -- '-o net1 -s 10.0.7.2/32 -j SNAT --to-source 203.0.113.21' "$snat_state"
[[ "$(grep -c -- '-t nat -D EXCLUSIVE_SNAT ' "$ipt_log")" == 1 ]]
# converged: a second run rewrites nothing
ensure_exclusive_snat_oif
[[ "$(wc -l < "$snat_state")" == 2 ]]
[[ "$(grep -c -- '-t nat -D EXCLUSIVE_SNAT ' "$ipt_log")" == 1 ]]

# FIP deletion recovers the full installed rule, including the new output-interface match.
printf '%s\n' '-d 203.0.113.20/32 -j DNAT --to-destination 10.0.7.1' > "$dnat_state"
if ! del_floating_ip 203.0.113.20; then exit 1; fi
[[ ! -s "$dnat_state" ]]
refute grep -qF -- '--to-source 203.0.113.20' "$snat_state"
grep -qxF -- '-o net1 -s 10.0.7.2/32 -j SNAT --to-source 203.0.113.21' "$snat_state"
grep -qxF -- '-D -d 203.0.113.20' "$conntrack_log"
del_calls=$(grep -c -- '-t nat -D EXCLUSIVE_' "$ipt_log")
if ! del_floating_ip 203.0.113.20; then exit 1; fi
[[ "$(grep -c -- '-t nat -D EXCLUSIVE_' "$ipt_log")" == "$del_calls" ]]
[[ "$(wc -l < "$conntrack_log")" == 1 ]]

# A FIP carries its own hairpin rule, written whether or not eip mode is active: switching modes
# is then only the jump rule in HAIRPIN_SNAT, never per-FIP bookkeeping.
: > "$snat_state"
: > "$dnat_state"
: > "$hairpin_fip_state"
: > "$ipt_log"
if ! ( add_floating_ip 203.0.113.30,10.0.7.3 ); then exit 1; fi
grep -qxF -- '-m conntrack --ctorigdst 203.0.113.30 -j SNAT --to-source 203.0.113.30 --random-fully' "$hairpin_fip_state"
# re-adding the same FIP converges: add_floating_ip returns early on its DNAT rule, and the
# hairpin rule is repaired on that path too (a gateway programmed before this feature existed).
# The iptables stub ignores a duplicate -A, so idempotency is asserted on the calls, not on the
# resulting state: exactly one -A, and every later call checking with -C first.
if ! ( add_floating_ip 203.0.113.30,10.0.7.3 ); then exit 1; fi
if ! ( add_floating_ip 203.0.113.30,10.0.7.3 ); then exit 1; fi
[[ "$(grep -c -- '-t nat -A HAIRPIN_FIP ' "$ipt_log")" == 1 ]]
[[ "$(grep -c -- '-t nat -C HAIRPIN_FIP ' "$ipt_log")" == 3 ]]
[[ "$(wc -l < "$hairpin_fip_state")" == 1 ]]
# a conflicting binding is still rejected, and adds no second rule
refute add_floating_ip 203.0.113.30,10.0.7.4
[[ "$(grep -c -- '-t nat -A HAIRPIN_FIP ' "$ipt_log")" == 1 ]]
[[ "$(wc -l < "$hairpin_fip_state")" == 1 ]]
# deleting the FIP takes its hairpin rule with it, idempotently: the second delete finds nothing
# through -C and issues no -D
: > "$ipt_log"
if ! del_floating_ip 203.0.113.30; then exit 1; fi
[[ ! -s "$hairpin_fip_state" ]]
[[ "$(grep -c -- '-t nat -D HAIRPIN_FIP ' "$ipt_log")" == 1 ]]
if ! del_floating_ip 203.0.113.30; then exit 1; fi
[[ ! -s "$hairpin_fip_state" ]]
[[ "$(grep -c -- '-t nat -D HAIRPIN_FIP ' "$ipt_log")" == 1 ]]

# QoS filter cleanup must distinguish the EIP class range (0x1-0x7ffe) from the NatGw range
# (0x8000-0xfeff) by value. tc prints classids without leading zeros, so a three-digit classid
# like 1:ab8 or 1:8cd is EIP-owned even though its first hex digit is in the NatGw range's set.
tc_log="$tmp_dir/tc.log"
tc_filters="$tmp_dir/tc.filters"
: > "$tc_log"
cat > "$tc_filters" <<'EOF'
filter parent 1: protocol ip pref 10 u32 fh 800::801 order 2049 key ht 800 bkt 0 flowid 1:ab8 not_in_hw
  match ip src 10.0.0.1/32
filter parent 1: protocol ip pref 10 u32 fh 800::802 order 2050 key ht 800 bkt 0 flowid 1:8cd not_in_hw
  match ip src 10.0.0.2/32
filter parent 1: protocol ip pref 20 u32 fh 800::803 order 2051 key ht 800 bkt 0 flowid 1:8005 not_in_hw
  match ip src 10.20.0.1/32
EOF
tc() {
    printf 'tc %s\n' "$*" >> "$tc_log"
    case "$*" in
        "-p filter show dev "*" parent 1:") cat "$tc_filters" ;;
        "filter del dev "*" parent 1: prio "*" handle "*" u32") : ;;
        "class del dev "*" classid 1:0x"*) : ;;
        *) echo "unexpected tc invocation: $*" >&2; return 1 ;;
    esac
}
delete_htb_filter_and_class "$VPC_INTERFACE" "10.0.0.1/32" "src" "eip"
grep -qF 'tc filter del dev eth0 parent 1: prio 10 handle 800::801 u32' "$tc_log"
grep -qF 'tc class del dev eth0 classid 1:0xab8' "$tc_log"
delete_htb_filter_and_class "$VPC_INTERFACE" "10.0.0.2/32" "src" "eip"
grep -qF 'tc class del dev eth0 classid 1:0x8cd' "$tc_log"
# the NatGw range keeps matching its four-digit classids
delete_htb_filter_and_class "$VPC_INTERFACE" "10.20.0.1/32" "src" "natgw"
grep -qF 'tc class del dev eth0 classid 1:0x8005' "$tc_log"
# cross-range identities are not touched at all
delete_calls="$(grep -c ' del ' "$tc_log")"
delete_htb_filter_and_class "$VPC_INTERFACE" "10.20.0.1/32" "src" "eip"
[[ "$(grep -c ' del ' "$tc_log")" == "$delete_calls" ]]

# Backward compatibility with filters installed by the previous script version.
# An old EIP classid at the very top of the EIP range (a collision-bumped 0x7fff) is still
# recognized as EIP-owned, and an EIP delete (which passes no priority) ignores the stored pref.
cat > "$tc_filters" <<'EOF'
filter parent 1: protocol ip pref 10 u32 fh 800::807 order 2055 key ht 800 bkt 0 flowid 1:7fff not_in_hw
  match ip src 10.0.0.3/32
EOF
delete_htb_filter_and_class "$VPC_INTERFACE" "10.0.0.3/32" "src" "eip"
grep -qF 'tc filter del dev eth0 parent 1: prio 10 handle 800::807 u32' "$tc_log"
grep -qF 'tc class del dev eth0 classid 1:0x7fff' "$tc_log"

# tc rewrites a requested priority 0 to pref 49152 when storing the filter, so an old NatGw
# filter carries no pref matching the rule's priority. With a single candidate for the identity
# the filter is still adopted and deleted.
cat > "$tc_filters" <<'EOF'
filter parent 1: protocol ip pref 49152 u32 fh 800::808 order 2056 key ht 800 bkt 0 flowid 1:8005 not_in_hw
  match ip dst 10.30.0.1/16
EOF
delete_htb_filter_and_class "$VPC_INTERFACE" "10.30.0.1/16" "dst" "natgw" "0"
grep -qF 'tc filter del dev eth0 parent 1: prio 49152 handle 800::808 u32' "$tc_log"
grep -qF 'tc class del dev eth0 classid 1:0x8005' "$tc_log"

# ... but when two candidates share the identity, the pref-mismatching one must not be guessed:
# nothing is deleted rather than adopting a foreign filter.
cat > "$tc_filters" <<'EOF'
filter parent 1: protocol ip pref 49152 u32 fh 800::809 order 2057 key ht 800 bkt 0 flowid 1:8005 not_in_hw
  match ip dst 10.30.0.1/16
filter parent 1: protocol ip pref 49152 u32 fh 800::80a order 2058 key ht 800 bkt 0 flowid 1:8006 not_in_hw
  match ip dst 10.30.0.1/16
EOF
delete_calls="$(grep -c ' del ' "$tc_log")"
delete_htb_filter_and_class "$VPC_INTERFACE" "10.30.0.1/16" "dst" "natgw" "0"
[[ "$(grep -c ' del ' "$tc_log")" == "$delete_calls" ]]

echo 'PASS: share-DNAT VIP address, hairpin and QoS class rules' 
