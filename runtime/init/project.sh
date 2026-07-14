#!/usr/bin/env bash
set -euo pipefail

export CODEX_HOME="$HOME/.codex"
export RUNTIME_HOME=/opt/codex-runtime
export WORKSPACE=/workspace

source "$RUNTIME_HOME/init/workspace.sh"
source "$RUNTIME_HOME/init/dependencies.sh"

init_workspace
init_dependencies
