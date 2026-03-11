// Package cp provides the command-line interface for the cp operation.
package cp

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
)

// Options defines the configuration for the cp command.
type Options struct {
	Source     string
	Dest       string
	Overwrite  bool
	IgnoreFile string

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for the command.
	Stdout io.Writer
	// Stderr is the error stream for the command.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
// It returns a [util.CommandError] if validation fails.
func (o *Options) Validate() error {
	if o.Source == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "source path is required")
	}
	if o.Dest == "" {
		return util.NewCommandErrorWithNameWithMessage(commandName, "destination path is required")
	}
	return nil
}
