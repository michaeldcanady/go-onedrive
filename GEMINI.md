# go-onedrive Project Context

go-onedrive (odc) is a CLI tool for interacting with OneDrive items as a unix‑style file system.

## Project Overview

- **Purpose:** Provide users with a terminal‑native way to browse, read, edit, and manipulate OneDrive items as if they were part of a local filesystem.
- **Main Technologies:**
  - **Runtime:** Go (>= 1.25.0)
  - **Testing:** `go test`
  - **Linting/Formatting:** `go vet`, `golangci-lint`
- **Architecture:** Monorepo structure with clear layering:
  - `internal2/interface/cli`: User‑facing terminal UI, input parsing, and display rendering.
  - `pkg/odc`: Entrypoint for the odc CLI.
  - `internal2/app`: Application layer containing orchestrating services.
  - `internal2/domain`: Domain data types, value objects, and interfaces.
  - `internal2/infra`: Infrastructure for interacting with Microsoft Graph, caching, HTTP, identity, and persistence.

## Building and Running

- **Install Dependencies:**  
  `go get`
- **Build CLI:**  
  `just build`

## Testing and Quality

- **Run all unit tests:**  
  `go test ./...`

- **Static analysis:**  
  `go vet ./...`  
  `golangci-lint run`

## Development Conventions

- **Contributions:** Follow the process in `CONTRIBUTING.md`.  
  Requires signing the Google CLA.
- **Commit Messages:** Follow the
  [Conventional Commits](https://www.conventionalcommits.org/) standard.
- **Coding Style:**  
  - Go code must follow idiomatic Go patterns and odc’s established layering rules.
- **Imports:**  
  - Group imports (stdlib → external → internal).  
  - Avoid unused imports; keep import blocks clean and consistent.

## Testing Conventions

- **Go tests:**  
  - Prefer table‑driven tests.  
  - Use mocks for Graph API, caching, and identity providers.  
  - Tests must be deterministic and must not hit the real network.  
  - Use realistic OneDrive metadata and content fixtures.

---

# **Skill Usage in This Project**

The odc project uses two specialized skills to ensure consistent, high‑quality contributions: **`docs-writer`** and **`code-writer`**.

These skills must be invoked based on the type of work being performed.

## When to use the `docs-writer` skill

Use the `docs-writer` skill **whenever the task involves documentation**, including:

- Writing new documentation in the `docs/` directory.
- Editing or reviewing existing `.md` files anywhere in the repo.
- Updating docs to reflect changes in code behavior.
- Improving clarity, structure, or consistency of documentation.
- Suggesting documentation updates when code changes make existing docs incomplete.

If the task touches **any documentation**, the `docs-writer` skill is required.

## When to use the `code-writer` skill

Use the `code-writer` skill **whenever the task involves Go code or runtime behavior**, including:

- Writing new Go files in `cmd/`, `pkg/odc`, or `internal2/`.
- Editing or refactoring existing Go code.
- Implementing new services, domain types, or infra adapters.
- Updating CLI commands or wiring in `interface/cli`.
- Modifying configuration files that affect runtime behavior (`.json`, `.yaml`, etc.).
- Writing or updating Go tests.
- Performing architectural changes across domain/app/infra layers.

If the task modifies **any Go code or runtime‑affecting configuration**, the `code-writer` skill is required.

## When both skills may be needed

Some tasks require **both** skills, such as:

- Adding a new CLI command (code) and documenting it (docs).
- Changing behavior in `internal2/app` and updating user‑facing docs.
- Introducing a new feature that requires both implementation and documentation.

In these cases, the assistant should:

1. Use the **`code-writer`** skill for all code changes.  
2. Use the **`docs-writer`** skill for all documentation changes.