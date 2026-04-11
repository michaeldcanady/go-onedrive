package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

const (
	providerName = "local"
)

// Provider implements the filesystem Service interface for the local filesystem.
type Provider struct {
	// log is the logger instance used for recording provider events.
	log logger.Logger
}

// NewProvider creates a new instance of the local filesystem provider.
func NewProvider(log logger.Logger) *Provider {
	return &Provider{
		log: log,
	}
}

func (p *Provider) Name() string {
	return providerName
}

// mapReadError translates errors from the local filesystem into a ReadError wrapping domain-specific errors.
func (p *Provider) mapReadError(err error, path string) error {
	if err == nil {
		return nil
	}

	wrapped := p.mapToDomainError(err, path)
	return NewReadError(path, wrapped)
}

// mapWriteError translates errors from the local filesystem into a WriteError wrapping domain-specific errors.
func (p *Provider) mapWriteError(err error, path string) error {
	if err == nil {
		return nil
	}

	wrapped := p.mapToDomainError(err, path)
	return NewWriteError(path, wrapped)
}

// mapGenericError translates errors from the local filesystem into domain-specific errors wrapped in an AppError.
func (p *Provider) mapGenericError(err error, path string) error {
	if err == nil {
		return nil
	}

	domainErr := p.mapToDomainError(err, path)

	// If mapToDomainError already returned one of our custom error types, return it directly.
	if domainErr != err {
		return domainErr
	}

	code := coreerrors.CodeInternal
	safeMsg := "an unexpected local filesystem error occurred"
	hint := ""

	appErr := coreerrors.NewAppError(code, domainErr, safeMsg, hint)
	if path != "" {
		appErr.WithContext(coreerrors.KeyPath, path)
	}

	return appErr
}

// mapToDomainError converts low-level OS errors into domain-specific error types.
func (p *Provider) mapToDomainError(err error, path string) error {
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return coreerrors.NewNotFoundError(path, err)
	} else if os.IsPermission(err) {
		return coreerrors.NewForbiddenError(path, err)
	} else if os.IsExist(err) {
		return coreerrors.NewConflictError(path, err)
	}

	return err
}

// Get retrieves metadata for a single item by its local path.
func (p *Provider) Get(ctx context.Context, uri *shared.URI) (shared.Item, error) {
	p.log.Debug("local.Get", logger.String("path", uri.Path))

	info, err := os.Stat(uri.Path)
	if err != nil {
		return shared.Item{}, p.mapGenericError(err, uri.Path)
	}
	return p.mapInfoToItem(uri.Path, info), nil
}

// List enumerates the contents of a directory on the local filesystem.
func (p *Provider) List(ctx context.Context, uri *shared.URI, opts shared.ListOptions) ([]shared.Item, error) {
	p.log.Debug("local.List", logger.String("path", uri.Path), logger.Bool("recursive", opts.Recursive))

	entries, err := os.ReadDir(uri.Path)
	if err != nil {
		return nil, p.mapGenericError(err, uri.Path)
	}

	var items []shared.Item
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		itemPath := filepath.Join(uri.Path, entry.Name())
		item := p.mapInfoToItem(itemPath, info)
		items = append(items, item)

		if opts.Recursive && entry.IsDir() {
			childURI := &shared.URI{
				Provider: uri.Provider,
				DriveRef: uri.DriveRef,
				Path:     itemPath,
			}
			children, err := p.List(ctx, childURI, opts)
			if err == nil {
				items = append(items, children...)
			}
		}
	}
	return items, nil
}

// ReadFile opens a read stream for a file's content on the local filesystem.
func (p *Provider) ReadFile(ctx context.Context, uri *shared.URI, opts shared.ReadOptions) (io.ReadCloser, error) {
	p.log.Debug("local.ReadFile", logger.String("path", uri.Path))

	f, err := os.Open(uri.Path)
	if err != nil {
		return nil, p.mapReadError(err, uri.Path)
	}
	return f, nil
}

// Stat returns metadata for a local file or directory.
func (p *Provider) Stat(ctx context.Context, uri *shared.URI) (shared.Item, error) {
	return p.Get(ctx, uri)
}

// WriteFile creates or updates a file on the local filesystem with the content from the reader.
func (p *Provider) WriteFile(ctx context.Context, uri *shared.URI, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	p.log.Debug("local.WriteFile", logger.String("path", uri.Path), logger.Bool("overwrite", opts.Overwrite))

	flags := os.O_WRONLY | os.O_CREATE
	if opts.Overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	f, err := os.OpenFile(uri.Path, flags, 0644)
	if err != nil {
		return shared.Item{}, p.mapWriteError(err, uri.Path)
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return shared.Item{}, p.mapWriteError(err, uri.Path)
	}

	return p.Get(ctx, uri)
}

// Mkdir creates a new folder on the local filesystem at the given path.
func (p *Provider) Mkdir(ctx context.Context, uri *shared.URI) error {
	p.log.Debug("local.Mkdir", logger.String("path", uri.Path))

	err := os.MkdirAll(uri.Path, 0755)
	if err != nil {
		return p.mapGenericError(err, uri.Path)
	}
	return nil
}

// Remove deletes an item from the local filesystem.
func (p *Provider) Remove(ctx context.Context, uri *shared.URI) error {
	p.log.Debug("local.Remove", logger.String("path", uri.Path))

	err := os.RemoveAll(uri.Path)
	if err != nil {
		return p.mapGenericError(err, uri.Path)
	}
	return nil
}

// Copy duplicates a file or folder on the local filesystem.
func (p *Provider) Copy(ctx context.Context, src, dst *shared.URI, opts shared.CopyOptions) error {
	p.log.Debug("local.Copy", logger.String("src", src.Path), logger.String("dst", dst.Path))

	r, err := p.ReadFile(ctx, src, shared.ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = p.WriteFile(ctx, dst, r, shared.WriteOptions{Overwrite: opts.Overwrite})
	return err
}

// Move relocates or renames a file or folder within OneDrive.
func (p *Provider) Move(ctx context.Context, src, dst *shared.URI) error {
	p.log.Debug("local.Move", logger.String("src", src.Path), logger.String("dst", dst.Path))

	err := os.Rename(src.Path, dst.Path)
	if err != nil {
		return p.mapGenericError(err, src.Path)
	}
	return nil
}

// Touch creates an empty file or updates the timestamp of an existing one.
func (p *Provider) Touch(ctx context.Context, uri *shared.URI) (shared.Item, error) {
	p.log.Debug("local.Touch", logger.String("path", uri.Path))

	f, err := os.OpenFile(uri.Path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return shared.Item{}, p.mapWriteError(err, uri.Path)
	}
	f.Close()

	now := time.Now()
	if err := os.Chtimes(uri.Path, now, now); err != nil {
		return shared.Item{}, p.mapWriteError(err, uri.Path)
	}

	return p.Get(ctx, uri)
}

func (p *Provider) mapInfoToItem(path string, info os.FileInfo) shared.Item {
	itemType := shared.TypeFile
	if info.IsDir() {
		itemType = shared.TypeFolder
	}

	return shared.Item{
		Path:       path,
		Name:       info.Name(),
		Size:       info.Size(),
		Type:       itemType,
		ModifiedAt: info.ModTime(),
	}
}
