package rm

import (
	"errors"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options provides the settings for the drive rm command.
type Options struct {
	// Path is the filesystem path of the item to remove.
	Path string `arg:"1"`
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
	// Recursive determines whether to remove directories and their contents.
	Recursive bool `flag:"recursive,short=r,desc='Remove directories and their contents recursively',default=false"`
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
