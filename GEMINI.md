# GEMINI.md

This file provides project-specific mandates and guidance for Gemini CLI when working on the `go-onedrive` project. These instructions take precedence over general defaults.

## Architectural Mandates

- **Three-Layer Pattern:** Strictly adhere to the Domain, Application, and Infrastructure layers in `internal2`.
  - **Domain:** No external imports (e.g., standard libraries or third-party APIs). Defines interfaces.
  - **App:** Orchestrates domain logic and implements service interfaces.
  - **Infra:** Implements technical details (DB, API, SDKs).
- **Dependency Injection:** Use the central `Container` in `internal2/app/di/container.go` for all service lifecycles and wiring.
  - Services must be lazily initialized using `sync.Once`.
  - Prefer constructor injection (e.g., `NewService(dep1, dep2)`) in the container.
- **Interface-First:** Always define or update interfaces in the `domain` layer before modifying implementations.

## Engineering Standards

- **Go 1.25:** Leverage Go 1.25 features where appropriate.
- **Error Handling:**
  - Define errors as package-level variables in `errors.go`.
  - Use `fmt.Errorf("...: %w", err)` for wrapping.
  - Use `errors.Join` for multiple errors.
- **Logging:**
  - Use the structured `logging.Logger` interface from `internal2/infra/common/logging`.
  - Always use `logger.WithContext(ctx)` to ensure correlation IDs are propagated.
- **Testing:**
  - **Parallelism:** Always use `t.Parallel()`.
  - **Test Packages:** Use the `_test` package suffix (e.g., `package auth_test`).
  - **Assertions:** Prefer `stretchr/testify/require` for setup and `assert` for outcomes.
  - **Mocks:** Use `testify/mock`. Define mocks in the test files or `mocks_test.go`.

## Documentation

- **Diátaxis Framework:** All documentation in `docs/` must follow the Diátaxis structure (Tutorials, How-to, Explanation, Reference).
- **Mermaid:** Use Mermaid diagrams for complex interactions or architecture overviews.
- **Internal Reference:** Keep `docs/developer/reference/domain-interfaces.md` updated when adding new core services.

## Validation Routine

Before marking a task as complete:
1. Run unit tests for the affected package: `go test -v ./internal2/app/your_package/...`
2. Run the full test suite: `go test ./...`
3. Check linting: `golangci-lint run ./internal2/...` (if available in the environment)
4. Verify build: `go build ./cmd/odc/main.go`
