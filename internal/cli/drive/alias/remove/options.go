package remove

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
)

// Options defines the configuration for the drive alias remove command.
// It encapsulates the name of the alias to be deleted and standard I/O streams.
type Options struct {
	// Alias is the friendly name of the drive alias to remove.
	Alias string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if the required Alias is missing.
func (o *Options) Validate() error {
	if o.Alias == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "alias is required")
	}
	return nil
}
