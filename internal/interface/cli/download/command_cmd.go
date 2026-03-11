// Package download provides the command-line interface for downloading OneDrive files to the local filesystem.
package download

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// DownloadCmd handles the execution logic for the 'download' command.
type DownloadCmd struct {
	util.BaseCommand
}

// NewDownloadCmd creates a new DownloadCmd instance with the provided dependency container.
func NewDownloadCmd(container didomain.Container) *DownloadCmd {
	return &DownloadCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the download command, transferring a OneDrive file to a local destination.
// It uses the domainfs.Reader interface to decouple from the full filesystem service.
func (c *DownloadCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting download command",
		domainlogger.String("src", opts.Source),
		domainlogger.String("dst", opts.Destination),
	)

	// Decouple by using the Reader interface.
	var fsSvc domainfs.Reader = c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	reader, err := fsSvc.ReadFile(ctx, opts.Source, domainfs.ReadOptions{})
	if err != nil {
		c.Log.Error("failed to read source file",
			domainlogger.String("path", opts.Source),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to read source file", err)
	}
	defer reader.Close()

	// Ensure destination directory exists.
	if err := os.MkdirAll(filepath.Dir(opts.Destination), 0o755); err != nil {
		return util.NewCommandError(c.Name, "failed to create destination directory", err)
	}

	destFile, err := os.Create(opts.Destination)
	if err != nil {
		return util.NewCommandError(c.Name, "failed to create destination file", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, reader); err != nil {
		c.Log.Error("failed to write destination file",
			domainlogger.String("path", opts.Destination),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to write destination file", err)
	}

	c.Log.Info("download completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "downloaded \"%s\" to \"%s\"", opts.Source, opts.Destination)

	return nil
}
