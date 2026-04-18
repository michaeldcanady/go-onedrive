package use

import (
	"io"
)

// Options provides the settings for the profile use command.
type Options struct {
	// Name is the name of the profile to switch to.
	Name string

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}
