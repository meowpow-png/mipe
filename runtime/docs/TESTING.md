# Testing

The runtime has two test suites:

| Suite       | What it checks                                          | Command                    |
|-------------|---------------------------------------------------------|----------------------------|
| Unit        | Go packages and bootstrap behavior in isolation         | `just test`                |
| Integration | The assembled runtime image and its Linux environment   | `just integration-test -v` |

Regular Go test runs do not start Docker. Integration tests only run when they are explicitly enabled.

## Unit Tests

Unit tests cover configuration loading, validation, bootstrap phases, file operations, process execution, and command-line behavior. They are the quickest way to check a change while working.

Run the full unit test suite with:

```bash
just test
```

Arguments after the recipe name are passed to `go test`, so normal Go options work as expected:

```bash
just test -v
just test -race
```

You can also run the underlying command directly:

```bash
go test ./...
```

## Coverage

Run the unit tests with coverage using:

```bash
just test-coverage
```

This prints a coverage summary and writes the full profile to `build/report/coverage.out`.

Integration tests run the runtime as a separate binary inside Docker, so they do not add to the Go coverage profile. Use unit tests to track line coverage and integration tests to check real container behavior.

## Integration Tests

Integration tests cover the parts that are difficult to reproduce in a normal Go test. This includes the container entrypoint, user and group IDs, file ownership, runtime configuration, project initialization, environment variables, working directories, and final command execution.

The suite builds the test image, generates fresh fixtures, and starts a separate container for each scenario. It captures the output and exit status, then removes the container when the test finishes.

The scenarios cover:

- Successful runtime startup
- Missing and failing initialization scripts
- Invalid configuration and environment overrides
- Missing or unwritable workspaces
- Ownership and persisted agent state
- Final command exit codes
- Missing or conflicting container identities

Docker must be running, with Buildx available. Run the suite in verbose mode to see image build progress and each test as it runs:

```bash
just integration-test -v
```

The recipe builds and loads `mipe-runtime-test:latest` before running the suite. If running the Go command directly, build the test image first:

```bash
docker buildx bake --load --provenance=false --sbom=false test
```

To run a specific test:

```bash
just integration-test -v -run TestInvalidConfiguration
```

The direct Go command is:

```bash
MIPE_INTEGRATION=1 go test -v ./integration
```

## Daily Workflow

Run focused unit tests while developing, then run the full unit suite before finishing a change:

```bash
just test
```

Also run the integration suite when changing the runtime image, entrypoint, bootstrap lifecycle, configuration, permissions, environment, or initialization behavior:

```bash
just integration-test -v
```
