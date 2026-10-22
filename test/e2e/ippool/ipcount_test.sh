#!/usr/bin/env bash
# Control-plane test of ippool spec.ipCount against a live cluster: a real Subnet, real IPPool
# objects and a real workload Pod, so the assertions come from what the controller actually
# committed to IPAM, OVN and the ippool status.
#
# Preconditions (not created here, they belong to the cluster):
#   - a kube-ovn controller supporting spec.ipCount; the test skips on an older CRD
#   - the default VPC accepts another subnet with the CIDR below
# Reserved ipcount-e2e resource names must not be used by other workloads.
#
# Run with: make e2e-ippool-ipcount
set -Eeuo pipefail

NS=${NS:-ipcount-e2e}
SUBNET=${SUBNET:-ipcount-e2e-subnet}
SUBNET_CIDR=${SUBNET_CIDR:-10.79.0.0/24}
SUBNET_GW=${SUBNET_GW:-10.79.0.1}
POOL=${POOL:-ipcount-e2e-pool}
POD=${POD:-ipcount-e2e-pod}
TIMEOUT=${TIMEOUT:-90}
# Reuse an image the cluster already runs, so the test pulls nothing.
TEST_IMAGE=${TEST_IMAGE:-$(kubectl -n kube-system get ds kube-ovn-cni -o jsonpath='{.spec.template.spec.containers[0].image}')}

fail() { echo "FAIL: $*" >&2; exit 1; }
step() { echo "== $*"; }

cleanup() {
  kubectl delete pod "$POD" ipcount-e2e-probe ipcount-e2e-restart -n "$NS" --ignore-not-found --wait=false >/dev/null 2>&1 || true
  kubectl delete ippool "$POOL" --ignore-not-found --timeout=60s >/dev/null 2>&1 || true
  kubectl delete ns "$NS" --ignore-not-found --wait=false >/dev/null 2>&1 || true
  kubectl delete subnet "$SUBNET" --ignore-not-found --timeout=60s >/dev/null 2>&1 || true
}
trap cleanup EXIT

# The shell equivalent of f.SkipVersionPriorTo: the feature is identified by the CRD field it
# needs, so the test skips instead of failing on a branch that does not carry it yet.
if [ -z "$(kubectl get crd ippools.kubeovn.io -o jsonpath='{.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.ipCount}')" ]; then
  echo "SKIP: the ippools CRD has no spec.ipCount"
  exit 0
fi

poolIPs() { kubectl get ippool "$POOL" -o jsonpath='{.status.ips}'; }
# the members as whitespace separated words, for exact comparisons - 10.79.0.2 is a substring of
# 10.79.0.20, so matching against the json array would report false hits
poolMembers() { kubectl get ippool "$POOL" -o jsonpath='{.status.ips[*]}'; }
poolReady() { kubectl get ippool "$POOL" -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}'; }
podIP() { kubectl get pod "$POD" -n "$NS" -o jsonpath='{.metadata.annotations.ovn\.kubernetes\.io/ip_address}'; }

# waitFor <description> <command...>: retries until the command succeeds or TIMEOUT expires.
waitFor() {
  local what="$1"; shift
  local deadline=$((SECONDS + TIMEOUT))
  until "$@"; do
    [ $SECONDS -lt $deadline ] || fail "timed out waiting for $what (pool ips: $(poolIPs), ready: $(poolReady))"
    sleep 2
  done
}

poolHasIPs() { [ "$(poolIPs)" = "$1" ]; }
poolHasMemberCount() { [ "$(poolMembers | wc -w)" -eq "$1" ]; }
poolIsReady() { [ "$(poolReady)" = "True" ]; }

setPoolIPCount() {
  kubectl patch ippool "$POOL" --type=merge -p "{\"spec\":{\"ips\":null,\"ipCount\":$1}}" >/dev/null
}

