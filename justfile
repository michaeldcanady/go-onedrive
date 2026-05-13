# Variables
export CONTAINER_TOOL := `if command -v podman >/dev/null 2>&1; then echo podman; else echo docker; fi`
export DOC_IMAGE := "sdk-docs"
export DOC_PORT := "8000"
export DOC_PATH := "."
export BINARY_NAME := "odc"

# Generate code
generate: generate-cli
    go generate ./...

# Generate CLI boilerplate from specs
generate-cli:
    go run ./cmd/spec-gen/main.go

# Run all tests
test:
    go test -v ./...

# Run all tests (alias for CI parity)
test-all: test

# Run unit tests (fast, no IO)
test-unit:
    go test -v ./...

# Run integration tests (domain interactions, mocks external APIs)
test-integration:
    go test -v -run Integration ./...

# Run functional tests (feature slice validation, DI wired)
test-functional:
    go test -v -run Functional ./...

# Run performance benchmarks
test-perf:
    go test -v -bench=. -run=^$ ./...

# Run E2E tests (binary against environment)
test-e2e: build
    go test -v -tags=e2e ./...

# Run a quick smoke test on the built binary
test-smoke: build
    ./{{BINARY_NAME}} --version

# Run linter
lint:
    golangci-lint run ./...

# Lint a specific package (optimized for hooks)
lint-pkg pkg:
    golangci-lint run {{pkg}}

# Run security checks
secure:
    govulncheck ./...

# Build the odc binary
build: build-cli build-plugins

# Build the odc cli
build-cli: generate
    go build -o {{BINARY_NAME}} ./cmd/odc/

# Build all storage and identity plugins
build-plugins: generate
    mkdir -p bin/plugins
    go build -o bin/plugins/storage-onedrive ./cmd/storage-plugin-onedrive/
    go build -o bin/plugins/storage-local ./cmd/storage-plugin-local/
    go build -o bin/plugins/storage-googledrive ./cmd/storage-plugin-googledrive/
    go build -o bin/plugins/identity-azure ./cmd/identity-plugin-azure/
    go build -o bin/plugins/identity-google ./cmd/identity-plugin-google/

# Run the odc binary
run *args: build
    ./{{BINARY_NAME}} {{args}}

# Clean up build artifacts and docs
clean: clean-docs
    rm -f {{BINARY_NAME}}

# Build the docs container image
build-docs:
    {{CONTAINER_TOOL}} build -f doc.dockerfile -t {{DOC_IMAGE}}

# Serve docs locally with live reload
serve-docs: build-docs
    {{CONTAINER_TOOL}} run --rm -p {{DOC_PORT}}:8000 {{DOC_IMAGE}}

# Build static site output (without serving)
generate-docs: build-docs
    {{CONTAINER_TOOL}} run --rm {{DOC_IMAGE}} mkdocs build --clean

# Clean up local build artifacts
clean-docs:
    rm -rf site

# Generate man pages
generate-man dir="./man": build
    ./{{BINARY_NAME}} docs man {{dir}}
