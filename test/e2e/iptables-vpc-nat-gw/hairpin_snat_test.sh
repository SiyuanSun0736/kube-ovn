#!/usr/bin/env bash
# Data-plane test of the vpc-nat-gw hairpin SNAT against real resources: a real VpcNatGateway Pod,
# real IptablesEIP/IptablesFIPRule objects and real workload Pods. The backends answer with the
# source address they see, so the assertions come from the kernel of the gateway under test.
#
# Preconditions (not created here, they belong to the cluster):
#   - vpc-nat-gw enabled: cm/ovn-vpc-nat-gw-config data.enable-vpc-nat-gw = "true"
#   - an external subnet plus its NetworkAttachmentDefinition (default ovn-vpc-external-network)
#     with at least 4 free addresses
#   - cm/ovn-vpc-nat-config data.image pointing at the vpc-nat-gateway image under test
#   - TEST_IMAGE contains python3; the gateway image contains nc (OpenBSD netcat)
# Reserved hairpin-e2e resource names must not be used by other workloads.
#
# Run with: make e2e-vpc-nat-gw-hairpin
set -Eeuo pipefail

NS=${NS:-hairpin-e2e}
GW_NAME=${GW_NAME:-hairpin-e2e-gw}
GW_NS=${GW_NS:-kube-system}
EXTERNAL_SUBNET=${EXTERNAL_SUBNET:-ovn-vpc-external-network}
EXTERNAL_NAD=${EXTERNAL_NAD:-kube-system/ovn-vpc-external-network}
VPC=${VPC:-hairpin-e2e-vpc}
SUBNET=${SUBNET:-hairpin-e2e-subnet}
SUBNET_CIDR=${SUBNET_CIDR:-10.77.0.0/24}
LAN_IP=${LAN_IP:-10.77.0.254}
BACKEND_IP=${BACKEND_IP:-10.77.0.11}
CLIENT_IP=${CLIENT_IP:-10.77.0.12}
VPC_INTERFACE=${VPC_INTERFACE:-eth0}
EXTERNAL_INTERFACE=${EXTERNAL_INTERFACE:-net1}
# Reuse an image the cluster already runs, so the test pulls nothing.
TEST_IMAGE=${TEST_IMAGE:-$(kubectl -n kube-system get ds kube-ovn-cni -o jsonpath='{.spec.template.spec.containers[0].image}')}

GW_POD=vpc-nat-gw-$GW_NAME-0
BACKEND_EIP_NAME=hairpin-backend-eip
CLIENT_EIP_NAME=hairpin-client-eip

fail() { echo "FAIL: $*" >&2; exit 1; }
step() { echo "== $*"; }

cleanup() {
    if [ -n "${local_server_pid:-}" ]; then
        gw_exec bash -c "kill $local_server_pid; rm -f /tmp/hairpin-local-test.log" >/dev/null 2>&1 || true
    fi
    kubectl delete --ignore-not-found --wait=false \
        iptables-fip-rules.kubeovn.io hairpin-backend-fip hairpin-client-fip >/dev/null 2>&1 || true
    sleep 2
    kubectl delete --ignore-not-found --wait=false \
        iptables-eips.kubeovn.io "$BACKEND_EIP_NAME" "$CLIENT_EIP_NAME" >/dev/null 2>&1 || true
    kubectl delete --ignore-not-found ns "$NS" --wait=false >/dev/null 2>&1 || true
    kubectl delete --ignore-not-found vpc-nat-gateways.kubeovn.io "$GW_NAME" --wait=false >/dev/null 2>&1 || true
    sleep 5
    kubectl delete --ignore-not-found subnets.kubeovn.io "$SUBNET" --wait=false >/dev/null 2>&1 || true
    kubectl delete --ignore-not-found vpcs.kubeovn.io "$VPC" --wait=false >/dev/null 2>&1 || true
}
trap cleanup EXIT

