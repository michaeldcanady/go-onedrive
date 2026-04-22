package rm

import (
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs"
)

// Options provides the settings for the drive rm command.
type Options struct {
	// Path is the filesystem path of the item to remove.
	Path string
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}
