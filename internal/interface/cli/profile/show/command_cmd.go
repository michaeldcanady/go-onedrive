// Package show provides the command-line interface for displaying detailed information about a OneDrive profile.
package show

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

// ShowCmd handles the execution logic for the 'profile show' command.
type ShowCmd struct {
	util.BaseCommand
}

// NewShowCmd creates a new ShowCmd instance with the provided dependency container.
func NewShowCmd(container didomain.Container) *ShowCmd {
	return &ShowCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the profile show command. It retrieves the metadata for the
// specified profile and displays its name, filesystem path, and configuration path.
// It uses the domainprofile.ProfileService interface to decouple from the full container.
func (c *ShowCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("showing profile details", domainlogger.String("profile", opts.Name))

	profileSvc := c.Container.Profile()
	if profileSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "profile service is nil")
	}

	p, err := profileSvc.Get(ctx, opts.Name)
	if err != nil {
		c.Log.Error("failed to retrieve profile",
			domainlogger.String("profile", opts.Name),
			domainlogger.Error(err),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.RenderMessage(opts.Stdout, "Name: %s\n", p.Name)
	c.RenderMessage(opts.Stdout, "Path: %s\n", p.Path)
	if p.ConfigurationPath != "" {
		c.RenderMessage(opts.Stdout, "Config: %s\n", p.ConfigurationPath)
	}

	c.Log.Info("profile show completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
