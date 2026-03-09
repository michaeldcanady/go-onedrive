package remove

import (
	"errors"
	"io"
)

type Options struct {
	Alias  string
	Stdout io.Writer
	Stderr io.Writer
}

func (o Options) Validate() error {
	if o.Alias == "" {
		return errors.New("alias is required")
	}
	return nil
}
