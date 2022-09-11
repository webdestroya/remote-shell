ARG BUILD_VERSION=master
ARG BUILD_SHA=devel

# DEBIAN BUILD
FROM golang:1.19 AS builder-deb
WORKDIR /tmp/gobuild

COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

ARG BUILD_VERSION
ARG BUILD_SHA
COPY . .
RUN go build \
  -ldflags "-linkmode external -extldflags -static -X main.buildVersion=${BUILD_VERSION} -X main.buildSha=${BUILD_SHA} -s -w" \
  -a \
  -o remote-shell




# ALPINE BUILD
FROM golang:1.19-alpine AS builder-alp
RUN apk add --update gcc g++

WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

ARG BUILD_VERSION
ARG BUILD_SHA
COPY . .
RUN go build \
  -ldflags "-linkmode external -extldflags -static -X main.buildVersion=${BUILD_VERSION} -X main.buildSha=${BUILD_SHA} -s -w" \
  -a \
  -o remote-shell




# ACTUAL IMAGE
FROM scratch
WORKDIR /cloud87
COPY --from=builder-deb /tmp/gobuild/remote-shell bin/remote-shell

WORKDIR /cloud87-alpine
COPY --from=builder-alp /tmp/gobuild/remote-shell bin/remote-shell