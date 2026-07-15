# Run Docker Compose
[arg("args", help="Arguments passed to docker compose")]
compose *args:
    @docker compose {{args}}

# Run Docker Compose for the Java runtime
[arg("args", help="Arguments passed to docker compose")]
compose-java *args:
    @just compose {{args}} java

# Run Docker Compose for the Node runtime
[arg("args", help="Arguments passed to docker compose")]
compose-node *args:
    @just compose {{args}} node

# Remove all Mipe images
docker-clean:
    @docker image rm -f \
        mipe-codex:21-java \
        mipe-codex:26-node 2>/dev/null || true
