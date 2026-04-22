package ls

import (
	"context"

	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Logger defines the interface required for logging within the ls command.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
	WithContext(ctx context.Context) logger.Logger
}

type ItemLister interface {
	List(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ListOptions) ([]pkgfs.Item, error)
}

type FormatCreator interface {
	Create(format formatting.Format) (formatting.OutputFormatter, error)
}

// URIFactory defines the interface for creating URIs.
type URIFactory interface {
	FromString(s string) (*fs.URI, error)
}

// Command executes the ls operation.
type Command struct {
	manager          ItemLister
	uriFactory       URIFactory
	formatterFactory FormatCreator
	log              Logger
}

// NewCommand initializes a new instance of the ls Command.
func NewCommand(m ItemLister, f URIFactory, ff FormatCreator, l Logger) *Command {
	return &Command{
		manager:          m,
		uriFactory:       f,
		formatterFactory: ff,
		log:              l,
	}
}

// Validate prepares and validates the options for the ls operation.
func (c *Command) Validate(ctx *CommandContext) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(ctx.Options.Path)
	if err != nil {
		return err
	}
	ctx.Options.URI = uri

	return ctx.Options.Validate()
}

// Execute lists items in a directory.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx).With(
		logger.String("path", ctx.Options.URI.String()),
		logger.Bool("recursive", ctx.Options.Recursive),
	)

	log.Info("starting ls operation")

	lsOpts := pkgfs.ListOptions{
		Recursive: ctx.Options.Recursive,
	}

	items, err := c.manager.List(ctx.Ctx, ctx.Options.URI, lsOpts)
	if err != nil {
		log.Error("list failed", logger.Error(err))
		return err
	}

	log.Debug("formatting output")
	formatter, err := c.formatterFactory.Create(ctx.Options.Format)
	if err != nil {
		log.Error("failed to create formatter", logger.Error(err))
		return err
	}

	itemsAny := make([]any, len(items))
	for i, item := range items {
		itemsAny[i] = item
	}

	if err := formatter.Format(ctx.Options.Stdout, itemsAny); err != nil {
		log.Error("format failed", logger.Error(err))
		return err
	}

	log.Info("ls completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the ls operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
