#
# Updates ownership of the home directory to match
# the local user, ensuring persistent runtime state
# has the correct permissions
#
init_permissions() {
    echo "Updating ownership..."

    chown -R "${LOCAL_UID}:${LOCAL_GID}" "$HOME"
}
