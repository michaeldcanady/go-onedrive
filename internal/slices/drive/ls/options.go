package ls

import (
	"io"
)

// Options provides the user-facing settings for the drive ls command.
type Options struct {
	// Path is the filesystem path to list.
	Path string
	// Recursive determines whether to list subdirectories.
	Recursive bool
	// Format is the output format (e.g., "short", "long", "json", "tree").
	Format string
	// SortField is the field to sort by (e.g., "Name", "Size", "ModifiedAt").
	SortField string
	// SortDescending determines whether to sort in descending order.
	SortDescending bool
	// All determines whether to include hidden files.
	All bool
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	// Path can be empty (defaults to current directory or root)
	return nil
}
