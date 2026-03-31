package use

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Handler executes the profile switch operation.
type Handler struct {
	profiles profile.Service
	state    state.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile switch Handler.
func NewHandler(p profile.Service, s state.Service, l logger.Logger) *Handler {
	return &Handler{
		profiles: p,
		state:    s,
		log:      l,
	}
}

// Handle executes the logic to switch the active profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("switching to profile", logger.String("name", opts.Name))

	exists, err := h.profiles.Exists(ctx, opts.Name)
	if err != nil {
		return fmt.Errorf("failed to check profile existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("profile %s does not exist", opts.Name)
	}

	if err := h.state.Set(state.KeyProfile, opts.Name, state.ScopeGlobal); err != nil {
		return fmt.Errorf("failed to set active profile: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "Switched to profile: %s\n", opts.Name)
	return nil
}
