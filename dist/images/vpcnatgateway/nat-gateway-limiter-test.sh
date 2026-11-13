#!/usr/bin/env bash
# Note: This test script is strictly for local and CI automated verification.
# It is excluded from the production gateway container image (never packaged).
set -Eeuo pipefail

script="$(cd "$(dirname "$0")" && pwd)/nat-gateway.sh"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

export NAT_GW_ENV_FILE="$tmp_dir/nat-gateway.env"
# Load script functions without command dispatcher
eval "$(sed '/^opt=\$1$/,$d' "$script")"

# 1. Test to_nft_rate conversion
assert_rate() {
    local input=$1
    local want=$2
    local got
    got=$(to_nft_rate "$input")
    if [ "$got" != "$want" ]; then
        echo "FAIL to_nft_rate '$input': got '$got', want '$want'" >&2
        exit 1
    fi
}

assert_rate "1/s" "1/second"
assert_rate "9/s" "9/second"
assert_rate "500/s" "500/second"
assert_rate "500/sec" "500/second"
assert_rate "500/second" "500/second"
assert_rate "100/m" "100/minute"
assert_rate "100/min" "100/minute"
assert_rate "100/minute" "100/minute"
assert_rate "200/h" "200/hour"
assert_rate "200/hour" "200/hour"
assert_rate "50/d" "50/day"
assert_rate "50/day" "50/day"

if to_nft_rate "500/invalid" 2>/dev/null; then
    echo "FAIL to_nft_rate accepted invalid unit" >&2
    exit 1
fi

# 2. Test is_stateless_active mode handling
if ! is_stateless_active "stateless"; then
    echo "FAIL: mode 'stateless' must be treated as active stateless" >&2
    exit 1
fi

if ! is_stateless_active "stateless-nft"; then
    echo "FAIL: mode 'stateless-nft' must be treated as active stateless" >&2
    exit 1
fi

if is_stateless_active "stateful"; then
    echo "FAIL: mode 'stateful' must not be treated as active stateless" >&2
    exit 1
fi

if is_stateless_active ""; then
    echo "FAIL: empty mode must not be treated as active stateless" >&2
    exit 1
fi

if is_stateless_active "unknown"; then
    echo "FAIL: unknown mode must not be treated as active stateless" >&2
    exit 1
fi

# 3. Test calc_bucket_expire_ms
assert_expire() {
    local rate=$1
    local burst=$2
    local want=$3
    local got
    got=$(calc_bucket_expire_ms "$rate" "$burst")
    if [ "$got" != "$want" ]; then
        echo "FAIL calc_bucket_expire_ms '$rate' '$burst': got '$got', want '$want'" >&2
        exit 1
    fi
}

assert_expire "1000/s" "2000" "10000"
assert_expire "20000/s" "40000" "10000"
assert_expire "1/minute" "1" "60000"
assert_expire "1/m" "5" "300000"
assert_expire "2/hour" "1" "1800000"
# Excessively large duration should be capped at 86400000ms (24h)
assert_expire "1/day" "10" "86400000"
assert_expire "1/hour" "40000" "86400000"
# Zero or invalid rate should safely fall back to baseline 10000ms without division by zero
assert_expire "0/s" "10" "10000"
assert_expire "invalid" "10" "10000"
# Zero or invalid burst should safely fall back to baseline 10000ms
assert_expire "1000/s" "0" "10000"
assert_expire "1000/s" "invalid" "10000"
assert_expire "1000/s" "" "10000"

# 4. Test rate pattern validation
rate_pattern="^[1-9][0-9]*/(s|sec|second|m|min|minute|h|hour|d|day)$"
if [[ "0/s" =~ $rate_pattern ]]; then
    echo "FAIL: rate_pattern must reject 0/s" >&2
    exit 1
fi
if ! [[ "1/s" =~ $rate_pattern ]]; then
    echo "FAIL: rate_pattern must accept 1/s" >&2
    exit 1
fi
if ! [[ "9/s" =~ $rate_pattern ]]; then
    echo "FAIL: rate_pattern must accept 9/s" >&2
    exit 1
fi
if ! [[ "100/s" =~ $rate_pattern ]]; then
    echo "FAIL: rate_pattern must accept 100/s" >&2
    exit 1
fi

