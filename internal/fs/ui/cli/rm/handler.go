package rm

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive rm operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive rm Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the rm operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(opts.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	opts.URI = uri

	return opts.Validate()
}

// Execute removes a file or directory at the specified path.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("path", opts.URI.String()),
	)

	log.Info("starting rm operation")

	log.Debug("delegating to filesystem manager for rm")
	if err := c.manager.Remove(ctx, opts.URI); err != nil {
		log.Error("rm failed", logger.Error(err))
		return fmt.Errorf("failed to remove %s: %w", opts.URI, err)
	}

	log.Info("rm completed successfully")
	fmt.Fprintf(opts.Stdout, "Removed %s\n", opts.URI)
	return nil
}

// Finalize performs any necessary cleanup after the rm operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
