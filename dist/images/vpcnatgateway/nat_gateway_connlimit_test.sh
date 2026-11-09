#!/usr/bin/env bash
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

echo "PASS: nat-gateway connlimit & stateless detection tests passed"
