package remove

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive alias remove operation.
type Handler struct {
	alias alias.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive alias remove Handler.
func NewHandler(alias alias.Service, l logger.Logger) *Handler {
	return &Handler{
		alias: alias,
		log:   l,
	}
}

// Handle deletes a drive alias.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("removing drive alias", logger.String("alias", opts.Alias))

	if err := h.alias.DeleteAlias(opts.Alias); err != nil {
		return fmt.Errorf("failed to remove drive alias: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "alias '%s' removed\n", opts.Alias)
	return nil
}
