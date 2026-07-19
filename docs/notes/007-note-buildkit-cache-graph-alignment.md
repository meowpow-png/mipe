# Note: BuildKit Cache Graph Alignment

## Context

After making the OCI layers reproducible, I investigated why CI could still rebuild package-install layers even when cache manifests were imported and local no-cache builds produced identical layer digests. The focus shifted from layer contents to the relationship between Dockerfile stages, Bake targets, and BuildKit cache records.

## Observations

BuildKit can import a cache manifest successfully while still rebuilding a Dockerfile instruction. A manifest only provides candidate records; the record must match the complete solve graph for the target being built.

The same Dockerfile stage behaved differently when built as a standalone Bake target and when reached as a dependency of a final image. In some runs, the standalone base target was cached while the final image solve rebuilt the corresponding stage. In other runs, the final image reused the stage while the separate cache-warming solve performed unnecessary work.

The relevant cache identity included more than the visible package-install instruction. It was affected by the parent stage, inherited build arguments, copied runtime files, target selection, and the complete set of stages reachable from the target. Identical output layer digests from independent local builds therefore did not prove that two CI solves would select the same cache record.

The original stage hierarchy made every dependency image inherit the Go-built `runtime` image. A change to the runtime parent consequently changed the parent filesystem for Node, Java, and Web stages, even when their own instructions and package versions were unchanged.

The web targets also initially omitted `PLAYWRIGHT_MCP_VERSION` from their effective build arguments. This created different cache keys for the standalone Web base and final Web targets despite identical-looking Dockerfile instructions.

## Analysis

I separated three concepts that had been conflated during the investigation:

- OCI reproducibility determines whether independent builds can produce the same layer bytes.
- BuildKit cache identity determines whether a particular solve can reuse a prior operation record.
- Cache backend configuration determines which candidate records are available to that solve.

Making the layer bytes reproducible solved only the first problem. BuildKit still evaluates the operation within its parent graph. A stable APT instruction cannot be reused when its parent stage or effective build arguments differ.

The cache-warming Bake also introduced a second graph for the same dependency stages. Its cache-only targets had their own root configuration, while final images evaluated the stages together with runtime copies and agent-specific arguments. Exporting a cache from the first graph did not guarantee that the second graph would select the same records.

The stable design is to create dependency stages from a stable `runtime-base`, keep Node, Java, and Web installation steps independent of the changing Go-built runtime, and place runtime copies after expensive dependency installation. The final Java and Web images then inherit the same base stages that CI exports. This aligns the cache-warming and final-image graphs while keeping runtime changes out of the APT cache keys.

The cache policy also became intentionally selective. Shared runtime and dependency stages are exported because they are expensive and reused across images. The six agent npm-install layers are rebuilt because their cache exports are large and volatile relative to their parallel rebuild cost.

## Conclusions

The remaining CI cache failures were caused by cache-graph misalignment, not by a contradiction between identical local layer digests and CI behavior. A matching layer digest describes the result; it does not identify the BuildKit operation record or guarantee that another target graph can reach it.

The important design rule is to keep expensive dependency operations below stable parents and ensure that the cache-exporting target is the same graph consumed by final images. Separate cache scopes provide storage isolation, but they do not make different target graphs interchangeable.

The final CI runs confirmed the result: runtime, Node, Java, and Web layers were restored in both the base and final-image solves. Only the intentionally uncached npm layers rebuilt.

## Next Steps

- Add a CI check that compares cache reuse for `runtime-base`, `node-base`, `java-base`, and `web-base` across two builds of the same commit
- Validate every Bake target with `docker buildx bake --print` and compare effective arguments for standalone and final-image solves
- Test a stable dependency-stage epoch against the commit-derived `SOURCE_DATE_EPOCH`, measuring Java cache reuse and OCI reproducibility across two commits
- Keep npm layers out of exported caches unless measured export time is lower than their parallel rebuild time

## References

- [CI Docker Build Caching](004-note-ci-docker-build-caching.md)
- [OCI Layer Reproducibility and CI Cache Misses](005-note-oci-layer-reproducibility.md)
- [APT Layer Reproducibility Solutions](006-note-apt-layer-reproducibility-solutions.md)
- [BuildKit reproducibility documentation](https://github.com/moby/buildkit/blob/master/docs/build-repro.md)
- [Docker cache backends](https://docs.docker.com/build/cache/backends/)
