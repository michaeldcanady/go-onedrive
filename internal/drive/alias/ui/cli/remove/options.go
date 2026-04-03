package remove

import "io"

// Options defines the configuration for the drive alias remove operation.
type Options struct {
	Alias  string
	Stdout io.Writer
}
