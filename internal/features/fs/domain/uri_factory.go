package fs

import (
	"fmt"
	"strings"
)

// URIFactory provides a structured way to create and validate URI objects.
type URIFactory struct {
	vfs interface {
		Resolve(absPath string) (string, string, error)
	}
}

// NewURIFactory initializes a new instance of the URIFactory.
func NewURIFactory(vfs interface {
	Resolve(absPath string) (string, string, error)
}) *URIFactory {
	return &URIFactory{
		vfs: vfs,
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
		mountPrefix := prefix
		if !strings.HasPrefix(mountPrefix, "/") {
			mountPrefix = "/" + mountPrefix
		}
		_, _, err := f.vfs.Resolve(mountPrefix)
		if err == nil {
			if !strings.HasPrefix(rest, "/") {
				rest = "/" + rest
			}
			return &URI{
				Provider: mountPrefix,
				Path:     rest,
			}, nil
		}
		return nil, fmt.Errorf("unknown mount point: %s", prefix)
	}

	path := input
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return &URI{
		Provider: DefaultProviderPrefix,
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

// FromMount creates a URI based on a mount prefix and a subpath.
func (f *URIFactory) FromMount(prefix, subpath string) (*URI, error) {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	_, _, err := f.vfs.Resolve(prefix)
	if err != nil {
		return nil, fmt.Errorf("mount point not found: %s", prefix)
	}

	if !strings.HasPrefix(subpath, "/") {
		subpath = "/" + subpath
	}

	return &URI{
		Provider: prefix,
		Path:     subpath,
	}, nil
}
