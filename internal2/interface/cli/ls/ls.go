package ls

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infrafiltering "github.com/michaeldcanady/go-onedrive/internal2/infra/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/sorting"
	"github.com/spf13/cobra"
)

const (
	allFlagLon   = "all"
	allFlagShort = "a"
	allFlagUsage = "show hidden items (names starting with '.')"

	formatFlagLong    = "format"
	formatFlagShort   = "f"
	formatFlagUsage   = "output format (e.g., json, yaml, long, short, tree)"
	formatFlagDefault = "short"

	loggerID    = "cli"
	commandName = "ls"

	filesOnlyFlagLong  = "files-only"
	filesOnlyFlagUsage = "show only files"

	foldersOnlyFlagLong  = "folders-only"
	foldersOnlyFlagUsage = "show only folders"

	sortFlagLong    = "sort"
	sortFlagUsage   = "sorts files by the specified field (e.g., name, size, modified)"
	sortFlagDefault = "name"

	recursiveFlagLong  = "recursive"
	recursiveFlagShort = "R"
	recursiveFlagUsage = ""
)

var (
	supportedFormats    = []string{"json", "yaml", "yml", "long", "short", "tree"}
	supportedProperties = []string{"name", "size", "modified"}
)

func CreateLSCmd(c di.Container) *cobra.Command {
	var (
		format       string
		includeAll   bool
		foldersOnly  bool
		filesOnly    bool
		recursive    bool
		sortProperty string
		sortOrder    = sorting.DirectionAscending
		sortOpts     = []sorting.SortingOption{sorting.WithDirection(sortOrder)}
		filterOpts   = []infrafiltering.FilterOption{}
	)

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),

		PreRunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
				cmd.SetContext(ctx)
			}

			logger, err := util.EnsureLogger(ctx, c, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command pre-run started",
				infralogging.String("command", commandName),

				infralogging.String("format", format),
				infralogging.Bool("include_all", includeAll),
				infralogging.Bool("folders_only", foldersOnly),
				infralogging.Bool("files_only", filesOnly),
				infralogging.String("sort_property", sortProperty),
				infralogging.String("sort_direction", sortOrder.String()),
				infralogging.Bool("recursive", recursive),
			)

			// Validate flags
			if foldersOnly && filesOnly {
				logger.Warn("invalid flag combination",
					infralogging.String("event", "validate_flags"),

					infralogging.Bool("folders_only", foldersOnly),
					infralogging.Bool("files_only", filesOnly),
				)
				return util.NewCommandErrorWithNameWithMessage(commandName, "can't use --folders-only and --files-only together")
			}

			if !slices.Contains(supportedFormats, format) {
				logger.Warn("unsupported format",
					infralogging.String("event", "validate_flags"),
					infralogging.String("format", format),
				)
				return util.NewCommandErrorWithNameWithMessage(commandName, fmt.Sprintf("unsupported format: %s; only supports: json, yaml/yml, long, or short", format))
			}

			if !slices.Contains(supportedProperties, sortProperty) {
				logger.Warn("unsupported sort property",
					infralogging.String("event", "validate_flags"),
					infralogging.String("sort_property", sortProperty),
				)
				return util.NewCommandErrorWithNameWithMessage(commandName, fmt.Sprintf("unsupported property: %s; only supports: name, size, or modified", sortProperty))
			}

			logger.Debug("registering formatters",
				infralogging.String("event", "register_formatters"),
			)

			// register formatters
			jsonFormatter := &formatting.JSONFormatter{}
			if err := c.Format().RegisterWithType("json", reflect.TypeOf([]domainfs.Item{}), jsonFormatter); err != nil {
				logger.Error("failed to register json formatter",
					infralogging.String("event", "register_formatters"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			yamlFormatter := &formatting.YAMLFormatter{}
			if err := c.Format().RegisterWithType("yaml", reflect.TypeOf([]domainfs.Item{}), yamlFormatter); err != nil {
				logger.Error("failed to register yaml formatter",
					infralogging.String("event", "register_formatters"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			humanLongFormatter := &formatting.HumanLongFormatter{}
			if err := c.Format().RegisterWithType("long", reflect.TypeOf([]domainfs.Item{}), humanLongFormatter); err != nil {
				logger.Error("failed to register long formatter",
					infralogging.String("event", "register_formatters"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			humanShortFormatter := &formatting.HumanShortFormatter{}
			if err := c.Format().RegisterWithType("short", reflect.TypeOf([]domainfs.Item{}), humanShortFormatter); err != nil {
				logger.Error("failed to register short formatter",
					infralogging.String("event", "register_formatters"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			treeFormatter := &formatting.TreeFormatter{}
			if err := c.Format().RegisterWithType("tree", reflect.TypeOf([]domainfs.Item{}), treeFormatter); err != nil {
				logger.Error("failed to register tree formatter",
					infralogging.String("event", "register_formatters"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			if format == "tree" && sortProperty != "" {
				logger.Debug("sort ignored for tree format",
					infralogging.String("event", "validate_flags"),
				)
			}

			logger.Debug("building filter options",
				infralogging.String("event", "build_filters"),
			)

			// Build filter options
			if includeAll {
				filterOpts = append(filterOpts, infrafiltering.IncludeAll())
			} else {
				filterOpts = append(filterOpts, infrafiltering.ExcludeHidden())
			}

			if filesOnly {
				filterOpts = append(filterOpts, infrafiltering.WithItemType(domainfs.ItemTypeFile))
			} else if foldersOnly {
				filterOpts = append(filterOpts, infrafiltering.WithItemType(domainfs.ItemTypeFolder))
			}

			filterOpt := infrafiltering.NewFilterOptions()
			if err := filterOpt.Apply(filterOpts); err != nil {
				logger.Error("failed to apply filter options",
					infralogging.String("event", "build_filters"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			optionFilter := infrafiltering.NewOptionsFilterer(*filterOpt)
			if err := c.Filter().RegisterWithType("option", reflect.TypeFor[[]domainfs.Item](), optionFilter); err != nil {
				logger.Error("failed to register filterer",
					infralogging.String("event", "build_filters"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			// Build sort options
			if sortProperty != "" {
				sortOpts = append(sortOpts, sorting.WithField(sortProperty))
			}

			logger.Info("command pre-run completed",
				infralogging.String("command", commandName),
			)

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
				cmd.SetContext(ctx)
			}

			logger, err := util.EnsureLogger(ctx, c, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command started",
				infralogging.String("command", commandName),

				infralogging.String("format", format),
				infralogging.Bool("include_all", includeAll),
				infralogging.Bool("folders_only", foldersOnly),
				infralogging.Bool("files_only", filesOnly),
				infralogging.String("sort_property", sortProperty),
				infralogging.String("sort_direction", sortOrder.String()),
				infralogging.Bool("recursive", recursive),
			)

			// Resolve path
			path := ""
			if len(args) > 0 {
				path = args[0]
			}

			logger.Debug("resolved path",
				infralogging.String("event", "resolve_path"),
				infralogging.String("path", path),
			)

			// Filesystem service
			fsSvc := c.FS()
			if fsSvc == nil {
				logger.Error("filesystem service is nil",
					infralogging.String("event", "resolve_fs_service"),
				)
				return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
			}

			logger.Debug("listing items",
				infralogging.String("event", "list_items"),
				infralogging.String("path", path),
				infralogging.Bool("recursive", recursive),
			)

			items, err := fsSvc.List(ctx, path, domainfs.ListOptions{
				Recursive: recursive,
			})
			if err != nil {
				logger.Error("failed to list items",
					infralogging.String("event", "list_items"),
					infralogging.Error(err),
					infralogging.String("path", path),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("items retrieved",
				infralogging.String("event", "list_items"),
				infralogging.Int("count", len(items)),
				infralogging.String("path", path),
			)

			// Filtering
			logger.Debug("applying filters",
				infralogging.String("event", "filter_items"),
			)

			if err := c.Filter().Filter("option", items); err != nil {
				logger.Error("filtering failed",
					infralogging.String("event", "filter_items"),
					infralogging.Error(err),
				)
				return util.NewCommandError(commandName, "failed to filter items", err)
			}

			logger.Info("items after filtering",
				infralogging.String("event", "filter_items"),
				infralogging.Int("count", len(items)),
			)

			// Sorting (skip for tree)
			if format != "tree" {
				logger.Debug("initializing sorter",
					infralogging.String("event", "sort_items"),
					infralogging.String("sort_property", sortProperty),
					infralogging.String("sort_direction", sortOrder.String()),
				)

				sorter, err := sorting.NewSorterFactory().Create(sortOpts...)
				if err != nil {
					logger.Error("failed to initialize sorter",
						infralogging.String("event", "sort_items"),
						infralogging.Error(err),
					)
					return util.NewCommandError(commandName, "failed to initialize sorter", err)
				}

				logger.Debug("sorting items",
					infralogging.String("event", "sort_items"),
				)

				items, err = sorter.Sort(items)
				if err != nil {
					logger.Error("sorting failed",
						infralogging.String("event", "sort_items"),
						infralogging.Error(err),
					)
					return util.NewCommandError(commandName, "failed to sort items", err)
				}
			}

			// Formatting
			logger.Debug("formatting output",
				infralogging.String("event", "format_items"),
				infralogging.String("format", format),
			)

			if err := c.Format().Format(cmd.OutOrStdout(), format, items); err != nil {
				logger.Error("formatting failed",
					infralogging.String("event", "format_items"),
					infralogging.Error(err),
				)
				return util.NewCommandError(commandName, "failed to format items", err)
			}

			logger.Info("command completed",
				infralogging.String("command", commandName),
				infralogging.Duration("duration", time.Since(start)),
				infralogging.Int("final_item_count", len(items)),
			)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&includeAll, allFlagLon, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatFlagLong, formatFlagShort, formatFlagDefault, formatFlagUsage)
	cmd.Flags().StringVar(&sortProperty, sortFlagLong, sortFlagDefault, sortFlagUsage)
	cmd.Flags().BoolVar(&foldersOnly, foldersOnlyFlagLong, false, foldersOnlyFlagUsage)
	cmd.Flags().BoolVar(&filesOnly, filesOnlyFlagLong, false, filesOnlyFlagUsage)
	cmd.Flags().BoolVarP(&recursive, recursiveFlagLong, recursiveFlagShort, false, recursiveFlagUsage)

	return cmd
}
