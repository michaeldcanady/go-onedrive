package logout

import "io"

// Options provides the user-facing settings for the auth logout command.
type Options struct {
	// Force specifies whether to clear all cached credentials for the profile.
	Force bool
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
	// Stderr is the destination for error messages.
	Stderr io.Writer
}

func (o *Options) Validate() error {
	return nil
}
