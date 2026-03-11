package delete

import (
	"io"
)

type Options struct {
	Name  string
	Force bool

	Stdout io.Writer
	Stderr io.Writer
}
