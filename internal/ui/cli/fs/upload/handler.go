package upload

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
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
func (c *Command) Validate(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	// Resolve URIs using the factory
	sourceURI, err := c.uriFactory.FromString(ctx.Options.Source)
	if err != nil {
		log.Error("invalid source uri", logger.String("uri", ctx.Options.Source), logger.Error(err))
		return fmt.Errorf("invalid source path: %w", err)
	}
	ctx.Options.SourceURI = sourceURI

	destinationURI, err := c.uriFactory.FromString(ctx.Options.Destination)
	if err != nil {
		log.Error("invalid destination uri", logger.String("uri", ctx.Options.Destination), logger.Error(err))
		return fmt.Errorf("invalid destination path: %w", err)
	}
	ctx.Options.DestinationURI = destinationURI

	return nil
}

// Execute uploads files and directories from the local filesystem to OneDrive.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("starting upload operation", logger.String("source", ctx.Options.SourceURI.String()),
		logger.String("destination", ctx.Options.DestinationURI.String()),
		logger.Bool("recursive", ctx.Options.Recursive))

	// Upload is essentially a Copy from local to remote.
	cpOpts := pkgfs.CopyOptions{
		Recursive: ctx.Options.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for upload (copy)")
	if err := c.manager.Copy(ctx.Ctx, ctx.Options.SourceURI, ctx.Options.DestinationURI, cpOpts); err != nil {
		log.Error("upload failed", logger.Error(err))
		return fmt.Errorf("failed to upload %s to %s: %w", ctx.Options.SourceURI, ctx.Options.DestinationURI, err)
	}

	log.Info("upload completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the upload operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	_, _ = fmt.Fprintf(ctx.Options.Stdout, "Uploaded %s to %s\n", ctx.Options.SourceURI, ctx.Options.DestinationURI)
	return nil
}
