package get

import (
	"io"
)

type Options struct {
	DriveIDOrAlias string

	Stdout io.Writer
	Stderr io.Writer
}
