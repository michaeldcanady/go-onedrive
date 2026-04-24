package list

import (
	"io"
)

type Options struct {
	Format string
	Stdout io.Writer
	Stderr io.Writer
}