# The server answers with the peer address it sees: that is the SNAT source under test. The body
# is indented to sit inside the YAML block scalar of the Pod below.
source_echo_server() { # port
    cat <<EOF
      import socket
      s = socket.socket()
      s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
      s.bind(("0.0.0.0", $1))
      s.listen(8)
      while True:
          c, a = s.accept()
          c.sendall(a[0].encode())
          c.close()
EOF
}
seen_source() { # dst port
    kubectl exec -n "$NS" hairpin-client -c main -- python3 -c '
import socket,sys
s=socket.create_connection((sys.argv[1],int(sys.argv[2])),5)
print(s.recv(64).decode())
' "$1" "$2"
}
gw_exec() { kubectl exec -n "$GW_NS" "$GW_POD" -- "$@"; }
gw_rules() { gw_exec bash -c 'iptables-save -t nat | grep -E "^-A (HAIRPIN_SNAT|EXCLUSIVE_SNAT) " || true'; }
# Test-only flush on the dedicated gateway, not a production migration step.
# A new connection must not reuse the conntrack entry of a previous case.
flush_conntrack() { gw_exec bash -c 'conntrack -F >/dev/null 2>&1 || true'; }

wait_for() { # seconds description command...
    local deadline=$(( SECONDS + $1 )) desc=$2
    shift 2
    until "$@" >/dev/null 2>&1; do
        [ "$SECONDS" -lt "$deadline" ] || fail "timed out waiting for $desc"
        sleep 3
    done
}

cleanup
step "creating vpc, subnet and nat gateway"
kubectl apply -f - >/dev/null <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: $NS
---
apiVersion: kubeovn.io/v1
kind: Vpc
metadata:
  name: $VPC
spec:
  namespaces:
  - $NS
  staticRoutes:
  - cidr: 0.0.0.0/0
    nextHopIP: $LAN_IP
    policy: policyDst
---
apiVersion: kubeovn.io/v1
kind: Subnet
metadata:
  name: $SUBNET
spec:
  vpc: $VPC
  cidrBlock: $SUBNET_CIDR
  protocol: IPv4
---
apiVersion: kubeovn.io/v1
kind: VpcNatGateway
metadata:
  name: $GW_NAME
spec:
  vpc: $VPC
  subnet: $SUBNET
  lanIp: $LAN_IP
  selector:
  - "kubernetes.io/os: linux"
EOF
wait_for 300 "gateway pod $GW_NS/$GW_POD" \
    kubectl -n "$GW_NS" wait --for=condition=Ready "pod/$GW_POD" --timeout=10s

step "creating workloads and EIPs"
kubectl apply -f - >/dev/null <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: hairpin-backend
  namespace: $NS
  annotations:
    ovn.kubernetes.io/logical_switch: $SUBNET
    ovn.kubernetes.io/ip_address: $BACKEND_IP
spec:
  containers:
  - name: main
    image: $TEST_IMAGE
    imagePullPolicy: IfNotPresent
    command:
    - python3
    - -c
    - |
$(source_echo_server 8080)
---
apiVersion: v1
kind: Pod
metadata:
  name: hairpin-client
  namespace: $NS
  annotations:
    ovn.kubernetes.io/logical_switch: $SUBNET
    ovn.kubernetes.io/ip_address: $CLIENT_IP
spec:
  containers:
  - name: main
    image: $TEST_IMAGE
    imagePullPolicy: IfNotPresent
    command: ["sleep","infinity"]
---
apiVersion: v1
kind: Pod
metadata:
  name: hairpin-ext-peer
  namespace: $NS
  annotations:
    k8s.v1.cni.cncf.io/networks: $EXTERNAL_NAD
    ovn.kubernetes.io/logical_switch: $SUBNET
spec:
  containers:
  - name: main
    image: $TEST_IMAGE
    imagePullPolicy: IfNotPresent
    command:
    - python3
    - -c
    - |
$(source_echo_server 9090)
---
apiVersion: kubeovn.io/v1
kind: IptablesEIP
metadata:
  name: $BACKEND_EIP_NAME
spec:
  natGwDp: $GW_NAME
  externalSubnet: $EXTERNAL_SUBNET
