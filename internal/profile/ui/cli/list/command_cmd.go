package list

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Handler executes the profile list operation.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile list Handler.
func NewHandler(p profile.Service, l logger.Logger) *Handler {
	return &Handler{
		profiles: p,
		log:      l,
	}
}

// Handle retrieves and displays all registered profiles.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("listing all profiles")

	profiles, err := h.profiles.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	if len(profiles) == 0 {
		fmt.Fprintln(opts.Stdout, "No profiles found.")
		return nil
	}

	fmt.Fprintln(opts.Stdout, "Profiles:")
	for _, p := range profiles {
		fmt.Fprintf(opts.Stdout, "- %s\n", p.Name)
	}

	return nil
}
