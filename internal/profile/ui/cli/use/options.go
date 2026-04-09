package use

import (
	"errors"
	"io"
	"strings"
)

// Options provides the settings for the profile use command.
type Options struct {
	// Name is the name of the profile to switch to.
	Name string

	// Stdout is the destination for standard output messages.
	Stdout io.Writer

	// Stderr is the destination for error messages.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	o.Name = strings.TrimSpace(o.Name)
	if o.Name == "" {
		return errors.New("profile name is required")
	}
	return nil
}
