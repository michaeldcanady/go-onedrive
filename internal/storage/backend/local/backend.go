package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Backend implements the fs.Backend and fs.AdvancedBackend interfaces for the local filesystem.
type Backend struct {
	root string
}

// NewBackend creates a new instance of the local filesystem backend.
func NewBackend(root string) *Backend {
	return &Backend{
		root: root,
	}
}

func (b *Backend) Name() string {
	return "local"
}

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

func (b *Backend) Stat(ctx context.Context, token, driveID, path string) (fs.Item, error) {
	info, err := os.Stat(b.fullPath(path))
	if err != nil {
		return fs.Item{}, b.mapError(err, path)
	}
	return b.mapInfoToItem(path, info), nil
}

func (b *Backend) List(ctx context.Context, token, driveID, path string) ([]fs.Item, error) {
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

func (b *Backend) Open(ctx context.Context, token, driveID, path string) (io.ReadCloser, error) {
	f, err := os.Open(b.fullPath(path))
	if err != nil {
		return nil, b.mapError(err, path)
	}
	return f, nil
}

func (b *Backend) Create(ctx context.Context, token, driveID, path string, r io.Reader) (fs.Item, error) {
	f, err := os.Create(b.fullPath(path))
	if err != nil {
		return fs.Item{}, b.mapError(err, path)
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return fs.Item{}, b.mapError(err, path)
	}

	return b.Stat(ctx, token, driveID, path)
}

func (b *Backend) Mkdir(ctx context.Context, token, driveID, path string) error {
	err := os.MkdirAll(b.fullPath(path), 0755)
	return b.mapError(err, path)
}

func (b *Backend) Remove(ctx context.Context, token, driveID, path string) error {
	err := os.RemoveAll(b.fullPath(path))
	return b.mapError(err, path)
}

func (b *Backend) Capabilities() fs.Capabilities {
	return fs.Capabilities{
		CanMove:      true,
		CanCopy:      true,
		CanRecursive: false,
	}
}

func (b *Backend) Move(ctx context.Context, token, driveID, src, dst string) error {
	err := os.Rename(b.fullPath(src), b.fullPath(dst))
	return b.mapError(err, src)
}

func (b *Backend) Copy(ctx context.Context, token, driveID, src, dst string) error {
	r, err := b.Open(ctx, token, driveID, src)
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = b.Create(ctx, token, driveID, dst, r)
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
