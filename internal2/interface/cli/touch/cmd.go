package touch

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Command struct {
	container di.Container
	logger    infralogging.Logger
}

// NewCmd creates a new Command instance with the provided dependency container.
func NewCmd(container di.Container) *Command {
	return &Command{
		container: container,
	}
}

// WithLogger allows injecting a logger into Command for testing.
func (c *Command) WithLogger(logger infralogging.Logger) *Command {
	c.logger = logger
	return c
}

func (c *Command) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if ctx == nil {
		ctx = context.Background()
	}

	if c.logger == nil {
		logger, err := util.EnsureLogger(c.container, loggerID)
		if err != nil {
			return util.NewCommandErrorWithNameWithError(commandName, err)
		}
		c.logger = logger
	}

	c.logger.Info("starting touch command")

	c.logger.Debug("resolving filesystem service")
	fsSvc := c.container.FS()
	if fsSvc == nil {
		c.logger.Error("filesystem service is nil")
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	c.logger.Info("initiating file creation",
		infralogging.String("path", opts.Path),
	)

	if _, err := fsSvc.Touch(ctx, opts.Path, fs.TouchOptions{}); err != nil {
		c.logger.Error("failed to touch file",
			infralogging.String("path", opts.Path),
			infralogging.Error(err),
		)
		return util.NewCommandError(commandName, "failed to touch file", err)
	}

	c.logger.Info("touch completed successfully",
		infralogging.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully touched \"%s\"\n", opts.Path)

	return nil
}