cleanup
step "create subnet $SUBNET ($SUBNET_CIDR) and namespace $NS"
kubectl create ns "$NS" >/dev/null
kubectl apply -f - >/dev/null <<EOF
apiVersion: kubeovn.io/v1
kind: Subnet
metadata:
  name: $SUBNET
spec:
  protocol: IPv4
  cidrBlock: $SUBNET_CIDR
  gateway: $SUBNET_GW
  excludeIps:
  - $SUBNET_GW
  namespaces:
  - $NS
EOF
waitFor "subnet $SUBNET" bash -c "[ \"\$(kubectl get subnet $SUBNET -o jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}')\" = True ]"

step "create ippool $POOL with ipCount 2, with an address set"
kubectl apply -f - >/dev/null <<EOF
apiVersion: kubeovn.io/v1
kind: IPPool
metadata:
  name: $POOL
spec:
  subnet: $SUBNET
  ipCount: 2
  enableAddressSet: true
EOF
waitFor "ippool $POOL to become ready" poolIsReady
ips=$(poolIPs)
[ "$ips" = '["10.79.0.2","10.79.0.3"]' ] || fail "expected the first two free addresses, got $ips"
[ "$(kubectl get ippool $POOL -o jsonpath='{.status.v4AvailableIPs}')" = 2 ] || fail "expected 2 available addresses"

step "the address set holds exactly the selected addresses"
# ovn rejects '-' in an address set name, so the controller maps it to '.'
asName=$(echo "$POOL" | tr '-' '.')
as=$(kubectl -n kube-system exec deploy/ovn-central -- ovn-nbctl --bare --columns=addresses find address_set "name=$asName" 2>/dev/null | tr -d '\r' | tr ' ' '\n' | grep . | sort | paste -sd,)
[ "$as" = "10.79.0.2,10.79.0.3" ] || fail "address set holds [$as], expected 10.79.0.2,10.79.0.3"

step "a pod bound to the pool gets one of its addresses"
kubectl apply -f - >/dev/null <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: $POD
  namespace: $NS
  annotations:
    ovn.kubernetes.io/logical_switch: $SUBNET
    ovn.kubernetes.io/ip_pool: $POOL
spec:
  containers:
  - name: pause
    image: $TEST_IMAGE
    command: ["sleep", "3600"]
  terminationGracePeriodSeconds: 0
EOF
waitFor "pod $POD to get an address" bash -c "[ -n \"\$(kubectl get pod $POD -n $NS -o jsonpath='{.metadata.annotations.ovn\.kubernetes\.io/ip_address}')\" ]"
podAddr=$(podIP)
case "$podAddr" in
  10.79.0.2|10.79.0.3) ;;
  *) fail "pod got $podAddr, which is not a member of the pool" ;;
esac

step "growing ipCount keeps the current addresses and appends new ones"
setPoolIPCount 4
waitFor "the pool to grow" poolHasIPs '["10.79.0.2","10.79.0.3","10.79.0.4","10.79.0.5"]'
waitFor "the grown pool to become ready" poolIsReady

step "shrinking releases the most recently added addresses first"
setPoolIPCount 2
waitFor "the pool to shrink" poolHasIPs '["10.79.0.2","10.79.0.3"]'
[ "$(podIP)" = "$podAddr" ] || fail "the pod lost its address on shrink: $(podIP)"

step "the released addresses are allocatable from the subnet again"
# the pool keeps 10.79.0.2 and 10.79.0.3, the gateway is excluded, so a pod allocating from the
# subnet must get 10.79.0.4 - the address the shrink just released - and not a pool member
kubectl -n "$NS" run ipcount-e2e-probe --image="$TEST_IMAGE" --restart=Never --command -- sleep 3600 >/dev/null
waitFor "the probe pod to get an address" bash -c "[ -n \"\$(kubectl get pod ipcount-e2e-probe -n $NS -o jsonpath='{.metadata.annotations.ovn\.kubernetes\.io/ip_address}')\" ]"
probeAddr=$(kubectl get pod ipcount-e2e-probe -n "$NS" -o jsonpath='{.metadata.annotations.ovn\.kubernetes\.io/ip_address}')
kubectl delete pod ipcount-e2e-probe -n "$NS" --wait=false >/dev/null
[ "$probeAddr" = 10.79.0.4 ] || fail "the probe pod got $probeAddr, expected the released 10.79.0.4"

