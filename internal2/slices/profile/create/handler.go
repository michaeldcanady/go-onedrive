package create

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/profile"
)

// Handler executes the profile creation operation.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile creation Handler.
func NewHandler(p profile.Service, l logger.Logger) *Handler {
	return &Handler{
		profiles: p,
		log:      l,
	}
}

// Handle executes the logic to create a new profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("creating new profile", logger.String("name", opts.Name))

	p, err := h.profiles.Create(ctx, opts.Name)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "Profile %s created successfully.\n", p.Name)
	return nil
}
