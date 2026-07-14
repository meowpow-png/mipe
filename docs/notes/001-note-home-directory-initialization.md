# Note: Codex Home Directory Initialization

## Context

While implementing the shared Devkit runtime, I wanted Codex to use a persistent home directory mounted at `/home/codex`. The goal was to keep the runtime configuration, authentication state, and caches inside a durable Docker volume that survives container recreation. During testing, however, Codex ignored the prepared configuration and reported its home directory as `/home/ubuntu`.

## Observations

The runtime successfully initialized `/home/codex/.codex/config.toml` and the directory was correctly owned by the local developer after startup. Debugging confirmed that immediately before executing Codex, the process environment contained:

- `HOME=/home/codex`
- `CODEX_HOME=/home/codex/.codex`

Despite this, once running, Codex reported its home directory as `/home/ubuntu` and attempted to locate its configuration there instead of under `/home/codex/.codex`.

## Analysis

I first verified that the runtime initialization behaved correctly by inspecting the filesystem immediately before the final `exec`. The expected configuration files were present under `/home/codex/.codex` and ownership was correct.

Next, I isolated the execution path by replacing the final `exec` with an interactive shell. The shell correctly inherited `HOME=/home/codex`, confirming that `gosu` preserved the environment and was not responsible for the behavior.

The issue was ultimately traced to the container environment itself. The Compose configuration did not define the `HOME` environment variable, relying instead on the entrypoint to export it. Defining `HOME=/home/codex` directly in the container environment immediately resolved the problem, allowing Codex to discover its configuration under `~/.codex` as intended.

## Conclusions

The container environment, not the runtime initialization, is the source of truth for the process home directory.

Although the entrypoint exported `HOME`, this was insufficient to guarantee consistent behavior for all processes. Defining `HOME=/home/codex` in Docker Compose ensured that every process launched inside the container, including Codex, consistently resolved its home directory to the persistent runtime volume.

The runtime now derives `CODEX_HOME` from `HOME` rather than defining the home directory itself.

## Next Steps

Keep `HOME=/home/codex` as part of the runtime container configuration for all Devkit images.

Continue treating the entrypoint as a consumer of the runtime environment rather than the component responsible for defining it. Future runtime paths should be derived from `HOME` instead of hardcoding `/home/codex`.

## References

-[ARCHITECTURE](../ARCHITECTURE.md)
