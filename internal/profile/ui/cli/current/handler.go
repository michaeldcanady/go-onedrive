package current

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Command executes the profile current operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile current Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile current operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute retrieves and displays the active profile.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Debug("fetching active profile")
	p, err := c.profile.GetActive(ctx)
	if err != nil {
		log.Error("failed to get active profile", logger.Error(err))
		return fmt.Errorf("failed to get active profile: %w", err)
	}

	log.Info("active profile retrieved successfully", logger.String("name", p.Name))
	fmt.Fprintf(opts.Stdout, "Active profile: %s\n", p.Name)
	return nil
}

// Finalize performs any necessary cleanup after the profile current operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
