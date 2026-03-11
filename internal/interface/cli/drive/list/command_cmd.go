// Package list provides the command-line interface for listing available OneDrive drives.
package list

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/common/formatting"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

// ListCmd handles the execution logic for the 'drive list' command.
type ListCmd struct {
	util.BaseCommand
}

// NewListCmd creates a new ListCmd instance with the provided dependency container.
func NewListCmd(container didomain.Container) *ListCmd {
	return &ListCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the drive list command. It retrieves all drives associated with
// the authenticated account, identifies the active drive and any aliases,
// and displays them in a formatted table.
// It uses specific domain services to decouple from the full container.
func (c *ListCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting drive list command")

	driveSvc := c.Container.Drive()
	if driveSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "drive service is nil")
	}

	drives, err := driveSvc.ListDrives(ctx)
	if err != nil {
		c.Log.Error("failed to retrieve drives",
			domainlogger.Error(err),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	activeDriveID, err := stateSvc.Get(domainstate.KeyDrive)
	if err != nil {
		c.Log.Warn("failed to get current drive",
			domainlogger.Error(err),
		)
		// Non-fatal, just won't show active marker
	}

	aliases, err := stateSvc.ListDriveAliases()
	if err != nil {
		c.Log.Warn("failed to list drive aliases",
			domainlogger.Error(err),
		)
		// Non-fatal, just won't show aliases
	}

	// Prepare alias lookup
	aliasMap := make(map[string]string)
	for alias, driveID := range aliases {
		aliasMap[driveID] = alias
	}

	columns := []formatting.Column[*domaindrive.Drive]{
		formatting.NewColumn(" ", func(item *domaindrive.Drive) string {
			if item.ID == activeDriveID {
				return "*"
			}
			return ""
		}),
		formatting.NewColumn("Alias", func(item *domaindrive.Drive) string {
			return aliasMap[item.ID]
		}),
		formatting.NewColumn("ID", func(item *domaindrive.Drive) string { return item.ID }),
		formatting.NewColumn("Name", func(item *domaindrive.Drive) string { return item.Name }),
		formatting.NewColumn("Type", func(item *domaindrive.Drive) string { return item.Type.String() }),
	}

	formatter := formatting.NewTableFormatter(columns...).WithTruncate(true)

	if err := formatter.Format(opts.Stdout, drives); err != nil {
		c.Log.Error("failed to format output", domainlogger.Error(err))
		return util.NewCommandError(c.Name, "failed to format output", err)
	}

	c.Log.Info("drive list completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
