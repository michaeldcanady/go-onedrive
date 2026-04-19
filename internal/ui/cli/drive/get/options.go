package get

import "io"

// Options defines the configuration for the drive get operation.
type Options struct {
	// IdentityID is the specific account to query the active drive for.
	IdentityID string
	// DriveRef identifies the target drive by ID, name, or alias.
	DriveRef string
	// Stdout is the destination for the formatted drive details.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
