package download

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type DownloadCmd struct {
	util.BaseCommand
}

func NewDownloadCmd(container didomain.Container) *DownloadCmd {
	return &DownloadCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *DownloadCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting download command",
		domainlogger.String("src", opts.Source),
		domainlogger.String("dst", opts.Destination),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	reader, err := fsSvc.ReadFile(ctx, opts.Source, domainfs.ReadOptions{})
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to read source file", err)
	}
	defer reader.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(opts.Destination), 0o755); err != nil {
		return util.NewCommandError(c.Name, "failed to create destination directory", err)
	}

	destFile, err := os.Create(opts.Destination)
	if err != nil {
		return util.NewCommandError(c.Name, "failed to create destination file", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, reader); err != nil {
		return util.NewCommandError(c.Name, "failed to write destination file", err)
	}

	c.Log.Info("download completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "Successfully downloaded \"%s\" to \"%s\"\n", opts.Source, opts.Destination)

	return nil
}
