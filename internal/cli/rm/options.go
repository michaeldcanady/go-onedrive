package rm

import (
	"errors"
	"io"
)

type Options struct {
	Path      string
	Recursive bool
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
}

func (o Options) Validate() error {
	if o.Path == "" {
		return errors.New("path is required")
	}
	return nil
}
