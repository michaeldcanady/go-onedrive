package upload

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive upload operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive upload Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the upload operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	log := c.log.WithContext(ctx)

	// Resolve URIs using the factory
	sourceURI, err := c.uriFactory.FromString(opts.Source)
	if err != nil {
		log.Error("invalid source uri", logger.String("uri", opts.Source), logger.Error(err))
		return fmt.Errorf("invalid source path: %w", err)
	}
	opts.SourceURI = sourceURI

	destinationURI, err := c.uriFactory.FromString(opts.Destination)
	if err != nil {
		log.Error("invalid destination uri", logger.String("uri", opts.Destination), logger.Error(err))
		return fmt.Errorf("invalid destination path: %w", err)
	}
	opts.DestinationURI = destinationURI

	return opts.Validate()
}

// Execute uploads files and directories from the local filesystem to OneDrive.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Info("starting upload operation", logger.String("source", opts.SourceURI.String()),
		logger.String("destination", opts.DestinationURI.String()),
		logger.Bool("recursive", opts.Recursive))

	// Upload is essentially a Copy from local to remote.
	cpOpts := fs.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for upload (copy)")
	if err := c.manager.Copy(ctx, opts.SourceURI, opts.DestinationURI, cpOpts); err != nil {
		log.Error("upload failed", logger.Error(err))
		return fmt.Errorf("failed to upload %s to %s: %w", opts.SourceURI, opts.DestinationURI, err)
	}

	log.Info("upload completed successfully")
	fmt.Fprintf(opts.Stdout, "Uploaded %s to %s\n", opts.SourceURI, opts.DestinationURI)
	return nil
}

// Finalize performs any necessary cleanup after the upload operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
