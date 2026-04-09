package use

import "io"

// Options defines the configuration for the drive use operation.
type Options struct {
	// DriveRef identifies the target drive by ID, name, or alias.
	DriveRef string
	// Stdout is the destination for the operation's output messages.
	Stdout io.Writer
}

func (o *Options) Validate() error {
	return nil
}
