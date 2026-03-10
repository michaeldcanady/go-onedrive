package upload

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type UploadCmd struct {
	util.BaseCommand
}

func NewUploadCmd(container di.Container) *UploadCmd {
	return &UploadCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *UploadCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting upload command",
		logger.String("src", opts.Source),
		logger.String("dst", opts.Destination),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if _, err := fsSvc.Upload(ctx, opts.Source, opts.Destination, fs.UploadOptions{}); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to upload item", err)
	}

	c.Log.Info("upload completed successfully",
		logger.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully uploaded \"%s\" to \"%s\"\n", opts.Source, opts.Destination)

	return nil
}