---
apiVersion: kubeovn.io/v1
kind: IptablesEIP
metadata:
  name: $CLIENT_EIP_NAME
spec:
  natGwDp: $GW_NAME
  externalSubnet: $EXTERNAL_SUBNET
EOF
for pod in hairpin-backend hairpin-client hairpin-ext-peer; do
    wait_for 300 "pod $NS/$pod" kubectl -n "$NS" wait --for=condition=Ready "pod/$pod" --timeout=10s
done
eip_ip() { kubectl get iptables-eips.kubeovn.io "$1" -o jsonpath='{.status.ip}'; }
wait_for 300 "eip $BACKEND_EIP_NAME" bash -c "[ -n \"\$(kubectl get iptables-eips.kubeovn.io $BACKEND_EIP_NAME -o jsonpath='{.status.ip}')\" ]"
wait_for 300 "eip $CLIENT_EIP_NAME" bash -c "[ -n \"\$(kubectl get iptables-eips.kubeovn.io $CLIENT_EIP_NAME -o jsonpath='{.status.ip}')\" ]"
BACKEND_EIP=$(eip_ip "$BACKEND_EIP_NAME")
CLIENT_EIP=$(eip_ip "$CLIENT_EIP_NAME")
EXT_PEER_IP=$(kubectl -n "$NS" get pod hairpin-ext-peer \
    -o jsonpath="{.metadata.annotations.${EXTERNAL_NAD#*/}\.${EXTERNAL_NAD%/*}\.kubernetes\.io/ip_address}")
[ -n "$EXT_PEER_IP" ] || fail "the external peer has no address on $EXTERNAL_NAD"
echo "   backend eip $BACKEND_EIP, client eip $CLIENT_EIP, external peer $EXT_PEER_IP"

step "the hairpin chain holds the per-FIP jump, and the lanIP fallback is off by default"
rules=$(gw_rules)
[ "$(grep -c '^-A HAIRPIN_SNAT ' <<< "$rules")" = 1 ] || fail "expected one hairpin rule, got: $rules"
grep -q -- "^-A HAIRPIN_SNAT -o $VPC_INTERFACE -m mark --mark 0x1/0x1 -m conntrack --ctstate DNAT -j HAIRPIN_FIP" \
    <<< "$rules" || fail "the per-FIP jump is not installed: $rules"

step "binding the FIPs"
kubectl apply -f - >/dev/null <<EOF
apiVersion: kubeovn.io/v1
kind: IptablesFIPRule
metadata:
  name: hairpin-backend-fip
spec:
  eip: $BACKEND_EIP_NAME
  internalIp: $BACKEND_IP
EOF
wait_for 120 "backend fip rule" bash -c "kubectl exec -n $GW_NS $GW_POD -- iptables-save -t nat | grep -q -- '-A EXCLUSIVE_DNAT -d $BACKEND_EIP/32'"
gw_rules | grep -q -- "^-A EXCLUSIVE_SNAT -s $BACKEND_IP/32 -o $EXTERNAL_INTERFACE -j SNAT --to-source $BACKEND_EIP" \
    || fail "the FIP egress rule is not scoped to $EXTERNAL_INTERFACE: $(gw_rules)"

step "a VPC client reaching the backend's FIP is SNAT'd to that FIP"
gw_exec bash -c "iptables-save -t nat | grep -- '-A HAIRPIN_FIP ' | grep -q -- '--to-source $BACKEND_EIP'" \
    || fail "the backend FIP has no hairpin rule"
flush_conntrack
src=$(seen_source "$BACKEND_EIP" 8080)
[ "$src" = "$BACKEND_EIP" ] || fail "backend saw $src, want its own FIP $BACKEND_EIP"
echo "   backend saw $src"

step "the client gets its own FIP: the hairpin source stays the dialed FIP, not the client EIP"
kubectl apply -f - >/dev/null <<EOF
apiVersion: kubeovn.io/v1
kind: IptablesFIPRule
metadata:
  name: hairpin-client-fip
spec:
  eip: $CLIENT_EIP_NAME
  internalIp: $CLIENT_IP
