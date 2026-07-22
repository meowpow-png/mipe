ARG NODE_VERSION=22.23.1

FROM golang:1.26-bookworm AS runtime-builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ["cmd/", "cmd/"]
COPY ["internal/", "internal/"]

ARG VERSION

RUN go build \
    -trimpath \
    -ldflags="-s -w -buildid= \
        -X github.com/meowpow-png/mipe/runtime/internal/build.Version=${VERSION}" \
    -o /out/mipe \
    ./cmd/mipe

FROM debian:bookworm-slim AS runtime-base

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        gosu \
        passwd \
        python3 \
        ca-certificates && \
    dpkg-query \
        --show \
        --showformat='${Package}: ${Version}\n' \
        bash \
        gosu \
        passwd \
        python3 \
        ca-certificates && \
    rm -rf /var/lib/apt/lists/* \
           /var/log/apt/history.log \
           /var/log/apt/term.log \
           /var/log/dpkg.log \
           /var/cache/ldconfig/aux-cache

FROM runtime-base AS runtime

COPY --from=runtime-builder /out/mipe /usr/local/bin/mipe

COPY docker/runtime/base/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/base/config/agent/ /opt/mipe/config/agent/
COPY hooks/ /opt/mipe/hooks/

COPY docker/runtime/base/entrypoint.sh /usr/local/bin/entrypoint

ENTRYPOINT ["entrypoint"]
CMD ["mipe", "--help"]

FROM runtime AS test

ARG LOCAL_UID=1000
ARG LOCAL_GID=1000

RUN mkdir -p /workspace && \
    chown "${LOCAL_UID}:${LOCAL_GID}" /workspace

FROM node:${NODE_VERSION}-bookworm-slim AS node-source

RUN printf 'node: %s\n' "$(node --version)" && \
    printf 'npm: %s\n' "$(npm --version)" && \
    printf 'corepack: %s\n' "$(corepack --version)"

FROM runtime-base AS node-base

ARG GIT_VERSION

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git="${GIT_VERSION}" && \
    dpkg-query \
        --show \
        --showformat='${Package}: ${Version}\n' \
        git && \
    rm -rf /var/lib/apt/lists/* \
           /var/log/apt/history.log \
           /var/log/apt/term.log \
           /var/log/dpkg.log \
           /var/cache/ldconfig/aux-cache

COPY --from=node-source /usr/local/ /usr/local/

FROM node-base AS java-base

ARG SOURCE_DATE_EPOCH
ARG TEMURIN_21_JDK_VERSION

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        gpg \
        faketime && \
    curl -fsSL https://packages.adoptium.net/artifactory/api/gpg/key/public | \
        gpg --dearmor -o /usr/share/keyrings/adoptium.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/adoptium.gpg] https://packages.adoptium.net/artifactory/deb bookworm main" \
        > /etc/apt/sources.list.d/adoptium.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        temurin-21-jdk="${TEMURIN_21_JDK_VERSION}" && \
    rm -f /etc/ssl/certs/adoptium/cacerts && \
    FAKE_TIME="$(date -u -d "@${SOURCE_DATE_EPOCH}" '+%Y-%m-%d %H:%M:%S')" && \
    faketime -f "@${FAKE_TIME}" \
        trust extract --format=java-cacerts /etc/ssl/certs/adoptium/cacerts && \
    apt-get purge -y --auto-remove faketime libfaketime && \
    rm -rf /var/lib/apt/lists/* \
           /var/log/apt/history.log \
           /var/log/apt/term.log \
           /var/log/dpkg.log \
           /var/log/alternatives.log \
           /var/cache/ldconfig/aux-cache

COPY --from=runtime /usr/local/bin/mipe /usr/local/bin/mipe
COPY --from=runtime /opt/mipe/config/ /opt/mipe/config/
COPY --from=runtime /opt/mipe/hooks/ /opt/mipe/hooks/
COPY --from=runtime /usr/local/bin/entrypoint /usr/local/bin/entrypoint

ENTRYPOINT ["entrypoint"]
CMD ["mipe", "--help"]

FROM node-base AS web-base

ARG PLAYWRIGHT_MCP_VERSION
ARG CHROMIUM_VERSION

ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV CHROME_BIN=/usr/bin/chromium

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        chromium="${CHROMIUM_VERSION}" \
        fonts-liberation && \
    npm install --global \
        @playwright/mcp@${PLAYWRIGHT_MCP_VERSION} && \
    rm -rf /var/lib/apt/lists/* \
           /var/log/apt/history.log \
           /var/log/apt/term.log \
           /var/log/dpkg.log \
           /var/log/alternatives.log \
           /var/cache/fontconfig/* \
           /var/cache/ldconfig/aux-cache \
           /etc/machine-id \
           /var/lib/dbus/machine-id \
           /root/.npm \
           /tmp/node-compile-cache

COPY --from=runtime /usr/local/bin/mipe /usr/local/bin/mipe
COPY --from=runtime /opt/mipe/config/ /opt/mipe/config/
COPY --from=runtime /opt/mipe/hooks/ /opt/mipe/hooks/
COPY --from=runtime /usr/local/bin/entrypoint /usr/local/bin/entrypoint

ENTRYPOINT ["entrypoint"]
CMD ["mipe", "--help"]

FROM node-base AS codex

ARG CODEX_VERSION
RUN npm install --global @openai/codex@${CODEX_VERSION} && \
    rm -rf /root/.npm /tmp/node-compile-cache

COPY --from=runtime /usr/local/bin/mipe /usr/local/bin/mipe
COPY --from=runtime /opt/mipe/config/ /opt/mipe/config/
COPY --from=runtime /opt/mipe/hooks/ /opt/mipe/hooks/
COPY --from=runtime /usr/local/bin/entrypoint /usr/local/bin/entrypoint

ENTRYPOINT ["entrypoint"]

WORKDIR /workspace

COPY docker/runtime/codex/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/codex/config/agent/ /opt/mipe/config/agent/

CMD ["mipe", "codex"]

FROM java-base AS codex-java

ARG CODEX_VERSION
RUN npm install --global @openai/codex@${CODEX_VERSION} && \
    rm -rf /root/.npm /tmp/node-compile-cache

WORKDIR /workspace

COPY docker/runtime/codex/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/codex/config/agent/ /opt/mipe/config/agent/

CMD ["mipe", "codex"]

FROM web-base AS codex-web

ARG CODEX_VERSION
RUN npm install --global @openai/codex@${CODEX_VERSION} && \
    rm -rf /root/.npm /tmp/node-compile-cache

WORKDIR /workspace

COPY docker/runtime/codex/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/codex/config/agent/ /opt/mipe/config/agent/
COPY docker/runtime/codex/config/agent/web.config.json /opt/mipe/config/config.json
COPY docker/runtime/codex/config/agent/web.config.toml /opt/mipe/config/agent/config.toml

CMD ["mipe", "codex"]

FROM node-base AS claude

ARG CLAUDE_VERSION
RUN npm install --global @anthropic-ai/claude-code@${CLAUDE_VERSION} && \
    rm -rf /root/.npm /tmp/node-compile-cache

COPY --from=runtime /usr/local/bin/mipe /usr/local/bin/mipe
COPY --from=runtime /opt/mipe/config/ /opt/mipe/config/
COPY --from=runtime /opt/mipe/hooks/ /opt/mipe/hooks/
COPY --from=runtime /usr/local/bin/entrypoint /usr/local/bin/entrypoint

ENTRYPOINT ["entrypoint"]

WORKDIR /workspace

COPY docker/runtime/claude/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/claude/config/agent/ /opt/mipe/config/agent/

CMD ["mipe", "claude", "--dangerously-skip-permissions"]

FROM java-base AS claude-java

ARG CLAUDE_VERSION
RUN npm install --global @anthropic-ai/claude-code@${CLAUDE_VERSION} && \
    rm -rf /root/.npm /tmp/node-compile-cache

WORKDIR /workspace

COPY docker/runtime/claude/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/claude/config/agent/ /opt/mipe/config/agent/

CMD ["mipe", "claude", "--dangerously-skip-permissions"]

FROM web-base AS claude-web

ARG CLAUDE_VERSION
RUN npm install --global @anthropic-ai/claude-code@${CLAUDE_VERSION} && \
    rm -rf /root/.npm /tmp/node-compile-cache

WORKDIR /workspace

COPY docker/runtime/claude/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/claude/config/agent/ /opt/mipe/config/agent/
COPY docker/runtime/claude/config/agent/web.config.json /opt/mipe/config/config.json
COPY docker/runtime/claude/config/agent/web.mcp.json /opt/mipe/config/agent/web.mcp.json

CMD ["mipe", "claude", "--dangerously-skip-permissions", "--mcp-config", "/opt/mipe/config/agent/web.mcp.json"]
