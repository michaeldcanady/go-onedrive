# Architecture Overview

The OneDrive CLI (`odc`) is built using a **Vertical Slice Architecture** (VSA).
This design focuses on organizing code by features rather than traditional
technical layers. This approach ensures that each command is self-contained, 
making the codebase easier to maintain, test, and extend.

## Why Vertical Slice Architecture?

In traditional layered architectures, a single feature (like `upload`) might be 
spread across multiple packages: `ui`, `application`, `domain`, and 
`infrastructure`. This often leads to "fragmented features," where a change in 
one place requires touching many files.

VSA addresses this by:

- **Feature Isolation:** Each CLI command is a self-contained "slice." 
  Modifying one command (e.g., `upload`) has a minimal "blast radius" and 
  is unlikely to break another (e.g., `ls`).
- **Improved Maintainability:** You can find all the code related to a 
  command—options parsing, business logic, and command definition—in a single 
  directory.
- **Reduced Coupling:** Slices depend on a stable "Core" but are independent 
  of each other.
- **Easier Testing:** Handlers are focused on a single responsibility, which 
  makes them straightforward to unit test using mocked core services.

## Core Principles

`odc` follows these fundamental architectural principles:

- **Organization by Feature:** Code is grouped by what the system does (e.g., 
  "Login," "List Files") rather than its technical role.
- **Self-Contained Slices:** Each slice contains the entry point (Cobra 
  command), the business logic (Command), and the specific options.
- **Shared Core:** Cross-cutting concerns like logging, configuration, 
  identity, and filesystem abstractions are kept in the `internal/` 
  directory as shared services.
- **Dependency Injection:** Services are wired in a central container and 
  passed to slices, ensuring decoupling and testability.

## Project Structure

The project is organized to support these principles:

- **`cmd/odc/`:** The entry point of the application.
- **`internal/`:** Contains the shared core services and the vertical slices.
    - **`internal/<feature>/ui/cli/<command>/`:** This is where a 
      vertical slice lives. It typically contains `command.go`, 
      `command_cmd.go`, and `options.go`.
    - **`internal/di/`:** Manages the Dependency Injection container.
    - **`internal/fs/`:** The filesystem abstraction layer.
    - **`internal/features/identity/`:** Authentication and identity management.
- **`pkg/`:** General-purpose utilities that are not specific to the 
  `odc` domain.

## Next steps

To learn more about how `odc` handles specific technical challenges, explore 
the other explanation guides:

- **[Dependency Injection](dependency-injection.md)**
- **[Configuration Management](configuration-management.md)**
- **[Filesystem Abstraction](filesystem-abstraction.md)**
