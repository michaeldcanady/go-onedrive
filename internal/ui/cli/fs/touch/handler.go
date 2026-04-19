package touch

import (
	"fmt"

	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
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
func (c *Command) Validate(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		log.Error("invalid path uri", logger.String("path", uri.String()), logger.Error(err))
		return fmt.Errorf("invalid path: %w", err)
	}
	ctx.Options.URI = uri

	return nil
}

// Execute creates a new empty file or updates the timestamp of an existing file.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("starting touch operation", logger.String("path", ctx.Options.URI.String()))

	log.Debug("delegating to filesystem manager for touch")
	if _, err := c.manager.Touch(ctx.Ctx, ctx.Options.URI); err != nil {
		log.Error("touch failed", logger.Error(err))
		return fmt.Errorf("failed to touch %s: %w", ctx.Options.URI, err)
	}

	log.Info("touch completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the touch operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Fprintf(ctx.Options.Stdout, "Touched %s\n", ctx.Options.URI)
	return nil
}
