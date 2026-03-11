package cp

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
)

// Options defines the configuration for the cp command.
// It encapsulates the source and destination paths, overwrite behavior,
// recursive copying for directories, and standard I/O streams.
type Options struct {
	// Source is the OneDrive path of the item to be copied.
	Source string

	// Dest is the OneDrive path where the item will be copied to.
	Dest string

	// Overwrite determines whether to overwrite an existing item at the destination.
	Overwrite bool

	// Recursive determines whether to copy subdirectories and their contents.
	// Required when copying folders.
	Recursive bool

	// IgnoreFile is the path to a .gitignore-style file for filtering items during recursive copy.
	IgnoreFile string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for successful operation messages.
	Stdout io.Writer
	// Stderr is the error stream for reporting command-specific issues.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if either Source or Dest is missing.
func (o *Options) Validate() error {
	if o.Source == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "source is required")
	}
	if o.Dest == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "destination is required")
	}
	return nil
}
