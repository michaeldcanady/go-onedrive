package mkdir

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

// WithLogger allows injecting a logger into UploadCmd for testing.
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

	c.logger.Info("starting mkdir command",
		infralogging.Bool("parent", opts.Parent),
	)

	c.logger.Debug("resolving filesystem service")
	fsSvc := c.container.FS()
	if fsSvc == nil {
		c.logger.Error("filesystem service is nil")
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	c.logger.Info("initiating directory creation",
		infralogging.String("path", opts.Path),
	)

	if err := fsSvc.Mkdir(ctx, opts.Path, fs.MKDirOptions{Parents: opts.Parent}); err != nil {
		c.logger.Error("failed to create directory",
			infralogging.String("path", opts.Path),
			infralogging.Error(err),
		)
		return util.NewCommandError(commandName, "failed to create directory", err)
	}

	c.logger.Info("mkdir completed successfully",
		infralogging.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully created \"%s\"\n", opts.Path)

	return nil
}
