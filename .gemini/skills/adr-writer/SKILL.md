---
name: adr-writer
description: Helps in drafting, editing, and managing Architectural Decision Records (ADRs). Use when an architect-level decision needs to be formally documented, reviewed, or updated within the project's docs/adr/ directory.
---

# Architectural Decision Record (ADR) Writer

This skill helps you create and manage formal Architectural Decision Records (ADRs) for the `go-onedrive` project, following the [MADR](https://adr.github.io/madr/) standard.

## When to Use This Skill

- You are about to make a significant architectural change.
- You have recently made a decision and need to capture the rationale.
- You need to update an existing ADR (e.g., status change).
- You are adding a new ADR for a proposed change.

## Standard ADR Workflow

1. **Propose**: Create a new ADR using the standard template.
2. **Document**: Fill in the rationale, consequences, and context.
3. **Review**: Present the draft for team discussion.
4. **Finalize**: Update the status to `accepted` once approved.

## Resources

- **Template**: [references/adr-template.md](references/adr-template.md)
- **Guidelines**: [references/adr-guidelines.md](references/adr-guidelines.md)
- **ADR Index**: [docs/adr/index.md](docs/adr/index.md) (Use this to find existing ADRs)

## Commands for the User

If you need to check existing ADRs or manage the status of one, remember that you can always use the standard git and file management tools.

If you are unsure where to start, ask: "Can you help me document a new architectural decision regarding [TOPIC]?"
