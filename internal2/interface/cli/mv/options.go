package mv

import (
	"errors"
	"io"
)

type Options struct {
	Source      string
	Destination string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
}

func (o Options) Validate() error {
	if o.Source == "" {
		return errors.New("source path is required")
	}
	if o.Destination == "" {
		return errors.New("destination path is required")
	}
	return nil
}
