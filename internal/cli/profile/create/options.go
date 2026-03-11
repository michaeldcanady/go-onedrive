package create

import (
	"io"
)

type Options struct {
	Name       string
	SetCurrent bool
	Force      bool

	Stdout io.Writer
	Stderr io.Writer
}
