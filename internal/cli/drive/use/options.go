package use

import "io"

// Options defines the configuration for the drive use command.
// It encapsulates the target drive identifier and standard I/O streams.
type Options struct {
	// DriveIDOrAlias is the ID or registered alias of the drive to switch to.
	DriveIDOrAlias string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}
