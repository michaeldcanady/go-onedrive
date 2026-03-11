package remove

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type RemoveCmd struct {
	util.BaseCommand
}

func NewRemoveCmd(container didomain.Container) *RemoveCmd {
	return &RemoveCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *RemoveCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("removing drive alias",
		domainlogger.String("alias", opts.Alias),
	)

	if err := c.Container.State().RemoveDriveAlias(opts.Alias); err != nil {
		c.Log.Error("failed to remove drive alias",
			domainlogger.Error(err),
		)
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to remove drive alias", err)
	}

	fmt.Fprintf(opts.Stdout, "Alias %q removed\n", opts.Alias)

	c.Log.Info("drive alias remove completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
