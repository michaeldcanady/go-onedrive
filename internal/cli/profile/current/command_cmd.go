// Package current provides the command-line interface for displaying the active OneDrive profile.
package current

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

// CurrentCmd handles the execution logic for the 'profile current' command.
type CurrentCmd struct {
	util.BaseCommand
}

// NewCurrentCmd creates a new CurrentCmd instance with the provided dependency container.
func NewCurrentCmd(container didomain.Container) *CurrentCmd {
	return &CurrentCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the profile current command. It retrieves the name of the
// active profile from the global state and displays it to the user.
// It uses the domainstate.Service interface to decouple from the full container.
func (c *CurrentCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("retrieving current profile")

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	name, err := stateSvc.Get(domainstate.KeyProfile)
	if err != nil {
		c.Log.Error("failed to get current profile from state", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.RenderMessage(opts.Stdout, "%s\n", name)

	c.Log.Info("profile current completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
