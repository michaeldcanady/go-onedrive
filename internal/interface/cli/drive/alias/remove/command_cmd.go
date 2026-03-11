// Package remove provides the command-line interface for deleting OneDrive drive aliases.
package remove

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

// RemoveCmd handles the execution logic for the 'drive alias remove' command.
type RemoveCmd struct {
	util.BaseCommand
}

// NewRemoveCmd creates a new RemoveCmd instance with the provided dependency container.
func NewRemoveCmd(container didomain.Container) *RemoveCmd {
	return &RemoveCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the drive alias remove command. It deletes a user-defined alias
// from the global state. It uses the domainstate.Service interface to decouple
// from the full container.
func (c *RemoveCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("removing drive alias",
		domainlogger.String("alias", opts.Alias),
	)

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	if err := stateSvc.RemoveDriveAlias(opts.Alias); err != nil {
		c.Log.Error("failed to remove drive alias",
			domainlogger.String("alias", opts.Alias),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to remove drive alias", err)
	}

	c.RenderSuccess(opts.Stdout, "alias %q removed", opts.Alias)

	c.Log.Info("drive alias remove completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
