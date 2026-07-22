# Note: Release-Candidate Isolation

## Context

After introducing image promotion, I evaluated how a release candidate should differ from the frequent development snapshot workflow. The goal was to validate the exact artifact that will become a release while keeping development builds fast and inexpensive.

This decision also determines how build metadata must cross from the RC build into the release workflow. That handoff is documented in [Build Metadata Handoff](011-note-build-metadata-handoff.md).

## Observations

Development CI uses GitHub Actions cache scopes to reduce the cost of rebuilding shared runtime, Node, Java, and Web stages. The RC workflow is separate from development CI, and its solves cannot reliably use the development workflow's cache records as their source of truth.

Trying to make one workflow conditionally support both cache policies added complexity to the Bake configuration and made it harder to reason about which cache records a release candidate actually used.

Release candidates have a different purpose from development snapshots. They are deliberate, versioned builds that must be tested and published with the metadata required for a later release. SBOM and provenance generation is therefore enabled for RCs even though it remains disabled for frequent development builds.

The stable release can promote the images produced by a successful RC. Rebuilding during release would create a second set of image descriptors and layers, weakening the connection between the artifact that was tested and the artifact that was released.

## Analysis

Development caching and release confidence are separate concerns. A cache hit is useful for a development snapshot, but a release candidate should not depend on whether another workflow's cache records happen to be available or compatible with its solve graph.

The stable design is to keep development and RC workflows separate. Development CI uses its dedicated cache policy and publishes development tags. RC CI performs a clean build, runs the required tests, publishes candidate tags, and attaches SBOM and provenance. The release workflow then promotes the validated RC images by immutable digest.

This separation also keeps attestation cost limited to deliberate release candidates and prevents release-specific conditions from spreading through the development build configuration.

## Conclusions

- Development CI is optimized for fast feedback and cache reuse
- RC CI is intentionally uncached with respect to development workflow caches
- RC CI is the point where release-grade attestations are generated
- Stable releases promote the validated RC images instead of rebuilding them
- Separate workflows make the cache, attestation, and publishing policies explicit

## Next Steps

- Keep RC cache behavior explicit and avoid introducing development cache imports without measuring their reliability
- Preserve the source commit and published image digests as the handoff from RC to release
- Revisit registry-backed caching only if clean RC builds become prohibitively expensive

## References

- [Package Pins and Optional Build Attestations](009-note-pins-and-attestations.md)
- [Docker build cache backends](https://docs.docker.com/build/cache/backends/)
- [Docker build attestations](https://docs.docker.com/build/metadata/attestations/)
