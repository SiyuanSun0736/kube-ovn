#!/usr/bin/env bash
# Tests for sync_nft_lanvip (the nft-lanvip-sync command) with a stubbed nft state machine.
set -Eeuo pipefail

script="$(cd "$(dirname "$0")" && pwd)/nat-gateway.sh"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

# Keep the script from reading an interface file left behind on the host.
export NAT_GW_ENV_FILE="$tmp_dir/nat-gateway.env"

# Load definitions without executing the command dispatcher at the end.
eval "$(sed '/^opt=\$1$/,$d' "$script")"

VPC_INTERFACE=eth0
vpc_addr=10.0.7.254
eip_addr=172.20.0.5

inited_calls=0
check_inited() { inited_calls=$((inited_calls+1)); }

# ---- stubbed nft state machine -------------------------------------------------
# state files: TABLE (touch), chains (one per line), sets, elements ("vip proto port chain"),
# rules-log ("table chain rule"), cmd-log (every batch line seen).
TABLE="$tmp_dir/table"
CHAINS="$tmp_dir/chains"
SETS="$tmp_dir/sets"
ELEMENTS="$tmp_dir/elements"
RULES_LOG="$tmp_dir/rules.log"
: > "$CHAINS"
: > "$SETS"
: > "$ELEMENTS"
: > "$RULES_LOG"

elem_key_from_batch() { # "add element ip T M { 1.2.3.4 . tcp . 80 : goto dnat-x }" -> "1.2.3.4 tcp 80 dnat-x"
    echo "$1" | sed -n 's/.*{ \(.*\) \. \(.*\) \. \(.*\) : goto \([^ ]*\) }.*/\1 \2 \3 \4/p'
}

nft_handle_batch_line() {
    local line="$1" name
    case "$line" in
        "add table ip $NFT_TABLE") touch "$TABLE" ;;
        "add chain ip $NFT_TABLE "*)
            name=$(echo "$line" | awk '{print $5}')
            if grep -qxF "$name" "$CHAINS"; then
                # real nft: `add chain` on an existing chain fails with EEXIST and the whole
                # `nft -f` batch aborts. Report a fatal (2) so the batch runner below stops.
                echo "Error: Could not process rule: File exists (add chain $name)" >&2
                return 2
            fi
            echo "$name" >> "$CHAINS"
            ;;
        "flush chain ip $NFT_TABLE "*)
            name=$(echo "$line" | awk '{print $5}')
            grep -vF " $name " "$RULES_LOG" > "$RULES_LOG.new" 2>/dev/null || true
            mv "$RULES_LOG.new" "$RULES_LOG"
            ;;
        "delete chain ip $NFT_TABLE "*)
            name=$(echo "$line" | awk '{print $5}')
            # model the kernel's referential integrity: a chain still jumped at cannot be deleted
            if grep -qE " (jump|goto) $name\$" "$RULES_LOG"; then
                echo "Error: Could not process rule: Resource busy (chain $name still referenced)" >&2
                return 1
            fi
            grep -vxF "$name" "$CHAINS" > "$CHAINS.new" || true
            mv "$CHAINS.new" "$CHAINS"
            # deleting a chain also drops its rules
            grep -vF " $name " "$RULES_LOG" > "$RULES_LOG.new" 2>/dev/null || true
            mv "$RULES_LOG.new" "$RULES_LOG"
            ;;
        "add set ip $NFT_TABLE "*)
            name=$(echo "$line" | awk '{print $5}')
            if grep -qxF "$name" "$SETS"; then
                # same EEXIST strictness as `add chain`
                echo "Error: Could not process rule: File exists (add set $name)" >&2
                return 2
            fi
            echo "$name" >> "$SETS"
            ;;
        "delete set ip $NFT_TABLE "*)
            name=$(echo "$line" | awk '{print $5}')
            grep -vxF "$name" "$SETS" > "$SETS.new" || true
            mv "$SETS.new" "$SETS"
            ;;
        "add element ip $NFT_TABLE "*)
            local key
            key=$(elem_key_from_batch "$line")
            [ -n "$key" ] || { echo "unparseable element: $line" >&2; return 1; }
            # upsert by vip.proto.port
            local vip proto port
            vip=$(echo "$key" | awk '{print $1}'); proto=$(echo "$key" | awk '{print $2}'); port=$(echo "$key" | awk '{print $3}')
            grep -v "^$vip $proto $port " "$ELEMENTS" > "$ELEMENTS.new" 2>/dev/null || true
            mv "$ELEMENTS.new" "$ELEMENTS"
            echo "$key" >>"$ELEMENTS"
            ;;
        "add rule ip $NFT_TABLE "*)
            name=$(echo "$line" | awk '{print $5}')
            echo "$NFT_TABLE $name ${line#add rule ip $NFT_TABLE $name }" >> "$RULES_LOG"
            ;;
        "add map ip $NFT_TABLE "*) ;;
        *)
            echo "unexpected nft batch line: $line" >&2
            return 1
            ;;
    esac
}

