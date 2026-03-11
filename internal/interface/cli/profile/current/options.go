package current

import "io"

// Options defines the configuration for the profile current command.
// It encapsulates the standard I/O streams.
type Options struct {
	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for displaying the current profile name.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}
