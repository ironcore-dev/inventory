FROM golang:1.18.3 as builder

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
COPY internal/ internal/
COPY pkg/ pkg/
COPY res/ res/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/inventory cmd/inventory/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/nic-updater cmd/nic-updater/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/benchmark cmd/benchmark/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/benchmark-scheduler cmd/benchmark-scheduler/main.go
RUN git clone https://github.com/axboe/fio.git && \
    cd fio/ && \
    ./configure --build-static && \
    make && \
    make install
RUN git clone https://github.com/ColinIanKing/stress-ng.git && \
    cd stress-ng/ && \
    make clean && \
  	STATIC=1 make


FROM amd64/busybox:1.35.0

WORKDIR /app

COPY --from=builder /build/bin/inventory .
COPY --from=builder /build/bin/nic-updater .
COPY --from=builder /build/bin/benchmark .
COPY --from=builder /build/bin/benchmark-scheduler .
COPY --from=builder /build/fio/fio /usr/local/bin/
COPY --from=builder /build/stress-ng/stress-ng /usr/local/bin/
COPY --from=builder /build/res/pci.ids ./res/

