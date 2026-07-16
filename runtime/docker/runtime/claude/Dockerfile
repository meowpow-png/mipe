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
COPY docker/claude/config.json /opt/mipe/config.json

CMD ["mipe", "--config", "/opt/mipe/config.json", "claude", "--dangerously-skip-permissions"]
