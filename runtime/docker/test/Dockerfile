FROM golang:1.26-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ["cmd/", "cmd/"]
COPY internal/ internal/

RUN go build \
    -trimpath \
    -ldflags="-s -w -buildid=" \
    -o /out/mipe \
    ./cmd/mipe

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        gosu && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/mipe /usr/local/bin/mipe
COPY docker/test/bin/ /bin/
COPY docker/test/config/ /opt/mipe/config

RUN chown -R 1000:1000 /test

ENTRYPOINT ["mipe"]
CMD ["/bin/integration-test"]
