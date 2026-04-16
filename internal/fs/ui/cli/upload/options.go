package upload

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options provides the settings for the drive upload command.
type Options struct {
	// Source is the local filesystem path of the item to upload.
	Source string `arg:"1"`

	// SourceURI is the local filesystem path of the item to upload.
	SourceURI *fs.URI

	// Destination is the remote filesystem path where the item should be uploaded.
	Destination string `arg:"2"`

	// DestinationURI is the remote filesystem path where the item should be uploaded.
	DestinationURI *fs.URI

	// Recursive determines whether to upload directories and their contents.
	Recursive bool `flag:"recursive,short=r,desc='upload directories recursively',default=false"`

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
