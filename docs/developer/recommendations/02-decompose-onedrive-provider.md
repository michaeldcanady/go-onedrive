# Recommendation: Decompose OneDrive Provider

## Goal
Break down the monolithic `onedrive.Provider` into smaller, more focused components to improve readability, maintainability, and testability.

## Current State
- `internal/fs/providers/onedrive/provider.go` is a large file (~600 lines) that handles:
    - URI expansion and path resolution.
    - Item mapping (Graph models to `shared.Item`).
    - Error mapping.
    - Small file uploads and downloads.
    - Large (resumable) file uploads.
    - Directory operations.

## Value
- **High**: Reduces the cognitive load for developers working on the OneDrive integration.
- Enables more granular unit testing of individual components (e.g., testing the upload logic without needing a full provider).
- Improves code reuse (e.g., the error mapper could be shared).

## Implementation Plan
1.  **Extract Error Mapping**: Move error mapping logic to a separate helper or utility within the package.
2.  **Extract Item Mapping**: Create a dedicated `mapper.go` for converting Microsoft Graph models to internal `fs.Item` structures.
3.  **Extract Transfer Logic**: Move `writeLargeFile` and related upload/download logic into a `TransferManager` or similar component.
4.  **Extract URI/Path Logic**: Move URI expansion and path manipulation to a specialized helper.
5.  **Refactor Provider**: Update the `Provider` struct to use these new components, keeping its role as an orchestrator that implements the `fs.Service` interface.

## Difficulty
- **Medium**: Requires careful extraction of logic and updating dependencies. Testing will be needed to ensure no regressions in complex operations like resumable uploads.
