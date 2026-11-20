#!/usr/bin/env bash
# Run only inside an isolated mount + network namespace, as documented below.
# unshare -mn -- bash -c 'mount --make-rprivate /; mount -t tmpfs tmpfs /run/netns; bash <this-file>'
set -Eeuo pipefail
if [[ $(readlink /proc/self/ns/net) == "$(readlink /proc/1/ns/net)" ]]; then
    echo 'Refusing to run in the host network namespace' >&2
    exit 1
fi
script="$(cd "$(dirname "$0")" && pwd)/nat-gateway.sh"
tmp=$(mktemp -d)
# cleanup 终止后台 Python 测试服务并销毁测试网络命名空间
cleanup() {
    for ns in sl-gw sl-lan sl-wan; do
        ip netns pids "$ns" 2>/dev/null | xargs -r kill 2>/dev/null || true
        ip netns del "$ns" 2>/dev/null || true
    done
    rm -rf "$tmp"
}
trap cleanup EXIT
for ns in sl-gw sl-lan sl-wan; do ip netns add "$ns"; ip -n "$ns" link set lo up; done
ip link add eth0 netns sl-gw type veth peer name eth0 netns sl-lan
ip link add net1 netns sl-gw type veth peer name eth0 netns sl-wan
ip -n sl-gw addr add 10.0.0.1/24 dev eth0
ip -n sl-gw addr add 192.0.2.1/24 dev net1
for host in 10 11 12; do ip -n sl-gw addr add "192.0.2.$host/32" dev net1; done
ip -n sl-lan addr add 10.0.0.2/24 dev eth0
ip -n sl-wan addr add 192.0.2.2/24 dev eth0
for dev in eth0 net1; do ip -n sl-gw link set "$dev" up; done
for ns in sl-lan sl-wan; do ip -n "$ns" link set eth0 up; done
ip -n sl-lan route add default via 10.0.0.1
# WAN 端刻意不配置私网子网路由：响应必须通过 NAT 转发
ip netns exec sl-gw sysctl -qw net.ipv4.ip_forward=1
export NAT_GW_ENV_FILE="$tmp/env"

# gw 在网关网络命名空间中执行 nat-gateway.sh
gw() { ip netns exec sl-gw bash "$script" "$@"; }
gw init eth0,net1 > "$tmp/init.log" 2>&1
gw hairpin-snat-set enable > /dev/null
gw stateless-init
gw stateless-init
cat > "$tmp/echo_server.py" <<'PY'
from http.server import BaseHTTPRequestHandler, HTTPServer
class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(self.client_address[0].encode())
    def log_message(self, *args):
        pass
HTTPServer(('0.0.0.0', 8080), Handler).serve_forever()
PY
for ns in sl-lan sl-wan; do ip netns exec "$ns" python3 "$tmp/echo_server.py" & done

# request 在指定命名空间发起 HTTP GET 请求并返回对端看到的客户端 IP
request() { ip netns exec "$1" curl --noproxy '*' --fail --silent --show-error --max-time 3 "http://$2:8080/"; }
for attempt in {1..20}; do
    if request sl-lan 10.0.0.2 >/dev/null 2>&1; then break; fi
    sleep .1
done

# assert_no_conntrack 校验网关命名空间中的 conntrack 表不存在被跟踪的 TCP 活跃连接
assert_no_conntrack() {
    local count
    count=$(ip netns exec sl-gw conntrack -L -p tcp 2>/dev/null | wc -l)
    [[ "$count" == 0 ]] || { echo "unexpected tracked NAT sessions: $count" >&2; return 1; }
}
# FIP ingress and egress preserve the remote address and bypass conntrack.
gw stateless-fip-add '192.0.2.10,10.0.0.2'
gw stateless-fip-add '192.0.2.10,10.0.0.2'
[[ $(request sl-wan 192.0.2.10) == 192.0.2.2 ]]
[[ $(request sl-lan 192.0.2.2) == 192.0.2.10 ]]
assert_no_conntrack
[[ $(request sl-lan 192.0.2.10) == 192.0.2.10 ]]
ip netns exec sl-gw conntrack -L -p tcp 2>/dev/null | grep -q 'dst=192.0.2.10'
gw stateless-fip-del '192.0.2.10'
gw stateless-fip-del '192.0.2.10'
! ip netns exec sl-gw nft list map ip kube_ovn_stateless_nat sl_fip_dnat | grep -q 192.0.2.10
ip netns exec sl-gw conntrack -F >/dev/null 2>&1
# DNAT rewrites both addresses and ports; hairpin remains tracked.
gw stateless-dnat-add '192.0.2.11,80,tcp,10.0.0.2,8080'
[[ $(ip netns exec sl-wan curl --noproxy '*' -fsS --max-time 3 http://192.0.2.11/) == 192.0.2.2 ]]
assert_no_conntrack
[[ $(ip netns exec sl-lan curl --noproxy '*' -fsS --max-time 3 http://192.0.2.11/) == 10.0.0.1 ]]
# UDP uses independent port rewrites; deleting TCP must preserve UDP.
cat > "$tmp/udp_server.py" <<'PYCODE'
import socket
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.bind(('0.0.0.0', 8081))
while True:
    data, peer = s.recvfrom(4096)
    s.sendto(data, peer)
PYCODE
ip netns exec sl-lan python3 "$tmp/udp_server.py" &
gw stateless-dnat-add '192.0.2.11,81,udp,10.0.0.2,8081'
gw stateless-dnat-del '192.0.2.11,80,tcp'
ip netns exec sl-wan python3 - <<'PYCODE'
import socket
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.settimeout(3)
s.sendto(b'stateless UDP', ('192.0.2.11', 81))
data, peer = s.recvfrom(4096)
assert data == b'stateless UDP' and peer == ('192.0.2.11', 81), (data, peer)
PYCODE
ip netns exec sl-gw conntrack -L -p udp 2>/dev/null | grep -q . && exit 1
gw stateless-dnat-del '192.0.2.11,81,udp'

! ip netns exec sl-gw nft list map ip kube_ovn_stateless_nat sl_dnat_in | grep -q 192.0.2.11
ip netns exec sl-gw conntrack -F >/dev/null 2>&1
# SNAT accepts /32 and bare host spellings for the same deletion identity.
gw stateless-snat-add '192.0.2.12,10.0.0.2/32'
# Outbound traffic from Pod is SNAT-ed to EIP
[[ $(request sl-lan 192.0.2.2) == 192.0.2.12 ]]
# Inbound TCP connection from public WAN to SNAT EIP must be dropped (security protection)
! request sl-wan 192.0.2.12 >/dev/null 2>&1
assert_no_conntrack
gw stateless-snat-del '192.0.2.12,10.0.0.2'
! ip netns exec sl-gw nft list map ip kube_ovn_stateless_nat sl_snat_in | grep -q 192.0.2.12
if gw stateless-snat-add '192.0.2.12,10.0.0.0/24'; then exit 1; fi
echo 'PASS: FIP/DNAT/SNAT bidirectional TCP, UDP DNAT, hairpin, replay, deletion, zero conntrack'
