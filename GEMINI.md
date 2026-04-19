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

The project follows a modular design, with core functionalities organized within the `internal/` directory. Key architectural components include:
- **`cmd/odc/`:** The entry point for the CLI application.
- **`internal/`:** Contains core domain services and logic (e.g., `core/fs`, `profile`, `drive`, `identity`, `mount`).
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

## Release Process

`odc` uses **Release Please** for versioning and changelog generation, and **GoReleaser** for automated releases and distribution.

### Automated Release Lifecycle

1.  **Pull Request:** When changes are pushed to `main`, **Release Please** automatically creates or updates a Release PR.
2.  **Versioning:** The version is determined by Conventional Commits.
3.  **Stages:** The release cycle follows `alpha` -> `beta` -> `rc` -> `full`.
    - To transition between stages, update the `prerelease-type` in `release-please-config.json`.
    - For a full release, set `"prerelease": false` in `release-please-config.json`.
4.  **Tagging:** Merging the Release PR triggers **Release Please** to create a GitHub Release and a corresponding tag (e.g., `v1.2.3-alpha.1`).
5.  **Distribution:** The new tag triggers the **Release** workflow, which runs **GoReleaser** to:
    - Build binaries for Linux, macOS, and Windows.
    - Generate `.deb`, `.rpm`, and `.apk` packages.
    - Publish to the **Homebrew Tap** (`michaeldcanady/homebrew-tap`) (skipped for pre-releases).
    - Publish to **Cloudsmith** for Linux package distribution.
    - Submit a manifest to **WinGet** (`microsoft/winget-pkgs`).

### Manual Tagging (Fallback)
1.  **Tag a new version:** `git tag -a vX.Y.Z -m "Release message"`
2.  **Push the tag:** `git push origin vX.Y.Z`
3.  The GitHub Action will automatically:
    - Build binaries for Linux, macOS, and Windows.
    - Generate `.deb`, `.rpm`, and `.apk` packages.
    - Publish to the **Homebrew Tap** (`michaeldcanady/homebrew-tap`).
    - Publish to **Cloudsmith** for Linux package distribution.
    - Submit a manifest to **WinGet** (`microsoft/winget-pkgs`).

### Prerequisites for First Release
Before the first release, ensure the following are configured:
1.  **Homebrew:** Create a repository named `homebrew-tap` in your GitHub account.
2.  **Cloudsmith:**
    - Create a free account at [Cloudsmith.com](https://cloudsmith.com/).
    - Create a repository named `odc`.
    - Add `CLOUDSMITH_API_KEY` to GitHub Secrets.
3.  **WinGet:**
    - Create a Personal Access Token (PAT) with `repo` scope.
    - Add `WINGET_GITHUB_TOKEN` to GitHub Secrets.

## Development Conventions

- **Go Version:** Go 1.25+ (as specified in `go.mod`).
- **Error Handling:** Prefer wrapping errors with `%w` and using package-level error variables.
- **Logging:** Utilize the structured logging abstraction provided by `internal/logger`.
- **Formatting:** Adhere to Go standards; run `go fmt ./...` before committing.
- **Linting:** Run `golangci-lint run` to ensure code quality.
- **CLI Command Validation:**
    - **Fail Fast:** Use Cobra's `PreRunE` to populate the command's `Options` struct and perform validation.
    - **Options Structs:** Every command should have an `Options` struct with a `Validate() error` method to centralize validation logic.
- **State vs. Configuration Management:**
    - **Session State:** Managed in-memory within domain services (e.g., `profile.Service` handles session-active profiles).
    - **Config (`internal/config`):** Represents user-defined settings (e.g., Azure Tenant ID, Redirect URIs, Mount Points). External-facing and user-modifiable.
    - **Domain-First Access Pattern:** Commands (UI layer) MUST NOT access internal persistence directly. They must use Domain Services (e.g., `profile.Service`, `drive.Service`, `mount.Service`) which encapsulate state interaction.
    - **Scoping:** The UI layer determines the scope of changes (e.g., `shared.ScopeGlobal` for persistence, `shared.ScopeSession` for transience) and passes it to the Domain Service methods.
- **Testing:** Unit tests are mandatory for new features and bug fixes. Leverage `stretchr/testify` for assertions and mocking, and use `t.Parallel()` where applicable.
- **Pull Request Process:** Follow the guidelines in `CONTRIBUTING.md`, including forking, branching, building, testing, linting, documentation updates, and submitting PRs.

## Logging Configuration

- **Default Output:** Logs are directed to a file (`./logs/app.log`) by default.
- **Directory:** The `./logs` directory must exist for logging to function correctly. Ensure it is created before running the application.

