package delete

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Handler executes the profile deletion operation.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile deletion Handler.
func NewHandler(
	p profile.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("profile-delete")
	return &Handler{
		profiles: p,
		log:      cliLog,
	}
}

// Handle executes the logic to delete a profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	log.Info("deleting profile", logger.String("name", opts.Name))

	if err := h.profiles.Delete(ctx, opts.Name); err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	fmt.Fprintf(opts.Stdout, "Profile %s deleted successfully.\n", opts.Name)
	return nil
}
