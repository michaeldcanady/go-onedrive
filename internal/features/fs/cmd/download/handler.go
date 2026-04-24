package download

import (
	"context"
	"fmt"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Logger defines the interface required for logging within the download command.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
	Info(msg string, fields ...logger.Field)
	With(fields ...logger.Field) logger.Logger
	WithContext(ctx context.Context) logger.Logger
}

type Manager interface {
	Copy(ctx context.Context, src *pkgfs.URI, dst *pkgfs.URI, opts pkgfs.CopyOptions) error
}

// URIFactory defines the interface for creating URIs.
type URIFactory interface {
	FromString(s string) (*fs.URI, error)
}

// Command executes the drive download operation.
type Command struct {
	manager    Manager
	uriFactory URIFactory
	log        Logger
}

// NewCommand initializes a new instance of the drive download Command.
func NewCommand(m fs.Service, f *fs.URIFactory, l Logger) *Command {
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
	ctx.SourceURI = sourceURI

	destinationURI, err := c.uriFactory.FromString(ctx.Options.Destination)
	if err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}
	ctx.DestinationURI = destinationURI

	return ctx.Options.Validate()
}

// Execute downloads files and directories from OneDrive to the local filesystem.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("source", ctx.SourceURI.String()),
		logger.String("destination", ctx.DestinationURI.String()),
		logger.Bool("recursive", ctx.Options.Recursive),
	)

	log.Info("starting download operation")

	// Download is essentially a Copy from remote to local.
	cpOpts := pkgfs.CopyOptions{
		Recursive: ctx.Options.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for download (copy)")
	if err := c.manager.Copy(ctx.Ctx, ctx.SourceURI, ctx.DestinationURI, cpOpts); err != nil {
		log.Error("download failed", logger.Error(err))
		return fmt.Errorf("failed to download %s to %s: %w", ctx.SourceURI, ctx.DestinationURI, err)
	}

	log.Info("download completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the download operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	fmt.Fprintf(ctx.Options.Stdout, "Downloaded %s to %s\n", ctx.SourceURI, ctx.DestinationURI)
	return nil
}
