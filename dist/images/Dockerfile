# syntax=docker/dockerfile:1
ARG VERSION=latest
ARG BASE_TAG=$VERSION

# Intermediate stage: download gobgp CLI binary.
# Keeps the base image untouched while upgrading third-party binaries.
FROM ubuntu:24.04 AS qiniu
ARG ARCH=amd64
ARG DEBIAN_FRONTEND=noninteractive
# renovate: datasource=github-releases depName=gobgp packageName=osrg/gobgp versioning=semver
ARG GOBGP_VERSION=4.7.0
ARG GOBGP_SHA256_AMD64=05d98ca0d7bbcb2f50a6b7b6ee51c5e5b5fd64d6a310ee807040ed9d7104d5e0
ARG GOBGP_SHA256_ARM64=9e582f0d5914ad706621d17047ce247879939aa9a359cb9f5a7446244adfec98
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates && \
    curl -sSf -L --retry 5 -o /tmp/gobgp.tar.gz \
    "https://github.com/osrg/gobgp/releases/download/v${GOBGP_VERSION}/gobgp_${GOBGP_VERSION}_linux_${ARCH}.tar.gz" && \
    case "$ARCH" in \
    amd64) EXPECTED_SHA256="$GOBGP_SHA256_AMD64" ;; \
    arm64) EXPECTED_SHA256="$GOBGP_SHA256_ARM64" ;; \
    *) echo "unsupported arch: $ARCH"; exit 1 ;; \
    esac && \
    echo "${EXPECTED_SHA256}  /tmp/gobgp.tar.gz" | sha256sum -c - && \
    tar -xzf /tmp/gobgp.tar.gz -C /usr/bin gobgp && \
    chmod +x /usr/bin/gobgp && \
    rm -f /tmp/gobgp.tar.gz && \
    apt-get purge -y --auto-remove curl && \
    rm -rf /var/lib/apt/lists/*

FROM ubuntu:24.04 AS eip-bpf
RUN apt-get update && apt-get install -y --no-install-recommends clang libbpf-dev linux-libc-dev make && \
    rm -rf /var/lib/apt/lists/*
COPY bpf /src/bpf
RUN make -C /src/bpf clean all

FROM kubeovn/kube-ovn-base:$BASE_TAG AS setcap

COPY *.sh /kube-ovn/
COPY kubectl-ko /kube-ovn/kubectl-ko
COPY 01-kube-ovn.conflist /kube-ovn/01-kube-ovn.conflist

COPY kube-ovn /kube-ovn/kube-ovn
COPY kube-ovn-cmd /kube-ovn/kube-ovn-cmd
COPY kube-ovn-daemon /kube-ovn/kube-ovn-daemon
COPY kube-ovn-controller /kube-ovn/kube-ovn-controller
RUN ln -s /kube-ovn/kube-ovn-cmd /kube-ovn/kube-ovn-monitor && \
    ln -s /kube-ovn/kube-ovn-cmd /kube-ovn/kube-ovn-speaker && \
    ln -s /kube-ovn/kube-ovn-cmd /kube-ovn/kube-ovn-webhook && \
    ln -s /kube-ovn/kube-ovn-cmd /kube-ovn/kube-ovn-leader-checker && \
    ln -s /kube-ovn/kube-ovn-cmd /kube-ovn/kube-ovn-ic-controller && \
    ln -s /kube-ovn/kube-ovn-controller /kube-ovn/kube-ovn-pinger && \
    setcap CAP_NET_BIND_SERVICE+eip /kube-ovn/kube-ovn-cmd && \
    setcap CAP_NET_RAW,CAP_NET_BIND_SERVICE+eip /kube-ovn/kube-ovn-controller && \
    setcap CAP_NET_ADMIN,CAP_NET_RAW,CAP_NET_BIND_SERVICE,CAP_SYS_ADMIN+eip /kube-ovn/kube-ovn-daemon

FROM kubeovn/kube-ovn-base:$BASE_TAG

COPY --chmod=0644 logrotate/* /etc/logrotate.d/
COPY grace_stop_ovn_controller /usr/share/ovn/scripts/grace_stop_ovn_controller

COPY --from=setcap /kube-ovn /kube-ovn
COPY --from=qiniu /usr/bin/gobgp /usr/bin/gobgp
COPY --from=eip-bpf /src/bpf/tc_eip.o /etc/kube-ovn/tc_eip.o
RUN /kube-ovn/iptables-wrapper-installer.sh --no-sanity-check

WORKDIR /kube-ovn
