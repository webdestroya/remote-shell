FROM ubuntu:bionic AS builder

RUN set -eux; \
  \
  savedAptMark="$(apt-mark showmanual)"; \
  apt-get update; \
  apt-get install -y --no-install-recommends \
    autoconf \
    automake \
    bzip2 \
    ca-certificates \
    g++ \
    gcc \
    make \
    patch \
    unzip \
    wget \
    xz-utils \
    zlib1g-dev \
    \
  ; \
  rm -rf /var/lib/apt/lists/*; \
#RUN set -eux; \
  \
  mkdir -p /cloud87; \
  \
  cd /; \
  wget -O jq.tar.gz "https://github.com/stedolan/jq/releases/download/jq-1.6/jq-1.6.tar.gz"; \
  mkdir -p /usr/src/jq; \
  tar -xzf jq.tar.gz -C /usr/src/jq --strip-components=1; \
  rm jq.tar.gz; \
  \
  cd /usr/src/jq; \
  \
  autoconf; \
  ./configure \
    --without-oniguruma \
    --enable-all-static \
  ; \
  make -j "$(nproc)" LDFLAGS=-all-static; \
  \
  mv /usr/src/jq/jq /jq; \
  \
  cd /; \
  \
  rm -rf /usr/src/jq; \
  \
  /jq --version; \
  \
# BUILDING DROPBEAR
  cd /; \
  wget -O dropbear.tar.bz2 "https://matt.ucc.asn.au/dropbear/releases/dropbear-2022.82.tar.bz2"; \
  mkdir -p /usr/src/dropbear; \
  tar -xjf dropbear.tar.bz2 -C /usr/src/dropbear --strip-components=1; \
  rm dropbear.tar.bz2; \
  cd /usr/src/dropbear; \
  \
  { \
    echo '#define DROPBEAR_DEFPORT "8722"'; \
    echo '#define DISABLE_SYSLOG 1'; \
    echo '#define DROPBEAR_SVR_PASSWORD_AUTH 0'; \
    echo '#define DROPBEAR_SFTPSERVER 0'; \
    echo '#define DROPBEAR_DSS 0'; \
    echo '#define DROPBEAR_SVR_AGENTFWD 0'; \
    echo '#define DROPBEAR_SVR_LOCALTCPFWD 0'; \
    echo '#define DROPBEAR_SVR_REMOTETCPFWD 0'; \
    echo '#define DROPBEAR_SHA1_HMAC 0'; \
    echo '#define DO_MOTD 0'; \
    echo '#define DSS_PRIV_FILENAME "/tmp/c87rs_dss_host_key"'; \
    echo '#define RSA_PRIV_FILENAME "/tmp/c87rs_rsa_host_key"'; \
    echo '#define ECDSA_PRIV_FILENAME "/tmp/c87rs_ecdsa_host_key"'; \
    echo '#define ED25519_PRIV_FILENAME "/tmp/c87rs_ed25519_host_key"'; \
  } > localoptions.h; \
  \
  autoconf; \
  ./configure \
    --enable-static \
  ; \
  \
  make PROGRAMS=dropbear; \
  \
  mv /usr/src/dropbear/dropbear /dropbear; \
  \
  cd /; \
  rm -rf /usr/src/dropbear; \
  /dropbear -V; \
  \
# cleanup of packages
  apt-mark auto '.*' > /dev/null; \
  apt-mark manual $savedAptMark > /dev/null; \
  apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \
  \
  cd /; \
  /dropbear -V

FROM scratch
COPY --from=builder /jq /cloud87/jq
COPY --from=builder /dropbear /cloud87/dropbear

# copy over the init script
COPY remote_shell.sh /cloud87/remote_shell_init

# # Really, dont use this image 
# ENTRYPOINT ["/bin/false"]

# CMD ["/bin/false"]
