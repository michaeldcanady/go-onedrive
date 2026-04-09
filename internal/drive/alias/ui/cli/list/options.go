package list

import "io"

// Options defines the configuration for the drive alias list operation.
type Options struct {
	// Stdout is the writer for standard output.
	Stdout io.Writer
}

func (o *Options) Validate() error {
	return nil
}
