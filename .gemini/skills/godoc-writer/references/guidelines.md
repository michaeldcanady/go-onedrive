# GoDoc Philosophy & Guidelines

Masterful Go documentation is about providing high signal-to-noise ratio. It should tell the reader what they can't already see from the code.

## Core Philosophies

1. **Symbol-First:** Every comment must be a complete sentence starting with the name of the symbol it describes.
   - *Good:* `// Service handles the orchestration of...`
   - *Bad:* `// The Service is...` or `// This struct handles...`

2. **Omit the Obvious:** If a function is named `Close()` and it's on a `File` struct, don't write `// Close closes the file.` Instead, describe side effects or return values.
   - *Masterful:* `// Close releases the underlying file descriptor. Any subsequent I/O returns ErrClosed.`

3. **Package-Level Docs:** Every package should have a package comment. If it's a large package, use a `doc.go` file.
   - Start with `// Package name ...`

4. **Constructors:** Document what the default state is if functional options aren't provided.

5. **Errors:** Document which specific error variables might be returned (e.g., "Returns ErrNotFound if the key does not exist").

## When NOT to Comment

- **Getters/Setters:** Unless they do something non-trivial like lazy initialization or cache invalidation.
- **Trivial Interface Implementations:** If a struct implements `io.Reader`, don't document the `Read` method unless it deviates from the standard `io.Reader` contract.
- **Internal/Private Symbols:** Use regular comments (`//`) for internal logic, but keep GoDocs (`//`) for the public API.

## Technical Standards

- **Paragraphs:** Separate paragraphs with a blank comment line (a line containing only `//`).
- **Preformatted Text:** Indent lines to show they are code or preformatted.
- **Links:** Use `[Name]` to link to other symbols in the same package, or `[*pkg.Name]` for other packages.
- **Deprecation:** Use `// Deprecated: Use [NewSymbol] instead.` as the first line of the comment.
- **Bugs:** Use `// BUG(who): description` for known issues.
