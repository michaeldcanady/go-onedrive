package download

import (
	"context"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs"
)

type CommandContext struct {
	Ctx     context.Context
	Options Options

	// SourceURI is the remote filesystem path of the item to download.
	SourceURI *fs.URI

	// DestinationURI is the local filesystem path where the item should be downloaded.
	DestinationURI *fs.URI
}
