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

func (p *Provider) mapError(err error, path string) error {
	if err == nil {
		return nil
	}

	var appErr *coreerrors.AppError
	if os.IsNotExist(err) {
		appErr = coreerrors.NewNotFound(err, "file or directory not found", "")
	} else if os.IsPermission(err) {
		appErr = coreerrors.NewForbidden(err, "permission denied", "Check your file system permissions.")
	} else if os.IsExist(err) {
		appErr = coreerrors.NewConflict(err, "file or directory already exists", "")
	} else {
		appErr = coreerrors.NewInternal(err, "an unexpected local filesystem error occurred", "")
	}

	return appErr.WithContext(coreerrors.KeyPath, path)
}

// Get retrieves metadata for a single item by its local path.
func (p *Provider) Get(ctx context.Context, path string) (shared.Item, error) {
	p.log.Debug("local.Get", logger.String("path", path))

	info, err := os.Stat(path)
	if err != nil {
		return shared.Item{}, p.mapError(err, path)
	}
	return p.mapInfoToItem(path, info), nil
}

// List enumerates the contents of a directory on the local filesystem.
func (p *Provider) List(ctx context.Context, path string, opts shared.ListOptions) ([]shared.Item, error) {
	p.log.Debug("local.List", logger.String("path", path), logger.Bool("recursive", opts.Recursive))

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, p.mapError(err, path)
	}

	var items []shared.Item
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		item := p.mapInfoToItem(filepath.Join(path, entry.Name()), info)
		items = append(items, item)

		if opts.Recursive && entry.IsDir() {
			children, err := p.List(ctx, filepath.Join(path, entry.Name()), opts)
			if err == nil {
				items = append(items, children...)
			}
		}
	}
	return items, nil
}

// ReadFile opens a read stream for a file's content on the local filesystem.
func (p *Provider) ReadFile(ctx context.Context, path string, opts shared.ReadOptions) (io.ReadCloser, error) {
	p.log.Debug("local.ReadFile", logger.String("path", path))

	f, err := os.Open(path)
	if err != nil {
		return nil, p.mapError(err, path)
	}
	return f, nil
}

// Stat returns metadata for a local file or directory.
func (p *Provider) Stat(ctx context.Context, path string) (shared.Item, error) {
	return p.Get(ctx, path)
}

// WriteFile creates or updates a file on the local filesystem with the content from the reader.
func (p *Provider) WriteFile(ctx context.Context, path string, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	p.log.Debug("local.WriteFile", logger.String("path", path), logger.Bool("overwrite", opts.Overwrite))

	flags := os.O_WRONLY | os.O_CREATE
	if opts.Overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	f, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return shared.Item{}, p.mapError(err, path)
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return shared.Item{}, p.mapError(err, path)
	}

	return p.Get(ctx, path)
}

// Mkdir creates a new folder on the local filesystem at the given path.
func (p *Provider) Mkdir(ctx context.Context, path string) error {
	p.log.Debug("local.Mkdir", logger.String("path", path))

	err := os.MkdirAll(path, 0755)
	return p.mapError(err, path)
}

// Remove deletes an item from the local filesystem.
func (p *Provider) Remove(ctx context.Context, path string) error {
	p.log.Debug("local.Remove", logger.String("path", path))

	err := os.RemoveAll(path)
	return p.mapError(err, path)
}

// Copy duplicates a file or folder on the local filesystem.
func (p *Provider) Copy(ctx context.Context, src, dst string, opts shared.CopyOptions) error {
	p.log.Debug("local.Copy", logger.String("src", src), logger.String("dst", dst))

	r, err := p.ReadFile(ctx, src, shared.ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = p.WriteFile(ctx, dst, r, shared.WriteOptions{Overwrite: opts.Overwrite})
	return err
}

// Move relocates or renames a file or folder within OneDrive.
func (p *Provider) Move(ctx context.Context, src, dst string) error {
	p.log.Debug("local.Move", logger.String("src", src), logger.String("dst", dst))

	err := os.Rename(src, dst)
	return p.mapError(err, src)
}

// Touch creates an empty file or updates the timestamp of an existing one.
func (p *Provider) Touch(ctx context.Context, path string) (shared.Item, error) {
	p.log.Debug("local.Touch", logger.String("path", path))

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return shared.Item{}, p.mapError(err, path)
	}
	f.Close()

	now := time.Now()
	if err := os.Chtimes(path, now, now); err != nil {
		return shared.Item{}, p.mapError(err, path)
	}

	return p.Get(ctx, path)
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
