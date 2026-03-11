package use

import "io"

// Options defines the configuration for the profile use command.
// It encapsulates the target profile name and standard I/O streams.
type Options struct {
	// Name is the unique name of the profile to switch to.
	Name string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}
