#
# Discovers and invokes project-specific dependency initialization
#
# The shared runtime does not install dependencies itself
# Projects may provide dependency initialization script to extend this phase
#
init_dependencies() {
    local script="$WORKSPACE/.codex/init/dependencies.sh"

    if [[ -f "$script" ]]; then
        echo "Initializing project dependencies..."

        source "$script"
        install_dependencies
    else
        echo "No project dependency initialization found, skipping."
    fi
}
