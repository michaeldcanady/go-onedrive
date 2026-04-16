package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/fs/providers"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger"
)

func init() {
	providers.Register(providers.Descriptor{
		Name: providerName,
		Factory: func(deps providers.Dependencies) (fs.Service, error) {
			return NewProvider(deps.Logger()), nil
		},
	})
}

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

	kind := fs.ErrInternal
	if os.IsNotExist(err) {
		kind = fs.ErrNotFound
	} else if os.IsPermission(err) {
		kind = fs.ErrForbidden
	} else if os.IsExist(err) {
		kind = fs.ErrConflict
	}

	return &fs.Error{
		Kind: kind,
		Err:  err,
		Path: path,
	}
}

// Get retrieves metadata for a single item by its local path.
func (p *Provider) Get(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	p.log.Debug("local.Get", logger.String("path", uri.Path))

	info, err := os.Stat(uri.Path)
	if err != nil {
		return fs.Item{}, p.mapError(err, uri.Path)
	}
	return p.mapInfoToItem(uri.Path, info), nil
}

// List enumerates the contents of a directory on the local filesystem.
func (p *Provider) List(ctx context.Context, uri *fs.URI, opts fs.ListOptions) ([]fs.Item, error) {
	p.log.Debug("local.List", logger.String("path", uri.Path), logger.Bool("recursive", opts.Recursive))

	entries, err := os.ReadDir(uri.Path)
	if err != nil {
		return nil, p.mapError(err, uri.Path)
	}

	var items []fs.Item
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		itemPath := filepath.Join(uri.Path, entry.Name())
		item := p.mapInfoToItem(itemPath, info)
		items = append(items, item)

		if opts.Recursive && entry.IsDir() {
			childURI := *uri // shallow copy
			childURI.Path = itemPath
			children, err := p.List(ctx, &childURI, opts)
			if err == nil {
				items = append(items, children...)
			}
		}
	}
	return items, nil
}

// ReadFile opens a read stream for a file's content on the local filesystem.
func (p *Provider) ReadFile(ctx context.Context, uri *fs.URI, opts fs.ReadOptions) (io.ReadCloser, error) {
	p.log.Debug("local.ReadFile", logger.String("path", uri.Path))

	f, err := os.Open(uri.Path)
	if err != nil {
		return nil, p.mapError(err, uri.Path)
	}
	return f, nil
}

// Stat returns metadata for a local file or directory.
func (p *Provider) Stat(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	return p.Get(ctx, uri)
}

// WriteFile creates or updates a file on the local filesystem with the content from the reader.
func (p *Provider) WriteFile(ctx context.Context, uri *fs.URI, r io.Reader, opts fs.WriteOptions) (fs.Item, error) {
	p.log.Debug("local.WriteFile", logger.String("path", uri.Path), logger.Bool("overwrite", opts.Overwrite))

	flags := os.O_WRONLY | os.O_CREATE
	if opts.Overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	f, err := os.OpenFile(uri.Path, flags, 0644)
	if err != nil {
		return fs.Item{}, p.mapError(err, uri.Path)
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return fs.Item{}, p.mapError(err, uri.Path)
	}

	return p.Get(ctx, uri)
}

// Mkdir creates a new folder on the local filesystem at the given path.
func (p *Provider) Mkdir(ctx context.Context, uri *fs.URI) error {
	p.log.Debug("local.Mkdir", logger.String("path", uri.Path))

	err := os.MkdirAll(uri.Path, 0755)
	return p.mapError(err, uri.Path)
}

// Remove deletes an item from the local filesystem.
func (p *Provider) Remove(ctx context.Context, uri *fs.URI) error {
	p.log.Debug("local.Remove", logger.String("path", uri.Path))

	err := os.RemoveAll(uri.Path)
	return p.mapError(err, uri.Path)
}

// Copy duplicates a file or folder on the local filesystem.
func (p *Provider) Copy(ctx context.Context, src, dst *fs.URI, opts fs.CopyOptions) error {
	p.log.Debug("local.Copy", logger.String("src", src.Path), logger.String("dst", dst.Path))

	r, err := p.ReadFile(ctx, src, fs.ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = p.WriteFile(ctx, dst, r, fs.WriteOptions{Overwrite: opts.Overwrite})
	return err
}

// Move relocates or renames a file or folder within the local filesystem.
func (p *Provider) Move(ctx context.Context, src, dst *fs.URI) error {
	p.log.Debug("local.Move", logger.String("src", src.Path), logger.String("dst", dst.Path))

	err := os.Rename(src.Path, dst.Path)
	return p.mapError(err, src.Path)
}

// Touch creates an empty file or updates the timestamp of an existing one.
func (p *Provider) Touch(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	p.log.Debug("local.Touch", logger.String("path", uri.Path))

	f, err := os.OpenFile(uri.Path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fs.Item{}, p.mapError(err, uri.Path)
	}
	f.Close()

	now := time.Now()
	if err := os.Chtimes(uri.Path, now, now); err != nil {
		return fs.Item{}, p.mapError(err, uri.Path)
	}

	return p.Get(ctx, uri)
}

func (p *Provider) mapInfoToItem(path string, info os.FileInfo) fs.Item {
	itemType := fs.TypeFile
	if info.IsDir() {
		itemType = fs.TypeFolder
	}

	return fs.Item{
		Path:       path,
		Name:       info.Name(),
		Size:       info.Size(),
		Type:       itemType,
		ModifiedAt: info.ModTime(),
	}
}
