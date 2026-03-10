package create

import (
	"io"
)

type Options struct {
	Name string

	Stdout io.Writer
	Stderr io.Writer
}
