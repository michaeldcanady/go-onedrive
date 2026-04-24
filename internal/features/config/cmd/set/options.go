package set

import (
	"io"
)

// Options provides the settings for the config set command.
type Options struct {
	// Key is the configuration key to set.
	Key string
	// Value is the configuration value to set.
	Value string
	// Stdout is the destination for standard output messages.
	Stdout io.Writer

	Stderr io.Writer
}

func NewOptions() *Options {
	return &Options{}
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
