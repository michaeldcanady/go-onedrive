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
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("removing item")

	log.Debug("resolving path provider")
	provider, subPath, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		log.Error("failed to resolve path", logger.Error(err))
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}
	log.Debug("resolved path provider", logger.String("provider", provider.Name()))

	log.Debug("requesting removal from provider")
	if err := provider.Remove(ctx, subPath); err != nil {
		log.Error("removal failed", logger.Error(err))
		return fmt.Errorf("failed to remove item at %s: %w", opts.Path, err)
	}

	log.Info("item removed successfully")
	fmt.Fprintf(opts.Stdout, "Item removed: %s\n", opts.Path)
	return nil
}
