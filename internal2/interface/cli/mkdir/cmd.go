package mkdir

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

	c.Log.Info("starting mkdir command",
		logger.Bool("parent", opts.Parent),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	c.Log.Info("initiating directory creation",
		logger.String("path", opts.Path),
	)

	if err := fsSvc.Mkdir(ctx, opts.Path, fs.MKDirOptions{Parents: opts.Parent}); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to create directory", err)
	}

	c.Log.Info("mkdir completed successfully",
		logger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "Successfully created \"%s\"\n", opts.Path)

	return nil
}
