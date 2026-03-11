// Package download provides the command-line interface for downloading files from OneDrive.
package download

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
)

// Options defines the configuration for the download command.
type Options struct {
	// Source is the path in OneDrive to download.
	Source string

	// Destination is the local path where the file should be saved.
	Destination string

	// Overwrite determines whether to replace an existing local file.
	Overwrite bool

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for the command.
	Stdout io.Writer
	// Stderr is the error stream for the command.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Source == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "source path is required")
	}
	if o.Destination == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "destination path is required")
	}
	return nil
}
