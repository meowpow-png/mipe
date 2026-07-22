# Note: Build Metadata Handoff

## Context

After separating release-candidate builds from development CI, I needed a reliable way to pass the RC build results to the release workflow. The release must promote the images that were actually built and tested, so it needs their source commit, CI run, repositories, and immutable digests.

This is a direct consequence of the workflow isolation described in [Release-Candidate Isolation](010-note-release-candidate-isolation.md). Once the build and release are separate workflow runs, the required image information must cross that boundary as an explicit artifact.

## Observations

The Docker Bake action exposes Buildx metadata as JSON output from the build step. In the workflow configuration used here, that metadata is not handed directly to the next step as a file suitable for artifact upload.

Passing the complete JSON document through GitHub step outputs and environment variables can exceed the runner's process argument or environment size limit. The failure appears before the manifest script runs, even though the build itself completed successfully.

The raw metadata contains more information than the release workflow needs. The release only needs the source identity, CI run identity, published image repositories, and validated `containerimage.digest` values for each target.

## Analysis

Buildx metadata is structured build output, not a release manifest. It must be reduced to the small set of fields required for promotion before it is passed between workflow steps or uploaded as an artifact.

The current Bake action interface does not provide a simple, direct metadata-file handoff for this use case. The practical workaround is to consume the metadata immediately, validate the expected target digests, and write a compact release-source manifest file. That file can then be uploaded as a workflow artifact and resolved by the release workflow.

The manifest also provides a stable boundary between Buildx and release logic. Buildx remains responsible for solving and publishing images; the release workflow consumes a deliberately defined record rather than depending on the complete shape or size of Buildx's metadata output.

## Conclusions

- Do not pass the complete Buildx metadata document through large environment variables
- Convert metadata into a compact, validated manifest immediately after the RC build
- Transfer the manifest as an artifact between the RC and release workflows
- Use the published image digests from the manifest as release inputs, subject to registry verification
- Prefer a native file-based metadata output in the future if the Bake action exposes one reliably

## Next Steps

- Keep the manifest schema versioned so release workflows can reject incompatible records
- Validate that every expected image target has a published digest before uploading the manifest
- Keep registry verification in the release workflow because metadata alone does not prove registry availability

## References

- [Reproducible Image Publishing](008-note-reproducible-image-publishing.md)
- [Docker buildx build metadata](https://docs.docker.com/reference/cli/docker/buildx/build/)
- [Docker Bake action](https://github.com/docker/bake-action)
