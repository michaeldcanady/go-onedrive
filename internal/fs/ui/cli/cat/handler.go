package cat

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive cat operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive cat Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the cat operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(opts.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	opts.URI = uri

	return opts.Validate()
}

// Execute displays the content of a file.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("path", opts.URI.String()),
	)

	log.Info("starting cat operation")

	log.Debug("delegating to filesystem manager for cat")
	reader, err := c.manager.ReadFile(ctx, opts.URI, fs.ReadOptions{})
	if err != nil {
		log.Error("cat failed", logger.Error(err))
		return fmt.Errorf("failed to open file %s: %w", opts.URI, err)
	}
	defer reader.Close()

	if _, err := io.Copy(opts.Stdout, reader); err != nil {
		log.Error("copying content failed", logger.Error(err))
		return fmt.Errorf("failed to read file %s: %w", opts.URI, err)
	}

	log.Info("cat completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the cat operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
