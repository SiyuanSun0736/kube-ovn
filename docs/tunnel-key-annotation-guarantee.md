# tunnel_key (VNI) annotation guarantee

Every non-hostNetwork pod on a non-vlan OVN VPC subnet carries a valid
`ovn.kubernetes.io/tunnel_key` annotation by the time its network is
configured (CNI ADD). Its value is the tunnel key of the pod's **VPC**, i.e.
the OVN logical router `Datapath_Binding.tunnel_key` recorded in
`Vpc.Status.TunnelKey`, not the subnet's logical switch key.

Why the VPC key: Cilium (native-vpc mode) scopes pod identities by this
value. A VPC is the IP uniqueness boundary (subnets inside one VPC cannot
overlap, different VPCs can), so a per-VPC scope lets Cilium resolve a peer
in any subnet of the same VPC with a single `(scope, IP)` lookup. A per-subnet
scope would need an extra subnet-to-subnet translation on every packet, whose
size grows with the square of the subnets per VPC. northd allocates switch and
router keys from one 24-bit namespace, so a VPC key never equals a subnet key.

`Subnet.Status.TunnelKey` (the logical switch key) is still synced before a
subnet is Ready and kept for observation; it no longer defines the pod
annotation. This includes both the default cluster-router VPC and
custom VPCs. The exact scope predicate is
`isOvnSubnet(subnet) && subnet.Spec.Vlan == ""`. VLAN/underlay and non-OVN
subnets keep the default
`status.tunnelKey=0` and do not carry the pod annotation. `Vpc.Status.TunnelKey`
is synced for every VPC because every VPC owns a logical router; VPCs imported
with `--enable-external-vpc` get it from the periodic external-VPC sync, which
only reads OVN SB. Cilium
(native-vpc mode) keys pod identities by this annotation, so a missing
annotation makes the pod fall back to the non-VPC scheme and collide with
overlapping VPC subnets.

The guarantee is upheld by two flows:

## Flow 1: normal pod creation

`handleAddOrUpdateVpc` (pkg/controller/vpc.go) reads the logical router key
from OVN SB and persists it in the same status patch that marks the VPC
`Standby`; subnets wait for `Standby`, so a VPC with subnets always has a key
in normal operation.

`reconcileAllocateSubnets` (pkg/controller/pod.go) refuses to complete the
allocation (requeues, persisting no pod annotations) while the resolved
subnet's VPC (`subnet.Spec.Vpc`) is unknown or its tunnel key is not valid
(`Vpc.Status.TunnelKey` is outside `1..16777215`), and writes the
tunnel_key annotation in the same atomic patch as `allocated=true`. For a
non-vlan OVN VPC subnet, the kube-ovn CNI server blocks until
`allocated=true` and the pod annotation is a valid value equal to
`Vpc.Status.TunnelKey` of the subnet's VPC. As an upgrade transition the
daemon also accepts the subnet's own key, which a controller predating the
VPC scope writes (see "Intended upgrade sequence"). Cilium (chained after
kube-ovn, the primary CNI) only runs once kube-ovn's CNI ADD has succeeded,
so by the time Cilium reads the pod, the annotation is correct.

For a multi-NIC pod, an earlier loop iteration may already have created
idempotent LSP/IP CR state before a later NIC hits the gate. This does not
weaken the CNI-visible guarantee: the pod annotation patch is committed only
after every requested OVN VPC NIC passes the gate, and retries safely reconcile
the earlier objects again.

Why flow 1 cannot produce a missing annotation: `reconcileAllocateSubnets`
returns an error without persisting pod annotations while any requested
OVN VPC subnet's VPC tunnel key is not ready; `tunnel_key` and `allocated=true`
are written in the same atomic patch; and the CNI server verifies both before
returning success. So under the current code a pod either carries the correct
annotation or is still waiting at CNI ADD - a missing annotation therefore
implies a legacy pod created before this guarantee existed, which is exactly
the population flow 2 repairs.

## Flow 2: kube-ovn-controller restart / legacy pods

- Pods created after the restart follow flow 1 (the gate lives in the
  allocation path and is unrelated to restarts).
- Pods whose persisted annotation is valid and matches their VPC status keep it.
- Legacy pods allocated before the VPC tunnel key was synced (or before
  this code existed) can have a missing, non-numeric or out-of-range
  annotation (valid OVN tunnel keys are `1..16777215`), or the subnet key
  written by controllers that predate the VPC scope. On startup,
  `InitIPAM` calls `repairPodTunnelKeyOnStartup`, which invokes
  `handleRepairTunnelKey` synchronously using the VPC's `Status.TunnelKey`.
  If it cannot finish, the dedicated queue retries indefinitely with
  exponential backoff capped at one minute and an aggregate 10-qps token
  bucket (VPC reconcile may still be syncing the key).

