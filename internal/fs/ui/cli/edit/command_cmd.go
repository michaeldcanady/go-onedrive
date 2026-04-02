package edit

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	featerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/editor"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive edit operation.
type Handler struct {
	fs     registry.Service
	editor editor.Service
	log    logger.Logger
}

// NewHandler initializes a new instance of the drive edit Handler.
func NewHandler(fs registry.Service, editor editor.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:     fs,
		editor: editor,
		log:    l,
	}
}

// Handle retrieves a file, opens it in an editor, and uploads changes.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("editing file")

	log.Debug("fetching metadata")
	item, err := h.fs.Get(ctx, opts.Path)
	if err != nil {
		log.Error("failed to get metadata", logger.Error(err))
		return fmt.Errorf("failed to get item metadata at %s: %w", opts.Path, err)
	}

	log.Debug("reading file for editing")
	r, err := h.fs.ReadFile(ctx, opts.Path, registry.ReadOptions{})
	if err != nil {
		log.Error("failed to read file", logger.Error(err))
		return fmt.Errorf("failed to read file at %s: %w", opts.Path, err)
	}
	defer r.Close()

	log.Debug("launching external editor")
	newData, _, err := h.editor.LaunchTempFile("odc-edit", ".txt", r)
	if err != nil {
		log.Error("editor launch failed", logger.Error(err))
		return fmt.Errorf("failed to edit file at %s: %w", opts.Path, err)
	}

	if newData == nil {
		log.Info("no changes detected, skipping upload")
		_, _ = fmt.Fprintln(opts.Stdout, "No changes detected.")
		return nil
	}

	log.Info("uploading changes", logger.Int("size", len(newData)))
	writeOpts := registry.WriteOptions{
		Overwrite: opts.Force,
	}
	// Use ETag for optimistic concurrency control if not forcing
	if !opts.Force {
		writeOpts.IfMatch = item.ETag
	}

	_, err = h.fs.WriteFile(ctx, opts.Path, bytes.NewReader(newData), writeOpts)
	if err != nil {
		if errors.Is(err, featerrors.ErrPrecondition) {
			log.Warn("conflict detected, upload aborted")
			return fmt.Errorf("conflict detected during upload: %w. Use --force to overwrite", err)
		}
		log.Error("failed to upload changes", logger.Error(err))
		return fmt.Errorf("failed to upload edited changes to %s: %w", opts.Path, err)
	}

	log.Info("edit completed successfully")
	_, _ = fmt.Fprintf(opts.Stdout, "successfully updated file \"%s\"\n", opts.Path)
	return nil
}
