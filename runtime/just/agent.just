# Run Codex in a disposable container
[arg("args", help="Arguments passed to container")]
codex-run *args:
    @just compose run --rm codex {{args}}

# Start a Bash shell in Codex container
codex-shell:
    @just compose run --rm --entrypoint bash codex

# Run Claude in a disposable container
[arg("args", help="Arguments passed to container")]
claude-run *args:
    @just compose run --rm claude {{args}}

# Start a Bash shell in Claude container
claude-shell:
    @just compose run --rm --entrypoint bash claude

# Run Codex Java in a disposable container
[arg("args", help="Arguments passed to container")]
codex-java-run *args:
    @just compose run --rm codex-java {{args}}

# Start a Bash shell in Codex Java container
codex-java-shell:
    @just compose run --rm --entrypoint bash codex-java

# Run Claude Java in a disposable container
[arg("args", help="Arguments passed to container")]
claude-java-run *args:
    @just compose run --rm claude-java {{args}}

# Start a Bash shell in Claude Java container
claude-java-shell:
    @just compose run --rm --entrypoint bash claude-java
