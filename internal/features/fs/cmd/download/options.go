package download

import (
	"errors"
	"io"

	"github.com/michaeldcanady/go-onedrive/pkg/validation"
)

// Options provides the settings for the drive download command.
type Options struct {
	// Source is the remote filesystem path of the item to download.
	Source string

	// Destination is the local filesystem path where the item should be downloaded.
	Destination string

	// Recursive determines whether to download directories and their contents.
	Recursive bool

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	p := validation.All(
		validation.PolicyFunc[Options](func(o Options) error {
			if o.Source == "" || o.Destination == "" {
				return errors.New("source and destination paths are required")
			}
			return nil
		}),
	)

	return p.Evaluate(*o)
}
