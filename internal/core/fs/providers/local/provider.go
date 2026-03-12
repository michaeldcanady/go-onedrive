package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
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

func (_ *Provider) Name() string {
	return providerName
}

// Get retrieves metadata for a single item by its local path.
func (p *Provider) Get(ctx context.Context, path string) (shared.Item, error) {
	p.log.Debug("local.Get", logger.String("path", path))

	info, err := os.Stat(path)
	if err != nil {
		return shared.Item{}, err
	}
	return p.mapInfoToItem(path, info), nil
}

// List enumerates the contents of a directory on the local filesystem.
func (p *Provider) List(ctx context.Context, path string, opts shared.ListOptions) ([]shared.Item, error) {
	p.log.Debug("local.List", logger.String("path", path), logger.Bool("recursive", opts.Recursive))

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
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

	return os.Open(path)
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
		return shared.Item{}, err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return shared.Item{}, err
	}

	return p.Get(ctx, path)
}

// Mkdir creates a new folder on the local filesystem at the given path.
func (p *Provider) Mkdir(ctx context.Context, path string) error {
	p.log.Debug("local.Mkdir", logger.String("path", path))

	return os.MkdirAll(path, 0755)
}

// Remove deletes an item from the local filesystem.
func (p *Provider) Remove(ctx context.Context, path string) error {
	p.log.Debug("local.Remove", logger.String("path", path))

	return os.RemoveAll(path)
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

	return os.Rename(src, dst)
}

// Touch creates an empty file or updates the timestamp of an existing one.
func (p *Provider) Touch(ctx context.Context, path string) (shared.Item, error) {
	p.log.Debug("local.Touch", logger.String("path", path))

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return shared.Item{}, err
	}
	f.Close()

	now := time.Now()
	if err := os.Chtimes(path, now, now); err != nil {
		return shared.Item{}, err
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
