package cat

import (
	"context"

	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
)

type CommandContext struct {
	Ctx     context.Context
	Options Options
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
}
