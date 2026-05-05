# Recommendation: Decompose OneDrive provider

## Goal
Break down the monolithic `onedrive.Provider` into smaller, more focused components to improve readability, maintainability, and testability

## Current state
- `internal/fs/providers/onedrive/provider.go` is a large file (~600 lines) that handles:
    - URI expansion and path resolution
    - Item mapping (Graph models to `shared.Item`)
    - Error mapping
    - Small file uploads and downloads
    - Large (resumable) file uploads
    - Directory operations

## Value
- **High**: Reduces the cognitive load for developers working on the OneDrive integration
- Enables more granular unit testing of individual components (for example, testing the upload logic without needing a full provider)
- Improves code reuse (for example, the error mapper could be shared)

## Implementation plan
1.  **Remove Error Mapping**: Move error mapping logic to a separate helper or utility within the package
2.  **Remove Item Mapping**: Create a dedicated `mapper.go` for converting Microsoft Graph models to internal `fs.Item` structures
3.  **Remove Transfer Logic**: Move `writeLargeFile` and related upload/download logic into a `TransferManager` or similar component
4.  **Remove URI/Path Logic**: Move URI expansion and path manipulation to a specialized helper
5.  **Refactor Provider**: Update the `Provider` struct to use these new components, keeping its role as an orchestrator that implements the `fs.Service` interface

## Difficulty
- **Medium**: Requires careful extraction of logic and updating dependencies. Testing will be needed to confirm no regressions in complex operations like resumable uploads
