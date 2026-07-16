FROM runtime

ARG LOCAL_UID=1000
ARG LOCAL_GID=1000

COPY docker/test/bin/ /bin/
COPY docker/test/config/ /opt/mipe/config/
COPY docker/test/init/ /workspace/.mipe/init/

RUN chown -R "${LOCAL_UID}:${LOCAL_GID}" /workspace

CMD ["--config", "/opt/mipe/config/config.json", "/bin/integration-test"]
