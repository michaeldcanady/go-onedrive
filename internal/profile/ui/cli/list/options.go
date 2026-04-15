package list

import "io"

// Options provides the user-facing settings for the profile list command.
type Options struct {
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
