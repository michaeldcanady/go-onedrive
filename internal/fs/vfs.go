package fs

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// VFS (Virtual FileSystem) orchestrates multiple backends via mount points.
type VFS struct {
	mounts map[string]fs.Backend
}

// NewVFS initializes a new Virtual FileSystem.
func NewVFS() *VFS {
	return &VFS{
		mounts: make(map[string]fs.Backend),
	}
}

// Mount associates a path prefix with a backend.
func (v *VFS) Mount(prefix string, backend fs.Backend) {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	v.mounts[prefix] = backend
}

// resolve finds the appropriate backend and relative path for a given absolute path.
func (v *VFS) resolve(absPath string) (fs.Backend, string, error) {
	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}

	// Find the longest matching prefix
	var bestPrefix string
	for prefix := range v.mounts {
		if strings.HasPrefix(absPath, prefix) {
			if len(prefix) > len(bestPrefix) {
				bestPrefix = prefix
			}
		}
	}

	if bestPrefix == "" {
		return nil, "", fmt.Errorf("no backend mounted for path: %s", absPath)
	}

	backend := v.mounts[bestPrefix]
	relPath := strings.TrimPrefix(absPath, bestPrefix)
	if relPath == "" {
		relPath = "/"
	}
	if !strings.HasPrefix(relPath, "/") {
		relPath = "/" + relPath
	}

	return backend, relPath, nil
}

// selectBackend returns the backend and relative path for a given URI.
func (v *VFS) selectBackend(uri *fs.URI) (fs.Backend, string, error) {
	if uri.Provider != "" {
		if backend, ok := v.mounts[uri.Provider]; ok {
			return backend, uri.Path, nil
		}
	}
	return v.resolve(uri.Path)
}

// Name returns the name of the VFS.
func (v *VFS) Name() string {
	return "vfs"
}

// Stat returns metadata for an item at the specified path.
func (v *VFS) Stat(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return fs.Item{}, err
	}
	return backend.Stat(ctx, relPath)
}

// Get is an alias for Stat for backward compatibility.
func (v *VFS) Get(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	return v.Stat(ctx, uri)
}

// List returns the children of a directory.
func (v *VFS) List(ctx context.Context, uri *fs.URI, opts fs.ListOptions) ([]fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return nil, err
	}
	return backend.List(ctx, relPath)
}

// ReadFile opens a read stream for a file's content.
func (v *VFS) ReadFile(ctx context.Context, uri *fs.URI, opts fs.ReadOptions) (io.ReadCloser, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return nil, err
	}
	return backend.Open(ctx, relPath)
}

// WriteFile creates or updates a file.
func (v *VFS) WriteFile(ctx context.Context, uri *fs.URI, r io.Reader, opts fs.WriteOptions) (fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return fs.Item{}, err
	}
	return backend.Create(ctx, relPath, r)
}

// Mkdir creates a new directory.
func (v *VFS) Mkdir(ctx context.Context, uri *fs.URI) error {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return err
	}
	return backend.Mkdir(ctx, relPath)
}

// Remove deletes an item.
func (v *VFS) Remove(ctx context.Context, uri *fs.URI) error {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return err
	}
	return backend.Remove(ctx, relPath)
}

// Touch creates an empty file.
func (v *VFS) Touch(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return fs.Item{}, err
	}
	// Simplified Touch using Create with empty reader
	return backend.Create(ctx, relPath, strings.NewReader(""))
}

// Copy duplicates an item, supporting cross-backend copy via streaming.
func (v *VFS) Copy(ctx context.Context, src, dst *fs.URI, opts fs.CopyOptions) error {
	srcBackend, srcRel, err := v.selectBackend(src)
	if err != nil {
		return err
	}

	dstBackend, dstRel, err := v.selectBackend(dst)
	if err != nil {
		return err
	}

	// Native copy if same backend and supported
	if srcBackend == dstBackend {
		if adv, ok := srcBackend.(fs.AdvancedBackend); ok && srcBackend.Capabilities().CanCopy {
			return adv.Copy(ctx, srcRel, dstRel)
		}
	}

	// Cross-backend copy: stream
	r, err := srcBackend.Open(ctx, srcRel)
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = dstBackend.Create(ctx, dstRel, r)
	return err
}

// Move relocates an item.
func (v *VFS) Move(ctx context.Context, src, dst *fs.URI) error {
	srcBackend, srcRel, err := v.selectBackend(src)
	if err != nil {
		return err
	}

	dstBackend, dstRel, err := v.selectBackend(dst)
	if err != nil {
		return err
	}

	// Native move if same backend and supported
	if srcBackend == dstBackend {
		if adv, ok := srcBackend.(fs.AdvancedBackend); ok && srcBackend.Capabilities().CanMove {
			return adv.Move(ctx, srcRel, dstRel)
		}
	}

	// Cross-backend move: Copy + Remove
	if err := v.Copy(ctx, src, dst, fs.CopyOptions{Overwrite: true}); err != nil {
		return err
	}
	return srcBackend.Remove(ctx, srcRel)
}

// Mounts returns a sorted list of mount points.
func (v *VFS) Mounts() []string {
	prefixes := make([]string, 0, len(v.mounts))
	for p := range v.mounts {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)
	return prefixes
}