# 5. Test is_syn_limit_disabled
for disabled_val in "0" "0/s" "disable" "disabled" "none" "off"; do
    if ! is_syn_limit_disabled "$disabled_val"; then
        echo "FAIL: is_syn_limit_disabled must return 0 for '$disabled_val'" >&2
        exit 1
    fi
done

for enabled_val in "100/s" "1000/s" "1/minute" "enable" "0/invalid" "0/foo"; do
    if is_syn_limit_disabled "$enabled_val"; then
        echo "FAIL: is_syn_limit_disabled must return 1 for '$enabled_val'" >&2
        exit 1
    fi
done

# 6. Test single-backend enforcement and independent src/dst SYN limit disablement
mock_log="$tmp_dir/mock.log"
check_inited() { return 0; }
nft() {
    if [ "${MOCK_NFT_FAIL:-0}" -eq 1 ]; then
        return 1
    fi
    if [ "${1:-}" = "-f" ]; then
        while IFS= read -r line; do
            echo "NFT_TX: $line" >> "$mock_log"
        done
        return 0
    fi
    if [ "${1:-}" = "list" ] && [ "${2:-}" = "chain" ]; then
        return 1
    fi
    echo "NFT: $*" >> "$mock_log"
    return 0
}
mock_iptables() {
    echo "IPT: $*" >> "$mock_log"
    if [[ "$*" == *"-C "* ]]; then
        return 1
    fi
    return 0
}
nft_cmd="nft"
iptables_cmd="mock_iptables"
iptables_save_cmd="true"
EXTERNAL_INTERFACE="net1"

# 6a. Auto mode with both nft and iptables present must only install nft rules (no double enforcement)
: > "$mock_log"
tcp_connlimit_set "50000" "stateful" >/dev/null
if grep -q -- "-A VPC_CONNLIMIT" "$mock_log"; then
    echo "FAIL: auto mode must not install iptables VPC_CONNLIMIT when nft succeeds" >&2
    exit 1
fi
if ! grep -q -- "NFT_TX: add rule ip kube-ovn conn-limit tcp flags & (fin|syn|rst|ack) == syn ct count over 50000 drop" "$mock_log"; then
    echo "FAIL: auto mode did not install expected nft conn-limit rule" >&2
    exit 1
fi
if ! grep -q -- 'NFT_TX: add rule ip kube-ovn conn-limit iifname != "net1" return' "$mock_log"; then
    echo "FAIL: nft conn-limit must filter out non-external interface traffic" >&2
    exit 1
fi

# 6b. Auto mode falls back to iptables when nft rule installation fails
: > "$mock_log"
MOCK_NFT_FAIL=1 tcp_connlimit_set "50000" "stateful" >/dev/null 2>&1
if ! grep -q -- 'IPT: -t mangle -I PREROUTING 1 -i net1 -p tcp --syn -j VPC_CONNLIMIT' "$mock_log"; then
    echo "FAIL: iptables connlimit PREROUTING jump must scope to external interface" >&2
    exit 1
fi
if ! grep -q -- "IPT: -t mangle -A VPC_CONNLIMIT -p tcp --syn -m connlimit --connlimit-above 50000" "$mock_log"; then
    echo "FAIL: auto mode must fall back to iptables when nft fails" >&2
    exit 1
fi

# 6c. Independent src_rate=0 keeps dst_rate active, and dst_rate=0 keeps src_rate active
: > "$mock_log"
tcp_syn_limit_set "0" "2000" "8000/s" "16000" >/dev/null
if grep -q "add set ip kube-ovn syn_src_meter" "$mock_log"; then
    echo "FAIL: src_rate=0 must not install syn_src_meter" >&2
    exit 1
fi
if ! grep -q "add set ip kube-ovn syn_dst_meter" "$mock_log"; then
    echo "FAIL: src_rate=0 must keep syn_dst_meter active" >&2
    exit 1
fi

: > "$mock_log"
tcp_syn_limit_set "300/s" "600" "0" "40000" >/dev/null
if ! grep -q "add set ip kube-ovn syn_src_meter" "$mock_log"; then
    echo "FAIL: dst_rate=0 must keep syn_src_meter active" >&2
    exit 1
fi
if grep -q "add set ip kube-ovn syn_dst_meter" "$mock_log"; then
    echo "FAIL: dst_rate=0 must not install syn_dst_meter" >&2
    exit 1
fi

echo "PASS: nat-gateway connlimit & stateless detection tests passed"
