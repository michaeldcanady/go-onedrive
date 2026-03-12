package touch

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Handler executes the drive touch operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive touch Handler.
func NewHandler(fs registry.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:  fs,
		log: l,
	}
}

// Handle creates an empty file or updates the timestamp of an existing one.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("touching item")

	log.Debug("resolving path provider")
	provider, subPath, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		log.Error("failed to resolve path", logger.Error(err))
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}
	log.Debug("resolved path provider", logger.String("provider", provider.Name()))

	log.Debug("requesting touch from provider")
	if _, err := provider.Touch(ctx, subPath); err != nil {
		log.Error("touch failed", logger.Error(err))
		return fmt.Errorf("failed to touch file at %s: %w", opts.Path, err)
	}

	log.Info("item touched successfully")
	fmt.Fprintf(opts.Stdout, "File touched: %s\n", opts.Path)
	return nil
}
