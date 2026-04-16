package download

import (
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
func (c *Command) Validate(ctx *CommandContext) error {
	// Resolve URIs using the factory
	sourceURI, err := c.uriFactory.FromString(ctx.Options.Source)
	if err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}
	ctx.Options.SourceURI = sourceURI

	destinationURI, err := c.uriFactory.FromString(ctx.Options.Destination)
	if err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}
	ctx.Options.DestinationURI = destinationURI

	return ctx.Options.Validate()
}

// Execute downloads files and directories from OneDrive to the local filesystem.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("source", ctx.Options.SourceURI.String()),
		logger.String("destination", ctx.Options.DestinationURI.String()),
		logger.Bool("recursive", ctx.Options.Recursive),
	)

	log.Info("starting download operation")

	// Download is essentially a Copy from remote to local.
	cpOpts := fs.CopyOptions{
		Recursive: ctx.Options.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for download (copy)")
	if err := c.manager.Copy(ctx.Ctx, ctx.Options.SourceURI, ctx.Options.DestinationURI, cpOpts); err != nil {
		log.Error("download failed", logger.Error(err))
		return fmt.Errorf("failed to download %s to %s: %w", ctx.Options.SourceURI, ctx.Options.DestinationURI, err)
	}

	log.Info("download completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the download operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Fprintf(ctx.Options.Stdout, "Downloaded %s to %s\n", ctx.Options.SourceURI, ctx.Options.DestinationURI)
	return nil
}
