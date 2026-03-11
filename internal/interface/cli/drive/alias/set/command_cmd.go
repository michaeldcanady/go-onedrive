// Package set provides the command-line interface for creating or updating OneDrive drive aliases.
package set

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

// SetCmd handles the execution logic for the 'drive alias set' command.
type SetCmd struct {
	util.BaseCommand
}

// NewSetCmd creates a new SetCmd instance with the provided dependency container.
func NewSetCmd(container didomain.Container) *SetCmd {
	return &SetCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the drive alias set command. It maps a user-defined alias to
// a specific OneDrive drive ID in the global state.
// It uses the domainstate.Service interface to decouple from the full container.
func (c *SetCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("setting drive alias",
		domainlogger.String("alias", opts.Alias),
		domainlogger.String("drive_id", opts.DriveID),
	)

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	if err := stateSvc.SetDriveAlias(opts.Alias, opts.DriveID); err != nil {
		c.Log.Error("failed to set drive alias",
			domainlogger.String("alias", opts.Alias),
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
