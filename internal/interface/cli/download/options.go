package download

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// Options defines the configuration for the download command.
// It encapsulates the remote source path, local destination path, and standard I/O streams.
type Options struct {
	// Source is the OneDrive path of the file to be downloaded.
	Source string

	// Overwrite specifies whether to overwrite the destination if it exists.
	Overwrite bool

	// Destination is the local filesystem path where the file will be saved.
	Destination string

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
