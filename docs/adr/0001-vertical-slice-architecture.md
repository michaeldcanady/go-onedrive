# 1. Use Vertical Slice Architecture for Feature Organization

Date: 2025-05-14

## Status

Status: Accepted

## Context

The `odc` project needs a maintainable and scalable way to organize code. Traditional layered architecture (e.g., separating all controllers, services, and repositories into distinct top-level directories) often leads to "shotgun surgery" where adding or modifying a single feature requires jumping between many layers.

## Decision

We have decided to organize code into **Vertical Slices** based on features within the `internal/features/` directory. Each slice (e.g., `identity`, `profile`, `drive`) should contain its own:
- Domain logic and interfaces
- Data access implementations (repositories)
- CLI commands and handlers
- Feature-specific DTOs or Protobuf definitions

Shared infrastructure and utilities are kept in `internal/core/` or `pkg/`.

## Consequences

## Benefits
- **Improved Cohesion:** All code related to a specific feature is located in one place.
- **Reduced Coupling:** Features are largely independent, making it harder for changes in one feature to break another.
- **Scalability:** New features can be added by creating a new directory under `internal/features/` without cluttering existing layers.
- **Easier Navigation:** Developers can find everything they need for a feature within a single directory tree.

## Trade-offs
- **Potential Duplication:** There is a risk of duplicating small utility logic between slices if not properly extracted to `pkg/`.
- **Architectural Discipline:** Requires discipline to ensure features don't start depending on each other's internal implementation details rather than their public interfaces.

## Links

- [Vertical Slice Architecture](https://github.com/jbogard/VerticalSliceArchitecture)
