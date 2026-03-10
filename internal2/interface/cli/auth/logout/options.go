package logout

import (
	"io"
)

type Options struct {
	Force bool

	Stdout io.Writer
	Stderr io.Writer
}
