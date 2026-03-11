package list

import (
	"io"
)

type Options struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (o Options) Validate() error {
	return nil
}
