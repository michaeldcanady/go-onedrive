# ADR-005: Validation Logic Stratification

## Status
Proposed

## Context
Validation logic is currently split between `Options.Validate()` and `handler.Validate()`, but the division of responsibilities is not always clear. This can result in redundant checks or missing validations if one layer assumes the other has performed a check.

## Decision
Explicitly stratify validation into two layers:
1.  **Syntactic Validation (`Options.Validate`):** Checks that are independent of application state (e.g., "is the path non-empty?", "is the count greater than zero?").
2.  **Semantic Validation (`handler.Validate`):** Checks that require external services or state (e.g., "does the source file exist?", "do we have write permissions to this drive?").

## Consequences
- **Pros:**
    - Clearer "Rules of Engagement" for developers.
    - Syntactic checks can run without initializing heavy dependencies.
    - Improved error messages: Distinguishes between "you gave me bad input" and "the operation is invalid in the current state."
- **Cons:**
    - Requires moving some code between files.
- **Impact:** Cleaner handlers and more robust input validation.

## Alternatives Considered
- **Option A:** Put everything in `handler.Validate`. Rejected as it makes the handler too bulky.
- **Option B:** Put everything in `Options.Validate`. Rejected as `Options` should not depend on services like `fs.Service`.
