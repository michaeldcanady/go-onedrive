package rm

import (
	"fmt"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
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
func (c *Command) Validate(ctx *CommandContext) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	ctx.Options.URI = uri

	return nil
}

// Execute removes a file or directory at the specified path.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("starting rm operation", logger.String("path", ctx.Options.URI.String()))

	log.Debug("delegating to filesystem manager for rm")
	if err := c.manager.Remove(ctx.Ctx, ctx.Options.URI); err != nil {
		log.Error("rm failed", logger.Error(err))
		return fmt.Errorf("failed to remove %s: %w", ctx.Options.URI, err)
	}

	log.Info("rm completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the rm operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	_, _ = fmt.Fprintf(ctx.Options.Stdout, "Removed %s\n", ctx.Options.URI)
	return nil
}
