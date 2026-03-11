package touch

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type Command struct {
	util.BaseCommand
}

// NewCmd creates a new Command instance with the provided dependency container.
func NewCmd(container didomain.Container) *Command {
	return &Command{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *Command) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting touch command",
		domainlogger.String("path", opts.Path),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if _, err := fsSvc.Touch(ctx, opts.Path, domainfs.TouchOptions{}); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to touch file", err)
	}

	c.Log.Info("touch completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully touched \"%s\"\n", opts.Path)

	return nil
}
