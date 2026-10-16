FROM alpine:edge

RUN set -ex \
    && apk update \
    && apk upgrade \
    && apk add --no-cache \
    bash \
    dumb-init \
    iproute2 \
    iptables iptables-legacy \
    nftables \
    iputils \
    tcpdump \
    netcat-openbsd \
    conntrack-tools

WORKDIR /kube-ovn
COPY --chmod=0755 nat-gateway.sh /kube-ovn/
COPY --chmod=0755 lb-svc.sh /kube-ovn/
RUN chmod +x /kube-ovn/*.sh

# Container startup command. The controller does not override this,
# keeping startup behavior decoupled and managed by the image.
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sleep", "infinity"]
