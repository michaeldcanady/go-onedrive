package cat

import (
	"context"
	"io"
	"time"

	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// CatCmd handles the execution of the 'cat' command, which reads and displays
// the contents of a file from OneDrive.
type CatCmd struct {
	util.BaseCommand
}

// NewCatCmd creates a new CatCmd instance with the provided dependency container.
func NewCatCmd(container didomain.Container) *CatCmd {
	return &CatCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the cat command. It initializes the command, retrieves the
// file contents from the filesystem service, and writes them to the specified output.
// It uses the domainfs.Reader interface to decouple from the full filesystem service.
func (c *CatCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting cat command",
		domainlogger.String("path", opts.Path),
	)

	// Decouple by using the Reader interface instead of the full Service.
	var fsSvc domainfs.Reader = c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	reader, err := fsSvc.ReadFile(ctx, opts.Path, domainfs.ReadOptions{})
	if err != nil {
		c.Log.Error("failed to read file",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to read file", err)
	}
	defer reader.Close()

	if _, err := io.Copy(opts.Stdout, reader); err != nil {
		c.Log.Error("failed to write output",
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to write output", err)
	}

	c.Log.Info("cat completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	return nil
}
