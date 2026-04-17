package use

import "io"

// Options defines the configuration for the drive use operation.
type Options struct {
	// DriveRef identifies the target drive by ID, name, or alias.
	DriveRef string
	// IdentityID is the specific account to scope this drive selection to (optional).
	IdentityID string
	// Stdout is the destination for the operation's output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
