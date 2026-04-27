package logout

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
)

// Command orchestrates the logout flow for the active profile.
type Command struct {
	config   config.Service
	identity identity.Service
	log      logger.Logger
}

// NewCommand initializes a new instance of the logout Command.
func NewCommand(
	cfg config.Service,
	id identity.Service,
	l logger.Logger,
) *Command {
	return &Command{
		config:   cfg,
		identity: id,
		log:      l,
	}
}

// Validate prepares and validates the options for the logout operation.
func (c *Command) Validate(ctx *CommandContext) error {
	return ctx.Options.Validate()
}

// Execute performs the logout operation.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)
	opts := ctx.Options

	log.Info("starting logout flow")

	log.Debug("loading profile configuration")
	cfg, err := c.config.GetConfig(ctx.Ctx)
	if err != nil {
		log.Error("failed to load configuration", logger.Error(err))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	provider := cfg.Auth.Provider
	if provider == "" {
		log.Debug("no provider specified, defaulting to microsoft")
		provider = "microsoft"
	}

	log.Info("logging out from provider", logger.String("provider", provider), logger.String("identity", opts.IdentityID))
	if err := c.identity.Logout(ctx.Ctx, provider, opts.IdentityID); err != nil {
		log.Error("logout failed", logger.Error(err))
		return fmt.Errorf("logout failed: %w", err)
	}

	log.Info("logout successful")
	fmt.Fprintln(opts.Stdout, "Logged out successfully.")

	return nil
}

// Finalize performs any necessary cleanup after the logout operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
