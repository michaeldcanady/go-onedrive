package download

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive download operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive download Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the download operation.
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

// Execute downloads files and directories from OneDrive to the local filesystem.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("source", opts.SourceURI.String()),
		logger.String("destination", opts.DestinationURI.String()),
		logger.Bool("recursive", opts.Recursive),
	)

	log.Info("starting download operation")

	// Download is essentially a Copy from remote to local.
	cpOpts := fs.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for download (copy)")
	if err := c.manager.Copy(ctx, opts.SourceURI, opts.DestinationURI, cpOpts); err != nil {
		log.Error("download failed", logger.Error(err))
		return fmt.Errorf("failed to download %s to %s: %w", opts.SourceURI, opts.DestinationURI, err)
	}

	log.Info("download completed successfully")
	fmt.Fprintf(opts.Stdout, "Downloaded %s to %s\n", opts.SourceURI, opts.DestinationURI)
	return nil
}

// Finalize performs any necessary cleanup after the download operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
