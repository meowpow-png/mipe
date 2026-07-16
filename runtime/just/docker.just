# Build docker images
[arg("targets", help="Bake targets to build")]
build-images *targets='default':
    @docker buildx bake {{targets}}

# Run base docker compose stack
[arg("args", help="Arguments passed to docker compose")]
compose *args:
    @docker compose {{args}}
