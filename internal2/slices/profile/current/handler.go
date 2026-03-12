package current

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
)

// Handler retrieves the name of the currently active profile.
type Handler struct {
	state state.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the profile current Handler.
func NewHandler(s state.Service, l logger.Logger) *Handler {
	return &Handler{
		state: s,
		log:   l,
	}
}

// Handle retrieves and displays the name of the current profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("retrieving current profile")

	profile, err := h.state.Get(state.KeyProfile)
	if err != nil {
		return fmt.Errorf("failed to get current profile: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "Current profile: %s\n", profile)
	return nil
}
