---
name: conventional-commit-writer
description: Assists in creating standardized conventional commit messages. Use when committing changes to ensure consistent versioning and changelog generation.
---

# Conventional Commit Writer

## Overview

This skill guides you through the process of writing high-quality conventional commit messages that follow the repository's standards.

## Workflow

1. **Gather Context**: Run `git status` and `git diff HEAD` to understand the changes.
2. **Evaluate Strategy**: Decide if the changes should be split into multiple atomic commits (see [references/strategy.md](references/strategy.md)).
3. **Review Style**: Run `git log -n 3` to see recent commit messages and match the project's tone and scope naming conventions.
4. **Determine Type**: Use [references/types.md](references/types.md) to choose the correct commit type (feat, fix, docs, etc.).
5. **Identify Scope**: If the change is localized to a specific feature slice or package (e.g., `mount`, `identity`, `di`), include it in parentheses.
6. **Draft Message**: Write a concise, imperative description (e.g., "add support for...", "fix bug in...").
7. **Identify Breaking Changes**: If the change breaks backward compatibility, add a `BREAKING CHANGE:` footer.

## Guidelines

- **Atomic Commits**: Each commit should do exactly one thing.
- **Imperative Mood**: Use "fix", not "fixed" or "fixes".
- **No Period**: Do not end the subject line with a period.
- **Lowercase**: Use lowercase for the subject line.
- **Conciseness**: Keep the subject line under 50-72 characters.

## References

- **[Commit Types](references/types.md)**: Detailed list of allowed commit types.
- **[Commit Strategy](references/strategy.md)**: Guidance on atomic commits and handling large changes.
- **[Examples](references/examples.md)**: Common commit message patterns.
