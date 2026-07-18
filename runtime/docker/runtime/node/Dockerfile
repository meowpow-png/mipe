ARG NODE_VERSION=22
FROM node:${NODE_VERSION}-bookworm-slim AS node

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        libatomic1 \
        python3 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=node /usr/local/ /usr/local/
