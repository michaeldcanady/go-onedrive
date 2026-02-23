package cat

import (
	"context"
	"io"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// CatCmd handles the execution logic for the 'cat' command.
type CatCmd struct {
	container di.Container
	logger    infralogging.Logger
}

// NewCatCmd creates a new CatCmd instance with the provided dependency container.
func NewCatCmd(container di.Container) *CatCmd {
	return &CatCmd{
		container: container,
	}
}

// WithLogger allows injecting a logger into CatCmd for testing.
func (c *CatCmd) WithLogger(logger infralogging.Logger) *CatCmd {
	c.logger = logger
	return c
}

// Run executes the cat lifecycle.
func (c *CatCmd) Run(ctx context.Context, opts Options) error {
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

	c.logger = c.logger.WithContext(ctx).With(logging.String("correlationID", util.CorrelationIDFromContext(ctx)))

	c.logger.Info("starting cat command", infralogging.String("path", opts.Path))

	fsSvc := c.container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	reader, err := fsSvc.ReadFile(ctx, opts.Path, fs.ReadOptions{})
	if err != nil {
		c.logger.Error("failed to read file", infralogging.Error(err))
		return util.NewCommandErrorWithNameWithMessage(commandName, "unable to read path contents")
	}
	defer reader.Close()

	_, err = io.Copy(opts.Stdout, reader)
	if err != nil {
		c.logger.Error("failed to write file contents", infralogging.Error(err))
		return util.NewCommandErrorWithNameWithMessage(commandName, "failed to write file contents")
	}

	c.logger.Info("cat command completed successfully",
		infralogging.Duration("duration", time.Since(start)),
	)

	return nil
}
