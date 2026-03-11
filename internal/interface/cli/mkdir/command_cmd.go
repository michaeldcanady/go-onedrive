// Package mkdir provides the command-line interface for creating directories in OneDrive.
package mkdir

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// Command handles the execution logic for the 'mkdir' command.
type Command struct {
	util.BaseCommand
}

// NewCmd creates a new Command instance with the provided dependency container.
func NewCmd(container didomain.Container) *Command {
	return &Command{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the mkdir command, creating a new directory at the specified path.
// It uses the domainfs.Writer interface to decouple from the full filesystem service.
func (c *Command) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting mkdir command",
		domainlogger.Bool("parent", opts.Parent),
	)

	// Decouple by using the Writer interface instead of the full Service.
	var fsSvc domainfs.Writer = c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	c.Log.Info("initiating directory creation",
		domainlogger.String("path", opts.Path),
	)

	if err := fsSvc.Mkdir(ctx, opts.Path, domainfs.MKDirOptions{Parents: opts.Parent}); err != nil {
		c.Log.Error("failed to create directory",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to create directory", err)
	}

	c.Log.Info("mkdir completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "created \"%s\"", opts.Path)

	return nil
}
