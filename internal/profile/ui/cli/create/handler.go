package create

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Command executes the profile create operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile create Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile create operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute creates a new profile.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Info("creating profile", logger.String("name", opts.Name))
	if _, err := c.profile.Create(ctx, opts.Name); err != nil {
		log.Error("failed to create profile", logger.String("name", opts.Name), logger.Error(err))
		return fmt.Errorf("failed to create profile %s: %w", opts.Name, err)
	}

	log.Info("profile created successfully", logger.String("name", opts.Name))
	fmt.Fprintf(opts.Stdout, "Profile '%s' created successfully.\n", opts.Name)
	return nil
}

// Finalize performs any necessary cleanup after the profile create operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
