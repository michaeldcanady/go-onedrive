package cat

import (
	"context"
	"fmt"
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Logger defines the interface required for logging within the cat command.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
	WithContext(ctx context.Context) logger.Logger
}

type Manager interface {
	ReadFile(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ReadOptions) (io.ReadCloser, error)
}

// URIFactory defines the interface for creating URIs.
type URIFactory interface {
	FromString(s string) (*fs.URI, error)
}

// Command executes the drive cat operation.
type Command struct {
	manager    Manager
	uriFactory URIFactory
	log        Logger
}

// NewCommand initializes a new instance of the drive cat Command.
func NewCommand(m fs.Service, f URIFactory, l Logger) *Command {
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
	ctx.URI = uri

	return ctx.Options.Validate()
}

// Execute displays the content of a file.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("path", ctx.URI.String()),
	)

	log.Info("starting cat operation")

	log.Debug("delegating to filesystem manager for cat")
	reader, err := c.manager.ReadFile(ctx.Ctx, ctx.URI, pkgfs.ReadOptions{})
	if err != nil {
		log.Error("cat failed", logger.Error(err))
		return fmt.Errorf("failed to open file %s: %w", ctx.URI, err)
	}
	defer reader.Close()

	if _, err := io.Copy(ctx.Options.Stdout, reader); err != nil {
		log.Error("copying content failed", logger.Error(err))
		return fmt.Errorf("failed to read file %s: %w", ctx.URI, err)
	}

	log.Info("cat completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the cat operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
