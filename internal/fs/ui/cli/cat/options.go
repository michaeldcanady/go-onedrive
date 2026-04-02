package cat

import (
	"errors"
	"io"
)

// Options provides the settings for the drive cat command.
type Options struct {
	// Path is the filesystem path of the file to read.
	Path string
	// Stdout is the destination for the file's content.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Path == "" {
		return errors.New("path is required")
	}
	return nil
}
