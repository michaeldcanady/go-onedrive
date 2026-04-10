package rm

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive rm operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive rm Handler.
func NewHandler(
	fs registry.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-rm")
	return &Handler{
		fs:  fs,
		log: cliLog,
	}
}

// Handle removes an item from the filesystem.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("removing item")

	log.Debug("requesting removal from provider")
	if err := h.fs.Remove(ctx, opts.Path); err != nil {
		wrapped := cli.WrapError(err, opts.Path)
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Info("item removed successfully")
	fmt.Fprintf(opts.Stdout, "Item removed: %s\n", opts.Path)
	return nil
}
