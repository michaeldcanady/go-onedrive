# Architectural Improvement Plan: Decoupling VFS and Identity

## Goal
Decouple filesystem (VFS) operations from internal identity management, moving towards an abstract interface-based design.

## Plan
1. **Define Interface**: Define an `Authorizer` interface in `internal/core/auth/` (or equivalent) that exposes the necessary methods for token retrieval.
2. **Implement**: Implement this interface within the `internal/features/identity` slice.
3. **Inject**: Update the `VFS` component to accept an `Authorizer` instance via its constructor instead of depending directly on the identity store.
4. **Update DI**: Configure the DI container to provide the `Authorizer` implementation to the VFS.

## Verification
- Confirm that `internal/features/fs/domain/vfs.go` no longer has direct imports from `internal/features/identity`.
- Run unit tests for `VFS` using a mocked `Authorizer`.
