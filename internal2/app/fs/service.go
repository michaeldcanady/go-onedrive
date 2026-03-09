package fs

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/app/fs/registry"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

var _ domainfs.Service = (*FileSystemManager)(nil)

type FileSystemManager struct {
	registry *registry.Registry
}

func NewFileSystemManager(registry *registry.Registry) *FileSystemManager {
	return &FileSystemManager{
		registry: registry,
	}
}

func parsePath(path string) (string, string) {
	prefix, rest, found := strings.Cut(path, ":")
	if !found {
		return "onedrive", path
	}

	rest = strings.TrimPrefix(rest, "//")

	return strings.ToLower(prefix), rest
}

func (m *FileSystemManager) Get(ctx context.Context, path string) (domainfs.Item, error) {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Get(ctx, path)
}

func (m *FileSystemManager) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return nil, err
	}
	return p.List(ctx, path, opts)
}

func (m *FileSystemManager) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Stat(ctx, path, opts)
}

func (m *FileSystemManager) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return nil, err
	}
	return p.ReadFile(ctx, path, opts)
}

func (m *FileSystemManager) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.WriteFile(ctx, path, r, opts)
}

func (m *FileSystemManager) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return err
	}
	return p.Mkdir(ctx, path, opts)
}

func (m *FileSystemManager) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return err
	}
	return p.Remove(ctx, path, opts)
}

func (m *FileSystemManager) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	srcProvider, srcPath := parsePath(src)
	pSrc, err := m.registry.Get(srcProvider)
	if err != nil {
		return err
	}

	dstProvider, dstPath := parsePath(dst)
	pDst, err := m.registry.Get(dstProvider)
	if err != nil {
		return err
	}

	if pSrc == pDst {
		return pSrc.Move(ctx, srcPath, dstPath, opts)
	}

	return errors.New("cross-provider move not supported yet")
}

func (m *FileSystemManager) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	provider, path := parsePath(dst)
	pDst, err := m.registry.Get(provider)
	if err != nil {
		return domainfs.Item{}, err
	}
	// src is local path in current implementation of Upload in odc
	return pDst.Upload(ctx, src, path, opts)
}

func (m *FileSystemManager) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	provider, path := parsePath(path)
	p, err := m.registry.Get(provider)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Touch(ctx, path, opts)
}
