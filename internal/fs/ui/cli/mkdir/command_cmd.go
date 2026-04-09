package mkdir

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive mkdir operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive mkdir Handler.
func NewHandler(
	fs registry.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-mkdir")
	return &Handler{
		fs:  fs,
		log: cliLog,
	}
}

// Handle creates a new directory.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("creating directory")

	log.Debug("requesting directory creation from provider")
	if err := h.fs.Mkdir(ctx, opts.Path); err != nil {
		h.log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	log.Info("directory created successfully")
	fmt.Fprintf(opts.Stdout, "Directory created: %s\n", opts.Path)
	return nil
}
