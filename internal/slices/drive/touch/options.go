package touch

import (
	"errors"
	"io"
)

// Options provides the settings for the drive touch command.
type Options struct {
	// Path is the filesystem path of the file to touch.
	Path string
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Path == "" {
		return errors.New("path is required")
	}
	return nil
}
