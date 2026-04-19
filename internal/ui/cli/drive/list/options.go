package list

import (
	"io"
	"os"
)

// Options defines the configuration for the drive list operation.
type Options struct {
	// IdentityID is the specific account to list drives for.
	IdentityID string
	// Stdout is the writer for standard output.
	Stdout io.Writer
	// Stderr is the writer for error output.
	Stderr io.Writer
}

// NewOptions creates a new instance of Options with default values.
func NewOptions() Options {
	return Options{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
