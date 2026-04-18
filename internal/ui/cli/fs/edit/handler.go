package edit

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
)

// Command executes the edit operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	editor     editor.Service
	log        logger.Logger
}

// NewCommand initializes a new instance of the edit Command.
func NewCommand(m fs.Service, f *fs.URIFactory, e editor.Service, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		editor:     e,
		log:        l,
	}
}

// Validate prepares and validates the options for the edit operation.
func (c *Command) Validate(ctx *CommandContext) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	ctx.Options.URI = uri

	return ctx.Options.Validate()
}

// Execute opens a file in the user's preferred editor for modification.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("path", ctx.Options.URI.String()),
	)

	log.Info("starting edit operation")

	svc := c.editor
	if ctx.Options.Editor != "" {
		svc = svc.WithIO(os.Stdin, ctx.Options.Stdout, os.Stderr)
	}

	reader, err := c.manager.ReadFile(ctx.Ctx, ctx.Options.URI, pkgfs.ReadOptions{})
	if err != nil {
		log.Error("failed to read file", logger.Error(err))
		return fmt.Errorf("failed to read file: %w", err)
	}

	prefix := "odc-edit-"
	suffix := filepath.Ext(ctx.Options.URI.Path)

	newData, _, err := svc.LaunchTempFile(prefix, suffix, reader)
	if err != nil {
		log.Error("editor session failed", logger.Error(err))
		return fmt.Errorf("editor session failed: %w", err)
	}

	if newData != nil {
		log.Info("file modified, writing changes")
		if _, err := c.manager.WriteFile(ctx.Ctx, ctx.Options.URI, bytes.NewReader(newData), pkgfs.WriteOptions{Overwrite: true}); err != nil {
			log.Error("failed to write changes", logger.Error(err))
			return fmt.Errorf("failed to write changes: %w", err)
		}
	} else {
		log.Info("no changes detected")
	}

	return nil
}

// Finalize performs any necessary cleanup after the edit operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
