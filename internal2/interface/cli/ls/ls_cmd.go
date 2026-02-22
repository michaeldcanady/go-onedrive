package ls

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/sorting"
	"github.com/spf13/cobra"
)

// LsCmd handles the core execution logic for the 'ls' command.
// It coordinates fetching, filtering, sorting, and formatting OneDrive items.
type LsCmd struct {
	container di.Container
	logger    infralogging.Logger
}

// NewLsCmd creates a new LsCmd instance with the provided dependency container.
func NewLsCmd(container di.Container) *LsCmd {
	return &LsCmd{
		container: container,
	}
}

// WithLogger allows injecting a logger into LsCmd.
func (c *LsCmd) WithLogger(logger infralogging.Logger) *LsCmd {
	c.logger = logger
	return c
}

// Run executes the full lifecycle of the ls command.
// 1. Ensures a logger is available.
// 2. Fetches items from the filesystem service.
// 3. Filters items based on user options.
// 4. Sorts items (unless in tree format).
// 5. Formats and writes the output to the command's stdout.
func (c *LsCmd) Run(ctx context.Context, cmd *cobra.Command, opts Options) error {
	start := time.Now()

	if ctx == nil {
		ctx = context.Background()
	}

	if c.logger == nil {
		logger, err := util.EnsureLogger(c.container, loggerID)
		if err != nil {
			return util.NewCommandErrorWithNameWithError(commandName, err)
		}
		c.logger = logger
	}

	c.logger.Info("starting ls command",
		infralogging.String("format", opts.Format),
		infralogging.Bool("includeAll", opts.IncludeAll),
		infralogging.Bool("foldersOnly", opts.FoldersOnly),
		infralogging.Bool("filesOnly", opts.FilesOnly),
		infralogging.String("sortProperty", opts.SortProperty),
		infralogging.String("sortDirection", opts.SortOrder.String()),
		infralogging.Bool("recursive", opts.Recursive),
	)

	fsSvc := c.container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	items, err := c.fetchItems(ctx, fsSvc, opts.Path, opts.Recursive)
	if err != nil {
		return util.NewCommandErrorWithNameWithError(commandName, err)
	}

	items, err = c.filterItems(items, opts)
	if err != nil {
		return util.NewCommandError(commandName, "failed to filter items", err)
	}

	if opts.Format != "tree" {
		items, err = c.sortItems(items, opts)
		if err != nil {
			return util.NewCommandError(commandName, "failed to sort items", err)
		}
	}

	if err := c.formatOutput(cmd, items, opts.Format); err != nil {
		return util.NewCommandError(commandName, "failed to format items", err)
	}

	c.logger.Info("ls command completed",
		infralogging.Duration("duration", time.Since(start)),
		infralogging.Int("finalItemCount", len(items)),
	)

	return nil
}

// fetchItems retrieves items from the OneDrive filesystem service.
// If the target path is a folder, it lists its children.
func (c *LsCmd) fetchItems(ctx context.Context, fsSvc domainfs.Service, path string, recursive bool) ([]domainfs.Item, error) {
	c.logger.Debug("path resolved", infralogging.String("path", path))

	item, err := fsSvc.Get(ctx, path)
	if err != nil {
		c.logger.Error("failed to get item", infralogging.String("error", err.Error()))
		return nil, err
	}

	items := []domainfs.Item{item}
	if item.Type == domainfs.ItemTypeFolder {
		c.logger.Debug("listing items from filesystem")
		listOpts := domainfs.ListOptions{
			Recursive: recursive,
		}
		items, err = fsSvc.List(ctx, path, listOpts)
		if err != nil {
			c.logger.Error("failed to list items", infralogging.String("error", err.Error()))
			return nil, err
		}
		c.logger.Info("items retrieved", infralogging.Int("count", len(items)))
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

	c.logger.Debug("initializing filterer")
	filterer, err := filtering.NewFilterFactory().Create(filterOpts...)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("applying filters")
	return filterer.Filter(items)
}

// sortItems sorts the items based on the property and direction specified in Options.
func (c *LsCmd) sortItems(items []domainfs.Item, opts Options) ([]domainfs.Item, error) {
	var sortOpts []sorting.SortingOption

	if opts.SortProperty != "" {
		sortOpts = append(sortOpts, sorting.WithField(opts.SortProperty))
	}
	sortOpts = append(sortOpts, sorting.WithDirection(opts.SortOrder))

	c.logger.Debug("initializing sorter")
	sorter, err := sorting.NewSorterFactory().Create(sortOpts...)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("sorting items")
	return sorter.Sort(items)
}

// formatOutput handles the final rendering of items to the user's terminal.
func (c *LsCmd) formatOutput(cmd *cobra.Command, items []domainfs.Item, format string) error {
	c.logger.Debug("initializing formatter", infralogging.String("format", format))
	formatter, err := formatting.NewFormatterFactory().Create(format)
	if err != nil {
		return err
	}

	c.logger.Debug("formatting output")
	return formatter.Format(cmd.OutOrStdout(), items)
}
