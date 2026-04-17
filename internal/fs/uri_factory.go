package fs

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
)

// URIFactory provides a structured way to create and validate URI objects.
type URIFactory struct {
	vfs interface {
		Mounts() []string
	}
	aliasSvc alias.Service
}

// NewURIFactory initializes a new instance of the URIFactory.
func NewURIFactory(vfs interface {
	Mounts() []string
}, aliasSvc alias.Service) *URIFactory {
	return &URIFactory{
		vfs:      vfs,
		aliasSvc: aliasSvc,
	}
}

// FromString parses a raw string and returns a structured URI object.
// It handles provider prefixes (e.g., "/local/tmp") and drive aliases (e.g., "work:/Documents").
func (f *URIFactory) FromString(input string) (*URI, error) {
	// 1. Check if input starts with a mount point
	mounts := f.vfs.Mounts()
	var bestPrefix string
	for _, mount := range mounts {
		if strings.HasPrefix(input, mount) {
			if len(mount) > len(bestPrefix) {
				bestPrefix = mount
			}
		}
	}

	if bestPrefix != "" {
		relPath := strings.TrimPrefix(input, bestPrefix)
		if relPath == "" {
			relPath = "/"
		}
		if !strings.HasPrefix(relPath, "/") {
			relPath = "/" + relPath
		}
		return &URI{
			Provider: bestPrefix,
			Path:     relPath,
		}, nil
	}

	prefix, rest, found := strings.Cut(input, ":")
	if !found {
		// Default to the onedrive provider
		path := input
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		return &URI{
			Provider: DefaultProviderPrefix,
			Path:     path,
		}, nil
	}

	// 2. Check if prefix is an alias
	driveID, err := f.aliasSvc.GetDriveIDByAlias(prefix)
	if err == nil {
		// If it's an alias, use the default provider and the drive ID
		if !strings.HasPrefix(rest, "/") {
			rest = "/" + rest
		}
		return &URI{
			Provider: DefaultProviderPrefix,
			DriveID:  driveID,
			Path:     rest,
		}, nil
	}

	// 3. Handle cases like "C:/path" on Windows or simple strings with colons that aren't providers/aliases.
	// For now, we return an error if a colon is present but no provider/alias matches.
	return nil, fmt.Errorf("unknown mount point or alias: %s", prefix)
}

// FromLocalPath creates a URI specifically for the local filesystem.
func (f *URIFactory) FromLocalPath(path string) (*URI, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return &URI{
		Provider: "/local",
		Path:     path,
	}, nil
}

// FromAlias creates a URI based on a drive alias and a subpath.
func (f *URIFactory) FromAlias(name, subpath string) (*URI, error) {
	driveID, err := f.aliasSvc.GetDriveIDByAlias(name)
	if err != nil {
		return nil, fmt.Errorf("alias not found: %s", name)
	}

	if !strings.HasPrefix(subpath, "/") {
		subpath = "/" + subpath
	}

	return &URI{
		Provider: DefaultProviderPrefix,
		DriveID:  driveID,
		Path:     subpath,
	}, nil
}
