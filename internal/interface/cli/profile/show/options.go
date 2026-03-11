package show

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// Options defines the configuration for the profile show command.
// It encapsulates the name of the profile to display and standard I/O streams.
type Options struct {
	// Name is the unique name of the profile whose details will be shown.
	Name string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for displaying profile details.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if the required Name is missing.
func (o *Options) Validate() error {
	if o.Name == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "name is required")
	}
	return nil
}
