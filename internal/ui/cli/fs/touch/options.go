package touch

import (
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
)

// Options provides the settings for the drive touch command.
type Options struct {
	// Path is the filesystem path of the file to touch.
	Path string
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}