nft() {
    case "$*" in
        "-f -"|"-f /dev/stdin")
            local rc=0 line line_rc
            while IFS= read -r line; do
                [ -z "$line" ] && continue
                printf 'batch %s\n' "$line" >> "$tmp_dir/cmd.log"
                nft_handle_batch_line "$line"
                line_rc=$?
                # fatal errors (EEXIST & friends) abort the rest of the batch, like the
                # single netlink transaction a real `nft -f` submits
                [ "$line_rc" -eq 2 ] && return 2
                [ "$line_rc" -ne 0 ] && rc=1
            done
            return $rc
            ;;
        "list table ip $NFT_TABLE")
            [ -f "$TABLE" ] || return 1
            {
                echo "table ip $NFT_TABLE {"
                while read -r c; do printf '\tchain %s {\n\t}\n' "$c"; done < "$CHAINS"
                while read -r s; do printf '\tset %s {\n\t}\n' "$s"; done < "$SETS"
                printf '\tmap %s {\n' "$NFT_SERVICES_MAP"
                printf '\t\ttype ipv4_addr . inet_proto . inet_service : verdict\n'
                if [ -s "$ELEMENTS" ]; then
                    printf '\t\telements = { '
                    local first=1
                    while read -r vip proto port chain; do
                        [ "$first" = 0 ] && printf ', '
                        printf '%s . %s . %s : goto %s' "$vip" "$proto" "$port" "$chain"
                        first=0
                    done < "$ELEMENTS"
                    printf ' }\n'
                fi
                printf '\t}\n}\n'
            }
            ;;
        "list map ip $NFT_TABLE $NFT_SERVICES_MAP")
            nft "list table ip $NFT_TABLE" || return 1
            ;;
        "list chain ip $NFT_TABLE "*)
            local want_chain
            want_chain=$(echo "$*" | awk '{print $5}')
            [ -f "$TABLE" ] || return 1
            grep -qxF "$want_chain" "$CHAINS"
            ;;
        "list set ip $NFT_TABLE "*)
            local want_set
            want_set=$(echo "$*" | awk '{print $5}')
            [ -f "$TABLE" ] || return 1
            grep -qxF "$want_set" "$SETS"
            ;;
        "delete element ip $NFT_TABLE "*)
            local inner vip proto port
            inner=$(echo "$*" | sed -n 's/.*{ \(.*\) \. \(.*\) \. \(.*\) }.*/\1 \2 \3/p')
            vip=$(echo "$inner" | awk '{print $1}'); proto=$(echo "$inner" | awk '{print $2}'); port=$(echo "$inner" | awk '{print $3}')
            grep -v "^$vip $proto $port " "$ELEMENTS" > "$ELEMENTS.new" 2>/dev/null || true
            mv "$ELEMENTS.new" "$ELEMENTS"
            printf 'delete-element %s\n' "$inner" >> "$tmp_dir/cmd.log"
            ;;
        *)
            echo "unexpected nft invocation: $*" >&2
            return 1
            ;;
    esac
}

conntrack() { printf 'conntrack %s\n' "$*" >> "$tmp_dir/cmd.log"; }

ip() {
    case "$*" in
        "-4 addr show dev $VPC_INTERFACE")
            printf 'inet %s/24 scope global %s\n' "$vpc_addr" "$VPC_INTERFACE"
            ;;
        *)
            echo "unexpected ip invocation: $*" >&2
            return 1
            ;;
    esac
}

assert_eq() {
    if [ "$1" != "$2" ]; then
        echo "ASSERT FAILED: $3 (want '$1', got '$2')" >&2
        exit 1
    fi
}
assert_file_contains() { grep -qF -- "$2" "$1" || { echo "ASSERT FAILED: $1 lacks '$2'" >&2; exit 1; }; }
assert_file_lacks() { ! grep -qF -- "$2" "$1" || { echo "ASSERT FAILED: $1 unexpectedly contains '$2'" >&2; exit 1; }; }

: > "$tmp_dir/cmd.log"

echo "== sync two identities (tcp none + udp clientip), plus a foreign EIP element that must survive"
# a pre-existing EIP identity of the exclusive feature: same table, different address
echo "$eip_addr tcp 443 dnat-e1a000000001" >> "$ELEMENTS"
echo "dnat-e1a000000001" >> "$CHAINS"
touch "$TABLE"

sync_nft_lanvip \
    "$vpc_addr,80,tcp,none,0,10.0.7.11:8080@10.0.7.12:8080" \
    "$vpc_addr,53,udp,clientip,600,10.0.7.13:53"

