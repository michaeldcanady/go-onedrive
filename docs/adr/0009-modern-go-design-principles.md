# 9. Modern Go Design Principles and Coding Standards

Date: 2026-05-06

## Status

Status: Accepted

## Context

As the `go-onedrive` (odc) project evolves, it is essential to establish clear design principles and coding standards to maintain consistency, improve readability, and ensure robust performance in a CLI environment. These principles focus on idiomatic Go patterns, effective dependency management, and clear architectural boundaries.

## Decision

We adopt the following design principles and coding standards:

### 1. Concurrency and Cancellation
- **Context Usage:** All services, repositories, and long-running operations MUST accept `context.Context` as their first argument. This enables consistent cancellation, timeout management, and request tracing.
- **Thread-Safety:** All CLI commands and domain services MUST be designed for thread-safety. Global state should be avoided. When shared state is necessary, use appropriate synchronization primitives (e.g., `sync.Mutex`, `sync.RWMutex`, `atomic`).

### 2. Dependency Management
- **Decoupling (Interfaces):** Avoid hard-dependencies on concrete implementations of infrastructure services (e.g., logging, database). Services should depend on internal interfaces.
- **Interfaces In, Structs Out:** Functions and methods should generally accept interfaces to facilitate testing and flexibility, while returning concrete types (structs or pointers to structs) to provide clear API surface and allow for optimizations.
- **Scoped Interfaces:** Prioritize the use of small, local interfaces over large, global ones. Follow the hierarchy: Local Small -> Local Large -> Global Small -> Global Large.

### 3. Code Organization and Readability
- **Human Comprehension:** Code should be written primarily for human readability. Clear intent, meaningful naming, and idiomatic patterns are prioritized over clever but obscure optimizations.
- **File Structure:** Break types into their own files. Logically group functions within files to ensure each file has a clear, focused purpose.
- **Modern Patterns:** Adhere to modern Go design patterns, such as functional options for complex configuration and table-driven patterns for testing.

### 4. Architectural Integrity
- **Vertical Slices:** Maintain the vertical slice architecture (ADR 0001) where all feature-related logic resides together.
- **Domain-First Access:** Strictly adhere to the domain-first access pattern (ADR 0004). Commands must never bypass domain services to interact with persistence or external SDKs.

### 5. Error Handling and Panic Recovery
- **Explicit Error Handling:** Errors must NEVER be ignored. Every function call that returns an error must have that error checked and handled appropriately.
- **Idiomatic Error Context:** Follow standard Go philosophies for error handling: wrap errors with context using `%w`, return errors early (guard clauses), and only handle an error once.
- **Panic Usage:** Panics should be used sparingly and ONLY for truly exceptional, unrecoverable states (e.g., programmer errors like out-of-bounds access if not checked).
- **Graceful Recovery:** All panics MUST be recovered at appropriate boundaries (e.g., command entry points, plugin boundaries, or within long-running goroutines). The application and its plugins must never crash due to an unhandled panic. They should perform a graceful exit, logging the incident and ensuring state integrity.

## Consequences

### Benefits
- **Consistent DX:** Developers have clear guidelines for adding new features or modifying existing ones.
- **Improved Testability:** Decoupled dependencies and interface-based design make unit and integration testing easier and more reliable.
- **Robustness:** Consistent use of `context.Context` and thread-safe design reduces the risk of leaks and race conditions.
- **Reliability:** Explicit error handling and panic recovery ensure the application is resilient to unexpected conditions and never crashes ungracefully.
- **Maintainability:** Readable code and logical file organization simplify long-term maintenance and onboarding of new contributors.

### Trade-offs
- **Initial Overhead:** Adhering to these standards may require more initial planning and boilerplate code (e.g., defining local interfaces, writing recovery logic).
- **Strictness:** Requires active discipline and code reviews to ensure standards are met across the codebase.

## Links

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
