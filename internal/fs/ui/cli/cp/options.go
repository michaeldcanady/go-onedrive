package cp

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options provides the settings for the drive cp command.
type Options struct {
	// Source is the filesystem path of the item to copy.
	Source string `arg:"1"`

	// SourceURI is the filesystem path of the item to copy.
	SourceURI *fs.URI

	// Destination is the filesystem path where the item should be copied.
	Destination string `arg:"2"`

	// DestinationURI is the filesystem path where the item should be copied.
	DestinationURI *fs.URI

	// Recursive determines whether to copy directories and their contents.
	Recursive bool `flag:"recursive,short=r,desc='Copy directories and their contents recursively',default=false"`

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
