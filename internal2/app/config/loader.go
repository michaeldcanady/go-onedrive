package config

import "github.com/michaeldcanady/go-onedrive/internal2/infra/config"

// Loader defines the interface for reading configuration files.
//
// A Loader implementation is responsible for opening the file at the
// given path, parsing its contents, and returning a config.Configuration3
// value.
type Loader interface {
	// Load takes in the provided path and returns the parsed config.Configuration3 or an error.
	Load(path string) (config.Configuration3, error)
}
