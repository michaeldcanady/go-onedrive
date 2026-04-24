# Architectural Improvement Plan: Simplify Command Boilerplate

## Goal
Reduce repetitive code in command handlers by creating shared CLI middleware/helpers.

## Plan
1. **Identify Patterns**: Analyze `cmd` handlers to identify redundant logic (e.g., URI/path parsing, formatter setup).
2. **Abstract**: Move identified redundant logic into a new shared package, e.g., `internal/core/cli/helpers`.
3. **Integrate**: Refactor existing command handlers to utilize these helpers.
4. **Validate**: Ensure subcommands remain functional and behave as expected.

## Verification
- Confirm that the amount of boilerplate code in `internal/features/*/cmd` is reduced.
- Verify that command behavior, CLI flags, and output formats remain identical.
