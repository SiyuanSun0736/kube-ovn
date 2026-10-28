#!/usr/bin/env bash
set -euo pipefail

CNI_SOCK=/run/openvswitch/kube-ovn-daemon.sock
OVS_SOCK=/run/openvswitch/db.sock
ENABLE_SSL=${ENABLE_SSL:-false}

function quit {
  rm -rf $CNI_SOCK
  exit 0
}
trap quit EXIT

if [[ -e "$CNI_SOCK" ]]
then
  echo "previous socket exists, remove and continue"
  rm ${CNI_SOCK}
fi

while true; do
  if [[ -e "$OVS_SOCK" ]]; then
    for component in ovsdb-server ovs-vswitchd; do
      echo "checking ${component} status"
      pid_file="/run/openvswitch/${component}.pid"
      if ! read -r pid _ < "$pid_file" 2>/dev/null || [[ -z "$pid" ]]; then
        echo "${component} is not ready (failed to read pid or pid is empty)"
        sleep 1
        continue 2
      fi
      if ! ovs-appctl -T 1 -t "/run/openvswitch/${component}.${pid}.ctl" version >/dev/null 2>&1; then
        echo "${component} is not ready (ovs-appctl failed)"
        sleep 1
        continue 2
      fi
    done
    break
  else
    echo "waiting for ovs ready"
    sleep 1
  fi
done

# update links to point to the iptables binaries
iptables -V

# If nftables not exist do not exit
set +e
iptables -P FORWARD ACCEPT
iptables-nft -P FORWARD ACCEPT

# Bypass conntrack for cross-node tunnel traffic (Geneve 6081 / VXLAN 4789, 8472)
for port in 6081 4789 8472; do
  for cmd in iptables iptables-nft; do
    if command -v "${cmd}" >/dev/null 2>&1; then
      "${cmd}" -w 5 -t raw -C PREROUTING -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null || \
        "${cmd}" -w 5 -t raw -A PREROUTING -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null
      "${cmd}" -w 5 -t raw -C OUTPUT -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null || \
        "${cmd}" -w 5 -t raw -A OUTPUT -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null
    fi
  done

  for cmd in ip6tables ip6tables-nft; do
    if command -v "${cmd}" >/dev/null 2>&1; then
      "${cmd}" -w 5 -t raw -C PREROUTING -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null || \
        "${cmd}" -w 5 -t raw -A PREROUTING -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null
      "${cmd}" -w 5 -t raw -C OUTPUT -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null || \
        "${cmd}" -w 5 -t raw -A OUTPUT -p udp -m udp --dport "${port}" -j NOTRACK 2>/dev/null
    fi
  done
done
set -e

./kube-ovn-daemon --ovs-socket=${OVS_SOCK} --bind-socket=${CNI_SOCK} "$@"
