package create

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Handler executes the profile creation operation.
type Handler struct {
	profiles profile.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the profile create Handler.
func NewHandler(
	p profile.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("profile-create")
	return &Handler{
		profiles: p,
		log:      cliLog,
	}
}


// Handle executes the logic to create a new profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	log.Info("creating new profile", logger.String("name", opts.Name))

	p, err := h.profiles.Create(ctx, opts.Name)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	fmt.Fprintf(opts.Stdout, "Profile %s created successfully.\n", p.Name)
	return nil
}
