// Package list provides the command-line interface for displaying registered OneDrive drive aliases.
package list

import (
	"context"
	"sort"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/common/formatting"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

// ListCmd handles the execution logic for the 'drive alias list' command.
type ListCmd struct {
	util.BaseCommand
}

// NewListCmd creates a new ListCmd instance with the provided dependency container.
func NewListCmd(container didomain.Container) *ListCmd {
	return &ListCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// aliasItem represents a single row in the alias list table.
type aliasItem struct {
	// Alias is the friendly name assigned to the drive.
	Alias string
	// DriveID is the unique identifier of the OneDrive drive.
	DriveID string
}

// Run executes the drive alias list command. It retrieves all registered aliases
// from the global state and displays them in a formatted table.
// It uses the domainstate.Service interface to decouple from the full container.
func (c *ListCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("listing drive aliases")

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	aliases, err := stateSvc.ListDriveAliases()
	if err != nil {
		c.Log.Error("failed to list drive aliases",
			domainlogger.Error(err),
		)
		return util.NewCommandError(c.Name, "failed to list drive aliases", err)
	}

	if len(aliases) == 0 {
		c.RenderInfo(opts.Stdout, "no drive aliases found")
		return nil
	}

	items := make([]aliasItem, 0, len(aliases))
	for k, v := range aliases {
		items = append(items, aliasItem{Alias: k, DriveID: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Alias < items[j].Alias
	})

	formatter := formatting.NewTableFormatter(
		formatting.NewColumn("ALIAS", func(i aliasItem) string { return i.Alias }),
		formatting.NewColumn("DRIVE ID", func(i aliasItem) string { return i.DriveID }),
	)

	if err := formatter.Format(opts.Stdout, items); err != nil {
		c.Log.Error("failed to format output", domainlogger.Error(err))
		return util.NewCommandError(c.Name, "failed to format output", err)
	}

	c.Log.Info("drive alias list completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
