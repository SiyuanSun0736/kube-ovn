# EIP TC bypass

The main Kube-OVN Dockerfile compiles `tc_eip.c` in a separate build stage
and installs `/etc/kube-ovn/tc_eip.o` for both amd64 and arm64 images.
The standalone build is `make -C dist/images/bpf` on Linux with clang,
libbpf headers and Linux UAPI headers installed.

Enable the speaker with `--node-route-eip-mode --enable-tc-bpf-redirect`.
It must run in the host network namespace with permission to load BPF programs,
manage TC filters, and write a mounted bpffs at `/sys/fs/bpf` (normally a
privileged speaker container with a host bpffs mount).

`--tc-bpf-interface` selects the physical ingress device; if omitted, the IPv4
default-route device is used. The gateway's external subnet independently selects
the macvlan master through its `ovn.kubernetes.io/nad-macvlan-master` annotation.
The host macvlan must use bridge mode. The gateway Pod's Multus network-status
must contain the external NAD and its actual MAC. Unresolved or invalid targets
are rejected, removed from the bypass map during reconciliation, and retried.
Inline VLAN headers are passed through without bypass.

`--tc-bpf-map-path` is the IPv4 pin filename, defaulting to
`/sys/fs/bpf/kube_ovn_eip_v4_map`. IPv6 uses the same filename with `_v4_map`
replaced by `_v6_map` (or `_v6_map` appended if that suffix is absent).
The speaker loads the object and pins the exact maps used by the attached program;
it does not discover arbitrary pinned maps by name.

The filter uses ingress priority 1, handle 1, and the name `kube-ovn-eip`.
An unrelated filter at priority 1 is an error, not something to overwrite.
Non-EIP traffic returns `TC_ACT_UNSPEC` (-1) so subsequent TC filters still run.
In direct-action mode, `TC_ACT_OK` (0) accepts the packet and ends TC
classification; it does not continue to lower-priority Cilium filters.
EIP packets are redirected through the host macvlan with the gateway Pod's MAC.
Local macvlan bridge delivery re-enters the master ingress; packets already
addressed to the Pod on that master return `TC_ACT_OK` to reach its macvlan RX
handler without a redirect loop or host IP-stack traversal.

Startup replaces old EIP attachments with empty maps and reconciles current
resources before announcing their BGP routes. Reconciliation retries loading or
attachment failures, updates changed targets, and deletes stale entries, including
entries whose EIPs moved to another node. Normal shutdown removes owned filters
and pins. To disable, keep node-route mode enabled and restart with
`--enable-tc-bpf-redirect=false`; the next speaker also cleans up attachments left
by an unclean exit. Do not run multiple node-route speakers on the same node.

## Kernel regression tests

Run these tests only inside a disposable privileged Linux container/network
namespace with bpffs, `ip`, and `ping` available. They create private interfaces
and namespaces, exercise actual BPF maps and program execution, and verify
macvlan delivery without a host route to the EIP.

```sh
make -C dist/images/bpf
KUBE_OVN_BPF_TEST=1 \
KUBE_OVN_BPF_OBJECT="$PWD/dist/images/bpf/tc_eip.o" \
go test ./pkg/speaker -run 'Test(TC|NatGwTarget)' -count=1 -v
```

Without `KUBE_OVN_BPF_TEST=1`, tests requiring kernel privileges are skipped.