EOF
wait_for 120 "client fip rule" bash -c "kubectl exec -n $GW_NS $GW_POD -- iptables-save -t nat | grep -q -- '-A EXCLUSIVE_SNAT -s $CLIENT_IP/32 -o $EXTERNAL_INTERFACE'"
flush_conntrack
src=$(seen_source "$BACKEND_EIP" 8080)
[ "$src" = "$BACKEND_EIP" ] || fail "backend saw $src, want the dialed FIP $BACKEND_EIP"
echo "   backend saw $src, not the client EIP"

step "egress is untouched: the same client reaching the external network keeps its EIP"
flush_conntrack
src=$(seen_source "$EXT_PEER_IP" 9090)
[ "$src" = "$CLIENT_EIP" ] || fail "external peer saw $src, want the client EIP $CLIENT_EIP"
echo "   external peer saw $src"

step "gateway-local INPUT services preserve the client address"
local_server_pid=$(gw_exec bash -c 'nc -n -lk -v -p 18080 </dev/null >/dev/null 2>/tmp/hairpin-local-test.log & echo $!')
wait_for 30 "gateway-local listener" gw_exec bash -c 'ss -ltn | grep -q :18080'
check_local_source() {
    gw_exec bash -c ': > /tmp/hairpin-local-test.log'
    kubectl exec -n "$NS" hairpin-client -c main -- python3 -c '
import socket,sys
s=socket.create_connection((sys.argv[1],18080),5)
s.sendall(b"local service test"); s.close()
' "$1"
    wait_for 30 "local service seeing the client address" gw_exec grep -F "Connection received on $CLIENT_IP " /tmp/hairpin-local-test.log
}
check_local_source "$LAN_IP"
# Port-scoped DNAT avoids changing the backend FIP's existing mapping.
gw_exec iptables -t nat -I DNAT_FILTER 1 -d "$BACKEND_EIP" -p tcp --dport 18080 -j DNAT --to-destination "$LAN_IP"
check_local_source "$BACKEND_EIP"
gw_exec iptables -t nat -D DNAT_FILTER -d "$BACKEND_EIP" -p tcp --dport 18080 -j DNAT --to-destination "$LAN_IP"
echo "   direct and DNAT-local services saw $CLIENT_IP"

step "a VIP with no EIP (ClusterIP, lanVIP) follows the lanIP fallback switch"
# Program via the real gateway script: this validates nft/iptables interoperability, not
# Service-controller feature gates (the test cluster may run an older controller).
CLUSTER_VIP=10.96.254.254
gw_exec bash /kube-ovn/nat-gateway.sh hairpin-snat-set enable >/dev/null
gw_exec bash /kube-ovn/nat-gateway.sh vip-addr-sync "$CLUSTER_VIP" >/dev/null
gw_exec bash /kube-ovn/nat-gateway.sh nft-dnat-map-add "$CLUSTER_VIP,8080,tcp,$BACKEND_IP:8080" >/dev/null
src=$(seen_source "$CLUSTER_VIP" 8080)
[ "$src" = "$LAN_IP" ] || fail "nft ClusterIP backend saw $src, want $LAN_IP"
gw_exec bash /kube-ovn/nat-gateway.sh nft-lanvip-sync "$LAN_IP,8080,tcp,$BACKEND_IP:8080" >/dev/null
src=$(seen_source "$LAN_IP" 8080)
[ "$src" = "$LAN_IP" ] || fail "nft lanVIP backend saw $src, want $LAN_IP"
gw_exec bash /kube-ovn/nat-gateway.sh nft-dnat-map-del "$CLUSTER_VIP,8080,tcp" >/dev/null
gw_exec bash /kube-ovn/nat-gateway.sh vip-addr-sync >/dev/null
gw_exec bash /kube-ovn/nat-gateway.sh nft-lanvip-sync >/dev/null
echo "   nft ClusterIP and lanVIP backends saw $LAN_IP"

