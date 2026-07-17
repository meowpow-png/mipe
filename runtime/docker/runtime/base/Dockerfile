FROM golang:1.26-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ["cmd/", "cmd/"]
COPY ["internal/", "internal/"]

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

RUN go build \
    -trimpath \
    -ldflags="-s -w -buildid= \
        -X github.com/meowpow-png/mipe/runtime/internal/build.Version=${VERSION} \
        -X github.com/meowpow-png/mipe/runtime/internal/build.Commit=${COMMIT} \
        -X github.com/meowpow-png/mipe/runtime/internal/build.Date=${BUILD_DATE}" \
    -o /out/mipe \
    ./cmd/mipe

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        gosu \
        passwd && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/mipe /usr/local/bin/mipe

COPY docker/runtime/base/config/agent/ /opt/mipe/config/agent/
COPY hooks/ /opt/mipe/hooks/

COPY docker/runtime/base/entrypoint.sh /usr/local/bin/entrypoint

ENTRYPOINT ["entrypoint"]
CMD ["mipe", "--help"]
