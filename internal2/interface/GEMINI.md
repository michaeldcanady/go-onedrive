# GEMINI.md - Interface Layer

This layer handles the interaction between the user (CLI) and the application services.

## Interface Mandates

- **Cobra Usage:** All commands must be implemented using the Cobra library.
- **Decoupling:** Commands must not implement business logic. They should only parse flags/args, call the appropriate `app` service, and format the output.
- **Dependency Injection:** Commands must receive the `di.Container` to access services.
- **Consistency:** Flag naming, help text style, and error messages should remain consistent across all commands.

## Engineering Standards

- **Output Formatting:** Use the `formatting` package to support multiple output types (json, yaml, table, etc.).
- **Error Handling:** Use `internal2/interface/cli/util` to create standardized CLI errors.
- **Validation:** Use `PreRunE` for flag and argument validation.
- **Logging:** Ensure the command ID is registered with the logger and correlation IDs are handled.

## Implementation Guide

1. Create a sub-package for the command in `internal2/interface/cli/`.
2. Define the command in a `Create[Name]Cmd(container di.Container)` function.
3. Wire flags in the `init` or factory function.
4. Implement the execution logic in `RunE`.
5. Use `formatter.Format()` to display results to the user.
