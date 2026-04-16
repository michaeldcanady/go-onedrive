package create

import (
	"io"
)

// Options provides the settings for the profile create command.
type Options struct {
	// Name is the name of the profile to create.
	Name string `arg:"1"`

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}
