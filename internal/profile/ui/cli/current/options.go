package current

import "io"

// Options provides settings for the profile current command.
type Options struct {
	// Stdout is the destination for standard output messages.
	Stdout io.Writer

	// Stderr is the destination for error messages.
	Stderr io.Writer
}

func (o *Options) Validate() error {
	return nil
}
