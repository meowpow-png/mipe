FROM runtime

ARG CODEX_VERSION
RUN npm install --global @openai/codex@${CODEX_VERSION}

WORKDIR /workspace

COPY config/codex/ /opt/mipe/config/
COPY docker/runtime/codex/config.json /opt/mipe/config.json

CMD ["mipe", "codex"]
