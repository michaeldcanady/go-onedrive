package show

import (
	"context"
	"fmt"
	"time"

	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "show"
	loggerID    = "cli"
)

type ShowCmd struct {
	util.BaseCommand
}

func NewShowCmd(container di.Container) *ShowCmd {
	return &ShowCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *ShowCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("showing profile details", logger.String("profile", opts.Name))

	p, err := c.Container.Profile().Get(ctx, opts.Name)
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return err
	}

	fmt.Fprintf(opts.Stdout, "Name: %s\n", p.Name)
	fmt.Fprintf(opts.Stdout, "Path: %s\n", p.Path)
	if p.ConfigurationPath != "" {
		fmt.Fprintf(opts.Stdout, "Config: %s\n", p.ConfigurationPath)
	}

	c.Log.Info("profile show completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateShowCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show details for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Name:   args[0],
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewShowCmd(container).Run(cmd.Context(), opts)
		},
	}
}
