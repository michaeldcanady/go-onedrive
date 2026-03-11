package set

import (
	"context"
	"fmt"
	"time"

	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type SetCmd struct {
	util.BaseCommand
}

func NewSetCmd(container di.Container) *SetCmd {
	return &SetCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *SetCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("setting drive alias",
		logger.String("alias", opts.Alias),
		logger.String("drive_id", opts.DriveID),
	)

	if err := c.Container.State().SetDriveAlias(opts.Alias, opts.DriveID); err != nil {
		c.Log.Error("failed to set drive alias",
			logger.Error(err),
		)
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to set drive alias", err)
	}

	fmt.Fprintf(opts.Stdout, "Alias %q set to drive %q\n", opts.Alias, opts.DriveID)

	c.Log.Info("drive alias set completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}
