package remove

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

type CommandContext struct {
	Ctx     context.Context
	Options *Options

	// Uri the mount point's path
	Uri *fs.URI
}

func NewCommandContext(ctx context.Context, opts *Options) *CommandContext {
	return &CommandContext{
		Ctx:     ctx,
		Options: opts,
	}
}
