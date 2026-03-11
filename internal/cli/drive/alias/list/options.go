package list

import "io"

// Options defines the configuration for the drive alias list command.
// It encapsulates the standard I/O streams.
type Options struct {
	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for displaying the alias table.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
