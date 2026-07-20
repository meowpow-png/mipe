# Runtime

Runtime bootstrap prepares the shared agent environment, initializes a project when requested, and starts the chosen command. Docker images provide base and AI agent environments, with various toolchain variants.

## Building

Build every local runtime image:

```bash
just build-images
```

Build selected agent variants:

```bash
just build-images codex
just build-images claude-java
```

Build a reproducible OCI image layout or compare two no-cache builds:

```bash
just build-oci runtime runtime-oci
just build-compare runtime
```

See [Testing](docs/TESTING.md#image-reproducibility) for when to use these checks.

## Running

Start the agent and toolchain you need:

```bash
just codex
just claude
just codex-java
just claude-java
just codex-web
just claude-web
```

Open an initialized shell instead of starting the agent directly:

```bash
just codex mipe bash
just claude-java mipe bash
```

Compose provides equivalent service commands:

```bash
docker compose run --rm codex
docker compose run --rm claude-java
```

Run an arbitrary command through the bootstrap with:

```bash
mipe [flags] <command> [args...]
```

Common flag uses:

```bash
just mipe --version
just codex mipe --debug codex
just claude mipe --config /opt/mipe/config/config.json claude
```

| Flag              | Description                                           |
|-------------------|-------------------------------------------------------|
| `--config <path>` | Load bootstrap configuration from the specified file. |
| `--debug`         | Enable debug logging.                                 |
| `--version`, `-v` | Print the Mipe version and exit.                      |

## Testing

Unit tests cover bootstrap behavior in isolation. Integration tests validate the assembled runtime image and Linux environment.

```bash
just test -v
just test-coverage
just integration-test -v
```

Useful direct commands:

```bash
go test ./...
docker buildx bake --load --provenance=false --sbom=false test
MIPE_INTEGRATION=1 go test -v ./integration
```

## Configuration

Bootstrap configuration describes agent environment, developer identity, and project workspace used for each runtime invocation. Values may be supplied as top-level JSON fields or through the `environment` map, and process environment variables take precedence.

| Field          | Required | Description                                                                                             |
|----------------|---------:|---------------------------------------------------------------------------------------------------------|
| `environment`  |       No | Map of configuration environment values, including overrides such as `MIPE_DEBUG` and `MIPE_LOG_FORMAT` |
| `agent_name`   |      Yes | Identifies the agent being initialized.                                                                 |
| `user_home`    |      Yes | Local developer's home directory; supplied to the agent as `HOME`                                       |
| `agent_home`   |       No | Agent-specific persistent home directory                                                                |
| `runtime_home` |      Yes | Location of the shared Mipe runtime configuration                                                       |
| `workspace`    |      Yes | Existing writable project directory used for initialization and command execution                       |
| `local_uid`    |      Yes | Numeric user ID used for project initialization and command execution                                   |
| `local_gid`    |      Yes | Numeric group ID used for project initialization and command execution                                  |

Mipe supplies configuration for each runtime environment. You can supply another file with `--config`, but prefer environment-variable overrides for per-run changes so supplied configuration remains the shared default.

| Variable              | Description                                                                                      |
|-----------------------|--------------------------------------------------------------------------------------------------|
| `AGENT_NAME`          | Required agent identifier                                                                        |
| `USER_HOME`           | Required local developer home; Mipe exports it to child commands as `HOME`                       |
| `AGENT_HOME`          | Optional persistent home for agent configuration and state; exported when set                    |
| `RUNTIME_HOME`        | Required location of shared Mipe runtime configuration; exported to child commands               |
| `WORKSPACE`           | Required writable project directory used for initialization and command execution                |
| `LOCAL_UID`           | Required numeric user ID for initialization and command execution                                |
| `LOCAL_GID`           | Required numeric group ID for initialization and command execution                               |
| `MIPE_DEBUG`          | Optional boolean debug toggle                                                                    |
| `MIPE_LOG_FORMAT`     | Optional log format: `console` (default) or `json`                                               |

Useful overrides:

```bash
LOCAL_UID="$(id -u)" LOCAL_GID="$(id -g)" just codex
MIPE_DEBUG=true just codex mipe bash
MIPE_LOG_FORMAT=json just claude mipe claude
```

## Documentation

- [Implementation](docs/IMPLEMENTATION.md)
- [Local deployment](docs/DEPLOYMENT.md)
- [Testing](docs/TESTING.md)
