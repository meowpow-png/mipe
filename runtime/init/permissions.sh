#
# Updates ownership of runtime and workspace files
# to match the local user, ensuring the development
# environment has the correct permissions
#
init_permissions() {
    echo "Updating ownership..."

    chown -R "${LOCAL_UID}:${LOCAL_GID}" "$HOME"
}
