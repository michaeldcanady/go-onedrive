package list

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

	c.Log.Info("listing profiles")

	profiles, err := c.Container.Profile().List(ctx)
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return err
	}

	for _, p := range profiles {
		fmt.Fprintf(opts.Stdout, "%s\n", p.Name)
	}

	c.Log.Info("profile list completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewListCmd(container).Run(cmd.Context(), opts)
		},
	}
}
