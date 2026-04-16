package cp

import (
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
)

// Options provides the settings for the drive cp command.
type Options struct {
	// Source is the filesystem path of the item to copy.
	Source string

	// SourceURI is the filesystem path of the item to copy.
	SourceURI *fs.URI

	// Destination is the filesystem path where the item should be copied.
	Destination string

	// DestinationURI is the filesystem path where the item should be copied.
	DestinationURI *fs.URI

	// Recursive determines whether to copy directories and their contents.
	Recursive bool

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}
