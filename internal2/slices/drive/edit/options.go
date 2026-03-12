package edit

import (
	"io"
	"os"
)

// Options defines the configuration for the drive edit operation.
type Options struct {
	// Path is the path to the item to edit.
	Path string
	// Force indicates whether to overwrite the item if it already exists.
	Force bool
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
