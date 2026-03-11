package touch

import (
	"context"
	"time"

	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Command struct {
	util.BaseCommand
}

// NewCmd creates a new Command instance with the provided dependency container.
func NewCmd(container di.Container) *Command {
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
		logger.String("path", opts.Path),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if _, err := fsSvc.Touch(ctx, opts.Path, fs.TouchOptions{}); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to touch file", err)
	}

	c.Log.Info("touch completed successfully",
		logger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "Successfully touched \"%s\"\n", opts.Path)

	return nil
}
