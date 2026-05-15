# Plugin Specification: `storage-local`

## Overview
The `storage-local` plugin treats the local filesystem as a storage backend for the VFS.

## Capabilities
- **File Operations:** Standard POSIX-like file operations on the host machine.
- **Integration:** Allows `odc` to manage local files and cloud files using the same interface.
- **Testing:** Useful for testing VFS logic without external network calls.

## Configuration Options
The following options can be set via `odc config set storage.local.<key> <value>`:
- `root_path`: The local directory to treat as the root of this storage instance (defaults to the current working directory).

## Interface
Implements the `StorageService` gRPC interface as defined in `specs/proto/storage.proto`.

## Behavior
- Operates with the permissions of the user running the `odc` process.
- Translates gRPC storage requests into local system calls (`os`, `io` packages in Go).
