// Package use provides the command-line interface for selecting the active OneDrive drive.
package use

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

// UseCmd handles the execution logic for the 'drive use' command.
type UseCmd struct {
	util.BaseCommand
}

// NewUseCmd creates a new UseCmd instance with the provided dependency container.
func NewUseCmd(container didomain.Container) *UseCmd {
	return &UseCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the drive use command. It resolves the provided drive ID or alias
// and updates the global state to make it the active drive for future operations.
// It uses specific domain services to decouple from the full container.
func (c *UseCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting drive use command", domainlogger.String("target", opts.DriveIDOrAlias))

	driveSvc := c.Container.Drive()
	if driveSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "drive service is nil")
	}

	resolvedDrive, err := driveSvc.ResolveDrive(ctx, opts.DriveIDOrAlias)
	if err != nil {
		c.Log.Error("failed to resolve drive",
			domainlogger.Error(err),
			domainlogger.String("target", opts.DriveIDOrAlias),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	if err := stateSvc.Set(domainstate.KeyDrive, resolvedDrive.ID, domainstate.ScopeGlobal); err != nil {
		c.Log.Error("failed to update current drive state",
			domainlogger.Error(err),
			domainlogger.String("driveID", resolvedDrive.ID),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.RenderSuccess(opts.Stdout, "now using drive: %s (%s)", resolvedDrive.Name, resolvedDrive.ID)

	c.Log.Info("drive use completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
