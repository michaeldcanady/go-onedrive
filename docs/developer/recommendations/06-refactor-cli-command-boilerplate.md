# Recommendation: Refactor CLI Command Boilerplate

## Goal
Reduce the repetitive boilerplate code in CLI command definitions and improve consistency.

## Current State
- Each command in `internal/**/ui/cli/` follows a similar structure: `command.go` (Cobra definition), `handler.go` (logic), `options.go` (flags/validation), and `command_context.go`.
- The `Create*Cmd` functions in `command.go` share common patterns for logger initialization and handler setup.

## Value
- **Medium**: Makes adding new commands faster and less error-prone.
- Improves code readability by focusing on the unique logic of each command.

## Implementation Plan
1.  **Identify Common Patterns**: Extract common setup logic like creating a logger for the command name and initializing a basic `CommandContext`.
2.  **Create Command Helper/Base**: Develop a helper function or a base struct in `internal/fs/ui/cli` (or a more shared package) to handle common Cobra setup.
3.  **Refactor Existing Commands**: Update existing commands to use the new helper, significantly reducing the size of `command.go` files.

## Difficulty
- **Low**: This refactoring is local to the UI layer and doesn't affect core domain logic.
