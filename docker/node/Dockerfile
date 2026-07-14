FROM node:26-bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        ca-certificates \
        curl \
        git \
        gosu \
        python3 && \
    npm install --global \
        @openai/codex@0.143.0 \
        playwright && \
    npx playwright install --with-deps chromium && \
    rm -rf /var/lib/apt/lists/*

COPY runtime/ /opt/codex-runtime/

RUN chmod +x \
    /opt/codex-runtime/bin/entrypoint.sh \
    /opt/codex-runtime/hooks/prevent_repo_scan.py

WORKDIR /workspace

ENTRYPOINT ["/opt/codex-runtime/bin/entrypoint.sh"]
CMD ["sleep", "infinity"]
