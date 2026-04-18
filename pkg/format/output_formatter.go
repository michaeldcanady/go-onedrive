package formatting

import (
	"io"
)

// OutputFormatter defines the interface for rendering a collection of items to a destination writer.
type OutputFormatter interface {
	// Format processes the items and writes the formatted representation to the provided writer.
	Format(w io.Writer, items []any) error
}
