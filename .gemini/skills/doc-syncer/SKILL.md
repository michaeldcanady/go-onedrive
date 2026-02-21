---
name: doc-syncer
description: Maintains synchronization between Go source code and developer documentation. Use when Go files in 'internal2/' are modified to identify and update corresponding documentation in 'docs/developer/' following the Diátaxis framework.
---

# Doc Syncer Skill

This skill provides a structured workflow for ensuring that the `odc` developer documentation stays in sync with changes to the Go source code.

## Workflow

### 1. Research: Identify Changed Files

Start by identifying the Go files that have been modified. Use standard Git commands:

```bash
git diff --name-only HEAD | grep "\.go$"
```

If multiple files are changed, group them by package to simplify documentation updates.

### 2. Strategy: Map Go Changes to Docs

For each changed Go package, determine the relevant documentation files using [mapping.md](references/mapping.md).

- **Domain Changes:** Focus on `reference/domain-interfaces.md` and related `explanation/` docs.
- **Application Changes:** Focus on `explanation/` and `tutorials/`.
- **Infrastructure Changes:** Focus on `reference/` (repositories) and `how-to/` guides.
- **Interface/CLI Changes:** Focus on `reference/` (patterns, error handling) and `how-to/add-subcommand.md`.

### 3. Execution: Update Documentation

Follow these Diátaxis-specific guidelines when applying updates:

- **Tutorials:** Update step-by-step instructions if the process has changed. Ensure the tutorial still achieves its specific learning goal.
- **How-to Guides:** Update specific instructions for tasks. Ensure any new flags or options are covered.
- **Explanations:** Update conceptual overviews and diagrams (Mermaid) if the architectural design or philosophy has shifted.
- **References:** Update API-like documentation, constants, interfaces, and technical specifications. This is usually the most important for Go interface changes.

### 4. Validation: Verify Documentation

After making changes:

1.  Check for broken links within the modified documentation.
2.  If new files were created, ensure they are registered in `mkdocs.yaml`.
3.  Propose the changes to the user, highlighting which code modifications triggered which documentation updates.

## Best Practices

- **Surgical Updates:** Only modify the parts of the documentation directly affected by the code change.
- **Mermaid Diagrams:** Always update related Mermaid diagrams in `explanation/` docs if the component interaction changes.
- **Consistency:** Maintain the tone and formatting of the existing documentation.
- **No Filler:** Avoid repetitive or obvious information. Focus on "Why" (Explanation) and "How" (Tutorial/How-to).
