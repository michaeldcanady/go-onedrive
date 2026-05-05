# 0008 Repository Pattern

Date: 2026-05-05

## Status

Status: Proposed

## Context

The `go-onedrive` project currently uses several different approaches for persistence across its feature slices:
- `profile` slice has a `BoltRepository` but it manages its own `*bolt.DB` connection and is tightly coupled to the service.
- `identity` slice uses an `IdentityRepository` interface and a `BoltRepository` that takes an external `*bolt.DB`.
- `config` and `mount` slices use `yaml_repository.go` and `config_repository.go` respectively.

This lack of consistency makes it difficult to:
1.  Unit test services without actual DB files.
2.  Standardize error handling for persistence failures.
3.  Manage database connection lifetimes centrally.

## Decision

We will adopt a standardized Repository pattern across all feature slices.

### Key Components

1.  **Repository Interface:** Each feature slice requiring persistence will define a clear interface for its data access needs.
2.  **Concrete Implementations:** Implementations (e.g., BoltDB, YAML, In-memory) will be provided and injected via the DI container.
3.  **Decoupled Lifecycle:** Repositories will NOT manage the opening or closing of their underlying data sources (like `*bolt.DB` or file handles). These will be managed by a higher-level storage service or the DI container itself.

### Implementation Details

- Feature-specific repositories (e.g., `ProfileRepository`, `IdentityRepository`) will be defined within their respective feature slices.
- The `internal/features/storage` service will be responsible for providing the raw database handles required by the repositories.

## Consequences

### Benefits
- **Testability:** Clearer boundaries for mocking persistence in unit tests.
- **Consistency:** Uniform data access patterns across the codebase.
- **Resource Management:** Better control over database file handles.
- **Flexibility:** Easier to swap storage backends (e.g., from BoltDB to another KV store).

### Trade-offs
- **Boilerplate:** Requires defining more interfaces and factory methods.

## Links

- [ADR 0002: bbolt for Persistent State](0002-bbolt-for-persistent-state.md)
- [ADR 0004: Domain-First Access Pattern](0004-domain-first-access-pattern.md)
