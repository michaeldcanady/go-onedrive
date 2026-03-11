package list

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type ListCmd struct {
	util.BaseCommand
}

func NewListCmd(container didomain.Container) *ListCmd {
	return &ListCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *ListCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("listing profiles")

	profiles, err := c.Container.Profile().List(ctx)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		c.RenderMessage(opts.Stdout, "%s\n", p.Name)
	}

	c.Log.Info("profile list completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
