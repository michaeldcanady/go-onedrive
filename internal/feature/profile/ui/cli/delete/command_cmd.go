package delete

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/feature/logger"
	"github.com/michaeldcanady/go-onedrive/internal/feature/profile"
)

// Handler executes the profile deletion operation.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile deletion Handler.
func NewHandler(p profile.Service, l logger.Logger) *Handler {
	return &Handler{
		profiles: p,
		log:      l,
	}
}

// Handle executes the logic to delete a profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("deleting profile", logger.String("name", opts.Name))

	if err := h.profiles.Delete(ctx, opts.Name); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "Profile %s deleted successfully.\n", opts.Name)
	return nil
}
