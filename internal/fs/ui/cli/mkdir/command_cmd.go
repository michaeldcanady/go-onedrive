package mkdir

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
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
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path.String()))

	log.Info("creating directory")

	log.Debug("requesting directory creation from provider")
	if err := h.fs.Mkdir(ctx, opts.Path); err != nil {
		wrapped := cli.WrapError(err, opts.Path.String())
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Info("directory created successfully")
	fmt.Fprintf(opts.Stdout, "Directory created: %s\n", opts.Path.String())
	return nil
}
