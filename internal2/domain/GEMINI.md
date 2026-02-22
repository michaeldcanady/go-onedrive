# GEMINI.md - Domain Layer

This file outlines the mandates and standards for the `internal2/domain` package, which is the core of the `go-onedrive` project.

## Domain Mandates

- **Dependency Isolation:** No external dependencies (SDKs, DBs, etc.) are allowed here. Only standard Go libraries.
- **Interface Definitions:** All core business services must be defined as interfaces in this layer.
- **Entities & Types:** Define fundamental data structures and custom types (e.g., `Item`, `Metadata`, `LoginOptions`).
- **Validation:** Business rules and logic should reside here where possible.

## Engineering Standards

- **Error Definition:** Define all package-level errors in `errors.go` or alongside the service interface.
- **Documentation:** Every interface and type MUST have a clear GoDoc comment.
- **Testing:** Unit tests should focus on business logic and validation rules, using mock dependencies where necessary.

## Implementation Guide

When adding a new feature:
1. Identify the core domain concepts.
2. Define the interface in `internal2/domain/[package]/service.go`.
3. Create the data types in `internal2/domain/[package]/[concept].go`.
4. Ensure no `infra` or `app` imports are introduced.
5. Provide high-level technical rationale for the interface design.