step "the fallback is a switch, the per-FIP rules are not"
[ "$(gw_rules | grep -c '^-A HAIRPIN_SNAT ')" = 2 ] || fail "the fallback is not on: $(gw_rules)"
gw_exec bash /kube-ovn/nat-gateway.sh hairpin-snat-set disable >/dev/null
rules=$(gw_rules)
[ "$(grep -c '^-A HAIRPIN_SNAT ' <<< "$rules")" = 1 ] || fail "none must leave only the per-FIP jump: $rules"
grep -q -- "^-A HAIRPIN_SNAT -o $VPC_INTERFACE -m mark --mark 0x1/0x1 -m conntrack --ctstate DNAT -j HAIRPIN_FIP" \
    <<< "$rules" || fail "none dropped the per-FIP jump: $rules"
# FIP traffic is unaffected by the switch: it is still SNAT'd to the dialed FIP
flush_conntrack
src=$(seen_source "$BACKEND_EIP" 8080)
[ "$src" = "$BACKEND_EIP" ] || fail "with the fallback off the backend saw $src, want $BACKEND_EIP"
echo "   fallback off again, FIP traffic kept $src"

step "legacy rules of an older script, and their convergence on the next init"
gw_exec bash -c "
set -e
iptables -t nat -F HAIRPIN_SNAT
iptables -t nat -A HAIRPIN_SNAT -m mark --mark 0x1/0x1 -o $VPC_INTERFACE -m conntrack --ctstate DNAT \
    --ctorigdst $BACKEND_EIP -j SNAT --to-source $BACKEND_EIP --random-fully
iptables -t nat -D EXCLUSIVE_SNAT -o $EXTERNAL_INTERFACE -s $CLIENT_IP -j SNAT --to-source $CLIENT_EIP
iptables -t nat -A EXCLUSIVE_SNAT -s $CLIENT_IP -j SNAT --to-source $CLIENT_EIP" >/dev/null
flush_conntrack
src=$(seen_source "$BACKEND_EIP" 8080)
[ "$src" = "$CLIENT_EIP" ] || fail "expected the legacy rules to SNAT to $CLIENT_EIP, got $src"
echo "   legacy rules reproduced (backend saw the client EIP $src)"

# This is the natGwScriptHostPath workflow: the script is replaced under a running Pod and init
# is run by hand. A Pod only ever runs init once, so this is the only path that meets old rules.
gw_exec bash /kube-ovn/nat-gateway.sh init "$VPC_INTERFACE,$EXTERNAL_INTERFACE" >/dev/null
rules=$(gw_rules)
[ "$(grep -c '^-A HAIRPIN_SNAT ' <<< "$rules")" = 1 ] || fail "expected one hairpin rule after init, got: $rules"
grep -q -- "^-A HAIRPIN_SNAT -o $VPC_INTERFACE -m mark --mark 0x1/0x1 -m conntrack --ctstate DNAT -j HAIRPIN_FIP" \
    <<< "$rules" || fail "init did not rebuild the per-FIP jump: $rules"
grep -q -- "^-A EXCLUSIVE_SNAT -s $CLIENT_IP/32 -o $EXTERNAL_INTERFACE " <<< "$rules" \
    || fail "init did not scope the FIP egress rule: $rules"
# the per-FIP chain is rebuilt from the FIP map in EXCLUSIVE_DNAT, without the controller
# replaying a single FIP
gw_exec bash -c "iptables-save -t nat | grep -- '-A HAIRPIN_FIP ' | grep -q -- '--to-source $CLIENT_EIP'" \
    || fail "init did not rebuild the per-FIP hairpin rules"
flush_conntrack
src=$(seen_source "$BACKEND_EIP" 8080)
[ "$src" = "$BACKEND_EIP" ] || fail "after init the backend saw $src, want $BACKEND_EIP"
flush_conntrack
src=$(seen_source "$EXT_PEER_IP" 9090)
[ "$src" = "$CLIENT_EIP" ] || fail "after init the external peer saw $src, want $CLIENT_EIP"
echo "   converged: per-FIP rules rebuilt, fallback back to its default, egress untouched"

echo 'PASS: hairpin_snat_test.sh'
