FROM golang:1.26-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ["cmd/", "cmd/"]
COPY ["internal/", "internal/"]

RUN go build \
    -trimpath \
    -ldflags="-s -w -buildid=" \
    -o /out/mipe \
    ./cmd/mipe

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        ca-certificates \
        gosu && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY --from=builder /out/mipe /usr/local/bin/mipe
COPY docker/test /test

RUN chown -R 1000:1000 /test

CMD ["/test/bin/integration-test"]
