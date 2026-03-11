package logout

import "io"

// Options defines the configuration for the logout command.
// It encapsulates the interactive flags and standard I/O streams.
type Options struct {
	// Force indicates whether to force logout even if errors occur during session cleanup.
	Force bool

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}
