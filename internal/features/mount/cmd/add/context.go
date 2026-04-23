package add

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

type CommandContext struct {
	Ctx          context.Context
	Options      *Options
	Uri          *fs.URI
	Type         string
	Identity     *identity.Account
	MountOptions map[string]string
}

func NewCommandContext(ctx context.Context, opts *Options) *CommandContext {
	return &CommandContext{
		Ctx:          ctx,
		Options:      opts,
		MountOptions: map[string]string{},
	}
}
