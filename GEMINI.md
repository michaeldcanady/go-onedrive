# go-onedrive (odc) Project Context

`go-onedrive` (cli name: `odc`) is a CLI tool designed to interact with OneDrive as a Unix-style file system, providing a terminal-native way to manage files (ls, cp, mv, rm, etc.).

## Project Overview

- **Purpose:** Provide a robust, coreutils-like experience for managing OneDrive files.
- **Main Technologies:**
    - **Runtime:** Go (>= 1.25.3)
    - **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
    - **API Client:** [Microsoft Graph SDK for Go](https://github.com/microsoftgraph/msgraph-sdk-go)
    - **Serialization:** [Kiota](https://github.com/microsoft/kiota-abstractions-go)
    - **Logging:** [Zap](https://github.com/uber-go/zap)
    - **Caching:** [bbolt](https://github.com/etcd-io/bbolt)
- **Architecture:** The project follows a modular design centered around dependency injection (`internal/di`) and slice-based command organization.
    - **`cmd/odc/`:** Entry point for the application.
    - **`internal/core/`:** Contains core domain services (auth, drive, fs, logger, profile, state).
        - **`fs/`:** Provides a filesystem abstraction for both local and OneDrive providers.
    - **`internal/slices/`:** Contains the CLI command implementations (auth, drive, profile, root).
    - **`pkg/`:** General-purpose utilities (cache, ignore patterns).
    - **`internal/di/`:** Dependency Injection container management.
    - **`internal/middleware/`:** Custom Graph API middleware (correlation IDs, logging).

## Building and Running

- **Install Dependencies:** `go get ./...`
- **Build CLI:** `just build` (outputs `odc` binary in root)
- **Run CLI:** `just run [args]`
- **Documentation:** `just serve-docs` (previews documentation locally via MkDocs)
- **Clean:** `just clean`

## Testing and Quality

- **Unit Tests:** `go test ./...`
- **Linting:** `golangci-lint run`
- **Formatting:** `go fmt ./...`
- **Conventions:**
    - Prefer table-driven tests for complex logic.
    - Mock external dependencies (Graph API, identity providers) using `stretchr/testify`.
    - Use structured logging via the `internal/core/logger` abstraction.

## Specialized Skills

The project uses several specialized Gemini skills located in `.gemini/skills/`:

- **`code-writer`:** Use for any Go source code changes. Follows specific templates for CLI commands and enforces strict engineering standards (errors, naming, DI).
- **`docs-writer`:** Use for all documentation changes (MkDocs in `docs/` and `.md` files). Enforces tone, formatting (80-char wrap), and structure (BLUF).
- **`principal-software-engineer`:** Use for major architectural decisions, system-wide refactors, or complex debugging.
- **`product-manager`:** Use for defining vision, strategy, and roadmaps.

## Development Conventions

- **CLI Commands:** All new commands must follow the template pattern in the `code-writer` skill (separating command factory, logic, and options).
- **Architecture:** Maintain clear boundaries between domain services (`core`) and CLI orchestration (`slices`).
- **Errors:** Return errors as the last return value and wrap them using `%w`.
- **Documentation:** Every heading in `.md` files should be followed by an introductory paragraph. Sentence case is used for all headings.
