package cat

import (
	"errors"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options provides the settings for the drive cat command.
type Options struct {
	// Path is the filesystem path of the file to read.
	Path *fs.URI
	// Stdout is the destination for the file's content.
	Stdout io.Writer
	// Stderr is the destination for error messages.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Path == nil {
		return errors.New("path is required")
	}
	return nil
}
