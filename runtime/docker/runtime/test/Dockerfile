FROM runtime

ARG LOCAL_UID=1000
ARG LOCAL_GID=1000

COPY docker/runtime/test/bin/ /bin/
COPY docker/runtime/test/config/ /opt/mipe/config/
COPY docker/runtime/test/init/ /workspace/.mipe/init/

RUN chown -R "${LOCAL_UID}:${LOCAL_GID}" /workspace

CMD ["--config", "/opt/mipe/config/config.json", "/bin/integration-test"]
