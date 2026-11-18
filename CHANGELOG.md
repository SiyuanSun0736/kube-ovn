# Changelog

## v1.17.0 (2026-11-07)

 * [3a5a5d231](https://github.com/kubeovn/kube-ovn/commit/3a5a5d231398241717464612f64948c161b50acb) Update module github.com/spf13/pflag to v1.1.0 (#7893)
 * [04c6136ff](https://github.com/kubeovn/kube-ovn/commit/04c6136ffba3166f18fe3537f1823eacfa536a43) chore: upgrade network-policy-api to v0.3.0 and CNP to v1beta1 (#7892)
 * [00bedc5ec](https://github.com/kubeovn/kube-ovn/commit/00bedc5ec6b70d35b163b8cc20b129b8b9722994) test(e2e): avoid assuming sequential IPs after deployment restart (#7890)
 * [e213e2f89](https://github.com/kubeovn/kube-ovn/commit/e213e2f89bcd62f1141008c05779b9ced91e3d2a) test(e2e): wait for IC dataplane convergence (#7889)
 * [da99aebe8](https://github.com/kubeovn/kube-ovn/commit/da99aebe88c0cf91352926b351a3b971f07a778e) fix: rebuild corrupted OVN members from Raft headers (#7888)
 * [3fa1c2a01](https://github.com/kubeovn/kube-ovn/commit/3fa1c2a0165f4383d200cd771be46878387ae019) test(e2e): wait for external gateway disable before EIP cleanup (#7887)
 * [d7dca9fdc](https://github.com/kubeovn/kube-ovn/commit/d7dca9fdca8d17367c445a7c53940fb5103429f6) fix(controller): refresh template load balancer updates on backend changes (#7884)
 * [3cde2bb63](https://github.com/kubeovn/kube-ovn/commit/3cde2bb6392a0922b0dcfc2941a5d087eda469aa) fix: serialize ProviderNetwork lifecycle cleanup (#7880)
 * [e663f36a1](https://github.com/kubeovn/kube-ovn/commit/e663f36a11859b19381ccef0dd8718854a08a71d) fix: rebuild invalid OVN members from live Raft cluster (#7871)
 * [fc738e18c](https://github.com/kubeovn/kube-ovn/commit/fc738e18c0ed136cc9c1df44058a0318b4b4ec21) fix(controller): reconcile security groups for retained pod ports (#7875)
 * [e4f984ae5](https://github.com/kubeovn/kube-ovn/commit/e4f984ae5813c21318bc85f0e1fd3f9af80bd434) fix(controller): run service load balancer gateways as non-root (#7876)
 * [393036c46](https://github.com/kubeovn/kube-ovn/commit/393036c46fa4094df1b080e690fd9af90cb87260) fix(controller): include VIPs in security groups without port security (#7879)
 * [715ffb575](https://github.com/kubeovn/kube-ovn/commit/715ffb57520547cfbec4ba1b4da7251a3c7c79ae) test: synchronize fake informer watches and fixtures (#7874)
 * [a7c5726a7](https://github.com/kubeovn/kube-ovn/commit/a7c5726a71023279933c229860a0b739159ede83) feat: make geneve tx offload configurable (#7873)
 * [039b6e835](https://github.com/kubeovn/kube-ovn/commit/039b6e835458e1a1396ea9912a2346c3869898f6) chore(deps): update dependency multus to v4.3.2 (#7877)
 * [4f70f33c4](https://github.com/kubeovn/kube-ovn/commit/4f70f33c432333ac15062391fa842e18ed94d59b) fix(deps): update module gopkg.in/k8snetworkplumbingwg/multus-cni.v4 to v4.3.2 (#7878)
 * [aa22460c8](https://github.com/kubeovn/kube-ovn/commit/aa22460c8cb152227cd6e88284ec1a838b4e9f68) fix(ovs): atomically create logical switch ports (#7865)
 * [fdc272147](https://github.com/kubeovn/kube-ovn/commit/fdc2721475a07606b776a2364de5911e2d52920f) refactor: migrate OVN/OVS consumers to shared libovsdb table layer (#7363)
 * [fe9b6122e](https://github.com/kubeovn/kube-ovn/commit/fe9b6122e965486c854d23f47f54f10ab848e37e) chore(deps): update helm/kind-action action to v1.15.2 (#7868)
 * [2b38f86b0](https://github.com/kubeovn/kube-ovn/commit/2b38f86b0b24833a991bcdb90c30717174507852) fix(deps): update module go.universe.tf/metallb to v0.16.2 (#7869)
 * [d3cccf24c](https://github.com/kubeovn/kube-ovn/commit/d3cccf24cb3427d27597bd1c4baa1df521208cb4) feat(ippool): allocate a fixed number of addresses from the subnet (#7702)
 * [097601062](https://github.com/kubeovn/kube-ovn/commit/09760106297100fdbcca290e6bbd7c5838025dea) Update dependency metallb to v0.16.2 (#7867)
 * [254ae76f8](https://github.com/kubeovn/kube-ovn/commit/254ae76f8ce889984c9d4126d0bac46a376a0f2b) test: wait for IPsec overlay connectivity convergence (#7860)
 * [2e9d43ff2](https://github.com/kubeovn/kube-ovn/commit/2e9d43ff23a7544a87b12147fdf4aa4445aa31e0) fix: preserve concurrent VPC status updates for external gateways (#7864)
 * [2dfa68ef0](https://github.com/kubeovn/kube-ovn/commit/2dfa68ef02eecd98795454c329965ea7b9211b78) fix: use service ports for pod load balancer DNAT (#7862)
 * [1c79eb8ad](https://github.com/kubeovn/kube-ovn/commit/1c79eb8ad3250abd78b8baf1aac03915815f8c55) test: cover netlink offload handling (#7859)
 * [277cb1513](https://github.com/kubeovn/kube-ovn/commit/277cb1513ff039bd18c1a436942d630d27e3eead) refactor: replace Felix IP sets with Netlink (#7852)
 * [eb9af3c85](https://github.com/kubeovn/kube-ovn/commit/eb9af3c85a7ece219f3c0b7dbabda04cc7c0eb08) fix(security-group): fail closed before port group exists (#7849)
 * [bacf72fbd](https://github.com/kubeovn/kube-ovn/commit/bacf72fbd94c5ba5f994a8e0396105ba6795e3ae) refactor: configure NIC offloads through netlink (#7855)
 * [c469fb185](https://github.com/kubeovn/kube-ovn/commit/c469fb1853332fe393aab39df591f398e02933c4) chore(deps): update go-iptables to rebased upstream HEAD (#7853)
 * [7e5560e85](https://github.com/kubeovn/kube-ovn/commit/7e5560e85f1a716c4977a9d0f1edd5013321a77f) fix(deps): update module github.com/ovn-kubernetes/libovsdb to v0.9.1 (master) (#7847)
 * [ef4e9538b](https://github.com/kubeovn/kube-ovn/commit/ef4e9538bf434dec8477076e9d5c0e7c277c1beb) feat(speaker): support explicit next hops for advertised routes (#7834)
 * [166dfc3d5](https://github.com/kubeovn/kube-ovn/commit/166dfc3d5cb7b71330fea97d407c54d6131b7deb) fix(ci): install Trivy directly from GitHub releases (#7854)
 * [1e7bf6c98](https://github.com/kubeovn/kube-ovn/commit/1e7bf6c98303e7adc4c33284eab97d7cfc4ee639) fix(ci): update kubectl-ko network inspect invocation (#7851)
 * [6c8432b9e](https://github.com/kubeovn/kube-ovn/commit/6c8432b9e6fff7b293720037ecf4d339a003101f) test: gate non-primary VPC egress gateway E2E by version (#7848)
 * [ad3927137](https://github.com/kubeovn/kube-ovn/commit/ad3927137dcca1eec8169ecc2d5cd0bfe0120ae3) feat(chart): make v2 Deployment tolerations configurable (#7844)
 * [2d70ca27a](https://github.com/kubeovn/kube-ovn/commit/2d70ca27aefaff89f682a61f6a337ebeb42645d8) fix(webhook): avoid stale Service conflict checks (#7842)
 * [db2102896](https://github.com/kubeovn/kube-ovn/commit/db210289683612d374d487e43783e712b6c02057) fix(deps): update golang.org/x/crypto (#7838)
 * [68b9a4f19](https://github.com/kubeovn/kube-ovn/commit/68b9a4f1957eee26516b8bf3e89cc08446fc05a1) feat(ko): share source references between trace and capture (#7826)
 * [4a29886cc](https://github.com/kubeovn/kube-ovn/commit/4a29886cc659e6903d6262def514367e83346bfc) feat(vpc-egress-gateway): support primary CNI egress (#7824)
 * [6f8d27fd3](https://github.com/kubeovn/kube-ovn/commit/6f8d27fd3684cc731830adaee0d1b6c41e9865db) test: isolate controller during IPsec disable probe (#7836)
 * [089606d47](https://github.com/kubeovn/kube-ovn/commit/089606d47d2ed8011fb4242fad41abdbec3fb6fc) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.34.0 (#7835)
 * [39ee464d4](https://github.com/kubeovn/kube-ovn/commit/39ee464d4027fa7bbae16f83a0a1f2d5f7b0ff91) Capture OVN database snapshots for HA E2E failures (#7832)
 * [0bb64d8b9](https://github.com/kubeovn/kube-ovn/commit/0bb64d8b989c0b3d619ccb881b2a6f04837fda17) Fix IPsec overlay SB state convergence check (#7823)
 * [2b3e29f62](https://github.com/kubeovn/kube-ovn/commit/2b3e29f62fa3c6aa894522889da099cd86bae07d) fix: support CNP conformance CRDs, priorities, and selectors (#7820)
 * [ad96b6aa5](https://github.com/kubeovn/kube-ovn/commit/ad96b6aa5367e9ada0d1ada5d898a8047829fd1f) docs: updated CHANGELOG.md (#7827)
 * [dfd83938e](https://github.com/kubeovn/kube-ovn/commit/dfd83938e61940b580e7f184ea29145fae77004b) chore(deps): update dependency gobgp to v4.11.0 (#7828)
 * [e07afca09](https://github.com/kubeovn/kube-ovn/commit/e07afca09f5fcb2107a5eaaf87825f653c3bacd1) fix(deps): update module github.com/osrg/gobgp/v4 to v4.11.0 (#7829)
 * [5bc57de87](https://github.com/kubeovn/kube-ovn/commit/5bc57de87bf99ced2f3d034eded465c9dc2f15ec) Update module github.com/cenkalti/backoff/v7 to v7.0.3 (#7831)
 * [5658834b9](https://github.com/kubeovn/kube-ovn/commit/5658834b96e1f4064bf6186d5adb6f19c2866009) feat(vpc-nat-gw): serve Services through nft share DNAT (#7451)
 * [da8e043b8](https://github.com/kubeovn/kube-ovn/commit/da8e043b81a578325e0378f24064e749b3023884) fix: honor legacy OVN attachment subnet selection (#7806)
 * [e701cc73c](https://github.com/kubeovn/kube-ovn/commit/e701cc73caab5929ae13fe3201a6b359ab34884b) test(e2e): wait for VIP load balancer attachments (#7818)
 * [a1a616c30](https://github.com/kubeovn/kube-ovn/commit/a1a616c301587afce2b477989cd7b0b2b3ebfad8) test: wait for OVN leaders after HA database recovery (#7819)
 * [86dc32242](https://github.com/kubeovn/kube-ovn/commit/86dc322426462bf954dc2bf00b3260f6648bbb15) feat(netpol): support named ports in ANP, BANP and CNP (#7807)
 * [00363e9c6](https://github.com/kubeovn/kube-ovn/commit/00363e9c69c21a7602a61af04a26621c37de148a) fix(controller): preserve externally managed OVN chassis in external VPC mode (#7684)
 * [9d95ca893](https://github.com/kubeovn/kube-ovn/commit/9d95ca89346c91517cd4c0ab1a26294c55bc8160) ci: remove Submariner E2E from x86 workflow (#7817)
 * [f9b8e6706](https://github.com/kubeovn/kube-ovn/commit/f9b8e6706818258cda78426d5753bc0f3c2be0ba) chore(deps): update dependency network-policy-api-cnp to v0.2.1 (master) (#7808)
 * [fa8f99062](https://github.com/kubeovn/kube-ovn/commit/fa8f9906274d590ae668ca8895f024953f902203) chore(deps): update helm/chart-releaser-action action to v1.9.0 (#7815)
 * [6d5a928a3](https://github.com/kubeovn/kube-ovn/commit/6d5a928a3434ed61f90d477e10209c2163a439bb) docs: update NEXT.md with Post-v1.16.0 features
 * [54d37495f](https://github.com/kubeovn/kube-ovn/commit/54d37495fbd4a546378d9cb769f80f81ba93d666) docs: updated CHANGELOG.md (#7794)
 * [cd93f7f22](https://github.com/kubeovn/kube-ovn/commit/cd93f7f22275cee7d5a835ca07545a40e9427748) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.4.2 (#7809)
 * [b7b28e317](https://github.com/kubeovn/kube-ovn/commit/b7b28e317e086adfbf565a25b57e3776e1ecdc2c) fix(deps): update module cel.dev/cel-go to v0.33.2 (#7811)
 * [371e6294b](https://github.com/kubeovn/kube-ovn/commit/371e6294bd77e375f6f20c2c1aa6af25e153c12f) docs: fix dead User Guide link in README (#7810)
 * [b51870336](https://github.com/kubeovn/kube-ovn/commit/b51870336ef8f2d5e48b64e294f76550439ed6ed) Update module github.com/brianvoe/gofakeit/v7 to v7.18.1 (#7798)
 * [829f8fffe](https://github.com/kubeovn/kube-ovn/commit/829f8fffeee1c7e149ddfb06c671bfd4a474d5de) feat(netpol): match named ports per destination pod (#7796)
 * [6c2c7c319](https://github.com/kubeovn/kube-ovn/commit/6c2c7c3195db8165bd10b0c770e64e41f302ffe2) fix(tproxy): preserve TCP probe failures for custom VPC pods (#7803)
 * [cd906bbc8](https://github.com/kubeovn/kube-ovn/commit/cd906bbc82530ae7991a48aa4a88907068ffde9f) fix(netpol): fail closed on ambiguous named ports (#7795)
 * [1f56c6751](https://github.com/kubeovn/kube-ovn/commit/1f56c67517181410f4c96b3bfdbe92abb33d4bdf) Add internal tenant network as default network for vpc egress gateway pod (#7425)
 * [37e8dfcd5](https://github.com/kubeovn/kube-ovn/commit/37e8dfcd53b3f0044beccf177211bd20f3df0336) fix(daemon): refresh underlay Service flows when OVS ofports change (#7763)
 * [7a56b34c5](https://github.com/kubeovn/kube-ovn/commit/7a56b34c5485ced198e03e15861b29f0c9c0948b) fix(deps): update kubernetes (master) (#7774)
 * [c088305b5](https://github.com/kubeovn/kube-ovn/commit/c088305b5865c57e379a93c2435457b18d35454a) test(controller): benchmark policy event queue coalescing (#7783)
 * [7f40800d9](https://github.com/kubeovn/kube-ovn/commit/7f40800d94f0671b4c67a85cafa036aac5f67d96) fix(controller): coalesce admin policy update deltas (#7782)
 * [f5a7e34f0](https://github.com/kubeovn/kube-ovn/commit/f5a7e34f0899bb497fb5e46f00981b1dd541e595) ci: allow cold-cache Webhook and HA E2E runs to finish (#7781)
 * [0e6bea527](https://github.com/kubeovn/kube-ovn/commit/0e6bea527cafb43ab9e2c2592d260cb1f7d98b48) fix(e2e): support legacy kubectl-ko context selection (#7779)
 * [a66894f2f](https://github.com/kubeovn/kube-ovn/commit/a66894f2f59aab429a5fe99a836196cf44f4e7e2) fix(controller): index NetworkPolicy label events (#7769)
 * [5c8be6588](https://github.com/kubeovn/kube-ovn/commit/5c8be65881fc5902df24b3ad9b1f8f8dc4260efb) fix(deps): update module google.golang.org/grpc to v1.85.0 (#7775)
 * [0cf6b6b78](https://github.com/kubeovn/kube-ovn/commit/0cf6b6b78b071843265ff86b0d3fb404a6aa0c09) fix(deps): update module cel.dev/cel-go to v0.33.1 (#7776)
 * [acc21c2a9](https://github.com/kubeovn/kube-ovn/commit/acc21c2a937652ebf2b7bbc7fe5f0f04c377c728) fix(deps): update golang.org/x/net to v0.62.0 (#7778)
 * [09543a006](https://github.com/kubeovn/kube-ovn/commit/09543a00630fa6a021a48d68788dffc9425ffc63) fix(controller): preserve successor subnet routes (#7771)
 * [3042c9491](https://github.com/kubeovn/kube-ovn/commit/3042c9491a63f9c4797465f5b2521b6c4a4f950a) fix(deps): update module github.com/moby/sys/mountinfo to v0.7.3 (#7766)
 * [8a587ae20](https://github.com/kubeovn/kube-ovn/commit/8a587ae206fb49683944ec01ea876dbd6e5ae17b) chore(deps): update helm/chart-releaser-action action to v1.8.0 (#7767)
 * [26c15a534](https://github.com/kubeovn/kube-ovn/commit/26c15a534ea39b39bc099ae5bebd75cc99db94d6) Fix shared kubeconfig race in OVN-IC E2E tests (#7772)
 * [baa1b2af6](https://github.com/kubeovn/kube-ovn/commit/baa1b2af69d2d215ed4059decd5b88a3d246d95b) Fix standalone VMI network loss after failed live migration (#7768)
 * [f74e28426](https://github.com/kubeovn/kube-ovn/commit/f74e284260b64c501f6fc439fc54b38cd7f459c2) Add libcap license notices to runtime image (#7758)
 * [200badc58](https://github.com/kubeovn/kube-ovn/commit/200badc58803c91ff54b41f2de486f9184a526bd) build: remove Pebble from base images (#7757)
 * [7a1f42c89](https://github.com/kubeovn/kube-ovn/commit/7a1f42c89444efd56ae158cbda67b00c8502988e) test(controller): benchmark NetworkPolicy event fan-out (#7760)
 * [f63f5939f](https://github.com/kubeovn/kube-ovn/commit/f63f5939fa4ba94cabaf2eca2da6cae45fa45c0f) fix(controller): invalidate selector cache on network changes (#7762)
 * [15644416f](https://github.com/kubeovn/kube-ovn/commit/15644416fc3f861b08ada134b8340cecdbe584af) fix(controller): cache repeated NetworkPolicy selector scans (#7756)
 * [23e26d471](https://github.com/kubeovn/kube-ovn/commit/23e26d471a3b126bcfabaaadfa64a7fa94b99ffb) fix(deps): update module github.com/google/cel-go to v0.32.0 (master) (#7754)
 * [a853df853](https://github.com/kubeovn/kube-ovn/commit/a853df853663b483374c4373d6cbd4f35e5b1d48) chore(deps): update dependency talos to v1.14.5 (#7753)
 * [e8a9ef7ba](https://github.com/kubeovn/kube-ovn/commit/e8a9ef7ba5af86275929e7f0ce9b8540b7d83cd7) fix(controller): answer ARP for switch_lb_vip via router port arp_proxy (#7626)
 * [22894164b](https://github.com/kubeovn/kube-ovn/commit/22894164b293a77f412d0537bf3b55ca04d72304) refactor: isolate node IPsec and scheduling privileges from CNI (#7598)
 * [22faabebf](https://github.com/kubeovn/kube-ovn/commit/22faabebf2ff0c7ffd91b1243cb46ef6a7d84b4e) ci: capture network E2E failures without Valgrind (#7723)
 * [db95b1d7e](https://github.com/kubeovn/kube-ovn/commit/db95b1d7ecc53fa8ae0cb787b9575e66cab1a718) ci: stop checking removed scheduled E2E workflow (#7749)
 * [e6e170cf4](https://github.com/kubeovn/kube-ovn/commit/e6e170cf4074e6af338715afddcf1415d2be618a) ci: remove release-1.13 workflows (#7748)
 * [ae8f327d0](https://github.com/kubeovn/kube-ovn/commit/ae8f327d0c5a4ada174394448cc1ed6c639912da) fix(deps): update module github.com/cert-manager/cert-manager to v1.21.3 (#7747)
 * [6dfd86d3a](https://github.com/kubeovn/kube-ovn/commit/6dfd86d3a72908820b11e9509faf6709f809a07a) chore(deps): update dependency cert-manager to v1.21.3 (#7746)
 * [11819ca17](https://github.com/kubeovn/kube-ovn/commit/11819ca17f6351757970189a8b08745e4edadf97) chore(deps): update dependency cilium to v1.20.4 (#7745)
 * [1bcefc510](https://github.com/kubeovn/kube-ovn/commit/1bcefc51092162229e0ca59acf2b49c31c6ace72) Include licenses for bundled image binaries (#7743)
 * [5da21f4cc](https://github.com/kubeovn/kube-ovn/commit/5da21f4ccf03a3bafb392d3d9b192ba2385963f6) feat(helm): allow externally managed webhook certificates (#7727)
 * [f2977afa5](https://github.com/kubeovn/kube-ovn/commit/f2977afa5f9046ef2e15e82f318f345e73cbf543) fix(deps): update module github.com/ovn-kubernetes/libovsdb to v0.9.0 (master) (#7737)
 * [4e580bd57](https://github.com/kubeovn/kube-ovn/commit/4e580bd57fe1546a36e8055ee99e461bc066a8a0) revert(ci): restore automatic x86 E2E triggers (#7742)
 * [260486bdb](https://github.com/kubeovn/kube-ovn/commit/260486bdb51c367ec85a4d1b7fca7c4571ac8017) chore(deps): update actions/setup-python action to v7.1.0 (#7740)
 * [6e5da01d7](https://github.com/kubeovn/kube-ovn/commit/6e5da01d750e7ba0269dbc8c0af75e669d1d9b05) docs: updated CHANGELOG.md (#7738)
 * [88a4303a2](https://github.com/kubeovn/kube-ovn/commit/88a4303a2b766eca7832fff24e07e768da5b5c52) feat(vip): drop pod static ip via vip and the kube_host_vm_vip type (#7721)
 * [53baa7574](https://github.com/kubeovn/kube-ovn/commit/53baa7574695b86501eb8212b4fbfb2837cb2e30) ovn: use standalone databases for single-instance central (#7725)
 * [444f2d43b](https://github.com/kubeovn/kube-ovn/commit/444f2d43b04bb205a872798d77ac0bdaf224690d) build: use Trivy default database repositories (#7729)
 * [cbb911ab2](https://github.com/kubeovn/kube-ovn/commit/cbb911ab2b12ce766a8bbed5b1a70842e34aab84) fix(deps): update module go.podman.io/image/v5 to v5.42.0 (master) (#7728)
 * [ee3634176](https://github.com/kubeovn/kube-ovn/commit/ee3634176f66a047a7ba6800cea70009ae12a3d3) fix(controller): clear FQDN policy resources when DNS resolver is disabled (#7715)
 * [df8cf991d](https://github.com/kubeovn/kube-ovn/commit/df8cf991d550906922739285fa099269812809e3) fix(controller): garbage collect LB resources when disabled (#7719)
 * [0e4d0fbf2](https://github.com/kubeovn/kube-ovn/commit/0e4d0fbf2000d37f9c79e151c792e27c2a9b9790) remove node external gateway reconciliation (#7670)
 * [626372ed0](https://github.com/kubeovn/kube-ovn/commit/626372ed08f272ec2a8f9711105d51ea2b6b8bf7) chore(deps): update dependency go to v1.27.3 (#7722)
 * [ec6a49e6e](https://github.com/kubeovn/kube-ovn/commit/ec6a49e6e74c89aefda563e6d9d67b03ab2a0fee) feat(cni): make periodic recompute optional (#7714)
 * [b7082f15b](https://github.com/kubeovn/kube-ovn/commit/b7082f15b130bfbd86ee05af29aa83b3840483e9) ci: use Trivy default database repositories (#7716)
 * [4c072cced](https://github.com/kubeovn/kube-ovn/commit/4c072cced1da94ed4003bf7e31b055e16bf88d66) fix(controller): only migrate the tier of kube-ovn ACLs (#7710)
 * [5582c44ed](https://github.com/kubeovn/kube-ovn/commit/5582c44eda81d802934788656be2943764c06109) fix: reconcile VLAN conflicts promptly and unblock subnet deletion (#7539)
 * [8d294e327](https://github.com/kubeovn/kube-ovn/commit/8d294e327628a57570f8d7f0933088a38ef839e1) docs: updated CHANGELOG.md (#7685)
 * [2095ad136](https://github.com/kubeovn/kube-ovn/commit/2095ad13656ee4c0d26e85ba33526ffd576dcd19) chore(deps): update dependency golangci/golangci-lint to v2.15.0 (#7708)
 * [51e743966](https://github.com/kubeovn/kube-ovn/commit/51e7439662f0d3b2b3d941a29a7a3e7553ee0cea) chore(deps): update dependency golangci-lint to v2.15.0 (master) (#7707)
 * [f0cb09f86](https://github.com/kubeovn/kube-ovn/commit/f0cb09f864d5174b66645ec49516bf18158a0a5f) chore(deps): update dependency talos to v1.14.4 (#7706)
 * [7648d2d7a](https://github.com/kubeovn/kube-ovn/commit/7648d2d7a5ea2f6f08954527a6fd8089921bfdee) feat(cni): optionally skip bundled CNI plugins (#7694)
 * [dc3b6f062](https://github.com/kubeovn/kube-ovn/commit/dc3b6f062d9940c594c16855a506350b66832569) fix: satisfy modernize for IPPool test literals (#7699)
 * [14ffa4cf0](https://github.com/kubeovn/kube-ovn/commit/14ffa4cf001d5f04dec89ea8ef4baa5b758de8b0) chore(deps): update dependency cilium to v1.20.3 (#7687)
 * [f904bc3dc](https://github.com/kubeovn/kube-ovn/commit/f904bc3dc9892096f9e40f752b1fb24831663e92) chore(deps): update dependency helm to v4.3.1 (#7688)
 * [67afdf051](https://github.com/kubeovn/kube-ovn/commit/67afdf051f166bb6d94d24f864869b160a507367) fix: run modernize from Go tool dependencies (#7693)
 * [de5519864](https://github.com/kubeovn/kube-ovn/commit/de5519864eec6595311381effe9e09322228a9f5) fix: route terminating ippools to the delete queue on informer list (#7691)
 * [91f5782e5](https://github.com/kubeovn/kube-ovn/commit/91f5782e53fb23b28c11ff3c97f56aa624c0fe5b) fix: preserve Service LB affinity and stale-connection handling (#7678)
 * [1f6875b5e](https://github.com/kubeovn/kube-ovn/commit/1f6875b5e0583e5e2d3d094e2bc5d51e9e2296bb) fix: keep routed Service traffic on its configured proxy path (#7679)
 * [524119bbc](https://github.com/kubeovn/kube-ovn/commit/524119bbccd5879fc96ccd1f567414c075dddc09) fix(deps): update module kubevirt.io/client-go to v1.9.1 (#7676)
 * [6108d4bdd](https://github.com/kubeovn/kube-ovn/commit/6108d4bdd95065f381dd5dd1a8d04679f41ecc84) Update module kubevirt.io/api to v1.9.1 (#7675)
 * [383052717](https://github.com/kubeovn/kube-ovn/commit/383052717a53b69a71c6daaeac7be1f0c535fe87) Update dependency kubevirt to v1.9.1 (#7674)
 * [8e58faedd](https://github.com/kubeovn/kube-ovn/commit/8e58faedd401315837bee43b558120a61037a61e) fix(deps): update module github.com/containerd/containerd/v2 to v2.4.2 (#7677)
 * [8ecdacfce](https://github.com/kubeovn/kube-ovn/commit/8ecdacfce6f5f90ee3469f80ec10cc78c7fe047c) fix: preserve VM and StatefulSet subnets when reusing IPs (#7672)
 * [775aafb3b](https://github.com/kubeovn/kube-ovn/commit/775aafb3bf1b9e0e950c491c1c43de9e4e5ec015) fix: reconcile Service load balancers in one NB transaction (#7660)
 * [2863e27cd](https://github.com/kubeovn/kube-ovn/commit/2863e27cd5e3ac16b02bd049f91327fdf8c81adc) docs: updated CHANGELOG.md (#7610)
 * [7c69a9891](https://github.com/kubeovn/kube-ovn/commit/7c69a98910cab44e9cc811a81531da0ee3d9ef3e) feat(ko): replace production image packet diagnostics with native agent (#7636)
 * [cf1de71d4](https://github.com/kubeovn/kube-ovn/commit/cf1de71d4aa5ea93b36744dd7facd721e1c1d54a) feat(ko): support kubectl-style workload exec (#7639)
 * [bd0e5d602](https://github.com/kubeovn/kube-ovn/commit/bd0e5d60255ecc0cfd9d484fdd2258a89c40a35b) fix(ovs): batch VPC and subnet resource operations atomically (#7661)
 * [bfa4468d3](https://github.com/kubeovn/kube-ovn/commit/bfa4468d3bdb20b386a4c0d5c068e10e1661304b) fix: preserve SecurityGroup rule hashes during forced sync (#7659)
 * [c18a0f182](https://github.com/kubeovn/kube-ovn/commit/c18a0f182daf119fb17aa6e0c969c22a24c2923f) ci: remove duplicate Router LB Rule E2E jobs (#7662)
 * [8745b7ca7](https://github.com/kubeovn/kube-ovn/commit/8745b7ca783e1ef6204f27fd3a38aefc578cc07c) fix: delete chassis_private together with chassis (#7664)
 * [fb9899d35](https://github.com/kubeovn/kube-ovn/commit/fb9899d35ca260e0bd2f3c88276f0e7dac571235) fix: preserve Service ClientIP affinity across IP families (#7666)
 * [cc0d748c3](https://github.com/kubeovn/kube-ovn/commit/cc0d748c3777ce0ef110ba4bdce125d87e242c6e) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.18.0 (#7658)
 * [e58a5093e](https://github.com/kubeovn/kube-ovn/commit/e58a5093e339afe45de0a7192a1d8483b3a492c9) fix: garbage collect orphaned U2O overlay address sets (#7652)
 * [67d51d80e](https://github.com/kubeovn/kube-ovn/commit/67d51d80e70ce88984eb407ebaaeadf889c1491c) fix: refresh cluster and admin policies on label changes (#7649)
 * [74d1d4efc](https://github.com/kubeovn/kube-ovn/commit/74d1d4efc0b0b0afbdb07adabadd91ebad2f7f07) fix: refresh network policies on pod and namespace label changes (#7648)
 * [03cd4b203](https://github.com/kubeovn/kube-ovn/commit/03cd4b2031ba6fbf7c349b0e89064eb06e49bdb7) fix: avoid libovsdb leader watcher for single endpoint (#7651)
 * [5e7843a62](https://github.com/kubeovn/kube-ovn/commit/5e7843a629fb4515af315eb625572f37fdc2aaf2) feat: add optional priorityClassName to VpcEgressGateway (#7635)
 * [3001bb1bc](https://github.com/kubeovn/kube-ovn/commit/3001bb1bc45c479c77185599723df886a7a8d833) fix(deps): use canonical ini module path for v1.67.3 (#7647)
 * [66c4da117](https://github.com/kubeovn/kube-ovn/commit/66c4da117cbd05f9ace31949a909d44c076f5f58) chore(deps): update module github.com/go-ini/ini to v1.67.3 (#7483)
 * [55b704080](https://github.com/kubeovn/kube-ovn/commit/55b70408066f82c70409c7ab79eac7512cf1d4d6) chore(deps): update actions/setup-go action to v7 (#7641)
 * [43de043a6](https://github.com/kubeovn/kube-ovn/commit/43de043a662290c3d94496de0f75aa2d5116bbd3) fix(deps): update module golang.org/x/tools to v0.52.0 (#7646)
 * [d120793ee](https://github.com/kubeovn/kube-ovn/commit/d120793ee8e0ec312d36d08853a71a2cc8ae7c9d) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.33.1 (#7637)
 * [d6bfb6907](https://github.com/kubeovn/kube-ovn/commit/d6bfb6907aa1156b3eaf87439143b17fd6548526) fix(deps): update module golang.org/x/mod to v0.42.0 (#7638)
 * [4f0b30a37](https://github.com/kubeovn/kube-ovn/commit/4f0b30a37e95ac700e5eb8ac77dfb5b202106a20) chore(deps): update actions/checkout action to v7 (#7640)
 * [be5f92b12](https://github.com/kubeovn/kube-ovn/commit/be5f92b12f7b4d0868d2fa28585cf2c5800f97d2) chore(deps): update github artifact actions (#7642)
 * [cef0acc68](https://github.com/kubeovn/kube-ovn/commit/cef0acc685a7368fde8e06ea050d26e293235be3) chore(deps): update golangci/golangci-lint-action action to v9 (#7643)
 * [2139b3545](https://github.com/kubeovn/kube-ovn/commit/2139b35454488d491918f3dd196b45ce5ed37642) feat: load BGP speaker parameters from yaml config file (#6933)
 * [8abffdcaa](https://github.com/kubeovn/kube-ovn/commit/8abffdcaa5b00b0dd75e9099dd878273b94025eb) feat(controller): add failure events for IP resources (#7200)
 * [c69e66c57](https://github.com/kubeovn/kube-ovn/commit/c69e66c57315775bb5924e9dcefa427fe04d26f0) fix: unblock subnet deletion with health-check VIPs (#7596)
 * [e936ea8c1](https://github.com/kubeovn/kube-ovn/commit/e936ea8c11b1877be19720725c6ffac776a195e8) refactor: use builtin new for pointer values (#7633)
 * [284264409](https://github.com/kubeovn/kube-ovn/commit/28426440974569e0b77dbbd8c4e8199ba5d6cdf5) fix(ci): report GitHub API errors from the PR router (#7618)
 * [4719c1a79](https://github.com/kubeovn/kube-ovn/commit/4719c1a79a43a369bbb244c9e5b745b59109bf1f) fix(deps): update golang (#7631)
 * [1908c5902](https://github.com/kubeovn/kube-ovn/commit/1908c590248464bbaea3c30524873e64d067b726) chore(deps): update helm/chart-testing-action action to v2.9.0 (#7630)
 * [c65e0cb61](https://github.com/kubeovn/kube-ovn/commit/c65e0cb61d1f733a595485c48bf52d728ddc8697) fix(ko): support VM and VMI packet capture across migrations (#7621)
 * [32acc1dee](https://github.com/kubeovn/kube-ovn/commit/32acc1dee07d86fee6d4ff0baf1d04efe7922592) fix(controller): switch_lb_vip uses own mac, not gateway mac (#7407)
 * [87217f336](https://github.com/kubeovn/kube-ovn/commit/87217f336efc343c5c620180b69a74a8f4c04a82) fix: wait for hosted OVN central StatefulSet rollout (#7619)
 * [90f32fa41](https://github.com/kubeovn/kube-ovn/commit/90f32fa41443da3fc718c183cf11edffeecd4c0b) fix: resolve HCP northd address families with bounded DNS lookups (#7342)
 * [c383f614f](https://github.com/kubeovn/kube-ovn/commit/c383f614feb07fb69110bea5725cacb571dd9921) feat(ko): rewrite kubectl-ko in Go with an independent node agent (#7606)
 * [afdb8f5f8](https://github.com/kubeovn/kube-ovn/commit/afdb8f5f89fd72d2da2b03e691ae55983d27e42f) fix(ci): execute x86 PR workflows at candidate HEAD (#7617)
 * [6d4d1fca7](https://github.com/kubeovn/kube-ovn/commit/6d4d1fca759090ffe4a5997b073192025e3f74ea) fix(deps): update golang (#7612)
 * [687119ed5](https://github.com/kubeovn/kube-ovn/commit/687119ed5837540b70d004f2595611a4e9b192fa) fix(deps): update module github.com/emicklei/go-restful/v3 to v3.14.0 (#7615)
 * [ac7e496a4](https://github.com/kubeovn/kube-ovn/commit/ac7e496a4dcbeef76396636038f66cf0972f4266) fix(deps): update module github.com/prometheus/client_golang to v1.25.0 (#7613)
 * [4f90bbbb0](https://github.com/kubeovn/kube-ovn/commit/4f90bbbb05e79879c0518b0918f070a15e5aa29b) fix(deps): update module github.com/moby/moby/client to v0.6.2 (#7614)
 * [47c262c02](https://github.com/kubeovn/kube-ovn/commit/47c262c0255222a4beeb2f99a7764f0a6acc36d7) fix(ovs): use atomic Wait+Insert to prevent duplicate load balancers (#7611)
 * [b8d7498e6](https://github.com/kubeovn/kube-ovn/commit/b8d7498e60fa39b059b349182bcdf9bab28f3328) ci: pin GitHub-hosted runners to Ubuntu 24.04 (#7534)
 * [15f466a74](https://github.com/kubeovn/kube-ovn/commit/15f466a74ab38316b27e16f6d1880d87d9d2dcf4) fix(pinger,monitor): allow monitoring/metrics components to find correct ctrl sockets (#7418)
 * [2c21bcca0](https://github.com/kubeovn/kube-ovn/commit/2c21bcca08577deb66d76f8a9537667a3f3f4883) docs: fix single-replica failover exec command (#7587)
 * [33e456008](https://github.com/kubeovn/kube-ovn/commit/33e45600821b0a87032c5373d96f4d2afff8d90d) chore(deps): update helm/kind-action action to v1.15.1 (#7607)
 * [0c2b90fd7](https://github.com/kubeovn/kube-ovn/commit/0c2b90fd74bb1331eacedef83854b0eddec7636a) fix: restore default NORMAL flow on provider bridges after ovs restart (#7546)
 * [0c4f9785f](https://github.com/kubeovn/kube-ovn/commit/0c4f9785f554c6c25ce975d9624cbf283f044edd) refactor: isolate CNI execution and runtime privileges (#7593)
 * [328c7b642](https://github.com/kubeovn/kube-ovn/commit/328c7b642aef664ef3e36e31ccb3e60991ab0247) fix(controller): requeue network policies when pod is routed (#7603)
 * [4b9b7866d](https://github.com/kubeovn/kube-ovn/commit/4b9b7866db036c155be524ca4f622e4d84ed3ba2) refactor: build ACL sampling as a standalone binary (#7602)
 * [2bb75ce50](https://github.com/kubeovn/kube-ovn/commit/2bb75ce500fb6ad5457f02f507cc04c7f4fb6919) chore(deps): update dependency gobgp to v4.10.0 (#7599)
 * [a26a6d17a](https://github.com/kubeovn/kube-ovn/commit/a26a6d17a62cb1d80b796da406c9c360e4164801) fix(deps): update module github.com/osrg/gobgp/v4 to v4.10.0 (#7600)
 * [fac394e1a](https://github.com/kubeovn/kube-ovn/commit/fac394e1a35f78d5ae1a4a457010f0a82cf872f5) fix: preserve OVN EIP status during allocation race (#7595)
 * [0ca8d9623](https://github.com/kubeovn/kube-ovn/commit/0ca8d9623180a7fc00ab550db8a89d0b81551b67) refactor: isolate TProxy namespace dialing in an optional CNI container (#7594)
 * [16bc4a937](https://github.com/kubeovn/kube-ovn/commit/16bc4a93781c85ff1f2f3e4bba7acd00b4dfd781) feat: upgrade ClusterNetworkPolicy to network-policy-api v0.2.0 (#7585)
 * [b4b0efd46](https://github.com/kubeovn/kube-ovn/commit/b4b0efd4699312edb397b3de9bf7386a97704d97) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.4.1 (#7588)
 * [50565c241](https://github.com/kubeovn/kube-ovn/commit/50565c24124e58226e4e8e81c00a09375df086b6) fix(deps): update module golang.org/x/tools to v0.51.0 (#7589)
 * [5bda1bb0c](https://github.com/kubeovn/kube-ovn/commit/5bda1bb0c5241022797da8d09917b4e077ee1a74) test: persist FRR daemon configuration files (#7584)
 * [b3a75e876](https://github.com/kubeovn/kube-ovn/commit/b3a75e8766f84ea040dbcac90ce531430e1200e5) fix(deps): update module sigs.k8s.io/controller-runtime to v0.25.2 (#7583)
 * [6471e875c](https://github.com/kubeovn/kube-ovn/commit/6471e875cde191115fb76f5f4e7cc008de8eabfd) fix(deps): update module github.com/moby/moby/client to v0.6.1 (#7582)
 * [4a234827f](https://github.com/kubeovn/kube-ovn/commit/4a234827fc7a3f1f6a9f6ae753911b5aadf8b610) fix(deps): update module github.com/moby/moby/api to v1.56.1 (#7581)
 * [c2fcd479e](https://github.com/kubeovn/kube-ovn/commit/c2fcd479eb8b7821274db65b69c0e2bb6b43e665) docs: updated CHANGELOG.md (#7579)
 * [9b254720a](https://github.com/kubeovn/kube-ovn/commit/9b254720addc0c2def82e134251d116034585b0a) fix(deps): update module github.com/cenkalti/backoff/v7 to v7.0.1 (#7580)
 * [a1c358f6f](https://github.com/kubeovn/kube-ovn/commit/a1c358f6f22704c9cb3b4263d4a441bc976cfbab) fix(controller): reconcile template VIPs and chassis variables atomically (#7575)
 * [bd71776c4](https://github.com/kubeovn/kube-ovn/commit/bd71776c443783475f3dda2c447963d303b89ebb) fix(ovs): remove stale load balancer health checks when disabled (#7573)
 * [dfc933576](https://github.com/kubeovn/kube-ovn/commit/dfc933576116e7bf3352937eeb1dcf30faa5face) fix(controller): clean family-scoped external service LBs (#7574)
 * [31bc9b872](https://github.com/kubeovn/kube-ovn/commit/31bc9b872ffeb28cca0c4c2292c390fc876a0b3d) fix(ovs): preserve per-port DHCP options during subnet reconciliation (#7577)
 * [6429a05cc](https://github.com/kubeovn/kube-ovn/commit/6429a05ccfd4caf3338bac06f7fd891caffa8cb8) fix(ovs): use atomic Wait+Insert to prevent duplicate logical switches and routers (#7542)
 * [d333928bc](https://github.com/kubeovn/kube-ovn/commit/d333928bc4663decd79d044bf686715d450b54d7) test: wait for informer service deletion (#7571)
 * [425005838](https://github.com/kubeovn/kube-ovn/commit/42500583812e99d1456f73422cd343debfd153af) fix(ovn): preserve router-owned ARP and ND with flood suppression (#7549)
 * [2fda8d268](https://github.com/kubeovn/kube-ovn/commit/2fda8d26872b877cd376b4af993a759f6ecd54a8) test(e2e): probe connectivity throughout live migration (#7513)
 * [203824b45](https://github.com/kubeovn/kube-ovn/commit/203824b4522ff3d4f9e51cb6c3ff327e698cf105) fix: preserve security groups when routing pods (#7565)
 * [eeb5b9f9c](https://github.com/kubeovn/kube-ovn/commit/eeb5b9f9c216d746e194e2e71df327960ac16ce5) chore(deps): update dependency submariner to v0.24.2 (#7561)
 * [a3e9ad2f1](https://github.com/kubeovn/kube-ovn/commit/a3e9ad2f1479ab2d606cda62968ca8a06ecf10f2) chore(deps): update dependency talos to v1.14.2 (#7562)
 * [fb9136130](https://github.com/kubeovn/kube-ovn/commit/fb9136130d5fc89994e95e0b2f36caa720b9e9ea) fix: support switch_lb_vip attachments with scoped load balancers (#7556)
 * [7e5c14928](https://github.com/kubeovn/kube-ovn/commit/7e5c14928b8dd70a6c355b9b69e17e2256accd77) docs: updated CHANGELOG.md (#7557)
 * [1cb3d20a4](https://github.com/kubeovn/kube-ovn/commit/1cb3d20a48316e763de4ac99b59d5fcb71386397) docs: updated CHANGELOG.md (#7521)
 * [ae1a05d71](https://github.com/kubeovn/kube-ovn/commit/ae1a05d71ceb04f3c0fbffa024f464f95989ed54) fix(daemon): do not require ovn-external-gw-config on external-gw nodes (#7452)
 * [1a1289378](https://github.com/kubeovn/kube-ovn/commit/1a128937895fe995863267bfdaecd9e1f48b7414) fix: reconcile attachSubnets for switch_lb_vip VIP (#6969)
 * [1f506959b](https://github.com/kubeovn/kube-ovn/commit/1f506959b969ac2d9040a332d05a37409e3f7f9e) fix: make EVPN VXLAN SVI forwarding reliable (#7514)
 * [ab8ead9dc](https://github.com/kubeovn/kube-ovn/commit/ab8ead9dc4570aa5685c9803ab5da39ae9d060fd) test: wait for FRR BGP readiness in E2E topology (#7526)
 * [da36be358](https://github.com/kubeovn/kube-ovn/commit/da36be3588cec7e4ea41208592303dfdb70e882c) fix(controller): Protect VMs from incorrect chassis settings (#7420)
 * [45fb2a0bd](https://github.com/kubeovn/kube-ovn/commit/45fb2a0bd01ca225f006a48847bcff4c506018ad) fix(e2e): preserve OVN string quotes in VIP queries (#7538)
 * [50c1d84f5](https://github.com/kubeovn/kube-ovn/commit/50c1d84f55d5414bda886d7848ed76c919938ec5) fix: fall back to docker RepoDigests for Kind KubeVirt aliases (#7537)
 * [f9ed176fa](https://github.com/kubeovn/kube-ovn/commit/f9ed176fa6d3052dc1a2d4990e2dc3bdea2eb4b8) fix: avoid concurrent VIP virtual port creation (#7527)
 * [16b24d03f](https://github.com/kubeovn/kube-ovn/commit/16b24d03f78ae0aaec62b2c31a6d7ec057d8c829) fix(controller): configure VMI migration ports while pending (#7415)
 * [fc4d57e05](https://github.com/kubeovn/kube-ovn/commit/fc4d57e05b13e84f2859c07d1fcfef1f2ad3e491) Update kubernetes (#7528)
 * [cfb0a18b3](https://github.com/kubeovn/kube-ovn/commit/cfb0a18b3685723a0cd25c26186c617daf8fff03) Update module github.com/containerd/containerd/v2 to v2.4.1 (#7529)
 * [96a2999a4](https://github.com/kubeovn/kube-ovn/commit/96a2999a4cc40a57a5ab2937404ecad71f974ae2) chore(deps): bump github.com/containerd/containerd/v2 (#7530)
 * [c986cf1d7](https://github.com/kubeovn/kube-ovn/commit/c986cf1d7d0bbf7fe3d0d4264eab5c6f5d36e207) Update dependency golangci-lint to v2.14.0 (#7531)
 * [f04893181](https://github.com/kubeovn/kube-ovn/commit/f04893181d6a9f891e058f8fcf87d492575a7dd1) Update module github.com/onsi/gomega to v1.44.0 (#7532)
 * [36f4cd019](https://github.com/kubeovn/kube-ovn/commit/36f4cd019951108da3cd21c77c0c0b7a45bccfc2) fix: use pure L3 hash for policy ECMP (#7505)
 * [700b6c9ae](https://github.com/kubeovn/kube-ovn/commit/700b6c9ae869f757be5665f17cc92fcb6f0e6fc1) ci: stop rerunning x86 E2E on master pushes (#7524)
 * [0f2eeed00](https://github.com/kubeovn/kube-ovn/commit/0f2eeed00e23807dd73b3c11d82b26cc159e2bfe) fix: avoid flooding IPv6 neighbor solicitations in underlay (#7515)
 * [6a5e52ffa](https://github.com/kubeovn/kube-ovn/commit/6a5e52ffa53a11bb2523a065fee9bedbe1bf86bc) refactor: rename EnableLb/EnableLbSvc to EnableOvnLB/EnablePodLbSvc (#7462)
 * [3986d2c6f](https://github.com/kubeovn/kube-ovn/commit/3986d2c6fab02ec4bcfa8d1512f287ba37415788) fix(controller): delete node router policies by owner on node deletion (#7480)
 * [2e4424c1a](https://github.com/kubeovn/kube-ovn/commit/2e4424c1aabfd25dc0b9bec8e2d92206ea0407be) ci: run conformance jobs on Ubuntu 24.04 (#7508)
 * [46883446f](https://github.com/kubeovn/kube-ovn/commit/46883446fa0c3328f75dcb6f3f53ba04508e177f) fix: keep a single OVN NB monitor for ACL sampling (#7507)
 * [7e27a0fdb](https://github.com/kubeovn/kube-ovn/commit/7e27a0fdbdc0fe0cbb003bb6a89ce4faa1c74a90) ovn: restore bcast_arp_req_flood with default false (#7503)
 * [4bce7830c](https://github.com/kubeovn/kube-ovn/commit/4bce7830c95081d4bd6b0a7fcc50045115645b58) Update dependency ubuntu to v26 (#7498)
 * [e23d73c0d](https://github.com/kubeovn/kube-ovn/commit/e23d73c0d1e7757b0b3ac90df7f6328651706103) docs: updated CHANGELOG.md (#7504)
 * [0008ac438](https://github.com/kubeovn/kube-ovn/commit/0008ac438ef746fafba1d09abac5528b0d54c00c) fix: reconcile distributed-subnet AAP VIP ports (#7491)
 * [71a014429](https://github.com/kubeovn/kube-ovn/commit/71a014429b7086a1eb9d53fa0cda87afe6ff02ee) fix: hash policy route ECMP on L4 ports (#7490)
 * [d819883bc](https://github.com/kubeovn/kube-ovn/commit/d819883bcd8b6c8811d171d65a9d79b0961bb511) fix: preserve IPv6 addresses for EVPN gateway VRF (#7494)
 * [71952054a](https://github.com/kubeovn/kube-ovn/commit/71952054a0cba63b69eada06e97f12e0f9ca2c99) fix: alias KubeVirt container disk digests in Kind (#7493)
 * [6b43b6e7b](https://github.com/kubeovn/kube-ovn/commit/6b43b6e7b2b56505989cdda9cad3b4309e9e5bf1) Add KubeVirt live migration network identity coverage (#7414)
 * [4b498fe1c](https://github.com/kubeovn/kube-ovn/commit/4b498fe1c369d238ffc645a2cc9511b053960e87) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.4.0 (#7492)
 * [53d3fceb3](https://github.com/kubeovn/kube-ovn/commit/53d3fceb3b135b8aa8965f46cee4f1ce85d1a42d) chore(deps): update dependency ubuntu to v26 (#7486)
 * [fbc5eff54](https://github.com/kubeovn/kube-ovn/commit/fbc5eff54a547df01d227e9066c9af9ac726d9a1) docs: updated CHANGELOG.md (#7482)
 * [7de2a1dc7](https://github.com/kubeovn/kube-ovn/commit/7de2a1dc7fa6b18bdcdef80fd6ef130724f852fc) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.17.1 (#7484)
 * [50c6d2a26](https://github.com/kubeovn/kube-ovn/commit/50c6d2a26b052ca4d46ac5499d59b080e6587b32) fix(deps): update module google.golang.org/grpc to v1.84.0 (#7485)
 * [4d05ee886](https://github.com/kubeovn/kube-ovn/commit/4d05ee8861efbe5c28e5122221fdb27f2a7631f9) chore(deps): update jlumbroso/free-disk-space action to v2 (#7487)
 * [c2eefd45c](https://github.com/kubeovn/kube-ovn/commit/c2eefd45c581b4c1a625fce62a6242e94385543c) fix: prevent service LB owner annotation spoofing (#7466)
 * [3ac2cdf28](https://github.com/kubeovn/kube-ovn/commit/3ac2cdf2847a6aa9ee45616c526b98925fea4617) fix(e2e): use legacy MetalLB assertions before v1.17 (#7481)
 * [ac64d0216](https://github.com/kubeovn/kube-ovn/commit/ac64d0216b45dd5b42e96d4e8d509cc1ec56b621) chore(users): add new kube-ovn users (#7478)
 * [37bf0ba24](https://github.com/kubeovn/kube-ovn/commit/37bf0ba2430b713254adbcd4762b2dc7cc7eb376) feat(controller): add maxMtu settings (#7442)
 * [3170803fc](https://github.com/kubeovn/kube-ovn/commit/3170803fc7d4c8fa6ef409a81b452f7dc859ce89) fix: bound initial OVN monitor setup (#7468)
 * [b34814097](https://github.com/kubeovn/kube-ovn/commit/b34814097ba7e13b450af000341842d175a1d568) fix(deps): update module github.com/onsi/gomega to v1.43.1 (master) (#7472)
 * [74301d183](https://github.com/kubeovn/kube-ovn/commit/74301d1833bba129e4b5a27b4c898da9c9abea8c) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.33.0 (#7473)
 * [ff718f194](https://github.com/kubeovn/kube-ovn/commit/ff718f194eec4035927810daff72d5ad5c8c0b6d) chore(deps): update jlumbroso/free-disk-space action to v1.3.2 (#7470)
 * [af87216a2](https://github.com/kubeovn/kube-ovn/commit/af87216a2a70d938e3546d41ed46bf22d967c12f) fix(deps): update module github.com/containerd/containerd/v2 to v2.4.0 (#7471)
 * [47985bf5c](https://github.com/kubeovn/kube-ovn/commit/47985bf5c3ae52414b7be45557b58b3f571c711a) docs: updated CHANGELOG.md (#7459)
 * [5dd485d23](https://github.com/kubeovn/kube-ovn/commit/5dd485d2334bda374885833fe55ad0e6b38561cf) fix(e2e): wait for vpc-external OVS bridges before docker disconnect (#7460)
 * [3f31a0a7d](https://github.com/kubeovn/kube-ovn/commit/3f31a0a7d278e9c1ce7972f84820fbf63bf77185) chore(deps): update dependency cilium to v1.20.2 (#7463)
 * [ed7569544](https://github.com/kubeovn/kube-ovn/commit/ed756954402b3c60db73977a43d408c75bb3fadd) chore(deps): update dependency talos to v1.14.1 (#7464)
 * [a1b576518](https://github.com/kubeovn/kube-ovn/commit/a1b5765184752932daf7213a036bb2429eab5d7e) fix(deps): update module go.podman.io/image/v5 to v5.41.2 (#7465)
 * [05af5081f](https://github.com/kubeovn/kube-ovn/commit/05af5081f9afa840b5c210623588c8cf05c365e2) fix: clear stale ip_port_mappings when reconciling Cluster mode (#7458)
 * [428ff6a38](https://github.com/kubeovn/kube-ovn/commit/428ff6a38589fe674c43d0ca66c2d2f8fe452c8b) fix: honor EndpointSlice node hints for PreferSameNode (#7457)
 * [0d127995a](https://github.com/kubeovn/kube-ovn/commit/0d127995aafc9ce31a9a4836b6fe7cd083023768) feat: add service-scoped OVN load balancers (#7378)
 * [06f57cc1e](https://github.com/kubeovn/kube-ovn/commit/06f57cc1eba5c0df5f7d527f37c59f0dc2a7d2a3) fix(e2e): skip vpc external lifecycle before setup on old versions (#7455)
 * [eea17edc4](https://github.com/kubeovn/kube-ovn/commit/eea17edc433b4578fac636affd9ef825d8954488) fix(deps): update module sigs.k8s.io/controller-runtime to v0.25.1 (#7450)
 * [f45c16c8c](https://github.com/kubeovn/kube-ovn/commit/f45c16c8ca7e7d77a929bf798821cc19d4e85825) support image pull secrets for vpc nat gw pod (#7326)
 * [2dd3adb10](https://github.com/kubeovn/kube-ovn/commit/2dd3adb104db481d3821221a741b0b7424aec57d) fix(vpc-external): keep external LRP chassis in sync (#7447)
 * [c1cc87dcb](https://github.com/kubeovn/kube-ovn/commit/c1cc87dcbf398eacebae2ab0bd808a9ce0aa9fa5) fix(vpc-nat-gw): apply rules to the live gateway Pods (#7448)
 * [c382d7ce0](https://github.com/kubeovn/kube-ovn/commit/c382d7ce045a4366dd97673a9fe96135f2ebf87d) fix(vpc-external): reconcile vpc when gw or pn changes (#6697)
 * [92412cc31](https://github.com/kubeovn/kube-ovn/commit/92412cc31fb42a084bd53686b17349c73b6bac5e) docs: updated CHANGELOG.md (#7444)
 * [d383c6ba6](https://github.com/kubeovn/kube-ovn/commit/d383c6ba64dc98db90ee9eca64d7a27edeb80fde) feat: enable nftable lb service by default (#7403)
 * [48d9a7677](https://github.com/kubeovn/kube-ovn/commit/48d9a767733fa1a3251defc663b895f1c2771abb) ci: share E2E Go cache across refs (#7438)
 * [481ea4233](https://github.com/kubeovn/kube-ovn/commit/481ea4233d1ad9ac1a61018b0475f59afd204279) ci: install Multus and KubeVirt from the PR checkout in conformance e2e (#7435)
 * [b8ea1c808](https://github.com/kubeovn/kube-ovn/commit/b8ea1c808f280e00a5086df74502f58b01196e4c) fix(daemon): flush nftables chains before deleting them (#7433)
 * [86e7790a9](https://github.com/kubeovn/kube-ovn/commit/86e7790a9a14f44bf58ef0c51c66d90be8158777) fix(vpc-nat-gw): track NAT references by UID (#7394)
 * [172396a11](https://github.com/kubeovn/kube-ovn/commit/172396a11a272a87b633b159d0c19d2822883538) ci: collect node logs when kubectl cp fails (#7428)
 * [5f29ecffa](https://github.com/kubeovn/kube-ovn/commit/5f29ecffa0eccca5e453146533f963788e770b98) test(e2e): tolerate asynchronous iptables chain creation (#7429)
 * [138c9499b](https://github.com/kubeovn/kube-ovn/commit/138c9499b545bddbc97ddfe1b39829b91116a102) test(e2e): preserve OVN interconnection spec order (#7431)
 * [31c5eefed](https://github.com/kubeovn/kube-ovn/commit/31c5eefedf6d5e4f7f8b2e30871703d7e78f4ebf) fix(deps): update module github.com/mdlayher/packet to v1.2.0 (#7430)
 * [4131a65b6](https://github.com/kubeovn/kube-ovn/commit/4131a65b65bab6d45b7f81ed283b1c1ff2b2dd7d) fix(dhcp): keep commas inside braced DHCP option values (#7427)
 * [f646c26d5](https://github.com/kubeovn/kube-ovn/commit/f646c26d57ef2e29680e1c62b8b55958503127c5) fix(test): persist FRR BGP topology configuration (#7426)
 * [fdd64b638](https://github.com/kubeovn/kube-ovn/commit/fdd64b638012b2e47cadfa3b48b60ce66229c73f) fix(vpc-nat-gw): make QoS lifecycle converge safely (#7404)
 * [0e8cbd829](https://github.com/kubeovn/kube-ovn/commit/0e8cbd8293c50d33b5f0161f092314c8095da543) chore(deps): update dependency cert-manager to v1.21.2 (#7423)
 * [421482b38](https://github.com/kubeovn/kube-ovn/commit/421482b38818405d4e7372181beeeffec523695e) fix: close orphaned NAT leaks in ovn fip/snat/dnat/eip controllers (#7405)
 * [82f86f9cb](https://github.com/kubeovn/kube-ovn/commit/82f86f9cbd1b15c07c801f87b6398264a2a31983) test(e2e): gate MetalLB VIP datapath assertions to v1.15+ (#7413)
 * [40465641b](https://github.com/kubeovn/kube-ovn/commit/40465641be16efcac4e5cc26a17307d3b1285112) fix(deps): update module kubevirt.io/client-go to v1.9.0 (#7380)
 * [70410fc38](https://github.com/kubeovn/kube-ovn/commit/70410fc381c4bb8027156547ea6e104b8eb00c24) fix(deps): update module github.com/cert-manager/cert-manager to v1.21.2 (#7409)
 * [7f8a99adc](https://github.com/kubeovn/kube-ovn/commit/7f8a99adc67c481a690a087610f19bcf559da837) chore(deps): update dependency helm to v4.3.0 (#7410)
 * [4fe5af48d](https://github.com/kubeovn/kube-ovn/commit/4fe5af48dcb862246fbd636baede65ad22918015) fix: use per-family match for dual-stack EIP/SNAT policies (#7401)
 * [6b0109b2a](https://github.com/kubeovn/kube-ovn/commit/6b0109b2a4e8db4166dcb111d0e7525054c29f35) fix(controller): close nat gateway workload lifecycle gaps (#7402)
 * [e452a5811](https://github.com/kubeovn/kube-ovn/commit/e452a5811c40d16eb589abaf292d315361a47c31) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.17.0 (#7389)
 * [6f556d170](https://github.com/kubeovn/kube-ovn/commit/6f556d1704af9b663208da2114388f8aad079c70) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.32.2 (#7406)
 * [00fd89ec5](https://github.com/kubeovn/kube-ovn/commit/00fd89ec5a334662e657c0c42dc523b89a7b7331) fix(crd): make vlanInterfaces item validation actually take effect (#7108)
 * [e794f6059](https://github.com/kubeovn/kube-ovn/commit/e794f605935afdcc7666a9e8b938561243d9c013) feat: support LB svc backed by vpc nat gw nft share dnat (#7047)
 * [8f3b8a0e9](https://github.com/kubeovn/kube-ovn/commit/8f3b8a0e9677d007c1274025943ae28bf22e8a6e) chore(deps): update dependency multus to v4.3.1 (#7396)
 * [e0996a0e8](https://github.com/kubeovn/kube-ovn/commit/e0996a0e86789d388f1fadeb0d8581a8b8c7d031) fix(deps): update module gopkg.in/k8snetworkplumbingwg/multus-cni.v4 to v4.3.1 (#7397)
 * [e0ea8b3e3](https://github.com/kubeovn/kube-ovn/commit/e0ea8b3e33b9f37d070cee6ebcce4cb80809ed4f) fix(deps): update golang (#7398)
 * [2d66feca6](https://github.com/kubeovn/kube-ovn/commit/2d66feca605fea12f37ebf219dc1fcc8071cd77d) fix(deps): update module sigs.k8s.io/structured-merge-diff/v6 to v7 (#7385)
 * [6731b1b8c](https://github.com/kubeovn/kube-ovn/commit/6731b1b8c73215d3d7065c53f397a240b950ee22) test: retry qos bandwidth measurements (#7390)
 * [758763cda](https://github.com/kubeovn/kube-ovn/commit/758763cda9fd5ac62ca08ced6b957df4ca8f2f48) ci: use kindest/node v1.37.0 and default kube-proxy to nftables (#7391)
 * [2e4d04115](https://github.com/kubeovn/kube-ovn/commit/2e4d0411596efdf5b9e915bbeae940abf5dda387) upgrade OVN to 26.03 and OVS to 3.7 (#7357)
 * [f49ec54d2](https://github.com/kubeovn/kube-ovn/commit/f49ec54d28230bbe8b3b6a407bbafd3eabead1ec) fix(ci): use Image Factory installer for Talos 1.14 (#7388)
 * [c7aea9cdc](https://github.com/kubeovn/kube-ovn/commit/c7aea9cdcb2c07897c785e56b01d9ba3f5577884) feat(controller): support dynamic nat gateway lan IP (#7237)
 * [3e78de00b](https://github.com/kubeovn/kube-ovn/commit/3e78de00bcca2dc91925713c67b8ee1d01e2944b) fix(deps): update module github.com/moby/moby/client to v0.6.0 (#7384)
 * [19263c5ec](https://github.com/kubeovn/kube-ovn/commit/19263c5ec55b613be85c4f675b1443c829807ba7) fix(deps): update module github.com/containernetworking/cni to v1.3.1 (#7379)
 * [e8203d9e0](https://github.com/kubeovn/kube-ovn/commit/e8203d9e0fa32e8105cc5449e56f126797840ec9) fix(deps): update module github.com/containerd/containerd/v2 to v2.3.5 (#7381)
 * [66987b892](https://github.com/kubeovn/kube-ovn/commit/66987b892ca7c5635080a89e818ae3fc6e54e757) chore(deps): update dependency talos to v1.14.0 (#7382)
 * [616b3fa26](https://github.com/kubeovn/kube-ovn/commit/616b3fa269278b710ed676a96ddd8cff257e10c5) fix: require localConnectIP for VPC peering (#7349)
 * [af90ea73c](https://github.com/kubeovn/kube-ovn/commit/af90ea73c1313a53de442127f2b6321a6471ffab) controller: record events for VPC resources (#7202)
 * [6bc15a51f](https://github.com/kubeovn/kube-ovn/commit/6bc15a51f9c67e2dd44a1017fc5c23a7a09c1f9d) fix(deps): update kubernetes (master) (#7335)
 * [82c046a59](https://github.com/kubeovn/kube-ovn/commit/82c046a59720309fe512a5b879841386765cf8a0) test(e2e): wait for live OVN DB cluster recovery (#7376)
 * [e00451782](https://github.com/kubeovn/kube-ovn/commit/e004517829dbcf346e2b880061c260193de2b5e9) controller: record subnet and IPPool lifecycle events (#7201)
 * [2471f1ae0](https://github.com/kubeovn/kube-ovn/commit/2471f1ae0d46bdea26e5ce3b4a557f78618c45a6) Support kube-proxy nftables mode (#7046)
 * [1381a3f36](https://github.com/kubeovn/kube-ovn/commit/1381a3f3623f10596a9bf6080cdffcbf6e75f8ac) test(e2e): wait for all OVN databases after recovery (#7373)
 * [931019b38](https://github.com/kubeovn/kube-ovn/commit/931019b3871a705cb863b78f64d98ca435798ba8) fix(deps): update module kubevirt.io/api to v1.9.0 (master) (#7103)
 * [50d3c5910](https://github.com/kubeovn/kube-ovn/commit/50d3c5910174ce1f1ec85adc391fe82dd8af5eba) chore(deps): update ubuntu docker tag to v26 (master) (#6683)
 * [3ace319ab](https://github.com/kubeovn/kube-ovn/commit/3ace319ab5c807d02c507b66d7eaac0cfbee1435) test(ci): accept trusted dispatch cache condition (#7372)
 * [199392042](https://github.com/kubeovn/kube-ovn/commit/199392042d3d7ce27e11b3211ccb5c7ef629a076) controller: record node reconciliation failures (#7203)
 * [7bd61ea49](https://github.com/kubeovn/kube-ovn/commit/7bd61ea49ff4d329ec91d13e77e9f4405f954a8c) controller: record VLAN resource events (#7199)
 * [05b88a8b4](https://github.com/kubeovn/kube-ovn/commit/05b88a8b4025cb5b4e40b98db0a349dcc649c421) fix(ci): rebuild base images for trusted PR dispatch (#7370)
 * [7aab10c5b](https://github.com/kubeovn/kube-ovn/commit/7aab10c5bdd29e0671e63f879b9c240f1a72759a) fix(ci): save trusted E2E Go cache (#7371)
 * [26b974533](https://github.com/kubeovn/kube-ovn/commit/26b97453311904266f4d57ed749730474e8abd30) chore(deps): update helm/kind-action action to v1.15.0 (#7367)
 * [e6a32a5ce](https://github.com/kubeovn/kube-ovn/commit/e6a32a5ceab2aa08b9aa724443b6d470728ba98f) feat: add routed subnet mode with /32 addressing and gateway hairpin (#7294)
 * [8145a9656](https://github.com/kubeovn/kube-ovn/commit/8145a96566a467cd8168cc09ceca5cb5e05bf153) Remove OVN version compatibility upgrade logic (#7358)
 * [6d5a9ef35](https://github.com/kubeovn/kube-ovn/commit/6d5a9ef352af8ea007d8b5c5df188a5e4731abe2) deps: update libovsdb to latest upstream (#7242)
 * [32f5ac7cd](https://github.com/kubeovn/kube-ovn/commit/32f5ac7cde87a85c17d7da62e77b023ddebb7f89) docs: updated CHANGELOG.md (#7356)
 * [30933dd64](https://github.com/kubeovn/kube-ovn/commit/30933dd6444bf111e1c502df72cb969b746637aa) chore(deps): update dependency go to v1.27.1 (#7355)
 * [f2e2cd505](https://github.com/kubeovn/kube-ovn/commit/f2e2cd5050a06a948c7adddf3ae707f7ffbeb5fc) fix(deps): update module github.com/osrg/gobgp/v4 to v4.9.0 (#7354)
 * [ada0334c3](https://github.com/kubeovn/kube-ovn/commit/ada0334c302eb4c1cbf43ff1a0d6d99c2f97a4a7) chore(deps): update dependency gobgp to v4.9.0 (#7353)
 * [203b25782](https://github.com/kubeovn/kube-ovn/commit/203b257826484551db4490f11a77b206f4a56f36) chore(e2e): bump test image tags (#7352)
 * [d47aeb4df](https://github.com/kubeovn/kube-ovn/commit/d47aeb4df5b449aae94bb7e00f2c429d053772a1) test(e2e): skip QoS finalizer cases before v1.17 (#7351)
 * [5e30b865f](https://github.com/kubeovn/kube-ovn/commit/5e30b865f72a12afde119feaec0ff69999c9d52d) fix(controller): reconcile missing node logical switch ports (#7322)
 * [7d39928dc](https://github.com/kubeovn/kube-ovn/commit/7d39928dc56589cc04158e1d3fde3fd789f80fb3) controller: propagate VIP release errors (#7185)
 * [d6de4646c](https://github.com/kubeovn/kube-ovn/commit/d6de4646c62b0c9357fafd57ea92a70cc1a151d9) chore(deps): update dependency submariner to v0.24.1 (#7345)
 * [894109052](https://github.com/kubeovn/kube-ovn/commit/894109052ed77cb36b8b2b9dff95162f732b8f5f) controller: extract stale pod network state (#7184)
 * [f2e011495](https://github.com/kubeovn/kube-ovn/commit/f2e011495c503cd03b192bb582127b429a855a5b) daemon: consolidate pod QoS reconciliation (#7183)
 * [59e858950](https://github.com/kubeovn/kube-ovn/commit/59e85895024a60876e1589065e45a73a269ea60b) test(bgp): avoid duplicate IPv6 source cleanup (#7344)
 * [8f6c10f49](https://github.com/kubeovn/kube-ovn/commit/8f6c10f49dfe6622b74e647f040024f1c05ff52b) chore(deps): update dependency golangci-lint to v2.13.2 (#7337)
 * [6488c381e](https://github.com/kubeovn/kube-ovn/commit/6488c381ec1b60c604dfa341f65d18ff0dbb989b) fix(deps): update module github.com/onsi/gomega to v1.43.0 (#7338)
 * [dc9064250](https://github.com/kubeovn/kube-ovn/commit/dc90642504f47072a6117e6303e1674b64d73edd) log: warning when using tun nic ip maybe keepalived vip (#7336)
 * [6e14daf95](https://github.com/kubeovn/kube-ovn/commit/6e14daf95a326d070099f91608cd2c0b4cad6772) test(metallb): harden underlay VIP e2e coverage
 * [583c37864](https://github.com/kubeovn/kube-ovn/commit/583c3786435e2df8506e189b70a90dd4f51511c1) fix(daemon): allow a full-mask address to be used as the encap IP (#7330)
 * [a897ebb26](https://github.com/kubeovn/kube-ovn/commit/a897ebb26a860e281d2d7c14cc92a313f8c51802) test(controller): use direct port group literals (#7334)
 * [70790c8ba](https://github.com/kubeovn/kube-ovn/commit/70790c8ba94db650458b3260e1060c7030689901) test(ovs): normalize numeric network policy ACL name (#7332)
 * [806c8f7f9](https://github.com/kubeovn/kube-ovn/commit/806c8f7f9e16b421113b1a9cc2a2ba013cb5af33) fix(ci): include ACL sampling in executor result (#7333)
 * [18f043bbd](https://github.com/kubeovn/kube-ovn/commit/18f043bbdf279b3dca5045f91c692de9f650511c) test(controller): update network policy delete request
 * [e403f70d1](https://github.com/kubeovn/kube-ovn/commit/e403f70d1e39cabca294fe0c5eb1e29006fa33bd) test: add ACL sampling end-to-end coverage (#7174)
 * [0ff5d756e](https://github.com/kubeovn/kube-ovn/commit/0ff5d756e652b30a8add604832cada05df00a466) feat: configure ACL sampling in install script (#7173)
 * [3583e69da](https://github.com/kubeovn/kube-ovn/commit/3583e69dac46cd3ff22a7f79f69d2ceaa6aee61b) docs: describe NetworkPolicy ACL sampling (#7172)
 * [d4936ef54](https://github.com/kubeovn/kube-ovn/commit/d4936ef54830ffd299578a8f48fe7812b3e17c3d) feat: configure ACL sampling in v2 chart (#7171)
 * [01d1beef6](https://github.com/kubeovn/kube-ovn/commit/01d1beef67aedf1a4a2632e655ce265f4913aa51) feat: add ACL sample debug commands (#7170)
 * [de6004cc5](https://github.com/kubeovn/kube-ovn/commit/de6004cc581711fea02778f899b75b0b900fa916) feat: listen for local ACL samples (#7169)
 * [a4ef051e8](https://github.com/kubeovn/kube-ovn/commit/a4ef051e82c632c8048224e0824390198b846643) feat: decode NetworkPolicy ACL samples (#7168)
 * [a0de9ff01](https://github.com/kubeovn/kube-ovn/commit/a0de9ff018482b6e7188b677edd2adc543f73974) feat: deliver ACL samples locally on nodes (#7167)
 * [2748fb1a0](https://github.com/kubeovn/kube-ovn/commit/2748fb1a038d8b800a352f3c5f6b3d641d82f87c) feat: report ACL sampling availability (#7166)
 * [f30c624b7](https://github.com/kubeovn/kube-ovn/commit/f30c624b745a1b6e268a988b0a272d8f6e521cc2) feat: clean up owned ACL sampling state (#7165)
 * [71a3a336d](https://github.com/kubeovn/kube-ovn/commit/71a3a336d1ea88b2ba30c784fdbe4209d712b395) feat: retry NetworkPolicy ACL sampling independently (#7164)
 * [9e8fefe95](https://github.com/kubeovn/kube-ovn/commit/9e8fefe95af3dcb40c08d78dd07241f3b5f76346) feat: attach sampling to NetworkPolicy ACLs (#7163)
 * [6da47f82b](https://github.com/kubeovn/kube-ovn/commit/6da47f82bb59fa50c8cb1262206c5f37d4345d07) feat: reconcile OVN ACL sampling objects (#7148)
 * [ea862c235](https://github.com/kubeovn/kube-ovn/commit/ea862c2355f15e20dbf62fa10807fae02911ce6c) feat: allocate stable ACL sample metadata (#7150)
 * [d76025915](https://github.com/kubeovn/kube-ovn/commit/d760259150bcc38b3ddfa3fc17e0cadee31b8c21) feat: add ACL sampling configuration (#7149)
 * [be2c88862](https://github.com/kubeovn/kube-ovn/commit/be2c88862e03e5098e3d12887f03d722a0ff8ab8) chore(deps): update kindest/node docker tag to v1.37.0 (#7328)
 * [2538cfa94](https://github.com/kubeovn/kube-ovn/commit/2538cfa94706d290552cb4291ec927e834fae25e) chore(deps): update dependency kind to v0.33.0 (#7327)
 * [45e41b483](https://github.com/kubeovn/kube-ovn/commit/45e41b483c665ea0dd32b4de876173b3f2f7260b) fix(deps): update module google.golang.org/grpc to v1.83.2 (#7324)
 * [89d0bb79e](https://github.com/kubeovn/kube-ovn/commit/89d0bb79e5a5ae665af098270114c6db73decabf) fix(deps): update module github.com/sirupsen/logrus to v1.10.2 (#7323)
 * [b1b448ef2](https://github.com/kubeovn/kube-ovn/commit/b1b448ef2f2db8c250acf6673b36f6f1ba5515a5) docs: updated CHANGELOG.md (#7317)
 * [4cd09ce46](https://github.com/kubeovn/kube-ovn/commit/4cd09ce469746e780863f6a7bd9b9fdf0860a125) test(bgp): fix IPv6 and dual-stack speaker E2E (#7316)
 * [951864a8f](https://github.com/kubeovn/kube-ovn/commit/951864a8fafcdf97443ae7f7f6a84871a0fd1b09) chore(deps): update dependency go to v1.27.0 (master) (#7283)
 * [1533681d7](https://github.com/kubeovn/kube-ovn/commit/1533681d7a6b3dbdd7e8870fa70f0dc6aa49a82a) fix(ci): run trusted E2E tests from PR head (#7318)
 * [b5f3c23a0](https://github.com/kubeovn/kube-ovn/commit/b5f3c23a002dbc39876e944b054dd1a3042e9471) fix(ci): use PR Go version in trusted E2E jobs
 * [ae1be3578](https://github.com/kubeovn/kube-ovn/commit/ae1be357852557e032be78079b5842389aea229e) test(bgp): run speaker E2E across IP families (#7314)
 * [5e4185e4b](https://github.com/kubeovn/kube-ovn/commit/5e4185e4b8580b22ec7dd690508a0f13c529ad7c) test(e2e): use version-specific QoS paths (#7311)
 * [5f57cdd67](https://github.com/kubeovn/kube-ovn/commit/5f57cdd67b602395dd6de945bb304666363724b6) fix(speaker): refresh announced next hops during reconcile (#7296)
 * [2399e17f3](https://github.com/kubeovn/kube-ovn/commit/2399e17f3aa4f8c0a603b7c1c5c81c918bf3e373) fix(ci): notify gate after trusted executor completion (#7312)
 * [e8179f7e1](https://github.com/kubeovn/kube-ovn/commit/e8179f7e1a4943be02f8eb3d5412b19fd4054329) fix(controller): tolerate TLS rotation database outage (#7299)
 * [785c2a874](https://github.com/kubeovn/kube-ovn/commit/785c2a874198db204b4407820cf374306bd28f3c) fix(ci): ignore unrelated E2E workflow completions (#7308)
 * [32e336822](https://github.com/kubeovn/kube-ovn/commit/32e336822dae57f947e37079208935b29ac3542e) fix(ci): reconcile completed E2E executors (#7307)
 * [1501d9741](https://github.com/kubeovn/kube-ovn/commit/1501d974189b992613ce80079a9e70ea3417bded) fix(ci): validate executors against the live base (#7306)
 * [f500a1cbb](https://github.com/kubeovn/kube-ovn/commit/f500a1cbbc3bfbbc1822d4030dc8abdcfb065fe8) fix(ci): resolve the live E2E target revision (#7305)
 * [c69814b8a](https://github.com/kubeovn/kube-ovn/commit/c69814b8a88914aeb643fe65ce90f2dad4792cf1) ci: separate approved groups from label provenance (#7304)
 * [059917610](https://github.com/kubeovn/kube-ovn/commit/059917610c5f8be476453aaa3dedc8bb25c44552) ci: refresh durable E2E approvals on base updates (#7303)
 * [f236b7075](https://github.com/kubeovn/kube-ovn/commit/f236b7075a1029354e8178abb403030b42b139db) ci: keep durable E2E approvals live (#7302)
 * [2e933d689](https://github.com/kubeovn/kube-ovn/commit/2e933d689e2dec9ac8f4a606759fc54cb7919c1c) ci: preserve controlled labels in x86 E2E approvals (#7298)
 * [e21a153a1](https://github.com/kubeovn/kube-ovn/commit/e21a153a18283ae01a7066976914c6bf29a89202) test(veg): cover named IPPool allocation (#7295)
 * [92187acc1](https://github.com/kubeovn/kube-ovn/commit/92187acc1cc2bb7d634ec2a78a670fb764cf3386) ci: simplify x86 E2E approval commands (#7297)
 * [96306858a](https://github.com/kubeovn/kube-ovn/commit/96306858a830c73d4923230a9c7bd10289dbba29) fix(controller): recycle NP/ANP/BANP OVN objects when features are disabled (#7269)
 * [4c3c6c62d](https://github.com/kubeovn/kube-ovn/commit/4c3c6c62df8e2da10335f1ffbebf9e70557d1851) feat(chart): make monitor replicas configurable (#7293)
 * [6f87f593d](https://github.com/kubeovn/kube-ovn/commit/6f87f593d87edbab332f9c8054fc1c770886e0a3) ci: skip the x86 E2E gate when the PR is already merged (#7281)
 * [988b62d4e](https://github.com/kubeovn/kube-ovn/commit/988b62d4e5ddaf28c96738333daeb47669598f04) fix: preserve VEG policies without dangling address sets (#7251)
 * [86ee7932d](https://github.com/kubeovn/kube-ovn/commit/86ee7932d3b23230cde8bf86bea86bf110ce7f92) feat(veg): add ippool option to vpc egress gateway (#7250)
 * [e8119e9ca](https://github.com/kubeovn/kube-ovn/commit/e8119e9ca80db843bb80102d1c6882b253fe75a1) fix(deps): update kubernetes (#7282)
 * [cdab5ca7d](https://github.com/kubeovn/kube-ovn/commit/cdab5ca7db5ab45cc9a7aa0a43b0c4ddf52d4436) chore(deps): update dependency golangci-lint to v2.13.1 (master) (#7284)
 * [70043375e](https://github.com/kubeovn/kube-ovn/commit/70043375e13c6eecf3931b92ecd4ce686c5fb4d4) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.16.0 (#7285)
 * [ad10f798d](https://github.com/kubeovn/kube-ovn/commit/ad10f798d4b46e7017ff727b9c88e7e046823f64) fix: keep OVN-IC routes and gateways consistent (#7265)
 * [e6ea05439](https://github.com/kubeovn/kube-ovn/commit/e6ea05439cfd146eca9afa88926856ff5be84c01) fix(ipam): stop random allocation of excludeIps after IPPool reconcile (#7275)
 * [c9df9cbbb](https://github.com/kubeovn/kube-ovn/commit/c9df9cbbb0b0dcec08c23fae276124308b6bbe5b) ci: post PR E2E checks with subprocess.run (#7280)
 * [9a19364a6](https://github.com/kubeovn/kube-ovn/commit/9a19364a6d3e0d974636cb4e341dd993c3119149) ci: stop duplicate automatic x86 E2E executors (#7279)
 * [e0d20f1d9](https://github.com/kubeovn/kube-ovn/commit/e0d20f1d9cb9e4e55e10e9cc5f45ffce72f8813c) ci: publish queued x86 E2E placeholders on PR checks (#7278)
 * [764164979](https://github.com/kubeovn/kube-ovn/commit/7641649792f538cfe496e9149f990501f50f9d39) fix(deps): update module github.com/stretchr/testify to v1.12.1 (#7260)
 * [ffb4e8287](https://github.com/kubeovn/kube-ovn/commit/ffb4e82879eaae39af238a267ab07d1eebdb97f5) ci: stream trusted x86 E2E job status onto PR checks (#7274)
 * [c988340e8](https://github.com/kubeovn/kube-ovn/commit/c988340e8e38d5ebef7fedffc2fe9e96b5482978) ci: publish trusted x86 E2E jobs as pull request checks (#7273)
 * [c4c5e5701](https://github.com/kubeovn/kube-ovn/commit/c4c5e5701c4aecdbb65b9af98310cea94d6f08d9) fix(ci): recognize approved gates after GitHub rewrites details_url (#7272)
 * [40819271e](https://github.com/kubeovn/kube-ovn/commit/40819271ed73f552801b8defa89365fab456d0f3) fix(ci): let the x86 E2E reducer post PR comments (#7270)
 * [b95a0da56](https://github.com/kubeovn/kube-ovn/commit/b95a0da5643e5982a34bb89d2c46db0e628bc067) fix(ci): allow the x86 E2E gate to post PR comments (#7268)
 * [c82351fba](https://github.com/kubeovn/kube-ovn/commit/c82351fbac06501f576e21ce096d7ff30a810e43) fix(ci): accept unbound /test e2e comments and post replies (#7267)
 * [e7cd3618d](https://github.com/kubeovn/kube-ovn/commit/e7cd3618dfd524e6e662dfe40ac4dac2bdfa3b3a) fix(deps): update module google.golang.org/grpc to v1.83.1 (#7261)
 * [823034b7f](https://github.com/kubeovn/kube-ovn/commit/823034b7f575c198f2a861f411396ebf3d43ee3f) fix(deps): update module github.com/sirupsen/logrus to v1.10.1 (#7258)
 * [4316cf624](https://github.com/kubeovn/kube-ovn/commit/4316cf6242bdd1584fd001306886454918161af0) chore(deps): update dependency talos to v1.13.9 (#7257)
 * [42375b0cc](https://github.com/kubeovn/kube-ovn/commit/42375b0cc42f5b957231d0bc80a589df6f5b6ef2) fix(ci): keep PR unit test, lint, and image build (#7255)
 * [78363e00a](https://github.com/kubeovn/kube-ovn/commit/78363e00a81fcdd653fb21f878ba92c7e7e51c4d) fix(ci): accept isolated executor refs in x86 E2E selection (#7254)
 * [f8b3c711f](https://github.com/kubeovn/kube-ovn/commit/f8b3c711f15a41a115c6ddbe846129e7aac0b31d) fix(ci): treat missing isolated executor refs as create (#7253)
 * [ead638494](https://github.com/kubeovn/kube-ovn/commit/ead638494dbc81f037f415af20ba4956202552cd) ci: enable selective x86 E2E execution (#7244)
 * [2ba4e3605](https://github.com/kubeovn/kube-ovn/commit/2ba4e360563cd0ec65f4168f8b5fa6c6ff31554e) Update dependency cilium to v1.20.1 (#7252)
 * [17b4ac6e4](https://github.com/kubeovn/kube-ovn/commit/17b4ac6e4b1179daece9dba4fe7eef29ef014c6d) fix(deps): update module github.com/stretchr/testify to v1.12.0 (#7246)
 * [507af9be2](https://github.com/kubeovn/kube-ovn/commit/507af9be23158026fcd06b28fe02ca3e2933a24d) ci: add comment-gated x86 E2E control (#7231)
 * [36d67b7ea](https://github.com/kubeovn/kube-ovn/commit/36d67b7eaa86d3090f1059966f37f38227a0d5a7) test: add BGP speaker unit and end-to-end coverage (#7241)
 * [fc131b9ab](https://github.com/kubeovn/kube-ovn/commit/fc131b9abccd977d5e2722a6d95ddc8bf2a634dc) add gzip to container image for logrotate compression (#7217)
 * [6eb352a17](https://github.com/kubeovn/kube-ovn/commit/6eb352a178d4bae0701b6f1cf1846b49979605d3) fix: reconcile terminating OVN EIP/FIP/DNAT/SNAT from the add handler on restart (#7220)
 * [bb94586e9](https://github.com/kubeovn/kube-ovn/commit/bb94586e9d78d228dc6fb0b6ea0160fe0483caba) fix(deps): update module github.com/sirupsen/logrus to v1.10.0 (#7234)
 * [765c99183](https://github.com/kubeovn/kube-ovn/commit/765c99183529d424ad5d08d9f58469e323770bc2) chore(deps): update dependency helm to v4.2.4 (#7232)
 * [ec921389c](https://github.com/kubeovn/kube-ovn/commit/ec921389c35793b38a272fd3cc60c81a00cde5a1) fix(deps): update golang (#7233)
 * [3a8b60406](https://github.com/kubeovn/kube-ovn/commit/3a8b60406e441c61195f8847b36d1cc20ed98ce5) chore: ignore built binary dist/images/vpc-egress-gateway-observer (#7229)
 * [cae651c3e](https://github.com/kubeovn/kube-ovn/commit/cae651c3e7c2b9c4b0ec993c52bb6995563c0e6b) e2e: trim the iptables eip qos specs to a natgw and an eip case (#7228)
 * [9d59c4fb4](https://github.com/kubeovn/kube-ovn/commit/9d59c4fb46c69566242d2b69c3f4db50c66e49c6) fix(daemon): cap provider VLAN internal port names (#7198)
 * [36cf19a5f](https://github.com/kubeovn/kube-ovn/commit/36cf19a5f64196f9b812f9031f50e9c066efe5c4) fix(e2e): run KubeVirt tests serially (#7227)
 * [07e4d76c8](https://github.com/kubeovn/kube-ovn/commit/07e4d76c861ef12aa55e92fc3480be8a6843f77c) fix(deps): update module go.podman.io/image/v5 to v5.41.1 (#7224)
 * [9ae065704](https://github.com/kubeovn/kube-ovn/commit/9ae065704e2982033ca1011a553e0ebc810ec92b) fix(deps): update module github.com/containerd/containerd/v2 to v2.3.4 (#7223)
 * [522cadf74](https://github.com/kubeovn/kube-ovn/commit/522cadf74dde983f917e34b1b32d46b7d188e2f1) daemon: record provider network status events (#7197)
 * [304654e4d](https://github.com/kubeovn/kube-ovn/commit/304654e4dff8748e5b4c2803ba3b0a87929a5ac9) docs: updated CHANGELOG.md (#7218)
 * [c912cc691](https://github.com/kubeovn/kube-ovn/commit/c912cc691698a58587d6f9af02b1ae28990b2d9c) fix: reconcile terminating iptables EIP/FIP/DNAT/SNAT from the add handler on restart (#7210)
 * [3c0175dcc](https://github.com/kubeovn/kube-ovn/commit/3c0175dcc80d32f3eb2a8131f16e2faf6b355471) test(ipam): wait for excluded pool reconciliation (#7214)
 * [943bcfa2b](https://github.com/kubeovn/kube-ovn/commit/943bcfa2b54025df99e3386e67fa7dff00101ecd) fix(gateway): add bounded BFD session recovery (#7209)
 * [cc572f788](https://github.com/kubeovn/kube-ovn/commit/cc572f788883c252312ffabea2edba4c9a0e267d) test(e2e): avoid platform-specific observer RSS limit (#7216)
 * [a7dc08397](https://github.com/kubeovn/kube-ovn/commit/a7dc08397772e704304914980d7ab365040c346a) fix: release QoS finalizer only after the referencing EIP/NatGw is truly gone (#7176)
 * [e57f7f02b](https://github.com/kubeovn/kube-ovn/commit/e57f7f02b9feb416195bb011993d51613f5116fd) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.32.1 (#7206)
 * [4b4abb9cb](https://github.com/kubeovn/kube-ovn/commit/4b4abb9cb30de78350a670cbb9c1bc019c3f6a8a) fix(deps): update module golang.org/x/mod to v0.39.0 (#7204)
 * [a9173cfd8](https://github.com/kubeovn/kube-ovn/commit/a9173cfd85a047965408adc15b05622ae896b23a) test(e2e): wait for TLS restoration before HA tests (#7195)
 * [73676e0bf](https://github.com/kubeovn/kube-ovn/commit/73676e0bfb999081571ec7a515ce96efd30fb1d8) test(e2e): use kubelet secret update timeout (#7191)
 * [8b19728b3](https://github.com/kubeovn/kube-ovn/commit/8b19728b328c025a2d800555c9501c779df962d1) Add observability for VPC egress gateways (#7116)
 * [c79de0bd5](https://github.com/kubeovn/kube-ovn/commit/c79de0bd5dae98006c05fe743cb3bb2f1d6eb2b0) ipam: align dual-stack family preservation (#7189)
 * [1aeebd1c1](https://github.com/kubeovn/kube-ovn/commit/1aeebd1c1e806fe2542b4bd2457268a7a04040b2) controller: preserve pod for failure events (#7188)
 * [1877c687d](https://github.com/kubeovn/kube-ovn/commit/1877c687dda0edf962896aefe60929b82b52663b) controller: deduplicate vpc egress gateway delete errors (#7186)
 * [abab3df20](https://github.com/kubeovn/kube-ovn/commit/abab3df20a681ed661738379278f196619c74d1a) ipam: restore counter views when preserving the other family's static address (#7187)
 * [27ab34cde](https://github.com/kubeovn/kube-ovn/commit/27ab34cde98698f5f51ed559c00309394f0d7983) feat: make controller leader election timings configurable (#7177)
 * [549864d2b](https://github.com/kubeovn/kube-ovn/commit/549864d2bae0573bb39953e3b9a6f7a78f9be7f6) docs: updated CHANGELOG.md (#7147)
 * [3c5097c19](https://github.com/kubeovn/kube-ovn/commit/3c5097c19475e68de70e00c47d815683aeedde67) fix(cni): refresh IPv6 neighbor caches (#7139)
 * [e79022da1](https://github.com/kubeovn/kube-ovn/commit/e79022da1a535d8cbea5b31d87bc25c642b80910) fix(ovn-central): tolerate transient duplicate leaders (#7138)
 * [7593540ba](https://github.com/kubeovn/kube-ovn/commit/7593540ba543853cc2391002a94b62c707bc2701) fix(gateway): bound bfdd status operations (#7137)
 * [d7483f920](https://github.com/kubeovn/kube-ovn/commit/d7483f9203686c801f27397e623a96ce33158ef1) fix(gateway): bound bfdd control requests (#7136)
 * [9b06a6b2d](https://github.com/kubeovn/kube-ovn/commit/9b06a6b2dc78dfcfb9208d0b9059d8d2b25bd441) fix(gateway): prevent probe disconnects from restarting bfdd (#7135)
 * [50556fa87](https://github.com/kubeovn/kube-ovn/commit/50556fa872c7a2c8005bb35b8315f5dc93db35ec) ci: pin Kamaji chart download (#7131)
 * [d187992d8](https://github.com/kubeovn/kube-ovn/commit/d187992d8a0b212f1876dad7f8e38eb314946820) fix: use fully qualified CRD names in cleanup.sh (#7118)
 * [46a7b80d6](https://github.com/kubeovn/kube-ovn/commit/46a7b80d66525c24114956d4962f6947441ce08e) Update dependency talos to v1.13.8 (#7129)
 * [b4370605d](https://github.com/kubeovn/kube-ovn/commit/b4370605d9d15db6e003ee403a97ccadf1de5710) chore(codegen): regenerate Kubernetes API code (#7128)
 * [9a8380cb1](https://github.com/kubeovn/kube-ovn/commit/9a8380cb10f57c1b8a00a4b81cc1ee5ffe28bac1) docs: updated CHANGELOG.md (#7125)
 * [9f7fb3583](https://github.com/kubeovn/kube-ovn/commit/9f7fb358365c20f0ab92a75c6ea241df00864cfc) ci: add hosted OVN central e2e coverage (#7091)
 * [e9774269a](https://github.com/kubeovn/kube-ovn/commit/e9774269aacd495e23ce1621df10cc6771c8b45d) feat(vpc-egress-gateway): accept bandwidth quantities (#7115)
 * [34d63fbcb](https://github.com/kubeovn/kube-ovn/commit/34d63fbcb133cd2d0530504392d7040d019ac139) fix(deps): update module github.com/osrg/gobgp/v4 to v4.8.0 (#7120)
 * [bdec52c1e](https://github.com/kubeovn/kube-ovn/commit/bdec52c1e85d1cee4b9a891c5d74fafe7f36a8c0) fix(ovn-central): recover from invalid raft headers (#7112)
 * [47a54b4a3](https://github.com/kubeovn/kube-ovn/commit/47a54b4a3a45c873e5d944a413bfeb8c87b8682f) chore(deps): update dependency gobgp to v4.8.0 (#7119)
 * [ec5394c50](https://github.com/kubeovn/kube-ovn/commit/ec5394c50b32e790557d3328a8e6261af449bcad) test(e2e): extend multus ip family coverage (#7093)
 * [0c123f66e](https://github.com/kubeovn/kube-ovn/commit/0c123f66e3d8fbeeeccf1cd063de487d8b4b74a0) chore(deps): update dependency kubevirt to v1.9.0 (#7101)
 * [55d35af2c](https://github.com/kubeovn/kube-ovn/commit/55d35af2cc95fc6c69b437f64130e17cbe95b4fb) fix: clarify command flag help text (#7094)
 * [76e679a39](https://github.com/kubeovn/kube-ovn/commit/76e679a39179cf64349f955799b72733c55b9d3f) fix(ipam): preserve IP pool allocation state (#7107)
 * [77e1ab9db](https://github.com/kubeovn/kube-ovn/commit/77e1ab9db00ca9a83c4d9b4ff1a3e61075518fcd) fix(deps): update module kubevirt.io/client-go to v1.9.0 (#7104)
 * [0c6988847](https://github.com/kubeovn/kube-ovn/commit/0c6988847befd42a03c3fe2b842fb04d3cc336b7) fix(deps): update module google.golang.org/grpc to v1.83.0 (#7102)
 * [5d69cd4c3](https://github.com/kubeovn/kube-ovn/commit/5d69cd4c3fa964462327c6cf0340d2e1c34827c2) chore(deps): update actions/stale action to v11 (#7090)
 * [3c6f8b807](https://github.com/kubeovn/kube-ovn/commit/3c6f8b807b1a2658edf42b0763d8fffd7683a925) chore(deps): update dependency cilium to v1.20.0 (#7097)
 * [67a6f9474](https://github.com/kubeovn/kube-ovn/commit/67a6f94742fe8c655e4e09fa9f35c3176eb31781) chore(deps): update dependency cert-manager to v1.21.1 (#7095)
 * [9f5c025e4](https://github.com/kubeovn/kube-ovn/commit/9f5c025e4be1725a6f7276becaf75c0ef21ef8e3) fix(deps): update module github.com/cert-manager/cert-manager to v1.21.1 (#7096)
 * [d7968a68e](https://github.com/kubeovn/kube-ovn/commit/d7968a68ee73691a0e1241c7a062c4fa77ab3201) fix: update Bingbing Zhang's GitHub username in MAINTAINERS (#7099)
 * [97cf8ffd7](https://github.com/kubeovn/kube-ovn/commit/97cf8ffd763c72b456c1275073256ad71eb59016) fix(deps): update module go.podman.io/image/v5 to v5.41.0 (#7098)
 * [4a05fdd4d](https://github.com/kubeovn/kube-ovn/commit/4a05fdd4dc250c39c387d8b1cb5fdf11937ba24c) test(e2e): cover per-pod IP family selection (#6985)
 * [c4731b5c6](https://github.com/kubeovn/kube-ovn/commit/c4731b5c6ee10b215a47e78c3c2f12934452dd69) docs: updated CHANGELOG.md (#7089)
 * [3d6b4f115](https://github.com/kubeovn/kube-ovn/commit/3d6b4f115a4a21dc17e958debf93cfcdd1d8534a) fix(deps): update module github.com/moby/moby/client to v0.5.1 (#7086)
 * [3042ef344](https://github.com/kubeovn/kube-ovn/commit/3042ef34481bf315acf09472475bf9157ce87065) feature(metallb): internal underlay pod traffic to LB service VIP node (#7040)
 * [0e942cae8](https://github.com/kubeovn/kube-ovn/commit/0e942cae8461432d1f7c82dc0881e0a0c0134473) ovs: skip address set update when the addresses are unchanged (#7084)
 * [dabf7b625](https://github.com/kubeovn/kube-ovn/commit/dabf7b625b43c54b460986c338be9a45f8e1915d) feat: infer pod subnet from named IPPool (#7059)
 * [b4b4a8801](https://github.com/kubeovn/kube-ovn/commit/b4b4a8801b5ab7123cd528138671a29d43f69f00) fix(e2e): ignore VPC egress gateway node order (#7081)
 * [27354cbfd](https://github.com/kubeovn/kube-ovn/commit/27354cbfda30b3bfe4094856b026530e7c1d25f6) controller: record vpc egress gateway events (#6988)
 * [42064b341](https://github.com/kubeovn/kube-ovn/commit/42064b341451c05c93c4fbd45cda7c4022209b50) controller, daemon: record pod network lifecycle events (#6989)
 * [6bd7e4854](https://github.com/kubeovn/kube-ovn/commit/6bd7e4854f65dcfe156119435b174a5c610f1d4d) fix(util): handle zero-compressed IPv6 addresses in IPv6ToLabelValue (#7077)
 * [066b1cf2c](https://github.com/kubeovn/kube-ovn/commit/066b1cf2c1e19e406ee0f43b5b0777cf8476c57e) feat(veg): allow co-located active-active replicas (#6789)
 * [8d22d3221](https://github.com/kubeovn/kube-ovn/commit/8d22d322141cc86cfd9338d7a58e10386bf27b12) fix(chart): correct v2 agent capability name (#7073)
 * [9b5317f36](https://github.com/kubeovn/kube-ovn/commit/9b5317f36551c9f3848721c6680467cfc102e421) fix(deps): update kubernetes (#7062)
 * [d7157ea47](https://github.com/kubeovn/kube-ovn/commit/d7157ea472cd2437cc6b20d88f47d27e0c8f7bcc) build(deps): bump github.com/google/cel-go from 0.28.1 to 0.29.0 (#7072)
 * [0c2d39058](https://github.com/kubeovn/kube-ovn/commit/0c2d39058d9ef37201e39a397c36956b837f1539) fix(deps): update module github.com/prometheus/client_golang to v1.24.1 (#7071)
 * [17be8af7a](https://github.com/kubeovn/kube-ovn/commit/17be8af7ad2494bc79908acc099520fa5f4a2593) fix(ovn-central): self-recover from stuck joining cluster state (#7039)
 * [0945764a3](https://github.com/kubeovn/kube-ovn/commit/0945764a36ac309c198bb29913d45b20b7524f69) fix(controller): extend BFD liveness probe timeout (#7060)
 * [29dd544b8](https://github.com/kubeovn/kube-ovn/commit/29dd544b89f781e24689368a77506dbd96f3cf98) Revert "fix: prevent subnet status update loop in patchSubnetStatus (#7058)" (#7063)
 * [13b9f2f8b](https://github.com/kubeovn/kube-ovn/commit/13b9f2f8bd64b961f03f8a9c1e4833a28ad61e22) fix: prevent subnet status update loop in patchSubnetStatus (#7058)
 * [1e3541810](https://github.com/kubeovn/kube-ovn/commit/1e3541810120903c55571a3edf0c977d4a5c50ac) test(gateway): limit bfdd zero-session recovery to v1.15+ (#7053)
 * [10565d4fa](https://github.com/kubeovn/kube-ovn/commit/10565d4fa73afaa911aa757c8429e5b2d1b880f2) fix(deps): update module github.com/prometheus/client_golang to v1.24.0 (#7056)
 * [a49921180](https://github.com/kubeovn/kube-ovn/commit/a49921180e9c074c4ae5889232f535cc0a320e2b) chore(deps): update actions/setup-python action to v7 (#7057)
 * [80d262ae5](https://github.com/kubeovn/kube-ovn/commit/80d262ae54e5dbcefe9cd88d66d9d57d87483231) fix(deps): update module github.com/go-logr/logr to v1.4.4 (#7055)
 * [c337e0b29](https://github.com/kubeovn/kube-ovn/commit/c337e0b29f8a58d06e79b09dddbe707b6ca3e1b2) chore(deps): update dependency talos to v1.13.7 (#7054)
 * [c15d21a54](https://github.com/kubeovn/kube-ovn/commit/c15d21a5482606eceb9c6f9c643b9d3577d1edc3) fix(gateway): recover missing BFD sessions (#7048)
 * [4b29ea03b](https://github.com/kubeovn/kube-ovn/commit/4b29ea03b267248ffa284dd1ed775e001c3fc5eb) chore(deps): update actions/checkout action to v7.0.1 (#7049)
 * [75ccf93de](https://github.com/kubeovn/kube-ovn/commit/75ccf93de83011284a70a93b4b71ebfc3c64ac4d) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.3.5 (#7050)
 * [d7785911d](https://github.com/kubeovn/kube-ovn/commit/d7785911df4921bdb9cf9eb0a1d12eca09ec4329) feat: support share type iptables DNAT with nft map-based LB (#6858)
 * [3d1be3557](https://github.com/kubeovn/kube-ovn/commit/3d1be355790ca4f83f2c78bec3fd580017e03255) fix(deps): update module google.golang.org/grpc to v1.82.1 (#7044)
 * [3486b3f8c](https://github.com/kubeovn/kube-ovn/commit/3486b3f8c9fdbb56cc23f6eacb9983cd6594170a) chore(deps): update dependency cilium to v1.19.6 (#7043)
 * [62860db91](https://github.com/kubeovn/kube-ovn/commit/62860db9169e9d162b5a543880df04d7eb8818a0) fix(deps): update module github.com/cenkalti/backoff/v6 to v7 (master) (#6956)
 * [db1e2d00a](https://github.com/kubeovn/kube-ovn/commit/db1e2d00ab812f93f285f506efb424aaf1e4abab) fix(chart): correct ovs-ipsec-keys volume indentation and Talos IPSEC docs (#7041)
 * [21df4f809](https://github.com/kubeovn/kube-ovn/commit/21df4f8096ab7f1ee60a3c38ec893fe20f0710ea) chore(deps): update actions/setup-go action to v7 (#7045)
 * [250dfda9a](https://github.com/kubeovn/kube-ovn/commit/250dfda9ab875a0607774bda2fc448c75f4ef3e1) perf(controller): skip irrelevant service and endpointslice updates (#7037)
 * [52ebbc10c](https://github.com/kubeovn/kube-ovn/commit/52ebbc10c317b130df4b5e18126670620b29d0f1) fix: recreate ACLs when ANP/BANP/CNP rule is renamed (#7038)
 * [cf5e46e27](https://github.com/kubeovn/kube-ovn/commit/cf5e46e27b259ebd8c41421febdeff5cc837f4b0) fix: install CNI config and binaries atomically (#7035)
 * [c944e314a](https://github.com/kubeovn/kube-ovn/commit/c944e314affd20761aa8dcda8295201b6bcf94ed) fix ipam only ip leak (#7034) (#7036)
 * [5843b67bd](https://github.com/kubeovn/kube-ovn/commit/5843b67bd77508d85ec099044f171c50ff7e850f) fix(deps): update module github.com/containerd/containerd/v2 to v2.3.3 (#7033)
 * [c0a616674](https://github.com/kubeovn/kube-ovn/commit/c0a61667493fa720f3834df058168398837f023f) ci: clean up chart releases and keep latest release correct (#7032)
 * [369e4b754](https://github.com/kubeovn/kube-ovn/commit/369e4b754a197f9b65d85556b72fec85e54807b8) docs: updated CHANGELOG.md (#7031)
 * [e92898096](https://github.com/kubeovn/kube-ovn/commit/e9289809648a8260b1badd6c711a4e3a920fee6a) chart: support HCP ovn-central scheduling (#7014)
 * [779ce06c4](https://github.com/kubeovn/kube-ovn/commit/779ce06c4e8a54a83582b73aec7dfb61addb5cdf) chore(deps): update dependency cert-manager to v1.21.0 (#7027)
 * [70b90aded](https://github.com/kubeovn/kube-ovn/commit/70b90aded4f3931956151b77fa73e9775b613d42) chore(deps): update dependency helm to v4.2.3 (#7028)
 * [3416b9c6a](https://github.com/kubeovn/kube-ovn/commit/3416b9c6a005df65c10e034c886b91508d337e18) fix(deps): update module golang.org/x/tools to v0.48.0 (#7029)
 * [dc9adfb98](https://github.com/kubeovn/kube-ovn/commit/dc9adfb98abeb4e8ae692cf41272f857255b8f8a) chore(deps): update dependency talos to v1.13.6 (#7026)
 * [0e284f556](https://github.com/kubeovn/kube-ovn/commit/0e284f5560297f565c93c1a80aba48e4096d6665) fix(deps): update module github.com/cert-manager/cert-manager to v1.21.0 (#7020)
 * [4a1717f29](https://github.com/kubeovn/kube-ovn/commit/4a1717f29a9e70d67d43fe2f2e8f9d9a79b30906) fix(deps): update module sigs.k8s.io/structured-merge-diff/v6 to v6.4.2 (#7018)
 * [11102db88](https://github.com/kubeovn/kube-ovn/commit/11102db88e3d450b41f5fc5d00e5764184b78835) fix(deps): update golang (#7019)
 * [9c3d07ee7](https://github.com/kubeovn/kube-ovn/commit/9c3d07ee745b10ad1838a9f5cd4028380535c4b3) fix: drop orphan centralized nat outgoing tcp packets (#6998)
 * [4341a6be1](https://github.com/kubeovn/kube-ovn/commit/4341a6be1c54becde321a72396b64c368ce11eb0) chore(deps): update dependency go to v1.26.5 (#7010)
 * [2017b3867](https://github.com/kubeovn/kube-ovn/commit/2017b3867b773760f09535fb12d2a3b30efda589) build: add OCI labels to images (#7007)
 * [6471a4656](https://github.com/kubeovn/kube-ovn/commit/6471a46569aaa9b62ea83bc5c49b9b061e5cbfc9) fix(deps): update module sigs.k8s.io/structured-merge-diff/v6 to v6.4.1 (#7009)
 * [4f9b1ed65](https://github.com/kubeovn/kube-ovn/commit/4f9b1ed65233d0e108d93aaaa79a6ba64ea17a43) docs: updated CHANGELOG.md (#7005)
 * [eefd2d24e](https://github.com/kubeovn/kube-ovn/commit/eefd2d24eb1623274e3682d4b2a196863f4cd467) chore(deps): update module golang.org/x/text to v0.39.0 (#7001)
 * [d04f5e4fd](https://github.com/kubeovn/kube-ovn/commit/d04f5e4fdd23ea3f1ef887bd343eb9d7488d2795) feat: support HCP ovn-central chart deployment (#6978)
 * [f29bd7a91](https://github.com/kubeovn/kube-ovn/commit/f29bd7a91ff6dce0ce54c1a00bffcd586076d100) fix: restore ip_src hash for policy route ECMP (#6953)
 * [8a9efc400](https://github.com/kubeovn/kube-ovn/commit/8a9efc400e16733cb645840d76528daca0b21095) fix: validate vlan id within range [0, 4095] (#6995)
 * [55bb00aac](https://github.com/kubeovn/kube-ovn/commit/55bb00aac34c9fabae3a07103611b145b6361495) fix: avoid wrapping nil error with %w in Errorf (#6994)
 * [22d08892f](https://github.com/kubeovn/kube-ovn/commit/22d08892f3bd2912670e912a8ebd01da4735d821) fix(kubectl-ko): honor KUBE_OVN_NS instead of hardcoded kube-system namespace (#6992)
 * [36e02eb39](https://github.com/kubeovn/kube-ovn/commit/36e02eb396bb108b81c97793a85aad46a723f016) fix: reject subnet gateway static IP allocation (#6980)
 * [5b2974ca8](https://github.com/kubeovn/kube-ovn/commit/5b2974ca8b5bfde9f9e49f51d7ce02f592f391fd) fix: require subnet for attachment provider (#6957)
 * [5a47910af](https://github.com/kubeovn/kube-ovn/commit/5a47910af204dffb5d95ac84ba3608209135f986) build(deps): bump golang.org/x/net from 0.48.0 to 0.55.0 in /test/anp (#6984)
 * [fec402dbd](https://github.com/kubeovn/kube-ovn/commit/fec402dbd001c4bf2d68c0f70997e7ba9e7b8450) feat: configure OVN DB server TLS options (#6967)
 * [c001e6c85](https://github.com/kubeovn/kube-ovn/commit/c001e6c85ab2bfc02c8da1b3ab7ed92b02bccd08) fix: improve BGP BFD metrics, logging, type safety and peer state check (#6941)
 * [6f1793c6b](https://github.com/kubeovn/kube-ovn/commit/6f1793c6bc81b515a7fa96f5dc36485e79301a1c) feat: support per-pod IP family selection (#6950)
 * [ad6c7a1ac](https://github.com/kubeovn/kube-ovn/commit/ad6c7a1ac983b2ae1f5e08369c47227f1a88dd27) fix veg init for mismatched IP families (#6958)
 * [5797a9e70](https://github.com/kubeovn/kube-ovn/commit/5797a9e705546cbc77f82b9d6d9ae04d468bf610) fix(deps): update module github.com/osrg/gobgp/v4 to v4.7.0 (#6973)
 * [bf0d892a9](https://github.com/kubeovn/kube-ovn/commit/bf0d892a98216a1632e89189026870266610727f) build(deps): bump golang.org/x/net from 0.48.0 to 0.55.0 in /test/cnp (#6974)
 * [9492fa1ed](https://github.com/kubeovn/kube-ovn/commit/9492fa1ede1fdda581d78e46f59464b086a43e32) chore(deps): update dependency gobgp to v4.7.0 (#6972)
 * [f29e959ef](https://github.com/kubeovn/kube-ovn/commit/f29e959ef5cd971d253c4b0f68b3d7d4c99722f2) fix(deps): update module google.golang.org/grpc to v1.82.0 (#6955)
 * [cc88e8099](https://github.com/kubeovn/kube-ovn/commit/cc88e80992fe6d44816f99b17263362b616293ae) docs: updated CHANGELOG.md (#6971)
 * [542b1f19f](https://github.com/kubeovn/kube-ovn/commit/542b1f19f6ae7b49b4a07e4b12c69166e70af040) fix: skip non-security-group port groups in sg gc (#6959)
 * [01a9c3c42](https://github.com/kubeovn/kube-ovn/commit/01a9c3c42c80405d86249ce1303893042c7e85d2) build: refresh kubectl golang.org/x deps and CNI plugin rebuilds (#6929) (#6951)
 * [91a995b1c](https://github.com/kubeovn/kube-ovn/commit/91a995b1c1dbcde0c985347c130b40fc9d9c8af3) fix vpc egress local policy without bfd (#6927)
 * [5e39ffb30](https://github.com/kubeovn/kube-ovn/commit/5e39ffb303b1b4568b5e65dc0f5b46f43f7cce58) fix: allow multi SNAT with same logical_ip (#6942)
 * [aab9792e5](https://github.com/kubeovn/kube-ovn/commit/aab9792e581cb2f3e672bedb0ece2234ea5418dd) build: include ovs and ovn commits in version (#6947)
 * [d0b301924](https://github.com/kubeovn/kube-ovn/commit/d0b301924fe542cdf1ecc3a250a47b30d9a4f942) fix(deps): update module github.com/cert-manager/cert-manager to v1.20.3 (#6935)
 * [077ffd84d](https://github.com/kubeovn/kube-ovn/commit/077ffd84d7ac7d35b3ac23d441c9edc29e563801) chore(deps): update dependency cert-manager to v1.20.3 (#6934)
 * [1a61c6746](https://github.com/kubeovn/kube-ovn/commit/1a61c6746d4e19c00db49eb6c08278724d706d86) fix(deps): update module github.com/prometheus-community/pro-bing to v0.9.1 (#6949)
 * [fddb9d0dd](https://github.com/kubeovn/kube-ovn/commit/fddb9d0ddd64bda8d193cab87aeb601a40807bd2) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.3.4 (#6948)
 * [d9e3c7327](https://github.com/kubeovn/kube-ovn/commit/d9e3c7327f77a1a1f04e4dcf7f7fd4fe18c9dce0) fix(daemon): make tproxy work with kube-ovn as secondary CNI (#6850)
 * [8c5dadb28](https://github.com/kubeovn/kube-ovn/commit/8c5dadb28cd1d58e6e43964fc7238081b87afb35) fix: unblock IPv6 OVN EIP deletion and harden underlay e2e against docker network races (#6872)
 * [cc244dcc1](https://github.com/kubeovn/kube-ovn/commit/cc244dcc1c825ae32427dd196f52b2c21098af84) feat: expose daemon host-tunnel-src flag in install.sh and helm charts (#6944)
 * [b7d23c18c](https://github.com/kubeovn/kube-ovn/commit/b7d23c18ca2c9e75d35307500039e800c023339a) feat: add kube-ovn TLS rotation e2e test (#6939)
 * [dde0481d0](https://github.com/kubeovn/kube-ovn/commit/dde0481d079b373bb360f3104f68cc4904cb770f) docs: updated CHANGELOG.md (#6940)
 * [e46da26e1](https://github.com/kubeovn/kube-ovn/commit/e46da26e14eb80fb10b58199d721136c9603ded5) fix(deps): update module golang.org/x/tools to v0.47.0 (#6936)
 * [a889faf67](https://github.com/kubeovn/kube-ovn/commit/a889faf67ea2a5c73bf0b8471a51bdb81da91e2b) e2e: add version check (#6930)
 * [a4c26e037](https://github.com/kubeovn/kube-ovn/commit/a4c26e037ed7cdd7bd4bb2ccf375282d0c4225a0) feat: Add kubeovn support for dynamic ip allocation/byo dhcp (#6885)
 * [f61946aa8](https://github.com/kubeovn/kube-ovn/commit/f61946aa87a6f96d2beefa935981b7d931c95b1e) chore(deps): update actions/setup-python action to v6.3.0 (#6926)
 * [6ca4c9cf2](https://github.com/kubeovn/kube-ovn/commit/6ca4c9cf2a03d88eca965e3e82d8b54d4e50bd07) feat: rotate kube-ovn TLS certificates (#6877)
 * [a9c0a221c](https://github.com/kubeovn/kube-ovn/commit/a9c0a221cbfc30f5e64b8057980c856e958493c8) feat(speaker): add GoBGP BFD support for fast failure detection (#6912)
 * [71737cb29](https://github.com/kubeovn/kube-ovn/commit/71737cb2975396f876a98fa7b52e5c65ad7c8cab) fix(ovs): improve dnat_and_snat Add with delete+create for pod swap (#6914)
 * [cec612aa3](https://github.com/kubeovn/kube-ovn/commit/cec612aa3c9e0146f6e6a4143c2a00101a29c06b) fix(deps): update module github.com/onsi/gomega to v1.42.1 (#6924)
 * [bfbed753e](https://github.com/kubeovn/kube-ovn/commit/bfbed753ecfddf7506eacea46dfc0dc7affb6713) chore(deps): update azure/setup-helm action to v5.0.1 (#6923)
 * [206565728](https://github.com/kubeovn/kube-ovn/commit/206565728490a247fa559444cbcb50c5b676e44d) chore(deps): update actions/cache action to v6 (#6922)
 * [f580c9f99](https://github.com/kubeovn/kube-ovn/commit/f580c9f9998826257d68978ac01409402f72ddb9) chore(deps): update dependency kwork to v0.8.0 (#6921)
 * [b660723b9](https://github.com/kubeovn/kube-ovn/commit/b660723b92e3ef8c9d73e543e2a56786119cbe54) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.32.0 (#6920)
 * [2cd4ad21c](https://github.com/kubeovn/kube-ovn/commit/2cd4ad21ce8d139fbe8a7484fb7e52809ae85036) chore(deps): update dependency talos to v1.13.5 (#6919)
 * [c75fe9f14](https://github.com/kubeovn/kube-ovn/commit/c75fe9f14ee4dec5817d4c32fbb2422bc9de61e0) docs: updated CHANGELOG.md (#6918)
 * [e2eaf3e96](https://github.com/kubeovn/kube-ovn/commit/e2eaf3e96ac5b3dde9e9f9670f8c34861e82e8ea) docs: updated CHANGELOG.md (#6917)
 * [8637c682a](https://github.com/kubeovn/kube-ovn/commit/8637c682a3116ef8441ff18a6270eb883678f65f) fix(crd): restore vpc egress gateway validations (#6913)
 * [b10f061c7](https://github.com/kubeovn/kube-ovn/commit/b10f061c72199a2861311c60d32e63873b10f6fa) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.3.3 (#6908)
 * [302810e0f](https://github.com/kubeovn/kube-ovn/commit/302810e0f6ba08420837994ca39f2385a365251c) chore(deps): update actions/checkout action to v7 (#6909)
 * [a37308a3f](https://github.com/kubeovn/kube-ovn/commit/a37308a3f9d92f64e74d6546b486efa35662efd5) fix(deps): update module github.com/cenkalti/backoff/v5 to v6 (master) (#6910)
 * [b1cccbd38](https://github.com/kubeovn/kube-ovn/commit/b1cccbd388bec6cf64e2a59582249bfee49629b7) fix(ovs): make DeleteNat idempotent when target NAT row is absent (#6911)
 * [1c6b8bce2](https://github.com/kubeovn/kube-ovn/commit/1c6b8bce21e6412f83eb183184017718377f60f1) fix(deps): update module github.com/k8snetworkplumbingwg/sriovnet to v1.3.0 (#6903)
 * [aaddcfb5a](https://github.com/kubeovn/kube-ovn/commit/aaddcfb5ad0741a9f698e0bdb7c8a342b2b7e124) fix(deps): update module github.com/containerd/containerd/v2 to v2.3.2 (#6905)
 * [9a834ad06](https://github.com/kubeovn/kube-ovn/commit/9a834ad063a6c35bcae991eef887cfded3e56662) fix(deps): update module github.com/moby/moby/client to v0.5.0 (#6906)
 * [ba6de188d](https://github.com/kubeovn/kube-ovn/commit/ba6de188d32243b3af8f85cbc52ec6dc7cc565b6) fix(ci): repair dpdk base image build (#6894)
 * [69bf350e7](https://github.com/kubeovn/kube-ovn/commit/69bf350e73f69b49d0a6c8db308381d6e9d645d5) fix: reconcile vpc bfd ha chassis on node changes (#6891)
 * [7be3a8138](https://github.com/kubeovn/kube-ovn/commit/7be3a8138646d2d738adf8c2ad388af4415120f6) chore(deps): update dependency kubevirt to v1.8.4 (#6887)
 * [5f7ed4db4](https://github.com/kubeovn/kube-ovn/commit/5f7ed4db42a6123a91527cd1adcb2eade9e6feba) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.3.2 (#6881)
 * [70fe40375](https://github.com/kubeovn/kube-ovn/commit/70fe403755756d771eca8af20ea6ced2f5b0df49) fix(controller): enable ARP learning on VPC routers that have peerings (#6882)
 * [9a059f91e](https://github.com/kubeovn/kube-ovn/commit/9a059f91ea1cfae93be422f6d64ce9e374c60bb0) fix(deps): update module kubevirt.io/api to v1.8.4 (#6888)
 * [df19a5a13](https://github.com/kubeovn/kube-ovn/commit/df19a5a1320d146ddf649087f447a9180cb3dc89) chore(deps): update dependency helm to v4.2.2 (#6889)
 * [d22d12508](https://github.com/kubeovn/kube-ovn/commit/d22d125080aa4bab396638bde2eaf42db85c78cf) fix(deps): update module kubevirt.io/client-go to v1.8.4 (#6890)
 * [9ff226e6d](https://github.com/kubeovn/kube-ovn/commit/9ff226e6d5da250c497e7f20579a82d3befc08ab) fix(subnet): allow all-ones host as gateway for IPv6 subnet (#6883)
 * [85b935b57](https://github.com/kubeovn/kube-ovn/commit/85b935b57b44b3711bdb713dcdf5f0feb17c86f6) chore(deps): update dependency cilium to v1.19.5 (#6880)
 * [18dad3818](https://github.com/kubeovn/kube-ovn/commit/18dad3818366a7795af51489b68c2100986fb5a6) fix(ci): correct cnp-domain e2e events artifact name (#6879)
 * [0eb7cd5a0](https://github.com/kubeovn/kube-ovn/commit/0eb7cd5a0d292eed71bd0921a3db44274ee92494) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.31.0 (#6874)
 * [e7d029b10](https://github.com/kubeovn/kube-ovn/commit/e7d029b1009e6d0eb1c58a6d4eea8f3c4c1977ea) fix(deps): update module github.com/onsi/gomega to v1.42.0 (#6875)
 * [7a6ceb4a3](https://github.com/kubeovn/kube-ovn/commit/7a6ceb4a37fb65ec6fe73df1c3fb601d6d7d941d) docs: updated CHANGELOG.md (#6876)
 * [97e784f65](https://github.com/kubeovn/kube-ovn/commit/97e784f651df965c0386ef06e903f994c8489257) feat(vpcnatgateway): add high-availability for nat gateways (#6661)
 * [bba147eeb](https://github.com/kubeovn/kube-ovn/commit/bba147eebc95ef54a85b8a348725223971a73fbd) chore(deps): update dependency helm to v4.2.1 (#6866)
 * [c756f91d4](https://github.com/kubeovn/kube-ovn/commit/c756f91d4cc9e5fc7ba25f155627380fa556424b) fix(deps): update golang to v0.46.0 (#6869)
 * [11ac4e649](https://github.com/kubeovn/kube-ovn/commit/11ac4e649c79398b836c9ecba4a7a2235a4062ec) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.30.0 (#6870)
 * [b522524d1](https://github.com/kubeovn/kube-ovn/commit/b522524d19d370d6826fe469aaba9eaee4a82b76) fix(deps): update module gopkg.in/k8snetworkplumbingwg/multus-cni.v4 to v4.3.0 (#6871)
 * [95886265b](https://github.com/kubeovn/kube-ovn/commit/95886265b543b76fa6f622d37a1c14b4a46a3302) fix(deps): update kubernetes (#6867)
 * [1553501d7](https://github.com/kubeovn/kube-ovn/commit/1553501d7e9ad1fb6c78e93f43e687edbde22ea4) chore(deps): update dependency multus to v4.3.0 (#6868)
 * [9949513fc](https://github.com/kubeovn/kube-ovn/commit/9949513fcb7f092667d7a68eb1a3d1c97fa7c538) fix(ovs): delete the actual duplicate policy instead of the kept one in batch dedup (#6849)
 * [5bf7cd877](https://github.com/kubeovn/kube-ovn/commit/5bf7cd877a50cbccb1d23fc025caf0a4cb2e3f45) fix: misc bug fixes from code audit (#6856)
 * [c16eb3219](https://github.com/kubeovn/kube-ovn/commit/c16eb3219d70eecca88350f592a07c4e66972b6e) feat(controller): add RouterLBRule CRD for L3/L4 load balancing on VPC router (#6563)
 * [4de891201](https://github.com/kubeovn/kube-ovn/commit/4de891201d2ae566d1413b818853bcf58daf64e2) fix(controller): skip service/endpointslice enqueueing when lb is disabled (#6859)
 * [97e10e17b](https://github.com/kubeovn/kube-ovn/commit/97e10e17b3bc8a913582eedd4565e9a1ec4c100a) fix(images): adapt OVN pinctrl to new OVS compose_nd_ns signature (#6857)
 * [763d34567](https://github.com/kubeovn/kube-ovn/commit/763d345670cdb947a9734429974412e082f2593f) fix(deps): update golang to v0.56.0 (#6852)
 * [00d052f17](https://github.com/kubeovn/kube-ovn/commit/00d052f178ac85e9685ab0ffa9f6cc593f5659ec) chore(deps): update dependency talos to v1.13.4 (#6851)
 * [94e76203a](https://github.com/kubeovn/kube-ovn/commit/94e76203ae9deb3604ed6df21dd3c83853197061) fix(deps): update module github.com/prometheus-community/pro-bing to v0.9.0 (#6853)
 * [d8bde8605](https://github.com/kubeovn/kube-ovn/commit/d8bde8605d25dae8b7977f1565d53ac91703ba94) fix(controller): correct EIP v6 label update in SNAT/DNAT status patch (#6846)
 * [4d9089020](https://github.com/kubeovn/kube-ovn/commit/4d908902054adc9e406f8eefeac93d00bba6c779) fix(controller): avoid nil queue panic in DNS name resolver handlers (#6845)
 * [f2f03f639](https://github.com/kubeovn/kube-ovn/commit/f2f03f639e222374de05558d9ef4dd9e2c0a1f13) fix(controller): avoid wiping priority-0 CNP map entries on zero-value key (#6844)
 * [1fa2b218d](https://github.com/kubeovn/kube-ovn/commit/1fa2b218de28bf8be67324dd73272e38f8a5c225) fix(controller): avoid nil pointer panic when attachment default subnet is gone (#6843)
 * [1299ad3be](https://github.com/kubeovn/kube-ovn/commit/1299ad3bef503bc70e7de6aaf466b7f00590ebf6) fix(ovs): check error before dereferencing logical router on static route deletion (#6842)
 * [a229dec3d](https://github.com/kubeovn/kube-ovn/commit/a229dec3d5150f55ca8344d0e296a8a02d52b726) fix(controller): re-enqueue VMI migration when target pod is not yet … (#6834)
 * [e67b5e586](https://github.com/kubeovn/kube-ovn/commit/e67b5e586487732ba391b77cfe11b25cff281677) test(e2e): require v1.17+ for non-primary CNI iptables cleanup test (#6840)
 * [bd2ffe6ee](https://github.com/kubeovn/kube-ovn/commit/bd2ffe6ee9c81a8ae5b07d823611042b62ba3034) fix(controller): avoid nil pointer panic in chassis and static route gc (#6835)
 * [9194883e1](https://github.com/kubeovn/kube-ovn/commit/9194883e10a55e9e56a20d929963164d127a6388) fix(controller): skip orphaned subnet instead of breaking namespace subnet loop (#6836)
 * [58f9223cc](https://github.com/kubeovn/kube-ovn/commit/58f9223cc2d33c16ca14f602a153b743b3b307e2) fix(deps): update golang (#6837)
 * [c5eabdb55](https://github.com/kubeovn/kube-ovn/commit/c5eabdb556ce37ed92fbed34d9d24bf02c9cc41d) perf(controller): skip redundant service update in checkServiceLBIPBelongToSubnet (#6828)
 * [7a475d22a](https://github.com/kubeovn/kube-ovn/commit/7a475d22a4eeadd08b0896001faed81bd6c19b30) perf(ovs): list NAT and static routes via single WhereCache scan (#6829)
 * [db385f28d](https://github.com/kubeovn/kube-ovn/commit/db385f28dc149078c9094c9cf660285d92186df0) perf(daemon): list subnets and pods once per gateway setter (#6827)
 * [ced939754](https://github.com/kubeovn/kube-ovn/commit/ced9397544c3cf76738eb2c648e2f77c2ba846df) perf(daemon): list nat PREROUTING once when ensuring kube-ovn rule (#6826)
 * [b78caec4c](https://github.com/kubeovn/kube-ovn/commit/b78caec4c450b47bcb5f8073dcfbbaf37d1c4e61) perf(controller): list VMs cluster-wide once per GC cycle (#6824)
 * [9fd6e90f3](https://github.com/kubeovn/kube-ovn/commit/9fd6e90f38ea6197a3e686ab20294df842a5a1b7) docs: updated CHANGELOG.md (#6833)
 * [fd4cdc93e](https://github.com/kubeovn/kube-ovn/commit/fd4cdc93ea90c7962afec4da75114ba0c7dcd7bc) docs: updated CHANGELOG.md (#6825)
 * [d51808150](https://github.com/kubeovn/kube-ovn/commit/d51808150c0dd06a17b9fea471acbd891fed64a8) add .claude to git ignore
 * [bbd863eb0](https://github.com/kubeovn/kube-ovn/commit/bbd863eb01f137685ebb3ad43569aabfb1bb1755) docs: updated CHANGELOG.md (#6821)
 * [7780155a3](https://github.com/kubeovn/kube-ovn/commit/7780155a3db6c3e77ac1d0dd468ade0344c6ccd0) chore(deps): update dependency kubevirt to v1.8.3 (#6818)
 * [209d1cf65](https://github.com/kubeovn/kube-ovn/commit/209d1cf657f9a9ca39e14e0761356903d95997ec) fix(deps): update module kubevirt.io/api to v1.8.3 (#6819)
 * [a4e53d71b](https://github.com/kubeovn/kube-ovn/commit/a4e53d71b6f487c657e30af7a7db2df530ee4169) feat: non-primary CNI not create any service-related iptables rules (#6462)
 * [5849eb0f6](https://github.com/kubeovn/kube-ovn/commit/5849eb0f6341b5cf1e3e4c9ff4a605e26bc69e0c) fix(deps): update module kubevirt.io/client-go to v1.8.3 (#6820)
 * [ceb7193b7](https://github.com/kubeovn/kube-ovn/commit/ceb7193b76760ec47232badf62e22c45084e08f5) feat: enable SSA applyconfigurations for go client (#6813)
 * [91d8563bf](https://github.com/kubeovn/kube-ovn/commit/91d8563bffed0989f66c9416149b64802edca27b) perf(ipam): allocation-free IP comparison in range engine (#6816)
 * [c76e9272e](https://github.com/kubeovn/kube-ovn/commit/c76e9272e21b6d81bc902d2d9a2d71e5b35dfcc5) perf(controller): index IP CRs by subnet for subnet status calculation (#6814)
 * [7a2b2a246](https://github.com/kubeovn/kube-ovn/commit/7a2b2a2463e40172a7f79432f843bdaecc4b70e2) PR introduces changes to correctly support multiple network interfaces from same subnet (#6766)
 * [4c435cc15](https://github.com/kubeovn/kube-ovn/commit/4c435cc154a2ae974ca7a5f6aaf2a25ee72e19d1) chore(deps): update golang to v1.26.4 (#6810)
 * [83984c104](https://github.com/kubeovn/kube-ovn/commit/83984c1040321d5b9cb8fef7120995cf2df1bfed) perf: reduce CPU and memory overhead on hot paths (#6811)
 * [62e1d531b](https://github.com/kubeovn/kube-ovn/commit/62e1d531b35c64b804c9ed5b3e74c11276b9927c) chore(deps): update kindest/node docker tag to v1.36.1 (#6806)
 * [35cd2a69e](https://github.com/kubeovn/kube-ovn/commit/35cd2a69e71ed35795fead35362d8053103d87cb) fix(controller): correct egress sync status field on ACL update failure (#6804)
 * [8961aec01](https://github.com/kubeovn/kube-ovn/commit/8961aec01d5f13038f29138fe05a9979d6cebffb) chore(deps): update dependency kind to v0.32.0 (#6807)
 * [7f13fcba2](https://github.com/kubeovn/kube-ovn/commit/7f13fcba24953b7415cdabbf8b1a4774497834c2) chore(deps): update actions/checkout action to v6.0.3 (#6805)
 * [07a3c31c9](https://github.com/kubeovn/kube-ovn/commit/07a3c31c9dd46b44ba907226a34352d6111a8acb) fix(deps): update module github.com/osrg/gobgp/v4 to v4.6.0 (#6802)
 * [fe6cfde59](https://github.com/kubeovn/kube-ovn/commit/fe6cfde59d6a00ed32f53418a68ff642ebb1ee19) chore(deps): update arp fork (#6796)
 * [61e679e67](https://github.com/kubeovn/kube-ovn/commit/61e679e67ddc505b0112429655882874bc67a496) chore(deps): update dependency gobgp to v4.6.0 (#6801)
 * [8bd070c1c](https://github.com/kubeovn/kube-ovn/commit/8bd070c1c6d4df79495223a6e3ae84b0b25e2e02) fix(controller): populate v6ip from cached EIP in handleAddVpcExternalSubnet (#6799)
 * [adf508427](https://github.com/kubeovn/kube-ovn/commit/adf508427105e6225ba32c215eab73c884cca85a) feat: single-replica ovn-central + PVC + LoadBalancer + Kamaji split-cluster deployment (#6793)
 * [f331bbb07](https://github.com/kubeovn/kube-ovn/commit/f331bbb071e069631000a92fb4f8979763f0fe65) fix(controller): skip in-pod cleanup when VPC NAT gateway CRD is deleted (#6784)
 * [2bfefecee](https://github.com/kubeovn/kube-ovn/commit/2bfefeceeb92fbeebdf30193deba4f9d08c21b1a) docs: updated CHANGELOG.md (#6797)
 * [252ba81d1](https://github.com/kubeovn/kube-ovn/commit/252ba81d1404325b61ef37e71dd650e261c7e07a) fix: restrict hook log file permissions (#6794)
 * [0ba96cacb](https://github.com/kubeovn/kube-ovn/commit/0ba96cacbca110d224b5435db1cfceef28e1c103) fix: report VEG workload network readiness (#6790)
 * [9ba19b264](https://github.com/kubeovn/kube-ovn/commit/9ba19b2648b7a455ff75f1b82685222173d45dc9) fix(veg): raise default bfdd CPU limit (#6788)
 * [512793750](https://github.com/kubeovn/kube-ovn/commit/512793750edba52ccd31466edcc9d08446d4328e) chore(deps): update dependency talos to v1.13.3 (#6786)
 * [3cab97cdb](https://github.com/kubeovn/kube-ovn/commit/3cab97cdb566d390308411b96577ed1353dde6d5) fix: avoid panic decoding invalid CA private keys (#6743)
 * [d327f73df](https://github.com/kubeovn/kube-ovn/commit/d327f73df1228440fd3679254a21eef42d34c236) fix(deps): update module go.universe.tf/metallb to v0.16.1 (#6792)
 * [221444dbb](https://github.com/kubeovn/kube-ovn/commit/221444dbb180cacbb35c9e6326d4ae6ca2329306) refactor(daemon): drop redundant Node DeepCopy in provider network sync (#6783)
 * [7e54fef90](https://github.com/kubeovn/kube-ovn/commit/7e54fef90020b91a1ef543e5677716738e0b1096) docs: updated CHANGELOG.md (#6781)
 * [9aa6246b9](https://github.com/kubeovn/kube-ovn/commit/9aa6246b9a7ecf50612eb7400cd75f0a33c10193) fix(controller): avoid panic in external VPC sync when no logical switch is found (#6777)
 * [fbed64995](https://github.com/kubeovn/kube-ovn/commit/fbed64995fe09c8266671eff41c32d373b27ad6c) fix(cniserver): probe unix socket from /livez to detect missing bind (#6778)
 * [e613bcee1](https://github.com/kubeovn/kube-ovn/commit/e613bcee17e3bfd621a1354b74223a9623533f77) fix(deps): update golang (#6776)
 * [62cb28f72](https://github.com/kubeovn/kube-ovn/commit/62cb28f721e370ce72f60de228300c3aed152cbb) chore(deps): update dependency metallb to v0.16.0 (#6763)
 * [c965e1998](https://github.com/kubeovn/kube-ovn/commit/c965e1998787a1c0931f31290e1542d6f59ed751) feat: allow use bfdd with kube-ovn bgp speaker (#6754)
 * [74b3f5695](https://github.com/kubeovn/kube-ovn/commit/74b3f5695ef2abed9b5daa5bc1bdc638ffea4e06) test(metallb): disable frr-k8s ServiceMonitor in kind install (#6773)
 * [98b3e5450](https://github.com/kubeovn/kube-ovn/commit/98b3e545022d7cf4c67ea0f0e01ea4062a8597d1) fix(deps): update module go.podman.io/image/v5 to v5.40.0 (#6770)
 * [b07e56494](https://github.com/kubeovn/kube-ovn/commit/b07e56494343195d4c228e13509750c52a1be706) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.3.1 (#6769)
 * [b57475d78](https://github.com/kubeovn/kube-ovn/commit/b57475d782a4abf9f2006b5b31d0add29182e585) fix(deps): update module github.com/containerd/containerd/v2 to v2.3.1 (#6765)
 * [bc72b1160](https://github.com/kubeovn/kube-ovn/commit/bc72b1160b5a58cedd5a011ad9352e29a3833b1b) fix(deps): update module go.universe.tf/metallb to v0.16.0 (#6764)
 * [d4559562d](https://github.com/kubeovn/kube-ovn/commit/d4559562d3ac1a0487e0096f8d224bcd2c1bcbd4) fix(deps): update module github.com/onsi/gomega to v1.41.0 (#6756)
 * [6ac134f01](https://github.com/kubeovn/kube-ovn/commit/6ac134f019a9a6f914f5f00aeadb089cb036abc6) docs: updated CHANGELOG.md (#6760)
 * [f77cd6311](https://github.com/kubeovn/kube-ovn/commit/f77cd6311861a1267de861597b4d19ecf4ebb20d) docs: updated CHANGELOG.md (#6759)
 * [e7707567d](https://github.com/kubeovn/kube-ovn/commit/e7707567d87e5ba1a8cc0faa91d74d1eddd77a1d) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.29.0 (#6755)
 * [af1ca3690](https://github.com/kubeovn/kube-ovn/commit/af1ca369020cfb1b8e59ccad1895f90cef331c6a) chore: apply modernize fixes flagged by latest gopls modernize (#6758)
 * [34feb6444](https://github.com/kubeovn/kube-ovn/commit/34feb64444042a20f093e48823ca255f98a17201) fix(deps): update kubernetes (master) (#6686)
 * [0a30fdc46](https://github.com/kubeovn/kube-ovn/commit/0a30fdc46457117e2ba2d59d1b9886274d2d92b6) feat(controller): distribute gateway chassis priority across VPCs (#6745)
 * [ca6231a54](https://github.com/kubeovn/kube-ovn/commit/ca6231a54047407da2a5df409e321c6cead0e0e0) feat: add u2o overlay only routing (#6736)
 * [9f4c41e89](https://github.com/kubeovn/kube-ovn/commit/9f4c41e89b51353ec4b65961b76683ae76633fe2) feat(cni): package ipvlan plugin in kube-ovn image (#6747)
 * [46689a97e](https://github.com/kubeovn/kube-ovn/commit/46689a97e04510acbbb3c5cc0c3abc31a79cdeae) fix(webhook): reject same-family duplicates in ip_address annotation (#6746)
 * [e56a9a807](https://github.com/kubeovn/kube-ovn/commit/e56a9a807b2424343074243d929c5ba185846fb4) chore(ovs): drop unused DeleteGatewayChassises helpers (#6742)
 * [ef8413511](https://github.com/kubeovn/kube-ovn/commit/ef8413511a4b737e38da5086cd0527e2562e7582) fix: default QoS interface based on direction when not specified (#6740)
 * [22388ad41](https://github.com/kubeovn/kube-ovn/commit/22388ad411ace7ab3ddc562d61b38c2ab4c02e4f) fix(daemon): close netns handles to prevent fd leaks (#6741)
 * [5a802e0f7](https://github.com/kubeovn/kube-ovn/commit/5a802e0f7989a2eaad286f66ba714e1bbf1ed184) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.15.0 (#6739)
 * [4f693c51d](https://github.com/kubeovn/kube-ovn/commit/4f693c51d11c0b74b079c3236bf134b893004a3f) chore(deps): update dependency golangci-lint to v2.12.2 (master) (#6695)
 * [35056ace1](https://github.com/kubeovn/kube-ovn/commit/35056ace1f09c07a4f1ba57e13ec03b04bbb801f) fix(controller): make BgpConf/EvpnConf informers optional (#6726) (#6737)
 * [3ade01c8f](https://github.com/kubeovn/kube-ovn/commit/3ade01c8fc3fb72a0ca9ccb39082c0a79c764845) fix(ipam): roll back partial state when random allocation fails (#6738)
 * [b87ecf49c](https://github.com/kubeovn/kube-ovn/commit/b87ecf49cd557b50b14a50b96b0ea69b737197c2) chore(deps): update dependency helm to v4.2.0 (#6734)
 * [85c06fee2](https://github.com/kubeovn/kube-ovn/commit/85c06fee2d3bd028e569b7811f2734e35821057b) fix(deps): update module google.golang.org/grpc to v1.81.1 (#6733)
 * [d2f12821f](https://github.com/kubeovn/kube-ovn/commit/d2f12821fb394ec3a773c2651756a87e25aff862) Allow set external-gateway-switch in chart (#6629)
 * [6c4ac7c81](https://github.com/kubeovn/kube-ovn/commit/6c4ac7c812d078a9c512e3f1163b6bd757d91920) chore(deps): update dependency cilium to v1.19.4 (#6730)
 * [48cd3ac99](https://github.com/kubeovn/kube-ovn/commit/48cd3ac997fce2b52866f0822a0eef345c91b808) docs: updated CHANGELOG.md (#6728)
 * [116a8b70e](https://github.com/kubeovn/kube-ovn/commit/116a8b70ea15a4d71be189a90c1579c60b6b58fd) docs: updated CHANGELOG.md (#6727)
 * [1b51dfa96](https://github.com/kubeovn/kube-ovn/commit/1b51dfa960a932d625606077cd20d4c0d785c850) chore(deps): update dependency talos to v1.13.2 (#6725)
 * [f64b67bc9](https://github.com/kubeovn/kube-ovn/commit/f64b67bc90b4efa678dcbb6a26086be421781ab1) chore(deps): update dependency submariner to v0.24.0 (#6724)
 * [c708e839d](https://github.com/kubeovn/kube-ovn/commit/c708e839d8f2bc509ffa9e896253b1d7e1e7f57c) Block vpc deletion while SNAT, DNAT or FIP rules are using it (#6718)
 * [268def680](https://github.com/kubeovn/kube-ovn/commit/268def680dfd5e381254f4e91e4a46f212d66c38) Revert "fix(ha): detect and recover from split-brain leader with incomplete cluster membership (#6363)" (#6721)
 * [d379ca38a](https://github.com/kubeovn/kube-ovn/commit/d379ca38afc0474f56048061ae94745ad1983b9b) speaker: cache startup-selected bgp source addresses (#6712)
 * [0d320b626](https://github.com/kubeovn/kube-ovn/commit/0d320b62698989751414e2672faa93b2220cf087) fix: use status variable instead of hardcoded ConditionTrue in loopOvn0Check (#6719)
 * [882758abe](https://github.com/kubeovn/kube-ovn/commit/882758abe2712b016a680659c5d4ac5c545b8b8a) test(iptables-vpc-nat-gw): skip image+annotations spec on v1.15- (#6714)
 * [25b872402](https://github.com/kubeovn/kube-ovn/commit/25b87240208227d052b94748f361c804b3f04c75) fix(deps): update golang (#6711)
 * [64f7433af](https://github.com/kubeovn/kube-ovn/commit/64f7433af51dcd1fd1ad6c145f4d9741d37247b9) chore(deps): update dependency go to v1.26.3 (#6704)
 * [ec8b28019](https://github.com/kubeovn/kube-ovn/commit/ec8b2801946fb317748ffa803dbccb1cbc689dc2) fix(image): rebase OVN patches for branch-25.03 ARP/ND localnet rename (#6707)
 * [00896432e](https://github.com/kubeovn/kube-ovn/commit/00896432e856ea337d64e735aea01cdee9ef02f8) feat(servicecidr): support K8s multiple ServiceCIDR (KEP-1880) (#6690)
 * [4b35ee767](https://github.com/kubeovn/kube-ovn/commit/4b35ee767e88e608f47b6f61a142045886bb3f5b) refactor(security): pass argv to exec instead of shell concatenation (#6700)
 * [e9f92e8d0](https://github.com/kubeovn/kube-ovn/commit/e9f92e8d00ccd0076fafff4400631d471357489e) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.3.0 (#6699)
 * [0c7c92c1e](https://github.com/kubeovn/kube-ovn/commit/0c7c92c1e221420fbf0d17eb5ebd3ab16010c2ea) fix(deps): update module google.golang.org/grpc to v1.81.0 (#6698)
 * [eb72911f8](https://github.com/kubeovn/kube-ovn/commit/eb72911f87b2f833a750025ff4e9de3e53c699cc) fix(deps): update module github.com/osrg/gobgp/v4 to v4.5.0 (#6696)
 * [310539bea](https://github.com/kubeovn/kube-ovn/commit/310539bea107bac0b3702c231a6c5d215c4ddfe2) chore(deps): update dependency gobgp to v4.5.0 (#6694)
 * [0acfdf863](https://github.com/kubeovn/kube-ovn/commit/0acfdf8630e37150dc4ffc5e6fd275fbdcd1a268) fix(deps): update module github.com/containerd/containerd/v2 to v2.3.0 (#6693)
 * [c5cfdbe43](https://github.com/kubeovn/kube-ovn/commit/c5cfdbe430628498896259a5691cffb6583aa3d5) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.28.3 (#6691)
 * [d4e827ffa](https://github.com/kubeovn/kube-ovn/commit/d4e827ffa9f3f5e5821994bd0c6be9f8f7817cc6) chore(deps): update dependency talos to v1.13.0 (master) (#6688)
 * [ecd9b02ce](https://github.com/kubeovn/kube-ovn/commit/ecd9b02cef341f28186da48d1c42fb60bacc9072) fix(deps): update module go.uber.org/zap to v1.28.0 (#6689)
 * [161358a18](https://github.com/kubeovn/kube-ovn/commit/161358a187bb454b62d4cf7a70c5be956079f2dc) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.28.2 (#6687)
 * [2320e01a6](https://github.com/kubeovn/kube-ovn/commit/2320e01a6d6ae2b19d4fb1fe3acab08d0f41fe64) fix(subnet): reject sub-1280 MTU on IPv6 subnets, warn at runtime (#6685)
 * [a775a110f](https://github.com/kubeovn/kube-ovn/commit/a775a110f2e10b54b5348a01080c24a36acf54d1) fix(deps): update kubernetes (master) (#6668)
 * [1992d6da0](https://github.com/kubeovn/kube-ovn/commit/1992d6da0ae37d66d14d571498c77e2176bdbe80) fix(vpcnatgw): order SHARED_SNAT rules by longest-prefix-first (#6680)
 * [794c3e10b](https://github.com/kubeovn/kube-ovn/commit/794c3e10bda73ff6a2e85b07f26741f62330781c) chore(deps): update dependency talos to v1.12.7 (#6682)
 * [b6c7e7cc1](https://github.com/kubeovn/kube-ovn/commit/b6c7e7cc1ae2f5ca704a0bcf2d7f00a093c645a6) chore(deps): bump vishvananda/netlink to main (af2a3ea5) (#6679)
 * [29998aa29](https://github.com/kubeovn/kube-ovn/commit/29998aa2914ad7703c772efb2c1edf0b24476df4) fix(release): clean up next-version base image tags (#6678)
 * [17438790a](https://github.com/kubeovn/kube-ovn/commit/17438790ac8b033e3e1fa94afc10a9a6cf3d0c13) docs: updated CHANGELOG.md (#6677)
 * [f2cf6628c](https://github.com/kubeovn/kube-ovn/commit/f2cf6628c7128d93a8c5245305955a10410fc9a0) test(iptables-vpc-nat-gw): gate custom namespace assertions on v1.17+ (#6676)
 * [25878a137](https://github.com/kubeovn/kube-ovn/commit/25878a137c4bd442a69a366e3b22e3cfc311dedc) fix(controller): keep VM LSP port-group memberships when a sibling pod is alive (#6666)
 * [28a800afd](https://github.com/kubeovn/kube-ovn/commit/28a800afdba7d4187a23e5a46b11a2c853fac316) fix(vpcnatgw): keep eth0 on primary CNI in non-primary CNI mode (#6675)
 * [388b35a70](https://github.com/kubeovn/kube-ovn/commit/388b35a70e5b4b41fa8cd7ed0023dc88ffc6bde9) feat: add ExternalSubnet column to iptables-eip printer (#6672)
 * [f6b64fb55](https://github.com/kubeovn/kube-ovn/commit/f6b64fb55d57130b19631adff8ef6122089ab22a) chore(deps): update aquasecurity/trivy-action action to v0.36.0 (#6673)
 * [4cf412152](https://github.com/kubeovn/kube-ovn/commit/4cf412152767a67d6c2ac535fe2efd53da8437a3) fix(lb): skip ct_flush on session-affinity UDP load balancers (#6671)
 * [5a15ad23d](https://github.com/kubeovn/kube-ovn/commit/5a15ad23d12afba20b4a2d79379864cca512acd8) fix(deps): update module github.com/moby/moby/client to v0.4.1 (#6656)
 * [e53ed45f3](https://github.com/kubeovn/kube-ovn/commit/e53ed45f346a8253de90f42e4a3c712f803e8c56) fix: docker timeout always (#6662)
 * [431926178](https://github.com/kubeovn/kube-ovn/commit/4319261783190b185fd2269da82e28fcf9b56ccc) ci: raise timeout for Iptables VPC NAT Gateway E2E to 45m (#6670)
 * [0b05be737](https://github.com/kubeovn/kube-ovn/commit/0b05be737d15aea3d139e71f8f24170e79cc7ae6) fix(speaker): require NODE_NAME and inject POD_IPS in v2 chart (#6669)
 * [5f11aca17](https://github.com/kubeovn/kube-ovn/commit/5f11aca1727010b9f64db877d519871279088de8) fix(deps): update module kubevirt.io/client-go to v1.8.2 (#6667)
 * [f384eeb8c](https://github.com/kubeovn/kube-ovn/commit/f384eeb8c5d11c9f4141f28ce2093b23202c5d39) fix: add env NODE_NAME in spicker pod (#6664)
 * [5a371fc05](https://github.com/kubeovn/kube-ovn/commit/5a371fc05dbc17ee6be5d9e83df1492c2cb2f1c9) fix(deps): update module kubevirt.io/api to v1.8.2 (#6658)
 * [751cd4763](https://github.com/kubeovn/kube-ovn/commit/751cd47633d220a5f5eb21948f7ffd105629590f) chore(deps): update dependency kubevirt to v1.8.2 (#6657)
 * [f3a7aea05](https://github.com/kubeovn/kube-ovn/commit/f3a7aea050f4bf3749f027f15b1aa1ecaf9861b9) fix(deps): update module github.com/moby/moby/api to v1.54.2 (#6655)
 * [b49190cbb](https://github.com/kubeovn/kube-ovn/commit/b49190cbb9d1ee17a4ce47a76005827b2ca42a3e) feat(qos): support configurable ingress/egress burst via annotation (#6628)
 * [f077e4c2a](https://github.com/kubeovn/kube-ovn/commit/f077e4c2ab1aea840e2f40d4389279d7b957bd0c) fix(deps): update kubernetes (master) (#6636)
 * [c8cd03c4b](https://github.com/kubeovn/kube-ovn/commit/c8cd03c4bb8ca144aedfd3af4aa0c5767c7db9e8) fix(slr): snapshot subnet/VPC in SwitchLBRuleInfo to avoid VIP leak (#6651)
 * [61114ca4f](https://github.com/kubeovn/kube-ovn/commit/61114ca4f37101e9319aa335fd75cf3f5581d1e1) chore(deps): update dependency cilium to v1.19.3 (master) (#6642)
 * [4841caf81](https://github.com/kubeovn/kube-ovn/commit/4841caf81ab2af5514df880eae8abe438f01e8df) chore(deps): bump github.com/moby/spdystream in /test/anp (#6648)
 * [902797d10](https://github.com/kubeovn/kube-ovn/commit/902797d10e89f88c8c6007f89e4f8a2c171baa2f) perf(controller): deep-trim pod informer cache (#6649)
 * [6435556fa](https://github.com/kubeovn/kube-ovn/commit/6435556fa6f977934dcb271fcb73524e5e885d52) fix(controller): stop overwriting VPC status on bootstrap (#6650)
 * [15b02bad3](https://github.com/kubeovn/kube-ovn/commit/15b02bad3ca8bf13d764b921af3a062875b2f4d3) perf(controller): use informer indexers for hot-path lookups (#6646)
 * [72cb59ac3](https://github.com/kubeovn/kube-ovn/commit/72cb59ac3987deb8d1bcec6f53b83dd188a64d76) fix(deps): update module github.com/puzpuzpuz/xsync/v4 to v4.5.0 (#6647)
 * [e1676b167](https://github.com/kubeovn/kube-ovn/commit/e1676b167467e3cf71c12da7ad1ada382a9018d3) fix(controller): serialize centralized subnet reconcile in node update path (#6627)
 * [c9689e42a](https://github.com/kubeovn/kube-ovn/commit/c9689e42abd5f337864b08f0a839890828f42174) chore(deps): bump github.com/moby/spdystream in /test/cnp (#6639)
 * [63cddcb7c](https://github.com/kubeovn/kube-ovn/commit/63cddcb7cbeedbc9dc38df016733f2b69424bfbf) chore(deps): bump github.com/moby/spdystream from 0.5.0 to 0.5.1 (#6638)
 * [2f6e02777](https://github.com/kubeovn/kube-ovn/commit/2f6e02777bca7b3ac390562b25ea6db3b60bc303) fix(deps): update module github.com/containerd/containerd/v2 to v2.2.3 (#6637)
 * [d63019311](https://github.com/kubeovn/kube-ovn/commit/d63019311d89cf9f267e3a6c574848de17abe4a3) fix(helm): fix pinger external-address and add v2 image overrides (#6609)
 * [6e50ec28d](https://github.com/kubeovn/kube-ovn/commit/6e50ec28da3757c1c6785279b85c56e90d45c0a6) feat(chart): ship PrometheusRule and Grafana dashboards in kube-ovn-v2 (#6630)
 * [c46fe6245](https://github.com/kubeovn/kube-ovn/commit/c46fe6245a356fd7052f345277df388153b56222) fix: limit vswitch client schema to core tables (#6614)
 * [2391e544f](https://github.com/kubeovn/kube-ovn/commit/2391e544fc434e831e5e098f57b95ded8a185343) feat(natgw): support set namespace in vpc-nat-gw and iptables-eip CRD (#6595)
 * [a6b77cad5](https://github.com/kubeovn/kube-ovn/commit/a6b77cad58e4edc50ac1554b70237a1bdbe4d4b5) fix(test): rename SlrInfo to SwitchLBRuleInfo in switch_lb_rule_test (#6626)
 * [663aca2f6](https://github.com/kubeovn/kube-ovn/commit/663aca2f603d8bb95afb4ae2bf26f739f46aab18) Rename SlrInfo -> SwitchLBRuleInfo (#6580)
 * [5fc7a333d](https://github.com/kubeovn/kube-ovn/commit/5fc7a333d9f1a17e65f81abf24536408fbdece1d) fix(deps): update module github.com/cert-manager/cert-manager to v1.20.2 (#6624)
 * [6ef769d64](https://github.com/kubeovn/kube-ovn/commit/6ef769d64138041538634202da6c953a6271f665) chore(deps): update dependency cert-manager to v1.20.2 (#6623)
 * [2b503b98b](https://github.com/kubeovn/kube-ovn/commit/2b503b98b4e9ae403100c1a0fd38d4a2929d64d2) fix(e2e): wait for BFD port status to be cleared in VPC egress gateway test (#6622)
 * [85e7f940f](https://github.com/kubeovn/kube-ovn/commit/85e7f940f4dd14186b809455f1e0adb78f435a51) fix(underlay): wait for localnet patch port before IPv6 DAD (#6621)
 * [c642da4da](https://github.com/kubeovn/kube-ovn/commit/c642da4da37ab4d1f46c798473809948fe8180d7) feat(lb): enable ct_flush for UDP load balancers (#6620)
 * [4ca547167](https://github.com/kubeovn/kube-ovn/commit/4ca547167e4a49def4542c55e8f13805380414a7) fix(daemon): sort map keys before iterating to stabilize iptables rules (#6619)
 * [050a62454](https://github.com/kubeovn/kube-ovn/commit/050a624543977b4387121f71cad67d06f44e8e0b) chore: update renovate branches to track release-1.16 (#6616)
 * [67b13536c](https://github.com/kubeovn/kube-ovn/commit/67b13536c5c939f6599a6b59dc0be882e7d82da6) docs: updated CHANGELOG.md (#6617)
 * [cb0c35e46](https://github.com/kubeovn/kube-ovn/commit/cb0c35e4684345465eeff83e26ffd331acf95bda) docs: add infos about v1 to v2 chart migration. (#6423)
 * [da1bae9b8](https://github.com/kubeovn/kube-ovn/commit/da1bae9b84a4f5d155e6f63cff746deaa5323448) chore(users): add kubestellar console to kube-ovn users
 * [1ec2857aa](https://github.com/kubeovn/kube-ovn/commit/1ec2857aa36bd35b1ef0539331178fdba9150128) chart: enable ovn-match-northd-version before ovs-ovn upgrade (#6606)
 * [4ac06fafa](https://github.com/kubeovn/kube-ovn/commit/4ac06fafa2ec3cb93a59a0cdf04e612b654173c8) fix(deps): update golang (#6611)
 * [90cf7d9ac](https://github.com/kubeovn/kube-ovn/commit/90cf7d9ac253ea30ff7775c1ba32ee65131fb535) chore(deps): update dependency helm to v4.1.4 (#6610)
 * [1fd1878a1](https://github.com/kubeovn/kube-ovn/commit/1fd1878a1c488dcdbb37bf491291164acf80330a) docs: updated CHANGELOG.md (#6605)
 * [3f0f96914](https://github.com/kubeovn/kube-ovn/commit/3f0f96914cdaa7d64c00b3d87127e9ea004e5685) update Next.md
 * [0fb2c0f10](https://github.com/kubeovn/kube-ovn/commit/0fb2c0f10d1aa1aa1ca458b63984538cc4280a6c) prepare for next release

### Contributors

 * 09mulcs
 * Aleksei Sviridkin
 * Andrei Kvapil
 * Andy Anderson
 * Apinant.u
 * Artem
 * Chanyeol Yoon
 * Clement Phu
 * Damien Chomat
 * Donghee Park
 * Gaurav Mehta
 * Hadi Abbasi
 * Hamid
 * Hariprasath
 * Harisudarsan
 * Immanuel Tikhonov
 * Jiri Jezek
 * Johann Schley
 * Kiarash
 * Maheswara Reddy Chennuru
 * Max Weiss
 * Mengxin Liu
 * Nick
 * Patryk Rostkowski
 * Paul
 * Pujitha Paladugu
 * Renuka Devi Rajendran
 * Rohan
 * SKALA NETWORKS
 * Vitaly
 * Vitaly Tyuta
 * andriishestakov
 * btray900
 * changluyi
 * clyi
 * dependabot[bot]
 * dmitriy-bty
 * github-actions[bot]
 * jimyag
 * jskrill
 * jxj
 * renovate[bot]
 * sarankumar
 * xieyanker
 * zbb88888
 * zhangzujian
 * zhaoxuan
 * 张祖建

## v1.16.0 (2026-04-09)

 * [7162bc4af](https://github.com/kubeovn/kube-ovn/commit/7162bc4afda1aee29c67addc2badc73b8c7227cf) fix(controller): prevalidate nat gw name and normalize long labels (#6604)
 * [2bbdc5ee8](https://github.com/kubeovn/kube-ovn/commit/2bbdc5ee8bdc031bde7f46596924e79f9402c816) docs: update NEXT.md with Post-v1.15.0 features (#6603)
 * [b95c867c2](https://github.com/kubeovn/kube-ovn/commit/b95c867c292da6a09851e045d0bcc5708bb7a05e) chore(deps): bump go.opentelemetry.io/otel/sdk from 1.40.0 to 1.43.0 (#6602)
 * [798ede99f](https://github.com/kubeovn/kube-ovn/commit/798ede99faad1982db8e36b2dba9d979ceec21fd) fix(deps): update module golang.org/x/sys to v0.43.0 (#6601)
 * [4420b520b](https://github.com/kubeovn/kube-ovn/commit/4420b520b5afc255c09b25ff9113cbb3541d0721) fix(controller): clean up orphaned LBHCs during multi-SLR deletion (#6596)
 * [558c6f874](https://github.com/kubeovn/kube-ovn/commit/558c6f874ebe0365e4d63f1d9afdeb56932cb8f3) fix(daemon): reorder SR-IOV representor lifecycle to prevent CmdAdd/CmdDel race (#6593)
 * [378a42282](https://github.com/kubeovn/kube-ovn/commit/378a42282c8ae7d72d6a2d9d5837a9d3d8a0e40b) fix(e2e): preload agnhost image into KIND nodes before running tests (#6594)
 * [6cf33fa5f](https://github.com/kubeovn/kube-ovn/commit/6cf33fa5f8119bd24d12ff303332e561e131901e) fix: add SplitTrimmed to prevent index out of range and empty token propagation (#6592)
 * [4f96a74a0](https://github.com/kubeovn/kube-ovn/commit/4f96a74a0b64f231ae06d0482a6e355610352e38) feat: support BGP and EVPN for vpc-egress-gateway (#6224)
 * [58fa36b54](https://github.com/kubeovn/kube-ovn/commit/58fa36b549c9e996e93631cea9a67bcb16665be9) docs: updated CHANGELOG.md (#6590)
 * [815b4e1dc](https://github.com/kubeovn/kube-ovn/commit/815b4e1dce0e486e75e2e88f021f69090b8e1924) chore(deps): update dependency go to v1.26.2 (#6586)
 * [2a63959ba](https://github.com/kubeovn/kube-ovn/commit/2a63959ba61b3be66acd9e26a4628819776f4410) fix(controller): preserve VPC policy route next hops during OVN diff (#6577)
 * [ee61003d5](https://github.com/kubeovn/kube-ovn/commit/ee61003d5ee8ead1893af8b6f447efc1ced52184) fix: add nil checks for net.ParseIP/ParseCIDR at input boundaries to prevent panic (#6581)
 * [d336c9801](https://github.com/kubeovn/kube-ovn/commit/d336c9801a884b417a182bfb5da441b3a9865946) feat(crd):auto check and generate CRDs in dev workflow (#6569)
 * [4d9e8e60d](https://github.com/kubeovn/kube-ovn/commit/4d9e8e60dfc0b4b6d7619975b673be143c4c643d) fix(deps): update module kernel.org/pub/linux/libs/security/libcap/cap to v1.2.78 (#6576)
 * [82be2d108](https://github.com/kubeovn/kube-ovn/commit/82be2d108008d8e3d64b55161d2a51d20efc4104) Update contact link for Kube-OVN Community
 * [075f3860f](https://github.com/kubeovn/kube-ovn/commit/075f3860f8ae584e59b3ab0f86818ad689e87b52) change slack channel
 * [02a0bbffb](https://github.com/kubeovn/kube-ovn/commit/02a0bbffbb0f6ce96b1c5a5eb2cfcb868eb88368) fix(chart-v2): fix extraEnvs indent in controller-deployment template (#6575)
 * [b2f395135](https://github.com/kubeovn/kube-ovn/commit/b2f395135db828237f394e469e9f22a035ef62ef) fix(deps): update module go.podman.io/image/v5 to v5.39.2 (#6573)
 * [814d525a2](https://github.com/kubeovn/kube-ovn/commit/814d525a284c2aa83208e1c983d4c2ba83378a9f) fix(deps): update module github.com/moby/moby/client to v0.4.0 (#6572)
 * [59cb1f51b](https://github.com/kubeovn/kube-ovn/commit/59cb1f51becff9af90847e04ad964837fcd9e965) Update USERS.md (#6562)
 * [c8f4dc901](https://github.com/kubeovn/kube-ovn/commit/c8f4dc901192b398f189d5d41ab06ae1538e159d) auto gen crd (#6548)
 * [eaaf7c6a1](https://github.com/kubeovn/kube-ovn/commit/eaaf7c6a1a24bf0e5f05ad7a1f042f01e4bc1fc6) fix(e2e): wait for second provider network bridge in underlay AfterEach (#6567)
 * [f06903105](https://github.com/kubeovn/kube-ovn/commit/f06903105f4285d990b273b207bfb2fc8d35ce86) fix(deps): update module google.golang.org/grpc to v1.80.0 (#6566)
 * [ecba16115](https://github.com/kubeovn/kube-ovn/commit/ecba1611540eb2f07aa01c912586fb6777bad02b) fix(deps): update module github.com/osrg/gobgp/v4 to v4.4.0 (#6565)
 * [4c122e416](https://github.com/kubeovn/kube-ovn/commit/4c122e4169581051fff459ed345a3203f4bd9232) fix(slr): handle orphaned ip_port_mappings and overwrites (#6558)
 * [779b0ec7c](https://github.com/kubeovn/kube-ovn/commit/779b0ec7c6bdd7a8fc17721690a4f31510f1b8b9) chore(deps): update dependency gobgp to v4.4.0 (#6564)
 * [178f22832](https://github.com/kubeovn/kube-ovn/commit/178f22832bb8eaa1086549b8efa064df8fcfb102) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.2.2 (#6561)
 * [133a2f72b](https://github.com/kubeovn/kube-ovn/commit/133a2f72b339043bb5a87475b1ff1b1b791e4e22) fix(daemon): filter link-local and loopback addresses in loopEncapIPCheck (#6559)
 * [67790efda](https://github.com/kubeovn/kube-ovn/commit/67790efdab989929831e80997c796a6f37f24f17) fix(daemon): match pod IPs to EGW by protocol instead of index position
 * [fd5f3ced1](https://github.com/kubeovn/kube-ovn/commit/fd5f3ced16b10afff1dd9c27f643f41cf8d408ac) fix(daemon): add missing NatOutgoing check in getLocalPodIPsNeedPR and getSubnetsNeedPR
 * [a082bb6c0](https://github.com/kubeovn/kube-ovn/commit/a082bb6c03b35f96b7fadf7ec55f99fcd4c51f0f) fix(daemon): prevent panic on ips index out of bounds in pod policy routing
 * [2378f5000](https://github.com/kubeovn/kube-ovn/commit/2378f50002836d84c72a7a507a1f5beb3628bbce) fix(daemon): fix rulespec in cleanObsoleteIptablesRules skipping "-A FORWARD" prefix
 * [ab1c2cdb1](https://github.com/kubeovn/kube-ovn/commit/ab1c2cdb194bc9600dd7e37b73d5bbf2bdd91f39) fix(daemon): add missing return after OVS bridge creation failure in setExGateway
 * [271b91ffa](https://github.com/kubeovn/kube-ovn/commit/271b91ffa31e26dd278131db3a0d4a17ab65482d) fix(daemon): remove dead protocol assignments in generateNatOutgoingPolicyChainRules
 * [6d2bceada](https://github.com/kubeovn/kube-ovn/commit/6d2bceadaaa79ccc330a458c0e0c3a7db40f6aae) test(e2e): add skip version for kubevirt NAD change test (#6556)
 * [848e2b0de](https://github.com/kubeovn/kube-ovn/commit/848e2b0de8bd5f035c0baee8db2a67c675e0d780) fix(deps): update module kubevirt.io/client-go to v1.8.1 (#6551)
 * [dcdb1afd0](https://github.com/kubeovn/kube-ovn/commit/dcdb1afd0b73356acf440ebcdf4c9ad3c413f1fc) chore(deps): update dependency kubevirt to v1.8.1 (#6552)
 * [42dae6722](https://github.com/kubeovn/kube-ovn/commit/42dae67221d4a0f2d0b8572f785a3f62782aa638) fix(deps): update module kubevirt.io/api to v1.8.1 (#6550)
 * [f1e071be2](https://github.com/kubeovn/kube-ovn/commit/f1e071be26601368beebbaf05b44d86dcf05fbf8) fix(daemon): harden protocol iteration and bounds checks in gateway (#6547)
 * [eca75a5d3](https://github.com/kubeovn/kube-ovn/commit/eca75a5d3049639b0d520f68b38b8e02d2ba37ca) fix(daemon): add TrimManagedFields to all daemon informer factories (#6546)
 * [ecf66bfb4](https://github.com/kubeovn/kube-ovn/commit/ecf66bfb46cc9a0f99bed77a7c8190ab2fdc83fd) fix(daemon): replace bulk netlink dumps with targeted queries (#6540)
 * [e07a2d2ae](https://github.com/kubeovn/kube-ovn/commit/e07a2d2aebc166730f4e360ff9e8aaa6ec814203) ci: add least-privilege permissions to GitHub Actions workflows (#6544)
 * [ac3666ef2](https://github.com/kubeovn/kube-ovn/commit/ac3666ef29f07a4bd377c892f512a5497bc828ea) fix(controller): handle missing logical switch in handleMcastQuerierChange (#6542)
 * [4146075dc](https://github.com/kubeovn/kube-ovn/commit/4146075dc8c2ad5b25d381aa3f172360fca1056b) fix(ci): tolerate pod disappearing during cleanup log collection (#6543)
 * [2f30f501b](https://github.com/kubeovn/kube-ovn/commit/2f30f501b33f57b23da170797de6cba7cbfbb51a) fix(deps): update module github.com/cert-manager/cert-manager to v1.20.1 (#6535)
 * [d208d2f9c](https://github.com/kubeovn/kube-ovn/commit/d208d2f9cdf2d0d0f0598d51ae47d44631526528) fix(e2e): use diff-based merge patch for DaemonSet instead of full object (#6539)
 * [99d183108](https://github.com/kubeovn/kube-ovn/commit/99d1831084aedc381e00b335752b4ab9cede11b9) fix: wrap EIP error shadowing and use uint64 for iptables counters (#6538)
 * [16456fe68](https://github.com/kubeovn/kube-ovn/commit/16456fe6831d3b7bd56417291fcea52b112fc121) fix(informer): return nil instead of error in VM indexer to prevent client-go panic (#6537)
 * [b54a61501](https://github.com/kubeovn/kube-ovn/commit/b54a615018f7455687f86ce20afc52d39edfbd88) fix(controller): prevent invalid ACL generation for named ports and IPBlock except CIDRs (#6536)
 * [ed67bac5f](https://github.com/kubeovn/kube-ovn/commit/ed67bac5feb9cb9f95bd85fbd6c30c465a7af3d3) chore(deps): update dependency cert-manager to v1.20.1 (#6534)
 * [09023e77a](https://github.com/kubeovn/kube-ovn/commit/09023e77a20b20e8a1812a3610757f11fd965830) fix(e2e): pin and preload kubevirt containerdisk image in KIND (#6533)
 * [b2d57d79e](https://github.com/kubeovn/kube-ovn/commit/b2d57d79e9ffefd8769a6e4c598fb4326f6282ff) feat(controller): support per-pod NIC granular DHCP control via annotations (#6475)
 * [023250be8](https://github.com/kubeovn/kube-ovn/commit/023250be8e03e99567284d82306fcecc5999710a) fix(controller): add finish logs for all gc functions (#6532)
 * [0074ae34f](https://github.com/kubeovn/kube-ovn/commit/0074ae34f0dd2206b5397a7a3e0c9070a0785887) fix(controller): defer subnet processing until vlan is fully ready (#6527)
 * [edd3de227](https://github.com/kubeovn/kube-ovn/commit/edd3de227d746afafb272ddbc978cd27fa9d3c47) update dependency (#6526)
 * [d3070d139](https://github.com/kubeovn/kube-ovn/commit/d3070d139b48e73c9a36f48dbdc3c898ba54d4c2) fix(controller): clean orphaned attachment IPs on KubeVirt NAD hotplug (#6519)
 * [0e0fa7d1b](https://github.com/kubeovn/kube-ovn/commit/0e0fa7d1b9be955e69af7c5da788188905ef8a3e) docs: updated CHANGELOG.md (#6528)
 * [ca9491418](https://github.com/kubeovn/kube-ovn/commit/ca9491418c3cb11c96dd6cbe0d137cc38356dcd2) fix: check northd endpoint by pod IP protocol in ovn leader checker (#6523)
 * [19546786e](https://github.com/kubeovn/kube-ovn/commit/19546786e25e3800d49932e1b6f2c79f5c59c460) fix(daemon): clean VLAN subinterfaces and label when excluding node from provider network (#6521)
 * [932fcccef](https://github.com/kubeovn/kube-ovn/commit/932fcccef211fc4e3e8130fa877d7489e7458350) fix(daemon): skip invalid policy routing rules when pod lacks matching protocol IP (#6522)
 * [2c28bfae3](https://github.com/kubeovn/kube-ovn/commit/2c28bfae39b289b28309c50d18c08bd385fedacb) chore(deps): upgrade k8s dependencies to v1.35.3 (#6518)
 * [b129a0b9e](https://github.com/kubeovn/kube-ovn/commit/b129a0b9e9d772e76460b3c0cc0d6b31d3609d0e) fix(deps): update module kubevirt.io/api to v1.8.0 (#6516)
 * [c9ee30249](https://github.com/kubeovn/kube-ovn/commit/c9ee30249ea90261eb100b225d87296f37f3e473) fix(deps): update module kubevirt.io/client-go to v1.8.0 (#6517)
 * [289868e76](https://github.com/kubeovn/kube-ovn/commit/289868e76fc6a79aa7557d8de4dcee9590489f90) feat(controller): use BigInt for SubnetStatus IP count fields (#6403)
 * [1ce02e4f2](https://github.com/kubeovn/kube-ovn/commit/1ce02e4f2383662c2c55d4763e32911f895dd112) fix(controller): use defer to protect vpcKeyMutex unlock in subnet handler (#6515)
 * [b0cf83470](https://github.com/kubeovn/kube-ovn/commit/b0cf83470c762d2e0cc2f7be4e5e88800e00d02c) fix(netpol): use parenthesis around ipblock expressions
 * [62c79dd07](https://github.com/kubeovn/kube-ovn/commit/62c79dd077d63be6193fdadbfedefa7b74515702) chore(ci): add test for netpol regression #6508
 * [f701ecc31](https://github.com/kubeovn/kube-ovn/commit/f701ecc311c825060dae7684baffe45f75665ffb) fix(daemon): use multiple KubeVirt labels for VM launcher pod lookup (#6512)
 * [8957bac67](https://github.com/kubeovn/kube-ovn/commit/8957bac677396cda441cf427090ab5ca265701ec) perf(daemon): avoid full netlink dump in getIfaceByIP (#6513)
 * [ed82f2201](https://github.com/kubeovn/kube-ovn/commit/ed82f22011f8f2cc287b765cb5cb1484c342fb43) chore(deps): update azure/setup-helm action to v5 (#6511)
 * [1221eec1f](https://github.com/kubeovn/kube-ovn/commit/1221eec1fdfafc1d5a90c2270ec2470c75fbd7c7) chore(charts): update readme and values on v2 chart
 * [119adbaa9](https://github.com/kubeovn/kube-ovn/commit/119adbaa9e653560dca41893a647458376d18768) fix(controller): filter empty VIP strings from SwitchLBRule annotation (#6507)
 * [3f2d212a5](https://github.com/kubeovn/kube-ovn/commit/3f2d212a54fe226f926bd31729cd680bd38eaa03) fix(ci): restart coredns pods after kube-ovn install in Talos test (#6504)
 * [fbf5c3d38](https://github.com/kubeovn/kube-ovn/commit/fbf5c3d3832c9a74b9872fe8781821482565fb44) ci: only fail on pod restarts caused by panic or SIGSEGV (#6506)
 * [c84bc03d7](https://github.com/kubeovn/kube-ovn/commit/c84bc03d724ae79e1ea236b5c1d1050ed14b2478) docs: updated CHANGELOG.md (#6505)
 * [0f2e7d6d3](https://github.com/kubeovn/kube-ovn/commit/0f2e7d6d3b07ee46b448fbcfe130eb22f0170763) fix(controller): remove dead code in retryDelDupChassis (#6503)
 * [456037cfa](https://github.com/kubeovn/kube-ovn/commit/456037cfa5822de54174d15cc631d7c28120cfa4) fix(controller): clone ExternalIDs map and skip empty pgName in node handling (#6502)
 * [8596ca271](https://github.com/kubeovn/kube-ovn/commit/8596ca271654b691558323a1ae7b7b9cbff36d70) chore(deps): update dependency golangci-lint to v2.11.4 (#6500)
 * [5005a4d82](https://github.com/kubeovn/kube-ovn/commit/5005a4d8211d63a026efc607fa2b1d0f316d18d7) docs: updated CHANGELOG.md (#6501)
 * [c84be7509](https://github.com/kubeovn/kube-ovn/commit/c84be75097719abb77a10f1dc4a589ab3251bd2d) fix(controller): fix variable shadowing in cleanDuplicatedChassis (#6496)
 * [e7451a9dc](https://github.com/kubeovn/kube-ovn/commit/e7451a9dc8e82a695bda0962b8d9cd9379c55743) fix(e2e): increase VM StartSync ready timeout from 2m to 3m (#6499)
 * [f1a627c73](https://github.com/kubeovn/kube-ovn/commit/f1a627c737f0414792a2e79ef8dba697b13c44a1) fix(controller): prevent dual-stack cross-deletion of local DNS cache policy routes (#6497)
 * [2eac83918](https://github.com/kubeovn/kube-ovn/commit/2eac83918831be92f67b0363174319952b30707a) fix(controller): handle npsLister.List error in checkAndUpdateNodePortGroup (#6498)
 * [9d8f3367c](https://github.com/kubeovn/kube-ovn/commit/9d8f3367cd70c7010e33961ce562e4511c94c104) fix(e2e): add retry to pre-corruption DB health check in HA test (#6495)
 * [40e4f5783](https://github.com/kubeovn/kube-ovn/commit/40e4f5783ccb453004c2f592d5affa0c7fb7f027) fix(controller): use exact match instead of substring match for VIP lookup in syncVirtualPort (#6491)
 * [9e707d7dd](https://github.com/kubeovn/kube-ovn/commit/9e707d7ddb5885d8bb6ad730594055e3d88f0f4b) refactor(controller): remove legacy svcMatchNetworkPolicies code (#6494)
 * [58f08b5a7](https://github.com/kubeovn/kube-ovn/commit/58f08b5a708c952ebb9b59d13a461a801c7e3dc0) fix(ipam): release allocated static IPs on partial allocation failure (#6493)
 * [fd8ea4143](https://github.com/kubeovn/kube-ovn/commit/fd8ea41433929f0b24bc8361a439690365e399d0) fix(controller): use consistent dual-stack check in subnet finalizer removal (#6492)
 * [745fb7f70](https://github.com/kubeovn/kube-ovn/commit/745fb7f7072453aca73522699483d9e915428b14) fix(controller): fix nil deref, cache mutation, and early return bugs in subnet controller (#6490)
 * [e8ba7b0f7](https://github.com/kubeovn/kube-ovn/commit/e8ba7b0f743975692490a92589c89a6e25242b93) fix(controller): prevent checkSubnetConflict from swallowing conflict errors (#6489)
 * [f6e1fb119](https://github.com/kubeovn/kube-ovn/commit/f6e1fb119c13a0f6d2ab0b1341413dde4e90fd17) fix(controller): handle restartable init containers in DeleteNamedPortByPod (#6486)
 * [e7df610cc](https://github.com/kubeovn/kube-ovn/commit/e7df610cc1283d0dc699798c9db1204b5e6544e8) refactor(controller): eliminate duplicate OVN/K8s queries in pod controller (#6485)
 * [9cc36dc88](https://github.com/kubeovn/kube-ovn/commit/9cc36dc886b11d60ff6fe1c381681f4e7b520f06) fix(controller): resolve variable shadowing causing IPAM key mismatch in per-interface static IP allocation (#6484)
 * [3c8fb0e18](https://github.com/kubeovn/kube-ovn/commit/3c8fb0e188cf71c5bed0689c21126d7e238fd4ff) fix(controller): prevent panics from nil deref and data races in pod controller (#6483)
 * [6035b84d5](https://github.com/kubeovn/kube-ovn/commit/6035b84d5e62924877bf7b40537f7ab9b6e2a79e) chore(deps): update dependency talos to v1.12.6 (#6481)
 * [924434330](https://github.com/kubeovn/kube-ovn/commit/9244343309d625072f17651626a6404996728d62) fix(daemon): create listen sockets before starting controller to avoid race (#6478)
 * [5a0009ad6](https://github.com/kubeovn/kube-ovn/commit/5a0009ad658b2cbffb98bb341c3f6159269d5064) fix(e2e): wait for OVS bridge cleanup before Docker disconnect in vlan subinterface tests (#6480)
 * [436d4c5b8](https://github.com/kubeovn/kube-ovn/commit/436d4c5b8f2c5515a0df440ae2da5c32a62c5c31) fix(helm): match control-plane nodes by label presence when value is empty (#6479)
 * [5d36c6068](https://github.com/kubeovn/kube-ovn/commit/5d36c606806afc767ddde4291e2758cf4ac470be) fix(deps): update module google.golang.org/grpc to v1.79.3 [security] (#6477)
 * [f584eb59b](https://github.com/kubeovn/kube-ovn/commit/f584eb59b8ab3a8cb7a28be02a30db96e68717f7) fix(daemon): sort encap IP list for deterministic ovs-vsctl output (#6473)
 * [610106b20](https://github.com/kubeovn/kube-ovn/commit/610106b2037f22adfb3337f8cc85981eb1bc2638) fix(controller): reuse addPolicyRouteToVpc to handle empty NextHopIP correctly (#6472)
 * [ecf61525c](https://github.com/kubeovn/kube-ovn/commit/ecf61525c1f3e9d8e0029df88201ee13e277f1a0) docs: updated CHANGELOG.md (#6471)
 * [a9bbb2180](https://github.com/kubeovn/kube-ovn/commit/a9bbb21801da17be5e48321e6b1fe10f2bd67ff3) fix(daemon): race condition between setExGateway and configProviderNic (#6469)
 * [beaf89ee5](https://github.com/kubeovn/kube-ovn/commit/beaf89ee5b4a848c90b15367b8fb1ccf2885abae) fix(daemon): handle tombstone objects in subnet and service delete handlers (#6470)
 * [42c450992](https://github.com/kubeovn/kube-ovn/commit/42c45099278eb01dc9a1bccafc0d9a85c86c9286) fix(deps): update module kubevirt.io/api to v1.7.2 (#6467)
 * [b2f2738b6](https://github.com/kubeovn/kube-ovn/commit/b2f2738b6343219df1f7ba664cd5b70b8b45ab63) chore(deps): update dependency kubevirt to v1.7.2 (#6466)
 * [f8e1f5830](https://github.com/kubeovn/kube-ovn/commit/f8e1f5830a46c107f87c184b29afdc3d135f2c76) Fix misleading log message on VPC Policy Route addition. (#6464)
 * [6bfb85848](https://github.com/kubeovn/kube-ovn/commit/6bfb858481d59c8e9b8da802a6ddc8c0ee8312a1) dist: add 25.03 ACL compatibility for OVN patch (#6463)
 * [7709e240c](https://github.com/kubeovn/kube-ovn/commit/7709e240cfd2acccccbf592351a83d5e8df27234) fix: FIP/DNAT/SNAT spec update handlers (#6412)
 * [009e0c5a0](https://github.com/kubeovn/kube-ovn/commit/009e0c5a09905ae62390acace63d885d205280b4) fix(controller): remove Service ClusterIP from NetworkPolicy address sets (#6461)
 * [a3b00701a](https://github.com/kubeovn/kube-ovn/commit/a3b00701a09efae4e282b7710fc447d12c79e211) fix(controller): backfill nat gateway lan ip safely (#6424)
 * [5b43fb24c](https://github.com/kubeovn/kube-ovn/commit/5b43fb24cef0b131c1722f535e20da41eb8928db) kind: enable OwnerReferencesPermissionEnforcement (#6456)
 * [a80b1185d](https://github.com/kubeovn/kube-ovn/commit/a80b1185d070b5ff9bca75a431dfa8d1e3f9c3d9) fix: restore NIC name during provider network cleanup in exchangeLinkName mode (#6455)
 * [f364d3302](https://github.com/kubeovn/kube-ovn/commit/f364d33025cee0ad707152448dfe209dc14af891) fix(deps): update module github.com/containernetworking/plugins to v1.9.1 (#6459)
 * [a494ab32e](https://github.com/kubeovn/kube-ovn/commit/a494ab32e8a41342293c86a96a72d26823d43d52) chore(deps): update dependency cni-plugin to v1.9.1 (#6458)
 * [2d641b9ad](https://github.com/kubeovn/kube-ovn/commit/2d641b9ad6721853e68c8e7e7302ae95c6ca34ac) fix: avoid error log for missing VPC annotation on underlay pods (#6454)
 * [169b14ac7](https://github.com/kubeovn/kube-ovn/commit/169b14ac7ce119f821c8a61b146ce7f36c89e569) fix: persist node IP annotation to API server during IPAM init (#6453)
 * [636f5166c](https://github.com/kubeovn/kube-ovn/commit/636f5166c77b7882f017ed58f2b0c4b37db5804f) Improve hairpin logic to support FIP, SNAT from any CIDR within the VPC and LoadBalancers (#6445)
 * [15c4ece9d](https://github.com/kubeovn/kube-ovn/commit/15c4ece9d5b1f8c36b8fe91561184e32ffef2966) fix: add webhook self-exclusion to prevent recovery deadlock (#6452)
 * [9c4837026](https://github.com/kubeovn/kube-ovn/commit/9c4837026264c26f19f7de0018d3c8300220fab1) fix: harden provider network cleanup to prevent orphaned bridges and routes (#6448)
 * [dc5ad3b07](https://github.com/kubeovn/kube-ovn/commit/dc5ad3b078865087d77e875f7015e7363d9ab20b) Split ECMP next-hops in custom VPC policy routes (#6405)
 * [b6014c0e1](https://github.com/kubeovn/kube-ovn/commit/b6014c0e12b053382ba0c54a368506cf936155f8) chore(deps): update dorny/paths-filter action to v4 (#6446)
 * [7cd20ce63](https://github.com/kubeovn/kube-ovn/commit/7cd20ce63a36da433c4c20edc0f523f785e1e122) kubectl-ko: fix trace failure when arping returns multiple MAC matches (#6444)
 * [2cf8fa40d](https://github.com/kubeovn/kube-ovn/commit/2cf8fa40d89d6a1554432c2b3428a07e40b95934) fix(deps): update module golang.org/x/tools to v0.43.0 (#6442)
 * [68dc63fa4](https://github.com/kubeovn/kube-ovn/commit/68dc63fa4b97ab00843e4b8c5a6a5c0ff783e9e2) chore(deps): update dependency submariner to v0.23.1 (#6441)
 * [caa98a5f8](https://github.com/kubeovn/kube-ovn/commit/caa98a5f8b04049483176a8db2f2562f952cc47d) chore(deps): update dependency helm to v4.1.3 (#6440)
 * [1acf8ee62](https://github.com/kubeovn/kube-ovn/commit/1acf8ee6236087cf1a507d5fdc6f95e32c783034) docs: updated CHANGELOG.md (#6438)
 * [a033c7c8a](https://github.com/kubeovn/kube-ovn/commit/a033c7c8acfa4da77e336aa9282c0145bdd94bd6) docs: updated CHANGELOG.md (#6437)
 * [4960f4c00](https://github.com/kubeovn/kube-ovn/commit/4960f4c001e60ee0c27806b9bee5c1448a95f833) docs: updated CHANGELOG.md (#6436)
 * [86692e6f2](https://github.com/kubeovn/kube-ovn/commit/86692e6f21009153901fec6b834156283caa5f85) feat: Extend SG API to have tiers, larger priority range, localAddress and port matches (#6330)
 * [078cbecd2](https://github.com/kubeovn/kube-ovn/commit/078cbecd205128dd9335cc58eae97dd76bf39ab6) controller: do not set blockOwnerDeletion for ip resources (#6435)
 * [838c18b2c](https://github.com/kubeovn/kube-ovn/commit/838c18b2c813914e8689f2e57a1a1a05f1c4e3bf) allow any eip share external subnet gw in one (native) vlan (#6426)
 * [5fd1241e8](https://github.com/kubeovn/kube-ovn/commit/5fd1241e8958d4fc38dc035d164b0d9fcefb9877) fix(ut): use Eventually for ovsdb cache queries in chassis tests (#6431)
 * [f2f3f632b](https://github.com/kubeovn/kube-ovn/commit/f2f3f632b3c9a1ca0bee6834fe7b4f39051f3f22) fix(deps): update golang (#6427)
 * [0e10d55fb](https://github.com/kubeovn/kube-ovn/commit/0e10d55fb13cbd13100378ef1b7c3e0ece45f55f) chart: remove unused values (#6418)
 * [dad7f3c62](https://github.com/kubeovn/kube-ovn/commit/dad7f3c6201775f4446e499c9403c7ca338491f6) fix: add ipv6_src to LB selection_fields for IPv6 session affinity (#6422)
 * [b391b1c10](https://github.com/kubeovn/kube-ovn/commit/b391b1c1023562b35b2785ca902179406618ddfe) fix(e2e): wait for namespace IPPool annotation before creating workloads to prevent race condition (#6421)
 * [01e5d9742](https://github.com/kubeovn/kube-ovn/commit/01e5d97427a08cfc57dc29cd05f95d5a4b420f5f) fix(helm): make v2 helm chart run upgrade logic within ovn (#6410)
 * [23626e7d0](https://github.com/kubeovn/kube-ovn/commit/23626e7d0837e9bab41822e64d3435db0aa07706) fix(deps): update module github.com/cert-manager/cert-manager to v1.20.0 (master) (#6416)
 * [13f996011](https://github.com/kubeovn/kube-ovn/commit/13f9960119e579e1fcf666ce19033b3d9dbda830) fix: scope ipBlock.except to its own CIDR in NetworkPolicy ACL rules (#6419)
 * [f96b05c5a](https://github.com/kubeovn/kube-ovn/commit/f96b05c5af1aef96ec010775b13075d7561e5282) fix: reconcile subnet and localnet network_name on vlan provider change (#6393)
 * [d5a9cae85](https://github.com/kubeovn/kube-ovn/commit/d5a9cae85b62b5675bb28a375839146b714f3698) chore(deps): update dependency talos to v1.12.5 (#6411)
 * [d57584d92](https://github.com/kubeovn/kube-ovn/commit/d57584d9242f565743df69b182901246b7cfedf1) fix(deps): update module github.com/containerd/containerd/v2 to v2.2.2 (#6417)
 * [6d6c3e0c6](https://github.com/kubeovn/kube-ovn/commit/6d6c3e0c66948dde280900f6746baa38d3891cd8) chore(deps): update dependency cert-manager to v1.20.0 (#6415)
 * [f364d9b8b](https://github.com/kubeovn/kube-ovn/commit/f364d9b8bd5b604fb28d9e3ccae4862bfaa2ca63) feat: support custom resources and bandwidth limits for VpcEgressGateway (#6407)
 * [a2010033f](https://github.com/kubeovn/kube-ovn/commit/a2010033f2e91bcefa6276d76da8f22fee8667c3) fix(e2e): preload pause image into kind cluster to prevent flaky test timeouts (#6414)
 * [f7995d2f1](https://github.com/kubeovn/kube-ovn/commit/f7995d2f19ce009eccc1c5a84e4bf3a9762169c1) fix(e2e): prevent OVS datapath race condition in underlay tests (#6413)
 * [07f73a489](https://github.com/kubeovn/kube-ovn/commit/07f73a48995cb5392401427b176076cd60a4acc3) chore(deps): update dependency golangci-lint to v2.11.3 (master) (#6402)
 * [78b2f742e](https://github.com/kubeovn/kube-ovn/commit/78b2f742e610be7fbacf73666035db5a36746377) fix: scope SLR deletion to its own VPC to avoid deleting VIPs across VPCs (#6391)
 * [6284594c4](https://github.com/kubeovn/kube-ovn/commit/6284594c4876fe8cd72a3638ca548bd978102f97) fix(helm): allow explicit override of ovs-ovn DaemonSet update strategy (#6406)
 * [0fa99c5f3](https://github.com/kubeovn/kube-ovn/commit/0fa99c5f3b9ffd13a01b3ea360d57aebc8b82e9b) chore(deps): update docker/setup-buildx-action action to v4 (#6396)
 * [d8544f97f](https://github.com/kubeovn/kube-ovn/commit/d8544f97fc6bd88afe347c639b67939584782608) chore(deps): update aquasecurity/trivy-action action to v0.35.0 (#6404)
 * [30a143e68](https://github.com/kubeovn/kube-ovn/commit/30a143e6858347da223915ab9e31d0807108f344) fix(deps): update module google.golang.org/grpc to v1.79.2 (#6400)
 * [9a968ef77](https://github.com/kubeovn/kube-ovn/commit/9a968ef773d331b98379224a2d9413edbbfa1899) fix(deps): update golang (#6409)
 * [cef09fee8](https://github.com/kubeovn/kube-ovn/commit/cef09fee81630fc6e29d1b2ddbdac190c9a37e8a) chore(deps): update dependency go to v1.26.1 (#6399)
 * [c86151b0d](https://github.com/kubeovn/kube-ovn/commit/c86151b0dafe424b1a17134b3710d15a29d87a87) fix(deps): update module github.com/moby/moby/client to v0.3.0 (#6395)
 * [5d4cc325f](https://github.com/kubeovn/kube-ovn/commit/5d4cc325fcd5f1d30cedfe11af213ffb321ce0e4) fix(deps): update module github.com/moby/moby/api to v1.54.0 (#6394)
 * [bfd45ab3a](https://github.com/kubeovn/kube-ovn/commit/bfd45ab3a3eac2f0029faf9d97ce296384a1ee6e) ci: wait for metallb-memberlist secret before speaker rollout (#6389)
 * [0567c9a62](https://github.com/kubeovn/kube-ovn/commit/0567c9a62d81d2c013bfd7086d9835c4bc65dafc) fix: add retry logic to checkLeader() and checkDaemonSet() in kubectl-ko (#6387)
 * [e19238eec](https://github.com/kubeovn/kube-ovn/commit/e19238eecea9c5d5f2f872b0dd52d645e79cdbf5) fix: kubevirt e2e subnet delete timeout and IP CRD label management (#6388)
 * [3986ad68d](https://github.com/kubeovn/kube-ovn/commit/3986ad68d038eb9017e432ea8ce2fdcaad6a0bac) fix(daemon): start IPSec service after certificates are configured in OVSDB (#6383)
 * [6e8dda15d](https://github.com/kubeovn/kube-ovn/commit/6e8dda15d65e8a6d2ff2d33a8b9b27e009e32a17) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.14.1 (#6386)
 * [8f155a3ec](https://github.com/kubeovn/kube-ovn/commit/8f155a3eca4e66d1adb8c0cb67e9aedb87603b09) feat:  enable ip forward should skip fastpath (#6384)
 * [50ac6c89c](https://github.com/kubeovn/kube-ovn/commit/50ac6c89c075a098bcd5d1f0a69df9b4fa2acd3e) ci: disable valgrind for metallb e2e and increase speaker rollout timeout (#6385)
 * [e5ab5a041](https://github.com/kubeovn/kube-ovn/commit/e5ab5a041e9cc25dc2059ae341cd562b23049478) fix(e2e): wait for OpenFlow rules and add curl retry in metallb underlay test (#6382)
 * [48aca4072](https://github.com/kubeovn/kube-ovn/commit/48aca40722e18e971a65c028706cfecddd3a4904) fix: protect ANP/BNP/CNP priority maps from concurrent access (#6373)
 * [74eeefa11](https://github.com/kubeovn/kube-ovn/commit/74eeefa11fc31d35071acaa9adcb80d98bac6da2) chore(deps): update aquasecurity/trivy-action action to v0.34.2 (#6380)
 * [6d6c5580c](https://github.com/kubeovn/kube-ovn/commit/6d6c5580c338539904be064c2969f586346107f4) fix: BGP speaker reads IPs from pod annotations instead of pod status (#6344)
 * [44f2e62e9](https://github.com/kubeovn/kube-ovn/commit/44f2e62e9c0e7fc359f6c4a83b2b37a085d0e0ac) fix(deps): update module github.com/osrg/gobgp/v4 to v4.3.0 (#6372)
 * [e6a657c8e](https://github.com/kubeovn/kube-ovn/commit/e6a657c8e7cc2f1750bd0d8c78633ebdc0445b70) fix(daemon): handle OVS bridge creation race condition with exchangeLinkName (#6374)
 * [503c357e9](https://github.com/kubeovn/kube-ovn/commit/503c357e91a2a1f4340e5d427c2c653470548670) test(e2e): increase VPC egress gateway curl timeout to 2 minutes (#6379)
 * [3543938b8](https://github.com/kubeovn/kube-ovn/commit/3543938b86c25c9ee07d4892e561cc1846d223e6) fix(e2e): stabilize ANP domain E2E tests and fix hasDomainNames override bug (#6376)
 * [75b79117c](https://github.com/kubeovn/kube-ovn/commit/75b79117c6638c9c2765b3bda2ca8de09aca26a0) allow from snat eip to fip eip (#6329)
 * [4293cbc39](https://github.com/kubeovn/kube-ovn/commit/4293cbc39211cb22476b6e71e2497ebe6451ed8c) docs: updated CHANGELOG.md (#6375)
 * [e5de0b374](https://github.com/kubeovn/kube-ovn/commit/e5de0b374fcc3fb174f5f190ef9ad3b241840faf) chore(deps): update dependency gobgp to v4.3.0 (#6371)
 * [3168803a5](https://github.com/kubeovn/kube-ovn/commit/3168803a56dc656e7306b2461db5e5ee0f496511) test(e2e): optimize test execution time by reducing unnecessary waits (#6369)
 * [605d61574](https://github.com/kubeovn/kube-ovn/commit/605d61574c3cd2a9da071ecff387c749e2e219ae) ci: wait for control-plane nodes before labeling in chart install (#6370)
 * [ede1d7123](https://github.com/kubeovn/kube-ovn/commit/ede1d7123afae810f3aeb473d5bed5f7d33c021b) fix(e2e): reuse OVSDB client in address set tests to avoid HA deadlock (#6368)
 * [43ff19a89](https://github.com/kubeovn/kube-ovn/commit/43ff19a8941f155ac4b391ca6deb6ec9b68b5662) ci: fix building base images
 * [e9b65ce70](https://github.com/kubeovn/kube-ovn/commit/e9b65ce70379da386528a7704e91a80841bc6fe5) fix: prevent subnet from getting permanently stuck when VLAN is not ready (#6352)
 * [b0d58bbb6](https://github.com/kubeovn/kube-ovn/commit/b0d58bbb64517b89e28930c97a7bd832e79f0435) cni-server: fix vlan id reference (#6365)
 * [e229c02ad](https://github.com/kubeovn/kube-ovn/commit/e229c02adee799ad54f5697bf5b4893101ed4499) fix(ha): detect and recover from split-brain leader with incomplete cluster membership (#6363)
 * [cf48799f4](https://github.com/kubeovn/kube-ovn/commit/cf48799f490670330c32337de81e569db0c53fea) fix(e2e): poll for ovs-ovn pod readiness instead of single snapshot check (#6364)
 * [c4eeb1d51](https://github.com/kubeovn/kube-ovn/commit/c4eeb1d516dd444bd84378716f7c7abc3c1b8142) fix(e2e): replace hard sleeps and unretried checks with WaitUntil polling (#6362)
 * [e43e1d994](https://github.com/kubeovn/kube-ovn/commit/e43e1d9943a47aaa2c67dca0f9a60280e54f0f78) fix(ci): redesign log upload steps in build-x86-image workflow (#6361)
 * [334e6c496](https://github.com/kubeovn/kube-ovn/commit/334e6c496099f3caf709441b51a1279345d8887e) fix(e2e): wait for endpoints and retry service trace in kubectl-ko test (#6358)
 * [f9865af39](https://github.com/kubeovn/kube-ovn/commit/f9865af39895515bf2b7a908e745518521498791) fix(e2e): wait for flow cleanup before checking service after deletion in metallb test (#6359)
 * [d9721475b](https://github.com/kubeovn/kube-ovn/commit/d9721475b42a1dfec1351835e7dbfae68be47073) ci: use larger runner for Talos installation test (#6360)
 * [058b3de59](https://github.com/kubeovn/kube-ovn/commit/058b3de595a5fa8b0322a6fd4602799d0f96bb47) docs: updated CHANGELOG.md (#6357)
 * [2dab3d6d9](https://github.com/kubeovn/kube-ovn/commit/2dab3d6d93b9b23db56f1a4afdc2c8a0b2177712) fix(deps): update module golang.org/x/net to v0.51.0 (#6353)
 * [1c654cbf7](https://github.com/kubeovn/kube-ovn/commit/1c654cbf7b700904e3c80b51e746ba618a9678dd) fix(e2e): poll for OVS bridge addresses in ovn-vpc-nat-gw test (#6355)
 * [bbd1abe68](https://github.com/kubeovn/kube-ovn/commit/bbd1abe682484cb58271373d330e5e6c479d27e2) chore(deps): update github artifact actions (#6356)
 * [5400b1e3f](https://github.com/kubeovn/kube-ovn/commit/5400b1e3fe5174d2b5491372156ca250a14fde00) fix(e2e): retry check-cluster after HA db corruption recovery (#6354)
 * [62408948e](https://github.com/kubeovn/kube-ovn/commit/62408948ef7bc8da21e03c73fb3caa4ee9e83a00) fix: add retry logic to kubectl-ko getOvnCentralPod for leader election (#6351)
 * [21141a3ef](https://github.com/kubeovn/kube-ovn/commit/21141a3ef399a2eae1ae961c6702fdb8021a7ac2) fix: use targeted patch in calcSubnetStatusIP to prevent U2O status overwrite (#6350)
 * [3dd3fff8f](https://github.com/kubeovn/kube-ovn/commit/3dd3fff8ff44922d7f559148a931ca43001e8b4c) fix(e2e): wait for DNSNameResolver ready before CNP connectivity checks (#6349)
 * [53b0946c5](https://github.com/kubeovn/kube-ovn/commit/53b0946c5913171c18772e1b0f35f287a46ec22b) fix: fix SLR e2e flaky failure caused by VIP not being cleaned up (#6347)
 * [9e1109c58](https://github.com/kubeovn/kube-ovn/commit/9e1109c58baf04517e4ecd2d90f941b976c33763) fix(e2e): wait for RAFT cluster convergence in HA db corruption test (#6348)
 * [d655c9002](https://github.com/kubeovn/kube-ovn/commit/d655c90021904effbabd87541c435e4cc40e2ece) fix(deps): update module github.com/cert-manager/cert-manager to v1.19.4 (#6334)
 * [21594f749](https://github.com/kubeovn/kube-ovn/commit/21594f7493d3a144f248210cbac97319fe3289fe) fix(e2e): fix checkPolicy logic bug and improve underlay test diagnostics (#6346)
 * [8e592b36a](https://github.com/kubeovn/kube-ovn/commit/8e592b36a6be720adc2c487498ce58440c3905c3) fix: add retry logic to vpc-egress-gateway e2e curl connectivity check (#6345)
 * [f5b527489](https://github.com/kubeovn/kube-ovn/commit/f5b527489c0a4e621926b93532fbb00367b414e9) fix: add retry logic to SLR e2e connectivity test to fix flaky failure (#6343)
 * [ca8b1c8e0](https://github.com/kubeovn/kube-ovn/commit/ca8b1c8e0ea46a4c0c2b5583eb97dd1792eca6b8) feat: introduce serviceMonitor and labels to services (#6340)
 * [4c8de6288](https://github.com/kubeovn/kube-ovn/commit/4c8de62880c4da38a0ac2d3cfa9b521211a84f15) chore(deps): update dependency cert-manager to v1.19.4 (#6332)
 * [344798172](https://github.com/kubeovn/kube-ovn/commit/344798172bebb8d74455dd7e2bb645257e0bc2dc) chore(deps): update dependency kubevirt to v1.7.1 (#6333)
 * [1c1cc9354](https://github.com/kubeovn/kube-ovn/commit/1c1cc9354321a19aa1619f490c2535efaaf7cfac) fix(deps): update module kubevirt.io/client-go to v1.7.1 (#6325)
 * [06f766c17](https://github.com/kubeovn/kube-ovn/commit/06f766c1747c10133158a28009bd339aa8aa0153) fix(deps): update module kubevirt.io/api to v1.7.1 (#6324)
 * [66ed19714](https://github.com/kubeovn/kube-ovn/commit/66ed19714d4b0923ed36f88a51f077d779ca3fbd) chore(deps): update kindest/node docker tag to v1.35.1 (#6317)
 * [4fe155dcf](https://github.com/kubeovn/kube-ovn/commit/4fe155dcf7366de1fb550a6d205a1013edea401c) fix(kind): set containerd LimitNOFILE to 4096 to avoid valgrind failure (#6337)
 * [42ff3555f](https://github.com/kubeovn/kube-ovn/commit/42ff3555fc17a2ba4c29b0a22568db852992d920) docs: update NEXT.md with Post-v1.15.0 changes (#6336)
 * [7ad2a5157](https://github.com/kubeovn/kube-ovn/commit/7ad2a51579ae7cff5ccb23ceec78d6b7925ec8f8) only build base for last 3 release branch
 * [820341d46](https://github.com/kubeovn/kube-ovn/commit/820341d46a6366b69dd9f18f10c52f583077283f) fix(ovn): update patches for compatibility with upstream OVN branch-25.03 (#6331)
 * [830362a7e](https://github.com/kubeovn/kube-ovn/commit/830362a7efcff7928a7c02fcf9cca61852545df6) fix(loadbalancers): ip_port_mapping is overwritten by new vips (#6302)
 * [53bad7e85](https://github.com/kubeovn/kube-ovn/commit/53bad7e85eeacca8d80937ab0c98e8a72ff84b57) fix: avoid OVN northd "No path for static route" warning on startup (#6323)
 * [6ef934de3](https://github.com/kubeovn/kube-ovn/commit/6ef934de39472e58e6b6274fc346274d3810408c) fix: skip unallocated pods during IPAM initialization to avoid spurious errors (#6327)
 * [0e8d11fd7](https://github.com/kubeovn/kube-ovn/commit/0e8d11fd7bdd0e72f448c1e06ff02406ff3c69cf) chore(deps): update actions/checkout action to v6 (#6322)
 * [0cad95f53](https://github.com/kubeovn/kube-ovn/commit/0cad95f5393bbe398d996291cf4669b68214d8ca) fix(ovn): add use-after-free fix for en-sync-sb to null check patch (#6326)
 * [7cb33fac6](https://github.com/kubeovn/kube-ovn/commit/7cb33fac6d3b1dc133191018b1a22d3f5d5844d5) feat: check for pod restarts after installation (#6311)
 * [73a05d0c9](https://github.com/kubeovn/kube-ovn/commit/73a05d0c93542952aea214f8dff8fd848ae03bbd) ci: upload kubectl-ko-log when check-restarts step fails (#6321)
 * [58b9f596e](https://github.com/kubeovn/kube-ovn/commit/58b9f596eaff0a0be8cfefdfe009a5ce675d96ab) feat(helm):Add affinity and nodeSelector support for ovs-ovn and ovs-ovn-dpdk DaemonSets using v2 chart (#6308)
 * [a067e089f](https://github.com/kubeovn/kube-ovn/commit/a067e089fc215216b62653f92286d6d1c714957e) chore(deps): update actions/setup-python action to v6.2.0 (#6319)
 * [0201cbb55](https://github.com/kubeovn/kube-ovn/commit/0201cbb55ae8edf963fe96d7fb5163c73873bdae) chore(deps): update helm/kind-action action to v1.14.0 (#6320)
 * [62b3d9b0c](https://github.com/kubeovn/kube-ovn/commit/62b3d9b0c168fd0d7804174f894a60d3e8eddd2f) fix(deps): update module gopkg.in/k8snetworkplumbingwg/multus-cni.v4 to v4.2.4 (#6318)
 * [5c6e746ff](https://github.com/kubeovn/kube-ovn/commit/5c6e746ffb57e36c6d383b13e4339e0317b19764) fix(ovn): add null check before sbrec delete calls to prevent crashes (#6316)
 * [82283a25b](https://github.com/kubeovn/kube-ovn/commit/82283a25b575b597e60d64a2dce75425fcefa1df) chore(deps): update dependency multus to v4.2.4 (#6314)
 * [bb4cbf37c](https://github.com/kubeovn/kube-ovn/commit/bb4cbf37c00378eb4e8650cd4a4ae940c695d3d0) chore(deps): update dependency talos to v1.12.4 (#6305)
 * [50144a391](https://github.com/kubeovn/kube-ovn/commit/50144a3913bc21274995ef61ffc36161a2bd7de7) fix(deps): update module google.golang.org/grpc to v1.79.1 (#6303)
 * [7d8b617bd](https://github.com/kubeovn/kube-ovn/commit/7d8b617bd0679860a00ceeb5f3a897b7612588df) chore(deps): update dependency submariner to v0.22.1 (#6298)
 * [ed37cf740](https://github.com/kubeovn/kube-ovn/commit/ed37cf740ade01fbfecef6b0ae575f6265298af2) chore(deps): update aquasecurity/trivy-action action to v0.34.1 (#6304)
 * [4c2be59aa](https://github.com/kubeovn/kube-ovn/commit/4c2be59aad7bf78126b93ab32a603eec44bb9c63) chore(deps): update actions/checkout action to v5.0.1 (#6313)
 * [5a6d74991](https://github.com/kubeovn/kube-ovn/commit/5a6d7499156ed5cf0ae3552a7d2967f20a8a3584) chore(deps): update dependency go to v1.26.0 (master) (#6292)
 * [23af56309](https://github.com/kubeovn/kube-ovn/commit/23af563099d0739f7d0bd14ef1197f539ddc1e70) chore(deps): update dependency golangci-lint to v2.10.1 (master) (#6299)
 * [3a28ee1b3](https://github.com/kubeovn/kube-ovn/commit/3a28ee1b39d1594d3dfe82a93c93b0c0bda89d68) fix(controller): prevent NPE in isOvnSubnet when subnet is nil (#6312)
 * [2332aff14](https://github.com/kubeovn/kube-ovn/commit/2332aff14da25f18701f6dfb198fa29ebc808f6e) add CLAUDE.md and link to AGENTS.md
 * [dedad4748](https://github.com/kubeovn/kube-ovn/commit/dedad474819f37903cf9fcc085a5ea2f70ab754b) feat: add GitHub Actions workflow for linting and testing Helm charts (#6184)
 * [e65ff4ef6](https://github.com/kubeovn/kube-ovn/commit/e65ff4ef6da05f6382b0562fb3e0b39844e95013) docs: updated CHANGELOG.md (#6300)
 * [7c0812df6](https://github.com/kubeovn/kube-ovn/commit/7c0812df657f29a075988c4e0a7e538af4f8aa5a) remove unnecessary logs
 * [436b519ec](https://github.com/kubeovn/kube-ovn/commit/436b519ec940f3fa7bd84f82ebbe2f96155ec763) refactor(controller): clean up redundant logic in vpc controller (#6295)
 * [a0e4135e0](https://github.com/kubeovn/kube-ovn/commit/a0e4135e0c1cb215996792f0d22bcc7ed796d157) fix: skip IPAM init for empty IP annotation and handle nil IPRangeList (#6296)
 * [8978c1818](https://github.com/kubeovn/kube-ovn/commit/8978c18181ed940ab85fd514719e8c840d5a53d1) networkpolicy: provider-scoped policies for multi-network pods (#6223)
 * [401160e77](https://github.com/kubeovn/kube-ovn/commit/401160e77b6754db6aed748e3fb104f148378921) fix duplicate description
 * [524b6c0a4](https://github.com/kubeovn/kube-ovn/commit/524b6c0a45e16e28f9373e7be3f894b4ceca96ab) fix(iptables-eip): look up subnet directly by name instead of resolving through NAD (#6287)
 * [503fbc8b7](https://github.com/kubeovn/kube-ovn/commit/503fbc8b7a971b6e1b392eec7060d9bdf2a3aaf7) chore(deps): update dependency helm to v4.1.1 (#6290)
 * [e819cb30b](https://github.com/kubeovn/kube-ovn/commit/e819cb30b5b5423f3ca19c0ffe3eb9b42e788c1a) fix(deps): update golang (#6288)
 * [7dfd0235c](https://github.com/kubeovn/kube-ovn/commit/7dfd0235c5cb537f60c4604a67352b20e6c4cefc) refactor(ovn_leader_checker): improve readability and naming (#6263)
 * [772531f8b](https://github.com/kubeovn/kube-ovn/commit/772531f8bda333020727345f28f928fef3b6fa7e) docs: updated CHANGELOG.md (#6284)
 * [08c83ddcf](https://github.com/kubeovn/kube-ovn/commit/08c83ddcfdcb4d1dc3bd66b8b063f2172fbf0b54) add skip test
 * [1adb28a8e](https://github.com/kubeovn/kube-ovn/commit/1adb28a8e68225d7f61baf4db9740fd6d2fe071b) perf: only enqueue subnet when not ready (#6268)
 * [45e0afead](https://github.com/kubeovn/kube-ovn/commit/45e0afeadfa557e99d56efeba43497b10cd4e9a4) fix(deps): update golang (#6281)
 * [2186c1342](https://github.com/kubeovn/kube-ovn/commit/2186c134268e417acae44bed60cb1adf84cb5a5a) fix: go modernize doesn't work (#6262)
 * [0f698339b](https://github.com/kubeovn/kube-ovn/commit/0f698339b553f907c1a58655e731688af41b90f8) chore(deps): update dependency talos to v1.12.3 (#6278)
 * [8f48fa9a5](https://github.com/kubeovn/kube-ovn/commit/8f48fa9a58c9a768df687c19a4ff1eb82536a363) controller: fix stale ip cr after subnet deletion (#6275)
 * [d43c8a270](https://github.com/kubeovn/kube-ovn/commit/d43c8a270f143d4b3349d7dff85f6772dd6e65d9) fix(deps): update module github.com/prometheus-community/pro-bing to v0.8.0 (#6277)
 * [066ce1208](https://github.com/kubeovn/kube-ovn/commit/066ce12088bec59fe8975dd4e866e17f7c06ab5e) docs: updated CHANGELOG.md (#6276)
 * [c76533220](https://github.com/kubeovn/kube-ovn/commit/c765332205dce5b8582fb36038e8cedabe378858) cni-server: add static fdb entry for subnets with u2o enabled (#6269)
 * [f3e7304cc](https://github.com/kubeovn/kube-ovn/commit/f3e7304cc5fa05248ff76b69e793d2313b2cebb6) controller: fix ips owned by subnet u2o ip and mcast querier ip deleted in gc (#6132)
 * [d55ce5ac3](https://github.com/kubeovn/kube-ovn/commit/d55ce5ac353dc0d4ebe55a35455eb7562a7f9761) fix(deps): update module go.podman.io/image/v5 to v5.39.1 (#6271)
 * [081cef09c](https://github.com/kubeovn/kube-ovn/commit/081cef09c82f1eca4c78932b2b5a6285095a8307) chore(deps): update dependency go to v1.25.7 (#6264)
 * [15c7c7984](https://github.com/kubeovn/kube-ovn/commit/15c7c7984b14460f6ce705caa3db69c9608d2cce) security: specify a storage limit for containers (#6259)
 * [872553774](https://github.com/kubeovn/kube-ovn/commit/8725537741614912363ccfef7304994a32589e5b) feat: add user-defined annotations support for nat gw sts template creation (#6256)
 * [f97672edf](https://github.com/kubeovn/kube-ovn/commit/f97672edfab2bc17a2f849226d6d8c7f5463613d) charts: add external-gateway-config-ns option for controller (#6211)
 * [193f231a1](https://github.com/kubeovn/kube-ovn/commit/193f231a19787d77d045e3da8470216a79867d55) refactor: extract pprof and metrics/health server helpers to pkg/metrics (#6260)
 * [f285541a4](https://github.com/kubeovn/kube-ovn/commit/f285541a479755ec7bd9caeb3f79b1e2f1377638) refactor: profiling, Retry, assignAddress, runIfLinkExists (#6258)
 * [93fca05a6](https://github.com/kubeovn/kube-ovn/commit/93fca05a6db26a320ed21708b843868806a7e4ce) ignore .context
 * [4b4df5fa6](https://github.com/kubeovn/kube-ovn/commit/4b4df5fa6f3d46dd3e63ca981ea7ba7380159773) fix: not gen vlan deepcopy (#6257)
 * [dc7afa4e6](https://github.com/kubeovn/kube-ovn/commit/dc7afa4e6458598dba9a1dcabda8fa992fdd1425) fix typo (#6253)
 * [6dfd15168](https://github.com/kubeovn/kube-ovn/commit/6dfd15168803f9c7ebd3093fa8406336bf91d766) fix(deps): update module github.com/cert-manager/cert-manager to v1.19.3 (#6251)
 * [70a2e3cd7](https://github.com/kubeovn/kube-ovn/commit/70a2e3cd78a6ed1835c224b25184eb05ff361997) chore(deps): update dependency cert-manager to v1.19.3 (#6250)
 * [2e089ff34](https://github.com/kubeovn/kube-ovn/commit/2e089ff3479527a6d14b6e34490e21a8703ddb6f) Fix qos not work (#6231)
 * [fb806189b](https://github.com/kubeovn/kube-ovn/commit/fb806189b7842609295a19e4c4444b565dd58f2b) docs: updated CHANGELOG.md (#6247)
 * [0ae88bbf6](https://github.com/kubeovn/kube-ovn/commit/0ae88bbf666567a00cf83cfd59819b1b180816b0) fix: Ensure Kubevirt optimizations do not use old migration state (#6239)
 * [9d4330ef9](https://github.com/kubeovn/kube-ovn/commit/9d4330ef9a843582a19b86a155bc8549c4663226) Set Kubevirt live-migrations multi-chassis options for all VM nics (#6241)
 * [98a3b1ecd](https://github.com/kubeovn/kube-ovn/commit/98a3b1ecd1316fe6cc8d9b40eab2cde66efdb969) fix localnet port should not be static "patch-localnet." (#6244)
 * [7cd3789df](https://github.com/kubeovn/kube-ovn/commit/7cd3789df9908104487b8b21a80112a62164fb50) rbac: replace wildcard with a clear list of allowed verbs (#6233)
 * [8f63d677e](https://github.com/kubeovn/kube-ovn/commit/8f63d677ebe1922a5bf5bdb3355a4f73bd1e4d68) fix(deps): update module github.com/onsi/gomega to v1.39.1 (#6237)
 * [9b814d74f](https://github.com/kubeovn/kube-ovn/commit/9b814d74fb3419c0c66bdc5e9eedaf898740ee4f) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.28.1 (#6238)
 * [64e05655f](https://github.com/kubeovn/kube-ovn/commit/64e05655f43d9eedce1bd5bf7a4f8b2aca3c278d) remove duplicate description (#6234)
 * [ff584e907](https://github.com/kubeovn/kube-ovn/commit/ff584e907facea2d35843b403421f01465fbcba5) refactor: do not read pid file when calling ovs/ovn appctl command (#6230)
 * [a0dd53de4](https://github.com/kubeovn/kube-ovn/commit/a0dd53de4d95ca6f4bf12ea6968d1926a21dc65b) docs: updated CHANGELOG.md (#6232)
 * [2238e2f6a](https://github.com/kubeovn/kube-ovn/commit/2238e2f6a2c27b57c27ce7f0939742d1836a9694) fix(vpcnatgw): nat gateway uses faulty default network (#6212)
 * [503af3dbe](https://github.com/kubeovn/kube-ovn/commit/503af3dbecc941c5d433655d623e52916b1549b6) docs: updated CHANGELOG.md (#6225)
 * [98c448618](https://github.com/kubeovn/kube-ovn/commit/98c448618eec4439f4c6d0f1eb432b175ad512b9) controller: fix vpc egress gateway forward policy not working for targets within the internal subnet's cidr (#6222)
 * [80aef1dbe](https://github.com/kubeovn/kube-ovn/commit/80aef1dbe15fccc8ee26bce5f7b25798681437a4) fix(deps): update module github.com/godbus/dbus/v5 to v5.2.2 (#6219)
 * [42393cf94](https://github.com/kubeovn/kube-ovn/commit/42393cf9475b938410a5358cbd615b0aa37bb3fc) controller: fix vpc egress gateway not working for targets within the internal subnet's cidr (#6218)
 * [7fade0694](https://github.com/kubeovn/kube-ovn/commit/7fade0694f922f61cf56786fede113848afdb3fa) e2e: fix OOM on ginkgo build (#6217)
 * [083389754](https://github.com/kubeovn/kube-ovn/commit/083389754097f3a5f6b02a8fc013ca860616383d) Update start-ic-db.sh to have gen_conn_addr function defined (#6216)
 * [87993313d](https://github.com/kubeovn/kube-ovn/commit/87993313df5b23dfd0fef66b79cc705c42ebbee4) fix(deps): update module github.com/moby/moby/client to v0.2.2 (#6213)
 * [af72723b2](https://github.com/kubeovn/kube-ovn/commit/af72723b2e4f3d58c405b07443bc8537047ceac7) fix: caching NAD CRD should before all kubeovn crds and pod (#6198)
 * [dd1880d50](https://github.com/kubeovn/kube-ovn/commit/dd1880d5048ab2e4d740d6d7b9173bf4ed88dd39) docs: updated CHANGELOG.md (#6209)
 * [cb157cfa9](https://github.com/kubeovn/kube-ovn/commit/cb157cfa9d67f2f4db1269a4b18eaf28f664dce7) fix(vpc-nat-gw): use consistent provider for all pod annotations (#6204)
 * [658725ff2](https://github.com/kubeovn/kube-ovn/commit/658725ff2f0ac5d62cb28386b42df589e5683e83) fix(deps): update module github.com/puzpuzpuz/xsync/v4 to v4.4.0 (#6206)
 * [9636d7247](https://github.com/kubeovn/kube-ovn/commit/9636d7247a492eb2f2a3e0ec1f08b6311cff68f6) fix: host no kube-ovn-speaker log (#6196)
 * [ab5106902](https://github.com/kubeovn/kube-ovn/commit/ab5106902442f64d43092a859678398361c3f4fa) docs: updated CHANGELOG.md (#6203)
 * [25cc1dd34](https://github.com/kubeovn/kube-ovn/commit/25cc1dd3467be6805d1e989f309fb6b07aa1ca1f) test: run ginkgo via go tool (#6191)
 * [62b1494ff](https://github.com/kubeovn/kube-ovn/commit/62b1494ffd83fa9d9029059057f70d4c0ae946d5) fix: skip NetworkUnavailable condition in non-primary CNI mode (#6195)
 * [bc813ee34](https://github.com/kubeovn/kube-ovn/commit/bc813ee34ae8274784b8a1bdfb0c66b1f9e29114) fix(vpcnatgw): cannot DNAT same EIP and same external port on two different protocols (#6201)
 * [64c3ac759](https://github.com/kubeovn/kube-ovn/commit/64c3ac759a844fc85a2f42dee14740a14c8d6903) chore(deps): update dependency helm to v4.1.0 (#6200)
 * [f87749dec](https://github.com/kubeovn/kube-ovn/commit/f87749decc0fdb00f23a38689a49da62d99d7e26) chore(deps): update dependency talos to v1.12.2 (#6199)
 * [83a7b37d2](https://github.com/kubeovn/kube-ovn/commit/83a7b37d29cf79a37f9417fb9ada91418d700e07) fix: cleanup auto-created vlan subinterfaces on interface switch (#6193)
 * [21d6cbbf6](https://github.com/kubeovn/kube-ovn/commit/21d6cbbf6c46d420b4a6d863845963d0f54ef4d5) fix master base image build err
 * [6d06763bf](https://github.com/kubeovn/kube-ovn/commit/6d06763bf6c47c7511c29ebe820e187587f29d49) Add vendor ID to transit switches (#6186)
 * [6a40d68ed](https://github.com/kubeovn/kube-ovn/commit/6a40d68ed9b688e29fab2b5dadc820082f2bac31) fix: enforce neighbor address validation in configuration (#6134)
 * [a52abe2d1](https://github.com/kubeovn/kube-ovn/commit/a52abe2d17e0b60b97bce92097ad5a544882fd9e) fix: webhook e2e (#6163)
 * [864f2c678](https://github.com/kubeovn/kube-ovn/commit/864f2c678521273b3269524531a0c69f62fbf9a5) issue#6181 lb的IPPortMappings不为空时，不允许忽略IPPortMapping (#6182)
 * [332875929](https://github.com/kubeovn/kube-ovn/commit/3328759297932d8b06d550e150698dbf9c1d3b8c) refactor: consolidate pod deletion logic for stateful and VM pods (#6183)
 * [cdf4feedd](https://github.com/kubeovn/kube-ovn/commit/cdf4feedddfc112b6b7c519411a74123f0ddc828) fix: acquire specified subnet from annotation (#6180)
 * [3979f55cf](https://github.com/kubeovn/kube-ovn/commit/3979f55cf9a2bc5bc15b453b0ba8c5861c04b20c) docs: updated CHANGELOG.md (#6175)
 * [6d7d5ebbb](https://github.com/kubeovn/kube-ovn/commit/6d7d5ebbb21bb8751ce14ff4c5205afd60681a52) chore(deps): update dependency helm to v4.0.5 (#6167)
 * [968dfbb14](https://github.com/kubeovn/kube-ovn/commit/968dfbb145b08444e6daffb5a6581c21e69abd5d) fix(deps): update module github.com/sirupsen/logrus to v1.9.4 (#6168)
 * [77e5e9c17](https://github.com/kubeovn/kube-ovn/commit/77e5e9c17cf569e0ce64f4f2bb84802691179b81) chore(deps): update dependency go to v1.25.6 (#6169)
 * [0ae0d13ec](https://github.com/kubeovn/kube-ovn/commit/0ae0d13ece50b33d22b3c8a7ce1fe996ddc1abde) add metallb underlay v6 dualcase (#6159)
 * [2238d4797](https://github.com/kubeovn/kube-ovn/commit/2238d47972c419036d04470cb2f836ee6362d1b7) feat: support static IP/MAC for multiple interfaces on same switch (#6060)
 * [d74b7c696](https://github.com/kubeovn/kube-ovn/commit/d74b7c696a2ff0dc9e352cebbfba23c4d479135e) fix(helm): rolling upgrade strategy for ovs is wrong (#6160)
 * [5d93ff826](https://github.com/kubeovn/kube-ovn/commit/5d93ff826c6c0c65bd39ae3dc9d758109b1c1de8) chore(deps): update dependency cilium to v1.18.6 (#6161)
 * [a8b356bbb](https://github.com/kubeovn/kube-ovn/commit/a8b356bbb78ca74e8ca0837c38620e1981258d7a) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.27.5 (#6158)
 * [45174b137](https://github.com/kubeovn/kube-ovn/commit/45174b137d031bd1c9d137ca977e9c1278e79e73) fix(deps): update module github.com/puzpuzpuz/xsync/v4 to v4.3.0 (#6156)
 * [80b830661](https://github.com/kubeovn/kube-ovn/commit/80b830661918cc5f17e6633be763bf9ea5753875) fix(deps): update golang (#6154)
 * [fc2966448](https://github.com/kubeovn/kube-ovn/commit/fc2966448614af38c9dc9f79b7bb590c10d9c811) chore(deps): update module golang.org/x/text to v0.33.0 (#6150)
 * [670b875e6](https://github.com/kubeovn/kube-ovn/commit/670b875e6074aec9b4b294ae69165034ab3d37f5) fix(deps): update golang (#6148)
 * [5d0dec44d](https://github.com/kubeovn/kube-ovn/commit/5d0dec44dd71b6e2e984421792d43b4f4f8d5125) some left crd add description (#6147)
 * [0661c5b07](https://github.com/kubeovn/kube-ovn/commit/0661c5b07140e26898cd2f3fa71f565437bd5914) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.27.4 (#6140)
 * [1fd404e5f](https://github.com/kubeovn/kube-ovn/commit/1fd404e5f80118d8d325c8f4a9c721b76446d20a) add kube-ovn cr description (#6133)
 * [d44e5738b](https://github.com/kubeovn/kube-ovn/commit/d44e5738b365e7f0ca76d820fe744ca46ed7af5b) feat(chart): introduce extraEnv variable for all components (#6142)
 * [74cd1481a](https://github.com/kubeovn/kube-ovn/commit/74cd1481af5d333e86fc26923eee065ead29c425) update dep check branch
 * [079fa2fc8](https://github.com/kubeovn/kube-ovn/commit/079fa2fc86ca0bbec81b99f6ee0c89ce35ab8a1c) chore(deps): update dependency golangci-lint to v2.8.0 (#6141)
 * [974e92b31](https://github.com/kubeovn/kube-ovn/commit/974e92b310cc809eab10c1672a1737ff7ac5c72d) fix(deps): update module github.com/onsi/gomega to v1.39.0 (#6143)
 * [fd2f14d2d](https://github.com/kubeovn/kube-ovn/commit/fd2f14d2d831b5a1caabea3c2d9fc3d37a9aff9a) fix(deps): update module golang.org/x/sys to v0.40.0 (#6144)
 * [66525ea12](https://github.com/kubeovn/kube-ovn/commit/66525ea128b8744c27eee34ca9f47a36f084a466) feat(chart): make DaemonSet update strategy configurable in values.yaml (#6136)
 * [a8300e338](https://github.com/kubeovn/kube-ovn/commit/a8300e33879c030d39e1bc4c281bc92c5c1d22d9) fiix panic if namespace is deleted before pod is processed (#6131)
 * [0c9250b2b](https://github.com/kubeovn/kube-ovn/commit/0c9250b2b9cc2ff554e2bcff3a535d1a5032e94f) ci: fix loading vpc-nat-gateway image into kind cluster (#6128)
 * [547ea5e17](https://github.com/kubeovn/kube-ovn/commit/547ea5e176ae54a1b36723263fc976aa066e7f36) ci: fix non-primary cni e2e test (#6127)
 * [b6a25a243](https://github.com/kubeovn/kube-ovn/commit/b6a25a24398401e9d13b05a2e64cf3eb04b1fd1a) chore(deps): update dependency submariner to v0.22.0 (master) (#6014)
 * [9223fa292](https://github.com/kubeovn/kube-ovn/commit/9223fa292c6797b6871915a8b7d0b78343176cf9) fix: add validation for iptables rules to prevent malformed commands in nat gateway pod (#6125)
 * [8afcfa8c0](https://github.com/kubeovn/kube-ovn/commit/8afcfa8c02d864c87d7c9f1d1fdad76b48aa6b58) sync vpc-dns crd
 * [77da4e19c](https://github.com/kubeovn/kube-ovn/commit/77da4e19c9d1a569ff2773dd51c4eb11686322c9) vpc-dns custom corefile
 * [cd48d86a5](https://github.com/kubeovn/kube-ovn/commit/cd48d86a5ef8f6e8e801bbbac1a6ec59d8b83a15) chore(deps): update dependency talos to v1.12.1 (#6126)
 * [5dea828d0](https://github.com/kubeovn/kube-ovn/commit/5dea828d06999b18e0e9a74263eb1e0ff0eafdd8) set transform on informers to remove objects' managed fields (#6119)
 * [5d79ed901](https://github.com/kubeovn/kube-ovn/commit/5d79ed90194a1909a0d3444660d640801f0c89ef) Add open flow sync refer to ovn-k8s (#6117)
 * [e02d4b9f4](https://github.com/kubeovn/kube-ovn/commit/e02d4b9f49efd1030eceadc4b7876c0160b1baa9) chore(deps): update dependency gobgp to v4.2.0 (#6121)
 * [0e102b54f](https://github.com/kubeovn/kube-ovn/commit/0e102b54f6d60cd8d96d1b58f66611b602c437e4) fix(deps): update module github.com/osrg/gobgp/v4 to v4.2.0 (#6122)
 * [ce91f14a6](https://github.com/kubeovn/kube-ovn/commit/ce91f14a6e51c86f9d0698e135fe367539646b36) docs: updated CHANGELOG.md (#6120)
 * [c9fa657b1](https://github.com/kubeovn/kube-ovn/commit/c9fa657b17deeb0bedda056c957cea3ae2f1e976) prepare for next release

### Contributors

 * Abhishek Pandey
 * Akbar Kusuma Negara
 * Aleksei Sviridkin
 * Asish Kumar
 * Chanyeol Yoon
 * Chris Chiu
 * Eli Treuherz
 * ILYAS ABDELLAOUI
 * Johann Schley
 * Juan Morete
 * Mengxin Liu
 * SKALA NETWORKS
 * SkalaNetworks
 * Tongruishan
 * Uğur Tafralı
 * Volodymyr Tsap
 * Zespre Chang
 * andrewlee1089
 * changluyi
 * clyi
 * dependabot[bot]
 * dnugmanov
 * donghee-park-rsk
 * github-actions[bot]
 * jimyag
 * jskrill
 * renovate[bot]
 * zbb88888
 * zhangzujian
 * 张祖建

## v1.15.10-alpha.3 (2026-09-18)

 * [3db986bac](https://github.com/kubeovn/kube-ovn/commit/3db986bacad0c93cdc701720cd7fce23de152329) fix(ci): attach assets before publishing releases (#93)

### Contributors

 * jimyag

## v1.15.10-alpha.2 (2026-09-18)

 * [7ecd260db](https://github.com/kubeovn/kube-ovn/commit/7ecd260dbbd5ea892c5818b666cb4613d9f9c19b) fix(ci): reuse stable base image for tagged releases (#92)

### Contributors

 * jimyag

## v1.15.10-alpha.1 (2026-09-18)

 * [cbe5bdca4](https://github.com/kubeovn/kube-ovn/commit/cbe5bdca48d4554452eebdd260a9332aa7256d23) fix(ci): allow tags from any release version (#89)
 * [836e677ad](https://github.com/kubeovn/kube-ovn/commit/836e677ad4aa226f26b20eafda6689315f808315) ci: publish tagged x86 releases to Qiniu GHCR (#88)
 * [765dd0cac](https://github.com/kubeovn/kube-ovn/commit/765dd0cacb02a1414358620b1e95a8386248d77a) Fix vm migrate with pg (#84)
 * [8688ff6a5](https://github.com/kubeovn/kube-ovn/commit/8688ff6a55b14786ef0a79531ee8a407a3d02320) fix(controller): configure VMI migration ports while pending (#83)
 * [e87ed1f71](https://github.com/kubeovn/kube-ovn/commit/e87ed1f71f1e6af996c45f357167cf3cffd56928) fix(daemon): stop GC from deleting OVS interfaces of live-migrating VMs (#82)
 * [e9ba148c8](https://github.com/kubeovn/kube-ovn/commit/e9ba148c8dabbc9b042b5af85f0aa05c668617cb) feat: backport nftable LB service support and session affinity (#80)
 * [e2b7ce1dc](https://github.com/kubeovn/kube-ovn/commit/e2b7ce1dc7eb4e88080cd84719c7966f70b1d8a0) fix: avoid rejecting initial QoS status reconciliation (#78)
 * [d61decb85](https://github.com/kubeovn/kube-ovn/commit/d61decb854ba1bea66e4299b41b2c0b87876cdba) fix(vpc-nat-gw): reject uppercase dnat protocols (#77)
 * [55ac2aebf](https://github.com/kubeovn/kube-ovn/commit/55ac2aebfd447798fd40cda3d7c90feebe442551) fix: release NAT claims when EIP is deleting (#76)
 * [7df728f04](https://github.com/kubeovn/kube-ovn/commit/7df728f04227f6caed3f09edec49d33ace772706) fix(vpc-nat-gw): reconcile refs when dependencies become available (#74)
 * [1f875b516](https://github.com/kubeovn/kube-ovn/commit/1f875b516bcd91697f19c02bb6273d0e3ec1647b) fix(vpc-nat-gw): reject bindings to terminating CRDs (#72)
 * [4bf41f1cb](https://github.com/kubeovn/kube-ovn/commit/4bf41f1cb6b07aa96ec8221c40d9a6f66b5b01d5) daemon: accept node internal IP as encap IP even with /32 mask (#71)
 * [70a003205](https://github.com/kubeovn/kube-ovn/commit/70a003205d1fe2cddf9e7cc2b48f6a37267bf123) controller: guarantee tunnel_key annotation before CNI ADD and restart (#70)
 * [8630fcc83](https://github.com/kubeovn/kube-ovn/commit/8630fcc835b7e5c12875345eb67916532dba1458) Fix eip qos nat left (#69)
 * [d58bf268e](https://github.com/kubeovn/kube-ovn/commit/d58bf268ed27e6f965e970711ec761739c1bef60) daemon: keepalive loop to auto-recover macvlan sub-interface for node-local EIP access (#68)
 * [82c241722](https://github.com/kubeovn/kube-ovn/commit/82c241722e18bf3ab03e8e29b9448ee9d42b5784) fix(controller): wait for subnet tunnel key before allocating pod IP (#65)
 * [413c60838](https://github.com/kubeovn/kube-ovn/commit/413c608383fe59f4e7da04bfc39f8b67b299e5af) Bump bgp to 4.7.0 (#64)
 * [ac51abc9d](https://github.com/kubeovn/kube-ovn/commit/ac51abc9d4f77edec3dfba5963491e810743a071) 7niu release 1.15 kubespray Guaranteed ovn pod (#63)
 * [025ef9fb3](https://github.com/kubeovn/kube-ovn/commit/025ef9fb3f5cb7f31f4c76d927e8ccea09c8a421) feat: add Namespace and ExternalSubnet columns to iptables-eip printer (#62)
 * [8d7c77df8](https://github.com/kubeovn/kube-ovn/commit/8d7c77df83e37c418937db580b4521282515821c) feat: support share type iptables DNAT with nft map-based LB (#60)
 * [915770540](https://github.com/kubeovn/kube-ovn/commit/9157705403822ce33b4cfbb4605bef13c46f01e4) fix: clean stale u2o overlay cidr address sets
 * [36aa6492e](https://github.com/kubeovn/kube-ovn/commit/36aa6492eab1eae611ef2e413e62ab9ad31ac86e) fix: enqueue subnet ACL and selector updates
 * [50ab7ce10](https://github.com/kubeovn/kube-ovn/commit/50ab7ce10078502b4f2ee859c8a59b589ba7cbc5) fix: tighten u2o stale policy cleanup
 * [60da5ed4a](https://github.com/kubeovn/kube-ovn/commit/60da5ed4a1dc2a1c4d68ed57351880cb946ed742) ci: normalize master e2e branch for 1.15 skip logic
 * [e05d8b456](https://github.com/kubeovn/kube-ovn/commit/e05d8b4566d49964824a23142c8a81de8e55720d) feat: backport u2o overlay only routing to release-1.15 (#6751)
 * [dc5bfe570](https://github.com/kubeovn/kube-ovn/commit/dc5bfe5706c5dd6028ce90ca2f4cbfc220c0befd) feat(speaker): add BGP/BFD Prometheus metrics and delta-based BFD error logging
 * [aa30601a7](https://github.com/kubeovn/kube-ovn/commit/aa30601a7d1d6a9588e1b1cce62051667f1c62f1) fix as review
 * [68bb760a8](https://github.com/kubeovn/kube-ovn/commit/68bb760a82f1f6f5a54ed608828240e353690a24) Feat(speaker): add GoBGP BFD support for fast failure detection (#6912)
 * [d81623d0f](https://github.com/kubeovn/kube-ovn/commit/d81623d0fbdbb7e8308868d689253ef188254a86) fix(controller): skip in-pod cleanup when VPC NAT gateway CRD is deleted (#6784)
 * [2820b96a7](https://github.com/kubeovn/kube-ovn/commit/2820b96a7ca4baa1e6d09953d6a864990cf2f077) fix fmt
 * [31ec88665](https://github.com/kubeovn/kube-ovn/commit/31ec8866513d57b7f27fa3641903e5710044a325) fix as review
 * [e6ded0c4b](https://github.com/kubeovn/kube-ovn/commit/e6ded0c4b2c3d6cdc82499ee9107ba0d1362a862) fix as reviews
 * [f6ccb1b6c](https://github.com/kubeovn/kube-ovn/commit/f6ccb1b6cc72fee355296dd426de55a5165970b2) feat(bgpspeaker): gate lb vip announcements by node label
 * [01540822a](https://github.com/kubeovn/kube-ovn/commit/01540822aa08d01155dcda175ea09c3fba3fafe1) fix reviews
 * [f27bb558b](https://github.com/kubeovn/kube-ovn/commit/f27bb558bcd0466b0ec2e3088032c2fad1857e3a) speaker: add whitelisted source address extension
 * [50e433feb](https://github.com/kubeovn/kube-ovn/commit/50e433febfff65786877b3e078aee05280346f5e) self review
 * [3fd98ca87](https://github.com/kubeovn/kube-ovn/commit/3fd98ca8722fa1db7418ce59cb737bafb4e97cdc) project ovn lb svc cluster ip
 * [6303f86ea](https://github.com/kubeovn/kube-ovn/commit/6303f86eacb3edf6bc97ffd2fbdd885746d0999d) fix webhook
 * [6484b5359](https://github.com/kubeovn/kube-ovn/commit/6484b5359a31862878a3a158f96b5dc821661003) fix xgo reviews
 * [c088baf90](https://github.com/kubeovn/kube-ovn/commit/c088baf90b1a51eaf7781961ced0b617b5634475) fix ut
 * [48109a164](https://github.com/kubeovn/kube-ovn/commit/48109a164003491bfcb8336c61ef88fb28938736) fix: lb reconcile
 * [7cc585bcb](https://github.com/kubeovn/kube-ovn/commit/7cc585bcb12c1c225c6be2a3f9f6dcb0dddf79fc) fix xgo reviews
 * [f4b8bc7b3](https://github.com/kubeovn/kube-ovn/commit/f4b8bc7b39b0efc0c748cf52dade7d0cc38d4c60) fix xgo reviews
 * [22ab31731](https://github.com/kubeovn/kube-ovn/commit/22ab31731b0afa0f5aadb8bd7dc2315f1a42dc91) remove not relevant np code
 * [d8c038202](https://github.com/kubeovn/kube-ovn/commit/d8c038202223f344c88a29d495d088c4cd93c8be) only handle LB add | del: no update
 * [eae73e380](https://github.com/kubeovn/kube-ovn/commit/eae73e380f408868544bda01bfd3fbf0f4e18715) simple
 * [04bff209c](https://github.com/kubeovn/kube-ovn/commit/04bff209cbc1b77fb40cb210d10f3cf5b739fcf8) fix as xgo reviews
 * [f40cbde8c](https://github.com/kubeovn/kube-ovn/commit/f40cbde8c76437dfddcc3f12e2765112fd35de60) fix as xgo reviews
 * [3035b8ea2](https://github.com/kubeovn/kube-ovn/commit/3035b8ea27b8fb20d46b7ec62212d9147a0de4b0) compat metal lb annotation
 * [339c26656](https://github.com/kubeovn/kube-ovn/commit/339c26656192451168f3371927556a05eeb7888a) LB: single node and ecmp all nodes
 * [2fb51cc0f](https://github.com/kubeovn/kube-ovn/commit/2fb51cc0fd2e8a16550228798544bbbf279dfbc0) speaker: update enable-lb-svc-announce for LB EIP
 * [756fb173a](https://github.com/kubeovn/kube-ovn/commit/756fb173a56e8fe4ad39374cf079f696be25d6c8) add log
 * [1da1f4aac](https://github.com/kubeovn/kube-ovn/commit/1da1f4aac37ead86d57aaddb4ef3dfd0b91624df) fix ut
 * [c2f3db856](https://github.com/kubeovn/kube-ovn/commit/c2f3db856889cda9b308c6624e1163439b4cda03) fix(ut): remove stale externalIPs fixtures and assertions from bgp-lb-vip tests fix(ut): remove stale ExternalIPs fixture; add log for empty VIP IP fix(controller): return nil when VIP has no IP yet in reconcileBgpLbVipServiceLocked
 * [f8fc74da6](https://github.com/kubeovn/kube-ovn/commit/f8fc74da61793837f39b248981c2065bdf7e1340) fix(ut): update TestReconcileBgpLbVipServiceLocked to match no-externalIPs behavior
 * [fb3bee2e6](https://github.com/kubeovn/kube-ovn/commit/fb3bee2e6c2b1e006837ce53aa95fc7a5526540b) fix(controller): align bgp-lb-vip service binding with MetalLB and remove stale externalIPs flow
 * [088f46f7a](https://github.com/kubeovn/kube-ovn/commit/088f46f7a58ec4770f6c8be226dbc84382383e0e) refactor(controller): extract OVN LB and BGP LB VIP logic into dedicated helpers
 * [1f516022e](https://github.com/kubeovn/kube-ovn/commit/1f516022ef9902c982000c58f30dbcc085a5c36c) docs(controller): add mode-matrix comments to service handler functions
 * [8121bc73b](https://github.com/kubeovn/kube-ovn/commit/8121bc73b22f657cce32e8ee2638f715cde06232) fix(controller): gate OVN LB operations in handleUpdateService behind EnableLb
 * [3bbad346e](https://github.com/kubeovn/kube-ovn/commit/3bbad346ea4d0c8011a8f9f3902bb93a4cf9ee62) refactor(config): clarify EnableLb flag and its usage in OVN load-balancer mode
 * [99017b462](https://github.com/kubeovn/kube-ovn/commit/99017b462b5f9463717fdc2e74347b9249b3a8f8) fix(controller): skip OVN LB deletion in handleDeleteService when EnableLb=false
 * [ea12bb572](https://github.com/kubeovn/kube-ovn/commit/ea12bb5720b559dd6990303b7e6faba89315dff5) chore: clean up controller-gen CRD formatting noise in install.sh and charts
 * [2e1f248e7](https://github.com/kubeovn/kube-ovn/commit/2e1f248e748c38274fc50085f9c9270638c3df02) fix as review
 * [bb6378969](https://github.com/kubeovn/kube-ovn/commit/bb6378969a110e6d0ba94566861e1448030af55c) feat(controller/speaker): add BGP LB EIP support via non-ovn subnet VIP
 * [ec3bff6e4](https://github.com/kubeovn/kube-ovn/commit/ec3bff6e4b576fd302728a2789ce8c85ceee8cb4) fix Dockerfile.ubuntu also
 * [a9a0aa96f](https://github.com/kubeovn/kube-ovn/commit/a9a0aa96faf9a10630432f48fc27af6727215403) fix: docker timeout always
 * [ac0df06f4](https://github.com/kubeovn/kube-ovn/commit/ac0df06f4965374d510ff7c5c5f3352714e539b2) fix as reviews
 * [a31efae81](https://github.com/kubeovn/kube-ovn/commit/a31efae81203d7b1b4541d513ebdf527cf1b399e) faster
 * [182cd3262](https://github.com/kubeovn/kube-ovn/commit/182cd3262691cd224de3fc4c6fbfea735fe7b865) feat(lint): add golangci-lint installation step to pre-build checks
 * [02e9f129b](https://github.com/kubeovn/kube-ovn/commit/02e9f129b40a6868d327a46882f40173282430a7) refactor(workflow): clean up runner-info job and enhance pre-build checks
 * [c2af28a34](https://github.com/kubeovn/kube-ovn/commit/c2af28a34b465aeb17c46b28888726580aa4c450) fix(review): address xgopilot review comments for PR #37
 * [dc425524e](https://github.com/kubeovn/kube-ovn/commit/dc425524ecc2409b73c6d4e50c2d6b014254b353) release v1.15.10
 * [55d5a1efb](https://github.com/kubeovn/kube-ovn/commit/55d5a1efb33f0ea97daad7a87e87a9801b83b709) Fix/upgrade 1.15 error (#6615)
 * [6b89d570e](https://github.com/kubeovn/kube-ovn/commit/6b89d570e4a16955886134699b1213b62d3666b2) ci: add e2e-approval environment for heavy jobs
 * [83e899f3e](https://github.com/kubeovn/kube-ovn/commit/83e899f3e205218ed425722ac9816374fe733732) fix(deps): sync k8s.io dependencies to v0.34.6 in test modules
 * [07e675444](https://github.com/kubeovn/kube-ovn/commit/07e67544473ad743f462f57fec0dbbd3682c4246) fix(deps): sync go versions in test modules and bump k8s.io/component-base to v0.35.3
 * [aa0084a25](https://github.com/kubeovn/kube-ovn/commit/aa0084a25eeadf7defe1ca69c7b78e1e19c38709) fix(deps): sync default versions in download-go-deps.sh with Dockerfile.base
 * [978b66e0c](https://github.com/kubeovn/kube-ovn/commit/978b66e0cfff911a0bd0497b64eb692e382afd53) fix(lint): use specific G204 rule for nosec in ipsec.go
 * [31117e01f](https://github.com/kubeovn/kube-ovn/commit/31117e01f3023c8f6c778ca6599021e37f803cf3) chore: sync trivy-action version to 0.34.2 in build-arm64-image
 * [a6f115d05](https://github.com/kubeovn/kube-ovn/commit/a6f115d050834fe3365f67740d5e77a044da5d57) fix(lint): resolve G703 path traversal warning in cniserver
 * [ed8335fad](https://github.com/kubeovn/kube-ovn/commit/ed8335fad5ed852dfad58f54e320a258c5076dad) chore: retain both kind-ghcr-pull and kind-dockerhub-pull
 * [2013dc3a5](https://github.com/kubeovn/kube-ovn/commit/2013dc3a54fc46538403c9ca82628a66635f4c79) chore: keep kind-install-enable-record-tunnel in kind.mk
 * [3ab9a1a63](https://github.com/kubeovn/kube-ovn/commit/3ab9a1a63a9dba7999fc69ddfa100caf1a0d06da) chore: keep kind-install-node-route-bgp-eip in kind.mk
 * [422e05779](https://github.com/kubeovn/kube-ovn/commit/422e05779e0301c1367c5fc5954d285293334a99) fix(lint): resolve G703 path traversal warning in cniserver
 * [24a42a5ca](https://github.com/kubeovn/kube-ovn/commit/24a42a5cafe381d3f696456b9b4d58c65813e7b5) fix(deps): update module golang.org/x/tools to v0.44.0 (#6613)
 * [b5a5cf7c8](https://github.com/kubeovn/kube-ovn/commit/b5a5cf7c8a8b1a27c4b470402fd8090a09e66ae2) fix(deps): update golang (#6607)
 * [a62647091](https://github.com/kubeovn/kube-ovn/commit/a62647091c4ccdb0c4396124a87d21ca3031aae1) fix(deps): update module golang.org/x/sys to v0.43.0 (#6597)
 * [4420d2e72](https://github.com/kubeovn/kube-ovn/commit/4420d2e72f9454192915ee7e698af6da7d4be5f2) chore(deps): update dependency go to v1.26.2 (#6587)
 * [52866dc76](https://github.com/kubeovn/kube-ovn/commit/52866dc762073df43110f6e33a9bf85a0c575fde) chore: go mod tidy
 * [1810fb139](https://github.com/kubeovn/kube-ovn/commit/1810fb1398becd888bea0206d7b32d49e6ae6301) deps: update k8s to v1.34.6 (#6579)
 * [008d91753](https://github.com/kubeovn/kube-ovn/commit/008d91753f15c87273d8b64228e10389f7eeaf07) fix(deps): update module google.golang.org/grpc to v1.79.3 [security] (#6477)
 * [79bce8e8c](https://github.com/kubeovn/kube-ovn/commit/79bce8e8cc2768a9c18cfd888e1e99b8cdf0af2e) fix(deps): update module golang.org/x/tools to v0.43.0 (#6432)
 * [b41edd7f5](https://github.com/kubeovn/kube-ovn/commit/b41edd7f5d0c7389f7fb0ec5aa712e5dfdb04fce) fix(deps): update golang (#6430)
 * [5dd9cbe44](https://github.com/kubeovn/kube-ovn/commit/5dd9cbe449fce31040982ff681b5fbbe1c4fd3c4) fix(deps): update golang (#6428)
 * [9d67c020f](https://github.com/kubeovn/kube-ovn/commit/9d67c020ff70d5771ba95abab9d9ee99cdfdc47c) fix(deps): update golang (#6408)
 * [26a8a97a2](https://github.com/kubeovn/kube-ovn/commit/26a8a97a2f63946c245d014e273c146e3c4cd5d8) chore(deps): update dependency go to v1.26.1 (release-1.15) (#6397)
 * [3876e6a3e](https://github.com/kubeovn/kube-ovn/commit/3876e6a3e2defacb08d1c1a042c241d15c7eb538) bump k8s to v1.34.5 (#6366)
 * [d1cf26755](https://github.com/kubeovn/kube-ovn/commit/d1cf267553013bd943386408c1712ead2c9cc023) fix(deps): update module golang.org/x/net to v0.51.0 (#6341)
 * [a55b011db](https://github.com/kubeovn/kube-ovn/commit/a55b011db61142465b13c079e13dd389b982339b) chore(deps): update dependency go to v1.26.0 (#6293)
 * [760aec4f4](https://github.com/kubeovn/kube-ovn/commit/760aec4f4f95ca59f333d175b471a0e0067dbefb) fix(deps): update golang (#6289)
 * [7949ec205](https://github.com/kubeovn/kube-ovn/commit/7949ec205dde4f6c9c7aac9d9075469270c2e421) update cert manager
 * [1fdebd01b](https://github.com/kubeovn/kube-ovn/commit/1fdebd01b8396a365e133e116e1c3a59473cc389) fix(deps): update module golang.org/x/sys to v0.41.0 (#6282)
 * [4ae4fc754](https://github.com/kubeovn/kube-ovn/commit/4ae4fc7542ed50a323c10791a25bd452d1f6bb64) chore(deps): update module golang.org/x/oauth2 to v0.35.0 (#6279)
 * [001b42cb4](https://github.com/kubeovn/kube-ovn/commit/001b42cb4237849d831371af3bc261024c8cc8b7) chore(deps): update dependency go to v1.25.7 (#6266)
 * [85cf19171](https://github.com/kubeovn/kube-ovn/commit/85cf19171bda5e403534a9d3850de01adc45a8db) fix: use ns/name
 * [96d357237](https://github.com/kubeovn/kube-ovn/commit/96d357237c22791543877400fb3ff2e5809f9f01) fix: preload hardcode vpc-nat-gateway:v1.14.25
 * [7845acfb4](https://github.com/kubeovn/kube-ovn/commit/7845acfb4b69e79f9d047b478fe0ec108a4ff50f) fix: nat gw e2e
 * [cc391e574](https://github.com/kubeovn/kube-ovn/commit/cc391e574f1ac9e99f56f751dfbd1574b01fa20b) refactor(metrics): remove unused dynamic cert cancel flow
 * [aba409b40](https://github.com/kubeovn/kube-ovn/commit/aba409b4038e3f39384f6d7227f4e2c87b830e33) fix(daemon): harden cni config installation path
 * [bbdf2deec](https://github.com/kubeovn/kube-ovn/commit/bbdf2deeccb1ef0542b52e36946370d468fb97a6) refactor(ovs): optimize lb ip port mapping updates
 * [11bc395d0](https://github.com/kubeovn/kube-ovn/commit/11bc395d0aa972a95dbfde2f5edab16f58db28fa) refactor(vpc): 优化子网列表获取以减少重复调用
 * [7c756e8c8](https://github.com/kubeovn/kube-ovn/commit/7c756e8c8f5062a868c61f74836e36f263af0153) refactor(controller): 移除冗余的队列添加逻辑
 * [0f4ed3099](https://github.com/kubeovn/kube-ovn/commit/0f4ed3099d96c2b85ecf478ab065d17aba030ca2) auto backfill eip ns
 * [5aea30474](https://github.com/kubeovn/kube-ovn/commit/5aea30474edd4991f8772855dbcf5f36f206cd12) fix(controller): prevalidate nat gw name and normalize long labels (#6604)
 * [1cd4e5182](https://github.com/kubeovn/kube-ovn/commit/1cd4e5182933f7b8f7ff250305b2bcbc0af94ef1) fix(ovs): rename loop-local protocol to portProto to avoid variable shadowing
 * [c5e5f0579](https://github.com/kubeovn/kube-ovn/commit/c5e5f057966faf2e7a40b49e4616c68dc65e6e69) fix(ovs): fix nil pointer dereference on port.Protocol in ACL generation
 * [a5e2b542b](https://github.com/kubeovn/kube-ovn/commit/a5e2b542b98f2ba8f5e5ccb458c77fadd3491c36) review: address PR review issues (P1/P2)
 * [eb4234b95](https://github.com/kubeovn/kube-ovn/commit/eb4234b958bf70e17cb23985ebc87ecaf9abdb08) fix(controller): remove unused perInterfaceIPAnnotationKey and perInterfaceMACAnnotationKey functions
 * [0480f332d](https://github.com/kubeovn/kube-ovn/commit/0480f332dd5eb148b2dde026368071b0fd11db4d) chore(deps): update dependency golangci-lint to v2.10.1 (master) (#6299)
 * [ac825c418](https://github.com/kubeovn/kube-ovn/commit/ac825c418d9c873a6366d364bc8fdfb00e10a508) feat: add Namespace field to VpcNatGateway and IptablesEIP (#6595)
 * [7b4ef6472](https://github.com/kubeovn/kube-ovn/commit/7b4ef6472a767b05fee1f89c6f507881d6de7c58) fix(controller): backfill nat gateway lan ip safely (#6424)
 * [c8617c84b](https://github.com/kubeovn/kube-ovn/commit/c8617c84bd508d5e1b64f9a30705b238db6bb271) Improve hairpin logic to support FIP, SNAT from any CIDR within the VPC and LoadBalancers (#6445)
 * [c15e31458](https://github.com/kubeovn/kube-ovn/commit/c15e31458f451545b68aeac617b962639b2623e2) allow any eip share external subnet gw in one (native) vlan (#6426)
 * [c376c40e5](https://github.com/kubeovn/kube-ovn/commit/c376c40e50bcdf25db29a731d1fc20652f61ed87) fix: FIP/DNAT/SNAT spec update handlers (#6412)
 * [716196ef9](https://github.com/kubeovn/kube-ovn/commit/716196ef95aecbb18f130ce69e5fe32f64dc8319) fix: add missing VlanInterfaces deep copy in ProviderNetworkSpec (#6257)
 * [c4ba3b927](https://github.com/kubeovn/kube-ovn/commit/c4ba3b92773950419a51802f0406648f77ff2e36) chore: align version to v1.15.9 (release-1.15 latest tag)
 * [5ff679f6d](https://github.com/kubeovn/kube-ovn/commit/5ff679f6dc87b8025605abcc36c5b800c7141e20) feat: support custom resources and bandwidth limits for VpcEgressGateway (#6407)
 * [d726f08a7](https://github.com/kubeovn/kube-ovn/commit/d726f08a7159df213e58ce0065fcfeaad34fc743) feat: introduce serviceMonitor and labels to services (#6340)
 * [82ca82484](https://github.com/kubeovn/kube-ovn/commit/82ca82484a6998a5e6a52c04bcd128f1181264bf) feat(helm):Add affinity and nodeSelector support for ovs-ovn and ovs-ovn-dpdk DaemonSets using v2 chart (#6308)
 * [553f75467](https://github.com/kubeovn/kube-ovn/commit/553f7546792151ddac3033c3af5a0fc0ea5f311f) refactor(controller): clean up redundant logic in vpc controller (#6295)
 * [bea7941ea](https://github.com/kubeovn/kube-ovn/commit/bea7941eaf6f3039ed7ddde386f9473ebe4dad00) refactor(ovn_leader_checker): improve readability and naming (#6263)
 * [0e7b8f1ff](https://github.com/kubeovn/kube-ovn/commit/0e7b8f1fffef7cae68ceb85011d15a25fc2be073) fix(controller): remove duplicate addOrUpdateSubnetQueue.Add in handleAddOrUpdateVpc
 * [087613d5a](https://github.com/kubeovn/kube-ovn/commit/087613d5a88d74c30214504ac6cad76193572af6) fix(controller): correct ExternalIDs filter logic in diffPolicyRouteWithLogical
 * [c13ca58bc](https://github.com/kubeovn/kube-ovn/commit/c13ca58bcd2254cc85727947eb578fef4ff2edec) chore(deps): upgrade ginkgo to v2.28.1 and gomega to v1.39.0
 * [e510a5634](https://github.com/kubeovn/kube-ovn/commit/e510a5634783686afa7458a975fa676d8812ed63) perf: only enqueue subnet when not ready (#6268)
 * [926bfec9b](https://github.com/kubeovn/kube-ovn/commit/926bfec9bafd68fb7e3c62485612ccc780b41dc2) security: specify a storage limit for containers (#6259)
 * [c9b2ad0dd](https://github.com/kubeovn/kube-ovn/commit/c9b2ad0dddcce536a50613c3c1ed516b6c33859c) fix(helm): make v2 helm chart run upgrade logic within ovn (#6410)
 * [88b1e2767](https://github.com/kubeovn/kube-ovn/commit/88b1e27674ba5df848e587ca383c54eab730b63a) Set Kubevirt live-migrations multi-chassis options for all VM nics (#6241)
 * [33ab615de](https://github.com/kubeovn/kube-ovn/commit/33ab615de32dd78dffa0c34d6c1a2a15f8e46935) cni-server: fix vlan id reference (#6365)
 * [1891ab1cd](https://github.com/kubeovn/kube-ovn/commit/1891ab1cd1c695df496ede95a174379fc7b063fb) fix(ovn): update patches for compatibility with upstream OVN branch-25.03 (#6331)
 * [1d284f54a](https://github.com/kubeovn/kube-ovn/commit/1d284f54aa6d6dcf182789557f488e2fb2542e98) controller: fix stale ip cr after subnet deletion (#6275)
 * [9e8b7a752](https://github.com/kubeovn/kube-ovn/commit/9e8b7a752fdeaf2313b0041307d966e08a8d396e) cni-server: add static fdb entry for subnets with u2o enabled (#6269)
 * [9143f0c1d](https://github.com/kubeovn/kube-ovn/commit/9143f0c1d082934dd4c30a4b3bc2d13271b115de) controller: fix ips owned by subnet u2o ip and mcast querier ip deleted in gc (#6132)
 * [751dcc0e7](https://github.com/kubeovn/kube-ovn/commit/751dcc0e788886336e09b1a4e97f55d16819405f) fix(chart-v2): fix extraEnvs indent in controller-deployment template (#6575)
 * [e1699f9d9](https://github.com/kubeovn/kube-ovn/commit/e1699f9d9fcbaf41037dc567f56d68ddf65f8cb7) fix(controller): preserve VPC policy route next hops during OVN diff (#6577)
 * [c43b88faf](https://github.com/kubeovn/kube-ovn/commit/c43b88faf90cd261452ccb421de944456f9fb83f) fix(controller): reuse addPolicyRouteToVpc to handle empty NextHopIP correctly (#6472)
 * [d8a2f8a4f](https://github.com/kubeovn/kube-ovn/commit/d8a2f8a4fdd4187857c080bdce6139457ed86950) Split ECMP next-hops in custom VPC policy routes (#6405)
 * [e66df8295](https://github.com/kubeovn/kube-ovn/commit/e66df8295352f1aee8425275d1c6cd9192819298) fix(deps): update module google.golang.org/grpc to v1.79.3 [security] (#6477)
 * [280bd6f91](https://github.com/kubeovn/kube-ovn/commit/280bd6f9151a050a17e524f9f49996ba56413025) fix: add ipv6_src to LB selection_fields for IPv6 session affinity (#6422)
 * [e64262402](https://github.com/kubeovn/kube-ovn/commit/e6426240299c32f4e92a790cf70aa284e7d96df5) fix: scope ipBlock.except to its own CIDR in NetworkPolicy ACL rules (#6419)
 * [5b40e3866](https://github.com/kubeovn/kube-ovn/commit/5b40e386695e8b63151e75bfc5ffd2632e6022b6) fix: reconcile subnet and localnet network_name on vlan provider change (#6393)
 * [be3869da0](https://github.com/kubeovn/kube-ovn/commit/be3869da03bb3e4161e6c6a3c6f84ef0d40e2d2e) fix: scope SLR deletion to its own VPC to avoid deleting VIPs across VPCs (#6391)
 * [2d2e098bb](https://github.com/kubeovn/kube-ovn/commit/2d2e098bb693dc0d377d59b04f5cd10b84faa979) fix: prevent subnet from getting permanently stuck when VLAN is not ready (#6352)
 * [6eab580e1](https://github.com/kubeovn/kube-ovn/commit/6eab580e184fa5d10d3ef4bb36bffc63621722c3) dist: add 25.03 ACL compatibility for OVN patch (#6463)
 * [09fdbcb37](https://github.com/kubeovn/kube-ovn/commit/09fdbcb37e66f1ecf52f959b3b0e746c5c837229) fix: BGP speaker reads IPs from pod annotations instead of pod status (#6344)
 * [7760a2958](https://github.com/kubeovn/kube-ovn/commit/7760a2958639e4d768ad589dbc8c8d0048267876) fix: use targeted patch in calcSubnetStatusIP to prevent U2O status overwrite (#6350)
 * [a8211e828](https://github.com/kubeovn/kube-ovn/commit/a8211e8280a24f1e6ace1e0c78b1f8529a2c9cb0) fix(daemon): handle tombstone objects in subnet and service delete handlers (#6470)
 * [6bb753ba9](https://github.com/kubeovn/kube-ovn/commit/6bb753ba920479087f72ffd7be7462656c676aee) fix(daemon): sort encap IP list for deterministic ovs-vsctl output (#6473)
 * [99966b409](https://github.com/kubeovn/kube-ovn/commit/99966b4095e07c3f282a6d44f88e053682288bac) fix(ipam): release allocated static IPs on partial allocation failure (#6493)
 * [9343d503c](https://github.com/kubeovn/kube-ovn/commit/9343d503c8b71911d3fc8bf5a2fd0d53830b71a1) fix(controller): handle restartable init containers in DeleteNamedPortByPod (#6486)
 * [8722427b5](https://github.com/kubeovn/kube-ovn/commit/8722427b567cffec61c249c1e7572df409153389) fix(controller): use exact match instead of substring match for VIP lookup in syncVirtualPort (#6491)
 * [9a0b422ae](https://github.com/kubeovn/kube-ovn/commit/9a0b422ae918ce9a4f5755c7bf5146a06dbcdc73) fix(controller): fix variable shadowing in cleanDuplicatedChassis (#6496)
 * [1ff1a09bd](https://github.com/kubeovn/kube-ovn/commit/1ff1a09bd0f6e3284133ccd83a87d726c8f1e8ac) fix(controller): prevent dual-stack cross-deletion of local DNS cache policy routes (#6497)
 * [dad2e0df9](https://github.com/kubeovn/kube-ovn/commit/dad2e0df9b11703cd351bba83f5c1449fb268ae7) fix(iptables-eip): look up subnet directly by name instead of resolving through NAD (#6287)
 * [3a0bdcb04](https://github.com/kubeovn/kube-ovn/commit/3a0bdcb04509474d73a43c772c835b40cb868617) fix(netpol): use parenthesis around ipblock expressions
 * [772fc80d4](https://github.com/kubeovn/kube-ovn/commit/772fc80d4a907ddfa520ff97dad3b44e7ee3ab1b) fix(ovn): add null check before sbrec delete calls to prevent crashes (#6316)
 * [de8cbc55a](https://github.com/kubeovn/kube-ovn/commit/de8cbc55a2fb403a21f20323856112c8c7407da5) fix(controller): prevent NPE in isOvnSubnet when subnet is nil (#6312)
 * [b02b547b4](https://github.com/kubeovn/kube-ovn/commit/b02b547b4447017f75664677122f549e37480ff0) feat(ovs): add OvnDatabaseControl helper for OVN nb/sb database ctl commands
 * [71a0bd282](https://github.com/kubeovn/kube-ovn/commit/71a0bd28213d1c3edb67da7e67912ba95d4b5725) fix(ha): detect and recover from split-brain leader with incomplete cluster membership (#6363)
 * [26adcbfb4](https://github.com/kubeovn/kube-ovn/commit/26adcbfb41e6a2f70db57e0f119f35f55e6d79ca) fix: check northd endpoint by pod IP protocol in ovn leader checker (#6523)
 * [10cb86f97](https://github.com/kubeovn/kube-ovn/commit/10cb86f97ed48232d9094ff29bc034f68076d935) fix(slr): handle orphaned ip_port_mappings and overwrites (#6558)
 * [49301baa6](https://github.com/kubeovn/kube-ovn/commit/49301baa6293dec4a0758964d8f7fe7028bfd5c2) fix: wrap EIP error shadowing and use uint64 for iptables counters (#6538)
 * [42ebc8270](https://github.com/kubeovn/kube-ovn/commit/42ebc8270d1837f396664fa5db8956fbc55e3ea0) fix(daemon): skip invalid policy routing rules when pod lacks matching protocol IP (#6522)
 * [5eca939e0](https://github.com/kubeovn/kube-ovn/commit/5eca939e03fc01d985a5ecb24ba5517f1ad5b9ad) fix(daemon): harden protocol iteration and bounds checks in gateway (#6547)
 * [1c90aed5c](https://github.com/kubeovn/kube-ovn/commit/1c90aed5c0237c3ce368fa095577c9b0086037ff) fix(cherry-pick): resolve #6536 conflicts on release-1.15
 * [1d3b7e0bb](https://github.com/kubeovn/kube-ovn/commit/1d3b7e0bbf3044286cce99aa43b6aa4593af686d) fix(controller): prevent invalid ACL generation for named ports and IPBlock except CIDRs (#6536)
 * [35e6dd446](https://github.com/kubeovn/kube-ovn/commit/35e6dd4467e5df8b337229113b7f2a9b0c5eb83e) fix: protect ANP/BNP/CNP priority maps from concurrent access (#6373)
 * [e81442f19](https://github.com/kubeovn/kube-ovn/commit/e81442f19814650fdfdf7bfd0d5532d67b593002) fix(controller): use defer to protect vpcKeyMutex unlock in subnet handler (#6515)
 * [0930df475](https://github.com/kubeovn/kube-ovn/commit/0930df475c03cf8bc023016b51b75b92c0948fb2) fix(informer): return nil instead of error in VM indexer to prevent client-go panic (#6537)
 * [7296acea6](https://github.com/kubeovn/kube-ovn/commit/7296acea66e809877bec28bcac1984c1b483ad02) fix(controller): fix nil deref, cache mutation, and early return bugs in subnet controller (#6490)
 * [226ffcca1](https://github.com/kubeovn/kube-ovn/commit/226ffcca152bf03e5538fb21061fa74df10bf27a) fix: old go code
 * [2706e6462](https://github.com/kubeovn/kube-ovn/commit/2706e6462355ec29a769010313fbafb7cfc3a52a) fix(controller): prevent panics from nil deref and data races in pod controller (#6483)
 * [359e24c5d](https://github.com/kubeovn/kube-ovn/commit/359e24c5d54b26b47d3dc238f11135cfb1a053aa) ci(workflow): 修正导出错误日志的条件
 * [ffbb769c7](https://github.com/kubeovn/kube-ovn/commit/ffbb769c75c0e5de3892e330083bc3974f4c2315) fix(ci): 修复CI脚本中日志获取失败时仍执行grep的问题
 * [7c81e0c7a](https://github.com/kubeovn/kube-ovn/commit/7c81e0c7af54a59c59e483f14558a1154d8acde1) ci: 根据测试结果有条件导出组件错误日志
 * [f98c47e3d](https://github.com/kubeovn/kube-ovn/commit/f98c47e3d29380726a5721d73db3b3de6117904b) ci: 改进CI脚本中的日志过滤条件
 * [2ce87a095](https://github.com/kubeovn/kube-ovn/commit/2ce87a0952bd74094583e2c8896a9dd5b9ebb3bb) fix(ci): 改进 CI 脚本中的日志收集逻辑
 * [7c79a1ad5](https://github.com/kubeovn/kube-ovn/commit/7c79a1ad548d3d7ffda2e78f4cab29e676798eb9) ci: 在检查脚本失败时导出组件错误日志
 * [451e6b111](https://github.com/kubeovn/kube-ovn/commit/451e6b11131fce725810e19521257d3b97bf7e95) fix as review
 * [9927c9ff0](https://github.com/kubeovn/kube-ovn/commit/9927c9ff083fb4dd6ba35706aa01a460be036f05) fix as review
 * [67c3b2a85](https://github.com/kubeovn/kube-ovn/commit/67c3b2a85a75802e1706e34b441d184248f4be2c) fix upload-artifact@v8
 * [55fb2917b](https://github.com/kubeovn/kube-ovn/commit/55fb2917b332a86eb8622a6647f16bbf431ce6de) fix as review
 * [039105e14](https://github.com/kubeovn/kube-ovn/commit/039105e142e95750514ed812eea66c926975e72b) simple
 * [77cbbf441](https://github.com/kubeovn/kube-ovn/commit/77cbbf441364619c8913aa31b995efe57b91f408) fix:kind image cache
 * [5d38f01f8](https://github.com/kubeovn/kube-ovn/commit/5d38f01f887fa0c7f71641978cf868103682a187) fix(ci): pull kindest-node from Docker Hub directly, no GHCR mirror needed
 * [7d873d357](https://github.com/kubeovn/kube-ovn/commit/7d873d3571bb4fd939a4ed363588addd8791391f) push only when not exist
 * [035d657b2](https://github.com/kubeovn/kube-ovn/commit/035d657b2fc7a759c30a1c042fed590f632e0efb) fix(ci): use user-controlled GHCR mirror for kindest-node to fix image pull failure
 * [fef09607d](https://github.com/kubeovn/kube-ovn/commit/fef09607dd9c63933a59a974861e6b4b0a0335a7) use master v1.35.1
 * [1b4735a0a](https://github.com/kubeovn/kube-ovn/commit/1b4735a0a3f0287ac75e0e9c04f65c59172e8ade) fail fast: if pull kind-node failed
 * [85005412a](https://github.com/kubeovn/kube-ovn/commit/85005412a2b4aaa4db1202fbe88e5cbd49a11700) use my github token
 * [e23811a4b](https://github.com/kubeovn/kube-ovn/commit/e23811a4b2a0b76b962e190b7a0bb04895730eb4) refine pr e2e workflow suites
 * [764b2612d](https://github.com/kubeovn/kube-ovn/commit/764b2612d287e17db1aa673b58edacaaf0db7b2f) more e2e
 * [6e1be00b7](https://github.com/kubeovn/kube-ovn/commit/6e1be00b73a6ec8e21e6a7b3b7fbc5a3045b69aa) add pr e2e workflow for 1.15 branch
 * [dde3972f1](https://github.com/kubeovn/kube-ovn/commit/dde3972f11bdcd2d09a1aabbf852c0078176ea25) feat: allow vpc nat gw image use Ubuntu 24.04
 * [3b558e3d0](https://github.com/kubeovn/kube-ovn/commit/3b558e3d04bfe328cdea61a647b25021d52ac113) fix: skip IPAM init for empty IP annotation and handle nil IPRangeList (#6296)
 * [d0963e4cb](https://github.com/kubeovn/kube-ovn/commit/d0963e4cbd3801c2fa0933fee6222188a96178ff) feat: add user-defined annotations support for nat gw sts template creation
 * [0bb4809d0](https://github.com/kubeovn/kube-ovn/commit/0bb4809d0ee2777c12e41d8e23904777d3266ca5) fix: Ensure Kubevirt optimizations do not use old migration state (#6239)
 * [3bb38e9bb](https://github.com/kubeovn/kube-ovn/commit/3bb38e9bb932f0e39c4d2b2d1b15e0602ccf6bc0) fix localnet port should not be static "patch-localnet." (#6244)
 * [0f99cb299](https://github.com/kubeovn/kube-ovn/commit/0f99cb2996734d35dd8d57f255951ca387831fe2) rbac: replace wildcard with a clear list of allowed verbs (#6233)
 * [9359ee849](https://github.com/kubeovn/kube-ovn/commit/9359ee849bf9770942c8b0e13cf0f74b9cd06530) remove duplicate description (#6234)
 * [6ec6e8fee](https://github.com/kubeovn/kube-ovn/commit/6ec6e8feef7bf93ae98c96d70d03158d5871b9b6) Add vendor ID to transit switches (#6186)
 * [32cea8504](https://github.com/kubeovn/kube-ovn/commit/32cea85044284e6efbd7ce46f2a12fed7db3f731) fix eip deletion not trigger qos again
 * [364f2e310](https://github.com/kubeovn/kube-ovn/commit/364f2e310824f1bd7946c4cf45b4c56f780c4453) go fmt
 * [58cb1e712](https://github.com/kubeovn/kube-ovn/commit/58cb1e712e0339d92580b8f3ca8d8e84184a63db) fix(qos): fix NAT Gateway QoS not working
 * [95ad29a08](https://github.com/kubeovn/kube-ovn/commit/95ad29a089ff0b622a2ec095d4e735200a8cd8c7) fix: skip NetworkUnavailable condition in non-primary CNI mode (#6195)
 * [af40be64b](https://github.com/kubeovn/kube-ovn/commit/af40be64b22428a93339dddf63f3997089cf2bf8) fix(vpcnatgw): nat gateway uses faulty default network (#6212)
 * [75d54d368](https://github.com/kubeovn/kube-ovn/commit/75d54d3680247459a5d36d4a658a21b60f2ffa27) controller: fix vpc egress gateway forward policy not working for targets within the internal subnet's cidr (#6222)
 * [a82289a91](https://github.com/kubeovn/kube-ovn/commit/a82289a91741335c97587ea7394c9ddb4f1ba281) controller: fix vpc egress gateway not working for targets within the internal subnet's cidr (#6218)
 * [437dcb258](https://github.com/kubeovn/kube-ovn/commit/437dcb258716d809aa757d49672a1d1f9ef7551c) Update start-ic-db.sh to have gen_conn_addr function defined (#6216)
 * [8da2d6176](https://github.com/kubeovn/kube-ovn/commit/8da2d6176a279fa1c36d91ed28e1c3804713915f) fix: backport lost util.TrimManagedFields
 * [3104d1d92](https://github.com/kubeovn/kube-ovn/commit/3104d1d92de3430d8d89318d13e19f37ed747fe7) subnet add macvlan type label
 * [6b573bb79](https://github.com/kubeovn/kube-ovn/commit/6b573bb79638b2b294041d6ac9c6eb30137b36fa) fix: cni eip cache mem leak
 * [264a2d933](https://github.com/kubeovn/kube-ovn/commit/264a2d9332b36a082632ec12de5a2104a6344859) fix: vpc nat gw img add netcat-openbsd
 * [5fd0ecdf3](https://github.com/kubeovn/kube-ovn/commit/5fd0ecdf3227c93084f44a5fd99078bc1ae505d9) feat: allow vpc nat gw pod mount its scripts form host
 * [fab4571ad](https://github.com/kubeovn/kube-ovn/commit/fab4571ad42fda9e81078918704d84b8e875d2cd) allow to observe pod tun | vni id
 * [43764c31b](https://github.com/kubeovn/kube-ovn/commit/43764c31b8ff3ac21598037c72be93910be75270) allow bgp via route of node local access eip
 * [92037a399](https://github.com/kubeovn/kube-ovn/commit/92037a39992a5ac24c1dba2a0556c654ad580cc9) feat(controller): add VPC NAT gateway and iptables EIP handling in daemon controller
 * [bf7371f37](https://github.com/kubeovn/kube-ovn/commit/bf7371f372ba0f9e8d7bba6bda8f13ae60eeb593) fix: caching NAD CRD should before all kubeovn crds and pod (#6198)
 * [2a6aab902](https://github.com/kubeovn/kube-ovn/commit/2a6aab902b9f538a9203ad0a7a9aee0b11a20090) fix 1.15 base image build
 * [dc57e83c8](https://github.com/kubeovn/kube-ovn/commit/dc57e83c894bc5469e14edc1e95a837bbb87430a) fix(vpc-nat-gw): use consistent provider for all pod annotations (#6204)
 * [995bed3e4](https://github.com/kubeovn/kube-ovn/commit/995bed3e42e6f6648ae15dee31fb192a7eebcdb2) fix(vpcnatgw): cannot DNAT same EIP and same external port on two different protocols (#6201)
 * [a41585dfd](https://github.com/kubeovn/kube-ovn/commit/a41585dfd7dc1a803fb707f25901482cb2e90167) fix: cleanup auto-created vlan subinterfaces on interface switch (#6193)

### Contributors

 * Asish Kumar
 * Eli Treuherz
 * ILYAS ABDELLAOUI
 * Johann Schley
 * Juan Morete
 * Mengxin Liu
 * SKALA NETWORKS
 * SkalaNetworks
 * Uğur Tafralı
 * Zespre Chang
 * andrewlee1089
 * changluyi
 * clyi
 * dnugmanov
 * jimyag
 * jskrill
 * renovate[bot]
 * zbb88888
 * 张祖建

## v1.15.10 (2026-04-10)

 * [a0e8fb099](https://github.com/kubeovn/kube-ovn/commit/a0e8fb099bcfd86a1be85b82bbb692c106a5142b) release v1.15.10
 * [0369b997b](https://github.com/kubeovn/kube-ovn/commit/0369b997b795c121312cb1dad23b7383afdd4839) Fix/upgrade 1.15 error (#6615)
 * [7e4b1ab31](https://github.com/kubeovn/kube-ovn/commit/7e4b1ab310abe3f70f6154a151648d153116185c) fix(deps): update module golang.org/x/tools to v0.44.0 (#6613)
 * [e9fdb8343](https://github.com/kubeovn/kube-ovn/commit/e9fdb83430611f22e63bf0e18d4b59eb99a157f1) fix(deps): update golang (#6607)
 * [02afa6e45](https://github.com/kubeovn/kube-ovn/commit/02afa6e45f68be0facc5a009fe6bdb84a8765bf2) fix(deps): update module golang.org/x/sys to v0.43.0 (#6597)
 * [cfd40d97a](https://github.com/kubeovn/kube-ovn/commit/cfd40d97ad9ca818f9ed83ca71cde29038b21972) prepare for next release

### Contributors

 * Mengxin Liu
 * changluyi
 * renovate[bot]

## v1.15.9 (2026-04-08)

 * [64eb81f53](https://github.com/kubeovn/kube-ovn/commit/64eb81f53816c54ab025650461c3c8655e31cbbc) release v1.15.9
 * [854164e4b](https://github.com/kubeovn/kube-ovn/commit/854164e4b88a02c9eae1a71becdeae7c558f6180) chore(deps): update dependency go to v1.26.2 (#6587)
 * [3a2264795](https://github.com/kubeovn/kube-ovn/commit/3a2264795f5229f3334e04d205c41fc985bf469d) fix(controller): preserve VPC policy route next hops during OVN diff (#6577)
 * [541993fd7](https://github.com/kubeovn/kube-ovn/commit/541993fd7584b2c277e827f3a2c92697cc00a05f) fix(controller): reuse addPolicyRouteToVpc to handle empty NextHopIP correctly (#6472)
 * [100933162](https://github.com/kubeovn/kube-ovn/commit/100933162e1ecc754d28da10fbe05d7016e90103) Split ECMP next-hops in custom VPC policy routes (#6405)
 * [1099892d2](https://github.com/kubeovn/kube-ovn/commit/1099892d2df24c49e0e46795c068986ff23f98e0) deps: update k8s to v1.34.6 (#6579)
 * [09a9e2316](https://github.com/kubeovn/kube-ovn/commit/09a9e2316f75bfe0fbab0ab91401ab5b996d40f4) fix(chart-v2): fix extraEnvs indent in controller-deployment template (#6575)
 * [72b6b3174](https://github.com/kubeovn/kube-ovn/commit/72b6b3174881eade15948c9d94c07c9356c6acae) ci: only fail on pod restarts caused by panic or SIGSEGV (#6506) (#6568)
 * [662100f7f](https://github.com/kubeovn/kube-ovn/commit/662100f7f8a3d51a5bc53df4d24ef78890fc5bca) fix(slr): handle orphaned ip_port_mappings and overwrites (#6558)
 * [8ff5589a7](https://github.com/kubeovn/kube-ovn/commit/8ff5589a709f5949cc588883e3581ba640c8739e) fix(ha): adapt OvnDatabaseControl call for release-1.15 compatibility
 * [ca652669a](https://github.com/kubeovn/kube-ovn/commit/ca652669a39e59864449d20a80413c87f013d884) fix: wrap EIP error shadowing and use uint64 for iptables counters (#6538)
 * [b54aa2140](https://github.com/kubeovn/kube-ovn/commit/b54aa2140294d439e12910f3dd00c43cfcf8fca9) fix(ha): detect and recover from split-brain leader with incomplete cluster membership (#6363)
 * [9a164334a](https://github.com/kubeovn/kube-ovn/commit/9a164334aeef5c5d760c65d5136a67e50851fa2a) fix(daemon): skip invalid policy routing rules when pod lacks matching protocol IP (#6522)
 * [9ca482dfc](https://github.com/kubeovn/kube-ovn/commit/9ca482dfc4187b52468def1e2be543ebe45c217e) fix(controller): prevent invalid ACL generation for named ports and IPBlock except CIDRs (#6536)
 * [3e7c07914](https://github.com/kubeovn/kube-ovn/commit/3e7c0791450f6a15e8c1fa8098579791f46baa7e) fix(daemon): harden protocol iteration and bounds checks in gateway (#6547)
 * [51c9bca8b](https://github.com/kubeovn/kube-ovn/commit/51c9bca8bb233d8174f64ab65b1f0cc105ff1247) fix: protect ANP/BNP/CNP priority maps from concurrent access (#6373)
 * [7daebcc27](https://github.com/kubeovn/kube-ovn/commit/7daebcc271d6387f68410f1328294733ef913ac0) fix(controller): use defer to protect vpcKeyMutex unlock in subnet handler (#6515)
 * [1a9acfa50](https://github.com/kubeovn/kube-ovn/commit/1a9acfa50d2b0a55f4f16b45722af61d1b781258) fix(informer): return nil instead of error in VM indexer to prevent client-go panic (#6537)
 * [f46c0b180](https://github.com/kubeovn/kube-ovn/commit/f46c0b1800dc9311f355330ed486b3e426c2285b) ci: install multus before kubevirt e2e tests
 * [5f2dec037](https://github.com/kubeovn/kube-ovn/commit/5f2dec037cad2e861fbb2333ddc92d20912fa5e3) Set Kubevirt live-migrations multi-chassis options for all VM nics (#6241)
 * [38965a6fe](https://github.com/kubeovn/kube-ovn/commit/38965a6fee0f2b42e4eb518ba1aad0cf9e509b76) prepare for next release

### Contributors

 * Asish Kumar
 * ILYAS ABDELLAOUI
 * Mengxin Liu
 * SKALA NETWORKS
 * Uğur Tafralı
 * jskrill
 * renovate[bot]
 * 张祖建

## v1.15.8 (2026-03-26)

 * [0060dd23c](https://github.com/kubeovn/kube-ovn/commit/0060dd23c36f26e52517b4808570b5bcd78f19d2) release v1.15.8
 * [d1fd99f50](https://github.com/kubeovn/kube-ovn/commit/d1fd99f505af92be811073dce4201c8ad6652749) fix: check northd endpoint by pod IP protocol in ovn leader checker (#6523)
 * [eeafcfdfd](https://github.com/kubeovn/kube-ovn/commit/eeafcfdfd918a8d3a59164a846f52ca6e78682a1) fix(netpol): use parenthesis around ipblock expressions
 * [5c7bd2852](https://github.com/kubeovn/kube-ovn/commit/5c7bd28524b2d7fb771fadef20a683c97409aa5e) fix: resolve unused imports and missing errors package in node_test.go
 * [4cbcae30f](https://github.com/kubeovn/kube-ovn/commit/4cbcae30f611647530d588a39e38cca17fdce7e4) Fix misleading log message on VPC Policy Route addition. (#6464)
 * [f7bdbea82](https://github.com/kubeovn/kube-ovn/commit/f7bdbea8212820d4125f246be3a068fe6c4ca7cf) fix(daemon): handle tombstone objects in subnet and service delete handlers (#6470)
 * [3b4479f73](https://github.com/kubeovn/kube-ovn/commit/3b4479f733cd6cb4580c9e049c71c0c3bf2d3bc5) fix(daemon): sort encap IP list for deterministic ovs-vsctl output (#6473)
 * [8952152f2](https://github.com/kubeovn/kube-ovn/commit/8952152f2f7238e64ebbad943cb9f689748ae729) fix(controller): prevent dual-stack cross-deletion of local DNS cache policy routes (#6497)
 * [e424a0f88](https://github.com/kubeovn/kube-ovn/commit/e424a0f88b358286187305a4bef3aafa9cce56aa) fix(controller): fix variable shadowing in cleanDuplicatedChassis (#6496)
 * [33d44c5a2](https://github.com/kubeovn/kube-ovn/commit/33d44c5a2160aefd163b9c5cfd6202bb9dc29c14) fix(controller): use exact match instead of substring match for VIP lookup in syncVirtualPort (#6491)
 * [f6cdf479e](https://github.com/kubeovn/kube-ovn/commit/f6cdf479e5286ac26740b278ca1153e8eadb5ef9) fix(controller): fix nil deref, cache mutation, and early return bugs in subnet controller (#6490)
 * [183db4685](https://github.com/kubeovn/kube-ovn/commit/183db468558f51276fed244800aebaf674637a1e) fix(controller): handle restartable init containers in DeleteNamedPortByPod (#6486)
 * [7bd041833](https://github.com/kubeovn/kube-ovn/commit/7bd041833855961573bb68c80bf33c1f8c1e4ab9) fix(controller): prevent panics from nil deref and data races in pod controller (#6483)
 * [e4fab4cb0](https://github.com/kubeovn/kube-ovn/commit/e4fab4cb0568be5136328912b6d062a4fdafb5d5) prepare for next release

### Contributors

 * Mengxin Liu
 * SkalaNetworks
 * andrewlee1089
 * changluyi

## v1.15.7 (2026-03-23)

 * [e7bd69c77](https://github.com/kubeovn/kube-ovn/commit/e7bd69c77f31cc0766e01401cfafb2e2831af503) release v1.15.7
 * [3ca483d00](https://github.com/kubeovn/kube-ovn/commit/3ca483d0055659ee2007135a2636a1334bc76f7f) fix(deps): update module google.golang.org/grpc to v1.79.3 [security] (#6477)
 * [6812bf03a](https://github.com/kubeovn/kube-ovn/commit/6812bf03a0de400b2b0204c5976d53142cebfc3d) fix(ipam): release allocated static IPs on partial allocation failure (#6493)
 * [18247bbf1](https://github.com/kubeovn/kube-ovn/commit/18247bbf1b8da9648a4de5dce415048af1e626ab) prepare for next release

### Contributors

 * Mengxin Liu
 * renovate[bot]
 * 张祖建

## v1.15.6 (2026-03-18)

 * [b32aee09f](https://github.com/kubeovn/kube-ovn/commit/b32aee09f65faf2fc38f13f1fbcb61fff0cbfb62) release v1.15.6
 * [a31ba277b](https://github.com/kubeovn/kube-ovn/commit/a31ba277bd4632a8e8d14f8c9e302fac77449e5e) dist: add 25.03 ACL compatibility for OVN patch (#6463)
 * [f3dca4b60](https://github.com/kubeovn/kube-ovn/commit/f3dca4b60a911ca79463583922d68e9d83958800) fix: BGP speaker reads IPs from pod annotations instead of pod status (#6344)
 * [a01efb4ed](https://github.com/kubeovn/kube-ovn/commit/a01efb4ed9bb277866e9c3af8a2fd64ba746daa2) kind: enable OwnerReferencesPermissionEnforcement (#6456)
 * [abb302d50](https://github.com/kubeovn/kube-ovn/commit/abb302d50efceeaa14d77ab603eacb68a5702060) controller: do not set blockOwnerDeletion for ip resources
 * [cdcd7034d](https://github.com/kubeovn/kube-ovn/commit/cdcd7034dd9dd7c356c4c1a06b3100d0650b547c) prepare for next release

### Contributors

 * Mengxin Liu
 * changluyi
 * zhangzujian

## v1.15.5 (2026-03-12)

 * [77eb581ff](https://github.com/kubeovn/kube-ovn/commit/77eb581ff6291e56184b3f14536a1833932cbd1c) release v1.15.5
 * [f7f135314](https://github.com/kubeovn/kube-ovn/commit/f7f135314297ac3c1a3cac21568df0790b99a69f) fix(deps): update module golang.org/x/tools to v0.43.0 (#6432)
 * [a0d787433](https://github.com/kubeovn/kube-ovn/commit/a0d787433f0e740b75d986849833815a0d929bc5) fix: use targeted patch in calcSubnetStatusIP to prevent U2O status overwrite (#6350)
 * [0fda9a69a](https://github.com/kubeovn/kube-ovn/commit/0fda9a69a3464f7b89ba3f7bd2b15049332efcf6) fix(deps): update golang (#6430)
 * [99876f673](https://github.com/kubeovn/kube-ovn/commit/99876f673af736523edda40c363d92a36a76de47) fix(helm): make v2 helm chart run upgrade logic within ovn (#6410)
 * [883121029](https://github.com/kubeovn/kube-ovn/commit/883121029aecf332ba05b05e8098d7043286e825) fix: add ipv6_src to LB selection_fields for IPv6 session affinity (#6422)
 * [9c89b8e9d](https://github.com/kubeovn/kube-ovn/commit/9c89b8e9d539623a3c2f6d142d6fa82c5707291a) fix: scope ipBlock.except to its own CIDR in NetworkPolicy ACL rules (#6419)
 * [7a6460130](https://github.com/kubeovn/kube-ovn/commit/7a64601308b0c7a7c72fd436c8079f2b8642050e) fix(deps): update golang (#6428)
 * [bacf2d842](https://github.com/kubeovn/kube-ovn/commit/bacf2d84254fbc35446852ed4c97677c2e79c033) fix: reconcile subnet and localnet network_name on vlan provider change (#6393)
 * [bcf256532](https://github.com/kubeovn/kube-ovn/commit/bcf256532061ff1d75fbeec1e45ce263f6726080) feat: support custom resources and bandwidth limits for VpcEgressGateway (#6407)
 * [40bfa3432](https://github.com/kubeovn/kube-ovn/commit/40bfa3432652f9cb0219c5ff83c8300de3a68d14) fix: scope SLR deletion to its own VPC to avoid deleting VIPs across VPCs (#6391)
 * [bdf022348](https://github.com/kubeovn/kube-ovn/commit/bdf022348e9e0ff0e0f9b5ede1f0a28523c6813b) chore(deps): update aquasecurity/trivy-action action to v0.35.0 (#6404)
 * [08fbad4f3](https://github.com/kubeovn/kube-ovn/commit/08fbad4f3d1520e1a8bd295e3b28e7cee6317752) fix(deps): update golang (#6408)
 * [02266c08f](https://github.com/kubeovn/kube-ovn/commit/02266c08f65e2d5ebed1d2f1946108c4716dcaf4) update trivy for arm build
 * [e1274ab2c](https://github.com/kubeovn/kube-ovn/commit/e1274ab2c11ccea2f342c9cafac09e2c4a15a0ef) chore(deps): update dependency go to v1.26.1 (release-1.15) (#6397)
 * [3bf23ab6a](https://github.com/kubeovn/kube-ovn/commit/3bf23ab6a0b2f665e55bb415e2f5fe78af606e43) feat(helm):Add affinity and nodeSelector support for ovs-ovn and ovs-ovn-dpdk DaemonSets using v2 chart (#6308)
 * [1d130e5b2](https://github.com/kubeovn/kube-ovn/commit/1d130e5b23d3cf7123c8ed91fae07bd57a1b2cf5) bump k8s to v1.34.5 (#6366)
 * [0720adb48](https://github.com/kubeovn/kube-ovn/commit/0720adb48c483b831b088b24806716982e43f361) fix(kind): set containerd LimitNOFILE to 4096 to avoid valgrind failure (#6337)
 * [95b2b80fc](https://github.com/kubeovn/kube-ovn/commit/95b2b80fcd879cb2ac61a82ef02448e91d889e2b) fix: prevent subnet from getting permanently stuck when VLAN is not ready (#6352)
 * [ca9282c0d](https://github.com/kubeovn/kube-ovn/commit/ca9282c0d4f02027f8b5d23ee6c536fc20fd3faa) cni-server: fix vlan id reference (#6365)
 * [b78c94c56](https://github.com/kubeovn/kube-ovn/commit/b78c94c565057024f7e839def4fa74702bd47030) fix(deps): update module golang.org/x/net to v0.51.0 (#6341)
 * [1b4dd3144](https://github.com/kubeovn/kube-ovn/commit/1b4dd3144552fb6edbca03966d85e5a9d0103eb6) feat: introduce serviceMonitor and labels to services (#6340)
 * [b67d8a480](https://github.com/kubeovn/kube-ovn/commit/b67d8a48066c13d2793afeb0f49e9b7adbbac100) fix(ovn): update patches for compatibility with upstream OVN branch-25.03 (#6331)
 * [7fee4a495](https://github.com/kubeovn/kube-ovn/commit/7fee4a495405cc3058589305523e30395ed34411) fix(loadbalancers): ip_port_mapping is overwritten by new vips (#6302)
 * [e3a1fd660](https://github.com/kubeovn/kube-ovn/commit/e3a1fd660f5b35e45d670ce98396ca0ed45170ca) fix(ovn): add null check before sbrec delete calls to prevent crashes (#6316)
 * [77bb39380](https://github.com/kubeovn/kube-ovn/commit/77bb39380bcd931aea1e1c792ff66839ef2ceac5) chore(deps): update dependency go to v1.26.0 (#6293)
 * [aeccd2d33](https://github.com/kubeovn/kube-ovn/commit/aeccd2d33c3ebc2ea58270e72b1c8441db52e883) fix lint issues
 * [be347caf1](https://github.com/kubeovn/kube-ovn/commit/be347caf17f110bc5369511fa025219f730762b6) fix(controller): prevent NPE in isOvnSubnet when subnet is nil (#6312)
 * [7ee1a5632](https://github.com/kubeovn/kube-ovn/commit/7ee1a5632b84f48ea2644db975de5b6f934adfd7) update .gitignore
 * [def05ac8e](https://github.com/kubeovn/kube-ovn/commit/def05ac8e0c56445b4caaa1840cb873aad637236) add CLAUDE.md and link to AGENTS.md
 * [103894d95](https://github.com/kubeovn/kube-ovn/commit/103894d95937c8739b1b97a593141139a09a6c73) prepare for next release

### Contributors

 * Johann Schley
 * Juan Morete
 * Mengxin Liu
 * SKALA NETWORKS
 * changluyi
 * renovate[bot]
 * 张祖建

## v1.15.4 (2026-02-12)

 * [fe5412f18](https://github.com/kubeovn/kube-ovn/commit/fe5412f180adba797d6cdf6818fe53928e10f1a7) release v1.15.4
 * [5b263ed93](https://github.com/kubeovn/kube-ovn/commit/5b263ed93d096b762450a09c26b1cde74efd139a) remove unnecessary logs
 * [8291be7a9](https://github.com/kubeovn/kube-ovn/commit/8291be7a94b0f43e687a79b7b3c4591c07e340f3) refactor(controller): clean up redundant logic in vpc controller (#6295)
 * [b60f10006](https://github.com/kubeovn/kube-ovn/commit/b60f10006130cc29357b7b70d09f8a4509835155) fix: skip IPAM init for empty IP annotation and handle nil IPRangeList (#6296)
 * [673317212](https://github.com/kubeovn/kube-ovn/commit/6733172126d9db1874671eb16571a0838054a97f) fix duplicate description
 * [c72287473](https://github.com/kubeovn/kube-ovn/commit/c72287473841c30cb0dd9ecd84f0d51e9e71de8b) refactor(ovn_leader_checker): improve readability and naming (#6263)
 * [26dab970f](https://github.com/kubeovn/kube-ovn/commit/26dab970f58f447aa1814d7ea887ef8a9cd267e2) fix(iptables-eip): look up subnet directly by name instead of resolving through NAD (#6287)
 * [25d473e73](https://github.com/kubeovn/kube-ovn/commit/25d473e73931892a8c71ab63537abc75e050c68e) fix(deps): update golang (#6289)
 * [9a1db503b](https://github.com/kubeovn/kube-ovn/commit/9a1db503b4331e803e2571a52d62a086ac6dee5f) prepare for next release

### Contributors

 * Mengxin Liu
 * Zespre Chang
 * renovate[bot]
 * zbb88888

## v1.15.3 (2026-02-09)

 * [c5c83b57d](https://github.com/kubeovn/kube-ovn/commit/c5c83b57da2cf57b9f52d53e453e628843b65732) release v1.15.3
 * [58aa0fd2e](https://github.com/kubeovn/kube-ovn/commit/58aa0fd2e149904f4f5116566dc75e031f2cbeda) update cert manager
 * [9a0fa3cf8](https://github.com/kubeovn/kube-ovn/commit/9a0fa3cf8ccde9de51a398d18520598a04f680e3) perf: only enqueue subnet when not ready (#6268)
 * [6ce9c3039](https://github.com/kubeovn/kube-ovn/commit/6ce9c30390e0354a0c31904f949ea5768b21774d) fix(deps): update module golang.org/x/sys to v0.41.0 (#6282)
 * [f0075adf0](https://github.com/kubeovn/kube-ovn/commit/f0075adf0dbc0bfba780e8cfbbdaa9e8da28ee58) chore(deps): update module golang.org/x/oauth2 to v0.35.0 (#6279)
 * [a3f04c743](https://github.com/kubeovn/kube-ovn/commit/a3f04c7435c5721715e39c879dc0b575ea4422dd) controller: fix stale ip cr after subnet deletion (#6275)
 * [b0d9b4316](https://github.com/kubeovn/kube-ovn/commit/b0d9b431673729893589231a1ff04694af57c4c1) cni-server: add static fdb entry for subnets with u2o enabled (#6269)
 * [3ab93d666](https://github.com/kubeovn/kube-ovn/commit/3ab93d666c82565a1dc08ba905d986a55ac033dd) controller: fix ips owned by subnet u2o ip and mcast querier ip deleted in gc (#6132)
 * [ae142dbf9](https://github.com/kubeovn/kube-ovn/commit/ae142dbf9b7960d271860749948c793c6b5c0127) security: specify a storage limit for containers (#6259)
 * [bf66a67bd](https://github.com/kubeovn/kube-ovn/commit/bf66a67bda190228c8e0c2793020a5977b177fb9) chore(deps): update dependency go to v1.25.7 (#6266)
 * [f68017379](https://github.com/kubeovn/kube-ovn/commit/f6801737942725a9bdcc996f538b91dca7aa0102) fix: Ensure Kubevirt optimizations do not use old migration state (#6239)
 * [ccf825412](https://github.com/kubeovn/kube-ovn/commit/ccf8254124dab659438e1f7c5902ebee8fc9da59) fix localnet port should not be static "patch-localnet." (#6244)
 * [5f2252890](https://github.com/kubeovn/kube-ovn/commit/5f2252890c2b6e690ce3dbb1a1726467c3c8d63d) rbac: replace wildcard with a clear list of allowed verbs (#6233)
 * [54001627e](https://github.com/kubeovn/kube-ovn/commit/54001627ef1c1126deb8abf2ecf9039d8f270cce) remove duplicate description (#6234)
 * [1fade93b2](https://github.com/kubeovn/kube-ovn/commit/1fade93b2e68d5f652e05bc1e6851f581e87e713) Add vendor ID to transit switches (#6186)
 * [8882096b9](https://github.com/kubeovn/kube-ovn/commit/8882096b939c9dd2634cfcc66045f7f008277851) prepare for next release

### Contributors

 * Eli Treuherz
 * Mengxin Liu
 * changluyi
 * jskrill
 * renovate[bot]
 * 张祖建

## v1.15.2 (2026-01-29)

 * [a4ecdaac8](https://github.com/kubeovn/kube-ovn/commit/a4ecdaac82c45396bb34d3c67569bb44c06cdc22) release v1.15.2
 * [0622723d3](https://github.com/kubeovn/kube-ovn/commit/0622723d30efca44dea6e8dfd7e7182f5ac7a017) fix: skip NetworkUnavailable condition in non-primary CNI mode (#6195)
 * [25acdabb4](https://github.com/kubeovn/kube-ovn/commit/25acdabb45397da11e4d59a3863d8713af067eef) fix(vpcnatgw): nat gateway uses faulty default network (#6212)
 * [b64138b90](https://github.com/kubeovn/kube-ovn/commit/b64138b9040bd901e5d62075f7d0d9b13f0dd810) controller: fix vpc egress gateway forward policy not working for targets within the internal subnet's cidr (#6222)
 * [c608e4841](https://github.com/kubeovn/kube-ovn/commit/c608e4841bec4d8eb6a03bf9e667fc1ddd755e2e) controller: fix vpc egress gateway not working for targets within the internal subnet's cidr (#6218)
 * [90eea0ed3](https://github.com/kubeovn/kube-ovn/commit/90eea0ed3b8511f88b224085597febda89d96ba3) Update start-ic-db.sh to have gen_conn_addr function defined (#6216)
 * [c0e9ec780](https://github.com/kubeovn/kube-ovn/commit/c0e9ec780a7dff9af773f01f9dd72553c9fd1b6c) fix: caching NAD CRD should before all kubeovn crds and pod (#6198)
 * [e3296ba1f](https://github.com/kubeovn/kube-ovn/commit/e3296ba1fd73f36f678674dcee2e981804bd52ff) prepare for next release

### Contributors

 * Mengxin Liu
 * SKALA NETWORKS
 * andrewlee1089
 * dnugmanov
 * zbb88888
 * 张祖建

## v1.15.1 (2026-01-25)

 * [02d360b90](https://github.com/kubeovn/kube-ovn/commit/02d360b9096cd6d1acdf24106b0bec7d3e1683b6) release v1.15.1
 * [b29a0ab9a](https://github.com/kubeovn/kube-ovn/commit/b29a0ab9ac457ae89fbfc937b54a3b94e680ad98) fix(vpc-nat-gw): use consistent provider for all pod annotations (#6204)
 * [b3d96f7cd](https://github.com/kubeovn/kube-ovn/commit/b3d96f7cdfd0d15e3d2f293d3735f9b180965931) fix(vpcnatgw): cannot DNAT same EIP and same external port on two different protocols (#6201)
 * [e22295bb1](https://github.com/kubeovn/kube-ovn/commit/e22295bb16b6bda49ef0fe88e323094db2c2b788) fix: cleanup auto-created vlan subinterfaces on interface switch (#6193)
 * [a8be30f83](https://github.com/kubeovn/kube-ovn/commit/a8be30f83f07d390b998995cae24a9da9e923f17) fix 1.15 base image build
 * [f6f0b9848](https://github.com/kubeovn/kube-ovn/commit/f6f0b9848081a635ec72c8f45ed31d8aafcb3bfb) fix: enforce neighbor address validation in configuration (#6134)
 * [0f605b825](https://github.com/kubeovn/kube-ovn/commit/0f605b82500d9c371ce7708390d01490128aff3d) fix: acquire specified subnet from annotation (#6180)
 * [508d9e816](https://github.com/kubeovn/kube-ovn/commit/508d9e81645a2bd27fff17da65a35c2ed64e6664) chore(deps): update dependency go to v1.25.6 (#6170)
 * [21aaa869b](https://github.com/kubeovn/kube-ovn/commit/21aaa869bee07b8429e0f79cd0e85c976d42462f) add metallb underlay v6 dualcase (#6159)
 * [7ff5bd3c4](https://github.com/kubeovn/kube-ovn/commit/7ff5bd3c4cb1ab855bb417dd8f59529f65e3d0d0) fix(helm): rolling upgrade strategy for ovs is wrong (#6160)
 * [2a1e6f570](https://github.com/kubeovn/kube-ovn/commit/2a1e6f570490370dcf55666a5b154fb2a690bbaa) fix(deps): update golang (#6155)
 * [05424e620](https://github.com/kubeovn/kube-ovn/commit/05424e6208fedbb6ac4efd5ed27b4092f9ed00d3) fix(deps): update golang (#6152)
 * [2c40054c5](https://github.com/kubeovn/kube-ovn/commit/2c40054c5baf3f688fa61037c9d3983141794897) fix(deps): update golang (#6149)
 * [d906cccb3](https://github.com/kubeovn/kube-ovn/commit/d906cccb3e37cca8ff87dc1ac5b44f41fdecd436) some left crd add description (#6147)
 * [4135e1f12](https://github.com/kubeovn/kube-ovn/commit/4135e1f121f77d67fa3353ec2fff4a54a15712a4) fix(deps): update module golang.org/x/sys to v0.40.0 (#6145)
 * [1cef6f91b](https://github.com/kubeovn/kube-ovn/commit/1cef6f91b502a5b9dc70bc24755ed4ecdd6f0587) add kube-ovn cr description (#6133)
 * [03efbe58f](https://github.com/kubeovn/kube-ovn/commit/03efbe58f8847a6a04f60dc78db2d1a870a457e2) feat(chart): introduce extraEnv variable for all components (#6142)
 * [65df1a7f2](https://github.com/kubeovn/kube-ovn/commit/65df1a7f2a85b233304e9ce41ab6f363b7b414b0) feat(chart): make DaemonSet update strategy configurable in values.yaml (#6136)
 * [d575fe803](https://github.com/kubeovn/kube-ovn/commit/d575fe803a10cda37662bdced918a3e7d9d68b15) bugfix kube-ovn image tag from v.1.14.0 to v1.15.0, add default labels (#6135)
 * [b00cbe47f](https://github.com/kubeovn/kube-ovn/commit/b00cbe47f60ac439bcc60860dd5e9aec1a2a9e7d) chore(deps): update dependency talos to v1.12.1 (#6126)
 * [59aa5cff3](https://github.com/kubeovn/kube-ovn/commit/59aa5cff32ab5d1e0e475898e1334eb674808475) ci: fix loading vpc-nat-gateway image into kind cluster (#6128)
 * [4a4d94e69](https://github.com/kubeovn/kube-ovn/commit/4a4d94e69659bd5eac5b036d61f6c999188c3c9f) ci: fix non-primary cni e2e test (#6127)
 * [26d9c0259](https://github.com/kubeovn/kube-ovn/commit/26d9c02590798248e35b0349bede7111e3208cf3) fix: add validation for iptables rules to prevent malformed commands in nat gateway pod (#6125)
 * [a8b25e78f](https://github.com/kubeovn/kube-ovn/commit/a8b25e78f0fd80750ec143f7f5e9f13022df4738) cleanup e2e tests
 * [78b6de652](https://github.com/kubeovn/kube-ovn/commit/78b6de652fcd4dec3ad330adbbacfe236661175a) Add open flow sync refer to ovn-k8s (#6117)
 * [6035ee206](https://github.com/kubeovn/kube-ovn/commit/6035ee2067350d2d064b0c24bd626942e4f4bee7) prepare for next release

### Contributors

 * Diego Gsponer
 * Johann Schley
 * Mengxin Liu
 * SKALA NETWORKS
 * changluyi
 * clyi
 * renovate[bot]
 * zbb88888
 * zhangzujian
 * 张祖建

## v1.15.0 (2026-01-01)

 * [304cd067f](https://github.com/kubeovn/kube-ovn/commit/304cd067f3abb17c46ce63aea579e4130cbb9a29) docs: updated CHANGELOG.md (#6113)
 * [4cafeecbe](https://github.com/kubeovn/kube-ovn/commit/4cafeecbec6e9605006998022831c66d2ea531be) refactor: add/update field selector to informer factory (#6091)
 * [2d5f4e502](https://github.com/kubeovn/kube-ovn/commit/2d5f4e5023d2051fe8ec878d00d2e09604278a0d) ovn db: join cluster with raft header backup (#6106)
 * [3baf28f13](https://github.com/kubeovn/kube-ovn/commit/3baf28f133c85864c81495266eabebe6cb63166c) better nat gw e2e (#6096)
 * [5e14dc115](https://github.com/kubeovn/kube-ovn/commit/5e14dc115e820087f3a645152813a3e320e354f1) remove unused (#6110)
 * [4e35b4687](https://github.com/kubeovn/kube-ovn/commit/4e35b4687816c3343a7989136304a1094d95ef9e) rename two gw nic routes variables (#6094)
 * [e5421f82f](https://github.com/kubeovn/kube-ovn/commit/e5421f82fe7737bc6f695c282e249c1213d43398) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.14.0 (#6111)
 * [a1d760668](https://github.com/kubeovn/kube-ovn/commit/a1d76066840e0d0c4caccd14e43fafdb3e55eb27) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.13.0 (#6109)
 * [788eae069](https://github.com/kubeovn/kube-ovn/commit/788eae06924edcb0888369c4daa22255ebc22746) docs: updated CHANGELOG.md (#6105)
 * [cde954cde](https://github.com/kubeovn/kube-ovn/commit/cde954cdebd0a5ac9de7c1f80d7dd9cd9479aee8) cni-server: do not limit read-only ovs-vsctl commands (#6104)
 * [38025b277](https://github.com/kubeovn/kube-ovn/commit/38025b2778fadf57d8a19aa4e8037c8365b0e39c) controller: fix LSP gc when no ls key in external-ids (#6103)
 * [ed66bf08d](https://github.com/kubeovn/kube-ovn/commit/ed66bf08d843a77b3271008da395f8f837601aae) fix ut (#6095)
 * [21bd8df16](https://github.com/kubeovn/kube-ovn/commit/21bd8df16e7280c744950498ab62ef4671cf8862) northd: set dl_src for packets redirected by router port (#6099)
 * [1410f79bb](https://github.com/kubeovn/kube-ovn/commit/1410f79bb7babd70d312b30c6cf89e942ca6b5dd) docs: updated CHANGELOG.md (#6101)
 * [f158b7bf7](https://github.com/kubeovn/kube-ovn/commit/f158b7bf795e79572621dfa870258b1ba13ed229) cni-server: fix ovn0 gateway check (#6098)
 * [c7296d713](https://github.com/kubeovn/kube-ovn/commit/c7296d71382309aad571eb5c1c3f8328a55df8cc) update next.md
 * [8948b6de2](https://github.com/kubeovn/kube-ovn/commit/8948b6de2f71c49bcf0ccb7dfb6fb37352d3cb09) fix(vm/sts): handle multiple network interfaces when cleaning resources on VM Pod deletion (#6089)
 * [68134fbc6](https://github.com/kubeovn/kube-ovn/commit/68134fbc689390cbc300c07c01bc04df29ea7e2c) fix(deps): update module google.golang.org/grpc to v1.78.0 (#6097)
 * [b25af2cc9](https://github.com/kubeovn/kube-ovn/commit/b25af2cc96370a755634fdf77e83ce75ee4afee9) fix(deps): update module github.com/containerd/nerdctl/v2 to v2.2.1 (#6090)
 * [faa14a995](https://github.com/kubeovn/kube-ovn/commit/faa14a995a2c267cb91680e6a10c86c8bf6953dc) ci: fix talos cluster initialization (#6093)
 * [5151eb9e4](https://github.com/kubeovn/kube-ovn/commit/5151eb9e4d75e151b097b3d0d5b5ea5ba2cb15a8) refactor: unify environment variable names using util package constants (#6061)
 * [769d1b7ba](https://github.com/kubeovn/kube-ovn/commit/769d1b7bad6ba29042118bbb0838d090a95121ff) add support for provider network vlan interfaces (#5949)
 * [2c0737090](https://github.com/kubeovn/kube-ovn/commit/2c0737090532832acb273cdccc0a91b7d656d43c) cni-server: partially revert #5973 to enable hostPID (#6087)
 * [d2409e4ac](https://github.com/kubeovn/kube-ovn/commit/d2409e4acd8462f75104f4afac1426b54fc9a237) crd: support new tolerations operators Gt and Lt introduced in k8s v1.35 (#6079)
 * [55b1cb394](https://github.com/kubeovn/kube-ovn/commit/55b1cb394c4caabd7dbb78a4f32fd72392a299f6) refactor: unify type constants (#6081)
 * [475cf5379](https://github.com/kubeovn/kube-ovn/commit/475cf5379ddab2677796ebe08387d160bcbe3539) remove windows support (#6080)
 * [fa35cbfa2](https://github.com/kubeovn/kube-ovn/commit/fa35cbfa2cf69c3193d6ea29991349c97583fad6) docs: updated CHANGELOG.md (#6082)
 * [2a2fc8b14](https://github.com/kubeovn/kube-ovn/commit/2a2fc8b149bd71ee448f4bf0001b9777b36edde2) fix: correct typo in Makefile and enhance IP pool status event logging (#6071)
 * [d537ded6b](https://github.com/kubeovn/kube-ovn/commit/d537ded6b0d15b576a9fdeb6c8e0cc531d997366) libovsdb: set inactivity timeout only when passed in value is greater than zero (#6078)
 * [990d22fde](https://github.com/kubeovn/kube-ovn/commit/990d22fde83e511234d4a429fedbf7172b7bb0cc) docs: updated CHANGELOG.md (#6077)
 * [aaff3bd6d](https://github.com/kubeovn/kube-ovn/commit/aaff3bd6d1a6ea3729e8356b08139697df5a6bce) skip vip e2e test for versions prior to v1.15
 * [fd1e9467f](https://github.com/kubeovn/kube-ovn/commit/fd1e9467f3a15922cffb12e988cf16b878481d57) skip e2e (#6076)
 * [3a4c77e34](https://github.com/kubeovn/kube-ovn/commit/3a4c77e34337b974723afdd1d26687c0007ff921) fix occasional migration failures caused by timing issues. (#6066)
 * [b73a8f64e](https://github.com/kubeovn/kube-ovn/commit/b73a8f64e20ad5a6fabc3db1f62bd1b6d9da2db7) chore(deps): update dependency cilium to v1.18.5 (#6072)
 * [f4d960ce4](https://github.com/kubeovn/kube-ovn/commit/f4d960ce4d5b007cab7d6b7d7c2318b90df90347) chore(deps): update dependency kind to v0.31.0 (#6073)
 * [d048c6adf](https://github.com/kubeovn/kube-ovn/commit/d048c6adf65fc80eefacfae1798283e09da1f17c) chore(deps): update kindest/node docker tag to v1.35.0 (#6074)
 * [2f57204f8](https://github.com/kubeovn/kube-ovn/commit/2f57204f85ac1c2553ae0c65a6429249b13812eb) fix(deps): update module github.com/containerd/containerd/v2 to v2.2.1 (#6075)
 * [2346e3582](https://github.com/kubeovn/kube-ovn/commit/2346e3582773497d8e4adf873809491f9e89145d) Add tests for VIP finalizer handling and subnet status updates (#6068)
 * [14de68217](https://github.com/kubeovn/kube-ovn/commit/14de682177139957671962a17984262a979c14bb) ovsdb: exit clustered ovsdb server if multiple raft leaders found (#6065)
 * [055891da4](https://github.com/kubeovn/kube-ovn/commit/055891da44780f4a2b352c1009a56fdf858cd02f) feat: add vendor tagging for OVN resources and migration support (#6001)
 * [3f881030c](https://github.com/kubeovn/kube-ovn/commit/3f881030cb6914648e85be7f20564a1bbdd89ed1) chore(deps): update kindest/node docker tag to v1.34.3 (#6067)
 * [d1044aaaf](https://github.com/kubeovn/kube-ovn/commit/d1044aaaf7a996b4adc9f05dde3a472619274008) fix(vpcnatgateway): bgp speaker would crash if no ipv4 is available (#6059)
 * [16106fc97](https://github.com/kubeovn/kube-ovn/commit/16106fc9748686b9ff2bbfd52f9a0675dc6ee75a) feat(natgw): give pod_ips to bgp speakre to prioritize IPv4 over IPv6 (#6058)
 * [d506cb3fa](https://github.com/kubeovn/kube-ovn/commit/d506cb3fab0193951c7584fdbfd28db17909d031) fix(helm): template and podaffinity (#6057)
 * [fc9d03b37](https://github.com/kubeovn/kube-ovn/commit/fc9d03b37be412941bb99b97025c33ad8266a9ce) 跟踪日志发现在eip绑定qos以及修改qos的过程中，有v4ip为空的情况也会执行nat网关的tc规则，v4ip为空规则解析错误，排除该场景。 (#6055)
 * [f2fe74720](https://github.com/kubeovn/kube-ovn/commit/f2fe74720b40ced6ad23f1483d336fb28b0c19fb) docs: updated CHANGELOG.md (#6052)
 * [8e7108730](https://github.com/kubeovn/kube-ovn/commit/8e71087306793380c160b8288bf4bf7a38386fc5) e2e: generate network policy api manifest url from go.mod (#6050)
 * [1c37f0021](https://github.com/kubeovn/kube-ovn/commit/1c37f00214bab47df8df75823f39f7906fc00fae) chore(deps): update dependency helm to v4.0.4 (#6051)
 * [bf6760ae1](https://github.com/kubeovn/kube-ovn/commit/bf6760ae14b5854afc5905088ba931addfe1d751) lint: skip generated files (#6049)
 * [73865d26d](https://github.com/kubeovn/kube-ovn/commit/73865d26d61004134e05d186599ff454e033177b) fix(deps): update module libovsdb (#6046)
 * [2b8ed7415](https://github.com/kubeovn/kube-ovn/commit/2b8ed7415ea99cbdfb0947d1039c6f2590f37b6d) chore(deps): update github artifact actions (#6048)
 * [5335a5912](https://github.com/kubeovn/kube-ovn/commit/5335a5912121b7d21fdb9a8e2fcc946175bbb2d2) chore(deps): update kindest/node docker tag to v1.34.2 (#6047)
 * [2733f7ada](https://github.com/kubeovn/kube-ovn/commit/2733f7ada626cac90ff3b6faf596d423defa1504) cni: disable ipv6 RA (#6045)
 * [823edaf56](https://github.com/kubeovn/kube-ovn/commit/823edaf564ee3518a6e9a9d7f21c0d2f3f6a5ec3) chore(deps): update dependency gosec to v2.22.11 (#6040)
 * [32df70ebc](https://github.com/kubeovn/kube-ovn/commit/32df70ebca80ec765e7ffb938b193e86820917ff) chore(deps): update dependency helm to v4.0.2 (#6041)
 * [b59c21bef](https://github.com/kubeovn/kube-ovn/commit/b59c21beffd23f0a71b2854845b449d2ec1402ee) chore(deps): update actions/cache action to v5 (#6044)
 * [6e21af232](https://github.com/kubeovn/kube-ovn/commit/6e21af232c911d5ae412a68e6dd1d32588c9792b) ip双栈的情况下gc interface失败 + sts的pod在经过多次重启后偶现chassis和pod所在节点不匹配的gc (#6033)
 * [7a3232ce6](https://github.com/kubeovn/kube-ovn/commit/7a3232ce6f4f5072a5b571529a5bfbdbd50632d2) fix e2e vlan interface (#6036)
 * [76bbf44bf](https://github.com/kubeovn/kube-ovn/commit/76bbf44bf1470f91bfe1c1cda75719e666a4ddb4) fix(deps): update kubernetes (master) (#6037)
 * [a6fcfbcf5](https://github.com/kubeovn/kube-ovn/commit/a6fcfbcf56bc6b9f7b3fb1c0728b016065523b8b) feat(charts/kube-ovn-v2): enable custom node affinity on deployments (#6025)
 * [526202664](https://github.com/kubeovn/kube-ovn/commit/5262026646c460436c76bfc2c17c22ad4e9d6c95) fix: double quoting causing nxdomain (#6039)
 * [0a1e96165](https://github.com/kubeovn/kube-ovn/commit/0a1e961655deca02abca8a332695e72f9742bc92) docs: updated CHANGELOG.md (#6035)
 * [8120178b3](https://github.com/kubeovn/kube-ovn/commit/8120178b3018113e6f31a70393dc32d763717d80) chore(deps): update actions/checkout action to v6 (#6038)
 * [7e936e755](https://github.com/kubeovn/kube-ovn/commit/7e936e755aa39a9746804b25b550939589fb3a1f) simplify subnet.go code (#6009)
 * [9e3780633](https://github.com/kubeovn/kube-ovn/commit/9e3780633c46a87fcaf61e111fa3695d9307bce9) Add support for ClusterNetworkPolicies, refactor and add more testing (#5921)
 * [cced83722](https://github.com/kubeovn/kube-ovn/commit/cced8372242abf1731e0340e331f0794fa6a9152) Fix empty subnet (#5964)
 * [09af0ba6a](https://github.com/kubeovn/kube-ovn/commit/09af0ba6a179d724deb6c3a20bb0d3fcd7fa1fdc) remove unused function (#6032)
 * [e2519ceaa](https://github.com/kubeovn/kube-ovn/commit/e2519ceaaa36c77694d9b8dc9fa471863551eb5a) chore(deps): update peter-evans/create-pull-request action to v8
 * [23f96fe52](https://github.com/kubeovn/kube-ovn/commit/23f96fe5242416bf63e796031cfbbbbd89ace9c3) Only GC LRPs that have the "vendor"="kube-ovn" external ID.
 * [adeb80a2f](https://github.com/kubeovn/kube-ovn/commit/adeb80a2f4ab7b2fd39f369e45a674c7dc94edd3) fix(deps): update module golang.org/x/tools to v0.40.0 (#6030)
 * [97abc081d](https://github.com/kubeovn/kube-ovn/commit/97abc081daa252cf6f9ed2d4d536110562f89b4a) chore(deps): update dependency cni-plugin to v1.9.0 (#6024)
 * [39a9dfaf1](https://github.com/kubeovn/kube-ovn/commit/39a9dfaf18adfe98d6fe7fb5264cfd4e53b9be60) fix(deps): update module github.com/containernetworking/plugins to v1.9.0 [security] (#6027)
 * [708e07a89](https://github.com/kubeovn/kube-ovn/commit/708e07a895e7f8c269a2d6a73f160c60d49bbb0a) chore(deps): update dependency cert-manager to v1.19.2 (#6028)
 * [cc0242fd1](https://github.com/kubeovn/kube-ovn/commit/cc0242fd16b1cf1a8e4ac747384244858fa2d9a6) fix(deps): update module github.com/cert-manager/cert-manager to v1.19.2 (#6023)
 * [31fc9817e](https://github.com/kubeovn/kube-ovn/commit/31fc9817e218754f1dea00886a7a9e8f696fd724) fix the crd with supported schema
 * [856988ae7](https://github.com/kubeovn/kube-ovn/commit/856988ae766bb4fa7ff0931c14625cf5e9848a91) ovn process skip non-ovn subnet (#6018)
 * [5307a238c](https://github.com/kubeovn/kube-ovn/commit/5307a238c25a30fb6093315dab5c819e8c4a16ac) enable_ipforward for vpc_nat_gateway
 * [bd9b79112](https://github.com/kubeovn/kube-ovn/commit/bd9b791128fbd1123f7b9ebc2f0a0e8f8f58ed79) fix(deps): update module go.universe.tf/metallb to v0.15.3 (#5993)
 * [89323dd16](https://github.com/kubeovn/kube-ovn/commit/89323dd16cad34ba32455b7b88e96b2b2bd308d2) chore(deps): update dependency metallb to v0.15.3 (#5992)
 * [497eb5860](https://github.com/kubeovn/kube-ovn/commit/497eb5860ba34c24933c2ff8eef802b07366aa9d) fix(deps): update module github.com/onsi/gomega to v1.38.3 (#6013)
 * [5306f607a](https://github.com/kubeovn/kube-ovn/commit/5306f607a615ed64671a4df7f14146a15f07ac74) fix(deps): update golang (#6012)
 * [856adea60](https://github.com/kubeovn/kube-ovn/commit/856adea6047d160528459eeb28bcffc17e70c3f6) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.27.3 (#6011)
 * [563e6d73e](https://github.com/kubeovn/kube-ovn/commit/563e6d73ee29e4c658a75be3c9e0e9a4c89ab6e5) update agent.md
 * [d883fdd77](https://github.com/kubeovn/kube-ovn/commit/d883fdd777380b604011f48f3894b62a205fa17e) make sure ip in subnet using range (#6005)
 * [68b9dce79](https://github.com/kubeovn/kube-ovn/commit/68b9dce79c40260831cd3a3f08559055704945e7) fix: handle empty IP and CIDR in GetIPAddrWithMask function (#6002)
 * [21d0ff8e8](https://github.com/kubeovn/kube-ovn/commit/21d0ff8e887e59fc930d5435024c32c1ce0f2f17) chore(deps): update dependency golangci-lint to v2.7.2
 * [03b477466](https://github.com/kubeovn/kube-ovn/commit/03b477466910fbb0770c667048f610ce5bcf3556) fix package name
 * [9a926596b](https://github.com/kubeovn/kube-ovn/commit/9a926596b41522a71571de995f33a4535c4f056b) chore(deps): update dependency golangci-lint to v2.7.1
 * [2d7db6fe1](https://github.com/kubeovn/kube-ovn/commit/2d7db6fe17ab9612b1d9f9e1de60a8bb8fbc4b72) perf: only tregger update when related annotations changed
 * [5685ea844](https://github.com/kubeovn/kube-ovn/commit/5685ea844ed04be794b48caff2800ffda36b7b9c) docs: updated CHANGELOG.md
 * [1762844fd](https://github.com/kubeovn/kube-ovn/commit/1762844fdcf12d2d5d479d18521dfe45aac79ede) skip e2e test prior 1.15
 * [af6d331b3](https://github.com/kubeovn/kube-ovn/commit/af6d331b31113c09aa0a10368ee416040887bdce) deps: migrate github.com/osrg/gobgp/v3 to github.com/osrg/gobgp/v4 (#5990)
 * [8fe809617](https://github.com/kubeovn/kube-ovn/commit/8fe8096173e017116b40a63991ab903fe7b82337) subnet: forbid subnet vpc change
 * [dfcea355f](https://github.com/kubeovn/kube-ovn/commit/dfcea355fc17e7326b21d22a7c01223e92f0e65f) add auto create vlan sub interface (#5966)
 * [a4c9b09ad](https://github.com/kubeovn/kube-ovn/commit/a4c9b09ad20df6d54b44c59ad90c91c867a64d5d) np: allow icmpv6 to ipv6 ready check
 * [c74c5c2fb](https://github.com/kubeovn/kube-ovn/commit/c74c5c2fbed936e5b9a37a4b121f241bd0d1bca1) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.12.1 (#5986)
 * [5cd3dbbed](https://github.com/kubeovn/kube-ovn/commit/5cd3dbbed0af7eda552af8e60f07690f6f6b1f1a) Dockerfile: fix git not found (#5985)
 * [830f3171a](https://github.com/kubeovn/kube-ovn/commit/830f3171a3f883717832b38cd11bd55e197374a2) Revert "update agent.md"
 * [a7e591fc0](https://github.com/kubeovn/kube-ovn/commit/a7e591fc0176aab7b17a449e60939d908a664164) update agent.md
 * [1f8cbfc62](https://github.com/kubeovn/kube-ovn/commit/1f8cbfc620a069ee916fb72af724230d54df7f5a) perf: remove sleep
 * [c70bf5804](https://github.com/kubeovn/kube-ovn/commit/c70bf5804e22a69563972645ad3b248ededd3025) perf: add concurrency for subnet handle
 * [355d46314](https://github.com/kubeovn/kube-ovn/commit/355d463141190813618c313cad2d16aa0afc6a80) perf: use separate worker to handle routes update.
 * [217fa7fbd](https://github.com/kubeovn/kube-ovn/commit/217fa7fbd7920dd5f23d56d6ea4b10a8cd25acb7) docs: updated CHANGELOG.md (#5983)
 * [2f57cfba8](https://github.com/kubeovn/kube-ovn/commit/2f57cfba84c8e455f0273cff7839685f87d3f532) chore(deps): update dependency go to v1.25.5 (#5980)
 * [77fdfd704](https://github.com/kubeovn/kube-ovn/commit/77fdfd7046472f01f3f0978f78493811e39e4bd9) vpc egress gateway: add support for specifying tolerations (#5978)
 * [a9db5bf32](https://github.com/kubeovn/kube-ovn/commit/a9db5bf32a041d4d577d73f0eb875caf29b6aa3c) chore(deps): update dependency gobgp to v4.1.0 (#5974)
 * [d870fa2ff](https://github.com/kubeovn/kube-ovn/commit/d870fa2ff7b0c5bc4af1d0a5d08b85b404203713) docs: updated CHANGELOG.md (#5979)
 * [ec26b6643](https://github.com/kubeovn/kube-ovn/commit/ec26b66435ff339fce958ea01708efe6d61f3c6d) security: only approve validated CSR (#5972)
 * [5610b551b](https://github.com/kubeovn/kube-ovn/commit/5610b551b42b3153b11b82ed72cab3af15347e94) security: disable hostPID for pinger and cni-server (#5973)
 * [af94696c7](https://github.com/kubeovn/kube-ovn/commit/af94696c7f1a251715602cc434a3bb0734945175) cni-server: reduce redundant logging (#5975)
 * [73db20497](https://github.com/kubeovn/kube-ovn/commit/73db20497339b34f65ee4027fddd1558b48d13c1) update next.md
 * [d5bfa7e42](https://github.com/kubeovn/kube-ovn/commit/d5bfa7e4218507766fba7e3ac37cb77236e3dd4c) security: set sticky bits on world-writable directories (#5971)
 * [67bc6e325](https://github.com/kubeovn/kube-ovn/commit/67bc6e3255dba0178c734dc4bdb2f0ddd44959bb) migrate deps github.com/docker/docker (#5969)
 * [d971af3ba](https://github.com/kubeovn/kube-ovn/commit/d971af3ba0a622f5df83301671ab43679d1487a5) security: use security context to enable the docker/default seccomp profile in pod definitions (#5970)
 * [3f0df188b](https://github.com/kubeovn/kube-ovn/commit/3f0df188b02bbb2ef780ba9e6b630855dd10cf46) feat: add gatewayNodeSelectors (#5956)
 * [75f642a93](https://github.com/kubeovn/kube-ovn/commit/75f642a93c3ef48969172d21e090d5443bc34d2e) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.12.0 (#5968)
 * [12b1ae76c](https://github.com/kubeovn/kube-ovn/commit/12b1ae76c475e37ff0eb2734ed534337136def03) add instructions for AI
 * [ca16857ea](https://github.com/kubeovn/kube-ovn/commit/ca16857ea799918661bb3324e2af036fcb7443a4) speaker: fix missing log directory (#5959)
 * [cb8c7e318](https://github.com/kubeovn/kube-ovn/commit/cb8c7e318e9697c4419ee890137e4e56570b67fd) fix(deps): update module kubevirt.io/api to v1.7.0 (#5961)
 * [3b3a9e8d6](https://github.com/kubeovn/kube-ovn/commit/3b3a9e8d61069ee91a01dbf55252e1277bc02c8b) fix: correct typos and comments in endpoint_slice.go (#5958)
 * [438279a67](https://github.com/kubeovn/kube-ovn/commit/438279a67e64786ee7637106359e9d5733e4cc01) e2e: skip node encap ip tests for version prior to v1.15 (#5962)
 * [ac58b2bf8](https://github.com/kubeovn/kube-ovn/commit/ac58b2bf81bc8edc087be00d8ea0223c44a52c3e) [feat] Add ovn acl log rate (#5938)
 * [934206133](https://github.com/kubeovn/kube-ovn/commit/93420613352c8153203c29306390759e89d031c1) fix: correct V6Eip field assignment in OvnDnatRule status (#5957)
 * [a57a76a64](https://github.com/kubeovn/kube-ovn/commit/a57a76a647c97adf8951a6ae424c0c33d25d407d) e2e: create ovn nb client from server-side ovsdb schema (#5947)
 * [0ae970e56](https://github.com/kubeovn/kube-ovn/commit/0ae970e569328f4a523735112403e3871ddbdc83) update gitignore
 * [80855d397](https://github.com/kubeovn/kube-ovn/commit/80855d39754d65434e232d2d12498d49c81aed46) feat: support selecting specific node IP for overlay subnet encapsula… (#5946)
 * [91e2fe412](https://github.com/kubeovn/kube-ovn/commit/91e2fe412be9355b7bc356082c46114328e50006) chore(deps): update dependency kubevirt to v1.7.0 (#5952)
 * [e274625e5](https://github.com/kubeovn/kube-ovn/commit/e274625e5c0c4e1865b28e5339aaf00c4d9abc1b) e2e: do not change join cidr and service cidr for ovn-ic tests (#5951)
 * [b2aa1cbaa](https://github.com/kubeovn/kube-ovn/commit/b2aa1cbaa7d85a1d874c934a597ebc6dcda62e41) fix(deps): update module kubevirt.io/client-go to v1.7.0 (#5954)
 * [588abd01b](https://github.com/kubeovn/kube-ovn/commit/588abd01bc1c0e7bed993f7c0cce7086b2e43b5d) cni-server: fix missing env POD_IPS
 * [03156e61e](https://github.com/kubeovn/kube-ovn/commit/03156e61e7c77d5dec35c077a12e21a494b16674) cni-server: configure sysctl parameters on demand (#5950)
 * [e3924e0e2](https://github.com/kubeovn/kube-ovn/commit/e3924e0e219d749a38a44a0fb9a2eb9e128ed4f0) fix: correct spelling of "deleting" in log messages (#5944)
 * [8d271ebc7](https://github.com/kubeovn/kube-ovn/commit/8d271ebc7cbbeb46fcdba79bd31aaee78dda81b3) e2e: skip IPPool AddressSet test for versions < v1.15 (#5943)
 * [3b7dd0068](https://github.com/kubeovn/kube-ovn/commit/3b7dd00686963d427433f0088293329045441c1b) fix: correct typos and function names in error messages (#5940)
 * [b7dc72ea7](https://github.com/kubeovn/kube-ovn/commit/b7dc72ea7799e63a5d35d335d2e5613f4f126ae7) chore(deps): update dependency helm to v4.0.1 (#5939)
 * [e0d44c836](https://github.com/kubeovn/kube-ovn/commit/e0d44c8363f4972d66cd37045904b8f78d3d8e70) fix: handle missing NAD during pod deletion (#5929)
 * [9fcab65e2](https://github.com/kubeovn/kube-ovn/commit/9fcab65e2cbc2f0c9383e99bbc90b67e49693bcb) fix: use the same cni dir as the mount (#5928)
 * [dc0e27051](https://github.com/kubeovn/kube-ovn/commit/dc0e270516f944153e3bdebcbe0de6f627ead387) IPPool sync to address set (#5920)
 * [288c00dd8](https://github.com/kubeovn/kube-ovn/commit/288c00dd84610e51297dd677af90b7dd108cf400) fix: correct error log message in getDBStatus function (#5934)
 * [d26c95b60](https://github.com/kubeovn/kube-ovn/commit/d26c95b60b1808423ec51ae02ca3878e46c22b51) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.11.0 (#5930)
 * [29adff305](https://github.com/kubeovn/kube-ovn/commit/29adff3058355274b7402eac7774c73a8c9b5e0b) metrics: split subnet CIDR block for dual-stack subnets (#5931)
 * [9e2982259](https://github.com/kubeovn/kube-ovn/commit/9e298225963d49037abc8c25d20f086e4b24a673) security: fix C-0190 (#5927)
 * [cecdad75c](https://github.com/kubeovn/kube-ovn/commit/cecdad75cfb6fc1f564d884859928dd69a96c77a) fix(deps): update module google.golang.org/grpc to v1.77.0 (#5923)
 * [cb7b10f39](https://github.com/kubeovn/kube-ovn/commit/cb7b10f3963e00a47a381a871737e8f841eb2c98) chore(deps): update actions/checkout action to v6 (#5926)
 * [d57a569d9](https://github.com/kubeovn/kube-ovn/commit/d57a569d95deb7d048fb46000074d60339f2c408) get pods just before using and add some log (#5922)
 * [7fddea018](https://github.com/kubeovn/kube-ovn/commit/7fddea018c9e57036309f0cfdceb7b0fe8e1e5b2) fix(deps): update module go.uber.org/zap to v1.27.1 (#5925)
 * [ab2c0a6b3](https://github.com/kubeovn/kube-ovn/commit/ab2c0a6b3615290310947ee255426ca4d379f942) chore(deps): update module golang.org/x/crypto to v0.45.0 (#5924)
 * [350855f2c](https://github.com/kubeovn/kube-ovn/commit/350855f2cf5650b4bdcf0c8b943f7700698354aa) feat(install): split kube-ovn-pinger, keep version in the head, reload multus,  enable control not del existing pod (#5909)
 * [7abd42181](https://github.com/kubeovn/kube-ovn/commit/7abd42181832ce35a4fdb352111b269fd8fd22f4) optimize log output (#5916)
 * [3faa8a614](https://github.com/kubeovn/kube-ovn/commit/3faa8a614ad4ddbd5f8f93f3cd9757aad23a7799) docs: updated CHANGELOG.md (#5917)
 * [41301275b](https://github.com/kubeovn/kube-ovn/commit/41301275bb69f4fdc5b92a7f83a6a2a4606cea8d) cni-server: set cni config file permission to 600 (#5906)
 * [d33693bdf](https://github.com/kubeovn/kube-ovn/commit/d33693bdfc7d17050ad22ff8e5e7238baca0a380) chore(deps): update dependency golangci-lint to v2.6.2 (#5911)
 * [6aa5ab908](https://github.com/kubeovn/kube-ovn/commit/6aa5ab908fdc67589eb64cdb5d6a3c7c640267ec) chore(deps): update dependency helm to v4 (#5908)
 * [7eb0771e0](https://github.com/kubeovn/kube-ovn/commit/7eb0771e0a62c1d23ce20eff17b8f120f73c7815) fix(deps): update golang (#5907)
 * [186605333](https://github.com/kubeovn/kube-ovn/commit/186605333fd140fcf570a2278d43a2ad63eebcf2) docs: updated CHANGELOG.md (#5905)
 * [a9bd54776](https://github.com/kubeovn/kube-ovn/commit/a9bd5477672f99f997c2896436a21674d229fdf8) cni-server: detect if dbus is available (#5896)
 * [1f57817b2](https://github.com/kubeovn/kube-ovn/commit/1f57817b207ce5e59977be48eeb2b4f37a57abe6) fix(deps): update kubernetes (#5900)
 * [1b89bd2c7](https://github.com/kubeovn/kube-ovn/commit/1b89bd2c77db74eb1127b104829736cb35d94a1a) fix(deps): update module kubevirt.io/client-go to v1.6.3 (#5901)
 * [146f528a2](https://github.com/kubeovn/kube-ovn/commit/146f528a2cb40e7892092056951b13ad66b2bf13) chore(deps): update dependency helm to v3.19.2 (#5899)
 * [24abf1924](https://github.com/kubeovn/kube-ovn/commit/24abf1924b73cd88144fca9a53525291ece90240) chore(deps): update dependency cilium to v1.18.4 (#5898)
 * [dad97d892](https://github.com/kubeovn/kube-ovn/commit/dad97d8923d65106b86cd8c30ea55cd9ebc5da05)  update chart-v2 release script
 * [163c02679](https://github.com/kubeovn/kube-ovn/commit/163c0267964f23b9899544cbe9915da9f132c580) fix(deps): update module kubevirt.io/api to v1.6.3 (#5893)
 * [cd7e28c10](https://github.com/kubeovn/kube-ovn/commit/cd7e28c107598f52aef4c8a6b6cc494a2d023cfd) chore(deps): update dependency kubevirt to v1.6.3 (#5891)
 * [cc9b557ce](https://github.com/kubeovn/kube-ovn/commit/cc9b557cebf8c0852be2c3a078badab1e6f80a8b) chore(deps): update dependency submariner to v0.21.2 (#5892)
 * [f72fc7955](https://github.com/kubeovn/kube-ovn/commit/f72fc79551b68f3fb23684e1281fb97c007b7e6a) chore(deps): update dependency helm to v3.19.1 (#5890)
 * [1cd2c382d](https://github.com/kubeovn/kube-ovn/commit/1cd2c382df8d6ea3fcd46f1bd55de45a0728e00d) server: add hostname and pod IPs into TLS certificate SAN (#5888)
 * [3d50f6fda](https://github.com/kubeovn/kube-ovn/commit/3d50f6fda9409aac6dbf5e8262e424fbf5efefa5) rollback permission
 * [4c00ddff1](https://github.com/kubeovn/kube-ovn/commit/4c00ddff1f75e2bbececed830d81d5b155766fe2) fix version trim
 * [cb7943032](https://github.com/kubeovn/kube-ovn/commit/cb79430321a4856028b4034c027b555fd50bd440) docs: updated CHANGELOG.md (#5887)
 * [c38e53bce](https://github.com/kubeovn/kube-ovn/commit/c38e53bce608b7a263ced541bb571dde73527447) vpc controller: only manage static routes with kube-ovn as vendor (#5838)
 * [3e168db62](https://github.com/kubeovn/kube-ovn/commit/3e168db62b842bad5fc0bf1ce4065efac83871b7) fix: update comments for NETWORK_TYPE and TUNNEL_TYPE in install.sh (#5878)
 * [ef748e964](https://github.com/kubeovn/kube-ovn/commit/ef748e964da6072846034fa3288ee5bb2b296ea6) fix(deps): update golang (master) (#5881)
 * [99a2ac4f0](https://github.com/kubeovn/kube-ovn/commit/99a2ac4f0b3be0b1e96b6b170017c02e91f046d6) fix: loadbalancerservice value (#5884)
 * [c0807033e](https://github.com/kubeovn/kube-ovn/commit/c0807033e3d4f985c4973afe3ad05ca3414cfeb9) Bump github.com/opencontainers/selinux from 1.12.0 to 1.13.0 (#5883)
 * [1e0a5cc65](https://github.com/kubeovn/kube-ovn/commit/1e0a5cc65df0c8f3945254d8d6721f1231fcd006) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.9.0 (#5882)
 * [6781af168](https://github.com/kubeovn/kube-ovn/commit/6781af168c816c2883a43bd063fbdc7387884f4f) fix migrate failed (#5873)
 * [c87482667](https://github.com/kubeovn/kube-ovn/commit/c8748266734af51da37e6c1f2f6a1881743627b6) move docs to the docs repo
 * [8b4da66a3](https://github.com/kubeovn/kube-ovn/commit/8b4da66a3f9beafa5171eed931aac0b2d18f7a97) fix(deps): update module github.com/containerd/containerd/v2 to v2.2.0 (#5875)
 * [85cfb2067](https://github.com/kubeovn/kube-ovn/commit/85cfb20679d9990b8ad55b1ba19635a32c88f853) fix(deps): update module github.com/containerd/containerd/v2 to v2.1.5 [security] (#5877)
 * [94f7c47fc](https://github.com/kubeovn/kube-ovn/commit/94f7c47fc19978cb59c0cbfba77be1ae3e0da06b) return empty nodeips during chart rendering dry-run (#5760)
 * [c959d53c7](https://github.com/kubeovn/kube-ovn/commit/c959d53c7af543a7a06e5b49cb8e7fc586b1ff11) fix release permission
 * [537457202](https://github.com/kubeovn/kube-ovn/commit/53745720204f6a435791c5c82bc46262f45e3613) feat: allow vpcTag to be specified to override tag (#5845)
 * [7c98ac5c7](https://github.com/kubeovn/kube-ovn/commit/7c98ac5c72d781a601eba36e2f00d07f10ef8c6a) chore(deps): update dependency go to v1.25.4 (#5869)
 * [986fd5394](https://github.com/kubeovn/kube-ovn/commit/986fd53943d6909a92af698a5928225a98f8b63c) fix(deps): update module github.com/docker/docker to v28.5.2+incompatible (#5867)
 * [e6ca23bf6](https://github.com/kubeovn/kube-ovn/commit/e6ca23bf606ccbdffa5a01f852d3e5678d18b461) chore(deps): update dependency submariner to v0.21.1 (#5860)
 * [a19f92efe](https://github.com/kubeovn/kube-ovn/commit/a19f92efe719385c07e7f18f3aedb3e9f79126b9) 修复vip无法正常固定ip (#5865)
 * [c64f61cf2](https://github.com/kubeovn/kube-ovn/commit/c64f61cf2452b940fb3d97e75704e796dedca2ec) docs: updated CHANGELOG.md (#5862)
 * [f0990293c](https://github.com/kubeovn/kube-ovn/commit/f0990293ce852d08d8eb906f2e35808cd9f95a04) chore(deps): update dependency golangci-lint to v2.6.1 (#5859)
 * [150082c31](https://github.com/kubeovn/kube-ovn/commit/150082c31a0b329ad18203734ef4b5190a9ab4a3) fix golanglint (#5858)
 * [626ef41ab](https://github.com/kubeovn/kube-ovn/commit/626ef41abea7bd0750cb8813bfffb89d567230d2) fix metallb e2e faile (#5856)
 * [741c99833](https://github.com/kubeovn/kube-ovn/commit/741c99833aef2102772e6ae62565508606a42a5e) chore(maintainers): add new maintainer (#5855)
 * [3a0e1ec55](https://github.com/kubeovn/kube-ovn/commit/3a0e1ec555e64da71390f1a168d584c2a2e2625f) chore(deps): update dependency gobgp to v4 (#5851)
 * [2d35661b2](https://github.com/kubeovn/kube-ovn/commit/2d35661b26e5a21bf9fda38f2452009adbc07475) chore(deps): update helm/kind-action action to v1.13.0 (#5850)
 * [4aa82ce86](https://github.com/kubeovn/kube-ovn/commit/4aa82ce861cc9d7880afc3ec0c1b1a39f25202bd) fix(deps): update module sigs.k8s.io/controller-runtime to v0.22.4 (#5849)
 * [b273e7eb8](https://github.com/kubeovn/kube-ovn/commit/b273e7eb8d6ccc084f9ffeb7869e9cdbc87a7359) chore(deps): update dependency multus to v4.2.3 (#5843)
 * [bcecc3a00](https://github.com/kubeovn/kube-ovn/commit/bcecc3a00bb07b93dfb6b36d416d88e5d6121f91) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.8.2 (#5847)
 * [ea0b4d7c2](https://github.com/kubeovn/kube-ovn/commit/ea0b4d7c2cd1f874d842046a36a19ed026c7753c) fix(deps): update module gopkg.in/k8snetworkplumbingwg/multus-cni.v4 to v4.2.3 (#5844)
 * [b1b9aeb92](https://github.com/kubeovn/kube-ovn/commit/b1b9aeb92e0964594b6f2c42ce04cfda4efc2a67) chore(deps): update dependency golangci-lint to v2.6.0 (master) (#5841)
 * [171505c50](https://github.com/kubeovn/kube-ovn/commit/171505c500a0868a0e6b5f546254316035f2345d) chore(deps): update github artifact actions (#5842)
 * [bb3032480](https://github.com/kubeovn/kube-ovn/commit/bb303248009d0c72194857f6510669aef1abcaa7) Tests: Added go unit tests and E2E tests for Kube-OVN Non-Primary CNI mode (#5732)
 * [638d31de7](https://github.com/kubeovn/kube-ovn/commit/638d31de727c7c5c6134a47033f29dfcbe8ce75a) Update the kubeovn_deny_all security group after the virtual machine … (#5831)
 * [cf2bd0036](https://github.com/kubeovn/kube-ovn/commit/cf2bd00363ebe7d1d5249b695ad26677c7924f96) Add OCI registry support for Helm charts (#5837)
 * [9681fc70d](https://github.com/kubeovn/kube-ovn/commit/9681fc70d401e8a03ced8eabe0c5bd43ecc82c66) clean migrate  state when migrate is done (#5833)
 * [669cb5674](https://github.com/kubeovn/kube-ovn/commit/669cb5674a2cfb614fc7ff75f10855ad7e327792) update next
 * [61fd35968](https://github.com/kubeovn/kube-ovn/commit/61fd3596812f689f89efb5e701d10477e739be71) modify route priority for /32 src route (#5832)
 * [646dbed49](https://github.com/kubeovn/kube-ovn/commit/646dbed494059f9c66689fc019e5361fd11c89cf) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.27.2 (#5830)
 * [3045d07e6](https://github.com/kubeovn/kube-ovn/commit/3045d07e6b36af6eaacfa868693db6ba3dedf0fb) chore(deps): update github artifact actions (#5828)
 * [3c4e2ae94](https://github.com/kubeovn/kube-ovn/commit/3c4e2ae947833f3a7d0d5224bc420ca5d42a9114) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.8.1 (#5827)
 * [80f2be4b9](https://github.com/kubeovn/kube-ovn/commit/80f2be4b98430c06dce164b25b68f11d6c7be48d) chore(deps): update dependency cilium to v1.18.3 (#5826)
 * [6ee1e17bd](https://github.com/kubeovn/kube-ovn/commit/6ee1e17bda473ff0fca5c2a6a8c3a8cf54eac59f) add skip target ipcidrs conntrack (#5821)
 * [7656d9514](https://github.com/kubeovn/kube-ovn/commit/7656d9514487688e27e0436f4d0f9ab5fa2b2c40) fix(deps): update module kernel.org/pub/linux/libs/security/libcap/cap to v1.2.77 (#5829)
 * [a074286b4](https://github.com/kubeovn/kube-ovn/commit/a074286b4d9b1e6b4eedb916613c0cef2ad6e03a) feat(chart): set mtu for the pod overlay network (#5824)
 * [1789a1b23](https://github.com/kubeovn/kube-ovn/commit/1789a1b23361ee22e889a7a5bf34ad006b1cd60c) add lock for gateway exec operations (#5817)
 * [77b7ec592](https://github.com/kubeovn/kube-ovn/commit/77b7ec592262a85e8733a76c601bf36f80bcd3e0) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.27.1 (#5818)
 * [373866907](https://github.com/kubeovn/kube-ovn/commit/3738669071b3567955e394fc9671141c63a9db3a) fix: prevent GC from deleting KubeVirt VM interfaces (#5814)
 * [1633e46c1](https://github.com/kubeovn/kube-ovn/commit/1633e46c193b49af57f2615163eb205d8fbe55db) remove unnecessary gateway check for lax networkpolicy (#5801)
 * [3d3256dab](https://github.com/kubeovn/kube-ovn/commit/3d3256daba30b6a00ab7c87b094fc890179d85c5) fix talos e2e failure on ipv4 (#5812)
 * [9216651a8](https://github.com/kubeovn/kube-ovn/commit/9216651a8717435201e049299e7509dbccbbba4c) fix(netpol): lax dhcp must be after lb for egress (#5810)
 * [2bf532745](https://github.com/kubeovn/kube-ovn/commit/2bf5327453ab0404a597987c2a677702e417d8d6) fix: nat rule finalizer (#5806)
 * [dd79c777f](https://github.com/kubeovn/kube-ovn/commit/dd79c777fe818ad04338bbcea9a3ae979315c0d8) fix(netpol): don't check if rules already exist in lax mode (#5802)
 * [d7e2057e6](https://github.com/kubeovn/kube-ovn/commit/d7e2057e65db3ec6241b832c614f88969d7b9557) fix: update Bingbing Zhang's company affiliation to Individual (#5808)
 * [361c09bac](https://github.com/kubeovn/kube-ovn/commit/361c09bac63563e6166c46ecb528e3db4d9e331d) rename in maintainers (#5807)
 * [49aa35dc1](https://github.com/kubeovn/kube-ovn/commit/49aa35dc1fcb904a9c3dabd6fef215d4d800a166) update next
 * [023eb9de5](https://github.com/kubeovn/kube-ovn/commit/023eb9de538af566c5e41680f02ff52cfb61853b) add e2e for networkpolicy lax mode (#5800)
 * [ddb55379e](https://github.com/kubeovn/kube-ovn/commit/ddb55379effdd8b50b863e056b3a3f23182e0c84) feat(netpol): authorize l3 protocols (#5745)
 * [bfd08097e](https://github.com/kubeovn/kube-ovn/commit/bfd08097e9c2bf9b01479fdaffdf18f8c7040e1f) chore(deps): update dependency gosec to v2.22.10 (#5795)
 * [9e1b327a1](https://github.com/kubeovn/kube-ovn/commit/9e1b327a12827fd0f79345482754620a553b17e7) fix(deps): update module github.com/cert-manager/cert-manager to v1.19.1 (#5796)
 * [e4881fd4c](https://github.com/kubeovn/kube-ovn/commit/e4881fd4c9edb6817fe47003b3da39ae74f30cbc) chore(deps): update dependency cert-manager to v1.19.1 (#5797)
 * [eea882c74](https://github.com/kubeovn/kube-ovn/commit/eea882c74fa4b4fa28d03c585275b04d2527f7e6) remove internal-port type interface (#5794)
 * [7907b518c](https://github.com/kubeovn/kube-ovn/commit/7907b518cc2a88b1e05de33fd0b57760a14350d0) fix(netpol): don't add default block twice for dualstacks (#5741)
 * [e57a2a03b](https://github.com/kubeovn/kube-ovn/commit/e57a2a03baa3f5b4294d3b6a2b0fd1aabaec04c0) chore(deps): update dependency go to v1.25.3 (master) (#5791)
 * [015c7f9f8](https://github.com/kubeovn/kube-ovn/commit/015c7f9f8da4edea11eefe351cce8f7a59f6cbb3) 1. add anp domain case (#5750)
 * [3a6640a6f](https://github.com/kubeovn/kube-ovn/commit/3a6640a6f219c36bc97bb7c5487f72ca48e518b4) skip endpointslice tests before 1.14
 * [509abdae4](https://github.com/kubeovn/kube-ovn/commit/509abdae486bf0ee9eace1be10fcbb24595a1037) fix anp e2e (#5790)
 * [f2c155c01](https://github.com/kubeovn/kube-ovn/commit/f2c155c01f9b304d8ecdaa61174441d707600859) gc interface if it's not exist in podLister (#5789)
 * [0dea495d3](https://github.com/kubeovn/kube-ovn/commit/0dea495d3846e025393b4e27e323d324a656aa8b) remove unnecessary dependencies (#5788)
 * [94b8838fc](https://github.com/kubeovn/kube-ovn/commit/94b8838fcf1f292d8c72ae9b4b9fbcb903dd3dae) fix(deps): update module sigs.k8s.io/controller-runtime to v0.22.3 (#5783)
 * [03a730e30](https://github.com/kubeovn/kube-ovn/commit/03a730e305c7223850954f3480396cd1325e9d1e) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.8.0 (#5787)
 * [acacaf320](https://github.com/kubeovn/kube-ovn/commit/acacaf320ff4b31732eb43f0db0fef7733f041d0) docs: updated CHANGELOG.md (#5786)
 * [88a74a072](https://github.com/kubeovn/kube-ovn/commit/88a74a0721fef2f3f8682d15d0c53e4b598d6aac) replace httprunner with net/http
 * [8b9320fb0](https://github.com/kubeovn/kube-ovn/commit/8b9320fb0507578bac03f8ebb880ddec55598232) Bump github.com/quic-go/quic-go from 0.52.0 to 0.54.1
 * [6e5440020](https://github.com/kubeovn/kube-ovn/commit/6e5440020dd6af0bb019e1568cb5c74fae0b7db7) fix sometimes e2e metallb err (#5761)
 * [7f4c722c7](https://github.com/kubeovn/kube-ovn/commit/7f4c722c74424989d7204f4f47970ea08efec085) fix(deps): update module google.golang.org/grpc to v1.76.0 (#5778)
 * [bab2a9f1b](https://github.com/kubeovn/kube-ovn/commit/bab2a9f1baf732373a80c4ce61b139890b9a1ada) chore(deps): update dependency cert-manager to v1.19.0 (#5770)
 * [43e9819a6](https://github.com/kubeovn/kube-ovn/commit/43e9819a65da446cca660bb477a1e5266ada21d8) fix(deps): update module github.com/docker/docker to v28.5.1+incompatible (#5776)
 * [9594cbaee](https://github.com/kubeovn/kube-ovn/commit/9594cbaee15c8fe7ab21b1b506f5be81e59ee70a) fix(deps): update module github.com/cert-manager/cert-manager to v1.19.0 (#5777)
 * [4df5b6ecc](https://github.com/kubeovn/kube-ovn/commit/4df5b6ecc153778dded42c224af2eba05f2494a9) fix(deps): update golang (#5769)
 * [8ecb00d5e](https://github.com/kubeovn/kube-ovn/commit/8ecb00d5e091dcadcd262aca68d775b1ba37600a) chore(deps): update github/codeql-action action to v4 (#5779)
 * [783361268](https://github.com/kubeovn/kube-ovn/commit/78336126822582a071ecfcd7e4d52b7902bd2ff7) fix(deps): update module sigs.k8s.io/controller-runtime to v0.22.2 (#5768)
 * [c6ebbd1b0](https://github.com/kubeovn/kube-ovn/commit/c6ebbd1b01e55aea934760859c72a10639eab0ce) fix(deps): update module google.golang.org/protobuf to v1.36.10 (#5765)
 * [5cc703883](https://github.com/kubeovn/kube-ovn/commit/5cc703883d966d41f6cacae5a27583a8ede48e2b) fix(deps): update module kubevirt.io/client-go to v1.6.2 (#5764)
 * [daced4243](https://github.com/kubeovn/kube-ovn/commit/daced4243cb3fcb51a7bed8dc096f3df81c73712) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.26.0 (#5766)
 * [b155e1958](https://github.com/kubeovn/kube-ovn/commit/b155e19589dc7b075bc40e05861f784403c9e6a0) fix(deps): update module github.com/docker/docker to v28.5.0+incompatible (#5767)
 * [ba6564bab](https://github.com/kubeovn/kube-ovn/commit/ba6564bab59b8d71850604de259ee0c3a3a890c6) fix(deps): update module kubevirt.io/api to v1.6.2 (#5763)
 * [6b955719a](https://github.com/kubeovn/kube-ovn/commit/6b955719a335d4f637bedde1c1d432af27de1c56) chore(deps): update dependency kubevirt to v1.6.2 (#5762)
 * [6054370c5](https://github.com/kubeovn/kube-ovn/commit/6054370c527472951ee790c3884b68a133dec27b) docs: updated CHANGELOG.md (#5759)
 * [772490fee](https://github.com/kubeovn/kube-ovn/commit/772490fee5473b7c995fc7f4e12af6df42cd4612) fix gabage u2o openflow (#5757)
 * [112d9ae97](https://github.com/kubeovn/kube-ovn/commit/112d9ae9740e6b2d8819c4728af9b8ba28ce5459) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.7.3 (#5754)
 * [7e9d36fb9](https://github.com/kubeovn/kube-ovn/commit/7e9d36fb984dd45cf740c23f6b2b1e7909335c73) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.7.2 (#5753)
 * [f0070b3c1](https://github.com/kubeovn/kube-ovn/commit/f0070b3c1a46ef5d2816d2169e1ef9237483da89) docs: updated CHANGELOG.md (#5751)
 * [d6329453d](https://github.com/kubeovn/kube-ovn/commit/d6329453df4f357719e5f88a3287b4489b636c6f) check ovn0 mac should be set expected (#5748)
 * [2bfb92ff9](https://github.com/kubeovn/kube-ovn/commit/2bfb92ff965d2394c327a06b31bd378f9bb6a247) chore(deps): update dependency gosec to v2.22.9 (#5746)
 * [bb56bb675](https://github.com/kubeovn/kube-ovn/commit/bb56bb675a2a0b56b4682467de10029f72a2a0a5) chore(deps): update dependency golangci-lint to v2.5.0 (master) (#5743)
 * [b1ffea0aa](https://github.com/kubeovn/kube-ovn/commit/b1ffea0aad055ba3591a9a4852a3f7e89b20e302) docs: updated CHANGELOG.md (#5744)
 * [6b8b22027](https://github.com/kubeovn/kube-ovn/commit/6b8b22027e87b42adbbb5a0a717099dd2808b78f) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.7.1 (#5742)
 * [efd240f83](https://github.com/kubeovn/kube-ovn/commit/efd240f83b4368eef17c02c779f3407526f55ee7) controller: fix lsp not deleted in gc (#5735)
 * [67651345f](https://github.com/kubeovn/kube-ovn/commit/67651345fea5213e17ea00c857ea142a2ab31142) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.7.0 (#5740)
 * [f9e2e5aaa](https://github.com/kubeovn/kube-ovn/commit/f9e2e5aaa1b78ac3bd9d296ad1c3d7d15560fd83) fix(deps): update module kubevirt.io/client-go to v1.6.1 (#5739)
 * [a425de45d](https://github.com/kubeovn/kube-ovn/commit/a425de45d916cdb009ba7a4755e86fbd3ed3d261) fix(deps): update module kubevirt.io/api to v1.6.1 (#5738)
 * [13f6919d3](https://github.com/kubeovn/kube-ovn/commit/13f6919d3e2cf75dde0721fd67a43002ad391d20) chore(deps): update dependency kubevirt to v1.6.1 (#5737)
 * [e14c95e81](https://github.com/kubeovn/kube-ovn/commit/e14c95e81b844cc2507442d812abbc9978e0d7e8) update next.md
 * [d831cb5cb](https://github.com/kubeovn/kube-ovn/commit/d831cb5cb5fd69eaba56c0328697b75b13b28b1a) docs: updated CHANGELOG.md (#5734)
 * [779d9e10e](https://github.com/kubeovn/kube-ovn/commit/779d9e10e6adb6299bd775e1813e74153cf26259) support for vpc nat gateway without eip (#5605)
 * [554039c0f](https://github.com/kubeovn/kube-ovn/commit/554039c0ffa5971d911a3c2ea516547508981c2e) ipam: fix allocating from ippool (#5731)
 * [fbfe62fe6](https://github.com/kubeovn/kube-ovn/commit/fbfe62fe6d56838163851af7284a43deac5744ad) chore(deps): update dependency cilium to v1.18.2 (#5730)
 * [8af452f8c](https://github.com/kubeovn/kube-ovn/commit/8af452f8cb477e017f3d9993fd46e8fc9e0ff294) chore(deps): update actions/setup-go action to v6 (#5725)
 * [80f1244cb](https://github.com/kubeovn/kube-ovn/commit/80f1244cb967cb13145f21c624bf5d5d0d2916de) chore(deps): update dependency helm to v3.19.0 (#5721)
 * [564a29e4e](https://github.com/kubeovn/kube-ovn/commit/564a29e4e71cc163058825a89faee007391a99f9) controller: fix ip allocation from IPPools (#5729)
 * [6f157356a](https://github.com/kubeovn/kube-ovn/commit/6f157356a5b0a2ff9a1bdc2f424f4dc1ec45489b) fix metallb underlay lflow rule is deleted unexpected (#5723)
 * [6000a7a2a](https://github.com/kubeovn/kube-ovn/commit/6000a7a2a23c156c79e605716d3b2173b3880391) fix: e2e change test url to google.com (#5726)
 * [caad4debe](https://github.com/kubeovn/kube-ovn/commit/caad4debe5d1ab7da4135a63557484f82cda68af) add fqdn anq to next
 * [8a4895291](https://github.com/kubeovn/kube-ovn/commit/8a48952911363181656c33c04af6dbf525e7cf5c) fix(deps): update module github.com/puzpuzpuz/xsync/v4 to v4.2.0 (#5722)
 * [4156a2670](https://github.com/kubeovn/kube-ovn/commit/4156a267000efe4d35847cd506ac1cebb365e694) Add imple egress wall (#5703)
 * [2e8753af3](https://github.com/kubeovn/kube-ovn/commit/2e8753af31b8ec41c88261f3356d8df761223275) add github runner cache gc
 * [9d97fe78d](https://github.com/kubeovn/kube-ovn/commit/9d97fe78d79f440c2a98f258a2e22c717939e8d0) remove roadmap.md with NEWS.md to better show what you will see next. (#5720)
 * [c5a9fab3b](https://github.com/kubeovn/kube-ovn/commit/c5a9fab3b1b5ee5a5af6d43cb18a95cea99559e9) ci: fix go cache restore keys (#5719)
 * [49e1beb0b](https://github.com/kubeovn/kube-ovn/commit/49e1beb0bcc938e08dd182cdaf7a5bf889dc5759) add branch name to cache to increase hit rate
 * [57a3e74bf](https://github.com/kubeovn/kube-ovn/commit/57a3e74bf404e1e9e477a6dba6a0c9ec2abbc5bc) use arp for ipve network ready check (#5716)
 * [52e04ef20](https://github.com/kubeovn/kube-ovn/commit/52e04ef20f36f39a4dc0be624990a8cbc24c31fa) fix(deps): update kubernetes (master) (#5355)
 * [e9efd0f90](https://github.com/kubeovn/kube-ovn/commit/e9efd0f90c592b1be25ed005f519d5841ae89de1) fix(deps): update module golang.org/x/tools to v0.37.0 (#5712)
 * [ba6251a43](https://github.com/kubeovn/kube-ovn/commit/ba6251a436f95a37bd99123939a647248782e5c8) fix(deps): update module google.golang.org/grpc to v1.75.1 (#5711)
 * [a63a945c7](https://github.com/kubeovn/kube-ovn/commit/a63a945c7d9a10c5f7190a034fa00145a57ee731) IPSec: Request new certificate if current certificate is not trusted (#5710)
 * [4129c0a1a](https://github.com/kubeovn/kube-ovn/commit/4129c0a1a6e64cfa425c6b4cb4bf190f3cdb7b9d) Feat: Added non primary cni mode support for Kube-OVN (#5618)
 * [b9b1a055b](https://github.com/kubeovn/kube-ovn/commit/b9b1a055b85f8c29caadeae12265b03029c4e8be) gc ip by different types (#5702)
 * [f0440da65](https://github.com/kubeovn/kube-ovn/commit/f0440da652e75ddc669b1ea99be2473355afbe45) fix(deps): update module google.golang.org/protobuf to v1.36.9 (#5706)
 * [05ae8cb70](https://github.com/kubeovn/kube-ovn/commit/05ae8cb70f03e74a852e192505a8c7d92df11ee2) fix(deps): update module golang.org/x/net to v0.44.0 (#5708)
 * [396fb1e48](https://github.com/kubeovn/kube-ovn/commit/396fb1e48e20337e8f7a4f8c445c48dba50259ce) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.6.0 (#5707)
 * [6455b199a](https://github.com/kubeovn/kube-ovn/commit/6455b199abadb529c97f56c53e41112455672129) docs: updated CHANGELOG.md (#5701)
 * [87d49b1e4](https://github.com/kubeovn/kube-ovn/commit/87d49b1e45a47ab1f1fe2fe1cdd793ad4a23639f) veg: fix init failure when the a cidr applied is a sub network of the internal subnet (#5685)
 * [02abcbb18](https://github.com/kubeovn/kube-ovn/commit/02abcbb180ab1a6474a2ef5468f943adb2bcb371) controller: fix marking ip resource as reserved incorrectly (#5698)
 * [5b2e3850b](https://github.com/kubeovn/kube-ovn/commit/5b2e3850bb802676694968dd59b3d0309fbf958f) fix(deps): update golang (#5699)
 * [167c0a7a7](https://github.com/kubeovn/kube-ovn/commit/167c0a7a7e63878663f3295bac5642ca2423f246) fix(deps): update golang (#5693)
 * [f2452aa15](https://github.com/kubeovn/kube-ovn/commit/f2452aa157328b26728435e7d2fd99318048af73) fix(deps): update module github.com/prometheus/client_golang to v1.23.2 (#5689)
 * [710c03980](https://github.com/kubeovn/kube-ovn/commit/710c03980b736141b5a1e6ff88359d3f6626dcef) chore(deps): update actions/stale action to v10 (#5690)
 * [53c843442](https://github.com/kubeovn/kube-ovn/commit/53c843442cec8bfbc2ba4067b3cf9ac460c1d4e4) fix(deps): update module github.com/docker/docker to v28.4.0+incompatible (#5682)
 * [089897fed](https://github.com/kubeovn/kube-ovn/commit/089897feda164e35b7ea6338ca2d81239c748f69) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.25.3 (#5683)
 * [d5f2f50ec](https://github.com/kubeovn/kube-ovn/commit/d5f2f50ecaaa8fd8308651b283509b338c2fd3d0) chore(deps): update actions/setup-go action to v6 (#5684)
 * [98668d671](https://github.com/kubeovn/kube-ovn/commit/98668d671437a01289bcb0adfc8282cb5ba12856) fix(deps): update module github.com/prometheus/client_golang to v1.23.1 (#5681)
 * [2affa4368](https://github.com/kubeovn/kube-ovn/commit/2affa4368ecf3dc04e67fceff3ac8bbb508066da) fix e2e (#5678)
 * [554bdf057](https://github.com/kubeovn/kube-ovn/commit/554bdf05772a852d7980688cf792b845e969df67) Fix e2e (#5677)
 * [a869778f3](https://github.com/kubeovn/kube-ovn/commit/a869778f338f1622191100d30cdf1214230169d7) fix e2e for stuck when delete subnet  (#5676)
 * [2508e5cf8](https://github.com/kubeovn/kube-ovn/commit/2508e5cf84a491744b95e2707951b1c496806bb3) crd: fix unrecognized format int32/int64 (#5671)
 * [08b9e2e9e](https://github.com/kubeovn/kube-ovn/commit/08b9e2e9e8effd82d2e66586bee6793c56c39da8) chore(deps): update dependency go to v1.25.1 (#5673)
 * [098a105c8](https://github.com/kubeovn/kube-ovn/commit/098a105c8e0ccc63892f368e4b61e7caf139047e) chore(deps): update aquasecurity/trivy-action action to v0.33.1 (#5672)
 * [ae9d05d5c](https://github.com/kubeovn/kube-ovn/commit/ae9d05d5cc1b223c4b6aee6eb304b7246ceff2d6) fix(deps): update module github.com/spf13/pflag to v1.0.10 (#5670)
 * [17aa56893](https://github.com/kubeovn/kube-ovn/commit/17aa568930ef00495096959a07d07800bf634830) ipam: fix ippool using/free ips after subnet update (#5668)
 * [c83d9ba32](https://github.com/kubeovn/kube-ovn/commit/c83d9ba3202683c4df071e4f8e734460f87a39a2) fix(deps): update module github.com/containernetworking/plugins to v1.8.0 (#5667)
 * [f7f3c4207](https://github.com/kubeovn/kube-ovn/commit/f7f3c4207d49be37b67ac5108364042bb46d81d3) chore(deps): update dependency cni-plugin to v1.8.0 (#5666)
 * [290d36a7d](https://github.com/kubeovn/kube-ovn/commit/290d36a7d84282c6651e093499c5f8ca8b334432) fix(deps): update module github.com/spf13/pflag to v1.0.9 (#5664)
 * [13f267542](https://github.com/kubeovn/kube-ovn/commit/13f267542c4f2069f0c3a9fd4c3997a88c8d7f60) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.5.1 (#5663)
 * [c01f03585](https://github.com/kubeovn/kube-ovn/commit/c01f03585b0f806ee198668dc34c338e188721d3) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.25.2 (#5644)
 * [5a1e50b5e](https://github.com/kubeovn/kube-ovn/commit/5a1e50b5ea6d5b0470636aa57a5697bb73c86887) chore(deps): update aquasecurity/trivy-action action to v0.33.0 (#5659)
 * [a3faa8916](https://github.com/kubeovn/kube-ovn/commit/a3faa8916d1daabcf4e104eccc46befbbd971ff6) chore(deps): update dependency kind to v0.30.0 (#5657)
 * [06563e9f9](https://github.com/kubeovn/kube-ovn/commit/06563e9f9f6260cb29209ff8f89ea29d20dd9037) release address when acquire a new static address (#5658)
 * [405350d6d](https://github.com/kubeovn/kube-ovn/commit/405350d6db8daf8d2f83c2e48d545e84a8d16057) fix(deps): update module github.com/stretchr/testify to v1.11.1 (#5656)
 * [fdc1d6836](https://github.com/kubeovn/kube-ovn/commit/fdc1d6836dfffe1c3aa5669ee38748df6a8f009a) fix vpc-dns annotations update revision (#5655)
 * [04c9bd668](https://github.com/kubeovn/kube-ovn/commit/04c9bd668b0bad90e546745913e430d41d86b8d4) Increase timeout
 * [f27cc976b](https://github.com/kubeovn/kube-ovn/commit/f27cc976b1d2d15f9c01ad1e34a3ef911e24c364) fix(deps): update module github.com/onsi/gomega to v1.38.2 (#5652)
 * [1bf5983cc](https://github.com/kubeovn/kube-ovn/commit/1bf5983cc141b0d230b61f27b8aa1c6eb481306a) handle delete final state unknown object in enqueue handler (#5649)
 * [ba68d3977](https://github.com/kubeovn/kube-ovn/commit/ba68d3977fb754e3f10dc8205d7adba1240ae086) Initialize annotations immediately before assignment in setVpcDNSRoute (#5648)
 * [00572c662](https://github.com/kubeovn/kube-ovn/commit/00572c662ffa7a14b86d7acfeb9e1cb79dca5dea) fix bug for issuse #5597: when slr update, the same vip will be delete too (#5646)
 * [a0be52045](https://github.com/kubeovn/kube-ovn/commit/a0be520456ab51d62131e7e580133ab647dd6bf5) fix(deps): update module github.com/onsi/gomega to v1.38.1 (#5645)
 * [dfcd941e6](https://github.com/kubeovn/kube-ovn/commit/dfcd941e6ee2b38dac943e961545b61e5e222160) fix(deps): update module github.com/stretchr/testify to v1.11.0 (#5647)
 * [bdb199782](https://github.com/kubeovn/kube-ovn/commit/bdb199782e165f1606b053afd220ab7bac76ddf9) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.25.0 (#5639)
 * [c6778e67a](https://github.com/kubeovn/kube-ovn/commit/c6778e67afe727e4345d69ce38a402f5a07ac2d1) fix(deps): update module google.golang.org/grpc to v1.75.0 (#5640)
 * [b46f67e2d](https://github.com/kubeovn/kube-ovn/commit/b46f67e2dd816efeafbfbc30e017682928f16cac) fix(deps): update module google.golang.org/protobuf to v1.36.8 (#5635)
 * [324d57541](https://github.com/kubeovn/kube-ovn/commit/324d57541428364cbf7a6a035e2752f42dba631b) docs: updated CHANGELOG.md (#5638)
 * [dcadedc6e](https://github.com/kubeovn/kube-ovn/commit/dcadedc6eff477f743ffeeda83834ae8db338dcb) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.4.0 (#5636)
 * [a8b7f1fe1](https://github.com/kubeovn/kube-ovn/commit/a8b7f1fe11142928929eb60fc57f9ebcc1086143) docs: updated CHANGELOG.md (#5634)
 * [5899d02c7](https://github.com/kubeovn/kube-ovn/commit/5899d02c7aec12fca14ef925da14619c98c00ffe) docs: updated CHANGELOG.md (#5633)
 * [0559aff56](https://github.com/kubeovn/kube-ovn/commit/0559aff564d8d754f13735128e4b1ceba0063d8d)  fix bug for issuse #5597 (#5616)
 * [98d445100](https://github.com/kubeovn/kube-ovn/commit/98d445100c98c7700389e1e1b5b80337ae30054f) chore(deps): update azure/setup-helm action to v4.3.1 (#5627)
 * [cc211be6f](https://github.com/kubeovn/kube-ovn/commit/cc211be6fb0a307cdca4a67a32ad05213bb96094) chore(deps): update dependency helm to v3.18.6 (#5628)
 * [951e783c1](https://github.com/kubeovn/kube-ovn/commit/951e783c18d97d4f3e7f6b7f238ae4e8adf249f9) docs: updated CHANGELOG.md (#5629)
 * [f0d1fcd90](https://github.com/kubeovn/kube-ovn/commit/f0d1fcd90e6a66a520bd97b5093b46756dd60c31) fix ovn ipsec when restart cni (#5603)
 * [bd921124f](https://github.com/kubeovn/kube-ovn/commit/bd921124fb8fec5b0b6a374ba358779156183489) feat(vpcnatgw): add the option to specify custom routes (#5608)
 * [548900765](https://github.com/kubeovn/kube-ovn/commit/5489007654742ab4bd8fd96d94ad5f5f15853605) fix static mac pod conflict with gateway mac (#5623)
 * [572c56c39](https://github.com/kubeovn/kube-ovn/commit/572c56c3935c31a1e941d158a50b7bb2e460ec0b) e2e: skip testcases of ServiceCIDR and IPAddress API for versions prior to 1.13 (#5624)
 * [267fb0268](https://github.com/kubeovn/kube-ovn/commit/267fb0268a409e5bac4e090bcebc6b4e70874d4a) feat(vpcnatgw): send gratuitous arp for nexthops at natgw initialization (#5607)
 * [d0ea14be8](https://github.com/kubeovn/kube-ovn/commit/d0ea14be83297ff2ab574a957ae53723af221940) chore(deps): update dependency go to v1.25.0 (master) (#5592)
 * [9d8b1c310](https://github.com/kubeovn/kube-ovn/commit/9d8b1c3101f9f4bdfb576bb22ab1e2453596e3f8) fix(deps): update module github.com/httprunner/httprunner/v5 to v5.0.0-20250818135001-47996ed22fb1 (#5619)
 * [3676597aa](https://github.com/kubeovn/kube-ovn/commit/3676597aa31bf11ed4772fe6de11dac8b42b2a83) fix(deps): update module go.uber.org/mock to v0.6.0 (#5621)
 * [93d5ab00c](https://github.com/kubeovn/kube-ovn/commit/93d5ab00c0934a25c1d7211687c137b853df813a) fix(deps): update module github.com/onsi/ginkgo/v2 to v2.24.0 (#5620)
 * [cf7c8a67b](https://github.com/kubeovn/kube-ovn/commit/cf7c8a67b82ac74e2a65bff03048325c00201470) chore(deps): update dependency cilium to v1.18.1 (#5613)
 * [c43b18a8b](https://github.com/kubeovn/kube-ovn/commit/c43b18a8bac3ade5efc8053c0f8e1b6fff5a0fac) fix(deps): update module github.com/emicklei/go-restful/v3 to v3.13.0 (#5614)
 * [e671f5da7](https://github.com/kubeovn/kube-ovn/commit/e671f5da73a18033fbce15587590821e80936679) bump k8s to v1.33.4 (#5182)
 * [37fd04cdb](https://github.com/kubeovn/kube-ovn/commit/37fd04cdbb24c75a6c769eca565228f0f847adcf) base: bump dpdk to 24.11.2 (LTS) (#5195)
 * [d7065e9a3](https://github.com/kubeovn/kube-ovn/commit/d7065e9a363783ba452ca56ed5e1e42cef7b7e8d) chore(deps): update dependency submariner to v0.21.0 (#5601)
 * [011bfaf49](https://github.com/kubeovn/kube-ovn/commit/011bfaf49e09ac8baf428b698bc367893342bee2) chore(deps): update dependency golangci-lint to v2.4.0 (#5610)
 * [3029f702d](https://github.com/kubeovn/kube-ovn/commit/3029f702d226e8fc6cc2e84ad2be4e90b437cac7) chore(deps): update dependency gosec to v2.22.8 (#5609)
 * [4337d077f](https://github.com/kubeovn/kube-ovn/commit/4337d077fdc5578ec5b2dda4840af87af2c64343) ci: ignore goveralls failure (#5606)
 * [d19a039d0](https://github.com/kubeovn/kube-ovn/commit/d19a039d001db1f570112c58b386b6785ac26ee4) chore(deps): update dependency helm to v3.18.5 (#5600)
 * [e22cb6af5](https://github.com/kubeovn/kube-ovn/commit/e22cb6af5ef569e487e6f99f97f87552b03a3b05) bump libovsdb to v0.8.1 (#5596)
 * [490794ce5](https://github.com/kubeovn/kube-ovn/commit/490794ce5848367bd6d37f6bc89123500e5e9803) simple vip (#5534)
 * [60deab89d](https://github.com/kubeovn/kube-ovn/commit/60deab89de8b836fec390ce26053c33b49da9e81) fix dpdk base image (#5595)
 * [8a7ef6e33](https://github.com/kubeovn/kube-ovn/commit/8a7ef6e33b3b80982d86a91442ccfd971deb256d) add backup script
 * [ab4ebc3c1](https://github.com/kubeovn/kube-ovn/commit/ab4ebc3c17e1007ebc6b84dbc825c7caa3d16871) reduce probability of maps conncurrent iteration and write (#5585)
 * [969154989](https://github.com/kubeovn/kube-ovn/commit/9691549893b3a9ff392acfe16189529eb0784d28) remove lsp when gw nodes change (#5591)
 * [c2a0e7bce](https://github.com/kubeovn/kube-ovn/commit/c2a0e7bce6b8253a9cb85a601c92cb8a69d18f84) chore(deps): update actions/checkout action to v5 (#5589)
 * [84f45f67a](https://github.com/kubeovn/kube-ovn/commit/84f45f67ad447c2e5b6158bf6659fce8dd7d8611) Support healthchecks for static endpoints in services (#5435)
 * [7bf3c404e](https://github.com/kubeovn/kube-ovn/commit/7bf3c404eab65063df707e53f074e0d2f55509a6) docs: updated CHANGELOG.md (#5586)
 * [ce2cfa90d](https://github.com/kubeovn/kube-ovn/commit/ce2cfa90df3c165d8b6a821e5aa0654fe0280cce) upgrade ovn to 25.03 (#5537)
 * [377fbd864](https://github.com/kubeovn/kube-ovn/commit/377fbd8647873290a6cf4ab6edcc3c84dd254704) ci: rebuild gobgp for arm64 (#5583)
 * [b26790b15](https://github.com/kubeovn/kube-ovn/commit/b26790b1555b651b6c2e5778e254cc308d17d34d) Use CA key name when scheduling certificate refresh (#5578)
 * [20af5be42](https://github.com/kubeovn/kube-ovn/commit/20af5be42b6db1a484f9c3aebe3740989294a2fe) chore(deps): update dependency go to v1.24.6 (#5582)
 * [d9f9932b4](https://github.com/kubeovn/kube-ovn/commit/d9f9932b46eebb318cd79de337bc68ae75b2eaa3) enable go version update
 * [b4edeef87](https://github.com/kubeovn/kube-ovn/commit/b4edeef87a55aeefc4ef97850b0f1507efd4772b) fix(deps): update golang (#5577)
 * [9f5ac84ac](https://github.com/kubeovn/kube-ovn/commit/9f5ac84ac146b7f365209a1b8cd618cfd4c6441b) fix(deps): update module google.golang.org/protobuf to v1.36.7 (#5576)
 * [e35b0fa89](https://github.com/kubeovn/kube-ovn/commit/e35b0fa89ccb3ef188fe2587d3dad6bfe247d4f7) vm static ip validate should use vm name (#5569)
 * [d9473a3ef](https://github.com/kubeovn/kube-ovn/commit/d9473a3ef3cb67254d062757f98e116e6bcbfb05) Fix the problem that if available ip is 0 but there is a value in excludeIPs, the fixed ip is used as the ip in excludeIPs but the error noAddressAvaliable is still reported. (#5561)
 * [edfe8afb9](https://github.com/kubeovn/kube-ovn/commit/edfe8afb90de59bb36b14b178f77954094dda0a7) chore(deps): update actions/download-artifact action to v5 (#5562)
 * [3e6dfd398](https://github.com/kubeovn/kube-ovn/commit/3e6dfd398f3069081ee4af522a8f3a54a0211e7c) chore(deps): update dependency go to v1.24.6 (#5563)
 * [f03180c43](https://github.com/kubeovn/kube-ovn/commit/f03180c4369f10fc8a5b1997e2e951e5e4589af6) Do not garbage collect IPs of stopped VMs using non-default multus networks (#5557)
 * [325b583df](https://github.com/kubeovn/kube-ovn/commit/325b583df354c2f6e1d4b499f699072cdadc9778) skip mv if config file already exists (#5558)
 * [13b7aec43](https://github.com/kubeovn/kube-ovn/commit/13b7aec4338ae34af3fa56939bb46f32d48085b2) skip nad without config (#5560)
 * [9b4425eed](https://github.com/kubeovn/kube-ovn/commit/9b4425eedd77f43674af491faa4553047f692762) Refactor Makefile structure (#5555)
 * [b31b02a5e](https://github.com/kubeovn/kube-ovn/commit/b31b02a5e4df1d14ed366ac8b98ebfc9053232bc) remove unnecessary files
 * [788694faf](https://github.com/kubeovn/kube-ovn/commit/788694faffd221e8e548624911483ae6a0235494) tolerate duplicate ACLs (#5552)
 * [8b2e6ed63](https://github.com/kubeovn/kube-ovn/commit/8b2e6ed6393649334f7916c543abdbbc5c467e5b) fix(helm): static master ips wouldn't work anymore (#5554)
 * [351401009](https://github.com/kubeovn/kube-ovn/commit/351401009eb57f56fff8ebca733817690a4c6099) disable external sync by default
 * [853869727](https://github.com/kubeovn/kube-ovn/commit/853869727cbf8918c7db87cf0915491d6528b9ca) do not handle update external vpcs
 * [7a9d00649](https://github.com/kubeovn/kube-ovn/commit/7a9d0064928c366211e57f90b1dd85e77f4a389b) chore(deps): update dependency golangci-lint to v2.3.1 (#5550)
 * [41f0f1fae](https://github.com/kubeovn/kube-ovn/commit/41f0f1fae0b18fb2ca246b3a38928e636c02237f) fix nat-gateway arping (#5546)
 * [d30cf16ae](https://github.com/kubeovn/kube-ovn/commit/d30cf16ae6945ccdbf6cd7eabb10b2dd955a69bf) fix vpcNatGw status cannot be updated (#5547)
 * [5a5fe24bc](https://github.com/kubeovn/kube-ovn/commit/5a5fe24bccabaa5ce3bbd5384aa8baa70a5a31cb) fix(deps): update module kubevirt.io/client-go to v1.6.0 (#5545)
 * [303b42fac](https://github.com/kubeovn/kube-ovn/commit/303b42fac7a721ed670c071a4b520c7dbc64d9ba) fix(deps): update module kubevirt.io/api to v1.6.0 (#5544)
 * [7fbcc5821](https://github.com/kubeovn/kube-ovn/commit/7fbcc58217bf8d176990a67e6b18a5d4c074beb0) fix(deps): update module github.com/prometheus/client_golang to v1.23.0 (#5540)
 * [9b85f8021](https://github.com/kubeovn/kube-ovn/commit/9b85f8021defe31fdd5cd531eebaa2e66f7cbbe8) chore(deps): update dependency kubevirt to v1.6.0 (#5539)
 * [0aa4c0026](https://github.com/kubeovn/kube-ovn/commit/0aa4c0026c9a049c8caf613abc2505c5869a334f) fix(deps): update module github.com/containerd/containerd/v2 to v2.1.4 (#5531)
 * [e175f73c1](https://github.com/kubeovn/kube-ovn/commit/e175f73c10c30c8564d955d8589ec922aa234dd2) Handle missing files gracefully in metric functions by returning early on os.IsNotExist errors (#5529)
 * [91cbc9cde](https://github.com/kubeovn/kube-ovn/commit/91cbc9cde56d2cff240237cce91335041dc950dc) support nodeselector for providerNetwork (#5518)
 * [821a2b2e2](https://github.com/kubeovn/kube-ovn/commit/821a2b2e247c7b442f67e1beaab29550d76ee448) chore(deps): update dependency cilium to v1.18.0 (#5525)
 * [4a3774231](https://github.com/kubeovn/kube-ovn/commit/4a3774231ae45f40ab096c6f7773df03bae151ae) fix(deps): update module github.com/docker/docker to v28.3.3+incompatible [security] (#5527)
 * [d3b5ecb3a](https://github.com/kubeovn/kube-ovn/commit/d3b5ecb3a9dd11cfb3a394a6995940f7726f0b9a) check underlay nic exist before config external bridge (#5520)
 * [0c170715c](https://github.com/kubeovn/kube-ovn/commit/0c170715c01a5ae5f8f83be2ac8015bdc4c72f56) chore(deps): update dependency multus to v4.2.2 (#5516)
 * [0361c59cb](https://github.com/kubeovn/kube-ovn/commit/0361c59cb2f44289e1e5f77963e09fbce05f7769) fix(deps): update module gopkg.in/k8snetworkplumbingwg/multus-cni.v4 to v4.2.2 (#5517)
 * [430f28a2d](https://github.com/kubeovn/kube-ovn/commit/430f28a2d0592ccbaf925bf3314909240a2eb972) fix(deps): update module github.com/cenkalti/backoff/v5 to v5.0.3 (#5513)
 * [6784a1ae1](https://github.com/kubeovn/kube-ovn/commit/6784a1ae180f6abd2d9b751c65710d4c2068f31a) docs: updated CHANGELOG.md (#5512)
 * [4ba008512](https://github.com/kubeovn/kube-ovn/commit/4ba008512d15b8c442060312134ea20e7ba9f3a6) fix NB Global not updated after OVN IC is disabled (#5511)
 * [5f150bc48](https://github.com/kubeovn/kube-ovn/commit/5f150bc4838b6bbeadc2b22838fcedee09503551) fix concurrent map race (#5510)
 * [7d39c9ee4](https://github.com/kubeovn/kube-ovn/commit/7d39c9ee45ce97d54721e7713d9dbcf1f13fcc37) fix(deps): update module github.com/onsi/gomega to v1.38.0 (#5509)
 * [6d9d57a3d](https://github.com/kubeovn/kube-ovn/commit/6d9d57a3d32093f8d66c640652caadc7be558ad2) fix(deps): update module google.golang.org/grpc to v1.74.2 (#5508)
 * [beee2e474](https://github.com/kubeovn/kube-ovn/commit/beee2e474d73628b93c06297f7c82dc5e4553799) chore(deps): update kindest/node docker tag to v1.33.2 (#5507)
 * [dbf373b00](https://github.com/kubeovn/kube-ovn/commit/dbf373b0057706664bc3206c2f99438b527b3c6c) chore(config): migrate config renovate.json (#5506)
 * [d085b310d](https://github.com/kubeovn/kube-ovn/commit/d085b310d0de83fc49c201fc4c1c0c5523fe97b9) fix cleanup stuck (#5505)
 * [3602e3adb](https://github.com/kubeovn/kube-ovn/commit/3602e3adba7d6d01ebedbc1d7d3db9caeb67b21c) panic install if labeled nodes are not found (#5504)
 * [e90f15785](https://github.com/kubeovn/kube-ovn/commit/e90f15785475e3b8cd364dbc1b22a856aa2aaeb8) fix: ip/vpc are not deleted by helm uninstall (#5501)
 * [ec8dd6f69](https://github.com/kubeovn/kube-ovn/commit/ec8dd6f69444ac656ea86fcda41741c8a77b4b0b) chore(deps): update dependency gosec to v2.22.7 (#5502)
 * [93f55db2d](https://github.com/kubeovn/kube-ovn/commit/93f55db2dbd04c39a26c596d6d1d51fd3dfa1d46) chore(deps): update dependency golangci-lint to v2.3.0 (#5503)
 * [5fdcffda8](https://github.com/kubeovn/kube-ovn/commit/5fdcffda81f46b2a17fafeece8c39d024b5e1489) chart: fix node selector (#5500)
 * [e4827e683](https://github.com/kubeovn/kube-ovn/commit/e4827e683e7f262c6c7446a781815bffa5383df1) fix(deps): update module google.golang.org/grpc to v1.74.1 (#5494)
 * [bcd7b57e2](https://github.com/kubeovn/kube-ovn/commit/bcd7b57e2e52f8f9b1bb3c780f82c0bd9d5de292) docs: updated CHANGELOG.md (#5497)
 * [09e130612](https://github.com/kubeovn/kube-ovn/commit/09e1306123653f7f3a462568dd2ee717743e616b) metrics: fix pinger_inconsistent_port_binding (#5496)
 * [a29592574](https://github.com/kubeovn/kube-ovn/commit/a29592574d142fc18f72090b1a578eb0018f3424) fix: only addOrUpdateSubnetQueue if the GatewayType is distributed instead of if it's not centralized (#5493)
 * [b9b48f7c4](https://github.com/kubeovn/kube-ovn/commit/b9b48f7c450e1a1dd6410fd086ff5d0d637c0b4c) bump k8s to v1.32.7 (#5492)
 * [9d39ee539](https://github.com/kubeovn/kube-ovn/commit/9d39ee539987a7b0a175753d333df6b8de096594) feat(endpoint_slice): refactor and support arbitrary providers (#5414)
 * [d7f63738e](https://github.com/kubeovn/kube-ovn/commit/d7f63738e3882c55d674f7f141056fc7972bf841) fix(deps): update module google.golang.org/grpc to v1.74.0 (#5491)
 * [5cca65ca4](https://github.com/kubeovn/kube-ovn/commit/5cca65ca4233018eb6b0ff74ff20ac57b0ed9a46) chore(deps): update dependency cilium to v1.17.6 (#5489)
 * [2ed3d0615](https://github.com/kubeovn/kube-ovn/commit/2ed3d06155a3bf22ef43fd65d729f89245953839) some fixes in ipsec backed with cert-manager (#5445)
 * [3e61e807a](https://github.com/kubeovn/kube-ovn/commit/3e61e807a6a37e5458ff5f277b0569d9e9f709ec) fix missing rbac rule for sa ovn-ovs (#5488)
 * [1039e1ae5](https://github.com/kubeovn/kube-ovn/commit/1039e1ae52e1bb953e655e2dfc908c471e628419) move kind node image to ghcr to avoid throttling (#5486)
 * [a4bf2baf1](https://github.com/kubeovn/kube-ovn/commit/a4bf2baf1a6c74420b3f3d4bd13e4759cc984ae8) move image to ghcr to avoid dockerhub throttle (#5485)
 * [cbb6e9cf5](https://github.com/kubeovn/kube-ovn/commit/cbb6e9cf54ed97e4a6a2b2cc236dbfb0214dba57) fix iptables nat gw e2e failed (#5484)
 * [f5ff076ee](https://github.com/kubeovn/kube-ovn/commit/f5ff076ee00e23e8216044933dfe93a23a8313e0) fix ipsec e2e (#5474)
 * [d08e3b655](https://github.com/kubeovn/kube-ovn/commit/d08e3b6556f51c7e3355232d495212e64dee82bf) sync chart version (#5483)
 * [5261c1a7f](https://github.com/kubeovn/kube-ovn/commit/5261c1a7f666ee6b94148f7cafb488718609106b) add adopters section
 * [da67e6237](https://github.com/kubeovn/kube-ovn/commit/da67e623712e3cfe44908002adcc1baafe123ced) add release image cleanup
 * [42af06dc8](https://github.com/kubeovn/kube-ovn/commit/42af06dc86db0a02481c52e5f79005e444fe4fdf) fix docs release
 * [83f3596df](https://github.com/kubeovn/kube-ovn/commit/83f3596df4c2ac07ef8a6935f8433f2b533cee32) docs: updated CHANGELOG.md (#5482)
 * [d7638b067](https://github.com/kubeovn/kube-ovn/commit/d7638b067da3251e2cbbbe8cdfaceadd5fccc2a2) controller: fix setting LR options on initialization (#5444)
 * [ec170958d](https://github.com/kubeovn/kube-ovn/commit/ec170958ddce1a0b3740a953065f86043f4c62ac) export ovs functions Get/Set/Remove (#5473)
 * [dea713cda](https://github.com/kubeovn/kube-ovn/commit/dea713cdabf7c461391ceca45096d1e5aff2d6f3) kubectl-ko: collect information about ipsec and xfrm (#5472)
 * [9ae6bb9c4](https://github.com/kubeovn/kube-ovn/commit/9ae6bb9c4d9c5d0636a3ada48488982306dae5b3) chore(deps): update dependency golangci-lint to v2.2.2 (#5475)
 * [5f3cd36d3](https://github.com/kubeovn/kube-ovn/commit/5f3cd36d34c7e6b5b5830a1227d3950d11c12ad4) fix(deps): update module golang.org/x/tools to v0.35.0 (#5476)
 * [6994d0e48](https://github.com/kubeovn/kube-ovn/commit/6994d0e48aa3669ec74e337b47ba51a2c50ebfcf) allow different provider use same vlan (#5471)
 * [6819cad37](https://github.com/kubeovn/kube-ovn/commit/6819cad37834bdfe6ce1ffe1cb62ef7f95219834) remove already has debs (#5447)
 * [7a808d153](https://github.com/kubeovn/kube-ovn/commit/7a808d1539096cd5ce5edb53eba03e478bb3338e) fix(deps): update module golang.org/x/net to v0.42.0 (#5467)
 * [836445e7c](https://github.com/kubeovn/kube-ovn/commit/836445e7c4aec8e3dfded0afd607a35d1c202bd1) feat: add pod_name and pod_namespace labels to interface metrics (#5463)
 * [1bb883261](https://github.com/kubeovn/kube-ovn/commit/1bb8832615b2d2a049c29e95a9fa4ff4da1769d5) fix(deps): update module golang.org/x/sys to v0.34.0 (#5461)
 * [bcb7166e1](https://github.com/kubeovn/kube-ovn/commit/bcb7166e1b20273d8632e7c5d4ad61dc0e8e7be5) group golang deps
 * [9e4ed4cad](https://github.com/kubeovn/kube-ovn/commit/9e4ed4cad8981e051506719150ee98feb3c737a1) fix(deps): update module golang.org/x/mod to v0.26.0 (#5450)
 * [6dc6a0740](https://github.com/kubeovn/kube-ovn/commit/6dc6a074059981cd00b7c01f98ac728169d8dcfd) fix(deps): update module github.com/docker/docker to v28.3.2+incompatible (#5449)
 * [7916c680f](https://github.com/kubeovn/kube-ovn/commit/7916c680feb8812e6719ef0bd227dcec974d7db6) refactor label selector (#5442)
 * [06d9117bc](https://github.com/kubeovn/kube-ovn/commit/06d9117bc616b59c0f303b11348bade6b54c7e53) chore(deps): update dependency helm to v3.18.4 (#5440)
 * [ba921323e](https://github.com/kubeovn/kube-ovn/commit/ba921323e75b9b93aa9b06aba58987a9ff0f58d5) chore(deps): update dependency go to v1.24.5 (master) (#5438)
 * [c880c8359](https://github.com/kubeovn/kube-ovn/commit/c880c8359ea5a3d3fd1e9cd5febc815b81d1a0d2) Use cert-manager to issue certificates for IPSec (#5365)
 * [8d4a8adb1](https://github.com/kubeovn/kube-ovn/commit/8d4a8adb1fb6b51561335182a7deb7f5f1f84bdb) ci: increase timeout of Multus-CNI conformance e2e (#5429)
 * [90a8c1b71](https://github.com/kubeovn/kube-ovn/commit/90a8c1b717ac2685aee1e0d11d7dd6cd0ebe82d7) replace Endpoint with EndpointSlice (#5437)
 * [5d6e9e7dc](https://github.com/kubeovn/kube-ovn/commit/5d6e9e7dc537d2b68cd8b013047abddab1946d6d) rbac: fix missing verbs (#5436)
 * [b16d5f387](https://github.com/kubeovn/kube-ovn/commit/b16d5f387b0d31af33d1e2260b047cfa0382a107) reorder rules to skip digest update
 * [56b36b180](https://github.com/kubeovn/kube-ovn/commit/56b36b18032f004288558a121d6226850449a2de) chore(deps): update aquasecurity/trivy-action action to v0.32.0 (#5433)
 * [6637dcbbf](https://github.com/kubeovn/kube-ovn/commit/6637dcbbf6b91938204ad9b667e74c5d8a505c62) docs: updated CHANGELOG.md (#5430)
 * [4605cd70b](https://github.com/kubeovn/kube-ovn/commit/4605cd70b9c7911a091685503859cba94bb34477) add hack script for cherry pick
 * [a5a200e81](https://github.com/kubeovn/kube-ovn/commit/a5a200e811abbf29ae7905bfedfdf5bf7b768114) fix setting LR option always_learn_from_arp_request (#5426)
 * [3ae1227f2](https://github.com/kubeovn/kube-ovn/commit/3ae1227f2d9452461b28c23816064dfcc133fd33) fix(deps): update module github.com/docker/docker to v28.3.1+incompatible (#5428)
 * [328b52899](https://github.com/kubeovn/kube-ovn/commit/328b528990edb11cc7c550b7e8142219a5d2e5f1) This allows caching of network attachment definitions and should improve (#5386)
 * [489b1eb89](https://github.com/kubeovn/kube-ovn/commit/489b1eb89f4f97dc07a2ad9e7cbbcf7d0f773e83) kubectl-ko: replace endpoint with endpointslice (#5425)
 * [266d3031a](https://github.com/kubeovn/kube-ovn/commit/266d3031a07e48892bcc655d06d92899b2c5f502) controller: set always_learn_from_arp_request to false only when LR is not connected to external network (#5419)
 * [dfb8f5a96](https://github.com/kubeovn/kube-ovn/commit/dfb8f5a96ef8a2a8b943124d7b6ac8d18372f3f5) fix parsing resolv.conf when systemd-resolved is running on the host (#5423)
 * [8cc5029b7](https://github.com/kubeovn/kube-ovn/commit/8cc5029b7ba32c329530b0ef60c28bed39380181) cni-server: fix inserting and deleting iptables rules (#5421)
 * [0109369f8](https://github.com/kubeovn/kube-ovn/commit/0109369f88005405cf5b1e8dcf7ff21a526984c6) chore(deps): update dependency cert-manager to v1.18.2 (#5422)
 * [d6a496a80](https://github.com/kubeovn/kube-ovn/commit/d6a496a806b8fbb3ad47dd10b486b5a8a9145eee) add manual release
 * [67f08f88f](https://github.com/kubeovn/kube-ovn/commit/67f08f88fef0eaf1efdbd8d200e8f7bdf9c0756c) update ovnv2 chart name (#5415)
 * [99acf9fee](https://github.com/kubeovn/kube-ovn/commit/99acf9feef5d1832c01c41e05d05d15c77bc9263) mock: generate code using go tool command (#5409)
 * [f4d1ca60e](https://github.com/kubeovn/kube-ovn/commit/f4d1ca60e1c5cb16231e022e6ee895e42174a226) track kind node version and release go version (#5412)
 * [f4d099824](https://github.com/kubeovn/kube-ovn/commit/f4d0998241d015722aeb10711aa3eb86e98f8ee4) chore(deps): update dependency golangci-lint to v2.2.1 (master) (#5407)
 * [90f2802a1](https://github.com/kubeovn/kube-ovn/commit/90f2802a1e90fabe9af39cabb8a12b505395340b) fix(deps): update module github.com/k8snetworkplumbingwg/network-attachment-definition-client to v1.7.7 (#5406)
 * [a7c6b543e](https://github.com/kubeovn/kube-ovn/commit/a7c6b543ea63c61a42570c06e184d61486f80061) update badges
 * [f39ec613a](https://github.com/kubeovn/kube-ovn/commit/f39ec613ad2b0c90443b7ff533a29500ef32d03a) simplify readme
 * [3d9290ccb](https://github.com/kubeovn/kube-ovn/commit/3d9290ccb5fc5be5bac98fe9d182ad2c4059c8c4) feat(endpoints): auto-detect provider to adjust ip mapping (#5404)
 * [024fa773f](https://github.com/kubeovn/kube-ovn/commit/024fa773f735d25c2fd57ca2634051bb45c2dea9) fix(deps): update module github.com/brianvoe/gofakeit/v7 to v7.3.0 (#5405)
 * [258200acc](https://github.com/kubeovn/kube-ovn/commit/258200accef6eb4856f704068df9da1c7c2b640c) fix matching loopback addresses (#5402)
 * [186a66e6d](https://github.com/kubeovn/kube-ovn/commit/186a66e6dbc8c086504bebd14e9a24abfe4a15c8) ci: deploy ipv6 talos cluster without ipv4 addresses (#5401)
 * [dfe951bb5](https://github.com/kubeovn/kube-ovn/commit/dfe951bb50ec1e437f6d426975a6d77f8a994c10) fix(deps): update module github.com/docker/docker to v28.3.0+incompatible (#5396)
 * [4f3be8921](https://github.com/kubeovn/kube-ovn/commit/4f3be8921dc3dafb20764f4eb0a450651b2d27f9) fix(services): SLRs could not target VM endpoints (#5397)
 * [f7c857c52](https://github.com/kubeovn/kube-ovn/commit/f7c857c522b181d46b0ad169eb585cd65233bdf9) del and add acls in one transaction (#5394)
 * [1f9a8e234](https://github.com/kubeovn/kube-ovn/commit/1f9a8e23429ea3797c61a000912b6c6f3c3ea8b9) Update USERS.md (#5399)
 * [cb85f9045](https://github.com/kubeovn/kube-ovn/commit/cb85f9045b03fe64d58b1e32e5eaeb862f366875) chore(deps): update docker.io/kubeovn/kube-ovn docker tag to v1.15.0 (#5398)
 * [537685f37](https://github.com/kubeovn/kube-ovn/commit/537685f374b8a6971a5f9aec7a5190ee1c6f4aa0) fix upstream not found
 * [8726274d1](https://github.com/kubeovn/kube-ovn/commit/8726274d15b45de75be9cdd8673149640a6188e7) docs: updated CHANGELOG.md (#5392)
 * [a3a957c79](https://github.com/kubeovn/kube-ovn/commit/a3a957c7989570605037a2a6db11af812dfff773) prepare for next release
 * [85a45ce30](https://github.com/kubeovn/kube-ovn/commit/85a45ce306585ce24e2ee23a752df850a81c034c) release vm ip when vm is stopped then deleted (#5390)
 * [904ac9afd](https://github.com/kubeovn/kube-ovn/commit/904ac9afde8e00ce23815b208072fd8fe7a2e719) wrong code (#5391)
 * [8194e60ff](https://github.com/kubeovn/kube-ovn/commit/8194e60ff6a4146c8d4a797329773f714372b2b3) vpc egress gateway: fix matching vpc (#5389)
 * [936150dd2](https://github.com/kubeovn/kube-ovn/commit/936150dd2f03c737037d9a25a1acf6788161128a) fix migrate down time (#5388)
 * [61f577ff3](https://github.com/kubeovn/kube-ovn/commit/61f577ff313094db4dc83d4501c2d6aa5492ff9d) fix duplicate acls because of parentkey (#5357)
 * [94c514e73](https://github.com/kubeovn/kube-ovn/commit/94c514e73bd688299124918b523c5a0d1549ad95) fix(slr): support switchlbrules targeting subnets with non-default providers (#5376)
 * [a2db6b178](https://github.com/kubeovn/kube-ovn/commit/a2db6b1784df4a1e126f178042b87091ac76724a) u2o keep src flow should consider multi-vlan on the same provider (#5385)
 * [ea758a79d](https://github.com/kubeovn/kube-ovn/commit/ea758a79dfffd9f8c707fa6778b61b1957f592e2) bump k8s to v1.32.6 (#5383)
 * [c386c7cdf](https://github.com/kubeovn/kube-ovn/commit/c386c7cdffbb65596525f5df1c4ccbac246f37f3) base: fix login shell of user sync (#5384)
 * [afdec4c35](https://github.com/kubeovn/kube-ovn/commit/afdec4c358172125a76a57a676d106fc635c86fb) fix(slr): deleting old entries from ipmapping to avoid clogging loadbalancer (#5380)
 * [5a66b6626](https://github.com/kubeovn/kube-ovn/commit/5a66b6626f2cf01203e5575f9767349ff5a923c7) feature: perform duplicate address detection if ipv6 address has a dadfailed flag (#5156)
 * [22ba1c035](https://github.com/kubeovn/kube-ovn/commit/22ba1c035359c3e710d6e9d70fe4c928ede65ff2) base: set login shell to /usr/sbin/nologin (#5382)
 * [f58127989](https://github.com/kubeovn/kube-ovn/commit/f58127989303570e385bf91a530806ad3ee8dd2a) chore(deps): update dependency cert-manager to v1.18.1 (#5373)
 * [7fb6efd88](https://github.com/kubeovn/kube-ovn/commit/7fb6efd88624eb8073b31c473a9ce49d1997c87c) fix(deps): update module github.com/containerd/containerd/v2 to v2.1.3 (#5381)
 * [79a799646](https://github.com/kubeovn/kube-ovn/commit/79a799646785ff265be2b64247e545d8d4a65e16) chore(deps): update dependency cilium to v1.17.5 (#5379)
 * [28ea8670b](https://github.com/kubeovn/kube-ovn/commit/28ea8670b810c95876b3e7e5154d64f43e84485f) ci: increase multus-cni resource limits (#5378)
 * [46c8045eb](https://github.com/kubeovn/kube-ovn/commit/46c8045eb782e0f17bd5cae8f1b3899f7a8637ea) fix(slr): switchlbrule doesn't support multi-homed/ipv6-first pods (#5375)
 * [016f52a01](https://github.com/kubeovn/kube-ovn/commit/016f52a01094b5d000429a8783ed3806dc77f253) vpc egress gateway: drop traffic if no nexthop is available (#5370)
 * [5e4e20aaf](https://github.com/kubeovn/kube-ovn/commit/5e4e20aaf8d2a75c31f31d6b670a1781b85c49a5) base: remove unused packages (#5368)
 * [501b4e8e1](https://github.com/kubeovn/kube-ovn/commit/501b4e8e1828ac58dacc5713fe638ea7c892e211) fix ip cr not delete with pod in the case of subnet not exist (#5364)
 * [2aeea8629](https://github.com/kubeovn/kube-ovn/commit/2aeea8629eafa13cfd65734c48f8bd98452b25c1) u2o keep src mac (#5192)
 * [9d398faaa](https://github.com/kubeovn/kube-ovn/commit/9d398faaac50a9ddab1c75e2d570b725d13caad0) fix: add validation to prevent subnet and VPC with the same name (#5371)
 * [574b33a2f](https://github.com/kubeovn/kube-ovn/commit/574b33a2ff31f3c883eeefc76dda034824423a06) chore(deps): update dependency helm to v3.18.3 (#5366)
 * [2865d3db8](https://github.com/kubeovn/kube-ovn/commit/2865d3db83affbbef0f3484a34d26e145fdda400) fix modernize lint (#5363)
 * [06540ebad](https://github.com/kubeovn/kube-ovn/commit/06540ebadd01bd257e890bb0b993147da9d5475b) chore(deps): update dependency gosec to v2.22.5 (#5362)
 * [29c25ab0c](https://github.com/kubeovn/kube-ovn/commit/29c25ab0cbd1af24f750b7d11ba2c6cab8a9117a) docs: updated CHANGELOG.md (#5361)
 * [33ff4cd73](https://github.com/kubeovn/kube-ovn/commit/33ff4cd73d9313ece91b1990ded038373510fb2a) fix(deps): update module github.com/containerd/containerd/v2 to v2.1.2 (#5354)
 * [1fceea23e](https://github.com/kubeovn/kube-ovn/commit/1fceea23ede605698a589f679b91cd28c348b25b) fix(slr): address family and familypolicy isn't correct (#5349)
 * [49bbed8d3](https://github.com/kubeovn/kube-ovn/commit/49bbed8d320e937bb23ceeab1a5315f74a86ee2d) controller: migrate acl tier after upgrade (#5351)
 * [82ef65ff9](https://github.com/kubeovn/kube-ovn/commit/82ef65ff99c9a993e5cf39011fc122017e0d9d8b) fix vpc egress gateway not applied to new pods (#5348)
 * [c40c6e38b](https://github.com/kubeovn/kube-ovn/commit/c40c6e38b78001da68dbb9c394120c565940be55) use group to reduce noisies
 * [cab9fa964](https://github.com/kubeovn/kube-ovn/commit/cab9fa964c7c7b0b693e0b5037e19d6376d7171e) add renovate label and ignore windows upgrade
 * [ddaf03eee](https://github.com/kubeovn/kube-ovn/commit/ddaf03eee7744a565cf5a4ee603ff3e22c91cef2) add return for enqueueUpdateIP error (#5350)
 * [1b4f2f558](https://github.com/kubeovn/kube-ovn/commit/1b4f2f5587935233d2dfc20133d65de7ef2e87bc) fix(deps): update module kubevirt.io/client-go to v1.5.2 (#5340)
 * [557f68652](https://github.com/kubeovn/kube-ovn/commit/557f686528c53346d5ec7bffb49c44b0665a7e8e) chore(deps): update dependency kubevirt to v1.5.2 (#5338)
 * [f1c5e4631](https://github.com/kubeovn/kube-ovn/commit/f1c5e4631d5be89a983b96741d77747928e65f10) docs: updated CHANGELOG.md (#5347)
 * [241398a1e](https://github.com/kubeovn/kube-ovn/commit/241398a1e0939defed0d59019da4822c8c7464f8) fix sts/vm lsp in incorrect port groups after rescheduled to another node (#5344)
 * [d4ddd23a9](https://github.com/kubeovn/kube-ovn/commit/d4ddd23a9d37a7cc00ed0b00355307247ad7a016) fix(deps): update module kubevirt.io/api to v1.5.2 (#5339)
 * [e30a970b2](https://github.com/kubeovn/kube-ovn/commit/e30a970b2869d8372bfe94e6c749589264c52b41) chore(deps): update dependency cert-manager to v1.18.0 (#5341)
 * [9ab91ada8](https://github.com/kubeovn/kube-ovn/commit/9ab91ada88dc8a4407cfdb5bbc0658d183445acd) fix(deps): update kubernetes packages to v0.33.1 (#5329)
 * [a5940a234](https://github.com/kubeovn/kube-ovn/commit/a5940a23477fbac8b3ea4ce991cdf2692fc36f66) base: build openssl fips module from deb source (#5331)
 * [4d6274bcd](https://github.com/kubeovn/kube-ovn/commit/4d6274bcdee5c2fe86d1d31b16d71d3de84dcee4) fix: vm has multi nic in the same subnet, but release all (#5336)
 * [cfcce171e](https://github.com/kubeovn/kube-ovn/commit/cfcce171e2a8b6e9b755877ab639d1184aada419) docs: updated CHANGELOG.md (#5335)
 * [3dc7a64ce](https://github.com/kubeovn/kube-ovn/commit/3dc7a64ce136251770b546e8b4d7399743abe5ad) merge dad check to 1.12 (#5334)
 * [bbd62e894](https://github.com/kubeovn/kube-ovn/commit/bbd62e89488ab180068a47e4380c523883f9d4de) fix fips not really replace (#5278)
 * [581425bc9](https://github.com/kubeovn/kube-ovn/commit/581425bc9186339b1e3644f186283eff0344aac5) netpol: fix missing ACL name (#5281)
 * [7624feb07](https://github.com/kubeovn/kube-ovn/commit/7624feb07ec7440dc0541445c44cea153480a493) fix: kubectl-ko: properly get the number of nodes (#5327)
 * [97ff89de2](https://github.com/kubeovn/kube-ovn/commit/97ff89de284a33f9ff941ace243ffb8c25b12787) fix(deps): update module gopkg.in/k8snetworkplumbingwg/multus-cni.v4 to v4.2.1 (#5330)
 * [dcdf899ea](https://github.com/kubeovn/kube-ovn/commit/dcdf899ea6e75672deb172a377ab4f572fd14a7b) chore(deps): update dependency multus to v4.2.1 (#5328)
 * [7ecf5339e](https://github.com/kubeovn/kube-ovn/commit/7ecf5339e5d82453a14ca0b34fef522e90f24a60) chore(deps): update dependency ubuntu to v24 (#5315)
 * [21286537a](https://github.com/kubeovn/kube-ovn/commit/21286537aa287edfcd885d24a7d908c9890fd7ae) chore(deps): update dependency metallb to v0.15.2 (#5313)
 * [a7e439f2c](https://github.com/kubeovn/kube-ovn/commit/a7e439f2ce39a5a2d1d06f28b08e061e58cb2d77) docs: updated CHANGELOG.md (#5325)
 * [eb395857f](https://github.com/kubeovn/kube-ovn/commit/eb395857fc9e04de3b6456abbfa522ef0f6b4ce5) bump go to 1.24.4 (#5321)
 * [c09b19d0c](https://github.com/kubeovn/kube-ovn/commit/c09b19d0cb033d734377e0c21d64fc26bca09a5d) fix(deps): update module go.universe.tf/metallb to v0.15.2 (#5314)
 * [3fb0dcc56](https://github.com/kubeovn/kube-ovn/commit/3fb0dcc569156c73f83a9d168fc6d1f479c07c1d) [BUG] Request Latency 99th Quantile series labels are not unique (#5320)
 * [96f05583b](https://github.com/kubeovn/kube-ovn/commit/96f05583b066e90b6a8cfce1ebf51731ecd57571) chore: add gomodTidy to postUpdateOptions
 * [59619e990](https://github.com/kubeovn/kube-ovn/commit/59619e99099044a73ca8cb5e6d5b8ef643435a19) chore: release branch automate security updates (#5311)
 * [7ad7db246](https://github.com/kubeovn/kube-ovn/commit/7ad7db2468b3ae473209e463ce6c099156870b0d) chore: remove Dependabot configuration file
 * [7e6728ce3](https://github.com/kubeovn/kube-ovn/commit/7e6728ce35dbcd7733be806ece498faa7f8bffe4) tls: rotate self-signed key/certificate (#5303)
 * [fe5161e1c](https://github.com/kubeovn/kube-ovn/commit/fe5161e1c402e20ee1ee34af429b3a918b2c103e) chore(deps): update dependency metallb to v0.15.0 (#5308)
 * [cc4c701d1](https://github.com/kubeovn/kube-ovn/commit/cc4c701d105d934966e54f358f3abe14abba554b) chore(deps): update aquasecurity/trivy-action action to v0.31.0 (#5307)
 * [3a96257a9](https://github.com/kubeovn/kube-ovn/commit/3a96257a936223e8728f80df9a810d38211e94a5) ovn fip spec distributed and support v4:v6 v6:v4 (#5283)
 * [82b73fece](https://github.com/kubeovn/kube-ovn/commit/82b73fece2476e474c18b86365ae97c55968e71f) fix podcidr route still use join ip as src ip (#5287)
 * [04d429faf](https://github.com/kubeovn/kube-ovn/commit/04d429faff27854130058b2c3b443e65b9e455d0) chore(deps): update dependency helm to v3.18.2 (#5305)
 * [bd7a40e7d](https://github.com/kubeovn/kube-ovn/commit/bd7a40e7d18ac81d691efabc9f8ff84a51859903) feat(controller): add database health check (#5294)
 * [f14f44113](https://github.com/kubeovn/kube-ovn/commit/f14f441135c87ed1d320b11ffd95124f8ce7a6f4) fix(deps): update module github.com/puzpuzpuz/xsync/v3 to v4 (#5302)
 * [05e9fbf19](https://github.com/kubeovn/kube-ovn/commit/05e9fbf19e24a314659894cc5ab54a58413c3e44) fix(deps): update module github.com/containerd/containerd to v2 (#5301)
 * [1244a4ff4](https://github.com/kubeovn/kube-ovn/commit/1244a4ff46a8d39d4f1de6194d9bb7f4bdf63dbb) fix(deps): update module github.com/cenkalti/backoff/v4 to v5 (#5300)
 * [ddf54e4e7](https://github.com/kubeovn/kube-ovn/commit/ddf54e4e740af971fc9396863f74917d47790de3) tls: add options to set tls min/max version and cipher suites (#5289)
 * [b6934cf08](https://github.com/kubeovn/kube-ovn/commit/b6934cf088485c63dfe7ef93852611e332a50c16) build(deps): bump github.com/docker/docker (#5295)
 * [baaf5d5f0](https://github.com/kubeovn/kube-ovn/commit/baaf5d5f064bb64c1dde3b9f6bd7471568ac1885) chore(deps): update dependency kind to v0.29.0 (#5297)
 * [e4d9a05d6](https://github.com/kubeovn/kube-ovn/commit/e4d9a05d68990e54a6eca5815fd4d7dab263eb6d) chore(deps): update dependency kwork to v0.7.0 (#5293)
 * [ba1ba78d7](https://github.com/kubeovn/kube-ovn/commit/ba1ba78d70dccf4d433c219bad7281810b6a8936) chore(deps): update dependency helm to v3.18.1 (#5292)
 * [4266f1875](https://github.com/kubeovn/kube-ovn/commit/4266f187539ab7efea63105df213cd5cfdedeefb) chore: Configure Renovate (#5253)
 * [2668ef876](https://github.com/kubeovn/kube-ovn/commit/2668ef876d9ac891b7b2b7f6748c879a657ce535) build(deps): bump github.com/docker/docker (#5290)
 * [9cbafeec8](https://github.com/kubeovn/kube-ovn/commit/9cbafeec86548dfab85fa438ac9ef466f5d10777) fix(bgpspeaker): nat gateway with embedded bgp speaker would crash (#5285)
 * [8583f3ccb](https://github.com/kubeovn/kube-ovn/commit/8583f3ccb189e1643e17ddade112838256e626e4) tls: specify cipher suites (#5288)
 * [6b0cef8d0](https://github.com/kubeovn/kube-ovn/commit/6b0cef8d041650ff5cf904c4514a68d4127353df) build(deps): bump github.com/go-logr/logr from 1.4.2 to 1.4.3 (#5286)
 * [c43c7a77e](https://github.com/kubeovn/kube-ovn/commit/c43c7a77e5c0c5249f04c835bff4511b4954c9c7) fix vlan check conflict and recover it  (#5282)
 * [5db981389](https://github.com/kubeovn/kube-ovn/commit/5db981389b04dbb80ef30ad708cdd92dacbc5b2f) base: add CAP_NET_ADMIN to traceroute (#5275)
 * [c8a6597c5](https://github.com/kubeovn/kube-ovn/commit/c8a6597c5544c3c2ee403bdfdefde28b594759f8) fix port name (#5274)
 * [e75c500ed](https://github.com/kubeovn/kube-ovn/commit/e75c500ed52cc7b61f82ab80d94bf511bb145b6e) any vpc can use any external subnet (#5268)
 * [543fd4583](https://github.com/kubeovn/kube-ovn/commit/543fd458302986d15d9ee589106d214d19767ab5) check vlan conflict (#5233)
 * [58bc11094](https://github.com/kubeovn/kube-ovn/commit/58bc1109447818193add3fcd35d41769f54ec885) docs: updated CHANGELOG.md (#5276)
 * [af90ec12c](https://github.com/kubeovn/kube-ovn/commit/af90ec12c44bbd9137f873bacdb1d14a77e4eb86) build(deps): bump google.golang.org/grpc from 1.72.1 to 1.72.2 (#5273)
 * [03ec76987](https://github.com/kubeovn/kube-ovn/commit/03ec76987351fb9be452a2ee05b2318a86c71c4c) some fixes for vpc egress gateway (#5271)
 * [f88c2df8f](https://github.com/kubeovn/kube-ovn/commit/f88c2df8f9b60a43bc89f6f1e04de3c76dc74707) base: update ovn patches (#5264)
 * [186ea9be4](https://github.com/kubeovn/kube-ovn/commit/186ea9be4ab9dbd27d7c8c4aa6a013df51e1379d) feature: traffic policy (#5263)
 * [4f060e2a5](https://github.com/kubeovn/kube-ovn/commit/4f060e2a5430c8b383a6e93bab3e05041d96b3cb) vpc egreess gateway: select workloads by namespace selector and pod selector (#5260)
 * [224d30cc6](https://github.com/kubeovn/kube-ovn/commit/224d30cc6ba5b787d2592762f89c6c716e854b39) chore(vpc-nat-gw): refactor before adding HA gw (#5197)
 * [5645113be](https://github.com/kubeovn/kube-ovn/commit/5645113beab86de6ecb0bc291f3232bcaba636b4) Feature/x requirements log permission (#5238) (#5258)
 * [eaee89f48](https://github.com/kubeovn/kube-ovn/commit/eaee89f48f918e14373a935b1981d1c67534ab61) fix Node access underlay pod failed when applying  network policy (#5256)
 * [32310fb7f](https://github.com/kubeovn/kube-ovn/commit/32310fb7fcaeb1f5a02ba2697e4bfd05087df49c) bump k8s to v1.32.5 (#5251)
 * [3f0edd93c](https://github.com/kubeovn/kube-ovn/commit/3f0edd93c3c110bb3a74c4906e50079ea6cea06c) cleanup: delete additional ConfigMaps and RoleBindings in cleanup script (#5255)
 * [33be6d513](https://github.com/kubeovn/kube-ovn/commit/33be6d51367a362518b7d819634beca7582904d3) fix: put iptables nat prerouting rule right before kube-proxy inserted one (#5232)
 * [f04a58b22](https://github.com/kubeovn/kube-ovn/commit/f04a58b22779cb1571772cd908ee3b50a59b1e14) fix handle pod elapsed time (#5250)
 * [073a5d7dd](https://github.com/kubeovn/kube-ovn/commit/073a5d7dde1cd5a6c56b1a0c054d735b089a095b) build(deps): bump google.golang.org/grpc from 1.72.0 to 1.72.1 (#5249)
 * [3bf5b1dd8](https://github.com/kubeovn/kube-ovn/commit/3bf5b1dd8c6dc4336b998fbb8953a3b1e8f0ca7d) Fix fisp arm (#5244)
 * [c4b013d87](https://github.com/kubeovn/kube-ovn/commit/c4b013d8786c26b50c630e59e708a8d9dd625e27) Update upgrade-ovs.sh to use POD_NAMESPACE variable for fetching update strategy (#5243)
 * [11b05dc22](https://github.com/kubeovn/kube-ovn/commit/11b05dc2220f95078ffd2fc1f51b2f621f1dfd80) fix slbr missing vip when svc is deleted
 * [8413ae0e9](https://github.com/kubeovn/kube-ovn/commit/8413ae0e91022ad6abcd6f1e723be649dcf19470) Refactor controller to use EndpointSlices instead of Endpoints
 * [3fede2ae0](https://github.com/kubeovn/kube-ovn/commit/3fede2ae074dc32179d91b541ede0d3b29837809) Fix fisp arm (#5241)
 * [aa2a32975](https://github.com/kubeovn/kube-ovn/commit/aa2a32975e04b9217ced84ba1aa629a8d9efb81f) Revert "fix arm (#5239)" (#5240)
 * [8bf693d8e](https://github.com/kubeovn/kube-ovn/commit/8bf693d8e2177ec81102f2284efea1888c7c301d) fix arm (#5239)
 * [a8a91f2d5](https://github.com/kubeovn/kube-ovn/commit/a8a91f2d5ca2aa7c0f4dc42d70c9c998a8f98234) complile openssl 3.5.0 to support fips (#5217)
 * [f2ccacd65](https://github.com/kubeovn/kube-ovn/commit/f2ccacd6571297bd132b24363e14798db6bdcb9a) docs: updated CHANGELOG.md (#5236)
 * [86bd74096](https://github.com/kubeovn/kube-ovn/commit/86bd74096764afbcb085652948abf34e1e491734) vpc egress gateway: set minimum replicas to zero (#5235)
 * [c7365c7d9](https://github.com/kubeovn/kube-ovn/commit/c7365c7d9982052161ca0c586b944c66eb2cb595) build(deps): bump github.com/vishvananda/netlink (#5234)
 * [3ae4be4c2](https://github.com/kubeovn/kube-ovn/commit/3ae4be4c27a1511f7eb68cf2546fc4765d1330e0) log show handle pod elapsed time (#5218)
 * [04c6845d0](https://github.com/kubeovn/kube-ovn/commit/04c6845d07d4fd46d8b13cfb60da7a6f434be696) Deleting a FIP sets the referenced EIP to be ready (#5141)
 * [7fcde6cab](https://github.com/kubeovn/kube-ovn/commit/7fcde6cab069bfd9dd09535e263f2a04e1f52bbe) simple vip (#5229)
 * [143781ea9](https://github.com/kubeovn/kube-ovn/commit/143781ea9a151109cb86aa740c5df75ff2be2bdd) We've seen instances of a VPC being left with a now-deleted subnet on (#5228)
 * [a7a3ef04f](https://github.com/kubeovn/kube-ovn/commit/a7a3ef04f94ebf8791bcefa5b01edca89aed005f) build(deps): bump github.com/osrg/gobgp/v3 from 3.36.0 to 3.37.0 (#5231)
 * [f29c6c8d5](https://github.com/kubeovn/kube-ovn/commit/f29c6c8d5ecaa95235b886482b1e15268636872e) ci: fix cilium chaining with underlay networking (#5226)
 * [fb758c7a3](https://github.com/kubeovn/kube-ovn/commit/fb758c7a3d455c1cff38f32c25dfb8ebb0760eba) build(deps): bump github.com/Microsoft/hcsshim from 0.12.9 to 0.13.0 (#5225)
 * [a86a25682](https://github.com/kubeovn/kube-ovn/commit/a86a256827d8a9c3b0e0f665267445de5c8e93ae) bump go to 1.24.3 (#5214)
 * [01e127ccb](https://github.com/kubeovn/kube-ovn/commit/01e127ccbe0d597ff4fcdc15f90eda64d32f9cbf) build(deps): bump golang.org/x/tools from 0.32.0 to 0.33.0 (#5212)
 * [986cfb645](https://github.com/kubeovn/kube-ovn/commit/986cfb645e91f6b8b182a0fbe408b853a91ba9f2) base: use local patch files (#5201)
 * [509ec599e](https://github.com/kubeovn/kube-ovn/commit/509ec599e9fe3ad7182cbd41bf59ae5d7b2b009c) vpc egress gateway: fix invalid route destination (#5202)
 * [8d75e5c3d](https://github.com/kubeovn/kube-ovn/commit/8d75e5c3d3fad91c2a8f634389f8af7e4aa1850e) build(deps): bump golang.org/x/sys from 0.32.0 to 0.33.0 (#5204)
 * [636aebc1a](https://github.com/kubeovn/kube-ovn/commit/636aebc1acb044b35d016278ca4d7da12612b317) fix (#5200)
 * [e624949db](https://github.com/kubeovn/kube-ovn/commit/e624949db863e90a0d276eea00a1174834e46407) fix nil pointer (#5199)
 * [8aaebc312](https://github.com/kubeovn/kube-ovn/commit/8aaebc31228f4bf8ce98a1433cdead3c1dac8200) build(deps): bump go.uber.org/mock from 0.5.1 to 0.5.2 (#5198)
 * [41cd03166](https://github.com/kubeovn/kube-ovn/commit/41cd031667768fe773a2cd980fc4adda4b592280) fix ip clean (#5193)
 * [8f16fc789](https://github.com/kubeovn/kube-ovn/commit/8f16fc789df5be6488c557a2163edf1a13c5af7a) fix installation script and chart for dpdk (#5194)
 * [056d7944f](https://github.com/kubeovn/kube-ovn/commit/056d7944ff78b63b8106626ac2d2db4f65c60059) some fixes for chart v2 (#5196)
 * [ce2065964](https://github.com/kubeovn/kube-ovn/commit/ce2065964ac2d7a5c2dad96205c9b35905cc3afb) feat(helm): new chart design (#4437)
 * [5ca36a36d](https://github.com/kubeovn/kube-ovn/commit/5ca36a36d49fe44fbd0e7de4826e2e901426c806) feat(helm): new chart design (#4437)
 * [05ac943f4](https://github.com/kubeovn/kube-ovn/commit/05ac943f4d06fd32d46559274c17d8a572bb5fd3) base: bump cni plugins to v1.7.1 (#5191)
 * [2f288ee10](https://github.com/kubeovn/kube-ovn/commit/2f288ee108190f76610fa22acf0183678e8f9e5b) build(deps): bump github.com/containernetworking/plugins (#5190)
 * [f87a51fb7](https://github.com/kubeovn/kube-ovn/commit/f87a51fb7139e3262d8f7636b1192d9142ae9c3b) Makefile: set major/minor version number for the master branch (#5186)
 * [43c6ee802](https://github.com/kubeovn/kube-ovn/commit/43c6ee80237b2aa86edbad8d4596ed5f15dde385) fix metallb e2e (#5183)
 * [d105d35f3](https://github.com/kubeovn/kube-ovn/commit/d105d35f3c17fbe392d3df65ea81d87649d36490) docs: updated CHANGELOG.md (#5181)
 * [b6158dae3](https://github.com/kubeovn/kube-ovn/commit/b6158dae3362e54f07d507e128b83800df2154e3) bump k8s to v1.32.4 (#5178)
 * [42233288a](https://github.com/kubeovn/kube-ovn/commit/42233288a168a8c2fe4999a68f231c7ecaf75a46) Makefile: fix installing dev version on talos (#5179)
 * [821234a02](https://github.com/kubeovn/kube-ovn/commit/821234a02237a1147fa5886c24882f9b420d3071) ci: ignore cni-server restarts caused by join network check failure (#5177)
 * [ab3d11ae7](https://github.com/kubeovn/kube-ovn/commit/ab3d11ae77c08b6fc354c7a8264bccf892d0d808) base: update ovn patch (#5175)
 * [6011b90ba](https://github.com/kubeovn/kube-ovn/commit/6011b90ba467775d7beb78c0517048eb654c9944) fix: clean up garbage lsp (#5172)
 * [05b7d839b](https://github.com/kubeovn/kube-ovn/commit/05b7d839b14cae3af41fca927282b2cb275d2514) build(deps): bump google.golang.org/grpc from 1.71.1 to 1.72.0 (#5171)
 * [8b26247fc](https://github.com/kubeovn/kube-ovn/commit/8b26247fca3c2c2d6d7970f03723773df0910081) docs: updated CHANGELOG.md (#5168)
 * [48db4e499](https://github.com/kubeovn/kube-ovn/commit/48db4e4997a19119372f91930687a7125f960862) controller: ensure ovn route policy is reconciled after node is initialized (#5166)
 * [3177e93bd](https://github.com/kubeovn/kube-ovn/commit/3177e93bd631e18811577a493a7a546bc09aa7f9) modernize: simplify code by using modern constructs (#5163)
 * [9987f5517](https://github.com/kubeovn/kube-ovn/commit/9987f551732db35e030a2cecd305fd09959fddfb) build(deps): bump github.com/docker/docker (#5165)
 * [ad810d6a8](https://github.com/kubeovn/kube-ovn/commit/ad810d6a8a2d93512456ad14876c9845283eb9d7) docs: updated CHANGELOG.md (#5164)
 * [449cc932d](https://github.com/kubeovn/kube-ovn/commit/449cc932d9064f28e4d94c91fc8ccdee51fad5bb) fix vpc nat gateway to correctly use subnet mapped attachment network (#5158)
 * [262b66867](https://github.com/kubeovn/kube-ovn/commit/262b668675104fcfce7cbbfb1b6a3a326f36a610) build(deps): bump github.com/docker/docker (#5160)
 * [91f264e7e](https://github.com/kubeovn/kube-ovn/commit/91f264e7e2f2cb1e7c6c635e7ed78ea2f487f997) support k8s host vm vip type (#5148)
 * [7e584096f](https://github.com/kubeovn/kube-ovn/commit/7e584096f411c6282d3599d762d8a9d220136bbc) pinger: add liveness/readiness probes (#5155)
 * [eb17163bb](https://github.com/kubeovn/kube-ovn/commit/eb17163bb41ae85fc9af215ea3e3808b6a8d5786) add dad e2e (#5146)
 * [80382fd1e](https://github.com/kubeovn/kube-ovn/commit/80382fd1ec161040833ec4514627854453863305) chart: fix ovs ipsec keys host path (#5137)
 * [3fa839c74](https://github.com/kubeovn/kube-ovn/commit/3fa839c74d80862120296fd9ede1586cc878f22f) build(deps): bump kernel.org/pub/linux/libs/security/libcap/cap (#5152)
 * [432426cdd](https://github.com/kubeovn/kube-ovn/commit/432426cdd7c1ba3aa8c2b75dbcabdf7c342d6248) ci: add tests for underlay installation on Talos (#5147)
 * [9204d85fd](https://github.com/kubeovn/kube-ovn/commit/9204d85fdfa3123f27c9cbd4e1d562b17e60f3a1) ci: bump golangci-lint to v2.1.1 (#5151)
 * [745615025](https://github.com/kubeovn/kube-ovn/commit/745615025befeec918d0f47d8009f633ff43ece1) ci: fix kvm/libvirt installation (#5150)
 * [8407e834f](https://github.com/kubeovn/kube-ovn/commit/8407e834fc4cd25a21b0f2151e83320194129645) Add Finalizer to FIP before programming FIP into the VPC NT Gateway (#5142)
 * [6b601709b](https://github.com/kubeovn/kube-ovn/commit/6b601709b50f9444b0fda01a35a354414d5d9625) Support /32 tunnel-source (#5144)
 * [7c8c730d8](https://github.com/kubeovn/kube-ovn/commit/7c8c730d8e24f53c965a80da30379cdde746a07e) ci: add installation test for Talos Linux (#5109)
 * [490ecc037](https://github.com/kubeovn/kube-ovn/commit/490ecc037f43805d32b278650230ec7957d38271) fix dbus/NetworkManager connection in Talos (#5140)
 * [db40fc449](https://github.com/kubeovn/kube-ovn/commit/db40fc44954f3628e5ac0ba6d261c323bcc5aa11) chart: fix local bin directory host path (#5136)
 * [b39caaccb](https://github.com/kubeovn/kube-ovn/commit/b39caaccbe4891f53330b584b1abf0471644f3de) base: update ovn patches (#5139)
 * [fbe9957c3](https://github.com/kubeovn/kube-ovn/commit/fbe9957c3f61ef537e79371a36f69585eca7282b) build(deps): bump github.com/prometheus/client_golang (#5133)
 * [1cc23652f](https://github.com/kubeovn/kube-ovn/commit/1cc23652f5e079b6ab7cf6234c31fcb2b3977b52) build(deps): bump go.uber.org/mock from 0.5.0 to 0.5.1 (#5134)
 * [17bd506d7](https://github.com/kubeovn/kube-ovn/commit/17bd506d7aa1bb829afe71999c7c701a3b8987b6) build(deps): bump github.com/prometheus-community/pro-bing (#5135)
 * [c78348c02](https://github.com/kubeovn/kube-ovn/commit/c78348c029e2c606f4f15159dae56ca225ae853a) build(deps): bump golang.org/x/tools from v0.31.0 to v0.32.0 (#5130)
 * [93483a59c](https://github.com/kubeovn/kube-ovn/commit/93483a59c4df468d3540f5171df35d9e6c7419f8) build(deps): bump github.com/onsi/ginkgo/v2 from 2.23.3 to 2.23.4 (#5128)
 * [f78099773](https://github.com/kubeovn/kube-ovn/commit/f78099773ef829077986fa3e14ab850775a22f3c) build(deps): bump github.com/containernetworking/cni from 1.2.3 to 1.3.0 (#5127)
 * [68c800172](https://github.com/kubeovn/kube-ovn/commit/68c800172ff5628187d09949863d199b8242823c) docs: updated CHANGELOG.md (#5123)
 * [645ac944e](https://github.com/kubeovn/kube-ovn/commit/645ac944e9ca97565a1ac61cbb5a7741d17df8b7) Introduce Swisscom to USERS.md (#5120)
 * [a38b004b3](https://github.com/kubeovn/kube-ovn/commit/a38b004b31ce0a1a3f9a6e224771f4bd5e614b15) gc: consider whether the sts pod is alive during lsp gc (#5122)
 * [ded57da84](https://github.com/kubeovn/kube-ovn/commit/ded57da84a608a476198356c1d3819406b6f5df5) build(deps): bump github.com/osrg/gobgp/v3 from 3.35.0 to 3.36.0 (#5121)
 * [26fa3c019](https://github.com/kubeovn/kube-ovn/commit/26fa3c019b7660ec3de1e6f80b21a0b78685cb0f) build(deps): bump github.com/onsi/gomega from 1.36.3 to 1.37.0 (#5119)
 * [3fe43f2c0](https://github.com/kubeovn/kube-ovn/commit/3fe43f2c0d034fcff19874af0c875179fe09495f) bump go to 1.24.2 (#5117)
 * [ca13d2894](https://github.com/kubeovn/kube-ovn/commit/ca13d2894352091906f11e59aaf7d155c52dadc0) build(deps): bump google.golang.org/grpc from 1.71.0 to 1.71.1 (#5116)
 * [d55a933a2](https://github.com/kubeovn/kube-ovn/commit/d55a933a2bb1e56046b23cbde55ed308b6d86578) ovn lb select the local chassis's backend prefer (#4894)
 * [c3e2307e4](https://github.com/kubeovn/kube-ovn/commit/c3e2307e40f84297c1aed02ed74737bf5b008a68) docs: updated CHANGELOG.md (#5113)
 * [0d227994e](https://github.com/kubeovn/kube-ovn/commit/0d227994e88dd07cec55816206f497fad31c076c) base: update ovs patches (#5111)
 * [c7a6a71ef](https://github.com/kubeovn/kube-ovn/commit/c7a6a71ef7148e1a6836da4f423a793f0396c32a) e2e: fix kubectl-ko trace test (#5108)
 * [6f88e1915](https://github.com/kubeovn/kube-ovn/commit/6f88e1915b1e5ae1baefce567850cf3ecfb1f735) feat(controller): skip appending VM LSPs if default Multus network is present (#5106)
 * [9fc5197fe](https://github.com/kubeovn/kube-ovn/commit/9fc5197fe9a011abea74c4447a4bd46dca2c5ece) build(deps): bump github.com/docker/docker (#5104)
 * [695db1531](https://github.com/kubeovn/kube-ovn/commit/695db15319e09cb9e4582b79da904e507bea0dc1) build(deps): bump github.com/onsi/gomega from 1.36.2 to 1.36.3 (#5103)
 * [0a6e5e38a](https://github.com/kubeovn/kube-ovn/commit/0a6e5e38a616348e66e27f1e30811bb91d5008d4) build(deps): bump github.com/rs/zerolog from 1.33.0 to 1.34.0 (#5101)
 * [25e392585](https://github.com/kubeovn/kube-ovn/commit/25e39258589b923503cf86c2e332a19bc9145f2d) build(deps): bump github.com/onsi/ginkgo/v2 from 2.23.2 to 2.23.3 (#5102)
 * [30451fe85](https://github.com/kubeovn/kube-ovn/commit/30451fe858cd617ad6a95144c4b4b78df8719fed) build(deps): bump github.com/docker/docker (#5100)
 * [51067fd2a](https://github.com/kubeovn/kube-ovn/commit/51067fd2a5390e6582e7ba1ef8aa90264c5c7d07) bump golangci-lint to v2.0.1 (#5097)
 * [af6c3f6ea](https://github.com/kubeovn/kube-ovn/commit/af6c3f6eaaafc31271c15a7d0988e8154ef7a443) build(deps): bump gopkg.in/k8snetworkplumbingwg/multus-cni.v4 (#5096)
 * [339df023f](https://github.com/kubeovn/kube-ovn/commit/339df023fbe9ae04c3d0b0fa13d49ef604874f8e) build(deps): bump github.com/onsi/ginkgo/v2 from 2.23.1 to 2.23.2 (#5094)
 * [5e624cf18](https://github.com/kubeovn/kube-ovn/commit/5e624cf18b6af088e3da8d99145958019a1d4076) build(deps): bump github.com/docker/docker (#5092)
 * [7b829a4ce](https://github.com/kubeovn/kube-ovn/commit/7b829a4ce4f909c3ef5ce76f7825bee35f0abfd9) build(deps): bump github.com/onsi/ginkgo/v2 from 2.23.0 to 2.23.1 (#5091)
 * [3310f23dd](https://github.com/kubeovn/kube-ovn/commit/3310f23dde1374e3941512ab7916f64163fcedc2) fix: egress network policy not work, when no pod hit matchlabel (#5089)
 * [c64143504](https://github.com/kubeovn/kube-ovn/commit/c64143504e3d944946224a12f308b9b53d03fecc) docs: updated CHANGELOG.md (#5086)
 * [7560a577d](https://github.com/kubeovn/kube-ovn/commit/7560a577db9bd3b7e0a968fb021a45d4c8824fe1) build(deps): bump github.com/containerd/containerd from 1.7.26 to 1.7.27 (#5085)
 * [fc14f7959](https://github.com/kubeovn/kube-ovn/commit/fc14f79598cd6987c3c1a1533ff8e8d5597536d2) build(deps): bump aquasecurity/trivy-action from 0.29.0 to 0.30.0 (#5081)
 * [ba9bcece1](https://github.com/kubeovn/kube-ovn/commit/ba9bcece1c06af831b36fbd9879e4dfe7a119dc0) bind to pod ips when env variable ENABLE_BIND_LOCAL_IP is set to true (#5049)
 * [6f4100d7f](https://github.com/kubeovn/kube-ovn/commit/6f4100d7f79cfa09b3a59f6e8f224a082b6ba75b) cni-server: check ipv6 dadfailed flag (#5042)
 * [b1c5f2565](https://github.com/kubeovn/kube-ovn/commit/b1c5f2565206d6b974f41020919d4884657d6868) [fix] When the Nat-gw pod container restarts unexpectedly, trigger nat-gw statefulset restart to restore the nat-gw pod configuration (#5072)
 * [3f261807e](https://github.com/kubeovn/kube-ovn/commit/3f261807edcbb6652975b783c305663d79429315) bump kubevirt to v1.5.0 (#5080)
 * [d650cfb6f](https://github.com/kubeovn/kube-ovn/commit/d650cfb6f1a5e98136e8890524e0de1f33471663) docs: updated CHANGELOG.md (#5079)
 * [f98582420](https://github.com/kubeovn/kube-ovn/commit/f98582420289fbde1041ec2e81338eb5c330e93d) bump go to 1.24.1 (#5076)
 * [56b54c7c2](https://github.com/kubeovn/kube-ovn/commit/56b54c7c253859bd7af5aec767200e439bc900a5) build(deps): bump golang.org/x/net from 0.35.0 to 0.36.0 in /test/anp (#5075)
 * [0f4fde3ba](https://github.com/kubeovn/kube-ovn/commit/0f4fde3baa19e766ce87690fa8367355cdb249bc) build(deps): bump k8s.io/kubernetes in the k8s-io group (#5074)
 * [7b11e93f9](https://github.com/kubeovn/kube-ovn/commit/7b11e93f9df2d185aebc078b40d6da328ef82e4f) feat: Improve subnet and VPC finalizer handling (#5071)
 * [58e83ce0f](https://github.com/kubeovn/kube-ovn/commit/58e83ce0f327b40a0f252ec3e07d4a761a23b020) feat: Make Kube-OVN namespace configurable with default value (#5069)
 * [d8a4a0590](https://github.com/kubeovn/kube-ovn/commit/d8a4a0590d68255de6d6fcdae134576bd1ac73a6) build(deps): bump sigs.k8s.io/controller-runtime from 0.20.2 to 0.20.3 (#5067)
 * [c26810950](https://github.com/kubeovn/kube-ovn/commit/c268109500f4a68bdf1ea64607f38efdce7830f9) Fix #5028: Orphaned subnets which reference a non-existent VPC cause new namespaces to never get correct annotations. (#5031)
 * [902d37c6f](https://github.com/kubeovn/kube-ovn/commit/902d37c6f34f012a2e131fa3a575b790d0198983) remove genev_sys_6081 when uninstall kube-ovn (#5066)
 * [b303fa119](https://github.com/kubeovn/kube-ovn/commit/b303fa119bd3f57783c97db78a18ca92498ddb69) add copy of the kubevirt informer from upstream (#5065)
 * [c5f5c5e0c](https://github.com/kubeovn/kube-ovn/commit/c5f5c5e0c7f7fbdec4ecdfd9715ce341c250d185) build(deps): bump golang.org/x/tools from 0.30.0 to 0.31.0 (#5059)
 * [bfef0c8a1](https://github.com/kubeovn/kube-ovn/commit/bfef0c8a18cf27dfd1caf03a475f3fa7b49861ba) build(deps): bump github.com/onsi/ginkgo/v2 from 2.22.2 to 2.23.0 (#5060)
 * [dad4cdfba](https://github.com/kubeovn/kube-ovn/commit/dad4cdfbae32d1b33fc4f130b6ee6ad5987b9acb) build(deps): bump golang.org/x/time from 0.10.0 to 0.11.0 (#5054)
 * [bb66e12f1](https://github.com/kubeovn/kube-ovn/commit/bb66e12f1cbfc9befb2de5d63d89871efe6095a2) build(deps): bump golang.org/x/sys from 0.30.0 to 0.31.0 (#5055)
 * [fb2027b28](https://github.com/kubeovn/kube-ovn/commit/fb2027b28baf8034d7a6257d4b53582354deb441) build(deps): bump google.golang.org/grpc from 1.70.0 to 1.71.0 (#5053)
 * [08c48a353](https://github.com/kubeovn/kube-ovn/commit/08c48a353628d3d72f69eea690a4c005a1e6c31f) build(deps): bump golang.org/x/mod from 0.23.0 to 0.24.0 (#5056)
 * [44a8f2045](https://github.com/kubeovn/kube-ovn/commit/44a8f20450c7b845d011b665a6d7afd4120a3f75) build(deps): bump kernel.org/pub/linux/libs/security/libcap/cap (#5052)
 * [7212f8574](https://github.com/kubeovn/kube-ovn/commit/7212f857427b8257b74e7a72073b12771e239bbb) build(deps): bump github.com/prometheus/client_golang (#5051)
 * [6dfea7466](https://github.com/kubeovn/kube-ovn/commit/6dfea74666c7ad81de39dedb1dc2bd2e3192324b) simple vip lable update and then update subnet status (#5036)
 * [fae0f6388](https://github.com/kubeovn/kube-ovn/commit/fae0f6388582ea7c451d5a2fa494ae1ebaee96c1) build(deps): bump github.com/osrg/gobgp/v3 from 3.34.0 to 3.35.0 (#5047)
 * [599c0c3d3](https://github.com/kubeovn/kube-ovn/commit/599c0c3d3319d3e51a13e96639695b863c512d5e) kubectl-ko: fix conntrack state (#5038)
 * [103ddc6f3](https://github.com/kubeovn/kube-ovn/commit/103ddc6f39bef62682eb7d1fbe1746c1eef58e91) build(deps): bump github.com/containerd/containerd from 1.7.25 to 1.7.26 (#5041)
 * [340c48a1c](https://github.com/kubeovn/kube-ovn/commit/340c48a1c8d2c52732fd795ef5b7e87113857239) build(deps): bump github.com/docker/docker (#5033)
 * [4a990cca2](https://github.com/kubeovn/kube-ovn/commit/4a990cca20ea460d26a93ed2958687d4592c4604) build(deps): bump github.com/emicklei/go-restful/v3 (#5032)
 * [aa4cdbea8](https://github.com/kubeovn/kube-ovn/commit/aa4cdbea8dfef15b2f879f21c3fe191bac0fd1a2) fix rerun the completed vim migration (#5020) (#5026)
 * [94c41b1da](https://github.com/kubeovn/kube-ovn/commit/94c41b1da256fcc1832798d169ee888eca3c7497) bfd: fix bfdd-control setting mintx/minrx (#5021)
 * [4005606ea](https://github.com/kubeovn/kube-ovn/commit/4005606ea768e600a1b59b40628aed6fc48613b8) feat(GC): Add check for GC disabled (#5005)
 * [87feaf5e2](https://github.com/kubeovn/kube-ovn/commit/87feaf5e252595c2e95b32826fd914746ae0d03b) build(deps): bump github.com/containernetworking/plugins from v1.6.0 to v1.6.2 (#5017)
 * [3c0d84902](https://github.com/kubeovn/kube-ovn/commit/3c0d849024edd02520f12404191b4422f0959772) fix logging in ovn lr/ls gc (#5014)
 * [91b969818](https://github.com/kubeovn/kube-ovn/commit/91b96981817d5d8d14c99659a0e1dfdc29766761) build(deps): bump github.com/prometheus/client_golang (#5016)
 * [a46d69506](https://github.com/kubeovn/kube-ovn/commit/a46d69506ea307f533d24bc70c0fca5fa7b5ef28) docs: updated CHANGELOG.md (#5012)
 * [1cf92b4c2](https://github.com/kubeovn/kube-ovn/commit/1cf92b4c2766fcb69257c1d4ee7c21c69c33dffc) build(deps): bump azure/setup-helm from 4.2.0 to 4.3.0 (#5011)
 * [7213dc4e3](https://github.com/kubeovn/kube-ovn/commit/7213dc4e3fc61f37befd18cc8ad0d0934848cd62) build(deps): bump github.com/puzpuzpuz/xsync/v3 from 3.5.0 to 3.5.1 (#5001)
 * [c10c6e99a](https://github.com/kubeovn/kube-ovn/commit/c10c6e99a3e1886f6762507c8c6f3408ee51d460) docs: updated CHANGELOG.md (#5000)
 * [a3e4ce730](https://github.com/kubeovn/kube-ovn/commit/a3e4ce730a92a54d3e84957dddfc71074cd51ad7) docs: updated CHANGELOG.md (#4999)
 * [0b45b667a](https://github.com/kubeovn/kube-ovn/commit/0b45b667a2dfeec159d1bb6a19ae60398e1c1666) bump k8s to v1.32.2 (#4992)
 * [f9904d26a](https://github.com/kubeovn/kube-ovn/commit/f9904d26a3b2c7304d70f573353d90acd2dc58e8) controller: respect named ports of restartable init containers (#4994)
 * [d4eb19d16](https://github.com/kubeovn/kube-ovn/commit/d4eb19d16220bd91013cf0b89c2acb931b783cff) build(deps): bump sigs.k8s.io/controller-runtime from 0.20.1 to 0.20.2 (#4993)
 * [48eec14ea](https://github.com/kubeovn/kube-ovn/commit/48eec14ea9c1ccfd8b83398c5093d99f868c2302) build(deps): bump k8s.io/kubernetes in the k8s-io group (#4990)
 * [b44ea855f](https://github.com/kubeovn/kube-ovn/commit/b44ea855f4f7d765da3ff0b117b5311e748da7c2) bump go to 1.24.0 (#4986)
 * [7a36465e3](https://github.com/kubeovn/kube-ovn/commit/7a36465e33a1da8d5d695643d6d5a27bcddf7693) iptables: fix subnet metrics rules (#4979)
 * [d7714f991](https://github.com/kubeovn/kube-ovn/commit/d7714f991d5dbcc5ff940ea60bac6943b4902910) fix superfluous response.WriteHeader (#4980)
 * [31601001c](https://github.com/kubeovn/kube-ovn/commit/31601001c985ef18af770cc54f6057d33d69eae4) e2e: skip some test cases for versions prior to v1.14 in cilium chaining (#4988)
 * [f909a3a7d](https://github.com/kubeovn/kube-ovn/commit/f909a3a7d15ff8789fe6b279653da1d3cfefebfd) remove twitter account
 * [7d7d2ddfd](https://github.com/kubeovn/kube-ovn/commit/7d7d2ddfd3fa463214f8dbd0751dcedd6fe9de81) fix security groups changed when vm is shut down (#4976)
 * [501231a03](https://github.com/kubeovn/kube-ovn/commit/501231a03b2995084c1cf0a3165cd4e64c33bd7f) bump golang.org/x/tools from v0.29.0 to v0.30.0 (#4983)
 * [89a220693](https://github.com/kubeovn/kube-ovn/commit/89a2206935eef5b27756e89a468a329c8e92eeab) ci: bump cilium to v1.17.0 (#4982)
 * [0407520b8](https://github.com/kubeovn/kube-ovn/commit/0407520b814b001a6c0eb239904944f11d148170) e2e: fix ip conflict (#4977)
 * [8beaeab1d](https://github.com/kubeovn/kube-ovn/commit/8beaeab1d6ac91b21711d2fbb7109d4662ab8eef) Revert "refactor: remove duplicated iptables subnet forward rules (#4860)" (#4978)
 * [242b04b72](https://github.com/kubeovn/kube-ovn/commit/242b04b727e60b39cf0690e1290a77cfd8c5d2e9) refactor: remove duplicated iptables subnet forward rules (#4860)
 * [57b373373](https://github.com/kubeovn/kube-ovn/commit/57b373373d3f0dd734cb6646b990c6424141afd9) use httpGet as liveness/readiness probe method (#4945)
 * [7f3031366](https://github.com/kubeovn/kube-ovn/commit/7f30313666ce03edcbc5d0e95b41f2c992db1bce) build(deps): bump google.golang.org/protobuf from 1.36.4 to 1.36.5 (#4970)
 * [30e97c7d6](https://github.com/kubeovn/kube-ovn/commit/30e97c7d619f29bf4ca40a236e202f4882bf600f) bump github.com/hashicorp/yamux from v0.1.1 to v0.1.2 (#4968)
 * [2a1239ffa](https://github.com/kubeovn/kube-ovn/commit/2a1239ffacf3abb914b2fefd3d4af73573e44238) controller: consider StatefulSet's start ordinal (#4967)
 * [61ae78798](https://github.com/kubeovn/kube-ovn/commit/61ae787989cee408c4ec48f95aa4603e25a92a89) bump go to 1.23.6 (#4962)
 * [48ea5af47](https://github.com/kubeovn/kube-ovn/commit/48ea5af479fceb8b52ad8f41ae3b68b5c14b2094) build(deps): bump golang.org/x/time from 0.9.0 to 0.10.0 (#4958)
 * [f5292fd21](https://github.com/kubeovn/kube-ovn/commit/f5292fd21fc512b6dc15d255ab2521bc78e6c2c5) build(deps): bump golang.org/x/sys from 0.29.0 to 0.30.0 (#4959)
 * [497c2e74b](https://github.com/kubeovn/kube-ovn/commit/497c2e74bf19ff51e25fa548d417e2b4b997f701) build(deps): bump github.com/osrg/gobgp/v3 from 3.33.0 to 3.34.0 (#4956)
 * [209b461ea](https://github.com/kubeovn/kube-ovn/commit/209b461eacf6e570cb780a5e90c6a985f76c7786) build(deps): bump github.com/prometheus-community/pro-bing (#4957)
 * [e65a7def2](https://github.com/kubeovn/kube-ovn/commit/e65a7def2a3117d0276ffbd3c2655815a5aff3a0) build(deps): bump github.com/brianvoe/gofakeit/v7 from 7.1.2 to 7.2.1 (#4954)
 * [2074e695e](https://github.com/kubeovn/kube-ovn/commit/2074e695e74ffa6a905e1925aacc2b48d0c5ea5a) build(deps): bump github.com/spf13/pflag from 1.0.5 to 1.0.6 (#4953)
 * [b92359055](https://github.com/kubeovn/kube-ovn/commit/b9235905561bef4d229baf8e438379b1bd552607) build(deps): bump github.com/evanphx/json-patch/v5 from 5.9.10 to 5.9.11 (#4951)
 * [e5a3c4d6c](https://github.com/kubeovn/kube-ovn/commit/e5a3c4d6c277fff2d3628cd7685043d1cdd60090) build(deps): bump github.com/evanphx/json-patch/v5 from 5.9.0 to 5.9.10 (#4949)
 * [b7786ad18](https://github.com/kubeovn/kube-ovn/commit/b7786ad187b756ec479000ac66038ed40ae45c0a) build(deps): bump github.com/golang/glog from 1.2.3 to 1.2.4 (#4950)
 * [a60e21b3d](https://github.com/kubeovn/kube-ovn/commit/a60e21b3d43722c532847167c250102627e6ca9d) build(deps): bump github.com/puzpuzpuz/xsync/v3 from 3.4.1 to 3.5.0 (#4948)
 * [de7ea70eb](https://github.com/kubeovn/kube-ovn/commit/de7ea70ebb5b9125cf7475909adf397173b6a813) build(deps): bump github.com/docker/docker (#4944)
 * [34ece0464](https://github.com/kubeovn/kube-ovn/commit/34ece0464deea380c48b5b6daa0b19b52cf49011) build(deps): bump google.golang.org/protobuf from 1.36.3 to 1.36.4 (#4947)
 * [6c059cef6](https://github.com/kubeovn/kube-ovn/commit/6c059cef6a110a47ee6076841e049de06229cbcc) build(deps): bump google.golang.org/grpc from 1.69.4 to 1.70.0 (#4946)
 * [29e989c94](https://github.com/kubeovn/kube-ovn/commit/29e989c9451b1859a2f51d1b691fd269bf50871e) build(deps): bump sigs.k8s.io/controller-runtime from 0.20.0 to 0.20.1 (#4943)
 * [903c532ff](https://github.com/kubeovn/kube-ovn/commit/903c532ff96f9a9b8055857c6b60cc6f4a9ec5d6) docs(cloudnull): Add Rackspace as a user (#4941)
 * [ead6e6b4e](https://github.com/kubeovn/kube-ovn/commit/ead6e6b4e61ed522cf89be67394b3d73bf504592) build(deps): bump github.com/prometheus-community/pro-bing (#4940)
 * [fe30a54cc](https://github.com/kubeovn/kube-ovn/commit/fe30a54cc2b7221fd3a0e56d0cf1fe3a37774ae0) build(deps): bump helm/chart-releaser-action from 1.6.0 to 1.7.0 (#4938)
 * [2f2b722eb](https://github.com/kubeovn/kube-ovn/commit/2f2b722ebae4a08241852798b91415ee6f84c426) build(deps): bump github.com/puzpuzpuz/xsync/v3 from 3.4.0 to 3.4.1 (#4937)
 * [6b9c9008a](https://github.com/kubeovn/kube-ovn/commit/6b9c9008af3cfdc87a51516610be53aabb727c45) ci: build arm64 images on arm64 hosted runners (#4936)
 * [aeecbd4e3](https://github.com/kubeovn/kube-ovn/commit/aeecbd4e3ac8589bb8fbfaae609e50672178a93a) bump go to 1.23.5 (#4935)
 * [e84536bdd](https://github.com/kubeovn/kube-ovn/commit/e84536bdd8616b97b511913c81980b9c72f004e2) controller: fix gateway nodes check (#4912)
 * [5c9058aec](https://github.com/kubeovn/kube-ovn/commit/5c9058aec4a3c2de6a6e7088fc838584ec6429f5) build(deps): bump sigs.k8s.io/controller-runtime from 0.19.4 to 0.20.0 (#4932)
 * [c0113468f](https://github.com/kubeovn/kube-ovn/commit/c0113468fa15f664a2d6583badaad14ed9debd0f) bump k8s to v1.32.1 (#4930)
 * [773b790dd](https://github.com/kubeovn/kube-ovn/commit/773b790dd6fe1da4c9f8caf995f36780f138df29) fix log (#4928)
 * [ccff9b09c](https://github.com/kubeovn/kube-ovn/commit/ccff9b09c4ee25897b90f8f7e5b7ae30447b8443) make sure gw pod exist before eip creation (#4924)
 * [fba416ade](https://github.com/kubeovn/kube-ovn/commit/fba416adea8edf669305e4330ad63ae756df390f) build(deps): bump google.golang.org/protobuf from 1.36.2 to 1.36.3 (#4927)
 * [af6a1b377](https://github.com/kubeovn/kube-ovn/commit/af6a1b37783905010990df38c95be4364a11b654) [performance] in large-scale clusters, init node router policy too slow (#4895)
 * [756cfd94b](https://github.com/kubeovn/kube-ovn/commit/756cfd94bc12c998faa70e2ec408e719d49adba2) build(deps): bump github.com/docker/docker (#4926)
 * [27f0754c5](https://github.com/kubeovn/kube-ovn/commit/27f0754c57c0612106f06352b3d356519d3449f0) remove dup func and fix ut (#4921)
 * [4dd28f558](https://github.com/kubeovn/kube-ovn/commit/4dd28f5588cca5cafaef31f0611f24f57a328839) build(deps): bump google.golang.org/grpc from 1.69.2 to 1.69.4 (#4923)
 * [d2608530d](https://github.com/kubeovn/kube-ovn/commit/d2608530d2778141c625d5d5bf4e7f3aa0fdfa93) controller: check condition NodeNetworkUnavailable when determining whether node is ready (#4917)
 * [f1d91c9bb](https://github.com/kubeovn/kube-ovn/commit/f1d91c9bb18ad46f8fd3e4e23c91f2dee54b895f) replace reflect.DeepEqual with slices.Equal and maps.Equal (#4918)
 * [e250afa56](https://github.com/kubeovn/kube-ovn/commit/e250afa56ff5d22019a8d7fa07ae7aa026ee3d9b) cni-server: set node NetworkUnavailable condition after join subnet gateway check (#4915)
 * [a7fff996b](https://github.com/kubeovn/kube-ovn/commit/a7fff996b5d0c66b6ea4428d0251844aad4a7d08) build(deps): bump sigs.k8s.io/controller-runtime from 0.19.3 to 0.19.4 (#4911)
 * [8bb726fb3](https://github.com/kubeovn/kube-ovn/commit/8bb726fb3b50c91d227f2d48c0ce98b345af5573) e2e: do not check mac annotation for versions prior to v1.13 (#4910)
 * [d5f620510](https://github.com/kubeovn/kube-ovn/commit/d5f6205107e378c0e487ec69c50cd2d0f461b6e0) ipam: check subnet's available ipv6 address count (#4903)
 * [eb3930813](https://github.com/kubeovn/kube-ovn/commit/eb3930813d4eb179b3fd780407bcc3f0aabf0c1c) base: bump cni plugins to v1.6.2 (#4904)
 * [672af54db](https://github.com/kubeovn/kube-ovn/commit/672af54dbe8a50aa0189b9b3a33c3203651b6ea6) build(deps): bump golang.org/x/tools from 0.28.0 to 0.29.0 (#4905)
 * [bf81d1523](https://github.com/kubeovn/kube-ovn/commit/bf81d1523d746c7851ce53846161764906de5a7b) build(deps): bump google.golang.org/protobuf from 1.36.1 to 1.36.2 (#4907)
 * [604413604](https://github.com/kubeovn/kube-ovn/commit/6044136041af2c8a527faa0e96e7f2e9ca721242) Increase code specification and readability in vxlan nic name and qos (#4896)
 * [d66af7ca5](https://github.com/kubeovn/kube-ovn/commit/d66af7ca5499d7c23affcf5b7007c0268ab3e560) build(deps): bump golang.org/x/time from 0.8.0 to 0.9.0 (#4901)
 * [cbb62ef56](https://github.com/kubeovn/kube-ovn/commit/cbb62ef5613846d5f6ac4d75af7c0aacb48db455) build(deps): bump golang.org/x/sys from 0.28.0 to 0.29.0 (#4900)
 * [44fcf08c9](https://github.com/kubeovn/kube-ovn/commit/44fcf08c90b2900d5701f035381572dd5c621f96) fix: kube-ovn-controller cannot be ready when ENABLE_METRICS is false (#4886)
 * [4d3dd9c7d](https://github.com/kubeovn/kube-ovn/commit/4d3dd9c7db1a73221a7da828a03ba449d08aecf4) controller: generate stable annotations for pod routes (#4889)
 * [ed8bea32f](https://github.com/kubeovn/kube-ovn/commit/ed8bea32fd8ae3846dbcb9c6ded86fcb10658394) e2e: fix duplicate random cidr (#4893)
 * [ac7cf401e](https://github.com/kubeovn/kube-ovn/commit/ac7cf401e2f414bacaf14acdf2b91f261a20bb82) build(deps): bump github.com/osrg/gobgp/v3 from 3.32.0 to 3.33.0 (#4891)
 * [3c6eeaf57](https://github.com/kubeovn/kube-ovn/commit/3c6eeaf57fd927e915161ff63603b501abcc637f) build(deps): bump github.com/onsi/ginkgo/v2 from 2.22.1 to 2.22.2 (#4890)
 * [48f309cd8](https://github.com/kubeovn/kube-ovn/commit/48f309cd8b834a185f539c8527751dc84ed33e54) update readme
 * [010b701f8](https://github.com/kubeovn/kube-ovn/commit/010b701f8e616ab5b778151af418349c1651f538) Support multiple IPPools in the namespace (#4777)
 * [5eb316b03](https://github.com/kubeovn/kube-ovn/commit/5eb316b036be3af8e809635d91f83c49abc7dc43) rename Makefile.e2e to e2e.mk (#4885)
 * [d10251fd7](https://github.com/kubeovn/kube-ovn/commit/d10251fd741033c7983fab2e3e5ac60a11be62b2) ipam: use ip provided by nad annotation when providing IPAM for other CNI plugins (#4883)
 * [bad1bdb3b](https://github.com/kubeovn/kube-ovn/commit/bad1bdb3b34075f0d9603915bdbd94f5906b4fc2) bump go modules used by ANP e2e tests (#4882)
 * [ee6f907be](https://github.com/kubeovn/kube-ovn/commit/ee6f907beb1914f1807840261b24b8d4290e83bf) add nodeSelector for lb-svc (#4793)
 * [38df57725](https://github.com/kubeovn/kube-ovn/commit/38df577256e3d6dab80a86b9179bcc39ff11f5ff) fix golangci-lint (#4880)
 * [c2bc18227](https://github.com/kubeovn/kube-ovn/commit/c2bc18227b1f8b3b0af439c52fbc33a3720602d6) remove unused function (#4877)
 * [9b0072304](https://github.com/kubeovn/kube-ovn/commit/9b00723041cef25d14ed0d1cc5b59cb01798b849) should be able to use mac and ip provided by k8s.v1.cni.cncf.io/networks annotation fix e2e version (#4878)
 * [93759b921](https://github.com/kubeovn/kube-ovn/commit/93759b921408fd52ddcf66e5d17a2b7b4399b3fa) pod should use mac and ips provider by multus firstly (#4800)
 * [79acd891d](https://github.com/kubeovn/kube-ovn/commit/79acd891dd97f8bd504d7bde2fba84c22eccec0b) build(deps): bump github.com/onsi/gomega from 1.36.1 to 1.36.2 (#4871)
 * [e04c2480e](https://github.com/kubeovn/kube-ovn/commit/e04c2480edc24af358c2a969f6b04eed387f72c8) skip node local dns ip conntrack when set acl (#4825)
 * [9aa076bb4](https://github.com/kubeovn/kube-ovn/commit/9aa076bb403108f3da4c54257012d8c29fdeb3b2) Add default subnet in custom vpc to the beginning of the list (#4826)
 * [cec47b529](https://github.com/kubeovn/kube-ovn/commit/cec47b5290207f8f943d68fd11a43f8e0ed515c8) keep dockerfile variable the same as download-go-deps.sh (#4863)
 * [bf5240087](https://github.com/kubeovn/kube-ovn/commit/bf5240087f37e697b427830f424ab38e38d8ed70) fix(controller/subnet): controller crashes on subnets if gateway is unspecified and netpol are disabled (#4848)
 * [8ecbae8b4](https://github.com/kubeovn/kube-ovn/commit/8ecbae8b4d36e33cd1ff3cf646ad6ab70cf7ac0b) clean up legacy iptables rules only when iptables/ip6_tables is loaded (#4855)
 * [b110300fa](https://github.com/kubeovn/kube-ovn/commit/b110300fa9c1769fb4f998b9f13efb9b4128a7bd) build(deps): bump helm/kind-action from 1.11.0 to 1.12.0 (#4868)
 * [454817da3](https://github.com/kubeovn/kube-ovn/commit/454817da3b78038c69310aa48c3c57316b67540f) build(deps): bump google.golang.org/protobuf from 1.36.0 to 1.36.1 (#4869)
 * [3245c91cf](https://github.com/kubeovn/kube-ovn/commit/3245c91cf34b0a7eed3484c74339852a214977cf) build(deps): bump github.com/onsi/ginkgo/v2 from 2.22.0 to 2.22.1 (#4861)
 * [6478007b7](https://github.com/kubeovn/kube-ovn/commit/6478007b742c76f4c831a4f933391d1dd0c871b2) fix: err is always nil (#4857)
 * [3e78b22da](https://github.com/kubeovn/kube-ovn/commit/3e78b22da960aebac76990ca2b5b229b690d508b) deps: bump golang.org/x/net to v0.33.0 (#4851)
 * [87016253a](https://github.com/kubeovn/kube-ovn/commit/87016253abcca997a3787b3156f0479ff5104566) build(deps): bump github.com/docker/docker (#4850)
 * [1caca2c56](https://github.com/kubeovn/kube-ovn/commit/1caca2c56df18013682adfd2f89497c569cb8532) fix gateway node check for centralized ecmp subnets (#4847)
 * [0a8adbf4c](https://github.com/kubeovn/kube-ovn/commit/0a8adbf4c9cdf517d9d45fa96a95f15b6069a8d9) remove unused function (#4821)
 * [96ad24588](https://github.com/kubeovn/kube-ovn/commit/96ad2458880a39993677e8166e2e7786dc0149bc) use cache.MetaObjectToName() to get namespaced name (#4842)
 * [a88306ac0](https://github.com/kubeovn/kube-ovn/commit/a88306ac099ccc5d69ae4419acdad757ac87fff5) use JSON merge patch to update labels/annotations (#4838)
 * [64acd012c](https://github.com/kubeovn/kube-ovn/commit/64acd012c1f0cc80a735daab1c4d85b488fe2bf9) fix getting subnet cidr by protocol (#4844)
 * [4875f23c6](https://github.com/kubeovn/kube-ovn/commit/4875f23c6be6e17b1650fee0ed839d2e37567b99) ci: wait for kubevirt crd to be created before creating CR (#4839)
 * [bcf771332](https://github.com/kubeovn/kube-ovn/commit/bcf771332a6dd1a4b712a359bf59414d52c71acd) build(deps): bump google.golang.org/protobuf from 1.35.2 to 1.36.0 (#4846)
 * [517f0e488](https://github.com/kubeovn/kube-ovn/commit/517f0e488e26c0f8a44539cd3ebc784f5374fe4d) e2e: add test case for specifying mac address when nad plugin is macvlan (#4836)
 * [9291e1f17](https://github.com/kubeovn/kube-ovn/commit/9291e1f17d17c3725f109e7c8c7b9a07295972f8) build(deps): bump helm/kind-action from 1.10.0 to 1.11.0 (#4837)
 * [6c109a9ef](https://github.com/kubeovn/kube-ovn/commit/6c109a9efe565a2f5116e58c2e3fd8f4816d614a) refactor: remove redundant policy route addition in node handling (#4835)
 * [77650e593](https://github.com/kubeovn/kube-ovn/commit/77650e5931c708dc39bba1792cbe93d9eb368613) fix  some normative issues 24.12.16 (#4833)
 * [ab4824e3a](https://github.com/kubeovn/kube-ovn/commit/ab4824e3a796747635bbc1135fccee0b9598f29e) docs: updated CHANGELOG.md (#4832)
 * [b2148a224](https://github.com/kubeovn/kube-ovn/commit/b2148a22492d8db7a360896e6504c89ad15b6049) docs: updated CHANGELOG.md (#4831)
 * [d8064d85e](https://github.com/kubeovn/kube-ovn/commit/d8064d85e9c285b766cf6d13a6ad07ef31901e44) build(deps): bump google.golang.org/grpc from 1.68.1 to 1.69.0 (#4830)
 * [bd3593da6](https://github.com/kubeovn/kube-ovn/commit/bd3593da679afee8e05e9fc99b5c6407a186bb2d) api: add scale subresource for vpc-egress-gateway (#4829)
 * [514a2b643](https://github.com/kubeovn/kube-ovn/commit/514a2b643e26b29e4102a1fe1b3a38938be082d1) build(deps): bump k8s.io/kubernetes from 1.31.4 to 1.32.0 (#4827)
 * [e2e6dd77c](https://github.com/kubeovn/kube-ovn/commit/e2e6dd77c8e7b8f858e4c77351298b03f627644e) cni: do not exit if the sysctl variable does not exist or can not be set (#4828)
 * [386290793](https://github.com/kubeovn/kube-ovn/commit/38629079338cd643470e3a059778b5cb88bebd36) add anp/banp unittests (#4774)
 * [7006e9cb8](https://github.com/kubeovn/kube-ovn/commit/7006e9cb8af3b0889d3b1c44474d93ba54cb2fa5) fix(helm): add get on crd for ovn-cr (#4816)
 * [a3d301c00](https://github.com/kubeovn/kube-ovn/commit/a3d301c0003442fc30fcb5ca22f3b411a28a74fa) build(deps): bump github.com/brianvoe/gofakeit/v7 from 7.0.0 to 7.1.2 (#4818)
 * [b5c977863](https://github.com/kubeovn/kube-ovn/commit/b5c97786350a65d5d1462686ccb33d1e3a04bfd9) docs: updated CHANGELOG.md (#4815)
 * [1218d0b20](https://github.com/kubeovn/kube-ovn/commit/1218d0b20b3dff0366e6b82c046a44afa702f280) add some unittests (#4790)
 * [35fcdc1fb](https://github.com/kubeovn/kube-ovn/commit/35fcdc1fbe606802ed32f282285bfaac0fa25f43) bump k8s to v1.31.4 (#4813)
 * [922b5578e](https://github.com/kubeovn/kube-ovn/commit/922b5578e6422ea8903327e2421b793b8af1a918) build(deps): bump github.com/docker/docker (#4812)
 * [57107d297](https://github.com/kubeovn/kube-ovn/commit/57107d297c7506d2946489ab0ae747e8341f27a4) build(deps): bump github.com/onsi/gomega from 1.36.0 to 1.36.1 (#4811)
 * [ac81ac29f](https://github.com/kubeovn/kube-ovn/commit/ac81ac29f9027fc23f069eb255ccdf20f91b59e7) fix issue 4803: The two names should have a containment relationship (#4807)
 * [bc6379ba5](https://github.com/kubeovn/kube-ovn/commit/bc6379ba57205cb365aee7f6be8193ee975417d1) ci: fix names of artifacts uploaded by vpc egress agteway e2e test (#4799)
 * [1844fe7d4](https://github.com/kubeovn/kube-ovn/commit/1844fe7d4c9925af5e3812f56c38002b81284091) base: fix underlay network break during upgrade from v1.12 (#4797)
 * [69b727140](https://github.com/kubeovn/kube-ovn/commit/69b727140ee1eeaf1c7eda702b0cb2ee398da7d8) build(deps): bump gopkg.in/k8snetworkplumbingwg/multus-cni.v4 (#4806)
 * [322fb8e90](https://github.com/kubeovn/kube-ovn/commit/322fb8e900c94d332e29e6e338b3cef72d9250d2) auto detect kubevirt install (#4791)
 * [6b6bf2d8b](https://github.com/kubeovn/kube-ovn/commit/6b6bf2d8bdd044db8fe0e1a1b5fd966b9f6177d7) build(deps): bump google.golang.org/grpc from 1.68.0 to 1.68.1 (#4796)
 * [6016e30c9](https://github.com/kubeovn/kube-ovn/commit/6016e30c993b4adc79df529ad8e6c660cf389684) build(deps): bump golang.org/x/sys from 0.27.0 to 0.28.0 (#4792)
 * [86ad84bd5](https://github.com/kubeovn/kube-ovn/commit/86ad84bd552a3c1a02f8421fe6a43d8a240db6b9) ci: report unittest coverage (#4783)
 * [6ee3b3636](https://github.com/kubeovn/kube-ovn/commit/6ee3b363685542b64d932b610fbe75c42d74bc73) bump go to 1.23.4 (#4781)
 * [f73fbaceb](https://github.com/kubeovn/kube-ovn/commit/f73fbaceb84cecf35d0c3df28b55f6ee18c86622) build(deps): bump github.com/osrg/gobgp/v3 from 3.31.0 to 3.32.0 (#4787)
 * [bd0bd0323](https://github.com/kubeovn/kube-ovn/commit/bd0bd0323d4d45eee2d8b7bd2466ac7a58ae92f7) build(deps): bump sigs.k8s.io/controller-runtime from 0.19.2 to 0.19.3 (#4788)
 * [e642eef8f](https://github.com/kubeovn/kube-ovn/commit/e642eef8f2ef595b6b5ba726d37b37a7365788b0) ci: bump kubevirt to v1.4.0 (#4782)
 * [e34d5624a](https://github.com/kubeovn/kube-ovn/commit/e34d5624a7135e916c11d1476adb33c2cb70bfc7) feature: VPC Egress Gateway (#4692)
 * [ae0fc894d](https://github.com/kubeovn/kube-ovn/commit/ae0fc894d4163bdef89dad393ed0ff482222f418) build(deps): bump kernel.org/pub/linux/libs/security/libcap/cap (#4779)
 * [7b633257b](https://github.com/kubeovn/kube-ovn/commit/7b633257b895fd45383a5f27563c639cabefca15) add kubevirt live migration optimize (#4773)
 * [a80490a4e](https://github.com/kubeovn/kube-ovn/commit/a80490a4e2ea673debe37d8dc572b885254b6998) build(deps): bump github.com/stretchr/testify from 1.9.0 to 1.10.0 (#4765)
 * [6df02a803](https://github.com/kubeovn/kube-ovn/commit/6df02a803ac4012061b9a823c8ad2e8b23251c09) build(deps): bump github.com/onsi/gomega from 1.35.1 to 1.36.0 (#4766)
 * [5c827f18e](https://github.com/kubeovn/kube-ovn/commit/5c827f18e575d3c0f3262d7833f40dfd17affdd0) docs: updated CHANGELOG.md (#4763)
 * [95ceeffac](https://github.com/kubeovn/kube-ovn/commit/95ceeffaca14ff6f342681bedb9c5a6d281fe162) remove capability SYS_MODULE (#4744)
 * [35f181ba7](https://github.com/kubeovn/kube-ovn/commit/35f181ba74bec18090cfb9f87d1cc2f66c637bc0) add observed generation into condition (#4751)
 * [6c55c613c](https://github.com/kubeovn/kube-ovn/commit/6c55c613c28f618570497a643400097b0c8f8fde) build(deps): bump k8s from 1.31.2 to 1.31.3 (#4754)
 * [623039c09](https://github.com/kubeovn/kube-ovn/commit/623039c09a0d273041c7ad5f50f82c8f25f90c28) ci: add release-1.13 to workflows (#4746)
 * [f6e2a253b](https://github.com/kubeovn/kube-ovn/commit/f6e2a253b9fbf8ffed28eb00edc5a36700e59def) build(deps): bump sigs.k8s.io/controller-runtime from 0.19.1 to 0.19.2 (#4757)
 * [be8439cd6](https://github.com/kubeovn/kube-ovn/commit/be8439cd626e4d1adfd95dcbb080f86bb1b3ddc6) build(deps): bump github.com/prometheus-community/pro-bing (#4758)
 * [06ac37d42](https://github.com/kubeovn/kube-ovn/commit/06ac37d42cfe14d5ab21a0847333ead9f6239662) build(deps): bump github.com/onsi/ginkgo/v2 from 2.21.0 to 2.22.0 (#4755)
 * [c3ae8059a](https://github.com/kubeovn/kube-ovn/commit/c3ae8059a865faf62e115f5a3e1ad4aaeb9a34df) split type definitions into separate files (#4750)
 * [17606a126](https://github.com/kubeovn/kube-ovn/commit/17606a1265572da83d44e8d5d0b3aa5b651fac7f) build(deps): bump aquasecurity/trivy-action from 0.28.0 to 0.29.0 (#4752)
 * [671ad64f5](https://github.com/kubeovn/kube-ovn/commit/671ad64f5654843c38865418d906b693042e8eec) add not found err check for lb-svc (#4748)
 * [6e158d6c8](https://github.com/kubeovn/kube-ovn/commit/6e158d6c8cd290272689db800a1a8db8ba6d7a92) update release script (#4749)
 * [623ab01ba](https://github.com/kubeovn/kube-ovn/commit/623ab01ba6d80cbdb24b6e23fed312393d083d00) vpc: add support for dedicated BFD LRP (#4717)
 * [bf9bea1e3](https://github.com/kubeovn/kube-ovn/commit/bf9bea1e3d659fa81779f6aa260ed17a1716a884) docs: updated CHANGELOG.md (#4743)
 * [6481c4a4e](https://github.com/kubeovn/kube-ovn/commit/6481c4a4eda95534e2e08f72937daf10673b0a3a) update release scripts
 * [3d92364b6](https://github.com/kubeovn/kube-ovn/commit/3d92364b67f07f0ea1f2dd1dd0ab9633cdf7ccad) prepare for next release
 * [ae4ce3770](https://github.com/kubeovn/kube-ovn/commit/ae4ce3770c71de6c2c9ba456d11bf30b651f83ab) add loop check for tunnel nic (#4736)
 * [1bf233b3c](https://github.com/kubeovn/kube-ovn/commit/1bf233b3c1ab9cb1eddae1be3ba1230000124ac7) fix vm ip lost and show more log about ip delete (#4737)
 * [b92366b0b](https://github.com/kubeovn/kube-ovn/commit/b92366b0b4d92cd5583ae7076dd95bec95a667f4) add adapt with nodes and networks for anp/banp egress rules (#4731)
 * [5fe31c998](https://github.com/kubeovn/kube-ovn/commit/5fe31c998cacee0f598177093887f208023a3faa) ut: fix duplicate lsp name (#4734)
 * [11933d11c](https://github.com/kubeovn/kube-ovn/commit/11933d11c750cf89db4d6e017447eb7cfe820590) build(deps): bump google.golang.org/protobuf from 1.35.1 to 1.35.2 (#4733)
 * [a9e73db72](https://github.com/kubeovn/kube-ovn/commit/a9e73db72a0ebdb88c1cdeb0e9dac23261abd5cd) kube-ovn-cni will panic if cidr is invalid (#4729)
 * [e9a6c3697](https://github.com/kubeovn/kube-ovn/commit/e9a6c3697fb70daadd25a33655f5799185ef8617) build(deps): bump kubevirt.io/api from 1.3.1 to 1.4.0 (#4732)
 * [393db6ab1](https://github.com/kubeovn/kube-ovn/commit/393db6ab15be4d3ddd39219526c6f2029da85bdb) lb-svc is compatible with dual stack subnets (#4724)
 * [ccb85730f](https://github.com/kubeovn/kube-ovn/commit/ccb85730fdd0d1239c42d9855a40b11c73b56269) [bugfix] Optimize gc method at  port group and node (#4722)
 * [92a0d954d](https://github.com/kubeovn/kube-ovn/commit/92a0d954d6633fdc6657392585e6bec2a3eda511) The eip is not cleaned after eip is deleted (#4718) (#4719)
 * [e806b5336](https://github.com/kubeovn/kube-ovn/commit/e806b5336b54ca594471b520e7b16904e78473a7) base: add ovs-sandbox (#4712)
 * [427ea72db](https://github.com/kubeovn/kube-ovn/commit/427ea72db65b7261a3d3ea2300325ef2f6742017) build(deps): bump kernel.org/pub/linux/libs/security/libcap/cap (#4725)
 * [f50da5390](https://github.com/kubeovn/kube-ovn/commit/f50da53908ee61a3b205f17420a1ebd629e10839) ci: bump golangci-lint to v1.62.0 (#4723)
 * [5c54e01c4](https://github.com/kubeovn/kube-ovn/commit/5c54e01c447ce63d1643c7f8e360b40244a1cde4) base: update ovn patches (#4721)
 * [2edbf3fe0](https://github.com/kubeovn/kube-ovn/commit/2edbf3fe0a6050c0f69daf7f67bb46ca54ba7864) ci: add connectivity e2e test (#4658)
 * [505495756](https://github.com/kubeovn/kube-ovn/commit/5054957568ad442279666e0cebca47427a9c9de3) ci: bump kind to v0.25.0 (#4720)
 * [9e55f4ae9](https://github.com/kubeovn/kube-ovn/commit/9e55f4ae9faab6ca4bc9797ce679d2ee0529b1c5) bump go to 1.23.3 (#4716)
 * [ad7138bce](https://github.com/kubeovn/kube-ovn/commit/ad7138bcee1a5b6f4cdf9dd056d7fc734bb9bb70) use metav1.LabelSelectorAsSelector() to convert label selector (#4709)
 * [9c4a239e9](https://github.com/kubeovn/kube-ovn/commit/9c4a239e97701349cb04891a264d6434895b5dd3) minor fixes in vpc pod probe e2e test (#4682)
 * [d60e328ed](https://github.com/kubeovn/kube-ovn/commit/d60e328edb0c9b4aeb227911e6a47cc35e95297d) Fix ovn nat e2e (#4715)
 * [f14836aac](https://github.com/kubeovn/kube-ovn/commit/f14836aacbbe7e1f353d5a0e648ac8ac0d95d4f7) base: bump gobgp to v3.31.0 (#4708)
 * [262961170](https://github.com/kubeovn/kube-ovn/commit/2629611708a0dda0a9c8ecdf21e125e8bea1a4d3) add ut in pkg/ovs 24.11.7 (#4713)
 * [aaa67e613](https://github.com/kubeovn/kube-ovn/commit/aaa67e613aea0976b0f080ed7fdcedd50d1f6c1f) ut: add unit test for pkg/pvs (#4711)
 * [e2f5c69dc](https://github.com/kubeovn/kube-ovn/commit/e2f5c69dc1046b9343366fe0621e6e75da6a5fcd) ut: fix duplicate LSP name (#4710)
 * [d5e00c246](https://github.com/kubeovn/kube-ovn/commit/d5e00c246361955048f3ad3ba4b0f1ed392d7279) add ut for ovn-nb-logical_router_policy.go (#4706)
 * [d81bb1365](https://github.com/kubeovn/kube-ovn/commit/d81bb13654c743a8aa3630df4c5a572314da7db7) feat(natgw): extended nexthop support to advertise ipv4 NRLI over ipv6 sessions (#4704)
 * [1dc93da82](https://github.com/kubeovn/kube-ovn/commit/1dc93da8270b5ac8950054c738d4464b74647736) ci: re-add ovn vpc nat gw e2e to push needs (#4703)
 * [74b5c495a](https://github.com/kubeovn/kube-ovn/commit/74b5c495ada57fd1fa1a2c6cf9d1a0a01f9ecf60) ut: add unit test (#4702)
 * [5734a2765](https://github.com/kubeovn/kube-ovn/commit/5734a2765da80125fbab5055876665e1c2e75cb3) Make NAT gateway generation prefix a variable (#4557)
 * [042bdff2b](https://github.com/kubeovn/kube-ovn/commit/042bdff2ba073e475834024549825bcf8ea11d46) build(deps): bump github.com/osrg/gobgp/v3 from 3.30.0 to 3.31.0 (#4700)
 * [63d715864](https://github.com/kubeovn/kube-ovn/commit/63d71586414381e4b62627b9bcadce5f4c84e409) ut: fix use of closed network connection (#4698)
 * [d71f0a2f5](https://github.com/kubeovn/kube-ovn/commit/d71f0a2f52cdd3a356767805bd2dbe4c25acacc0) fix concurrent issue in ovn-nb_test.go (#4697)
 * [5db327096](https://github.com/kubeovn/kube-ovn/commit/5db327096035a7c5bc335e2bc928d1ea326366bb) add ut for ovn-nb-logical_router_port_test.go (#4696)
 * [6842df3c9](https://github.com/kubeovn/kube-ovn/commit/6842df3c940832f68be421c924de94d6b9e86408) add process for lb-svc ports update (#4660)
 * [451a463f4](https://github.com/kubeovn/kube-ovn/commit/451a463f4912b33e617ff074962e074b3a56fa8f) add WaitUntil check for ns (#4695)
 * [e4d1d3836](https://github.com/kubeovn/kube-ovn/commit/e4d1d3836b14ee45e50efb5319f4b79f5a3667fa) add ut for ovn-nb-logical_router_policy (#4693)
 * [1f0eb6f93](https://github.com/kubeovn/kube-ovn/commit/1f0eb6f93f4dc4d6edb1968f56dfa79faa270b88) ut: add unit test for ovn-nb-logical_router.go (#4689)
 * [b057e4ce6](https://github.com/kubeovn/kube-ovn/commit/b057e4ce6d3d29cb71c11656f51da29380ad0261) build(deps): bump github.com/onsi/gomega from 1.35.0 to 1.35.1 (#4688)
 * [63492716d](https://github.com/kubeovn/kube-ovn/commit/63492716d834a0658a2c0e6f0dbc806c318c7f6b) add fail nb client ut case in ovn-nb-logical_switch_port_test.go (#4685)
 * [f7f7fa877](https://github.com/kubeovn/kube-ovn/commit/f7f7fa877ff6fa6b6596894836915cf76fa9a963) base: compile OpenBFDD locally (#4683)
 * [ef4bb5a63](https://github.com/kubeovn/kube-ovn/commit/ef4bb5a631413377edbdadb7bf09c09cee4476fa) fix: cache eip should not update (#4673)
 * [4c8a49b4c](https://github.com/kubeovn/kube-ovn/commit/4c8a49b4c1db433d87edaeeb7373528ceda9cc97) build(deps): bump github.com/onsi/gomega from 1.34.2 to 1.35.0 (#4678)
 * [d90ae7c94](https://github.com/kubeovn/kube-ovn/commit/d90ae7c94d8b0d982a4d53caa592cdcb90974b52) Revert "build(deps): bump github.com/Microsoft/hcsshim from 0.12.7 to 0.12.9 …" (#4681)
 * [768880820](https://github.com/kubeovn/kube-ovn/commit/768880820e7c135756e34bde81d765e1f7ed329c) build(deps): bump github.com/Microsoft/hcsshim from 0.12.7 to 0.12.9 (#4677)
 * [bb398bf2b](https://github.com/kubeovn/kube-ovn/commit/bb398bf2bdc9cd417ebf1d65ea534690abee5cb8) build(deps): bump github.com/onsi/ginkgo/v2 from 2.20.2 to 2.21.0 (#4679)
 * [45c54b6c5](https://github.com/kubeovn/kube-ovn/commit/45c54b6c5da803d04e38471dae7416a3f8bdb7a8) add ut for subnet.go and ovn-nb-acl.go (#4674)
 * [f966ecc30](https://github.com/kubeovn/kube-ovn/commit/f966ecc306ab2c13ce57ae8dd519d5fc18a0e39d) add e2e for subnet with namespaceSelector (#4665)
 * [3d92683af](https://github.com/kubeovn/kube-ovn/commit/3d92683af5de2587c222d156a437e8e3d8d9744c) ut: add ovs and  ovn nb and cidr check (#4671)
 * [535048dcc](https://github.com/kubeovn/kube-ovn/commit/535048dcc3a759bc6f268756503cf4ba283d9ccf) add unit for ovn-nb-logical_router_route.go (#4672)
 * [bc2eb2b80](https://github.com/kubeovn/kube-ovn/commit/bc2eb2b80f2acb0d1d23cab3d57174ab9dbd42ed) build(deps): bump kernel.org/pub/linux/libs/security/libcap/cap (#4667)
 * [406c3794b](https://github.com/kubeovn/kube-ovn/commit/406c3794b732edf8b67e32aeb6e2eb78d1cd0be2) add unit test for ovn-nb-nat.go (#4664)
 * [a6e5276e1](https://github.com/kubeovn/kube-ovn/commit/a6e5276e111c790f09ada33e38f504bc5575f672) ut: add ovn nat (#4650)
 * [44cd6e5c2](https://github.com/kubeovn/kube-ovn/commit/44cd6e5c2ff111e6f75ac4ace8df129e70a8e087) add unit test for net.go (#4659)
 * [8e9f9479e](https://github.com/kubeovn/kube-ovn/commit/8e9f9479e610be1c0df4b44189c732de4c5dc275) add ut for ovn-nb-logical_router_port (#4651)
 * [1ffd88dad](https://github.com/kubeovn/kube-ovn/commit/1ffd88dad4584a9736da61ddf182d78f735d6411) base: fix ovn patch for setting ether dst addr after dnat (#4662)
 * [ab2c21293](https://github.com/kubeovn/kube-ovn/commit/ab2c21293ac5972e4047a32589f2b157dd7b8cb8) build(deps): bump sigs.k8s.io/controller-runtime from 0.19.0 to 0.19.1 (#4661)
 * [eb51ad211](https://github.com/kubeovn/kube-ovn/commit/eb51ad2111ea4f1d8b2c7d65b99a29e50cd46c46) bump k8s to v1.31.2 (#4654)
 * [62ab30a7d](https://github.com/kubeovn/kube-ovn/commit/62ab30a7d416100e4004c565d34d983f16ff9bb2) makefile: remove target kind-install-debug (#4656)
 * [92f6800bf](https://github.com/kubeovn/kube-ovn/commit/92f6800bf07a533dbdb0d86b2fa2aa95fb82e4af) build(deps): bump k8s.io/kubernetes in the k8s-io group (#4653)
 * [5cfb6b4c4](https://github.com/kubeovn/kube-ovn/commit/5cfb6b4c46854fee5fb8ef11e43f589dce70022d) ut: add logical_switch and port (#4644)
 * [0706c41ae](https://github.com/kubeovn/kube-ovn/commit/0706c41aea938336ff8a0e2d3c25a153e83c7ace) ut: add unit for ovn-nb-load_balancer (#4646)
 * [7003bca4f](https://github.com/kubeovn/kube-ovn/commit/7003bca4fafbd4b82f0c4ade14294887dd8d21ea) ci: ignore ovn vpc-nat-gw conformance e2e test (#4649)
 * [a7a8a9235](https://github.com/kubeovn/kube-ovn/commit/a7a8a9235566398f37f1512fb9266a303c6168fe) fix: bfd should use multi arch (#4647)
 * [c8288e89b](https://github.com/kubeovn/kube-ovn/commit/c8288e89b0e6960c235e1093bfeea56f529e29be) fix sometimes restart kube-ovn-cni ,ipsec not start (#4641)
 * [a9cd9a1a8](https://github.com/kubeovn/kube-ovn/commit/a9cd9a1a807535a8625c7b0237437edd47040d3e) fix: udp bad checksum on VXLAN interface (#4639)
 * [c414e6b96](https://github.com/kubeovn/kube-ovn/commit/c414e6b967a3b985365f2e546b2afbd51d258f0c) add e2e test for accessing nodeport when ovs/ovn component is down (#4642)
 * [d16bb52da](https://github.com/kubeovn/kube-ovn/commit/d16bb52da711da83222d404ca538186a16070f2b) add ovs vsctl , gw chassis, lb healthcheck ut (#4611)
 * [2ac7e2c0f](https://github.com/kubeovn/kube-ovn/commit/2ac7e2c0f589793ea0530b13a28ca4ddc6678afc) ut: fail fast (#4630)
 * [d02cffc28](https://github.com/kubeovn/kube-ovn/commit/d02cffc287780840a8c80c7aeaf97b1a13cb5ab4) use ovn-appctl for ovn components (#4638)
 * [79c13434b](https://github.com/kubeovn/kube-ovn/commit/79c13434bbdd9e805bdeef3e15a381195bd7872a) fix 'invalid memory address or nil pointer dereference' for namespace (#4636)
 * [7665b2b68](https://github.com/kubeovn/kube-ovn/commit/7665b2b6867fec44f88ce290bc612d76741a80f4) chart: fix missing CRD properties (#4637)
 * [ad7c89059](https://github.com/kubeovn/kube-ovn/commit/ad7c8905952be59e563a85f100703b4188c274ec) fix: CreateGatewayACL creates duplicate ndACL (#4592)
 * [aa1e67958](https://github.com/kubeovn/kube-ovn/commit/aa1e679586100754f48f07c27abb43c487549948) Fix accidental logical switch name conflict in ovn-nb-logical_switch_port_test.go (#4635)
 * [03aacacff](https://github.com/kubeovn/kube-ovn/commit/03aacacff0d8265d52a8cb3f10b88b80c55d7bc0) Add ut for ovn-nb-logical_switch_port.go date 24.10.16pm (#4625)
 * [c52c8c0cc](https://github.com/kubeovn/kube-ovn/commit/c52c8c0cc0267a46bcbc7d8be37cbde549e094c6) bump go.uber.org/mock to v0.5.0 (#4632)
 * [e0e8083f1](https://github.com/kubeovn/kube-ovn/commit/e0e8083f160c626cd4f30beddc21bb1031def1a1) docs: updated CHANGELOG.md (#4634)
 * [c9afcafbe](https://github.com/kubeovn/kube-ovn/commit/c9afcafbe4ad5e433cbd5313cdc9e1a1ebf191be) docs: updated CHANGELOG.md (#4633)
 * [ff6cc124a](https://github.com/kubeovn/kube-ovn/commit/ff6cc124adb0a413dd3e716ba6dfd9b7aa9c73dd) Refactor network policy matching logic (#4626)
 * [18f6ac918](https://github.com/kubeovn/kube-ovn/commit/18f6ac918447d968a331592afc86a18d474e1f92) team device not set unmanage (#4629)
 * [642973dbb](https://github.com/kubeovn/kube-ovn/commit/642973dbb3ca9e9da0a866e997db5468ff0aeab3) fix: once ovn-encap-csum is set false,it can't be reset to true (#4623)
 * [1c19fc70b](https://github.com/kubeovn/kube-ovn/commit/1c19fc70bcb69d0da105dba8dabb684caaca570f) build(deps): bump aquasecurity/trivy-action from 0.27.0 to 0.28.0 (#4628)
 * [d1af97bdd](https://github.com/kubeovn/kube-ovn/commit/d1af97bddea7c75aa4216033aa9c0b68fc789bba) bump go to 1.23.2 (#4412)
 * [0dbe49e68](https://github.com/kubeovn/kube-ovn/commit/0dbe49e6830200592e6a1acd53e1163efdef8277) ovn-ic: do not restart ovn-controller (#4624)
 * [77abe0c04](https://github.com/kubeovn/kube-ovn/commit/77abe0c048866c398f6baaf49bed76f75cf0b1eb) fix: ut dial tcp longer (#4614)
 * [1f227958d](https://github.com/kubeovn/kube-ovn/commit/1f227958deba7239ddd71ac41de4fa9e7ea63c3d) fix permission issue in kube-ovn installation script and chart (#4613)
 * [ea80cf25f](https://github.com/kubeovn/kube-ovn/commit/ea80cf25ffb2bab78991dd911a04cc57b3a8ff95) docs: updated CHANGELOG.md (#4620)
 * [5df224b16](https://github.com/kubeovn/kube-ovn/commit/5df224b16fee9cca8b06e01c8367efd915b35305) docs: updated CHANGELOG.md (#4619)
 * [aa8dce0bf](https://github.com/kubeovn/kube-ovn/commit/aa8dce0bff383a0b7f93f7f3aa0202aaee0c4893) build(deps): bump github.com/prometheus/client_golang (#4617)
 * [f9a5a3b97](https://github.com/kubeovn/kube-ovn/commit/f9a5a3b974098b7500e6c249fff6ec523304ca13) build(deps): bump github.com/k8snetworkplumbingwg/network-attachment-definition-client (#4618)
 * [63df5fe7d](https://github.com/kubeovn/kube-ovn/commit/63df5fe7d62cf07101145ea760291ed93beb9561) fix: ut (#4609)
 * [5cb13da09](https://github.com/kubeovn/kube-ovn/commit/5cb13da0958ec58e54bf67222df519df0ec6b571) fix memory overflow, add mac_binding related options to router (#4603)
 * [a6f13a67e](https://github.com/kubeovn/kube-ovn/commit/a6f13a67e25bdd86ff52351a8c06bc80ec7511e8) libovsdb: improve performance for listing logical router policies (#4599)
 * [08d8a546f](https://github.com/kubeovn/kube-ovn/commit/08d8a546fea22758685106b65d375d7eeb58c0d5) Add unit test in pkg/ovs/ovn-nb-logical_switch_port.go (#4602)
 * [071eb8e72](https://github.com/kubeovn/kube-ovn/commit/071eb8e7299ee1e8590e56590b606528d2c5490f) fix: sb ut failed (#4606)
 * [27379a662](https://github.com/kubeovn/kube-ovn/commit/27379a6625deca02cec448b5e34440f2b08d7bfb) vsctl ut (#4601)
 * [e8a9db09d](https://github.com/kubeovn/kube-ovn/commit/e8a9db09d5be04bd98283df8e65789080b737d36) fix logging (#4600)
 * [1251e7827](https://github.com/kubeovn/kube-ovn/commit/1251e78275d4894a3c85678e94847cbc118932dd) fix: dns overwrites dhcp when enabling dhcp and setting dhcp options at the same time (#4597)
 * [42ebd7cdd](https://github.com/kubeovn/kube-ovn/commit/42ebd7cdd406f88c34c6da606a8272a6a958a8ec) build(deps): bump aquasecurity/trivy-action from 0.26.0 to 0.27.0 (#4598)
 * [e91a7426f](https://github.com/kubeovn/kube-ovn/commit/e91a7426f713a09d0f4c4624afa76373e3a7f830) add namespaceSelector field in subnet (#4520)
 * [a30d88d03](https://github.com/kubeovn/kube-ovn/commit/a30d88d0307edff1ae2c37a6f61a102cd6bf7e78) chart: fix value of env OVN_DB_IPS (#4595)
 * [a9ca636d1](https://github.com/kubeovn/kube-ovn/commit/a9ca636d135a2ee3e49180153e75b6bac7e1e362) cni: always enable ipv6 (#4435)
 * [2e169a3f0](https://github.com/kubeovn/kube-ovn/commit/2e169a3f0e1c97c0edccf1e2f0e12e74688d70cd) ut: add ovn ic (#4593)
 * [8a3f2ea5a](https://github.com/kubeovn/kube-ovn/commit/8a3f2ea5a90c9815e74c9aad1ecff5126d86695f) chart: add missing sa imagePullSecrets (#4594)
 * [39067cd07](https://github.com/kubeovn/kube-ovn/commit/39067cd0775c7f122e3d53f73e0c4af94df1f271) add: fail case ut (#4589)
 * [4f80f33ce](https://github.com/kubeovn/kube-ovn/commit/4f80f33ce2a6c5b2197b68134390c5876f3fbe47) build(deps): bump aquasecurity/trivy-action from 0.25.0 to 0.26.0 (#4590)
 * [73a7ec69a](https://github.com/kubeovn/kube-ovn/commit/73a7ec69ac9a57d8a836aac660478a29aaee0f5f) build(deps): bump google.golang.org/protobuf from 1.34.2 to 1.35.1 (#4588)
 * [4cd5b0e5c](https://github.com/kubeovn/kube-ovn/commit/4cd5b0e5ca56c156154ab2d9472043e0eba57344) build(deps): bump aquasecurity/trivy-action from 0.24.0 to 0.25.0 (#4587)
 * [bd766ce4a](https://github.com/kubeovn/kube-ovn/commit/bd766ce4a8864b720f5919ccea817162b7c027ff) bump go to 1.22.8 (#4584)
 * [004877f00](https://github.com/kubeovn/kube-ovn/commit/004877f00eee4ffa5e195d55cc261aeb9578f8b2) Fix error log not printed nearby in util package (#4574)
 * [1d19ec090](https://github.com/kubeovn/kube-ovn/commit/1d19ec0904b9ff28686024d0c0a54e12f0fbe703) Fix error log not printed nearby in ipam package (#4585)
 * [e00511de3](https://github.com/kubeovn/kube-ovn/commit/e00511de3bc843b1f0a91511f0a81ff53a55b794) gc for security group (#4559)
 * [b2830a5d5](https://github.com/kubeovn/kube-ovn/commit/b2830a5d55643945cf78046315dcd6548db3990f) fix: test_list_chassis_with_no_entries (#4586)
 * [a676e630f](https://github.com/kubeovn/kube-ovn/commit/a676e630fbcda65ff95d207c9c56164863e80c75) add values for OVSDB_CON_TIMEOUT and OVSDB_INACTIVITY_TIMEOUT (#4568)
 * [edd2ccf9a](https://github.com/kubeovn/kube-ovn/commit/edd2ccf9a2d7cd6f72218cc91c6cd25061e9c125) Add genev_sys_6081 and vxlan_sys_4789 to cilium devices (#4575)
 * [4d5cdc7cd](https://github.com/kubeovn/kube-ovn/commit/4d5cdc7cd86720b654f1e55bad3e0f918c97fce2) fix slice init length (#4579)
 * [03e5c2d1a](https://github.com/kubeovn/kube-ovn/commit/03e5c2d1a95c9dbc816a0a1ba453ba3f6f52e3c1) ci: bump cilium, multus-cni, cert-manager and submariner (#4572)
 * [6f3d36bc0](https://github.com/kubeovn/kube-ovn/commit/6f3d36bc0d10dbfe1c8ac8fc3b349ccc2e323df5) ci: set trivy db repository to public.ecr.aws/aquasecurity/trivy-db:2 (#4570)
 * [38298e458](https://github.com/kubeovn/kube-ovn/commit/38298e45897d1f59282f7968a93e2e274e2b62a7) build(deps): bump github.com/osrg/gobgp/v3 from 3.29.0 to 3.30.0 (#4576)
 * [57cf0fdec](https://github.com/kubeovn/kube-ovn/commit/57cf0fdecff41ba4a3ab842aa0d13687e6580682) build(deps): bump golang.org/x/sys from 0.25.0 to 0.26.0 (#4581)
 * [1865f8dd6](https://github.com/kubeovn/kube-ovn/commit/1865f8dd61703db6c94f1915389d750a2c0e8404) build(deps): bump golang.org/x/time from 0.6.0 to 0.7.0 (#4580)
 * [18f4dbace](https://github.com/kubeovn/kube-ovn/commit/18f4dbace97779dccb3ec844700b6f5d715f85e7) build(deps): bump gopkg.in/k8snetworkplumbingwg/multus-cni.v4 (#4583)
 * [5d5480edf](https://github.com/kubeovn/kube-ovn/commit/5d5480edf2d2942b061a924518b82c2a9240be6e) build(deps): bump google.golang.org/grpc from 1.67.0 to 1.67.1 (#4577)
 * [519370471](https://github.com/kubeovn/kube-ovn/commit/519370471e9ca44184d65240e1130fbf32f374d4) Check network mask when init subnet (#4494)
 * [e90a58a6b](https://github.com/kubeovn/kube-ovn/commit/e90a58a6bc7071e4ee1626ad44cf9aef92adac0a) Fix error log not printed nearby in ovs package (#4562)
 * [3d12fa70d](https://github.com/kubeovn/kube-ovn/commit/3d12fa70da3b45ee4c741f4ebeda0e77c501338b) build(deps): bump github.com/Microsoft/hcsshim from 0.12.6 to 0.12.7 (#4561)
 * [62b80ca62](https://github.com/kubeovn/kube-ovn/commit/62b80ca62e373fe3b10ac9247868602f94fe3ad1) docs: updated CHANGELOG.md (#4571)
 * [49cba10a7](https://github.com/kubeovn/kube-ovn/commit/49cba10a77c18d9d67abcd1505a649ba0ceebdf8) ut: ovn sb suite (#4549)
 * [84203a086](https://github.com/kubeovn/kube-ovn/commit/84203a0861d59204288413fc3655acf5846f72c5) fix ovn0 add ipv6 local-link address (#4547)
 * [856dc10a1](https://github.com/kubeovn/kube-ovn/commit/856dc10a14ba8dc0219b98e0aa7c999614dd0175) build(deps): bump github.com/docker/docker (#4556)
 * [2f4a6c247](https://github.com/kubeovn/kube-ovn/commit/2f4a6c247f590327a511ef01e726c33d0e0454b2) docs: updated CHANGELOG.md (#4555)
 * [d38fef188](https://github.com/kubeovn/kube-ovn/commit/d38fef1881d04ee6af7db82b9d880a5ffb4663bd) fix kube-ovn-cni capabilities (#4550)
 * [514bc77de](https://github.com/kubeovn/kube-ovn/commit/514bc77de78fc09166f3ce89c3a68e10e9f385e1) docs: updated CHANGELOG.md (#4548)
 * [5a70de9c7](https://github.com/kubeovn/kube-ovn/commit/5a70de9c7c24dceefd45b8548047f0b8eccf2dad) allow user to set vxlan_sys_4789 tx off (#4543)
 * [2baa5bda7](https://github.com/kubeovn/kube-ovn/commit/2baa5bda791e2a9e0fddd5687e9b69f090b4fb97) add utils ut (#4544)
 * [73ce0ec8e](https://github.com/kubeovn/kube-ovn/commit/73ce0ec8ed20128319d182af9c3bd9d64849f814) build(deps): bump github.com/docker/docker (#4541)
 * [d38914a49](https://github.com/kubeovn/kube-ovn/commit/d38914a49f9dd883caadec78e445fa3d938beaa8) build(deps): bump google.golang.org/grpc from 1.66.2 to 1.67.0 (#4542)
 * [996da8152](https://github.com/kubeovn/kube-ovn/commit/996da8152ce9b97cc117e3213bbb2abc24d11d47) add arp ut (#4535)
 * [b90491266](https://github.com/kubeovn/kube-ovn/commit/b90491266d8e5df4296839f0088ccf09978f7eec) base: rebuild go binary deps from source (#4524)
 * [8ba16af2d](https://github.com/kubeovn/kube-ovn/commit/8ba16af2d3172ec7c266f27e0e6cc05cd2328814) fix speaker metrics (#4523)
 * [412f79904](https://github.com/kubeovn/kube-ovn/commit/412f79904de9d666aa55995f89bf6e99f88b0d89) [fix]: delete the NIC regardless of whether the Pod was found or not. (#4500)
 * [2fa78ab8d](https://github.com/kubeovn/kube-ovn/commit/2fa78ab8df74159277c7add9d48013c2ad9cfc11) add ut for ip.go (#4536)
 * [b406ce91d](https://github.com/kubeovn/kube-ovn/commit/b406ce91dd48eba075bc8e337df47482d7f94463) docs: updated CHANGELOG.md (#4534)
 * [94749490a](https://github.com/kubeovn/kube-ovn/commit/94749490a7ee2f3e7a56ad1eeda8478b4ad7256f) bump multus to v4.1.1 (#4532)
 * [5769b072b](https://github.com/kubeovn/kube-ovn/commit/5769b072b24b3da88619a4470c45a07a3f6ad6ac) check whether subnet exists before deleting vpc (#4533)
 * [a2b22ec2b](https://github.com/kubeovn/kube-ovn/commit/a2b22ec2b6277c1a40b6f6e63ba30f15dfb750ad) fix bugs and enhance unit test coverage for all functions in pkg/util/validator.go (#4505)
 * [640f6835e](https://github.com/kubeovn/kube-ovn/commit/640f6835e9f52f4ac8232e353a36c8a16369f503) build(deps): bump github.com/prometheus/client_golang (#4531)
 * [7f3c719ca](https://github.com/kubeovn/kube-ovn/commit/7f3c719cada9877655cbf2d10214594747317ff6) ut:add ut for func DialTCP and DialAPIServer in apk/util/k8s.go (#4515)
 * [c06857477](https://github.com/kubeovn/kube-ovn/commit/c06857477180a095d5ea1044808d9be07f040eb5) build(deps): bump gopkg.in/k8snetworkplumbingwg/multus-cni.v4 (#4529)
 * [217e4860b](https://github.com/kubeovn/kube-ovn/commit/217e4860bf1bc444e78e63f75a439e1660a0d8b7) add ut (#4525)
 * [4b06bd07b](https://github.com/kubeovn/kube-ovn/commit/4b06bd07b2971f7deff8bed7465bf4bb28197f33) add pod patch and exec ut (#4521)
 * [e2697fc91](https://github.com/kubeovn/kube-ovn/commit/e2697fc91a8a06cd240abca97ddc0e397a0d101a) Utils net ut (#4514)
 * [56f786977](https://github.com/kubeovn/kube-ovn/commit/56f786977f710c7c02780a26c2be4c8c4526b05f) add ginkgo ut coverage (#4510)
 * [a3c6676ab](https://github.com/kubeovn/kube-ovn/commit/a3c6676ab8e13f13ec55e9d9ad34f65c95fa69bd) fix mistake to announce ipv4 (#4497)
 * [544f2addc](https://github.com/kubeovn/kube-ovn/commit/544f2addcbcfeadc5e74ea7e5162f0d1097c260e) ut: add unit for dhcp_options logical_switch and ovn-nb_global (#4508)
 * [d032f05dd](https://github.com/kubeovn/kube-ovn/commit/d032f05dd3331c77e7862665163174c6301bf548) build(deps): bump github.com/k8snetworkplumbingwg/network-attachment-definition-client from 1.7.2 to 1.7.3 (#4522)
 * [f95e06c9a](https://github.com/kubeovn/kube-ovn/commit/f95e06c9a6d4f1640ee9c1c6603554812c090a5d) fix: upgrade submodule ginkgo (#4509)
 * [f3061ef67](https://github.com/kubeovn/kube-ovn/commit/f3061ef67cb9a3786852f20237d69909b9cfef51) fix: mac arm64 run netlink ut (#4511)
 * [bf75a65ee](https://github.com/kubeovn/kube-ovn/commit/bf75a65eedc5aa75eecaa679a2f85bafee91d6e0) bump k8s to v1.31.1 (#4512)
 * [14484fd35](https://github.com/kubeovn/kube-ovn/commit/14484fd35dfd74098dbe5750a2745e701f0ced48) build(deps): bump google.golang.org/grpc from 1.66.1 to 1.66.2 (#4517)
 * [6bf158781](https://github.com/kubeovn/kube-ovn/commit/6bf1587816e4f862075c9a451a5a9f79c52ff8fb) build(deps): bump github.com/k8snetworkplumbingwg/network-attachment-definition-client (#4518)
 * [2b24b2b4c](https://github.com/kubeovn/kube-ovn/commit/2b24b2b4c6b642fab05da3a53dad904005ffad7f) fix mcast e2e error (#4503)
 * [979513112](https://github.com/kubeovn/kube-ovn/commit/979513112c86b3c1f1f97ffc34f45da2b0bde8c4) fix: err log (#4507)
 * [486cde80a](https://github.com/kubeovn/kube-ovn/commit/486cde80a7357fb95de078cc556156c0359bdfbb) make: add ut cover (#4490)
 * [3e843e81d](https://github.com/kubeovn/kube-ovn/commit/3e843e81dcb60d8da15569934f9e4dcc7f54978d) build(deps): bump github.com/containerd/containerd from 1.7.21 to 1.7.22 (#4501)
 * [f16258302](https://github.com/kubeovn/kube-ovn/commit/f16258302203b32e0fadac462a7fd2fca12e218b) build(deps): bump google.golang.org/grpc from 1.66.0 to 1.66.1 (#4502)
 * [1cf7e5086](https://github.com/kubeovn/kube-ovn/commit/1cf7e50866230091fe158d9e8ac9b17185a2e67f) fix: gobgp CVE (#4498)
 * [dbcbc056a](https://github.com/kubeovn/kube-ovn/commit/dbcbc056a17854dc877db7c7917fb6f245563211) add mcast querier ip for multicast (#4375)
 * [b0bb154f4](https://github.com/kubeovn/kube-ovn/commit/b0bb154f438ac5b8686ac0cb34114f52e1d6360a) bump golangci-lint to v1.61.0 (#4496)
 * [0542c92ad](https://github.com/kubeovn/kube-ovn/commit/0542c92addedfd33420ec531cca0c92503e437bf) build(deps): bump github.com/docker/docker (#4495)
 * [3653e74dc](https://github.com/kubeovn/kube-ovn/commit/3653e74dc69338eea23287556279d7c9eb34358e) fix: support ptp networks with ipv4 /31 netmask and ipv6 /127 netmask  (#4425)
 * [5391b8015](https://github.com/kubeovn/kube-ovn/commit/5391b80157444e73cdecd12a3694595303326c5f) tproxy: support named port (#4487)
 * [2d1152262](https://github.com/kubeovn/kube-ovn/commit/2d11522625639c4a526865e30d0b99e83ef24932) fix: ipam ut name conflicts (#4489)
 * [baaaed0b5](https://github.com/kubeovn/kube-ovn/commit/baaaed0b570c12df90561d4e9aaf04f324e93412) add gobgp client cmd (#4460)
 * [70c7b979c](https://github.com/kubeovn/kube-ovn/commit/70c7b979c7a07776fe1eaafda3a6d8b2b898b3ee) ut: add unit test for bfd acl and address_set (#4461)
 * [827c50f23](https://github.com/kubeovn/kube-ovn/commit/827c50f23d95c0d3cfabd652da6cbacf06c25230) Add ipam (#4451)
 * [2214d2702](https://github.com/kubeovn/kube-ovn/commit/2214d270239ac63e1a1d5e343df4e4f7346b7fdc) Add ipam ut (#4455)
 * [6d4011f9a](https://github.com/kubeovn/kube-ovn/commit/6d4011f9ad1e5a7ba85f5c04bd80d6cbcdb3460a) bump go to 1.22.7 (#4482)
 * [d76bf6db6](https://github.com/kubeovn/kube-ovn/commit/d76bf6db68cc609a5cb3c6e8e4eeb3cebf2a0333) fix 解决部署遇到的不成功问题 (#4484)
 * [b9c3c2137](https://github.com/kubeovn/kube-ovn/commit/b9c3c21378356cc7c714dcf7c63dcde59549d74a) fix: arping reply may duplicate (#4477)
 * [fd5ce3169](https://github.com/kubeovn/kube-ovn/commit/fd5ce316947114cb963b804d3f45a20387c02eec) ci: bump cilium to v1.16.1 (#4337)
 * [ff2300bfb](https://github.com/kubeovn/kube-ovn/commit/ff2300bfbd340666f43d4f06f77212811fa363e8) fix bug: dpdk场景下添加删除ProviderNetwork失败 (#4466)
 * [1104627f4](https://github.com/kubeovn/kube-ovn/commit/1104627f44386f7959f56c17ea3059002771ab09) build(deps): bump github.com/prometheus/client_golang (#4480)
 * [1b46c42a1](https://github.com/kubeovn/kube-ovn/commit/1b46c42a13c719f5ba6e11c54ce6f706494c6438) build(deps): bump golang.org/x/mod from 0.20.0 to 0.21.0 (#4481)
 * [33272f110](https://github.com/kubeovn/kube-ovn/commit/33272f110457a7fc2407cc63036fc96d2a5ffccc) fix: kubectl-ko using kube-ovn-cni pod for nsenter (#4478)
 * [b3075e316](https://github.com/kubeovn/kube-ovn/commit/b3075e316319a9ebff3816bd2d0fbdd6dfdae852) remove incorrect error logging (#4473)
 * [c55e97aef](https://github.com/kubeovn/kube-ovn/commit/c55e97aefb2651a9065e16c1f636e1281585d2f5) ci: bump golangci-lint to v1.60.3 (#4474)
 * [4e5cc7bb9](https://github.com/kubeovn/kube-ovn/commit/4e5cc7bb94fb37d1313ea80fa7e45f016ac60572) build(deps): bump golang.org/x/sys from 0.24.0 to 0.25.0 (#4476)
 * [e99013eca](https://github.com/kubeovn/kube-ovn/commit/e99013ecaf6af0344ad6aa95b72f51d3e2bbf0e4) add anp/banp e2e case (#4347)
 * [be408f895](https://github.com/kubeovn/kube-ovn/commit/be408f895deb3e5e0a453423380f912abb189d89) when using vf, set pod_nic_type to sriov (#4463)
 * [31311598b](https://github.com/kubeovn/kube-ovn/commit/31311598b6a3bb3c99771b48f8e81acf110a779a) fix: create and parse coredns template (#4445)
 * [938ddada8](https://github.com/kubeovn/kube-ovn/commit/938ddada8e14a2fb55a8c1ed01b1b5269e945b0e) kubectl-ko trace using ovs-ovn instead of kube-ovn-cni (#4471)
 * [2cfbc91d8](https://github.com/kubeovn/kube-ovn/commit/2cfbc91d8716eea11850caf3e09b821e3ca62b81) build(deps): bump github.com/opencontainers/runc from 1.1.13 to 1.1.14 (#4469)
 * [286c34419](https://github.com/kubeovn/kube-ovn/commit/286c3441998fd65bb7d77a19cb7d6d786d51c249) build(deps): bump peter-evans/create-pull-request from 6 to 7 (#4467)
 * [ee6e590d0](https://github.com/kubeovn/kube-ovn/commit/ee6e590d0b34c7c8ba6eeb6cd3416d6f98fbb8f4) build(deps): bump github.com/onsi/gomega from 1.34.1 to 1.34.2 (#4458)
 * [50016edf3](https://github.com/kubeovn/kube-ovn/commit/50016edf3b34874532973e5dcb9e588f3c20b44f) build(deps): bump github.com/onsi/ginkgo/v2 from 2.20.1 to 2.20.2 (#4459)
 * [6ecb05e5a](https://github.com/kubeovn/kube-ovn/commit/6ecb05e5a94a826b75c1592c602c64fcdffdd5c7) docs: updated CHANGELOG.md (#4452)
 * [80df67b38](https://github.com/kubeovn/kube-ovn/commit/80df67b38fef49f3b2ae05959ddab1360a50c35f) metrics: do not export information if a subnet is not validated (#4444)
 * [365d17dbc](https://github.com/kubeovn/kube-ovn/commit/365d17dbc3109bd645185291cd1d35974859fdc9) add ut (#4448)
 * [07ec2a337](https://github.com/kubeovn/kube-ovn/commit/07ec2a337f650ba605f21c1c5715edd6278f4ea0) build(deps): bump google.golang.org/grpc from 1.65.0 to 1.66.0 (#4449)
 * [0df9de1f9](https://github.com/kubeovn/kube-ovn/commit/0df9de1f91dc2f75d94f1ec998308f74a1396e8d) build(deps): bump github.com/docker/docker (#4450)
 * [add87ee6e](https://github.com/kubeovn/kube-ovn/commit/add87ee6e8fbe29e75aaa425ca66fc487a181c44) ut: add unit test for ovn.go and ovs/util.go (#4440)
 * [72cbfa5e5](https://github.com/kubeovn/kube-ovn/commit/72cbfa5e51440ebeaa863fb8e38a75b4fe5f65d5) Add ipam ut (#4442)
 * [509a660f8](https://github.com/kubeovn/kube-ovn/commit/509a660f8ee7ab4664ff3bde92f180ca3ea753fd) build(deps): bump github.com/containerd/containerd from 1.7.20 to 1.7.21 (#4443)
 * [2e17cb4dc](https://github.com/kubeovn/kube-ovn/commit/2e17cb4dc542a95616bf8bb1f113a0bd76f044e0) fix: init default vpc (#4429)
 * [f7da30552](https://github.com/kubeovn/kube-ovn/commit/f7da30552ad35d398709d4db1ea4926830810f27) build(deps): bump github.com/vishvananda/netlink from 1.2.1 to 1.3.0 (#4438)
 * [aa08ad652](https://github.com/kubeovn/kube-ovn/commit/aa08ad652a500034fb845627b59baad250907c13) build(deps): bump github.com/prometheus/client_golang (#4439)
 * [b197ad34e](https://github.com/kubeovn/kube-ovn/commit/b197ad34e20bac6ce0e38110f2ed40ded63c3a19) build(deps): bump kubevirt.io/api from 1.3.0 to 1.3.1 (#4431)
 * [6f669c705](https://github.com/kubeovn/kube-ovn/commit/6f669c7050487b354618df733adcdf291b3c53db) feat: allow setting default subnet for custom vpc (#4171)
 * [c1a367b52](https://github.com/kubeovn/kube-ovn/commit/c1a367b52910039b792700559ffdc150322fe539) Add ipam ut (#4424)
 * [b21c96c3c](https://github.com/kubeovn/kube-ovn/commit/b21c96c3c3e8a8a8a682580c6b985beb8425a162) bump github.com/vishvananda/netlink to v1.2.1 (#4434)
 * [748aaf4bd](https://github.com/kubeovn/kube-ovn/commit/748aaf4bd7990c9be4f4aeab83c7220d6d70ed65) bump kubevirt to v1.3.1 (#4433)
 * [41a23b88b](https://github.com/kubeovn/kube-ovn/commit/41a23b88b3c1ccb9f7caf1459fd6d6b1712b8f7d) add unit test for ovn-nb-port_group (#4427)
 * [f93c39336](https://github.com/kubeovn/kube-ovn/commit/f93c3933678c6ddb88c9d50f19a4c22b808ff7ad) add unit test for NewAnpACLMatch (#4386)
 * [6355a337a](https://github.com/kubeovn/kube-ovn/commit/6355a337a58e6d0c4b0d1686da1ca872f412b9c2) build(deps): bump github.com/onsi/ginkgo/v2 from 2.20.0 to 2.20.1 (#4432)
 * [5d5b17236](https://github.com/kubeovn/kube-ovn/commit/5d5b17236504d0e6ad263ff01e46be0ac4abfe90) vpc-nat-gateway: use iptables-legacy for centos 7 (#4428)
 * [6bfbfe815](https://github.com/kubeovn/kube-ovn/commit/6bfbfe815532287a2b8f0f11f09bc4d74d0d5c77) feat(bgp): logger & enhance nadProvider config (#4352)
 * [46b8f569c](https://github.com/kubeovn/kube-ovn/commit/46b8f569c09571e95abeabcf5df4f4e526a804a8) netpol: add allow acl rules for u2o logical gateway (#4420)
 * [226601a64](https://github.com/kubeovn/kube-ovn/commit/226601a647af667377b6999457f3cd8b772e43fe) build(deps): bump github.com/prometheus/client_golang (#4423)
 * [7bb3cb588](https://github.com/kubeovn/kube-ovn/commit/7bb3cb588d33ccd61302bdd8b110a5a4ba2d7f75) build(deps): bump github.com/Microsoft/hcsshim from 0.12.5 to 0.12.6 (#4422)
 * [e5eba5f10](https://github.com/kubeovn/kube-ovn/commit/e5eba5f105db7f037e02776d801734f9fe3100bb) vpc-nat-gateway: do not add routes for underlay subnets (#4416)
 * [0ec244071](https://github.com/kubeovn/kube-ovn/commit/0ec244071cc6916ffa5f5353fc4b0d134e6fec5d) Makefile: simplify underlay u2o installation (#4419)
 * [b077d2f46](https://github.com/kubeovn/kube-ovn/commit/b077d2f46f87b39570be0c0654c8414aac367c32) bump kind to v0.24.0 (#4411)
 * [9e2cc200b](https://github.com/kubeovn/kube-ovn/commit/9e2cc200bad2ee107cc9d06dd52e9ef3b7098fb9) bump k8s to v1.31.0 (#4403)
 * [1c1e1a520](https://github.com/kubeovn/kube-ovn/commit/1c1e1a520a95a9e51b91f0d604498bd3f7c74fe7) simplicify if debug (#4402)
 * [9c80fc8cc](https://github.com/kubeovn/kube-ovn/commit/9c80fc8ccc14b744d0cc7ed7b9e4fe7358841110) fix kube-ovn-cni run fail on docker (#4409)
 * [8a8c06c74](https://github.com/kubeovn/kube-ovn/commit/8a8c06c747a1220e8f56ed397a6cb87c5bb79786) auto fix lint issues in local dev (#4385)
 * [25473b58a](https://github.com/kubeovn/kube-ovn/commit/25473b58a96d66824db396207b51bcf549256579) ci: bump golangci-lint to v1.60.1 (#4407)
 * [a02733fa4](https://github.com/kubeovn/kube-ovn/commit/a02733fa49ba32d17960e0798fbaae0352b41171) docs: updated CHANGELOG.md (#4401)
 * [ac1cc2c6b](https://github.com/kubeovn/kube-ovn/commit/ac1cc2c6bc95e05e0c848d3a838aaafb0fff806d) build(deps): bump github.com/docker/docker (#4400)
 * [71acfadb4](https://github.com/kubeovn/kube-ovn/commit/71acfadb49c845e20ff895e70ec4185c71693eca) build(deps): bump sigs.k8s.io/controller-runtime from 0.18.4 to 0.18.5 (#4398)
 * [06542452d](https://github.com/kubeovn/kube-ovn/commit/06542452db03e6094b14bca6c263299624478a2c) fix: uppercase ip address is not recommended (#4372)
 * [c9c2d3c61](https://github.com/kubeovn/kube-ovn/commit/c9c2d3c610638056a26ac7b71acc4632ccdd21e6) increase health probe timeout (#4397)
 * [0cf8c0f4e](https://github.com/kubeovn/kube-ovn/commit/0cf8c0f4ee761a628ebafe0f389404f1bdee881c) Support yusur smartnic (#4393)
 * [6c48efe58](https://github.com/kubeovn/kube-ovn/commit/6c48efe586639349fbdc100ffff64010e10ef202) add log (#4390)
 * [966b42ef1](https://github.com/kubeovn/kube-ovn/commit/966b42ef1bcae935cdecffcf93ee85ce35249463) cni-server: do not set sysctl variables (#4395)
 * [881644c5a](https://github.com/kubeovn/kube-ovn/commit/881644c5acbf744c84ddb8f497eaf8d74bfa0344) ci: fix git diff (#4388)
 * [db820a377](https://github.com/kubeovn/kube-ovn/commit/db820a37768c8733e90236307dfeb88faf2bc1d5) ipsec: fix chart installation failure (#4387)
 * [cbe4e37ce](https://github.com/kubeovn/kube-ovn/commit/cbe4e37cef588602a529a08384a084774e08f58b) change cidr to avoid conflict (#4392)
 * [5fc1c7717](https://github.com/kubeovn/kube-ovn/commit/5fc1c7717adfc19b39cb42e920f1bb0625a93d9a) docs: updated CHANGELOG.md (#4391)
 * [ef40eb019](https://github.com/kubeovn/kube-ovn/commit/ef40eb019b77fef3b1605512553b83e15f700ba1) increase the default probe interval for large cluster (#4389)
 * [33d7a282d](https://github.com/kubeovn/kube-ovn/commit/33d7a282d9ad5c46a62330a932bea626a9fd4401) serve pprof in metrics server if metrics server listens on 0.0.0.0 (#4384)
 * [f17033e56](https://github.com/kubeovn/kube-ovn/commit/f17033e5646164547477a06a488ec3d49dc6a867) fix EOF during TLS handshake caused by health check (#4381)
 * [ae88ad811](https://github.com/kubeovn/kube-ovn/commit/ae88ad8113e1fd799dd12b2791791775a5b78864) ut: add unit test for SGLostACL (#4378)
 * [502a90efe](https://github.com/kubeovn/kube-ovn/commit/502a90efeed87537e8f440b42dc6ca3ef1903309) refactor ipsec function (#4334)
 * [b32eb6ea6](https://github.com/kubeovn/kube-ovn/commit/b32eb6ea620554a0a042c1cfe3900c93f226f977) replace protocol check in netpol update (#4359)
 * [610ac19ea](https://github.com/kubeovn/kube-ovn/commit/610ac19eab12e50027c21ec4955313ed175b321d) fix controller-runtime logger not set (#4380)
 * [c3b3041a1](https://github.com/kubeovn/kube-ovn/commit/c3b3041a129cba1832ab487c033e747efa967a16) build(deps): bump golang.org/x/sys from 0.23.0 to 0.24.0 (#4383)
 * [5baabfb5f](https://github.com/kubeovn/kube-ovn/commit/5baabfb5f2e05e64f92d5ac405aea368d91f0166) fix: should delete subnet before vpc (#4370)
 * [96e4aed19](https://github.com/kubeovn/kube-ovn/commit/96e4aed19cab227db8bf559a4c1d180105fba708) fix: use route replace to avoid confilcts (#4371)
 * [818649cbb](https://github.com/kubeovn/kube-ovn/commit/818649cbbaf590805f38283c8fec405619f3f57f) Fix DISABLE_MODULES_MANAGEMENT (#4365)
 * [0faf8a038](https://github.com/kubeovn/kube-ovn/commit/0faf8a0387bde757f2fdb75665e7c330305ad7d8) fix: correct the priority of ACL associated with security group (#4366)
 * [ef4c017e2](https://github.com/kubeovn/kube-ovn/commit/ef4c017e2f655dc68cfd49a6c48012b62f60b02b) bump multus to v4.1.0 (#4377)
 * [33399c2dc](https://github.com/kubeovn/kube-ovn/commit/33399c2dc328983794837f311569181aff0df396) bump go to 1.22.6 (#4373)
 * [d50778168](https://github.com/kubeovn/kube-ovn/commit/d5077816805ab18e1c574faa9ad041b394973643) build(deps): bump golang.org/x/mod from 0.19.0 to 0.20.0 (#4367)
 * [9c4c6e0ce](https://github.com/kubeovn/kube-ovn/commit/9c4c6e0ce0d73cec74afda3ff09f2aced84a9aed) build(deps): bump golang.org/x/sys from 0.22.0 to 0.23.0 (#4368)
 * [631ff2ff1](https://github.com/kubeovn/kube-ovn/commit/631ff2ff1fe16c2b8d0de410efdf31fb9bb796af) build(deps): bump golang.org/x/time from 0.5.0 to 0.6.0 (#4369)
 * [ed16ce5a2](https://github.com/kubeovn/kube-ovn/commit/ed16ce5a2b8fb4288c88d87fefb174530796a007) vpc dns: remove unused variable (#4335)
 * [5181842ef](https://github.com/kubeovn/kube-ovn/commit/5181842ef55cfc228ea27e7ccfe440ec734fca88) build(deps): bump github.com/osrg/gobgp/v3 from 3.28.0 to 3.29.0 (#4362)
 * [55865b05a](https://github.com/kubeovn/kube-ovn/commit/55865b05ab7cf074e04c340d1c00aac946e8b111) build(deps): bump github.com/k8snetworkplumbingwg/network-attachment-definition-client (#4361)
 * [204451948](https://github.com/kubeovn/kube-ovn/commit/204451948d3f1f46dcac85ab18e4d74e6dc6ed08) fix lr-lb dnat with multiple distributed gateway ports (#4351)
 * [e4c67ca93](https://github.com/kubeovn/kube-ovn/commit/e4c67ca93351daab935d0b7efa1d9cd3495cad23) fix: load-balancer for dnat cleaned after restarting kube-ovn-controller (#4357)
 * [df2d066b4](https://github.com/kubeovn/kube-ovn/commit/df2d066b4487879afab308966060b29e9816b2d8) build(deps): bump github.com/onsi/gomega from 1.34.0 to 1.34.1 (#4350)
 * [386732c68](https://github.com/kubeovn/kube-ovn/commit/386732c68d469d9d8fc95677b2091806dc93ec6a) cni-server: disable udp-fragmentation-offload (#4342)
 * [fa429125c](https://github.com/kubeovn/kube-ovn/commit/fa429125cbc87df486a8face28e63eab49152b25) reduce image size by merging layers (#4346)
 * [b3b775e91](https://github.com/kubeovn/kube-ovn/commit/b3b775e911a5875ce85831ff915a7b2030352b04) reduce image sice by merging webhook binary into kube-ovn-cmd (#4345)
 * [6261a3ed3](https://github.com/kubeovn/kube-ovn/commit/6261a3ed3d2f4717ae09f5fd3b3d23019bf3c355) move initDefaultVpc into initOVN (#4344)
 * [be2ddb920](https://github.com/kubeovn/kube-ovn/commit/be2ddb920163c1d63aeb7233c50dd3229f0b3bbb) refactor eip and snat in default vpc with lr policy (#4343)
 * [28844842d](https://github.com/kubeovn/kube-ovn/commit/28844842d4f10375e9303c3d7fe263aa67c9423a) Add BGP capabilities to the NAT GW (#4285)
 * [cd3da988d](https://github.com/kubeovn/kube-ovn/commit/cd3da988d00b74b8d59b7212eb47a6f5aca0969b) build(deps): bump github.com/onsi/ginkgo/v2 from 2.19.0 to 2.19.1 (#4340)
 * [4ce36b654](https://github.com/kubeovn/kube-ovn/commit/4ce36b654d2f5c21f110cc871b2d397d7e0fb166) build(deps): bump github.com/onsi/gomega from 1.33.1 to 1.34.0 (#4341)
 * [41260c8c0](https://github.com/kubeovn/kube-ovn/commit/41260c8c040f6a10350c2fdfac9f6261112b1e5e) add acl action log annotation for netpol/anp/banp (#4338)
 * [006f74e80](https://github.com/kubeovn/kube-ovn/commit/006f74e80215e1c5a513ba9fd6c8e0bd7694b11d) Add ut (#4296)
 * [22a78ea80](https://github.com/kubeovn/kube-ovn/commit/22a78ea80d5166ec791457f32ef3b3b2e30bc982) build(deps): bump github.com/prometheus-community/pro-bing (#4339)
 * [9316e17cd](https://github.com/kubeovn/kube-ovn/commit/9316e17cd4917b11680d9314edb7bd46582ebf94) security: run as non-root user (#4326)
 * [d80de4599](https://github.com/kubeovn/kube-ovn/commit/d80de4599f5e2ce081117e25e51c8a5173763faf) build(deps): bump github.com/docker/docker (#4336)
 * [60040887a](https://github.com/kubeovn/kube-ovn/commit/60040887a6209452a50a9386462392eb6377c9d8) crd: fix missing fields (#4333)
 * [af95247e4](https://github.com/kubeovn/kube-ovn/commit/af95247e4608713daadec6028ecc3a0bc4670da9) pinger: fix process not terminated on sigkill (#4329)
 * [05984fc94](https://github.com/kubeovn/kube-ovn/commit/05984fc949f046bcf6e11948e61191fbd81ba166) fix: scripts (#4291)
 * [9229c4551](https://github.com/kubeovn/kube-ovn/commit/9229c45515fb32da63b56b14ac1aca21733a30de) build(deps): bump github.com/containernetworking/cni from 1.2.2 to 1.2.3 (#4328)
 * [4c5ec03e7](https://github.com/kubeovn/kube-ovn/commit/4c5ec03e7f5d524470c5291a947accec03adabc3) build(deps): bump github.com/docker/docker (#4327)
 * [1ba9e7c9e](https://github.com/kubeovn/kube-ovn/commit/1ba9e7c9ebb3f4b466e08b53b3b809358e3f432e) add new linter (#4322)
 * [d009a27f1](https://github.com/kubeovn/kube-ovn/commit/d009a27f185646ede2d80bc617598cd4796a88c6) add support for interface activation strategy (#4294)
 * [a21de7eaf](https://github.com/kubeovn/kube-ovn/commit/a21de7eaf42c872fd76bdf2997f3c0ba2f123759) ci: bump submariner to 0.18.0 (#4287)
 * [be7fd1741](https://github.com/kubeovn/kube-ovn/commit/be7fd17412bfb756208df00c8fa1df0082d69904) fix dialing https server (#4324)
 * [7537eb937](https://github.com/kubeovn/kube-ovn/commit/7537eb9376c39430f88e92c947ddbb8c98796801) refactor metrics (#4313)
 * [d9e761cbf](https://github.com/kubeovn/kube-ovn/commit/d9e761cbf422167b48f8fb2ca8dcf2f49a4b5d32) metrics: fix missing rbac for sa ovn (#4312)
 * [6824a58e2](https://github.com/kubeovn/kube-ovn/commit/6824a58e27b89a3ccc6dd4a3b1eae3d1000919a4) bump kubevirt to v1.3.0 (#4311)
 * [35fcef891](https://github.com/kubeovn/kube-ovn/commit/35fcef8912b92b111278300ac01bd8c8c0a09da9) update update-codegen-docker.sh (#4319)
 * [b2aec9258](https://github.com/kubeovn/kube-ovn/commit/b2aec9258804e0fba83ec9a31e97b0663b3d9084) ci: add timeout for executing cleanup.sh (#4321)
 * [0606d63b2](https://github.com/kubeovn/kube-ovn/commit/0606d63b2174232af7cbe6125635ecd3fd6e834c) e2e: skip north gateway test for versions prior to v1.12 (#4320)
 * [3a0a9104e](https://github.com/kubeovn/kube-ovn/commit/3a0a9104ed52cc4293eb85c39d522c2448639d5e) metrics: add support for secure serving (#4297)
 * [77d6f9985](https://github.com/kubeovn/kube-ovn/commit/77d6f9985ee98b2e962c38aaf084e94ec603ebc7) bump k8s to v1.30.3 (#4305)
 * [087e3323b](https://github.com/kubeovn/kube-ovn/commit/087e3323b9e7a90593c87da0fc12696585edc220) build(deps): bump k8s.io/kubernetes in the k8s-io group (#4309)
 * [109fd0a89](https://github.com/kubeovn/kube-ovn/commit/109fd0a892475e8662e8dfde171d4bf51e8da4f3) build(deps): bump github.com/moby/sys/mountinfo from 0.7.1 to 0.7.2 (#4310)
 * [2159ff5ef](https://github.com/kubeovn/kube-ovn/commit/2159ff5efe2511276a301cf3bf1db5b68cdce7ef) fix map concurrent read and write crash (#4302)
 * [26c04b1bf](https://github.com/kubeovn/kube-ovn/commit/26c04b1bfd7111ef1e26492369b05f3e9e71dd99) security: run as unprivileged (#3040)
 * [9e928d639](https://github.com/kubeovn/kube-ovn/commit/9e928d639265792c4f049d9aadddc1cd671812c6) add anp/banp code (#4290)
 * [1f58f75dc](https://github.com/kubeovn/kube-ovn/commit/1f58f75dcb28c34ba8c787a6bf9132490600d4ec) fix e2e
 * [83fba9ea3](https://github.com/kubeovn/kube-ovn/commit/83fba9ea38549ec8933054f6d9c1edb63dae23d5) Use route policy to reimplement northGateway (#4289)
 * [4f0d43d30](https://github.com/kubeovn/kube-ovn/commit/4f0d43d308adacd0e78b048b2beab77f71801f50) build(deps): bump github.com/Microsoft/hcsshim from 0.12.4 to 0.12.5 (#4299)
 * [e3222733d](https://github.com/kubeovn/kube-ovn/commit/e3222733df840c06e5c0684ae786f0087fba8947) base: bump ubuntu to 24.04 (#4293)
 * [ee78b5c3b](https://github.com/kubeovn/kube-ovn/commit/ee78b5c3b0e722c3f8ea35efd5bae3d7fae8a72c) underlay: set trunks of host nic port (#4282)
 * [b17975854](https://github.com/kubeovn/kube-ovn/commit/b17975854204aa916bd00dc57b7d9717b556cdda) do not set gatewayType for non ovn subnet and custom vpc subnet (#4283)
 * [8a9f4c1b1](https://github.com/kubeovn/kube-ovn/commit/8a9f4c1b1cec24dec94cef5ed3462da4eeecb3e3) fix using route table name as vpc name (#4284)
 * [f21c9c308](https://github.com/kubeovn/kube-ovn/commit/f21c9c308df1b065db0e2529d62dbf5c91859c6e) fix ovn lb not updated due to service update failure (#4280)
 * [78fe1f920](https://github.com/kubeovn/kube-ovn/commit/78fe1f920b920492412c7e7820cc840b307e1936) vpc-nat-gw: configure routes for subnets by pod annotations (#4278)
 * [e043e970f](https://github.com/kubeovn/kube-ovn/commit/e043e970faa860ce4d98e4c12a2efa20e6a44c5b) docs: updated CHANGELOG.md (#4276)
 * [e5029d569](https://github.com/kubeovn/kube-ovn/commit/e5029d569f481397222ae62a4033fde5488a2083) fix: set dhcp gateway to U2OInterconnectionIP when enabling dhcp and u2o (#4228)
 * [45a16a713](https://github.com/kubeovn/kube-ovn/commit/45a16a7137243cf818b0ddbb603c7c19400b1fdc) ci: run go mod tidy before building kubectl (#4274)
 * [51839ce3e](https://github.com/kubeovn/kube-ovn/commit/51839ce3e970b721f9009673fa3326bfa7185083) build(deps): bump aquasecurity/trivy-action from 0.23.0 to 0.24.0 (#4275)
 * [bc61393b2](https://github.com/kubeovn/kube-ovn/commit/bc61393b2e1475768a77bc639cf87061d3e0dce3) enable nat gw by default (#4273)
 * [e56a4cd13](https://github.com/kubeovn/kube-ovn/commit/e56a4cd13ca84dbb8d09f84c4e90012035f92f88) klog: set log file max size to 200MB (#4272)
 * [e74e6cad3](https://github.com/kubeovn/kube-ovn/commit/e74e6cad3fa337e0edd921b13365e437fc3bebe3) logrotate: set file size limit to 100M (#4271)
 * [0900ace50](https://github.com/kubeovn/kube-ovn/commit/0900ace50892a10305f76b92c13bae01f06ac89c) remove unused environment variable LOG_ROTATE (#4270)
 * [6b5356512](https://github.com/kubeovn/kube-ovn/commit/6b5356512b758c7c11ab75645eabc40497a02ab4) fix e2e (#4269)
 * [6ed2463ed](https://github.com/kubeovn/kube-ovn/commit/6ed2463ed10e47a295ceb2bb29a4b79729a34d72) ci: disable cgo when building kubelet and cni plugins (#4268)
 * [7006d9dd6](https://github.com/kubeovn/kube-ovn/commit/7006d9dd6efaf0bd2b0a95112e25df219bd5eb15) when router is deleted return success for static route deletion (#4266)
 * [cfd2bd70a](https://github.com/kubeovn/kube-ovn/commit/cfd2bd70a1faa1b49b1ed3899f68d29a38f504a2) ipam: fix ip not released for non-ovn subnets (#4265)
 * [e25bf15f7](https://github.com/kubeovn/kube-ovn/commit/e25bf15f7441f49e88dd247cb33feeeb53bc4c86) fix vpcnatgw image is not synced (#4264)
 * [5f8cef9d5](https://github.com/kubeovn/kube-ovn/commit/5f8cef9d5523dd3f0c309ead96944901bdfea80f) cni-server: break gateway check after pod is terminated/missing (#4248)
 * [c7310a606](https://github.com/kubeovn/kube-ovn/commit/c7310a606133776525e0b73867c0574ee375801f) do not create iptables rule for setting tcp mss (#4260)
 * [989d854cf](https://github.com/kubeovn/kube-ovn/commit/989d854cf60aa8c160971a1403f6fb87062b2776) fix invalid subnet not sync route (#4258)
 * [e4cf12180](https://github.com/kubeovn/kube-ovn/commit/e4cf12180487a2bfd1935a3799fab2c32bda36f4) do not allocate mac for non-ovn subnets (#4242)
 * [929e3362b](https://github.com/kubeovn/kube-ovn/commit/929e3362bac99e9c58452fd82921db120806d27f) refactor lb svc e2e (#4236)
 * [7de95baa6](https://github.com/kubeovn/kube-ovn/commit/7de95baa6fc2f3fbbb7eb61372d9ac5504b7f027) build(deps): bump golang.org/x/sys from 0.21.0 to 0.22.0 (#4257)
 * [082b1488f](https://github.com/kubeovn/kube-ovn/commit/082b1488f064d86a45bdfe81863cfbe2f0a4c001) build(deps): bump golang.org/x/mod from 0.18.0 to 0.19.0 (#4256)
 * [877862820](https://github.com/kubeovn/kube-ovn/commit/8778628206119c0ec4fb84fe144b30dd1f0a3bca) build kubectl and cni plugins from source if vuln found in the base image (#4253)
 * [607818c9c](https://github.com/kubeovn/kube-ovn/commit/607818c9c43d04186c29e34bee22bd2970701a3c) add retry check u2o openflow filter rule (#4255)
 * [11c4c7640](https://github.com/kubeovn/kube-ovn/commit/11c4c7640c09aeebe0aa6323f73f902c09f997e3) e2e case "should drop ARP/ND request from localnet port to LRP" (#4251)
 * [2b71a56ac](https://github.com/kubeovn/kube-ovn/commit/2b71a56ac78b1836df8b3af245751a435a33a332) check both sts name and UID when handling pod deletion (#4238)
 * [41755bfc0](https://github.com/kubeovn/kube-ovn/commit/41755bfc035068522ebb55773b9a14b2fde09ec1) do not exit if chassis is not found (#4246)
 * [d0e147603](https://github.com/kubeovn/kube-ovn/commit/d0e147603345dc9e1e33e077a1641a990bf27b63) build(deps): bump github.com/osrg/gobgp/v3 from 3.27.0 to 3.28.0 (#4244)
 * [09dd08fd7](https://github.com/kubeovn/kube-ovn/commit/09dd08fd7498a4efe793f9cb4926ee83c75feb8c) build(deps): bump github.com/docker/docker (#4245)
 * [d9b5e4f61](https://github.com/kubeovn/kube-ovn/commit/d9b5e4f61e218622a2352da09b3c0a99418e8d15) force deleted subnet left ip, can only be handled by human (#4157)
 * [b614b1bc2](https://github.com/kubeovn/kube-ovn/commit/b614b1bc21f34b53100de0d83825889fc6aa3745) vpc-nat-gateway: configure initial routes by pod annotations (#4240)
 * [751fa71f5](https://github.com/kubeovn/kube-ovn/commit/751fa71f50e5cb9afe99682f1164c6bd5165b273) vpc-dns: set route to apiserver by pod annotation (#4239)
 * [627e2fdbb](https://github.com/kubeovn/kube-ovn/commit/627e2fdbbb61be4eb55a19c785f54050a7c3aa58) update node labels/annotations by json merge patch (#4230)
 * [d857ffbbf](https://github.com/kubeovn/kube-ovn/commit/d857ffbbff3f26e3ca9098ab8fe97539a431107b) vpc-nat-gateway: print messgae to stderr (#4237)
 * [9bcee4116](https://github.com/kubeovn/kube-ovn/commit/9bcee4116e2901584a8ea7b586ae6f6e27ef3955) fix: init if nat gw exists (#4232)
 * [cc9ca11b5](https://github.com/kubeovn/kube-ovn/commit/cc9ca11b58449f5b5c294b637b29c8ccbdb0567f) lb svc: fix dead lock (#4235)
 * [06fd80b45](https://github.com/kubeovn/kube-ovn/commit/06fd80b4538721bb502190d10f176d1eb7d5da32) ci: dump/upload events on e2e test failure (#4231)
 * [cbad90e1b](https://github.com/kubeovn/kube-ovn/commit/cbad90e1b330b44d4dd4acbc7b37b887d1082bda) ci: collect and upload audit logs on e2e failure (#4227)
 * [eebe7038b](https://github.com/kubeovn/kube-ovn/commit/eebe7038b89e73de59d9aa2b5548edc335e90596) docs: updated CHANGELOG.md (#4234)
 * [3e894e23b](https://github.com/kubeovn/kube-ovn/commit/3e894e23b7070eca6ef652dc67fa82e28c727533) build(deps): bump github.com/docker/docker (#4233)
 * [fa8b8ef2c](https://github.com/kubeovn/kube-ovn/commit/fa8b8ef2cb24fd75c120214e9afb0c7a9be4af8f) e2e: retry traceroute (#4225)
 * [735868470](https://github.com/kubeovn/kube-ovn/commit/73586847099ee68d8a0772be14fd1b0857356bb3) fix ipv6 service ip not added to ovn lb vips due to pod cache not synced (#4223)
 * [816a4bfad](https://github.com/kubeovn/kube-ovn/commit/816a4bfad0c58e80be0cdcae003014ecbc7271dc) build(deps): bump github.com/containernetworking/cni from 1.2.1 to 1.2.2 (#4222)
 * [c81e91c9f](https://github.com/kubeovn/kube-ovn/commit/c81e91c9fbc12141af585698d33e29907f170318) build(deps): bump github.com/docker/docker (#4221)
 * [be0cf5387](https://github.com/kubeovn/kube-ovn/commit/be0cf53871407faabc1f976298e041da0e2ccb95) e2e: exec ovn-nbctl via kubectl exec svc/ovn-nb (#4218)
 * [42cb4daa3](https://github.com/kubeovn/kube-ovn/commit/42cb4daa3596ac1c6822a07f1acf6e973d614fde) e2e: add curl timeouts (#4216)
 * [81f1927f3](https://github.com/kubeovn/kube-ovn/commit/81f1927f3ee14a2cf6c4335b9114abd92a3501be) e2e: fix logging (#4220)
 * [c5f03705d](https://github.com/kubeovn/kube-ovn/commit/c5f03705df56cf546e64762cd4386ae4b395d7a9) pinger: fix warning logs (#4213)
 * [2e80b7efe](https://github.com/kubeovn/kube-ovn/commit/2e80b7efe6e4f8b2433b4bce184765c43b0c4bc8) delete unused scripts (#4214)
 * [80b1d62bb](https://github.com/kubeovn/kube-ovn/commit/80b1d62bb7fb79e15f04290c6e2aa17cbe05a9ab) base: clean ipset deb files (#4215)
 * [be5277aec](https://github.com/kubeovn/kube-ovn/commit/be5277aec1c56491bb707ffa3545c5778bc98ed9) ic: ensure db file is fixed (#4211)
 * [acf78727b](https://github.com/kubeovn/kube-ovn/commit/acf78727b30d70162d162c1a53b82fc127eab0c9) fix getting service cluster ips (#4206)
 * [a11abec4e](https://github.com/kubeovn/kube-ovn/commit/a11abec4e4d5c38d32e43b867f6077de2ca26c88) add unit test for ip_range (#4125)
 * [ccff4ed94](https://github.com/kubeovn/kube-ovn/commit/ccff4ed942a71516e5218cfbcd69b94409f9ccc9) fix: nil pointer when subnet is not ready (#4190)
 * [6dde684a9](https://github.com/kubeovn/kube-ovn/commit/6dde684a91ad171842846ab4aaa08ee2dabd85b5) e2e: fix getting ipv4 address from link (#4177)
 * [bbf29e7b3](https://github.com/kubeovn/kube-ovn/commit/bbf29e7b3ff5b8321a48466b94b838d7771a4c54) docs: updated CHANGELOG.md (#4203)
 * [f9e53dc59](https://github.com/kubeovn/kube-ovn/commit/f9e53dc598d874abca205e0898080c122be9df0b) remove unused annotation (#4191)
 * [f63172d12](https://github.com/kubeovn/kube-ovn/commit/f63172d12d4046227c86f3a3579aa10366c78c8e) add e2e test cases for nat outgoing routes (#4187)
 * [2942337fa](https://github.com/kubeovn/kube-ovn/commit/2942337fa0ef9b38eec721a752520efc2da47adf) e2e: add ginkgo.GinkgoHelper() (#4192)
 * [e1ddaafb8](https://github.com/kubeovn/kube-ovn/commit/e1ddaafb803104737fae1f71730138456da7390e) build(deps): bump k8s.io/klog/v2 in the k8s-io group (#4202)
 * [183b2e606](https://github.com/kubeovn/kube-ovn/commit/183b2e606eec377d1a49e2ac6dfc9be1eccbe993) fix vm not running after changing the subnet (#4199)
 * [9684aee46](https://github.com/kubeovn/kube-ovn/commit/9684aee46d059c4458c406956d8b5e507655c074) pinger: reset interface_rx_multicast_packets (#4198)
 * [5c6029e58](https://github.com/kubeovn/kube-ovn/commit/5c6029e58bf6869cbc1f649dc8251d1e41aecac5) fix kube-ovn-cni crash for newly added nodes , due to old legacy event in deleteNodeQueue (#4197)
 * [e16f3809a](https://github.com/kubeovn/kube-ovn/commit/e16f3809a90d4702b4a50c0b1700f4ccca14753c) build(deps): bump k8s.io/klog/v2 from 2.120.1 to 2.130.0 (#4184)
 * [6b73daeb8](https://github.com/kubeovn/kube-ovn/commit/6b73daeb8ca2e0ff97dacc6a270911b27764ce60) build(deps): bump github.com/containernetworking/cni from 1.2.0 to 1.2.1 (#4189)
 * [ee65995a9](https://github.com/kubeovn/kube-ovn/commit/ee65995a947646c14a4a0a7b7cee054e3565f36f) build(deps): bump github.com/containernetworking/plugins (#4188)
 * [5566029f6](https://github.com/kubeovn/kube-ovn/commit/5566029f683ad63e9712a8d7d746cf357ca7310f) base: bump cni plugins to v1.5.1 (#4185)
 * [b57964cab](https://github.com/kubeovn/kube-ovn/commit/b57964cabc6523fb07a8104072a74921dc46f4e1) fix u2o arp/nd e2e
 * [e0bcc2967](https://github.com/kubeovn/kube-ovn/commit/e0bcc2967422bcb5d9b95dc43ee7c666d52e12d9) update ovn upgrade porcess (#4132)
 * [7fe3f6e6c](https://github.com/kubeovn/kube-ovn/commit/7fe3f6e6c81e7f4659b7cb0ffb37692e51ed6d92) add helper function PodIPs() (#4178)
 * [2f3c5d48b](https://github.com/kubeovn/kube-ovn/commit/2f3c5d48bc5a96275ac80ad5d2290b34ab55557d) northd: skip arp/nd request for lrp addresses from localnet ports (#4158)
 * [bbbb0d509](https://github.com/kubeovn/kube-ovn/commit/bbbb0d509c69d8b3578cb1df6a25d918d4f2fca6) ci: check pod crashes on installation/e2e failure (#4160)
 * [3e933249e](https://github.com/kubeovn/kube-ovn/commit/3e933249eb4a79ae46296c1fe838b01a5ede6c9c) fix change vm subnet e2e (#4114)
 * [f4efa2c86](https://github.com/kubeovn/kube-ovn/commit/f4efa2c861b6dea14ada63fdcf1ce04698a56fbb) fix ipam init (#4155)
 * [a639831a3](https://github.com/kubeovn/kube-ovn/commit/a639831a38ba8467ba26cd08b3b2ed576f3217f4) use IsZero() to check whether the DeletionTimestamp is nil (#4176)
 * [8562d70d9](https://github.com/kubeovn/kube-ovn/commit/8562d70d95de9242f6569befa87172358f6f6f04) ci: push amd64 legacy image only for branches >= v1.13 (#4179)
 * [3917b8b5f](https://github.com/kubeovn/kube-ovn/commit/3917b8b5f81a8baadc3f8855edf3bac2b7fcbd0e) ut add ip range split (#4115)
 * [e6445d816](https://github.com/kubeovn/kube-ovn/commit/e6445d816a9b03b2ab6c48aac8a4826d0008f241) fix reconcile routes (#4168)
 * [b722ffb15](https://github.com/kubeovn/kube-ovn/commit/b722ffb15a859b4f35d9bbf224f658e0b6e743e5) add support for amd64 legacy cpu (#4170)
 * [0edbf4e90](https://github.com/kubeovn/kube-ovn/commit/0edbf4e909c9961736487bd77e222a549012cbeb) build(deps): bump github.com/docker/docker from 26.1.4+incompatible to 27.0.0+incompatible (#4173)
 * [10871e85a](https://github.com/kubeovn/kube-ovn/commit/10871e85abbda5d48fa47fb76c3b73a2567d3c3a) docs: updated CHANGELOG.md (#4175)
 * [49c9ff1d8](https://github.com/kubeovn/kube-ovn/commit/49c9ff1d890fd525c066192c70fb1c3f87b42c9c) e2e: fix retrieving docker network information (#4169)
 * [6ac97fbb2](https://github.com/kubeovn/kube-ovn/commit/6ac97fbb21b1aaa58a11262038f09d69c4406aab) docs: updated CHANGELOG.md (#4166)
 * [3da2fa298](https://github.com/kubeovn/kube-ovn/commit/3da2fa29880f1370638275bd3a46037ccba12bf5) fix lr policy recreation after kube-ovn-controller restarts (#4123)
 * [b6b789081](https://github.com/kubeovn/kube-ovn/commit/b6b78908169765f3a8599d76abfe9bcd8e954d72) bump k8s to v1.30.2 (#4164)
 * [5132e77ee](https://github.com/kubeovn/kube-ovn/commit/5132e77ee9534047ee8ef86fa51375490f7d8cc5) base: bump kubectl to v1.30.2 (#4163)
 * [ccd5132d0](https://github.com/kubeovn/kube-ovn/commit/ccd5132d088ed54e73d4958f1615e99ee91261f9) ci: fix retrieving docker network subnet/gateway (#4161)
 * [98447feb6](https://github.com/kubeovn/kube-ovn/commit/98447feb6f63a78e4cb3bc4925978d1eeaa33af6) update ovn patches (#4143)
 * [9c3ff7d36](https://github.com/kubeovn/kube-ovn/commit/9c3ff7d368556699eb7ff0ee33ee18d4bb973207) fix u2o e2e in 1.9 (#4154)
 * [351ba96f7](https://github.com/kubeovn/kube-ovn/commit/351ba96f7901088d188c352ad5e85421a94eb11f) speaker: add option to set node IP to override pod IP (#4140)
 * [37b7c46ad](https://github.com/kubeovn/kube-ovn/commit/37b7c46ad5d83aa8a4ca94260ab227d3feb4ee0e) Drop u2o arp request (#4135)
 * [f5b1d2abc](https://github.com/kubeovn/kube-ovn/commit/f5b1d2abcda157f447a98d715902b2a0a6e544ec) add ovn0 default route (#4127)
 * [5ed1c9f66](https://github.com/kubeovn/kube-ovn/commit/5ed1c9f669f5c0494915df55ec9e136e68e4dcd6) build(deps): bump google.golang.org/protobuf from 1.34.1 to 1.34.2 (#4146)
 * [51d46cdf1](https://github.com/kubeovn/kube-ovn/commit/51d46cdf187676607d5f021c3e1e43bf4c6cfb21) ovn: add support for conditionally skipping conntrack (#4131)
 * [815bfc280](https://github.com/kubeovn/kube-ovn/commit/815bfc2807fe13109d43d782a4c5e3628bf831c4) Ipam ut (#4112)
 * [50b0b8252](https://github.com/kubeovn/kube-ovn/commit/50b0b8252de6e181dc0567a53bcfb0064975b87a) build(deps): bump github.com/Microsoft/hcsshim from 0.12.3 to 0.12.4 (#4141)
 * [7220a5ee2](https://github.com/kubeovn/kube-ovn/commit/7220a5ee29a84ccfcd84eb34f8fbd48f6f2cab87) allow ip change name in migration (#4133)
 * [d6b41942a](https://github.com/kubeovn/kube-ovn/commit/d6b41942a393aa4b88f8a14fbbeaafa69061d641) distinguish-portSecurity-with-security-group (#4134)
 * [442cb0d5b](https://github.com/kubeovn/kube-ovn/commit/442cb0d5b91f8d40683b75736736e8740daea9c3) build(deps): bump github.com/docker/docker (#4136)
 * [5253c69e0](https://github.com/kubeovn/kube-ovn/commit/5253c69e0830b2eb64aa81b67dc629b84fa2c6dc) build(deps): bump sigs.k8s.io/controller-runtime from 0.18.3 to 0.18.4 (#4137)
 * [426b5d0a6](https://github.com/kubeovn/kube-ovn/commit/426b5d0a633e5a5e2c15492b04f5b059f1921515) fix ipam subnet concurrent map iteration and map write (#4126)
 * [ab8665191](https://github.com/kubeovn/kube-ovn/commit/ab8665191aa51aa4197985501bfd257729e637d0) trivy: ignore unfixed CVEs (#4129)
 * [9650f03b5](https://github.com/kubeovn/kube-ovn/commit/9650f03b5bdf04964aa6a391fef0a2b42f89a181) bump ovn to v24.03 - part 1 (#4083)
 * [7acdbab5a](https://github.com/kubeovn/kube-ovn/commit/7acdbab5a6e1681e90df5a0a29a749085de8c559) bump go to 1.22.4 (#4121)
 * [120bc8346](https://github.com/kubeovn/kube-ovn/commit/120bc83466fc29fe5f9ad41e9f3b2fe52953d7f1) build(deps): bump golang.org/x/sys from 0.20.0 to 0.21.0 (#4118)
 * [a222b682b](https://github.com/kubeovn/kube-ovn/commit/a222b682b818e6f94c2815d5d4d81deff26ed221) Enable inactivity check on ovndb connection (#4006)
 * [51880ef57](https://github.com/kubeovn/kube-ovn/commit/51880ef57ae87c8dde0b010ff68128cdf81b855c) build(deps): bump golang.org/x/mod from 0.17.0 to 0.18.0 (#4120)
 * [9b309190a](https://github.com/kubeovn/kube-ovn/commit/9b309190a1751a075ac74c515ebe77c51ea9286b) build(deps): bump github.com/osrg/gobgp/v3 from 3.26.0 to 3.27.0 (#4119)
 * [78f268fee](https://github.com/kubeovn/kube-ovn/commit/78f268fee444914d5e02f54cb818ad3d266f6efb) fix e2e timeout (#4116)
 * [96190f4cd](https://github.com/kubeovn/kube-ovn/commit/96190f4cd02fb6abed6b4201806d637ef2f6c57b) fix: IP Add/Sub overflow or underflow (#4111)
 * [bc0308e40](https://github.com/kubeovn/kube-ovn/commit/bc0308e40338ec43f3cab6598d9d0d98bf208e6a) delete unused dpdk dockerfile (#4113)
 * [cf7c98262](https://github.com/kubeovn/kube-ovn/commit/cf7c98262f383afe6756e9af10dfd09b7208e151) ovn vpc nat gw e2e only run after 1.12 (#4110)
 * [d62aace22](https://github.com/kubeovn/kube-ovn/commit/d62aace2284c1d8680fcd6b8c7cafa9daf9070fe) u2o: fix ip cr's mac field does not match with lrp mac (#4107)
 * [db821d221](https://github.com/kubeovn/kube-ovn/commit/db821d2210f3921af250897c930147d1c672d779) Makefile: clean generate ssl files and ko logs (#4106)
 * [d5f90bf4a](https://github.com/kubeovn/kube-ovn/commit/d5f90bf4a76e302c79ea90421d6ce5cb2c567aa5) add unit test for util (#4060)
 * [e67d5f12b](https://github.com/kubeovn/kube-ovn/commit/e67d5f12b750ca4c04b632c3dabd94b047d1f1e5) add reinit process for lb-svc pod (#4092)
 * [bcc96b141](https://github.com/kubeovn/kube-ovn/commit/bcc96b1413a39e64dc16a4e69f5eb1f0eb54b5cb) build(deps): bump github.com/emicklei/go-restful/v3 (#4103)
 * [55dbfbb8a](https://github.com/kubeovn/kube-ovn/commit/55dbfbb8acbef4b3b7992ac97b8d1251559c2f42) Fix change subnet (#4088)
 * [5fe053eb8](https://github.com/kubeovn/kube-ovn/commit/5fe053eb8b434622fe395206366309fb99c4c3d3) e2e: remove unused variables (#4098)
 * [1947da9e0](https://github.com/kubeovn/kube-ovn/commit/1947da9e09e5ea46c706342f2e812c762c5f02d7) Makefile: run kubectl-ko script when collecting logs (#4100)
 * [8270d802c](https://github.com/kubeovn/kube-ovn/commit/8270d802c362db4bfd86a97916ee524467227e57) e2e: bump images (#4101)
 * [2eab0593a](https://github.com/kubeovn/kube-ovn/commit/2eab0593ace8765957dd7480bc1da7428ed7bb40) install.sh: waiting for deleted kube-ovn-pinger to disapper (#4096)
 * [368d19043](https://github.com/kubeovn/kube-ovn/commit/368d190430532a66110e0ae0c9fe130b9fb7a973) fix mac conflict (#4095)
 * [9c981274b](https://github.com/kubeovn/kube-ovn/commit/9c981274b50fe1221a9ddc5bc7eac094f6807795) fix: add ip_reserved label for vip (#4093)
 * [202923114](https://github.com/kubeovn/kube-ovn/commit/2029231144785b9dfbdde0ace9ba981901d48625) build(deps): bump github.com/onsi/ginkgo/v2 from 2.18.0 to 2.19.0 (#4085)
 * [afe5446c4](https://github.com/kubeovn/kube-ovn/commit/afe5446c42d03cf3da3483bd60e988b3b43b2f40) build(deps): bump sigs.k8s.io/controller-runtime from 0.18.2 to 0.18.3 (#4084)
 * [e1310e170](https://github.com/kubeovn/kube-ovn/commit/e1310e1705096d80e7ea78d8d47fae6c728e56d8) install.sh: wait for all kube-ovn-pinger pods to be ready (#4082)
 * [13564d9be](https://github.com/kubeovn/kube-ovn/commit/13564d9bebcb106f09ab41ad4ba8cb9f235bd688) ci: print all the previous logs for restarted pods (#4081)
 * [9f67d358c](https://github.com/kubeovn/kube-ovn/commit/9f67d358caa98a09f4daaf75734870e040e2f86f) opt: replace ovn-sbctl with ovsdb-client (#4075)
 * [aebe363ee](https://github.com/kubeovn/kube-ovn/commit/aebe363eea382d351cd37d6dccf16afdc18bc21a) ovs: get controllerrevision with option --ignore-not-found (#4058)
 * [7148bcc36](https://github.com/kubeovn/kube-ovn/commit/7148bcc36b3d91555d5441813e9979ba73c35f38) fix exit on error (#4080)
 * [4bf54e35e](https://github.com/kubeovn/kube-ovn/commit/4bf54e35e38eea37e94b4e84352caf3be92d9ecd) delete lease on cleanup (#4079)
 * [f1daf5c86](https://github.com/kubeovn/kube-ovn/commit/f1daf5c86b7a7214dd8e0267d887d7edb7475b02) base: aoivd unnecessary env variables (#4070)
 * [796d5f089](https://github.com/kubeovn/kube-ovn/commit/796d5f089347243db21e30243a7ac032c3fbe132) ci: downgrade node image to v1.29.2 (#4069)
 * [23f963200](https://github.com/kubeovn/kube-ovn/commit/23f9632002afac8a4c2abb0819aed49f96691bf8) fix crypto/rand: argument to Int is <= 0 (#4077)
 * [e31784afa](https://github.com/kubeovn/kube-ovn/commit/e31784afa1c099aa2d665f1494fcbb0467a0a652) docs: updated CHANGELOG.md (#4072)
 * [2a3cdbd36](https://github.com/kubeovn/kube-ovn/commit/2a3cdbd361a1ecf1442d2659f1fd5abdff39c3cd) fix: gateway should not be network address and broadcast address (#4043)
 * [08e3c56ea](https://github.com/kubeovn/kube-ovn/commit/08e3c56ead51540c440bf9383be4bb2e7cafba51) support ctr to generate ssl certs (#4068)
 * [66f61445a](https://github.com/kubeovn/kube-ovn/commit/66f61445aeb0eb10a8e384d4a9d2f9eb4b2d510d) e2e: skip SCTP connectivity tests for versions prior to 1.12 (#4065)
 * [d76545b8d](https://github.com/kubeovn/kube-ovn/commit/d76545b8d26280f8137c441cddec857c146740c7) fix: should update subnet status after change vm subnet (#4061)
 * [3a7ee3431](https://github.com/kubeovn/kube-ovn/commit/3a7ee3431627b8945c3067a240da2df39bcf672e) bump ginkgo to v2.18.0 (#4062)
 * [1903bd8ce](https://github.com/kubeovn/kube-ovn/commit/1903bd8cef57a7b1cdea395ed43335b5a747fc8e) ci: fix scheduled e2e (#4057)
 * [72180104d](https://github.com/kubeovn/kube-ovn/commit/72180104df098b90e4778da1bec44bd625f459b5) ci: build e2e binaries and free disk space on necessary (#4059)
 * [f00ff138f](https://github.com/kubeovn/kube-ovn/commit/f00ff138fab75c466555fc303ba95704de0ace4e) crd: add subnet name pattern (#4054)
 * [bf1cb76cd](https://github.com/kubeovn/kube-ovn/commit/bf1cb76cd9dc4e592b5513d3c14e1832a8c32c8d) optimize code (#4049)
 * [d30c1e7be](https://github.com/kubeovn/kube-ovn/commit/d30c1e7be5af8b0b0fcba0ed86247d60868b501e) add err log (#4046)
 * [f98308791](https://github.com/kubeovn/kube-ovn/commit/f98308791273f8dc2d0a798c8e01ed4f566f68f7) build(deps): bump github.com/containernetworking/plugins (#4055)
 * [94ce0acc5](https://github.com/kubeovn/kube-ovn/commit/94ce0acc5f536abe995efd47b8d6d4b6f7064c2c) docs: updated CHANGELOG.md (#4050)
 * [bb46f5713](https://github.com/kubeovn/kube-ovn/commit/bb46f5713fbcd66fa16fe772d8d2df775b7ff3b7) wait for all pods to be deleted before deleting serviceaccount/clusterrole/clusterrolebinding (#4035)
 * [96cd668e8](https://github.com/kubeovn/kube-ovn/commit/96cd668e8d425ee6b54e97c430fcab748a455358) fix: add route for underlay subnet which enables u2o and disables LB (#4039)
 * [e907ec51e](https://github.com/kubeovn/kube-ovn/commit/e907ec51e249a799dfe144177019671e1b69cb76) add pure arm64 build target (#4044)
 * [a0c2a4344](https://github.com/kubeovn/kube-ovn/commit/a0c2a4344b8f211c502dc1e1c161d94bb778a51b) log deleting iptables rule (#4031)
 * [87c007413](https://github.com/kubeovn/kube-ovn/commit/87c007413dc17d6bd7af2568090b9b37d968eaad) build(deps): bump github.com/docker/docker (#4037)
 * [ef3cef030](https://github.com/kubeovn/kube-ovn/commit/ef3cef03058610caf393a3d4f774ec0e0f04ff66) uninstall.sh: delete OVN-POSTROUTING rule in mangle table (#4034)
 * [4995dcc05](https://github.com/kubeovn/kube-ovn/commit/4995dcc05d9123378f1c7555c6350a8ea819e751) cleanup.sh: remove sa/clusterrole/clusterrolebinding (#4024)
 * [e0aa56628](https://github.com/kubeovn/kube-ovn/commit/e0aa566287fb0353871fce0a85f79edcc51582e2) do not use exec for start scripts with trap quit EXIT (#4025)
 * [3015e5dab](https://github.com/kubeovn/kube-ovn/commit/3015e5dab42bdf6d4fbabfcec2e5c95ca2c034be) bump k8s to 1.30.1 (#4028)
 * [2f92a6fc6](https://github.com/kubeovn/kube-ovn/commit/2f92a6fc6fe4542e31ea904a5c9a966567193e9f) fix node annotations not updated when initializing the default provider-network (#4030)
 * [2245d6404](https://github.com/kubeovn/kube-ovn/commit/2245d6404dcf2762d10b4f4489fcc50148eab8a6) set mac in U2OInterconnection ip resources (#4008)
 * [fa66889b1](https://github.com/kubeovn/kube-ovn/commit/fa66889b1154076c466933c9f0f2ed313c7fbf62) rollback residual link and port (#4012)
 * [51602bd8c](https://github.com/kubeovn/kube-ovn/commit/51602bd8cc6695ef3392bd56cca4f287536ac879) build(deps): bump google.golang.org/grpc from 1.63.2 to 1.64.0 (#4027)
 * [a1815db05](https://github.com/kubeovn/kube-ovn/commit/a1815db0537482a53e9d430467b8f5a5cd8cba37) fix add ip eip trigger subnet status count ip (#4023)
 * [85e28bc1e](https://github.com/kubeovn/kube-ovn/commit/85e28bc1e6e2b40de5b5aff6df951f4111783cc0) bump gosec to 2.20.0 (#4021)
 * [9236508a8](https://github.com/kubeovn/kube-ovn/commit/9236508a8ea145a0f788c71328d468ace080772f) remove unused yamls (#4022)
 * [178299a32](https://github.com/kubeovn/kube-ovn/commit/178299a3248f935c5a7683ab98d6390b8536e874) reconcile iptable to the first position after iptables -t nat -F (#3995)
 * [6924c4e6a](https://github.com/kubeovn/kube-ovn/commit/6924c4e6aac50f3f06f5b01b4d4956a7af82d8b3) fix container args (#4020)
 * [6b048eb49](https://github.com/kubeovn/kube-ovn/commit/6b048eb497bdb3b7692874bee3005899d2a3f88f) Append new DHCP options to existing options (#3997)
 * [18a477418](https://github.com/kubeovn/kube-ovn/commit/18a477418b7e57548078beb69aeac5cbbf54d8c1) Make dhcpOptions can accept multiple addresses (#4003)
 * [5994860b6](https://github.com/kubeovn/kube-ovn/commit/5994860b61742c53b72e5bfe16e6a0c743df600d) ci: bump k8s to v1.30.0 (#4019)
 * [b8d77e242](https://github.com/kubeovn/kube-ovn/commit/b8d77e24222630f73af7403bb09bbfd59a6f54c1) fix lsp not updated correctly when logical switch is changed (#4015)
 * [a57d04115](https://github.com/kubeovn/kube-ovn/commit/a57d0411566c517f2031ebb04c9facd9a353519d) base: set entrypoint to  dumb-init (#4018)
 * [67848da17](https://github.com/kubeovn/kube-ovn/commit/67848da173514561221325b95146c0ed26c6f718) fix: u2o dual stack uses the corresponding svc ip (#4013)
 * [35e37ad7e](https://github.com/kubeovn/kube-ovn/commit/35e37ad7eef6ef9ef2f2bba4113ef84ab927a98c) fix policy route not deleted after subnet is deleted (#4016)
 * [4d7ec749a](https://github.com/kubeovn/kube-ovn/commit/4d7ec749a6e7878fcd84a1b0dc9340aa6945d281) fix: Resolved the hidden issue with zombie processes (#4004)
 * [4d8c39199](https://github.com/kubeovn/kube-ovn/commit/4d8c39199cb1ebe7d59500f2a611bc9db1094593) docs: updated CHANGELOG.md (#4017)
 * [64452939c](https://github.com/kubeovn/kube-ovn/commit/64452939c9f4e56281c5db3bb31a1ad533b06108) fix lsp not updating addresses (#4011)
 * [41cb044c9](https://github.com/kubeovn/kube-ovn/commit/41cb044c9c15de6bc3919f03e6a358eda822c085) simplify file reading (#4010)
 * [aaf5736a8](https://github.com/kubeovn/kube-ovn/commit/aaf5736a8a3de45bc24861da07622a4044d46fb4) fix: close file (#4007)
 * [9d19242c6](https://github.com/kubeovn/kube-ovn/commit/9d19242c6e9fce48e2bd1a45b48db4a8573a95ce) fix controller-runtime logger not set (#4005)
 * [3b73df6e8](https://github.com/kubeovn/kube-ovn/commit/3b73df6e80e2b345d418768c272eb9678e979b06) update meeting info
 * [e2dee3e42](https://github.com/kubeovn/kube-ovn/commit/e2dee3e4278ebe87d600874e10bd7626e1e27f50) fix node gc (#3992)
 * [060231bf7](https://github.com/kubeovn/kube-ovn/commit/060231bf750b2bd686e77b213ee544cd5ceb1b66) build(deps): bump github.com/docker/docker (#4000)
 * [dd4de4ac1](https://github.com/kubeovn/kube-ovn/commit/dd4de4ac1b3ea6a00afa64c16700d2c061aa8b43) e2e: fix skip cases for k8s network conformance tests (#3994)
 * [5bcdf2878](https://github.com/kubeovn/kube-ovn/commit/5bcdf2878ee6e95c8a04d2a3adad1f90cb29d152) Disable pod use NodeSwitch subnet (#3991)
 * [8c548754c](https://github.com/kubeovn/kube-ovn/commit/8c548754c92d9f088ac34b621aa650652133df07) update roadmap
 * [0068ce3e7](https://github.com/kubeovn/kube-ovn/commit/0068ce3e7cf0d073dff85a1253e4542e43678f27) bump go to 1.22.3 (#3989)
 * [b64f9fe59](https://github.com/kubeovn/kube-ovn/commit/b64f9fe591817c1db46bc2fd690e52b2c3646435) subnet print vlan as vpc (#3983)
 * [5f6483871](https://github.com/kubeovn/kube-ovn/commit/5f64838711ec0c537e18db802693ce2328609362) Revert "build(deps): bump github.com/kubeovn/felix (#3984)" (#3988)
 * [e3c729cd5](https://github.com/kubeovn/kube-ovn/commit/e3c729cd5cdaba1e028e72f3e6885cf0bbed3d64) build(deps): bump github.com/onsi/ginkgo/v2 from 2.17.2 to 2.17.3 (#3985)
 * [3c82e51c9](https://github.com/kubeovn/kube-ovn/commit/3c82e51c9e193e048d7837942b1507a74715301e) build(deps): bump golangci/golangci-lint-action from 5 to 6 (#3986)
 * [33bda4f28](https://github.com/kubeovn/kube-ovn/commit/33bda4f286516cf15b4286447f4b91ade032db5d) build(deps): bump github.com/kubeovn/felix (#3984)
 * [63d0aade6](https://github.com/kubeovn/kube-ovn/commit/63d0aade6c9f201a6f3efec71b844e84334f8ec8) bump k8s to v1.30.0 (#3965)
 * [776bd8373](https://github.com/kubeovn/kube-ovn/commit/776bd83733ac1d84cf221fa29a012c491b2c5a40) remove vendor name info (#3978)
 * [81b40c30e](https://github.com/kubeovn/kube-ovn/commit/81b40c30ef51f14871b44c453c26e5437b60c3c5) support ovn eip nats dualstack (#3822)
 * [02f876fb9](https://github.com/kubeovn/kube-ovn/commit/02f876fb986d8f5ef47239c71572b8f81eb38aa4) docs: updated CHANGELOG.md (#3982)
 * [828402dff](https://github.com/kubeovn/kube-ovn/commit/828402dff10cffe24446fc774c964b8f4fe1bcb9) build(deps): bump google.golang.org/protobuf from 1.34.0 to 1.34.1 (#3981)
 * [6656a1fed](https://github.com/kubeovn/kube-ovn/commit/6656a1fed4cf5a671d843ef4ef8f31faba0bb6a6) build(deps): bump golang.org/x/sys from 0.19.0 to 0.20.0 (#3980)
 * [5fa123a44](https://github.com/kubeovn/kube-ovn/commit/5fa123a4482818f9dc751bd96baddcf28a8e041c) ipam: fix IPRangeList clone (#3979)
 * [8a38cfab8](https://github.com/kubeovn/kube-ovn/commit/8a38cfab81d12ee7a6596694c755737374fcebd6) build(deps): bump github.com/docker/docker (#3974)
 * [ef799c77b](https://github.com/kubeovn/kube-ovn/commit/ef799c77b4fae4d851b0bf233e9f32fdb6e4e665) chart: fix kubeVersion to allow for patterns that match sub versions (#3975)
 * [afdedc6a7](https://github.com/kubeovn/kube-ovn/commit/afdedc6a72102b10c7a0f6769aebe182fc13dea7) build(deps): bump github.com/osrg/gobgp/v3 from 3.25.0 to 3.26.0 (#3973)
 * [1ea18d8c9](https://github.com/kubeovn/kube-ovn/commit/1ea18d8c9da4ea03047dd9211da95713a97d4428) build(deps): bump github.com/onsi/gomega from 1.33.0 to 1.33.1 (#3972)
 * [b50a7ec26](https://github.com/kubeovn/kube-ovn/commit/b50a7ec26ec21c6995e75cb6a567e46f628e7952) build(deps): bump google.golang.org/protobuf from 1.33.0 to 1.34.0 (#3971)
 * [50b101b08](https://github.com/kubeovn/kube-ovn/commit/50b101b08b5bbc3ea1361cea78eda9db0c5199c1) docs: updated CHANGELOG.md (#3969)
 * [2e8556d50](https://github.com/kubeovn/kube-ovn/commit/2e8556d5050f3a84e1b608e6cd73035ff8480625) change subnet check using ip location to avoid ipam init err (#3962)
 * [fa2770bf2](https://github.com/kubeovn/kube-ovn/commit/fa2770bf2f136d06d45a2620fe727811376e2fd3) fix subnet acl with same net allow (#3961)
 * [dd9a25dbf](https://github.com/kubeovn/kube-ovn/commit/dd9a25dbf9cf4fc5f34a60c112ba42aebc6ccf2c) add patch permission for cni clusterrole (#3959)
 * [18d1a1558](https://github.com/kubeovn/kube-ovn/commit/18d1a1558080a8a23c48cd5a53d5e8f722152015) use jinjanator instead of j2cli (#3957)
 * [7cddf7f7d](https://github.com/kubeovn/kube-ovn/commit/7cddf7f7d08458e0b83de3184acf5171c70f0ee6) fix index out of range (#3958)
 * [4c5061961](https://github.com/kubeovn/kube-ovn/commit/4c50619616e5e51d613cb01f3cb0388ae94ed35c) Fix finalizer (#3955)
 * [77f011c09](https://github.com/kubeovn/kube-ovn/commit/77f011c0968f080bd5458aec6fc536fbd9050855) build(deps): bump golangci/golangci-lint-action from 4 to 5 (#3952)
 * [e56a6367d](https://github.com/kubeovn/kube-ovn/commit/e56a6367d7ac7178660fd342b30c5d4ac452441b) fix nil pointer dereference (#3951)
 * [7e0afbd72](https://github.com/kubeovn/kube-ovn/commit/7e0afbd729caa3587b768854bc04e2ac43f79c9e) build(deps): bump helm/kind-action from 1.9.0 to 1.10.0 (#3949)
 * [c0a08ab3a](https://github.com/kubeovn/kube-ovn/commit/c0a08ab3a8e47bc3a49b648ed4a7c67713a6696b) build(deps): bump github.com/docker/docker (#3948)
 * [609559bb6](https://github.com/kubeovn/kube-ovn/commit/609559bb6f901dd837c55a845f3dcc1fcc2285f0) rename templates
 * [ff83a0e16](https://github.com/kubeovn/kube-ovn/commit/ff83a0e16236ca0ac60d2243896632d409d941e4) add issue template config
 * [f99cd2b36](https://github.com/kubeovn/kube-ovn/commit/f99cd2b3695fc6c3088d285d96142ea35fba30ea) modify issue templates
 * [a44335253](https://github.com/kubeovn/kube-ovn/commit/a44335253ff4caba20c30b17cea4724374045e54) check kube ovn podrestarts (#3925)
 * [015564d34](https://github.com/kubeovn/kube-ovn/commit/015564d3499ad1422a6e24606eca147d827810c3) docs: updated CHANGELOG.md (#3946)
 * [152d54da3](https://github.com/kubeovn/kube-ovn/commit/152d54da35f8a5e161cc2920dda2d5d6e018f63f) fix: lower camel case (#3942)
 * [740e01f4b](https://github.com/kubeovn/kube-ovn/commit/740e01f4bfe94eda561446b560b6887f1159b50c) build(deps): bump github.com/Microsoft/hcsshim from 0.12.2 to 0.12.3 (#3944)
 * [cd9df5d26](https://github.com/kubeovn/kube-ovn/commit/cd9df5d26352e4f335dd51f4e991300359f2f3fd) drop both IPv4 and IPv6 traffic in networkpolicy drop acl (#3940)
 * [cb6f22992](https://github.com/kubeovn/kube-ovn/commit/cb6f22992a544cc964b1ebe5d3424ca15efa954a) build(deps): bump github.com/onsi/gomega from 1.32.0 to 1.33.0 (#3938)
 * [7ffb791ae](https://github.com/kubeovn/kube-ovn/commit/7ffb791aeeb6dd443047bd8fca16f5b960021a48) add metric for subnet info (#3932)
 * [ce5c01704](https://github.com/kubeovn/kube-ovn/commit/ce5c01704e15edf48b4ab8c8c4a73168e7558973) build(deps): bump github.com/docker/docker (#3935)
 * [05b2d5b58](https://github.com/kubeovn/kube-ovn/commit/05b2d5b58c68693cbeb73a4ccbcb8a7e281162b1) cni-server: set sysctl variables only when the env variables are passed in (#3929)
 * [6b92ab638](https://github.com/kubeovn/kube-ovn/commit/6b92ab638d3c27fb5998a51a62d473b0c48a41bb) ovn: check whether db file is fixed (#3928)
 * [b1d87f172](https://github.com/kubeovn/kube-ovn/commit/b1d87f172b8ad216fae6720fa5f961c380354d12) bump k8s to v1.29.4 (#3927)
 * [a5bf8612e](https://github.com/kubeovn/kube-ovn/commit/a5bf8612e28d110c981f7844ece5df9349e90986) append filepath when compare cni config (#3923)
 * [f6534cf24](https://github.com/kubeovn/kube-ovn/commit/f6534cf2419324cbb1c8732c4ec6e06e9b1243b6) build(deps): bump azure/setup-helm from 4.1.0 to 4.2.0 (#3926)
 * [1d693287b](https://github.com/kubeovn/kube-ovn/commit/1d693287b90133387844a804370092939e35f757) cni: add mtu in add result (#3922)
 * [8cc1a0f73](https://github.com/kubeovn/kube-ovn/commit/8cc1a0f730beadfea0db4fc8f6429596cfb1ed65) build(deps): bump kubevirt.io/api from 1.1.1 to 1.2.0 (#3918)
 * [23c8df489](https://github.com/kubeovn/kube-ovn/commit/23c8df489ebf133f2a6dcc56706db5ea802a39e6) re calculate subnet using ips while inconsistency detected (#3920)
 * [f536d4805](https://github.com/kubeovn/kube-ovn/commit/f536d4805939d879bd3d6b16c951fba9c2653323) update ovn monitor (#3903)
 * [f81256213](https://github.com/kubeovn/kube-ovn/commit/f81256213594c62e54c254f6c64896d3af9bdccd) add monitor for sysctl para (#3913)
 * [98580d68e](https://github.com/kubeovn/kube-ovn/commit/98580d68e876977623ebdfa557c939ac3f1ef008) update form link
 * [44466906d](https://github.com/kubeovn/kube-ovn/commit/44466906dd917cefc9a889579a5fa1111b5a0dbb) build(deps): bump github.com/containernetworking/cni (#3917)
 * [fe157d3df](https://github.com/kubeovn/kube-ovn/commit/fe157d3df185c693304402f9176ccee764b01405) refactor kubevirt vm e2e (#3914)
 * [062a55881](https://github.com/kubeovn/kube-ovn/commit/062a5588162ed0a3938c49c90382c96f54ca785e) fix subnet provider validation (#3915)
 * [748aac73a](https://github.com/kubeovn/kube-ovn/commit/748aac73a9ef231e8218c7ff764081171f36ebae) bump go toolchain to go1.22.2 (#3905)
 * [f5aef5a7d](https://github.com/kubeovn/kube-ovn/commit/f5aef5a7d1547154b57e0b1caad50e03281901cd) build(deps): bump github.com/k8snetworkplumbingwg/network-attachment-definition-client from 1.6.0 to 1.7.0 (#3906)
 * [52fed946b](https://github.com/kubeovn/kube-ovn/commit/52fed946bc26733747ac7f7d8107fab740f909c3) build(deps): bump github.com/docker/docker (#3907)
 * [85ac74124](https://github.com/kubeovn/kube-ovn/commit/85ac74124e889b571187bce0c358ba2ddc32a9cf) support specifying routes when providing IPAM for other CNI plugins (#3904)
 * [e0fffa6b5](https://github.com/kubeovn/kube-ovn/commit/e0fffa6b53b3d61b384700b50d0faa5ef1a7cebb) fix br-external not init because of  no permission after ovn-nat-gw configmap created (#3902)
 * [a6e9ef461](https://github.com/kubeovn/kube-ovn/commit/a6e9ef461b3a353294d3d6ab549cafa7f1ea2d22) build(deps): bump sigs.k8s.io/controller-runtime from 0.17.2 to 0.17.3 (#3901)
 * [629c9df59](https://github.com/kubeovn/kube-ovn/commit/629c9df59ec4df94108d99a58ea59166b870433a) build(deps): bump google.golang.org/grpc from 1.63.0 to 1.63.2 (#3900)
 * [859e109b5](https://github.com/kubeovn/kube-ovn/commit/859e109b522b125c45dd7b0b1f6b1cb9be3574eb) Fix init sg (#3890)
 * [2fa5df222](https://github.com/kubeovn/kube-ovn/commit/2fa5df222b927fcc6f0a6a1319818db5cd364e30) distinguish portSecurity with security group (#3862)
 * [28aa65d18](https://github.com/kubeovn/kube-ovn/commit/28aa65d18fb1e9fa44d1a4faee1eca130110d7ce) skip conntrack when access node dns ip (#3894)
 * [019bacd36](https://github.com/kubeovn/kube-ovn/commit/019bacd361e0d54ee65acf025328b68a46a75236) build(deps): bump github.com/osrg/gobgp/v3 from 3.24.0 to 3.25.0 (#3897)
 * [f983f5b49](https://github.com/kubeovn/kube-ovn/commit/f983f5b497ffc2a3476242c9340ef6a9f696c262) build(deps): bump golang.org/x/sys from 0.18.0 to 0.19.0 (#3896)
 * [b292e01f8](https://github.com/kubeovn/kube-ovn/commit/b292e01f861530ebde0bf875bf9c014ade227613) build(deps): bump golang.org/x/mod from 0.16.0 to 0.17.0 (#3895)
 * [9207cbab2](https://github.com/kubeovn/kube-ovn/commit/9207cbab29f0dcd62cd41a9bb66f0655ac3b8484) build(deps): bump google.golang.org/grpc from 1.62.1 to 1.63.0 (#3898)
 * [b7a2c0a83](https://github.com/kubeovn/kube-ovn/commit/b7a2c0a83d16ab4885a720faaaf96f7046637d2f) chart: fix ovs-ovn update strategy (#3887)
 * [a48c597eb](https://github.com/kubeovn/kube-ovn/commit/a48c597eb610e6d499aa56ca65abaf070f38962d) fix ic e2e (#3831)
 * [f8dee83c1](https://github.com/kubeovn/kube-ovn/commit/f8dee83c1af8adbacd2f18a93a7b713dbf6db3ec) build(deps): bump github.com/Microsoft/hcsshim from 0.12.1 to 0.12.2 (#3891)
 * [64f6f2ba6](https://github.com/kubeovn/kube-ovn/commit/64f6f2ba6cf670e56c124af0bb931f6be6feeb64) Fix index out of range in controller security_group (#3845)
 * [8d17662fe](https://github.com/kubeovn/kube-ovn/commit/8d17662fe74e137c0135d3bafc5b56e5cbfa04a3) fix: ipam invalid memory address or nil pointer dereference (#3889)
 * [d954b6358](https://github.com/kubeovn/kube-ovn/commit/d954b63580e7e71599cfe740ed7857b1a960cac6) add tracepath (#3884)
 * [e1a8a39bd](https://github.com/kubeovn/kube-ovn/commit/e1a8a39bd7e5c66bcc4a9de8f521e5ec3a078cd4) change northd probe interval to 5s (#3885)
 * [31968e810](https://github.com/kubeovn/kube-ovn/commit/31968e8100f1364a4928b210bf9d0bb6573eba18) docs: updated CHANGELOG.md (#3883)
 * [57f08ad94](https://github.com/kubeovn/kube-ovn/commit/57f08ad94f662cfe01bb21812cc3662fc4cbbf51) ovn: reduce down time during upgrading from version 21.06 (#3881)
 * [3328b9012](https://github.com/kubeovn/kube-ovn/commit/3328b901273d3103cacf0b5049ba46d42d08b941) compile binaries with debug symbols for debug images (#3871)
 * [b9116899c](https://github.com/kubeovn/kube-ovn/commit/b9116899c47c13a65451e7c7476551b0db6892ce) ovn: update patch for skipping ct (#3879)
 * [74fafdb62](https://github.com/kubeovn/kube-ovn/commit/74fafdb624987afaf29d2a76deb212b8933a6c14) docs: updated CHANGELOG.md (#3877)
 * [f1bbab62a](https://github.com/kubeovn/kube-ovn/commit/f1bbab62a8b80692c3a04fc6f37ede7331800ce5) build(deps): bump github.com/cenkalti/backoff/v4 from 4.2.1 to 4.3.0 (#3875)
 * [40614b5a3](https://github.com/kubeovn/kube-ovn/commit/40614b5a3da72a748b622a72e2aea27c8cf56775) ci: fix memory leak reporting caused by ovn-controller crashes (#3873)
 * [982b13688](https://github.com/kubeovn/kube-ovn/commit/982b13688bfbca9fc4e46c5ba9d5333499b01453) fix reference to nil pointer (#3872)
 * [1173594e2](https://github.com/kubeovn/kube-ovn/commit/1173594e2887bc66cffa7eeceb6d30a55d619a78) build(deps): bump github.com/onsi/ginkgo/v2 from 2.17.0 to 2.17.1 (#3866)
 * [485d07f6f](https://github.com/kubeovn/kube-ovn/commit/485d07f6f5fc8b4137ad3fecdfe84434e6c5e265) build(deps): bump github.com/Microsoft/hcsshim from 0.12.0 to 0.12.1 (#3867)
 * [b48c519f7](https://github.com/kubeovn/kube-ovn/commit/b48c519f73a839713c92b8d2bbe01a8e8fd18743) fix: when give ipv4 cidr ipv6 gateway, the gateway will extend infinitely (#3865)
 * [fbec723a4](https://github.com/kubeovn/kube-ovn/commit/fbec723a4c989540332cd1530a583f7ff7afced4) docs: updated CHANGELOG.md (#3864)
 * [54c95e484](https://github.com/kubeovn/kube-ovn/commit/54c95e484d6852748b014d0f158135f4a2bc5ffa) exclude vip as encap ip (#3859)
 * [2f48cc7f5](https://github.com/kubeovn/kube-ovn/commit/2f48cc7f58794f0bf63ba4f41b004b5cf201ba47) docs: updated CHANGELOG.md (#3858)
 * [05c67d87c](https://github.com/kubeovn/kube-ovn/commit/05c67d87cf74535387607c256285f74508ed5bdf) fix: aap get wrong lsp when using multus (#3827)
 * [cef6f0e92](https://github.com/kubeovn/kube-ovn/commit/cef6f0e92dd957000391a1f2eac0590cf3454df2) chart: fix missing ENABLE_IC (#3851)
 * [a4d5c06b9](https://github.com/kubeovn/kube-ovn/commit/a4d5c06b93394a570e1deefd18960225d234a250) delete unused mount volumes (#3838)
 * [d3e477f44](https://github.com/kubeovn/kube-ovn/commit/d3e477f44a8815690215e9f3303d6d6987deae41) bump k8s to v1.29.3 (#3833)
 * [07ac371a2](https://github.com/kubeovn/kube-ovn/commit/07ac371a2239eb91ba7db3963aae07d0c370e120) iptables: reject access to invalid service port only for TCP (#3843)
 * [eac020818](https://github.com/kubeovn/kube-ovn/commit/eac020818f998e7b6ca4cd7ecadfde7fe63677d0) ci: bump helm to v3.14.3 (#3829)
 * [d019ee273](https://github.com/kubeovn/kube-ovn/commit/d019ee27317faf368491a4440eaf600bf59a0cdd) build(deps): bump github.com/docker/docker (#3854)
 * [c5c121d53](https://github.com/kubeovn/kube-ovn/commit/c5c121d53e4905910ec773f1a8a67a162a48b59e) build(deps): bump github.com/docker/docker (#3849)
 * [b9810c80e](https://github.com/kubeovn/kube-ovn/commit/b9810c80eaecac182f6ea9167cbdbb46a89f4d09) u2o: add policy route for u2o when disabling lb (#3813)
 * [18d88d93e](https://github.com/kubeovn/kube-ovn/commit/18d88d93ec662a6274b47769b9b486ae45ab88c1) bump github.com/onsi/ginkgo/v2 to v2.17.0 (#3837)
 * [aaf1464f2](https://github.com/kubeovn/kube-ovn/commit/aaf1464f2e3b8508e96aff23a04413b693a1bb06) fix ip is a substring in node name (#3832)
 * [b28663b58](https://github.com/kubeovn/kube-ovn/commit/b28663b58cba265591d8de0e1eabb4fb495b6582) base: bump cni plugins and kubectl (#3826)
 * [b65941147](https://github.com/kubeovn/kube-ovn/commit/b65941147949a2af90fd1ec61f4fa83bf30d84ad) add ovs flow num dashboard (#3821)
 * [73e1c42d1](https://github.com/kubeovn/kube-ovn/commit/73e1c42d175db04c92d3c44dae800dbb7d03f0e6) build(deps): bump github.com/containernetworking/plugins (#3825)
 * [ad025ff1a](https://github.com/kubeovn/kube-ovn/commit/ad025ff1a663538882e7e85f7eb36aea62f35c59) build(deps): bump github.com/emicklei/go-restful/v3 (#3824)
 * [d29d9b36f](https://github.com/kubeovn/kube-ovn/commit/d29d9b36f1f97432f4d64b6936f7eb8f46717da3) build(deps): bump github.com/prometheus-community/pro-bing (#3820)
 * [be6cc86d9](https://github.com/kubeovn/kube-ovn/commit/be6cc86d9543f821e76bb823c2c2764fad567d13) fix cmd inject (#3809)
 * [a2977e47f](https://github.com/kubeovn/kube-ovn/commit/a2977e47fe74f2e4b14a7b06bf678baa98261c8c) kubectl-ko: fix subnet diagnose failure (#3808)
 * [f745776e1](https://github.com/kubeovn/kube-ovn/commit/f745776e1bceb55d9a3cd7e2771f14c0fe9fe1de) Makefile: fix clab bgp configuration (#3818)
 * [0314e82d6](https://github.com/kubeovn/kube-ovn/commit/0314e82d68db54b5f1ad24799d61cf7bbc257d2b) Makefile: fix cilium installation (#3817)
 * [a8dc9ab95](https://github.com/kubeovn/kube-ovn/commit/a8dc9ab953cb94a80bfeb1ba4f5c708f86e16301) add sts e2e case for pr (#3800)
 * [9e3fdabf3](https://github.com/kubeovn/kube-ovn/commit/9e3fdabf39eaed64a0e6771e635b643a94aef0ea) build(deps): bump github.com/docker/docker (#3811)
 * [61b6b5e40](https://github.com/kubeovn/kube-ovn/commit/61b6b5e403ddf8f9f2d700e1d08a836aab786a11) pinger: do not setup metrics if disabled (#3806)
 * [4f37bbcab](https://github.com/kubeovn/kube-ovn/commit/4f37bbcabbecff1370f031fba040be602b38e4a4) remove unused environment variables (#3804)
 * [4a9a1f821](https://github.com/kubeovn/kube-ovn/commit/4a9a1f8214223f273850bb53b19c82407de23621) Fix the failure to enable multi-network card traffic mirroring for newly created pods (#3805)
 * [a4c022809](https://github.com/kubeovn/kube-ovn/commit/a4c022809a2c53719ac553f021b1af174d4efada) build(deps): bump google.golang.org/protobuf from 1.32.0 to 1.33.0 (#3802)
 * [307c05dd7](https://github.com/kubeovn/kube-ovn/commit/307c05dd73b9bf1e29009d3a0a467e3f8949f154) build(deps): bump google.golang.org/grpc from 1.62.0 to 1.62.1 (#3803)
 * [b0678af05](https://github.com/kubeovn/kube-ovn/commit/b0678af053de33b39e703746fef4684a06c811c6) build(deps): bump github.com/onsi/ginkgo/v2 from 2.15.0 to 2.16.0 (#3797)
 * [7dda62d6f](https://github.com/kubeovn/kube-ovn/commit/7dda62d6fa8c763fbd382ca94bd72429519cf9fc) build(deps): bump github.com/Microsoft/hcsshim from 0.11.4 to 0.12.0 (#3796)
 * [116c24caa](https://github.com/kubeovn/kube-ovn/commit/116c24caace0e78620d78f80fd73055165cbc5c1) VM live migrate (#3767)
 * [5afcc3ec6](https://github.com/kubeovn/kube-ovn/commit/5afcc3ec69824ee6f5019dcb619fd6ac2e42aebb) refactor fastpath (#3794)
 * [7b937cb6b](https://github.com/kubeovn/kube-ovn/commit/7b937cb6b4f648836ad4b558236b9f50ca0adcc1) add missing metallb-cr.yaml.j2
 * [3ee81e23a](https://github.com/kubeovn/kube-ovn/commit/3ee81e23a3cf607751c4a28b3d6aa357408e0efd) Makefile: add target kind-install-metallb (#3795)
 * [e08e3a2c1](https://github.com/kubeovn/kube-ovn/commit/e08e3a2c1cafe9a7f21132c0c6281df988b0c9fb) fix: go toolchain config (#3790)
 * [31e0c8d51](https://github.com/kubeovn/kube-ovn/commit/31e0c8d51f220a638b48806f7b8b49513008aea6) build(deps): bump golang.org/x/sys from 0.17.0 to 0.18.0 (#3792)
 * [0951b36ae](https://github.com/kubeovn/kube-ovn/commit/0951b36aecea1cf412fd81358f7908dad3a63b77) build(deps): bump golang.org/x/mod from 0.15.0 to 0.16.0 (#3793)
 * [498a62de6](https://github.com/kubeovn/kube-ovn/commit/498a62de6ff8d73dfe191a5cb77a33f4f0a8c446) build(deps): bump azure/setup-helm from 4.0.0 to 4.1.0 (#3791)
 * [47ea57b72](https://github.com/kubeovn/kube-ovn/commit/47ea57b72666dd928aaa2850cdc110b11a0c86c7) Refactor build-go targets to use -trimpath flag (#3789)
 * [3e8329f7f](https://github.com/kubeovn/kube-ovn/commit/3e8329f7f61211993aa53f0cc973927df6707fd9) fix incorrect variable assignment (#3787)
 * [e264ed841](https://github.com/kubeovn/kube-ovn/commit/e264ed8414ea99dc76e198df23a201b0f9ac845e) Makefile: bump deepflow to v6.4 (#3783)
 * [b2a6b3e97](https://github.com/kubeovn/kube-ovn/commit/b2a6b3e977179751d03aea0a414fe50e928a8e8f) bump github.com/osrg/gobgp/v3 to v3.24.0 (#3785)
 * [25ebadc4e](https://github.com/kubeovn/kube-ovn/commit/25ebadc4e981536dcbd8a4cd4c0cebc463cfd113) build(deps): bump github.com/stretchr/testify from 1.8.4 to 1.9.0 (#3784)
 * [78234970e](https://github.com/kubeovn/kube-ovn/commit/78234970e8981501817aca3deb6e6b8da5c7afc9) docs: updated CHANGELOG.md (#3782)
 * [83ae77ff5](https://github.com/kubeovn/kube-ovn/commit/83ae77ff526d0891ce1cb3eb18d076fee3e02a5b) docs: updated CHANGELOG.md (#3779)
 * [da6b52213](https://github.com/kubeovn/kube-ovn/commit/da6b522133d61c4c28c77176856acc84c3b1a35f) fix sts pod's logical switch port do not update externa_id vendor and ls (#3777)
 * [aa6e36d0d](https://github.com/kubeovn/kube-ovn/commit/aa6e36d0d45a5bbb2d34e54f7c05baf73ea8a9a9) optimize nat gw (#3775)
 * [6037ee83a](https://github.com/kubeovn/kube-ovn/commit/6037ee83a0b4478e542097030f7289c84e4eca53) update files generated by k8s codegen (#3771)
 * [e85263fee](https://github.com/kubeovn/kube-ovn/commit/e85263fee0122d8f3ec25aba78875406e0e75d61) fix dot imports (#3773)
 * [ebb0f347c](https://github.com/kubeovn/kube-ovn/commit/ebb0f347c08bf3fe96811b0e0c9528b1a1de894c) fix crd (#3772)
 * [eeda6dcd4](https://github.com/kubeovn/kube-ovn/commit/eeda6dcd4b948ebec842932b7335e385faf94812) build(deps): bump github.com/emicklei/go-restful/v3 (#3768)
 * [6280d42e8](https://github.com/kubeovn/kube-ovn/commit/6280d42e87b7bad45bcf204d8fb8800d403ae5fc) docs: updated CHANGELOG.md (#3766)
 * [6e292a89b](https://github.com/kubeovn/kube-ovn/commit/6e292a89bd31602eee60785948f52bb58fc2d6e4) refactor ovn clusterrole (#3755)
 * [924d99192](https://github.com/kubeovn/kube-ovn/commit/924d9919236b4da2e3f308110f1a1917337dbc1f) e2e: fix subnet finalizers validation (#3765)
 * [5520490cd](https://github.com/kubeovn/kube-ovn/commit/5520490cd6c3503891bcfe6c7e4902d98bd304d9) kube-ovn-controller: fix finalizers migration (#3764)
 * [461d353c2](https://github.com/kubeovn/kube-ovn/commit/461d353c2587bafea4de0d04824a66fd0ab05662) remove trailing spaces (#3751)
 * [2a804103a](https://github.com/kubeovn/kube-ovn/commit/2a804103a109f3e1f82e0d9831fb9a0536fd8f01) if startOVNIC firstly, and setazName secondly, the ovn-ic-db may sync the old azname (#3759)
 * [21e821e7c](https://github.com/kubeovn/kube-ovn/commit/21e821e7cbcbc97dcc930e15c7e17429b5ea29a0) Add missing commas, correction of spelling errors (#3761)
 * [608d032da](https://github.com/kubeovn/kube-ovn/commit/608d032da9af37957bbf2cdda360171e42f30b6f) kube-ovn-controller: optimize finalizers migration (#3754)
 * [26c0affa0](https://github.com/kubeovn/kube-ovn/commit/26c0affa00c1ee88d08f46c40d34ddcb53dc0c9d) build(deps): bump google.golang.org/grpc from 1.61.1 to 1.62.0 (#3757)
 * [aa7d0d4c6](https://github.com/kubeovn/kube-ovn/commit/aa7d0d4c6884fd92b4c868077d054178ff8a1ca0) build(deps): bump github.com/parnurzeal/gorequest from 0.2.16 to 0.3.0 (#3756)
 * [d631f9de1](https://github.com/kubeovn/kube-ovn/commit/d631f9de114abe5aa696f194d45c1c568f77d7c3) log near err (#3753)
 * [d739b4121](https://github.com/kubeovn/kube-ovn/commit/d739b41217c44c123ca3c809212a0687e1f00c9b) fix: should use domain-qualified finalizer name (#3748)
 * [0405b92eb](https://github.com/kubeovn/kube-ovn/commit/0405b92eb6faa9177d51d4c6805e02eaada72481) Fix slr (#3746)
 * [fd3ae1d1c](https://github.com/kubeovn/kube-ovn/commit/fd3ae1d1c478afe5e67e0aa159bd9cd384d6819f) add helper script for release (#3752)
 * [26d7a9c45](https://github.com/kubeovn/kube-ovn/commit/26d7a9c4525fba43d5a9dbad2974a876f0047cee) docs: updated CHANGELOG.md (#3750)
 * [4594a2540](https://github.com/kubeovn/kube-ovn/commit/4594a2540881d69c34ff945c2cd8a51d1b1aac80) ovn-ic-server recover wait longer time (#3749)
 * [882c64ead](https://github.com/kubeovn/kube-ovn/commit/882c64ead05121d49690d8fe714896105520abdc) fix kubevirt vm e2e (#3747)
 * [71fa9c48c](https://github.com/kubeovn/kube-ovn/commit/71fa9c48c1da7b2320248ee6bc73c0c901d89e25) ci: bump azure/setup-helm to v4.0.0 (#3743)
 * [63f2c9964](https://github.com/kubeovn/kube-ovn/commit/63f2c9964b3e08c393ea4bb08b52f9544032dfc5) base: install libmnl0 instead of libmnl-dev (#3745)
 * [835e787ba](https://github.com/kubeovn/kube-ovn/commit/835e787ba071d53ecc6c3ed7bbc6d2d7af5d38c7) ci: collect ko logs for all kind clusters (#3744)
 * [e2ec8f118](https://github.com/kubeovn/kube-ovn/commit/e2ec8f118ebc663fa2d6c174a31c9eb090c91d5c) ci: fix ovn ic log file name (#3742)
 * [2217162d2](https://github.com/kubeovn/kube-ovn/commit/2217162d2b248bdf8552496f6c6fc0e21d072bb8) ci: bump cilium to v1.15.1 (#3735)
 * [b0f0bb71c](https://github.com/kubeovn/kube-ovn/commit/b0f0bb71cd72c017d32aaa1af779f9bf1e845464) remove invalid ovs build option (#3733)
 * [f95ca28fb](https://github.com/kubeovn/kube-ovn/commit/f95ca28fb0754e57477933f49cbc0cd3530a8211) add directory to charts (#3734)
 * [457e6bafb](https://github.com/kubeovn/kube-ovn/commit/457e6bafbbf23786ff02bad0899db47bfffcebce) bump k8s to v1.29.2 (#3725)
 * [8228a1e21](https://github.com/kubeovn/kube-ovn/commit/8228a1e21fe42d569439bcd6eaacf908966eafda) dpdk: remove unnecessary ovn patch (#3736)
 * [cc935bd79](https://github.com/kubeovn/kube-ovn/commit/cc935bd792eaaea49a9fa5a261778b322d83d9c7) remove util.UniqString() and util.ContainsString() (#3732)
 * [9abe9217f](https://github.com/kubeovn/kube-ovn/commit/9abe9217f5df24cc25756f28e1611754f7a09d96) fix kubevirt e2e (#3729)
 * [d265d7844](https://github.com/kubeovn/kube-ovn/commit/d265d7844d8a217d04b699cc18c3b871bc5726af) bump libovsdb (#3731)
 * [b98d6cee4](https://github.com/kubeovn/kube-ovn/commit/b98d6cee4c9c4adf23d3fd4aee1c01814a4c6eec) Makefile: bump kwok to v0.5.1 (#3730)
 * [7891205a3](https://github.com/kubeovn/kube-ovn/commit/7891205a3d2426f617eeef2716ef8a11b7f92e1f) fix typo (#3727)
 * [a7e70db93](https://github.com/kubeovn/kube-ovn/commit/a7e70db93fccd9bbedfd19b0a40f936d2397bb2d) update chart relase action workflow (#3728)
 * [3552db384](https://github.com/kubeovn/kube-ovn/commit/3552db38497ee97838b4c737ea515d2c218c422d) Fix: Resolve issue with skipped execution of sg annotations (#3700)
 * [28781c373](https://github.com/kubeovn/kube-ovn/commit/28781c37344d301ad331e803d4fed631838fded6) remove copilot4pr as the experiment ended
 * [8334fb3ac](https://github.com/kubeovn/kube-ovn/commit/8334fb3acf831b09e8ada81486d74e6c85844536) install.sh: wait ovn-ic-controller to be ready (#3704)
 * [704168dec](https://github.com/kubeovn/kube-ovn/commit/704168dec3032e48a5703871944a80be4b1acc12) fix: all gw nodes (#3724)
 * [40230ef74](https://github.com/kubeovn/kube-ovn/commit/40230ef74db8f89a154c599fe1c24cb94fc50e06) bump go to v1.22 (#3718)
 * [2a0a0be1f](https://github.com/kubeovn/kube-ovn/commit/2a0a0be1f545b5f88179ab71b8a748d41d248dbb) fix lint errors (#3719)
 * [01999c632](https://github.com/kubeovn/kube-ovn/commit/01999c6323838271a6184c6fd5b3a70f6dcf05cb) ovn: remove unnecessary patch (#3720)
 * [77db8f1b7](https://github.com/kubeovn/kube-ovn/commit/77db8f1b7edabb8613562eafd5c7103d04103c12) Add Ænix to the Adopters list (#3711)
 * [0e44d5e7a](https://github.com/kubeovn/kube-ovn/commit/0e44d5e7a2eb890591d88cfcfd705ba4563db20d) build(deps): bump golangci/golangci-lint-action from 3 to 4 (#3712)
 * [64ee14c52](https://github.com/kubeovn/kube-ovn/commit/64ee14c521235698191e46001ea754ec365a4f6e) build(deps): bump helm/kind-action from 1.8.0 to 1.9.0 (#3713)
 * [d7ad7f2d1](https://github.com/kubeovn/kube-ovn/commit/d7ad7f2d16d35a26fd18ff8010467c933b8d9fcd) build(deps): bump github.com/docker/docker (#3705)
 * [2e2fa7daa](https://github.com/kubeovn/kube-ovn/commit/2e2fa7daaa37cd3f07b08dbba26bb80b7a2b4be1) build(deps): bump golang.org/x/mod from 0.14.0 to 0.15.0 (#3706)
 * [a28af4aa1](https://github.com/kubeovn/kube-ovn/commit/a28af4aa1459812c3a194805dd267f7270526a93) fix: ip delete use wrong key, ip delete only once (#3698)
 * [d3ef8922a](https://github.com/kubeovn/kube-ovn/commit/d3ef8922aa27f44b9ba762c6d345f0d4ea6414fb) fix typo
 * [54f72e32c](https://github.com/kubeovn/kube-ovn/commit/54f72e32c0c28f5d0f3dd41362e7d1e89dc36c65) ci: bump kind to v0.21.0 (#3694)
 * [b6b1fccc3](https://github.com/kubeovn/kube-ovn/commit/b6b1fccc3a988ede4fb8e2b96f5d7d1ad6c4f1fb) update chart release workflow (#3691)
 * [36eba82e0](https://github.com/kubeovn/kube-ovn/commit/36eba82e0316b38faa865ff5d4b33c141860ac79) remove unused (#3696)
 * [a487f6108](https://github.com/kubeovn/kube-ovn/commit/a487f61081f4034b4a739a8a7d63e75f57053d82) kube-ovn-controller: remove unused codes (#3692)
 * [f2be2d58f](https://github.com/kubeovn/kube-ovn/commit/f2be2d58f3ce834e9140482e65624154ef2d7bb3) fix: security group base acl direction (#3690)
 * [0270012a7](https://github.com/kubeovn/kube-ovn/commit/0270012a704ab656f565092e67fc5567f82bf42f) remove fip controller (#3684)
 * [3716edd4a](https://github.com/kubeovn/kube-ovn/commit/3716edd4a08db4699567f5ebc511a04d41a92e3b) build(deps): bump github.com/osrg/gobgp/v3 from 3.22.0 to 3.23.0 (#3688)
 * [24acf0463](https://github.com/kubeovn/kube-ovn/commit/24acf046310079a11c7b4a7b7eb82ecf5218f972) build(deps): bump github.com/docker/docker (#3689)
 * [028e7ab36](https://github.com/kubeovn/kube-ovn/commit/028e7ab3602a33a06ae04dc810b92e13f0f63f73) build(deps): bump peter-evans/create-pull-request from 5 to 6 (#3685)
 * [c797ccd14](https://github.com/kubeovn/kube-ovn/commit/c797ccd1450f11312b6780ab0e1962ec5c7aefa4) build(deps): bump github.com/opencontainers/runc from 1.1.10 to 1.1.12 (#3686)
 * [0d6b4adda](https://github.com/kubeovn/kube-ovn/commit/0d6b4adda75c4a52961914e215d61cb80e974f2e) Compatible with controller deployment methods before kube-ovn 1.11.16 (#3677)
 * [6db8c840b](https://github.com/kubeovn/kube-ovn/commit/6db8c840b159e0d710487ff1c0b80ed7118eba8e) Support ip creation (#3667)
 * [38db6c34d](https://github.com/kubeovn/kube-ovn/commit/38db6c34d889f90968054e98ceae93fda4377498) set after genev_sys_6081 started (#3680)
 * [a79551c01](https://github.com/kubeovn/kube-ovn/commit/a79551c016201ca138a607e750ee0eb7133f3d61) ovn: add nb option version_compatibility (#3671)
 * [18f8f8ce4](https://github.com/kubeovn/kube-ovn/commit/18f8f8ce4ee71d5e58eafeea22856baaedf992bc) Makefile: fix install/upgrade chart (#3678)
 * [b0e813aa7](https://github.com/kubeovn/kube-ovn/commit/b0e813aa7856a2a0947d02a03cd3606f84aa7529) build(deps): bump github.com/evanphx/json-patch/v5 from 5.8.1 to 5.9.0 (#3676)
 * [f11820f8c](https://github.com/kubeovn/kube-ovn/commit/f11820f8cbc6baf6be77882fe9afbd1c0fc84527) ovn: do not send direct traffic between lports to conntrack (#3663)
 * [8f2cef01b](https://github.com/kubeovn/kube-ovn/commit/8f2cef01b6d6b37dff1d00175f722664ef9cd894) ovn-ic-ecmp refactor (#3632)
 * [ab4f1d861](https://github.com/kubeovn/kube-ovn/commit/ab4f1d8611236f407944e20dd084df458a565b56) fix getting vip with an empty name (#3673)
 * [5239bda19](https://github.com/kubeovn/kube-ovn/commit/5239bda19438b3b324d04ee8880614642289d61f) build(deps): bump dorny/paths-filter from 2 to 3 (#3672)
 * [0c82109ce](https://github.com/kubeovn/kube-ovn/commit/0c82109cefbca4f0aa5e92218957487b8fc1ad09) build(deps): bump github.com/docker/docker (#3670)
 * [c11ce4d58](https://github.com/kubeovn/kube-ovn/commit/c11ce4d5837d0eccc315ac5bcbbe3d9da5256184) replace github.com/golang/mock with go.uber.org/mock/mockgen (#3664)
 * [82986a64b](https://github.com/kubeovn/kube-ovn/commit/82986a64b14db83c0de920e789595601700a9930) build(deps): bump google.golang.org/grpc from 1.60.1 to 1.61.0 (#3669)
 * [01aabdb65](https://github.com/kubeovn/kube-ovn/commit/01aabdb650d9f46e67eb3d4681b4b70e0308ea05) fix 409 (#3662)
 * [0667fc009](https://github.com/kubeovn/kube-ovn/commit/0667fc009ef01c0605df6215c2c505e4b0ab717d) fix nil pointer (#3661)
 * [0e70364a1](https://github.com/kubeovn/kube-ovn/commit/0e70364a10bd4a86de1eb6a5d5ad62c3952c9443) docs: updated CHANGELOG.md (#3659)
 * [d4a815d83](https://github.com/kubeovn/kube-ovn/commit/d4a815d8369530360f652be2532f6701ad0786e4) chart: fix parsing image tag when the image url contains a port (#3644)
 * [b2dead02e](https://github.com/kubeovn/kube-ovn/commit/b2dead02e1992f61debe315fe0ef0b62c5d9567a) ovs: reduce cpu utilization (#3650)
 * [42e91567c](https://github.com/kubeovn/kube-ovn/commit/42e91567c336cd0edfb4637a0515caec47b33879) remove static route for custom vpc (#3647)
 * [d92e929e3](https://github.com/kubeovn/kube-ovn/commit/d92e929e3c293e70993cdab506fac52842f326c7) bump go modules (#3654)
 * [3104a6f0a](https://github.com/kubeovn/kube-ovn/commit/3104a6f0a366e8cc3a758c1c78880f886aa0e594) kube-ovn-monitor and kube-ovn-pinger export pprof path (#3649)
 * [8838b08d2](https://github.com/kubeovn/kube-ovn/commit/8838b08d227615eff7dc9829465908cd43d0cf89) build(deps): bump github.com/onsi/gomega from 1.31.0 to 1.31.1 (#3652)
 * [ea9cec2a9](https://github.com/kubeovn/kube-ovn/commit/ea9cec2a9f2c7956fd62b86566719c47ce22e8a0) build(deps): bump github.com/docker/docker (#3651)
 * [a69e6a3bc](https://github.com/kubeovn/kube-ovn/commit/a69e6a3bc46d36e8b056e3b2b6ce5bcb6868f9c9) build(deps): bump github.com/onsi/gomega from 1.30.0 to 1.31.0 (#3641)
 * [6da2734f0](https://github.com/kubeovn/kube-ovn/commit/6da2734f068b24ef79d3ab28ff1d2cc4e365b759) build(deps): bump actions/cache from 3 to 4 (#3643)
 * [44e248404](https://github.com/kubeovn/kube-ovn/commit/44e248404cfe374f94e655cd77a32f6ea509db3f) build(deps): bump the k8s-io group with 1 update (#3638)
 * [9f262fc32](https://github.com/kubeovn/kube-ovn/commit/9f262fc3252489ab796cc8d11eddc749bec794fb) build(deps): bump github.com/onsi/ginkgo/v2 from 2.14.0 to 2.15.0 (#3642)
 * [5f9ba94f5](https://github.com/kubeovn/kube-ovn/commit/5f9ba94f59ffc79ef667c51b0de9935a047e1e59) build(deps): bump github.com/evanphx/json-patch/v5 from 5.8.0 to 5.8.1 (#3633)
 * [f13ec5ae7](https://github.com/kubeovn/kube-ovn/commit/f13ec5ae7357cacf84c528891401722e81f94ac3) e2e for vip dual stack with security group (#3630)
 * [d5685ff78](https://github.com/kubeovn/kube-ovn/commit/d5685ff7831cd11b03a6dbbf1c07f927bc3f5463) SYSCTL_IPV4_IP_NO_PMTU_DISC change default to 0 (#3596)
 * [e1643a0c0](https://github.com/kubeovn/kube-ovn/commit/e1643a0c0cd72e57cd0949f26a9f9c0128a24045) remove duplicated function (#3629)
 * [790d9cb15](https://github.com/kubeovn/kube-ovn/commit/790d9cb159afd8e3f0cf7d7145e5f4b7d981f86e) build(deps): bump github.com/evanphx/json-patch/v5 from 5.7.0 to 5.8.0 (#3628)
 * [d487d711d](https://github.com/kubeovn/kube-ovn/commit/d487d711db6e9ff4c21d782c6e7e6f5e6c2bc2bc) fix: static one ip in dualstack subnet (#3625)
 * [45b00710c](https://github.com/kubeovn/kube-ovn/commit/45b00710cfc9e2b77e42bc297ab15908ab7dfa12) support vip dual stack (#3617)
 * [8707ff750](https://github.com/kubeovn/kube-ovn/commit/8707ff750fb244024a595eacce473d9a0f59f365) build(deps): bump github.com/onsi/ginkgo/v2 from 2.13.2 to 2.14.0 (#3624)
 * [bf7eb65d0](https://github.com/kubeovn/kube-ovn/commit/bf7eb65d0f46769f46adee43f3a4c8e992b7744f) build(deps): bump k8s.io/klog/v2 from 2.110.1 to 2.120.0 (#3622)
 * [bafdeae81](https://github.com/kubeovn/kube-ovn/commit/bafdeae81c3bd8de38c1efb8bbaa6a6dc7f84216) fix: empty  loop (#3615)
 * [2641a5932](https://github.com/kubeovn/kube-ovn/commit/2641a593225ef80f8afcc5f50a29cf0b335fe2ca) chart: fix ovs-ovn upgrade (#3613)
 * [524fa50c9](https://github.com/kubeovn/kube-ovn/commit/524fa50c9e153b36d3492b89179c21a6f28339ce) build(deps): bump github.com/emicklei/go-restful/v3 (#3619)
 * [9a77310e0](https://github.com/kubeovn/kube-ovn/commit/9a77310e0f6fc3ae2a5bb4f5f9f3d4604598d41c) fix: subnet status usingIp not consistent with availableIPRange and usingIPRange (#3608)
 * [cb28bed29](https://github.com/kubeovn/kube-ovn/commit/cb28bed294c6ec7aa1c89c78e49045cc34f5542f) build(deps): bump golang.org/x/sys from 0.15.0 to 0.16.0 (#3607)
 * [7c6c9127b](https://github.com/kubeovn/kube-ovn/commit/7c6c9127b49a8e8bb590b38721487a050149446a) build(deps): bump github.com/emicklei/go-restful/v3 (#3606)
 * [032314131](https://github.com/kubeovn/kube-ovn/commit/0323141312aca21bc5d174f323b40635ca9f8e46) ci: install kube-ovn with helm chart (#3599)
 * [83ca99d39](https://github.com/kubeovn/kube-ovn/commit/83ca99d3909173b61458cdd9de2320f685b1f906) ovs: increase cpu limit to 2 cores (#3530)
 * [5589e8ed8](https://github.com/kubeovn/kube-ovn/commit/5589e8ed870f333d70e3e371cc6efa4debf0e18e) update policy route when subnet cidr is changed (#3587)
 * [5d9a24a4d](https://github.com/kubeovn/kube-ovn/commit/5d9a24a4ded4f078dabc6ae926e2426cac0e8879) update ipset to v7.17 (#3601)
 * [1cd058f77](https://github.com/kubeovn/kube-ovn/commit/1cd058f7788199208415c3d318efcc9a41254813) fix ko (#3600)
 * [d4d2cbd22](https://github.com/kubeovn/kube-ovn/commit/d4d2cbd229d7de96c8d5c359fdf59b41b4076a4d) build(deps): bump github.com/osrg/gobgp/v3 from 3.21.0 to 3.22.0 (#3603)
 * [f7bd2520f](https://github.com/kubeovn/kube-ovn/commit/f7bd2520f34b6c6c236302f299f38c09519a55f0) add: MASTER_NODES_LABEL (#3597)
 * [d8e51e225](https://github.com/kubeovn/kube-ovn/commit/d8e51e225851afc7c7a2e9ed8f18d6378514c72a) chart: generate SSL certificates (#3598)
 * [37d832b37](https://github.com/kubeovn/kube-ovn/commit/37d832b37f7ab96017170f9848e0dbd63e314abb) fix: lrp type eip no finalizer (#3584)
 * [f2e494cae](https://github.com/kubeovn/kube-ovn/commit/f2e494caea9632c153e02360916e789ed97d450b) fix: ovn eip type (#3590)
 * [a91bc8ed0](https://github.com/kubeovn/kube-ovn/commit/a91bc8ed01af4e84042bd243363a21e15a6f4bfa) do not count ips in excludeIPs as available and using IPs (#3582)
 * [0052a409c](https://github.com/kubeovn/kube-ovn/commit/0052a409c374553021058a1c8214bb2200256df5) chart: support Talos Linux (#3593)
 * [09e2a0d5d](https://github.com/kubeovn/kube-ovn/commit/09e2a0d5db411f107f64129530c0e0bced4fc51e) chart: Autodiscover master nodes (#3591)
 * [897ff9f2b](https://github.com/kubeovn/kube-ovn/commit/897ff9f2badbfc40ffd685dd4411b3cb698e8eec) bump k8s to v1.28.5 (#3594)
 * [b3ea8b9c4](https://github.com/kubeovn/kube-ovn/commit/b3ea8b9c430029771d9eab75d706e2d5f4c1e9ad) fix static build for CNI-plugin (#3585)
 * [80cb018ee](https://github.com/kubeovn/kube-ovn/commit/80cb018ee59679e1b428110a600b67fbc7745a62) fix: duplicate ip deletion (#3589)
 * [f60bb3d5b](https://github.com/kubeovn/kube-ovn/commit/f60bb3d5b8ef09c50f199c8718a3876232897cee) fix security issue (#3588)
 * [44654a985](https://github.com/kubeovn/kube-ovn/commit/44654a985ffbcad0ebc5706899a04626c20e1dbc) chart: fix specifying namespace (#3592)
 * [b823c9bb9](https://github.com/kubeovn/kube-ovn/commit/b823c9bb9da15b9485a4bcb7dac289f902ab7342) build(deps): bump github.com/prometheus/client_golang (#3586)
 * [b0df19270](https://github.com/kubeovn/kube-ovn/commit/b0df19270aec959185665393224d93ed743269ff) ovn0 ipv6 addr gen mode set 0 (#3580)
 * [f9b3aa49a](https://github.com/kubeovn/kube-ovn/commit/f9b3aa49a46c10b89f566fe090e492b5ac4dd279) fix: G307 (#3575)
 * [fea6165eb](https://github.com/kubeovn/kube-ovn/commit/fea6165eb4d23ebd79c5996c979a0ee8dc5dc5ad) add bfd check from lrp to  underlay gw (#3441)
 * [f7ded7dff](https://github.com/kubeovn/kube-ovn/commit/f7ded7dff2a69ad22527e4bfef497a7c18c3d4a1) fix: add err log (#3572)
 * [71007d8d1](https://github.com/kubeovn/kube-ovn/commit/71007d8d17b3cbfc466b11de1d15a1a89f8c995a) build(deps): bump google.golang.org/protobuf from 1.31.0 to 1.32.0 (#3571)
 * [c6fefd4c1](https://github.com/kubeovn/kube-ovn/commit/c6fefd4c18b18e48d4c34d4b0b09af4cb6793873) fix: calculate usings before removed finalizer (#3570)
 * [ef35f62f0](https://github.com/kubeovn/kube-ovn/commit/ef35f62f04990d9ea2f817e6aa80cb3d48b3a8dd) fix speaker log (#3487)
 * [a106e7fe0](https://github.com/kubeovn/kube-ovn/commit/a106e7fe0690185513dd0eaa339cce01ed6bbe6b) fix: cleanup.sh (#3569)
 * [a4c2237c4](https://github.com/kubeovn/kube-ovn/commit/a4c2237c4b1a2f1225d8f9c70e755e202604098c) Makefile: fix kwok installation (#3561)
 * [69c3d6c74](https://github.com/kubeovn/kube-ovn/commit/69c3d6c74b1a9699fce767386d79421377c2daee) fix u2o infinity recycle (#3563)
 * [f6c83bad4](https://github.com/kubeovn/kube-ovn/commit/f6c83bad45f54e03482874c37434e87a3f869c50) roll back to v3 due to massive request timeout https://github.com/actions/download-artifact/issues/249
 * [82e11ecc0](https://github.com/kubeovn/kube-ovn/commit/82e11ecc0086d46aba79b6bc50c41a787bee500a) Fix 1.12 ipam deletion (#3554)
 * [72e90fd16](https://github.com/kubeovn/kube-ovn/commit/72e90fd16e9f9460f47cb06502e97f377d36e268) htbqos is not used anymore (#3559)
 * [51dcf7f79](https://github.com/kubeovn/kube-ovn/commit/51dcf7f79e6a37433cc8872a4f48f772ad16f27b) delete lsp and ipam with ip (#3540)
 * [67a83bf6a](https://github.com/kubeovn/kube-ovn/commit/67a83bf6a57010d1917e2527c8ac0d475e1400eb) skip subnet mtu e2e prior to v1.9 (#3557)
 * [a4545549d](https://github.com/kubeovn/kube-ovn/commit/a4545549de679c8d79b0038db24385a6bb8c7864) add np prefix to networkpolicy name when networkpolicy's name starts with number (#3551)
 * [4850f0a7e](https://github.com/kubeovn/kube-ovn/commit/4850f0a7ed5cd048cad33acdfb9feeb214b86843) do not calculate subnet.spec.excludeIPs as availableIPs (#3550)
 * [794a52a81](https://github.com/kubeovn/kube-ovn/commit/794a52a8130639bfa89883c64c836896bca21790) Fix: centos8 reached EOL and ovn 21.03 tarball unmatch with ovs 2.15 tarball issue (#3546)
 * [b09d249ba](https://github.com/kubeovn/kube-ovn/commit/b09d249baaece56dda23dfbee20a99ed9ba38dde) build(deps): bump google.golang.org/grpc from 1.60.0 to 1.60.1 (#3547)
 * [8081f437c](https://github.com/kubeovn/kube-ovn/commit/8081f437c777eb3dfe2c60f11ab78894d6fd77cb) keep ovs cpu:mem 1:1 and show pod status in kube-system namespace before diagnose (#3537)
 * [2f76694f8](https://github.com/kubeovn/kube-ovn/commit/2f76694f8d4d53d6189236267c6575a6602bf32d) fix ovn ic not clean lsp and lrp when az name contains "-" (#3542)
 * [65cf9ce47](https://github.com/kubeovn/kube-ovn/commit/65cf9ce472e7434a3d6b69e0ca71560a444cb3cd) fix: apply changes to the latest version (#3514)
 * [4e64005c2](https://github.com/kubeovn/kube-ovn/commit/4e64005c24e63e203346dba884a129b4e4c72b53) build(deps): bump golang.org/x/crypto from 0.16.0 to 0.17.0 (#3544)
 * [1712eda5a](https://github.com/kubeovn/kube-ovn/commit/1712eda5a8d752fb2c36e960141c9c4c5200b083) Revert "ovn-central: check raft inconsistency from nb/sb logs (#3532)"
 * [f25e20daa](https://github.com/kubeovn/kube-ovn/commit/f25e20daa7da04cd72519adc4049e67cfcdb47c8) SgPorts should be updated when pod deleted (#3516)
 * [44705bb89](https://github.com/kubeovn/kube-ovn/commit/44705bb898b9b1b334b235f79cba42956b406e96) bump go modules (#3505)
 * [2c7acdbb6](https://github.com/kubeovn/kube-ovn/commit/2c7acdbb6aa6273f6cab3d59a95f97353dfa4cac) ovn-central: check raft inconsistency from nb/sb logs (#3532)
 * [2c3eacb72](https://github.com/kubeovn/kube-ovn/commit/2c3eacb72b8faced363316768129a0798533b5f3) build(deps): bump actions/download-artifact from 3 to 4 (#3535)
 * [f07636b22](https://github.com/kubeovn/kube-ovn/commit/f07636b22fdb0086297ac387e6a95d976969cdcf) build(deps): bump actions/upload-artifact from 3 to 4 (#3534)
 * [a8ee70b32](https://github.com/kubeovn/kube-ovn/commit/a8ee70b32bd3d35d76307110ef2618ce915fa4b7) fix chassis gc (#3525)
 * [508535a43](https://github.com/kubeovn/kube-ovn/commit/508535a438c92d2021b5f88a9f0a4dc48a574e6b) Makefile: fix clab bgp cleanup (#3524)
 * [b0e254dea](https://github.com/kubeovn/kube-ovn/commit/b0e254deaee4ef4b373be1a253c9ec655400109c) docs: updated CHANGELOG.md (#3523)
 * [dbb946989](https://github.com/kubeovn/kube-ovn/commit/dbb9469893359762b19d3e07753acf9e9a17133c) build(deps): bump github/codeql-action from 2 to 3 (#3522)
 * [a1020196a](https://github.com/kubeovn/kube-ovn/commit/a1020196af71bb41f19083ea38d84d58ce83a4af) build(deps): bump google.golang.org/grpc from 1.59.0 to 1.60.0 (#3517)
 * [f49d93b0a](https://github.com/kubeovn/kube-ovn/commit/f49d93b0a9eb0b229d0bfdea7dd8ab531b5abf44) netpol: do not create unused address set for services (#3208)
 * [1f9dbf8d4](https://github.com/kubeovn/kube-ovn/commit/1f9dbf8d4f73b389eaa40b0b63dfe2bf849ef5c3) fix: gc keep lb for dnat after restart controller (#3508)
 * [869a9bc3c](https://github.com/kubeovn/kube-ovn/commit/869a9bc3c6718e92d6daecf748fb34af29386cad) ci: fix ovn-vpc-nat-gw e2e failed (#3503)
 * [7e285581a](https://github.com/kubeovn/kube-ovn/commit/7e285581aed29bf61a8b67cf4f00a43a421d9472) cni-server: set sysctl variable net.ipv4.ip_no_pmtu_disc to 1 by default (#3504)
 * [079e6c8dc](https://github.com/kubeovn/kube-ovn/commit/079e6c8dceaed82684cf8af4172018adab162b4b) build(deps): bump actions/stale from 8 to 9 (#3507)
 * [a999b009f](https://github.com/kubeovn/kube-ovn/commit/a999b009f31cccb8d1f698a800df33ea4489805e) fix: skip unnecessary del port request (#3501)
 * [23419b8e1](https://github.com/kubeovn/kube-ovn/commit/23419b8e1ddcd7d149817479ff3652a3b125161b) update image
 * [d243606a2](https://github.com/kubeovn/kube-ovn/commit/d243606a268947aa5b431f6672b6de2f4ac434a3) fix: duplicate gw nodes (#3500)
 * [7e5085020](https://github.com/kubeovn/kube-ovn/commit/7e508502067e273ae7740bb7c1828f603575d4ae) build(deps): bump actions/setup-go from 4 to 5 (#3499)
 * [260cca8de](https://github.com/kubeovn/kube-ovn/commit/260cca8de5b384c03fe1cd1ce44baf172ed2c8be) fix typo (#3494)
 * [6a939d928](https://github.com/kubeovn/kube-ovn/commit/6a939d928adc6a26ecc956b644cb339a3e6d4b70) fix: lost gc lsp in previous pr (#3493)
 * [874da5bd9](https://github.com/kubeovn/kube-ovn/commit/874da5bd9ca0ed610707357d24fbe4173bd55f8e) delete String() function (#3488)
 * [487bcd036](https://github.com/kubeovn/kube-ovn/commit/487bcd0365c07986e05856c4e008412e135a34ce) build(deps): bump github.com/containernetworking/plugins (#3489)
 * [be730193b](https://github.com/kubeovn/kube-ovn/commit/be730193ba423b7ba68f586d05318b6bc04bd237) e2e aap for security group (#3459)
 * [984f227f5](https://github.com/kubeovn/kube-ovn/commit/984f227f525712296c81bcf6ebe7c283e2a44183) iptables drop invalid rst (#3484)
 * [7436dc4de](https://github.com/kubeovn/kube-ovn/commit/7436dc4de1e570fd245b957466e080d5072499e2) add part of e2es for release-1.12-mc (#3478)
 * [f7ff76b8d](https://github.com/kubeovn/kube-ovn/commit/f7ff76b8d9d3442338250a9c8f2781f248b3a95f) make vpc_nat_gateway redo eip fip snat dnat iptables rules after k8s cluster reboot (#3267)
 * [8628642e0](https://github.com/kubeovn/kube-ovn/commit/8628642e0a03b1df01bea25fafd3da0f605a354a) fix: github action cannot ping 114 (#3486)
 * [39e920141](https://github.com/kubeovn/kube-ovn/commit/39e920141e6f8b01c421946143e5c05c4892b858) fix: check chassis before creation (#3482)
 * [fc7697704](https://github.com/kubeovn/kube-ovn/commit/fc7697704cea33412ba36ffcd4e87cb7f9e1faac) build(deps): bump github.com/osrg/gobgp/v3 from 3.20.0 to 3.21.0 (#3481)
 * [19269984a](https://github.com/kubeovn/kube-ovn/commit/19269984a6e8f9d5cc7ca4bfd06596b2d995920e) fix ovn eip not calculated (#3477)
 * [91dea3587](https://github.com/kubeovn/kube-ovn/commit/91dea35873a549455f8c43f805889361ba4c75d0) fix: calculate subnet before handle finalizer (#3469)
 * [f25700396](https://github.com/kubeovn/kube-ovn/commit/f25700396973f5287e1e603c760b0cbcb1ec2732) fix: ipam clean all pod nic ip address and mac even if just delete a nic (#3453)
 * [c6fe51733](https://github.com/kubeovn/kube-ovn/commit/c6fe51733b3891271ba9224910e83ee9504f1e2c) schedule kube-ovn-controller on the kube-ovn-master node (#3479)
 * [6233434a8](https://github.com/kubeovn/kube-ovn/commit/6233434a8f5a14e6f6ce2e29411840b9fb670d39) delete vm's lsp and release ipam.ip (#3476)
 * [ab0cf2214](https://github.com/kubeovn/kube-ovn/commit/ab0cf221482c85bcc604c79409c2bd34ef75a191) build(deps): bump github.com/onsi/ginkgo/v2 from 2.13.1 to 2.13.2 (#3475)
 * [45cee5780](https://github.com/kubeovn/kube-ovn/commit/45cee57808f15c603e0419f1b101c9141644170f) ci: remove branch release-1.12-cmss
 * [b1de14c0f](https://github.com/kubeovn/kube-ovn/commit/b1de14c0fcbf2d08e6021e98d23eef4e840e2cf7) fix: add static route for custom vpc when create subnet with (#3462)
 * [f536b24a1](https://github.com/kubeovn/kube-ovn/commit/f536b24a1419b5e9fd1c895bd812621b74c525f7) build(deps): bump golang.org/x/time from 0.4.0 to 0.5.0 (#3463)
 * [3351db34e](https://github.com/kubeovn/kube-ovn/commit/3351db34eaeac862809b7c9a23ed66afbf4618a7) build(deps): bump golang.org/x/sys (#3464)
 * [07798be1e](https://github.com/kubeovn/kube-ovn/commit/07798be1ebc9231e49acd67d8b87774c26b27c26) kube-ovn-cni: fix pinger result when timeout is reached (#3457)
 * [e334b60c3](https://github.com/kubeovn/kube-ovn/commit/e334b60c3df21b961147d55fb6523604b88c7f92) ovs-healthcheck: ignore error when log file does not exist (#3456)
 * [f8435b319](https://github.com/kubeovn/kube-ovn/commit/f8435b319980e9e287e97f2154f8bdb567ec3e01) ipam: fix duplicate allocation after cidr expansion (#3455)
 * [be68461ff](https://github.com/kubeovn/kube-ovn/commit/be68461ffa6ff7ab2d3142f81b2481a04d5fa497) feat:  support dpdk interface hotplug for kubevirt (#3438)
 * [d970e528f](https://github.com/kubeovn/kube-ovn/commit/d970e528fda0f87b724558b15d07c95e07cb8350) fix e2e install failed (#3450)
 * [ea14a7161](https://github.com/kubeovn/kube-ovn/commit/ea14a71615d41c02462e689a8a5b754cd1166b89) ci: fix scheduled workflow runs of building base images (#3449)
 * [b6a931fc2](https://github.com/kubeovn/kube-ovn/commit/b6a931fc2dc0c5cb1f1829ef43ae29c2d25d0114) readd assigned ip addresses to ipam when subnet has been changed (#3447)
 * [54949dbc2](https://github.com/kubeovn/kube-ovn/commit/54949dbc219ed6db541b8187370bd98f8a27efe7) base: fix missing CFLAGS -fPIC for arm64 (#3428)
 * [14f136fe3](https://github.com/kubeovn/kube-ovn/commit/14f136fe33228479b3a93457807ae6c61af9ab16) support allow address pairs for SecurityGroup (#3442)
 * [0c4fc68fc](https://github.com/kubeovn/kube-ovn/commit/0c4fc68fcd159b9d8919e2a21f9dd00b48f6dbdb) ci: support specifying branch (#3444)
 * [4622372ee](https://github.com/kubeovn/kube-ovn/commit/4622372ee8da97c279f6adc048178963845fbc71) OvnEip using default external subnet when not set (#3445)
 * [64a5cf523](https://github.com/kubeovn/kube-ovn/commit/64a5cf523de88d1dee28ddb476a82120c9201df2) ci: fix deepflow installation (#3443)
 * [c48863d82](https://github.com/kubeovn/kube-ovn/commit/c48863d82db6ab6f8ed8d25c3741a59bbfe5f034) fix: release VF assigned to the Pod back to the host network namespace (#3440)
 * [05a9a162c](https://github.com/kubeovn/kube-ovn/commit/05a9a162c0751a3274727065b009e5dc6c161155) fix:  multus network status not find dpdk interface name (#3432)
 * [2ba4c78a1](https://github.com/kubeovn/kube-ovn/commit/2ba4c78a1866ce184ce87611aeaf3f18901e9269) e2e for extra external subnets (#3435)
 * [68e96e77d](https://github.com/kubeovn/kube-ovn/commit/68e96e77d5206c0a42ceeb1f72469262b2f77993) Correcting spelling errors in word static (#3434)
 * [cffb893cf](https://github.com/kubeovn/kube-ovn/commit/cffb893cf22069077ed00ed6a75cae43ffb9515c) ci: fix missing environment variables (#3430)
 * [024b95e8f](https://github.com/kubeovn/kube-ovn/commit/024b95e8fd39c0111b4dc77e892698f66c6e1971) e2e: fix VIP test failure in v1.12
 * [02b60df1b](https://github.com/kubeovn/kube-ovn/commit/02b60df1baacc4e2a5b7b7825e83b89c4e35495b) base: fix dpdk build failure (#3426)
 * [38fefbe7a](https://github.com/kubeovn/kube-ovn/commit/38fefbe7a0aa2ac712c7d42617e7da2eeb578480) bump k8s to v1.28.4 (#3423)
 * [f52f2ef52](https://github.com/kubeovn/kube-ovn/commit/f52f2ef52215dced0e1691196a51bd1a78c8d3c7) base: fix ovn build failure (#3340)
 * [43f9b64e0](https://github.com/kubeovn/kube-ovn/commit/43f9b64e0707f6439107d2896bccfad5decdc3a5) ci: fix scheduled base image build workflow (#3415)
 * [965360e32](https://github.com/kubeovn/kube-ovn/commit/965360e32ac0e4929115e1504b8865fc38245753) fix:  lsp dhcp options set failed when subnet dhcp option is enabled (#3422)
 * [62b0e21e3](https://github.com/kubeovn/kube-ovn/commit/62b0e21e3f050416b40f9d5302838e1db75cbc04) chart: move CRD into directory crds (#3420)
 * [f60ad1fd2](https://github.com/kubeovn/kube-ovn/commit/f60ad1fd269dbadcd47daae337b6c3486ae5fd1e) trivy: ignore CVE-2023-5528
 * [4452eeaaa](https://github.com/kubeovn/kube-ovn/commit/4452eeaaa268cab474f5e8a7a3c74b9c4b0adbe1) base: fix ovn-northd/ovn-controller not creating pidfile in arm64 (#3413)
 * [e2e3b4f97](https://github.com/kubeovn/kube-ovn/commit/e2e3b4f975068d5690e2939b8082547975119353) fix kube-ovn-monitor probe (#3409)
 * [f7e216096](https://github.com/kubeovn/kube-ovn/commit/f7e2160962a91a9b2cbf70ab61588bf64ba01dae) ci: fix dpdk jobs (#3405)
 * [bbdf6c6f9](https://github.com/kubeovn/kube-ovn/commit/bbdf6c6f9b8a73b0f315bf3a8f5f7329e08c537e) support ovn ic ecmp (#3348)
 * [57b7d92fb](https://github.com/kubeovn/kube-ovn/commit/57b7d92fb5513c987c2988814ab86184572f1917) e2e for allowed address pair (#3394)
 * [df98a0a89](https://github.com/kubeovn/kube-ovn/commit/df98a0a89a904d0493914488fa5a7b06fb9fbd7f) ci: free disk space for all x86 jobs (#3406)
 * [819700b71](https://github.com/kubeovn/kube-ovn/commit/819700b7162360f505beb9838a655a3afc47bc6b) build(deps): bump k8s.io/klog/v2 from 2.100.1 to 2.110.1 (#3353)
 * [34cb63229](https://github.com/kubeovn/kube-ovn/commit/34cb63229d6b2b6f79d2b1306c92adfb65ab69ae) ci: free disk space (#3404)
 * [dd73d833a](https://github.com/kubeovn/kube-ovn/commit/dd73d833ac9ed442ee27d0a27c201a7cf2bdc2ee) build(deps): bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.45.0 to 0.46.0 (#3402)
 * [ca5ebcd71](https://github.com/kubeovn/kube-ovn/commit/ca5ebcd71c858fbbf0695c146135aaffd78351f7) build(deps): bump github.com/onsi/ginkgo/v2 from 2.13.0 to 2.13.1 (#3401)
 * [77785601a](https://github.com/kubeovn/kube-ovn/commit/77785601a5d52f10fb47804d2e803fec6c02ed14) fix: wrong usage about DeepEqual (#3396)
 * [8d3248936](https://github.com/kubeovn/kube-ovn/commit/8d3248936382998c694a4a919fbc1705b05472a1) build(deps): bump github.com/onsi/gomega from 1.29.0 to 1.30.0 (#3398)
 * [7d22425bb](https://github.com/kubeovn/kube-ovn/commit/7d22425bba829f7af8cb08f6dce340980307c04b) subnet support config mtu (#3367)
 * [a76445ed4](https://github.com/kubeovn/kube-ovn/commit/a76445ed4e0d0e1c9b7007f91ba16189a1538bc9) fix dualStack network checkgw raise panic (#3392)
 * [182c075a4](https://github.com/kubeovn/kube-ovn/commit/182c075a49e230426092796ef63091db888ad9e4) fix dpdk workflow (#3384)
 * [35b103b51](https://github.com/kubeovn/kube-ovn/commit/35b103b51167069faf6bdf256d6bfb5f59384718) support aap through vip crd using selector (#3359)
 * [08d8321e0](https://github.com/kubeovn/kube-ovn/commit/08d8321e0c65a440a48cf0569d75a8484e260204) fix:  gc delete multus ip cr and lsp setting when enable keep vm ip (#3378)
 * [a6b7045d3](https://github.com/kubeovn/kube-ovn/commit/a6b7045d3a7fe9e3291be7c493e469596ee6b541) add kube-ovn-controller nodeAffinity prefer not on ic gateway (#3376)
 * [ebbe6d95f](https://github.com/kubeovn/kube-ovn/commit/ebbe6d95f624072f85f67b98801935a31ff0ff37) Add clean ic az db (#3372)
 * [c5e877372](https://github.com/kubeovn/kube-ovn/commit/c5e87737223637e0a3e7a57f7fdc9184c8c9ea7c) build(deps): bump github.com/moby/sys/mountinfo from 0.6.2 to 0.7.1 (#3389)
 * [75bbeb202](https://github.com/kubeovn/kube-ovn/commit/75bbeb2024923bf35147e5584b9993a7bee3aea8) fix: externalID map should not include external_ids (#3385)
 * [92b386037](https://github.com/kubeovn/kube-ovn/commit/92b3860379deb841d913403f2a6c32758ea75f2e) fix dependency issues
 * [c3fd09d28](https://github.com/kubeovn/kube-ovn/commit/c3fd09d2835d70dc6eb06d5b883a139acea2b377) build(deps): bump golang.org/x/mod from 0.13.0 to 0.14.0 (#3380)
 * [4d9d8b36b](https://github.com/kubeovn/kube-ovn/commit/4d9d8b36b360f914639daf5969d303d708b801c9) build(deps): bump golang.org/x/sys from 0.13.0 to 0.14.0 (#3381)
 * [9604c8193](https://github.com/kubeovn/kube-ovn/commit/9604c8193738880ac1409d2bc56634484b7cad98) build(deps): bump golang.org/x/time from 0.3.0 to 0.4.0 (#3383)
 * [b8e63dba8](https://github.com/kubeovn/kube-ovn/commit/b8e63dba86eb1e9670bf6da827903e2c05e1ee66) docs: updated CHANGELOG.md (#3379)
 * [5ae3d25da](https://github.com/kubeovn/kube-ovn/commit/5ae3d25da24cf7c5c0ecbfe0b378d9b57983a460) kube-ovn-dpdk building need its dpdk base img (#3371)
 * [cc1f6adb4](https://github.com/kubeovn/kube-ovn/commit/cc1f6adb48cdbce33748f61367e772ae6fe690ba) delete check for existing ip cr (#3361)
 * [3bb2b3078](https://github.com/kubeovn/kube-ovn/commit/3bb2b3078fce72046ac45881361bb8a8c6fcf59b) fix IP residue after changing subnet of vm in some scenarios (#3370)
 * [20125e8fe](https://github.com/kubeovn/kube-ovn/commit/20125e8fe9fb2d10bc5dd0fb26fcbf6082527ef3) fix start-ovs-dpdk-v2.sh syntax error (#3373)
 * [b1d2d4dc2](https://github.com/kubeovn/kube-ovn/commit/b1d2d4dc2324261af8b05283a81845189c6b8e5d) build(deps): bump github.com/Microsoft/hcsshim from 0.11.2 to 0.11.4 (#3374)
 * [435fda5a9](https://github.com/kubeovn/kube-ovn/commit/435fda5a949257805e6c5adeb56809a77ff5bee4) build(deps): bump helm/chart-releaser-action from 1.5.0 to 1.6.0 (#3375)
 * [99ad03e90](https://github.com/kubeovn/kube-ovn/commit/99ad03e90c2f1cb6654b2be45439b1d82334d2bf) Add cmss work flow (#3358)
 * [2b200e8fa](https://github.com/kubeovn/kube-ovn/commit/2b200e8fac725126e647f21f0a35b757a7c439c0) kube-ovn-controller: fix memory growth caused by unused workqueue (#3366)
 * [31ae1fd57](https://github.com/kubeovn/kube-ovn/commit/31ae1fd578412b774559499ac6abe59a9f9de4eb) build(deps): bump github.com/osrg/gobgp/v3 from 3.19.0 to 3.20.0 (#3362)
 * [6b82297c7](https://github.com/kubeovn/kube-ovn/commit/6b82297c702989b7fb352ef84842c270208d340e) add e2e case for netpol (#3355)
 * [91994f925](https://github.com/kubeovn/kube-ovn/commit/91994f925c4aab5978a6034380402d9cd9448fba) build(deps): bump actions/stale from 7 to 8 (#3352)
 * [66f365c8e](https://github.com/kubeovn/kube-ovn/commit/66f365c8ed03a810daf2f0d20ec09936b6ec85d0) Signed-off-by: changluyi <clyi@alauda.io> (#3349)
 * [ce9975bd4](https://github.com/kubeovn/kube-ovn/commit/ce9975bd4cffcc44babb802c465aef4d34043cfe) Allow same subnet traffic (#3344)
 * [0a98e42e2](https://github.com/kubeovn/kube-ovn/commit/0a98e42e241aa674c2e3202a1b5cb0812f7a7569) add stale bot
 * [68747f904](https://github.com/kubeovn/kube-ovn/commit/68747f90440790ecdb88d5e12ef07de5046bba94) add compact switch (#3337)
 * [ac770718f](https://github.com/kubeovn/kube-ovn/commit/ac770718fe6540503aa28b9ed9e6b60eb0a50bb0) Add Layer 2 forwarding for subnet ports again (#3300)
 * [7484a2b76](https://github.com/kubeovn/kube-ovn/commit/7484a2b76d4010a3d099fe3199e2c7a129ebf1f8) build(deps): bump github.com/docker/docker (#3346)
 * [6db0ce87e](https://github.com/kubeovn/kube-ovn/commit/6db0ce87ebdcb157bc3ad868b4cadeda81248c84) fix lr-lb dnat with multiple distributed gateway ports (#3345)
 * [b95f196ab](https://github.com/kubeovn/kube-ovn/commit/b95f196ab4e52c3e2e76da09f5aa99bd187a3bbb) ci: bump various versions (#3339)
 * [920db9a84](https://github.com/kubeovn/kube-ovn/commit/920db9a84ef06b5aa650dc126faa989fb322f9b5) build(deps): bump github.com/onsi/gomega from 1.28.1 to 1.29.0 (#3343)
 * [f68a3da15](https://github.com/kubeovn/kube-ovn/commit/f68a3da15a6ea04a63f5b61699c32ff315457333) feat:  dpdk-22.11.1 support by kube-ovn (#3332)
 * [42572d2f7](https://github.com/kubeovn/kube-ovn/commit/42572d2f73ed6f129cc0a79bd1597c48ecbc64e3) docs: updated CHANGELOG.md (#3334)
 * [352b161ae](https://github.com/kubeovn/kube-ovn/commit/352b161ae651e9e2325eb4ff5a57c6d0366da16a) bump k8s to v1.28.3 (#3324)
 * [d974d89b4](https://github.com/kubeovn/kube-ovn/commit/d974d89b45b3d44122d5f53f74cb52a6232b41cf) profile: fix fd leak (#3328)
 * [1606f3c8a](https://github.com/kubeovn/kube-ovn/commit/1606f3c8a14a5ff797b966782dca20e98b873ace) Nat reuse router port external ip (#3313)
 * [9a2375393](https://github.com/kubeovn/kube-ovn/commit/9a23753934cc7b11263f12ea1345c6707a76e442) build(deps): bump github.com/onsi/gomega from 1.28.0 to 1.28.1 (#3326)
 * [9a1e662c7](https://github.com/kubeovn/kube-ovn/commit/9a1e662c74af5f903bece162c5b2fab194087683) kube-ovn-controller: fix ovn ic log directory not mounted to hostpath (#3322)
 * [4d74b584e](https://github.com/kubeovn/kube-ovn/commit/4d74b584e93521029d3cd0d23a79e1e69e3fd656) fix golang lint error (#3323)
 * [cc6bd8835](https://github.com/kubeovn/kube-ovn/commit/cc6bd8835edfd697b13df8c6063bf870c6cdd365) build(deps): bump sigs.k8s.io/controller-runtime from 0.16.2 to 0.16.3 (#3315)
 * [7db25b22f](https://github.com/kubeovn/kube-ovn/commit/7db25b22f326f017f1b1a5052e2e3cb13af68276) build(deps): bump the k8s-io group with 1 update (#3314)
 * [3fe3ec665](https://github.com/kubeovn/kube-ovn/commit/3fe3ec665cae5dccd5dfd8dfa156083c3cd58ef0) add load balancer health check (#3216)
 * [4b37fa9b4](https://github.com/kubeovn/kube-ovn/commit/4b37fa9b4aaa674a8fe5cd17c6cc3911e185ddec) build(deps): bump google.golang.org/grpc from 1.58.3 to 1.59.0 (#3310)
 * [2c30a2183](https://github.com/kubeovn/kube-ovn/commit/2c30a2183d6c2f221b208694e26feae9a10eb5a2) build(deps): bump github.com/Microsoft/hcsshim from 0.11.1 to 0.11.2 (#3309)
 * [a94eade9c](https://github.com/kubeovn/kube-ovn/commit/a94eade9c59229f9ec7552f193f8e9fdc150efc1) support vpc configuration of multiple external network segments through label and crd (#3264)
 * [e6086d281](https://github.com/kubeovn/kube-ovn/commit/e6086d281a2aec252683f0b8a1243654703c309c) sync subnet to vpc while switching between custom VPC and default VPC (#3218)
 * [99fb1899b](https://github.com/kubeovn/kube-ovn/commit/99fb1899b4f8ca000d32036538070f6ca73f0aae) security: ignore kubectl cve (#3305)
 * [fd9f3a19b](https://github.com/kubeovn/kube-ovn/commit/fd9f3a19bba7509c747e2532181ba5816006e459) Don't enqueue VPC update when DeletionTimestamp is zero (#3302)
 * [60ef1ad1f](https://github.com/kubeovn/kube-ovn/commit/60ef1ad1f31d2974cb9462670b0ad26b1015be6f) Revert update base image to ubuntu:23.10 (#3289) (#3299)
 * [ce8012093](https://github.com/kubeovn/kube-ovn/commit/ce801209307643525ab3e021c697a3fbdb30b6a5) add base rules for allowing vrrp packets (#3293)
 * [9f93981d6](https://github.com/kubeovn/kube-ovn/commit/9f93981d676813a6d10bfc9b96cf290f8931066d) build(deps): bump google.golang.org/grpc from 1.58.2 to 1.58.3 (#3295)
 * [01df39645](https://github.com/kubeovn/kube-ovn/commit/01df3964521c7a1fb028944d6468f6b67d6ced0b) build(deps): bump golang.org/x/net from 0.16.0 to 0.17.0 (#3296)
 * [b6192f37a](https://github.com/kubeovn/kube-ovn/commit/b6192f37a15d5a35c0c472c301688dced23ea20f) add concurrency limiter to ovs-vsctl (#3288)
 * [286e6340d](https://github.com/kubeovn/kube-ovn/commit/286e6340dcd5da3a01b9b5f67d1daa61b94c1c4e) update base image to ubuntu:23.10 (#3289)
 * [e5af67e75](https://github.com/kubeovn/kube-ovn/commit/e5af67e752dbdbfe969ee3ae7e9034566dda5786) build(deps): bump github.com/onsi/ginkgo/v2 from 2.12.1 to 2.13.0 (#3290)
 * [33bc9a626](https://github.com/kubeovn/kube-ovn/commit/33bc9a626759673ea9b2ffc3aa1fb737636bf008) support custom vpc dns its deployment replicas (#3286)
 * [f9edd6692](https://github.com/kubeovn/kube-ovn/commit/f9edd66926466c83a0a4ff42a18946e414e724e7) webhook: fix ip validation when pod is annotated with an ippool name (#3284)
 * [4af196059](https://github.com/kubeovn/kube-ovn/commit/4af196059283b296b79b434df9d8771eef66199d) webhook: use dedicated port for health probe (#3285)
 * [51722079c](https://github.com/kubeovn/kube-ovn/commit/51722079c9c0abe3f3571b005f352456088d894f) bgp: support multiple peer addresses (#3283)
 * [c493a95bf](https://github.com/kubeovn/kube-ovn/commit/c493a95bf6491ce1cfbe5df81b67063db2087e8b) bgp: support announce policy (#3282)
 * [177d6e1cf](https://github.com/kubeovn/kube-ovn/commit/177d6e1cf901006388c916271f233fed4438033e) bump go modules (#3279)
 * [7ec458958](https://github.com/kubeovn/kube-ovn/commit/7ec4589589d2c57774ae3e136ee1d6114e8229a8) dump cpu/mem profile into file on signal SIGUSR1/SIGUSR2 (#3262)
 * [3f0755d71](https://github.com/kubeovn/kube-ovn/commit/3f0755d710b9b94a3e109d10616c8da15b32c5f7) ovs: load kernel module ip_tables only when it exists (#3281)
 * [5707dd03e](https://github.com/kubeovn/kube-ovn/commit/5707dd03e792b02e4fc2dc48708c8234e1413ff5) kubectl-ko: fix help message (#3277)
 * [ce2c72216](https://github.com/kubeovn/kube-ovn/commit/ce2c722164a9ee3799c5ec931a81a7c771dc7101) update directory name in charts readme (#3276)
 * [7ae309337](https://github.com/kubeovn/kube-ovn/commit/7ae309337afd7fd0c84a1e858402344431ca5040) fix ovn build failure (#3275)
 * [71491211b](https://github.com/kubeovn/kube-ovn/commit/71491211b6f90aec3a5b32b1233ac4c3b70cf8ed) build(deps): bump golang.org/x/mod from 0.12.0 to 0.13.0 (#3272)
 * [9076c3163](https://github.com/kubeovn/kube-ovn/commit/9076c3163dfc4a614cfa4a98d89d66785898229e) build(deps): bump golang.org/x/sys from 0.12.0 to 0.13.0 (#3271)
 * [dac982d1a](https://github.com/kubeovn/kube-ovn/commit/dac982d1a4b5fb4c493113968794790eba2ebddd) build(deps): bump github.com/osrg/gobgp/v3 from 3.18.0 to 3.19.0 (#3270)
 * [6ee71b52b](https://github.com/kubeovn/kube-ovn/commit/6ee71b52b8c2ba04ae1607901777d77faa5dbf8d) build(deps): bump github.com/onsi/gomega from 1.27.10 to 1.28.0 (#3268)
 * [cb46bfdd8](https://github.com/kubeovn/kube-ovn/commit/cb46bfdd8c44e82d9e47ef1094535a82dee24ce6) build(deps): bump github.com/prometheus/client_golang (#3266)
 * [aee41bd68](https://github.com/kubeovn/kube-ovn/commit/aee41bd687c53e7a7c0023e5477e128ca7477101) bump k8s to v1.28.2 (#3159)
 * [f20343297](https://github.com/kubeovn/kube-ovn/commit/f2034329793283457c5ac856fa7b051ff225f488) pinger: increase packet send interval (#3259)
 * [55f254e43](https://github.com/kubeovn/kube-ovn/commit/55f254e43e0058cbcd2fb90f2afac3cae6a24e1a) add init container in vpc-nat-gateway statefulset for init (#3254)
 * [4394bad9e](https://github.com/kubeovn/kube-ovn/commit/4394bad9e5f370ea69e06cc1ab5490d2abf6b81e) lrp should use chassis name instead of uuid (#3258)
 * [4e176a6ae](https://github.com/kubeovn/kube-ovn/commit/4e176a6aeb41ce116586df76a6fad0d4d9f3c7f6) docs: updated CHANGELOG.md (#3252)
 * [e25bd2ed7](https://github.com/kubeovn/kube-ovn/commit/e25bd2ed74c21b1091fe14b8c163ed781cf5d7bd) fix: for existing nic, no need to set the port type to internal (#3243)
 * [0ceb92408](https://github.com/kubeovn/kube-ovn/commit/0ceb924083307117bb00f5429c8827bc74e80a2c) adjust vip prints as ip (#3248)
 * [ff9b57b1c](https://github.com/kubeovn/kube-ovn/commit/ff9b57b1c44209e19c9177609f223acc5ed08361) add dpdk probe (#3151)
 * [b0836845f](https://github.com/kubeovn/kube-ovn/commit/b0836845ffd2fccad0dbe2e38c66826b9cae2893) add wrk in kube-ovn:test (#3233)
 * [6cde93263](https://github.com/kubeovn/kube-ovn/commit/6cde9326357709f8a1cff7d68202605234ced2b1) build(deps): bump google.golang.org/grpc from 1.58.1 to 1.58.2 (#3251)
 * [f016bb37a](https://github.com/kubeovn/kube-ovn/commit/f016bb37a5cf46d02129ff300ccf46c2e5927eef) build(deps): bump github.com/Microsoft/hcsshim from 0.11.0 to 0.11.1 (#3245)
 * [53c291966](https://github.com/kubeovn/kube-ovn/commit/53c2919666a9122bbd6e619a3345774dd7c0ad96) update policy route nexthops para (#3224)
 * [9a408a50a](https://github.com/kubeovn/kube-ovn/commit/9a408a50afc095583a814766b280188a1740ba55) build(deps): bump the k8s-io group with 1 update (#3242)
 * [233e46a25](https://github.com/kubeovn/kube-ovn/commit/233e46a2585b5a20d37c4b40f99072d6a608feb6) fix goproxy Denial of Service vulnerability (#3240)
 * [28127e229](https://github.com/kubeovn/kube-ovn/commit/28127e2290a33ef7d0bef55780f56af32820dce3) build(deps): bump github.com/cyphar/filepath-securejoin (#3239)
 * [2e6a0c9a6](https://github.com/kubeovn/kube-ovn/commit/2e6a0c9a6cca65c71d7bce1ed737ce9051f7365a) build(deps): bump github.com/onsi/ginkgo/v2 from 2.11.0 to 2.12.1 (#3237)
 * [c0b968ea1](https://github.com/kubeovn/kube-ovn/commit/c0b968ea13e1679f1b5fa30ff6d6862a02ff58d3) build(deps): bump github.com/docker/docker (#3234)
 * [5f1d3fe66](https://github.com/kubeovn/kube-ovn/commit/5f1d3fe66c6e59b745d19803b7dd802c7a5ff0a9) build(deps): bump github.com/evanphx/json-patch/v5 from 5.6.0 to 5.7.0 (#3235)
 * [f26f0c357](https://github.com/kubeovn/kube-ovn/commit/f26f0c3572df1da73a961d085d7da4f2ffc8047f) build(deps): bump github.com/osrg/gobgp/v3 from 3.17.0 to 3.18.0 (#3238)
 * [6d5fd3739](https://github.com/kubeovn/kube-ovn/commit/6d5fd3739d01a476c84b1900cb640535fab75b1c) update dependabot
 * [423e099b6](https://github.com/kubeovn/kube-ovn/commit/423e099b64dad206083338e33f230cac91178b36) build(deps): bump google.golang.org/grpc from 1.57.0 to 1.58.1 (#3229)
 * [ae1e0ef99](https://github.com/kubeovn/kube-ovn/commit/ae1e0ef99cb028c4bf21b94717eb94b621244f95) build(deps): bump github.com/Microsoft/hcsshim from 0.10.0 to 0.11.0 (#3228)
 * [ef5b250cd](https://github.com/kubeovn/kube-ovn/commit/ef5b250cd219b1e44658e46f6a7cb459e716a252) build(deps): bump golang.org/x/sys from 0.11.0 to 0.12.0 (#3232)
 * [4a0a7393f](https://github.com/kubeovn/kube-ovn/commit/4a0a7393f89c72afe582c032c246c252517982d8) fix typos (#3215)
 * [eb747615a](https://github.com/kubeovn/kube-ovn/commit/eb747615a8e94a0e7c101081a213e8d9990ca462) chart: remove subnet finalizers before subnets are deleted (#3213)
 * [494f47581](https://github.com/kubeovn/kube-ovn/commit/494f475811f7d0ae0472b3d0e078cd5f03d49e53) build(deps): bump docker/setup-buildx-action from 2 to 3 (#3206)
 * [b04cb3cc7](https://github.com/kubeovn/kube-ovn/commit/b04cb3cc7078a784b5277c1d91d8007254bc4b2a) add special handling for the route policy of the default VPC (#3194)
 * [a17fb7d74](https://github.com/kubeovn/kube-ovn/commit/a17fb7d74c3cd67fcd14614e4497cd1f4f1ef3ff) build(deps): bump docker/setup-qemu-action from 2 to 3 (#3207)
 * [8a0ec964e](https://github.com/kubeovn/kube-ovn/commit/8a0ec964e1158a447bd09a55e07de8a5bb295203) fix add static route to wrong table of ovn (#3195)
 * [4e4f27f29](https://github.com/kubeovn/kube-ovn/commit/4e4f27f2922e9e756e02c88d46da33b760c7164b) kubectl-ko: add new command ovn-trace for tracing ovn lflows only (#3202)
 * [2cba74b39](https://github.com/kubeovn/kube-ovn/commit/2cba74b3914aff8ee582f520b95c4e809585a42d) netpol: fix duplicate default drop acl (#3197)
 * [d766711fb](https://github.com/kubeovn/kube-ovn/commit/d766711fb82b4a24dd081d2ea8a211cd5219ccc5) add log to help find conflict ip owner (#3191)
 * [11408af46](https://github.com/kubeovn/kube-ovn/commit/11408af4660835c8776b29e80890bee0974be543) suuport user custom log location (#3186)
 * [8a3cf0349](https://github.com/kubeovn/kube-ovn/commit/8a3cf034961207f55fe128e4de70cdbcf92de6bf) add golang lint (#3154)
 * [371d7c3c4](https://github.com/kubeovn/kube-ovn/commit/371d7c3c40c471f143060c48d8fa8d4a55d1d844) underlay: fix ip/route tranfer when the nic is managed by NetworkManager (#3184)
 * [8490ae67d](https://github.com/kubeovn/kube-ovn/commit/8490ae67dea0fb5dbad6c5bee380a29fe1926ade) ci: wait for terminating ovs-ovn pod to disappear (#3160)
 * [d8b06ecbe](https://github.com/kubeovn/kube-ovn/commit/d8b06ecbe8cf558b5db788159de8898c830d51b4) fix ovn build (#3166)
 * [089175485](https://github.com/kubeovn/kube-ovn/commit/089175485190c2f9b2a7ab3c61e3005251600a66) build(deps): bump actions/checkout from 3 to 4 (#3183)
 * [f71246f80](https://github.com/kubeovn/kube-ovn/commit/f71246f80d85b403f39df08add7f87aea0168a35) chart: fix ovs-ovn upgrade (#3164)
 * [0505ca338](https://github.com/kubeovn/kube-ovn/commit/0505ca3384929c6b0820abd8961dfc177fbe048b) subnet: fix deleting lr policy on node deletion (#3176)
 * [093730364](https://github.com/kubeovn/kube-ovn/commit/093730364c69acf908dac636905b53ab9941a46e) Add tolerations to pinger daemonSet to allow schedule on not ready nodes. (#3173)
 * [58efd1daf](https://github.com/kubeovn/kube-ovn/commit/58efd1daf2f70af3557d29a7558c06886f6a3b03) ci/test: bump various versions (#3162)
 * [5172b62b7](https://github.com/kubeovn/kube-ovn/commit/5172b62b79ac1e64c633c676de46bfdaa89d68e1) kubectl-ko: get ovn db leaders only on necessary (#3158)
 * [a5c8510cd](https://github.com/kubeovn/kube-ovn/commit/a5c8510cd72633cc1e7163c2d227c679f63ec3a0) underlay: fix NetworkManager operation (#3147)
 * [3a67ea0d5](https://github.com/kubeovn/kube-ovn/commit/3a67ea0d5bd9429b91f130a7f1a581eaf5313a28) enable set --ovn-northd-n-threads (#3150)
 * [c753d4e44](https://github.com/kubeovn/kube-ovn/commit/c753d4e4483c5d94fd36cf5eb970e61ca1676856) build(deps): bump github.com/emicklei/go-restful/v3 (#3156)
 * [4d4577122](https://github.com/kubeovn/kube-ovn/commit/4d457712260ac92d5fcf7e2cc5314d5bbe1fdb55) Fix max unavailable (#3149)
 * [46cee0374](https://github.com/kubeovn/kube-ovn/commit/46cee0374a80e3b9e431f64b4277e0827c0e335d) base: remove ovn patch for skipping ct (#3141)
 * [bcbae34eb](https://github.com/kubeovn/kube-ovn/commit/bcbae34ebaeecbadb91b020da8214a3d5b069ebc) sbctl chassis operation  replace with libovsdb (#3119)
 * [a10d94aeb](https://github.com/kubeovn/kube-ovn/commit/a10d94aeb387c622d4629081b694af4a62f5e23f) Enable set probe (#3145)
 * [d6e5aba00](https://github.com/kubeovn/kube-ovn/commit/d6e5aba009990dc71aee7cc18c1bdf7520073ac0) support recreate a backup pod with full annotation (#3144)
 * [2f6f67515](https://github.com/kubeovn/kube-ovn/commit/2f6f67515dab6f2aa1c657c87e4af0e6ffce5614) fix ovn nat not clean (#3139)
 * [d0063e3cc](https://github.com/kubeovn/kube-ovn/commit/d0063e3cca6597544f17fc3d930bc881d1263000) add mulicast snoop switch (#3129)
 * [e541a5797](https://github.com/kubeovn/kube-ovn/commit/e541a5797aeb32ffa41efe3bbefc68c07da98b80) ovn: do not send direct traffic between lports to conntrack (#3131)
 * [563bca688](https://github.com/kubeovn/kube-ovn/commit/563bca6881e041df1d2a2d36b8efbe44d3fe5d5a) bump go version to 1.21 (#3137)
 * [bbb0e102d](https://github.com/kubeovn/kube-ovn/commit/bbb0e102d202b6766c6a01366cac6afde5c3c207) delete append externalIds process in initIPAM (#3134)
 * [01c5d362a](https://github.com/kubeovn/kube-ovn/commit/01c5d362ad494f49fb4830da923ee2ddc7fca4b1) add probe (#3133)
 * [51ceddb66](https://github.com/kubeovn/kube-ovn/commit/51ceddb66362fd35c78c575b24a80c0afbb91ec0) ci: fix version comparison (#3127)
 * [795e026ad](https://github.com/kubeovn/kube-ovn/commit/795e026ad14f1ed6cf11cce30016c3828c5e4aea) fix ci (#3125)
 * [294e5f996](https://github.com/kubeovn/kube-ovn/commit/294e5f9969557a6b041ecae20c814c183f8ea555) move unnecessary init process after startWorkers (#3124)
 * [1a720192d](https://github.com/kubeovn/kube-ovn/commit/1a720192ddbf3639bd5371a4ffbd604405da3cb3) add e2e test for ovn db recover (#3118)
 * [99f37b398](https://github.com/kubeovn/kube-ovn/commit/99f37b3989ad3ab0446db58c79d567d9a8e4d8ee) prepare for next release
 * [e1a1b78b0](https://github.com/kubeovn/kube-ovn/commit/e1a1b78b0c3dc6314a44ff69d9bfc9ff7cc4dc7f) docs: updated CHANGELOG.md (#3122)
 * [15861418b](https://github.com/kubeovn/kube-ovn/commit/15861418b017b9573d71d970882073c08aa0a391) update changelog
 * [6cf53101f](https://github.com/kubeovn/kube-ovn/commit/6cf53101f620caafd7d08b424e1408e967adaca3) build(deps): bump sigs.k8s.io/controller-runtime from 0.15.0 to 0.15.1 (#3120)
 * [cd1202ca4](https://github.com/kubeovn/kube-ovn/commit/cd1202ca400a64bb5c3d68784eeb23586444d405) ovn: fix corrupted database file on start (#3112)
 * [02f8c630b](https://github.com/kubeovn/kube-ovn/commit/02f8c630b076e9a97cd8788714813a5c1858f634) some fixes in e2e (#3116)
 * [d8fa8395e](https://github.com/kubeovn/kube-ovn/commit/d8fa8395e4150bb6c31ec0dce2b447d0b8461c6d) controller: fix vpc update (#3117)
 * [b5b25ffd0](https://github.com/kubeovn/kube-ovn/commit/b5b25ffd086eed2d2f72cc32b3a5684cf002e6ea) increase event burst size (#3115)
 * [c8031f6ef](https://github.com/kubeovn/kube-ovn/commit/c8031f6ef60d25095a7a54b532aa6f0871cc39c0) build(deps): bump golang.org/x/sys from 0.10.0 to 0.11.0 (#3114)
 * [6ba997d20](https://github.com/kubeovn/kube-ovn/commit/6ba997d20b32555e02906a635ef00ce980cbb52b) 简化 ovn eip 类型 (#3107)
 * [a0c5e3895](https://github.com/kubeovn/kube-ovn/commit/a0c5e3895907bc9100a31d2d7d7039b2eed7a664) fix u2o policy route allocate too many openflows cause oom (#3099)
 * [a9fdbf92e](https://github.com/kubeovn/kube-ovn/commit/a9fdbf92ea09a2cfff184ecae4c72c632f0d213b) Fix relevant annotations are not deleted in hotnoplug nic process (#3108)
 * [3c6d6bc05](https://github.com/kubeovn/kube-ovn/commit/3c6d6bc051bc9a88dd1c55bec14685d42e99cd11) ovn: delete the db file if the node with new empty db file cannot join cluster for more than 120s (#3101)
 * [914bf6130](https://github.com/kubeovn/kube-ovn/commit/914bf6130e4f65bc58458e9dd61d72aa495f85f8) get all chassis once (#3103)
 * [42e0574c2](https://github.com/kubeovn/kube-ovn/commit/42e0574c2eb820ea8b1265dd633687ab92de5941) distinguish nat ip for central subnet with ecmp and active-standby (#3100)
 * [a27ce4c33](https://github.com/kubeovn/kube-ovn/commit/a27ce4c334db6d49d6c70dc8375667b5b9ead019) build(deps): bump github.com/osrg/gobgp/v3 from 3.16.0 to 3.17.0 (#3105)
 * [68dc1c380](https://github.com/kubeovn/kube-ovn/commit/68dc1c38023e4967e9644c89e72e63c605df8591) add log near err (#3098)
 * [c6c472a0e](https://github.com/kubeovn/kube-ovn/commit/c6c472a0e295dfd4189df9fd8cd47a69827fdb04) iptables: reject access to invalid service port when kube-proxy works in IPVS mode (#3059)
 * [f8835ef56](https://github.com/kubeovn/kube-ovn/commit/f8835ef5674af736136713e8867af9a59f173ff7) Ovn nat 1 (#3095)
 * [5704dae0c](https://github.com/kubeovn/kube-ovn/commit/5704dae0c7be8a86d70cd81f03602e7e5942b703) skip ok pod (#3090)
 * [18580edff](https://github.com/kubeovn/kube-ovn/commit/18580edffa990a65b61e3e49d356e0a4c603721e) ipam: return error for invalid ip range (#3088)
 * [a7e7a83d1](https://github.com/kubeovn/kube-ovn/commit/a7e7a83d12542f23a42e8f48ceef4a7c6b8fd9c7) some fixes in e2e (#3094)
 * [882187431](https://github.com/kubeovn/kube-ovn/commit/88218743189db426fd6d97272a670762b71ba1ae) bug_fix if only one port bind to the sg, then unbind the port to the sg ,it will not enforce in port_group (#3092)
 * [4c1161e92](https://github.com/kubeovn/kube-ovn/commit/4c1161e92f5d2ee2aecc4b4f88a4f9fc49a1e7fe) fix .status.default when initializing the default vpc (#3086)
 * [fa91428bc](https://github.com/kubeovn/kube-ovn/commit/fa91428bc2882e9d9dd82b78f6eb3dcc91500cd4) fix repeate set chassis (#3083)
 * [68a798f47](https://github.com/kubeovn/kube-ovn/commit/68a798f47aa4a9051a41dd7020d11ab1eda34910) build(deps): bump google.golang.org/grpc from 1.56.2 to 1.57.0 (#3085)
 * [da1648ccd](https://github.com/kubeovn/kube-ovn/commit/da1648ccd64c3436669d643e2d7a1b474d1d7cd9) fix go fmt
 * [88b5912fc](https://github.com/kubeovn/kube-ovn/commit/88b5912fce998957704fb689340911c213a8e902) fix kube-ovn-speaker log (#3081)
 * [89544c352](https://github.com/kubeovn/kube-ovn/commit/89544c35233fadb34ce8e0236d712ccbd6127b81) remove FOSSA status card
 * [ac648680e](https://github.com/kubeovn/kube-ovn/commit/ac648680ebc6b2655903bae5cb50c09ae90b0c12) cni-server: fix ovn mappings for vpc nat gateway (#3075)
 * [0fe9dcb27](https://github.com/kubeovn/kube-ovn/commit/0fe9dcb276e5dd9f7f0ca31b568cb40d3aabc49a) fix kube-ovn-speaker (#3076)
 * [827a5a273](https://github.com/kubeovn/kube-ovn/commit/827a5a2730f969477b60ae466e5aa4ca175073c1) build(deps): bump github.com/Microsoft/hcsshim from 0.9.10 to 0.10.0 (#3079)
 * [38cd9203a](https://github.com/kubeovn/kube-ovn/commit/38cd9203a56bd65f3a596775fa0a4f3ca75173e7) ovn client: fix sb chassis existence check (#3072)
 * [038ff7def](https://github.com/kubeovn/kube-ovn/commit/038ff7defef82eac3ab35fae55f034e0279d9e6c) e2e: fix switch lb rule test (#3071)
 * [e14ebbd58](https://github.com/kubeovn/kube-ovn/commit/e14ebbd58af2ef0599b6bfaa3dd171d6cc0388d0) bump github.com/docker/docker to v24.0.5 (#3073)
 * [90c830574](https://github.com/kubeovn/kube-ovn/commit/90c830574c9f7f6786201e4fb3ce2951accff7e9) iptables: add --random-fully to SNAT rules (#3066)
 * [1350397ec](https://github.com/kubeovn/kube-ovn/commit/1350397ec97c51f356be168f89f2fe49fd5511ab) update lint tmeout
 * [ff6d03d0e](https://github.com/kubeovn/kube-ovn/commit/ff6d03d0ec22260490a234b55bd775243a9c1f9d) build(deps): bump github.com/onsi/gomega from 1.27.9 to 1.27.10 (#3069)
 * [76b01496f](https://github.com/kubeovn/kube-ovn/commit/76b01496f3d5d4e03ee54b31a0ac1bf94ed3b6fa) bump k8s to v1.27.4 (#3063)
 * [d8e59ab61](https://github.com/kubeovn/kube-ovn/commit/d8e59ab619ec097a25ef0fbe26122ede73d4fd91) e2e: do not import pkg/daemon (#3055)
 * [20a6526c4](https://github.com/kubeovn/kube-ovn/commit/20a6526c407b5dbd69330c119681efb2a66903ca) build(deps): bump github.com/onsi/gomega from 1.27.8 to 1.27.9 (#3065)
 * [af335ba81](https://github.com/kubeovn/kube-ovn/commit/af335ba81a9789fde3fa9b7a4a41146ba00fd114) build(deps): bump github.com/Microsoft/hcsshim from 0.9.9 to 0.9.10 (#3061)
 * [976a32b01](https://github.com/kubeovn/kube-ovn/commit/976a32b01a84824f615211ad20ae94833c5b576b) ci: fix multus installation (#3062)
 * [0d1599ffd](https://github.com/kubeovn/kube-ovn/commit/0d1599ffdcca1428647f63e44db396aabf8081e2) add srl connectivity test (#3056)
 * [42f35a358](https://github.com/kubeovn/kube-ovn/commit/42f35a3588e90986668b0100d99424862d1c6ada) ipam: fix ippool with single dual-stack address (#3054)
 * [2ba3b8e54](https://github.com/kubeovn/kube-ovn/commit/2ba3b8e54cd06bc3b66c029d31ba9a8744fd4f61) controller: skip VIP gc if LB not found (#3048)
 * [52232b5e3](https://github.com/kubeovn/kube-ovn/commit/52232b5e33f8bd678d516602d8ef81d364c9757b) keep vm vip when enableKeepVmIP is true (#3053)
 * [ed58b2103](https://github.com/kubeovn/kube-ovn/commit/ed58b210302269f2d09b27087001b74fe50cd131) cni: reduce memory usage (#3047)
 * [3be1e84c2](https://github.com/kubeovn/kube-ovn/commit/3be1e84c2142dd89412d13a4896b40a7f91c0296) set genev_sys_6081 tx checksum off (#3045)
 * [4e761156c](https://github.com/kubeovn/kube-ovn/commit/4e761156c1ede9e44e85dbf2a03aca4f816a8712) fix vpc lb init (#3046)
 * [f4f80415a](https://github.com/kubeovn/kube-ovn/commit/f4f80415a06fa73d030e1becfb1a9d479eafdca6) custom vpc pod support tcp http probe with tproxy method (#3024)
 * [494209d3c](https://github.com/kubeovn/kube-ovn/commit/494209d3cf3e9177c2d311d6b7c063e67810393a) change log (#3042)
 * [b40c35b8c](https://github.com/kubeovn/kube-ovn/commit/b40c35b8c4076622886f63c6dae344764cc2a05d) Makefile: add deepflow and kwok installation (#3036)
 * [5a0686b27](https://github.com/kubeovn/kube-ovn/commit/5a0686b2700070d04fd31d50ffea1e3367202b50) windows: fix ovn patches (#3035)
 * [e3b7439d4](https://github.com/kubeovn/kube-ovn/commit/e3b7439d40c41b88525d87fb66a7cf1c03d2881a) ci: pin go version to 1.20.5 (#3034)
 * [97a3e1bc9](https://github.com/kubeovn/kube-ovn/commit/97a3e1bc9b192a76a11e338b609876d2d8c77165) static ip in exclude-ips can be allocated normally when subnet's availableIPs is 0 (#3031)
 * [9d88e4970](https://github.com/kubeovn/kube-ovn/commit/9d88e4970718cc377914470d9059a63c548ed900) pinger: use fully qualified domain name (#3032)
 * [f3833f91b](https://github.com/kubeovn/kube-ovn/commit/f3833f91b42d46692f618ec0b3f4976e95b85739) feat:  suport kubevirt nic hotplug (#3013)
 * [62f332897](https://github.com/kubeovn/kube-ovn/commit/62f332897a26f014031f7dbab1a8566f6d28fd22) fix lrp eip not clean (#3026)
 * [047af4a21](https://github.com/kubeovn/kube-ovn/commit/047af4a21607125b43b15a6f8ea6b213b931b1bf) build(deps): bump helm/kind-action from 1.7.0 to 1.8.0 (#3029)
 * [e01e616ee](https://github.com/kubeovn/kube-ovn/commit/e01e616ee0be12c54ab1ac54b03df9f002c55456) update maintainer
 * [ea9c1f1e5](https://github.com/kubeovn/kube-ovn/commit/ea9c1f1e5dad1d4342d851d37315631b696d6048) uninstall.sh: fix ipset name (#3028)
 * [9e6dc6366](https://github.com/kubeovn/kube-ovn/commit/9e6dc6366091c443bba60a0bc46fbedda29b0636) build(deps): bump github.com/docker/docker (#3027)
 * [3dd7f4ab7](https://github.com/kubeovn/kube-ovn/commit/3dd7f4ab788797831444148c4485869f61aca5b7) replace ovn legacy client with libovsdb (#3018)
 * [c5bfdb464](https://github.com/kubeovn/kube-ovn/commit/c5bfdb4645a7e2899c8b38ae51d01ff5c87ce2b2) install.sh: fix duplicate resources apply (#3023)
 * [2e4fb05cf](https://github.com/kubeovn/kube-ovn/commit/2e4fb05cf015751439115426df38ca1d7bad06f0) build(deps): bump github.com/docker/docker (#3019)
 * [aefaef5ab](https://github.com/kubeovn/kube-ovn/commit/aefaef5abeb6d6120fc5948e3112532cfa305e54) build(deps): bump google.golang.org/grpc from 1.56.1 to 1.56.2 (#3020)
 * [1f1fb82e4](https://github.com/kubeovn/kube-ovn/commit/1f1fb82e4fcaa4dd22d65d54843c79f2e26b09bc) ovn: fix cluster connections when SSL is enabled (#3001)
 * [18560c964](https://github.com/kubeovn/kube-ovn/commit/18560c964caf4b4a7f007f812c76f8a885665e34) cleanup.sh: wait for provier-networks to be deleted before deleting kube-ovn-cni (#3006)
 * [9304ae5a4](https://github.com/kubeovn/kube-ovn/commit/9304ae5a460f8dbff0e785aa703adef899de4109) kube-ovn-controller: fix workqueue metrics (#3011)
 * [d4153885e](https://github.com/kubeovn/kube-ovn/commit/d4153885e06248b42dbbab9c4b1e8aaf17dd4b0e) ci: fix go cache key (#3015)
 * [5a5f66ebf](https://github.com/kubeovn/kube-ovn/commit/5a5f66ebf55f38a0c4dfedf2de753dc6fe4300af) fix vlan subnet use logical gw can not access outside cluster node (#3007)
 * [18fd55ddf](https://github.com/kubeovn/kube-ovn/commit/18fd55ddf4a45c350a2561a4b2976999cc5b86c4) build(deps): bump github.com/prometheus-community/pro-bing (#3016)
 * [269f460d7](https://github.com/kubeovn/kube-ovn/commit/269f460d72167ba5018a94e95788ccac9588ba37) fix vpc already delete while delete policy route (#3005)
 * [e744d76ea](https://github.com/kubeovn/kube-ovn/commit/e744d76eac8d366b836abd0d7a39c4d5cbd127c2) make compatible with simplicified enable-eip-snat-cm (#3009)
 * [2a6525300](https://github.com/kubeovn/kube-ovn/commit/2a6525300c996afa2cdede47fc6f36398115244c) build(deps): bump golang.org/x/sys from 0.9.0 to 0.10.0 (#3012)
 * [d5f89bce7](https://github.com/kubeovn/kube-ovn/commit/d5f89bce73841b76993c43a94820e9c04fd0fe52) subnet: fix nat outgoing policy rule (#3003)
 * [8358a91e8](https://github.com/kubeovn/kube-ovn/commit/8358a91e8408c133adf20cfc4a889931a890184a) build(deps): bump github.com/osrg/gobgp/v3 from 3.15.0 to 3.16.0 (#3010)
 * [fe924e9f4](https://github.com/kubeovn/kube-ovn/commit/fe924e9f49cf00e419d943fa8453319d4d87d0ad) fix subnet finalizer (#3004)
 * [12366937b](https://github.com/kubeovn/kube-ovn/commit/12366937b1177496b2469d8e87bc67abe4daa4b2) chart: fix readOnly in volumes (#3002)
 * [d5462b10c](https://github.com/kubeovn/kube-ovn/commit/d5462b10c9d26dc99e287b5cc111e5c4f94cc57b) libovsdb: various bug fixes (#2998)
 * [af04530e9](https://github.com/kubeovn/kube-ovn/commit/af04530e9aa4db98ef25ac4f0a681dfcb3be0fd4) choose subnet by pod's annotation in networkpolicy (#2987)
 * [5c455499b](https://github.com/kubeovn/kube-ovn/commit/5c455499b35d91f2c976520927bc70822f727a65) IPPool: fix missing support for CIDR (#2982)
 * [f2d063a85](https://github.com/kubeovn/kube-ovn/commit/f2d063a853938a81c5be8ab3ce762a6b4e35d99e) kubectl ko performance enhance (#2975)
 * [d5d196e7c](https://github.com/kubeovn/kube-ovn/commit/d5d196e7c1fdc39ab50ce0fd87c5c5340f00d5d3) fix deleting old sb chassis for a re-added node (#2989)
 * [30cd09e6e](https://github.com/kubeovn/kube-ovn/commit/30cd09e6e792ca0cf648614155fae423379c9644) add e2e for new ippool feature (#2981)
 * [5fdf1f9ee](https://github.com/kubeovn/kube-ovn/commit/5fdf1f9ee08955d1c29b1079c3cbb47c62ffbd8a) underlay: fix NetworkManager syncer for virtual interfaces (#2988)
 * [1bb51239c](https://github.com/kubeovn/kube-ovn/commit/1bb51239cf8e26055cf230315edb0e09f7bc7792) underlay: does not set a device managed to no if it has VLAN managed by NM (#2986)
 * [3793e9937](https://github.com/kubeovn/kube-ovn/commit/3793e99370206a410743e1c96839e11ad0647b4f) build(deps): bump google.golang.org/protobuf from 1.30.0 to 1.31.0 (#2985)
 * [6a5bfe46d](https://github.com/kubeovn/kube-ovn/commit/6a5bfe46dd6b2e9124e005386191672e4b82f3e6) support helm install  hybrid_dpdk ovs-ovn (#2980)
 * [dc40a8cb6](https://github.com/kubeovn/kube-ovn/commit/dc40a8cb6df0d3fbba171069c1484e3c87ae9157) add unittest for IPAM (#2977)
 * [daa436d3c](https://github.com/kubeovn/kube-ovn/commit/daa436d3c1d4ea099faff713be0affb34ff121cb) IPAM: fix subnet mutex not released when static IP is out of range (#2979)
 * [65fd8a4aa](https://github.com/kubeovn/kube-ovn/commit/65fd8a4aa4d596363131ba710a50bfbc052e4071) fix initialization check of vpc nat gateway configuration (#2978)
 * [e558702de](https://github.com/kubeovn/kube-ovn/commit/e558702de597f99eb0591cb8e2bf31fd084cfd0b) refactor: make qos test cases parallel (#2957)
 * [27685a173](https://github.com/kubeovn/kube-ovn/commit/27685a173a7f537e4ffd1539a3a655c8354ee59d) IPAM: add support for ippool (#2958)
 * [41b4f266d](https://github.com/kubeovn/kube-ovn/commit/41b4f266dc565538b1711e7dd876cfd34c8ab894) build(deps): bump google.golang.org/grpc from 1.56.0 to 1.56.1 (#2974)
 * [57b01b4ad](https://github.com/kubeovn/kube-ovn/commit/57b01b4ad210d1c3845867529c211ba1b2960e2b) ovn ic support dual (#2970)
 * [7a14cf219](https://github.com/kubeovn/kube-ovn/commit/7a14cf21958a50e9c5b40fa1fe836294b64f7bee) base: fix ovn patches (#2971)
 * [a5194e668](https://github.com/kubeovn/kube-ovn/commit/a5194e6682c9b2603b7857024193d4ea59efec87) build(deps): bump github.com/onsi/ginkgo/v2 from 2.10.0 to 2.11.0 (#2968)
 * [a5e63c720](https://github.com/kubeovn/kube-ovn/commit/a5e63c7209c61491179d7b272e87e37592f0b3f4) add detail comment (#2969)
 * [13256fab7](https://github.com/kubeovn/kube-ovn/commit/13256fab7002db60de4f2df0c89bb986ef4f9dc8) 1. add host multicast perf (#2965)
 * [33b6df120](https://github.com/kubeovn/kube-ovn/commit/33b6df1202ce53c02a8707895113f22beae2b1ca) cni-server: reconcile ovn0 routes periodically (#2963)
 * [e4f682675](https://github.com/kubeovn/kube-ovn/commit/e4f682675319eb80289bbf68a008619278ff13f6) uninstall.sh: flush and delete iptables chain OVN-MASQUERADE (#2961)
 * [9fbebd340](https://github.com/kubeovn/kube-ovn/commit/9fbebd3406cacf4f10c3afbfc55b25b407e84f90) fix e2e failed (#2960)
 * [5de251721](https://github.com/kubeovn/kube-ovn/commit/5de251721684e01cf50511c9ea59faa4a0e6d60c) u2o specify u2oip from v1.9 (#2934)
 * [30ea6d6c4](https://github.com/kubeovn/kube-ovn/commit/30ea6d6c4679b1b206aec0a9e8bb55b1b25fe226) underlay: sync NetworkManager IP config to OVS bridge (#2949)
 * [27a2f300c](https://github.com/kubeovn/kube-ovn/commit/27a2f300c918fccc588728e6fe07b02253b4f5ca) chore: USERS.md (#2955)
 * [1c29580e9](https://github.com/kubeovn/kube-ovn/commit/1c29580e9ee1e52943956cd96ed746adcf51a07e) bump k8s version to v1.27.3 (#2953)
 * [c0730acb7](https://github.com/kubeovn/kube-ovn/commit/c0730acb7c998c8257714c8dd9253628d36bb92a) ci: fix build-base strategy (#2950)
 * [f52b15091](https://github.com/kubeovn/kube-ovn/commit/f52b15091d21a18a0e1928c55788f035f206edf7) e2e: add qos policy test cases (#2924)
 * [d8739d29b](https://github.com/kubeovn/kube-ovn/commit/d8739d29bdead176da5d4594fdebc7a3f836ec53) typo (#2952)
 * [de9c96036](https://github.com/kubeovn/kube-ovn/commit/de9c9603679bcb2957becaca1f78e54c3a006168) build(deps): bump google.golang.org/grpc from 1.55.0 to 1.56.0 (#2951)
 * [f4b3c0fe1](https://github.com/kubeovn/kube-ovn/commit/f4b3c0fe190676237d4ca3c6b484e4be6e91c783) build(deps): bump github.com/prometheus/client_golang (#2948)
 * [e13c2005d](https://github.com/kubeovn/kube-ovn/commit/e13c2005d65a8c7040fd873a640968daaf37c304) Revert "nm not managed only in the change provide nic name case (#2754)" (#2944)
 * [765dc8d7f](https://github.com/kubeovn/kube-ovn/commit/765dc8d7f74560e70d7b41c40576a4150b1fbffa) add permision for test-server.sh (#2942)
 * [a9d0b4be5](https://github.com/kubeovn/kube-ovn/commit/a9d0b4be58d631d5871d7df122d7fa27d8bc134a) Kubectl ko diagnose perf (#2915)
 * [8f414f724](https://github.com/kubeovn/kube-ovn/commit/8f414f724911c00d4722c7ff7dc0e127540301cd) build(deps): bump golang.org/x/sys from 0.8.0 to 0.9.0 (#2940)
 * [88f706e41](https://github.com/kubeovn/kube-ovn/commit/88f706e412dc0392224f9ef47a34e05a08061142) controller: fix DHCP MTU when the default network mode is underlay (#2941)
 * [ea56b5603](https://github.com/kubeovn/kube-ovn/commit/ea56b56034b2a3e9a5de150765e6e98f02c3b686) e2e: fix u2o case (#2931)
 * [c1c716f18](https://github.com/kubeovn/kube-ovn/commit/c1c716f1834c8b174732136c438153cfe5ecc235) add err log to help find conflict ip owner (#2939)
 * [1f27076f7](https://github.com/kubeovn/kube-ovn/commit/1f27076f7c8bb1b6b1452dab4bd00438a5258d85) support set the mtu of  dhcpv4_options (#2930)
 * [f1d2011a9](https://github.com/kubeovn/kube-ovn/commit/f1d2011a9b95c4c94bb9993b967af4355df13084) modify lb-svc dnat port error (#2927)
 * [d7edac78a](https://github.com/kubeovn/kube-ovn/commit/d7edac78a115e6ec8505f3729e344c48845cd885) fix race condition in gateway check logs (#2928)
 * [fc7c16ae7](https://github.com/kubeovn/kube-ovn/commit/fc7c16ae7d006e6a951eb051cba3af4461598296) add subnet.spec.u2oInterconnectionIP (#2921)
 * [6105d57b4](https://github.com/kubeovn/kube-ovn/commit/6105d57b404555ad2a1345a6139221ed4d40cfe3) disable ai review
 * [8773ea3d4](https://github.com/kubeovn/kube-ovn/commit/8773ea3d460351eacadf5e6d70dba8cd815e48be) e2e: fix waiting deployment to be restarted (#2909)
 * [56927913e](https://github.com/kubeovn/kube-ovn/commit/56927913e2b227e1c5850cb6b3d83e1cbb0e4c4c) make conformance with underlay  pn vlan subnet has no gw (#2908)
 * [0356a63fe](https://github.com/kubeovn/kube-ovn/commit/0356a63fe092bb7a33bdabb4d92016d2076b7c49) fix: natgw init check command not work (#2923)
 * [3a8e13ee3](https://github.com/kubeovn/kube-ovn/commit/3a8e13ee3b137906d0a7430e88017ef772a1cced) fix issue 2916 (#2917)
 * [517d3791a](https://github.com/kubeovn/kube-ovn/commit/517d3791a9363dd3635ca37b40cbd0cdd2906931) add sync map to fix cocurrent write (#2918)
 * [dff950b18](https://github.com/kubeovn/kube-ovn/commit/dff950b18c790820354f1b22c52aeda5d4cf8929) cni-server: clear iptables mark before doing masquerade (#2919)
 * [d043a2d2e](https://github.com/kubeovn/kube-ovn/commit/d043a2d2edfb18a874f5c73b2354c2e8dd40d539) build(deps): bump github.com/onsi/ginkgo/v2 from 2.9.7 to 2.10.0 (#2913)
 * [525b0b76d](https://github.com/kubeovn/kube-ovn/commit/525b0b76d4ddad8d248bda0b845ecbc21fb944a1) build(deps): bump github.com/onsi/gomega from 1.27.7 to 1.27.8 (#2914)
 * [3616d3de0](https://github.com/kubeovn/kube-ovn/commit/3616d3de045d9a812a89c43d6af7c18a83fe11c3) For eip created without spec.V4ip this field (#2912)
 * [ace0b970b](https://github.com/kubeovn/kube-ovn/commit/ace0b970b22c42a342b46d94d3bdd613dcb8e45f) match outgoing interface when perform snat (#2911)
 * [d61a2ad64](https://github.com/kubeovn/kube-ovn/commit/d61a2ad64d8285ef3c68956facb047b8288b5529) libovsdb: ignore not found error when listing objects with a filter (#2900)
 * [78f923a9f](https://github.com/kubeovn/kube-ovn/commit/78f923a9f21301705a202580ecdd4090a2c02513) build(deps): bump github.com/sirupsen/logrus from 1.9.2 to 1.9.3 (#2903)
 * [0e27e0caa](https://github.com/kubeovn/kube-ovn/commit/0e27e0caa26767e38b395baefe65b6e0dcf044a1) build(deps): bump github.com/osrg/gobgp/v3 from 3.14.0 to 3.15.0 (#2904)
 * [fd92c2a87](https://github.com/kubeovn/kube-ovn/commit/fd92c2a87ec6979f306104d7bafb4ecf0643862e) fix base build
 * [668287af8](https://github.com/kubeovn/kube-ovn/commit/668287af888224a48438de9315b011d1010f465d) fix build base ci
 * [a17461409](https://github.com/kubeovn/kube-ovn/commit/a1746140969f2fbe7d57b0f8a6fe0aff307fe5f0) fix build base ci
 * [2f52d9297](https://github.com/kubeovn/kube-ovn/commit/2f52d9297eef625e9a85fabd8a06642e13d99763) refactor IPAM (#2896)
 * [db51370fd](https://github.com/kubeovn/kube-ovn/commit/db51370fdc4ef8ee5a61732893878208b74a926c) add e2e u2o vpc version check (#2901)
 * [6acecb604](https://github.com/kubeovn/kube-ovn/commit/6acecb6046b8ee929582c69831593f34f476bfff) kube-ovn-controller: fix subnet update (#2882)
 * [35aa8b403](https://github.com/kubeovn/kube-ovn/commit/35aa8b4035800e620f7a699c1759bf8a09cdeb22) Supporting user-defined kubelet directory (#2893)
 * [3883a744c](https://github.com/kubeovn/kube-ovn/commit/3883a744c7f799b4c2f40abd260d4ff74854a96a) ci: use latest golangci-lint
 * [efe3ee34f](https://github.com/kubeovn/kube-ovn/commit/efe3ee34fc4ab30b01aa4349502f774ce6165bb1) underlay: do not delete patch ports created by ovn-controller (#2851)
 * [04c64f0ae](https://github.com/kubeovn/kube-ovn/commit/04c64f0aee60ed2d988460fe2986db03c61c7b61) update pr-review
 * [aa1ffaa3a](https://github.com/kubeovn/kube-ovn/commit/aa1ffaa3a5fc349c623e7e51ab9d0c5e1655ee3c) auto build base for release branches
 * [fe4eec778](https://github.com/kubeovn/kube-ovn/commit/fe4eec7783655cc6da120a06eecf6bad194205d2) Add natoutgoing policy rules (#2883)
 * [bbe04e89d](https://github.com/kubeovn/kube-ovn/commit/bbe04e89dac9b69c452f95c3a036b9b94bae19ba) pin golangci-lint version
 * [0c5f90553](https://github.com/kubeovn/kube-ovn/commit/0c5f90553f65aa083fe2cf0853ad8a4d54113c36) skip case 'connect to NodePort service with external traffic policy set to Local from other nodes' (#2895)
 * [93f027f3e](https://github.com/kubeovn/kube-ovn/commit/93f027f3e602b56db6934e1e344a3f8363a61110) refactor subnet gateway (#2872)
 * [33c526237](https://github.com/kubeovn/kube-ovn/commit/33c52623700d0284a5701ee4b4f9446a3fef00bc) update webhook check (#2878)
 * [a123be78e](https://github.com/kubeovn/kube-ovn/commit/a123be78eb0f37ee686d41fa1a3ea8d67eff93ce) skip pr-review as run out openai quota
 * [5c2c94873](https://github.com/kubeovn/kube-ovn/commit/5c2c948731b8cba961c95fa3f29d115bc8b1565e) skip kubectl cve
 * [403c2dcd7](https://github.com/kubeovn/kube-ovn/commit/403c2dcd7a70cc46d36dcf9dff6763b4d9701417) build(deps): bump github.com/onsi/ginkgo/v2 from 2.9.5 to 2.9.7 (#2890)
 * [589d0b6f8](https://github.com/kubeovn/kube-ovn/commit/589d0b6f8221953d9d3a4f749d6f4a6113cb831e) e2e: multiple external network (#2884)
 * [79521c318](https://github.com/kubeovn/kube-ovn/commit/79521c3185b988892f0e2cc7aea360b8c59ba9af) build(deps): bump github.com/stretchr/testify from 1.8.3 to 1.8.4 (#2885)
 * [71253fe5f](https://github.com/kubeovn/kube-ovn/commit/71253fe5fa925aec2d280f76f5aab4307e5829a1) fix vip str format (#2879)
 * [6b5345ff2](https://github.com/kubeovn/kube-ovn/commit/6b5345ff2aae3a6c6072593d444674bb0752909b) ci: fix valgrind result analysis (#2853)
 * [7c80a1358](https://github.com/kubeovn/kube-ovn/commit/7c80a13580d845c3d7b43366156f19902e185e12) ovs: fix memory leak in qos (#2871)
 * [9f39621a8](https://github.com/kubeovn/kube-ovn/commit/9f39621a8d57627e4b45cba5058d91b79e7c61e7) feat: vpc nat gw e2e (#2866)
 * [e68b983b1](https://github.com/kubeovn/kube-ovn/commit/e68b983b1041f01dcdafa72184c8725fbab20ce6) build(deps): bump github.com/docker/docker (#2875)
 * [056b4cf83](https://github.com/kubeovn/kube-ovn/commit/056b4cf83366a021da6fae13611b6969ef1fe40f) fix gc nil pointer (#2858)
 * [32b85219c](https://github.com/kubeovn/kube-ovn/commit/32b85219c3b55d33041551ac7d7bc2803a372f63) bump k8s to v1.27.2 (#2861)
 * [a80b3754c](https://github.com/kubeovn/kube-ovn/commit/a80b3754c39e02e01f42fe8c92eeba68a7ab29a7) add e2e test for slr (#2841)
 * [20b20366d](https://github.com/kubeovn/kube-ovn/commit/20b20366da888c15aa25730cb5eff205c5b63843) Move docs to new website (#2862)
 * [24d9dfeef](https://github.com/kubeovn/kube-ovn/commit/24d9dfeef0ceaf805326226b762ffa5407154b99) build(deps): bump gopkg.in/k8snetworkplumbingwg/multus-cni.v4 (#2860)
 * [83a47a734](https://github.com/kubeovn/kube-ovn/commit/83a47a734ff0f92f741a644b2424de73573716ed) update dependabot
 * [d6202bc7a](https://github.com/kubeovn/kube-ovn/commit/d6202bc7addf94abb9b89e42d060dbd770b42e66) refactor clusterrole for kube-ovn (#2833)
 * [b1c77ad73](https://github.com/kubeovn/kube-ovn/commit/b1c77ad737b3e2f94d92b97c13400b2b7863f164) some fixes in CI/e2e (#2856)
 * [a94fb0b21](https://github.com/kubeovn/kube-ovn/commit/a94fb0b2177a1cfc7f5e34f18c68dfaf240f51df) manage ovn bfd with libovsdb (#2812)
 * [d9a038ce9](https://github.com/kubeovn/kube-ovn/commit/d9a038ce9d15f39ffad938e3a8905c9ac42bd9e6) update the volumeMounts premission (#2852)
 * [d642f5b58](https://github.com/kubeovn/kube-ovn/commit/d642f5b583ba8225714c5c6533e3bb80a17b989d) fix vip lsp not clean (#2848)
 * [a1cf2b39b](https://github.com/kubeovn/kube-ovn/commit/a1cf2b39b894ceb09b36858563d2c1ba9ff814cf) U2o support custom vpc (#2831)
 * [2068d879c](https://github.com/kubeovn/kube-ovn/commit/2068d879c26f52b729f61e7e67cc89b7e2bfe309) kubectl-ko: fix trace when u2oInterconnection is enabled (#2836)
 * [6ee56d089](https://github.com/kubeovn/kube-ovn/commit/6ee56d089e477057223cd0ce1bd28f82ddcb721e) ci: detect ovs/ovn memory leak (#2839)
 * [776567739](https://github.com/kubeovn/kube-ovn/commit/776567739b40b057a5c244c4d3b192978d442f15) iptables: always do SNAT for access from other nodes to nodeport with external traffic policy set to Local (#2844)
 * [175fb262e](https://github.com/kubeovn/kube-ovn/commit/175fb262e9dc2a39879134e7304665f5ce3b71e2) fix underlay access to node through ovn0 (#2842)
 * [98392b3a1](https://github.com/kubeovn/kube-ovn/commit/98392b3a11fda41cb1229320dc9dcc45e1dd7a65) build(deps): bump github.com/docker/docker (#2843)
 * [da944a3ea](https://github.com/kubeovn/kube-ovn/commit/da944a3eabbc7144a991dbc2b5005acd296e1903) adapt vpc dns in master (#2822)
 * [c7b7a0a5f](https://github.com/kubeovn/kube-ovn/commit/c7b7a0a5f21f56d8849e8363381ebf8378c6c2f5) bump go dependencies (#2820)
 * [94d7cc867](https://github.com/kubeovn/kube-ovn/commit/94d7cc867208e1f2b4d555d57ae56e23fb2bd7b6) fix MTU when subnet is using logical gateway (#2834)
 * [486c61aef](https://github.com/kubeovn/kube-ovn/commit/486c61aefe221ca11f0d192470fd50f948d33623) refactor image builds (#2818)
 * [a7fd9ddf9](https://github.com/kubeovn/kube-ovn/commit/a7fd9ddf97480832480de202463bb0dbb13b5596) build(deps): bump github.com/stretchr/testify from 1.8.2 to 1.8.3 (#2832)
 * [667f5a7c6](https://github.com/kubeovn/kube-ovn/commit/667f5a7c6ef7c9858f6b082d7e5b8104363e2f25) build(deps): bump github.com/onsi/gomega from 1.27.6 to 1.27.7 (#2830)
 * [853abd9d9](https://github.com/kubeovn/kube-ovn/commit/853abd9d9d4270d46b5ee0456523f087970c68ba) vip support create arp proxy logical switch port (#2817)
 * [46bdd01ab](https://github.com/kubeovn/kube-ovn/commit/46bdd01abb0eeed0addd2b8656404235b72a6d91) build(deps): bump github.com/sirupsen/logrus from 1.9.0 to 1.9.2 (#2828)
 * [e988089ee](https://github.com/kubeovn/kube-ovn/commit/e988089ee95618a7915f25bea633887108174eae) build(deps): bump github.com/docker/docker (#2827)
 * [3b8c9edc5](https://github.com/kubeovn/kube-ovn/commit/3b8c9edc55e083468ed98217c1b688ce926ec2ff) add route for service ip range when init vpc-nat-gw (#2821)
 * [4f015f6d3](https://github.com/kubeovn/kube-ovn/commit/4f015f6d36f3d4a0a4b21cd5d971165cbff48f30) do not allocate MAC address when kube-ovn is called as an IPAM plugin (#2816)
 * [a30daea4f](https://github.com/kubeovn/kube-ovn/commit/a30daea4f0a571d0db4e1950a8744cf8476255d8) Iptables nat support share eip (#2805)
 * [0466edce6](https://github.com/kubeovn/kube-ovn/commit/0466edce660cb9ed041ac91e4ae96644bf0089fe) fix typos (#2815)
 * [fca6c9d51](https://github.com/kubeovn/kube-ovn/commit/fca6c9d51151864acff266a03807ab1f3b884ff1) fix some typos (#2814)
 * [630104d55](https://github.com/kubeovn/kube-ovn/commit/630104d55a570b99e98952217cec2d7766de3bd4) add iperf to test group multicast (#2796)
 * [2ba3846b0](https://github.com/kubeovn/kube-ovn/commit/2ba3846b097ec30e4b4bf257c12b4b038480197a) add available check for northd enpoint (#2799)
 * [253358ea8](https://github.com/kubeovn/kube-ovn/commit/253358ea83e4f3f01c2ec826a809f42743e58ea6) manage ovn lr static route with libovsdb (#2804)
 * [781b47d91](https://github.com/kubeovn/kube-ovn/commit/781b47d916c3f1bcff9b4f27f10f8b2bd81714a1) add support of user-defined endpoints to SwitchLBRule (#2777)
 * [74221a6e4](https://github.com/kubeovn/kube-ovn/commit/74221a6e48179528a88b4b24496d375745b78cb2) e2e: fix test container not removed (#2800)
 * [6ddd03bf2](https://github.com/kubeovn/kube-ovn/commit/6ddd03bf2f62ab0c99bb68cf05d903ba3c9466d2) manage ovn lr policy with libovsdb (#2788)
 * [7350db5fa](https://github.com/kubeovn/kube-ovn/commit/7350db5fad4d6a1f185e0fb2a52d5f61f324a29b) build(deps): bump github.com/docker/distribution (#2797)
 * [8f43028a2](https://github.com/kubeovn/kube-ovn/commit/8f43028a24040af1d4d6e65ec528f6d5b9e5d2f3) fix handedeletePod repeat 4 times (#2789)
 * [c8af3dd36](https://github.com/kubeovn/kube-ovn/commit/c8af3dd366d79ab6f9e6fc9c4ad3eb6261e6b604) fix cleanup order (#2792)
 * [b9542ad35](https://github.com/kubeovn/kube-ovn/commit/b9542ad3509d92af8ca0989f37adedac37303587) fix missing main route table for the default vpc (#2785)
 * [1511573dc](https://github.com/kubeovn/kube-ovn/commit/1511573dc4a05a04289f7b6a287c350a1c57399f) add ovn DVR fip e2e (#2780)
 * [0127e10a2](https://github.com/kubeovn/kube-ovn/commit/0127e10a20c53fb665a327cbeb8990bad2f8e023) build(deps): bump github.com/containernetworking/plugins (#2784)
 * [100227bee](https://github.com/kubeovn/kube-ovn/commit/100227bee623dadbc9a9ff3f65992e7353b5cd43) add key lock for more resources (#2781)
 * [16db5082d](https://github.com/kubeovn/kube-ovn/commit/16db5082d54cdde6a1f793578d7d953fb620128e) bump cni plugins to v1.3.0 (#2786)
 * [08e2e66f2](https://github.com/kubeovn/kube-ovn/commit/08e2e66f2e9db182a40b55a8f31fbb561a4dd99f) replace util.DefaultVpc with c.config.ClusterRouter (#2782)
 * [e1154acf6](https://github.com/kubeovn/kube-ovn/commit/e1154acf679fb5bd58f9dbdfac4b84ae1649a995) fix static route recreation after kube-ovn-controller restarts (#2778)
 * [e7190e6ae](https://github.com/kubeovn/kube-ovn/commit/e7190e6ae383db35bffccc153dfe86041ec783b7) clean up code about static routes (#2779)
 * [b1a339b77](https://github.com/kubeovn/kube-ovn/commit/b1a339b77b53ac4acc89c1b170a5396707cc8930) Reorder cleanup step by put subnet and vpc to the last to avoid conflict (#2776)
 * [a2b789cca](https://github.com/kubeovn/kube-ovn/commit/a2b789cca2a8b2106b60d9782e122e989450aff3) optimize kube-ovn-controller logic (#2771)
 * [3b2b07166](https://github.com/kubeovn/kube-ovn/commit/3b2b07166e99bf99f47b2f469f36f486c63551e1) use rate limiting queue with delaying for pod deletion events (#2774)
 * [04e4d2582](https://github.com/kubeovn/kube-ovn/commit/04e4d2582f9ae39df50b100b337ed057c66df571) fix underlay subnet kubectl ko trace error (#2773)
 * [9b1de4816](https://github.com/kubeovn/kube-ovn/commit/9b1de48169218d37d52941f8bf8280ff63f8aea1) feat: natgw qos (#2753)
 * [11b171e14](https://github.com/kubeovn/kube-ovn/commit/11b171e146f0a5672e8e9064a2f60820dc715cac) build(deps): bump github.com/docker/docker (#2770)
 * [62d8122c8](https://github.com/kubeovn/kube-ovn/commit/62d8122c828ead3344c3e4faf91ad7b527bddb1d) fix ip statistics in subnet status (#2769)
 * [d3d01762c](https://github.com/kubeovn/kube-ovn/commit/d3d01762c1640d4e6469e88ee5d9c6a16b07fb45) informer: wait for cache sync before adding event handlers (#2768)
 * [a23dd865a](https://github.com/kubeovn/kube-ovn/commit/a23dd865aa7eec08e315d025850cdc231478639b) build(deps): bump github.com/scylladb/go-set (#2766)
 * [e2bf60f76](https://github.com/kubeovn/kube-ovn/commit/e2bf60f767e9ce3e55f3f250f0479e2725bea64d) support disable arp check ip conflict in vlan provider network (#2760)
 * [c55cbd6e0](https://github.com/kubeovn/kube-ovn/commit/c55cbd6e0d1da7df3e69e5ac18b3f195edd7d3a6) replace string map with string set (#2765)
 * [99be9cb0e](https://github.com/kubeovn/kube-ovn/commit/99be9cb0effc22c1edf0e2ace3ba181279bf2bc1) cni-server: wait ovs-vswitchd to be running (#2759)
 * [1933ed872](https://github.com/kubeovn/kube-ovn/commit/1933ed872c37b386a0c60cd7caa745b51b39464a) kubectl-ko: support trace for pod with host network (#2761)
 * [bf1a3d7c4](https://github.com/kubeovn/kube-ovn/commit/bf1a3d7c424927c97543c46e207f650ba8069809) libovsdb: fix potential duplicate addresses (#2763)
 * [5585d447a](https://github.com/kubeovn/kube-ovn/commit/5585d447a874f40ed3c4dcd458c929efa2a8d784) ci: run kube-ovn e2e for underlay (#2762)
 * [cf1748c64](https://github.com/kubeovn/kube-ovn/commit/cf1748c64fdcf4d95ea4663191d42a1fe8e1cfc4) kubectl-ko: fix pod tracing in underlay (#2757)
 * [6db99d531](https://github.com/kubeovn/kube-ovn/commit/6db99d531f77e9373bffe13c71213ac3ebd4ba82) When Subnet spec.vpc is updated, the status in VPC should also be updated. (#2756)
 * [328a8911d](https://github.com/kubeovn/kube-ovn/commit/328a8911d557cb5deb5c901f10a3630dd728b273) ovn-nbctl: remove unused functions (#2755)
 * [86a07a30e](https://github.com/kubeovn/kube-ovn/commit/86a07a30eeafc2621fd33105218bcbc81251d83d) add route table option in static route for subnet (#2748)
 * [f6414ce18](https://github.com/kubeovn/kube-ovn/commit/f6414ce18e82ed140e636ee847b184ae0d8f1683) replace acl/address_set function call with ovnClient (#2648)
 * [c77f36818](https://github.com/kubeovn/kube-ovn/commit/c77f368187b5be35fe27a97ef60a09375e8ca226) nm not managed only in the change provide nic name case (#2754)
 * [cc1be3ee7](https://github.com/kubeovn/kube-ovn/commit/cc1be3ee79c1e85564515be6ef0f998ad45c53ad) support node local dns cache (#2733)
 * [d7fa2a491](https://github.com/kubeovn/kube-ovn/commit/d7fa2a491e3d1cffff20aa4d7ffb79ac924cb1a0) build(deps): bump google.golang.org/grpc from 1.54.0 to 1.55.0 (#2752)
 * [faff1e627](https://github.com/kubeovn/kube-ovn/commit/faff1e627c94edfbf722633e9c64a17702236a08) build(deps): bump golang.org/x/sys from 0.7.0 to 0.8.0 (#2751)
 * [bdd201b16](https://github.com/kubeovn/kube-ovn/commit/bdd201b16b24436ffa8947db5e8c8a222ba9564f) update eip qos procees, replace qosLabelEIP with natLabelEip (#2736)
 * [d1711acd2](https://github.com/kubeovn/kube-ovn/commit/d1711acd2d27b74815b2ef22a985c64b3f63d58f) refresh nat gw image before using it (#2743)
 * [353df49ac](https://github.com/kubeovn/kube-ovn/commit/353df49ac64134a483b1c258e8d64a886ad8820c) build(deps): bump github.com/prometheus/client_golang (#2745)
 * [91400eccc](https://github.com/kubeovn/kube-ovn/commit/91400eccc6269de60ea7a1553e97e0dce4f8701d) Using full repo name to avoid short-name error in podman (#2746)
 * [fa404a061](https://github.com/kubeovn/kube-ovn/commit/fa404a0614ac1bb1ad582818f2f364a8cddbe891) build(deps): bump github.com/osrg/gobgp/v3 from 3.13.0 to 3.14.0 (#2738)
 * [7eed8341a](https://github.com/kubeovn/kube-ovn/commit/7eed8341a7c96f182374d4e165df82c0c11a7225) add policy route when use old active gateway node for centralized subnet (#2722)
 * [66615b6d8](https://github.com/kubeovn/kube-ovn/commit/66615b6d86d50864650dd7e2f3dd988c3cd3933c) feat: support for multiple external network (#2725)
 * [f8328bdbf](https://github.com/kubeovn/kube-ovn/commit/f8328bdbf282b975999d964567e5d858f3a0e0ad) build(deps): bump github.com/docker/docker (#2732)
 * [6198f691b](https://github.com/kubeovn/kube-ovn/commit/6198f691b5102f51becbe21cdba397a28c8d2801) build(deps): bump github.com/Microsoft/hcsshim from 0.9.8 to 0.9.9 (#2731)
 * [2a015e5c4](https://github.com/kubeovn/kube-ovn/commit/2a015e5c461ef8178ecc24c2419949fee32fb00c) base: remove patch for fixing ofpbuf memory leak (#2715)
 * [a01f9606d](https://github.com/kubeovn/kube-ovn/commit/a01f9606da5121ae1fa1d4fee721ab6376bae5f7) fix recover db failed using method in (#2711)
 * [a6d2a53cc](https://github.com/kubeovn/kube-ovn/commit/a6d2a53cc53080b2e1324e13291e4858cc1db4f9) refactor: improve performance by using cache (#2713)
 * [7dbfd2be0](https://github.com/kubeovn/kube-ovn/commit/7dbfd2be0bc1ee06b4f52baa9871ef0e6b2367f0) For dualstack and ipv6 the default ipv6 range should be same with the ipv4 cidr. (#2708)
 * [15780bfbb](https://github.com/kubeovn/kube-ovn/commit/15780bfbb611bfef95c6f2ccfcc9858a4adf9bbd) feat: support dynamically changing qos for EIP (#2671)
 * [d865b48d1](https://github.com/kubeovn/kube-ovn/commit/d865b48d1ea8ae19e141c49799bff8f8487f3264) base: refactor dockerfile (#2696)
 * [53bfcf447](https://github.com/kubeovn/kube-ovn/commit/53bfcf447140f6c02b406ca666466875b5fbb567) kubectl-ko: add support for tracing nodes (#2697)
 * [f5fee4c93](https://github.com/kubeovn/kube-ovn/commit/f5fee4c936d3e398c08c791093f671a1b0c19c50) cni-server: do not perform ipv4 conflict detection during VM live migration (#2693)
 * [942b87d13](https://github.com/kubeovn/kube-ovn/commit/942b87d13a21c60666ae1075f0bfaba95e28e9e3) fix: iptables nat gw e2e not clean sts eth0 net1 ip (#2698)
 * [236574c75](https://github.com/kubeovn/kube-ovn/commit/236574c750edad001626c02170c21a280e50ca2c) Add random fully when nat (#2681)
 * [9e3f70c1e](https://github.com/kubeovn/kube-ovn/commit/9e3f70c1e2a69d1dbd5b924b00196b09d3592b12) replace StrategicMergePatchType with MergePatchType (#2694)
 * [b59bfd331](https://github.com/kubeovn/kube-ovn/commit/b59bfd331dff52b1fbde28894d28f97e6d6d2a1d) ci: fix scheduled vpc nat gateway e2e (#2692)
 * [d469235f9](https://github.com/kubeovn/kube-ovn/commit/d469235f967a0d519d3d9cbc4e710bc042eb00aa) ovn-controller: do not send GARP on localnet for Kube-OVN ports (#2690)
 * [7db85edd0](https://github.com/kubeovn/kube-ovn/commit/7db85edd06e472f9629cb4a0a6a76029115929e5) netpol: fix enqueueing network policy after LSP creation (#2687)
 * [aba724431](https://github.com/kubeovn/kube-ovn/commit/aba724431824034c8b712d03e3988d93e603523e) add tcp mem collector (#2683)
 * [07a6d4cae](https://github.com/kubeovn/kube-ovn/commit/07a6d4cae2b9854e99c55b702a6bdf406bb99ce8) fix manifest yamls (#2689)
 * [1d6a0fe4b](https://github.com/kubeovn/kube-ovn/commit/1d6a0fe4b4546e44030748df1e2959f9c2662324) attach node name label in ip cr (#2680)
 * [233dc61ec](https://github.com/kubeovn/kube-ovn/commit/233dc61ec23e17daac29bb8b2820048eaeef8c56) adapt ippool annotation (#2678)
 * [095dca268](https://github.com/kubeovn/kube-ovn/commit/095dca2682ae0e745eee4e6f0801c113818263b2) netpol: fix packet drop casued by incorrect address set deletion (#2677)
 * [3dc36c8c3](https://github.com/kubeovn/kube-ovn/commit/3dc36c8c3898aad19d15c005df0b3534ce17254c) fix kubectl ko using ovn-central pod that not in a good status (#2676)
 * [9c5523f7f](https://github.com/kubeovn/kube-ovn/commit/9c5523f7fb9c37e11bb9636ceb6c9f4e9d928183) add nat gw e2e (#2639)
 * [a9993dac3](https://github.com/kubeovn/kube-ovn/commit/a9993dac38dc1dc0d92065ae8149312e4d662ab7) add workflows for release chart (#2672)
 * [4399963e7](https://github.com/kubeovn/kube-ovn/commit/4399963e796744136f0dc96318bd0bddeaf3157e) build(deps): bump github.com/Microsoft/go-winio from 0.6.0 to 0.6.1 (#2663)
 * [d6b0c28df](https://github.com/kubeovn/kube-ovn/commit/d6b0c28df328afc01f7dd2f0a7d933f295f2048c) remove auto update k8s and cadvisor
 * [b57f36ff3](https://github.com/kubeovn/kube-ovn/commit/b57f36ff3a64db751a8bdc5376ca8940da57f3de) build(deps): bump k8s.io/sample-controller from 0.26.3 to 0.26.4 (#2675)
 * [a33adde2f](https://github.com/kubeovn/kube-ovn/commit/a33adde2f2291cd82704c1a41b05fcb868dc7f5b) ignore k8s major and minor dependencies as they always break build.
 * [68f813e0f](https://github.com/kubeovn/kube-ovn/commit/68f813e0fb5617cd524c7ecb87be9bc080b796cf) rename charts (#2667)
 * [933d76e3f](https://github.com/kubeovn/kube-ovn/commit/933d76e3fa419c343e586cff16c2bc11903cf1e3) ipam update condition refactor (#2651)
 * [05e725160](https://github.com/kubeovn/kube-ovn/commit/05e725160535b4e92efd442f69949d965806b639) fix LSP existence check (#2657)
 * [f84343e84](https://github.com/kubeovn/kube-ovn/commit/f84343e84b1bcd7ac0abf874c29ee90ea716dd0d) fix network policy issues (#2652)
 * [148f1bf4a](https://github.com/kubeovn/kube-ovn/commit/148f1bf4a5ca19bfba1af6e627ecc864e029406a) Resolve SetLoadBalancerAffinityTimeout not being effective (#2647)
 * [0b5fc5d36](https://github.com/kubeovn/kube-ovn/commit/0b5fc5d36ac7c52e38150934c1f91c03f12ee804) broadcast free arp when pod is setup (#2638)
 * [dc31cbd20](https://github.com/kubeovn/kube-ovn/commit/dc31cbd2038fa4ca4947d465bc59040f62500676) delete sync user (#2629)
 * [7e872fbe5](https://github.com/kubeovn/kube-ovn/commit/7e872fbe53c7420f0a90a10c23a6f92033be098b) fix: eip qos (#2632)
 * [ddf28fc2f](https://github.com/kubeovn/kube-ovn/commit/ddf28fc2f020bffda95a9d206e42769b6888eaeb) fix: make webhook port configurable. (#2631)
 * [c53d58dac](https://github.com/kubeovn/kube-ovn/commit/c53d58dac3d240205add1de6e0f6d5d89a9842c1) support ovn ipsec  (#2616)
 * [53bf75d21](https://github.com/kubeovn/kube-ovn/commit/53bf75d218a6c0cf994ee49adcdca8dc4d67baa0) feat: add support for EIP QoS (#2550)
 * [1fc5d8534](https://github.com/kubeovn/kube-ovn/commit/1fc5d85346a38593803fb0c7d8c344e4b5d6c4ed) libovsdb: fix race condition in OVN LB operations (#2625)
 * [cfff2db3d](https://github.com/kubeovn/kube-ovn/commit/cfff2db3ded9a5f5c3bce5fc7e4d19e116fe68f8) fix IPAM allocation caused by incorrect pod annotations patch (#2624)
 * [3e67e8935](https://github.com/kubeovn/kube-ovn/commit/3e67e89354ba62269bceaac08869bc2e4717f4f3) ci: deploy multus in thick mode (#2628)
 * [1caaea2ab](https://github.com/kubeovn/kube-ovn/commit/1caaea2ab5b04a5d8573da5b0dfb1a41ea2dd601) libovsdb: use monitor_cond as the monitor method (#2627)
 * [c0ab8351c](https://github.com/kubeovn/kube-ovn/commit/c0ab8351c9f5776553f7f10772531ba162b0e173) ci: fix multus installation (#2622)
 * [84a910b02](https://github.com/kubeovn/kube-ovn/commit/84a910b0235ccecf24950222dfacb64fdfd3de37) ovs: fix dpif-netlink ofpbuf memory leak (#2620)
 * [42a86869e](https://github.com/kubeovn/kube-ovn/commit/42a86869e3fc3d12682f395c65c9b6a38a9a85de) Optimized tolerations code in vpc-nat-gw (#2613)
 * [1e8e38288](https://github.com/kubeovn/kube-ovn/commit/1e8e38288c9accd275ab65168ec82d82962c0606) replace port_group function call with ovnClient (#2608)
 * [9b577403b](https://github.com/kubeovn/kube-ovn/commit/9b577403ba071c267d606affb9fcdd438b841e67) reduce test binary size and add missing webhook build (#2610)
 * [949eb8b73](https://github.com/kubeovn/kube-ovn/commit/949eb8b73f5eb4d2a560ff9cae9cc75875fcb807) fix: ovneip print column and finalizer (#2593)
 * [5babe8e69](https://github.com/kubeovn/kube-ovn/commit/5babe8e69fd31c114e844fc1da9a07ef96e2a9d8) add affinity to vpc-nat-gw (#2609)
 * [6bf15d4a4](https://github.com/kubeovn/kube-ovn/commit/6bf15d4a4cb3edc890501daf090ad4095b804e8e) ci: fix multus installation (#2604)
 * [8629d6340](https://github.com/kubeovn/kube-ovn/commit/8629d6340d7f100c2f0795c180ad342db94f61f8) update .gitignore (#2600)
 * [254598fbd](https://github.com/kubeovn/kube-ovn/commit/254598fbdca25920a5f72ade1e60722344f8a57b) bump go modules (#2603)
 * [602b16052](https://github.com/kubeovn/kube-ovn/commit/602b16052ff139bd2bf32fbafe91164730e2b4cd) build(deps): bump peter-evans/create-pull-request from 4 to 5 (#2606)
 * [787616f1c](https://github.com/kubeovn/kube-ovn/commit/787616f1c2fde385b47ceca9f8f622b27287901e) build(deps): bump github.com/docker/docker (#2605)
 * [62b8761d0](https://github.com/kubeovn/kube-ovn/commit/62b8761d0031ed73279f92db1b7ee4d0f8861046) build(deps): bump golang.org/x/sys from 0.6.0 to 0.7.0 (#2607)
 * [d2523f466](https://github.com/kubeovn/kube-ovn/commit/d2523f466a258a14e03766ec2b2c7a2204c0daaf) cut invalid OVN_NB_DAEMON to make log more readable (#2601)
 * [4c7ddc68d](https://github.com/kubeovn/kube-ovn/commit/4c7ddc68d952eedd44c58730ec8bf7cf514a3a30) unittest: fix length assertion (#2597)
 * [7ba428d71](https://github.com/kubeovn/kube-ovn/commit/7ba428d7101cc2a25e196aa1c3e76a33ee65bf33) use copilot to generate pr content
 * [1a474fd98](https://github.com/kubeovn/kube-ovn/commit/1a474fd988345e87ce5170d81c9570fc5df4f133) replace lb function call with ovnClient (#2598)
 * [a73deb47a](https://github.com/kubeovn/kube-ovn/commit/a73deb47a0b6799a9030912fbc58b6d8bc53f6af) build(deps): bump github.com/osrg/gobgp/v3 from 3.12.0 to 3.13.0 (#2596)
 * [2fb1f95a9](https://github.com/kubeovn/kube-ovn/commit/2fb1f95a9aa960f3991872b40f771bbfe9f551ce) Merge handleAddPod with handleUpdatePod. (#2563)
 * [9399c1e1f](https://github.com/kubeovn/kube-ovn/commit/9399c1e1f64e2ae9df9004e0135d154f43201c4c) fix log (#2586)
 * [da323a52d](https://github.com/kubeovn/kube-ovn/commit/da323a52d103ee1f6cc6c95bfa37c7bb02df9d29) fix: ovn snat and fip delete (#2584)
 * [048e9315f](https://github.com/kubeovn/kube-ovn/commit/048e9315fa58c1d2347aabba74e05c4633b28b0b) underlay: get address/route before setting nm managed to no (#2592)
 * [5d036cd56](https://github.com/kubeovn/kube-ovn/commit/5d036cd560fa0671f13e72572592d182278b6682) update chart description (#2582)
 * [6d50bdc36](https://github.com/kubeovn/kube-ovn/commit/6d50bdc36167737ebb655e8359aa004030dfde83) iptables: use the same mode with kube-proxy (#2535)
 * [09477984f](https://github.com/kubeovn/kube-ovn/commit/09477984f6f4ad8007a611035de12c6b96eacd6b) ci: bump kind image to v1.26.3 (#2581)
 * [5b7bdccb0](https://github.com/kubeovn/kube-ovn/commit/5b7bdccb0099d1dff46b68999acbef97c4838d75) fix: invalid memory address (#2585)
 * [cba9c16e6](https://github.com/kubeovn/kube-ovn/commit/cba9c16e652b8a4a1e398f12f9c0ac43633cfe1e) kubectl ko change solution to collect logs to path kubectl-ko-log (#2575)
 * [bb268618f](https://github.com/kubeovn/kube-ovn/commit/bb268618fa0485a817147f50c78d843d7b9e0b59) if one item is removed, do not requeue (#2578)
 * [5aad7c535](https://github.com/kubeovn/kube-ovn/commit/5aad7c535da3c4bf87d3b7199d8ca15595e76ac1) build(deps): bump github.com/onsi/gomega from 1.27.5 to 1.27.6 (#2579)
 * [a9d66220d](https://github.com/kubeovn/kube-ovn/commit/a9d66220da1624d637d1d36235689111f3d7c396) fix vpc dns when ovn-default is dualstack (#2576)
 * [279717ca7](https://github.com/kubeovn/kube-ovn/commit/279717ca7d1a94d0a32974693be8226402556539) move the vpc-nat generic configurations into one single ConfigMap (#2574)
 * [887df215a](https://github.com/kubeovn/kube-ovn/commit/887df215a42a676782505895322a557ade855615) feat: add ovn dnat (#2565)
 * [02a868734](https://github.com/kubeovn/kube-ovn/commit/02a868734adb662381bb3a6b0241cd7d64d2257d) Fix kubectl ko log loss when restart deployment or ds (#2531)
 * [1d1f5fabb](https://github.com/kubeovn/kube-ovn/commit/1d1f5fabbeccf8ce3523266b95d3780704f7da84) add wait until (#2569)
 * [c0e843fd4](https://github.com/kubeovn/kube-ovn/commit/c0e843fd42ff2eaffe0af0770056f9b8b5638089) do no review dependency update
 * [a7ccd1ae0](https://github.com/kubeovn/kube-ovn/commit/a7ccd1ae0c73fe45ec1f1559f141714845787096) build(deps): bump github.com/opencontainers/runc from 1.1.4 to 1.1.5 (#2572)
 * [5dce9cd2c](https://github.com/kubeovn/kube-ovn/commit/5dce9cd2c80fbd6b80093988216fd1672155ddc6) move ipam.subnet.mutex to caller (#2571)
 * [9fba0b548](https://github.com/kubeovn/kube-ovn/commit/9fba0b548c5489ec183f9740c945ec7d7cec0aaf) build(deps): bump sigs.k8s.io/controller-runtime from 0.14.5 to 0.14.6 (#2568)
 * [3f7997b31](https://github.com/kubeovn/kube-ovn/commit/3f7997b319083e1f0b87a1739804849f299a6c2a) fix: memory leak in IPAM caused by leftover map keys (#2566)
 * [1e9f35299](https://github.com/kubeovn/kube-ovn/commit/1e9f3529980c6ab37f68dc00f5d812cc8e102d55) build(deps): bump github.com/docker/docker (#2567)
 * [8e03e97b8](https://github.com/kubeovn/kube-ovn/commit/8e03e97b84101c331d097236b22f444750850a18) fix ovn-bridge-mappings deletion (#2564)
 * [e19620b0a](https://github.com/kubeovn/kube-ovn/commit/e19620b0aad1abbb61b5046789aa6c406d590f5e) fix lrp deletion after upgrade (#2548)
 * [ed9283489](https://github.com/kubeovn/kube-ovn/commit/ed92834898238c7e9ee2696e5255d348311a6864) fix gw label for vpc update field (#2562)
 * [642fa92a2](https://github.com/kubeovn/kube-ovn/commit/642fa92a2420305db232e1d8db16f984ad630d29) update CRD in helm chart (#2560)
 * [1a41369d6](https://github.com/kubeovn/kube-ovn/commit/1a41369d64482dff27cff5acc87d0aa4015c3302) fix CRD indent in install.sh (#2559)
 * [f955143f7](https://github.com/kubeovn/kube-ovn/commit/f955143f7cad8cfc29ed5923117e138c9f326335) fix update snat rules not effect correctly (#2554)
 * [fd6ec3d8a](https://github.com/kubeovn/kube-ovn/commit/fd6ec3d8af336b13a11b1fe8edab04d90d224b74) fix go mod list (#2556)
 * [b4e7e2e87](https://github.com/kubeovn/kube-ovn/commit/b4e7e2e87983a0ded5f796de5cc8f9e6bb396da1) do not set device unmanaged if NetworkManager is not running (#2549)
 * [fe1b4ac62](https://github.com/kubeovn/kube-ovn/commit/fe1b4ac6242fe9ef5ea871d61de06b22e6ea7c0c) update review bot
 * [f9eb0ca41](https://github.com/kubeovn/kube-ovn/commit/f9eb0ca41ecc0d31bbdeadf970cc01e91ba69db7) build(deps): bump github.com/onsi/gomega from 1.27.4 to 1.27.5 (#2551)
 * [955cf0ffb](https://github.com/kubeovn/kube-ovn/commit/955cf0ffba719f5bf9fc7a4184ebcb1ef4c3ef93) underlay: fix network manager operation (#2546)
 * [b8fc9d9a0](https://github.com/kubeovn/kube-ovn/commit/b8fc9d9a05bc4bd438f677bb644075b73f3c4aef) controller: fix apiserver connection timeout on startup (#2545)
 * [2ae8a9afc](https://github.com/kubeovn/kube-ovn/commit/2ae8a9afc1de5e9951f5715914747a8e0b468ef4) fix update fip rules not effect correctly (#2540)
 * [98dc2f251](https://github.com/kubeovn/kube-ovn/commit/98dc2f251aafccf7c2791e88f3b626e6b41b3292) fix lsp deletion failure when external-ids:ls is empty (#2544)
 * [6b9cdd33a](https://github.com/kubeovn/kube-ovn/commit/6b9cdd33a1edf44a6f579eb4c895f9d829f96e41) sync parameters to charts from install script (#2526)
 * [8c49fc012](https://github.com/kubeovn/kube-ovn/commit/8c49fc012f8506243acfcfdb8307c890d08f6576) underlay: delete altname after renaming the link (#2539)
 * [2a81f4043](https://github.com/kubeovn/kube-ovn/commit/2a81f40433e6b8c31eaff1d6789994fd3f1221f1) failed to delete ovn-fip or ovn-snat (#2534)
 * [17807e555](https://github.com/kubeovn/kube-ovn/commit/17807e5553039141716543ba40cd7f8974dd573a) fix encap_ip will be lost when we restart the ovs-dpdk node (#2543)
 * [829e74c2f](https://github.com/kubeovn/kube-ovn/commit/829e74c2fce41510628b745fa8d09c3a804a49a0) fix service fail (#2537)
 * [bd91f8b86](https://github.com/kubeovn/kube-ovn/commit/bd91f8b86c6ce1bd000692b91cb04da916770936) Add speaker param check (#2538)
 * [7e6feabe8](https://github.com/kubeovn/kube-ovn/commit/7e6feabe89ea87b4f95c1e376886cabd56024225) feat: support nic-hotplug to a running pod. (#2521)
 * [bbe1f3e88](https://github.com/kubeovn/kube-ovn/commit/bbe1f3e884d5d8cb8b27215412e0d60cfa6ccf5c) build(deps): bump google.golang.org/grpc from 1.53.0 to 1.54.0 (#2541)
 * [ae51a6561](https://github.com/kubeovn/kube-ovn/commit/ae51a6561121195f7a6341885714514f4dc3e8a1) fix update dnat rules not effect correctly (#2518)
 * [569b576ac](https://github.com/kubeovn/kube-ovn/commit/569b576ac4e4c42890721909996071147ca0fcd9) underlay: fix link name exchange (#2516)
 * [e97109590](https://github.com/kubeovn/kube-ovn/commit/e97109590322bcc1d035e53d9406965138960dc3) add vip to webhook e2e (#2525)
 * [30d30bfea](https://github.com/kubeovn/kube-ovn/commit/30d30bfead7ab67242b00faa63bf62ced4f83710) fix submariner e2e (#2519)
 * [9eda4859f](https://github.com/kubeovn/kube-ovn/commit/9eda4859fda67897b429034441f59e365b6723d6) fix lsp gc after upgrade (#2513)
 * [0b8964c9e](https://github.com/kubeovn/kube-ovn/commit/0b8964c9e0a6106c3ad05e9d0ae8677cc5b9da26) fix: ovn-fip creation failure due to an excessively long label (#2529)
 * [cc8a11d72](https://github.com/kubeovn/kube-ovn/commit/cc8a11d72021fab878b1be9302f0c46186ca5dc8) add sleep (#2523)
 * [416cc7727](https://github.com/kubeovn/kube-ovn/commit/416cc7727ff85ec01c1a4260d857ced651a4730c) when restart deployment kube-ovn-controller the kubectl ko log loss (#2508)
 * [e7085dec7](https://github.com/kubeovn/kube-ovn/commit/e7085dec79de60b83f25a85a9fa7a385ba192a72) optimize e2e framework (#2492)
 * [4b59bdfc3](https://github.com/kubeovn/kube-ovn/commit/4b59bdfc362e63e0a0ef9ce35819f63b48492ee4) fix ovs patches (#2506)
 * [1138c2cff](https://github.com/kubeovn/kube-ovn/commit/1138c2cff7233be4b39ed3265923cb2236d75021) fix subnet iprange not correct (#2505)
 * [0ebb67853](https://github.com/kubeovn/kube-ovn/commit/0ebb67853164682e5a7fc719c4bab325d4682730) bump k8s to v1.26.3 (#2514)
 * [6fb799232](https://github.com/kubeovn/kube-ovn/commit/6fb7992326a8a844b0081fc2c3f8cfde596b3abd) add kubevirt multus nic lsp before gc process (#2504)
 * [3fc6d8e34](https://github.com/kubeovn/kube-ovn/commit/3fc6d8e3495f9032f119773df02d672eafc3e417) update slack link
 * [46d9edbd2](https://github.com/kubeovn/kube-ovn/commit/46d9edbd28dbfa5f942c420e097f3870ed2305d3) docs: updated CHANGELOG.md (#2515)
 * [36329d54c](https://github.com/kubeovn/kube-ovn/commit/36329d54cb15a40885173526093e2631d29b9633) optimize ovs upgrade script (#2512)
 * [f8aabdf58](https://github.com/kubeovn/kube-ovn/commit/f8aabdf583b80d7bef228265ee3f26a666cb975e) ci: change to pull_request_target
 * [089d8cd20](https://github.com/kubeovn/kube-ovn/commit/089d8cd20b0ac734ddfcb2f60d1666411b9f400d) ci: add openai to review the code (#2511)
 * [ee5e59a94](https://github.com/kubeovn/kube-ovn/commit/ee5e59a9436e85982aeb3fbe14040577a2f39302) add support of user-defined image name for vpc-dns (#2502)
 * [20e702226](https://github.com/kubeovn/kube-ovn/commit/20e702226357cd15474e43ec14c8adee67c34e25) build(deps): bump google.golang.org/protobuf from 1.29.1 to 1.30.0 (#2500)
 * [b6913c522](https://github.com/kubeovn/kube-ovn/commit/b6913c522c32749ef0559a2c61e6ad6840a30d9f) build(deps): bump github.com/Microsoft/hcsshim from 0.9.7 to 0.9.8 (#2499)
 * [443dd58bf](https://github.com/kubeovn/kube-ovn/commit/443dd58bf8fb51005cfacf9ec4f9f78a274cceef) replace lr/ls/lrp/lsp function call with ovnClient (#2477)
 * [599ed2349](https://github.com/kubeovn/kube-ovn/commit/599ed2349ec37b3765e1aa90f49e9f3eacda3847) ci: fix go cache (#2498)
 * [0606f90df](https://github.com/kubeovn/kube-ovn/commit/0606f90df816fac12c81caaa9357a2f2822eceb3) add skip (#2491)
 * [1aff1c4f5](https://github.com/kubeovn/kube-ovn/commit/1aff1c4f5a9ba282e0aae1edef8a85182fab4792) ensure address label is correct before deleting it (#2487)
 * [e9dd28929](https://github.com/kubeovn/kube-ovn/commit/e9dd28929a529d8cb73cea5bb7742b4c56c085e4) fix scheduled submariner e2e (#2469)
 * [c66a93ac7](https://github.com/kubeovn/kube-ovn/commit/c66a93ac71e44a30cf92ddaed7076f540e7f98e2) build(deps): bump actions/setup-go from 3 to 4 (#2490)
 * [a8aede741](https://github.com/kubeovn/kube-ovn/commit/a8aede741c5c08763853a08e87101ab2f1f764a6) build(deps): bump github.com/onsi/gomega from 1.27.3 to 1.27.4 (#2489)
 * [70a220a01](https://github.com/kubeovn/kube-ovn/commit/70a220a013c12f20559898f6f0b352dc56fc8769) add some sleep wait iptables clean (#2488)
 * [0b8e53463](https://github.com/kubeovn/kube-ovn/commit/0b8e534639c7b504e2b2a203de44a930141b8cde) Add kubectl ko log (#2451)
 * [c3620cd09](https://github.com/kubeovn/kube-ovn/commit/c3620cd09bcc43c874506476198b4f0371bc620c) fix: gw configmap may not exist (#2484)
 * [a31235b11](https://github.com/kubeovn/kube-ovn/commit/a31235b114bafa1a08d655d942ac52aea4317d3e) fix ovs qos e2e for versions prior to v1.12 (#2483)
 * [1470c10d3](https://github.com/kubeovn/kube-ovn/commit/1470c10d30aea08f44b16733bdce5b9c27083536) add node to addNodeQueue if required annations are missing (#2481)
 * [47a705577](https://github.com/kubeovn/kube-ovn/commit/47a705577c2cd3c08a5443a1d6f36d78c0f4748f) Add jitter support to netem qos parameters (#2476)
 * [b15fc51bb](https://github.com/kubeovn/kube-ovn/commit/b15fc51bbbdc20dbcfe61b76ce0856efa4f81ea5) build(deps): bump google.golang.org/protobuf from 1.29.0 to 1.29.1 (#2480)
 * [d1cd3dddf](https://github.com/kubeovn/kube-ovn/commit/d1cd3dddfee6ff61e748eefbf987d1334f3364cd) fix ovs-ovn startup/restart (#2467)
 * [b26784f1d](https://github.com/kubeovn/kube-ovn/commit/b26784f1db1b5bdb230148f704b9e4cd78e66d67) fix changging the stopped vm's subnets, the vm cann't start normally (#2463)
 * [7e2e437d5](https://github.com/kubeovn/kube-ovn/commit/7e2e437d57341a80b157eddb0f7b6f1a13a0f11c) build(deps): bump github.com/onsi/gomega from 1.27.2 to 1.27.3 (#2475)
 * [5b07ccbbc](https://github.com/kubeovn/kube-ovn/commit/5b07ccbbc97504d4598ee4a0310b9eda4a32ed5b) when we delete the pod，it's no need to update the sgs assign to pod (#2465)
 * [3fd564b77](https://github.com/kubeovn/kube-ovn/commit/3fd564b77208d9f0cb45a74e00ef0c4c0ab50017) fix libovsdb issues (#2462)
 * [0689a7293](https://github.com/kubeovn/kube-ovn/commit/0689a7293d6c993ef15db83d6c7a89f2bb36ac7f) fix ips CR not found due to etcd error (#2472)
 * [e368a20ed](https://github.com/kubeovn/kube-ovn/commit/e368a20ed1d0fa87a5992216a5228de85070576c) wait for subnet lb (#2471)
 * [0ecd9aff9](https://github.com/kubeovn/kube-ovn/commit/0ecd9aff989696c45b1b9faf57694843090072a0) chore: update base periodically to resolve security issues. (#2470)
 * [5387acf40](https://github.com/kubeovn/kube-ovn/commit/5387acf40925bca7efba3d0f6571b4724bec330d) do not delete external switch if it is created by provider network vlan subnet (#2449)
 * [282706d67](https://github.com/kubeovn/kube-ovn/commit/282706d67fa5d98a468dff48811fb21e939cdda1) add upgrade compatibility (#2468)
 * [482167a92](https://github.com/kubeovn/kube-ovn/commit/482167a921b5670a56493efbfa8f39f708a7dff5) ci: fix ovn-ic installation (#2456)
 * [2bce5080f](https://github.com/kubeovn/kube-ovn/commit/2bce5080f355eaa2dbd5ea1fdc43eaad4c993f2d) Fixed：Prevents grep from prematurely exiting the shell script if it cannot find a pattern (#2466)
 * [4d850e017](https://github.com/kubeovn/kube-ovn/commit/4d850e017c61ef8e515a45df221c026e0047a5f4) add install for webhook (#2460)
 * [f17b4348b](https://github.com/kubeovn/kube-ovn/commit/f17b4348bdbce1890faf2a17b392110fe9d5af7e) e2e add some debug info and sleep (#2439)
 * [8df83cb14](https://github.com/kubeovn/kube-ovn/commit/8df83cb14469203ed1b27b46cefc2edea0fac968) do not set subnet's vlan empty on failure (#2445)
 * [7ae8db6cf](https://github.com/kubeovn/kube-ovn/commit/7ae8db6cf553797067628381050f69aba5b04b5a) wait subnet lb clear in set subnet EnableLb to false e2e (#2450)
 * [674cc2907](https://github.com/kubeovn/kube-ovn/commit/674cc2907dbfeced1c2fdefcd0009c71dd8e8c71) build(deps): bump github.com/emicklei/go-restful/v3 (#2458)
 * [e4c089abe](https://github.com/kubeovn/kube-ovn/commit/e4c089abe1cdcd7ecb63015a5f3340fbafb40da4) ci(Mergify): configuration update (#2457)
 * [0444a2b20](https://github.com/kubeovn/kube-ovn/commit/0444a2b20c1c89df7e4fcb94db1f521e3c395233) kube-ovn-speaker support IPv6/Dual (#2455)
 * [790c7cc29](https://github.com/kubeovn/kube-ovn/commit/790c7cc296aea61d5e6b8a70ad65d628c0455414) replace nb_global function call with ovnClient (#2454)
 * [0d129742f](https://github.com/kubeovn/kube-ovn/commit/0d129742ff3e31f51dfd962e3eae915f1878cd32) build(deps): bump google.golang.org/protobuf from 1.28.1 to 1.29.0 (#2452)
 * [b399cca68](https://github.com/kubeovn/kube-ovn/commit/b399cca68ddab0a2f331b614d14fb8af8ce0b396) fix parsing logical router static routes (#2443)
 * [9df323d7e](https://github.com/kubeovn/kube-ovn/commit/9df323d7ed65a781a7838cc5b2b4c010fb123962) base: fix ovn patches (#2444)
 * [3259b9120](https://github.com/kubeovn/kube-ovn/commit/3259b9120bfa26742245369e1bbad11d472535d7) prepare for libovsdb replacement (#1978)
 * [d71a314d5](https://github.com/kubeovn/kube-ovn/commit/d71a314d5dbd5199d134e4ef42ab6c4cbee312ad) support auto change external bridge (#2437)
 * [f4bef89e2](https://github.com/kubeovn/kube-ovn/commit/f4bef89e20859d4daab6264780c0c0b8186fe1f7) fix ovn-speaker router bug (#2433)
 * [497260efe](https://github.com/kubeovn/kube-ovn/commit/497260efe4abfa2ceef95f9954b26dec57c9f7a3) ovs: change update strategy to RollingUpdate (#2422)
 * [c84479be2](https://github.com/kubeovn/kube-ovn/commit/c84479be2e0ad190fd4cfca4da1fcc25ffc6b914) add kubevirt install (#2430)
 * [e9017e2ab](https://github.com/kubeovn/kube-ovn/commit/e9017e2ab558fb03132ff2a67c0a48b1673fc46e) e2e: wait for subnet to meet specified condition (#2431)
 * [810f7b994](https://github.com/kubeovn/kube-ovn/commit/810f7b994fb937d35534c2452f60e21caa3f4b9c) delete all  invalid ovn lb  strategy and prevent invalid multiple endpoint reconsile (#2419)
 * [25fef7cce](https://github.com/kubeovn/kube-ovn/commit/25fef7cceea1a715369da5e97a61489a4e6099c4) add sumbarier case (#2416)
 * [a99ceb20d](https://github.com/kubeovn/kube-ovn/commit/a99ceb20dd67ed5ea52556f7f3b5adc15b6e8cd7) iptables-rules upgrade compatible (#2429)
 * [570338473](https://github.com/kubeovn/kube-ovn/commit/570338473c3a48e51b0d9d1d2155e7b0724b09bd) add log (#2423)
 * [824f2e0ab](https://github.com/kubeovn/kube-ovn/commit/824f2e0ab9475ee40570659c3717e3026eaf478c) check subnet gateway after wait (#2428)
 * [86c01e6ba](https://github.com/kubeovn/kube-ovn/commit/86c01e6ba6c924393fb551f350733911c6f2de6f) fix chart install/upgrade e2e (#2426)
 * [322eab3b9](https://github.com/kubeovn/kube-ovn/commit/322eab3b927ea2cc47310fe056ef039bd976a866) ci: fix cilium chaining e2e (#2391)
 * [79367647b](https://github.com/kubeovn/kube-ovn/commit/79367647b84607ca531d03084fb19efc4552e614) build(deps): bump golang.org/x/sys from 0.5.0 to 0.6.0 (#2427)
 * [dc5148bb8](https://github.com/kubeovn/kube-ovn/commit/dc5148bb86f66ce9574152f0eacaf3f79714c182) resolve e2e error in v1.12.0 (#2425)
 * [541b641f4](https://github.com/kubeovn/kube-ovn/commit/541b641f431d1c57932c05b36df079e846d7e0a0) update test server and test results (#2421)
 * [980507053](https://github.com/kubeovn/kube-ovn/commit/9805070534fc8496ca296b114f28029263d221e0) Modify the pod scheduling of vpcdns (#2420)
 * [83ab70ffb](https://github.com/kubeovn/kube-ovn/commit/83ab70ffb61241ef7b6467927de56ee96daf1334) e2e: double parallel test nodes in ci (#2411)
 * [fd3bee6ee](https://github.com/kubeovn/kube-ovn/commit/fd3bee6ee22780d1d5739ce2b7ec156daf6c938d) fix scheduled e2e (#2417)
 * [5cd8649bc](https://github.com/kubeovn/kube-ovn/commit/5cd8649bc6eefdccc255955f11a47de75b74d252) build(deps): bump sigs.k8s.io/controller-runtime from 0.14.4 to 0.14.5 (#2415)
 * [68d2ebfa7](https://github.com/kubeovn/kube-ovn/commit/68d2ebfa77f78f662ed29a41b87d21949a87ec04) build(deps): bump github.com/osrg/gobgp/v3 from 3.11.0 to 3.12.0 (#2414)
 * [8f6c21ce1](https://github.com/kubeovn/kube-ovn/commit/8f6c21ce18d06866acd3191238c6314c23842139) build(deps): bump k8s.io/klog/v2 from 2.90.0 to 2.90.1 (#2413)
 * [d837d978e](https://github.com/kubeovn/kube-ovn/commit/d837d978e951aad90d696052e843f811fd8adfe1) bump go modules (#2408)
 * [8fbc5dd17](https://github.com/kubeovn/kube-ovn/commit/8fbc5dd17dee8e08a732e85a9290e38defd58c07) e2e: fix random conflict in parallel processes (#2410)
 * [cedcbbc84](https://github.com/kubeovn/kube-ovn/commit/cedcbbc84efd5e14dacdb5726276b6e99059a51e) fix_base_sg_rule (#2401)
 * [4a28cfb3d](https://github.com/kubeovn/kube-ovn/commit/4a28cfb3df675b82a3dcd0c56d3062d2a7617770) build(deps): bump k8s.io/sample-controller from 0.26.1 to 0.26.2 (#2403)
 * [d30935e0c](https://github.com/kubeovn/kube-ovn/commit/d30935e0c2f2f48222afef44fd2253d2325545af) build(deps): bump github.com/onsi/gomega from 1.27.1 to 1.27.2 (#2396)
 * [645908f66](https://github.com/kubeovn/kube-ovn/commit/645908f66fc276660c455a94d7956487dc15f52a) Support bfd management (#2382)
 * [b1a09bafc](https://github.com/kubeovn/kube-ovn/commit/b1a09bafc7b3df1fa0a3e0c06425f6e0679eacad) remove unused param (#2393)
 * [d24455196](https://github.com/kubeovn/kube-ovn/commit/d244551962ea359dfc1055f2496958ce30519f17) update ipv6 security-group remote group name (#2389)
 * [db435dcc3](https://github.com/kubeovn/kube-ovn/commit/db435dcc320acadd1d578a98876a2bae7451c143) Fix routeregexp ipv6 (#2395)
 * [8a63d2801](https://github.com/kubeovn/kube-ovn/commit/8a63d2801411c82ab6d98c6ff05a7f615c6fb78e) ci: fix ref name check (#2390)
 * [42e6a3022](https://github.com/kubeovn/kube-ovn/commit/42e6a302267a3457fbb037212dae7292d38fbe10) add support of user-defined kubelet directory (#2388)
 * [282644e92](https://github.com/kubeovn/kube-ovn/commit/282644e929dfeb340df6aab3b9277ce6fcd6f8d1) support 1.11 (#2387)
 * [2d1c12529](https://github.com/kubeovn/kube-ovn/commit/2d1c12529ee194102d4e62a8639400299f34a772) ci: skip netpol e2e automatically for push events (#2379)
 * [109704d0a](https://github.com/kubeovn/kube-ovn/commit/109704d0a654749436935f312facebc2db4817f6) ci: make path filter more accurate (#2381)
 * [770224375](https://github.com/kubeovn/kube-ovn/commit/77022437575d112f81002955a681509815dcff4f) build(deps): bump github.com/stretchr/testify from 1.8.1 to 1.8.2 (#2386)
 * [9737c3902](https://github.com/kubeovn/kube-ovn/commit/9737c3902fc6e13893fe35bbd5292502f9f6b423) Fix comment format (#2383)
 * [01e558052](https://github.com/kubeovn/kube-ovn/commit/01e558052ee004e5e854e9b5f2b8c3e51b147bdd) fix: ovs-ovn should reboot now (#2297)
 * [5e0c305f6](https://github.com/kubeovn/kube-ovn/commit/5e0c305f6cd1cdd5db17e14c4a11c2c85db92b90) fix service dual stack add/del cluster ips not change ovn nb (#2367)
 * [ff8361165](https://github.com/kubeovn/kube-ovn/commit/ff83611654b5af8e0a4fd525b82c3f6c512dbaee) ci: fix path filter for windows build (#2378)
 * [4f3f4e741](https://github.com/kubeovn/kube-ovn/commit/4f3f4e741219c4f0969f2ff008e264dd8053562b) e2e: run specs in parallel (#2375)
 * [ffbb15248](https://github.com/kubeovn/kube-ovn/commit/ffbb15248045d512f1947a965cdd1c957c3392f5) add base sg rules for ports (#2365)
 * [db9f9272a](https://github.com/kubeovn/kube-ovn/commit/db9f9272af57f21f62914f7665d6b65bc24a805e) accelerate cleanup (#2376)
 * [50df652c2](https://github.com/kubeovn/kube-ovn/commit/50df652c26c7e95428ff80efb0bc432f963ff86d) update ovnnb model (#2371)
 * [f68044bce](https://github.com/kubeovn/kube-ovn/commit/f68044bce43cca4359e9e3c046c13cf03d854260) docs: updated CHANGELOG.md (#2373)
 * [8a1814a83](https://github.com/kubeovn/kube-ovn/commit/8a1814a83e652272249a0a14d5aa2456d913a02f) fix changelog workflow (#2372)
 * [a1a528b75](https://github.com/kubeovn/kube-ovn/commit/a1a528b75978417d55c1f6431712a6bff9375c23) build(deps): bump github.com/Microsoft/hcsshim from 0.9.6 to 0.9.7 (#2370)
 * [ee53dfe17](https://github.com/kubeovn/kube-ovn/commit/ee53dfe179aeb344b9964485a802c4ee734c22db) Add gateway monitor metrics and event (#2345)
 * [c061ae182](https://github.com/kubeovn/kube-ovn/commit/c061ae182e6cf98783a20876bf1a451fed8ceb35) ci: fix default branch test (#2369)
 * [4a0829a71](https://github.com/kubeovn/kube-ovn/commit/4a0829a71ef442bdf42b9aad5a74a26ce7d30fd5) fix github actions workflows (#2363)
 * [62834eb12](https://github.com/kubeovn/kube-ovn/commit/62834eb1286c213ec0ec10293cadd8617522b58b) Fixed iptables creation failure due to an excessively long label (#2366)
 * [c5d8ebacc](https://github.com/kubeovn/kube-ovn/commit/c5d8ebacc6c4f468609802e6e570cc19744df5ec) use existing node switch cidr instead of the configured one (#2359)
 * [092aa0830](https://github.com/kubeovn/kube-ovn/commit/092aa083074fced40771da6135b10fb1e1090bd0) Do not wait pod deletion one by one to accelerate install (#2360)
 * [1974f8b11](https://github.com/kubeovn/kube-ovn/commit/1974f8b1169e60e5b3cb5b7be789b8ac939eb951) Change log level (#2362)
 * [13f345da4](https://github.com/kubeovn/kube-ovn/commit/13f345da485ee22c2f61be80307189c14a1a68e6) change log level (#2356)
 * [5bd51760d](https://github.com/kubeovn/kube-ovn/commit/5bd51760d19e5cf662d14e88ad14efa0dd04fae3) build(deps): bump github.com/onsi/gomega from 1.27.0 to 1.27.1 (#2357)
 * [3b466d2de](https://github.com/kubeovn/kube-ovn/commit/3b466d2de935101d5df42b564038a9f7c5c03cc7) simplify github actions workflows (#2338)
 * [8fe8bc58c](https://github.com/kubeovn/kube-ovn/commit/8fe8bc58cf008330a2fa9e5bb9e96092d4ee0622) update go version to v1.20 (#2312)
 * [90f504c76](https://github.com/kubeovn/kube-ovn/commit/90f504c760d620b9b23918cb4f80975bda7ce5c1) build(deps): bump golang.org/x/net from 0.6.0 to 0.7.0 (#2353)
 * [a9753c34b](https://github.com/kubeovn/kube-ovn/commit/a9753c34b435ec73f08b0b1f50759b9fe6e91d66) build(deps): bump github.com/onsi/gomega from 1.26.0 to 1.27.0 (#2349)
 * [6e21a93e9](https://github.com/kubeovn/kube-ovn/commit/6e21a93e965b385006388fe918a4fdacef969ee9) chore: no need to wait 30 seconds before kube-ovn-cni get ready. (#2339)
 * [f8b97e72f](https://github.com/kubeovn/kube-ovn/commit/f8b97e72f550fc2af2cb661ae2cd665a108b42bf) do not remove link local route on ovn0 (#2341)
 * [79584c43b](https://github.com/kubeovn/kube-ovn/commit/79584c43b8d10c31550332855272dbfcfa109fcc) fix encap ip when the tunnel interface has multiple addresses (#2340)
 * [156d59767](https://github.com/kubeovn/kube-ovn/commit/156d5976745d95709776be0eb1401e7b8aaf91d9) fix legacy network policy err (#2313)
 * [9c51bd9eb](https://github.com/kubeovn/kube-ovn/commit/9c51bd9ebcb16cb60e6941e93cb2016d6bae9bcb) enqueue endpoint when handling service add event (#2337)
 * [cdf549976](https://github.com/kubeovn/kube-ovn/commit/cdf54997653b20e207dc3c4bad9c9e91e9057fff) Add neighbor-address format check for kube-ovn-speaker (#2335)
 * [b0b46948e](https://github.com/kubeovn/kube-ovn/commit/b0b46948e69f493a34df1806b96197370919b12c) add ovnext0 inside ns on gw node for ecmp static route with bfd (#2237)
 * [4ca994bf8](https://github.com/kubeovn/kube-ovn/commit/4ca994bf85670da0f7fdc839cd3c150fc9ec609b) OVN LB: add support for SCTP protocol (#2331)
 * [ea14e91fa](https://github.com/kubeovn/kube-ovn/commit/ea14e91fa79385f3401c55593735c7cf98914668) fix getting service backends in dual-stack clusters (#2323)
 * [937d3ced4](https://github.com/kubeovn/kube-ovn/commit/937d3ced4a55e4ecc0aa70743abdc1b79ab91692) e2e: skip case of switching session affinity (#2328)
 * [eb2b36a52](https://github.com/kubeovn/kube-ovn/commit/eb2b36a52898bf6b8cbdd7f952be905d38b5ff1b) fix k8s networking dns e2e (#2325)
 * [1c97f58a3](https://github.com/kubeovn/kube-ovn/commit/1c97f58a327ff2e13105ae1e2aba0c913c30e49d) Add the bgp router-id format check (#2316)
 * [f7f2375fa](https://github.com/kubeovn/kube-ovn/commit/f7f2375fa83ac91e6c49d62da686d7e773e82f2a) perform the gateway check but ignore the result when the annotation of subnet is ‘disableGatewayCheck=true’ to make sure of the first network packet (#2290)
 * [0bd7c7e59](https://github.com/kubeovn/kube-ovn/commit/0bd7c7e598761ed06f877cbe5d0cb0a46ad07f60) perf: use empty struct to reduce memory usage (#2327)
 * [b2eaea006](https://github.com/kubeovn/kube-ovn/commit/b2eaea0064abd74e9a997984d743e085f9567868) split netpol cases (#2322)
 * [40b5890ae](https://github.com/kubeovn/kube-ovn/commit/40b5890aefc14376b5e18b471cbeccef94ee88e4) feat: support default service session stickiness timeout (#2311)
 * [83685b5a9](https://github.com/kubeovn/kube-ovn/commit/83685b5a9b4d5b98c69d01848104e5a4f51d4c55) feat: configure routes via pod annotation (#2307)
 * [c8d443ef4](https://github.com/kubeovn/kube-ovn/commit/c8d443ef46044fda8690e833ed508e10a24281b6) build(deps): bump github.com/docker/docker (#2320)
 * [4e2fe3103](https://github.com/kubeovn/kube-ovn/commit/4e2fe3103c709f5f78f669cdc1a151f78dbc21c2) e2e: do not test versions prior to 1.11 for ovn-ic update (#2319)
 * [0d2aa03cd](https://github.com/kubeovn/kube-ovn/commit/0d2aa03cd6ae536e0bf3e1bbe8d5f2ae5a06fe74) ovndb: use Local_Config to configure listen addresses (#2299)
 * [87bacf5ff](https://github.com/kubeovn/kube-ovn/commit/87bacf5ff55ebc3b0fec551ce4a8f957a43c5568) chore: improve the list style in Markdown (#2315)
 * [8c1edc807](https://github.com/kubeovn/kube-ovn/commit/8c1edc807299f5e40f0ccd09b70aecf745045b84) fix egress node and gateway acl should apply after lb. (#2310)
 * [22cc93370](https://github.com/kubeovn/kube-ovn/commit/22cc933708b09c8c5048f649715ea4ea3753e705) fix kube-ovn-controller crash on startup (#2305)
 * [b6eb7ce21](https://github.com/kubeovn/kube-ovn/commit/b6eb7ce219651a6737bf5428d74c0b8f5ecc19ae) build(deps): bump google.golang.org/grpc from 1.52.3 to 1.53.0 (#2308)
 * [5ca2a5c89](https://github.com/kubeovn/kube-ovn/commit/5ca2a5c89c7fa8d60b7abfc790656209421a8555) build(deps): bump golang.org/x/sys from 0.4.0 to 0.5.0 (#2309)
 * [eb31a1783](https://github.com/kubeovn/kube-ovn/commit/eb31a17833c033b30f26b6b78596d5894387239b) ignore e2e for subnet enableEcmp before v1.12.0 (#2306)
 * [f81c43a12](https://github.com/kubeovn/kube-ovn/commit/f81c43a122b9a5328abca11c700c51297896aff5) fix u2o code err (#2300)
 * [993fefaab](https://github.com/kubeovn/kube-ovn/commit/993fefaab49d048c0b110346d1b361abf1e11018) set join subnet.spec.enableLb to nil (#2304)
 * [d1d109727](https://github.com/kubeovn/kube-ovn/commit/d1d109727641cda6c5b2d0bdff91342ecd515aeb) fix image tag in helm chart (#2302)
 * [77cf5e9b2](https://github.com/kubeovn/kube-ovn/commit/77cf5e9b22031d31ab56c8414ea3e2112c3bc77a) update trivy deprecated arg and the ignored CVE. (#2296)
 * [9b85bbac8](https://github.com/kubeovn/kube-ovn/commit/9b85bbac8668b2d1aac16cbcb90e2efe2837811b) move enableEcmp to subnet (#2284)
 * [87eacf593](https://github.com/kubeovn/kube-ovn/commit/87eacf59362bbaf3ef323bd55652a6442c117593) build(deps): bump sigs.k8s.io/controller-runtime from 0.14.3 to 0.14.4 (#2301)
 * [971add05a](https://github.com/kubeovn/kube-ovn/commit/971add05a436600ef48b2d16d1eb18b6527463b9) fix gosec ci installation (#2295)
 * [ac72f7717](https://github.com/kubeovn/kube-ovn/commit/ac72f7717927022eb0fa288423d6394a4bc74163) delete htb qos priority (#2288)
 * [36da29cb9](https://github.com/kubeovn/kube-ovn/commit/36da29cb96cad0eccdbd65608057cf7fb908037b) build(deps): bump sigs.k8s.io/controller-runtime from 0.14.2 to 0.14.3 (#2292)
 * [ea1df9641](https://github.com/kubeovn/kube-ovn/commit/ea1df9641066a2f6102c566d76ae7553af3d8ae7) ovn northd: fix connection inactivity probe (#2286)
 * [54984d67a](https://github.com/kubeovn/kube-ovn/commit/54984d67a20404ec1b8ce14efb6f6d834711c292) fix ct new config error (#2289)
 * [3f0a50088](https://github.com/kubeovn/kube-ovn/commit/3f0a500885ebadd1184924457f82e927effaee6f) fix wrong network interface name in gateway check (#2282)
 * [74a7da88a](https://github.com/kubeovn/kube-ovn/commit/74a7da88aebbad8302e07b650d56df83973db3c1) build(deps): bump github.com/docker/docker (#2287)
 * [20e576991](https://github.com/kubeovn/kube-ovn/commit/20e5769912d0f0a9ee2fd758d5d9b9db6819573f) Improve webhook (#2278)
 * [f0d915136](https://github.com/kubeovn/kube-ovn/commit/f0d915136b96bc14d5c646818da30d88eb44fb41) add named port support (#2273)
 * [9985ee5c6](https://github.com/kubeovn/kube-ovn/commit/9985ee5c6b95bc248fc141ce4e508511d91cd881) fix access from node to overlay pods when network policy ingress exists (#2279)
 * [2b3834009](https://github.com/kubeovn/kube-ovn/commit/2b38340091bbb64c8c622c06bcddcf08a5473053) move enableLb to subnet (#2276)
 * [5712485d7](https://github.com/kubeovn/kube-ovn/commit/5712485d7e94af72c95ec3716d374f4b63a76f83) build(deps): bump github.com/osrg/gobgp/v3 from 3.10.0 to 3.11.0 (#2280)
 * [805f83ea0](https://github.com/kubeovn/kube-ovn/commit/805f83ea026d71fa4104fa21a53132251f0d3a70) add V4/V6UsingIPRange and V4/V6AvailableIPRange in subnet status (#2268)
 * [0c74034d0](https://github.com/kubeovn/kube-ovn/commit/0c74034d056b0f90fc577dc3cb0079f147f0eca4) skip u2o test case before 1.9 (#2274)
 * [eddf18d86](https://github.com/kubeovn/kube-ovn/commit/eddf18d8601a116b15d502f59a06d81fffd6b6e4) fix network break on kube-ovn-cni startup (#2272)
 * [26c506d89](https://github.com/kubeovn/kube-ovn/commit/26c506d89b309dfa6cd7b3e18fb4a96f26bd6df8) bump go modules (#2267)
 * [e10d076ef](https://github.com/kubeovn/kube-ovn/commit/e10d076ef8a004ee4a5af83c3e58df3b084dd5b9) fix setting mtu for ovs internal port (#2247)
 * [155768a3e](https://github.com/kubeovn/kube-ovn/commit/155768a3e4d4215542af4896acbd59f6636cc067) bump ovs/ovn versions (#2254)
 * [281242ef3](https://github.com/kubeovn/kube-ovn/commit/281242ef3deed8d617aa469f6fe84e3ac5de9678) use node ip instead of ovn0 ip when accessing overlay pod/svc from host network (#2243)
 * [ea8523c66](https://github.com/kubeovn/kube-ovn/commit/ea8523c6667bd04cebcb31f3b45970459aa147d8) build(deps): bump google.golang.org/grpc from 1.52.1 to 1.52.3 (#2265)
 * [cb20b12e3](https://github.com/kubeovn/kube-ovn/commit/cb20b12e36d45b36fbdbaf215ce97164a77b5042) build(deps): bump google.golang.org/grpc from 1.52.0 to 1.52.1 (#2264)
 * [69546ffb8](https://github.com/kubeovn/kube-ovn/commit/69546ffb819729f87bce228e9fc9b67114d55324) build(deps): bump k8s.io/klog/v2 from 2.80.1 to 2.90.0 (#2262)
 * [4d9177f7d](https://github.com/kubeovn/kube-ovn/commit/4d9177f7dc5975de9898387f8c280cc901ae52e7) build(deps): bump github.com/onsi/gomega from 1.25.0 to 1.26.0 (#2263)
 * [cc4bfd544](https://github.com/kubeovn/kube-ovn/commit/cc4bfd5445f749da36cb84d97e3202a85b96f87d) build(deps): bump k8s.io/sample-controller from 0.26.0 to 0.26.1 (#2260)
 * [8a6ac128b](https://github.com/kubeovn/kube-ovn/commit/8a6ac128b60609638ce37edf57c5b0f75f2fc9c6) build(deps): bump github.com/docker/docker (#2259)
 * [b33086f74](https://github.com/kubeovn/kube-ovn/commit/b33086f74cda2a3e39968bd613248ca29527ed3a) egress networkpolicy acl add option apply-after-lb (#2251)
 * [625a6854b](https://github.com/kubeovn/kube-ovn/commit/625a6854bc891f101abee55316c7d9ac546991a8) ovn db: add support for listening on pod ip (#2235)
 * [6969dcd80](https://github.com/kubeovn/kube-ovn/commit/6969dcd80d7472a305d489ad7161bc081274038f) update cni plugin to 1.2.0 (#2255)
 * [1f9957099](https://github.com/kubeovn/kube-ovn/commit/1f99570994822c7f4ade37a637fda6a936847629) build(deps): bump github.com/onsi/gomega from 1.24.2 to 1.25.0 (#2257)
 * [486e8ee25](https://github.com/kubeovn/kube-ovn/commit/486e8ee254dd2c33de63da318d044305615c865b) clean up legacy u2o implement (#2248)
 * [5e684e9d8](https://github.com/kubeovn/kube-ovn/commit/5e684e9d82ff759b93f9da848017d73625b2aee6) eip status状态切换缓慢 (#2256)
 * [1049d2459](https://github.com/kubeovn/kube-ovn/commit/1049d2459cf63a0a47573bea28207ce42f8d5eb6) build(deps): bump github.com/containernetworking/plugins (#2253)
 * [9092956fc](https://github.com/kubeovn/kube-ovn/commit/9092956fc0571f5dc9fe3a79b03cc3d4bb56cc71) fix vip create (#2245)
 * [dc731efdc](https://github.com/kubeovn/kube-ovn/commit/dc731efdc2e9f82bc9d870a70c44df851c6f4eec) improve webhook functions for vpc and subnet (#2241)
 * [dfb1cc2b8](https://github.com/kubeovn/kube-ovn/commit/dfb1cc2b8821aa7f2dfb1f71083310e964fd3be9) fix syntax errors (#2240)
 * [e65498020](https://github.com/kubeovn/kube-ovn/commit/e65498020c2c9548f423bc11a01e3a89a65eb0d8) add release-1.11 to scheduled e2e (#2238)
 * [6adf82678](https://github.com/kubeovn/kube-ovn/commit/6adf82678fd97c021d248481d177945bb4230dfe) fix webhook (#2236)
 * [3f5bd39ba](https://github.com/kubeovn/kube-ovn/commit/3f5bd39ba7e97fe8f77107e7fe9284ab705978d3) fix: ovnic del old AZ after establish the new as name (#2229)
 * [b0c17afdc](https://github.com/kubeovn/kube-ovn/commit/b0c17afdce512c2a878d36721a33afb755ebb4e3) prepare for next release
 * [91db26f1a](https://github.com/kubeovn/kube-ovn/commit/91db26f1a108ab3a8fcadb841635f44d89cf3237) build(deps): bump google.golang.org/grpc from 1.51.0 to 1.52.0 (#2234)
 * [a49d18198](https://github.com/kubeovn/kube-ovn/commit/a49d18198e03afe290d73e6b2969200f090e6455) Update CHANGELOG.md for v1.11.0
 * [59bc50f73](https://github.com/kubeovn/kube-ovn/commit/59bc50f734bd8b3fa7ead27ab22ddf1574e77c1c) feat: add helm upgrade e2e (#2222)
 * [25f897378](https://github.com/kubeovn/kube-ovn/commit/25f89737827c66666b9317b5b4392b9071dcc251) fix: now route with connected/static will all be sync (#2231)
 * [c2467d219](https://github.com/kubeovn/kube-ovn/commit/c2467d219d6bfcf082a2f19aa9978d8bc5a8818f) add enable-metrics arg to disable metrics (#2232)
 * [67024ec5f](https://github.com/kubeovn/kube-ovn/commit/67024ec5f0cdaa2a9b910de593a452c5ccb481da) add u2o test case (#2203)
 * [f5d801109](https://github.com/kubeovn/kube-ovn/commit/f5d8011098ef0a22513b53adb2779929d0b7d978) add more args to break test server
 * [f5b9eef8f](https://github.com/kubeovn/kube-ovn/commit/f5b9eef8f3dc7759a678af477ed6585b2be5c234) add release-1.8/1.9/1.10 to scheduled e2e (#2224)
 * [ab5a2c825](https://github.com/kubeovn/kube-ovn/commit/ab5a2c825acb5d20ae65b417814d85adcafdbf1e) cni-server: fix waiting for routed annotation (#2225)
 * [6fd9ea0cb](https://github.com/kubeovn/kube-ovn/commit/6fd9ea0cb390671080361f136039117620e2f6f2) build(deps): bump golang.org/x/sys from 0.3.0 to 0.4.0 (#2223)
 * [cbde65e24](https://github.com/kubeovn/kube-ovn/commit/cbde65e24c2b7a74b6a7e1b0ffd832c5800958ac) feature: detect ipv4 address conflict in underlay (#2208)
 * [64d6f24f5](https://github.com/kubeovn/kube-ovn/commit/64d6f24f5ed4296cafccf6c2a2c3fec159a150d8) fix git ref name in e2e (#2218)
 * [b0cd45c63](https://github.com/kubeovn/kube-ovn/commit/b0cd45c63886c7f38b2ca9fdbabac0d19ddb27c3) fix e2e for v1.8 (#2216)
 * [5843892b0](https://github.com/kubeovn/kube-ovn/commit/5843892b007b1554be0970e616cabfcfb88abd9d) some fixes for e2e testing (#2207)
 * [b6a11789c](https://github.com/kubeovn/kube-ovn/commit/b6a11789c3ae5196a597d71decbf952cbc7688e8) build(deps): bump github.com/osrg/gobgp/v3 from 3.9.0 to 3.10.0 (#2209)
 * [4f08d9413](https://github.com/kubeovn/kube-ovn/commit/4f08d941306c3fd1087ef4ae537eaf0c9ed35fe3) distinguish ippool process for dualstack and normal ippool situation (#2204)
 * [098a8212d](https://github.com/kubeovn/kube-ovn/commit/098a8212dae2c3698c77b249f8e279cb90955e51) u2o feature (#2189)
 * [c0d76fd8b](https://github.com/kubeovn/kube-ovn/commit/c0d76fd8b494a28c68581d54fcce5eb8a2e18e4c) ovn nb and sb can't bind lan ip in ssl (#2200)
 * [1489b65cb](https://github.com/kubeovn/kube-ovn/commit/1489b65cb15f7792434351bb01dc60d29807da3a) build(deps): bump sigs.k8s.io/controller-runtime from 0.14.0 to 0.14.1 (#2199)
 * [16002a282](https://github.com/kubeovn/kube-ovn/commit/16002a282a9e77059611dfa52641ff0119a91ca6) local ip bind to service (#2195)
 * [1407eba24](https://github.com/kubeovn/kube-ovn/commit/1407eba24b9c78231e371dcf4f0f36e729dce5b5) refactor e2e testing (#2078)
 * [86fab667b](https://github.com/kubeovn/kube-ovn/commit/86fab667b414a5ed470ae8ad7ae9e64020844ac9) fix: ovs gc just for pod if (#2187)
 * [1a43c6dee](https://github.com/kubeovn/kube-ovn/commit/1a43c6deed0b00e40b48cc39ff093e781c31314e) update docs link in install.sh (#2196)
 * [02feb9a91](https://github.com/kubeovn/kube-ovn/commit/02feb9a91543c31f72aaa31134008442c5eaa95d) fix lr policy for default subnet with logical gateway enabled (#2177)
 * [3e129fe1b](https://github.com/kubeovn/kube-ovn/commit/3e129fe1b51a51b61ee6fb5213fb97674eb31ee3) sync delete pod process from release-1.9 (#2190)
 * [b6e507065](https://github.com/kubeovn/kube-ovn/commit/b6e507065cf408ac508f4801dab6654ade39302f) fix: update helm 1.11.0 (#2182)
 * [3fb825c8f](https://github.com/kubeovn/kube-ovn/commit/3fb825c8f723c06f4ad20033f9c6dee1cb506146) reserve pod eip static route when update vpc (#2185)
 * [159fd9f02](https://github.com/kubeovn/kube-ovn/commit/159fd9f0290d32bd1dc4d86a6877492c8401066f) ignore conflict check for pod ip crd (#2188)
 * [4d6ad644f](https://github.com/kubeovn/kube-ovn/commit/4d6ad644f8772936afe202169db84aec8f2ff74a) remove unused subnet status fields (#2178)
 * [484fe97ab](https://github.com/kubeovn/kube-ovn/commit/484fe97ab3e214055d8152d518a1d92f082f2116) fix:react leader elect (#2167)
 * [c914fe78d](https://github.com/kubeovn/kube-ovn/commit/c914fe78de5860868ad560c4e631c44d2f7d8be0) fix base/windows build (#2172)
 * [6a8fc2f3b](https://github.com/kubeovn/kube-ovn/commit/6a8fc2f3b5e3e51c2e2074108388ed06bb2c7b11) add metric interface_rx_multicast_packets (#2156)
 * [2b5e28ff9](https://github.com/kubeovn/kube-ovn/commit/2b5e28ff9735df0829eae63c5d947b43dd64f2a3) build(deps): bump github.com/onsi/gomega from 1.24.1 to 1.24.2 (#2168)
 * [0992f36f8](https://github.com/kubeovn/kube-ovn/commit/0992f36f89743e4bc2f06b4290ddcd7f8098f691) update wechat link
 * [d45a04407](https://github.com/kubeovn/kube-ovn/commit/d45a04407532dce7e9e85087835f5e05e961b1c4) build(deps): bump github.com/Microsoft/hcsshim from 0.9.5 to 0.9.6 (#2161)
 * [adecee76f](https://github.com/kubeovn/kube-ovn/commit/adecee76fb680cd24a30358affc28224053136ad) ci: refactor previous push multi arch (#2164)
 * [5e4955c91](https://github.com/kubeovn/kube-ovn/commit/5e4955c918026b2d21abaf5fbeecb8cfdcb56280) security: we should check all the vulnerabilities that can be fixed (#2163)
 * [502a25bfb](https://github.com/kubeovn/kube-ovn/commit/502a25bfb7c014e86801d8232d3a45e0d8d198f3) An error occurred when netpol was added in double-stack mode (#2160)
 * [dbbbddc16](https://github.com/kubeovn/kube-ovn/commit/dbbbddc166e5d09f815cea1fe909fa30dc1bf600) add process for delete networkpolicy start with number (#2157)
 * [26f407fc0](https://github.com/kubeovn/kube-ovn/commit/26f407fc073c9436e322d36212aabeeda1f09d54) security remove private key (#2159)
 * [57457bd46](https://github.com/kubeovn/kube-ovn/commit/57457bd46d7216e897c1816d412023edcb07c488) add scheduled e2e testing (#2144)
 * [5444126aa](https://github.com/kubeovn/kube-ovn/commit/5444126aad775a9779880f06967b491e862412f5) northd: fix race condition in health check (#2154)
 * [755a46a69](https://github.com/kubeovn/kube-ovn/commit/755a46a6981984112e7580d449631f466fce643e) add check for subnet cidr (#2153)
 * [c627468ab](https://github.com/kubeovn/kube-ovn/commit/c627468abbca5544f43b3af3f3e367f9727cd56b) delete nc cmd in image (#2148)
 * [207a52cdf](https://github.com/kubeovn/kube-ovn/commit/207a52cdfbff68572e1d26f64e59e8a4cef09299) bump k8s to v1.26 (#2152)
 * [a4a8b5ad6](https://github.com/kubeovn/kube-ovn/commit/a4a8b5ad619d0e55d799d6c93f77a8228f6c152a) add benchmark test for ipam (#2123)
 * [4b1e78c22](https://github.com/kubeovn/kube-ovn/commit/4b1e78c22839d1cb04a930cc250879eb1469c2ff) update: add YuDong Wang into MAINTAINERS (#2147)
 * [39ee1e7cc](https://github.com/kubeovn/kube-ovn/commit/39ee1e7cc4513c19c73dfdc028cd17927e8b4102) build(deps): bump k8s.io/sample-controller from 0.25.4 to 0.25.5 (#2146)
 * [7aa9bdbcc](https://github.com/kubeovn/kube-ovn/commit/7aa9bdbcc59417e46861f7e497815360fc0b504a) delete nc in base image (#2141)
 * [aab79cb8e](https://github.com/kubeovn/kube-ovn/commit/aab79cb8e1e95b8cccc640de9ac9cac491e9942b) update go modules (#2142)
 * [fa32177d2](https://github.com/kubeovn/kube-ovn/commit/fa32177d242849242469c0f6f69716aa78edfc62) delete ip crd base on podName (#2143)
 * [4072eb762](https://github.com/kubeovn/kube-ovn/commit/4072eb76269c8ddf409874e704d5d89a5b0d50a0) fix vpc spec external not true after init external gw (#2140)
 * [51907e025](https://github.com/kubeovn/kube-ovn/commit/51907e0254e9a7243c7f330dad2f256c65eb007b) refactor ipam unit test (#2126)
 * [ad56e98fb](https://github.com/kubeovn/kube-ovn/commit/ad56e98fb5375faa2fdec66b7cf4de65165b52ef) build(deps): bump github.com/k8snetworkplumbingwg/network-attachment-definition-client (#2139)
 * [012ab59e1](https://github.com/kubeovn/kube-ovn/commit/012ab59e164dc743d047023c22b0efca72401d90) some optimization for provider network status update (#2135)
 * [c410d8b49](https://github.com/kubeovn/kube-ovn/commit/c410d8b498b9ccf49152c6528be905b9a6069e37) simplify iptables eip nat (#2137)
 * [ef4e75554](https://github.com/kubeovn/kube-ovn/commit/ef4e75554d7a1495ca35fb3a95ec89cf57cf5459) kind: support to specify api server address/port (#2134)
 * [9bbf5e439](https://github.com/kubeovn/kube-ovn/commit/9bbf5e4399df6442708d110e0b63abcae3dfeea9) kubectl-ko: fix registry/version (#2133)
 * [2156ef0d7](https://github.com/kubeovn/kube-ovn/commit/2156ef0d7ec42004f5930b2512a4ad4a1e0efa3d) check if subnet cidr is correct (#2136)
 * [f58c88fc2](https://github.com/kubeovn/kube-ovn/commit/f58c88fc244ef7aaa2b28358b813488f6b5e1204)  fix: sometimes alloc ipv6 address failed sometimes ipam.GetStaticAddress return NoAvailableAddress (#2132)
 * [27d22b7fb](https://github.com/kubeovn/kube-ovn/commit/27d22b7fb24b2e7eca178014ef18abdf4abed6a3) fix: delete static route should consider dualstack (#2130)
 * [9b38bf7f6](https://github.com/kubeovn/kube-ovn/commit/9b38bf7f632a84be334b5b4f5cf4d799d1d93e2f) build(deps): bump github.com/osrg/gobgp/v3 from 3.8.0 to 3.9.0 (#2121)
 * [f9f63cae4](https://github.com/kubeovn/kube-ovn/commit/f9f63cae481ddc1c76a2ac09a576616a4e22c613) build(deps): bump github.com/Wifx/gonetworkmanager from 0.4.0 to 0.5.0 (#2122)
 * [67b4dc1b9](https://github.com/kubeovn/kube-ovn/commit/67b4dc1b9933005cd64d30fcc2b052386f564fe7) build(deps): bump golang.org/x/time from 0.2.0 to 0.3.0 (#2120)
 * [78584b7cb](https://github.com/kubeovn/kube-ovn/commit/78584b7cb1b45973fb658fe2e5f36de97aeb4a6a) fix: vlan gw clean in 2 scene (#2117)
 * [b8e15e197](https://github.com/kubeovn/kube-ovn/commit/b8e15e1977a6e39afc3dd407120af9c5797c3c6b) optimize provider network (#2099)
 * [66e96b8e7](https://github.com/kubeovn/kube-ovn/commit/66e96b8e75d5b7b0a0a773a7c525e711fdb9c4fc) build(deps): bump golang.org/x/sys from 0.2.0 to 0.3.0 (#2119)
 * [625e31732](https://github.com/kubeovn/kube-ovn/commit/625e317323734b8179b9641df3316b149b3a3029) fix removing default static route in default vpc (#2116)
 * [141c4c355](https://github.com/kubeovn/kube-ovn/commit/141c4c3556e4753ea89c9fcb565e5ccffe557111) fix: eip deletion (#2118)
 * [86f75c835](https://github.com/kubeovn/kube-ovn/commit/86f75c835de89b86b9cc23c6267834b498c9411e) fix: ecmp route keep delete and recreate  (#2083)
 * [15fd547b3](https://github.com/kubeovn/kube-ovn/commit/15fd547b30ab4e09915130fccdef5e7f9cd7502a) fix policy route for subnets with logical gateway (#2108)
 * [c7549d41c](https://github.com/kubeovn/kube-ovn/commit/c7549d41c048ee5e1ebe68506432d1bfe3315789) build(deps): bump github.com/emicklei/go-restful/v3 from 3.9.0 to 3.10.1 (#2113)
 * [c42dae317](https://github.com/kubeovn/kube-ovn/commit/c42dae3176efc64604226bd608f6b0c6a4eb6108) refactor function name isIPAssignedToPod to isIPAssignedToOtherPod (#2096)
 * [c52f384ee](https://github.com/kubeovn/kube-ovn/commit/c52f384ee92a457307b87bdf5617f981b3d4aec0) build(deps): bump github.com/onsi/gomega from 1.24.0 to 1.24.1 (#2111)
 * [fc80d5922](https://github.com/kubeovn/kube-ovn/commit/fc80d592282c47220c7a26a0e5a52de211fdfb2b) fix: logical gw underlay gw subnet not clean (#2114)
 * [5862b0205](https://github.com/kubeovn/kube-ovn/commit/5862b0205bc229a3f3717690ffc7e9c3ffb79251) build(deps): bump github.com/osrg/gobgp/v3 from 3.6.0 to 3.8.0 (#2110)
 * [4b4bdb3c3](https://github.com/kubeovn/kube-ovn/commit/4b4bdb3c3d58a9c81aa25777fac7a905558fa827) build(deps): bump sigs.k8s.io/controller-runtime from 0.12.3 to 0.13.1 (#2109)
 * [684d1c757](https://github.com/kubeovn/kube-ovn/commit/684d1c757bfa67cc3ce20cb913f389049351d384) fix go mod (#2107)
 * [8ac8cc4e0](https://github.com/kubeovn/kube-ovn/commit/8ac8cc4e0535e814cb5e155b912e11e7a612689d) build(deps): bump github.com/onsi/ginkgo/v2 from 2.3.1 to 2.5.1 (#2103)
 * [12f2f404d](https://github.com/kubeovn/kube-ovn/commit/12f2f404d0fdd492eeac2860812793936bcce72a) build(deps): bump k8s.io/sample-controller from 0.24.4 to 0.25.4 (#2101)
 * [5caec7039](https://github.com/kubeovn/kube-ovn/commit/5caec7039d9b81c3f2e47342a415833dd733078a) build(deps): bump github.com/Microsoft/go-winio from 0.5.2 to 0.6.0 (#2104)
 * [e2eae04cb](https://github.com/kubeovn/kube-ovn/commit/e2eae04cba360a20d84d5ca07a30c5928e8936d7) build(deps): bump google.golang.org/grpc from 1.49.0 to 1.51.0 (#2102)
 * [8f4bf43ae](https://github.com/kubeovn/kube-ovn/commit/8f4bf43ae6ce347ba63b69ed40351294d56ed225) build(deps): bump github.com/Microsoft/hcsshim from 0.9.4 to 0.9.5 (#2100)
 * [47fe3eef8](https://github.com/kubeovn/kube-ovn/commit/47fe3eef820b4e0b8e5547abac39a8b29b221fb8) Create dependabot.yml
 * [5bed4af14](https://github.com/kubeovn/kube-ovn/commit/5bed4af148a13e3aa479cf0b62723a6a84bb7416) replace klog.Fatalf with klog.ErrorS and klog.FlushAndExit (#2093)
 * [2652bcfec](https://github.com/kubeovn/kube-ovn/commit/2652bcfec145a821df9b5677a758ea422dd4824a) fix: slow vip finalizer operation (#2092)
 * [4486e7fe7](https://github.com/kubeovn/kube-ovn/commit/4486e7fe74cb325e1e67d5303a8b4676eb317ea6) ko-trace: support ARP request/reply (#2046)
 * [a659f2e53](https://github.com/kubeovn/kube-ovn/commit/a659f2e53eaeec9783dd8fa5c6b32cc5eea99e31) fix: cni response missing sandbox field (#2089)
 * [88500fa5a](https://github.com/kubeovn/kube-ovn/commit/88500fa5ad96b9ac89180130339634310df43c6c) check if externalIds map is nil when add node as gw for centralized subnet (#2088)
 * [47d3872c1](https://github.com/kubeovn/kube-ovn/commit/47d3872c108f6b4c3c13cbf5ff9e4c7a8eb3eb9b) fix: del createIPS (#2087)
 * [d16163852](https://github.com/kubeovn/kube-ovn/commit/d161638529c7ffe15f0f875074a888c71f32b0e7) fix: add opts for ips del (#2079)
 * [4da9e4e5e](https://github.com/kubeovn/kube-ovn/commit/4da9e4e5e3bb28bcc1b7209c371335b21ba8ce6d) fix ovs bridge not deleted cause by port link not found (#2084)
 * [7344578e2](https://github.com/kubeovn/kube-ovn/commit/7344578e2105b4a0a3010c1439482079eaef1d23) fix libovsdb issues (#2070)
 * [9c292c006](https://github.com/kubeovn/kube-ovn/commit/9c292c006946f02d628161ac524f96c797065e3b) ipset: fix unknown ipset data attribute from kernel (#2086)
 * [def110817](https://github.com/kubeovn/kube-ovn/commit/def110817a7bdfb2842db851afba72a013f60fb7) fix: vpc lrp reset after restart kube-ovn-controller (#2074)
 * [0c668432d](https://github.com/kubeovn/kube-ovn/commit/0c668432dc926f48de9cf3f9617c4a20a7d63dd6) fix: add del bash for redundant ips (#2063)
 * [1c334c8d5](https://github.com/kubeovn/kube-ovn/commit/1c334c8d580addb749380dc521b2bb7095011fb5) refactor: add unknown config to logic switch port (#2067)
 * [419c385bb](https://github.com/kubeovn/kube-ovn/commit/419c385bb7e23939ca2b3daf4a319399118e699a) ovs-dpdk supports adding bond for multi-NICs (#2064)
 * [aef4cd3a6](https://github.com/kubeovn/kube-ovn/commit/aef4cd3a6324bfe4a65c02deb1d4331ab3aba9ed) fix OVN LS/LB gc (#2069)
 * [8aa724eba](https://github.com/kubeovn/kube-ovn/commit/8aa724ebae2583898cc8da6d4b05bf8b7e57c3de) fix: vip ipam not recover all (#2071)
 * [514b76664](https://github.com/kubeovn/kube-ovn/commit/514b7666431db4288094318b686449512fa97a47) bug-fix: make kind-reload invalid (#2068)
 * [657dbf60c](https://github.com/kubeovn/kube-ovn/commit/657dbf60cb5e05a617c5d87f55e3245afa1b45a4) remove no need params svcasname (#2057)
 * [1fcfbc423](https://github.com/kubeovn/kube-ovn/commit/1fcfbc423b015efb425a282e7d97758ba6498b93) Fix:hybrid-dpdk with vxlan tunnel mode，The OVS node does not create a VXLAN tunnel to the OVS-DPDK node (#2065)
 * [a7ed4429f](https://github.com/kubeovn/kube-ovn/commit/a7ed4429fde5a728f294ccd4f84b1b35561db054) update ipv6 address for vpc peer (#2060)
 * [db4fd6290](https://github.com/kubeovn/kube-ovn/commit/db4fd629032d6e4406ce3e017178f7321a601265) perf: reduce controller init time (#2054)
 * [34f426177](https://github.com/kubeovn/kube-ovn/commit/34f4261772f898a5a985f45e3df63789c881138e) reflactor note (#2053)
 * [b22e66ad7](https://github.com/kubeovn/kube-ovn/commit/b22e66ad7daaf32f9212db6a73a89ff2ef5a67f9) fix: replace replace with add to override existing route (#2052)
 * [c08741582](https://github.com/kubeovn/kube-ovn/commit/c08741582c82009f2fb7d94efe7ab49681c5980c) refactor Makefile (#1901)
 * [ea22c1ac0](https://github.com/kubeovn/kube-ovn/commit/ea22c1ac00b2f7b8c46b9bac49fbef49c878f1e3) pass klog verbosity to libovsdb (#2048)
 * [8a29023d2](https://github.com/kubeovn/kube-ovn/commit/8a29023d2b532be4df51a097227cdc05a449b5ee) ovs: fix reaching resubmit limit in underlay (#2038)
 * [db796b436](https://github.com/kubeovn/kube-ovn/commit/db796b4368bb75f1d95fa3a12c2164c1242716b2) sync crd yamls (#2040)
 * [f8f5a4c33](https://github.com/kubeovn/kube-ovn/commit/f8f5a4c33ae2049844ad413327286bbee89d330e) add helm and e2e test (#2020)
 * [5f40c2221](https://github.com/kubeovn/kube-ovn/commit/5f40c22210ce03ef38ec8fa7bfb4557c1449fae0) fix: add unitest (#2030)
 * [686110517](https://github.com/kubeovn/kube-ovn/commit/686110517e03ad63496127addebd9ea5b059c042) fix: pod not add finalizer after add iptables fip (#2041)
 * [75da1603a](https://github.com/kubeovn/kube-ovn/commit/75da1603afb8a51dd4fbb3c996eb2d5454d5b094) feat: ovn eip snat fip (#2029)
 * [730756056](https://github.com/kubeovn/kube-ovn/commit/73075605622d5f8558104ad380793fc82da2c74c) fix: vpc and vpc nat gw not clean (#2032)
 * [79a5ef340](https://github.com/kubeovn/kube-ovn/commit/79a5ef3405b6f5e9314db8716eba391355c9e2a5) update CHANGELOG.md
 * [e7fb20896](https://github.com/kubeovn/kube-ovn/commit/e7fb208965cff4344203c592255a6eb3186f0e80) fix pinger namespace error (#2034)
 * [4abd912f8](https://github.com/kubeovn/kube-ovn/commit/4abd912f875fd2c664b7e7ac69573089bc98d582) iptables: avoid duplicate logging (#2028)
 * [3fb645c9b](https://github.com/kubeovn/kube-ovn/commit/3fb645c9b9b16d16a473a4f7b9c2f3949881ec10) fix: gateway route should stay still when node is pingable (#2011)
 * [92b9c8c3c](https://github.com/kubeovn/kube-ovn/commit/92b9c8c3c9aadc4d1aeab3040b2c74e44725c8c7) update np name with character prefix (#2024)
 * [394978550](https://github.com/kubeovn/kube-ovn/commit/3949785505572474bf4750b622fadfda964c6885) bump kind and node image versions (#2023)
 * [56992c867](https://github.com/kubeovn/kube-ovn/commit/56992c867fb4b5ea47544ffb99ec56eb86253c4d) fix ovn nb/sb health check (#2019)
 * [bf93c458e](https://github.com/kubeovn/kube-ovn/commit/bf93c458ef3d66ebb17ffec3a39e7e3f515a72b8) fix ovs fdb for the local bridge port (#2014)
 * [578301543](https://github.com/kubeovn/kube-ovn/commit/578301543024e989e50f992524220a2f63bdb826) fix go version
 * [ad7cfe87f](https://github.com/kubeovn/kube-ovn/commit/ad7cfe87f5166dffdb097fcaa2ef769bd692b70f) perf: add debug info for perf trace (#2017)
 * [16a958363](https://github.com/kubeovn/kube-ovn/commit/16a9583639ddb0f1a7a602d838493163a518f8d6) fix: not append finalizer (#2012)
 * [688fd5e21](https://github.com/kubeovn/kube-ovn/commit/688fd5e2164651499582c1c79aea09955c51bbca) do not need to delete pg when update networkpolicy (#1959)
 * [c4d8a2f3e](https://github.com/kubeovn/kube-ovn/commit/c4d8a2f3ea26cd0bd16bddf53b2990fae9c39440) test: add test-server to collect packet lost during upgrade (#2010)
 * [f89908e79](https://github.com/kubeovn/kube-ovn/commit/f89908e7941b8edf656503f310aaac04bfd60870) support create iptables fip and eip automatically if pod enable fip  (#1993)
 * [ab80fd88b](https://github.com/kubeovn/kube-ovn/commit/ab80fd88b90935d9745210655fd5422b8a070c55) ci: upgrade deprecated actions (#2004)
 * [de5ef5116](https://github.com/kubeovn/kube-ovn/commit/de5ef5116c2f4388d9bbab32e537923776b54f43) fix: make ip deletion the same as creation (#2002)
 * [bfcc952c3](https://github.com/kubeovn/kube-ovn/commit/bfcc952c32acbfbd5018a0bac51e66bc7c27295e) fix: Add support for Mellanox NIC (#1999)
 * [f4c977f14](https://github.com/kubeovn/kube-ovn/commit/f4c977f1443b15ebe1da99ebb23e40fb813cab88) fix: nat gw not enqueue its resources (#1996)
 * [32f65f81d](https://github.com/kubeovn/kube-ovn/commit/32f65f81d46db912dd5f19ee877eb3aff6d9fcd9) fix: delete fiprule failed at first time (#1998)
 * [eaa936b38](https://github.com/kubeovn/kube-ovn/commit/eaa936b389104733395d87541fde62fab0cc0390) fix typo (#1994)
 * [dd3790ac1](https://github.com/kubeovn/kube-ovn/commit/dd3790ac1da2af5c09c65ed1fecb80c0d5a00a33) feat: now interface for containerd could be inspected (#1987)
 * [fee5bfd3d](https://github.com/kubeovn/kube-ovn/commit/fee5bfd3db394c65c96f3edb04810616e55aaf0f) fix: snat conntrack race (#1985)
 * [e1f7d72c1](https://github.com/kubeovn/kube-ovn/commit/e1f7d72c15b61b083590672745c79b5d8a2d903a) add check of write to ovn sb db for ovn-controller (#1989)
 * [892aa759d](https://github.com/kubeovn/kube-ovn/commit/892aa759d475586d6915278d0cd9343ca4c67c76) fix grep matching device in routes (#1986)
 * [113f62f6f](https://github.com/kubeovn/kube-ovn/commit/113f62f6fc332776c7aeee0b535ebb408adc04ee) delete pod after TerminationGracePeriodSeconds (#1984)
 * [87996f754](https://github.com/kubeovn/kube-ovn/commit/87996f75465e8c8ba1d7d5a7d5d908ec11e0033f) ovs: fix waiting flows in underlay networking (#1983)
 * [eea788869](https://github.com/kubeovn/kube-ovn/commit/eea7888699e3dae8790b0c931be61ab07094ed0e) feature: support default vpc use nat gw pod as cust vpc (#1979)
 * [3d2c7a59a](https://github.com/kubeovn/kube-ovn/commit/3d2c7a59a8cd2eff4f901fd18996a305c0523287) ovn db: recover automatically on startup if db corruption is detected (#1980)
 * [9ff3b9c05](https://github.com/kubeovn/kube-ovn/commit/9ff3b9c05983b7b88b8fb54a526ee197bd0326cf) fix: modify src route priority (#1973)
 * [57c75c1e2](https://github.com/kubeovn/kube-ovn/commit/57c75c1e2c328391014b8d8e690d79addc560dd4) upgrade ovs-ovn pod by generation version instead of chart version (#1960)
 * [b9e98e524](https://github.com/kubeovn/kube-ovn/commit/b9e98e52473a6bdec2cea27bfdf18f309f81943b) avoid concurrent subnet status update (#1976)
 * [ea854d460](https://github.com/kubeovn/kube-ovn/commit/ea854d460a106f6fc53ea34d429d2076530afaf0) fix metrics name (#1977)
 * [15f676f64](https://github.com/kubeovn/kube-ovn/commit/15f676f64b62c40e55ab756a7610830b595d740c) add vm pod to ipam by ip when initIPAM (#1974)
 * [afe06d819](https://github.com/kubeovn/kube-ovn/commit/afe06d819d159a2bde66343c58354979b37dece2) validate nbctl socket path in start-controller.sh (#1971)
 * [3796a5828](https://github.com/kubeovn/kube-ovn/commit/3796a5828f4f95205cdd745748fae2bf2c464b88) skip CVE-2022-3358 (#1972)
 * [80aab2ea0](https://github.com/kubeovn/kube-ovn/commit/80aab2ea046a62619c9a8956a42ba7526aeed516) fix version mismatch between the Ginkgo CLI and the imported package (#1967)
 * [b7863bdb1](https://github.com/kubeovn/kube-ovn/commit/b7863bdb1c6cee1d732930a8cd27389d3c018ca2) ovs: fix mac learning in environments with hairpin enabled (#1943)
 * [95939ca4d](https://github.com/kubeovn/kube-ovn/commit/95939ca4d9a652b151ea57320d7e67800f80a921) fix: add  default deny acl (#1935)
 * [de3d65c01](https://github.com/kubeovn/kube-ovn/commit/de3d65c01a1a4430590fe083e90494ce17425ecd) Fix registry for ovn-central container in install.sh (#1951)
 * [c8d22d2cf](https://github.com/kubeovn/kube-ovn/commit/c8d22d2cf9c0644035d6d5a0cf6b454ac5f59a3e) ovs: add fdb update logging (#1941)
 * [f1f6642b9](https://github.com/kubeovn/kube-ovn/commit/f1f6642b939202968abb427dae7ec05ca8498e1e) add chart version check when upgrade ovs-ovn pod (#1942)
 * [73fde2cef](https://github.com/kubeovn/kube-ovn/commit/73fde2cefd8c02d4bd440a582cfab633ac6cb575) fix underlay e2e testing (#1929)
 * [38956b6ca](https://github.com/kubeovn/kube-ovn/commit/38956b6cad4ec99a330d5abb3d16d520f76f802d) set leader flag when get leader (#1939)
 * [af6973fe3](https://github.com/kubeovn/kube-ovn/commit/af6973fe362cb9f6a47ca6b4f0726754c46247cf) set ovsdb-server vlog level to avoid warnings caused by ovs-vsctl (#1937)
 * [a3292078d](https://github.com/kubeovn/kube-ovn/commit/a3292078df22377509c8d480acde5934cc611a3f) fix: UpdateNatRule will error when logicalIP, externalIP is different protocol; replace : to \\: when IPv6 in ovs cli.
 * [76541ef12](https://github.com/kubeovn/kube-ovn/commit/76541ef12fb4b0dbc62060efd0fea5ffa08c0211) fix: noAllowLiveMigration port can't sync vips (#7)
 * [474206be5](https://github.com/kubeovn/kube-ovn/commit/474206be5e37c5271345f1e07638f673e4d40025) fix: add pod not update vip virtual port
 * [596741bcd](https://github.com/kubeovn/kube-ovn/commit/596741bcd31f8e04c374df6b2e08def8930773f0) fix: delete chassis (#1927)
 * [395a35549](https://github.com/kubeovn/kube-ovn/commit/395a35549a287aaabaad950c5ba3e9c5a5261bd5) fix: pod mistaken ls label (#1925)
 * [797100edf](https://github.com/kubeovn/kube-ovn/commit/797100edf42b2731a31fb2268ab82888a8077946) ignore pod without lsp when add pod to port-group (#1924)
 * [1a49e738f](https://github.com/kubeovn/kube-ovn/commit/1a49e738f214d96722181f69f0c0e9fc29bca0e4) add network partition check in ovn probes (#1923)
 * [16c0ed9fd](https://github.com/kubeovn/kube-ovn/commit/16c0ed9fdf5550f6a089a86cb7d873a60f190a2c) fix: fip unbind can't take effect immediately when conntrack record exists (#1922)
 * [606e6f62d](https://github.com/kubeovn/kube-ovn/commit/606e6f62d0b79d1fd5d30a37c0ae1fab45122b0c) No need to change deivceID to sriov_netdevice. (#1904)
 * [76dd9afa7](https://github.com/kubeovn/kube-ovn/commit/76dd9afa79e1c9c719f13dd815cd9e39b84b1917) update ns annotation when subnet cidr changed (#1921)
 * [8d1ce4206](https://github.com/kubeovn/kube-ovn/commit/8d1ce42066f076c98086446ef751c7f3923e35e7) fix EIP/SNAT on dynamic Pod annotation (#1918)
 * [4882c354e](https://github.com/kubeovn/kube-ovn/commit/4882c354e70513e568bc5029e2a59901e9e16498) fix: eip and nat crd can delete even if nat gw pod deleted and ipatab… (#1917)
 * [d8886d133](https://github.com/kubeovn/kube-ovn/commit/d8886d1334152e836b2e0129916d42a3663b541a) fix missing crd (#1909)
 * [8d2991e37](https://github.com/kubeovn/kube-ovn/commit/8d2991e371da857583d8a17a5858440e84e8b09a) Nat gw support toleration (#1907)
 * [b3bbe1d4f](https://github.com/kubeovn/kube-ovn/commit/b3bbe1d4f3fc8a0a5f5bc5222d097eca1e8e0174) Update USERS.md (#1908)
 * [dbe4ebb31](https://github.com/kubeovn/kube-ovn/commit/dbe4ebb31496d959768dd6d451939b11f831856a) fix typo (#1897)
 * [8d15497a8](https://github.com/kubeovn/kube-ovn/commit/8d15497a8f66e2acf4ea64edd4c3ff91db22ef13) fix: Make the /sys directory in ovs-ovn-dpdk pod writable (#1899)
 * [5befab467](https://github.com/kubeovn/kube-ovn/commit/5befab467404775e1f10948ce0979df792b19ac2) fix: failed to add eip (#1898)
 * [7fae28aec](https://github.com/kubeovn/kube-ovn/commit/7fae28aecf7119f8ed3f777defac7294f699be4c) fix: gatewaynode might be null (#1896)
 * [57edfe41d](https://github.com/kubeovn/kube-ovn/commit/57edfe41d82ef00d5f58e3e2755f2a5f300f3dd7) ci: increase golangci-lint timeout (#1894)
 * [07f12c2df](https://github.com/kubeovn/kube-ovn/commit/07f12c2df1ce79b18143ae9a9e47de65685c319e) update Go to version 1.19 (#1892)
 * [e5878ff9f](https://github.com/kubeovn/kube-ovn/commit/e5878ff9f00887d37a1fb51dc9f5b14ab38e6f87) fix: api rollback (#1895)
 * [82db50fe7](https://github.com/kubeovn/kube-ovn/commit/82db50fe71bf6e9d40bd9bf033615ee939675297) ci: use concurrency to ensure that only a single workflow (#1850)
 * [83b867ab5](https://github.com/kubeovn/kube-ovn/commit/83b867ab54a9e4d6165d0c71bd510cac7d2feabe) kubectl-ko: turn off pipefail for ovn leader check (#1891)
 * [959a64dca](https://github.com/kubeovn/kube-ovn/commit/959a64dcab1822fb7b206f3b55ec26492ee7e80a) kubectl-ko: fix trace for KubeVirt VM (#1802)
 * [10fd33304](https://github.com/kubeovn/kube-ovn/commit/10fd33304ff3e04b27f2ccc72b50207c900035e4) fix duplicate logs for leader election (#1886)
 * [13ebb855d](https://github.com/kubeovn/kube-ovn/commit/13ebb855dd8494af6e63ce9eda4bc0f25cec4453) fix setting ether dst addr for dnat (#1881)
 * [b14850356](https://github.com/kubeovn/kube-ovn/commit/b1485035635ea2812a45e5d4e6ec54fe9864b590) change the prtocol string to const (#1887)
 * [7e56931f3](https://github.com/kubeovn/kube-ovn/commit/7e56931f36df58281cdf887529c6cd99e18c8aa3) refactor iptables rules (#1868)
 * [14898dd3c](https://github.com/kubeovn/kube-ovn/commit/14898dd3c70025c149a00dcdda293add16bf1a9d) cni should handler unmont volume, when delete pod. (#1873)
 * [031ed0317](https://github.com/kubeovn/kube-ovn/commit/031ed0317e595c80dce6cd4eb3a27487041e0543) delete and recreate netem qos when update process (#1872)
 * [dedd5aaab](https://github.com/kubeovn/kube-ovn/commit/dedd5aaabd1e95daf8a20ba73b0dac26384e346c) feat: check configuration (#1832)
 * [97b41127e](https://github.com/kubeovn/kube-ovn/commit/97b41127e14c9b9df5ae930f73f3ded81f49bd25) security: conform to gosec G114 (#1860)
 * [2dfb6e72a](https://github.com/kubeovn/kube-ovn/commit/2dfb6e72ace6ccf7fb6e78857cdf22954b1946ec) update CHANGELOG.md
 * [b6450b2f2](https://github.com/kubeovn/kube-ovn/commit/b6450b2f27fc80b27c61116085053f26cd7cdccd) e2e: add timeout for waiting resources to be ready (#1871)
 * [e1656752a](https://github.com/kubeovn/kube-ovn/commit/e1656752a12e43674dbef73056384f1a2995eecc) upgrade to Ginkgo v2 (#1861)
 * [0adecb0cb](https://github.com/kubeovn/kube-ovn/commit/0adecb0cb7217e77210b6afcf8b7bab530ae98b5) feat: reduce downtime by increasing arp cache timeout
 * [c24f678b8](https://github.com/kubeovn/kube-ovn/commit/c24f678b8b63907af3b5961674376229da38b56e) feat: reduce wait time by counting the flow num.
 * [05611aa75](https://github.com/kubeovn/kube-ovn/commit/05611aa758de93adba97b09c539438f790cd908f) fix: missing stop_ovn_daemon args
 * [6ab837c4e](https://github.com/kubeovn/kube-ovn/commit/6ab837c4e70da188e4133690177eb1d4b21ec5cf) fix: nat gw pod should set default gw to net1 so that to access public (#1864)
 * [b08765dc0](https://github.com/kubeovn/kube-ovn/commit/b08765dc092210d7dc72cb3b284df071dd2e5ebd) delete log severity for drop acl when update networkpolicy (#1863)
 * [e13c4ef13](https://github.com/kubeovn/kube-ovn/commit/e13c4ef13df25f092afd5a3bf0484931013444bf) ovs: fix log file descriptor leak in monitor process (#1855)
 * [8a235e9e8](https://github.com/kubeovn/kube-ovn/commit/8a235e9e8dc84adbd89aaca978c30ce37328c47b) fix: dnat port not use whole words to check (#1854)
 * [34e02ebb5](https://github.com/kubeovn/kube-ovn/commit/34e02ebb545908fd59e879fac15f91dfafa855d9) fix ovs-ovn logging (#1848)
 * [2a2a32f99](https://github.com/kubeovn/kube-ovn/commit/2a2a32f99ea591a59cbb737852870b47c72e84ac) fix ovn dhcp not work with ovs-dpdk (#1853)
 * [44e412509](https://github.com/kubeovn/kube-ovn/commit/44e4125093672f0923efd10069d85c90a9e5383d) docs: Update USERS.md (#1851)
 * [51f491f27](https://github.com/kubeovn/kube-ovn/commit/51f491f272914a5afa70c068cdf65bf1433320c6) fix: multus macvlan ipvlan use kube-ovn ipam，but  ip not inited in init-ipam (#1843)
 * [8ef6c01cf](https://github.com/kubeovn/kube-ovn/commit/8ef6c01cfb47bdc81c8e3f51a2b4e7420181e36c) fix underlay e2e (#1828)
 * [c33276fd5](https://github.com/kubeovn/kube-ovn/commit/c33276fd5ec83243b79417a0ad796c8473fdc48a) fix arping error log (#1841)
 * [69cf5ca59](https://github.com/kubeovn/kube-ovn/commit/69cf5ca5991945cbc8c015c1c562c88bdbaf73be) ko: fix kube-proxy check (#1842)
 * [5012ff3ef](https://github.com/kubeovn/kube-ovn/commit/5012ff3efc34b759c9852ee9b45e6de3979e8add) base: use patch from OVN upstream (#1844)
 * [07773cb08](https://github.com/kubeovn/kube-ovn/commit/07773cb08212cc6d5cc2556301317c5400fed977) ci: switch environment to ubuntu-20.04 (#1838)
 * [eded55166](https://github.com/kubeovn/kube-ovn/commit/eded55166f6d13c47e2c3563ba10c0dfed4eb739) ci: split image builds to speed up jobs (#1807)
 * [656bd46ca](https://github.com/kubeovn/kube-ovn/commit/656bd46ca595bfc480f24ca575bd564709596291) ci: update Go cache to speed up jobs (#1829)
 * [57d74bff3](https://github.com/kubeovn/kube-ovn/commit/57d74bff3bc847db36473ca2f8ce64cfc1a27b87) windows: fix ovs/ovn versions and patches (#1830)
 * [babd80219](https://github.com/kubeovn/kube-ovn/commit/babd802194c97d877d1479236c722c4bbc33a861) 修改 ovs-ovn-dpdk 容器镜像编译打包，解决容器中 ovs 运行不正常：无法添加物理网卡，无法创建 vhostuserclient port 问题 (#1831)
 * [0ed5c924c](https://github.com/kubeovn/kube-ovn/commit/0ed5c924c5bd8a78ff91f2ffb4a134be7b69b76b) support adding routes in underlay Pods for access to overlay Pods (#1762)
 * [013549ab9](https://github.com/kubeovn/kube-ovn/commit/013549ab94ee5e4e74d0f6c9f88af81812d2bc70) update centralized subnet gateway ready patch operation (#1827)
 * [9937ef876](https://github.com/kubeovn/kube-ovn/commit/9937ef8764e24a89b7c39ad8855315f3e7ca0521) remove pod security policy (#1822)
 * [725957a7b](https://github.com/kubeovn/kube-ovn/commit/725957a7ba9c46a415b0cc21da7d90d9c2c003cb) fix duplicate log for tunnel interface decision (#1823)
 * [2e64133ce](https://github.com/kubeovn/kube-ovn/commit/2e64133ce186763d4c5c1f09566a2811da9cd606) update ovs/ovn version to fix hardware offload (#1821)
 * [9c87d9d13](https://github.com/kubeovn/kube-ovn/commit/9c87d9d136f86a5834a9408f763d374ee542a7d8) fix: use full longest word to match full ip about dnat (#1825)
 * [385064a0c](https://github.com/kubeovn/kube-ovn/commit/385064a0ca1e2927ee71c76afb97e8e1c9d0627a) update centralize subnet gatewayNode until gw is ready (#1814)
 * [4944db750](https://github.com/kubeovn/kube-ovn/commit/4944db750ec6981406a34485254a66e63d03adde) initialize IPAM from IP CR with empty PodType for sts Pods (#1812)
 * [d41c043d8](https://github.com/kubeovn/kube-ovn/commit/d41c043d825cd85990f36e5f0bc9587030bd82db) feat: add editable ovn-ic (#1795)
 * [ddcdfb9ee](https://github.com/kubeovn/kube-ovn/commit/ddcdfb9eedbaca142160be947e70aaefb5932849) kubectl-ko: fix missing env-check (#1804)
 * [8b2df5880](https://github.com/kubeovn/kube-ovn/commit/8b2df588095dac7db6838b6bd5abb348d72125bb) kubectl-ko: fix destination mac (#1801)
 * [e8816e964](https://github.com/kubeovn/kube-ovn/commit/e8816e964a6cb7794973b3b9e0ccc833367bc3e8) fix cilium e2e (#1759)
 * [bc3804151](https://github.com/kubeovn/kube-ovn/commit/bc380415163b972a7bfa4f505e472b1c0ebd376e) abort kube-ovn-controller on leader change (#1797)
 * [26e77eadf](https://github.com/kubeovn/kube-ovn/commit/26e77eadfb93e49185769999f57edf734ddc2d18) avoid invalid ovn-nbctl daemon socket path (#1799)
 * [e7064062b](https://github.com/kubeovn/kube-ovn/commit/e7064062b7eb319c37361bb9f27a5866b6de9f41) update CHANGELOG.md
 * [70f1b141c](https://github.com/kubeovn/kube-ovn/commit/70f1b141ca7f5c58e64a46c9dd8f89a0e1d82df8) Perf/memleak (#1791)
 * [225da25e7](https://github.com/kubeovn/kube-ovn/commit/225da25e746b6890e708c17891b291327dabac01) delete htb qos when releated annotation is deleted (#1788)
 * [8db05a1fa](https://github.com/kubeovn/kube-ovn/commit/8db05a1fa600c468b395020d1bf27589a9fde126) Fix nag gw gc (#1783)
 * [277f6f697](https://github.com/kubeovn/kube-ovn/commit/277f6f697708cfad21625f3dad25eac09134f54c) fix iptables for services with external traffic policy set to Local (#1773)
 * [42812b92e](https://github.com/kubeovn/kube-ovn/commit/42812b92ef1766de0d5d80db3935231d277be052) perf: reduce metrics labels (#1784)
 * [9ffb9d229](https://github.com/kubeovn/kube-ovn/commit/9ffb9d229e8e4cbbe5055965c8a649522228bb36) northd: remove lookup_arp_ip actions (#1780)
 * [9efd4bb2f](https://github.com/kubeovn/kube-ovn/commit/9efd4bb2f88760b1fb16f53a248efe85ef888620) fix: 5ms is too short for eip and nats creation (#1781)
 * [80425b7cf](https://github.com/kubeovn/kube-ovn/commit/80425b7cf8649f35486aa5eaee27069c55512fe7) Lb-svc supports custom VPCs (#1779)
 * [cd00ddb64](https://github.com/kubeovn/kube-ovn/commit/cd00ddb642a799fb896f756d3bb56fc06b7bd3b1) fix ovnic e2e (#1763)
 * [51bf142fc](https://github.com/kubeovn/kube-ovn/commit/51bf142fcadd934be2f9b648185db1cf0b0eece6) fix iptables for service traffic when external traffic policy set to local (#1728)
 * [916600f6d](https://github.com/kubeovn/kube-ovn/commit/916600f6dd94236373807f0f553697c42d9acdd2) set sysctl variables on cni server startup (#1758)
 * [c10c71184](https://github.com/kubeovn/kube-ovn/commit/c10c71184f5afbfc53f36017d411c89cd2f05971) fix: add omitempty to subnet spec (#1765)
 * [2dd46e693](https://github.com/kubeovn/kube-ovn/commit/2dd46e6931a726e1eec85c3234d0e7a81adf9deb) perf: replace jemalloc to reduce memory usage (#1764)
 * [35513157e](https://github.com/kubeovn/kube-ovn/commit/35513157e44b6e9de65f12b29a444ecadef65c70) avoid patch interface deletion & recreation during restart (#1741)
 * [e7ce68bbf](https://github.com/kubeovn/kube-ovn/commit/e7ce68bbf710f0780e82f01d98b9c93b21755983) feature: support exchange link names of OVS bridge and provider nic in underlay networks (#1736)
 * [5254731f3](https://github.com/kubeovn/kube-ovn/commit/5254731f3e86e24aa764d3e09226b44f91431f03) dpdk-v2 ，--with-hybrid-dpdk 修改 Dockerfile.base-dpdk 解决 编译安装 ovs-dpdk 正常运行 (#1754)
 * [3717c6cb6](https://github.com/kubeovn/kube-ovn/commit/3717c6cb620a9e9e2cec97b2a95831ad321982e0) fix: Adjust order for Log and output err when get NatRule faild. (#1751)
 * [f7358337d](https://github.com/kubeovn/kube-ovn/commit/f7358337d4d0ff8643bf2221efabfb4f6068bd1f) only support IPv4 snat in vpc-nat-gw when internal subnet is dual (#1747)
 * [99c53d296](https://github.com/kubeovn/kube-ovn/commit/99c53d2965317a9968df15dcb3151a7587bad81a) update README.md
 * [7c4293ebd](https://github.com/kubeovn/kube-ovn/commit/7c4293ebd29ce00ae9d4dcff6aed963c08c58539) docs: update USERS.md (#1743)
 * [a8ed4bce6](https://github.com/kubeovn/kube-ovn/commit/a8ed4bce6a8acceabf6565edabc056ff69f7ef7a) style: import group ordering. (#1742)
 * [bcafb10c8](https://github.com/kubeovn/kube-ovn/commit/bcafb10c89c02a52107e2a909960821ace7b2b46) enqueue subnets after vpc update (#1722)
 * [3af851624](https://github.com/kubeovn/kube-ovn/commit/3af8516249d45d002ba941447e133bb592a5c122) do not add subnet not processed by kube-ovn to vpc (#1735)
 * [7ac7b5927](https://github.com/kubeovn/kube-ovn/commit/7ac7b59272ed4d42df79d1b7dbc9a98bbd96b9c5) dpdk-v2 ，--with-hybrid-dpdk qemu 创建 sock 权限问题 (#1739)
 * [96683cb13](https://github.com/kubeovn/kube-ovn/commit/96683cb1397584a675ed47ae2abe7301156888d9) fix: const extGw may expired, after subnet updated, so use ipam subne… (#1730)
 * [861cf05d4](https://github.com/kubeovn/kube-ovn/commit/861cf05d4bcd33bdebcf5abd39ded6925b9f0110) fix service not working when a node's IPv6 address is before the IPv4 address (#1724)
 * [daaddba2f](https://github.com/kubeovn/kube-ovn/commit/daaddba2f6e16de31a1755dad59a3b85fa0c41f9) update pr template
 * [ce40f7ed0](https://github.com/kubeovn/kube-ovn/commit/ce40f7ed000da1ddfbb75e297eaec7dd9f48b99c) fix: If pod has snat or eip, also need delete staticRoute when delete pod. (#1731)
 * [3ff586ecd](https://github.com/kubeovn/kube-ovn/commit/3ff586ecd80fb984b2361711b4b69856f0ca3360) optimize lrp create for subnet in vpc (#1712)
 * [f582a11bd](https://github.com/kubeovn/kube-ovn/commit/f582a11bd2b0d79c16a1356e47ac27b1717b759c) fix: cancel delete staticroute when it's used by NatRule (#1733)
 * [3c7588bf4](https://github.com/kubeovn/kube-ovn/commit/3c7588bf44008b92ea7313f100b8d6f3dad1f0e1) fix: wrong info when update subnet from dual to ipv4 or ipv6. (#1726)
 * [d1b291ed6](https://github.com/kubeovn/kube-ovn/commit/d1b291ed659b079b0b7fd5e39f8cb7b60f4a43db) change: add newline at end of file (#1717)
 * [26260c917](https://github.com/kubeovn/kube-ovn/commit/26260c91722026c2cdb11f8ea79c4eb15f9672cc) add kernel prerequisite for Rocky Linux 8.6 (#1713)
 * [2fd8e41ee](https://github.com/kubeovn/kube-ovn/commit/2fd8e41eea42945151892aaf94780ed5524ecd4a) Add CODE_STYLE.md (#1711)
 * [3b9111f90](https://github.com/kubeovn/kube-ovn/commit/3b9111f90476bf9cd2fbc436f0716b783bb042c1) Change system-cluster-critical to system-node-critical to prevent pods of DaemonSet type from being (#1709)
 * [e8877f1d7](https://github.com/kubeovn/kube-ovn/commit/e8877f1d7e9c7bc1009ebfe10748c0b70ba568e1) Develop custom vpc-dns (#1662)
 * [ac01a603e](https://github.com/kubeovn/kube-ovn/commit/ac01a603e65eb2655ef8bb61b1c9ac0b91cfcd23) fix CVE-2022-30065 (#1710)
 * [4b13888d0](https://github.com/kubeovn/kube-ovn/commit/4b13888d053d94cd869cd7f115e0cd3ee0233775) fix: add and set ENABLE_KEEP_VM_IP=true to keep vm ip (#1702)
 * [fe18db301](https://github.com/kubeovn/kube-ovn/commit/fe18db301f0281e28c27bd7b6ecd3f427c115337) update CHANGELOG.md
 * [1ab550563](https://github.com/kubeovn/kube-ovn/commit/1ab5505631a53634f1810fd8fa2d100e9797518c) fix overlay MTU in vxlan/stt tunnels (#1693)
 * [c9d9923ef](https://github.com/kubeovn/kube-ovn/commit/c9d9923ef83c1179f54d1d3fbf60f70cdd53065a) fix: response has no gw when create nic without default route (#1703)
 * [b52655fc7](https://github.com/kubeovn/kube-ovn/commit/b52655fc7c202707b881248f2263e7a4a580dbe4) add note in install.sh for install --with-hybrid-dpdk(dpdk-v2). (#1699)
 * [4530a435c](https://github.com/kubeovn/kube-ovn/commit/4530a435c9acad95c0cf8bb5d8dd353bc01dfe02) ignore ovsdb-server/compact error: not storing a duplicate snapshot (#1691)
 * [c4e91cbdf](https://github.com/kubeovn/kube-ovn/commit/c4e91cbdfbaa4c15168a71060bdc7655b7c89671) Get latest vpc data from apiserver instead of cache (#1684)
 * [45bc2f7e9](https://github.com/kubeovn/kube-ovn/commit/45bc2f7e9bfda18b4303b09e9099f2e96f486d48) support kubernetes v1.24 (#1553)
 * [48d914e79](https://github.com/kubeovn/kube-ovn/commit/48d914e7961a3dda8679222ac7addc09b0aec437) update priority range in htb qos (#1688)
 * [41bdcd05c](https://github.com/kubeovn/kube-ovn/commit/41bdcd05c6d95e73d4679490f767c88aec50763b) fix: clean vip eip snat dant fip in cleanup.sh (#1690)
 * [c2cea8850](https://github.com/kubeovn/kube-ovn/commit/c2cea8850d2b9cfd7c136f68120bc139c0a54fa3) update: add Bingbing Zhang into MAINTAINERS (#1687)
 * [a38dbb5d4](https://github.com/kubeovn/kube-ovn/commit/a38dbb5d4d214dd58b89e1eb83fee1dbdcb7aa13) fix: move away words that is considered offensive after k8s v1.20.0 (#1682)
 * [76edad080](https://github.com/kubeovn/kube-ovn/commit/76edad0800b72b2fe3d4bd3eb342ac2123c76f03) update CHANGELOG.md
 * [e4575c882](https://github.com/kubeovn/kube-ovn/commit/e4575c882d34d3d32f0c6c1c607358d8a37d8d9c) add upgrade-ovs script (#1681)
 * [b3c32210d](https://github.com/kubeovn/kube-ovn/commit/b3c32210d5936c45a8c3443b68168ec77cd52ae4) fix: change ovn-ic static route to policy (#1670)
 * [db3b9f7f1](https://github.com/kubeovn/kube-ovn/commit/db3b9f7f1f23f1f6aa7e71dbb4355703d596aa90) Delete pod if subnet of the pod's owner(sts/vm) updated (#1678)
 * [533859fea](https://github.com/kubeovn/kube-ovn/commit/533859fea55b7c62ce46c609f2183ee06d4bf08b) Develop switch lb rule (#1656)
 * [8cb6e36d5](https://github.com/kubeovn/kube-ovn/commit/8cb6e36d5eca4eb6ae1603093ebb1ec76fb181f4) do not snat packets only for subnets with distributed gateway when external traffic policy is set to local (#1616)
 * [86ecfd2d7](https://github.com/kubeovn/kube-ovn/commit/86ecfd2d7586030d0cfb378d4cc4e682a5752163) refactor: extract external routes from eip func, make it the same as … (#1671)
 * [62ddfda6d](https://github.com/kubeovn/kube-ovn/commit/62ddfda6dbeeb708b21f14231203674672ad0052) add loadbalancer service (#1611)
 * [b174b4f45](https://github.com/kubeovn/kube-ovn/commit/b174b4f455ea4a73da38c06d0d2b5963b1bf39f4) bgp: consolidate service check and use service const (#1674)
 * [24786f486](https://github.com/kubeovn/kube-ovn/commit/24786f486ce92efd733fb29e0653b79c820d9d03) security: disable pprof by default (#1672)
 * [3191d8c88](https://github.com/kubeovn/kube-ovn/commit/3191d8c8817a28ead49be803b8f7ac6164cb4f3c) fix bgp: sync service cache (#1673)
 * [f12266827](https://github.com/kubeovn/kube-ovn/commit/f12266827277518de2480e5df22a3dcb6c0b5c14) fix iptables for direct routing (#1578)
 * [2b615f33a](https://github.com/kubeovn/kube-ovn/commit/2b615f33a7f9da403c9a83564b5a8719e28a0b60) feature: support pod use static vip (#1650)
 * [d69024850](https://github.com/kubeovn/kube-ovn/commit/d690248504b7481e1ee7096df6ab24b9e5fb5cb6)  fix: kubectl-ko does't work when ovn-nb, ovn-sb and ovn-northd master slave Switchover (#1669)
 * [b4f890108](https://github.com/kubeovn/kube-ovn/commit/b4f8901089845bb6a63472efe790901ba7991bc0) mount modules for auto load ip6tables moudles (#1665)
 * [72fccd7e2](https://github.com/kubeovn/kube-ovn/commit/72fccd7e23b3745912ae182aa35066dfdad508c2) update docs links
 * [bcec46a3c](https://github.com/kubeovn/kube-ovn/commit/bcec46a3c7c21342831b966d7c81869ae483400b) fix: subnet failed when create without protocol (#1653)
 * [eddce7597](https://github.com/kubeovn/kube-ovn/commit/eddce7597410dabea5c37ebda2399c439c8e64e3) ignore pod not scheduled when reconcile subnet (#1666)
 * [a158f8d6b](https://github.com/kubeovn/kube-ovn/commit/a158f8d6bca4536ba98617c7f92fbc4add756a29) fix libovsdb (#1664)
 * [200c53176](https://github.com/kubeovn/kube-ovn/commit/200c53176770ea0e3a11f3a052a2d7baac68dcbc) fix ovs-ovn not running on newly added nodes (#1661)
 * [4616285c6](https://github.com/kubeovn/kube-ovn/commit/4616285c6161be8790c4af6bbaa9b73645ea6ad8) fix get security group name by external_ids (#1663)
 * [993ae20c8](https://github.com/kubeovn/kube-ovn/commit/993ae20c8c6c5a85be9a09a00a3917e4c439ce2c) fix:can not delete pod with sriov vf (#1654)
 * [9a180f59e](https://github.com/kubeovn/kube-ovn/commit/9a180f59e0fb976931e9e24cfeb6042f962e64bc) add policy route when add subnet (#1655)
 * [a93f211fd](https://github.com/kubeovn/kube-ovn/commit/a93f211fd48677029bac54c83a56446cabaadb66) update CHANGELOG.md
 * [cf1c2017f](https://github.com/kubeovn/kube-ovn/commit/cf1c2017f759ce5eea9fbfdd66750647d3a01f47) fix: no need routed when use v1.multus-cni.io/default-network (#1652)
 * [bcfbe8c6a](https://github.com/kubeovn/kube-ovn/commit/bcfbe8c6a31f0c4db7fc26ea7704e894ca3878b6) docs: add GOVERNANCE.md and SECURITY.md
 * [c2b9eeb4e](https://github.com/kubeovn/kube-ovn/commit/c2b9eeb4ea7b40776b4082d80f7ada210119ac32) fix: should go on check ip after occupied ip (#1649)
 * [9d0cefb5b](https://github.com/kubeovn/kube-ovn/commit/9d0cefb5b3205275a2dbc40fd1c31aed3ebbd829) set ether dst addr for dnat on logical switch (#1512)
 * [a9d5e50d5](https://github.com/kubeovn/kube-ovn/commit/a9d5e50d59f7c63681b170e756dd3fbcdf8d9444) docs: update README.md
 * [d71d1f696](https://github.com/kubeovn/kube-ovn/commit/d71d1f696423c31adc9de0dc4956862df459250d) CI: delete resources in order to avoid a long time waiting for subnet deletions. (#1643)
 * [fe1ebe060](https://github.com/kubeovn/kube-ovn/commit/fe1ebe060a38c2c4400d50bd067854a0230e7c72) ci: fix golangci-lint (#1639)
 * [89514c9c8](https://github.com/kubeovn/kube-ovn/commit/89514c9c834ac86f8a56f183bef70e6425e4aef5) Update install.sh (#1645)
 * [4c8c0a395](https://github.com/kubeovn/kube-ovn/commit/4c8c0a395e05c40a304e8aa99d260206d075dc77) fix: make sure pod annotation switch is the first choice to allocate ip,  and fix vpc nat sts not delete (#1640)
 * [ae8bc1b4c](https://github.com/kubeovn/kube-ovn/commit/ae8bc1b4ca7874e2341e4f0f8eaa4adc587f9c25) docs: update docs link
 * [0f91a61e5](https://github.com/kubeovn/kube-ovn/commit/0f91a61e52d182d455452d5a93ecc4f1e7d6f71c) set networkpolicy log default to false (#1633)
 * [5c857350b](https://github.com/kubeovn/kube-ovn/commit/5c857350b35e3bc92d0b9a804a697ad4fe45ab96) update policy route when join subnet cidr changed (#1638)
 * [9ac8797b3](https://github.com/kubeovn/kube-ovn/commit/9ac8797b327ff64b8f7b9c88c407f424985bf087) fix: diskfull may lead to wrong raft status for ovs db (#1635)
 * [011eba28e](https://github.com/kubeovn/kube-ovn/commit/011eba28ec7ad0760336a65dbf83e7137b209eed) ci: update trivy options (#1637)
 * [3d82780ea](https://github.com/kubeovn/kube-ovn/commit/3d82780ea10f329b14c2e5f6de086697c4fa1b4d) fix no interface report to multus cni, missing in k8s.v1.cni.cncf.io/network[s]-status (#1636)
 * [c953b8f37](https://github.com/kubeovn/kube-ovn/commit/c953b8f37a29f152808ac6ccd6e2a26b709af434) change vp gw pod workload from deployment to statefulset (#1630)
 * [6238695f7](https://github.com/kubeovn/kube-ovn/commit/6238695f74692d1a2fad56485053542e1319a678) increase initial delay of ovs-ovn liveness probe (#1634)
 * [de99e8260](https://github.com/kubeovn/kube-ovn/commit/de99e8260b8f487888299c5d6f366e8541344ad7) fix: cleanup should ignore patch failed (#1626)
 * [8c946a4c0](https://github.com/kubeovn/kube-ovn/commit/8c946a4c01c2addaa472cee6576fe9f4e9a0a889) delete "allow" policy route on subnet deletion (#1628)
 * [75ece0e06](https://github.com/kubeovn/kube-ovn/commit/75ece0e06a719c79440d3b87fe65a8955e0b9453) wait ovn-central pods running before delete ovs-ovn pods (#1627)
 * [f55e32b3e](https://github.com/kubeovn/kube-ovn/commit/f55e32b3eb4f6bdf9cd030bfd2a63011fe19b65d) vip, eip support ipv6 vip count (#1624)
 * [b458bf78c](https://github.com/kubeovn/kube-ovn/commit/b458bf78c733a9089c70c40b5b675e6991c46977) ci: auto changelog now (#1625)
 * [fce2cf172](https://github.com/kubeovn/kube-ovn/commit/fce2cf172f456b21551f4e3fa4110d87a40c3ec6) get dbstatus for all ovn-central pod (#1619)
 * [be66af4be](https://github.com/kubeovn/kube-ovn/commit/be66af4be1fb011ef6da6b8ed71847e08794ed3e) refactor: use ConfigMap resourceVersion to check if ovn-vpc-nat-gw-config changed (#1617)
 * [c54cafa5d](https://github.com/kubeovn/kube-ovn/commit/c54cafa5d4a89ab89c6b2b25500e0a8cc212235e) fix controller exit before process pod update event (#1621)
 * [0ac7e7b98](https://github.com/kubeovn/kube-ovn/commit/0ac7e7b98c9b18cb466957e34c8f28ca9675f19b) docs: update ROADMAP.md
 * [d6915bd64](https://github.com/kubeovn/kube-ovn/commit/d6915bd643b41349b38ef8b7352a9f66e7e3460e) fix acl log name too long (exceed 63) (#1612)
 * [0b840398c](https://github.com/kubeovn/kube-ovn/commit/0b840398c7371aa53a1c05cb0da5f53248fe50dd) docs: Add High-level design of ovn-spekaer (#1609)
 * [7ebf36ddb](https://github.com/kubeovn/kube-ovn/commit/7ebf36ddbe9858a6483225f44b51ebc5dac20fe9) docs: Fix allowed subnets (#1610)
 * [b2f65cd13](https://github.com/kubeovn/kube-ovn/commit/b2f65cd132a07328853f24507f24ecf37dafd617) add cni log Prevent "for loop time" approximately health check time (#1606)
 * [71af95319](https://github.com/kubeovn/kube-ovn/commit/71af95319bffd86e4844c24b76769239efaee3f4) docs:Add Usage of ovn-speaker for passivemode and ebgp-multihop (#1605)
 * [aa456623b](https://github.com/kubeovn/kube-ovn/commit/aa456623b645e5cb75250c4c537da2c3982df246) update static ip docs (#1607)
 * [1400d5b5c](https://github.com/kubeovn/kube-ovn/commit/1400d5b5cf2cfc7a6ff05d9b90225c5506123b6e) Modify the next hop calculation method for kube-ovn-speaker (#1604)
 * [bc86ec86a](https://github.com/kubeovn/kube-ovn/commit/bc86ec86a966d3786424c2ad45637e534f3b382c) fix static ip error in dual stack (#1598)
 * [7fdb4bc68](https://github.com/kubeovn/kube-ovn/commit/7fdb4bc68f207468035e4ef38bff091f561e5869) ci: build amd64 images without avx512 (#1584)
 * [b7453796e](https://github.com/kubeovn/kube-ovn/commit/b7453796e9933f51df23c32bfe60f50b675bc7b5) Add ebgp-multihop function for kube-ovn-speaker (#1601)
 * [6dd2f0ae4](https://github.com/kubeovn/kube-ovn/commit/6dd2f0ae44e4485ca4e2a9bb32502f88de347a1b) monitor dns in cilium e2e (#1597)
 * [75cfe4147](https://github.com/kubeovn/kube-ovn/commit/75cfe4147b40f525a8f3087f3e3a3459bb79937d) Add passivemode for kube-ovn-speaker (#1600)
 * [05bddc6ae](https://github.com/kubeovn/kube-ovn/commit/05bddc6aeed819fa7bc9d2ac202d038d4e457c15) Bump github.com/emicklei/go-restful/v3 from 3.7.4 to 3.8.0 (#1599)
 * [5d54ade92](https://github.com/kubeovn/kube-ovn/commit/5d54ade92bca16eddc1d84c699d5f02a38f0c0cc) docs: fix the kind name (#1593)
 * [4053d46b2](https://github.com/kubeovn/kube-ovn/commit/4053d46b2d9c00edcdc873ad7f4384145149d1bb) Support CNI VESION command (#1596)
 * [34dc29d6b](https://github.com/kubeovn/kube-ovn/commit/34dc29d6b40923370a1725340593ef7530cd0162) update ovs health check, delete connection to ovn sb db (#1588)
 * [d891a84b5](https://github.com/kubeovn/kube-ovn/commit/d891a84b5fee73df91c41a4c10526bcaddd0983c) fix ovn-ic doc err (#1590)
 * [5507440b7](https://github.com/kubeovn/kube-ovn/commit/5507440b79602cf5aa233c441ec3fe4629064982) fix: all cluster pod will be in podadd queue (#1587)
 * [990e291e7](https://github.com/kubeovn/kube-ovn/commit/990e291e7c6f34fdffc9d6f21d2f2531c9be8a26) feat: add args for gc/inspect interval (#1572)
 * [ea2686bf1](https://github.com/kubeovn/kube-ovn/commit/ea2686bf1b0d5dc9429b8808d2b23ccc18886a98) fix: Do not Recreate Logical_Router_Port when Vpc recreated (#1570)
 * [c651a2fce](https://github.com/kubeovn/kube-ovn/commit/c651a2fce98f39581157416445ecb7b0e2a3810f) optimized initialization and gc for the chassis (#1511)
 * [f5d5b0beb](https://github.com/kubeovn/kube-ovn/commit/f5d5b0bebdda54461d046d6fdcbdaa021a78c0fc) fix pod could not be ready (#1562)
 * [337f6e05a](https://github.com/kubeovn/kube-ovn/commit/337f6e05a6b12277eafd86689a7d6604a9cc44ae) Fix incorrect usage info of 'argExternalGatewayNet' (#1567)
 * [d5535f2ba](https://github.com/kubeovn/kube-ovn/commit/d5535f2ba43c0e9a40998aa9e76876fd63a6215a) fix: delete pod panic when delete vm or statefulset. (#1565)
 * [a83305061](https://github.com/kubeovn/kube-ovn/commit/a83305061a91d2729b5251fc4779578d2844408b) fix: clean CRDs introduced by new vpc-nat-gateway (#1563)
 * [fb1f59b10](https://github.com/kubeovn/kube-ovn/commit/fb1f59b107b43ab5583f4c7cfa453ff2d17459e3) do not gc vm pod lsp when vm still exists (#1558) (#1561)
 * [a4c2ef3a9](https://github.com/kubeovn/kube-ovn/commit/a4c2ef3a9d264d3f11a171d7e7599c862968c8e4) do not delete static routes on controller startup (#1560)
 * [a13cd1e29](https://github.com/kubeovn/kube-ovn/commit/a13cd1e29eea78b0d01a8000a5fe40176a8a21e3) update alpine to v3.16 (#1559)
 * [8b9138ae0](https://github.com/kubeovn/kube-ovn/commit/8b9138ae0e27c348117c18ed8e39a0e0345cb515) fix VPC document (#1554)
 * [308e9ecd0](https://github.com/kubeovn/kube-ovn/commit/308e9ecd0d3130d076aa219a0cb0f66fc1dc0033) replace ovn-nbctl daemon with libovsdb in frequent operations (#1544)
 * [434837746](https://github.com/kubeovn/kube-ovn/commit/4348377461487f286c2bfc8a2c0e6e72bd517ba4) fix exec cmd in vpc nat gateway (#1556)
 * [5c17eeaf6](https://github.com/kubeovn/kube-ovn/commit/5c17eeaf68344d43648d8058d09e4c783c8e66bf) CNI: do not return route if nic is not eth0 (#1555)
 * [ed8ed00eb](https://github.com/kubeovn/kube-ovn/commit/ed8ed00ebfeac0bd68e79b2daf15cdc51e1dcc28) do not nat packets for incoming traffic when (#1552)
 * [8f8734bdb](https://github.com/kubeovn/kube-ovn/commit/8f8734bdbd0b53365d85d0ceb2ae2ab728af20fe) add kubeovn 1.9.2 charts (#1539)
 * [d0d9f4eae](https://github.com/kubeovn/kube-ovn/commit/d0d9f4eae467e2b63ae21bbdbf0beef6b5bec73f) fix: opt kubectl-ko install solution (#1550)
 * [9d950f672](https://github.com/kubeovn/kube-ovn/commit/9d950f67247c82862620055d12c12779c6061c28) always set mac address to sriov vf (#1551)
 * [4f1e81219](https://github.com/kubeovn/kube-ovn/commit/4f1e81219f844de841fb0f3dd069c0f9827046ef) use leases for leader election (#1529)
 * [ae56306b4](https://github.com/kubeovn/kube-ovn/commit/ae56306b4e5bf6e40469dc30b485771283abb9ae) fix: fix db-check bug (#1541)
 * [f6b24444b](https://github.com/kubeovn/kube-ovn/commit/f6b24444bb7384057ebbf486a9a232ae55d3d2ff) bump version to v1.11.0 (#1545)
 * [24791f450](https://github.com/kubeovn/kube-ovn/commit/24791f450dbdef54261d1b436a796d15a5258352) exit kube-ovn-controller on stopped leading (#1536)
 * [39e5f0a12](https://github.com/kubeovn/kube-ovn/commit/39e5f0a12219526e6011a636915c5876dc000853) fix: update check script for restart ovs-ovn after rebuild ovsdb (#1534)
 * [c84a14c65](https://github.com/kubeovn/kube-ovn/commit/c84a14c652afdf75f78087d7d1be1e5e582d5368) tmp cancel cilium external svc test (#1531)
 * [5d173010d](https://github.com/kubeovn/kube-ovn/commit/5d173010d8707effffdae3fd2d9ddd25a382249c) remove name for default drop acl in networkpolicy (#1522)
 * [16d28f755](https://github.com/kubeovn/kube-ovn/commit/16d28f755b22704427c297918c01119955ed6e6d) release 1.10.0
 * [bcdb33886](https://github.com/kubeovn/kube-ovn/commit/bcdb338864fd35bf43110a97e8515cd0373d64d3) use inc-engine/recompute instead of deprecated recompute (#1528)
 * [120947669](https://github.com/kubeovn/kube-ovn/commit/1209476696394494982d64ce294580e1751b51fd) update kind to v0.13.0 (#1530)
 * [673138f28](https://github.com/kubeovn/kube-ovn/commit/673138f284a26b26535ce14450a6192c3cd77077) move dumb-init from base images to kube-ovn image (#1527)
 * [ad6826d9f](https://github.com/kubeovn/kube-ovn/commit/ad6826d9f1a4b883e1881870d3a535144fa5b286) fix installing dumb-init in arm64 image (#1525)
 * [4eebabc1e](https://github.com/kubeovn/kube-ovn/commit/4eebabc1e18a69a80412dc98991a8290f6e89a4f) optimize ovs request in cni (#1518)
 * [7a3f73d56](https://github.com/kubeovn/kube-ovn/commit/7a3f73d566e358bf9ba328d6f290927ccd5369b7) optimize node port-group check (#1514)
 * [b7c01d438](https://github.com/kubeovn/kube-ovn/commit/b7c01d438e92af1a9eeccd90a2ebb55d9462c4b9) logic optimization (#1521)
 * [65ee71b4b](https://github.com/kubeovn/kube-ovn/commit/65ee71b4ba1f3a77d146ed43aa27cd60371f69af) fix defunct ovn-nbctl daemon (#1523)
 * [ebe003701](https://github.com/kubeovn/kube-ovn/commit/ebe00370173585ca2428968bda978796e30132e5) fix arm image (#1524)
 * [354d6c3ef](https://github.com/kubeovn/kube-ovn/commit/354d6c3ef8592e3d7506c3e9cea3be0ca1559bdc) fix: keep vm's and statefulset's ips when user specified subnet (#1520)
 * [6021e5288](https://github.com/kubeovn/kube-ovn/commit/6021e5288e28caaf506da3906fb48dda1337b0c8) feature: add doc for tunning packages (#1513)
 * [8e72f2e1f](https://github.com/kubeovn/kube-ovn/commit/8e72f2e1ff7f0cac1a2984e7fc9b40e54bc77a7a) add document for windows support (#1515)
 * [d7ef43b3e](https://github.com/kubeovn/kube-ovn/commit/d7ef43b3e8e0916877b19aa4b351c06adf718102) reduce ovs-ovn restart downtime (#1516)
 * [7b8aa1241](https://github.com/kubeovn/kube-ovn/commit/7b8aa12410c986eed7d5e41aea969abff81dabf1) finish basic windows support (#1463)
 * [ecc8268fe](https://github.com/kubeovn/kube-ovn/commit/ecc8268fe25706962ce1b33eb73c65f342339f2b) refactor logical router routes (#1500)
 * [516036241](https://github.com/kubeovn/kube-ovn/commit/51603624190f0271b73979ed13a0436faa4fb58e) add netem qos when create pod (#1510)
 * [5158dd9d2](https://github.com/kubeovn/kube-ovn/commit/5158dd9d2d96d2e19f8826b403f1c4a6d5299ce6) handle the case of error node cidr (#1509)
 * [1285b0398](https://github.com/kubeovn/kube-ovn/commit/1285b03983a8add91e8442a1fef4211691df0594) fix: ovs trace flow always ends with controller action (#1508)
 * [694286902](https://github.com/kubeovn/kube-ovn/commit/694286902e595ee61f39b1ba78c94944a82e6a7c) add qos e2e test (#1505)
 * [f214ee202](https://github.com/kubeovn/kube-ovn/commit/f214ee202e02591b7ed23320839b203a162dbf4b) optimize IPAM initialization (#1498)
 * [367d6b746](https://github.com/kubeovn/kube-ovn/commit/367d6b74612c5ce12bb5ecea0accb0dc2ef5dcdf) test: fix flaky test (#1506)
 * [79ad4fcf4](https://github.com/kubeovn/kube-ovn/commit/79ad4fcf44066a8be2809bb6f3991fba84b972a1) docs: update README.md
 * [85d09ccd0](https://github.com/kubeovn/kube-ovn/commit/85d09ccd05dd74e437bd5cbd937ac3ce36262c0c) synchronize yamls with installation script (#1504)
 * [63dc5219c](https://github.com/kubeovn/kube-ovn/commit/63dc5219cbfc40d09ddb4d5c7737f27a424e4dc0) feature: svc of multiple clusters (#1491)
 * [011eacf63](https://github.com/kubeovn/kube-ovn/commit/011eacf63b901c93a8c8b65cc7d7bbc42a616d78) use OVS branch-2.17 (#1495)
 * [afc9ef622](https://github.com/kubeovn/kube-ovn/commit/afc9ef62295711399f576ec7d79b43d39fef9723) Update USERS.md (#1496)
 * [b057404ba](https://github.com/kubeovn/kube-ovn/commit/b057404baa655afba6db3ce57244dcb1c2f8f142) update document for mellanox hardware offload (#1494)
 * [fb3c3e6e8](https://github.com/kubeovn/kube-ovn/commit/fb3c3e6e8a7e26c848f2ccac786c0a4ec78f29ad) Feature iptables eip nats splits (#1437)
 * [0c95402ed](https://github.com/kubeovn/kube-ovn/commit/0c95402edd68fc95bf4c973bed49a6ce1f274254) Update USERS.md (#1493)
 * [08a7d5b61](https://github.com/kubeovn/kube-ovn/commit/08a7d5b61ed59973e193dd86adc35ff3d08613d4) update github actions (#1489)
 * [ad28dca06](https://github.com/kubeovn/kube-ovn/commit/ad28dca06fef8b1f31bc448caea4e4566070c50a) update USER.md (#1492)
 * [0db632268](https://github.com/kubeovn/kube-ovn/commit/0db63226817ff865570c203ddb3c57ca66b610fc) fix: add empty chassis check in ovn db (#1484)
 * [d631f8f8b](https://github.com/kubeovn/kube-ovn/commit/d631f8f8b8838846184eb678dc2c934377697258) feat: lsp forwarding external Layer-2 packets (#1487)
 * [d4d700ecd](https://github.com/kubeovn/kube-ovn/commit/d4d700ecdfb020a4bbb12851a5023edc36c5dbc6) base: add back kubectl (#1485)
 * [59e4ae738](https://github.com/kubeovn/kube-ovn/commit/59e4ae73879a1d6cfa95905df66d5cbb02a6fab8) delete ipam record when gc lsp (#1483)
 * [73405b2ad](https://github.com/kubeovn/kube-ovn/commit/73405b2ad2577ae5ec42521b4c827e91954ee4fd) fix: wrong vpc-nat-gateway arm image (#1482)
 * [881622d47](https://github.com/kubeovn/kube-ovn/commit/881622d47dbe61340d99614620464d421c7613cc) fix pod annotation may override by patch (#1480)
 * [e772ee95e](https://github.com/kubeovn/kube-ovn/commit/e772ee95ecdc96e8c3c6fc5fbb35d54a3d4671f5) add acl doc (#1476)
 * [6ef72e75d](https://github.com/kubeovn/kube-ovn/commit/6ef72e75db2309c1090fc58306e400cc938fff47) fix: workqueue_depth should show count not rate (#1478)
 * [5ba5c5264](https://github.com/kubeovn/kube-ovn/commit/5ba5c5264c28e9c59e5d67977588163af0a073be) add delete ovs pods after restore nb db (#1474)
 * [73f9d15fc](https://github.com/kubeovn/kube-ovn/commit/73f9d15fcc1bb7a32a1e137a3c26deffffa5fbde) delete monitor noexecute toleration (#1473)
 * [abaebea47](https://github.com/kubeovn/kube-ovn/commit/abaebea4790d7b9490eb5fa8a962fc4dd3302031) add env-check (#1464)
 * [1d6d46532](https://github.com/kubeovn/kube-ovn/commit/1d6d46532690f8e85e6726939233ab9a65c413a1) Support kubevirt vm live migrate for pod static ip (#1468)
 * [54cab3aa2](https://github.com/kubeovn/kube-ovn/commit/54cab3aa2f0bd2c5ca28fe883ff50afbc8ee802a) fix routes for packets from Pods to other nodes
 * [ba8c5937e](https://github.com/kubeovn/kube-ovn/commit/ba8c5937e8e4205a5f19768d739472033866666e) add manual compile method for ubuntu20.04 (#1461)
 * [7848d71fb](https://github.com/kubeovn/kube-ovn/commit/7848d71fbf415d1b43a0ecd4cdd8cef760efcb9d) append metrics (#1465)
 * [4f0b19766](https://github.com/kubeovn/kube-ovn/commit/4f0b197663d3582f0a1861591f262adf6b31e880) Annotation network_type always is geneve
 * [6ddba02af](https://github.com/kubeovn/kube-ovn/commit/6ddba02af015ee74a3d6d195dcde0efd3eee3081) masquerade packets from Pods to service IP
 * [3d18b8d3f](https://github.com/kubeovn/kube-ovn/commit/3d18b8d3f3322a50d36358e328c15df3338e2dad) update OVS and OVN for windows
 * [39cdfc5c0](https://github.com/kubeovn/kube-ovn/commit/39cdfc5c0df09d2e6114b03730ba77209e95f426) windows support for cni server
 * [75d8f4de3](https://github.com/kubeovn/kube-ovn/commit/75d8f4de335dd7e259192b8c0aca7e9bdeae924f) add kube-ovn-controller switch for EIP and SNAT
 * [8ac3e0c01](https://github.com/kubeovn/kube-ovn/commit/8ac3e0c019645a138ed6e79dd3988fc69d587589) docs: add USERS.md (#1454)
 * [8c214bc91](https://github.com/kubeovn/kube-ovn/commit/8c214bc91129481a69c4d815d8b832183d9165ec) update topology pic
 * [cd5c591cd](https://github.com/kubeovn/kube-ovn/commit/cd5c591cd521628dfe7aba4a3bd503b689977ed3) feature: add sb/nb db check bash script (#1441)
 * [fc5f7190a](https://github.com/kubeovn/kube-ovn/commit/fc5f7190ae5c9363422abba098c0746d44c4a632) add routed check in circulation (#1446)
 * [aa7565193](https://github.com/kubeovn/kube-ovn/commit/aa756519386eabb2225c2d669c706f14e8bbf6c1) modify init ipam by ip crd only for sts pod (#1448)
 * [3a5ead6d7](https://github.com/kubeovn/kube-ovn/commit/3a5ead6d751b9a0cbd3626270691dfb6acc0c46d) base: refactor ovn/ovs build (#1444)
 * [430511668](https://github.com/kubeovn/kube-ovn/commit/43051166883c0f8ade64add8b1049267fc4b578b) log: show the reason if get gw node failed (#1443)
 * [8f1e85ae6](https://github.com/kubeovn/kube-ovn/commit/8f1e85ae6fd4f2f89319edbee91c9e42eadb57c7) add doc for #1358 (#1440)
 * [0c0a03081](https://github.com/kubeovn/kube-ovn/commit/0c0a03081965e9642136a54a0e5f67158d5016ab) prepare windows support for cni server
 * [88b074984](https://github.com/kubeovn/kube-ovn/commit/88b0749846bb1ca49480fee75b6661313e4dc69d) modify webhook img to independent image (#1442)
 * [3dbfa4de2](https://github.com/kubeovn/kube-ovn/commit/3dbfa4de2899ef8d219085e07a3ab96f1c5e2b09) update alpine to fix CVE-2022-1271
 * [03af744f1](https://github.com/kubeovn/kube-ovn/commit/03af744f11e2b8686d774ae11c35752fec7085d2) fix adding key to delete Pod queue
 * [0ea24dcf2](https://github.com/kubeovn/kube-ovn/commit/0ea24dcf234502e0ca5d7104a7fe6549183a2137) fix IPAM initialization
 * [b26a06e7a](https://github.com/kubeovn/kube-ovn/commit/b26a06e7aacd9790008bc6b2a0d6c54042f51ecb) temporary cancel the external2cluater  e2e test for cilium (#1428)
 * [94bc20878](https://github.com/kubeovn/kube-ovn/commit/94bc20878860979aa3d4aaad1cbc0222a212e9a4) ignore all link local unicast addresses/routes
 * [9be57346b](https://github.com/kubeovn/kube-ovn/commit/9be57346b2388214adfe45c74703ba561418a825) fix error handling for netlink.AddrDel
 * [87164cc95](https://github.com/kubeovn/kube-ovn/commit/87164cc9531926eda12e52adfec5b2595ae04114) replace pod name when create ip crd (#1425)
 * [e7c69ba58](https://github.com/kubeovn/kube-ovn/commit/e7c69ba58d7cc16e68eec872c71e2f493e6474e0) add webhook vaildate the vpc resource whether can be deleted. (#1423)
 * [c9a58886a](https://github.com/kubeovn/kube-ovn/commit/c9a58886a818ac14d85cad42b722c9ae5535d11c) We are looking forward to your PR! (#1422)
 * [743ce241e](https://github.com/kubeovn/kube-ovn/commit/743ce241e1497245d1b70791c87e76940415b19a) support alloc static ip from any subnet after ns supports multi subnets (#1417)
 * [d3f6431f2](https://github.com/kubeovn/kube-ovn/commit/d3f6431f234b3310b4cc0f9604c36415ab404288) fix provider-networks status
 * [48e0c4ed7](https://github.com/kubeovn/kube-ovn/commit/48e0c4ed78d701f63d8f6fd2e6439086df387116) build ovs/ovn for windows in ci
 * [3b4ac99ac](https://github.com/kubeovn/kube-ovn/commit/3b4ac99ac9dc2ee11b44502ddd59808f12603a54) cilium e2e: deploy k8s without kube-proxy
 * [902315ed5](https://github.com/kubeovn/kube-ovn/commit/902315ed50a9699ae52ba6ec715eb500666861c8) windows support for CNI
 * [f2baa2f7f](https://github.com/kubeovn/kube-ovn/commit/f2baa2f7fd634f8a1da4eae2d0d5e550f75fee90) add simple e2e for multus integration
 * [e3693436c](https://github.com/kubeovn/kube-ovn/commit/e3693436c7972455452f525f66fb068115189306) update e2e testing
 * [60bf81a35](https://github.com/kubeovn/kube-ovn/commit/60bf81a35ecee8a0a5405d5dd39a040e0685ff39) recover ips CR on IPAM initialization
 * [8e1cd4687](https://github.com/kubeovn/kube-ovn/commit/8e1cd4687a00f17321c9dfe5870dca60b558354b) docs: update ROADMAP.md and MAINTAINERS
 * [19ecaeee2](https://github.com/kubeovn/kube-ovn/commit/19ecaeee27d3052802195ba4a85900bd5be99664) create ip crd in kube-ovn-controller (#1413)
 * [25abbce7d](https://github.com/kubeovn/kube-ovn/commit/25abbce7d83d37ef755e059c161fc84888a41088) add condition for triggering the deletion of redundant chassises in sbdb (#1411)
 * [a378fad24](https://github.com/kubeovn/kube-ovn/commit/a378fad2469da2916d254227bf9a0e682bcbb78f) fix: do not recreate port for terminating pods (#1409)
 * [9587ad41a](https://github.com/kubeovn/kube-ovn/commit/9587ad41a96f319f2dbfad17c8df8a6da2f7e21c) update cni version to 1.0
 * [df83c5fb7](https://github.com/kubeovn/kube-ovn/commit/df83c5fb7bf3d9376b9ba3ce1fa22e6e44b61ce9) update underlay environment requirements
 * [ff695aa36](https://github.com/kubeovn/kube-ovn/commit/ff695aa36a7011b475f33107a104226f2ca38b95) avoid frequent ipset update
 * [f475736c6](https://github.com/kubeovn/kube-ovn/commit/f475736c6f90b753d4a673e4829a87c80fab596a) add reset for kube-ovn-monitor metrics (#1403)
 * [87d6839dd](https://github.com/kubeovn/kube-ovn/commit/87d6839dda10a5f921d307171c6dec0cb9702607) fix: The underlay physical gateway config by external-gw-addr when use snat&eip (#1400)
 * [d36a0d8d7](https://github.com/kubeovn/kube-ovn/commit/d36a0d8d74ebe568ce55d3f4c21bae7b6f5a9283) add custom acls for subnet (#1395)
 * [3206a7a2a](https://github.com/kubeovn/kube-ovn/commit/3206a7a2ae88a2e03e61f888b92aa433da7c8564) check the cidr format whether is correct (#1396)
 * [a33d519b2](https://github.com/kubeovn/kube-ovn/commit/a33d519b24f66cba7b92ddde1408a0bda2a284ce) optimize docs due to frequently asked question. (#1393)
 * [7bd25c639](https://github.com/kubeovn/kube-ovn/commit/7bd25c639118fc8bcc5f679986c94ac0e7e75cd9) adding IP Protocol enumeration to CRD can reduce the kube-ovn Controller judgment logic (#1391)
 * [dcc7971ae](https://github.com/kubeovn/kube-ovn/commit/dcc7971ae09bc22e46dca4895e12fc50007879ea) change the wechat qcode
 * [677690d51](https://github.com/kubeovn/kube-ovn/commit/677690d51f3e4ecbf8868e752a64f3f356c8eb47) append vm deletion check (#1390)
 * [0d663ebe6](https://github.com/kubeovn/kube-ovn/commit/0d663ebe67d8b206be4137fe9bb629b3f9ebd354) We should handle the case where the subnet protocol is handled (#1373)
 * [7289e87c8](https://github.com/kubeovn/kube-ovn/commit/7289e87c8ab380c7842906ddfd8e5fc0082c22ce) VIP is decoupled from port security (#1389)
 * [12907270b](https://github.com/kubeovn/kube-ovn/commit/12907270bda18366bf591403ed4a8ebde4d69a0f) chore: reduce image size (#1388)
 * [5e108fe87](https://github.com/kubeovn/kube-ovn/commit/5e108fe873eed45f12746f63b473b55b808f523c) docs: update the maintainer and roadmap (#1387)
 * [fe7cbe1ba](https://github.com/kubeovn/kube-ovn/commit/fe7cbe1ba2ffed4adacfceb9368d4302d2e0c233) ci: update kind and k8s
 * [ea60cdf71](https://github.com/kubeovn/kube-ovn/commit/ea60cdf712e064356990bc908501e58959077a44) fix external egress gateway
 * [22cb15c51](https://github.com/kubeovn/kube-ovn/commit/22cb15c513ba94aa25597cbce8b0a396d70a0980) add missing link scope routes in vpc-nat-gateway
 * [5571619da](https://github.com/kubeovn/kube-ovn/commit/5571619da26fe3a45660037444d91a9016a7cb63) update nodeips for restore cmd in ko plugin
 * [33180a1c7](https://github.com/kubeovn/kube-ovn/commit/33180a1c7648500f16dfe55b22bbc7776f4e5115) increase memory limit of ovn-central
 * [aa24894ea](https://github.com/kubeovn/kube-ovn/commit/aa24894ea3b35fc7de50213344c001506b1bc7f8) fix range loop
 * [1f24d64d9](https://github.com/kubeovn/kube-ovn/commit/1f24d64d942e0655417f7d4be16d4e5dee98b7c0) fix probe error
 * [c621853ab](https://github.com/kubeovn/kube-ovn/commit/c621853abfd2de8a0050c59d940f38b253287cb0) update script to add restore plugin cmd
 * [dd4a5e0d6](https://github.com/kubeovn/kube-ovn/commit/dd4a5e0d62a284d347f986cffec2478c577cae2a) support dpdk (#1317)
 * [8ad9e8382](https://github.com/kubeovn/kube-ovn/commit/8ad9e8382b513b02d117be09301f1a38bddf18b6) Use camel case instead of snake case
 * [9f3426ee8](https://github.com/kubeovn/kube-ovn/commit/9f3426ee82611edc991ed814a1c8cfd24d35e14e) add detail error when failed to create resource
 * [44dae1f70](https://github.com/kubeovn/kube-ovn/commit/44dae1f704ed049126a07e85f53c9a54ddb8ef9e) add restore process for ovn nb db
 * [c4bb24543](https://github.com/kubeovn/kube-ovn/commit/c4bb24543a4b612661e80248d9cd562ee4dbb1c1) add reset porocess for ovs interface metrics
 * [8e8da1958](https://github.com/kubeovn/kube-ovn/commit/8e8da19585cc83ac6f67a4d4841c272c790d3727) fix SNAT/PR on Pod startup
 * [e9a4bd5c7](https://github.com/kubeovn/kube-ovn/commit/e9a4bd5c79823c6e2e67b13a221326da7d95bb51) optimize kube-ovn-monitor yaml
 * [b11ffa31c](https://github.com/kubeovn/kube-ovn/commit/b11ffa31c6f00616af8f70a4c62d2b7b4dc7d289) Update subnet.go
 * [0b43fc804](https://github.com/kubeovn/kube-ovn/commit/0b43fc8042de74cb11bbce8a0823cc048f8449c6) feat: add webhook to check subnet deletion.
 * [218377849](https://github.com/kubeovn/kube-ovn/commit/218377849857dbc53e5023ea658afe2e71deacf6) modify ipam v6 release ip problem
 * [1264684c6](https://github.com/kubeovn/kube-ovn/commit/1264684c69ddaea9e72e4a1cf2f57e50714e0013) skip ping gateway for pods during live migration
 * [0da84f83f](https://github.com/kubeovn/kube-ovn/commit/0da84f83f481db3fd3d597750ef0c891cd6b6c25) don't check conflict for migration pod with only static mac
 * [89aa2413d](https://github.com/kubeovn/kube-ovn/commit/89aa2413d9f6f6d4a9c19b9c01416363361a3dd4) add service cidr when init kubeadm
 * [bfcb0331e](https://github.com/kubeovn/kube-ovn/commit/bfcb0331eca84c37b5345247d1372fea0669a8ca) docs: add provide and ns spec for multus crd
 * [4f987b10a](https://github.com/kubeovn/kube-ovn/commit/4f987b10a203fd28e774cb87e1304e5943d184b8) update flag parse in webhook
 * [7354d0c30](https://github.com/kubeovn/kube-ovn/commit/7354d0c3005092660c6074a1ac75e0297f9d320f) fix usage of ovn commands
 * [ffd5c8448](https://github.com/kubeovn/kube-ovn/commit/ffd5c844854efc6b18f01e9a64ad872609260f63) add check for pod update process
 * [fe7a6e030](https://github.com/kubeovn/kube-ovn/commit/fe7a6e030947a874dd747b269450ff3682666804) log: rotate all logs in kube-ovn-cni and add compress
 * [024d1684b](https://github.com/kubeovn/kube-ovn/commit/024d1684b7e2ce9328621b516c56e14755930f3d) keep ip for kubevirt pod
 * [8c0b358d0](https://github.com/kubeovn/kube-ovn/commit/8c0b358d08c61e26c913e60695420c5549378280) docs: add integration with Corigine OVS offload
 * [07c531208](https://github.com/kubeovn/kube-ovn/commit/07c531208c806542457a94219ebc78b9c1f6d16f) fix OVS bridge with bond port in mode 6
 * [baeb3af41](https://github.com/kubeovn/kube-ovn/commit/baeb3af415464b9f53bf865f6fc65d49b0e0e4b3) fix: continue of deletion for del pod failed when can't found vpc or subnet (#1335)
 * [8e204be47](https://github.com/kubeovn/kube-ovn/commit/8e204be4759804dd90bf89c1e403fa83154f136f) feat: support DHCP
 * [8393f322f](https://github.com/kubeovn/kube-ovn/commit/8393f322f4ff5feccaa40d11d11112e59af50cf3) Fix usage of ovn commands
 * [bb7b5e56b](https://github.com/kubeovn/kube-ovn/commit/bb7b5e56b0c37a5437c4617e82caf9e8734bc09d) resync provider network status periodically
 * [62642ea8e](https://github.com/kubeovn/kube-ovn/commit/62642ea8efc56294e7a350e70dc8e58de9e9bc28) Revert "resync provider network status periodically"
 * [6ba89e8c0](https://github.com/kubeovn/kube-ovn/commit/6ba89e8c0ea9ebad9917932979d0feeab9e075a6) use const instead the string
 * [d8ba8d038](https://github.com/kubeovn/kube-ovn/commit/d8ba8d038ec1db6bdd04c06ae51f2964c4674799) when update gateway info, we should append old to new deploy
 * [cc124556d](https://github.com/kubeovn/kube-ovn/commit/cc124556d9b0b264531fb24358ae45008e56aef6) resync provider network status periodically
 * [c53b28b1e](https://github.com/kubeovn/kube-ovn/commit/c53b28b1e5d8a25ecc4e966343e6f28ca7dacee9) fix underlay subnet in custom VPC
 * [c4a807b1d](https://github.com/kubeovn/kube-ovn/commit/c4a807b1d4c2b7df9852b3f1e74c93365ef6ebaa) fix ips update
 * [3269bad93](https://github.com/kubeovn/kube-ovn/commit/3269bad932a630415addb233eefc888d2760a9ba) kube-ovn CNI配置文件名字可配置 (#1318)
 * [491abaa88](https://github.com/kubeovn/kube-ovn/commit/491abaa88e2b8c88b6e81967886ad69df33f32ab) delete the logic of repeated enqueueing
 * [31c0b0759](https://github.com/kubeovn/kube-ovn/commit/31c0b07597e059e4bb2f4e67ae1a8dd3ef44e4ff) add log to file, update upgrade script
 * [61c5ebb83](https://github.com/kubeovn/kube-ovn/commit/61c5ebb8399be1323fcdcea6e7c7c8e2b2797bc7) Temporarily comment out the compile and upload of the centos8 compile container.
 * [aef6595f5](https://github.com/kubeovn/kube-ovn/commit/aef6595f58cba068e5f61ae0b1f29f15f9e4fbb3) Revert "Temporarily comment out the compile and upload of the centos8 compile…"
 * [79a268738](https://github.com/kubeovn/kube-ovn/commit/79a26873882398f33e44a7d9a8926e02438b16e7) Temporarily comment out the compile and upload of the centos8 compile container.
 * [1fd27d7c0](https://github.com/kubeovn/kube-ovn/commit/1fd27d7c036a9d06681c5bea4105f66ae2cc747e) feat: add webhook for subnet update validation
 * [6ab8e3698](https://github.com/kubeovn/kube-ovn/commit/6ab8e36980f02baa86164a2aa3f971f3e885a2c1) optimized decision logic
 * [af0baa0ca](https://github.com/kubeovn/kube-ovn/commit/af0baa0ca66e5bcc7143dfd747b88098f2db4f03) Use camel case instead of snake case
 * [b6764e0bc](https://github.com/kubeovn/kube-ovn/commit/b6764e0bc6f5c9effad18a689a275f5894732cda) append add cidr and excludeIps annotation for namespace
 * [a34bb3538](https://github.com/kubeovn/kube-ovn/commit/a34bb353881285a897b31469bbd8faab0a40a3e1) feat: vpc peering connection
 * [9c5556c80](https://github.com/kubeovn/kube-ovn/commit/9c5556c80ba9bf5dfb70c1e7c6bf331539cdea3e) Remove excess code
 * [273eb844b](https://github.com/kubeovn/kube-ovn/commit/273eb844be70a7332ad2f6422ee0c521c4765ec6) chore: show install options when installing (#1293)
 * [d5e342c06](https://github.com/kubeovn/kube-ovn/commit/d5e342c068a743b4a940bf983d0a36b41786616c) feat: update provider network via node annotation
 * [e9c9b1cef](https://github.com/kubeovn/kube-ovn/commit/e9c9b1cef55107e5fd0f6af75ad68d1d77c8cf4c) add container compile and insmod
 * [a90b06a8d](https://github.com/kubeovn/kube-ovn/commit/a90b06a8d8f8077ca15e0a6d767cde35d489c303) add policy route for centralized subnet
 * [2a39f793b](https://github.com/kubeovn/kube-ovn/commit/2a39f793b6674548628e075ee3a3972d1b1b069a) fix: replace ecmp dp_hash with hash by src_ip (#1289)
 * [0fd564e40](https://github.com/kubeovn/kube-ovn/commit/0fd564e400452193a9a299a318b728efe3aad828) Use go to rerimplement ovn-is-leader.sh (#1243)
 * [432c4070e](https://github.com/kubeovn/kube-ovn/commit/432c4070e966ba3a22b59fae2a6417603f071815) fix: only log matched svc with np (#1287)
 * [cb1a698a2](https://github.com/kubeovn/kube-ovn/commit/cb1a698a254c2d2b1f53fe0fa9d68d1cb2b82790) feat: Replace command health check with k8s tcpSocket check (#1251)
 * [b220f0c6e](https://github.com/kubeovn/kube-ovn/commit/b220f0c6ee0652f9b677ecd2c4bafea60a9b8162) add 'virtual' port for vip (#1278)
 * [36c43c486](https://github.com/kubeovn/kube-ovn/commit/36c43c48653cf3782f03fe373b253e24f6e96ec2) skip the missing of kube-dns (#1286)
 * [dad0ef626](https://github.com/kubeovn/kube-ovn/commit/dad0ef62615fda516ac1ccab615aa9b16c9b9657) fix: check if taint exists before un-taint
 * [9365a62d3](https://github.com/kubeovn/kube-ovn/commit/9365a62d39dabd3d3aba802d39482d5fbede103e) add policy route for distributed subnet in default vpc
 * [a5ca73c8a](https://github.com/kubeovn/kube-ovn/commit/a5ca73c8a88519265a90f1b23be0e69b2bdcc102) ci: add retry to fix flaky test
 * [4fdca7146](https://github.com/kubeovn/kube-ovn/commit/4fdca714654b9265b8e20549693b49bdbb0d0087) set up tunnel correctly in hybrid mode
 * [7f8f322ba](https://github.com/kubeovn/kube-ovn/commit/7f8f322bac7740c9092695a76540b22609cd2563) check static route conflict
 * [e7bf87b89](https://github.com/kubeovn/kube-ovn/commit/e7bf87b89f2ebb235246c2a03acb636b31d8e833) fix: https://github.com/kubeovn/kube-ovn/issues/1271#issue-1108813998
 * [017e51252](https://github.com/kubeovn/kube-ovn/commit/017e5125207a5d276c8c0a6437eec03eb47f1482) transfer IP/route earlier in OVS startup
 * [ee2ccf1b9](https://github.com/kubeovn/kube-ovn/commit/ee2ccf1b93193e0cdc7fee64251e68d6e4f135cd) delete unused constant
 * [4022bd577](https://github.com/kubeovn/kube-ovn/commit/4022bd577cbe142a264c8f7544332711c271d95f) add metric for ovn nb/sb db status
 * [fdcc833a3](https://github.com/kubeovn/kube-ovn/commit/fdcc833a3e7a1478f1c0eac44cc3668dfd1ac5d1) add gateway check after update subnet
 * [f40e26ad7](https://github.com/kubeovn/kube-ovn/commit/f40e26ad78c375f131c0cbe8c7f4c77fd32449fb) we should first see if a condition is not going to be met
 * [3ae628cb8](https://github.com/kubeovn/kube-ovn/commit/3ae628cb8bec67852712d2f854afcc918acd53d1) add judge before use slices index
 * [47625c52c](https://github.com/kubeovn/kube-ovn/commit/47625c52c1d8262ded65671a1c325aeef2980caf) prevent multiple namespace reconcile
 * [4455c8692](https://github.com/kubeovn/kube-ovn/commit/4455c8692e306db226d2779df9bc6a3a74c51839) prevent multiple namespace reconcile
 * [6b60a5876](https://github.com/kubeovn/kube-ovn/commit/6b60a5876caacc68273fb858e0f0b408c34858fd) fix: validate statefulset pod by name
 * [fa02cb216](https://github.com/kubeovn/kube-ovn/commit/fa02cb2161b1d7ec8312569d5b84998fbb72aaca) fix golang and base image versions
 * [f210b9340](https://github.com/kubeovn/kube-ovn/commit/f210b93403240a13cbe8d2a704ba0338d088dd79) add back centralized subnet active-standby mode
 * [2557c5167](https://github.com/kubeovn/kube-ovn/commit/2557c51670b091d950859dbabcf2a660bf8ebb96) support to add multiple subnets for a namespace
 * [c230ed8a1](https://github.com/kubeovn/kube-ovn/commit/c230ed8a1b80181e055d6fb5d5e11a329166b79c) prepare for next release
 * [f95a90eb3](https://github.com/kubeovn/kube-ovn/commit/f95a90eb3ee579d01069bf610fcd184d70f22c4e) Support only configure static mac_address
 * [e4d48df38](https://github.com/kubeovn/kube-ovn/commit/e4d48df38d6ed16acb77d92d66686df7d40f55ea) prepare for release 1.9.0
 * [c830594dc](https://github.com/kubeovn/kube-ovn/commit/c830594dc5b9575531b34eea358bad019d0ff3a5) fix: liveMigration with IPv6
 * [e52b68976](https://github.com/kubeovn/kube-ovn/commit/e52b689764cad9166f6b499b1757b3e76ee4a765) update networkpolicy port process
 * [851ad0ce6](https://github.com/kubeovn/kube-ovn/commit/851ad0ce6874fe2ec1dff3ecb8ded3079ce27f18) Add args to configure port ln-ovn-external
 * [5d95d6285](https://github.com/kubeovn/kube-ovn/commit/5d95d62857b3a4ffdcabe2b8ae945d48d9ef1249) update check for delete statefulset pod
 * [695f45320](https://github.com/kubeovn/kube-ovn/commit/695f45320200a29948264edac201675892cd8e4d) ignore hostnetwork pod when initipam
 * [4b98d15fb](https://github.com/kubeovn/kube-ovn/commit/4b98d15fb07ab23adefae6fb23aee365e15db18a) kubectl-ko: support trace Pods being created
 * [63bc25ea8](https://github.com/kubeovn/kube-ovn/commit/63bc25ea84da3adafc16bc9a1467adb6930aa9b1) add dnsutils for base image
 * [6318d0049](https://github.com/kubeovn/kube-ovn/commit/6318d004990743c87d88ff7885d01bb1e36fd858) Add new arg to configure ns of ExternalGatewayConfig
 * [715229204](https://github.com/kubeovn/kube-ovn/commit/71522920498015847c12455728fc38d17eaab5b5) update scripts for 1.8.2
 * [960f02c15](https://github.com/kubeovn/kube-ovn/commit/960f02c15bfb2665d9d2589713a6fbbab9958a69) Optimized decision logic
 * [8974f6a37](https://github.com/kubeovn/kube-ovn/commit/8974f6a3712d14d6fd5674b449dd62f903c22f98) add svc cidr in ovs LB for optimization
 * [0192a9ae8](https://github.com/kubeovn/kube-ovn/commit/0192a9ae8851b14a81fbf6dad7b4ebb006a4c71e) add doc for gateway pod in default vpc
 * [1f9dc754c](https://github.com/kubeovn/kube-ovn/commit/1f9dc754c9fdebe089699a07fec2fc3e76e1dc12) optimize log for node port-group
 * [36d6b00a8](https://github.com/kubeovn/kube-ovn/commit/36d6b00a8fae03b1df57116befc3b101a5f348dd) fix iptables rules and service e2e
 * [8dc938d83](https://github.com/kubeovn/kube-ovn/commit/8dc938d83d923ad1ffdcbbf72e559dc9497ddeed) add kubectl-ko to docker image
 * [c4cc8f0d9](https://github.com/kubeovn/kube-ovn/commit/c4cc8f0d9b43d9ed4b5f21879263eab1e235cd61) fix: invalid syntax error
 * [a4f4cb490](https://github.com/kubeovn/kube-ovn/commit/a4f4cb490ff257063f859cead802013018517563) fix pod tolerations
 * [8611de822](https://github.com/kubeovn/kube-ovn/commit/8611de8229214b732c858175bc2822a25bfcd02b) modify pod's process of update for use multus cni as default cni
 * [5ab83ba42](https://github.com/kubeovn/kube-ovn/commit/5ab83ba42b41e4dbab912be345a42adc0fdefdd1) fix installation script
 * [09ef9be09](https://github.com/kubeovn/kube-ovn/commit/09ef9be09c3fc1809dc0d09209502cbe156c7682) add log for ecmp route
 * [791b00f42](https://github.com/kubeovn/kube-ovn/commit/791b00f42c660f1afbb0ce1d4714344558774e12) fix: ipv6 traffic still go into ct
 * [55e6a8ca3](https://github.com/kubeovn/kube-ovn/commit/55e6a8ca326c73a62c91bfd6622a200d8366d1e7) append check for centralized subnet nat process
 * [58a44fb2b](https://github.com/kubeovn/kube-ovn/commit/58a44fb2b7d90086af77095aa20cfcc48e353a36) move chassis judge to the end of node processing
 * [9f0c42fae](https://github.com/kubeovn/kube-ovn/commit/9f0c42fae734db095a12921555eb8ae140bac192) change nbctl args 'wait=sb' to 'no-wait'
 * [6f3567055](https://github.com/kubeovn/kube-ovn/commit/6f35670556dc8fa924d0b41198433e8afc78084a) use different ip crd with provider suffix for pod multus nic
 * [f7b595dcf](https://github.com/kubeovn/kube-ovn/commit/f7b595dcf67fc2222c387058b88c811e6ed1116d) fix service cidr in dual stack cluster
 * [c510b4397](https://github.com/kubeovn/kube-ovn/commit/c510b43972bb8ee217e05856c8e0022f6c93c86b) add healthcheck cmd to probe live and ready
 * [e14bc40c5](https://github.com/kubeovn/kube-ovn/commit/e14bc40c512324808deab974d2235fa5a61b5ba1) delete frequently log
 * [bde98e757](https://github.com/kubeovn/kube-ovn/commit/bde98e7571069d50caead0e4b7968d4e2947feb4) support running ovn-ic e2e on macOS
 * [727ea53a8](https://github.com/kubeovn/kube-ovn/commit/727ea53a809a1846c6d701cab5bca26b151a9313) pinger: fix getting empty PodIPs
 * [205a0c021](https://github.com/kubeovn/kube-ovn/commit/205a0c021e3be8d8df0931bbbbba29f392fa8220) fix cni deepcopy
 * [650ea6d3c](https://github.com/kubeovn/kube-ovn/commit/650ea6d3c5693b9e11f1b016c936b4a95a1b11a4) add cilium e2e
 * [46ba84eef](https://github.com/kubeovn/kube-ovn/commit/46ba84eefebd9fc15e6749c253b4dd4fd8cf1840) filter used qos when delete qos
 * [1de284ebc](https://github.com/kubeovn/kube-ovn/commit/1de284ebc8504e53cd4a2ec0388d407731226671) add protocol check when subnet is dual-stack
 * [1f4a247dd](https://github.com/kubeovn/kube-ovn/commit/1f4a247ddf8801f32b5b177dd8408a5ea0827a60) lint: make go-lint happy
 * [91f3fa4b9](https://github.com/kubeovn/kube-ovn/commit/91f3fa4b9b715faa09ad7e52f547e59bb9b920a3) some fixes
 * [d57bc1d72](https://github.com/kubeovn/kube-ovn/commit/d57bc1d72e66cfa0822cdb73b216af74fe4bb7d7) compatible with OVN 20.06
 * [9116425a4](https://github.com/kubeovn/kube-ovn/commit/9116425a42fdb3fb6cc8812d04cc0aa900c0b385) use multus-cni as default cni to assign ip
 * [d18323a49](https://github.com/kubeovn/kube-ovn/commit/d18323a4906626eb2dbb51a5d81085cbc34c86b9) some fixes
 * [668c21256](https://github.com/kubeovn/kube-ovn/commit/668c2125646cc1c2d29893457885dc2d245170d7) perf: jemalloc and ISA optimization
 * [5c08d28da](https://github.com/kubeovn/kube-ovn/commit/5c08d28da0eb7fa6981e71ca65c59f8dedeb42ed) fix: check np switch
 * [365715556](https://github.com/kubeovn/kube-ovn/commit/36571555676739e4e2e6949cd37d13999b5d9175) fix: port security
 * [e713bdf0f](https://github.com/kubeovn/kube-ovn/commit/e713bdf0f281415536eeab623385fd534e694d6c) fix nat rule
 * [d8e84cf06](https://github.com/kubeovn/kube-ovn/commit/d8e84cf06fc232ae0e954dfb20e3aa08ce0de30e) When netpol is added to a workload, the workload's POD can be accessed using service
 * [51365b41d](https://github.com/kubeovn/kube-ovn/commit/51365b41d1745f45c401943b75589ad1779efaf4) when update subnet's execpt ip,we should filter repeat ip
 * [5aacec592](https://github.com/kubeovn/kube-ovn/commit/5aacec592b8f63c03cf3b2d41b0292a48a084a62) update wechat image
 * [6c8fa978b](https://github.com/kubeovn/kube-ovn/commit/6c8fa978b008a75727ce97e1d326d2b5a7c096df) fix: do not reuse released ip after subnet updated
 * [e4648cc81](https://github.com/kubeovn/kube-ovn/commit/e4648cc81c2a8cc9bf11157379ee5e99d5c97e0c) update: update 1.7-1.8 script
 * [b1f8332c6](https://github.com/kubeovn/kube-ovn/commit/b1f8332c68770fe1db8a1dae2bdb84ef601bc195) perf: do not send traffic to ct if not designate to svc
 * [178cf7b87](https://github.com/kubeovn/kube-ovn/commit/178cf7b87fbaecd4868620c7fdbd9e5a3fa89145) fix: add back the leader check
 * [7be43c974](https://github.com/kubeovn/kube-ovn/commit/7be43c9748ba4470d7d40c63f39fd52352b0913a) fix port_security
 * [e596c3c4b](https://github.com/kubeovn/kube-ovn/commit/e596c3c4b0f5118db38bb308dce9a21bf25ea952) sync live migration vm port
 * [e8b1ff5b4](https://github.com/kubeovn/kube-ovn/commit/e8b1ff5b4ec877840c1647c4544b23c0a38ac252) docs: add f5 ces integration docs
 * [7058d5680](https://github.com/kubeovn/kube-ovn/commit/7058d5680f9338e3d8ee0dd711a535ab2faba368) update Go modules
 * [84dbb102b](https://github.com/kubeovn/kube-ovn/commit/84dbb102b3782abe57f6cd66e51abff9f219ac4e) update delete operation for statefulset pod
 * [e9e2c9111](https://github.com/kubeovn/kube-ovn/commit/e9e2c9111e26b5c7b4814104ba94a76949b3f548) chore: update klog to v2 which embed log rotation
 * [fafd55552](https://github.com/kubeovn/kube-ovn/commit/fafd5555284ee8d766109f1a785cf043cd4e1715) fix: add kube-ovn-cni prob timeout
 * [490590a46](https://github.com/kubeovn/kube-ovn/commit/490590a4694b753f9733559492bc8380b4a2680a) append add db compact for nb and sb db
 * [4fb302f5e](https://github.com/kubeovn/kube-ovn/commit/4fb302f5ebbed2b363c3bafb9bb6157328641969) deleting all chassises which are not nodes
 * [c49a74040](https://github.com/kubeovn/kube-ovn/commit/c49a740401323cc09d2e09e6d8abc00bbdca827b) add db compact for nb and sb db
 * [3b7ec06c7](https://github.com/kubeovn/kube-ovn/commit/3b7ec06c766d59598a7a178f376f85682a4be84d) add vendor param for fix list LR
 * [ae23d3dfd](https://github.com/kubeovn/kube-ovn/commit/ae23d3dfd500820cb2c5ff79baca58f234f1efb6) fix LB: skip service without cluster IP
 * [df3d3977b](https://github.com/kubeovn/kube-ovn/commit/df3d3977b9f19941060d0d7a4af63ef99c4ef494) add webhook with cert-manager issued certificate
 * [2be11269a](https://github.com/kubeovn/kube-ovn/commit/2be11269a9646f4cc124b8c6f2178dd4fd289bbe) security: update base ubuntu image
 * [eb3647176](https://github.com/kubeovn/kube-ovn/commit/eb3647176ddbed8f61aedb688d6d2c4800e60446) add pod in default vpc to node port-group
 * [ea300d2bf](https://github.com/kubeovn/kube-ovn/commit/ea300d2bf179d58f51465c5e767d487e606f7894) fix pinger's compatibility for k8s v1.16
 * [3837b0a23](https://github.com/kubeovn/kube-ovn/commit/3837b0a231659728d51d3ad64570c565d66bfed4) check IPv4 gateway by resolving gateway MAC in underlay subnets
 * [75604b5d9](https://github.com/kubeovn/kube-ovn/commit/75604b5d94ef5ff9ada5f5533ebaa157894adf8c) add nodeSelector for vpc-nat-gateway pod
 * [fac6c725a](https://github.com/kubeovn/kube-ovn/commit/fac6c725acf646e5f3e5fd3fc3a797402b7a10df) do not send multicast packets to conntrack
 * [c3004bbc2](https://github.com/kubeovn/kube-ovn/commit/c3004bbc2aa9c7d020fead19da390144c6d18162) Revert "support to set NB_Global option mcast_privileged"
 * [2802b94dd](https://github.com/kubeovn/kube-ovn/commit/2802b94ddb002ee434236e564ce18a6e23350850) add ip address for lsp
 * [28a93927a](https://github.com/kubeovn/kube-ovn/commit/28a93927af69358c35116e68a53ebe703c1326ae) fix: no need to set address for ls to lr port
 * [2048007a2](https://github.com/kubeovn/kube-ovn/commit/2048007a2c8a0e7fa5f691639782c69ce7e3aae1) add sg acl check when init
 * [b9abee715](https://github.com/kubeovn/kube-ovn/commit/b9abee71542b934f36f6cfb569f46d0bc2f79201) cleanup command flags
 * [54a3b913e](https://github.com/kubeovn/kube-ovn/commit/54a3b913ef3c1481318512e3aac3fd81b21bcab8) replace port-group named address-set with port-group since there's no ip set for lsp when create lsp
 * [743502cd7](https://github.com/kubeovn/kube-ovn/commit/743502cd744c3c397665b60e81bfdd2817e5beeb) support to set NB_Global option mcast_privileged
 * [a5f0256a9](https://github.com/kubeovn/kube-ovn/commit/a5f0256a979041adfead4c7ed47b72dd0b828a72) add networkpolicy support for attachment cni
 * [45f64bfaf](https://github.com/kubeovn/kube-ovn/commit/45f64bfaf83ff131f88b0121f8e5951d8b852b5e) add process for pod attachment nic with subnet in default vpc
 * [49e9197e5](https://github.com/kubeovn/kube-ovn/commit/49e9197e57db258a9baebcdbc3ba997599e27b6d) fix security group
 * [60e896f89](https://github.com/kubeovn/kube-ovn/commit/60e896f8967e0d9a3f4f0f1fa1c3a73083ca5b14) fix the duplicate call about strings.Split
 * [c9f5f4b46](https://github.com/kubeovn/kube-ovn/commit/c9f5f4b46cf9969816b471fec1fa0fb79fa86933) deepcopy fix steps
 * [e0cb19aa7](https://github.com/kubeovn/kube-ovn/commit/e0cb19aa73f95bba16872bd532ef4b6fd0d0a1c2) fix: do not nat route traffic
 * [4e4d95d5f](https://github.com/kubeovn/kube-ovn/commit/4e4d95d5f31f5c36af2c8d137e39081b6128ec50) fix: Skip MAC address Settings when PCI addresse is unavailable
 * [adce05c74](https://github.com/kubeovn/kube-ovn/commit/adce05c74ea6e89bee155e15b1ff43a1ab3c170e) add ovn-ic e2e
 * [3b6b5034d](https://github.com/kubeovn/kube-ovn/commit/3b6b5034d9900d820e93ec08a3f18c5e3298bc88) other CNI can be used as the default network
 * [841f907b0](https://github.com/kubeovn/kube-ovn/commit/841f907b0d0e8b40e9a0d186cafc33152e5b818a) fix: move macvlan binary to host
 * [52ec0af47](https://github.com/kubeovn/kube-ovn/commit/52ec0af4721478bb4d0bc4d5aff09960891ca9a2) Revert "ci: init kind cluster before build finish"
 * [a85993254](https://github.com/kubeovn/kube-ovn/commit/a859932540fca9fc5a9b8d862e43f22cf04291e2) fix ko trace
 * [1dd66a770](https://github.com/kubeovn/kube-ovn/commit/1dd66a770612544342c9af0d0d962f23a818640a) add ovn-ic HA deploy
 * [bc3ce0bbf](https://github.com/kubeovn/kube-ovn/commit/bc3ce0bbf11c0b6bcb9d621e0cf5c9d94e1b775a) fix node address set name
 * [cbed28204](https://github.com/kubeovn/kube-ovn/commit/cbed282046713526f15d55ef6bba920051237f64) update cni init image
 * [a648bfc6a](https://github.com/kubeovn/kube-ovn/commit/a648bfc6a23f6ced4c5688b1689e772ea2511e64) chore: update kind k8s to 1.22 and remove pre 1.16 support
 * [a1d56e973](https://github.com/kubeovn/kube-ovn/commit/a1d56e9737b256da80da09b3aef83efb28f75945) do not set bridge-nf-call-iptables
 * [738c76126](https://github.com/kubeovn/kube-ovn/commit/738c76126e5f2eea8c57073436ab6330d7b31029) use logical router policy for accessing node
 * [6719ee242](https://github.com/kubeovn/kube-ovn/commit/6719ee242b798352a99d12d26b034f58aa1fe401) ci: init kind cluster before build finish
 * [61817bf41](https://github.com/kubeovn/kube-ovn/commit/61817bf417278b553aeec981c84275efac6a623c) reduce qos query with ovs-vsctl cmd
 * [1776c4471](https://github.com/kubeovn/kube-ovn/commit/1776c4471447ec2e34fb373019a2c55ff82e1f1b) fix read-only pointer in vlan and provider-network
 * [329228d4a](https://github.com/kubeovn/kube-ovn/commit/329228d4a9b11e45bfe122e44089f1f4f66ca9f9) fix: trace in custom vpc
 * [a9c0a4aa4](https://github.com/kubeovn/kube-ovn/commit/a9c0a4aa499da7c12991e7bfaa1047291f4e7ada) fix read-only pointer in vlan and provider-network
 * [62df34162](https://github.com/kubeovn/kube-ovn/commit/62df34162ef72710a3e532d8fb489cbb78ffc594) update docs
 * [a546ba954](https://github.com/kubeovn/kube-ovn/commit/a546ba9545d22bb2d366e94e84e4fb5914964dc5) fix LB in dual stack cluster
 * [eb63f72ec](https://github.com/kubeovn/kube-ovn/commit/eb63f72ec1db5cd6a5394f8717bb13514fe965f1) fix: check allocated annotation in update handler
 * [55b8b8ac2](https://github.com/kubeovn/kube-ovn/commit/55b8b8ac2d543eb99055451b2327d9fe19478049) support using logical gateway in underlay subnet
 * [ef424d731](https://github.com/kubeovn/kube-ovn/commit/ef424d731a07c791943b587ca0e65014dec439df) docs: optimize cilium integration docs
 * [a09e84d0c](https://github.com/kubeovn/kube-ovn/commit/a09e84d0c4ac028ce6de7944a60ff1bca0fed3c3) fix: ensure all kube-ovn components deleted before annotate pods
 * [e7aeb96ec](https://github.com/kubeovn/kube-ovn/commit/e7aeb96ec789ca45ac1f10ef1945f01cb17dd26e) fix bug: logical switch ts not ready
 * [dc4e693f1](https://github.com/kubeovn/kube-ovn/commit/dc4e693f1f90396a9094ce4e3e0136176a5c4b01) Fix unpopulated CPU charts
 * [003723e58](https://github.com/kubeovn/kube-ovn/commit/003723e5881364a27afb44667c2bfb4359c2c2e7) Revert "get default subnet"
 * [418feb1ba](https://github.com/kubeovn/kube-ovn/commit/418feb1ba8ea96c804265b6966f0dba28d81debd) add htbqoses rbac
 * [850e42186](https://github.com/kubeovn/kube-ovn/commit/850e42186e373c1456a8c612a10cf143d248cf67) feat: pod can use multiple nic with the same subnet
 * [5840d5093](https://github.com/kubeovn/kube-ovn/commit/5840d5093ff6c98710849ba078b02350b4419a78) add error detail
 * [e6377caef](https://github.com/kubeovn/kube-ovn/commit/e6377caefbe46869f3000fef05affe993aa053ca) add check switch for default subnet's gateway
 * [b5b6c3267](https://github.com/kubeovn/kube-ovn/commit/b5b6c32678a696efff4e4b10fb84b40b327b5a7d) get default subnet
 * [fbafca410](https://github.com/kubeovn/kube-ovn/commit/fbafca410eb4f79b0252029d1dee117aef8d1dd9) remove node chassis annotation on cleanup
 * [348eaf367](https://github.com/kubeovn/kube-ovn/commit/348eaf3670720181c1e5a82ca10a1bc02a3c1732) update: add 1.7 to 1.8 update scripts
 * [f934613d3](https://github.com/kubeovn/kube-ovn/commit/f934613d36f0e756201df03ac5e1bd1e5cf23f9f) base: add macvlan to help vpc setup
 * [cd1dda1e5](https://github.com/kubeovn/kube-ovn/commit/cd1dda1e5449f6debd5a6e95e87bf25067bb825d) fix: delete vpc-nat-gw deployment
 * [50eddac3e](https://github.com/kubeovn/kube-ovn/commit/50eddac3e95b97136399647b390dcc614257f7fa) ko: check ovsdb storage status
 * [20670e877](https://github.com/kubeovn/kube-ovn/commit/20670e8772d46e1a9a29714d2e18d24727a655b8) fix cleanup.sh and uninstall.sh
 * [b31c4d19f](https://github.com/kubeovn/kube-ovn/commit/b31c4d19f4761a846167d352c4606fd612b1441c) use constant instead a string
 * [86f63f26d](https://github.com/kubeovn/kube-ovn/commit/86f63f26dbc319e4b65263029ac140f8e4466ccd) fix: check and load ip_tables module
 * [3bfd82b7e](https://github.com/kubeovn/kube-ovn/commit/3bfd82b7e08ca9cf1694649c7dd43d5c4234aa62) fix: multus-cni subnet allocation
 * [e5ed1ace0](https://github.com/kubeovn/kube-ovn/commit/e5ed1ace04509669b2331d26fb8797f645275cde) docs: add svg
 * [17ff6c55f](https://github.com/kubeovn/kube-ovn/commit/17ff6c55f97c484cd37408db678224b2f3b4d648) chore: update install
 * [ce97b94c8](https://github.com/kubeovn/kube-ovn/commit/ce97b94c8e7c4c8f18b3702a1d4b757e7720259f) integrate Cilium into Kube-OVN
 * [fda0c17b4](https://github.com/kubeovn/kube-ovn/commit/fda0c17b4066f9d439a85b5ca47f9cc526db614d) fix kubectl-ko diagnose
 * [3f8a2b0ed](https://github.com/kubeovn/kube-ovn/commit/3f8a2b0ed4551290e047829e197f99a5750ea46d) change inspection logic from manually adding lsp to just readding pod queue
 * [01ca82f9c](https://github.com/kubeovn/kube-ovn/commit/01ca82f9ce006dbae6a418fe1bfe8fa8b139d7f0) fix pinger in dual stack cluster
 * [0ba64dea9](https://github.com/kubeovn/kube-ovn/commit/0ba64dea916c1898411962c9fc0567724252b863) add e2e testing for dual stack underlay
 * [7f27a05d5](https://github.com/kubeovn/kube-ovn/commit/7f27a05d53690aff302593b5e8a29fea68116df1) fix pinger and monitor in underlay networking
 * [6a56f8bb5](https://github.com/kubeovn/kube-ovn/commit/6a56f8bb5694336d6758f4567ca6b0220e256fa0) fix kubectl plugin ko
 * [2c9fe438d](https://github.com/kubeovn/kube-ovn/commit/2c9fe438de87f1ef2748224782db96f54247c3dc) adjust the location of the log
 * [86ee933aa](https://github.com/kubeovn/kube-ovn/commit/86ee933aa16c747a56cee523092d030afc042bc3) ci: push vpc-nat-gateway
 * [f459ca975](https://github.com/kubeovn/kube-ovn/commit/f459ca975fc3d31edb95d3f703b6ec892c85e320) replace api for get lsp id by name
 * [0a533984f](https://github.com/kubeovn/kube-ovn/commit/0a533984f0238c21594f334f3153308a5a61f22d) docs：revise vpc.md
 * [788478991](https://github.com/kubeovn/kube-ovn/commit/788478991bdcf5f5b3a9ad6fe883a0de8e50bf32) grafana: optimize grafana dashboard
 * [168a7c976](https://github.com/kubeovn/kube-ovn/commit/168a7c9768a3660ce2c1c1472721b324ef2c1798) In netpol egress rules, except rule should be set to != and should not be ==
 * [d7edf24bf](https://github.com/kubeovn/kube-ovn/commit/d7edf24bf98d8229e4e71827e7b5d0f26c7b295a) ci: add vpc-nat-gateway build
 * [5cd32df80](https://github.com/kubeovn/kube-ovn/commit/5cd32df809a895d0df5b8bae1f0ce4e1d772c31a) Update OVN to version 21.06
 * [dd36d61c3](https://github.com/kubeovn/kube-ovn/commit/dd36d61c3b70e20111d8246dcd7eb9b4a770bf05) modify kube-ovn as multus-cni problem
 * [d17f61513](https://github.com/kubeovn/kube-ovn/commit/d17f6151310e7024a2d8a04e89abdff586648832) support to set htb qos priority
 * [c20e01110](https://github.com/kubeovn/kube-ovn/commit/c20e01110e0f2dcc138aac2f13793815fafd5ae2) perf: add fastpath module for 4.x kernel
 * [ff5d3df3f](https://github.com/kubeovn/kube-ovn/commit/ff5d3df3f898834709cd73dafbd898e2ceb8b04d) add inspection
 * [3e9f9a99c](https://github.com/kubeovn/kube-ovn/commit/3e9f9a99c0bde8fc1857e7a3d73a0a542231d4b5) perf: add stt section and update benchmark
 * [d38423271](https://github.com/kubeovn/kube-ovn/commit/d3842327175880b73a09bfb86a95ddbaf876ebbf) feat: optimize log
 * [4c6c29a36](https://github.com/kubeovn/kube-ovn/commit/4c6c29a36a0300eac7646c590cabe144df915e35) fix: init node with wrong ipamkey and lead conflict
 * [47255a105](https://github.com/kubeovn/kube-ovn/commit/47255a10591ab69ff8566b33ece98875f6f1fff1) fix installation scripts
 * [fd7454870](https://github.com/kubeovn/kube-ovn/commit/fd745487040ff067d2b9f161658555201f223857) fix getting LSP UUID by name
 * [1f5719a59](https://github.com/kubeovn/kube-ovn/commit/1f5719a59afc1acf009eb7e3438806c780908bf6) fix StatefulSet down scale
 * [5bccd8453](https://github.com/kubeovn/kube-ovn/commit/5bccd8453419cb4477dc925555f686c1bf3a7705) fix vpc policy route
 * [acb82de0d](https://github.com/kubeovn/kube-ovn/commit/acb82de0dd34152e537036208e147ec86c5fed6a) docs: update roadmap
 * [87f9b863e](https://github.com/kubeovn/kube-ovn/commit/87f9b863ea2e666aa27f1fc99fb6670aff94d6f2) refactor: mute ovn0 ping log and add ping details
 * [a99c4200c](https://github.com/kubeovn/kube-ovn/commit/a99c4200cfcd3ef233e81bb778cfc0eac5356c27) fix: wrong link for iptables
 * [52b01c017](https://github.com/kubeovn/kube-ovn/commit/52b01c017bd11c72b0134d7e092e5cb6636590ff) fix IPAM for StatefulSet
 * [51511e633](https://github.com/kubeovn/kube-ovn/commit/51511e633030eb92b6cc29f28338e890ff10f6d7) append externalIds for pod and node when upgrade
 * [391f7014e](https://github.com/kubeovn/kube-ovn/commit/391f7014e5c617b8f29723b460e23cf54779c4a5) feature: LoadBalancer for custom VPC
 * [7fd8cf44f](https://github.com/kubeovn/kube-ovn/commit/7fd8cf44fb4532fc182652939c0fbed1cf13e8d8) feat: support vip
 * [25f634fbc](https://github.com/kubeovn/kube-ovn/commit/25f634fbcbcf54f8f1dfde6f90aed4cd31b65479) fix VPC document
 * [97a5b2a33](https://github.com/kubeovn/kube-ovn/commit/97a5b2a337a511917f600303cf5e0823a4a7b942) fix init ipam
 * [71fcbf121](https://github.com/kubeovn/kube-ovn/commit/71fcbf12153330db0ddae6a1ba8baed628eafcb1) fix: gc lb
 * [2b154b1a5](https://github.com/kubeovn/kube-ovn/commit/2b154b1a52fcbb873988fb11900120273a1e4163) Update prometheus.md
 * [1e766f9cd](https://github.com/kubeovn/kube-ovn/commit/1e766f9cd76229189a957a78f8cc5d17a3e771cf) feat: support VLAN subnet in VPC
 * [4c013a3e1](https://github.com/kubeovn/kube-ovn/commit/4c013a3e186273d01b33cd0812a95925471f7a1a) ci: push dev image to separate repo
 * [39c8a19c1](https://github.com/kubeovn/kube-ovn/commit/39c8a19c177808f8c7060c30b2cc903cb53f1ee7) fix: kubeclient timeout
 * [edaf41e04](https://github.com/kubeovn/kube-ovn/commit/edaf41e041ef527abd141a4bbafca6cb445c90fa) fix: serialize pod add/delete order
 * [78a77f797](https://github.com/kubeovn/kube-ovn/commit/78a77f797c43b62086817b5a421789cdedbf6191) perf: increase ovn-nb timeout
 * [5937ccbf7](https://github.com/kubeovn/kube-ovn/commit/5937ccbf7aba3d3aee7d5ce49fe7cae1f062a088) fix gc lsp statistic for multiple subnet
 * [c71620ce1](https://github.com/kubeovn/kube-ovn/commit/c71620ce129baf8cc2434129ba965440eebb3ed2) fix: re-check ns annotation to avoid annotations lost
 * [d40d57015](https://github.com/kubeovn/kube-ovn/commit/d40d5701522c6ca43e771556848dcd6191337af6) perf: do not diagnose external access
 * [871c1493b](https://github.com/kubeovn/kube-ovn/commit/871c1493b251792bb8cbf5ed6c8b9dbd7f4aee9c) feature: vpc support policy route
 * [90b1a2ea4](https://github.com/kubeovn/kube-ovn/commit/90b1a2ea4205ebcdcb392760f1af12456543e801) reactor: remove ovn ipam options
 * [7f43f25c1](https://github.com/kubeovn/kube-ovn/commit/7f43f25c1cf2f65e0e658894aa4a5312add2f5ec) perf: switch's router port's addresses to "router"
 * [8dbe8f946](https://github.com/kubeovn/kube-ovn/commit/8dbe8f94678a8f5c7ec6259477d7a9e9dd711057) lint: make staticcheck happy
 * [8ad46dad4](https://github.com/kubeovn/kube-ovn/commit/8ad46dad4f96855fe413b18a1e3d368b418c5c3f) fix e2e testing
 * [5a126378a](https://github.com/kubeovn/kube-ovn/commit/5a126378a3466bb7ce715174c702e0e9454e85b6) prepare for next release
 * [5b70c81d3](https://github.com/kubeovn/kube-ovn/commit/5b70c81d3b7707d8cc1efe3632df28822ecf347c) fix variable referrence
 * [42fed929e](https://github.com/kubeovn/kube-ovn/commit/42fed929e61bed1fc1e0e654d4b379ab9317d601) fix typos
 * [f59aff27e](https://github.com/kubeovn/kube-ovn/commit/f59aff27e39f038eb2a0bca50294f95e072e9a20) refactor: reuse waitNetworkReady to check ovn0 and slightly improve the installation speed
 * [ea723d6dc](https://github.com/kubeovn/kube-ovn/commit/ea723d6dc9afed1e246407f8a34499c481d67a6c) fix nat-outgoing/policy-routing on pod startup
 * [7c5fed654](https://github.com/kubeovn/kube-ovn/commit/7c5fed6547c8695056de3f12a5ca7e0754b37d39) fix adding OVN routes in dual stack Kubernetes
 * [80a037ee5](https://github.com/kubeovn/kube-ovn/commit/80a037ee557cf4c3ffe60845a52ae8f1f00196f8) release: prepare for 1.8
 * [f59bfb864](https://github.com/kubeovn/kube-ovn/commit/f59bfb864a261e8f0174c8940850389d79ae02a1) add update process and adding label to ls/lsp/lr
 * [e09d99b3b](https://github.com/kubeovn/kube-ovn/commit/e09d99b3b4e8e46f7af795fafeaa725bf8bb0d0b) fix: VLAN CIDR conflict check
 * [e6b8341e6](https://github.com/kubeovn/kube-ovn/commit/e6b8341e6b9cf337104aff383ef27279faf51bc7) security: update base image
 * [29422965e](https://github.com/kubeovn/kube-ovn/commit/29422965e5cf0b9a62143c66205f3f5bf059f7e8) update provider network CRD
 * [25b151c85](https://github.com/kubeovn/kube-ovn/commit/25b151c8587eb63778e0f87664ceb54801797687) fix external-vpc
 * [44a8b4f69](https://github.com/kubeovn/kube-ovn/commit/44a8b4f6912b5869c9eb19ec92d3b340fbb57821) perf: use link alias to filter packet
 * [e9984fe0c](https://github.com/kubeovn/kube-ovn/commit/e9984fe0c047f34d0be196b54062a6be2c450504) security: fix CVE-2021-3538
 * [d41c5e9b8](https://github.com/kubeovn/kube-ovn/commit/d41c5e9b840e740db398564d5c2d0b303457804d) add print columns for subnet/vpc/vpc-nat-gw crd
 * [730e4f171](https://github.com/kubeovn/kube-ovn/commit/730e4f17165d853055557e3a9b442655aa21fed5) improve support for dual-stack
 * [2439c86e6](https://github.com/kubeovn/kube-ovn/commit/2439c86e6e4945073f023dc51b3ab49a5fa2be19) feat: suport vm live migration
 * [c148a5ac8](https://github.com/kubeovn/kube-ovn/commit/c148a5ac861f27247847df0114956a510a1162e6) initialize ipsets on cni server startup
 * [10613e872](https://github.com/kubeovn/kube-ovn/commit/10613e8727b33e60ba15f2f57869c4884a029c2a) delete residual ovs internal ports
 * [361d4bbe0](https://github.com/kubeovn/kube-ovn/commit/361d4bbe026784ea4314ee262c12745ad7e6e982) simplify vlan implement
 * [6fde0a566](https://github.com/kubeovn/kube-ovn/commit/6fde0a56686d543eabc6ba5d5368b5ba38adc768) fix: ovn-northd svc flip flop
 * [b11060565](https://github.com/kubeovn/kube-ovn/commit/b110605650819ee4aadce6c21b22c0cddf852f24) add container run command for runtime containerd
 * [42e212ca2](https://github.com/kubeovn/kube-ovn/commit/42e212ca24467174c56e71e9ca497dba6843371c) fix subnet conflict check for node address
 * [3d2c6eb96](https://github.com/kubeovn/kube-ovn/commit/3d2c6eb96e191e3f8346cf97fe086e241bf2f405) feat: read interface in installation from environment
 * [35acf424e](https://github.com/kubeovn/kube-ovn/commit/35acf424ee659eee1ee0872afa36bc5925974988) update encap ip by node annotation periodic
 * [13b2080a0](https://github.com/kubeovn/kube-ovn/commit/13b2080a03c922e766d11ef8a8739296c0689e5b) fix ipset on pod creation/deletion
 * [f415b1ba5](https://github.com/kubeovn/kube-ovn/commit/f415b1ba54836e1750667f40c1c86f529682f6b2) add ready status for provider network
 * [092838495](https://github.com/kubeovn/kube-ovn/commit/0928384955266226aef1119fcd8e90ebf4f34ccb) avoid Pod IP to be the same with node internal IP
 * [70fbbecc4](https://github.com/kubeovn/kube-ovn/commit/70fbbecc43aa9db9d015eef4a8fbacd60d47de68) remove subnet's `spec.underlayGateway` field
 * [96b0c1185](https://github.com/kubeovn/kube-ovn/commit/96b0c11859661bd6c377ae22add2e163aebe3d1d) add support for custom routes
 * [45aafca21](https://github.com/kubeovn/kube-ovn/commit/45aafca2104a94ab56b3ed03a3e8df67abc462bc) Add missing metadata directive in VpcNatGateway example
 * [0380d64ce](https://github.com/kubeovn/kube-ovn/commit/0380d64cea659085aabb1b77dbcc17e3a5a18ac6) use util.hostNameEnv instead KUBE_NODE_NAME
 * [38e04f34a](https://github.com/kubeovn/kube-ovn/commit/38e04f34aaadfbef889fae2f31e5df280e5429b8) chore: change wechat image
 * [5df9fdd4c](https://github.com/kubeovn/kube-ovn/commit/5df9fdd4c7594af4e02c5100e68e74c95eb94c35) fix typo
 * [4a7dd7340](https://github.com/kubeovn/kube-ovn/commit/4a7dd7340dc9e4ae024c4c670c7382d8c2e269ce) perf: add fastpath and tuning guide
 * [3d8cdb6cb](https://github.com/kubeovn/kube-ovn/commit/3d8cdb6cb0177b5212b235bf850fe706fe61eb70) update node labels and provider network's status.readyNodes when provider network is not initialized successfully in a node
 * [8596ddc90](https://github.com/kubeovn/kube-ovn/commit/8596ddc901a55a04474a941f06e1a67f6dd8e644) fix issues in underlay networking
 * [7724990dd](https://github.com/kubeovn/kube-ovn/commit/7724990dda02f599a2b9bb153e42d034591e6059) add external vpc switch
 * [ffef618db](https://github.com/kubeovn/kube-ovn/commit/ffef618db5aed74330627e6336a5b7f5aeaca525) update versions in docs and yamls
 * [6e8d5c80a](https://github.com/kubeovn/kube-ovn/commit/6e8d5c80a809f11f794841ba7881038030548079) update Go to version 1.16
 * [3deb57708](https://github.com/kubeovn/kube-ovn/commit/3deb57708fcf95f0b6a46dfc23c18fb970914341) fix IPv6-related issues
 * [2e4922d56](https://github.com/kubeovn/kube-ovn/commit/2e4922d560da774590694002cb95b8cfa6adad3b) ci: use stable version
 * [dcda11d6e](https://github.com/kubeovn/kube-ovn/commit/dcda11d6e19d9508ecd4a23f8b41857a9fe81fc4) fix: bad udp checksum when access nodeport
 * [f12e5ee58](https://github.com/kubeovn/kube-ovn/commit/f12e5ee587e26c1474351b9f49e8889e870d1680) fix port-security, address parameters should be merged into one
 * [f03d43503](https://github.com/kubeovn/kube-ovn/commit/f03d435038bdc8f7eeb131c0bbd3628a4dcc2af6) docs: optimize description
 * [b5b5bdb89](https://github.com/kubeovn/kube-ovn/commit/b5b5bdb89e717670a3f5f10e2c758148a73bdefc) ensure provider nic is up
 * [b5bbed38a](https://github.com/kubeovn/kube-ovn/commit/b5bbed38a3f175982f9ba25b14a57cd56e7ba611) fix uninstall.sh
 * [3ba5168ce](https://github.com/kubeovn/kube-ovn/commit/3ba5168cec17e251fff40254df057aaa837fdd16) some optimizations
 * [9ae0b3c35](https://github.com/kubeovn/kube-ovn/commit/9ae0b3c351e60539fcdf89cfbff6628e50f76c0a) fix gofmt lint
 * [410d93290](https://github.com/kubeovn/kube-ovn/commit/410d932900eab3184bfddd017c197b189a11a391) fix multi-nic.md
 * [5e9e41ac5](https://github.com/kubeovn/kube-ovn/commit/5e9e41ac5555e4521680221b5bb73dd81fc620cc) fix dual stack cluster created by kind
 * [386d6160f](https://github.com/kubeovn/kube-ovn/commit/386d6160f2970ef0ea498680f96e881f32b4f0a1) remove external egress gateway from additionalPrinterColumns
 * [70ae50efb](https://github.com/kubeovn/kube-ovn/commit/70ae50efb3ef53458b11a519efdf962c1b04a396) fix default bind socket of cni server
 * [56025edea](https://github.com/kubeovn/kube-ovn/commit/56025edea4628b755b6b6116772be1bc1efa0492) if the string of ip is empty,program will die
 * [9492f63f6](https://github.com/kubeovn/kube-ovn/commit/9492f63f6bd590153a17ea30f43509f8208a1e99) if the string of ip  is empty,program will die
 * [324dce2e3](https://github.com/kubeovn/kube-ovn/commit/324dce2e3b38f71461fd9a9b372aacb1cd54ca55) fix underlay networking on node reboot
 * [f7077d58e](https://github.com/kubeovn/kube-ovn/commit/f7077d58ee08a9864725a5b1962c6e9bb54a6033) add judge before use the index about cidrBlocks and ips
 * [f25b1ae2b](https://github.com/kubeovn/kube-ovn/commit/f25b1ae2b544af576842d0bf538b0cf90a9e0d32) add validation check function
 * [bda102a77](https://github.com/kubeovn/kube-ovn/commit/bda102a772483b84903bc4ac7532472e88608a78) docs: add wechat qcode
 * [14ccbeb37](https://github.com/kubeovn/kube-ovn/commit/14ccbeb377578ea18ea69b5c467445f5f08f13e3) feat: security group
 * [992a09d35](https://github.com/kubeovn/kube-ovn/commit/992a09d353a4bdb012aea0a591248b8a523a8f2f) delete subnet AvailableIPs and UsingIPs para
 * [057ade92e](https://github.com/kubeovn/kube-ovn/commit/057ade92edbdca16000598f65dd5afe6fc063c89) fix: panic when node has nil annotations
 * [59869daaf](https://github.com/kubeovn/kube-ovn/commit/59869daaf1246621c9dad409e0549483bf5a4a35) append pod/exec resource for vpc nat gw
 * [3ed2fe26d](https://github.com/kubeovn/kube-ovn/commit/3ed2fe26dfd913155335d2a97414426c8e72ed7e) update comment for SetInterfaceBandwidth
 * [e1caa5941](https://github.com/kubeovn/kube-ovn/commit/e1caa5941e55dd14f5e30dee5446900fb5138766) update qos process
 * [80e5e2bab](https://github.com/kubeovn/kube-ovn/commit/80e5e2babbe01cd15863837bb8ee74087bad9c91) fix LoadBalancer in custom VPCs
 * [bb1146eeb](https://github.com/kubeovn/kube-ovn/commit/bb1146eeb5887f8f3856e929626ef6b939b9e9bf) Support Pod annotations control port mirroring
 * [4c4b09007](https://github.com/kubeovn/kube-ovn/commit/4c4b09007b5d0ce1cbf8a361c30bcf5591a22dcf) fix docs
 * [a04d964dc](https://github.com/kubeovn/kube-ovn/commit/a04d964dcf1ef99a4a8f3ed62449a39a5cecc6f4) externalOvnRouters is ok with 0
 * [9524c93f3](https://github.com/kubeovn/kube-ovn/commit/9524c93f3b55a40518f8198a2d91ac369ebd49f4) delete attachment ips
 * [6dd6a51d7](https://github.com/kubeovn/kube-ovn/commit/6dd6a51d7d53a9d9caa175108e330e8dff43f93e) fix external_ids:pod_netns
 * [cbe8ae689](https://github.com/kubeovn/kube-ovn/commit/cbe8ae689fe0b8118cabe0e2d4fe683602026d3d) add switch for network policy support
 * [dc56d2385](https://github.com/kubeovn/kube-ovn/commit/dc56d23853ff6b0fc38a5f514b029387be0b5787) fix subnet e2e
 * [e3daee830](https://github.com/kubeovn/kube-ovn/commit/e3daee8306ab1d71957f01db5de8890e5d31194c) ignore empty strings when counting lbs
 * [81ce45c2e](https://github.com/kubeovn/kube-ovn/commit/81ce45c2e89a31c064eba54eec46cc729129f25c) fix iptables
 * [e9ea6a0f9](https://github.com/kubeovn/kube-ovn/commit/e9ea6a0f98220453fc3b36a581cc695445f7a503) fix issue #944
 * [1cb57358a](https://github.com/kubeovn/kube-ovn/commit/1cb57358ac0dbd36f7fa9fc6eda71b86f621fd12) fix openstackonkubernetes doc bugs
 * [fcdb0106b](https://github.com/kubeovn/kube-ovn/commit/fcdb0106b263afffba370987763d07f7486d3490) add switch for gateway connectivity check
 * [4dc4624f9](https://github.com/kubeovn/kube-ovn/commit/4dc4624f9bcf75ca296fc9c84c03590163c31505) fix cleanup.sh
 * [4fb974071](https://github.com/kubeovn/kube-ovn/commit/4fb974071601402230e5d5c2a45491ec9fa8df4c) security: fix CVE-2021-33910
 * [41b6429c4](https://github.com/kubeovn/kube-ovn/commit/41b6429c4334bdff5c2c85d17d9775484800958a) delete ecmp route when node is deleted
 * [5bd96ac71](https://github.com/kubeovn/kube-ovn/commit/5bd96ac718927a511f0111c3ceb20350f6a9effb) fix: if nftables not exists do no exit
 * [6c5efbc30](https://github.com/kubeovn/kube-ovn/commit/6c5efbc30e0ad3e7f38b46af6b455b0d8abcb118) update wechat contract method
 * [e449b8eaf](https://github.com/kubeovn/kube-ovn/commit/e449b8eaf156c387c9619016ca9f107a1c1353c4) delete overlapped var subnet
 * [2427a4b3a](https://github.com/kubeovn/kube-ovn/commit/2427a4b3a05471239f29bceb1b9e6f0444fa5788) add designative nat ip process for centralized subnet
 * [1595eac56](https://github.com/kubeovn/kube-ovn/commit/1595eac56740e388b718548204906cd7b5f9a4dc) fix ipsets
 * [7e24e7d6f](https://github.com/kubeovn/kube-ovn/commit/7e24e7d6f78430fee4fd3951adf9bb006ee0f07b) update underlay e2e testing
 * [27c649a50](https://github.com/kubeovn/kube-ovn/commit/27c649a50a40d4c96ed73262bbaaabbc47bc6dcc) match chassis until timeout
 * [df76038a0](https://github.com/kubeovn/kube-ovn/commit/df76038a0a0475a269a3829e14875a9e847a6e45) fix CRD provider-networks.kubeovn.io
 * [d1c7a2ee3](https://github.com/kubeovn/kube-ovn/commit/d1c7a2ee3664ba1b63220e50cb42034a8408e687) fix: set vf mac
 * [949c28c25](https://github.com/kubeovn/kube-ovn/commit/949c28c25d31322e4f7b910c2395911a571c28a4) update qos ingress_policing_burst
 * [8a05bdc88](https://github.com/kubeovn/kube-ovn/commit/8a05bdc88d51c236b8a20040282a808b51ccada2) add field defaultNetworkType in configmap ovn-config
 * [1810dfc32](https://github.com/kubeovn/kube-ovn/commit/1810dfc32f22e42654fab674ae125927e14c512a) keep subnet's vlan empty if not specified
 * [4e28600d4](https://github.com/kubeovn/kube-ovn/commit/4e28600d4a42844832bfb2cb70f30b51dea0b21b) delete ecmp route when node is not ready
 * [d145f5759](https://github.com/kubeovn/kube-ovn/commit/d145f5759a3245dcce407cadcd5271034fe9a224) add del learned routes when remove ovnic
 * [6499e5859](https://github.com/kubeovn/kube-ovn/commit/6499e5859f92c0ba58f266aa308c795a2c52ba3b) [kubectl-ko] support trace in underlay networking
 * [23d84f0a7](https://github.com/kubeovn/kube-ovn/commit/23d84f0a7e8ef9e83777b003986f2e1bbdf11a38) fix subnet available IPs
 * [eced6bacd](https://github.com/kubeovn/kube-ovn/commit/eced6bacdcc66cb523a0ed47271061f8ce654056) fix bug for deleting ovn-ic lrp failed
 * [a4abbb2e8](https://github.com/kubeovn/kube-ovn/commit/a4abbb2e8af05b04e0b5c0702439b7aa5966b019) add node internal ip into ovn-ic advertise blacklist
 * [2ec0aa749](https://github.com/kubeovn/kube-ovn/commit/2ec0aa7494805e472182b52b0c5fb643899e083a) underlay/vlan network refactoring
 * [ead2c65f0](https://github.com/kubeovn/kube-ovn/commit/ead2c65f00263f653dc093abbaf5f2d64710eff3) chore: update ovn to 21.03
 * [651a634d6](https://github.com/kubeovn/kube-ovn/commit/651a634d6b78c9d719947fcaee6aa6cda4ec84a0) security: fix CVE-2021-3121
 * [8cff68519](https://github.com/kubeovn/kube-ovn/commit/8cff685191f6911f477d4632d43a409311213da1) list ls with label to avoid listing ts failure
 * [3fd9c7acd](https://github.com/kubeovn/kube-ovn/commit/3fd9c7acd32e1f361ab59b90db3c58a88d70b8dd) Update log error
 * [0fe67258a](https://github.com/kubeovn/kube-ovn/commit/0fe67258a9e65eaf46671425592ba73d271cbcb9) delete the process of ip crd delete in cni delete request
 * [9049fc725](https://github.com/kubeovn/kube-ovn/commit/9049fc725eb66761461268fde65a6c0a3d7673af) update networkpolicy process
 * [a5b22a21e](https://github.com/kubeovn/kube-ovn/commit/a5b22a21e9ab6a210a5a0e4f24367b4181e55e12) modify func name Additonal to Additional
 * [0cd5dcfec](https://github.com/kubeovn/kube-ovn/commit/0cd5dcfec78517ee738d6620fb5cf3ae18e9eda0) fix uninstall.sh execution in OVS pods
 * [b4ce83a2c](https://github.com/kubeovn/kube-ovn/commit/b4ce83a2cf2df2467cd015682ac2812f8b5dabbc) perf: enable tx offload again as upstream already fix it
 * [9ca47b658](https://github.com/kubeovn/kube-ovn/commit/9ca47b6584a99866e7fee94381285e34df9cd1fa) label lr, ls and lsp, and add label filter when gc
 * [37a045a31](https://github.com/kubeovn/kube-ovn/commit/37a045a31d2581a0373538a04a9362cf7cb158b0) security: add go build security options
 * [bdf91846b](https://github.com/kubeovn/kube-ovn/commit/bdf91846bd646d084487f7637d79a2ff9ec3ab6a) feat: ko support cluster operations status/kick/backup
 * [efdce4645](https://github.com/kubeovn/kube-ovn/commit/efdce464567c2a9705b3f6f70629662a17d48367) docs: update docs about vlan/internal-port/kubeconfig
 * [ced434053](https://github.com/kubeovn/kube-ovn/commit/ced434053d9b5bf280de476b3a93988a3eabb78b) add judge before use slices's index
 * [3d98d7626](https://github.com/kubeovn/kube-ovn/commit/3d98d7626ff7db93927b3c25f80531400e1dceff) update kind to version v0.11.1
 * [e1e63cfaf](https://github.com/kubeovn/kube-ovn/commit/e1e63cfaff83fea058bfb6ad584c849a427c619a) adapt to vfio-pci driver
 * [205f57124](https://github.com/kubeovn/kube-ovn/commit/205f571245c939820faf27e292c383ecaecb7397) fix IP/route transfer on node reboot
 * [a3cac539f](https://github.com/kubeovn/kube-ovn/commit/a3cac539f6bb13d52135641b6132a79098e7d2ca) add master check when a node adding to a cluster and config sb/nb address
 * [b98afeefc](https://github.com/kubeovn/kube-ovn/commit/b98afeefcd415ec6bb236b36bf58ff7797ae8035) update installation scripts
 * [2d750cbf1](https://github.com/kubeovn/kube-ovn/commit/2d750cbf1226dd5e8e6f3ed8d0dc762a3fa59883) enable hw-offload
 * [64b9abae6](https://github.com/kubeovn/kube-ovn/commit/64b9abae6083eb0305cd85d9ff1ddb3631a93cab) do not delete statefulset pod when update pod
 * [4359c1980](https://github.com/kubeovn/kube-ovn/commit/4359c1980534d3b6e38a977ac327b07c5da35cbe) fix: node route should filter out 'vpc'
 * [744e6577b](https://github.com/kubeovn/kube-ovn/commit/744e6577ba7d705ec94acc9a7a53dfd0c33907c6) feat: lb switch
 * [7ec2f994d](https://github.com/kubeovn/kube-ovn/commit/7ec2f994dd06332eef4db6ff4507ef9a0b6929fe) docs: show openstack docs and docker image status
 * [5484387f2](https://github.com/kubeovn/kube-ovn/commit/5484387f24f60af534b55c44bfe1e5dde4c6b0fe) fix: clean up gateway chassis list for external gw
 * [acc95f1d6](https://github.com/kubeovn/kube-ovn/commit/acc95f1d67332abeb15bb67244af15179111a772) add doc for openstack/kubernetes hybrid deploy
 * [e2973c4f0](https://github.com/kubeovn/kube-ovn/commit/e2973c4f04b32528a2769db200e9c9f7c8f56f94) configure OVS internal port after dummy interface
 * [8608b7e5a](https://github.com/kubeovn/kube-ovn/commit/8608b7e5ab97273adfca03ec9de05f1144937aec) some fixes in vlan initialization
 * [872340c87](https://github.com/kubeovn/kube-ovn/commit/872340c87cb2e0db84fcd52f9669b99333a1d21f) clean up vpc service
 * [fde899145](https://github.com/kubeovn/kube-ovn/commit/fde89914587de175191517877dbbd33960bff0ba) feat: vpc load balancer
 * [8ed91be4a](https://github.com/kubeovn/kube-ovn/commit/8ed91be4a2fcfe969a55721ed18ea419120bbdb6) fix: lsp may lost when server pressure is high
 * [42fbe86e7](https://github.com/kubeovn/kube-ovn/commit/42fbe86e7d0d66a3f8e8ebcaa81ef82f2875be58) fix: check crds when controller start
 * [a5fef59b5](https://github.com/kubeovn/kube-ovn/commit/a5fef59b57a840b1248b581b3d00959914055fd2) start evpc ph1
 * [31ee8c104](https://github.com/kubeovn/kube-ovn/commit/31ee8c104396b57c3cc0b9963e8d5f462d5aa691) start evpc ph1
 * [44db142e5](https://github.com/kubeovn/kube-ovn/commit/44db142e511c7d655eee655d98387dd88701bcb8) ci: retry arm build when failed
 * [96c139851](https://github.com/kubeovn/kube-ovn/commit/96c139851d3ddd4f35f90729a9cc813f0d28526e) update ecmp notes
 * [8c169322f](https://github.com/kubeovn/kube-ovn/commit/8c169322fcbd168033cfacaa63a641fb799549c8) add interface name in cni response
 * [aa88e2a21](https://github.com/kubeovn/kube-ovn/commit/aa88e2a21b8749e6531479162aceaf0f174e7cfa) add nicType for offload
 * [eb387428e](https://github.com/kubeovn/kube-ovn/commit/eb387428e114bee62591c6850a2d21a5973f8300) 1.Support to specify node nic name 2.Delete extra blank lines
 * [cb8cc6454](https://github.com/kubeovn/kube-ovn/commit/cb8cc6454f19b55432a8378402cd1f29ae308522) ignore update pod nic annotation when is not nil
 * [3a4347b92](https://github.com/kubeovn/kube-ovn/commit/3a4347b92b343cc9fbc0bd8f841213b681cca9c9) set default UnderlayGateway to true in vlan mode
 * [a0d78920f](https://github.com/kubeovn/kube-ovn/commit/a0d78920fb1aaabd08d38992254be1e7da87bb3c) unify logical entity list funcs (#863)
 * [9e563d847](https://github.com/kubeovn/kube-ovn/commit/9e563d8476fd46c4bd23e0183a0f6f042437c034) ci: remove dpdk ci
 * [e48a0894e](https://github.com/kubeovn/kube-ovn/commit/e48a0894e0e82cd672b3de53682b5b8049cc4601) correct vlan e2e testing
 * [f690085d8](https://github.com/kubeovn/kube-ovn/commit/f690085d89ea26821dba38175da185ba22f951f1) fix: remove rollout check
 * [2b2df3dc1](https://github.com/kubeovn/kube-ovn/commit/2b2df3dc1c925de42b1674c34e56014b0ee0ebf2) adapt internal tcpdump
 * [2531779a4](https://github.com/kubeovn/kube-ovn/commit/2531779a464221bf63c309d043de14096697fab3) update docker buildx install method
 * [eef1b0aa8](https://github.com/kubeovn/kube-ovn/commit/eef1b0aa88408c80cc5120fa5b3299bd0848b053) fix: remove wait ovn sb
 * [2e59e81c9](https://github.com/kubeovn/kube-ovn/commit/2e59e81c9113b217351205e985c264b3acc69a0a) fix: ci issues
 * [df47c489e](https://github.com/kubeovn/kube-ovn/commit/df47c489e86305a92485e98f367f5b1d208e5fae) fix: cleanup kube-ovn-monitor resource
 * [598cffdd7](https://github.com/kubeovn/kube-ovn/commit/598cffdd7f2ee931fd247e96a0d12426d0a5c62f) fix multi-nic.md
 * [f4b75bd04](https://github.com/kubeovn/kube-ovn/commit/f4b75bd04fee12b8e55896ba7f325ef6974e455d) fix: acl overlay issues
 * [2fe4fe1d4](https://github.com/kubeovn/kube-ovn/commit/2fe4fe1d4d4e79c332dea720b85284ae7b8bfd29) ci: split ovn/ovs into base image
 * [db2b7b069](https://github.com/kubeovn/kube-ovn/commit/db2b7b069a1f99a65aeeb69f758a46cd05baddde) add judge before use slices's index
 * [3e259ae90](https://github.com/kubeovn/kube-ovn/commit/3e259ae909877e7e787a2533de325f1d3a24cd16) update version to v1.7 in docs
 * [eb54dc037](https://github.com/kubeovn/kube-ovn/commit/eb54dc037ea75cc58ba01aab97be9d340eeff274) update master version to v1.8.0

### Contributors

 * Abdul Basit
 * Abhishek Pandey
 * Alex Jones
 * Amin Mohammadian
 * Andrei Kvapil
 * Anton Patsev
 * Bryan Lee
 * Caleb Tallquist
 * Chris
 * Congqi Zhao
 * Cookie Wang
 * Copilot
 * Daviddcc
 * DiMalovanyy
 * Fudankenshin
 * Guangyu Suo
 * Hargrove Wang
 * Joachim Hill-Grannec
 * Johann Schley
 * Kaihang Zhang
 * Karol Szwaj
 * Kevin Carter
 * KillMaster9
 * Longchuanzheng
 * Mengxin Liu
 * MengxinLiu
 * Miika Petäjäniemi
 * Money Liu
 * Navid
 * Nico Wang
 * Noah
 * Paul
 * QEDQCD
 * Qinghao Huang
 * Rick
 * Robin Lee
 * Ruijian Zhang
 * SKALA NETWORKS
 * Samuel Liu
 * ShaPoHun
 * Six
 * SkalaNetworks
 * Tobias
 * Usman Malik
 * Vishal Mohan
 * Wang Bo
 * Xiaobo Liu
 * Zespre Chang
 * Zhao Congqi
 * amoy-xuhao
 * andrewlee1089
 * aoc55
 * azee
 * bob199x
 * bobz965
 * bogdan-cehash
 * caohuilong
 * carezkh
 * changluyi
 * chestack
 * cmdy
 * coldzerofear
 * cui fliter
 * dependabot[bot]
 * dnugmanov
 * dolibali
 * fanriming
 * feixiang43
 * fsl
 * github-actions[bot]
 * gongysh2004
 * guangwu
 * gugu
 * guofen
 * hackeren
 * halfcrazy
 * huangjunwei
 * huangsq
 * hzma
 * isabelleatkins
 * jeffy
 * jimyag
 * jizhixiang
 * jyjiangkai
 * lanyujie
 * lfpython
 * lhalbert
 * liqd
 * liuzhen21
 * liyh-yusur
 * long.wang
 * luoyunhe
 * lut777
 * mengyu0829
 * mingo
 * narutoqq
 * netdever
 * pengbinbin
 * pengbinbin1
 * pengpeng wu
 * qiutingjun
 * renovate[bot]
 * runzhliu
 * shane
 * singchia
 * smartVan
 * vseeker
 * wang_yudong
 * wangchl01
 * wanglei
 * wangwangyusur
 * wangwangyusur288
 * wangyd1988
 * wenwenxiong
 * wujixin
 * xiaoyie
 * xieyanker
 * xuhao
 * xujunjie
 * xujunjie-cover
 * yi.luo
 * zbb88888
 * zcq98
 * zhang.zujian
 * zhangzujian
 * zhouhui-Corigine
 * 刘睿华
 * 夜微澜
 * 尚墨
 * 张祖建
 * 罗云鹤
 * 范日明
 * 袁又袁
 * 马洪贞

## v1.7.0 (2021-06-03)

 * [907b34d2b](https://github.com/kubeovn/kube-ovn/commit/907b34d2b7264436461c3b7d5fc02609e0a446d1) prepare for release v1.7.0
 * [ab727c98b](https://github.com/kubeovn/kube-ovn/commit/ab727c98be5dbe848dafc44025e766553508049e) diagnose: check sa related resource
 * [9bd2e9f87](https://github.com/kubeovn/kube-ovn/commit/9bd2e9f87dfe85062020c9ce957cff64e90242ec) fix: do not nat route traffic
 * [3bd14945b](https://github.com/kubeovn/kube-ovn/commit/3bd14945b4f2d998e939a1b3bd4d10b3b7535364) fix: release ip addresses even if pods not found
 * [f47941838](https://github.com/kubeovn/kube-ovn/commit/f47941838bd5acf9c1e87fb12cb01a0cff1ac688) fix typo
 * [2a2160d0c](https://github.com/kubeovn/kube-ovn/commit/2a2160d0c663b0d19bcd00347fe5239ea785ffa2) docs: add description of custom kubeconfig
 * [3dd99a791](https://github.com/kubeovn/kube-ovn/commit/3dd99a7910830c604d7bef30ff24415214b099d7) fix: add address_set to avoid error message
 * [ba40fd67e](https://github.com/kubeovn/kube-ovn/commit/ba40fd67e4a4c7bc6cf99c63549b97e8d65c5c0a) optimize Makefile
 * [cb95f4e6d](https://github.com/kubeovn/kube-ovn/commit/cb95f4e6db63a3476e7b11249f6ac25b5e387316) update vlan document
 * [31a96f219](https://github.com/kubeovn/kube-ovn/commit/31a96f219d162705345e5a76c278b79946aada58) add label to avoid deleting other
 * [6cd6b34b0](https://github.com/kubeovn/kube-ovn/commit/6cd6b34b0f1e836bc6cdf65532676dcce7e88390) delete unused log
 * [347340106](https://github.com/kubeovn/kube-ovn/commit/347340106c365edecfdca3756b1ea1461715cf05) add ovs internal-port for pod network interface
 * [9e715623d](https://github.com/kubeovn/kube-ovn/commit/9e715623deab30b7b72e21f03e001474fe77ed4a) support underlay mode with single nic
 * [d6c96d07d](https://github.com/kubeovn/kube-ovn/commit/d6c96d07dd56049ef708f8d156dadc57e3997e92) support underlay mode with single nic
 * [c1d3fc3cf](https://github.com/kubeovn/kube-ovn/commit/c1d3fc3cfe51d886e427923caa75b9fd9b4f4df7) fix: add node to pod allow acl
 * [ed49cd497](https://github.com/kubeovn/kube-ovn/commit/ed49cd4978c81982eb0dcabd3db4ef937b7af533) traffic rate for multus nic
 * [1b00190f3](https://github.com/kubeovn/kube-ovn/commit/1b00190f3025fa66ef0ba1c3a4b8145eeb9a98d8) add ovs internal-port for pod network interface
 * [775aec6c6](https://github.com/kubeovn/kube-ovn/commit/775aec6c6b24cfd9ab6fbcc46943556a1d282810) Add maintainers
 * [59847bc10](https://github.com/kubeovn/kube-ovn/commit/59847bc108cd37dc624da2c7d1332079263f6aea) add e2e tests for external egress gateway
 * [a0006ebf0](https://github.com/kubeovn/kube-ovn/commit/a0006ebf01b7a8425a72f8f1f51864585fde01e0) fix e2e testing on macOS
 * [0ff3d6bb7](https://github.com/kubeovn/kube-ovn/commit/0ff3d6bb78f6371790ca647cb2febc982210b9cb) ci: fix lint and scan error
 * [33e0ec27c](https://github.com/kubeovn/kube-ovn/commit/33e0ec27c302151b8c86fa47d17e13a0184df8d8) fix: check if provider network exists
 * [9e53d4cc8](https://github.com/kubeovn/kube-ovn/commit/9e53d4cc8b9382e9ef51471e5bf40250ebea31b0) update subnet document
 * [a2e4fec4f](https://github.com/kubeovn/kube-ovn/commit/a2e4fec4f943232b6fe6cac0d23a0e047dca11e3) rename ExternalGateway to ExternalEgressGateway
 * [1ccaec9af](https://github.com/kubeovn/kube-ovn/commit/1ccaec9af8fc86065276ecba48c7871c210ab5d5) fix installation doc
 * [34fb47594](https://github.com/kubeovn/kube-ovn/commit/34fb47594208c9a15b9edeb57668227b86feec48) fix: forward policy to accept
 * [bbbd091f6](https://github.com/kubeovn/kube-ovn/commit/bbbd091f6e0819405096ec5fd2f4ae77b5a21487) ci: fix lint error
 * [28cf4cc2d](https://github.com/kubeovn/kube-ovn/commit/28cf4cc2d2e69dda06b9ac2a1c61dab9db8925a8) traffic rate for multus nic
 * [0dcf69300](https://github.com/kubeovn/kube-ovn/commit/0dcf69300b864438685db904ce73679691d6d5d3) refactor: optimize service.go and subnet.go
 * [7719fc2ad](https://github.com/kubeovn/kube-ovn/commit/7719fc2ad42c542a8188c93cfb7c3cb260c8da2e) Check and Fetch all ValidatePodNetwork errors
 * [123ead482](https://github.com/kubeovn/kube-ovn/commit/123ead48272728ced58c721a86c19e34527d51d0) add judge about nic address
 * [17fe23027](https://github.com/kubeovn/kube-ovn/commit/17fe230273fdac9a979f2c3e523b5940f42e1286) implement new feature: external gateway
 * [01686e3e7](https://github.com/kubeovn/kube-ovn/commit/01686e3e77add612cdf1fc4d9fac93547c8c58e1) start_ic should run regardless of ts port
 * [c733c7e40](https://github.com/kubeovn/kube-ovn/commit/c733c7e40498f8eac88ca0d4f19b24d65a4c555e) add judge before use index
 * [ba709afb5](https://github.com/kubeovn/kube-ovn/commit/ba709afb5dbb11dc1041fcdf9f0c4e0bcb78ccee) specify ovs ops on diff nodes
 * [07089205d](https://github.com/kubeovn/kube-ovn/commit/07089205dc0f1031890f98886cbfb20fd6360050) fix mss rule
 * [4458a4d72](https://github.com/kubeovn/kube-ovn/commit/4458a4d7254010fe36028c46adf3d24c0eaf3e4b) Get node info from listerv1.NodeLister(index)
 * [19a7aed98](https://github.com/kubeovn/kube-ovn/commit/19a7aed98a20bdff6ba6cfa4f04bc541bd7dabc3) Clean up the wrong log
 * [27fe348a7](https://github.com/kubeovn/kube-ovn/commit/27fe348a7a7aa0a8e9807ff551754432ad911952) refactor: optimize subnet.go
 * [ddfd06b25](https://github.com/kubeovn/kube-ovn/commit/ddfd06b25e2bc35ca73814c690292f1d807ce521) Optimise the redundancy code
 * [bd55c104f](https://github.com/kubeovn/kube-ovn/commit/bd55c104fa532c030a0fad6512e02b2450b7a7f0) Handler the parse config error before used
 * [bd3f13dcc](https://github.com/kubeovn/kube-ovn/commit/bd3f13dccd3360efe3c6c9173d5d992b246f0e76) ci: remove 3-master e2e
 * [9e827e7b5](https://github.com/kubeovn/kube-ovn/commit/9e827e7b52f5cc2dcb08d9e2210d8ad2c9ced27f) Remove the unnecessary rm command
 * [587bbcdbb](https://github.com/kubeovn/kube-ovn/commit/587bbcdbb366e9b7dfae70e314acf11c9b29212d) Use localtime when the kube-ovn installed
 * [a52a38d06](https://github.com/kubeovn/kube-ovn/commit/a52a38d064e1ad83ae98b9255320b8750b02a95a) Fix the different time from container and host
 * [436e788be](https://github.com/kubeovn/kube-ovn/commit/436e788be597d039712fb32333d5b480d7a6da7e) add issue template
 * [5fc3cfb18](https://github.com/kubeovn/kube-ovn/commit/5fc3cfb184568cfc1b20c2b4a881881f0fe98224) add bgp doc
 * [f16fcb9ad](https://github.com/kubeovn/kube-ovn/commit/f16fcb9ada97e8060090b786b478c783e210ecaa) support afisafis
 * [d94af379b](https://github.com/kubeovn/kube-ovn/commit/d94af379bbc322c91c2b14a61e5dcf002f4565f9) feat: support graceful restart
 * [26a027252](https://github.com/kubeovn/kube-ovn/commit/26a0272527e5ad79652607ca20b2bf0bd7f6f486) fix: del might panic if duplicate delete
 * [41226d86f](https://github.com/kubeovn/kube-ovn/commit/41226d86f80d41d3fcacfc22d3e5b3bf9c9c920f) fix: lr-route for eip using nic-ip, and not external gateway addr.
 * [d176dac7c](https://github.com/kubeovn/kube-ovn/commit/d176dac7cc42aaa25ec5c1a8f43f595d24920f61) feat: support announce service ip
 * [136571d16](https://github.com/kubeovn/kube-ovn/commit/136571d167152b0d870e3b84d84da87de9734f1a) Fix some minor nits for docs
 * [2781a47b4](https://github.com/kubeovn/kube-ovn/commit/2781a47b46895e0c2e3d94514863a87f7ad564e0) add bpg options in bgp.md
 * [1b788902c](https://github.com/kubeovn/kube-ovn/commit/1b788902c388d7916dc711e44cf09cb1c09867fe) add Opstk&K8s ic doc
 * [cc8438161](https://github.com/kubeovn/kube-ovn/commit/cc84381617906003edfabe7ab9f92b91db3c4c57) add holdtime function
 * [b9e963393](https://github.com/kubeovn/kube-ovn/commit/b9e96339328c7340e297b34d8bf45dac38a45425) fix: do not re-generate ts port
 * [610f132bd](https://github.com/kubeovn/kube-ovn/commit/610f132bdbdb81be7b1bc1efc86ccbe132102200) fix: ignore root path doc ci
 * [bd1e0975f](https://github.com/kubeovn/kube-ovn/commit/bd1e0975f7dea75566ae7328b1fd3faef514e176) fix: do not gc learned routes
 * [be2048be3](https://github.com/kubeovn/kube-ovn/commit/be2048be3e636e7bb0cfc2105254ef3f23a68482) feat: add vxlan in README.md
 * [cbb2ddd47](https://github.com/kubeovn/kube-ovn/commit/cbb2ddd476f2e173f4f3e6dd440cc3da2d58dcb3) fix: get_leader_ip always return fist node ip
 * [03f597ce4](https://github.com/kubeovn/kube-ovn/commit/03f597ce47ed434ce5b1a28f5a4d67c464f68711) fix: remove tty error notification
 * [cc353bbca](https://github.com/kubeovn/kube-ovn/commit/cc353bbca9d547ac30e3ca39e816ed8838c999fd) fix ovn nb reconnect
 * [af2709dfa](https://github.com/kubeovn/kube-ovn/commit/af2709dfa32d79887c8573572a21475684afa2c1) add docs for 'multus ovn network'
 * [ffc20a915](https://github.com/kubeovn/kube-ovn/commit/ffc20a9158cb92cce59e105a0eff294ea940fd48) add vpc nat gateway docs
 * [a1ae937ac](https://github.com/kubeovn/kube-ovn/commit/a1ae937ac580bf9c529be99a56d8f45a8859c136) fix: static route for default multus network
 * [0489a72ad](https://github.com/kubeovn/kube-ovn/commit/0489a72adf650b8bf91fe1cf8df279e26ca1bafc) feat: support vxlan tunnel
 * [77f65449c](https://github.com/kubeovn/kube-ovn/commit/77f65449c1e5a1a5c90a256c10829c2a51ba53d4) append delete ovn-monitor in ovn.yaml
 * [c5ee49e82](https://github.com/kubeovn/kube-ovn/commit/c5ee49e8225f17f83393cd67ca666a33c49b6d32) split ovn-central and ovn-monitor
 * [e0890f727](https://github.com/kubeovn/kube-ovn/commit/e0890f727327afd792aed64521aa5c78a384ef41) Fix mount the systemid path
 * [fc92fbc2e](https://github.com/kubeovn/kube-ovn/commit/fc92fbc2eadfd2f83f16aac04ebf32dfe2d0fac7) handle update deployment vpc-nat-gw
 * [686681ef3](https://github.com/kubeovn/kube-ovn/commit/686681ef3603519588158a9b5e80f49dbac291b6) refactor: remove function genNatGwDeployment's return error
 * [064c38510](https://github.com/kubeovn/kube-ovn/commit/064c38510219e8fb81c1f0e65352ff15e35e6603) Update crd vpc-nat-gateways.kubeovn.io for pre-1.16
 * [a0dfea1b0](https://github.com/kubeovn/kube-ovn/commit/a0dfea1b0f32d9e5931f569409f0d0a2de079577) fix incorrect method for gateway node judgement
 * [86c99c378](https://github.com/kubeovn/kube-ovn/commit/86c99c378bb346cefc4a47005ec16c79150800a8) Fix the 'multus how to use' link
 * [1acb49921](https://github.com/kubeovn/kube-ovn/commit/1acb499213170a1a9d3f3d7c08a5374758455912) fix multi nic
 * [9c5ca0a04](https://github.com/kubeovn/kube-ovn/commit/9c5ca0a04b08b68d3d68c5dadc82e970b2fa1d10) fix duplicate imports
 * [b4750853d](https://github.com/kubeovn/kube-ovn/commit/b4750853db9ae994dc7a7599fa6e506885416984) fix: compatible with JSON format
 * [2a2cd27a6](https://github.com/kubeovn/kube-ovn/commit/2a2cd27a6a4f41125c4800a85ad87181f5ccb65d) fix: leader may change during startup, use cluster connection to set options.
 * [aad815484](https://github.com/kubeovn/kube-ovn/commit/aad8154848f71dab233f16eec3c160b2b2e8fa5c) fix SNAT on pod startup
 * [388119a75](https://github.com/kubeovn/kube-ovn/commit/388119a75ffe65896203136e5f4b6de22e2d3d23) fix development guide
 * [2efdac9a4](https://github.com/kubeovn/kube-ovn/commit/2efdac9a4b3a74b782f4f0531525ed311a0a6a33) fix gofmt
 * [c264bec18](https://github.com/kubeovn/kube-ovn/commit/c264bec18408064ddbfc7322af50c8f2cc688ea5) fix: configure nic failed when ifname empty
 * [763f8bcf7](https://github.com/kubeovn/kube-ovn/commit/763f8bcf79b1ce885b72d2b99bd9850fcdf941de) fix: port does not support vlan tag update
 * [a60764ea9](https://github.com/kubeovn/kube-ovn/commit/a60764ea98e695d91826057e143840b38d0c7663) fix build dev image
 * [faa7bc6a5](https://github.com/kubeovn/kube-ovn/commit/faa7bc6a5bc4e7d4bb2924bc4a4bb0c421437394) support hybrid mode for geneve and vlan
 * [d8472ba77](https://github.com/kubeovn/kube-ovn/commit/d8472ba77be74afae2bf1294d7f2f39265504cf3) remove extra space
 * [f9c836b6d](https://github.com/kubeovn/kube-ovn/commit/f9c836b6df232008743f7a2f5dc58a1bc54b8087) fix: compatible with no norhtd svc
 * [bbed09d36](https://github.com/kubeovn/kube-ovn/commit/bbed09d3634002299a2da36789aeb4ee5e329ef8) fix chassis check for node
 * [dfdf5f8b0](https://github.com/kubeovn/kube-ovn/commit/dfdf5f8b0fefc78d51e0b80b8b16f586924690a8) optimization for ovn/ovs status metric
 * [9e82ca3d5](https://github.com/kubeovn/kube-ovn/commit/9e82ca3d594ce169dc575f0eddf8f09c1918c162) fix: release norhtd lock when power off
 * [1fbfad523](https://github.com/kubeovn/kube-ovn/commit/1fbfad523fd9e95e9470181c76a5315ff291e4f9) add single node e2e
 * [f9ae6258e](https://github.com/kubeovn/kube-ovn/commit/f9ae6258e21e87a632bcb19f4acc6cc847612598) fix get pod attachment net
 * [0632e253a](https://github.com/kubeovn/kube-ovn/commit/0632e253a182793cd68c9bac841f56c83e0bb51d) support ovn defautl attach net
 * [2c1a8aa64](https://github.com/kubeovn/kube-ovn/commit/2c1a8aa6425e43055d5c3271ac96083a36a1da2f) add network-attachment-definitions clusterRole
 * [808a3a93b](https://github.com/kubeovn/kube-ovn/commit/808a3a93b3b4572f207480a5c94b9781939b0689) feat: multus ovn nic
 * [28e14188a](https://github.com/kubeovn/kube-ovn/commit/28e14188ae40ced07979f4e99fa231063ff9769d) update node ip when upgrade to dualstack
 * [0265747dc](https://github.com/kubeovn/kube-ovn/commit/0265747dc8bbb578eea25bc9f1963bcf2ebb6f5e) add details for prerequisite
 * [3e42f6842](https://github.com/kubeovn/kube-ovn/commit/3e42f68421b5e73a204fd08f0cd04d07c1ee11ab) Add Ecmp Static Route for centralized gateway
 * [b72e9d50b](https://github.com/kubeovn/kube-ovn/commit/b72e9d50bcfa8547233137a6a94c828172ad0ace) fix: disable offload if geneve port exists
 * [f4e665b95](https://github.com/kubeovn/kube-ovn/commit/f4e665b95f26927ec2aceecf4a709f27de4c03a9) disable offload for genev_sys_6081
 * [acade01b4](https://github.com/kubeovn/kube-ovn/commit/acade01b4c41b05332e0c6b3d923fdad1428fa48) refactor: optimize ovn command when error exists
 * [5251c2722](https://github.com/kubeovn/kube-ovn/commit/5251c2722310f89573337151d166a83d2e5232a1) add net-attach-def ClusterRole
 * [5126aedd1](https://github.com/kubeovn/kube-ovn/commit/5126aedd1c75680354bb7e44ea02fa01dcd6be53) add lsp with external_id
 * [ec7f7425b](https://github.com/kubeovn/kube-ovn/commit/ec7f7425b30d7c6cb27ebecb2586cb0d23917671) feat: multus ovn nic
 * [19e23d141](https://github.com/kubeovn/kube-ovn/commit/19e23d1419e04b6efc872fa29238e8e362a74015) fix: check ovn0 status
 * [c02afc00c](https://github.com/kubeovn/kube-ovn/commit/c02afc00c9c7835fcc267ff44de6e393f136e3c5) livenessprobe fail if ovn nb/ovn sb not running
 * [983831e0d](https://github.com/kubeovn/kube-ovn/commit/983831e0dfeb2c252bcc6d314d0f4d9df597ce9b) fix: disable checksum offload for ovn0 to prevent kernel issue
 * [d9f166b71](https://github.com/kubeovn/kube-ovn/commit/d9f166b7132a63619b305af7eed416915911af73) ignore ip6tabels check for v4 hostIP
 * [680802d61](https://github.com/kubeovn/kube-ovn/commit/680802d618956c48f6815820f2edb50c15b0644e) improve the code style of [import group ordering]
 * [8e38a79da](https://github.com/kubeovn/kube-ovn/commit/8e38a79da6c1c97bedbd5ff0103dae7bfc0ee60e) fix wrong sequence
 * [1e0d77c35](https://github.com/kubeovn/kube-ovn/commit/1e0d77c35f313d4d775963cf8dae2670349e1c92) update arm64 build
 * [638a03ac4](https://github.com/kubeovn/kube-ovn/commit/638a03ac4ef1edd42d57825901e9f1035f105e36) fix: restart ovn-controller to force update flows
 * [14784fbb9](https://github.com/kubeovn/kube-ovn/commit/14784fbb9c12b2c079319410ecb88cc7296ce209) fix: disable checksum validation
 * [a04dcfb62](https://github.com/kubeovn/kube-ovn/commit/a04dcfb626c6ad85083f40703fab0ea02a7735a8) Use public network effective image
 * [24095d7fe](https://github.com/kubeovn/kube-ovn/commit/24095d7fe927376d6bc315cf629303da3594dfc9) update usingips check when update finalizer for subnet
 * [54ef1af2a](https://github.com/kubeovn/kube-ovn/commit/54ef1af2afbb85584c16a1981d158372df2fc8fa) fix dependency
 * [717688d66](https://github.com/kubeovn/kube-ovn/commit/717688d66a20e0f13bd77dbd230eb7b7db4c6b61) Update vendor.
 * [496fc4ddd](https://github.com/kubeovn/kube-ovn/commit/496fc4dddbf1b863555029905a0f60fe3ed6afb5) trim space the port_binding's output
 * [00fdac83c](https://github.com/kubeovn/kube-ovn/commit/00fdac83cdfebeacccba8a8bb6ff21fe5dbaba58) refactor: remove unnecessary config logic
 * [b06dad21b](https://github.com/kubeovn/kube-ovn/commit/b06dad21b8c89c68d8d571293a03e8f712e53af2) update maintainers
 * [e5d9584e3](https://github.com/kubeovn/kube-ovn/commit/e5d9584e3c5c34f3a167877a8ffdf9698ab12785) docs: deprecated webhook
 * [92cc4ed37](https://github.com/kubeovn/kube-ovn/commit/92cc4ed376b62cdbeb1634bc5b4b2e5c60820be2) fix: add missing ovn-ic binary
 * [c0349e4fc](https://github.com/kubeovn/kube-ovn/commit/c0349e4fcb4a5afc278f2a1e72443f98da715e3f) chore: change action name
 * [1a448eccc](https://github.com/kubeovn/kube-ovn/commit/1a448eccce0d65f7cf48e0895f4a012e74ef70ff) chore update artworks
 * [537588c39](https://github.com/kubeovn/kube-ovn/commit/537588c397caec5b1ee5a32eb0150a436e58c9d6) fix: delete chassis_private when delete node
 * [a50fb1817](https://github.com/kubeovn/kube-ovn/commit/a50fb1817b7a3b4b8cbcf86e35b11583af64c91a) Add 'kubectl ko trace' command's default namespace
 * [fad9473d9](https://github.com/kubeovn/kube-ovn/commit/fad9473d90f360cf797955fa78b90aa3ab57e077) Add 'kubectl ko trace' command's default namespace
 * [77c92ca8c](https://github.com/kubeovn/kube-ovn/commit/77c92ca8c5ec5dd78268333a0d141687dceb470c) perf: reclaim heap memory after compaction.
 * [f3df58aea](https://github.com/kubeovn/kube-ovn/commit/f3df58aeafcb85d7107fa0a4974ac3d71f5f7466) remove the old script
 * [b69f389cf](https://github.com/kubeovn/kube-ovn/commit/b69f389cf6ea337da30fcc012eae208196f85abd) docs: add CNCF description
 * [08b95e747](https://github.com/kubeovn/kube-ovn/commit/08b95e747bd96d054a2ffd4371a0f0e7328de062) fix: gc not exist node error
 * [9f6614613](https://github.com/kubeovn/kube-ovn/commit/9f6614613cedfdbb30c268be96b3c68fc0755bf8) perf: use new option to decrease ovn-sb size
 * [9dc069083](https://github.com/kubeovn/kube-ovn/commit/9dc0690831cdc75d143b2a2eab20f3d6d739507c) fix: return err
 * [8bd446083](https://github.com/kubeovn/kube-ovn/commit/8bd4460832a207a62074cd60e8b8f731fa666b17) docs: add faq section
 * [482e6f71a](https://github.com/kubeovn/kube-ovn/commit/482e6f71a078396f51cd29fd223d3cb017e1a65c) add vpc nat gateway Dockerfile
 * [b0e983f04](https://github.com/kubeovn/kube-ovn/commit/b0e983f04521ffa9bf5fcac394331a3e11ab9774) feat: vpc nat gateway
 * [951e31ea3](https://github.com/kubeovn/kube-ovn/commit/951e31ea3135c72c0472cdcdb86d8a8864cf8b1a) add node address allocate check when init
 * [215c8f45b](https://github.com/kubeovn/kube-ovn/commit/215c8f45bbc3f77ba5928f50c897eea53aa8c2a6) update upgrade for ovn-default and join subnet
 * [a537985d9](https://github.com/kubeovn/kube-ovn/commit/a537985d9d5d4dea6208c1104f1800bf8dd281a8) fix: lint error
 * [d0d3e89cc](https://github.com/kubeovn/kube-ovn/commit/d0d3e89cc82159383c3e22e0a1fb62f815134c65) fix: add missing ovn-ic-db schema
 * [98651014c](https://github.com/kubeovn/kube-ovn/commit/98651014c8994ed08b4de0f9afaf32b9fe8c4384) update subnet ip num calculate
 * [d6bb03bd5](https://github.com/kubeovn/kube-ovn/commit/d6bb03bd527f03aa26a8b87d5a809cd1a80a603d) fix: masq traffic to ovn0 from other nodes
 * [0a7024f97](https://github.com/kubeovn/kube-ovn/commit/0a7024f97bba8511a84185474d907a236f71a93a) refactor: reduce duplicated GetNodeInternalIP function
 * [ac294669c](https://github.com/kubeovn/kube-ovn/commit/ac294669c7c7e2a24d2af6910494b6b8d175091b) chore: update go version
 * [0e9c717d0](https://github.com/kubeovn/kube-ovn/commit/0e9c717d05ebc26a01c2c6c3961a17316c675260) chore: move build dependency from alauda to kubeovn
 * [64fac57a9](https://github.com/kubeovn/kube-ovn/commit/64fac57a9c8626c916ff6ffdb103bd5f765ed729) feat: support set default gateway in install script
 * [ca71de3cf](https://github.com/kubeovn/kube-ovn/commit/ca71de3cfb8de87413de482c155df9927316bb4b) docs: fix typos
 * [582cb9ce6](https://github.com/kubeovn/kube-ovn/commit/582cb9ce6acb323fbe8e8dd19651901a24249b7f) Update install-pre-1.16.sh
 * [62fc20efe](https://github.com/kubeovn/kube-ovn/commit/62fc20efea653fc7463840c746ec7e70013fcf2f) Update install.sh
 * [87859ac19](https://github.com/kubeovn/kube-ovn/commit/87859ac1951acbe9c4331cdeb80e7eb10948c3e8) go import repo change to kubeovn
 * [1152744e9](https://github.com/kubeovn/kube-ovn/commit/1152744e9f47626310601a4c63a7af61e3fe6bf6) feat: vpc nat gateway
 * [298138e41](https://github.com/kubeovn/kube-ovn/commit/298138e41d5d36f0400633bc3183f191de0c74cf) Resolving typo.
 * [4701fcb37](https://github.com/kubeovn/kube-ovn/commit/4701fcb3741213c4fbcd6e5e1a151fdcd4c6d3c2) filter repeat exclude ips
 * [e3931f0ea](https://github.com/kubeovn/kube-ovn/commit/e3931f0ea926e91c36fb2ac2699cb4e080f87b58) modify ip count for dual
 * [a4ddb3604](https://github.com/kubeovn/kube-ovn/commit/a4ddb3604584dd5b558686e663face1cb2a15584) docs: add ARCHITECTURE.MD
 * [9eee6f938](https://github.com/kubeovn/kube-ovn/commit/9eee6f938da6d79aac59171b6f53a108d123bede) refactor: reduce duplicated function
 * [a7b687a05](https://github.com/kubeovn/kube-ovn/commit/a7b687a05fadf8acf7b25d4dd994f96a7aa67a5a) fix: add dpdk pod name
 * [d32b423b6](https://github.com/kubeovn/kube-ovn/commit/d32b423b68ffdb8d49394f65c99c9342fd73d699) Update cleanup.sh
 * [9faaff574](https://github.com/kubeovn/kube-ovn/commit/9faaff574f68bf0d763c3a29cfe4b31f26c0b659) Update cleanup.sh
 * [df065f94b](https://github.com/kubeovn/kube-ovn/commit/df065f94b4a835e8f6a1dc17d7e73570f4ec80ee) test: add service e2e
 * [60e49f5a0](https://github.com/kubeovn/kube-ovn/commit/60e49f5a0ff999a98e9edc733ad6b548bed224ab) modify test problem
 * [2dbcb76fa](https://github.com/kubeovn/kube-ovn/commit/2dbcb76fa31cce633a63f18d1ff0a2a592e16bf9) fix: kube-proxy check
 * [512044cb3](https://github.com/kubeovn/kube-ovn/commit/512044cb389df7257adca989f604e3a00da0ddd4) ovn-central: set default db addr same with leader node to fix nb and sb error 'bind: Address already in use'
 * [c755ef232](https://github.com/kubeovn/kube-ovn/commit/c755ef232313476a24111621cda522119aec7749) fix: reset ovn0 addr
 * [a168c2826](https://github.com/kubeovn/kube-ovn/commit/a168c28264e12dfcbec823c6756c76cb1c22f974) tests: add e2e for ofctl/dpctl/appctl
 * [f6dc58a5d](https://github.com/kubeovn/kube-ovn/commit/f6dc58a5dcabd0731db9abe1a5f92af7c396892e) ci: replace image
 * [b1d03370e](https://github.com/kubeovn/kube-ovn/commit/b1d03370e61d81a0d72a4f0f2b4917fbdcffbc59) docs: clarify dpdk usage scenario
 * [21d9940b9](https://github.com/kubeovn/kube-ovn/commit/21d9940b9d13257699291e516166ee02df1a45d2) ci: update kind version and set timeout
 * [8b833ee52](https://github.com/kubeovn/kube-ovn/commit/8b833ee52ec0139dcd75b55556b8fe9ce51703a0) Update install-pre-1.16.sh
 * [4b6f0eedb](https://github.com/kubeovn/kube-ovn/commit/4b6f0eedbed57a64c8b025f06a220f0a9670161a) Update install.sh
 * [f6f88501f](https://github.com/kubeovn/kube-ovn/commit/f6f88501f9876f9b95c0549abcc22a43733e18e6) refactor: remove duplicated call
 * [473cdc487](https://github.com/kubeovn/kube-ovn/commit/473cdc48766dea1a152b2c430069014cc4064804) Update kubectl-ko
 * [1ca17686b](https://github.com/kubeovn/kube-ovn/commit/1ca17686b9493f2e64b72d7b5e3d85f7e2df08d0) Fix missing square brackets in curl ipv6
 * [136336b21](https://github.com/kubeovn/kube-ovn/commit/136336b2124593ebedd8e5663f164b7ecf9b5981) Modify the health check for kube-proxy port, compatible with ipv6
 * [98a56dece](https://github.com/kubeovn/kube-ovn/commit/98a56dece7c20235cd806ed2cdef50cebb07de24) Update controller.go
 * [c52c067b9](https://github.com/kubeovn/kube-ovn/commit/c52c067b946b4aac7eef8cbb3002eb4aeefc9045) Fix: remove IsNotFound when get configmap external gateway
 * [74fa7729d](https://github.com/kubeovn/kube-ovn/commit/74fa7729d806b4129edb8604ef2effa0b79c13ca) Fix: check kube-proxy's 10256 port healthz
 * [d594554d7](https://github.com/kubeovn/kube-ovn/commit/d594554d73f83d7d5b9e601533c8fb27f225dc47) fix: ip6tables check error
 * [b17f23732](https://github.com/kubeovn/kube-ovn/commit/b17f23732c65807e84a7d520c5c64c6869e3ee8e) Add MAINTAINERS file
 * [2783c134f](https://github.com/kubeovn/kube-ovn/commit/2783c134f16a9ea9c172afd30206868025d677ec) add vpcs && vpcs/status clusterRole
 * [31e1226e0](https://github.com/kubeovn/kube-ovn/commit/31e1226e01080cdd4ce27ef7979a18ba2736439b) Update install-pre-1.16.sh
 * [f1efaa7f2](https://github.com/kubeovn/kube-ovn/commit/f1efaa7f256f0bb936a65b5a1a5b8b38d12d84ff) delete connect to ovsdb for ovn-monitor
 * [f69ae44be](https://github.com/kubeovn/kube-ovn/commit/f69ae44bec65b67e2ab61689de0aa27da7143b41) cni-bin-dir,cni-conf-dir configurable Fix https://github.com/alauda/kube-ovn/issues/655
 * [f5999b3b5](https://github.com/kubeovn/kube-ovn/commit/f5999b3b580f129c25c8882f136d6f35d25875c7) Update install.sh
 * [e13448aac](https://github.com/kubeovn/kube-ovn/commit/e13448aac6c524c3ee1bfa42f27b2569418e6a2c) Error: unknown command "ko" for "kubectl"
 * [7d56483a3](https://github.com/kubeovn/kube-ovn/commit/7d56483a3e76983a96eff5d894b3ad92df186c37) Fix: wrong split in FindLoadbalancer function
 * [34776b8a7](https://github.com/kubeovn/kube-ovn/commit/34776b8a7986225cae89be0b30773201578a84bb) vlan nic support regex
 * [f23093c4f](https://github.com/kubeovn/kube-ovn/commit/f23093c4f94194be1416d0aea760d0861c77b9fe) fix underlay gateway flood logs
 * [4a9901aa6](https://github.com/kubeovn/kube-ovn/commit/4a9901aa68f95ff8207cb45fb271c4d83095d563) fix: check required module before start
 * [8d4694f83](https://github.com/kubeovn/kube-ovn/commit/8d4694f835e724cac2cb23cfedbf47a1646d2044) docs: add underlay docs
 * [3713b2537](https://github.com/kubeovn/kube-ovn/commit/3713b2537f6bea24f4a1d87f64016265d5f5727e) chore: update ovn to 20.12 and ovs to 2.15
 * [1ab871307](https://github.com/kubeovn/kube-ovn/commit/1ab871307bf2c98b692bfff3f7d3e1db314dab34) prepare for next release
 * [a94803d32](https://github.com/kubeovn/kube-ovn/commit/a94803d3208faec094cafb7dfa80c6cb6a286218) fix: make sure northd leader change
 * [03487cf28](https://github.com/kubeovn/kube-ovn/commit/03487cf28ac435c1723b858f03471464a689504d) fix: make sure ovn-central is updated one by one
 * [9d3b78a31](https://github.com/kubeovn/kube-ovn/commit/9d3b78a319fcd623c8df02df1f87a1b7926e205e) fix: restart when init ping failed
 * [6e09c77de](https://github.com/kubeovn/kube-ovn/commit/6e09c77deeba8ff2643f3309b19ee0d109c6f838) fix: increase raft timer to avoid leader flap
 * [87aa15cbf](https://github.com/kubeovn/kube-ovn/commit/87aa15cbf9494d68d66d0e1494e7362c7f39aa76) pass golangci-lint
 * [134ea89d0](https://github.com/kubeovn/kube-ovn/commit/134ea89d0e196584547171d11c78caecda059cf7) add golangci-lint to github actions
 * [d325e7e0b](https://github.com/kubeovn/kube-ovn/commit/d325e7e0b606ca141b05a169ce3ea1f993782006) fix pod terminating not recycle ip when controller not ready
 * [87af4ca9f](https://github.com/kubeovn/kube-ovn/commit/87af4ca9ff374fc8ae1bcbe8af84e73532f2e7c3) fix: add new iptable cleanup commands
 * [d287063bf](https://github.com/kubeovn/kube-ovn/commit/d287063bf7e9a69fa740a6ea35f007e5878e78b8) modify static gw changed problem
 * [fcf3be190](https://github.com/kubeovn/kube-ovn/commit/fcf3be190e141fcc6099d34bcefd4313ae71ae81) Fix wait pod network ready take long time
 * [0b4e4458e](https://github.com/kubeovn/kube-ovn/commit/0b4e4458e51033f096eb229624a36fd18b39ce13) fix: when address is empty, skip route/nat deletion
 * [ed0e9ba22](https://github.com/kubeovn/kube-ovn/commit/ed0e9ba22a7e65fc1bc0fc50c220a9957a59ed44) fix: update ipam cidr when subnet changed
 * [06816efb5](https://github.com/kubeovn/kube-ovn/commit/06816efb55507304c44a33b90cf1e140533bc47a) modify test problem for dual-stack upgrade
 * [d47ccb678](https://github.com/kubeovn/kube-ovn/commit/d47ccb678692e441a774d11477269a4c4e430544) release: 1.6.0
 * [b8f221bf7](https://github.com/kubeovn/kube-ovn/commit/b8f221bf7d47b2190acfd716878e1b5aa441a409) docs: add docs for vpc
 * [12cf140b1](https://github.com/kubeovn/kube-ovn/commit/12cf140b167755bcb7a29981f5962ff369689694) fix typo
 * [b13cb7bf8](https://github.com/kubeovn/kube-ovn/commit/b13cb7bf8f34516a9fe9cf64eb0d56b14644c7d1) ci: update go version to 1.15
 * [7f9eefedb](https://github.com/kubeovn/kube-ovn/commit/7f9eefedb1267d9d18059c638b23361d1c198891) Fix: replace the command to run the script via 'sh' with 'bash'
 * [076ab28f8](https://github.com/kubeovn/kube-ovn/commit/076ab28f80c46826c5237da406eb18eb38d4bb54) Fix the default mtu parameter's describe
 * [8e6086678](https://github.com/kubeovn/kube-ovn/commit/8e6086678b923eb032a8b95d13a9bf214b1f38e8) modify network policy process
 * [171dcff6d](https://github.com/kubeovn/kube-ovn/commit/171dcff6dd3b27899e06ae752f8fa34896b159de) upgrade for subnet from single protocol to dual-stack
 * [bbc68577b](https://github.com/kubeovn/kube-ovn/commit/bbc68577b91ab26cdec5208c02dc165fe73a8222) add network policy adapt for dual-stack
 * [c01766cfa](https://github.com/kubeovn/kube-ovn/commit/c01766cfaef9554b6acbb435d81050951c97a1de) feat: update ovn to 20.09
 * [315831aa0](https://github.com/kubeovn/kube-ovn/commit/315831aa0f5baacac396266d009f746565b79db0) docs: prepare docs for 1.6.0 release
 * [a1e7974fa](https://github.com/kubeovn/kube-ovn/commit/a1e7974fa950420d5e2520942335f6284f161bdc) perf: add pprof to pinger
 * [627956e95](https://github.com/kubeovn/kube-ovn/commit/627956e95e2c440b4b75709dfdc4e33050209815) doc for dual-stack
 * [02751bf42](https://github.com/kubeovn/kube-ovn/commit/02751bf42d196e6e5542b5284953a554eb83e857) Update the container nic name use the CNI_IFNAME parameter which passed by kubelet
 * [14f36814f](https://github.com/kubeovn/kube-ovn/commit/14f36814f7899402f62ef85941ac066b1f2312dc) ci: enable docker experimental feature
 * [9a785fc9e](https://github.com/kubeovn/kube-ovn/commit/9a785fc9eeb744ddcc3d8eb4d98cd419ca26910b) ci: build multi arch image
 * [03ff96e66](https://github.com/kubeovn/kube-ovn/commit/03ff96e66b8316348fa50ac7371cc27616464caa) <fix>(np) fix mulit np rule and gateway bug
 * [20f3fcb17](https://github.com/kubeovn/kube-ovn/commit/20f3fcb178c0bf300ad6b792de53cc9aab9218fd) fix start-db.sh echo message
 * [52b39d764](https://github.com/kubeovn/kube-ovn/commit/52b39d764dacece2dec406db485e46181f3bd7d3) fix: iface check error
 * [072870b16](https://github.com/kubeovn/kube-ovn/commit/072870b16e6858d3690b10db975b2f28da7e7b7b) fix: add missing ping due to deb build
 * [efdd3913b](https://github.com/kubeovn/kube-ovn/commit/efdd3913b22a2bed4a98b8882897403869fc82aa) fix: find iface by full match first then regex match
 * [f922ef752](https://github.com/kubeovn/kube-ovn/commit/f922ef752883131d5d70abdaa8a2bb4b3235ef32) fix: livenessProb/readinessProb might conflict when run logrotate at same time
 * [f1fe2b2ea](https://github.com/kubeovn/kube-ovn/commit/f1fe2b2ead442a58331efe8fb1b47ac7bc4858f6) modify subnet and ip crds
 * [a2d76df7c](https://github.com/kubeovn/kube-ovn/commit/a2d76df7c0fea1b64289024a4cf637ef8657e2c2) modify service vip parse error
 * [8aa5d0a4f](https://github.com/kubeovn/kube-ovn/commit/8aa5d0a4f978a85a460fe2adcc05b4f6b1a39dd5) update vendor
 * [44381c74e](https://github.com/kubeovn/kube-ovn/commit/44381c74e23437b03fdd31ce6b30f0f5a2c29005) update client-go
 * [96c1c1003](https://github.com/kubeovn/kube-ovn/commit/96c1c1003535a87166bbf1cf6fbccc9a99a99cc6) fix: np with multiple rules
 * [87e6ded04](https://github.com/kubeovn/kube-ovn/commit/87e6ded0411e1a27bc18b0d807c5eacff028bb58) modify loop error for get metrics
 * [1e2a74774](https://github.com/kubeovn/kube-ovn/commit/1e2a747749b8a702cf59304f00d1549778ff0e34) diagnose: add more diagnose info
 * [aea12bae5](https://github.com/kubeovn/kube-ovn/commit/aea12bae5a744b1b6c688ed5343cfab729dcd802) ci: trigger action when yamls change
 * [7bd6bf39f](https://github.com/kubeovn/kube-ovn/commit/7bd6bf39f793698e1b429018fb7a1b10cb19e192) fix: ha e2e failed
 * [56774aafd](https://github.com/kubeovn/kube-ovn/commit/56774aafd2d98d4589b2e4fd18ea758b7e6cf66f) fix: allow traffic to gateway
 * [a78c2661b](https://github.com/kubeovn/kube-ovn/commit/a78c2661bdc559e6b00521ae7ee62399ac633f05) fix: cni-server default encap ip use right interface ip
 * [7d31e617a](https://github.com/kubeovn/kube-ovn/commit/7d31e617a42719e68e338d60a278774e51f5146b) feat: change default build image to ubuntu
 * [e2cd78715](https://github.com/kubeovn/kube-ovn/commit/e2cd7871583f812271089c5002743460246e6242) add build for dualstack
 * [ddda63320](https://github.com/kubeovn/kube-ovn/commit/ddda633204486877d8176476d1bd0470a84c3ecc) feat: distributed eip
 * [a6fef94a8](https://github.com/kubeovn/kube-ovn/commit/a6fef94a823be183eb80ca14be4703300e2c5add) Add CNI modify for dualstack
 * [a54bfc284](https://github.com/kubeovn/kube-ovn/commit/a54bfc2840f84d6fcad315a7ba0fd05ded6c7d12) Debian: Add debian docker image support
 * [8a01cb1c2](https://github.com/kubeovn/kube-ovn/commit/8a01cb1c211e3becefc46262c0fc78257be47b02) Add adaption for dualstack, part of daemon process.
 * [9738af184](https://github.com/kubeovn/kube-ovn/commit/9738af184b7604013ef166a2fdf6ec1043d624a9) chore: reduce binary size
 * [6483d6e33](https://github.com/kubeovn/kube-ovn/commit/6483d6e330cb4f46a9e260d0d65add82c757ab91) modify build problem
 * [dab50b33d](https://github.com/kubeovn/kube-ovn/commit/dab50b33d5c7ac8e51ccb470a257ba4bb3a332fd) Append ip monitor to document
 * [344288195](https://github.com/kubeovn/kube-ovn/commit/344288195ea1a74df7f0b87fd32fed517b59fc89) license: fix felix dir
 * [2ef665683](https://github.com/kubeovn/kube-ovn/commit/2ef6656832d1c091acd33d83472efef7e871c886) feat: support advertise subnet route
 * [ecbd01a65](https://github.com/kubeovn/kube-ovn/commit/ecbd01a65546239c2786e0fe00b611f0c17fbb01) Add IP Num Alert
 * [d64e69317](https://github.com/kubeovn/kube-ovn/commit/d64e693178ce90a3dcec71a2ecfefb40b411ec4a) Add adaption for dualstack, part of controller process.
 * [7246037b1](https://github.com/kubeovn/kube-ovn/commit/7246037b16227deb35d4e675d056de519d432228) convert ip to string
 * [2aecb3d9b](https://github.com/kubeovn/kube-ovn/commit/2aecb3d9b3890b7d16dad0d17e5cb7d2c80699dd) add pod static ip validate
 * [b58e01b61](https://github.com/kubeovn/kube-ovn/commit/b58e01b6185ccccb02ca722e2a6f25498cbe60a9) chore: add COC and roadmap
 * [7bbdc00f2](https://github.com/kubeovn/kube-ovn/commit/7bbdc00f29e9011731d0f72b900afb5d8f77b3eb) fix: move felix to self repo to remove bird license
 * [d2b570cfb](https://github.com/kubeovn/kube-ovn/commit/d2b570cfbdd6b62314c4b646538305decaedf60a) Add license scan report and status
 * [86584b95d](https://github.com/kubeovn/kube-ovn/commit/86584b95d7ea4498ac61b026b46fa004c459823c) fix: default network
 * [ccea68bfa](https://github.com/kubeovn/kube-ovn/commit/ccea68bfa249b8d382c21e7133af4e84c6288d5b) release for 1.5.2
 * [07347501a](https://github.com/kubeovn/kube-ovn/commit/07347501aa62d191fb08824bf76cdbe1c7590f58) fix: ovn-ic support ssl
 * [4d8b186a3](https://github.com/kubeovn/kube-ovn/commit/4d8b186a35da4fcd0dc6d0f0c5959d292311061b) fix: nat rules can be modified
 * [f535460ff](https://github.com/kubeovn/kube-ovn/commit/f535460ff720851bee6586840a786b9d1cd23d1a) fix: remove svc cidr routes
 * [e3082cd77](https://github.com/kubeovn/kube-ovn/commit/e3082cd7727ffd47aeae6af0ea47dabbd397a80e) ci: specify ubuntu version to make github action happy
 * [f6cce9a0a](https://github.com/kubeovn/kube-ovn/commit/f6cce9a0afabe91904cd2cd13a9733c877cdd9e3) fix: specify exec container to mute warning message
 * [2215c05f4](https://github.com/kubeovn/kube-ovn/commit/2215c05f45f9f42c0fffe57a342d73d01f5da103) feat: remove cluster ip dependency for ovn/ovs components
 * [a9747b316](https://github.com/kubeovn/kube-ovn/commit/a9747b3161f51844cf3a052c6c589c0d4f580e9d) fix: add resources limits to avoid eviction
 * [005711961](https://github.com/kubeovn/kube-ovn/commit/005711961cf3b5cc1a11c06558d1b0b9f21d69ba) fix: vpc static route manage
 * [8deb5d8d0](https://github.com/kubeovn/kube-ovn/commit/8deb5d8d05babbe20881f492a775fc86d08b7b8d) fix: validate vpc subnet
 * [256ac6c53](https://github.com/kubeovn/kube-ovn/commit/256ac6c5350c5e8a67657a06872928c86b491363) Fix external-address config description
 * [ccda611a5](https://github.com/kubeovn/kube-ovn/commit/ccda611a5c3d8cc196bfeaab1985350f0168b7d9) Fix the problem of confusion between old and new versions of crd
 * [f2f648011](https://github.com/kubeovn/kube-ovn/commit/f2f6480112272ce221bf0bd4da3010633a15a541) fix: ovn-central check if it exits in NODE_IPS
 * [5b973a89e](https://github.com/kubeovn/kube-ovn/commit/5b973a89e501dd2e1c27cc0f78ca7267322d2827) fix: check ipv6 requirement before start
 * [86941a8a8](https://github.com/kubeovn/kube-ovn/commit/86941a8a87885765ceaafa5d9a07cf100723004e) feat: add ovs/ovn log rotation
 * [ef41733c5](https://github.com/kubeovn/kube-ovn/commit/ef41733c558ffd68c7b7ec79c91f93afbef727de) add node ping total count metric
 * [5e6bd9112](https://github.com/kubeovn/kube-ovn/commit/5e6bd9112b69ea1928830084961d4b3286381ccd) diagnose: add ovs-vsctl show to diagnose results
 * [7301e9926](https://github.com/kubeovn/kube-ovn/commit/7301e992665c6b18f47384ddc1e0fb36c07e8274) fix: nat rules
 * [6026028a4](https://github.com/kubeovn/kube-ovn/commit/6026028a48bdffb9252af001e3df110479287fa2) fix: masq other nodes to local pod to avoid nodrport triangle traffic
 * [d41110ec1](https://github.com/kubeovn/kube-ovn/commit/d41110ec10d8887c44c707e51dd4b1a8c6823221) Update install.sh to allow dpdk limits configuration (#546)
 * [a128d7fcc](https://github.com/kubeovn/kube-ovn/commit/a128d7fcc4a752f7deffa8aeb3d6c26bbf0eb76f) format
 * [b6ad17b53](https://github.com/kubeovn/kube-ovn/commit/b6ad17b537389e64f306d4837538b4c0d0ef0d59) test: e2e uses IPVS cluster by default
 * [f6951cf5d](https://github.com/kubeovn/kube-ovn/commit/f6951cf5d98db035df3471364bc396682c87f703) chore: update go version to 1.15
 * [1f703c3d3](https://github.com/kubeovn/kube-ovn/commit/1f703c3d3f49f887de951b4fe0b057a44e7f1fe6) fix: tolerate all taints
 * [f8ace73c1](https://github.com/kubeovn/kube-ovn/commit/f8ace73c190930b21ce005c8a7ad9a3d4b0ace7d) feature: add vpc static route
 * [f62cb4ebf](https://github.com/kubeovn/kube-ovn/commit/f62cb4ebfe52ef3de8f52b2f5c84acaee74e705b) fix: cleanup script error
 * [3bac21f7e](https://github.com/kubeovn/kube-ovn/commit/3bac21f7ecaa674eaa386e8d55f010a20aaac101) docs: modify eip config description
 * [1f07d96b7](https://github.com/kubeovn/kube-ovn/commit/1f07d96b7081b0496fb339635385f09f069dc6de) security: remove sqlite to mute cve warning
 * [015bc6259](https://github.com/kubeovn/kube-ovn/commit/015bc6259b7cb9e32bca94fa3a781371cf6a6c0a) test: add e2e for kubectl-ko
 * [aa86e406f](https://github.com/kubeovn/kube-ovn/commit/aa86e406fe17f9cf9c76803514e7b4c373e3bb8b) feat: pinger can return exit code when failed
 * [2cf855ecd](https://github.com/kubeovn/kube-ovn/commit/2cf855ecd7b4550fc4a273aeac7c3eb9d6641aed) fix: nat traffic that from host to svc
 * [cbe0ad55f](https://github.com/kubeovn/kube-ovn/commit/cbe0ad55f41c47a4d8e5f963802e7051f4561ab6) docs: new feat for disable-ic, regex iface and pod bind subnet
 * [5dbaf2d3a](https://github.com/kubeovn/kube-ovn/commit/5dbaf2d3a528e0bf10febcb9877bea7d51bcd003) sync the default subnet of ns by vpc's status
 * [dd2234f4e](https://github.com/kubeovn/kube-ovn/commit/dd2234f4e65d8fd31f2cad84aa055d12b3ae46e9) fix: devault vpc lb/dns
 * [32c49c1b9](https://github.com/kubeovn/kube-ovn/commit/32c49c1b93842130e60f74dd477aad0aad9c5a30) fix: shutdown vpc workqueue
 * [67076d62f](https://github.com/kubeovn/kube-ovn/commit/67076d62f0e3f26af1a864678e660edadfaf5464) fix: subnet CIDRConflict
 * [d5b819b03](https://github.com/kubeovn/kube-ovn/commit/d5b819b03710f149ae4c8d66687e87ea09589827) fix: subnet bind to ns
 * [921190ef8](https://github.com/kubeovn/kube-ovn/commit/921190ef8534b0d591e9f89c6eb9ca4b07d2fbc8) feature: add vpc crd
 * [b5ecac95c](https://github.com/kubeovn/kube-ovn/commit/b5ecac95cbb27b39250f870dd6e6c885ca7dae79) Release and gc the resources in vpc
 * [15eca9dca](https://github.com/kubeovn/kube-ovn/commit/15eca9dca39aa497cc3670bcb453aff9d020acdc) fix: gc logic router
 * [91fec5631](https://github.com/kubeovn/kube-ovn/commit/91fec5631ca0ecbdd3d677282a25e8029169ecad) gc and clean vpc
 * [7a0e28b98](https://github.com/kubeovn/kube-ovn/commit/7a0e28b98bfcbf90ee6766bf89a72d85df51428c) Remove the VPC while removing the default subnet
 * [99217cec1](https://github.com/kubeovn/kube-ovn/commit/99217cec142f3e6943280e34a0be85926664ae7c) feature: support custom vpc
 * [9d821bce0](https://github.com/kubeovn/kube-ovn/commit/9d821bce095b59b17cee1d86f365fa5032d74fcf) chore: refactor log
 * [240cd800a](https://github.com/kubeovn/kube-ovn/commit/240cd800a7545356d723f28c09e3cfdee5d8fe87) feat: iface support regexp
 * [94b6b1b59](https://github.com/kubeovn/kube-ovn/commit/94b6b1b59b089023864a70bf5c55334926a0abdf) feat: support disable interconnection for specific subnet
 * [652190c35](https://github.com/kubeovn/kube-ovn/commit/652190c359fa7bf0f37db317dc4f4680e70c1fb5) modify review problems
 * [7285581a5](https://github.com/kubeovn/kube-ovn/commit/7285581a535b3e9207d1a8231fba9e7fa852d4cc) docs: v1.5.1 changelog
 * [47f0acbbd](https://github.com/kubeovn/kube-ovn/commit/47f0acbbd9ca240abb2946b6a5d39d5fa271c0bd) perf: accelerate ic and ex gw update
 * [bafac87ee](https://github.com/kubeovn/kube-ovn/commit/bafac87ee1c885dfbefeee26db8fc4d5364ef835) fix: missing version date
 * [8ef12007c](https://github.com/kubeovn/kube-ovn/commit/8ef12007cc3b77612ba132fef4c0864a8aa92ec6) fix: check multicast and loopback subnet
 * [3b20abb0e](https://github.com/kubeovn/kube-ovn/commit/3b20abb0e7770d6cf6722f144b4b307ab5caac82) monitor: refactor grafana dashboard
 * [f9cbaea5b](https://github.com/kubeovn/kube-ovn/commit/f9cbaea5b3ea9adc6744064fde8b5841fc51c0d0) docs: do not allow install to namespace other than kube-system
 * [559e2cd8c](https://github.com/kubeovn/kube-ovn/commit/559e2cd8cbe25370b6745b90787781d147959a29) update review problems for ovn_monitor
 * [1c356a365](https://github.com/kubeovn/kube-ovn/commit/1c356a3659885d38cf1ac840b6e6d7e237f99967) monitor: add more dashboard
 * [aa7b20d75](https://github.com/kubeovn/kube-ovn/commit/aa7b20d75a0ff848cf99439ac36ce48d840fdb5d) chore: add vendor
 * [97d64f934](https://github.com/kubeovn/kube-ovn/commit/97d64f9342d77463fd7476d9511f4e924303ff8e) Updated Dockerfile.dpdk1911 to use Centos8 and DPDK19.11.4
 * [b4aa989da](https://github.com/kubeovn/kube-ovn/commit/b4aa989da6b93012e0c17410547cdb85970a4331) fix: CodeQL scan warning
 * [a27e17603](https://github.com/kubeovn/kube-ovn/commit/a27e176039549b00d8d914c387f565013f5315d3) fix: ipt wrong order and add cluster route
 * [9eb96dd78](https://github.com/kubeovn/kube-ovn/commit/9eb96dd788d3d5fbbabbadb9a867257e719a8be5) opt: only allow specifies default subnet
 * [0da634e8f](https://github.com/kubeovn/kube-ovn/commit/0da634e8fe592d16223257cd67cc5248331a21aa) chore: reduce image size
 * [93bf54235](https://github.com/kubeovn/kube-ovn/commit/93bf5423534d8f775be65b9fd7252d6a05677879) feature: Support for namespace binding multiple subnets
 * [e37159c23](https://github.com/kubeovn/kube-ovn/commit/e37159c23dea126b91aae598ce03c67fcc23935f) docs: fix multi nic subnet options
 * [c35a159b9](https://github.com/kubeovn/kube-ovn/commit/c35a159b974e321432e562b94a82ebd324271e7d) docs: add pinger/controller/cni metrics
 * [7f5b42374](https://github.com/kubeovn/kube-ovn/commit/7f5b423742d743ff4a9213df19c685f49c3532f6) fix: add default ssl var for compatibility
 * [59b706964](https://github.com/kubeovn/kube-ovn/commit/59b706964218988a5d6e5fb7623436a3a8a831df) Add monitor doc
 * [bb130cac2](https://github.com/kubeovn/kube-ovn/commit/bb130cac2cc0c2cc961f94750a72bc510d7d78fe) fix: ipv6 network format when update subnet
 * [dc62d1050](https://github.com/kubeovn/kube-ovn/commit/dc62d10506a7a9f9ff2f2f9d8b6514d14b3d008f) fix: ipv6 len mismatch
 * [6088851d4](https://github.com/kubeovn/kube-ovn/commit/6088851d49c5d21f37117a85c8e3a3b69b6a37f1) chore: add version info
 * [88001376d](https://github.com/kubeovn/kube-ovn/commit/88001376db58d8640ac67e294b00a25c57613cea) metrics: add ovs client latency metrics
 * [3cafd5f8a](https://github.com/kubeovn/kube-ovn/commit/3cafd5f8a5d57c3f77ee9ecd0b0ec4bbcbee09aa) Add OVN/OVS Monitor
 * [89567776a](https://github.com/kubeovn/kube-ovn/commit/89567776afcc203a4ed4aac76468d7b84e5968ce) docs: performance test method
 * [0c975e343](https://github.com/kubeovn/kube-ovn/commit/0c975e343bae483882d36b4c8480337fb6c971c8) fix: wrong port porto for udp
 * [f3759b78e](https://github.com/kubeovn/kube-ovn/commit/f3759b78e0d186fac611a2637b79317b63d3c7e4) docs: add descriptions of local files
 * [b46acd6c3](https://github.com/kubeovn/kube-ovn/commit/b46acd6c38f136ee0ca3e9f53265f239480cab81) ci: add github code scan
 * [2444d51ae](https://github.com/kubeovn/kube-ovn/commit/2444d51aefb075356b7b6665c593ffdfb83f19db) fix: do not adv join cidr when enable ovn-ic
 * [292bf4caf](https://github.com/kubeovn/kube-ovn/commit/292bf4caf51e2acba1226d64504a2c16300d0cb2) perf: remove default acl rules
 * [20e82c393](https://github.com/kubeovn/kube-ovn/commit/20e82c39388c7dcefebff4e6bb3ba720b3df5fb2) prepare for next release
 * [9324491cc](https://github.com/kubeovn/kube-ovn/commit/9324491cc02ee4e7b798f565f13bf39d52969205) fix: use internal IP when node connect pod
 * [c1870c1ac](https://github.com/kubeovn/kube-ovn/commit/c1870c1acda04d567d7aa2c40fa0bd3f0bdbeadb) ci: change to docker buildx action
 * [a1976650c](https://github.com/kubeovn/kube-ovn/commit/a1976650c25566fe37bc6274cf2f7c3db95dab47) fix: delete pod when marked with deletionTimestamp
 * [c3c4f1c5b](https://github.com/kubeovn/kube-ovn/commit/c3c4f1c5b00b62e7475e5defacac14adbb3bda07) fix: remove not alive pod in pg
 * [c0a34b842](https://github.com/kubeovn/kube-ovn/commit/c0a34b842eb4cf13121fec2f9f69c579019d6b84) release: prepare for release 1.5.0
 * [955484579](https://github.com/kubeovn/kube-ovn/commit/955484579832c8386a69baf3e82d231b8de7614d) perf: use podLister to optimize k8s calls
 * [6635f9302](https://github.com/kubeovn/kube-ovn/commit/6635f9302c33824dc126592e722c7cc4a6a08bb0) chore: enable ssl to default ci tests
 * [5f29fc307](https://github.com/kubeovn/kube-ovn/commit/5f29fc307c78c30d22a248972b072a562a50906e) security: change ovsdb file access to 600
 * [0e0a6887a](https://github.com/kubeovn/kube-ovn/commit/0e0a6887ac4c0e550aa7830b86070f18bb7bd0f8) docs: improve hw-offload
 * [a1a215dcc](https://github.com/kubeovn/kube-ovn/commit/a1a215dcc72314202879833d4ab68953ea23d04b) feat: support db ssl communication
 * [e7a88c113](https://github.com/kubeovn/kube-ovn/commit/e7a88c1133857eb1d8d3d2f100dd8a9b8f9059fb) diagnose: show nb/sb/node info
 * [090624fd9](https://github.com/kubeovn/kube-ovn/commit/090624fd9e00e1a2798fc3ac5671925f8dac82ee) fix: pinger diagnose should use cmd args
 * [fae393e3b](https://github.com/kubeovn/kube-ovn/commit/fae393e3b491153a5e49f23eea8168e3982bab51) fix: ipv6 get portmap failed again
 * [b74189fee](https://github.com/kubeovn/kube-ovn/commit/b74189feea76dc7acac069e6d0e85699666410e1) fix: ipv6 get portmap failed
 * [f1c2f9952](https://github.com/kubeovn/kube-ovn/commit/f1c2f99528f264941b50d97de0c574db2ba3aa67) fix: delay mv cni conf to when cniserver is ready
 * [98bb7510c](https://github.com/kubeovn/kube-ovn/commit/98bb7510cd0c6aef5d02aa7fe7524bf5786b65fe) chore: update kind and kube-ovn-cni updateStrategy
 * [646404216](https://github.com/kubeovn/kube-ovn/commit/646404216812dfff8461a2e6e9c019ed68f86a51) monitor: add cni grafana dashboard
 * [38adc18f1](https://github.com/kubeovn/kube-ovn/commit/38adc18f18cb8744b381fc67af62d297e80daa2c) monitor: add more kube-ovn-cni metrics
 * [36e9091d1](https://github.com/kubeovn/kube-ovn/commit/36e9091d10ff3c6799629fbb871ed9a69a915058) feat: update pinger dashboard
 * [ab736d8f5](https://github.com/kubeovn/kube-ovn/commit/ab736d8f523967ad31a43b5b0f7e566728e20c7d) fix: issues with vlan underlay gateway
 * [2e5f0ecbe](https://github.com/kubeovn/kube-ovn/commit/2e5f0ecbe48bcbe6c8fa98d550d2f45ab27592f3) feat: set more metadata to interface external_ids
 * [77c4a5f29](https://github.com/kubeovn/kube-ovn/commit/77c4a5f2932c66d304d2f87332838df7d8953d63) feat: grace stop ovn-controller
 * [ebfc1530d](https://github.com/kubeovn/kube-ovn/commit/ebfc1530dfaeab15127ffb0f234c9756a27cd294) refactor: fix bridge-mappings and refactor vlan code
 * [729ed3c7a](https://github.com/kubeovn/kube-ovn/commit/729ed3c7a502460aeadfb0cd9785832ad44f797c) fix: allow mirror config update
 * [84bb3c838](https://github.com/kubeovn/kube-ovn/commit/84bb3c8380c28f3b323dca5b0bbc91e8a77ec66e) fix: cleanup v6 iptables and ipset
 * [da4937172](https://github.com/kubeovn/kube-ovn/commit/da4937172279e42eda90eb4b8109eec263ba869a) docs: add gateway docs and optimize others
 * [ece4219f7](https://github.com/kubeovn/kube-ovn/commit/ece4219f7bc54da03f9678095335b5a8b0a3accb) feat: integrate ovn sfc
 * [2b2e7a9a7](https://github.com/kubeovn/kube-ovn/commit/2b2e7a9a788ee7cc11d44f7b6f821aec6f9086dc) feat: support pod snat
 * [7a60b569c](https://github.com/kubeovn/kube-ovn/commit/7a60b569c9e2f431eb98c948061cb35e796bfd87) prepare for next release
 * [e9933619e](https://github.com/kubeovn/kube-ovn/commit/e9933619e913f4bcd2136168faf2d78f9b007629) fix: ovn-ic-db restart failed
 * [115c12668](https://github.com/kubeovn/kube-ovn/commit/115c126684972af06f2eb0019bc25031004045f3) fix: stop ovn-ic when disabled
 * [e98614444](https://github.com/kubeovn/kube-ovn/commit/e98614444005451184e260e96ef0032ee8f9ae98) fix: use nodeName as chassis hostname
 * [0f973a5ad](https://github.com/kubeovn/kube-ovn/commit/0f973a5ad897f2c6b70eee404772852d911bcce1) prepare for 1.4 release
 * [78ab9b1e2](https://github.com/kubeovn/kube-ovn/commit/78ab9b1e2d21baa0c765529a2aebc048017af04b) fix: do not gc learned routes
 * [3ddb9614d](https://github.com/kubeovn/kube-ovn/commit/3ddb9614d9e1fee9b08f1eb5de612fc23750336c) chore: add psp
 * [f847e5be6](https://github.com/kubeovn/kube-ovn/commit/f847e5be666dfe7a8ff790013f539686aac56b9b) perf: apply udp improvement
 * [a8f0d2285](https://github.com/kubeovn/kube-ovn/commit/a8f0d2285e44ff0c035cc149e0ead9c6b2c4d08e) chore: sync pre-1.16 install.sh
 * [0918e9a20](https://github.com/kubeovn/kube-ovn/commit/0918e9a2073c16fb78908b5aed79ae9116e768db) ci: use go 1.15
 * [f43a10272](https://github.com/kubeovn/kube-ovn/commit/f43a102729c28cc97cc9f6e900598f05f90e782f) fix: add prob timeout to wait script finish
 * [c5ca0b1b7](https://github.com/kubeovn/kube-ovn/commit/c5ca0b1b70c1d56cf21daec590633210913bdc32) resolve review problem
 * [28d5a8aab](https://github.com/kubeovn/kube-ovn/commit/28d5a8aab4a75e983afb7d80bba026efa99524d0) chore: suppress verbose logs
 * [df54b0d13](https://github.com/kubeovn/kube-ovn/commit/df54b0d130394dc3be8a02932db9db170bde91b7) fix: do not gc ic logical_switch
 * [b9ab4d661](https://github.com/kubeovn/kube-ovn/commit/b9ab4d661f48ded7a9f5bc8bfb93e86ec690f8c1) fix: only gc VIF type logical_switch_port
 * [731fef992](https://github.com/kubeovn/kube-ovn/commit/731fef9924cb97ec618e266bb12078ccb1da38ec) docs: update docs
 * [e9ae40a98](https://github.com/kubeovn/kube-ovn/commit/e9ae40a98c0fa6214db4c1e7741d10c99d8f6884) chore: add back lflow reduction optimization
 * [022c79037](https://github.com/kubeovn/kube-ovn/commit/022c79037ba0adaf880b8ebf560d115885db7036) chore: update ovs to 2.14.0
 * [8e93c054e](https://github.com/kubeovn/kube-ovn/commit/8e93c054e252ce89dd0bc5480b64f26a1f248f22) fix: remove duplicated gcLogicalSwitch
 * [c3b7457a1](https://github.com/kubeovn/kube-ovn/commit/c3b7457a12728108e399705287c81df53ae577f7) fix: modify src-ip route priority
 * [e0096f9b9](https://github.com/kubeovn/kube-ovn/commit/e0096f9b9fc7c269c2da0d2d31fc0522c46da10f) fix: missing session lb to logical switch
 * [6fbcc198a](https://github.com/kubeovn/kube-ovn/commit/6fbcc198a1e7360955d8a15f839f406c8574bb32) feat: ovn-ic integration
 * [0ea62c163](https://github.com/kubeovn/kube-ovn/commit/0ea62c16375a890f4d6ea94bd1bf2f761b6a020c) fix:resolve gosec check problem
 * [b2d0393b2](https://github.com/kubeovn/kube-ovn/commit/b2d0393b2eb140ecbd0e38a6a96a0b7c79e5dbe4) feat: do not perform masq on external traffic
 * [4e1ad1260](https://github.com/kubeovn/kube-ovn/commit/4e1ad1260d1fa40542d2a23ba8f19f5e649aacfe) chore: fix patch failure
 * [a7c460a48](https://github.com/kubeovn/kube-ovn/commit/a7c460a48efe275bf82f6fac2f86f402f0c551c2) fix: subnet acl might conflict if allowSubnets and subnet cidr cover each other
 * [0dd85e469](https://github.com/kubeovn/kube-ovn/commit/0dd85e469e62c56ea67cc6d3c2b4bb83af67def6) feat: acl log drop packets
 * [6d048632f](https://github.com/kubeovn/kube-ovn/commit/6d048632fdbbea7c4c1258b80c46b7f021b14e66) chore: remove juju log dependency
 * [9535c26ba](https://github.com/kubeovn/kube-ovn/commit/9535c26ba2d65ad2c92e63a32685b6c9d7a9bf93) feat: gw switch from overlay to underlay
 * [4b0955808](https://github.com/kubeovn/kube-ovn/commit/4b095580830573e4c09d0debc6c480aeff1ced29) chore: prepare for 1.4 release
 * [c9d07e1d5](https://github.com/kubeovn/kube-ovn/commit/c9d07e1d5c3b23584993900807a862d11ffbf038) fix: prevent update failed logs
 * [a98ec5bd1](https://github.com/kubeovn/kube-ovn/commit/a98ec5bd19393587991853dde0ff6ac191613eb4) fix: ko use external-ids to find related nic
 * [1cad39cef](https://github.com/kubeovn/kube-ovn/commit/1cad39cef27c303737c122d406ecb0fbd994d68f) fix: forward accept rules
 * [45d30713d](https://github.com/kubeovn/kube-ovn/commit/45d30713d07ffff760127cf70a806bf5fc4446d9) chore: add build date
 * [c99532349](https://github.com/kubeovn/kube-ovn/commit/c995323497a1c9b0369e1b9164c9a83c5fa0a1c2) release: update 1.3.0 docs
 * [34627a669](https://github.com/kubeovn/kube-ovn/commit/34627a66953fceca9e905a53f1b642ba3a6ce687) fix: call appendMssRule function to resolved mss according problem
 * [bb961ae5c](https://github.com/kubeovn/kube-ovn/commit/bb961ae5c63c07cc597a54c3bcdbb80bd2b4c002) dpdk: add kmod, pdump and proc-info tools
 * [cf47ee1b9](https://github.com/kubeovn/kube-ovn/commit/cf47ee1b91db83551ae50fc3bf7a1d61dfb88236) fix: ci image tags
 * [467681796](https://github.com/kubeovn/kube-ovn/commit/46768179604d756186057d20475f637a239e51c8) chore: optimize dpdk build
 * [5c1076875](https://github.com/kubeovn/kube-ovn/commit/5c1076875ef7d8d330eb08311eeb1c8df34cedc5) docs: add hw-offload docs and resolve some issues
 * [e64c61327](https://github.com/kubeovn/kube-ovn/commit/e64c61327d12b66e7c0ea0aa20949b288dcfbec1) fix: if sriov device, do not delete the host nic
 * [f55c3fba2](https://github.com/kubeovn/kube-ovn/commit/f55c3fba21ba0f1f513ae834497eed76d623a17f) fix: use keymutex to serialize pod add/delete operation
 * [d438574d4](https://github.com/kubeovn/kube-ovn/commit/d438574d43fd8c0d108d54fed16a94bcb0b0f8c3) feat: assign a pod as the gw
 * [1806a5722](https://github.com/kubeovn/kube-ovn/commit/1806a57226d87152fead314f3df596510e3a8db9) ci: add arm build to normal ci process
 * [5aed1ef1c](https://github.com/kubeovn/kube-ovn/commit/5aed1ef1ca6a28ebc68fc7ca994b0a02a7f3319d) ci: add unfixed cve
 * [19201a369](https://github.com/kubeovn/kube-ovn/commit/19201a36999041309f4eea2cf12adda266d07da4) ci: arm64 build accelerate
 * [63fbc0086](https://github.com/kubeovn/kube-ovn/commit/63fbc00869c02ef62a1d861a9f297687bcd7937d) chore: add logs to sriov interface
 * [82140c931](https://github.com/kubeovn/kube-ovn/commit/82140c9313be9e50daf94f8c299674249007c4a1) ci: add ipv6 install e2e
 * [c3814c725](https://github.com/kubeovn/kube-ovn/commit/c3814c725f783aecf869f81f560fe5b463e0a86d) feat: recycle lsp at runtime
 * [3f9d7c928](https://github.com/kubeovn/kube-ovn/commit/3f9d7c928dfcfbebaff28409849bcac0e18bdeff) fix: qos error
 * [e460541dd](https://github.com/kubeovn/kube-ovn/commit/e460541dd3fdc938d43b17b2c71e81e5a7d03b8c) fix: variable error
 * [de9493f2f](https://github.com/kubeovn/kube-ovn/commit/de9493f2f03e7a365120a7e3a2285bff48f2753d) ci: modify cache usage
 * [1994e5c30](https://github.com/kubeovn/kube-ovn/commit/1994e5c30b1609a891b2f3222e25039e9dd41378) ci: save ci time
 * [5c4d5a3cf](https://github.com/kubeovn/kube-ovn/commit/5c4d5a3cfa64479c154ab3379f3f6dc3175ac546) chore: use j2 to render different kind.yaml
 * [d1a184ef9](https://github.com/kubeovn/kube-ovn/commit/d1a184ef9901055c79b0ba97dce2dadda4ba20aa) fix: set qlen for ovn0
 * [a2d969e8f](https://github.com/kubeovn/kube-ovn/commit/a2d969e8f99f406ca1651142cd9d8a6633e0773d) prepare for 1.3 release
 * [3a018a86e](https://github.com/kubeovn/kube-ovn/commit/3a018a86e9dde7308f76a1fb30902184a0150cb9) chore: update build.sh
 * [be7c68f23](https://github.com/kubeovn/kube-ovn/commit/be7c68f23fbdb18e166479608eb2b4d663bce8c3) fix: log error
 * [31723f669](https://github.com/kubeovn/kube-ovn/commit/31723f669e54f662d9d56c9833c7386a05fe768f) chore: check ovn-sb connectivity from ovn-ovs probe
 * [d017f1f2d](https://github.com/kubeovn/kube-ovn/commit/d017f1f2d0e6a0a0eef96e8ddcfa8b907a7007f1) fix: available ips calculation issues
 * [309c80800](https://github.com/kubeovn/kube-ovn/commit/309c808005ae3d38923569c9e9d0b39c0a082a66) perf: add hw offload
 * [4b8faede4](https://github.com/kubeovn/kube-ovn/commit/4b8faede48168fb21ac39cbb3084a40bd7a3777f) docs: add gateway qos doc
 * [32a9af2bd](https://github.com/kubeovn/kube-ovn/commit/32a9af2bde14322acd4a77661566f75e1aebf1a4) ci: remove master taint
 * [3865220df](https://github.com/kubeovn/kube-ovn/commit/3865220df86b2af451642037f9a12c1ad618bf28) chore: update cni dependencies
 * [8e032392d](https://github.com/kubeovn/kube-ovn/commit/8e032392d62d5a8c4b9545f66ec78ba214cbb9db) feat: session service
 * [34b7cba75](https://github.com/kubeovn/kube-ovn/commit/34b7cba751209975063c155fac2ecbf919348a6e) Revert "perf: use policy-route to replace src-ip route"
 * [1d13d5c38](https://github.com/kubeovn/kube-ovn/commit/1d13d5c38ab51d4c5829a3a76ce7eaa0552d8ed6) Revert "fix: ipv6 policy route"
 * [658136408](https://github.com/kubeovn/kube-ovn/commit/658136408bc58902645b17ac7b8d56ff07a0c09a) Revert "fix: reset address_set when delete subnet"
 * [e6817a651](https://github.com/kubeovn/kube-ovn/commit/e6817a65163ede10b76e5229896af0c6fafe5e1a) fix: reset address_set when delete subnet
 * [dbc968ca5](https://github.com/kubeovn/kube-ovn/commit/dbc968ca57823a5b939f6ae133581bb2c8a7f7f4) test: statefulset without ippool
 * [9440a11ff](https://github.com/kubeovn/kube-ovn/commit/9440a11ff09ab344224a0e7c0a97ceaa9c83cdbb) match apps/* statefulset
 * [ca122027b](https://github.com/kubeovn/kube-ovn/commit/ca122027b0dcf94b736306a982fc48bd2cec9981) fix: ipv6 policy route
 * [54acd0c35](https://github.com/kubeovn/kube-ovn/commit/54acd0c3576b7a09b97d6ffb34b2ed4da582260a) feat: support gw qos
 * [b8f032487](https://github.com/kubeovn/kube-ovn/commit/b8f032487624e5a189a7e0675a66640316a12673) perf: use policy-route to replace src-ip route
 * [83dc420e4](https://github.com/kubeovn/kube-ovn/commit/83dc420e4f8c21a1014de6a98bc7feb6ce0c0f2e) Solve the problem of non-standard statefulset creation mode
 * [32e6d5728](https://github.com/kubeovn/kube-ovn/commit/32e6d57282dd679b685d75ba0695244fe25bcf76) fix: arm64 build missing env
 * [c93f0d848](https://github.com/kubeovn/kube-ovn/commit/c93f0d848d99f20d0ecdee178521f545b45cfcef) action: use commit as image tag
 * [732b240c0](https://github.com/kubeovn/kube-ovn/commit/732b240c07a1134774e27606b1d40192587eddc8) Add libatomic to docker image
 * [9d5294bb9](https://github.com/kubeovn/kube-ovn/commit/9d5294bb96f4a0d19dcc742e9412e7f8631c3c9a) chore: save disk space when building
 * [4b1f5244c](https://github.com/kubeovn/kube-ovn/commit/4b1f5244c3ce141a0c74f15893b7615bec996777) chore: change crd form v1beta1 to v1
 * [e6fb0fcb9](https://github.com/kubeovn/kube-ovn/commit/e6fb0fcb9947c36d36d4bae1e3e64bca3add3c9c) kubectl-ko: add ovs-tracing info
 * [61aa3ba2b](https://github.com/kubeovn/kube-ovn/commit/61aa3ba2b47071de74ce10bd42eff43eb22520eb) pinger: add metrics to resolve external address
 * [ef0f3b279](https://github.com/kubeovn/kube-ovn/commit/ef0f3b279edcc62c520144b9d79f7138d27666d8) chore: update ovn to 20.06
 * [961f5f1ab](https://github.com/kubeovn/kube-ovn/commit/961f5f1ab42de8f40c28640565dbb894d53ef79b) update changelog
 * [85f2e0e0f](https://github.com/kubeovn/kube-ovn/commit/85f2e0e0f85e8b9c88730c5d38eada3f91c19ab0) fix some version in docs
 * [f989bdd8f](https://github.com/kubeovn/kube-ovn/commit/f989bdd8f449ec65381fd7fdfee7ecf95d3195a8) fix: rename variable
 * [990bf9839](https://github.com/kubeovn/kube-ovn/commit/990bf9839cc061a1adab0978e9d9e32222ade2ae) fix: minor fix
 * [8d7045b3a](https://github.com/kubeovn/kube-ovn/commit/8d7045b3a11c34cad82382b95d273dd6af1a869f) feat: use never used address first to reduce conflict
 * [db2516c2d](https://github.com/kubeovn/kube-ovn/commit/db2516c2d899803518880d7036aff4154619b2e9) ci: use tmpfs to accelerate e2e
 * [792723761](https://github.com/kubeovn/kube-ovn/commit/792723761535fb01cf20391af0b0d9f61a26da67) fix: create/delete order might lead ip conflict
 * [b27d75455](https://github.com/kubeovn/kube-ovn/commit/b27d754551ab336961f951540a905de4e4009a54) ci: do not push image when pr
 * [a1f53e67d](https://github.com/kubeovn/kube-ovn/commit/a1f53e67d76da37322dac333cf015e0bf07ac615) clean up all white noise
 * [a4f40370d](https://github.com/kubeovn/kube-ovn/commit/a4f40370d19c0b2712a8d9fefda388e4b6eb1a84) security: update yum repo
 * [270c825cd](https://github.com/kubeovn/kube-ovn/commit/270c825cd599c7ffa52632c7f897744d3722a860) fix node's annottaions overwrited incorrectly
 * [5adc5a440](https://github.com/kubeovn/kube-ovn/commit/5adc5a440a06aec9715c3e1f232b84f4332e8013) Fix typo in multi-nic.md
 * [3ac92a157](https://github.com/kubeovn/kube-ovn/commit/3ac92a157d2bb3015963a9e3304e002cdc184d98) Userspace-CNI updates in dpdk.md
 * [76e72b7ea](https://github.com/kubeovn/kube-ovn/commit/76e72b7ea09dab0a6299d9510a50e2b527ca92a2) Remove empty lines from DPDK Dockerfile
 * [9b5c018ac](https://github.com/kubeovn/kube-ovn/commit/9b5c018aca9e7968b6e98d595eca1dfe13ed9f5a) security: update loopback to fix CVE
 * [bd1f2acf7](https://github.com/kubeovn/kube-ovn/commit/bd1f2acf7c511af1cf670f6119cf7c4ea735368c) Make OVS-DPDK start script more robust
 * [3bfc39f82](https://github.com/kubeovn/kube-ovn/commit/3bfc39f82050e16e37a20a843185ca800fa7940e) Reduce DPDK image size
 * [4917afe98](https://github.com/kubeovn/kube-ovn/commit/4917afe98ed89302c1481d25fb1622d79acf54f1) fix: add back privilege for ipv6
 * [8121afd60](https://github.com/kubeovn/kube-ovn/commit/8121afd6037a3414b2ca1a0ff8f3f004d47695a7) Config support for OVS-DPDK
 * [ad30e6870](https://github.com/kubeovn/kube-ovn/commit/ad30e6870a2edc35fd6f624d317fcd5bbbeb2ed3) security: add trivy scan and fix image CVEs
 * [06256a09b](https://github.com/kubeovn/kube-ovn/commit/06256a09bce362754b680fa8dc1eba6a93e65830) docs: modify arm build
 * [9d2e64a40](https://github.com/kubeovn/kube-ovn/commit/9d2e64a40bab54394e60f90a2f3e8bee38dd35f2) docs: update development
 * [bd9757681](https://github.com/kubeovn/kube-ovn/commit/bd975768146fc27dc7ffcdecb2de2ceef9507482) refactor: use ovs.Exec replace raw command
 * [32024ba88](https://github.com/kubeovn/kube-ovn/commit/32024ba88963a56593cd112717ea433abf6cfc30) chore: add gosec to audit code security
 * [1db9046d6](https://github.com/kubeovn/kube-ovn/commit/1db9046d6803245d5ff3a5cbdbaff6e0a7c6794f) prepare for next release
 * [aa72ba6ce](https://github.com/kubeovn/kube-ovn/commit/aa72ba6ce070ff7406fec3c0d0bccebdd08508ca) fix: arm build
 * [628f5c5e1](https://github.com/kubeovn/kube-ovn/commit/628f5c5e154e0070f9f6de7dccfb3ae9d3f3a796) fix: change version in install.sh
 * [280a1bd39](https://github.com/kubeovn/kube-ovn/commit/280a1bd39ea045ece2e8ef05259a023846873076) chore: prepare for release 1.2
 * [4342187d8](https://github.com/kubeovn/kube-ovn/commit/4342187d8f72ff18511cbfe7ebb3f70cf357a064) chore: prepare for release 1.2
 * [4a52bb43d](https://github.com/kubeovn/kube-ovn/commit/4a52bb43df3eb9fbbd1f56aa43f02085c87363b9) DPDK doc update and small image reduction
 * [b055cc68a](https://github.com/kubeovn/kube-ovn/commit/b055cc68af72111431130376ada2494102eb70fe) Add OVS-DPDK support, for issue 104
 * [f7fdd2dca](https://github.com/kubeovn/kube-ovn/commit/f7fdd2dcae463837bc98f51b08cb682372d1762e) fix: pod get deleted between configure nb and patch pod
 * [e13dc5ac0](https://github.com/kubeovn/kube-ovn/commit/e13dc5ac0ec194f90302f31a09c37492de4beab4) fix: native vlan and delete subnet issues
 * [44b5a6a7c](https://github.com/kubeovn/kube-ovn/commit/44b5a6a7cf53646b14f02524e3060238e0c0fce6) fix: trigger github action when dist dir change
 * [3a2ee051c](https://github.com/kubeovn/kube-ovn/commit/3a2ee051cba6d1f55f81abc2e37e6167ad35fedf) fix: update ovn patch
 * [6e1589cc5](https://github.com/kubeovn/kube-ovn/commit/6e1589cc50c009b07293bb83f726483586ac59ca) chore: improve log
 * [00f984890](https://github.com/kubeovn/kube-ovn/commit/00f9848902d14c4c1fec406001fd3c115f2bb300) fix: gc lsp for pod that not alive
 * [701e9efdc](https://github.com/kubeovn/kube-ovn/commit/701e9efdc3cde047b7336511c1a5f61538c13681) feat: support underlay without vlan encap
 * [83ad499fd](https://github.com/kubeovn/kube-ovn/commit/83ad499fdda2f32b7a1b6d0774ceddedf91aa5ac) chore: optimize kube-ovn-cni log
 * [84b6cdcf6](https://github.com/kubeovn/kube-ovn/commit/84b6cdcf64e7bd39aede73f87fc0d070869f6dfb) fix: gc node lsp
 * [7aafd9447](https://github.com/kubeovn/kube-ovn/commit/7aafd94477588c8d4fc022a7d0d6ce8676e8bd87) chore: remove vagrant
 * [92ccf729f](https://github.com/kubeovn/kube-ovn/commit/92ccf729f46adaee18ab00102aef496cfe2a6999) fix: dst route policy might be empty
 * [6c89a0467](https://github.com/kubeovn/kube-ovn/commit/6c89a046755837dc93e2c73cc258efc4a675ba4e) feat: in vlan mode if physical gateway exists, no need to create a virtual one
 * [1d5c6958a](https://github.com/kubeovn/kube-ovn/commit/1d5c6958aa83e83f4ad4bfa22ed8a58be69900c8) perf: add amd64 compile flags back
 * [b0f0947d0](https://github.com/kubeovn/kube-ovn/commit/b0f0947d010b81f486148c293cd22e323e728922) fix: init ipam before gc, other wise routes will be deleted
 * [dbc23c5e9](https://github.com/kubeovn/kube-ovn/commit/dbc23c5e99057a547cdd216487771095c14cf60a) fix: patch ovn to lower src-ip route priority to work with ovn-ic
 * [5a7638200](https://github.com/kubeovn/kube-ovn/commit/5a7638200612b179414342252cf1d0fd27de78ec) fix: return early if allocation is not ready
 * [b03c37684](https://github.com/kubeovn/kube-ovn/commit/b03c37684ded0dd89df7f3ea84b985aa09459a12) chore: remove networks crd
 * [2853438c4](https://github.com/kubeovn/kube-ovn/commit/2853438c4ad3632082ad9cb2d8b40d9325fb20a1) perf: remove more stale lflow
 * [0665f2e8f](https://github.com/kubeovn/kube-ovn/commit/0665f2e8fee7ff6b5a5f5d68731b7d15d967c723) ci: run ut and e2e in github action
 * [e71b68c01](https://github.com/kubeovn/kube-ovn/commit/e71b68c01e6045aa291fc60d398c531d44895b43) fix: check svc and endpoint protocol
 * [508eb7a28](https://github.com/kubeovn/kube-ovn/commit/508eb7a28ff3138010ac73808a95c3ca427f9cdb) perf: reduce lflow count
 * [5f8b9b406](https://github.com/kubeovn/kube-ovn/commit/5f8b9b406427d90fb5a3aca097839172aa7c57b1) fix: when podName or namespace contains dot, lsp cannot be deleted correctly
 * [27c72560c](https://github.com/kubeovn/kube-ovn/commit/27c72560c6caa58b1a125d020c7af06c19584c73) fix: wrong subnet status
 * [f0b17a69b](https://github.com/kubeovn/kube-ovn/commit/f0b17a69b0f1b0945e015232713170c01e75f56b) feat: change pod route when update gateway type
 * [13283daf3](https://github.com/kubeovn/kube-ovn/commit/13283daf39c4f1a7231d88f2e5d5e5c0ea72adf1) feat: refactor subnet and allow cidr change
 * [23821d6cb](https://github.com/kubeovn/kube-ovn/commit/23821d6cb256ea8a9adaf71d5af389d12e3b8207) fix: use kubectl to avoid tls handshake error
 * [e647cc6c0](https://github.com/kubeovn/kube-ovn/commit/e647cc6c05960b77a11781a2de98a51303a89ccd) chore: reduce logs
 * [aef4336d1](https://github.com/kubeovn/kube-ovn/commit/aef4336d1c13e67422508f5a504ac3016fb51122) feat: only show error log of kube-ovn-controller
 * [a9ab0bc25](https://github.com/kubeovn/kube-ovn/commit/a9ab0bc2531ca472c9c3f3d0608cd138db9f8bd2) fix: map concurrent panic
 * [2dd13b232](https://github.com/kubeovn/kube-ovn/commit/2dd13b2328a056cf343914f635da56cffb56057b) fix: ipv6 related issues
 * [86c443e71](https://github.com/kubeovn/kube-ovn/commit/86c443e71d72e624518ca71d285e115979cf626a) fix: validate if subnet cidr conflicts with svc ip
 * [eb4cb1b3a](https://github.com/kubeovn/kube-ovn/commit/eb4cb1b3a2ef58f566002b982a656f2bab00295d) fix: validate if node address conflict with subnet cidr
 * [7f595ee03](https://github.com/kubeovn/kube-ovn/commit/7f595ee030c6cc2e4a6099d47f592659c0a66eae) feat: github action
 * [1046b5727](https://github.com/kubeovn/kube-ovn/commit/1046b5727767a804f2001b9342cabf5211c4e696) fix: wait node annotations ready before handle pods
 * [7a0151cc2](https://github.com/kubeovn/kube-ovn/commit/7a0151cc21919d1b89b076c7262169418386d1c7) fix: check ovn-nbctl socket in new dir
 * [0dc76768b](https://github.com/kubeovn/kube-ovn/commit/0dc76768bfba4105ecece940e29a00142b02721a) fix: error log found in scale test
 * [047159438](https://github.com/kubeovn/kube-ovn/commit/047159438c7c500b54dc5e5aac2f6c8375605acb) fix: concurrent panic
 * [da14eaebb](https://github.com/kubeovn/kube-ovn/commit/da14eaebbe4950a369db0200e9531f20d79e0024) feat: use bgp to announce pod ip
 * [909b5a00d](https://github.com/kubeovn/kube-ovn/commit/909b5a00da781f56123f77add7016afc56a3a8be) release 1.1.1
 * [ab834b5a1](https://github.com/kubeovn/kube-ovn/commit/ab834b5a1d817a0dd9f7cd929c5ab6c74d32db6b) fix: labels might be nil
 * [0c0824db9](https://github.com/kubeovn/kube-ovn/commit/0c0824db916f5e36e1437fd4acdb347476a2f86a) fix: ping output format
 * [ce27fb31c](https://github.com/kubeovn/kube-ovn/commit/ce27fb31c44ce78f78fdcded6020fc1576ff561f) monitor: make graph more sensitive to changes
 * [9b05fccf2](https://github.com/kubeovn/kube-ovn/commit/9b05fccf2942a453ff818beb6e7e1a5bee72e074) docs: update vlan docs
 * [d0544d896](https://github.com/kubeovn/kube-ovn/commit/d0544d8969f7ece8fa352f5c81a78ba03b3fbd9b) docs: update docs
 * [28aef840e](https://github.com/kubeovn/kube-ovn/commit/28aef840e09a9bbf35883eb621c72f55f438691d) feat: improve install/uninstall
 * [8d8536564](https://github.com/kubeovn/kube-ovn/commit/8d853656496dd94d54f5fc681e2c8b18c5c468ef) refactor: refactor cni-server
 * [d99ffff0e](https://github.com/kubeovn/kube-ovn/commit/d99ffff0e39dfe71bf22e5de606bb42caa2f8813) refactor: controller refactor
 * [8f1f0135d](https://github.com/kubeovn/kube-ovn/commit/8f1f0135dd7188359502db41bcb801b86c051293) feat: modify install.sh for vlan type network
 * [cfe9d2761](https://github.com/kubeovn/kube-ovn/commit/cfe9d27614859cf73ea442a4c28d7d88b3e1c213) feat(vlan): vlan network type
 * [edd0ea81e](https://github.com/kubeovn/kube-ovn/commit/edd0ea81e935c9d5a76f1a7fb1c6871ea1a49511) feat(vlan): vlan network type
 * [c63accf42](https://github.com/kubeovn/kube-ovn/commit/c63accf42e109537b7e26d5b0dbadcb34e701349) fix: yaml indent and ovn central dir
 * [5bc84d7b4](https://github.com/kubeovn/kube-ovn/commit/5bc84d7b4a5b1f825f6b27d6c50e42879846b0c3) docs: chinese wechat info
 * [feaec4ddd](https://github.com/kubeovn/kube-ovn/commit/feaec4ddd63d816a22f5b25cc6bf568b2183fd76) fix: fork go-ping and apply patches
 * [58f73b33f](https://github.com/kubeovn/kube-ovn/commit/58f73b33f4062e25a730ad2570ac328b2291b2d5) chore: update kind node to 1.18 and ginkgo
 * [d274a979a](https://github.com/kubeovn/kube-ovn/commit/d274a979a786847f8eda80e8573083a26da1ef86) docs: add arm build steps
 * [d061fc3ca](https://github.com/kubeovn/kube-ovn/commit/d061fc3ca98e7a8f795c04648782b36c7881478e) fix: mount etc/origin/ovn to ovs-ovn
 * [f8d6fd5c1](https://github.com/kubeovn/kube-ovn/commit/f8d6fd5c1a922c178a2e43a42d6a211ba7865cd4) add support for multi-arch build
 * [953f5be7b](https://github.com/kubeovn/kube-ovn/commit/953f5be7b345a805e6080516fc5ba769a7251299) docs: change the cidr to avoid misunderstanding
 * [5c5b9e08f](https://github.com/kubeovn/kube-ovn/commit/5c5b9e08f4733c931e1eb82c433d67c87038987f) feat: diagnose check if dns/kubernetes svc exist
 * [7c6d6784f](https://github.com/kubeovn/kube-ovn/commit/7c6d6784f417dd5b9bf71376b30be95d3b5d9c98) OVS local interface table mac_in_use row is lower case, but pod annotation store mac in Upper case.
 * [b53a21537](https://github.com/kubeovn/kube-ovn/commit/b53a2153727a8bd6c0e7a1a1fa89a20cd2bc2fd7) prepare for 1.2
 * [0d60df326](https://github.com/kubeovn/kube-ovn/commit/0d60df326555aa06b1a90299872d0a0b73923e78) fix: separate log for no address and wrong address
 * [a4106b2d7](https://github.com/kubeovn/kube-ovn/commit/a4106b2d7eab2346bef84c9afa05d2c35accf291) docs: format docs
 * [de9b003da](https://github.com/kubeovn/kube-ovn/commit/de9b003da042aab313ee2ae9d4f069231c5df846) release 1.1.0
 * [4511a16b6](https://github.com/kubeovn/kube-ovn/commit/4511a16b68cc6919e170e55b357537d57ebb922f) feat: use buildx to reduce image size
 * [370689e79](https://github.com/kubeovn/kube-ovn/commit/370689e79633f05810eea404c15e9831cca6a2e6) test: check host route when add/del a subnet
 * [0df863b67](https://github.com/kubeovn/kube-ovn/commit/0df863b67680b5baebc62b5af72a9bc8bd4f7df0) [DO NOT REVIEW] vendor update: introduce klogr and do some tidy
 * [eeba4c019](https://github.com/kubeovn/kube-ovn/commit/eeba4c01918166017f06f434ca9312cdb3953094) [webhook] init logger for controller-runtime
 * [ae187152a](https://github.com/kubeovn/kube-ovn/commit/ae187152abd6004dcaa0958b7fd2128eaf3ce8ab) test: add node test
 * [e1038d222](https://github.com/kubeovn/kube-ovn/commit/e1038d222d35700c71cc5749d46181d0103235c1) fix: acl and qos issues
 * [a4c81ba78](https://github.com/kubeovn/kube-ovn/commit/a4c81ba78999532fc6a11bb48765ad59abd4075c) feat: expose iface in install.sh
 * [b6967f57b](https://github.com/kubeovn/kube-ovn/commit/b6967f57bd58407576eef6964e11f0d8ace7227e) fix: remove auto checksums
 * [dbc85075a](https://github.com/kubeovn/kube-ovn/commit/dbc85075ab7822b0f7478a5b528792683e339c7a) perf: offload udp checksum if possible
 * [bdb23691c](https://github.com/kubeovn/kube-ovn/commit/bdb23691cce02e8c7a1dab88dac5ab3070a3f2ae) release v1.0.1
 * [cdf4de3f9](https://github.com/kubeovn/kube-ovn/commit/cdf4de3f905029f8617c8657d7d1768df03840bd) perf: add x86 optimization CFLAGS
 * [131181c24](https://github.com/kubeovn/kube-ovn/commit/131181c248681602eb75cd2722f1b30119e943e4) chore: add scripts to build ovs
 * [2b5dd72b6](https://github.com/kubeovn/kube-ovn/commit/2b5dd72b61ef8d877ea43d1ee36079990daaa2bf) fix: lost route when subnet add and is not ready
 * [9032ac843](https://github.com/kubeovn/kube-ovn/commit/9032ac8439422b342c75f684d909607ed8cb5345) fix: ip prefix might be empty
 * [d1654e152](https://github.com/kubeovn/kube-ovn/commit/d1654e152d676598dfb01881b8334711a93da094) chore: reduce image size
 * [464e991ee](https://github.com/kubeovn/kube-ovn/commit/464e991ee8e1aac9785264b6d6bab70ea2b075de) chore: modify nodeSelector label to support k8s 1.17
 * [2814a1d5f](https://github.com/kubeovn/kube-ovn/commit/2814a1d5f846f5aa6ed240daf9d77956dc499ad5) fix: use ovn-appctl to do recompute
 * [0eaedd99f](https://github.com/kubeovn/kube-ovn/commit/0eaedd99fffe3c20b2e3e54c0b13f2749c52ea49) docs: multi nic
 * [dd1923c32](https://github.com/kubeovn/kube-ovn/commit/dd1923c3292c14a2d109b9b3b247481cc3af2a22) feat: ip cr support multi-nic
 * [b2ce6f08e](https://github.com/kubeovn/kube-ovn/commit/b2ce6f08e09f14102f5fac90608c530667cae59d) fix: update in svc 1.1.1.1 may del svc 1.1.1.10
 * [20bb7a784](https://github.com/kubeovn/kube-ovn/commit/20bb7a7842b52f5e6a4a801f23cc3b3eedc52997) feat: add cni side logical to support ipam for multi-nic
 * [1319eb5d4](https://github.com/kubeovn/kube-ovn/commit/1319eb5d436298cd402d8f2a26faf0ddee13b144) feat: add basic allocation function for multus-cni
 * [8f6997a93](https://github.com/kubeovn/kube-ovn/commit/8f6997a93584b45d5acd036f619268549af43afd) fix: only delete pod that restart policy is Always
 * [3a2de9cdc](https://github.com/kubeovn/kube-ovn/commit/3a2de9cdc94eaa2da174a25fd14e0e4e3a44805e) perf: only enqueue updatePod when needed
 * [0f7b9d4ce](https://github.com/kubeovn/kube-ovn/commit/0f7b9d4ceef86e75f35ad78bd06119cfb41d0318) fix: add iptables to accept container traffic
 * [bdd021c03](https://github.com/kubeovn/kube-ovn/commit/bdd021c0347c7adb207d4684e3c631ea97830f75) feat: check kube-proxy and coredns in diagnose
 * [502f18cf2](https://github.com/kubeovn/kube-ovn/commit/502f18cf28a68fb3a739149480e8d549e87c3421) feat: add label param in install script
 * [5a1cf3719](https://github.com/kubeovn/kube-ovn/commit/5a1cf3719832ad306e02062da40eabfbbb0200b5) perf: recycle ip and lsp for pod that in failed or succeeded phase
 * [d19685841](https://github.com/kubeovn/kube-ovn/commit/d1968584156cf77180a0cd5d834e434d937eea63) fix: add inactivity_probe back
 * [417a001b6](https://github.com/kubeovn/kube-ovn/commit/417a001b660f5ba8fb04060cddea9b7b06177245) feat: check if crds exist in diagnose
 * [e65a9d091](https://github.com/kubeovn/kube-ovn/commit/e65a9d091d669de597905cea0629d5995bdb26fa) fix: gc static routes
 * [91829d240](https://github.com/kubeovn/kube-ovn/commit/91829d2404672926364b980cd6281fdbcb9a02ec) fix: still delete lsp if pod not in ipam
 * [7d22430d4](https://github.com/kubeovn/kube-ovn/commit/7d22430d4b93902f7e74210af8ab629001486f1b) fix: delete chassis from sb when delete node
 * [5f5df34e9](https://github.com/kubeovn/kube-ovn/commit/5f5df34e9d721d54bb5408fd1cd460a0919dda7e) fix: missing label selector
 * [9822dba99](https://github.com/kubeovn/kube-ovn/commit/9822dba992442cfa669f944817642b4239bbeb79) feat: add one script installer
 * [479437a35](https://github.com/kubeovn/kube-ovn/commit/479437a35667e228ba605d1c8800867b1294a377) fix: cleanup in offline environment
 * [e707eb96d](https://github.com/kubeovn/kube-ovn/commit/e707eb96d95cb28468e04ecd47a7f0a877296c3b) feat: diagnose check ds/deployment status
 * [3c786f57f](https://github.com/kubeovn/kube-ovn/commit/3c786f57f77d808f310905c7f184b3d4fad5bc77) refactor: the ipam now has lock itself no need for ippool queue
 * [9211486be](https://github.com/kubeovn/kube-ovn/commit/9211486be0a91aa7ea8f55d4d7020276a0f3b42a) fix: if pod is evicted, recycle address
 * [2546deaf7](https://github.com/kubeovn/kube-ovn/commit/2546deaf7fea0721de546946a3e473cb5a109aa9) fix: use uuid to fetch vip
 * [51f06bd64](https://github.com/kubeovn/kube-ovn/commit/51f06bd64ca07a9e958f7da4d32848a3b9495881) refactor ipam
 * [2336dc75f](https://github.com/kubeovn/kube-ovn/commit/2336dc75f546b6fa13700d7359284b477dff47c9) release 1.0.0
 * [7d918f560](https://github.com/kubeovn/kube-ovn/commit/7d918f5600468e0ac5f834b768b78e4cc9d429c1) refactor pod controller
 * [866db995f](https://github.com/kubeovn/kube-ovn/commit/866db995f89a41f9de3ed9e66fd5ace2f6da6071) merge images into one
 * [8296a9e7e](https://github.com/kubeovn/kube-ovn/commit/8296a9e7e4b4deaaae9a28412af961a8f50c96a4) fix:enablebash alias option in Dockerfile CMD scripts
 * [68d87ec23](https://github.com/kubeovn/kube-ovn/commit/68d87ec2327068f4bd5549080d0b9cf2d94dfe50) webhook: use global variables to avoid repeated map constructing
 * [cf2784adf](https://github.com/kubeovn/kube-ovn/commit/cf2784adff77f30237fd1bcf9e66100810a7a313) remove useless fields in webhook.yaml
 * [657b5a295](https://github.com/kubeovn/kube-ovn/commit/657b5a295c9ad7d59677e03620d45262e668eeae) remove leader-election for webhook manager
 * [2bcf0d284](https://github.com/kubeovn/kube-ovn/commit/2bcf0d284f3edf8a041bc4289293c4970a75ffea) feat: update to 20.03.0 ovn
 * [282387945](https://github.com/kubeovn/kube-ovn/commit/282387945fa246f2d2a1f22f0fe377f583d797d8) prepare for 1.0
 * [a036b37b7](https://github.com/kubeovn/kube-ovn/commit/a036b37b790ecf545aedd7c94fd4d60b666cc7db) fix: add back missing lsp gc
 * [44d53c24a](https://github.com/kubeovn/kube-ovn/commit/44d53c24ae740b8053f4d44c7d3befc66878b033) fix: delete lb if it has no backend
 * [b8498a836](https://github.com/kubeovn/kube-ovn/commit/b8498a8364dac6839049cbd3af58ff0553592448) metrics: expose cni operation metrics
 * [a75f99917](https://github.com/kubeovn/kube-ovn/commit/a75f99917742711661dba4db54f780051607ea32) refactor: refactor server.go
 * [c88221ee4](https://github.com/kubeovn/kube-ovn/commit/c88221ee41bfe6b460694b74241c4bda30cde480) fix: disable ovn-nb inactivity_probe
 * [957654f98](https://github.com/kubeovn/kube-ovn/commit/957654f9845aa2a02daaef58dffd132331a457d4) fix: wait for container network ready before cni return
 * [870d20b0a](https://github.com/kubeovn/kube-ovn/commit/870d20b0a0c95c8386aaadac2992132ca128ceca) refactor: refactor controller.go
 * [2885419d8](https://github.com/kubeovn/kube-ovn/commit/2885419d83c9e03103a6ee01244b798860ccdb66) ovn: pick upstream performance patch
 * [115987397](https://github.com/kubeovn/kube-ovn/commit/1159873975902d7d2cb94b63bc026758f7edd5ce) docs: add the development guide and fix the lint
 * [0be255161](https://github.com/kubeovn/kube-ovn/commit/0be2551616b2cbc9424a85902b7d1b0539b0ef12) docs: add companies using kube-ovn section
 * [d56552b89](https://github.com/kubeovn/kube-ovn/commit/d56552b896ef0060c96cb9711dda9cbfdce07677) docs: add community information
 * [8edd02255](https://github.com/kubeovn/kube-ovn/commit/8edd0225548c4c964d71ef0a60e5c7f8403a2792) fix: alleviate ping lost
 * [632bbc5e1](https://github.com/kubeovn/kube-ovn/commit/632bbc5e17565396043d615760845721c48285f5) refactor: refactor ovn-nbctl.go
 * [8aafa415e](https://github.com/kubeovn/kube-ovn/commit/8aafa415ef4a882363f7324b177149d8f5ba9f6a) docs: modify the readme
 * [60ce76592](https://github.com/kubeovn/kube-ovn/commit/60ce76592a20c890287d9613eaa0fb2d7772ddfb) fix: pinger percentage error
 * [276a28cf8](https://github.com/kubeovn/kube-ovn/commit/276a28cf8af18246d5953a062457469ce500b7e5) fix: add kube-ovn types to default scheme
 * [998a9e634](https://github.com/kubeovn/kube-ovn/commit/998a9e63449158b04c87defe6282ac4717747db4) refactor: cniserver
 * [a5d339b2b](https://github.com/kubeovn/kube-ovn/commit/a5d339b2b5132a0ce01dc416e0042204d7c7f239) docs: update docs
 * [dc92afa3d](https://github.com/kubeovn/kube-ovn/commit/dc92afa3dcfa9022da7c34526289bc65e8d354fd) fix: add a periodically recompute to ovn-controller to avoid inconsistency
 * [8488ae2a6](https://github.com/kubeovn/kube-ovn/commit/8488ae2a65f3361bd554025967294c3bb4faefa4) fix: add timeout to pinger access ovs/ovn
 * [ff1ff1457](https://github.com/kubeovn/kube-ovn/commit/ff1ff1457906480c7502b408707777df68676b83) fix: when subnet cidr conflict requeue the subnet
 * [e31a08ec3](https://github.com/kubeovn/kube-ovn/commit/e31a08ec3fa67d2ab6ea4281f8a952b57f4d2e3d) fix: add runGateway to wait.Until
 * [182390737](https://github.com/kubeovn/kube-ovn/commit/1823907373f12cb812bc2616cf68b2e4497ed8c8) fix: restart nbctl-daemon if not response
 * [839308e08](https://github.com/kubeovn/kube-ovn/commit/839308e080243fafa1a40d7a3161a3c75f4d402e) feat: display controller log in kubectl-ko diagnose
 * [8e6c3d62d](https://github.com/kubeovn/kube-ovn/commit/8e6c3d62d529e4c6e1c9b818b6c37f6f43629230) refactor: separate normal check and ovn specific check
 * [c97831817](https://github.com/kubeovn/kube-ovn/commit/c97831817d6fd962ef8e07c32fcbcb364d4fce2d) fix: do not return not found err
 * [f19e55964](https://github.com/kubeovn/kube-ovn/commit/f19e5596453b5ba97b97693d54add4f4fc4561e2) fix: move components to kube-system ns and add priorityClass
 * [a5d298dbf](https://github.com/kubeovn/kube-ovn/commit/a5d298dbfb45ae33059ce5e3d5d8ae9b633e0ee8) feat: cniserver check allocated annotation before configure pod network
 * [8f72b7ebb](https://github.com/kubeovn/kube-ovn/commit/8f72b7ebb7437a95759aa9b628145b8044c8e4d5) fix: set ovn-openflow-probe-interval
 * [3838a46d1](https://github.com/kubeovn/kube-ovn/commit/3838a46d146148bd8b9d7b7c0cc914f984a8ab84) pinger: add port binds check between local ovs and ovn-sb
 * [f8248cec4](https://github.com/kubeovn/kube-ovn/commit/f8248cec44caf653c0ade0537c14e5f6f4f34cd2) fix: if cidr block not ends with zero, reformat it
 * [dff1d6485](https://github.com/kubeovn/kube-ovn/commit/dff1d64857359919e1405f2dd11f8b3782e68fec) fix: resync iptables
 * [40fab55f2](https://github.com/kubeovn/kube-ovn/commit/40fab55f27081b87303dd45cf019195f19cac06d) update version
 * [920053c5d](https://github.com/kubeovn/kube-ovn/commit/920053c5dfddf85ea0d81f9c5d1d4303c488a9d5) pinger: add timeout for dns resolve
 * [513d2bd9f](https://github.com/kubeovn/kube-ovn/commit/513d2bd9f55a19446a65bb7c8aa14543c5c931be) e2e: add basic framework and tests for e2e
 * [9747d5405](https://github.com/kubeovn/kube-ovn/commit/9747d5405f8c8552a33f89a0f35ffc916af7d11d) docs: update changelog
 * [adf5071e1](https://github.com/kubeovn/kube-ovn/commit/adf5071e1dc58dae33fafee7913429e2b93e2cfe) fix: address in ep might be empty
 * [182bb1513](https://github.com/kubeovn/kube-ovn/commit/182bb1513017de14327fd681fa940f2965f8287c) fix: cniserver wait ovs ready
 * [518c0a781](https://github.com/kubeovn/kube-ovn/commit/518c0a7817e174a353d3d2bfd9f4d7f3dc633af0) fix: wrong deletion in gc lb and portgroup
 * [2492a1661](https://github.com/kubeovn/kube-ovn/commit/2492a166147d008dcaa4fbd4aa71267618c1cdfb) ovn: add memory patch to slow down memory increase
 * [d0bd71fd7](https://github.com/kubeovn/kube-ovn/commit/d0bd71fd7acde3e2315491d8491cc21875bb0656) fix: wait default and node logical switch ready
 * [23cad463b](https://github.com/kubeovn/kube-ovn/commit/23cad463b9f6fe07a8668dcc1008dba70666cc2b) fix: podSelector in networkpolicy should only consider pods in the same ns
 * [ca5539f09](https://github.com/kubeovn/kube-ovn/commit/ca5539f09cbff947bec69487e5e8e70b377b0400) fix: do not add unallocated pod to port-group
 * [d5ed1ee73](https://github.com/kubeovn/kube-ovn/commit/d5ed1ee73151a67612a87372850e9b7da01a4d66) release 0.10.0
 * [3c62ea29e](https://github.com/kubeovn/kube-ovn/commit/3c62ea29e5688a7a8d019e942040c05995e35126) ovn: pick up commit from upstream
 * [4c966c371](https://github.com/kubeovn/kube-ovn/commit/4c966c371713474edecec051151330761c716806) feat: pinger support check an address out of cluster.
 * [f00960789](https://github.com/kubeovn/kube-ovn/commit/f009607895f0a4610fa3ca72c16d7ea8d05a604c) chore: double quote shell variables
 * [83364b52e](https://github.com/kubeovn/kube-ovn/commit/83364b52e1a74d9cf3805620599ff581a7a8cbbe) fix: cluster mode db will generate lots listen error log
 * [d9e1cd1ca](https://github.com/kubeovn/kube-ovn/commit/d9e1cd1cad29de427bfdeea25aea71092c51a613) fix: gc logical_switch_port form listing pods and nodes
 * [a5dc8bb9b](https://github.com/kubeovn/kube-ovn/commit/a5dc8bb9b3a2b00e8e700cf75e923395268fa50f) fix: some init and cleanup bugs
 * [a5eb5e7f3](https://github.com/kubeovn/kube-ovn/commit/a5eb5e7f3bb259f53c50fba35b225bd160ee3694) fix: ovn-cluster mode
 * [a6f0dd143](https://github.com/kubeovn/kube-ovn/commit/a6f0dd143b8bc0b47b313b54247a810f3b13b929) feat: exclude_ips can be changed dynamically
 * [d9c594343](https://github.com/kubeovn/kube-ovn/commit/d9c594343a3727d2c15a2f199eb1b0624605bd7d) update ovn to 2.12.0-1
 * [06eceb3b3](https://github.com/kubeovn/kube-ovn/commit/06eceb3b31fc389248309746e64664e541f78012) feat: use label to select leader to avoid pod status misleading
 * [aa53c7dd4](https://github.com/kubeovn/kube-ovn/commit/aa53c7dd4122dd2e7b1b9d20bd07e6393f2b352e) fix: ip conflict when use ippool
 * [590443301](https://github.com/kubeovn/kube-ovn/commit/59044330170e47c0a3b6f60857de8ac91405fab3) docs: add v0.9.1 changelog
 * [5efbea9f3](https://github.com/kubeovn/kube-ovn/commit/5efbea9f31ccd6b47bc3b30b4e5f3a27df2602ea) fix: block subnet deletion when there any ip in use
 * [a1dc8c11f](https://github.com/kubeovn/kube-ovn/commit/a1dc8c11fb95b1c60d52ae7ecb39df26932d5a0f) plugin: kubectl plugin now expose ovs-vsctl to each node
 * [d3c6a71c5](https://github.com/kubeovn/kube-ovn/commit/d3c6a71c50caa98b6b667ca71b7004c5e0a9b266) fix: nbctl need timeout to avoid hang infinitely
 * [77e589034](https://github.com/kubeovn/kube-ovn/commit/77e58903458c1765a9541bc3e4415e5ee52877fe) perf: as lr-route-add with --may-exist will replace exist route, no need for another delete
 * [d4a51bdc6](https://github.com/kubeovn/kube-ovn/commit/d4a51bdc67b4dbafd6a4d3e91c71242d47407250) perf: when controller restart skip pod already create lsp
 * [7617fa797](https://github.com/kubeovn/kube-ovn/commit/7617fa79776ffa04a38fc7c8d8a607f3e44c0ff6) fix: when delete node recycle related ip/route resource
 * [f4e874769](https://github.com/kubeovn/kube-ovn/commit/f4e8747693c1220597616cbbb4653696c45864ed) fix typo in start-ovs.sh
 * [9b88e0841](https://github.com/kubeovn/kube-ovn/commit/9b88e08419471e186e29d25dde47e5a53c7ef068) perf: skip evicted pod when enqueueAddPod and enqueueUpdatePod
 * [e48186240](https://github.com/kubeovn/kube-ovn/commit/e481862402cd5af579ce44391a34c1b4f9af44b8) fix: use ep.subset.port.name to infer target port number
 * [0d8ae20c3](https://github.com/kubeovn/kube-ovn/commit/0d8ae20c3eb42dfaf115c8ec0b67ba7d720be5f0) fix: if no available address delete pod might failed related to #155
 * [bbd4257d0](https://github.com/kubeovn/kube-ovn/commit/bbd4257d0236e69d7f0b480e22161ca10ea4ae3c) kind: support reload kube-ovn component in kind cluster
 * [d0479e903](https://github.com/kubeovn/kube-ovn/commit/d0479e9035aa63d64864182176a8ee293872af02) perf: filter pod in informer list-watch and disable resync
 * [61a7a7b9c](https://github.com/kubeovn/kube-ovn/commit/61a7a7b9c12d8ff364366f94eecfdd9b8479e41b) fix: index out of range err when create lsp
 * [623661ef2](https://github.com/kubeovn/kube-ovn/commit/623661ef28db2444f1ef6a4242357e59c251e449) prepare for next release
 * [1643c7f0d](https://github.com/kubeovn/kube-ovn/commit/1643c7f0da3c472e5265584960ad379adfe9f8bd) kind: support to install kube-ovn in kind
 * [9611599f1](https://github.com/kubeovn/kube-ovn/commit/9611599f1d447f911d20b6fad5b21a28de960bee) fix: mount /var/run/netns that kind will use it to store network ns files
 * [53db261a1](https://github.com/kubeovn/kube-ovn/commit/53db261a12444f89b3520a1b2528c7f08063ce6b) release: v0.9.0
 * [1984cbe80](https://github.com/kubeovn/kube-ovn/commit/1984cbe8010bc52d49989caa336076b2675110e3) feat: when use nodelocaldns do not nat the address
 * [446999f4c](https://github.com/kubeovn/kube-ovn/commit/446999f4cd022e4540f0239a346bd17178f4f9af) docs: add description about relation of cidr and static ip allocation
 * [6f1854f9f](https://github.com/kubeovn/kube-ovn/commit/6f1854f9f58cc94cfc74b71eebf89f2770e0509b) Check the short name of kubernetes services which is independant of the cluster domain name.
 * [c6f8efeb0](https://github.com/kubeovn/kube-ovn/commit/c6f8efeb059208c766fe3585558c8f8c38439b58) fix: some grafana modification
 * [40144160d](https://github.com/kubeovn/kube-ovn/commit/40144160d2926628d27a7766d3332dd7569342f9) fix: add missing cap
 * [7c464d692](https://github.com/kubeovn/kube-ovn/commit/7c464d692f7e3a80c47aa642034645b0ea70e28d) chore: update ovn and other minor fix
 * [ac5371521](https://github.com/kubeovn/kube-ovn/commit/ac53715217327651264af87b93e627c9908d39a9) fix re-annotate namespaces when subnet deleted
 * [fe2f2612a](https://github.com/kubeovn/kube-ovn/commit/fe2f2612a60164b1bbe049a4c5dbdcef984fc6d6) fix: add ingress_policing_burst to accurate limit ingress bandwidth
 * [20b2c83de](https://github.com/kubeovn/kube-ovn/commit/20b2c83deab57801aaf2fce8b475463cefc991d2) fix: network unreachable when add egress qos for pod
 * [758dbc1ce](https://github.com/kubeovn/kube-ovn/commit/758dbc1ce7c14b93bbf613f7385ea2e7f4bbb7c5) fix: err when add egress qos
 * [bdfd351d7](https://github.com/kubeovn/kube-ovn/commit/bdfd351d738f849f26b186c398da11c9453080bc) fix: remove privilege=true from long run container
 * [0859da1fe](https://github.com/kubeovn/kube-ovn/commit/0859da1fe5a2d0195a9ebd81a2c723e4b3be6d98) perf: optimize pod add
 * [3718851df](https://github.com/kubeovn/kube-ovn/commit/3718851df4087e92cdc2d68d2dd4592413c61742) fix: add keepalive to ovn-controller
 * [6ad981064](https://github.com/kubeovn/kube-ovn/commit/6ad981064336a19c441c4eceff2153066ba53273) feat: add controller metrics
 * [b87ed0eec](https://github.com/kubeovn/kube-ovn/commit/b87ed0eec24dde544770635e318c710a81adb0c7) If pod have not a status.PodIP skip add/del static route
 * [b9108fba4](https://github.com/kubeovn/kube-ovn/commit/b9108fba4d0d6f26f7940d55f8c7236c4f5f3858) fix: ippool pod static route might lost during leader election
 * [a2e24de62](https://github.com/kubeovn/kube-ovn/commit/a2e24de62a553e791506a7b8679f72f94e161dea) fix: static route might lost during leader election
 * [8202a1886](https://github.com/kubeovn/kube-ovn/commit/8202a188663e674eced3348ec1dac90d423930c2) feat: add grafana config and modify metrics.
 * [cae0ef274](https://github.com/kubeovn/kube-ovn/commit/cae0ef274857ada317cfe5092278c6f1b53e675c) fix: only keep the last iface-id
 * [f3528f237](https://github.com/kubeovn/kube-ovn/commit/f3528f237bcfc07b3b459bdd0a47a15077dd66d2) fix: add missing gc
 * [3791ba29d](https://github.com/kubeovn/kube-ovn/commit/3791ba29d0adf66bd4511a3e35efc207cff2659c) fix: gc resource when start controller
 * [f970615bb](https://github.com/kubeovn/kube-ovn/commit/f970615bb1a583b10f506a2f70510f45e293a08f) fix: watch will break if timeout is set
 * [ef285b213](https://github.com/kubeovn/kube-ovn/commit/ef285b213de5086980dc70eeb62542370e8e4427) feat: pinger add apiserver check metrics
 * [d33685e64](https://github.com/kubeovn/kube-ovn/commit/d33685e6403aea8838ab4f6ca2ccd04dadd9002e) fix: avoid conflict when init
 * [6b57f61b3](https://github.com/kubeovn/kube-ovn/commit/6b57f61b34c18ee5815db7e595bca327967d8f68) release v0.8.0
 * [6ed722f94](https://github.com/kubeovn/kube-ovn/commit/6ed722f9439bd5987f232b6ab584b5c61ceefc33) fix: loss might be negative number
 * [7c0517b51](https://github.com/kubeovn/kube-ovn/commit/7c0517b5120770ff72ff35b45376457ae85084e4) feat: pinger prometheus support
 * [e23bd552f](https://github.com/kubeovn/kube-ovn/commit/e23bd552f46bdb76065f869d4130533c670c55a9) feat: support pinger
 * [d837aa122](https://github.com/kubeovn/kube-ovn/commit/d837aa122b9d37728e3db349b1e5e9e1f9580c03) chore: update ovs/ovn
 * [4246cb74d](https://github.com/kubeovn/kube-ovn/commit/4246cb74de513fb1267b512e19eedcb2e8b969e4) feat: gateway ha
 * [e27c9e54d](https://github.com/kubeovn/kube-ovn/commit/e27c9e54df996c82319f91f8b255706b2ae8781c) chore: remove ovs-ipsec and update go to 1.13
 * [ba3084ebc](https://github.com/kubeovn/kube-ovn/commit/ba3084ebc827b3fba3dbac367f92f705151445db) feat: add kubectl plugin
 * [54a465d16](https://github.com/kubeovn/kube-ovn/commit/54a465d161f9351c49e103261c716ba2b3568e68) docs: add comparison
 * [38be68d6f](https://github.com/kubeovn/kube-ovn/commit/38be68d6f2a8c50c6116413a3f95d668de43f42e) fix: pod should be accessed from node when acl applied
 * [e62f0ab04](https://github.com/kubeovn/kube-ovn/commit/e62f0ab0412561d8607aca7843b899d2a1aac514) enable portmap by default to support hostport
 * [80de8e58d](https://github.com/kubeovn/kube-ovn/commit/80de8e58d53a5e98fdbce89030637535e63b1ea8) feat: add port security to pod port
 * [4849f0562](https://github.com/kubeovn/kube-ovn/commit/4849f05621dac0220d6db10c8bd3f52907035ee7) feat: add node switch allocated ip cr
 * [34e8406ee](https://github.com/kubeovn/kube-ovn/commit/34e8406ee17ef7d28229887e9fd6557e782cdee8) prepare for next release
 * [933fd8d25](https://github.com/kubeovn/kube-ovn/commit/933fd8d2505cc813f0db18f05d4535ac7d530b15) release: bump v0.7.0
 * [7e2bdf522](https://github.com/kubeovn/kube-ovn/commit/7e2bdf522015a7678f7a05b0f6b7cbba3af5c72d) fix: add default excludeIps and check kern version
 * [31544abb5](https://github.com/kubeovn/kube-ovn/commit/31544abb510f615172dc19eaeb96e342b80de222) fix: deal with ipv6 connection str
 * [0f8f2aad7](https://github.com/kubeovn/kube-ovn/commit/0f8f2aad7492fa295f680faeb4c65e14b5ed8a2a) fix missing condition when subnet is private
 * [d37da1bc4](https://github.com/kubeovn/kube-ovn/commit/d37da1bc48c130221c7f3631a7e7e2d8b4b25948) add subnet status
 * [4a5c54983](https://github.com/kubeovn/kube-ovn/commit/4a5c5498345af803e6ebd16d7bed648134855be1) fix: acl related issues
 * [62a395e6f](https://github.com/kubeovn/kube-ovn/commit/62a395e6f242c96f4af81be6f93781bcb7508326) Revert "add subnet status field"
 * [b8f1d9ef0](https://github.com/kubeovn/kube-ovn/commit/b8f1d9ef0fef618ee14bb33dd18bce560ac084a5) add missing subnets/status operation permission
 * [6c119ad1e](https://github.com/kubeovn/kube-ovn/commit/6c119ad1e2d784ab5b6397e722c8302867b29d38) Update cleanup.sh
 * [b08ece4fb](https://github.com/kubeovn/kube-ovn/commit/b08ece4fbca124d97ed83f9c6ac195ae84a4f64e) feat: add exclude_ips annotation to namespace
 * [a2774ed0e](https://github.com/kubeovn/kube-ovn/commit/a2774ed0e44a48e17efa0ea0bbb5b32277ea5430) fix: use pg-del to remove pg and acl, check if ports is empty before set pg
 * [422c6dc00](https://github.com/kubeovn/kube-ovn/commit/422c6dc0001108bb4bf394baf8b903ee1baa30ea) add subnet status
 * [fde683eaf](https://github.com/kubeovn/kube-ovn/commit/fde683eaf8e8574f6d53867478b955310a248594) feat: add subnet annotation to ns and automatically unbind ns from subnet.
 * [948e13063](https://github.com/kubeovn/kube-ovn/commit/948e130638f8870b7919d58229d1e6efd18d7c5d) docs: add cn docs link
 * [5278e1051](https://github.com/kubeovn/kube-ovn/commit/5278e1051139d986e78f04c32a1a95074787f61a) feat: add default values to subnet
 * [ea451a1aa](https://github.com/kubeovn/kube-ovn/commit/ea451a1aace30bbf35b93ddc14978c5dcab32322) write back subnet name to ip label
 * [1c7121dbb](https://github.com/kubeovn/kube-ovn/commit/1c7121dbb9ade6aa0edb362050223cbed2ad8b6d) chore: enable mirror in yaml and modify docs
 * [db9783a3a](https://github.com/kubeovn/kube-ovn/commit/db9783a3a5d045814250110d9d111a54277a622e) fix: duplicate import in network_policy.go
 * [8a57747ef](https://github.com/kubeovn/kube-ovn/commit/8a57747ef794ea6fb9c570992e0d023e01bc3a76) fix: improve cni-conf name priority
 * [5f1436bef](https://github.com/kubeovn/kube-ovn/commit/5f1436bef28b59fcc8a6a2b1b66b6b0a746633c2) fix: wait subnet ready before start worker.
 * [661387efe](https://github.com/kubeovn/kube-ovn/commit/661387efea058278baf503add43d53ef4cc4d03d) fix: check ls exists before handle it
 * [9e05f5333](https://github.com/kubeovn/kube-ovn/commit/9e05f53333afd65c6cb122200a3a6d21b1aa8c6d) docs: add more installation tools.
 * [dccb93c7b](https://github.com/kubeovn/kube-ovn/commit/dccb93c7b347613dcae28c59a347651d2188582d) docs: add support os and notes.
 * [c6a160b3f](https://github.com/kubeovn/kube-ovn/commit/c6a160b3f09df6a20fd48abe9a5ec338fa5792ef) Update subnet.md
 * [31ad00bd7](https://github.com/kubeovn/kube-ovn/commit/31ad00bd75a3e1fe4aaa4f7993b56e9ec1b94163) feat: add ip info to ip crd
 * [ad7b5c2f6](https://github.com/kubeovn/kube-ovn/commit/ad7b5c2f67d885c8e7ad6c1067d8b5ddc3d663cf) feat: update logo
 * [44c3077c7](https://github.com/kubeovn/kube-ovn/commit/44c3077c7b0153944861acf2c20b1363819cb39c) feat: add logo
 * [55d7fd6f7](https://github.com/kubeovn/kube-ovn/commit/55d7fd6f791357d666cca3bc24d11e3e96485acf) feat: reserve vport for statefulset pod
 * [7a3c8a6a8](https://github.com/kubeovn/kube-ovn/commit/7a3c8a6a8f7be35af6ab07d0363d2484dd4462d5) docs: add crd installation
 * [aa016c1b0](https://github.com/kubeovn/kube-ovn/commit/aa016c1b0220f51a3a273acc44bd5c9d48aadde2) fix: modify default header length
 * [85b406907](https://github.com/kubeovn/kube-ovn/commit/85b4069077dcc52600ccdca26bda24e9b5a723f8) fix: do not create exist logical switch
 * [362943660](https://github.com/kubeovn/kube-ovn/commit/3629436602367fcd8e4077a292859da59eb335a5) chore: prepare for next release
 * [463d6253a](https://github.com/kubeovn/kube-ovn/commit/463d6253a4e5362496d2b9fd5f90d7c39ecf87e4) docs: add crd/ipv6 docs and bump version 0.6.0
 * [103c23af6](https://github.com/kubeovn/kube-ovn/commit/103c23af64b08998768770b512b894e609456528) fix build error
 * [9d173ba0c](https://github.com/kubeovn/kube-ovn/commit/9d173ba0c28dc84b322e4ea729edf6d7972f9100) feat: support ipv6-only mode
 * [05566017c](https://github.com/kubeovn/kube-ovn/commit/05566017cb41850fe1958c0ac392017f6f4bbeab) add webhook docs
 * [766cec9b8](https://github.com/kubeovn/kube-ovn/commit/766cec9b8435ae3a2b2f37f1fdd0056c9c582ebf) add admission webhook for static ip
 * [2abeacb4a](https://github.com/kubeovn/kube-ovn/commit/2abeacb4a85b7db613dcfe96106ddda3a57b059b) docs: add support platform version
 * [ed7264ea9](https://github.com/kubeovn/kube-ovn/commit/ed7264ea91cafe6201ad48559fc9ac7dfed46dea) feat: use subnet crd to manage logical switch
 * [1e5c9f6cb](https://github.com/kubeovn/kube-ovn/commit/1e5c9f6cbdc5b120da2a528e99883796281fd5a6) Use k8s hostname, fix #60
 * [873672952](https://github.com/kubeovn/kube-ovn/commit/873672952e945cfc90e9e1cba0630d6551cab644) fix: remove dependency on cluster-admin
 * [e0864a03d](https://github.com/kubeovn/kube-ovn/commit/e0864a03d012a311e5ff3ae854d86a7cafabc790) chore: use go mod to replace dep
 * [96ec620d9](https://github.com/kubeovn/kube-ovn/commit/96ec620d924568a67d482b003cba61627930c64a) docs: update mirror feature to readme
 * [855d834fa](https://github.com/kubeovn/kube-ovn/commit/855d834fa9721ebad93ed26d3fdfb88d8a0a8c39) feat: support traffic mirror
 * [d1c3ea853](https://github.com/kubeovn/kube-ovn/commit/d1c3ea853f935414eb0f34baadce8b377996ff62) prepare for next release
 * [782e04bee](https://github.com/kubeovn/kube-ovn/commit/782e04bee9922016c2332602145f45cd186fc7c3) chore: bump v0.5.0
 * [a27f83394](https://github.com/kubeovn/kube-ovn/commit/a27f8339456a5264611af4dbbe5482e42fdd1dc1) fix: wrong mtu
 * [447071676](https://github.com/kubeovn/kube-ovn/commit/4470716766a6aa3da3ac6453aae2992846301c58) feat: support user define iface and mtu
 * [f8d8e186f](https://github.com/kubeovn/kube-ovn/commit/f8d8e186fee449dc7f1549a046566719c19c8047) fix: remove mask field from ip annotation
 * [550904047](https://github.com/kubeovn/kube-ovn/commit/5509040475fe64297c5e68a7e8f45b4dcc61be9e) feat: auto assign gw for controller config and expose more cmd args
 * [48da0fe17](https://github.com/kubeovn/kube-ovn/commit/48da0fe177dab436450650d1e346431be349bd03) feat: add pprof and use it as probe
 * [8984c90b5](https://github.com/kubeovn/kube-ovn/commit/8984c90b567750195b72b694d1a230bfdceb1c8a) feat: set kernel args when start cniserver
 * [208a1dfc7](https://github.com/kubeovn/kube-ovn/commit/208a1dfc797552d24bc385e2873198e8ae9f851f) feat: support network policy
 * [c8d208fbd](https://github.com/kubeovn/kube-ovn/commit/c8d208fbdfdf00e18653f9cba735b35cd03f74fd) prepare for next release
 * [5a2cb093c](https://github.com/kubeovn/kube-ovn/commit/5a2cb093c6eddb1c450b74f1c942c153341a40bb) bump version to v0.4.1
 * [f8e8b0011](https://github.com/kubeovn/kube-ovn/commit/f8e8b00114af09ca62a8e00eb28623640a73ec87) fix: manual static ip allocation and automatic allocation should use different ip validation
 * [031924d10](https://github.com/kubeovn/kube-ovn/commit/031924d10dbcda61deb217fa4f3efe1b13748f6b) Fix json: cannot unmarshal string into Go value of type request.PodResponse https://github.com/alauda/kube-ovn/issues/33
 * [24259dbf9](https://github.com/kubeovn/kube-ovn/commit/24259dbf95af068990e0d1ebf52971a0d928c0d4) fix: use ovsdb-client to get leader info
 * [3541b6cf7](https://github.com/kubeovn/kube-ovn/commit/3541b6cf747dc33d7667acb467e5a1a0f86c8dbc) fix: use default-gw as default-exclude-ips and expose args to docs
 * [69c485380](https://github.com/kubeovn/kube-ovn/commit/69c48538074d0e002c444207fbd68995f63b8293) to cleanup all created resources, not only kube-ovn namespace.
 * [9361bb438](https://github.com/kubeovn/kube-ovn/commit/9361bb4383ccbea67c63858150366c439fb92360) prepare for next release
 * [509bf4a41](https://github.com/kubeovn/kube-ovn/commit/509bf4a4168f95e99e86ce8d4375967b50b66a64) feat: bump version to 0.4.0
 * [2e4145197](https://github.com/kubeovn/kube-ovn/commit/2e414519755441d90e0947a4337bd4b72fbf29d1) feat: support expose pod ip to external network
 * [8992bbe38](https://github.com/kubeovn/kube-ovn/commit/8992bbe3833e1f9ac80b0aa95fafcd7d173d0000) fix: check conflict subnet cidr
 * [0f9d1e4be](https://github.com/kubeovn/kube-ovn/commit/0f9d1e4be28947269394d44747f7bea06ece8911) fix: start informer when controller is leader
 * [71c15d65b](https://github.com/kubeovn/kube-ovn/commit/71c15d65b0ed3ce09f0296fa897c8a4c1957445d) feat: validate namespace/pod annotations
 * [89491b570](https://github.com/kubeovn/kube-ovn/commit/89491b5702ac3dbffca10c7b7e38e0fe2165c48a) fix: wait node-gw info ready
 * [0d86393d2](https://github.com/kubeovn/kube-ovn/commit/0d86393d298725de08b28401cbf094811eee0e93) fix: use ovn/ovs-ctl to health check
 * [278ccfe5e](https://github.com/kubeovn/kube-ovn/commit/278ccfe5e94e208fa2f8e812c01ea4a9fa992070) feat: remove finalizer dependency improve svc performance
 * [8f962673a](https://github.com/kubeovn/kube-ovn/commit/8f962673a7ee83bdb9706ebc2548885f4661aca5) fix: reuse node ip and mac annotation
 * [b8f85143b](https://github.com/kubeovn/kube-ovn/commit/b8f85143b6571462845fc4487966ea33cbfbcbe6) Add ha for ovn dbs and simplify makefile
 * [3c617451e](https://github.com/kubeovn/kube-ovn/commit/3c617451efe78b46acddcf012e00d049bf722171) feat: merge ovn-nbctl request
 * [b5ac7da4a](https://github.com/kubeovn/kube-ovn/commit/b5ac7da4a5c9a1f50a1065440abb0660440d2aaf) feat: separate ip pool pod and add parallelism to workers
 * [ce105dffb](https://github.com/kubeovn/kube-ovn/commit/ce105dffb1a6634bc414ed33030ee61ef507c106) Mute logrus log for ipset Dont need to change the vendored code.
 * [657470c8f](https://github.com/kubeovn/kube-ovn/commit/657470c8f27de71451c81639e64d148b56354829) Fix klog cant use V module The side affect of this commit is glog's V module not work.
 * [5429f51be](https://github.com/kubeovn/kube-ovn/commit/5429f51bea2157f321ada9400463df85db41f342) feat: use ovn macam to allocate mac for static ip pod
 * [5a8958cdb](https://github.com/kubeovn/kube-ovn/commit/5a8958cdb18d4907668db16191e98abc1d9b1cba) feat: update ovn to 2.11.1
 * [ca036f9e7](https://github.com/kubeovn/kube-ovn/commit/ca036f9e7eb9a65b3641fb98c6c744454880897a) Add vagrantfile
 * [660c0570c](https://github.com/kubeovn/kube-ovn/commit/660c0570cf2130b9bddb2d5706af5fc4caad3b41) fix: use tag version yaml url
 * [bc66671c5](https://github.com/kubeovn/kube-ovn/commit/bc66671c511f64cdd8c9d76588f8586464434b89) chore: fix go-report golint issues
 * [12a4bec93](https://github.com/kubeovn/kube-ovn/commit/12a4bec93e7b655e439af8e2903e895d62a7b047) ha for kube-ovn-controller
 * [b7d0f599e](https://github.com/kubeovn/kube-ovn/commit/b7d0f599e153f2c39cc96af59425c57db3e62ad1) cleanup unused code
 * [756831d72](https://github.com/kubeovn/kube-ovn/commit/756831d7213846d8c9676775a02a6d505c67dfd5) docs: add network topology
 * [c05594873](https://github.com/kubeovn/kube-ovn/commit/c055948732089450464c747842c5914703762d1b) chore: Minor updates to gateway.md
 * [21e34e9f2](https://github.com/kubeovn/kube-ovn/commit/21e34e9f217a64ef7bd8c078150d5c2e199542e8) chore: Gateway documentation touch-ups
 * [aa0b2b7c1](https://github.com/kubeovn/kube-ovn/commit/aa0b2b7c147c3108c0e1315870a83393efcb6d7a) chore: QoS documentation touch-ups
 * [3ec0098a4](https://github.com/kubeovn/kube-ovn/commit/3ec0098a4e0ab10b6924a58bcf18d04538b29ef1) chore: Subnet Isolation documentation touch-ups
 * [524845e93](https://github.com/kubeovn/kube-ovn/commit/524845e93c7197e3b61e50949012e5ab4dd30e14) chore: Static IP documentation touch-up
 * [b510016c2](https://github.com/kubeovn/kube-ovn/commit/b510016c2916b91b16930a2c36a6e41122b2c0a1) chore: Subnet documentation touch-ups
 * [524f7d3ff](https://github.com/kubeovn/kube-ovn/commit/524f7d3ff3ce02b4f6261c4e33b0137a935b4431) chore: Installation Guide touch-ups
 * [a1995d037](https://github.com/kubeovn/kube-ovn/commit/a1995d0374f47e8ab932fa09ed83fc690ec6bcb0) chore: README touch-up.
 * [79c0642eb](https://github.com/kubeovn/kube-ovn/commit/79c0642eba7b8146269262ce2d011f1c725ef1a3) docs: bump version
 * [cb2f50da4](https://github.com/kubeovn/kube-ovn/commit/cb2f50da4639aa02e6bde5acc26874d3a08a47ed) fix: acl rule error
 * [1a6f492ad](https://github.com/kubeovn/kube-ovn/commit/1a6f492ad603cdc1ac34992fb502f790d86c9fdb) fix: init node gw before run controller
 * [75c514a19](https://github.com/kubeovn/kube-ovn/commit/75c514a19a633c72068d094eae7fc31d25a5ec25) fix: external dns issues
 * [130688927](https://github.com/kubeovn/kube-ovn/commit/130688927a4282bb8c9c662f09a74bbefca77ef2) feat: use daemon ovn-nbctl to improve performance and cleanup unused dns code
 * [24cda418c](https://github.com/kubeovn/kube-ovn/commit/24cda418cd50af8173bff0e3698ddf7780bb1c53) Implement centralized gateway.
 * [890934f47](https://github.com/kubeovn/kube-ovn/commit/890934f473772c26453c49bca93a1fcb57dbc962) chore: migrate from bitbucket to github
 * [adf655cb9](https://github.com/kubeovn/kube-ovn/commit/adf655cb95e57f4b9e9921dd1074dfa517a88fc5) remove dns from ls and bump new version
 * [ca21c6cb1](https://github.com/kubeovn/kube-ovn/commit/ca21c6cb1f50fb3b313f1104cadc5c0baf7deb73) make filter table forward chain default accept
 * [cd0ddf10a](https://github.com/kubeovn/kube-ovn/commit/cd0ddf10a39cf3e95bcff11dd9b0d37075968891) ipset exclude cluster service ip range
 * [1d753c8e1](https://github.com/kubeovn/kube-ovn/commit/1d753c8e1d189bb9a05500365be53618207b0e1c) fix: lb bugs
 * [cb91d9842](https://github.com/kubeovn/kube-ovn/commit/cb91d9842e9cf5fdfa9817a4da94c868cf786a6c) read cidr from ns annotation
 * [e99983329](https://github.com/kubeovn/kube-ovn/commit/e9998332991445ed4a9664071af1601fc6b57af6) fix: remove dns table from nodeswitch and remove unused other_config:namespace
 * [049cab2c1](https://github.com/kubeovn/kube-ovn/commit/049cab2c114f3fdca4b0b4b18fcd1a3515c2d855) fix pod has no ip
 * [170c3c634](https://github.com/kubeovn/kube-ovn/commit/170c3c63433cb75c8d522dcd7256e09bf059dcd4) Distributed gateway implement
 * [cebb8dfda](https://github.com/kubeovn/kube-ovn/commit/cebb8dfda9f7ee72a3b66d92066a60c9fe0a10ac) fix: clean lost interface.
 * [4367ba079](https://github.com/kubeovn/kube-ovn/commit/4367ba079c268d163b08c282af510ccc4ee0beb0) feat: support subnet isolation
 * [1fe8c916d](https://github.com/kubeovn/kube-ovn/commit/1fe8c916dcb6f678fa8dcde28ff570ce63a456c5) feat: support dynamic qos
 * [e04bc0939](https://github.com/kubeovn/kube-ovn/commit/e04bc09397a922e60e5b13be7192f80816624d98) fix: ovn restart issues
 * [014f1dcf5](https://github.com/kubeovn/kube-ovn/commit/014f1dcf531e8eb6be43a9d6d28e6cc952111f38) fix: ovn restart issues
 * [3e78ddc37](https://github.com/kubeovn/kube-ovn/commit/3e78ddc375a51a96bfb711a96d70912c10dafd60) fix: validate namespace switch annotations
 * [44eafc505](https://github.com/kubeovn/kube-ovn/commit/44eafc505a05bfebea49a7eaf8f09f98f4c7a885) fix lint && add docker build
 * [cb3e01a4e](https://github.com/kubeovn/kube-ovn/commit/cb3e01a4ef7c599ac78e40e802fd4a7346001dba) feat: update yaml, add readiness/liveness probe, add pass shell args
 * [004deefd9](https://github.com/kubeovn/kube-ovn/commit/004deefd9dd061e73d9a54ee2721afab4ee8ecf2) feat: support qos
 * [d37264e45](https://github.com/kubeovn/kube-ovn/commit/d37264e4547ef51e3bd533aba33525a2121aa0e6) feat: add simple gateway implementation

### Contributors

 * Amye Scavarda Perrin
 * Bruce Ma
 * Emma Kenny
 * Gary
 * Haocheng Liu
 * JinLin Fu
 * Kai Chen
 * Mengxin Liu
 * MengxinLiu
 * Patryk Strusiewicz-Surmacki
 * QIANSHUANGYANG [钱双洋]
 * Sébastien BERNARD
 * Wan Junjie
 * Xiang Dai
 * Yan Wei
 * Yan Zhu
 * Your Name
 * caoyingjun
 * chestack
 * ckji
 * cmj
 * danieldin95
 * emmakenny
 * fanbin
 * fangtian
 * feixiang
 * fossabot
 * ftiannew
 * halfcrazy
 * hzma
 * laik
 * linruichao
 * luoyunhe1
 * lut777
 * pengbinbin1
 * qsyqian
 * sayicui
 * shuangyang.qian
 * wangyudong
 * withlin
 * wiwen
 * xieyanker
 * zhangzujian
 * 范日明
 * 马洪贞

