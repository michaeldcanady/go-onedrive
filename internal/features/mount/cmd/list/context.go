package list

import (
	"context"

	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
)

type CommandContext struct {
	Ctx     context.Context
	Options *Options
	Format  formatting.Format
}

func NewCommandContext(ctx context.Context, opts *Options) *CommandContext {
	return &CommandContext{
		Ctx:     ctx,
		Options: opts,
	}
}
