package mkdir

import (
	"context"
	"fmt"

	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
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
	log := h.log.WithContext(ctx).With(logger.String("path", opts.URI.String()))

	log.Info("creating directory")

	log.Debug("requesting directory creation from provider")
	if err := h.fs.Mkdir(ctx, opts.URI.String()); err != nil {
		log.Error("directory creation failed", logger.Error(err))
		return fmt.Errorf("failed to create directory at %s: %w", opts.URI.String(), err)
	}

	log.Info("directory created successfully")
	fmt.Fprintf(opts.Stdout, "Directory created: %s\n", opts.URI.String())
	return nil
}
