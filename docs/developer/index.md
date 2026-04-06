# Developer Guide

Welcome to the `odc` developer guide. This section is intended for contributors and anyone interested in understanding the internal workings of the OneDrive CLI.

## Overview

`odc` is built using a **Vertical Slice Architecture**, which organizes code by features rather than layers. This approach makes it easier to add new commands and maintain existing ones without navigating complex, cross-cutting layers.

## Key Technologies

- **Go:** The primary programming language.
- **Cobra:** CLI framework for command management and flag parsing.
- **Microsoft Graph SDK for Go:** Interface for communicating with the OneDrive/Graph API.
- **bbolt:** Local key-value store for profile management and state persistence.
- **zap:** Structured logging.

## Getting Started as a Developer

1.  **[Environment Setup](tutorials/setup.md):** Configure your local development environment.
2.  **[Architecture Overview](explanation/architecture.md):** Understand the core design principles of `odc`.
3.  **[Adding a New Command](how-to/add-subcommand.md):** A step-by-step guide to extending the CLI.

---

## Technical Reference

- **[Architecture Overview](explanation/architecture.md):** Understand the core design principles of `odc`.
- **[Dependency Injection](explanation/dependency-injection.md):** How services are wired and managed.
- **[Configuration Management](explanation/configuration-management.md):** How user settings and application state are handled.
- **[Domain Interfaces](reference/domain-interfaces.md):** Core interfaces that define the system's behavior.
- **[Error Handling](reference/cli-error-handling.md):** Standards for handling and reporting errors.
# - **[Testing Standards](how-to/testing.md):** How to write and run unit and integration tests.
