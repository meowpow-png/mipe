# Note: CI Docker Build Caching

## Context

While building the runtime images in CI, I investigated why repeated builds without source changes did not consistently reuse Docker layers. The problem was most visible in the `java-base` and `web-base` stages, especially around `apt-get` and package installation. The goal was to understand the cache semantics well enough to decide which images should use CI cache storage and which should simply be rebuilt.

## Observations

CI imports and exports BuildKit cache manifests. Re-running the same workflow without new commits therefore does not by itself guarantee that every `RUN` instruction will be reported as `CACHED`; the cache must contain a matching record for the exact stage, instruction, inputs, arguments, and parent filesystem.

The pinned Node image digest remained identical between runs, and the runtime image was generally stable. This ruled out a changing Node base image or a changing runtime image timestamp as the primary explanation for the unstable Java and Web layers.

The build originally used separate Dockerfiles and Bake targets with target-to-target dependencies. I then consolidated the image stages into `docker/Dockerfile` and removed the target contexts. This made the dependency graph explicit in one file, but it did not make every target share one cache record: each Bake target still represents a separate solve and can have its own cache import and export configuration.

Repeated log comparisons showed that runtime, node-base, Codex, and Claude were the most reliable cache candidates. Java and Web package-install layers were repeatedly rebuilt, and exporting their dedicated caches added enough overhead that removing those exports improved CI time by roughly thirty seconds.

Cache mounts were also tested for APT and npm directories. A cache mount is a BuildKit-managed directory available while a `RUN` instruction executes; it is not the same as a cached image layer and does not make the instruction itself appear as `CACHED`. On ephemeral CI builders, its benefit depends on the mount contents being persisted and restored by the selected cache backend.

The registry cache backend was replaced with the GitHub Actions cache backend because these caches are only needed during CI and release images are published to Docker Hub. Separate cache scopes were retained for independent image families to avoid concurrent jobs overwriting one shared cache reference.

## Analysis

A BuildKit layer is reusable only when the instruction and all relevant inputs match. Relevant inputs include the parent stage, copied files, build arguments, environment, base-image digest, and the exact Dockerfile instruction. A cache manifest being imported proves that cache metadata was available; it does not prove that a matching layer exists for every target.

The runtime, Java, and Web stages do not have identical cache behavior even though they are defined in the same Dockerfile. Java and Web add large package-install steps after `node-base`, so their cache keys include the resulting parent filesystem and the package-install instruction. Their package managers also access external repositories, but changing repository contents is not required for a miss: a missing export, a different solve scope, or a changed parent record is sufficient.

The earlier sequential Bake experiment made the result worse because it changed the solve and cache-import behavior rather than merely ordering identical work. A shared cache reference was intentionally avoided because concurrent CI runs can race while exporting it. Distinct scopes provide isolation, but they also mean that a target can reuse another scope only when that scope is explicitly imported.

The current cache policy reflects observed value rather than an assumption that every layer must be cached. Runtime and node-base are expensive shared foundations. Codex and Claude have useful, isolated package-install caches. Java and Web package layers have not demonstrated reliable reuse and currently cost more to export than they save. Building those stages without dedicated cache exports is therefore a predictable performance tradeoff, not evidence that the Dockerfile is ignoring cache metadata.

## Conclusions

The inconsistent Java and Web cache hits were not caused primarily by the pinned Node digest changing between immediate reruns, nor by the runtime image being rebuilt with a new date or time.

The important distinction is between a cache manifest, a reusable image layer, and a cache mount. Importing a manifest only makes candidate records available. A layer is reused only when its complete BuildKit cache key matches. A cache mount accelerates package-manager work inside a rebuilt instruction but does not replace layer caching.

Unifying the Dockerfile removed unnecessary file and target indirection, but it could not guarantee cache reuse across separate Bake solves. Sequential builds likewise did not solve the underlying cache-key and scope behavior.

The practical CI configuration is to use GitHub Actions cache scopes for runtime, node-base, Codex, and Claude, while leaving Java and Web without dedicated cache exports until their cache behavior or cost profile justifies another experiment.

## Next Steps

- Keep cache scopes separate for concurrent CI builds and import shared foundation scopes only where a target consumes them
- Treat `java-base` and `web-base` as ordinary CI build stages unless logs show reliable layer reuse
- When investigating future misses, compare the full BuildKit command, target, imported scopes, build arguments, parent digest, and copied-file inputs before changing the Dockerfile
- Keep release image publication independent from CI cache storage

## References

- [Docker cache backends](https://docs.docker.com/build/cache/backends/)
- [GitHub Actions cache backend](https://docs.docker.com/build/cache/backends/gha/)
- [Optimize cache usage](https://docs.docker.com/build/cache/optimize/)
- [GitHub Actions dependency caching](https://docs.github.com/en/actions/reference/workflows-and-actions/dependency-caching)