InitIPAM first attempts every eligible repair synchronously, using persisted
pod annotations (`podProvidersNeedingTunnelKeyRepair`) before NAD/default-
subnet resolution. If VPC status is already valid, the pod is fixed before
InitIPAM continues. Only an unavailable key or transient API error is queued.
That asynchronous fallback is necessary because InitIPAM runs before
`startWorkers`; a missing key may require the VPC worker, so waiting for it
inside InitIPAM would block the worker that produces it. On the first start
after upgrading to the VPC scope, every OVN VPC pod is in this state: the VPC
keys are not persisted yet, so all repairs are deferred and land once the VPC
workers sync the keys. When a VPC key first reaches the informer cache, pods
waiting for it (any requested provider not allocated yet, or an annotation
needing repair) are re-enqueued with a plain `Add`, bypassing their backoff.
The VPC informer callback only signals; one goroutine runs the pod sweep, and
a burst of VPCs syncing their keys (as on the first upgrade) coalesces into a
few sweeps instead of one per VPC.
The repair queue's limiter (per-item exponential backoff of 1 s to 1 min,
combined with a shared 10 qps / burst 100 bucket) paces only retries after a
failed repair, so it does not bound the backfill rate. The backfill is one
annotation patch per pod when the keys of all its providers' VPCs are ready (a
multi-NIC pod whose VPCs sync in batches is patched once per batch), run by `--worker-num` (default 3) repair workers in
parallel under the controller's client-side API limit (1000 qps, burst 2000);
its duration therefore depends on the worker count and the API server's patch
latency rather than on a fixed rate. Coalescing only merges notifications that
overlap a pending sweep, so the number of sweeps depends on how the VPC events
arrive; each sweep still lists every subnet once per pod with attachments. Measured: 8 pods in about 3 s after the key
sync on a single-node cluster. There is no other periodic task or pod-update
fallback. Detection and successful repair are
logged at Warning level, including what Cilium needs afterwards: a restart for
a backfilled (previously missing) key; a changed key is followed by running
endpoints (Cilium before qiniu/cilium#70 needs a pod recreation instead).

The deferred asynchronous repair is fallback recovery only; it is not part of
normal pod allocation. Normal allocation writes IP/network, tunnel_key and
`allocated=true` in one pod annotation patch.

Repair is multi-NIC aware and driven by the per-provider annotations the
allocation wrote (`allocated` + `logical_switch`). A provider is repaired
only when its logical switch resolves to a non-vlan OVN VPC subnet and its
key is missing, invalid or different from its VPC's `Status.TunnelKey`.
Providers managed by another CNI have no `logical_switch`; OVN
vlan/underlay subnets have one but are filtered by `isOvnVpcSubnet`. Both keep no pod tunnel_key annotation;
startup repair removes stale keys written by older controller versions.
The subnet is never guessed from namespace/default fallbacks, because
writing a wrong VNI is worse than a missing one (nothing would correct it
afterwards); a `logical_switch` that does not resolve is counted by
`pod_tunnel_key_repair_skipped_total` as the anomaly it is. Progress is
not all-or-nothing: ready providers are patched in one call even when
another provider's VPC key is still 0 (the handler then returns an
error to requeue for the remainder).

## Annotation protection

The annotation defines which VPC a pod's identity belongs to in Cilium, so a
user who can update pods must not be able to rewrite it: moving a pod into
another VPC's scope could take over the identity of an address in that VPC.
The ValidatingAdmissionPolicy `kube-ovn-tunnel-key-annotation` (installed by
both charts and install.sh when the API is available, Kubernetes >= 1.30)
rejects any pod UPDATE, including through the `pods/status` subresource (which
also persists annotations), that adds, changes or removes a
`*.kubernetes.io/tunnel_key` annotation unless the request comes from the
kube-ovn-controller service account (`system:serviceaccount:<namespace>:ovn`).
Writing the same value and changing other annotations or labels stay allowed.

Pod creation is not covered by the policy. For pods kube-ovn allocates, the
allocation patch overwrites any supplied value, and CNI ADD requires the key of
the VPC of the persisted logical switch. A creator who pre-sets `allocated`,
`logical_switch` and a matching key is the pre-existing trust question of all
kube-ovn network annotations (static IPs, subnets), outside this policy.
Without the policy (older Kubernetes) the annotation of a running pod is
unprotected. To repair an annotation manually, delete the
binding temporarily:

```bash
kubectl delete validatingadmissionpolicybinding kube-ovn-tunnel-key-annotation
```

## Failure mode and diagnosis

