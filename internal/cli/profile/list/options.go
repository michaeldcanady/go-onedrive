package list

import "io"

// Options defines the configuration for the profile list command.
// It encapsulates the standard I/O streams.
type Options struct {
	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for displaying the profile list.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}
