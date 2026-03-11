// Package rm provides the command-line interface for removing items in OneDrive.
package rm

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// Command handles the execution logic for the 'rm' command.
type Command struct {
	util.BaseCommand
}

// NewCmd creates a new Command instance with the provided dependency container.
func NewCmd(container didomain.Container) *Command {
	return &Command{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the rm command, removing an item at the specified path.
// It uses the domainfs.Manager interface to decouple from the full filesystem service.
func (c *Command) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting rm command",
		domainlogger.String("path", opts.Path),
		domainlogger.Bool("recursive", opts.Recursive),
	)

	// Decouple by using the Manager interface instead of the full Service.
	var fsSvc domainfs.Manager = c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if err := fsSvc.Remove(ctx, opts.Path, domainfs.RemoveOptions{Recursive: opts.Recursive}); err != nil {
		c.Log.Error("failed to remove item",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to remove item", err)
	}

	c.Log.Info("rm completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "removed \"%s\"", opts.Path)

	return nil
}
