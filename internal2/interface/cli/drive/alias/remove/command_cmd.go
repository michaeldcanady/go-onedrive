package remove

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type RemoveCmd struct {
	util.BaseCommand
}

func NewRemoveCmd(container di.Container) *RemoveCmd {
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
		logger.String("alias", opts.Alias),
	)

	if err := c.Container.State().RemoveDriveAlias(opts.Alias); err != nil {
		c.Log.Error("failed to remove drive alias",
			logger.Error(err),
		)
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to remove drive alias", err)
	}

	fmt.Fprintf(opts.Stdout, "Alias %q removed\n", opts.Alias)

	c.Log.Info("drive alias remove completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}
