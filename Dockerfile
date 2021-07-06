FROM golang:1.15.6 as builder

WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum

ARG GOPRIVATE
ARG GIT_USER
ARG GIT_PASSWORD
RUN if [ ! -z "$GIT_USER" ] && [ ! -z "$GIT_PASSWORD" ]; then \
        printf "machine github.com\n \
            login ${GIT_USER}\n \
            password ${GIT_PASSWORD}\n \
            \nmachine api.github.com\n \
            login ${GIT_USER}\n \
            password ${GIT_PASSWORD}\n" \
            >> ${HOME}/.netrc;\
    fi
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
COPY res/ res/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/inventory cmd/inventory/main.go

FROM busybox:1.32.0

WORKDIR /app
COPY --from=builder /build/bin/inventory .
COPY --from=builder /build/res/pci.ids ./res/

ENTRYPOINT ["/app/inventory"]
