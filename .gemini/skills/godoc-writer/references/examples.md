# GoDoc Examples: Before & After

## Structs

### Before (Verbose & Redundant)
```go
// This is the Client struct which is used to make API calls to the server.
// It contains the configuration and the HTTP client.
type Client struct {
    // Config is the configuration for the client
    Config *Config
    // HTTPClient is the client used for requests
    HTTPClient *http.Client
}
```

### After (Masterful)
```go
// Client coordinates communication with the OneDrive API.
// It is safe for concurrent use by multiple goroutines.
type Client struct {
    Config     *Config
    HTTPClient *http.Client
}
```

## Methods

### Before (Noisy)
```go
// GetDriveByID retrieves a drive from the API using the provided ID.
// It returns a pointer to the Drive and an error if one occurred.
func (s *Service) GetDriveByID(ctx context.Context, id string) (*Drive, error) { ... }
```

### After (High Signal)
```go
// GetDriveByID returns the drive metadata for the given unique identifier.
// It returns [ErrNotFound] if the ID does not map to an existing drive.
func (s *Service) GetDriveByID(ctx context.Context, id string) (*Drive, error) { ... }
```

## Package Level

### Before (Missing Context)
```go
package drive
```

### After (Masterful `doc.go`)
```go
// Package drive provides high-level abstractions for interacting with
// Microsoft Graph Drive resources, including support for large file
// uploads and recursive synchronization.
package drive
```
