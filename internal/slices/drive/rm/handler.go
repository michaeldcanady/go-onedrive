package rm

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Handler executes the drive rm operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive rm Handler.
func NewHandler(fs registry.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:  fs,
		log: l,
	}
}

// Handle removes an item from the filesystem.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("removing item", logger.String("path", opts.Path))

	provider, path, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}

	if err := provider.Remove(ctx, path); err != nil {
		return fmt.Errorf("failed to remove item at %s: %w", path, err)
	}

	fmt.Fprintf(opts.Stdout, "Item removed: %s\n", opts.Path)
	return nil
}
