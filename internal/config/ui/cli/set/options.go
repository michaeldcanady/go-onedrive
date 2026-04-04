package set

import (
	"errors"
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
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Key == "" {
		return errors.New("configuration key is required")
	}
	return nil
}
