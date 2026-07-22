# Note: Bootstrap Configuration and Workspace Resolution

## Context

While decoupling bootstrap configuration from the agent runtime environment, I found that some runtime paths had multiple sources of truth.

Mipe relied on environment variables to locate its own configuration, while the configured workspace controlled initialization without determining where the agent actually started. The goal was to make `config.json` the primary source of runtime configuration, allow intentional environment overrides, and apply the resolved configuration consistently.

## Observations

The bootstrap used `HOME` both as its own configuration input and as the value passed to the agent. This meant the container's `HOME` could unintentionally override the configured home directory.

Configuration discovery had a similar problem. Without `--config`, Mipe looked for `$RUNTIME_HOME/config.json`. Removing `RUNTIME_HOME` from Docker Compose therefore prevented Mipe from finding the file that defined `runtime_home`.

Workspace handling was also inconsistent. The configured `workspace` was used to locate initialization scripts, but dependency initialization and the agent inherited the container's working directory. For example, configuring `workspace` as `/work` while the image used `WORKDIR /workspace` caused Mipe to initialize one directory and launch the agent from another.

Changing the configured workspace never moved or copied project files. Docker volume mounts remained a deployment responsibility.

## Analysis

Bootstrap configuration is now resolved before constructing the agent environment.

The bootstrap now reads `USER_HOME` from `config.json`, with an optional `USER_HOME` environment override. The process `HOME` is no longer used during configuration. When launching the agent, Mipe exports the resolved `USER_HOME` as `HOME` and does not expose `USER_HOME`.

Default configuration discovery is now independent of `RUNTIME_HOME`. Unless `--config` is provided, Mipe loads `/opt/mipe/config.json`. This removes the circular dependency where `RUNTIME_HOME` was needed just to find the configuration that defined it.

The configured workspace is now the single source of truth for execution. Mipe requires the directory to already exist and be writable, runs dependency initialization from that directory, and launches the agent from it. It does not copy or relocate project files.

The Dockerfile `WORKDIR` remains a sensible default, but once Mipe resolves the workspace, the configured path becomes authoritative.

## Conclusions

Bootstrap configuration and the agent runtime are now separate concerns. `USER_HOME` is used only by the bootstrap, while the agent receives a standard `HOME` environment variable.

The default configuration file is always loaded from `/opt/mipe/config.json` unless an explicit `--config` path is provided. `RUNTIME_HOME` is no longer required to discover configuration.

The configured workspace now controls both initialization and agent execution. The project must already be mounted at that location, and Mipe validates that the directory exists and is writable.

Environment variables may still override values from `config.json`, but this is now an intentional configuration mechanism rather than a deployment requirement.

## Next Steps

- Keep project volume mount aligned with the configured `workspace`. If the workspace changes from `/workspace` to `/work`, update the mount destination accordingly
- Keep `/opt/mipe/config.json` in runtime images unless an explicit `--config` path is supplied

No further runtime changes are required for this design. If workspace creation is ever needed, it should be treated as a separate provisioning step rather than part of bootstrap initialization.

## References

- [Codex Home Directory Initialization](001-note-home-directory-initialization.md)
