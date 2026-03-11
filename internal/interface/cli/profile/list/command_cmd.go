// Package list provides the command-line interface for listing all OneDrive profiles.
package list

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

// ListCmd handles the execution logic for the 'profile list' command.
type ListCmd struct {
	util.BaseCommand
}

// NewListCmd creates a new ListCmd instance with the provided dependency container.
func NewListCmd(container didomain.Container) *ListCmd {
	return &ListCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the profile list command. It retrieves all configured profiles
// and displays their names to the user.
// It uses the domainprofile.ProfileService interface to decouple from the full container.
func (c *ListCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("listing profiles")

	profileSvc := c.Container.Profile()
	if profileSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "profile service is nil")
	}

	profiles, err := profileSvc.List(ctx)
	if err != nil {
		c.Log.Error("failed to list profiles", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	for _, p := range profiles {
		c.RenderMessage(opts.Stdout, "%s\n", p.Name)
	}

	c.Log.Info("profile list completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
