package fs

import (
	"fmt"
	"path"
	"strings"
)

// URI represents a parsed URI with its components: provider, drive reference, and path.
type URI struct {
	// Provider is the name of the provider, e.g. "local", "onedrive", etc.
	Provider string
	// DriveRef is the drive reference, e.g. Drive ID for OneDrive, bucket name for S3/GCS, etc.
	DriveRef string
	// Path is the path within the DriveRef.
	Path string
}

// String returns a string representation of the URI.
func (u *URI) String() string {
	if u.DriveRef != "" {
		return fmt.Sprintf("%s:%s:%s", u.Provider, u.DriveRef, u.Path)
	}
	return fmt.Sprintf("%s:%s", u.Provider, u.Path)
}

// ManagerPath returns the path format expected by the FileSystemManager (provider:path).
// Note: For OneDrive, the path might already contain the drive ID if it was parsed from a 3-part URI.
func (u *URI) ManagerPath() string {
	if u.Provider == DefaultProviderPrefix && u.DriveRef != "" && u.DriveRef != "me" {
		return fmt.Sprintf("%s:%s:%s", u.Provider, u.DriveRef, u.Path)
	}
	return fmt.Sprintf("%s:%s", u.Provider, u.Path)
}

// ParseURI parses a raw string into a [URI] structure.
//
// It supports:
//   - provider:drive_ref:path (e.g., onedrive:me:/Documents)
//   - provider:path (e.g., local:/tmp, onedrive:/Documents)
//   - path (e.g., /Documents -> defaults to onedrive)
func ParseURI(uri string) (*URI, error) {
	if uri == "" {
		return nil, NewInvalidURIError(uri, "empty URI", nil)
	}

	parts := strings.Split(uri, ":")
	var u URI

	switch len(parts) {
	case 3:
		u.Provider = parts[0]
		u.DriveRef = parts[1]
		u.Path = parts[2]
	case 2:
		u.Provider = parts[0]
		u.Path = parts[1]
		// Special case: if provider looks like an absolute path, it might be a raw path.
		// But usually prefixes don't start with /.
		if strings.HasPrefix(u.Provider, "/") {
			u.Provider = DefaultProviderPrefix
			u.Path = uri
		}
	case 1:
		u.Provider = DefaultProviderPrefix
		u.Path = uri
	default:
		return nil, NewInvalidURIError(uri, "too many components", nil)
	}

	// Normalize path separators to "/" and clean the path.
	u.Path = strings.ReplaceAll(u.Path, "\\", "/")
	u.Path = path.Clean(u.Path)

	if ok, err := ContainsIllegalChars(u.Path); ok {
		return nil, NewInvalidURIError(uri, "illegal characters in path", err)
	}

	return &u, nil
}

// GetPath extracts the path component from the given URI string.
func GetPath(uri string) (string, error) {
	u, err := ParseURI(uri)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}

// GetDriveRef extracts the drive reference component from the given URI string.
func GetDriveRef(uri string) (string, error) {
	u, err := ParseURI(uri)
	if err != nil {
		return "", err
	}
	return u.DriveRef, nil
}

// GetProvider extracts the provider component from the given URI string.
func GetProvider(uri string) (string, error) {
	u, err := ParseURI(uri)
	if err != nil {
		return "", err
	}
	return u.Provider, nil
}
