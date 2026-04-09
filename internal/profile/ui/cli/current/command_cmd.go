package current

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Handler retrieves the name of the currently active profile.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile current Handler.
func NewHandler(
	p profile.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("profile-current")
	return &Handler{
		profiles: p,
		log:      cliLog,
	}
}

// Handle retrieves and displays the name of the current profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	log.Info("retrieving current profile")

	p, err := h.profiles.GetActive(ctx)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	fmt.Fprintf(opts.Stdout, "Current profile: %s\n", p.Name)
	return nil
}
