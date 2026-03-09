package local

import (
	"context"
	"io"
	"os"
	"path/filepath"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

var _ domainfs.Service = (*Provider)(nil)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Get(ctx context.Context, path string) (domainfs.Item, error) {
	info, err := os.Stat(path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.mapInfoToItem(path, info), nil
}

func (p *Provider) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []domainfs.Item
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

func (p *Provider) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	return p.Get(ctx, path)
}

func (p *Provider) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	return os.Open(path)
}

func (p *Provider) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	flag := os.O_WRONLY | os.O_CREATE
	if opts.Overwrite {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_EXCL
	}

	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return domainfs.Item{}, err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return domainfs.Item{}, err
	}

	return p.Get(ctx, path)
}

func (p *Provider) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	if opts.Parents {
		return os.MkdirAll(path, 0755)
	}
	return os.Mkdir(path, 0755)
}

func (p *Provider) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	return os.RemoveAll(path)
}

func (p *Provider) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	return os.Rename(src, dst)
}

func (p *Provider) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	f, err := os.Open(src)
	if err != nil {
		return domainfs.Item{}, err
	}
	defer f.Close()

	return p.WriteFile(ctx, dst, f, domainfs.WriteOptions{Overwrite: opts.Overwrite})
}

func (p *Provider) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return domainfs.Item{}, err
	}
	f.Close()
	return p.Get(ctx, path)
}

func (p *Provider) mapInfoToItem(path string, info os.FileInfo) domainfs.Item {
	itemType := domainfs.ItemTypeFile
	if info.IsDir() {
		itemType = domainfs.ItemTypeFolder
	}

	return domainfs.Item{
		Path:     path,
		Name:     info.Name(),
		Size:     info.Size(),
		Type:     itemType,
		Modified: info.ModTime(),
	}
}
