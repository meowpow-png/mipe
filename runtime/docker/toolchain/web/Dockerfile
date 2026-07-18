ARG NODE_VERSION=22
FROM node:${NODE_VERSION}-bookworm-slim AS node

FROM runtime

COPY --from=node /usr/local/ /usr/local/

ARG PLAYWRIGHT_MCP_VERSION

ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV CHROME_BIN=/usr/bin/chromium

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        chromium \
        fonts-liberation && \
    npm install --global \
        @playwright/mcp@${PLAYWRIGHT_MCP_VERSION} && \
    rm -rf /var/lib/apt/lists/*
