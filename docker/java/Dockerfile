FROM eclipse-temurin:21-jdk-noble

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        ca-certificates \
        curl \
        git \
        gosu \
        python3 && \
    rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://deb.nodesource.com/setup_26.x | bash - && \
    apt-get update && \
    apt-get install -y --no-install-recommends nodejs && \
    npm install --global @openai/codex@0.143.0 && \
    rm -rf /var/lib/apt/lists/*

COPY runtime/ /opt/codex-runtime/

RUN chmod +x \
    /opt/codex-runtime/bin/entrypoint.sh \
    /opt/codex-runtime/hooks/prevent_repo_scan.py

WORKDIR /workspace

ENTRYPOINT ["/opt/codex-runtime/bin/entrypoint.sh"]
CMD ["sleep", "infinity"]
