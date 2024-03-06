FROM golang:1.22.1 as builder

ARG GOARCH

ENV GOPRIVATE='github.com/onmetal/*'

WORKDIR /build
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

COPY hack/ hack/
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN --mount=type=ssh --mount=type=secret,id=github_pat \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  GITHUB_PAT_PATH=/run/secrets/github_pat ./hack/setup-git-redirect.sh \
  && mkdir -p -m 0600 ~/.ssh \
  && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts \
  && go mod download

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
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GO111MODULE=on go build -a -o bin/benchmark cmd/benchmark/main.go
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GO111MODULE=on go build -a -o bin/benchmark-scheduler cmd/benchmark-scheduler/main.go

FROM debian:testing-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends fio stress-ng\
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/bin/inventory .
COPY --from=builder /build/bin/nic-updater .
COPY --from=builder /build/bin/benchmark .
COPY --from=builder /build/bin/benchmark-scheduler .
COPY --from=builder /build/res/pci.ids ./res/
