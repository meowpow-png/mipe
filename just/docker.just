# Run Docker Compose
[arg("args", help="Arguments passed to docker compose")]
compose *args:
    @docker compose {{args}}

# Remove all Mipe Docker resources
docker-clean:
    @just docker-clean-images
    @just docker-clean-volumes

# Remove all Mipe images
docker-clean-images:
    @ids=$(docker image ls -q --filter "reference=mipe-*"); \
    [ -n "$ids" ] && docker image rm -f $ids || true

# Remove all Mipe volumes
docker-clean-volumes:
    @ids=$(docker volume ls -q --filter "name=^mipe_"); \
    [ -n "$ids" ] && docker volume rm -f $ids || true
