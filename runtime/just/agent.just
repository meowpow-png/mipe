# Run Codex in a disposable container
[arg("args", help="Arguments passed to container")]
codex-run *args:
    @just compose run --rm codex {{args}}

# Start a Bash shell in Codex container
[arg("args", help="Arguments passed to Bash")]
codex-shell *args:
    @just compose run --rm codex mipe bash {{args}}

# Run Claude in a disposable container
[arg("args", help="Arguments passed to container")]
claude-run *args:
    @just compose run --rm claude {{args}}

# Start a Bash shell in Claude container
[arg("args", help="Arguments passed to Bash")]
claude-shell *args:
    @just compose run --rm claude mipe bash {{args}}

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

# Run Codex Web in a disposable container
[arg("args", help="Arguments passed to container")]
codex-web-run *args:
    @just compose run --rm codex-web {{args}}

# Start a Bash shell in Codex Web container
codex-web-shell:
    @just compose run --rm --entrypoint bash codex-web

# Run Claude Web in a disposable container
[arg("args", help="Arguments passed to container")]
claude-web-run *args:
    @just compose run --rm claude-web {{args}}

# Start a Bash shell in Claude Web container
claude-web-shell:
    @just compose run --rm --entrypoint bash claude-web
