package list

import (
	"io"
)

// Options defines the configuration for the drive list operation.
type Options struct {
	// Stdout is the writer for standard output.
	Stdout io.Writer
	// Stderr is the writer for error output.
	Stderr io.Writer
}

func (o *Options) Validate() error {
	return nil
}
