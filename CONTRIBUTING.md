# Contributing to go-onedrive

Thank you for your interest in contributing to the OneDrive CLI (`odc`)! This project aims to provide a robust, coreutils-like experience for managing OneDrive files from the terminal.

## Development Environment

The easiest way to get started is by using the provided **VS Code DevContainer**. It comes pre-configured with Go 1.25, `just`, `podman`, and all necessary extensions.

If you prefer a manual setup, please refer to the [Developer Setup Guide](docs/developer/tutorials/setup.md).

## Project Structure

We follow a modular architectural pattern in the `internal` directory:
- **Core (`internal/core`):** Core domain services, filesystem abstractions, and identity management.
- **Slices (`internal/slices`):** CLI command definitions (Cobra) organized by functionality.
- **Dependency Injection (`internal/di`):** Service container and dependency resolution.
- **Middleware (`internal/features/middleware`):** API communication and logging middleware.

## Coding Standards

- **Go Version:** Go 1.25+
- **Error Handling:** Use package-level error variables and wrap errors using `%w`.
- **Logging:** Use the project's structured logging abstraction (`logging.Logger`).
- **Formatting:** Run `go fmt ./...` before committing.
- **Linting:** We use `golangci-lint`. Ensure your changes pass `golangci-lint run`.

See the full [Coding Standards](docs/developer/reference/coding-standards.md) for more details.

## Testing Requirements

- All new features and bug fixes **must** include unit tests.
- Use `stretchr/testify` for assertions and mocking.
- Run `t.Parallel()` in all tests where possible.
- Aim for high test coverage in the `app` and `domain` layers.

Run tests with: `go test ./...`

## Pull Request Process

1. Fork the repository and create your branch from `main`.
2. Ensure your code builds and passes all tests and linting.
3. Update the documentation in the `docs/` directory if your change affects the user or developer experience.
4. Submit a Pull Request with a clear description of the changes and the problem they solve.

## Documentation

Documentation is built using `mkdocs`. You can preview it locally using `just serve-docs`.
