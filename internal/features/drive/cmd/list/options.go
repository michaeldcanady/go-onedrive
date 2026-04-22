package list

import (
	"io"
)

// Options defines the configuration for the drive list operation.
type Options struct {
	// IdentityID is the specific account to list drives for.
	IdentityID string
	// Format is the output format.
	Format string
	// Stdout is the writer for standard output.
	Stdout io.Writer
	// Stderr is the writer for error output.
	Stderr io.Writer
}

// NewOptions creates a new instance of Options with default values.
func NewOptions() Options {
	return Options{}
}
