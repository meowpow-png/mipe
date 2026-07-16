FROM mipe-runtime

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        nodejs \
        npm \
        python3 && \
    npm install --global @openai/codex@0.143.0 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

CMD ["codex"]
