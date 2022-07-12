FROM registry.access.redhat.com/ubi9/ubi as builder

ARG DNF_FLAGS="\
  -y \
  --nodocs \
  --releasever 9 \
  --setopt install_weak_deps=false \
  --installroot \
"
ARG DNF_PACKAGES="\
  openssl \
  coreutils-single \
  glibc-minimal-langpack \
"
ARG ROOTFS="/rootfs"
ARG TARGETARCH

RUN set -ex \
     && mkdir -p ${ROOTFS} \
     && dnf install ${DNF_FLAGS} ${ROOTFS} ${DNF_PACKAGES} \
     && dnf clean all ${DNF_FLAGS} ${ROOTFS} \
     && rm -rf ${ROOTFS}/var/cache/* \
    && echo

FROM scratch
COPY --from=builder /rootfs/ /
COPY ./client-linux-${TARGETARCH} /usr/local/bin/client
ENTRYPOINT ["/usr/local/bin/client"]
CMD ["version"]
