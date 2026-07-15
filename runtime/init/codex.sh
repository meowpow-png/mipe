#
# Initializes Codex runtime by preparing user's Codex
# home and installing the shared runtime configuration
#
init_codex() {
    echo "Initializing Codex..."

    mkdir -p "$CODEX_HOME"

    cp -r \
        "$RUNTIME_HOME/config/." \
        "$CODEX_HOME/"
}
