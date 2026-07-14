#
# Initializes the project workspace by preparing
# Devkit-managed resources required during development
#
init_workspace() {
    echo "Initializing workspace..."

    mkdir -p "$WORKSPACE/.codex/hooks"

    cp -a \
        "$RUNTIME_HOME/hooks/." \
        "$WORKSPACE/.codex/hooks/"
}