step "the selection survives a controller restart, restored from status.ips"
# rollout restart replaces the pod, so rollout status only returns once the restarted
# controller is up - waiting for any ready pod could match the one being deleted and let
# the assertions below run against the old process
kubectl -n kube-system rollout restart deployment/kube-ovn-controller >/dev/null
waitFor "the controller to come back" kubectl -n kube-system rollout status deployment/kube-ovn-controller --timeout=10s
# the restore reserves the recorded addresses and the reconcile activates them again, so
# a pod of the pool gets one of them once the pool is ready
waitFor "the restored pool to become ready" poolIsReady
kubectl apply -f - >/dev/null <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: ipcount-e2e-restart
  namespace: $NS
  annotations:
    ovn.kubernetes.io/logical_switch: $SUBNET
    ovn.kubernetes.io/ip_pool: $POOL
spec:
  containers:
  - name: pause
    image: $TEST_IMAGE
    command: ["sleep", "3600"]
  terminationGracePeriodSeconds: 0
EOF
waitFor "the pod after the restart to get an address" bash -c "[ -n \"\$(kubectl get pod ipcount-e2e-restart -n $NS -o jsonpath='{.metadata.annotations.ovn\.kubernetes\.io/ip_address}')\" ]"
restartAddr=$(kubectl get pod ipcount-e2e-restart -n "$NS" -o jsonpath='{.metadata.annotations.ovn\.kubernetes\.io/ip_address}')
kubectl delete pod ipcount-e2e-restart -n "$NS" --wait=false >/dev/null
found=no
for member in $(poolMembers); do [ "$member" != "$restartAddr" ] || found=yes; done
[ "$found" = yes ] || fail "after the restart the pod got $restartAddr, not a member of $(poolIPs)"

setPoolIPCount 3
waitFor "the pool to grow after the restart" poolHasIPs '["10.79.0.2","10.79.0.3","10.79.0.4"]'

step "switching to spec.ips clears the recorded addresses"
kubectl patch ippool "$POOL" --type=merge -p '{"spec":{"ipCount":null,"ips":["10.79.0.20"]}}' >/dev/null
waitFor "the recorded addresses to be cleared" bash -c "[ -z \"\$(kubectl get ippool $POOL -o jsonpath='{.status.ips}')\" ]"
waitFor "the spec.ips pool to become ready" poolIsReady

step "switching back to spec.ipCount keeps the current address and tops up"
setPoolIPCount 2
waitFor "the pool to select again" poolIsReady
waitFor "the pool to hold two addresses" poolHasMemberCount 2
[ "$(kubectl get ippool "$POOL" -o jsonpath='{.status.ips[0]}')" = 10.79.0.20 ] ||
  fail "expected the address the pool already holds to be kept first, got $(poolIPs)"
# the pod still holds its address from outside the pool, so the top up must skip it
for member in $(poolMembers); do
  [ "$member" != "$podAddr" ] || fail "the selection took $podAddr, which the pod is still using: $(poolIPs)"
done

step "deleting the pool returns its addresses to the subnet"
kubectl delete pod "$POD" -n "$NS" --wait=false >/dev/null
kubectl delete ippool "$POOL" --timeout=60s >/dev/null
waitFor "the pool to disappear" bash -c "! kubectl get ippool $POOL >/dev/null 2>&1"
as=$(kubectl -n kube-system exec deploy/ovn-central -- ovn-nbctl --bare --columns=name find address_set "name=$asName" 2>/dev/null | tr -d '\r\n ')
[ -z "$as" ] || fail "the address set of the deleted pool is still there: $as"

echo "PASS: ipcount_test.sh"
