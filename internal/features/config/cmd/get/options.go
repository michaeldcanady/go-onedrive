package get

import "io"

// Options provides the settings for the config get command.
type Options struct {
	// Key is the configuration key to retrieve (optional).
	Key string
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
