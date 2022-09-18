FROM alpine AS builder

RUN set -eux; \
  apk add --no-cache \
    ca-certificates \
  ;

# Static vars
ENV GITHUB_REPO="webdestroya/remote-shell"

# Build Args
ARG RSHELL_VERSION=master

WORKDIR /rshell

RUN set -eux; \
  \
  osArch="$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)"; \
  \
  wget -O rshell.tar.gz "https://github.com/${GITHUB_REPO}/releases/download/v${RSHELL_VERSION}/remote-shell_${RSHELL_VERSION}_linux_${osArch}.tar.gz"; \
  \
  mkdir -p /cloud87; \
  tar -xzf rshell.tar.gz -C /cloud87;\
  \
  /cloud87/remote-shell -version;

# ACTUAL IMAGE
FROM scratch
WORKDIR /cloud87
COPY --from=builder /cloud87 .
