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

The recipe builds and loads `mipe-runtime-test:local` before running the suite. If running the Go command directly, build the test image first:

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

## Image Reproducibility

When changing image construction, verify that independent no-cache builds produce same image. This catches nondeterministic files, timestamps, and metadata before they cause CI cache misses or republish unchanged layers.

Build an OCI image layout for inspection:

```bash
just build-oci runtime runtime-oci
```

Layout is written under `.mipe/`. Recipe uses configured fixed epoch and timestamp rewriting, so OCI manifests and layers are suitable for reproducibility checks.

For normal validation, run an automated full-image comparison with `diffoci`. Use affected Bake target in place of `runtime`; any difference is a reproducibility regression.

```bash
just build-compare runtime
```

Compare exported OCI layouts when checking files and metadata. Build same target twice with unchanged inputs, then compare both layouts.

```bash
just build-oci runtime runtime-a
just build-oci runtime runtime-b
diff -ru .mipe/runtime-a .mipe/runtime-b
```

Inspect layer identifiers when narrowing down a difference. Build two tagged local images, then compare their layer lists.

```bash
VERSION="$(bash scripts/get-go-build-version.sh)" \
  docker buildx bake runtime --load --provenance=false --sbom=false \
    --set runtime.tags=runtime-a --no-cache
VERSION="$(bash scripts/get-go-build-version.sh)" \
  docker buildx bake runtime --load --provenance=false --sbom=false \
    --set runtime.tags=runtime-b --no-cache
diff -u \
  <(docker image inspect runtime-a --format '{{json .RootFS.Layers}}') \
  <(docker image inspect runtime-b --format '{{json .RootFS.Layers}}')
```

Matching layer lists are a quick signal, not a complete image comparison. Use `just build-compare` or OCI-layout comparison before accepting a change.

## Daily Workflow

Run focused unit tests while developing, then run the full unit suite before finishing a change:

```bash
just test
```

Also run the integration suite when changing the runtime image, entrypoint, bootstrap lifecycle, configuration, permissions, environment, or initialization behavior:

```bash
just integration-test -v
```
