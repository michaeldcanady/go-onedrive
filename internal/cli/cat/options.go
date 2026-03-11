// Package cat provides the command-line interface for reading OneDrive file contents.
package cat

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
)

// Options defines the configuration for the cat command.
// It encapsulates the path of the file to be read and the standard I/O streams.
type Options struct {
	// Path is the OneDrive path to the file whose contents will be read.
	Path string

	// Stdin is the input stream for the command, primarily used for piped input.
	Stdin io.Reader
	// Stdout is the output stream where the file contents will be written.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if the required Path is missing.
func (o *Options) Validate() error {
	if o.Path == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "path is required")
	}
	return nil
}
