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

func (m *FileSystemManager) getProvider(ctx context.Context, name, fullPath string) (domainfs.Service, string, error) {
	p, err := m.registry.Get(name)
	if err == nil {
		_, subPath := parsePath(fullPath)
		return p, subPath, nil
	}

	// Fallback to onedrive provider
	p, err = m.registry.Get("onedrive")
	if err != nil {
		return nil, "", err
	}
	return p, fullPath, nil
}

func (m *FileSystemManager) Get(ctx context.Context, path string) (domainfs.Item, error) {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Get(ctx, subPath)
}

func (m *FileSystemManager) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return nil, err
	}
	return p.List(ctx, subPath, opts)
}

func (m *FileSystemManager) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Stat(ctx, subPath, opts)
}

func (m *FileSystemManager) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return nil, err
	}
	return p.ReadFile(ctx, subPath, opts)
}

func (m *FileSystemManager) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.WriteFile(ctx, subPath, r, opts)
}

func (m *FileSystemManager) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return err
	}
	return p.Mkdir(ctx, subPath, opts)
}

func (m *FileSystemManager) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return err
	}
	return p.Remove(ctx, subPath, opts)
}

func (m *FileSystemManager) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	srcProviderName, _ := parsePath(src)
	pSrc, srcSubPath, err := m.getProvider(ctx, srcProviderName, src)
	if err != nil {
		return err
	}

	dstProviderName, _ := parsePath(dst)
	pDst, dstSubPath, err := m.getProvider(ctx, dstProviderName, dst)
	if err != nil {
		return err
	}

	if pSrc == pDst {
		return pSrc.Move(ctx, srcSubPath, dstSubPath, opts)
	}

	return errors.New("cross-provider move not supported yet")
}

func (m *FileSystemManager) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	providerName, _ := parsePath(dst)
	pDst, dstSubPath, err := m.getProvider(ctx, providerName, dst)
	if err != nil {
		return domainfs.Item{}, err
	}
	// src is local path in current implementation of Upload in odc
	return pDst.Upload(ctx, src, dstSubPath, opts)
}

func (m *FileSystemManager) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	providerName, _ := parsePath(path)
	p, subPath, err := m.getProvider(ctx, providerName, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Touch(ctx, subPath, opts)
}
