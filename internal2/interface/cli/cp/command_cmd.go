package cp

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type CpCmd struct {
	util.BaseCommand
}

func NewCpCmd(container di.Container) *CpCmd {
	return &CpCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CpCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting cp command",
		logger.String("src", opts.Source),
		logger.String("dst", opts.Dest),
		logger.Bool("overwrite", opts.Overwrite),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if err := fsSvc.Copy(ctx, opts.Source, opts.Dest, fs.CopyOptions{Overwrite: opts.Overwrite}); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to copy item", err)
	}

	c.Log.Info("cp completed successfully",
		logger.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully copied \"%s\" to \"%s\"\n", opts.Source, opts.Dest)

	return nil
}
