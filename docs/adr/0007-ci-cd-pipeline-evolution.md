# Design Document: CI/CD Pipeline & Automation Evolution

## 1. Introduction
This document outlines the strategy for enhancing the CI/CD pipelines and automation scripts for the `go-onedrive` (odc) project. The goal is to improve security, ensure performance stability, and support the upcoming transition to a plugin-based architecture with OpenTelemetry integration.

## 2. Current State Assessment

### 2.1 Existing Infrastructure
- **CI/CD**: GitHub Actions for testing, linting, and releases.
- **Release Management**: Release Please and GoReleaser.
- **Automation**: `justfile` for basic build and documentation tasks.
- **Testing**: Unit tests and some integration tests running on Linux, macOS, and Windows.

### 2.2 Identified Gaps
- **Security**: Lack of automated vulnerability scanning (Go-specific), secret scanning, and SBOM generation.
- **Performance**: Benchmarks exist but are not executed or tracked in CI.
- **E2E Testing**: No full-scenario E2E tests or validation of distributed packages.
- **Plugin Support**: Pipeline is not yet configured to handle multi-binary plugin builds.
- **Telemetry**: No validation of OpenTelemetry instrumentation.

## 3. Proposed Improvements

### 3.1 Security Guardrails
We will implement a "Security-First" approach by adding the following to the CI pipeline:
- **`govulncheck`**: Integrate into the `go.yaml` workflow to detect known vulnerabilities in dependencies.
- **Secret Scanning**: Add `gitleaks` or a similar action to prevent accidental credential leakage.
- **SBOM Generation**: Configure GoReleaser to generate SPDX/CycloneDX SBOMs for all releases.
- **CodeQL**: Enable GitHub CodeQL analysis for deep static analysis.

### 3.2 Performance Benchmarking Pipeline
- **Continuous Benchmarking**: Add a new job to `go.yaml` that runs `go test -bench`.
- **Regression Detection**: Use `github-action-benchmark` to store results and provide visual feedback on performance changes over time.
- **Performance Thresholds**: Define critical paths (e.g., VFS resolution) and fail builds if performance degrades beyond a set percentage.

### 3.3 Enhanced E2E and Multi-Platform Validation
- **Scenario Testing**: Develop a suite of E2E tests using a mock server (or a test OneDrive account with secrets) that mimics real-world usage.
- **Smoke Tests for Packages**: Add a job to validate that the `.deb`, `.rpm`, and `.apk` packages produced by GoReleaser can be installed and run basic commands in a containerized environment.
- **Windows Release Restoration**: Re-enable and validate Windows builds in `.goreleaser.yaml`.

### 3.4 Support for Plugin-Based Architecture
As the project moves to a plugin-based model (e.g., `storage-plugin-onedrive`):
- **Matrix Builds for Plugins**: Update the build job to iterate over all directories in `cmd/` to ensure all plugins are built and tested.
- **Plugin Compatibility Tests**: Create a test suite that verifies the core CLI can load and communicate with plugins using the defined RPC/proto interface.

### 3.5 OpenTelemetry Integration Validation
- **Span Verification**: Implement a "telemetry smoke test" that runs a CLI command and verifies that the expected OTel spans are produced (using a local collector or OTLP-compatible mock).
- **Metric Tracking**: Ensure that basic metrics (command execution time, error rates) are consistently reported.

### 3.6 Automation Parity (`justfile`)
Update the `justfile` to include commands that mirror CI steps, allowing developers to catch issues early:
- `just lint`: Runs `golangci-lint`.
- `just secure`: Runs `govulncheck` and `gitleaks`.
- `just bench`: Runs benchmarks locally.
- `just test-all`: Runs unit, integration, and E2E tests.

### 3.7 Documentation Integrity
- **Man Page Validation**: Add a CI step to ensure that `just generate-man` succeeds and doesn't produce empty files.
- **MkDocs Build Check**: Integrate `just generate-docs` into the PR workflow to prevent documentation build regressions.

## 4. Implementation Roadmap

### Phase 1: Security & Automation Hardening (Immediate)
- Update `justfile` with `lint`, `secure`, and `test-all`.
- Add `govulncheck` and `gitleaks` to `go.yaml`.
- Enable SBOM generation in `.goreleaser.yaml`.

### Phase 2: Performance & Package Validation (Short-term)
- Implement `github-action-benchmark` in CI.
- Add containerized smoke tests for Linux packages.
- Restore Windows builds in GoReleaser.

### Phase 3: Plugin & Telemetry Support (Mid-term)
- Refactor CI to support multi-plugin builds.
- Add E2E scenario tests for plugin interactions.
- Integrate OTel span validation in CI.

## 5. Conclusion
By implementing these improvements, the `go-onedrive` project will achieve a high level of operational excellence, ensuring that it remains secure, performant, and reliable as it scales in complexity.
