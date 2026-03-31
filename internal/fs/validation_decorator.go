package fs

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// ValidationDecorator is a middleware that validates requests before they reach the underlying service.
type ValidationDecorator struct {
	next   Service
	logger logger.Logger
}

// NewValidationDecorator creates a new ValidationDecorator.
func NewValidationDecorator(next Service, l logger.Logger) Service {
	return &ValidationDecorator{
		next:   next,
		logger: l,
	}
}

func (vd *ValidationDecorator) Name() string {
	return vd.next.Name()
}

// Get validates the path and then calls the next service's Get method.
func (vd *ValidationDecorator) Get(ctx context.Context, path string) (Item, error) {
	if err := vd.validatePath(path, "get"); err != nil {
		return Item{}, err
	}
	return vd.next.Get(ctx, path)
}

// List validates the path and then calls the next service's List method.
func (vd *ValidationDecorator) List(ctx context.Context, path string, opts ListOptions) ([]Item, error) {
	if err := vd.validatePath(path, "list"); err != nil {
		return nil, err
	}
	return vd.next.List(ctx, path, opts)
}

// ReadFile validates the path and then calls the next service's ReadFile method.
func (vd *ValidationDecorator) ReadFile(ctx context.Context, path string, opts ReadOptions) (io.ReadCloser, error) {
	if err := vd.validatePath(path, "read"); err != nil {
		return nil, err
	}
	return vd.next.ReadFile(ctx, path, opts)
}

// Stat validates the path and then calls the next service's Stat method.
func (vd *ValidationDecorator) Stat(ctx context.Context, path string) (Item, error) {
	if err := vd.validatePath(path, "stat"); err != nil {
		return Item{}, err
	}
	return vd.next.Stat(ctx, path)
}

// WriteFile validates the path and then calls the next service's WriteFile method.
func (vd *ValidationDecorator) WriteFile(ctx context.Context, path string, r io.Reader, opts WriteOptions) (Item, error) {
	if err := vd.validatePath(path, "write"); err != nil {
		return Item{}, err
	}
	return vd.next.WriteFile(ctx, path, r, opts)
}

// Mkdir validates the path and then calls the next service's Mkdir method.
func (vd *ValidationDecorator) Mkdir(ctx context.Context, path string) error {
	if err := vd.validatePath(path, "mkdir"); err != nil {
		return err
	}
	return vd.next.Mkdir(ctx, path)
}

// Remove validates the path and then calls the next service's Remove method.
func (vd *ValidationDecorator) Remove(ctx context.Context, path string) error {
	if err := vd.validatePath(path, "remove"); err != nil {
		return err
	}
	return vd.next.Remove(ctx, path)
}

// Touch validates the path and then calls the next service's Touch method.
func (vd *ValidationDecorator) Touch(ctx context.Context, path string) (Item, error) {
	if err := vd.validatePath(path, "touch"); err != nil {
		return Item{}, err
	}
	return vd.next.Touch(ctx, path)
}

// Copy validates source and destination paths and then calls the next service's Copy method.
func (vd *ValidationDecorator) Copy(ctx context.Context, src, dst string, opts CopyOptions) error {
	if err := vd.validatePath(src, "copy source"); err != nil {
		return err
	}
	if err := vd.validatePath(dst, "copy destination"); err != nil {
		return err
	}
	return vd.next.Copy(ctx, src, dst, opts)
}

// Move validates source and destination paths and then calls the next service's Move method.
func (vd *ValidationDecorator) Move(ctx context.Context, src, dst string) error {
	if err := vd.validatePath(src, "move source"); err != nil {
		return err
	}
	if err := vd.validatePath(dst, "move destination"); err != nil {
		return err
	}
	return vd.next.Move(ctx, src, dst)
}

// validatePath checks for common path issues like trailing slashes or illegal characters.
func (vd *ValidationDecorator) validatePath(p string, operation string) error {
	// Disallow trailing slashes unless it's the root path "/"
	if strings.HasSuffix(p, "/") && p != "/" {
		return fmt.Errorf("path '%s' has a trailing slash, which is not allowed for operation '%s'", p, operation)
	}

	// Disallow illegal characters (e.g., ':', '#', '?', '*', '[', ']', '\') - Windows illegal chars
	// Note: This list might need to be expanded based on specific provider limitations.
	illegalChars := []string{":", "#", "?", "*", "[", "]", "\\"}
	for _, char := range illegalChars {
		if strings.Contains(p, char) {
			return fmt.Errorf("path '%s' contains illegal character '%s' for operation '%s'", p, char, operation)
		}
	}

	vd.logger.Debug("path validation successful", logger.String("path", p), logger.String("operation", operation))
	return nil
}
