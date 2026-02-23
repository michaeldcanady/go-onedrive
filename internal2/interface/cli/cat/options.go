package cat

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// Options defines the configuration for the cat command.
type Options struct {
	// Path is the OneDrive path to the file to be read.
	Path string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for the command.
	Stdout io.Writer
	// Stderr is the error stream for the command.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Path == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "path is required")
	}
	return nil
}
