# ADR Guidelines

## When to write an ADR?

- You are about to make a change that affects the architecture, dependency stack, or fundamental design patterns of `odc`.
- You are adding a new feature that introduces significant new abstractions.
- You are choosing between multiple technologies or implementation strategies.

## How to write an ADR?

1. **Keep it concise**: Focus on the rationale and trade-offs.
2. **Be specific**: Use clear, actionable language.
3. **Be honest**: Explicitly list the negative consequences/trade-offs of your decision.
4. **Use the Template**: Always start with `references/adr-template.md`.
5. **Update Status**: Ensure the status is clearly marked.
6. **File Location**: Save in `docs/adr/` with a descriptive, numeric filename (e.g., `docs/adr/0001-di-container-refactor.md`).
