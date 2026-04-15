package ls

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the ls operation.
type Command struct {
	manager          fs.Service
	uriFactory       *fs.URIFactory
	formatterFactory *formatting.FormatterFactory
	log              logger.Logger
}

// NewCommand initializes a new instance of the ls Command.
func NewCommand(m fs.Service, f *fs.URIFactory, ff *formatting.FormatterFactory, l logger.Logger) *Command {
	return &Command{
		manager:          m,
		uriFactory:       f,
		formatterFactory: ff,
		log:              l,
	}
}

// Validate prepares and validates the options for the ls operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	// Resolve URI using the factory
	uri, err := c.uriFactory.FromString(opts.Path)
	if err != nil {
		return err
	}
	opts.URI = uri

	return opts.Validate()
}

// Execute lists items in a directory.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx).With(
		logger.String("path", opts.URI.String()),
		logger.Bool("recursive", opts.Recursive),
	)

	log.Info("starting ls operation")

	lsOpts := fs.ListOptions{
		Recursive: opts.Recursive,
	}

	items, err := c.manager.List(ctx, opts.URI, lsOpts)
	if err != nil {
		log.Error("list failed", logger.Error(err))
		return err
	}

	log.Debug("formatting output")
	formatter, err := c.formatterFactory.Create(opts.Format)
	if err != nil {
		log.Error("failed to create formatter", logger.Error(err))
		return err
	}

	itemsAny := make([]any, len(items))
	for i, item := range items {
		itemsAny[i] = item
	}

	if err := formatter.Format(opts.Stdout, itemsAny); err != nil {
		log.Error("format failed", logger.Error(err))
		return err
	}

	log.Info("ls completed successfully")
	return nil
}

// Finalize performs any necessary cleanup after the ls operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
