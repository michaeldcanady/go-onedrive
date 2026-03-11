// Package cp provides the command-line interface for copying items in OneDrive.
package cp

import (
	"context"
	"os"
	"time"

	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
)

// CpCmd handles the execution logic for the 'cp' command.
type CpCmd struct {
	util.BaseCommand
}

// NewCpCmd creates a new CpCmd instance with the provided dependency container.
func NewCpCmd(container didomain.Container) *CpCmd {
	return &CpCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the cp command, copying an item from a source path to a destination path.
// It uses the domainfs.Reader and domainfs.Manager interfaces to decouple from the full filesystem service.
func (c *CpCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting cp command",
		domainlogger.String("src", opts.Source),
		domainlogger.String("dst", opts.Dest),
		domainlogger.Bool("overwrite", opts.Overwrite),
		domainlogger.String("ignoreFile", opts.IgnoreFile),
	)

	// Decouple by using specific interfaces.
	var reader domainfs.Reader = c.Container.FS()
	var manager domainfs.Manager = c.Container.FS()
	if reader == nil || manager == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	matcher, err := c.loadIgnoreMatcher(ctx, opts.IgnoreFile)
	if err != nil {
		c.Log.Warn("failed to load ignore file", domainlogger.String("path", opts.IgnoreFile), domainlogger.Error(err))
		matcher = nil
	}

	item, err := reader.Get(ctx, opts.Source)
	if err != nil {
		c.Log.Error("failed to get source item",
			domainlogger.String("path", opts.Source),
			domainlogger.Error(err),
		)
		return err
	}

	if item.Type == domainfs.ItemTypeFolder && !opts.Recursive {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "recursive is required for folders")
	} else if item.Type == domainfs.ItemTypeFile && opts.Recursive {
		c.RenderWarning(opts.Stdout, "disabling recursive since path is a file")
		opts.Recursive = false
	}

	copyOpts := domainfs.CopyOptions{
		Overwrite: opts.Overwrite,
		Recursive: opts.Recursive,
		Matcher:   matcher,
	}

	if err := manager.Copy(ctx, opts.Source, opts.Dest, copyOpts); err != nil {
		c.Log.Error("failed to copy item",
			domainlogger.String("src", opts.Source),
			domainlogger.String("dst", opts.Dest),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to copy item", err)
	}

	c.Log.Info("cp completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "copied \"%s\" to \"%s\"", opts.Source, opts.Dest)

	return nil
}

// loadIgnoreMatcher attempts to create an ignore matcher from the provided file path.
func (c *CpCmd) loadIgnoreMatcher(ctx context.Context, ignorePath string) (domainfs.IgnoreMatcher, error) {
	if ignorePath == "" {
		return nil, nil
	}

	f, err := os.Open(ignorePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	factory := c.Container.IgnoreMatcherFactory()
	if factory == nil {
		return nil, nil
	}

	return factory.CreateMatcher(ctx, f)
}
