package touch

import (
	"context"
	"fmt"

	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
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
	log := h.log.WithContext(ctx).With(logger.String("path", opts.URI.String()))

	log.Info("touching item")

	log.Debug("requesting touch from provider")
	if _, err := h.fs.Touch(ctx, opts.URI); err != nil {
		log.Error("touch failed", logger.Error(err))
		return fmt.Errorf("failed to touch file at %s: %w", opts.URI.String(), err)
	}

	log.Info("item touched successfully")
	fmt.Fprintf(opts.Stdout, "File touched: %s\n", opts.URI.String())
	return nil
}
