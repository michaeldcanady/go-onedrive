package domain

// Loader defines the interface for reading configuration files.
//
// A Loader implementation is responsible for opening the file at the
// given path, parsing its contents, and returning a Configuration
// value.
type Loader interface {
	// Load takes in the provided path and returns the parsed Configuration or an error.
	Load(path string) (Configuration, error)
}
