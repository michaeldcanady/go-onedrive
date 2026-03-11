// Package edit provides the command-line interface for editing OneDrive files in a local text editor.
package edit

import (
	"context"
	"errors"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainerrors "github.com/michaeldcanady/go-onedrive/internal/common/errors"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal/editor/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// ConflictHandler defines an interface for resolving upload conflicts during editing.
type ConflictHandler interface {
	// HandleConflict is called when a precondition failure (conflict) occurs during upload.
	// It should return whether to remove the temp file, the final path used, and any error.
	HandleConflict(ctx context.Context, path string, content []byte, tmpPath string) (bool, string, error)
}

// EditCmd handles the execution logic for the 'edit' command.
type EditCmd struct {
	util.BaseCommand
	editorSvc       domaineditor.Service
	conflictHandler ConflictHandler
}

// NewEditCmd creates a new EditCmd instance with the provided dependency container.
func NewEditCmd(container didomain.Container) *EditCmd {
	return &EditCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// WithLogger allows injecting a logger for testing.
func (c *EditCmd) WithLogger(log domainlogger.Logger) *EditCmd {
	c.Log = log
	return c
}

// WithEditor allows injecting an editor service for testing.
func (c *EditCmd) WithEditor(svc domaineditor.Service) *EditCmd {
	c.editorSvc = svc
	return c
}

// WithConflictHandler allows injecting a conflict handler for testing.
func (c *EditCmd) WithConflictHandler(handler ConflictHandler) *EditCmd {
	c.conflictHandler = handler
	return c
}

// Run executes the edit command. It downloads the file to a temporary location,
// opens it in the user's editor, and uploads the changes back to OneDrive.
// It uses domainfs.Reader and domainfs.Writer interfaces to decouple from the full filesystem service.
func (c *EditCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting edit command",
		domainlogger.String("path", opts.Path),
	)

	// Decouple by using specific interfaces.
	var reader domainfs.Reader = c.Container.FS()
	var writer domainfs.Writer = c.Container.FS()
	if reader == nil || writer == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	editorSvc := c.editorSvc
	if editorSvc == nil {
		editorSvc = c.Container.Editor()
	}
	if editorSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "editor service is nil")
	}

	c.Log.Debug("fetching metadata", domainlogger.String("path", opts.Path))
	item, err := reader.Get(ctx, opts.Path)
	if err != nil {
		c.Log.Error("failed to get item metadata",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to get item metadata", err)
	}

	c.Log.Debug("reading file for editing", domainlogger.String("path", opts.Path))
	r, err := reader.ReadFile(ctx, opts.Path, domainfs.ReadOptions{})
	if err != nil {
		c.Log.Error("failed to read file for editing",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to read file for editing", err)
	}
	defer r.Close()

	c.Log.Debug("launching editor")
	newData, tmpPath, err := editorSvc.LaunchTempFile("odc-edit", ".txt", r)
	if err != nil {
		c.Log.Error("failed to edit file",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to edit file", err)
	}

	if newData == nil {
		c.Log.Info("no changes detected, skipping upload")
		c.RenderInfo(opts.Stdout, "No changes detected.")
		return nil
	}

	c.Log.Info("uploading changes", domainlogger.String("path", opts.Path))
	_, err = writer.WriteFile(ctx, opts.Path, util.NewReader(newData), domainfs.WriteOptions{
		Overwrite: opts.Force,
		IfMatch:   item.ETag,
	})

	if err != nil && (errors.Is(err, domainerrors.ErrPrecondition)) {
		c.Log.Warn("conflict detected during upload", domainlogger.Error(err))
		if c.conflictHandler != nil {
			_, finalPath, hErr := c.conflictHandler.HandleConflict(ctx, opts.Path, newData, tmpPath)
			if hErr != nil {
				return util.NewCommandError(c.Name, "conflict resolution failed", hErr)
			}
			c.RenderSuccess(opts.Stdout, "updated file \"%s\"", finalPath)
			return nil
		}
	}

	if err != nil {
		c.Log.Error("failed to upload edited changes",
			domainlogger.String("path", opts.Path),
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to upload edited changes", err)
	}

	c.Log.Info("edit completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)

	c.RenderSuccess(opts.Stdout, "updated file \"%s\"", opts.Path)

	return nil
}
