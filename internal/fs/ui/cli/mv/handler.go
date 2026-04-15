package mv

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive mv operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive mv Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the move operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	// Resolve URIs using the factory
	sourceURI, err := c.uriFactory.FromString(opts.Source)
	if err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}
	opts.SourceURI = sourceURI

	destinationURI, err := c.uriFactory.FromString(opts.Destination)
	if err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}
	opts.DestinationURI = destinationURI

	return opts.Validate()
}

// Execute moves an item from the source to the destination.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("source", opts.SourceURI.String()),
		logger.String("destination", opts.DestinationURI.String()),
	)

	log.Info("starting move operation")

	log.Debug("delegating to filesystem manager for move")
	if err := c.manager.Move(ctx, opts.SourceURI, opts.DestinationURI); err != nil {
		log.Error("move failed", logger.Error(err))
		return fmt.Errorf("failed to move %s to %s: %w", opts.SourceURI, opts.DestinationURI, err)
	}

	log.Info("move completed successfully")
	fmt.Fprintf(opts.Stdout, "Moved %s to %s\n", opts.SourceURI, opts.DestinationURI)
	return nil
}

// Finalize performs any necessary cleanup after the move operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
