FROM runtime AS codex-base

ARG CODEX_VERSION
RUN npm install --global @openai/codex@${CODEX_VERSION}

WORKDIR /workspace

COPY docker/runtime/codex/config/ /opt/mipe/config/
COPY docker/runtime/codex/config/config.json /opt/mipe/config.json

FROM codex-base AS web

COPY docker/runtime/codex/config/web.config.toml \
    /opt/mipe/config/config.toml

CMD ["mipe", "codex"]

FROM codex-base AS default

CMD ["mipe", "codex"]
