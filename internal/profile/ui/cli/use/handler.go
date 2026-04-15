package use

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Command executes the profile use operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile use Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile use operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute sets the active profile.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Info("switching active profile", logger.String("name", opts.Name))
	if err := c.profile.SetActive(ctx, opts.Name, state.ScopeGlobal); err != nil {
		log.Error("failed to switch profile", logger.String("name", opts.Name), logger.Error(err))
		return fmt.Errorf("failed to switch to profile %s: %w", opts.Name, err)
	}

	log.Info("switched active profile successfully", logger.String("name", opts.Name))
	fmt.Fprintf(opts.Stdout, "Switched to profile '%s'.\n", opts.Name)
	return nil
}

// Finalize performs any necessary cleanup after the profile use operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
