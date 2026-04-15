package touch

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive touch operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive touch Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the touch operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	log := c.log.WithContext(ctx)

	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(opts.Path)
	if err != nil {
		log.Error("invalid path uri", logger.String("path", uri.String()), logger.Error(err))
		return fmt.Errorf("invalid path: %w", err)
	}
	opts.URI = uri

	return opts.Validate()
}

// Execute creates a new empty file or updates the timestamp of an existing file.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Info("starting touch operation", logger.String("path", opts.URI.String()))

	log.Debug("delegating to filesystem manager for touch")
	if _, err := c.manager.Touch(ctx, opts.URI); err != nil {
		log.Error("touch failed", logger.Error(err))
		return fmt.Errorf("failed to touch %s: %w", opts.URI, err)
	}

	log.Info("touch completed successfully")
	fmt.Fprintf(opts.Stdout, "Touched %s\n", opts.URI)
	return nil
}

// Finalize performs any necessary cleanup after the touch operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
