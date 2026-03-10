package rm

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
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

	c.Log.Info("starting rm command",
		logger.String("path", opts.Path),
	)

	if opts.Permanent {
		c.RenderWarning(opts.Stdout, "This action will permanently delete \"%s\" and cannot be undone.", opts.Path)
	}

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if err := fsSvc.Remove(ctx, opts.Path, fs.RemoveOptions{
		Permanent: opts.Permanent,
	}); err != nil {
		return util.NewCommandError(c.Name, "failed to move item", err)
	}

	c.Log.Info("rm completed successfully",
		logger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "Successfully removed \"%s\"\n", opts.Path)

	return nil
}
