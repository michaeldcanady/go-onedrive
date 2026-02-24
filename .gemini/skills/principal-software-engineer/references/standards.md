# Engineering Standards & Quality

## Architectural Principles
- **Separation of Concerns:** Maintain strict boundaries between Domain, App, and Infra layers.
- **Dependency Inversion:** Depend on abstractions, not implementations.
- **Single Responsibility:** Every module or function should do one thing well.
- **YAGNI (You Ain't Gonna Need It):** Avoid "just-in-case" abstractions or features.

## Quality Standards
- **Test Coverage:** Every bug fix must include a reproduction test. Every new feature must have table-driven unit tests.
- **Error Handling:** Wrap errors with context. Never swallow errors. Use domain-specific error types where appropriate.
- **Code Reviews:** Focus on maintainability, security, and adherence to established patterns over stylistic nitpicks.
- **Performance:** For CLI tools, prioritize startup time and efficient I/O (e.g., lazy loading, streaming).

## Verification Checklist
- [ ] Does this follow the project's layering rules?
- [ ] Is the code idiomatic Go?
- [ ] Are there enough tests to verify the behavior and edge cases?
- [ ] Is the documentation updated?
