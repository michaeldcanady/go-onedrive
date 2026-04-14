package fs

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/michaeldcanady/go-onedrive/internal/concurrency"
)

const (
	// defaultConcurrency sets the maximum number of simultaneous file operations during recursive actions.
	defaultConcurrency = 5
)

// FileSystemManager orchestrates operations across multiple filesystem providers.
// It handles path resolution and provides high-level logic for cross-provider and recursive operations.
type FileSystemManager struct {
	// registry is the source for resolving providers based on path prefixes.
	registry ServiceRegistry
}

// ServiceRegistry is a subset of the Registry interface needed by the Manager.
type ServiceRegistry interface {
	Get(provider string) (Service, error)
}

// NewFileSystemManager initializes a new instance of the FileSystemManager.
func NewFileSystemManager(registry ServiceRegistry) *FileSystemManager {
	return &FileSystemManager{
		registry: registry,
	}
}

func (m *FileSystemManager) Name() string {
	return "fs_manager"
}

// Get retrieves metadata for an item by its structured URI.
func (m *FileSystemManager) Get(ctx context.Context, uri *URI) (Item, error) {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return Item{}, err
	}
	return p.Get(ctx, uri)
}

// Stat returns metadata for an item at the specified URI.
func (m *FileSystemManager) Stat(ctx context.Context, uri *URI) (Item, error) {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return Item{}, err
	}
	return p.Stat(ctx, uri)
}

// List returns the children of a directory at the specified URI.
func (m *FileSystemManager) List(ctx context.Context, uri *URI, opts ListOptions) ([]Item, error) {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return nil, err
	}
	return p.List(ctx, uri, opts)
}

// ReadFile opens a read stream for a file's content.
func (m *FileSystemManager) ReadFile(ctx context.Context, uri *URI, opts ReadOptions) (io.ReadCloser, error) {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return nil, err
	}
	return p.ReadFile(ctx, uri, opts)
}

// WriteFile creates or updates a file with the content from the provided reader.
func (m *FileSystemManager) WriteFile(ctx context.Context, uri *URI, r io.Reader, opts WriteOptions) (Item, error) {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return Item{}, err
	}
	return p.WriteFile(ctx, uri, r, opts)
}

// Mkdir creates a new directory at the specified URI.
func (m *FileSystemManager) Mkdir(ctx context.Context, uri *URI) error {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return err
	}
	return p.Mkdir(ctx, uri)
}

// Remove deletes an item from its respective provider.
func (m *FileSystemManager) Remove(ctx context.Context, uri *URI) error {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return err
	}
	return p.Remove(ctx, uri)
}

// Touch creates an empty file or updates the timestamp of an existing one.
func (m *FileSystemManager) Touch(ctx context.Context, uri *URI) (Item, error) {
	p, err := m.registry.Get(uri.Provider)
	if err != nil {
		return Item{}, err
	}
	return p.Touch(ctx, uri)
}

// Copy duplicates an item from a source URI to a destination URI, supporting cross-provider copy.
func (m *FileSystemManager) Copy(ctx context.Context, src, dst *URI, opts CopyOptions) error {
	srcItem, err := m.Stat(ctx, src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if srcItem.Type == TypeFolder && !opts.Recursive {
		return fmt.Errorf("omitting directory '%s'", src.String())
	}

	if opts.Recursive {
		pool := concurrency.NewWorkerPool(concurrency.WithCapacity(defaultConcurrency))
		return m.copyRecursive(ctx, src, dst, opts, pool)
	}

	return m.copySingle(ctx, srcItem, src, dst, opts)
}

func (m *FileSystemManager) copySingle(ctx context.Context, srcItem Item, src, dst *URI, opts CopyOptions) error {
	pSrc, err := m.registry.Get(src.Provider)
	if err != nil {
		return err
	}

	pDst, err := m.registry.Get(dst.Provider)
	if err != nil {
		return err
	}

	// If same provider, delegate to it for potentially optimized copy.
	if pSrc == pDst {
		return pSrc.Copy(ctx, src, dst, opts)
	}

	// Cross-provider copy: stream data from source to destination.
	r, err := pSrc.ReadFile(ctx, src, ReadOptions{})
	if err != nil {
		return fmt.Errorf("failed to read source for cross-provider copy: %w", err)
	}
	defer r.Close()

	writeOpts := WriteOptions{
		Overwrite: opts.Overwrite,
		Size:      srcItem.Size,
	}

	if _, err := pDst.WriteFile(ctx, dst, r, writeOpts); err != nil {
		return fmt.Errorf("failed to write destination for cross-provider copy: %w", err)
	}

	return nil
}

func (m *FileSystemManager) copyRecursive(ctx context.Context, src, dst *URI, opts CopyOptions, pool *concurrency.WorkerPool) error {
	item, err := m.Stat(ctx, src)
	if err != nil {
		return err
	}

	if item.Type == TypeFile {
		pool.Submit(func() {
			fileOpts := opts
			fileOpts.Recursive = false
			if err := m.copySingle(ctx, item, src, dst, fileOpts); err != nil {
				// error handling
			}
		})
		return nil
	}

	if err := m.Mkdir(ctx, dst); err != nil {
		return err
	}

	children, err := m.List(ctx, src, ListOptions{
		Recursive: opts.Recursive,
	})
	if err != nil {
		return err
	}

	for _, child := range children {
		childSrc := m.Join(src, child.Name)
		childDst := m.Join(dst, child.Name)
		if err := m.copyRecursive(ctx, childSrc, childDst, opts, pool); err != nil {
			return err
		}
	}

	return nil
}

// Move relocates or renames an item, supporting cross-provider move via copy and delete.
func (m *FileSystemManager) Move(ctx context.Context, src, dst *URI) error {
	pSrc, err := m.registry.Get(src.Provider)
	if err != nil {
		return err
	}

	pDst, err := m.registry.Get(dst.Provider)
	if err != nil {
		return err
	}

	// If same provider, delegate to it for potentially optimized move.
	if pSrc == pDst {
		return pSrc.Move(ctx, src, dst)
	}

	// Cross-provider move: Copy + Delete.
	if err := m.Copy(ctx, src, dst, CopyOptions{Overwrite: true}); err != nil {
		return err
	}

	return pSrc.Remove(ctx, src)
}

// Join combines a base URI and a name, returning a new URI.
func (m *FileSystemManager) Join(base *URI, name string) *URI {
	newURI := *base // shallow copy
	newURI.Path = path.Join(base.Path, name)
	return &newURI
}
