package ls

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/sorting"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
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
	h.log.Info("listing directory", logger.String("path", opts.Path))

	provider, path, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}

	items, err := provider.List(ctx, path, shared.ListOptions{
		Recursive: opts.Recursive,
	})
	if err != nil {
		return fmt.Errorf("failed to list items at %s: %w", path, err)
	}

	// 1. Filtering
	filterFactory := filtering.NewFilterFactory()
	filterOpts := []filtering.FilterOption{}
	if opts.All {
		filterOpts = append(filterOpts, filtering.IncludeAll())
	}
	filterer, err := filterFactory.Create(filterOpts...)
	if err != nil {
		return fmt.Errorf("failed to create filterer: %w", err)
	}
	items, err = filterer.Filter(items)
	if err != nil {
		return fmt.Errorf("failed to filter items: %w", err)
	}

	// 2. Sorting
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
	items, err = sorter.Sort(items)
	if err != nil {
		return fmt.Errorf("failed to sort items: %w", err)
	}

	// Convert items to []any for formatter
	anyItems := make([]any, len(items))
	for i, v := range items {
		anyItems[i] = v
	}

	// 3. Formatting
	formatterFactory := formatting.NewFormatterFactory()
	formatter, err := formatterFactory.Create(opts.Format)
	if err != nil {
		return fmt.Errorf("failed to create formatter: %w", err)
	}

	return formatter.Format(opts.Stdout, anyItems)
}
