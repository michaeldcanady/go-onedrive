package delete

import (
	"io"
)

// Options provides the settings for the profile delete command.
type Options struct {
	// Name is the name of the profile to delete.
	Name string

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}
