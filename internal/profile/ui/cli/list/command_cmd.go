package list

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Handler executes the profile list operation.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile list Handler.
func NewHandler(
	p profile.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("profile-list")
	return &Handler{
		profiles: p,
		log:      cliLog,
	}
}

// Handle retrieves and displays all registered profiles.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	log.Info("listing all profiles")

	profiles, err := h.profiles.List(ctx)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	if len(profiles) == 0 {
		log.Info("no profiles found")
		fmt.Fprintln(opts.Stdout, "No profiles found.")
		return nil
	}

	log.Info("found profiles", logger.Int("count", len(profiles)))
	fmt.Fprintln(opts.Stdout, "Profiles:")
	for _, p := range profiles {
		fmt.Fprintf(opts.Stdout, "- %s\n", p.Name)
	}

	return nil
}
