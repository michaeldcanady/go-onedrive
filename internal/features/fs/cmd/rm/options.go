package rm

import (
	"errors"
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/validation"
)

// Options provides the settings for the drive rm command.
type Options struct {
	// Path is the filesystem path of the item to remove.
	Path string
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	p := validation.All(
		validation.PolicyFunc[Options](func(o Options) error {
			if o.Path == "" {
				return errors.New("path is required")
			}
			return nil
		}),
	)

	return p.Evaluate(*o)
}
