// Package upload provides the command-line interface for uploading local files to OneDrive.
package upload

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// UploadCmd handles the execution logic for the 'upload' command.
type UploadCmd struct {
	util.BaseCommand
}

// NewUploadCmd creates a new UploadCmd instance with the provided dependency container.
func NewUploadCmd(container didomain.Container) *UploadCmd {
	return &UploadCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the upload command, transferring a local file or directory to a OneDrive path.
// It uses the domainfs.Writer interface to decouple from the full filesystem service.
func (c *UploadCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting upload command",
		domainlogger.String("src", opts.Source),
		domainlogger.String("dst", opts.Destination),
	)

	// Decouple by using the Writer interface.
	var fsSvc domainfs.Writer = c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	if _, err := fsSvc.Upload(ctx, opts.Source, opts.Destination, domainfs.UploadOptions{}); err != nil {
		c.Log.Error("failed to upload item",
			domainlogger.String("src", opts.Source),
			domainlogger.String("dst", opts.Destination),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to upload item", err)
	}

	c.Log.Info("upload completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "uploaded \"%s\" to \"%s\"", opts.Source, opts.Destination)

	return nil
}
