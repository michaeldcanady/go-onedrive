package login

import "io"

// Options defines the configuration for the login command.
// It encapsulates the interactive flags and standard I/O streams.
type Options struct {
	// Force indicates whether to force a new login, even if a valid token exists.
	Force bool

	// ShowToken indicates whether to print the acquired access token to stdout.
	ShowToken bool

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages and token display.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}
