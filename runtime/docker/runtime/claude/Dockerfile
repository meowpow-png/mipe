ARG NODE_VERSION=22
FROM node:${NODE_VERSION}-bookworm-slim AS node

FROM runtime AS claude-base

COPY --from=node /usr/local/ /usr/local/

ARG CLAUDE_VERSION
RUN npm install --global @anthropic-ai/claude-code@${CLAUDE_VERSION}

WORKDIR /workspace

COPY docker/runtime/claude/config/config.json /opt/mipe/config/config.json
COPY docker/runtime/claude/config/agent/ /opt/mipe/config/agent/

FROM claude-base AS web

COPY docker/runtime/claude/config/agent/web.config.json /opt/mipe/config/config.json
COPY docker/runtime/claude/config/agent/web.mcp.json \
    /opt/mipe/config/agent/web.mcp.json

CMD ["mipe", \
     "claude", \
     "--dangerously-skip-permissions", \
     "--mcp-config", \
     "/opt/mipe/config/agent/web.mcp.json"]

FROM claude-base AS default

CMD ["mipe", "claude", "--dangerously-skip-permissions"]
