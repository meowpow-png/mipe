# Build Docker images
[arg("targets", help="Bake targets to build")]
build-images *targets='default':
    @VERSION=$(bash scripts/get-go-build-version.sh) \
     docker buildx bake --load --provenance=false --sbom=false {{targets}}

# Build reproducible OCI image layouts
[arg("target", help="Bake target to build")]
[arg("output", help="Output OCI image name")]
build-oci target output:
     @rm -rf .mipe/{{output}}
     @VERSION=$(bash scripts/get-go-build-version.sh) \
     docker buildx bake {{target}} \
       --provenance=false \
       --sbom=false \
       --set '{{target}}.output=type=oci,dest='.mipe/{{output}}',tar=false,rewrite-timestamp=true' \
       --no-cache

# Builds and compare two no-cache builds
build-compare target:
    @trap 'docker image rm -f local-a local-b >/dev/null 2>&1' EXIT; \
    just build-tagged-image {{target}} local-a && \
    just build-tagged-image {{target}} local-b && \
    diffoci diff docker://local-a docker://local-b || test $? -eq 1

# Show latest available package versions
[arg("target", help="Docker image to inspect")]
show-versions target:
	@docker run --rm {{target}} bash -eu -c '\
		apt-get update >/dev/null; \
		apt-cache policy \
			chromium \
			temurin-21-jdk \
			python3 \
			bash \
			gosu \
			ca-certificates'

[private]
build-tagged-image target tag:
    @VERSION=$(bash scripts/get-go-build-version.sh) \
     docker buildx bake {{target}} \
      --load \
      --provenance=false \
      --sbom=false \
      --set '{{target}}.tags={{tag}}' \
      --no-cache
