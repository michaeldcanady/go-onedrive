package mkdir

import (
	"errors"
	"io"
)

// Options provides the settings for the drive mkdir command.
type Options struct {
	// Path is the filesystem path of the directory to create.
	Path string
	// Stdout is the destination for standard output messages.
	Stdout io.Writer

	// Stderr is the destination for error messages.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Path == "" {
		return errors.New("path is required")
	}
	return nil
}
