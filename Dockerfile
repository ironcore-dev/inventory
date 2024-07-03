FROM golang:1.23rc1 AS builder

ARG GOARCH


WORKDIR /build
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

COPY hack/ hack/

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY res/ res/

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GO111MODULE=on go build -a -o bin/inventory cmd/inventory/main.go
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GO111MODULE=on go build -a -o bin/nic-updater cmd/nic-updater/main.go

FROM debian:testing-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends fio stress-ng\
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/bin/inventory .
COPY --from=builder /build/bin/nic-updater .
COPY --from=builder /build/res/pci.ids ./res/
