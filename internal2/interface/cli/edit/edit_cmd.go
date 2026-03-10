package edit

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainerrors "github.com/michaeldcanady/go-onedrive/internal2/domain/common/errors"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type ConflictHandler interface {
	HandleConflict(ctx context.Context, path string, content []byte, tmpPath string) (bool, string, error)
}

type EditCmd struct {
	util.BaseCommand
	editorSvc       editor.Service
	conflictHandler ConflictHandler
}

func NewEditCmd(container di.Container) *EditCmd {
	return &EditCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// WithLogger allows injecting a logger for testing.
func (c *EditCmd) WithLogger(log logger.Logger) *EditCmd {
	c.Log = log
	return c
}

// WithEditor allows injecting an editor service for testing.
func (c *EditCmd) WithEditor(svc editor.Service) *EditCmd {
	c.editorSvc = svc
	return c
}

// WithConflictHandler allows injecting a conflict handler for testing.
func (c *EditCmd) WithConflictHandler(handler ConflictHandler) *EditCmd {
	c.conflictHandler = handler
	return c
}

func (c *EditCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting edit command",
		logger.String("path", opts.Path),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	editorSvc := c.editorSvc
	if editorSvc == nil {
		editorSvc = c.Container.Editor()
	}
	if editorSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "editor service is nil")
	}

	c.Log.Debug("fetching metadata", logger.String("path", opts.Path))
	item, err := fsSvc.Get(ctx, opts.Path)
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to get item metadata", err)
	}

	c.Log.Debug("reading file for editing", logger.String("path", opts.Path))
	reader, err := fsSvc.ReadFile(ctx, opts.Path, fs.ReadOptions{})
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to read file for editing", err)
	}
	defer reader.Close()

	c.Log.Debug("launching editor")
	newData, tmpPath, err := editorSvc.LaunchTempFile("odc-edit", ".txt", reader)
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to edit file", err)
	}

	if newData == nil {
		c.Log.Info("no changes detected, skipping upload")
		fmt.Fprintln(opts.Stdout, "No changes detected.")
		return nil
	}

	c.Log.Info("uploading changes", logger.String("path", opts.Path))
	_, err = fsSvc.WriteFile(ctx, opts.Path, util.NewReader(newData), fs.WriteOptions{
		Overwrite: true,
		IfMatch:   item.ETag,
	})

	if err != nil && (errors.Is(err, domainerrors.ErrPrecondition)) {
		c.Log.Warn("conflict detected during upload", logger.Error(err))
		if c.conflictHandler != nil {
			removeTemp, finalPath, hErr := c.conflictHandler.HandleConflict(ctx, opts.Path, newData, tmpPath)
			if hErr != nil {
				return util.NewCommandError(c.Name, "conflict resolution failed", hErr)
			}
			if removeTemp {
				// No-op for now as LaunchTempFile handles it via defer os.Remove
			}
			fmt.Fprintf(opts.Stdout, "File \"%s\" updated successfully.\n", finalPath)
			return nil
		}
	}

	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandError(c.Name, "failed to upload edited changes", err)
	}

	c.Log.Info("edit completed successfully",
		logger.Duration("duration", time.Since(start)),
	)

	fmt.Fprintf(opts.Stdout, "File \"%s\" updated successfully.\n", opts.Path)

	return nil
}
