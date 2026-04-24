package remove

import (
	"io"
)

type Options struct {
	Stdout io.Writer
	Stderr io.Writer

	// Path the mount point's path
	Path string
}
