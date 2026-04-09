package touch

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive touch operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive touch Handler.
func NewHandler(
	fs registry.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-touch")
	return &Handler{
		fs:  fs,
		log: cliLog,
	}
}

// Handle creates an empty file or updates the timestamp of an existing one.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("touching item")

	log.Debug("requesting touch from provider")
	if _, err := h.fs.Touch(ctx, opts.Path); err != nil {
		h.log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	log.Info("item touched successfully")
	fmt.Fprintf(opts.Stdout, "File touched: %s\n", opts.Path)
	return nil
}
