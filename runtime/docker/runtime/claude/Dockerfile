FROM runtime

ARG CLAUDE_VERSION
RUN npm install --global @anthropic-ai/claude-code@${CLAUDE_VERSION}

WORKDIR /workspace

COPY config/claude/ /opt/mipe/config/
COPY docker/runtime/claude/config.json /opt/mipe/config.json

CMD ["mipe", "claude", "--dangerously-skip-permissions"]
