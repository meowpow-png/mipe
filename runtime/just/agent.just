# Run Codex in Docker container
[arg("args", help="Command arguments")]
codex *args:
    @just compose run --rm codex {{args}}

# Run Claude in Docker container
[arg("args", help="Command arguments")]
claude *args:
    @just compose run --rm claude {{args}}

# Run Codex with Java toolchain in Docker container
[arg("args", help="Command arguments")]
codex-java *args:
    @just compose run --rm codex-java {{args}}

# Run Claude with Java toolchain in Docker container
[arg("args", help="Command arguments")]
claude-java *args:
    @just compose run --rm claude-java {{args}}

# Run Codex with web toolchain in Docker container
[arg("args", help="Command arguments")]
codex-web *args:
    @just compose run --rm codex-web {{args}}

# Run Claude with web toolchain in Docker container
[arg("args", help="Command arguments")]
claude-web *args:
    @just compose run --rm claude-web {{args}}
