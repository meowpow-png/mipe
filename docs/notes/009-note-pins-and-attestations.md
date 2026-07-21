# Note: Package Pins and Optional Build Attestations

## Context

After making image layers reproducible, I evaluated exact package pins and image attestations. The goal was to improve build determinism without adding unnecessary cost to the development snapshot workflow.

## Observations

The versions of the directly installed APT packages are supplied through Docker Bake and consumed by the Dockerfile as explicit build arguments. The current pins cover Temurin, Chromium, and Git.

Exact APT pins improve reproducibility, but builds still depend on the pinned package versions remaining available in the configured repositories. A package being valid today does not guarantee that a future rebuild can retrieve it from a live repository.

The `show-versions` Just recipe reports installed and candidate package versions from built images. The Dockerfile also prints selected package versions while building the shared runtime and Node-based stages.

CI currently disables both SBOM and provenance generation. A trial of enabling them increased the final image build time by approximately two minutes. The additional time came from SBOM scanning each final image.

## Analysis

Package pinning makes the build inputs explicit but does not provide indefinite package availability. Strict reproducibility therefore requires repository retention, a snapshot repository, or a process that updates pins before the current versions disappear.

SBOM and provenance are optional metadata for published images. If enabled, they are attached to the image index as separate manifests. They do not change the runnable image layers or BuildKit cache records, but they do change the top-level image index digest.

Generating attestations for every development snapshot provides complete metadata but imposes a recurring build cost. The current workflow publishes development images on every relevant `dev` build, so that cost is paid frequently without providing much additional value for short-lived snapshots.

An eventual release workflow could promote an immutable development image without rebuilding its layers and generate SBOM and provenance artifacts for the GitHub release. Retagging alone would not generate attestations; those would need to be produced and uploaded as a separate step.

## Conclusions

- Keep exact package versions in Bake and install them with pinned APT versions
- Treat live repository retention as an operational constraint of every APT pin
- Keep SBOM and provenance generation disabled for frequent development snapshots for now
- Reconsider attestations when a release workflow or selected-snapshot process exists
- Use an SBOM as the authoritative component and version inventory when one is generated

## Next Steps

- Decide whether pinned package updates require a snapshot repository or an automated refresh process
- Add release-image promotion and attestation generation once a release workflow exists
- Record the promoted source digest so a release tag identifies the exact development image that was released

## References

- [Reproducible Image Publishing](008-note-reproducible-image-publishing.md)
- [Docker build attestations](https://docs.docker.com/build/metadata/attestations/)
- [Docker SBOM attestations](https://docs.docker.com/build/metadata/attestations/sbom/)
- [Docker Buildx image inspection](https://docs.docker.com/reference/cli/docker/buildx/imagetools/inspect/)
