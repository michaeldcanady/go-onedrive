package current

import (
	"io"
)

type Options struct {
	Stdout io.Writer
	Stderr io.Writer
}
