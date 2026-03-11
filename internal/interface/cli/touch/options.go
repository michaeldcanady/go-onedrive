package touch

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// Options defines the configuration for the touch command.
// It encapsulates the target path and standard I/O streams.
type Options struct {
	// Path is the OneDrive path of the file to be touched.
	Path string

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
