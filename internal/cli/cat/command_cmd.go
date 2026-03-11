package cat

import (
	"context"
	"io"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type CatCmd struct {
	util.BaseCommand
}

func NewCatCmd(container didomain.Container) *CatCmd {
	return &CatCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CatCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting cat command",
		domainlogger.String("path", opts.Path),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	reader, err := fsSvc.ReadFile(ctx, opts.Path, domainfs.ReadOptions{})
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to read file", err)
	}
	defer reader.Close()

	if _, err := io.Copy(opts.Stdout, reader); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to write output", err)
	}

	c.Log.Info("cat completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	return nil
}
