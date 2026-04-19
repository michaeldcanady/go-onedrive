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

- **`Get(ctx context.Context, uri *URI) (Item, error)`**: Retrieves a single 
  item by its structured URI.
- **`List(ctx context.Context, uri *URI, opts ListOptions) ([]Item, error)`**:
  Returns the immediate children of the specified directory URI.
- **`ReadFile(ctx context.Context, uri *URI, opts ReadOptions) (io.ReadCloser, error)`**:
  Provides an `io.ReadCloser` for the content of the file at the specified URI.
- **`Stat(ctx context.Context, uri *URI) (Item, error)`**: Returns metadata 
  for an item at the specified URI.

### `fs.Writer`

Defines operations for modifying and creating items in the filesystem.

- **`WriteFile(ctx context.Context, uri *URI, r io.Reader, opts WriteOptions) (Item, error)`**:
  Uploads or updates a file with the content from the provided reader at the specified URI.
- **`Mkdir(ctx context.Context, uri *URI) error`**: Creates a new directory 
  at the specified URI.
- **`Touch(ctx context.Context, uri *URI) (Item, error)`**: Creates a new 
  empty file or updates the modification time of an existing one at the specified URI.

### `fs.Manager`

Defines operations for higher-level filesystem management and item manipulation.

- **`Remove(ctx context.Context, uri *URI) error`**: Deletes an item from 
  the filesystem at the specified URI.
- **`Copy(ctx context.Context, src, dst *URI, opts CopyOptions) error`**:
  Duplicates an item from a source URI to a destination URI.
- **`Move(ctx context.Context, src, dst *URI) error`**: Relocates an item 
  from a source URI to a destination URI.

---

## Drive Service

The Drive Service manages OneDrive drive-specific operations.

### `drive.Service`

- **`ListDrives(ctx context.Context, identityID string) ([]Drive, error)`**: Retrieves all 
  OneDrive drives accessible to the user, optionally scoped to an identity.
- **`ResolveDrive(ctx context.Context, driveRef string, identityID string) (Drive, error)`**:
  Identifies a drive by its ID or name, optionally scoped to an identity.
- **`ResolvePersonalDrive(ctx context.Context, identityID string) (Drive, error)`**: Retrieves 
  the user's primary personal OneDrive drive, optionally scoped to an identity.

---

## Other Key Interfaces

- **`identity.Authenticator`**: Manages user authentication and token retrieval.
- **`profile.Service`**: Handles user profiles and their associated configurations.
- **`config.Service`**: Manages application-level configuration.
- **`logger.Service`**: Defines the structured logging interface.
