package cp

import (
	"errors"
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/validation"
)

// Options provides the settings for the drive cp command.
type Options struct {
	// Source is the filesystem path of the item to copy.
	Source string

	// SourceURI is the filesystem path of the item to copy.
	SourceURI *fs.URI

	// Destination is the filesystem path where the item should be copied.
	Destination string

	// DestinationURI is the filesystem path where the item should be copied.
	DestinationURI *fs.URI

	// Recursive determines whether to copy directories and their contents.
	Recursive bool

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	p := validation.All(
		validation.PolicyFunc[Options](func(o Options) error {
			if o.Source == "" || o.Destination == "" {
				return errors.New("source and destination are required")
			}
			return nil
		}),
	)

	return p.Evaluate(*o)
}
