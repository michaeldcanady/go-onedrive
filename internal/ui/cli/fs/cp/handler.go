package cp

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Command executes the drive cp operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive cp Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the copy operation.
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

	return nil
}

// Execute copies an item from the source to the destination.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("starting copy operation",
		logger.String("source", ctx.Options.SourceURI.String()),
		logger.String("destination", ctx.Options.DestinationURI.String()),
		logger.Bool("recursive", ctx.Options.Recursive),
	)

	cpOpts := pkgfs.CopyOptions{
		Recursive: ctx.Options.Recursive,
		Overwrite: true, // Default to overwrite for now, can be a flag later.
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := c.manager.Copy(ctx.Ctx, ctx.Options.SourceURI, ctx.Options.DestinationURI, cpOpts); err != nil {
		log.Error("copy failed", logger.Error(err))
		return fmt.Errorf("failed to copy %s to %s: %w", ctx.Options.SourceURI, ctx.Options.DestinationURI, err)
	}

	log.Info("copy completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the copy operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Fprintf(ctx.Options.Stdout, "Copied %s to %s\n", ctx.Options.SourceURI, ctx.Options.DestinationURI)
	return nil
}
