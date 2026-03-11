package list

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/formatting"
	"github.com/spf13/cobra"
)

const (
	commandName = "list"
	loggerID    = "cli"
)

type ListCmd struct {
	util.BaseCommand
}

func NewListCmd(container di.Container) *ListCmd {
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
		c.Log.Warn("failed to retrieve drives", logger.Error(err))
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

	columns := []formatting.Column[*drive.Drive]{
		formatting.NewColumn(" ", func(item *drive.Drive) string {
			if item.ID == activeDriveID {
				return "*"
			}
			return ""
		}),
		formatting.NewColumn("Alias", func(item *drive.Drive) string {
			return aliasMap[item.ID]
		}),
		formatting.NewColumn("ID", func(item *drive.Drive) string { return item.ID }),
		formatting.NewColumn("Name", func(item *drive.Drive) string { return item.Name }),
		formatting.NewColumn("Type", func(item *drive.Drive) string { return string(item.Type) }),
	}

	formatter := formatting.NewTableFormatter(columns...).WithTruncate(true)

	if err := formatter.Format(opts.Stdout, drives); err != nil {
		c.Log.Warn("failed to format output", logger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("drive list completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists available drives",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewListCmd(container).Run(cmd.Context(), opts)
		},
	}
}
