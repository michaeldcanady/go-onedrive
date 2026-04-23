package add

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
)

type MountAdder interface {
	// AddMount adds or updates a mount point in the configuration.
	AddMount(ctx context.Context, m mount.MountConfig) error
}

type AccountGetter interface {
	GetAccount(ctx context.Context, identityID string) (*identity.Account, error)
}

type Command struct {
	log        logger.Logger
	mountSvc   MountAdder
	accountSvc AccountGetter
}

func NewCommand(mountSvc MountAdder, accountSvc AccountGetter, l logger.Logger) *Command {
	return &Command{
		log:        l,
		mountSvc:   mountSvc,
		accountSvc: accountSvc,
	}
}

func (c *Command) Validate(ctx *CommandContext) error {
	for _, opt := range ctx.Options.MountOptions {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid option format: %s. Expected key=value", opt)
		}
		ctx.MountOptions[parts[0]] = parts[1]
	}

	if account, err := c.accountSvc.GetAccount(ctx.Ctx, ctx.Options.IdentityID); err != nil {
		return fmt.Errorf("invalid account %s: %w", ctx.Options.IdentityID, err)
	} else {
		ctx.Identity = account
	}

	return nil
}

func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)
	log.Info("starting mount add operation")

	c.mountSvc.AddMount(ctx.Ctx, mount.MountConfig{
		Path:       ctx.Uri.String(),
		Type:       ctx.Type,
		IdentityID: ctx.Identity.ID,
		Options:    ctx.MountOptions,
	})

	log.Info("mount add completed successfully")
	return nil
}

func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
