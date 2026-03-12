package login

import "io"

// Options provides the user-facing settings for the auth login command.
type Options struct {
	// ShowToken determines whether the acquired access token is printed to stdout.
	ShowToken bool
	// Force specifies whether to re-authenticate regardless of existing credentials.
	Force bool
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
	// Stderr is the destination for error messages.
	Stderr io.Writer
}
