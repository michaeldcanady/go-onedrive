package ls

import (
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

		PreRunE: func(_ *cobra.Command, _ []string) error {
			// Validate flags
			if foldersOnly && filesOnly {
				return util.NewCommandErrorWithNameWithMessage(commandName, "can't use --folders-only and --files-only together")
			}

			if !slices.Contains(supportedFormats, format) {
				return util.NewCommandErrorWithNameWithMessage(commandName, fmt.Sprintf("unsupported format: %s; only supports: json, yaml/yml, long, or short", format))
			}

			if !slices.Contains(supportedProperties, sortProperty) {
				return util.NewCommandErrorWithNameWithMessage(commandName, fmt.Sprintf("unsupported property: %s; only supports: name, size, or modified", sortProperty))
			}

			// register formatters
			jsonFormatter := &formatting.JSONFormatter{}
			if err := c.Format().RegisterWithType("json", reflect.TypeOf([]domainfs.Item{}), jsonFormatter); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			yamlFormatter := &formatting.YAMLFormatter{}
			if err := c.Format().RegisterWithType("yaml", reflect.TypeOf([]domainfs.Item{}), yamlFormatter); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			humanLongFormatter := &formatting.HumanLongFormatter{}
			if err := c.Format().RegisterWithType("long", reflect.TypeOf([]domainfs.Item{}), humanLongFormatter); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			humanShortFormatter := &formatting.HumanShortFormatter{}
			if err := c.Format().RegisterWithType("short", reflect.TypeOf([]domainfs.Item{}), humanShortFormatter); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			treeFormatter := &formatting.TreeFormatter{}
			if err := c.Format().RegisterWithType("tree", reflect.TypeOf([]domainfs.Item{}), treeFormatter); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			if format == "tree" && sortProperty != "" {
				// TODO: warn user that sort can't be used with tree format
			}

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
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}
			optionFilter := infrafiltering.NewOptionsFilterer(*filterOpt)
			if err := c.Filter().RegisterWithType("option", reflect.TypeFor[[]domainfs.Item](), optionFilter); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			// Build sort options
			if sortProperty != "" {
				sortOpts = append(sortOpts, sorting.WithField(sortProperty))
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()

			logger, err := util.EnsureLogger(c, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("starting ls command",
				infralogging.String("format", format),
				infralogging.Bool("includeAll", includeAll),
				infralogging.Bool("foldersOnly", foldersOnly),
				infralogging.Bool("filesOnly", filesOnly),
				infralogging.String("sortProperty", sortProperty),
				infralogging.String("sortDirection", sortOrder.String()),
			)

			// Filesystem service
			fsSvc := c.FS()
			if fsSvc == nil {
				logger.Error("filesystem service is nil")
				return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
			}

			// Resolve path
			path := ""
			if len(args) > 0 {
				path = args[0]
			}
			logger.Debug("path resolved", infralogging.String("path", path))

			logger.Debug("listing items from filesystem")
			items, err := fsSvc.List(cmd.Context(), path, domainfs.ListOptions{
				Recursive: recursive,
			})
			if err != nil {
				logger.Error("failed to list items", infralogging.String("error", err.Error()))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}
			logger.Info("items retrieved", infralogging.Int("count", len(items)))

			// Filtering
			logger.Debug("applying filters")
			if err := c.Filter().Filter("option", items); err != nil {
				logger.Error("failed to filter items", infralogging.String("error", err.Error()))
				return util.NewCommandError(commandName, "failed to filter items", err)
			}
			logger.Info("items after filtering", infralogging.Int("count", len(items)))

			if format != "tree" {
				// Sorting
				logger.Debug("initializing sorter")
				sorter, err := sorting.NewSorterFactory().Create(sortOpts...)
				if err != nil {
					logger.Error("failed to initialize sorter", infralogging.String("error", err.Error()))
					return util.NewCommandError(commandName, "failed to initialize sorter", err)
				}

				logger.Debug("sorting items")
				items, err = sorter.Sort(items)
				if err != nil {
					logger.Error("failed to sort items", infralogging.String("error", err.Error()))
					return util.NewCommandError(commandName, "failed to sort items", err)
				}
			}

			// Formatting
			logger.Debug("formatting output", infralogging.String("format", format))
			if err := c.Format().Format(cmd.OutOrStdout(), format, items); err != nil {
				logger.Error("failed to format items", infralogging.String("error", err.Error()))
				return util.NewCommandError(commandName, "failed to format items", err)
			}

			logger.Info("ls command completed",
				infralogging.Duration("duration", time.Since(start)),
				infralogging.Int("finalItemCount", len(items)),
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
