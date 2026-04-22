package upload

import (
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
)

// Options provides the settings for the drive upload command.
type Options struct {
	// Source is the local filesystem path of the item to upload.
	Source string

	// SourceURI is the local filesystem path of the item to upload.
	SourceURI *fs.URI

	// Destination is the remote filesystem path where the item should be uploaded.
	Destination string

	// DestinationURI is the remote filesystem path where the item should be uploaded.
	DestinationURI *fs.URI

	// Recursive determines whether to upload directories and their contents.
	Recursive bool

	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}
