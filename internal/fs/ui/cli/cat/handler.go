package cat

import (
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive cat operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	log        logger.Logger
}

// NewCommand initializes a new instance of the drive cat Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the cat operation.
func (c *Command) Validate(ctx *CommandContext) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	ctx.Options.URI = uri

	return ctx.Options.Validate()
}

// Execute displays the content of a file.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("path", ctx.Options.URI.String()),
	)

	log.Info("starting cat operation")

	log.Debug("delegating to filesystem manager for cat")
	reader, err := c.manager.ReadFile(ctx.Ctx, ctx.Options.URI, fs.ReadOptions{})
	if err != nil {
		log.Error("cat failed", logger.Error(err))
		return fmt.Errorf("failed to open file %s: %w", ctx.Options.URI, err)
	}
	defer reader.Close()

	if _, err := io.Copy(ctx.Options.Stdout, reader); err != nil {
		log.Error("copying content failed", logger.Error(err))
		return fmt.Errorf("failed to read file %s: %w", ctx.Options.URI, err)
	}

	log.Info("cat completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the cat operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
