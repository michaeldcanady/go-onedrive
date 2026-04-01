# go-onedrive (odc) Project Context

`go-onedrive` (cli name: `odc`) is a CLI tool designed to interact with OneDrive as a Unix-style file system, providing a terminal-native way to manage files (ls, cp, mv, rm, etc.).

## Project Overview

- **Purpose:** Provide a robust, coreutils-like experience for managing OneDrive files.
- **Main Technologies:**
    - **Runtime:** Go (>= 1.25.3)
    - **CLI Framework:** Cobra (inferred from `cmd/odc/main.go` and `spf13/cobra` dependency)
    - **API Client:** Microsoft Graph SDK for Go (`msgraph-sdk-go`)
    - **Serialization:** Kiota (`microsoft/kiota-*` dependencies)
    - **Logging:** Zap (`go.uber.org/zap`)
    - **State/Profile Persistence:** bbolt (`go.etcd.io/bbolt`)
    - **Task Automation:** Just (`justfile`)
    - **Documentation:** MkDocs (`mkdocs.yaml`)
    - **Cloud Auth:** Azure SDK (`azcore`, `azidentity`)

## Architecture

The project follows a modular design, with core functionalities organized within the `internal/feature/` directory. Key architectural components include:
- **`cmd/odc/`:** The entry point for the CLI application.
- **`internal/feature/`:** Contains core domain services and feature-specific logic (e.g., `fs`, `state`, `profile`, `drive`, `identity`).
- **`internal/di/`:** Manages the Dependency Injection container for wiring services.
- **`pkg/`:** Houses general-purpose utilities and shared components.
- **`docs/`:** Contains project documentation managed by MkDocs.

## Building and Running

- **Install Dependencies:** `go get ./...` (dependencies are listed in `go.mod`)
- **Build CLI:** `just build` (outputs `odc` binary in the root directory)
- **Run CLI:** `just run [args]`
- **Serve Documentation:** `just serve-docs` (previews documentation locally via MkDocs)
- **Build Static Docs:** `just generate-docs`
- **Clean:** `just clean` (removes build artifacts and docs)

## Development Conventions

- **Go Version:** Go 1.25+ (as specified in `go.mod`).
- **Error Handling:** Prefer wrapping errors with `%w` and using package-level error variables.
- **Logging:** Utilize the structured logging abstraction provided by `internal/logger`.
- **Formatting:** Adhere to Go standards; run `go fmt ./...` before committing.
- **Linting:** Run `golangci-lint run` to ensure code quality.
- **Testing:** Unit tests are mandatory for new features and bug fixes. Leverage `stretchr/testify` for assertions and mocking, and use `t.Parallel()` where applicable.
- **Pull Request Process:** Follow the guidelines in `CONTRIBUTING.md`, including forking, branching, building, testing, linting, documentation updates, and submitting PRs.

## Logging Configuration

- **Default Output:** Logs are directed to a file (`./logs/app.log`) by default.
- **Directory:** The `./logs` directory must exist for logging to function correctly. Ensure it is created before running the application.

