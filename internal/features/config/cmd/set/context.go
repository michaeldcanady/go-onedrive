package set

import (
	"context"
)

type CommandContext struct {
	Ctx     context.Context
	Options *Options
}

func NewCommandContext(ctx context.Context, opts *Options) *CommandContext {
	return &CommandContext{
		Ctx:     ctx,
		Options: opts,
	}
}
