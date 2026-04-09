package get

import "io"

// Options defines the configuration for the drive get operation.
type Options struct {
	// DriveRef identifies the target drive by ID, name, or alias.
	DriveRef string
	// Stdout is the destination for the formatted drive details.
	Stdout io.Writer
}

func (o *Options) Validate() error {
	return nil
}
