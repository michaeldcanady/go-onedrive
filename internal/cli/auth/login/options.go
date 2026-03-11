package login

import (
	"io"
)

type Options struct {
	ShowToken bool
	Force     bool

	Stdout io.Writer
	Stderr io.Writer
}
