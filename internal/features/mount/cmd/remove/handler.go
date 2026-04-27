package remove

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

type MountRemover interface {
	// RemoveMount removes a mount point from the configuration.
	RemoveMount(ctx context.Context, path string) error
}

type URIFactory interface {
	FromString(input string) (*fs.URI, error)
}

type Command struct {
	mountSvc   MountRemover
	uriFactory URIFactory
	log        logger.Logger
}

func NewCommand(mountSvc MountRemover, uriFactory URIFactory, l logger.Logger) *Command {
	return &Command{
		mountSvc:   mountSvc,
		uriFactory: uriFactory,
		log:        l,
	}
}

func (c *Command) Validate(ctx *CommandContext) error {
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("failed to parse uri %s: %w", ctx.Options.Path, err)
	}

	ctx.Uri = uri

	return nil
}

func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)
	log.Info("starting mount remove operation")

	if err := c.mountSvc.RemoveMount(ctx.Ctx, ctx.Uri.String()); err != nil {
		return fmt.Errorf("failed to remove mount %s: %w", ctx.Uri.String(), err)
	}

	log.Info("mount remove completed successfully")
	return nil
}

func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
