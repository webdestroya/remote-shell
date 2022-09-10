FROM golang:1.19 AS builder
WORKDIR /tmp/gobuild
COPY . .
RUN go build \
  -ldflags "-linkmode external -extldflags -static" \
  -a \
  -o remote_shell

FROM scratch
WORKDIR /cloud87
# RUN mkdir -p lib bin keys
# RUN mkdir -p bin
COPY --from=builder /tmp/gobuild/remote_shell bin/remote_shell