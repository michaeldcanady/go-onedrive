package delete

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "delete"
	loggerID    = "cli"
)

type DeleteCmd struct {
	util.BaseCommand
}

func NewDeleteCmd(container di.Container) *DeleteCmd {
	return &DeleteCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *DeleteCmd) Run(ctx context.Context, opts Options, force bool) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	name := strings.ToLower(strings.TrimSpace(opts.Name))
	if name == "" {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	if name == profile.DefaultProfileName {
		return util.NewCommandErrorWithNameWithMessage(
			c.Name,
			"cannot delete the default profile",
		)
	}

	current, err := c.Container.State().GetCurrentProfile()
	if err != nil {
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	// If deleting the active profile, confirm unless forced
	if current == name && !force {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("You are deleting the active profile %q. Continue", name),
			IsConfirm: true,
			Stdout:    util.NewNopWriteCloser(opts.Stdout),
		}

		_, err := prompt.Run()
		if err != nil {
			fmt.Fprintln(opts.Stdout, "Aborted.")
			return nil
		}

		c.Log.Info("deleting current profile; switching to default")

		if err := c.Container.State().SetCurrentProfile(profile.DefaultProfileName); err != nil {
			return util.NewCommandErrorWithNameWithError(
				c.Name,
				fmt.Errorf("failed to switch to default profile: %w", err),
			)
		}
	}

	// Delete the profile directory
	if err := c.Container.Profile().Delete(ctx, name); err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	fmt.Fprintf(opts.Stdout, "Deleted profile %q\n", name)

	c.Log.Info("profile delete completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateDeleteCmd(container di.Container) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Name:   args[0],
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			return NewDeleteCmd(container).Run(cmd.Context(), opts, force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")

	return cmd
}
