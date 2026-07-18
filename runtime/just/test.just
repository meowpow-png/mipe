# Run Go tests
[arg("args", help="Arguments passed to Go test")]
test *args:
    @go test {{args}} ./...

# Run Go tests with coverage
[arg("args", help="Arguments passed to Go test")]
test-coverage *args:
    @mkdir -p build/report
    @go test -coverprofile=build/report/coverage.out {{args}} ./...
    @go tool cover -func=build/report/coverage.out

# Run integration tests
[arg("args", help="Arguments passed to Go test")]
integration-test *args:
	@docker buildx bake --load --provenance=false --sbom=false test
	@MIPE_INTEGRATION=1 go test {{args}} ./integration/...
