package set

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type SetCmd struct {
	util.BaseCommand
}

func NewSetCmd(container didomain.Container) *SetCmd {
	return &SetCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *SetCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("setting drive alias",
		domainlogger.String("alias", opts.Alias),
		domainlogger.String("drive_id", opts.DriveID),
	)

	if err := c.Container.State().SetDriveAlias(opts.Alias, opts.DriveID); err != nil {
		c.Log.Error("failed to set drive alias",
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to set drive alias", err)
	}

	c.RenderSuccess(opts.Stdout, "alias %q set to drive %q", opts.Alias, opts.DriveID)

	c.Log.Info("drive alias set completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