assert_eq 3 "$(wc -l < "$ELEMENTS" | tr -d ' ')" "two lanVIP identities + one foreign"
assert_file_contains "$ELEMENTS" "$vpc_addr tcp 80 "
assert_file_contains "$ELEMENTS" "$vpc_addr udp 53 "
assert_file_contains "$ELEMENTS" "$eip_addr tcp 443 dnat-e1a000000001"
# postrouting jump + per-identity snat rules
assert_file_contains "$RULES_LOG" " $NFT_POSTROUTING_CHAIN oifname \"$VPC_INTERFACE\" jump $NFT_LANVIP_SNAT_CHAIN"
assert_file_contains "$RULES_LOG" "ct original ip daddr $vpc_addr ct original proto-dst 80 snat to $vpc_addr fully-random"
assert_file_contains "$RULES_LOG" "ct original ip daddr $vpc_addr ct original proto-dst 53 snat to $vpc_addr fully-random"
assert_file_contains "$RULES_LOG" "meta mark and 0x1 == 0x1"
# clientip identity got its affinity sets
assert_eq 1 "$(grep -c '^aff-' "$SETS")" "one affinity set for the udp identity"
# numeric chains exist for both identities, plus the foreign one
assert_eq 3 "$(grep -c '^dnat-' "$CHAINS")" "two lanVIP identity chains + the foreign one"
assert_eq 2 "$(grep -c "^$NFT_TABLE $NFT_LANVIP_SNAT_CHAIN " "$RULES_LOG")" "two snat rules"

echo "== re-sync the identical set (chains already exist => still succeeds, SNAT chain rebuilt)"
: > "$RULES_LOG"
sync_nft_lanvip \
    "$vpc_addr,80,tcp,none,0,10.0.7.11:8080@10.0.7.12:8080" \
    "$vpc_addr,53,udp,clientip,600,10.0.7.13:53"
assert_eq 3 "$(wc -l < "$ELEMENTS" | tr -d ' ')" "repeat sync: identities unchanged"
assert_eq 2 "$(grep -c "^$NFT_TABLE $NFT_LANVIP_SNAT_CHAIN " "$RULES_LOG")" "repeat sync: snat chain rebuilt with both rules"
assert_file_contains "$RULES_LOG" " $NFT_POSTROUTING_CHAIN oifname \"$VPC_INTERFACE\" jump $NFT_LANVIP_SNAT_CHAIN"

echo "== shrink to one identity"
sync_nft_lanvip "$vpc_addr,53,udp,clientip,600,10.0.7.13:53"
assert_file_lacks "$ELEMENTS" "$vpc_addr tcp 80 "
assert_file_contains "$ELEMENTS" "$vpc_addr udp 53 "
assert_file_contains "$ELEMENTS" "$eip_addr tcp 443 dnat-e1a000000001"
assert_eq 2 "$(grep -c '^dnat-' "$CHAINS")" "stale tcp chain garbage-collected; udp + foreign chains stay"
assert_file_contains "$tmp_dir/cmd.log" "delete-element $vpc_addr tcp 80"
assert_file_contains "$tmp_dir/cmd.log" "conntrack -D -d $vpc_addr -p tcp --dport 80"
# the snat chain was rebuilt with only the udp rule
assert_eq 1 "$(grep -c "^$NFT_TABLE $NFT_LANVIP_SNAT_CHAIN " "$RULES_LOG")" "snat chain shrunk to one rule"
assert_file_contains "$RULES_LOG" "ct original proto-dst 53 "

echo "== wipe (zero args)"
sync_nft_lanvip
assert_eq 1 "$(wc -l < "$ELEMENTS" | tr -d ' ')" "only the foreign EIP identity survives the wipe"
assert_file_contains "$ELEMENTS" "$eip_addr tcp 443 dnat-e1a000000001"
assert_eq 1 "$(grep -c '^dnat-' "$CHAINS" || true)" "lanVIP chains gone; foreign chain stays (its element still references it)"

echo "== validation rejects garbage"
# failure paths use `exit` inside the function, so run them in a subshell
if (sync_nft_lanvip "$vpc_addr,80,tcp,none,0,10.0.x.1:80" 2>/dev/null); then
    echo "ASSERT FAILED: invalid backend ip accepted" >&2
    exit 1
fi
if (sync_nft_lanvip "$vpc_addr,80,tcp,none,0,10.0.7.11:70000" 2>/dev/null); then
    echo "ASSERT FAILED: invalid backend port accepted" >&2
    exit 1
fi
if (sync_nft_lanvip "$vpc_addr,80,sctp,none,0,10.0.7.11:80" 2>/dev/null); then
    echo "ASSERT FAILED: invalid protocol accepted" >&2
    exit 1
fi
if (sync_nft_lanvip "$vpc_addr,80,tcp,bogus,0,10.0.7.11:80" 2>/dev/null); then
    echo "ASSERT FAILED: invalid affinity accepted" >&2
    exit 1
fi

echo "== table absent + zero args is a no-op"
rm -f "$TABLE"
: > "$ELEMENTS"
sync_nft_lanvip >/dev/null
assert_eq 0 "$(wc -l < "$ELEMENTS" | tr -d ' ')" "still empty"

echo "PASS: nat_gateway_lanvip_test.sh"
