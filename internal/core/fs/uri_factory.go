package fs

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/alias"
	st "github.com/michaeldcanady/go-onedrive/internal/storage"
)

// URIFactory provides a structured way to create and validate URI objects.
type URIFactory struct {
	vfs interface {
		Resolve(absPath string) (string, string, error)
	}
	aliasSvc alias.Service
}

// NewURIFactory initializes a new instance of the URIFactory.
func NewURIFactory(vfs interface {
	Resolve(absPath string) (string, string, error)
}, aliasSvc alias.Service) *URIFactory {
	return &URIFactory{
		vfs:      vfs,
		aliasSvc: aliasSvc,
	}
}

// FromString parses a raw string and returns a structured URI object.
func (f *URIFactory) FromString(input string) (*URI, error) {
	if strings.HasPrefix(input, "/") {
		prefix, relPath, err := f.vfs.Resolve(input)
		if err == nil {
			return &URI{
				Provider: prefix,
				Path:     relPath,
			}, nil
		}
	}

	prefix, rest, found := strings.Cut(input, ":")
	if found {
		driveID, err := f.aliasSvc.GetDriveIDByAlias(context.Background(), prefix)
		if err == nil {
			if !strings.HasPrefix(rest, "/") {
				rest = "/" + rest
			}
			return &URI{
				Provider: st.DefaultProviderPrefix,
				DriveID:  driveID,
				Path:     rest,
			}, nil
		}
		return nil, fmt.Errorf("unknown mount point or alias: %s", prefix)
	}

	path := input
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return &URI{
		Provider: st.DefaultProviderPrefix,
		Path:     path,
	}, nil
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
	driveID, err := f.aliasSvc.GetDriveIDByAlias(context.Background(), name)
	if err != nil {
		return nil, fmt.Errorf("alias not found: %s", name)
	}

	if !strings.HasPrefix(subpath, "/") {
		subpath = "/" + subpath
	}

	return &URI{
		Provider: st.DefaultProviderPrefix,
		DriveID:  driveID,
		Path:     subpath,
	}, nil
}
