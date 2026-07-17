# Build docker images
[arg("targets", help="Bake targets to build")]
build-images *targets='default':
    @docker buildx bake {{targets}}

# Build versioned Docker images for development
[arg("targets", help="Bake targets to build")]
build-dev *targets='default':
    @VERSION=$(git describe --tags --always --dirty) \
     COMMIT=$(git rev-parse --short HEAD) \
     BUILD_DATE=$(date -u +%FT%TZ) \
     docker buildx bake {{targets}}

# Run base docker compose stack
[arg("args", help="Arguments passed to docker compose")]
compose *args:
    @docker compose {{args}}
