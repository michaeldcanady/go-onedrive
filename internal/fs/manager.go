package fs

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/michaeldcanady/go-onedrive/internal/concurrency"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
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
	log      logger.Logger
}

// ServiceRegistry is a subset of the Registry interface needed by the Manager.
type ServiceRegistry interface {
	Resolve(ctx context.Context, path *URI) (Service, *URI, error)
}

// NewFileSystemManager initializes a new instance of the FileSystemManager.
func NewFileSystemManager(registry ServiceRegistry, l logger.Logger) *FileSystemManager {
	return &FileSystemManager{
		registry: registry,
		log:      l,
	}
}

func (m *FileSystemManager) Name() string {
	return "fs_manager"
}

// Get retrieves metadata for an item by its path, resolving the appropriate provider.
func (m *FileSystemManager) Get(ctx context.Context, uri *URI) (Item, error) {
	log := m.log.WithContext(ctx)
	log.Debug("fs manager: get", logger.String("path", uri.String()))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		if _, ok := err.(*errors.UnregisteredProviderError); !ok {
			log.Error("failed to resolve provider for get", logger.String("path", uri.String()), logger.Error(err))
			return Item{}, err
		}

		log.Error("an unknown error occurred", logger.String("path", uri.String()), logger.Error(err))
		return Item{}, errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.Get(ctx, resolvedURI)
}

// Stat returns metadata for an item at the specified path.
func (m *FileSystemManager) Stat(ctx context.Context, uri *URI) (Item, error) {
	log := m.log.WithContext(ctx)
	log.Debug("fs manager: stat", logger.String("path", uri.String()))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		if _, ok := err.(*errors.UnregisteredProviderError); !ok {
			log.Error("failed to resolve provider for stat", logger.String("path", uri.String()), logger.Error(err))
			return Item{}, err
		}

		log.Error("an unknown error occurred", logger.String("path", uri.String()), logger.Error(err))
		return Item{}, errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.Stat(ctx, resolvedURI)
}

// List returns the children of a directory at the specified path.
func (m *FileSystemManager) List(ctx context.Context, uri *URI, opts ListOptions) ([]Item, error) {
	log := m.log.WithContext(ctx)
	log.Debug("fs manager: list", logger.String("path", uri.String()), logger.Bool("recursive", opts.Recursive))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		if _, ok := err.(*errors.UnregisteredProviderError); !ok {
			log.Error("failed to resolve provider for list", logger.String("path", uri.String()), logger.Error(err))
			return nil, err
		}

		log.Error("an unknown error occurred", logger.String("path", uri.String()), logger.Error(err))
		return nil, errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.List(ctx, resolvedURI, opts)
}

// ReadFile opens a read stream for a file's content.
func (m *FileSystemManager) ReadFile(ctx context.Context, uri *URI, opts ReadOptions) (io.ReadCloser, error) {
	log := m.log.WithContext(ctx)
	log.Debug("fs manager: read file", logger.String("path", uri.String()))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		if _, ok := err.(*errors.UnregisteredProviderError); !ok {
			log.Error("failed to resolve provider for read", logger.String("path", uri.String()), logger.Error(err))
			return nil, err
		}

		log.Error("fs manager: failed to resolve provider for read", logger.String("path", uri.String()), logger.Error(err))
		return nil, errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.ReadFile(ctx, resolvedURI, opts)
}

// WriteFile creates or updates a file with the content from the provided reader.
func (m *FileSystemManager) WriteFile(ctx context.Context, uri *URI, r io.Reader, opts WriteOptions) (Item, error) {
	log := m.log.WithContext(ctx)
	log.Info("fs manager: write file", logger.String("path", uri.String()), logger.Bool("overwrite", opts.Overwrite))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		log.Error("fs manager: failed to resolve provider for write", logger.String("path", uri.String()), logger.Error(err))
		return Item{}, errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.WriteFile(ctx, resolvedURI, r, opts)
}

// Mkdir creates a new directory at the specified path.
func (m *FileSystemManager) Mkdir(ctx context.Context, uri *URI) error {
	log := m.log.WithContext(ctx)
	log.Info("fs manager: mkdir", logger.String("path", uri.String()))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		log.Error("fs manager: failed to resolve provider for mkdir", logger.String("path", uri.String()), logger.Error(err))
		return errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.Mkdir(ctx, resolvedURI)
}

// Remove deletes an item from its respective provider.
func (m *FileSystemManager) Remove(ctx context.Context, uri *URI) error {
	log := m.log.WithContext(ctx)
	log.Info("fs manager: remove", logger.String("path", uri.String()))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		log.Error("fs manager: failed to resolve provider for remove", logger.String("path", uri.String()), logger.Error(err))
		return errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.Remove(ctx, resolvedURI)
}

// Touch creates an empty file or updates the timestamp of an existing one.
func (m *FileSystemManager) Touch(ctx context.Context, uri *URI) (Item, error) {
	log := m.log.WithContext(ctx)
	log.Info("fs manager: touch", logger.String("path", uri.String()))
	p, resolvedURI, err := m.registry.Resolve(ctx, uri)
	if err != nil {
		log.Error("fs manager: failed to resolve provider for touch", logger.String("path", uri.String()), logger.Error(err))
		return Item{}, errors.NewInvalidInput(err, "could not resolve filesystem provider", "Check if the provider prefix is correct and registered.").WithContext(errors.KeyPath, uri.String())
	}
	return p.Touch(ctx, resolvedURI)
}

