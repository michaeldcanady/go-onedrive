package delete

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/michaeldcanady/go-onedrive/internal/profile/infra"
)

type DeleteCmd struct {
	util.BaseCommand
}

func NewDeleteCmd(container didomain.Container) *DeleteCmd {
	return &DeleteCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *DeleteCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	name := strings.ToLower(strings.TrimSpace(opts.Name))
	if name == "" {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	if name == infra.DefaultProfileName {
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
	if current == name && !opts.Force {
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

		if err := c.Container.State().SetCurrentProfile(infra.DefaultProfileName); err != nil {
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
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
