package add

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

type MountAdder interface {
	// AddMount adds or updates a mount point in the configuration.
	AddMount(ctx context.Context, m mount.MountConfig) error
}

type AccountGetter interface {
	GetAccount(ctx context.Context, identityID string) (*identity.Account, error)
}

type URIFactory interface {
	FromString(input string) (*fs.URI, error)
}

type Command struct {
	log        logger.Logger
	mountSvc   MountAdder
	accountSvc AccountGetter
	uriFactory URIFactory
}

func NewCommand(mountSvc MountAdder, accountSvc AccountGetter, uriFactory URIFactory, l logger.Logger) *Command {
	return &Command{
		log:        l,
		mountSvc:   mountSvc,
		accountSvc: accountSvc,
		uriFactory: uriFactory,
	}
}

func mountOptions(_ *Command, ctx *CommandContext) error {
	if len(ctx.Options.MountOptions) <= 0 {
		return nil
	}

	for _, opt := range ctx.Options.MountOptions {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid option format: %s. Expected key=value", opt)
		}
		ctx.MountOptions[parts[0]] = parts[1]
	}
	return nil
}

func parseURI(cmd *Command, ctx *CommandContext) error {
	uri, err := cmd.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("failed to parse uri %s: %w", ctx.Options.Path, err)
	}
	ctx.Uri = uri

	return err
}

func validateIdentity(cmd *Command, ctx *CommandContext) error {
	account, err := cmd.accountSvc.GetAccount(ctx.Ctx, ctx.Options.IdentityID)
	if err != nil {
		return fmt.Errorf("invalid account %s: %w", ctx.Options.IdentityID, err)
	}
	ctx.Identity = account
	return nil
}

func (c *Command) Validate(ctx *CommandContext) error {
	if err := mountOptions(c, ctx); err != nil {
		return err
	}

	if err := parseURI(c, ctx); err != nil {
		return err
	}
	ctx.Type = ctx.Options.Type

	if err := validateIdentity(c, ctx); err != nil {
		return err
	}

	return nil
}

func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)
	log.Info("starting mount add operation")

	if err := c.mountSvc.AddMount(ctx.Ctx, mount.MountConfig{
		Path:       ctx.Uri.String(),
		Type:       ctx.Type,
		IdentityID: ctx.Identity.ID,
		Options:    ctx.MountOptions,
	}); err != nil {
		log.Error("mount add failed", logger.Error(err))
		return fmt.Errorf("failed to add mount: %w", err)
	}

	log.Info("mount add completed successfully")
	return nil
}

func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
