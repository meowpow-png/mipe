# Run base Docker Compose stack
[arg("args", help="Arguments passed to Docker Compose")]
compose *args:
    @docker compose {{args}}
