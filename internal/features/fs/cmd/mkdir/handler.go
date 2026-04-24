package mkdir

import (
	"context"
	"fmt"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Logger defines the interface required for logging within the mkdir command.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
	Info(msg string, fields ...logger.Field)
	With(fields ...logger.Field) logger.Logger
	WithContext(ctx context.Context) logger.Logger
}

// URIFactory defines the interface for creating URIs.
type URIFactory interface {
	FromString(s string) (*fs.URI, error)
}

type ItemCreator interface {
	Mkdir(ctx context.Context, uri *fs.URI) error
}

// Command executes the drive mkdir operation.
type Command struct {
	manager    ItemCreator
	uriFactory URIFactory
	log        Logger
}

// NewCommand initializes a new instance of the drive mkdir Command.
func NewCommand(m ItemCreator, f URIFactory, l Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		log:        l,
	}
}

// Validate prepares and validates the options for the mkdir operation.
func (c *Command) Validate(ctx *CommandContext) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	ctx.Options.URI = uri

	return ctx.Options.Validate()
}

// Execute creates a new directory at the specified path.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("starting mkdir operation", logger.String("path", ctx.Options.URI.String()))

	log.Debug("delegating to filesystem manager for mkdir")
	if err := c.manager.Mkdir(ctx.Ctx, ctx.Options.URI); err != nil {
		log.Error("mkdir failed", logger.Error(err))
		return fmt.Errorf("failed to create directory %s: %w", ctx.Options.URI, err)
	}

	log.Info("mkdir completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the mkdir operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Fprintf(ctx.Options.Stdout, "Created directory %s\n", ctx.Options.URI)
	return nil
}
