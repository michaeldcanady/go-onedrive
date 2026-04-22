package list

import (
	"io"
	"os"

	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
)

type Options struct {
	Format string
	Stdout io.Writer
}

func NewOptions() *Options {
	return &Options{
		Format: formatting.FormatShort.String(),
		Stdout: os.Stdout,
	}
}

func (o *Options) Validate() error {
	return nil
}
