FROM runtime

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        nodejs \
        npm \
        python3 && \
    npm install --global @openai/codex@0.144.5 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY config/codex/ /opt/mipe/config/
COPY docker/runtime/codex/config.json /opt/mipe/config.json

CMD ["mipe", "codex"]
