package edit

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
)

// Options defines the configuration for the edit command.
// It encapsulates the target path and standard I/O streams.
type Options struct {
	// Path is the OneDrive path of the file to be edited.
	Path string
	// Force force upload even if it differs from the expected version.
	Force bool
	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
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
