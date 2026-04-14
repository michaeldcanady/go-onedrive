package download

import (
	"context"
	"errors"
	"fmt"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive download operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive download Handler.
func NewHandler(m shared.Service, l logger.Logger) *Handler {
	return &Handler{
		manager: m,
		log:     l,
	}
}

// Handle downloads an item from the remote filesystem to the local destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.SourceURI.String()),
		logger.String("destination", opts.DestinationURI.String()),
	)

	log.Info("starting download")

	if opts.DestinationURI.Provider == "" {
		log.Debug("adding local prefix to destination")
		opts.DestinationURI.Provider = "local"
	} else if opts.DestinationURI.Provider != "local" {
		log.Warn("destination provider is not local")
		return errors.New("destination provider is not local")
	}

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := h.manager.Copy(ctx, opts.SourceURI.String(), opts.DestinationURI.String(), cpOpts); err != nil {
		log.Error("download failed", logger.Error(err))
		return fmt.Errorf("failed to download %s to %s: %w", opts.SourceURI.String(), opts.DestinationURI.String(), err)
	}

	log.Info("download completed successfully")
	fmt.Fprintf(opts.Stdout, "Downloaded %s to %s\n", opts.SourceURI.String(), opts.DestinationURI.String())
	return nil
}
