# Note: Reproducible Image Publishing

## Context

After aligning the BuildKit cache graphs, I investigated why repeated builds could still produce different image layers and trigger unnecessary GHCR uploads. The focus shifted from cache-record selection to the stability of generated files, build timestamps, image versions, and registry publishing behavior.

## Observations

Using the latest commit timestamp as `SOURCE_DATE_EPOCH` made otherwise identical image layers receive different timestamps on every commit. A fixed epoch made the timestamp rewriting stable across builds and commits.

The build version originally contained the commit SHA, which forced the Go binary and its runtime layer to change on every commit. The version is now derived from the contents that affect the Go build: `go.mod`, `go.sum`, and production Go source files. Test files are excluded from the version hash and Docker build context.

Independent no-cache builds initially differed in one large agent npm layer. The differing files were generated npm cache indexes, npm logs, and Node compile-cache files. Removing these generated files at the end of each agent npm installation made the OCI manifests and all layer digests identical.

BuildKit debug logs showed that GHCR may still receive a repeated upload for an identical small image config blob. The image and layer digests can remain unchanged while the registry receives another upload request. No further large layer uploads appeared after the generated npm files were removed.

## Analysis

I separated three additional concepts from the cache-graph discussion:

- Reproducible build inputs determine whether independent builds produce the same files.
- Stable timestamps and content-derived versions determine whether those files produce the same image descriptors.
- Registry deduplication determines whether an already-known blob is reused without transmitting it again.

Caching can hide nondeterministic build output, but it cannot make that output reproducible. Conversely, identical image digests prove that the image contents are unchanged, but do not prove that the exporter skipped every registry request.

The stable design is to keep the epoch fixed, derive the Go build version from its actual inputs, and remove generated package-manager and runtime caches from published layers. Agent npm layers can remain outside exported BuildKit caches when their parallel rebuild cost is lower than their cache export cost, provided those rebuilds are deterministic.

## Conclusions

The remaining image differences were caused by generated npm and Node cache files, not by `SOURCE_DATE_EPOCH` or the Dockerfile dependency graph. Cleaning those files restored reproducible agent images and prevented large unchanged layers from being republished as different blobs.

Repeated manifest or small config uploads are separate from layer reproducibility. They may still occur when BuildKit publishes an unchanged image, but they are negligible compared with retransmitting large layer contents.

The important design rule is to make every published layer deterministic independently of whether BuildKit reuses its cache record. Stable timestamps, content-derived application versions, and cleanup of generated build artifacts provide that guarantee.

## Next Steps

- Keep fixed `SOURCE_DATE_EPOCH` unchanged unless the reproducibility policy intentionally changes
- Keep Go build version based on the production Go inputs rather than the commit timestamp or SHA
- Treat any future differing layer digest as a reproducibility failure and inspect its files before changing cache policy
- Continue monitoring GHCR debug output for unexpected large layer uploads

## References

- [BuildKit Cache Graph Alignment](007-note-buildkit-cache-graph-alignment.md)
- [Docker BuildKit configuration](https://docs.docker.com/build/buildkit/configure/)
- [Docker buildx build documentation](https://docs.docker.com/reference/cli/docker/buildx/build/)
- [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec/blob/main/spec.md)
