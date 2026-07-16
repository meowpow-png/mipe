# Mipe

## Project Dependencies

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

## Deployment

The project workspace must be writable by the configured local developer user. 
Project initialization executes with the developer user's permissions and requires 
write access to the workspace to install dependencies and generate project files.
