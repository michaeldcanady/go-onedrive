package formatting

import (
	"io"
)

type OutputFormatter[T any] interface {
	Format(w io.Writer, items []T) error
}
