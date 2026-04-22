package mkdir

import (
	"errors"
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs"
)

// Options provides the settings for the drive mkdir command.
type Options struct {
	// Path is the filesystem path of the directory to create.
	Path string
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
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
