#!/usr/bin/env bash
set -euo pipefail

LOCAL_UID="${LOCAL_UID:?LOCAL_UID is required}"
LOCAL_GID="${LOCAL_GID:?LOCAL_GID is required}"

groupadd --gid "$LOCAL_GID" dev 2>/dev/null || true

useradd \
    --uid "$LOCAL_UID" \
    --gid "$LOCAL_GID" \
    --create-home \
    --home-dir /home/dev \
    --shell /bin/bash \
    dev 2>/dev/null || true

exec mipe "$@"
