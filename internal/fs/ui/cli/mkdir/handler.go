package mkdir

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive mkdir operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive mkdir Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the mkdir operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(opts.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	opts.URI = uri

	return opts.Validate()
}

// Execute creates a new directory at the specified path.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Info("starting mkdir operation", logger.String("path", opts.URI.String()))

	log.Debug("delegating to filesystem manager for mkdir")
	if err := c.manager.Mkdir(ctx, opts.URI); err != nil {
		log.Error("mkdir failed", logger.Error(err))
		return fmt.Errorf("failed to create directory %s: %w", opts.URI, err)
	}

	log.Info("mkdir completed successfully")
	fmt.Fprintf(opts.Stdout, "Created directory %s\n", opts.URI)
	return nil
}

// Finalize performs any necessary cleanup after the mkdir operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
