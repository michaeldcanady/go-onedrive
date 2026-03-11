package upload

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type UploadCmd struct {
	util.BaseCommand
}

func NewUploadCmd(container didomain.Container) *UploadCmd {
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
		domainlogger.String("src", opts.Source),
		domainlogger.String("dst", opts.Destination),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if _, err := fsSvc.Upload(ctx, opts.Source, opts.Destination, domainfs.UploadOptions{}); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to upload item", err)
	}

	c.Log.Info("upload completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully uploaded \"%s\" to \"%s\"\n", opts.Source, opts.Destination)

	return nil
}
