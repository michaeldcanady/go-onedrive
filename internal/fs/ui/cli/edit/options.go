package edit

import (
	"io"
	"os"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options defines the configuration for the drive edit operation.
type Options struct {
	// Path is the path to the item to edit.
	Path string `arg:"1"`
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
	// Force indicates whether to overwrite the item if it already exists.
	Force bool `flag:"force,short=f,desc='Indicates whether to overwrite the item if it already exists.',default=false"`
	// Editor is the editor to use for editing the file.
	Editor string `flag:"editor,desc='Editor to use (overrides config and environment)',default="`
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
	// No additional validation needed for edit options at this time
	return nil
}
