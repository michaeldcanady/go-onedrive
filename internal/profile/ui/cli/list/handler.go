package list

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Command executes the profile list operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile list Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile list operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute retrieves and displays all profiles.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Debug("fetching all profiles")
	profiles, err := c.profile.List(ctx)
	if err != nil {
		log.Error("failed to list profiles", logger.Error(err))
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	log.Debug("fetching active profile for marking")
	current, _ := c.profile.GetActive(ctx)

	log.Info("profiles retrieved successfully", logger.Int("count", len(profiles)))
	fmt.Fprintln(opts.Stdout, "Available profiles:")
	for _, p := range profiles {
		prefix := "  "
		if current.Name != "" && p.Name == current.Name {
			prefix = "* "
		}
		fmt.Fprintf(opts.Stdout, "%s%s\n", prefix, p.Name)
	}
	return nil
}

// Finalize performs any necessary cleanup after the profile list operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
