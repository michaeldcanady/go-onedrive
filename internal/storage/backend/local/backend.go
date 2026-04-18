package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger"
)

// Backend implements the fs.Backend and fs.AdvancedBackend interfaces for the local filesystem.
type Backend struct {
	root string
	log  logger.Logger
}

// NewBackend creates a new instance of the local filesystem backend.
// The root parameter defines the base directory for all operations.
func NewBackend(root string, log logger.Logger) *Backend {
	return &Backend{
		root: root,
		log:  log,
	}
}

func (b *Backend) Name() string {
	return "local"
}

// fullPath joins the backend root with the provided relative path.
func (b *Backend) fullPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return filepath.Join(b.root, path)
}

func (b *Backend) mapError(err error, path string) error {
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

func (b *Backend) Stat(ctx context.Context, path string) (fs.Item, error) {
	b.log.Debug("local.Stat", logger.String("path", path))
	info, err := os.Stat(b.fullPath(path))
	if err != nil {
		return fs.Item{}, b.mapError(err, path)
	}
	return b.mapInfoToItem(path, info), nil
}

func (b *Backend) List(ctx context.Context, path string) ([]fs.Item, error) {
	b.log.Debug("local.List", logger.String("path", path))
	entries, err := os.ReadDir(b.fullPath(path))
	if err != nil {
		return nil, b.mapError(err, path)
	}

	var items []fs.Item
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		items = append(items, b.mapInfoToItem(filepath.Join(path, entry.Name()), info))
	}
	return items, nil
}

func (b *Backend) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	b.log.Debug("local.Open", logger.String("path", path))
	f, err := os.Open(b.fullPath(path))
	if err != nil {
		return nil, b.mapError(err, path)
	}
	return f, nil
}

func (b *Backend) Create(ctx context.Context, path string, r io.Reader) (fs.Item, error) {
	b.log.Debug("local.Create", logger.String("path", path))
	f, err := os.Create(b.fullPath(path))
	if err != nil {
		return fs.Item{}, b.mapError(err, path)
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return fs.Item{}, b.mapError(err, path)
	}

	return b.Stat(ctx, path)
}

func (b *Backend) Mkdir(ctx context.Context, path string) error {
	b.log.Debug("local.Mkdir", logger.String("path", path))
	err := os.MkdirAll(b.fullPath(path), 0755)
	return b.mapError(err, path)
}

func (b *Backend) Remove(ctx context.Context, path string) error {
	b.log.Debug("local.Remove", logger.String("path", path))
	err := os.RemoveAll(b.fullPath(path))
	return b.mapError(err, path)
}

func (b *Backend) Capabilities() fs.Capabilities {
	return fs.Capabilities{
		CanMove:      true,
		CanCopy:      true,
		CanRecursive: false, // VFS handles recursion for local for now
	}
}

func (b *Backend) Move(ctx context.Context, src, dst string) error {
	b.log.Debug("local.Move", logger.String("src", src), logger.String("dst", dst))
	err := os.Rename(b.fullPath(src), b.fullPath(dst))
	return b.mapError(err, src)
}

func (b *Backend) Copy(ctx context.Context, src, dst string) error {
	b.log.Debug("local.Copy", logger.String("src", src), logger.String("dst", dst))
	// Simplified local copy via streaming for now, or could use os level if available
	r, err := b.Open(ctx, src)
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = b.Create(ctx, dst, r)
	return err
}

func (b *Backend) mapInfoToItem(path string, info os.FileInfo) fs.Item {
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
