package download

import (
	"errors"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options provides the settings for the drive download command.
type Options struct {
	// Source is the remote filesystem path of the item to download.
	Source string

	// SourceURI is the remote filesystem path of the item to download.
	SourceURI *fs.URI

	// Destination is the local filesystem path where the item should be downloaded.
	Destination string

	// DestinationURI is the local filesystem path where the item should be downloaded.
	DestinationURI *fs.URI

	// Recursive determines whether to download directories and their contents.
	Recursive bool

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
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
