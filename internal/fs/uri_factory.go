package fs

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
)

// URIFactory provides a structured way to create and validate URI objects.
type URIFactory struct {
	registry interface {
		RegisteredNames() ([]string, error)
	}
	aliasSvc alias.Service
}

// NewURIFactory initializes a new instance of the URIFactory.
func NewURIFactory(registry interface {
	RegisteredNames() ([]string, error)
}, aliasSvc alias.Service) *URIFactory {
	return &URIFactory{
		registry: registry,
		aliasSvc: aliasSvc,
	}
}

// FromString parses a raw string and returns a structured URI object.
// It handles provider prefixes (e.g., "local:/tmp") and drive aliases (e.g., "work:/Documents").
func (f *URIFactory) FromString(input string) (*URI, error) {
	prefix, rest, found := strings.Cut(input, ":")
	if !found {
		// Default to the onedrive provider with the full input as the path
		return &URI{
			Provider: DefaultProviderPrefix,
			Path:     input,
		}, nil
	}

	// 1. Check if prefix is a registered provider
	names, err := f.registry.RegisteredNames()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve registered providers: %w", err)
	}

	for _, name := range names {
		if prefix == name {
			return &URI{
				Provider: name,
				Path:     rest,
			}, nil
		}
	}

	// 2. Check if prefix is an alias
	driveID, err := f.aliasSvc.GetDriveIDByAlias(prefix)
	if err == nil {
		// If it's an alias, use the default provider and the drive ID
		return &URI{
			Provider: DefaultProviderPrefix,
			DriveID:  driveID,
			Path:     rest,
		}, nil
	}

	// 3. Handle cases like "C:/path" on Windows or simple strings with colons that aren't providers/aliases.
	// For now, we return an error if a colon is present but no provider/alias matches.
	return nil, fmt.Errorf("unknown provider or alias: %s", prefix)
}

// FromLocalPath creates a URI specifically for the local filesystem.
func (f *URIFactory) FromLocalPath(path string) (*URI, error) {
	return &URI{
		Provider: "local",
		Path:     path,
	}, nil
}

// FromAlias creates a URI based on a drive alias and a subpath.
func (f *URIFactory) FromAlias(name, subpath string) (*URI, error) {
	driveID, err := f.aliasSvc.GetDriveIDByAlias(name)
	if err != nil {
		return nil, fmt.Errorf("alias not found: %s", name)
	}

	return &URI{
		Provider: DefaultProviderPrefix,
		DriveID:  driveID,
		Path:     subpath,
	}, nil
}
