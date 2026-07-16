# Run Docker Compose
[arg("args", help="Arguments passed to docker compose")]
compose *args:
    @docker compose {{args}}

# Remove all Mipe images
docker-clean:
    @ids=$(docker image ls -q --filter "reference=mipe-*"); \
    [ -n "$ids" ] && docker image rm -f $ids || true
