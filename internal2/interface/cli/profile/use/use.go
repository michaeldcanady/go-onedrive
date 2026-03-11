package use

import (
	"context"
	"fmt"
	"strings"
	"time"

	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
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

	name := strings.TrimSpace(opts.Name)
	if name == "" {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	name = strings.ToLower(name)

	c.Log.Info("setting current profile", logger.String("profile", name))

	// Validate profile exists
	p, err := c.Container.Profile().Get(ctx, name)
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	// Persist as current profile
	if err := c.Container.State().SetCurrentProfile(p.Name); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	fmt.Fprintf(opts.Stdout, "Active profile set to %q\n", p.Name)

	c.Log.Info("profile use completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateUseCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Set current profile",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Name:   args[0],
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewUseCmd(container).Run(cmd.Context(), opts)
		},
	}
}
