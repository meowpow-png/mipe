#!/usr/bin/env bash
set -euo pipefail

export CODEX_HOME="$HOME/.codex"
export RUNTIME_HOME=/opt/codex-runtime

source "$RUNTIME_HOME/init/codex.sh"
source "$RUNTIME_HOME/init/permissions.sh"

init_codex
init_permissions

gosu "${LOCAL_UID}:${LOCAL_GID}" \
    env HOME="$HOME" CODEX_HOME="$CODEX_HOME" \
    "$RUNTIME_HOME/init/project.sh"

exec gosu "${LOCAL_UID}:${LOCAL_GID}" \
    env HOME="$HOME" CODEX_HOME="$CODEX_HOME" \
    "$@"
