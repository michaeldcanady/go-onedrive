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
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile switch Handler.
func NewHandler(p profile.Service, l logger.Logger) *Handler {
	return &Handler{
		profiles: p,
		log:      l,
	}
}

// Handle executes the logic to switch the active profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("switching to profile", logger.String("name", opts.Name))

	if err := h.profiles.SetActive(ctx, opts.Name, state.ScopeGlobal); err != nil {
		return fmt.Errorf("failed to set active profile: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "Switched to profile: %s\n", opts.Name)
	return nil
}
