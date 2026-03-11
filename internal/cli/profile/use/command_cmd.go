// Package use provides the command-line interface for selecting the active OneDrive profile.
package use

import (
	"context"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

// UseCmd handles the execution logic for the 'profile use' command.
type UseCmd struct {
	util.BaseCommand
}

// NewUseCmd creates a new UseCmd instance with the provided dependency container.
func NewUseCmd(container didomain.Container) *UseCmd {
	return &UseCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the profile use command. It validates that the requested profile
// exists and updates the global state to make it the active profile for future operations.
// It uses specific domain services to decouple from the full container.
func (c *UseCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	name := strings.TrimSpace(opts.Name)
	if name == "" {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	name = strings.ToLower(name)

	c.Log.Info("setting current profile", domainlogger.String("profile", name))

	profileSvc := c.Container.Profile()
	if profileSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "profile service is nil")
	}

	// Validate profile exists
	p, err := profileSvc.Get(ctx, name)
	if err != nil {
		c.Log.Error("failed to retrieve profile",
			domainlogger.String("profile", name),
			domainlogger.Error(err),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	// Persist as current profile
	if err := stateSvc.Set(domainstate.KeyProfile, p.Name, domainstate.ScopeGlobal); err != nil {
		c.Log.Error("failed to update current profile state",
			domainlogger.String("profile", p.Name),
			domainlogger.Error(err),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.RenderSuccess(opts.Stdout, "active profile set to %q", p.Name)

	c.Log.Info("profile use completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
