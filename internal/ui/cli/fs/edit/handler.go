package edit

import (
	"context"
	"fmt"
	"io"
	"os"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Logger defines the interface required for logging within the edit command.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
	Info(msg string, fields ...logger.Field)
	With(fields ...logger.Field) logger.Logger
	WithContext(ctx context.Context) logger.Logger
}

type Manager interface {
	ReadFile(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ReadOptions) (io.ReadCloser, error)
	WriteFile(ctx context.Context, uri *pkgfs.URI, r io.Reader, opts pkgfs.WriteOptions) (pkgfs.Item, error)
}

type Editor interface {
	CreateSession(ctx context.Context, remoteURI *fs.URI, reader io.Reader) (*editor.Session, error)
	Open(ctx context.Context, session *editor.Session) error
	Cleanup(ctx context.Context, session *editor.Session) error
	Modified(session *editor.Session) (bool, error)
	NewContent(session *editor.Session) (io.ReadCloser, error)
	WithOptions(opts ...editor.Option) editor.Service
}

// URIFactory defines the interface for creating URIs.
type URIFactory interface {
	FromString(s string) (*fs.URI, error)
}

// Command executes the edit operation.
type Command struct {
	manager    Manager
	uriFactory URIFactory
	editor     Editor
	log        Logger
}

// NewCommand initializes a new instance of the edit Command.
func NewCommand(m fs.Service, f *fs.URIFactory, e editor.Service, l Logger) *Command {
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
