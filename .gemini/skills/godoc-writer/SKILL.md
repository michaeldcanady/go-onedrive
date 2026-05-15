---
name: godoc-writer
description: Masterful Go documentation generator. Use when adding or improving GoDoc for packages, structs, methods, constants, and variables. Prioritizes high-signal, idiomatic comments following Effective Go and Go Code Review Comments philosophies.
---

# Godoc Writer

## Overview

The `godoc-writer` skill transforms bare or verbose Go code into a masterfully documented API. It focuses on providing essential context—the "why" and "how"—while omitting redundant information that is already evident from the code itself.

## Workflow

1.  **Analyze Context:** Read the file and surrounding package to understand the symbol's role. Check for existing documentation patterns.
2.  **Reference Guidelines:** Review [guidelines.md](references/guidelines.md) for core GoDoc philosophies (Symbol-First, Omit the Obvious, etc.).
3.  **Generate Documentation:** Draft concise, complete-sentence comments for the target symbols.
4.  **Refine:** Compare against [examples.md](references/examples.md) to ensure the documentation is high-signal and idiomatic.

## Masterful Documentation Patterns

### Package Level
Ensure every package has a summary. Use a `doc.go` file for complex packages.
- Pattern: `// Package name ...`

### Structs & Interfaces
Focus on the *purpose* and *concurrency safety*.
- Pattern: `// Name [purpose/coordination]. [Concurrency details].`

### Methods & Functions
Focus on *behavior*, *side effects*, and *error conditions*.
- Pattern: `// Name [behavior/returns]. [Side effects]. [Error conditions].`

## Resources

- **[guidelines.md](references/guidelines.md):** Deep dive into GoDoc philosophies and technical standards.
- **[examples.md](references/examples.md):** Before & After transformations for common Go symbols.
