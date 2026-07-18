ARG NODE_VERSION=22
FROM node:${NODE_VERSION}-bookworm-slim AS node

FROM runtime

COPY --from=node /usr/local/ /usr/local/

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        gpg \
        curl && \
    curl -fsSL https://packages.adoptium.net/artifactory/api/gpg/key/public | \
        gpg --dearmor -o /usr/share/keyrings/adoptium.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/adoptium.gpg] https://packages.adoptium.net/artifactory/deb bookworm main" \
        > /etc/apt/sources.list.d/adoptium.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        temurin-21-jdk && \
    rm -rf /var/lib/apt/lists/*
