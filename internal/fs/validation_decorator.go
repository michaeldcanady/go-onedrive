package fs

import (
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// ValidationDecorator wraps a filesystem service and performs pre-operation validation on all paths.
type ValidationDecorator struct {
	next   Service
	logger logger.Logger
}

// NewValidationDecorator initializes a new instance of the ValidationDecorator.
func NewValidationDecorator(next Service, logger logger.Logger) *ValidationDecorator {
	return &ValidationDecorator{
		next:   next,
		logger: logger,
	}
}

func (vd *ValidationDecorator) Name() string {
	return vd.next.Name()
}

// Get validates the URI before retrieving metadata.
func (vd *ValidationDecorator) Get(ctx context.Context, uri *URI) (Item, error) {
	if err := vd.validateURI(uri, "Get"); err != nil {
		return Item{}, err
	}
	return vd.next.Get(ctx, uri)
}

// Stat validates the URI before retrieving metadata.
func (vd *ValidationDecorator) Stat(ctx context.Context, uri *URI) (Item, error) {
	if err := vd.validateURI(uri, "Stat"); err != nil {
		return Item{}, err
	}
	return vd.next.Stat(ctx, uri)
}

// List validates the URI before retrieving children.
func (vd *ValidationDecorator) List(ctx context.Context, uri *URI, opts ListOptions) ([]Item, error) {
	if err := vd.validateURI(uri, "List"); err != nil {
		return nil, err
	}
	return vd.next.List(ctx, uri, opts)
}

// ReadFile validates the URI before opening a read stream.
func (vd *ValidationDecorator) ReadFile(ctx context.Context, uri *URI, opts ReadOptions) (io.ReadCloser, error) {
	if err := vd.validateURI(uri, "ReadFile"); err != nil {
		return nil, err
	}
	return vd.next.ReadFile(ctx, uri, opts)
}

// WriteFile validates the URI before uploading content.
func (vd *ValidationDecorator) WriteFile(ctx context.Context, uri *URI, r io.Reader, opts WriteOptions) (Item, error) {
	if err := vd.validateURI(uri, "WriteFile"); err != nil {
		return Item{}, err
	}
	return vd.next.WriteFile(ctx, uri, r, opts)
}

// Mkdir validates the URI before creating a directory.
func (vd *ValidationDecorator) Mkdir(ctx context.Context, uri *URI) error {
	if err := vd.validateURI(uri, "Mkdir"); err != nil {
		return err
	}
	return vd.next.Mkdir(ctx, uri)
}

// Remove validates the URI before deletion.
func (vd *ValidationDecorator) Remove(ctx context.Context, uri *URI) error {
	if err := vd.validateURI(uri, "Remove"); err != nil {
		return err
	}
	return vd.next.Remove(ctx, uri)
}

// Touch validates the URI before creation or update.
func (vd *ValidationDecorator) Touch(ctx context.Context, uri *URI) (Item, error) {
	if err := vd.validateURI(uri, "Touch"); err != nil {
		return Item{}, err
	}
	return vd.next.Touch(ctx, uri)
}

// Copy validates both source and destination URIs.
func (vd *ValidationDecorator) Copy(ctx context.Context, src, dst *URI, opts CopyOptions) error {
	if err := vd.validateURI(src, "Copy.Source"); err != nil {
		return err
	}
	if err := vd.validateURI(dst, "Copy.Destination"); err != nil {
		return err
	}
	return vd.next.Copy(ctx, src, dst, opts)
}

// Move validates both source and destination URIs.
func (vd *ValidationDecorator) Move(ctx context.Context, src, dst *URI) error {
	if err := vd.validateURI(src, "Move.Source"); err != nil {
		return err
	}
	if err := vd.validateURI(dst, "Move.Destination"); err != nil {
		return err
	}
	return vd.next.Move(ctx, src, dst)
}

// validateURI checks for common URI/path issues.
func (vd *ValidationDecorator) validateURI(uri *URI, operation string) error {
	if err := ValidatePathSyntax(uri.Path); err != nil {
		vd.logger.Error("URI validation failed",
			logger.String("path", uri.Path),
			logger.String("operation", operation),
			logger.Error(err))
		return err
	}

	vd.logger.Debug("URI validation successful",
		logger.String("path", uri.Path),
		logger.String("operation", operation))
	return nil
}
