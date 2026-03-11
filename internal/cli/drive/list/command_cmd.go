package list

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/common/formatting"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
)

type ListCmd struct {
	util.BaseCommand
}

func NewListCmd(container didomain.Container) *ListCmd {
	return &ListCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *ListCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting drive list command")

	drives, err := c.Container.Drive().ListDrives(ctx)
	if err != nil {
		c.Log.Warn("failed to retrieve drives", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	stateSvc := c.Container.State()
	activeDriveID, _ := stateSvc.GetCurrentDrive()
	aliases, _ := stateSvc.ListDriveAliases()

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
		formatting.NewColumn("Type", func(item *domaindrive.Drive) string { return string(item.Type) }),
	}

	formatter := formatting.NewTableFormatter(columns...).WithTruncate(true)

	if err := formatter.Format(opts.Stdout, drives); err != nil {
		c.Log.Warn("failed to format output", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("drive list completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
