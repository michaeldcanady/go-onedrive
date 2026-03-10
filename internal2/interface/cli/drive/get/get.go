package get

import (
	"context"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "get"
	loggerID    = "cli"
)

type GetCmd struct {
	util.BaseCommand
}

func NewGetCmd(container di.Container) *GetCmd {
	return &GetCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *GetCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	id := strings.ToLower(strings.TrimSpace(opts.DriveIDOrAlias))
	if id == "" {
		c.Log.Warn("id is empty", logger.String("command", c.Name))
		return util.NewCommandErrorWithNameWithMessage(c.Name, "id is empty")
	}

	c.Log.Info("retrieving drive details", logger.String("target", id))

	drive, err := c.Container.Drive().ResolveDrive(ctx, id)
	if err != nil {
		c.Log.Warn("failed to retrieve drive", logger.Error(err), logger.String("target", id))
		c.RenderError(opts.Stderr, err)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	formatter := NewTableFormatter(driveIDColumn, driveNameColumn, driveOwnerColumn, driveReadOnlyColumn, driveTypeColumn)
	if err := formatter.Format(opts.Stdout, []*domaindrive.Drive{drive}); err != nil {
		c.Log.Warn("failed to format output", logger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("drive get completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateGetCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id|alias>",
		Short: "Get information of named drive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveIDOrAlias: args[0],
				Stdout:         cmd.OutOrStdout(),
				Stderr:         cmd.ErrOrStderr(),
			}

			return NewGetCmd(container).Run(cmd.Context(), opts)
		},
	}
}