// Copy duplicates an item from a source path to a destination path, supporting cross-provider copy.
func (m *FileSystemManager) Copy(ctx context.Context, src, dst *URI, opts CopyOptions) error {
	log := m.log.WithContext(ctx).With(logger.String("src", src.String()), logger.String("dst", dst.String()))
	log.Info("fs manager: copy", logger.Bool("recursive", opts.Recursive))

	srcItem, err := m.Stat(ctx, src)
	if err != nil {
		return err
	}

	if srcItem.Type == TypeFolder && !opts.Recursive {
		return errors.NewInvalidInput(nil, fmt.Sprintf("omitting directory '%s'", src.String()), "Use -r or --recursive to copy directories.")
	}

	if opts.Recursive {
		pool := concurrency.NewWorkerPool(concurrency.WithCapacity(defaultConcurrency))
		return m.copyRecursive(ctx, src, dst, opts, pool)
	}

	return m.copySingle(ctx, srcItem, src, dst, opts)
}

func (m *FileSystemManager) copySingle(ctx context.Context, srcItem Item, src, dst *URI, opts CopyOptions) error {
	log := m.log.WithContext(ctx).With(logger.String("src", src.String()), logger.String("dst", dst.String()))

	pSrc, resolvedSrc, err := m.registry.Resolve(ctx, src)
	if err != nil {
		return errors.NewInvalidInput(err, "could not resolve source filesystem provider", "Check if the source provider prefix is correct and registered.").WithContext(errors.KeyPath, src.String())
	}

	pDst, resolvedDst, err := m.registry.Resolve(ctx, dst)
	if err != nil {
		return errors.NewInvalidInput(err, "could not resolve destination filesystem provider", "Check if the destination provider prefix is correct and registered.").WithContext(errors.KeyPath, dst.String())
	}

	// If same provider, delegate to it for potentially optimized copy.
	if pSrc == pDst {
		log.Debug("fs manager: delegating single copy to provider", logger.String("provider", pSrc.Name()))
		return pSrc.Copy(ctx, resolvedSrc, resolvedDst, opts)
	}

	// Cross-provider copy: stream data from source to destination.
	log.Info("fs manager: performing cross-provider copy", logger.String("src_provider", pSrc.Name()), logger.String("dst_provider", pDst.Name()))
	r, err := pSrc.ReadFile(ctx, resolvedSrc, ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	writeOpts := WriteOptions{
		Overwrite: opts.Overwrite,
		Size:      srcItem.Size,
	}

	if _, err := pDst.WriteFile(ctx, resolvedDst, r, writeOpts); err != nil {
		return err
	}

	return nil
}

func (m *FileSystemManager) copyRecursive(ctx context.Context, src, dst *URI, opts CopyOptions, pool *concurrency.WorkerPool) error {
	log := m.log.WithContext(ctx).With(logger.String("src", src.String()), logger.String("dst", dst.String()))

	item, err := m.Stat(ctx, src)
	if err != nil {
		return err
	}

	if item.Type == TypeFile {
		pool.Submit(func() {
			fileOpts := opts
			fileOpts.Recursive = false
			if err := m.copySingle(ctx, item, src, dst, fileOpts); err != nil {
				log.Error("fs manager: failed to copy file during recursive operation", logger.String("file", src.String()), logger.Error(err))
			}
		})
		return nil
	}

	log.Debug("fs manager: recursive copy - creating directory", logger.String("path", dst.String()))
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
	log := m.log.WithContext(ctx).With(logger.String("src", src.String()), logger.String("dst", dst.String()))
	log.Info("fs manager: move")

	pSrc, resolvedSrc, err := m.registry.Resolve(ctx, src)
	if err != nil {
		return errors.NewInvalidInput(err, "could not resolve source filesystem provider", "Check if the source provider prefix is correct and registered.").WithContext(errors.KeyPath, src.String())
	}

	pDst, resolvedDst, err := m.registry.Resolve(ctx, dst)
	if err != nil {
		return errors.NewInvalidInput(err, "could not resolve destination filesystem provider", "Check if the destination provider prefix is correct and registered.").WithContext(errors.KeyPath, dst.String())
	}

	// If same provider, delegate to it for potentially optimized move.
	if pSrc == pDst {
		log.Debug("fs manager: delegating move to provider", logger.String("provider", pSrc.Name()))
		return pSrc.Move(ctx, resolvedSrc, resolvedDst)
	}

	// Cross-provider move: Copy + Delete.
	log.Info("fs manager: performing cross-provider move", logger.String("src_provider", pSrc.Name()), logger.String("dst_provider", pDst.Name()))
	if err := m.Copy(ctx, src, dst, CopyOptions{Overwrite: true}); err != nil {
		return err
	}

	return pSrc.Remove(ctx, resolvedSrc)
}

// Join combines a base URI and a name, returning a new URI.
func (m *FileSystemManager) Join(base *URI, name string) *URI {
	newPath := path.Join(base.Path, name)
	return &URI{
		Provider: base.Provider,
		DriveRef: base.DriveRef,
		Path:     newPath,
	}
}
