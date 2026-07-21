# Mipe

[![CI](https://github.com/meowpow-png/mipe/actions/workflows/ci.yml/badge.svg)](https://github.com/meowpow-png/mipe/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/meowpow-png/mipe/branch/dev/graph/badge.svg?token=MbG8tNgD2G)](https://codecov.io/github/meowpow-png/mipe)
![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
![OpenAI Codex](https://img.shields.io/badge/OpenAI-Codex-412991?logo=openai&logoColor=white)
![Anthropic Claude](https://img.shields.io/badge/Anthropic-Claude-D97757?logo=anthropic&logoColor=white)

There are only so many times you can copy the same AI development setup into a new repository before admitting it probably belongs somewhere else. That's how this repository happened.

## What is this?

The boring stuff nobody wants to rebuild for every AI-assisted project.

Imagine every project came with its own AI workspace instead of borrowing whatever happened to be installed on your machine.

Clone the repository, start the container, and your AI agent finds the tools, configuration, prompts, and project context waiting for it. No machine setup, no copy-pasted bootstrap scripts, no wondering what you forgot this time.

## Is this another AI agent?

No.

Mipe doesn't replace Codex, Claude, or whatever comes next. It's the environment around the agent, not the agent itself. Think of it like a kitchen. The AI agent does the cooking, Mipe makes sure the kitchen is ready when it walks in.

## Why should I use it?

Because your AI setup is part of your project, not your laptop.

Instead of one global environment shared across everything, each project gets its own tools, configuration, initialization, and runtime.

Open your project directory, start a container, and get to work.

Your host machine stays clean, and every project gets its own isolated environment. Install whatever tools the project needs, experiment freely, and throw the whole environment away when you're done. The next project starts with a clean slate.

## Available Images

The available runtime images are listed below.

| Image                                          | Tags                                 |
|------------------------------------------------|--------------------------------------|
| `ghcr.io/meowpow-png/mipe-runtime-codex`       | `latest`, `dev-latest`, `dev-<hash>` |
| `ghcr.io/meowpow-png/mipe-runtime-claude`      | `latest`, `dev-latest`, `dev-<hash>` |
| `ghcr.io/meowpow-png/mipe-runtime-codex-java`  | `latest`, `dev-latest`, `dev-<hash>` |
| `ghcr.io/meowpow-png/mipe-runtime-claude-java` | `latest`, `dev-latest`, `dev-<hash>` |
| `ghcr.io/meowpow-png/mipe-runtime-codex-web`   | `latest`, `dev-latest`, `dev-<hash>` |
| `ghcr.io/meowpow-png/mipe-runtime-claude-web`  | `latest`, `dev-latest`, `dev-<hash>` |

> [!NOTE] 
> Use `dev-<hash>` tag to pin a specific Mipe runtime build. The hash is shared across all image variants and identifies the runtime binary version, independent of bundled dependency versions.

Each variant includes `mipe-runtime` and the following components:

| Variant       | Components                                                            |
|---------------|-----------------------------------------------------------------------|
| `codex`       | `python3`, `git`, `codex`, `node`                                     |
| `claude`      | `python3`, `git`, `claude-code`, `node`                               |
| `codex-java`  | `python3`, `git`, `codex`, `node`, `temurin-21-jdk`                   |
| `claude-java` | `python3`, `git`, `claude-code`, `node`, `temurin-21-jdk`             |
| `codex-web`   | `python3`, `git`, `codex`, `node`, `chromium`, `playwright-mcp`       |
| `claude-web`  | `python3`, `git`, `claude-code`, `node`, `chromium`, `playwright-mcp` |

## Quickstart

**Requirements:**

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install)

The following example starts [Codex CLI](https://github.com/openai/codex) with Mipe.

Create `compose.yaml`:

```yaml
services:
  codex:
    image: ghcr.io/meowpow-png/mipe-runtime-codex:latest
    volumes:
      - codex-home:/home/dev/.codex
      - .:/workspace:Z
    stdin_open: true
    tty: true

volumes:
  codex-home:
```

Start Codex:

```bash
docker compose run --rm codex
```

If your host user isn't `1000`, pass your UID and GID:

```bash
docker compose run --rm \
    -e LOCAL_UID="$(id -u)" \
    -e LOCAL_GID="$(id -g)" \
    codex
```

Authenticate when prompted. The `codex-home` volume preserves your authentication, settings, and caches between runs.

> [!TIP]
> Replace `:latest` with `:dev-latest` to try the latest development build.

For runtime configuration, agent variants, and advanced usage, see [Runtime](runtime/README.md) documentation.

## Project setup

Projects can prepare themselves automatically before each AI session.

Create the following file in your project:

```text
.mipe/init/setup.sh
```

Mipe runs this script as the project user before starting the agent. It can install project dependencies, customize the agent environment, or generate files the agent needs. Keep the commands repeatable because the script runs when a session starts.

For example, a project can override the image defaults with your own global agent instructions and configuration:

```bash
setup_project() {
    [ -f /mipe/config/AGENTS.md ] &&
        install -Dm644 /mipe/config/AGENTS.md "$AGENT_HOME/AGENTS.md"

    [ -f /mipe/config/config.toml ] &&
        install -Dm644 /mipe/config/config.toml "$AGENT_HOME/config.toml"
}
```

This requires mounting your host `~/.mipe` directory into the container: 

```yaml
volumes:
  - ~/.mipe:/mipe:ro
```

The agent can then combine project-specific instructions from the workspace with your shared configuration. Use this to customize or override the default agent configuration provided by the image.

It can also generate project context before the agent starts:

```bash
setup_project() {
    npm run generate
    git ls-files > .mipe/project-files.txt
}
```

Use this when the agent depends on generated metadata, indexes, or other project context that should be refreshed for each session.

For dependency setup, a Node.js project image might run:

```bash
setup_project() {
    npm ci
}
```

Alternatively, a Java project using Maven Wrapper might run:

```bash
setup_project() {
    ./mvnw dependency:go-offline
}
```

Use these dependency hooks when sessions run in fresh or disposable environments, and the project should prepare itself before the agent starts.

## Documentation

- [Runtime](runtime/README.md) — Build, run, and configure Mipe
- [Concept](docs/CONCEPT.md) — Project motivation and vision
- [Roadmap](docs/ROADMAP.md) — Planned milestones and release goals
- [Examples](runtime/examples/README.md) — Complete example projects

## Contributing

Contributions are welcome. Let's be honest, if you found yourself here, you were probably about to change something anyway. You might as well open a pull request.

Start by reading the [Contributing](CONTRIBUTING.md) guide.

## License

Licensed under Apache License 2.0. See [LICENSE](LICENSE).
