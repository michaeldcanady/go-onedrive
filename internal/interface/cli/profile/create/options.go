package create

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// Options defines the configuration for the profile create command.
// It encapsulates the profile name, flags for overwriting and selecting the profile,
// and standard I/O streams.
type Options struct {
	// Name is the unique name for the new profile.
	Name string

	// Force indicates whether to overwrite the profile if it already exists.
	Force bool

	// SetCurrent indicates whether to set the newly created profile as the active one.
	SetCurrent bool

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
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
