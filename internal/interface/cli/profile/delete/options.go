package delete

import "io"

// Options defines the configuration for the profile delete command.
// It encapsulates the profile name, force flag, and standard I/O streams.
type Options struct {
	// Name is the unique name of the profile to be deleted.
	Name string

	// Force indicates whether to bypass confirmation when deleting the active profile.
	Force bool

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages and prompts.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}
