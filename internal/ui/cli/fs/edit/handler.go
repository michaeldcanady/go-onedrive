package edit

import (
	"fmt"
	"os"

	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Command executes the edit operation.
type Command struct {
	manager    fs.Service
	uriFactory *fs.URIFactory
	editor     editor.Service
	log        logger.Logger
}

// NewCommand initializes a new instance of the edit Command.
func NewCommand(m fs.Service, f *fs.URIFactory, e editor.Service, l logger.Logger) *Command {
	return &Command{
		manager:    m,
		uriFactory: f,
		editor:     e,
		log:        l,
	}
}

// Validate prepares and validates the options for the edit operation.
func (c *Command) Validate(ctx *CommandContext) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	ctx.Options.URI = uri

	return ctx.Options.Validate()
}

// Execute opens a file in the user's preferred editor for modification.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("path", ctx.Options.URI.String()),
	)

	log.Info("starting edit operation")

	svc := c.editor
	var opts []editor.Option
	if ctx.Options.Editor != "" {
		opts = append(opts, editor.WithEditor(ctx.Options.Editor))
	}
	opts = append(opts, editor.WithIO(os.Stdin, ctx.Options.Stdout, os.Stderr))
	svc = svc.WithOptions(opts...)

	reader, err := c.manager.ReadFile(ctx.Ctx, ctx.Options.URI, pkgfs.ReadOptions{})
	if err != nil {
		log.Error("failed to read file", logger.Error(err))
		return fmt.Errorf("failed to read file: %w", err)
	}
	defer reader.Close()

	session, err := svc.CreateSession(ctx.Ctx, ctx.Options.URI, reader)
	if err != nil {
		log.Error("failed to create editor session", logger.Error(err))
		return fmt.Errorf("failed to create editor session: %w", err)
	}
	defer svc.Cleanup(ctx.Ctx, session)

	log.Debug("editor session created",
		logger.String("id", session.ID),
		logger.String("local_uri", session.LocalURI.String()))

	if err := svc.Open(ctx.Ctx, session); err != nil {
		log.Error("editor session failed", logger.Error(err))
		return fmt.Errorf("editor session failed: %w", err)
	}

	modified, err := svc.Modified(session)
	if err != nil {
		log.Error("failed to check for modifications", logger.Error(err))
		return fmt.Errorf("failed to check for modifications: %w", err)
	}

	if modified {
		log.Info("file modified, writing changes")
		newContent, err := svc.NewContent(session)
		if err != nil {
			log.Error("failed to open modified content", logger.Error(err))
			return fmt.Errorf("failed to open modified content: %w", err)
		}
		defer newContent.Close()

		if _, err := c.manager.WriteFile(ctx.Ctx, ctx.Options.URI, newContent, pkgfs.WriteOptions{Overwrite: true}); err != nil {
			log.Error("failed to write changes", logger.Error(err))
			return fmt.Errorf("failed to write changes: %w", err)
		}
	} else {
		log.Info("no changes detected")
	}

	return nil
}

// Finalize performs any necessary cleanup after the edit operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
