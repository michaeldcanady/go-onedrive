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

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Validate() error {
	return nil
}
