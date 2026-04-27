package get

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/pkg/validation"
)

// Options defines the configuration for the drive get operation.
type Options struct {
	// IdentityID is the specific account to query the active drive for.
	IdentityID string
	// DriveRef identifies the target drive by ID, name, or alias.
	DriveRef string
	// Stdout is the destination for the formatted drive details.
	Stdout io.Writer
	// Stderr is the destination for standard error messages.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	p := validation.Required(func(o Options) string { return o.DriveRef }, "drive reference")

	return p.Evaluate(*o)
}
