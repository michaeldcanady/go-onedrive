package manager

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
)

const (
	providerName = "fs_manager"
	// defaultConcurrency sets the maximum number of simultaneous file operations during recursive actions.
	defaultConcurrency = 5
)

// FileSystemManager orchestrates operations across multiple filesystem providers.
// It handles path resolution and provides high-level logic for cross-provider and recursive operations.
type FileSystemManager struct {
	// registry is the source for resolving providers based on path prefixes.
	registry registry.Service
}

// NewFileSystemManager initializes a new instance of the FileSystemManager.
func NewFileSystemManager(registry registry.Service) *FileSystemManager {
	return &FileSystemManager{
		registry: registry,
	}
}

func (m *FileSystemManager) Name() string {
	return providerName
}

// Get retrieves metadata for an item by its path, resolving the appropriate provider.
func (m *FileSystemManager) Get(ctx context.Context, path string) (shared.Item, error) {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return shared.Item{}, err
	}
	return p.Get(ctx, subPath)
}

// Stat returns metadata for an item at the specified path.
func (m *FileSystemManager) Stat(ctx context.Context, path string) (shared.Item, error) {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return shared.Item{}, err
	}
	return p.Stat(ctx, subPath)
}

// List returns the children of a directory at the specified path.
func (m *FileSystemManager) List(ctx context.Context, path string, opts shared.ListOptions) ([]shared.Item, error) {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return nil, err
	}
	return p.List(ctx, subPath, opts)
}

// ReadFile opens a read stream for a file's content.
func (m *FileSystemManager) ReadFile(ctx context.Context, path string, opts shared.ReadOptions) (io.ReadCloser, error) {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return nil, err
	}
	return p.ReadFile(ctx, subPath, opts)
}

// WriteFile creates or updates a file with the content from the provided reader.
func (m *FileSystemManager) WriteFile(ctx context.Context, path string, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return shared.Item{}, err
	}
	return p.WriteFile(ctx, subPath, r, opts)
}

// Mkdir creates a new directory at the specified path.
func (m *FileSystemManager) Mkdir(ctx context.Context, path string) error {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return err
	}
	return p.Mkdir(ctx, subPath)
}

// Remove deletes an item from its respective provider.
func (m *FileSystemManager) Remove(ctx context.Context, path string) error {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return err
	}
	return p.Remove(ctx, subPath)
}

// Touch creates an empty file or updates the timestamp of an existing one.
func (m *FileSystemManager) Touch(ctx context.Context, path string) (shared.Item, error) {
	p, subPath, err := m.registry.Resolve(ctx, path)
	if err != nil {
		return shared.Item{}, err
	}
	return p.Touch(ctx, subPath)
}

// Copy duplicates an item from a source path to a destination path, supporting cross-provider copy.
func (m *FileSystemManager) Copy(ctx context.Context, src, dst string, opts shared.CopyOptions) error {
	srcItem, err := m.Stat(ctx, src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if srcItem.Type == shared.TypeFolder && !opts.Recursive {
		return fmt.Errorf("omitting directory '%s'", src)
	}

	if opts.Recursive {
		sem := make(chan struct{}, defaultConcurrency)
		return m.copyRecursive(ctx, src, dst, opts, sem)
	}

	return m.copySingle(ctx, srcItem, src, dst, opts)
}

func (m *FileSystemManager) copySingle(ctx context.Context, srcItem shared.Item, src, dst string, opts shared.CopyOptions) error {
	pSrc, srcSubPath, err := m.registry.Resolve(ctx, src)
	if err != nil {
		return err
	}

	pDst, dstSubPath, err := m.registry.Resolve(ctx, dst)
	if err != nil {
		return err
	}

	// If same provider, delegate to it for potentially optimized copy.
	if pSrc == pDst {
		return pSrc.Copy(ctx, srcSubPath, dstSubPath, opts)
	}

	// Cross-provider copy: stream data from source to destination.
	r, err := pSrc.ReadFile(ctx, srcSubPath, shared.ReadOptions{})
	if err != nil {
		return fmt.Errorf("failed to read source for cross-provider copy: %w", err)
	}
	defer r.Close()

	writeOpts := shared.WriteOptions{
		Overwrite: opts.Overwrite,
		Size:      srcItem.Size, // Pass source item size to WriteOptions
	}

	if _, err := pDst.WriteFile(ctx, dstSubPath, r, writeOpts); err != nil {
		return fmt.Errorf("failed to write destination for cross-provider copy: %w", err)
	}

	return nil
}

func (m *FileSystemManager) copyRecursive(ctx context.Context, src, dst string, opts shared.CopyOptions, sem chan struct{}) error {
	item, err := m.Stat(ctx, src)
	if err != nil {
		return err
	}

	if item.Type == shared.TypeFile {
		sem <- struct{}{}
		defer func() { <-sem }()

		fileOpts := opts
		fileOpts.Recursive = false
		return m.copySingle(ctx, item, src, dst, fileOpts)
	}

	// It's a folder, ensure destination exists.
	_ = m.Mkdir(ctx, dst)

	children, err := m.List(ctx, src, shared.ListOptions{
		Recursive: opts.Recursive,
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(children))

	for _, child := range children {
		wg.Add(1)
		go func(c shared.Item) {
			defer wg.Done()
			childSrc := m.Join(src, c.Name)
			childDst := m.Join(dst, c.Name)
			if err := m.copyRecursive(ctx, childSrc, childDst, opts, sem); err != nil {
				errs <- err
			}
		}(child)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// Move relocates or renames an item, supporting cross-provider move via copy and delete.
func (m *FileSystemManager) Move(ctx context.Context, src, dst string) error {
	pSrc, srcSubPath, err := m.registry.Resolve(ctx, src)
	if err != nil {
		return err
	}

	pDst, dstSubPath, err := m.registry.Resolve(ctx, dst)
	if err != nil {
		return err
	}

	// If same provider, delegate to it for potentially optimized move.
	if pSrc == pDst {
		return pSrc.Move(ctx, srcSubPath, dstSubPath)
	}

	// Cross-provider move: Copy + Delete.
	if err := m.Copy(ctx, src, dst, shared.CopyOptions{Overwrite: true}); err != nil {
		return err
	}

	return pSrc.Remove(ctx, srcSubPath)
}

// Join combines a base path and a name, preserving the provider prefix if present.
func (m *FileSystemManager) Join(base, name string) string {
	prefix, subPath, found := strings.Cut(base, ":")
	if !found {
		return path.Join(base, name)
	}

	subPath = strings.TrimPrefix(subPath, "//")
	return prefix + ":" + path.Join(subPath, name)
}
