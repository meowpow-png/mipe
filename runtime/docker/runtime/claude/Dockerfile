FROM runtime

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        nodejs \
        npm \
        python3 && \
    npm install --global @anthropic-ai/claude-code@2.1.210 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY config/claude/ /opt/mipe/config/
COPY docker/runtime/claude/config.json /opt/mipe/config.json

CMD ["mipe", "claude", "--dangerously-skip-permissions"]
