# Configuring `odc` with the Builder pattern

This plan describes how to refactor the construction of the application's
`Config` struct using the Builder pattern. The goal is to provide a clear,
fluent API for configuring `odc` from multiple sources (defaults, files,
environment, and CLI flags).

Currently, the configuration initialization is a complex, multi-step process
spread across several functions. This makes it difficult to understand the
precedence of settings and ensures that defaults are consistently applied. A
dedicated `ConfigBuilder` provides a structured way to assemble these various
sources.

## Objectives

The primary goals for this refactoring include:

*   **Fluent Configuration API:** Create a clear, readable process for
    overriding default settings with session-specific values.
*   **Explicit Precedence:** Clearly define the order in which configuration
    sources (e.g., CLI > Env > File > Default) are applied.
*   **Robust Defaults:** Ensure that the `Config` struct is always valid, even
    if external sources are missing or incomplete.
*   **Maintainability:** Centralize the logic for parsing and merging
    configuration from disparate sources into a single, dedicated component.

## Proposed infrastructure (`internal/config/builder`)

The `ConfigBuilder` maintains the intermediate state of the configuration as
it is being constructed.

```go
type ConfigBuilder struct {
    config Config
}

func (b *ConfigBuilder) FromDefaults() *ConfigBuilder { ... }
func (b *ConfigBuilder) FromFile(path string) *ConfigBuilder { ... }
func (b *ConfigBuilder) FromEnv() *ConfigBuilder { ... }
func (b *ConfigBuilder) Build() (Config, error) { ... }
```

## Implementation steps

1.  **Define the builder:** Create the `ConfigBuilder` struct and its fluent
    methods in `internal/config/builder/`.
2.  **Migrate existing logic:** Move the logic for loading settings from files
    and the environment into the builder's methods.
3.  **Update CLI commands:** Update the `PreRunE` blocks to use the builder for
    applying command-line overrides to the configuration.
4.  **Simplify config initialization:** Replace the current multi-step
    initialization in `NewDefaultContainer` with a single, clear builder chain.

## Next steps

After refactoring the configuration system, we can easily add new configuration
sources, such as remote configuration servers or per-directory `.odc.yaml`
settings, simply by adding new methods to the `ConfigBuilder`.
