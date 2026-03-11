package current

import (
	"context"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

type CurrentCmd struct {
	util.BaseCommand
}

func NewCurrentCmd(container didomain.Container) *CurrentCmd {
	return &CurrentCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CurrentCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("retrieving current profile")

	name, err := c.Container.State().Get(domainstate.KeyProfile)
	if err != nil {
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.RenderMessage(opts.Stdout, "%s\n", name)

	c.Log.Info("profile current completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
