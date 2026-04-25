package upload

import (
	"errors"
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/validation"
)

// Options provides the settings for the drive upload command.
type Options struct {
	// Source is the local filesystem path of the item to upload.
	Source string

	// SourceURI is the local filesystem path of the item to upload.
	SourceURI *fs.URI

	// Destination is the remote filesystem path where the item should be uploaded.
	Destination string

	// DestinationURI is the remote filesystem path where the item should be uploaded.
	DestinationURI *fs.URI

	// Recursive determines whether to upload directories and their contents.
	Recursive bool

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	p := validation.All(
		validation.PolicyFunc[Options](func(o Options) error {
			if o.Source == "" {
				return errors.New("source path is required")
			}
			if o.Destination == "" {
				return errors.New("destination path is required")
			}
			return nil
		}),
	)

	return p.Evaluate(*o)
}
