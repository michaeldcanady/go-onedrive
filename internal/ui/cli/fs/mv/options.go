package mv

import (
	"errors"
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs"
)

// Options provides the settings for the drive mv command.
type Options struct {
	// Source is the filesystem path of the item to move.
	Source string

	// SourceURI is the filesystem path of the item to move.
	SourceURI *fs.URI

	// Destination is the filesystem path where the item should be moved.
	Destination string

	// DestinationURI is the filesystem path where the item should be moved.
	DestinationURI *fs.URI

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Source == "" || o.SourceURI == nil {
		return errors.New("source path is required")
	}
	if o.Destination == "" || o.DestinationURI == nil {
		return errors.New("destination path is required")
	}
	return nil
}
