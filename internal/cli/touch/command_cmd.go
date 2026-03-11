// Package touch provides the command-line interface for creating empty files or updating timestamps in OneDrive.
package touch

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// Command handles the execution logic for the 'touch' command.
type Command struct {
	util.BaseCommand
}

// NewCmd creates a new Command instance with the provided dependency container.
func NewCmd(container didomain.Container) *Command {
	return &Command{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the touch command, creating a new empty file or updating the timestamp
// of an existing one. It uses the domainfs.Writer interface to decouple from the full
// filesystem service.
func (c *Command) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting touch command",
		domainlogger.String("path", opts.Path),
	)

	// Decouple by using the Writer interface instead of the full Service.
	var fsSvc domainfs.Writer = c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if _, err := fsSvc.Touch(ctx, opts.Path, domainfs.TouchOptions{}); err != nil {
		c.Log.Error("failed to touch file",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to touch file", err)
	}

	c.Log.Info("touch completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "touched \"%s\"", opts.Path)

	return nil
}
