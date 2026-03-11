// Package mv provides the command-line interface for moving or renaming items in OneDrive.
package mv

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// Command handles the execution logic for the 'mv' command.
type Command struct {
	util.BaseCommand
}

// NewCmd creates a new Command instance with the provided dependency container.
func NewCmd(container didomain.Container) *Command {
	return &Command{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the mv command, moving or renaming an item from a source path to a destination path.
// It uses the domainfs.Manager interface to decouple from the full filesystem service.
func (c *Command) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting mv command",
		domainlogger.String("src", opts.Source),
		domainlogger.String("dst", opts.Destination),
	)

	// Decouple by using the Manager interface instead of the full Service.
	var fsSvc domainfs.Manager = c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if err := fsSvc.Move(ctx, opts.Source, opts.Destination, domainfs.MoveOptions{}); err != nil {
		c.Log.Error("failed to move item",
			domainlogger.String("src", opts.Source),
			domainlogger.String("dst", opts.Destination),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to move item", err)
	}

	c.Log.Info("mv completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "moved \"%s\" to \"%s\"", opts.Source, opts.Destination)

	return nil
}
