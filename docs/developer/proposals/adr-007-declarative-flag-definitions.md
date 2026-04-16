# ADR-007: Declarative Flag Definitions

## Status
Proposed

## Context
Flags are currently registered manually in each `Create<Cmd>Cmd` function (e.g., `cmd.Flags().BoolVarP(...)`). This leads to repetitive boilerplate and makes it easy to forget to update a flag's description or default value in the documentation vs. the code.

## Decision
Use struct tags in the `Options` struct to declaratively define CLI flags. A helper function will then reflect on the struct to register the flags with Cobra automatically.

```go
type Options struct {
    Recursive bool `flag:"recursive" short:"r" desc:"List items recursively" default:"false"`
}
```

## Consequences
- **Pros:**
    - Single Source of Truth: Flag definitions are colocated with the data they populate.
    - Reduced Boilerplate: No need for long lists of `cmd.Flags()` calls.
    - Automatic Documentation: Flag metadata can be extracted to generate CLI reference docs automatically.
- **Cons:**
    - Uses reflection, which can be slightly less performant (though negligible for CLI startup).
    - Harder to implement complex flag logic (e.g., mutually exclusive flags) purely through tags.
- **Impact:** Significant simplification of `Create<Cmd>Cmd` functions.

## Alternatives Considered
- **Option A:** Manual registration. Rejected as it is error-prone and repetitive.
- **Option B:** External configuration file. Rejected as it decouples flag definitions from the code too much.
