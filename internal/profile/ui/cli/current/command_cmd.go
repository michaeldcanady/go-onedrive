package current

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Handler retrieves the name of the currently active profile.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile current Handler.
func NewHandler(p profile.Service, l logger.Logger) *Handler {
	return &Handler{
		profiles: p,
		log:      l,
	}
}

// Handle retrieves and displays the name of the current profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("retrieving current profile")

	p, err := h.profiles.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current profile: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "Current profile: %s\n", p.Name)
	return nil
}
