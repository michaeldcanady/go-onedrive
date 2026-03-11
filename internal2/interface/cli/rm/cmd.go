package rm

import (
	"context"
	"fmt"
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

	// TODO: need to open bug with MS Graph SDk https://github.com/microsoftgraph/msgraph-sdk-go/issues/980
	if opts.Permanent {
		if !opts.Quiet {
			c.RenderWarning(opts.Stdout, "This action will permanently delete \"%s\" and cannot be undone.", opts.Path)
		}
		if !opts.Force {
			if opts.Quiet {
				// TODO: add appropriate error message
				return util.NewCommandErrorWithNameWithMessage(c.Name, "")
			}
			confirmed, err := c.PromptConfirm(opts.Stdout, "Are you sure you want to proceed")
			if err != nil {
				return util.NewCommandError(c.Name, "failed to get confirmation", err)
			}
			if !confirmed {
				fmt.Fprintln(opts.Stdout, "Aborted.")
				return nil
			}
		}
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

	if !opts.Quiet {
		c.RenderSuccess(opts.Stdout, "Successfully removed \"%s\"\n", opts.Path)
	}
	return nil
}
