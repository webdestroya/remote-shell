FROM golang:1.19 AS builder
WORKDIR /tmp/gobuild

COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

ARG BUILD_VERSION=master
ARG BUILD_SHA=devel

COPY . .
RUN go build \
  -ldflags "-linkmode external -extldflags -static -X main.buildVersion=${BUILD_VERSION} -X main.buildSha=${BUILD_SHA} -s -w" \
  -a \
  -o remote-shell


# ACTUAL IMAGE
FROM scratch
WORKDIR /cloud87
COPY --from=builder /tmp/gobuild/remote-shell bin/remote-shell
