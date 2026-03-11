package upload

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// Options defines the configuration for the upload command.
// It encapsulates the local source path, remote destination path, and standard I/O streams.
type Options struct {
	// Source is the local filesystem path of the item to be uploaded.
	Source string

	// Destination is the OneDrive path where the item will be uploaded.
	Destination string

	// Overwrite determines whether to overwrite an existing item at the destination.
	Overwrite bool

	// Stdin is the input stream for the command.
	Stdin io.Reader

	// Stdout is the output stream for successful operation messages.
	Stdout io.Writer

	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if either the Source or Destination is missing.
func (o *Options) Validate() error {
	if o.Source == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "source is required")
	}
	if o.Destination == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "destination is required")
	}
	return nil
}
