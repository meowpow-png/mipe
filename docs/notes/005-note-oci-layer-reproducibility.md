# Note: OCI Layer Reproducibility and CI Cache Misses

## Context

I investigated why identical BuildKit builds can miss Dockerfile instruction caches in CI. The investigation was performed locally to separate OCI layer reproducibility from BuildKit cache import/export behavior. The focus was the `java-base` and `web-base` portion of `runtime/docker/Dockerfile`, where CI had shown unreliable reuse.

## Observations

I built the `codex-java` and `codex-web` targets twice each with cache disabled, using the same source tree, Dockerfile, build arguments, and pinned version arguments. Each build was exported as an OCI image archive.

The top-level OCI index and image manifest differed between repeated builds. The Debian base layer and the copied Node source layer were identical. The first differing layer in both target families was the runtime package-install instruction:

```dockerfile
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        gosu \
        passwd && \
    rm -rf /var/lib/apt/lists/*
```

The Java-specific JDK installation layer and the web-specific Chromium/font installation layer also differed. Later layers, including the Codex npm installation and configuration copies, consequently differed because they were based on changed parent filesystem snapshots.

The differing APT layer archives contained the same primary package files and file sizes, but directories, locks, package status files, and logs had build-time mtimes. The web package layer showed the same behavior across installed system files. Some later package-install layers also contained byte differences in generated files such as package logs, certificate data, and linker caches.

The OCI image metadata contained build-time creation annotations as well. These explain differing image metadata, but they do not explain the first filesystem-layer divergence.

## Analysis

I first verified the build inputs and used representative targets that exercise the unstable Java and Web stages. For each target, I produced two independent no-cache OCI archives rather than comparing BuildKit progress output or cache records.

I read only the OCI index and image manifest metadata initially, then compared the ordered compressed layer digest and size lists. I read each image configuration only after a layer mismatch was found. The configuration history supplied the `created_by` descriptions and allowed each layer position to be mapped back to its Dockerfile instruction; empty history entries for `ARG`, `ENV`, `ENTRYPOINT`, `CMD`, and `WORKDIR` were accounted for when aligning history with filesystem layers.

I then extracted the first differing compressed layer from each archive and compared archive listings. This showed the same paths and package payload sizes with different wall-clock metadata on the runtime APT layer. I repeated the listing comparison for the Web package-install layer and found the same timestamp pattern. Finally, I compared extracted later package layers while ignoring filesystem metadata. That showed that some generated files also differed in bytes, including package logs, certificate-related data, and linker caches; these differences occur in package-install steps after the initial APT divergence and should not be treated as evidence that the base image or BuildKit cache manifest changed.

This establishes that the cache misses are reproducibility-related: the package-install instructions produce different filesystem snapshots across otherwise identical executions. A cache manifest can make prior records available, but it cannot make a record match when the instruction's parent snapshot or generated output is different.

## Conclusions

The primary cause is nondeterministic filesystem state produced by the APT-containing `RUN` instructions, especially wall-clock mtimes on package-manager-created files and directories. Some package steps also generate content whose bytes vary, such as logs and system caches.

The first confirmed divergent instruction is the shared runtime APT installation. The Java JDK installation and Web Chromium/font installation introduce additional independent nondeterminism. Later npm and configuration layers are affected transitively; the local evidence does not show that npm dependency resolution is the first cause of the cache misses.

The pinned Debian and Node inputs were not the cause of the immediate repeated-build divergence. Image creation metadata differs too, but the filesystem-layer comparison demonstrates that the issue exists before final image metadata is assembled.

## Next Steps

When applying a fix in CI, first make the shared runtime APT step reproducible or isolate its generated state, then verify the Java and Web package-install steps separately. Repeat the paired OCI comparison after each change and confirm that the first divergent layer has become identical before evaluating BuildKit cache reuse in CI.

## References

- [CI Docker Build Caching](004-note-ci-docker-build-caching.md)
