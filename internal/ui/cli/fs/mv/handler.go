package mv

import (
	"fmt"

	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
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

// Execute moves an item from the source to the destination.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("source", ctx.Options.SourceURI.String()),
		logger.String("destination", ctx.Options.DestinationURI.String()),
	)

	log.Info("starting move operation")

	log.Debug("delegating to filesystem manager for move")
	if err := c.manager.Move(ctx.Ctx, ctx.Options.SourceURI, ctx.Options.DestinationURI); err != nil {
		log.Error("move failed", logger.Error(err))
		return fmt.Errorf("failed to move %s to %s: %w", ctx.Options.SourceURI, ctx.Options.DestinationURI, err)
	}

	log.Info("move completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the move operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Fprintf(ctx.Options.Stdout, "Moved %s to %s\n", ctx.Options.SourceURI, ctx.Options.DestinationURI)
	return nil
}
