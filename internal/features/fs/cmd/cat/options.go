package cat

import (
	"errors"
	"io"

	"github.com/michaeldcanady/go-onedrive/pkg/validation"
)

// Options provides the settings for the drive cat command.
type Options struct {
	// Path is the filesystem path of the file to read.
	Path string
	// Stdout is the destination for the file's content.
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
