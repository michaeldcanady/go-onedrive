package mkdir

import "io"

// Options defines the configuration for the mkdir command.
// It encapsulates the flags that control the mkdir behavior.
type Options struct {
	// Path
	Path string
	// Parent
	Parent bool
	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for the command.
	Stdout io.Writer
	// Stderr is the error stream for the command.
	Stderr io.Writer
}

func (o *Options) Validate() error {
	return nil
}
