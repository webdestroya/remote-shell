FROM golang:1.19 AS builder
WORKDIR /tmp/gobuild

COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

COPY . .
RUN go build \
  -ldflags "-linkmode external -extldflags -static" \
  -a \
  -o remote_shell

FROM scratch
WORKDIR /cloud87
COPY --from=builder /tmp/gobuild/remote_shell bin/remote_shell