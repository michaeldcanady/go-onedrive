package ls

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/sorting"
)

// LsCmd handles the core execution logic for the 'ls' command.
// It coordinates fetching, filtering, sorting, and formatting OneDrive items.
type LsCmd struct {
	util.BaseCommand
}

// NewLsCmd creates a new LsCmd instance with the provided dependency container.
func NewLsCmd(container di.Container) *LsCmd {
	return &LsCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the full lifecycle of the ls command.
func (c *LsCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting ls command",
		logger.String("format", opts.Format),
		logger.Bool("includeAll", opts.IncludeAll),
		logger.Bool("foldersOnly", opts.FoldersOnly),
		logger.Bool("filesOnly", opts.FilesOnly),
		logger.String("sortProperty", opts.SortProperty),
		logger.String("sortDirection", opts.SortOrder.String()),
		logger.Bool("recursive", opts.Recursive),
	)

	fsSvc := c.Container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "filesystem service is nil")
	}

	items, err := c.fetchItems(ctx, fsSvc, opts.Path, opts.Recursive)
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	items, err = c.filterItems(items, opts)
	if err != nil {
		return util.NewCommandError(c.Name, "failed to filter items", err)
	}

	if opts.Format != "tree" {
		items, err = c.sortItems(items, opts)
		if err != nil {
			return util.NewCommandError(c.Name, "failed to sort items", err)
		}
	}

	if err := c.formatOutput(items, opts); err != nil {
		return util.NewCommandError(c.Name, "failed to format items", err)
	}

	c.Log.Info("ls command completed",
		logger.Duration("duration", time.Since(start)),
		logger.Int("finalItemCount", len(items)),
	)

	return nil
}

// fetchItems retrieves items from the OneDrive filesystem service.
func (c *LsCmd) fetchItems(ctx context.Context, fsSvc domainfs.Service, path string, recursive bool) ([]domainfs.Item, error) {
	item, err := fsSvc.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	items := []domainfs.Item{item}
	if item.Type == domainfs.ItemTypeFolder {
		listOpts := domainfs.ListOptions{
			Recursive: recursive,
		}
		items, err = fsSvc.List(ctx, path, listOpts)
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

// filterItems applies inclusion/exclusion filters to the retrieved items.
func (c *LsCmd) filterItems(items []domainfs.Item, opts Options) ([]domainfs.Item, error) {
	var filterOpts []filtering.FilterOption

	if opts.IncludeAll {
		filterOpts = append(filterOpts, filtering.IncludeAll())
	} else {
		filterOpts = append(filterOpts, filtering.ExcludeHidden())
	}

	if opts.FilesOnly {
		filterOpts = append(filterOpts, filtering.WithItemType(domainfs.ItemTypeFile))
	} else if opts.FoldersOnly {
		filterOpts = append(filterOpts, filtering.WithItemType(domainfs.ItemTypeFolder))
	}

	filterer, err := filtering.NewFilterFactory().Create(filterOpts...)
	if err != nil {
		return nil, err
	}

	return filterer.Filter(items)
}

// sortItems sorts the items based on the property and direction specified in Options.
func (c *LsCmd) sortItems(items []domainfs.Item, opts Options) ([]domainfs.Item, error) {
	var sortOpts []sorting.SortingOption

	if opts.SortProperty != "" {
		sortOpts = append(sortOpts, sorting.WithField(opts.SortProperty))
	}
	sortOpts = append(sortOpts, sorting.WithDirection(opts.SortOrder))

	sorter, err := sorting.NewSorterFactory().Create(sortOpts...)
	if err != nil {
		return nil, err
	}

	return sorter.Sort(items)
}

// formatOutput handles the final rendering of items to the user's terminal.
func (c *LsCmd) formatOutput(items []domainfs.Item, opts Options) error {
	formatter, err := formatting.NewFormatterFactory().Create(opts.Format)
	if err != nil {
		return err
	}

	return formatter.Format(opts.Stdout, items)
}
