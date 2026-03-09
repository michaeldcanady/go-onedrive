package set

import (
	"errors"
	"io"
)

type Options struct {
	Alias   string
	DriveID string
	Stdout  io.Writer
	Stderr  io.Writer
}

func (o Options) Validate() error {
	if o.Alias == "" {
		return errors.New("alias is required")
	}
	if o.DriveID == "" {
		return errors.New("drive-id is required")
	}
	return nil
}
