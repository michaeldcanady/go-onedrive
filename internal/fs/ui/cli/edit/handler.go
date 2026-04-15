package edit

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/editor"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
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
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(opts.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	opts.URI = uri

	return opts.Validate()
}

// Execute opens a file in the user's preferred editor for modification.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("path", opts.URI.String()),
	)

	log.Info("starting edit operation")

	svc := c.editor
	if opts.Editor != "" {
		svc = svc.WithIO(os.Stdin, opts.Stdout, os.Stderr)
	}

	reader, err := c.manager.ReadFile(ctx, opts.URI, fs.ReadOptions{})
	if err != nil {
		log.Error("failed to read file", logger.Error(err))
		return fmt.Errorf("failed to read file: %w", err)
	}

	prefix := "odc-edit-"
	suffix := filepath.Ext(opts.URI.Path)

	newData, _, err := svc.LaunchTempFile(prefix, suffix, reader)
	if err != nil {
		log.Error("editor session failed", logger.Error(err))
		return fmt.Errorf("editor session failed: %w", err)
	}

	if newData != nil {
		log.Info("file modified, writing changes")
		if _, err := c.manager.WriteFile(ctx, opts.URI, bytes.NewReader(newData), fs.WriteOptions{Overwrite: true}); err != nil {
			log.Error("failed to write changes", logger.Error(err))
			return fmt.Errorf("failed to write changes: %w", err)
		}
	} else {
		log.Info("no changes detected")
	}

	return nil
}

// Finalize performs any necessary cleanup after the edit operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
