# Run Docker Compose
[arg("args", help="Arguments passed to docker compose")]
compose *args:
    @docker compose {{args}}

# Remove all Mipe images
docker-clean:
    @docker image rm -f \
        mipe-codex:21-java \
        mipe-codex:26-node 2>/dev/null || true
