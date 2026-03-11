// Package delete provides the command-line interface for removing OneDrive profiles.
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
	infraprofile "github.com/michaeldcanady/go-onedrive/internal/profile/infra"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

// DeleteCmd handles the execution logic for the 'profile delete' command.
type DeleteCmd struct {
	util.BaseCommand
}

// NewDeleteCmd creates a new DeleteCmd instance with the provided dependency container.
func NewDeleteCmd(container didomain.Container) *DeleteCmd {
	return &DeleteCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the profile delete command. It removes the specified profile's
// configuration and data. If the profile being deleted is the active one,
// it prompts for confirmation or switches to the default profile.
// It uses specific domain services to decouple from the full container.
func (c *DeleteCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	name := strings.ToLower(strings.TrimSpace(opts.Name))
	if name == "" {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	if name == infraprofile.DefaultProfileName {
		return util.NewCommandErrorWithNameWithMessage(
			c.Name,
			"cannot delete the default profile",
		)
	}

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	current, err := stateSvc.Get(domainstate.KeyProfile)
	if err != nil {
		c.Log.Warn("failed to retrieve current profile from state", domainlogger.Error(err))
		// Continue anyway, we just might miss the confirmation logic
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
			c.RenderInfo(opts.Stdout, "aborted")
			return nil
		}

		c.Log.Info("deleting current profile; switching to default")

		if err := stateSvc.Set(domainstate.KeyProfile, infraprofile.DefaultProfileName, domainstate.ScopeGlobal); err != nil {
			return util.NewCommandErrorWithNameWithError(
				c.Name,
				fmt.Errorf("failed to switch to default profile: %w", err),
			)
		}
	}

	profileSvc := c.Container.Profile()
	if profileSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "profile service is nil")
	}

	// Delete the profile directory
	if err := profileSvc.Delete(ctx, name); err != nil {
		c.Log.Error("failed to delete profile",
			domainlogger.String("profile", name),
			domainlogger.Error(err),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.RenderSuccess(opts.Stdout, "deleted profile %q", name)

	c.Log.Info("profile delete completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
