FROM runtime AS claude-base

ARG CLAUDE_VERSION
RUN npm install --global @anthropic-ai/claude-code@${CLAUDE_VERSION}

WORKDIR /workspace

COPY docker/runtime/claude/config/ /opt/mipe/config/
COPY docker/runtime/claude/config/config.json /opt/mipe/config.json

FROM claude-base AS web

COPY docker/runtime/claude/config/web.mcp.json \
    /opt/mipe/config/web.mcp.json

CMD ["mipe", \
     "claude", \
     "--dangerously-skip-permissions", \
     "--mcp-config", \
     "/opt/mipe/config/web.mcp.json"]

FROM claude-base AS default

CMD ["mipe", "claude", "--dangerously-skip-permissions"]
