# Domain Interfaces

This document outlines the core domain interfaces that define the behavior and
interactions within `odc`.

## Overview

`odc` uses interfaces to decouple the implementation details of different
services (like OneDrive, Local File System, or Mock providers) from the core
application logic. This approach facilitates testing, extensibility, and
maintenance.

## Filesystem Service

The Filesystem Service is the primary interface for all file and directory
operations.

### `fs.Service`

The `fs.Service` interface is a composite interface that includes `Namer`,
`Reader`, `Writer`, and `Manager`.

```go
type Service interface {
	Namer
	Reader
	Writer
	Manager
}
```

### `fs.Reader`

Defines operations for retrieving item information and file contents.

- **`Get(ctx context.Context, path string) (Item, error)`**: Retrieves a single 
  item by its path.
- **`List(ctx context.Context, path string, opts ListOptions) ([]Item, error)`**:
  Returns the immediate children of the specified directory path.
- **`ReadFile(ctx context.Context, path string, opts ReadOptions) (io.ReadCloser, error)`**:
  Provides an `io.ReadCloser` for the content of the file at the specified path.
- **`Stat(ctx context.Context, path string) (Item, error)`**: Returns metadata 
  for an item at the specified path.

### `fs.Writer`

Defines operations for modifying and creating items in the filesystem.

- **`WriteFile(ctx context.Context, path string, r io.Reader, opts WriteOptions) (Item, error)`**:
  Uploads or updates a file with the content from the provided reader.
- **`Mkdir(ctx context.Context, path string) error`**: Creates a new directory 
  at the specified path.
- **`Touch(ctx context.Context, path string) (Item, error)`**: Creates a new 
  empty file or updates the modification time of an existing one.

### `fs.Manager`

Defines operations for higher-level filesystem management and item manipulation.

- **`Remove(ctx context.Context, path string) error`**: Deletes an item from 
  the filesystem.
- **`Copy(ctx context.Context, src, dst string, opts CopyOptions) error`**:
  Duplicates an item from a source path to a destination path.
- **`Move(ctx context.Context, src, dst string) error`**: Relocates an item 
  from a source path to a destination path.

---

## Drive Service

The Drive Service manages OneDrive drive-specific operations.

### `drive.Service`

- **`ListDrives(ctx context.Context) ([]Drive, error)`**: Retrieves all 
  OneDrive drives accessible to the user.
- **`ResolveDrive(ctx context.Context, driveRef string) (Drive, error)`**:
  Identifies a drive by its ID, name, or alias.
- **`ResolvePersonalDrive(ctx context.Context) (Drive, error)`**: Retrieves 
  the user's primary personal OneDrive drive.
- **`GetActive(ctx context.Context, identityID string) (Drive, error)`**: Retrieves the 
  currently active drive.
- **`SetActive(ctx context.Context, driveID string, identityID string, scope shared.Scope) error`**:
  Marks a specific drive as the active one with the given scope.

---

## Other Key Interfaces

- **`identity.Authenticator`**: Manages user authentication and token retrieval.
- **`profile.Service`**: Handles user profiles and their associated configurations.
- **`config.Service`**: Manages application-level configuration.
- **`logger.Service`**: Defines the structured logging interface.
