package current

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "current"
	loggerID    = "cli"
)

type CurrentCmd struct {
	util.BaseCommand
}

func NewCurrentCmd(container di.Container) *CurrentCmd {
	return &CurrentCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CurrentCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("retrieving current profile")

	name, err := c.Container.State().GetCurrentProfile()
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	fmt.Fprintf(opts.Stdout, "%s\n", name)

	c.Log.Info("profile current completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateCurrentCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewCurrentCmd(container).Run(cmd.Context(), opts)
		},
	}
}
