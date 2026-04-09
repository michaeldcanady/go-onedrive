package cp

import (
	"errors"
	"io"
)

// Options provides the settings for the drive cp command.
type Options struct {
	// Source is the filesystem path of the item to copy.
	Source string
	// Destination is the filesystem path where the item should be copied.
	Destination string
	// Recursive determines whether to copy directories and their contents.
	Recursive bool
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
	// Stderr is the destination for error messages.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Source == "" {
		return errors.New("source path is required")
	}
	if o.Destination == "" {
		return errors.New("destination path is required")
	}
	return nil
}
