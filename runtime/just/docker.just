# Build Docker images
[arg("targets", help="Bake targets to build")]
build-images *targets='default':
    @docker buildx bake --provenance=false -sbom=false {{targets}}

# Build versioned Docker images for development
[arg("targets", help="Bake targets to build")]
build-dev *targets='default':
    @VERSION=$(git describe --tags --always --dirty) \
     COMMIT=$(git rev-parse --short HEAD) \
     BUILD_DATE=$(date -u +%FT%TZ) \
     docker buildx bake --provenance=false --sbom=false {{targets}}

# Run base Docker Compose stack
[arg("args", help="Arguments passed to Docker Compose")]
compose *args:
    @docker compose {{args}}
