# Continuous Integration

CI verifies runtime code, builds container images, runs integration tests, and publishes development images from `dev`. It keeps shared build layers warm so routine changes do not pay full dependency-install cost.

## Workflow

Pipeline runs in one job:

1. Checks out source, sets up Go and Buildx, then calculates runtime build version
2. Runs unit tests with coverage and builds `runtime` plus `test` Bake targets
3. Runs integration tests against loaded test image
4. Builds shared toolchain base targets to refresh BuildKit caches
5. Builds runtime plus all agent images. On `dev`, it also publishes them to GHCR
6. Uploads the unit-test coverage report to Codecov

```mermaid
flowchart TD
    source["Runtime source"] --> verify["Go unit tests"]
    verify --> testImage["Runtime and test image"]
    testImage --> integration["Integration tests"]
    integration --> bases["Toolchain base images"]
    bases --> images["Runtime and agent images"]
    images -->|dev only| ghcr[("GHCR")]
    images --> codecov[("Codecov")]
```

Builds use plain BuildKit output, making cache hits and misses visible in job logs.

## Caching

CI stores BuildKit records in GitHub Actions cache scopes. Each final image imports every scope used by its dependency chain. For example, Java images import Java, Node, and runtime caches; web images import web, Node, and runtime caches.

```mermaid
flowchart LR
    testImage["Runtime and test image build"] -->|exports| runtimeCache[("runtime cache")]
    baseBuild["Base-cache build"] -->|exports| nodeCache[("node-base cache")]
    baseBuild -->|exports| javaCache[("java-base cache")]
    baseBuild -->|exports| webCache[("web-base cache")]

    runtimeCache -->|imports| finalImages["Final image build"]
    nodeCache -->|imports| finalImages
    javaCache -->|imports| finalImages
    webCache -->|imports| finalImages
```

Scopes follow image dependency graph. Shared dependency installs sit below stable parents, and runtime copies happen after expensive base image setup. Cache-warming targets therefore produce records final image builds can reuse.

CI exports shared runtime and toolchain caches because many images consume them. It does not export any agent npm-install layers. Rebuilding those layers in parallel costs less than exporting and restoring their large, volatile cache records.

> [!NOTE]
> Cache availability does not guarantee a hit. BuildKit also needs matching inputs, build arguments, parent filesystem, and solve graph. Reproducible layers make those records dependable across runs.

Changes rebuild only affected image families:

- Codex changes rebuild Codex images; Claude changes rebuild Claude images.
- Java changes rebuild Java images; Playwright MCP changes rebuild web images.
- Node.js changes rebuild every agent image.
- MIPE runtime changes rebuild final agent layers while Node, Java, and web dependency layers remain cacheable.

## Publishing

Pushes to `dev` publish development images. Release-candidate tags publish candidate images; stable releases promote those images without rebuilding them. CI authenticates to `ghcr.io` with GitHub-provided credentials and publishes under `ghcr.io/<owner>/mipe-runtime`.

Publishing rewrites timestamps in exported layers. Without a fixed `SOURCE_DATE_EPOCH`, that would give unchanged files a new timestamp on every build or commit, changing layer digests for no useful reason. Fixed epoch keeps those timestamps stable. 

Together with pinned build inputs from Bake and generated-cache cleanup, unchanged layers keep the same digest and GHCR can reuse existing blobs instead of receiving large duplicate uploads.

## Releasing

Runtime releases are promoted from a tested release candidate; they are never rebuilt for publication. This ensures the images users receive are the exact images that passed the candidate build, including their SBOM and provenance.

```mermaid
flowchart TD
    commit["Release commit on dev"] --> rcTag["Tag runtime/vX.Y.Z-rc.N"]
    rcTag --> rcBuild["Build, test, publish candidate images"]
    rcBuild --> rcImages["Candidate images\nvX.Y.Z-rc.N"]
    rcImages --> approval{"Candidate approved?"}
    approval -->|no| fix["New fix commit on dev\nand a new RC tag"]
    fix --> rcTag
    approval -->|yes| stableTag["Tag the same commit\nruntime/vX.Y.Z"]
    stableTag --> promote["Promote the tested image digests"]
    promote --> release["Publish GitHub Release\nand verify image tags"]
```

The candidate and stable Git tags point to the same commit, ensuring stable publication uses the exact candidate artifact.

### Release Candidates

A candidate tag such as `runtime/v0.1.0-rc.1` builds and publishes every runtime image with the corresponding candidate image tag:

```text
ghcr.io/<owner>/mipe-runtime-codex:v0.1.0-rc.1
```

The runtime inside those images is built as the intended final version (`0.1.0`). The `-rc.1` suffix is only the temporary registry tag used while the candidate is being evaluated.

### Stable publication

After a candidate succeeds, its matching stable tag triggers publication. Stable publication verifies the candidate, then applies two tags to each already-published image digest:

- `vX.Y.Z` — immutable versioned reference
- `latest` — current stable runtime release

```mermaid
flowchart LR
    candidate["Candidate image\nvX.Y.Z-rc.N"] --> digest["Immutable image digest\nwith SBOM and provenance"]
    digest --> versioned["Stable image tag\nvX.Y.Z"]
    digest --> latest["Stable image tag\nlatest"]
```

The release also creates a GitHub Release named `Runtime vX.Y.Z`. Its notes are the matching `CHANGELOG.md` section. Publication then confirms that both stable image tags resolve to the expected digest and that the GitHub Release exists. See [Contributing](../CONTRIBUTING.md#releases) for the contributor release procedure.

## Versioning

CI uses three complementary identifiers to describe a build. The build version identifies the runtime source revision and is embedded into the binary, for example:

```text
dev-4a8d2c1
```

Image tags provide convenient registry references for pulling images. They may move over time as newer builds are published, for example:

```text
ghcr.io/<owner>/mipe-runtime:dev-latest
ghcr.io/<owner>/mipe-runtime:dev-4a8d2c1
```

Release candidates use a final runtime build version such as `0.1.0`, while their image tags remain candidate-specific, such as `v0.1.0-rc.1`. Stable publication adds `v0.1.0` and `latest` to that same image digest.

Image digests uniquely identify the published OCI image contents. Unlike tags, digests are immutable and always refer to exactly one published image, for example:

```text
ghcr.io/<owner>/mipe-runtime@sha256:abc123...
```

Development build version is computed by a local script. Its hash includes `go.mod`, `go.sum`, and production Go files under `cmd/` and `internal/`. Release candidates instead use their intended Semantic Version.

Tests do not contribute to the build version and are excluded from the Docker build context. Test-only changes still run checks but do not rebuild the runtime binary layer. Production Go source or module changes produce a new version and invalidate layers that depend on the binary.

## Triggering

The development CI pipeline supports manual runs and pushes to `main` or `dev`. Pushes run only when changes affect CI configuration or runtime build inputs: runtime commands, internal code, integration tests, Go modules, Docker files, hooks, Bake files, or Compose configuration. Release pipelines are triggered by runtime candidate and stable tags.

Path filtering avoids image builds for unrelated repository changes. Manual dispatch remains available for cache checks, publishing validation, and other CI investigations.
