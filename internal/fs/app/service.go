package app

import (
	"context"
	"io"
	"path"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/fs/registry"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
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

func (m *FileSystemManager) resolve(_ context.Context, path string) (domainfs.Service, string, error) {
	providerName, subPath := util.ParsePath(path)

	p, err := m.registry.Get(providerName)
	if err == nil {
		return p, subPath, nil
	}

	// Fallback to onedrive provider if the prefix was not a registered provider
	// but only if there was no prefix (handled by parsePath) or if the prefix lookup failed.
	p, err = m.registry.Get("onedrive")
	if err != nil {
		return nil, "", err
	}
	return p, path, nil
}

func (m *FileSystemManager) Get(ctx context.Context, path string) (domainfs.Item, error) {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Get(ctx, subPath)
}

func (m *FileSystemManager) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return nil, err
	}
	return p.List(ctx, subPath, opts)
}

func (m *FileSystemManager) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Stat(ctx, subPath, opts)
}

func (m *FileSystemManager) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return nil, err
	}
	return p.ReadFile(ctx, subPath, opts)
}

func (m *FileSystemManager) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.WriteFile(ctx, subPath, r, opts)
}

func (m *FileSystemManager) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return err
	}
	return p.Mkdir(ctx, subPath, opts)
}

func (m *FileSystemManager) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return err
	}
	return p.Remove(ctx, subPath, opts)
}

func (m *FileSystemManager) Copy(ctx context.Context, src, dst string, opts domainfs.CopyOptions) error {
	if opts.Recursive {
		return m.copyRecursive(ctx, src, dst, opts)
	}

	pSrc, srcSubPath, err := m.resolve(ctx, src)
	if err != nil {
		return err
	}

	pDst, dstSubPath, err := m.resolve(ctx, dst)
	if err != nil {
		return err
	}

	if pSrc == pDst {
		return pSrc.Copy(ctx, srcSubPath, dstSubPath, opts)
	}

	// Cross-provider copy
	r, err := pSrc.ReadFile(ctx, srcSubPath, domainfs.ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = pDst.WriteFile(ctx, dstSubPath, r, domainfs.WriteOptions{
		Overwrite: opts.Overwrite,
	})
	return err
}

func (m *FileSystemManager) copyRecursive(ctx context.Context, src, dst string, opts domainfs.CopyOptions) error {
	item, err := m.Get(ctx, src)
	if err != nil {
		return err
	}

	if opts.Matcher != nil && opts.Matcher.ShouldIgnore(item.Path, item.Type == domainfs.ItemTypeFolder) {
		return nil
	}

	if item.Type == domainfs.ItemTypeFile {
		// Just a regular copy for files
		opts.Recursive = false
		return m.Copy(ctx, src, dst, opts)
	}

	// It's a folder, ensure destination exists
	if err := m.Mkdir(ctx, dst, domainfs.MKDirOptions{Parents: true}); err != nil {
		// We might want to ignore "already exists" errors here if they occur
	}

	children, err := m.List(ctx, src, domainfs.ListOptions{})
	if err != nil {
		return err
	}

	for _, child := range children {
		// Use provider-aware join
		childSrc := m.joinPath(src, child.Name)
		childDst := m.joinPath(dst, child.Name)

		if err := m.copyRecursive(ctx, childSrc, childDst, opts); err != nil {
			return err
		}
	}

	return nil
}

func (m *FileSystemManager) joinPath(base, name string) string {
	provider, subPath := util.ParsePath(base)
	joined := path.Join(subPath, name)
	if provider == "" {
		return joined
	}
	return provider + ":" + joined
}

func (m *FileSystemManager) Move(ctx context.Context, src, dst string, opts domainfs.MoveOptions) error {
	pSrc, srcSubPath, err := m.resolve(ctx, src)
	if err != nil {
		return err
	}

	pDst, dstSubPath, err := m.resolve(ctx, dst)
	if err != nil {
		return err
	}

	if pSrc == pDst {
		return pSrc.Move(ctx, srcSubPath, dstSubPath, opts)
	}

	// Cross-provider move: Copy + Delete
	if err := m.Copy(ctx, src, dst, domainfs.CopyOptions{Overwrite: true}); err != nil {
		return err
	}

	return pSrc.Remove(ctx, srcSubPath, domainfs.RemoveOptions{})
}

func (m *FileSystemManager) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	pDst, dstSubPath, err := m.resolve(ctx, dst)
	if err != nil {
		return domainfs.Item{}, err
	}
	// src is local path in current implementation of Upload in odc
	return pDst.Upload(ctx, src, dstSubPath, opts)
}

func (m *FileSystemManager) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	p, subPath, err := m.resolve(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}
	return p.Touch(ctx, subPath, opts)
}
