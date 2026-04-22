package list

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
)

type Handler struct {
	profile profile.Service
	log     logger.Logger
}

func NewHandler(p profile.Service, l logger.Logger) *Handler {
	return &Handler{
		profile: p,
		log:     l,
	}
}

// Validate ensures options and environment are correct.
func (h *Handler) Validate(ctx *CommandContext) error {
	if ctx.Options.Stdout == nil {
		return fmt.Errorf("stdout must not be nil")
	}
	return nil
}

// Execute performs the core business logic.
func (h *Handler) Execute(ctx *CommandContext) error {
	log := h.log.WithContext(ctx.Ctx)

	log.Debug("fetching all profiles")
	profiles, err := h.profile.List(ctx.Ctx)
	if err != nil {
		log.Error("failed to list profiles", logger.Error(err))
		return fmt.Errorf("failed to list profiles: %w", err)
	}
	ctx.Profiles = profiles

	log.Debug("fetching active profile")
	active, err := h.profile.GetActive(ctx.Ctx)
	if err != nil {
		log.Error("failed to retrieve active profile", logger.Error(err))
		return fmt.Errorf("failed to retrieve active profile: %w", err)
	}

	ctx.Active = &active

	return nil
}

// Finalize handles presentation or cleanup.
func (h *Handler) Finalize(c *CommandContext) error {
	fmt.Fprintln(c.Options.Stdout, "Available profiles:")
	for _, p := range c.Profiles {
		prefix := "  "
		if c.Active != nil && p.Name == c.Active.Name {
			prefix = "* "
		}
		fmt.Fprintf(c.Options.Stdout, "%s%s\n", prefix, p.Name)
	}
	return nil
}