The allocation gate is intentionally fail-closed. If the VPC tunnel key of a
non-vlan OVN VPC subnet cannot be observed (for example, OVN SB is
unreachable or northd has not created the logical router's
`Datapath_Binding`), new pods on that subnet are not allocated with a guessed
or zero VNI:

- the pod remains in `ContainerCreating`;
- kube-ovn CNI waits for `allocated=true`, fails after its wait loop, and the
  kubelet retries CNI ADD;
- the pod receives an `AcquireAddressFailed` event containing
  `vpc <vpc> tunnel key of subnet <name> not observed on allocation`.

Check the VPC status with:

```bash
kubectl get vpc <vpc> -o jsonpath='{.status.tunnelKey}'
```

A value outside `1..16777215` is invalid and triggers fail-closed
re-synchronization. Check OVN SB reachability and whether the logical router
has a `Datapath_Binding` in SB (`ovn-sbctl find Datapath_Binding
external_ids:name=<vpc>`).
VLAN/underlay subnets are not subject to this gate and keep the default
TunnelKey value.

## Intended upgrade sequence

Upgrading from a subnet-scoped release:

0. Protect the annotation, then upgrade Cilium. A Cilium that follows
   tunnel_key changes on running endpoints (qiniu/cilium#70) trusts the
   annotation, so the protection must be enforced first (see "Annotation
   protection"):

   1. Apply the policy of the new release and check that it is enforced. It
      works with the current controller: only kube-ovn-controller writes the
      annotation, and the policy exempts its service account (`ovn`).

      ```bash
      # namespace: where kube-ovn-controller runs (its service account is exempt)
      helm template kube-ovn ./charts/kube-ovn --set namespace=kube-system \
        --api-versions admissionregistration.k8s.io/v1/ValidatingAdmissionPolicy \
        -s templates/tunnel-key-policy.yaml | kubectl apply -f -
      # must be denied (Forbidden): only kube-ovn-controller may change the key
      kubectl -n <ns> annotate pod <pod> ovn.kubernetes.io/tunnel_key=4095 \
        --overwrite --dry-run=server
      ```

      Kubernetes < 1.30 cannot install the policy. There, use a Cilium that
      follows key changes only if no untrusted user can update pods (RBAC);
      otherwise keep the older Cilium and recreate the pods after step 3.
   2. Upgrade Cilium to a release that follows tunnel_key changes on running
      endpoints. An older Cilium keeps existing endpoints in their subnet
      scope, even across a restart, so they would need to be recreated (see
      "Compatibility and rolling upgrade").
1. Apply the CRD (`Vpc.status.tunnelKey`) and RBAC (`vpcs` for kube-ovn-cni).
2. Upgrade kube-ovn-cni on every node. The new daemon accepts both the VPC key
   and the subnet key an older controller writes, so CNI ADD keeps working
   while the controller is still old. An old daemon accepts only the subnet
   key, so upgrading the controller first blocks CNI ADD of new pods on every
   node whose daemon is not upgraded yet.
3. Upgrade kube-ovn-controller. Flow 2 rewrites every OVN VPC pod annotation
   to its VPC key (see the backfill duration notes above), and Cilium moves
   each running endpoint to the new scope; no pod is recreated.

Every CNI ADD waits for
`allocated=true`; when the persisted logical_switch resolves to a non-vlan OVN
VPC subnet, it additionally waits for a valid annotation. If backfill has not
completed within the wait loop, ADD fails and kubelet retries.

## Compatibility and rolling upgrade

The Vpc API gains an optional `status.tunnelKey` integer; the Subnet API is
unchanged and pod annotation keys keep their existing provider-based format.
The kube-ovn-cni ClusterRole gains `get/list/watch` on `vpcs`: the daemon
waits for its VPC cache at startup, so apply the CRD and RBAC before rolling
out the new daemon.

Switching the annotation value from the subnet key to the VPC key changes
the identity scope of every OVN VPC pod in Cilium. Cilium with
qiniu/cilium#70 follows the change on running endpoints, at runtime and on
restore if the agent was down: the endpoint's VNI, identity, ipcache entry,
CiliumEndpoint and datapath configuration all move to the VPC key. Measured
(single node and two-node kind, subnet keys 11/12 -> VPC key 10):

- the controller rewrote every OVN VPC pod annotation within ~3-6 s of the VPC
  key sync, and every endpoint moved within ~5-12 s of its annotation;
- no CiliumEndpoint ever paired the new VNI with an identity of the old scope,
  and no entry stayed under a subnet key;
- afterwards all pods of the VPC with the same labels share one identity and
  cross-subnet security-group rules match.

Until both ends of a connection have moved, a lookup of one in the other's
scope misses and security-group rules that select endpoints deny the traffic,
for a few seconds per VPC.

An older Cilium adopts a key only when an endpoint is created; an existing
endpoint keeps the scope it was admitted with, including across a Cilium
restart. With it, a pod created after the backfill cannot reach not-yet-
recreated pods of its VPC, even in the same subnet, until every old pod is
recreated; protect the annotation and upgrade Cilium first (step 0) instead.

Two behavior changes are intentional:

- VLAN/underlay and non-OVN subnets now converge to
  `status.tunnelKey=0` and no pod tunnel_key annotation. Older controller
  versions may have populated those values even though native-vpc identity
  does not consume them; startup reconciliation removes those stale values.
- The new daemon fails CNI ADD for a provider whose logical_switch resolves
  to an OVN VPC subnet until `allocated=true` and the matching tunnel_key are
  visible. This converts a previously possible legacy race into a retryable
  CNI failure. A provider without logical_switch is not classified as VPC and
  keeps the old no-key behavior (required by NoDefaultEIP/IPAM-only networks).

For rolling upgrades follow "Intended upgrade sequence": annotation protection
and Cilium first, then the kube-ovn daemon before the controller. The daemon's acceptance of the subnet
key is a transition aid; it can be removed once no controller writes subnet
keys. Rolling the controller back to a subnet-scoped version rewrites the
annotations to subnet keys on its next startup, and Cilium endpoints follow
them back (an older Cilium needs the pods recreated again). An older controller
does not know `Vpc.status.tunnelKey` and drops it when it replaces a VPC
status; the new controller re-syncs it on its next start; an older controller
may also repopulate VLAN tunnel-key status, which a subsequent new-controller
startup cleans again.

Non-primary-CNI remains compatible: non-vlan OVN attachment providers use the
same per-provider annotation contract. IPAM-only and NoDefaultEIP attachments
without a persisted logical_switch bypass the VPC key gate, including legacy
provider naming.

## Known gaps (documented, not closed by code)

- `Vpc.Status.TunnelKey` and `Subnet.Status.TunnelKey` are not periodically
  compared with OVN SB. This is intentional because a `Datapath_Binding` key
  is stable during the normal lifecycle. After a destructive SB rebuild or an
  inconsistent DR restore, northd may assign a different key while the
  persisted status still contains an otherwise valid old value;
  `reconcileVpcTunnelKey` and `reconcileSubnetTunnelKey` will not re-read SB
  until that status is invalidated. Recover every affected VPC (pod
  annotations) and subnet (observation only) as follows:

  1. Clear its persisted key:

     ```bash
     kubectl patch vpc <name> --subresource=status --type=merge \
       -p '{"status":{"tunnelKey":0}}'
     kubectl patch subnet <name> --subresource=status --type=merge \
       -p '{"status":{"tunnelKey":0}}'
     ```

  2. Restart kube-ovn-controller. InitIPAM detects affected pods; after VPC
     workers re-read the new SB key, the deferred repair queue updates their
     annotations. Wait for the startup repair Warning logs to complete.
  3. Cilium moves the running endpoints of the affected pods to the new key.
     With Cilium before qiniu/cilium#70, recreate the affected pods instead:
     an existing endpoint does not adopt a changed key, not even after a
     Cilium restart.

  A controller restart without first invalidating status is insufficient.
- Repair progress is observable via `pod_tunnel_key_repair_patch_total`
  (annotation patches that added, corrected or removed one or more tunnel_key
  annotations; a multi-NIC pod can be patched more than once) and
  `pod_tunnel_key_repair_skipped_total` (repairs skipped because the
  logical_switch subnet no longer exists) - see
  pkg/controller/tunnel_key_metrics.go.
- Repair keying is per provider parsed from the `*.kubernetes.io/allocated`
  annotation suffix (the same `*.kubernetes.io/*` templates as the
  allocation path). Cilium currently only recognizes the primary NIC
  (default provider, `ovn.kubernetes.io/*` annotations), so a mismatch on
  a multus secondary NIC would be harmless today; it only matters if
  Cilium starts keying secondary NICs.
- A pod whose `logical_switch` subnet no longer exists is skipped (counted
  by `pod_tunnel_key_repair_skipped_total`) and not retried: the repair
  queue forgets it, and there is no periodic or pod-update resync. It is
  reconsidered on the next controller restart - acceptable because the pod
  usually cannot survive its subnet's deletion anyway, and accidental manual
  deletion is outside this guarantee.
- The upgrade sequence above is operator discipline, not enforced by code:
  upgrading the controller before the daemons blocks CNI ADD on nodes with an
  old daemon, and restarting cilium before the backfill completes leaves
  already-created endpoints on the non-VPC scheme until the next cilium restart. This affects
  only endpoints that already exist; CNI ADD itself is fail-closed on the key.
