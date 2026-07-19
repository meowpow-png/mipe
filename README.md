# Mipe

## Getting Started

### Project Dependencies

To configure project initialization, create the following file:

```text
.mipe/init/dependencies.sh
```

The script must define an `install_dependencies` function containing the commands required to prepare this project before development begins.

Example:

```bash
#!/usr/bin/env bash
set -euo pipefail

install_dependencies() {
    go mod download
}
```

### Development Builds

Development builds are available through GitHub Container Registry (GHCR).

> [!WARNING]
> Development builds are intended for testing, evaluation, and early adoption of new functionality. They may introduce breaking changes at any time and should not be considered stable.

## Deployment

The project workspace must be writable by the configured local developer user. 
Project initialization executes with the developer user's permissions and requires 
write access to the workspace to install dependencies and generate project files.

### Running Mipe

Run the locally built bootstrap without rebuilding a container image:

```bash
just mipe bash
```

Use the container-backed runtime when validating image behavior:

```bash
just mipe-docker bash
```
