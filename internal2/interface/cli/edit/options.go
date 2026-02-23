// Package edit provides the command-line interface for editing OneDrive files locally.
package edit

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// Options defines the configuration for the edit command.
// It encapsulates the OneDrive path of the file to be edited.
type Options struct {
	// Path is the OneDrive path of the file to edit.
	Path string

	// Force determines whether to overwrite the file even if it changed in the cloud.
	Force bool

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for the command.
	Stdout io.Writer
	// Stderr is the error stream for the command.
	Stderr io.Writer
}

// Validate ensures that the provided options are valid.
// It returns a [util.CommandError] if the required path is missing.
func (o *Options) Validate() error {
	if o.Path == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "file path is required")
	}
	return nil
}
