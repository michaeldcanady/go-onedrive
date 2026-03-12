package edit

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/editor"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/core/errors"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
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
	h.log.Info("editing file", logger.String("path", opts.Path))

	provider, path, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}

	h.log.Debug("fetching metadata", logger.String("path", path))
	item, err := provider.Get(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to get item metadata at %s: %w", path, err)
	}

	h.log.Debug("reading file for editing", logger.String("path", path))
	r, err := provider.ReadFile(ctx, path, shared.ReadOptions{})
	if err != nil {
		return fmt.Errorf("failed to read file at %s: %w", path, err)
	}
	defer r.Close()

	h.log.Debug("launching editor")
	newData, _, err := h.editor.LaunchTempFile("odc-edit", ".txt", r)
	if err != nil {
		return fmt.Errorf("failed to edit file at %s: %w", path, err)
	}

	if newData == nil {
		h.log.Info("no changes detected, skipping upload")
		_, _ = fmt.Fprintln(opts.Stdout, "No changes detected.")
		return nil
	}

	h.log.Info("uploading changes", logger.String("path", path))
	writeOpts := shared.WriteOptions{
		Overwrite: opts.Force,
	}
	// Use ETag for optimistic concurrency control if not forcing
	if !opts.Force {
		writeOpts.IfMatch = item.ETag
	}

	_, err = provider.WriteFile(ctx, path, bytes.NewReader(newData), writeOpts)
	if err != nil {
		if errors.Is(err, coreerrors.ErrPrecondition) {
			return fmt.Errorf("conflict detected during upload: %w. Use --force to overwrite", err)
		}
		return fmt.Errorf("failed to upload edited changes to %s: %w", path, err)
	}

	_, _ = fmt.Fprintf(opts.Stdout, "successfully updated file \"%s\"\n", opts.Path)
	return nil
}
