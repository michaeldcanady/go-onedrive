package use

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type UseCmd struct {
	util.BaseCommand
}

func NewUseCmd(container didomain.Container) *UseCmd {
	return &UseCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *UseCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting drive use command", domainlogger.String("target", opts.DriveIDOrAlias))

	resolvedDrive, err := c.Container.Drive().ResolveDrive(ctx, opts.DriveIDOrAlias)
	if err != nil {
		c.Log.Warn("failed to resolve drive",
			domainlogger.Error(err),
			domainlogger.String("target", opts.DriveIDOrAlias),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	if err := c.Container.State().SetCurrentDrive(resolvedDrive.ID); err != nil {
		c.Log.Warn("failed to update current drive state",
			domainlogger.Error(err),
			domainlogger.String("driveID", resolvedDrive.ID),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	fmt.Fprintf(opts.Stdout, "Now using drive: %s (%s)\n", resolvedDrive.Name, resolvedDrive.ID)

	c.Log.Info("drive use completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
