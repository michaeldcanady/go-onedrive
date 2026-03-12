package mkdir

import (
	"context"
	"fmt"
	"path"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Handler executes the drive mkdir operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive mkdir Handler.
func NewHandler(fs registry.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:  fs,
		log: l,
	}
}

// Handle creates a new directory.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("creating directory")

	log.Debug("resolving path provider")
	provider, subPath, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		log.Error("failed to resolve path", logger.Error(err))
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}
	log.Debug("resolved path provider", logger.String("provider", provider.Name()))

	log.Debug("requesting directory creation from provider")
	if err := provider.Mkdir(ctx, subPath); err != nil {
		log.Error("directory creation failed", logger.Error(err))
		return fmt.Errorf("failed to create directory at %s: %w", path.Join(provider.Name(), subPath), err)
	}

	log.Info("directory created successfully")
	fmt.Fprintf(opts.Stdout, "Directory created: %s\n", opts.Path)
	return nil
}
