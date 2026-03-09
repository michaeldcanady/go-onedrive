package cp

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// CpCmd handles the execution logic for the 'cp' command.
type CpCmd struct {
	container di.Container
	logger    infralogging.Logger
}

// NewCpCmd creates a new CpCmd instance with the provided dependency container.
func NewCpCmd(container di.Container) *CpCmd {
	return &CpCmd{
		container: container,
	}
}

// WithLogger allows injecting a logger into CpCmd for testing.
func (c *CpCmd) WithLogger(logger infralogging.Logger) *CpCmd {
	c.logger = logger
	return c
}

// Run executes the cp lifecycle.
func (c *CpCmd) Run(ctx context.Context, opts Options) error {
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

	c.logger.Info("starting cp command",
		infralogging.String("source", opts.Source),
		infralogging.String("dest", opts.Dest),
		infralogging.Bool("overwrite", opts.Overwrite),
	)

	fsSvc := c.container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service not found")
	}

	copyOpts := domainfs.CopyOptions{
		Overwrite: opts.Overwrite,
	}

	if err := fsSvc.Copy(ctx, opts.Source, opts.Dest, copyOpts); err != nil {
		c.logger.Error("failed to copy", infralogging.Error(err))
		return util.NewCommandErrorWithNameWithError(commandName, err)
	}

	c.logger.Info("cp completed successfully",
		infralogging.Duration("duration", time.Since(start)),
	)

	return nil
}
