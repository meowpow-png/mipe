#
# Initializes Codex runtime by preparing user's Codex
# home and installing the shared runtime configuration
#
init_codex() {
    echo "Initializing Codex..."

    mkdir -p "$CODEX_HOME"

    cp \
        "$RUNTIME_HOME/config/config.toml" \
        "$CODEX_HOME/config.toml"
}
