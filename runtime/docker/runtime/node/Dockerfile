ARG NODE_VERSION=22
FROM node:${NODE_VERSION}-bookworm-slim AS node

FROM debian:bookworm-slim

RUN echo hello > /hello.txt

COPY --from=node /usr/local/ /usr/local/
