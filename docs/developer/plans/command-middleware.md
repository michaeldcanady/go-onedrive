# Command Middleware with the Chain of Responsibility pattern

This plan describes how to refactor common CLI command validation and execution
logic into a middleware system using the Chain of Responsibility pattern.
The goal is to eliminate redundant boilerplate in Cobra command definitions.

Most CLI commands in `odc` follow a similar pattern in their `PreRunE` and
`RunE` blocks: validating arguments, parsing URIs, checking provider
registries, and configuring logging. Middleware allows us to define these
steps as reusable "links" in a processing chain.

## Objectives

The primary goals for this refactoring include:

*   **Reduce boilerplate:** Simplify Cobra command definitions by moving
    common validation logic into reusable middleware.
*   **Clearer intent:** Make it obvious what a command requires (e.g., one path
    argument, a specific provider) by expressing it as a sequence of middleware
    calls.
*   **Centralized validation:** Ensure that common constraints like path syntax
    or argument counts are checked consistently across all commands.
*   **Testable components:** Each middleware "link" can be unit-tested in
    isolation, reducing the need for exhaustive end-to-end testing for simple
    validations.

## Proposed infrastructure (`internal/fs/ui/cli/middleware`)

Middleware is implemented as a function that takes a Cobra command and
returns another function or directly executes validation.

```go
type Middleware func(cmd *cobra.Command, args []string) error

func Chain(mw ...Middleware) Middleware {
    return func(cmd *cobra.Command, args []string) error {
        for _, m := range mw {
            if err := m(cmd, args); err != nil {
                return err
            }
        }
        return nil
    }
}
```

## Potential middlewares

Several common validations can be extracted into middleware:

*   **`RequireArgs(n)`**: Ensures the exact number of required arguments is
    present.
*   **`WithParsedURI(index, target)`**: Parses the argument at `index` as a
    `shared.URI` and stores it in the provided target.
*   **`WithValidPath(index)`**: Validates the syntax of the path at the
    specified argument index.
*   **`WithRegisteredProvider(index)`**: Verifies that the provider prefix in
    the URI at `index` is registered in the container.

## Implementation steps

1.  **Create the middleware engine:** Implement the `Chain` and standard
    middleware functions in `internal/fs/ui/cli/middleware/`.
2.  **Define domain middlewares:** Implement validations for URI parsing, path
    syntax, and provider checks.
3.  **Refactor CLI commands:** Update existing commands like `ls`, `edit`, and
    `upload` to use the middleware chain in their `PreRunE` blocks.
4.  **Simplify command logic:** Move argument processing from `PreRunE` into
    the chain, making the command definitions cleaner and more declarative.

## Next steps

After refactoring `PreRunE` logic, we can explore "Execution Middleware" for the
`RunE` phase, which could handle cross-cutting CLI concerns like output
formatting, automatic error reporting, or terminal cleanup after an interactive
session.
