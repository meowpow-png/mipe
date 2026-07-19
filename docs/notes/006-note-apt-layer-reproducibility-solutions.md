# Note: APT Layer Reproducibility Solutions

## Context

I investigated official ways to make the Docker image layers produced by the APT installation steps reproducible. The immediate goal was to determine whether repeated local `--no-cache` builds could produce identical OCI layer digests, then identify which parts of the solution should be used in CI.

## Observations

The shared `runtime` APT instruction was the first source of nondeterminism. Before the change, repeated builds differed in APT/dpkg-generated logs and in filesystem timestamps. The affected files included:

- `/var/log/apt/history.log`
- `/var/log/apt/term.log`
- `/var/log/dpkg.log`

The downstream `node-base`, Java, Web, Codex, and Claude targets inherit this runtime layer. Java and Web add their own package-install instructions and may therefore introduce additional independent differences.

The `node-base` copy instruction did not introduce a reproducibility problem after the shared runtime layer was fixed.

## Analysis

I reviewed Debian/APT documentation, Debian reproducible-build guidance, Docker documentation, and BuildKit’s reproducibility documentation.

APT’s documented date settings control repository Release-file validation, such as `Check-Date` and `Check-Valid-Until`; they do not normalize timestamps or generated log contents from `apt-get install`. Disabling those checks would weaken package freshness and replay protection and is not a reproducibility solution.

The official BuildKit mechanism is `SOURCE_DATE_EPOCH`, which fixes image metadata and provides a reference timestamp. BuildKit’s image/OCI exporter additionally supports `rewrite-timestamp=true`, which rewrites file timestamps in exported layers. These mechanisms do not make APT/dpkg log contents deterministic, so the three generated logs must be removed at the end of the same APT `RUN` instruction. Removing them in a later layer would leave the nondeterministic bytes in the earlier layer.

Debian reproducible-build guidance also identifies snapshot repositories as the mechanism for reproducing historical package versions. This addresses changing package/repository contents, but it does not by itself normalize install-time timestamps or generated logs.

### Tested Solutions

I updated the shared runtime APT instruction to remove the three APT/dpkg logs in the same layer:

```dockerfile
rm -rf /var/lib/apt/lists/* \
       /var/log/apt/history.log \
       /var/log/apt/term.log \
       /var/log/dpkg.log
```

I then performed two manual no-cache OCI builds of `runtime` with a fixed Git-derived `SOURCE_DATE_EPOCH` and `rewrite-timestamp=true`. The two builds produced identical OCI indexes, manifests, config digests, and ordered layer digests.

I repeated the same test for `node-base`. Its OCI indexes, manifests, config digests, and ordered layer digests were also identical. Its own Node copy layer was stable, confirming that the runtime fix propagates successfully to this downstream target.

### Build Configuration

`SOURCE_DATE_EPOCH` can be represented as a Bake variable and passed as a build argument, with the value supplied dynamically by the local environment or CI.

`rewrite-timestamp=true` is an exporter option, not a Dockerfile instruction. It must be applied to an image/OCI export. The Docker exporter used by `--load` performs unpacking, and combining it directly with timestamp rewriting can fail with:

```text
exporter option "rewrite-timestamp" conflicts with "unpack"
```

The local Just path therefore intentionally keeps `--load` and does not enable timestamp rewriting. CI is intended to supply the fixed timestamp and image exporter override, but that exact CI Bake/action combination remains unvalidated after the exporter conflict was discovered.

## Conclusions

The confirmed solution for reproducible OCI output is:

1. Remove APT/dpkg-generated nondeterministic logs in the same APT `RUN` layer
2. Supply a stable `SOURCE_DATE_EPOCH` value
3. Export with BuildKit’s `rewrite-timestamp=true` option
4. Pin package sources or use Debian snapshots when reproducibility of package versions is also required

No APT option was found that replaces the same-layer removal of the generated logs. The runtime Dockerfile change is shared by all downstream stages; `node-base` requires no additional Dockerfile change.

## Remaining Work

- Test the Java-specific JDK APT layer after applying the same generated-state treatment
- Test the Web Chromium/font APT layer, including generated files such as linker caches and certificate-related data
- Validate the final CI Bake action configuration with `SOURCE_DATE_EPOCH`, `type=image,rewrite-timestamp=true`, and image loading/publishing requirements
- If exact package-version reproduction is required, test Debian and Adoptium snapshot/version pinning separately

## References

- [CI Docker Build Caching](004-note-ci-docker-build-caching.md)
- [OCI Layer Reproducibility and CI Cache Misses](005-note-oci-layer-reproducibility.md)
- [BuildKit reproducibility documentation](https://github.com/moby/buildkit/blob/master/docs/build-repro.md)
- [Docker reproducible builds with GitHub Actions](https://docs.docker.com/build/ci/github-actions/reproducible-builds/)
- [Docker Bake file reference](https://docs.docker.com/build/bake/reference/)
- [Dockerfile reference](https://docs.docker.com/reference/dockerfile)
- [APT configuration reference](https://manpages.debian.org/unstable/apt/apt.conf.5.en.html)
- [SOURCE_DATE_EPOCH guidance](https://reproducible-builds.org/docs/source-date-epoch/)
