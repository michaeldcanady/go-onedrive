package list

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
)

type MountLister interface {
	ListMounts(ctx context.Context) ([]mount.MountConfig, error)
}

type FormatCreator interface {
	Create(format formatting.Format) (formatting.OutputFormatter, error)
}

type Command struct {
	mounts           MountLister
	formatterFactory FormatCreator
	log              logger.Logger
}

func NewCommand(m MountLister, ff FormatCreator, l logger.Logger) *Command {
	return &Command{
		mounts:           m,
		formatterFactory: ff,
		log:              l,
	}
}

func (c *Command) Validate(ctx *CommandContext) error {
	ctx.Format = formatting.NewFormat(ctx.Options.Format)
	return nil
}

func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)
	log.Info("starting mount list operation")

	items, err := c.mounts.ListMounts(ctx.Ctx)
	if err != nil {
		log.Error("list failed", logger.Error(err))
		return err
	}

	formatter, err := c.formatterFactory.Create(ctx.Format)
	if err != nil {
		log.Error("failed to create formatter", logger.Error(err))
		return err
	}

	itemsAny := make([]any, len(items))
	for i, item := range items {
		itemsAny[i] = item
	}

	if err := formatter.Format(ctx.Options.Stdout, itemsAny); err != nil {
		log.Error("format failed", logger.Error(err))
		return err
	}

	log.Info("mount list completed successfully")
	return nil
}

func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
