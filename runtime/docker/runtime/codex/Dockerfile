ARG NODE_VERSION=22
FROM node:${NODE_VERSION}-bookworm-slim AS node

FROM runtime AS codex-base

COPY --from=node /usr/local/ /usr/local/

ARG CODEX_VERSION
RUN npm install --global @openai/codex@${CODEX_VERSION}

WORKDIR /workspace

COPY docker/runtime/codex/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/codex/config/agent/ /opt/mipe/config/agent/

FROM codex-base AS web

COPY docker/runtime/codex/config/agent/web.config.json /opt/mipe/config/config.json
COPY docker/runtime/codex/config/agent/web.config.toml \
    /opt/mipe/config/agent/config.toml

CMD ["mipe", "codex"]

FROM codex-base AS default

CMD ["mipe", "codex"]
