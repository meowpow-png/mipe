FROM runtime AS codex-base

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
