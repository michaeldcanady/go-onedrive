package set

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive alias set operation.
type Handler struct {
	alias alias.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive alias set Handler.
func NewHandler(alias alias.Service, l logger.Logger) *Handler {
	return &Handler{
		alias: alias,
		log:   l,
	}
}

// Handle creates or updates a drive alias.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("setting drive alias", logger.String("alias", opts.Alias), logger.String("driveID", opts.DriveID))

	if err := h.alias.SetAlias(opts.Alias, opts.DriveID); err != nil {
		return fmt.Errorf("failed to set drive alias: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "alias '%s' set to drive '%s'\n", opts.Alias, opts.DriveID)
	return nil
}
