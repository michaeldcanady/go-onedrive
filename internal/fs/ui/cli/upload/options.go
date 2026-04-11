package upload

import (
	"errors"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options provides the settings for the drive upload command.
type Options struct {
	// Source is the local filesystem path of the item to upload.
	Source *fs.URI
	// Destination is the remote filesystem path where the item should be uploaded.
	Destination *fs.URI
	// Recursive determines whether to upload directories and their contents.
	Recursive bool
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Source == nil {
		return errors.New("source path is required")
	}
	if o.Destination == nil {
		return errors.New("destination path is required")
	}
	return nil
}
