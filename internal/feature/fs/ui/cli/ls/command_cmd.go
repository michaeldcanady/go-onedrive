package ls

import (
	"context"
	"fmt"

	registry "github.com/michaeldcanady/go-onedrive/internal/feature/fs"
	"github.com/michaeldcanady/go-onedrive/internal/feature/fs/filtering"
	"github.com/michaeldcanady/go-onedrive/internal/feature/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/feature/fs/sorting"
	"github.com/michaeldcanady/go-onedrive/internal/feature/logger"
)

// Handler executes the drive ls operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive ls Handler.
func NewHandler(fs registry.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:  fs,
		log: l,
	}
}

// Handle retrieves, filters, sorts, and displays the contents of a directory.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("listing directory")

	log.Debug("fetching items from provider")
	items, err := h.fs.List(ctx, opts.Path, registry.ListOptions{
		Recursive: opts.Recursive,
	})
	if err != nil {
		log.Error("failed to list items", logger.Error(err))
		return fmt.Errorf("failed to list items at %s: %w", opts.Path, err)
	}
	log.Info("retrieved items from provider", logger.Int("count", len(items)))

	// 1. Filtering
	log.Debug("preparing filter")
	filterFactory := filtering.NewFilterFactory()
	filterOpts := []filtering.FilterOption{}
	if opts.All {
		filterOpts = append(filterOpts, filtering.IncludeAll())
	}
	filterer, err := filterFactory.Create(filterOpts...)
	if err != nil {
		return fmt.Errorf("failed to create filterer: %w", err)
	}

	log.Debug("filtering items")
	items, err = filterer.Filter(items)
	if err != nil {
		log.Error("failed to filter items", logger.Error(err))
		return fmt.Errorf("failed to filter items: %w", err)
	}
	log.Debug("items filtered", logger.Int("count", len(items)))

	// 2. Sorting
	log.Debug("preparing sorter")
	sortFactory := sorting.NewSorterFactory()
	sortOpts := []sorting.SortingOption{}
	if opts.SortField != "" {
		sortOpts = append(sortOpts, sorting.WithField(opts.SortField))
	}
	if opts.SortDescending {
		sortOpts = append(sortOpts, sorting.WithDirection(sorting.DirectionDescending))
	}
	sorter, err := sortFactory.Create(sortOpts...)
	if err != nil {
		return fmt.Errorf("failed to create sorter: %w", err)
	}

	log.Debug("sorting items")
	items, err = sorter.Sort(items)
	if err != nil {
		log.Error("failed to sort items", logger.Error(err))
		return fmt.Errorf("failed to sort items: %w", err)
	}
	log.Debug("items sorted")

	// Convert items to []any for formatter
	anyItems := make([]any, len(items))
	for i, v := range items {
		anyItems[i] = v
	}

	// 3. Formatting
	log.Debug("formatting output", logger.String("format", opts.Format.String()))
	formatterFactory := formatting.NewFormatterFactory()
	formatter, err := formatterFactory.Create(opts.Format)
	if err != nil {
		return fmt.Errorf("failed to create formatter: %w", err)
	}

	return formatter.Format(opts.Stdout, anyItems)
}
