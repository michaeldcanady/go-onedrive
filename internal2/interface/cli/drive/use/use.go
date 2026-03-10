package use

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "use"
	loggerID    = "cli"
)

type UseCmd struct {
	util.BaseCommand
}

func NewUseCmd(container di.Container) *UseCmd {
	return &UseCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *UseCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting drive use command", logger.String("target", opts.DriveIDOrAlias))

	resolvedDrive, err := c.Container.Drive().ResolveDrive(ctx, opts.DriveIDOrAlias)
	if err != nil {
		c.Log.Warn("failed to resolve drive",
			logger.Error(err),
			logger.String("target", opts.DriveIDOrAlias),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	if err := c.Container.State().SetCurrentDrive(resolvedDrive.ID); err != nil {
		c.Log.Warn("failed to update current drive state",
			logger.Error(err),
			logger.String("driveID", resolvedDrive.ID),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	fmt.Fprintf(opts.Stdout, "Now using drive: %s (%s)\n", resolvedDrive.Name, resolvedDrive.ID)

	c.Log.Info("drive use completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateUseCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use [drive-id|alias]",
		Short: "Sets the active drive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveIDOrAlias: args[0],
				Stdout:         cmd.OutOrStdout(),
				Stderr:         cmd.ErrOrStderr(),
			}

			return NewUseCmd(container).Run(cmd.Context(), opts)
		},
	}
}
