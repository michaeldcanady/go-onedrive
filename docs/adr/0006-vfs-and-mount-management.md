# 6. Virtual File System (VFS) and Mount Management

Date: 2026-04-22

## Status

Status: Accepted

## Context

Users of `odc` need to interact with multiple storage accounts and different types of storage (cloud and local) simultaneously. A unified way to address these different locations is required, similar to how a Unix-like operating system mounts different filesystems into a single global tree.

## Decision

We have decided to implement a Virtual File System (VFS) layer that manages multiple "mounts."
- **Mount Definition:** A mount maps a logical path (e.g., `/work`) to a specific backend provider instance (e.g., a specific OneDrive account) and a path within that provider.
- **Unified URI:** The VFS supports a unified URI scheme (e.g., `odc://profile/path/to/file`) and local relative/absolute paths by resolving them against active mounts.
- **Provider Abstraction:** All storage operations from the CLI layer go through the VFS, which dispatches the request to the appropriate backend provider based on the path.
- **Mount Persistence:** Mount configurations are stored in the persistent state (bbolt) to be available across sessions.

## Consequences

### Benefits
- **Unified Interface:** CLI commands (ls, cp, mv) work consistently regardless of whether the target is OneDrive, a local disk, or another plugin-provided backend.
- **Multi-Account Support:** Allows users to "mount" multiple OneDrive accounts and move files between them easily.
- **Abstraction:** Hides the complexity of different API structures (e.g., OneDrive DriveIDs and ItemIDs) behind a familiar path-based hierarchy.

### Trade-offs
- **Path Resolution Overhead:** Every operation requires a lookup in the mount table to determine the target provider.
- **Complexity in Symlinks/Cross-Mount Ops:** Moving files across different mounts (and thus different providers) requires implementing a "copy and delete" strategy rather than a simple "move" operation.

## Links

- [ADR 0002: Use bbolt for Persistent State Storage](0002-bbolt-for-persistent-state.md)
- [ADR 0005: External Plugin Architecture using gRPC and go-plugin](0005-external-plugin-architecture.md)
