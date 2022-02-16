FROM golang:1.15.6 as builder

ARG GOPRIVATE

WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum

COPY hack/ hack/
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN --mount=type=ssh --mount=type=secret,id=github_pat GITHUB_PAT_PATH=/run/secrets/github_pat ./hack/setup-git-redirect.sh \
  && mkdir -p -m 0600 ~/.ssh \
  && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts \
  && go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
COPY res/ res/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/inventory cmd/inventory/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/nic-updater cmd/nic-updater/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/benchmark cmd/benchmark/main.go

FROM busybox:1.32.0

WORKDIR /app
COPY --from=builder /build/bin/inventory .
COPY --from=builder /build/bin/nic-updater .
COPY --from=builder /build/bin/benchmark .
COPY --from=builder /build/res/pci.ids ./res/

ENTRYPOINT ["/app/inventory"]
