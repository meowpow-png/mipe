#!/usr/bin/env bash
set -euo pipefail

# Default local developer identity
LOCAL_UID="${LOCAL_UID:-1000}"
LOCAL_GID="${LOCAL_GID:-1000}"

# Ensure local developer group exists with the expected GID
if getent group dev >/dev/null; then
    ACTUAL_GID="$(getent group dev | cut -d: -f3)"

    if [ "$ACTUAL_GID" != "$LOCAL_GID" ]; then
        echo "Group 'dev' exists with GID $ACTUAL_GID, expected $LOCAL_GID" >&2
        exit 1
    fi
else
    groupadd --gid "$LOCAL_GID" dev
fi

# Ensure local developer user exists with the expected UID
if id dev >/dev/null 2>&1; then
    ACTUAL_UID="$(id -u dev)"

    if [ "$ACTUAL_UID" != "$LOCAL_UID" ]; then
        echo "User 'dev' exists with UID $ACTUAL_UID, expected $LOCAL_UID" >&2
        exit 1
    fi
else
    useradd \
        --uid "$LOCAL_UID" \
        --gid "$LOCAL_GID" \
        --create-home \
        --home-dir /home/dev \
        --shell /bin/bash \
        dev
fi

exec "$@"
