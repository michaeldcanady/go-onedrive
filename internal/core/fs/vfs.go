package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// VFS (Virtual FileSystem) orchestrates multiple backends via mount points.
type VFS struct {
	mounts   map[string]fs.Backend
	identity identity.Service
}

// NewVFS initializes a new Virtual FileSystem.
func NewVFS(idService identity.Service) *VFS {
	return &VFS{
		mounts:   make(map[string]fs.Backend),
		identity: idService,
	}
}

// Get implements the fs.Reader interface for the VFS.
func (v *VFS) Get(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	return v.Stat(ctx, uri)
}

// Mount associates a path prefix with a backend.
func (v *VFS) Mount(prefix string, backend fs.Backend) {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	v.mounts[prefix] = backend
}

// Resolve identifies the best matching mount point and relative path for a given absolute path.
func (v *VFS) Resolve(absPath string) (string, string, error) {
	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}

	var bestPrefix string
	for prefix := range v.mounts {
		if strings.HasPrefix(absPath, prefix) {
			if prefix == "/" || len(absPath) == len(prefix) || absPath[len(prefix)] == '/' {
				if len(prefix) > len(bestPrefix) {
					bestPrefix = prefix
				}
			}
		}
	}

	if bestPrefix == "" {
		return "", "", fmt.Errorf("no backend mounted for path: %s", absPath)
	}

	relPath := strings.TrimPrefix(absPath, bestPrefix)
	if relPath == "" {
		relPath = "/"
	}
	if !strings.HasPrefix(relPath, "/") {
		relPath = "/" + relPath
	}

	return bestPrefix, relPath, nil
}

func (v *VFS) resolve(absPath string) (fs.Backend, string, error) {
	prefix, relPath, err := v.Resolve(absPath)
	if err != nil {
		return nil, "", err
	}
	return v.mounts[prefix], relPath, nil
}

func (v *VFS) selectBackend(uri *fs.URI) (fs.Backend, string, error) {
	if uri.Provider != "" {
		if backend, ok := v.mounts[uri.Provider]; ok {
			return backend, uri.Path, nil
		}
	}
	return v.resolve(uri.Path)
}

func (v *VFS) Name() string {
	return "vfs"
}

func (v *VFS) Stat(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return fs.Item{}, err
	}
	token, err := v.getToken(ctx, uri.Provider)
	if err != nil {
		return fs.Item{}, err
	}
	return backend.Stat(ctx, token, uri.DriveID, relPath)
}

func (v *VFS) List(ctx context.Context, uri *fs.URI, opts fs.ListOptions) ([]fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return nil, err
	}
	token, err := v.getToken(ctx, backend.Name())
	if err != nil {
		return nil, err
	}
	return backend.List(ctx, token, uri.DriveID, relPath)
}

func (v *VFS) ReadFile(ctx context.Context, uri *fs.URI, opts fs.ReadOptions) (io.ReadCloser, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return nil, err
	}
	token, err := v.getToken(ctx, uri.Provider)
	if err != nil {
		return nil, err
	}
	return backend.Open(ctx, token, uri.DriveID, relPath)
}

func (v *VFS) WriteFile(ctx context.Context, uri *fs.URI, r io.Reader, opts fs.WriteOptions) (fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return fs.Item{}, err
	}
	token, err := v.getToken(ctx, uri.Provider)
	if err != nil {
		return fs.Item{}, err
	}
	return backend.Create(ctx, token, uri.DriveID, relPath, r)
}

func (v *VFS) Mkdir(ctx context.Context, uri *fs.URI) error {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return err
	}
	token, err := v.getToken(ctx, uri.Provider)
	if err != nil {
		return err
	}
	return backend.Mkdir(ctx, token, uri.DriveID, relPath)
}

func (v *VFS) Remove(ctx context.Context, uri *fs.URI) error {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return err
	}
	token, err := v.getToken(ctx, uri.Provider)
	if err != nil {
		return err
	}
	return backend.Remove(ctx, token, uri.DriveID, relPath)
}

func (v *VFS) Touch(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	backend, relPath, err := v.selectBackend(uri)
	if err != nil {
		return fs.Item{}, err
	}
	token, err := v.getToken(ctx, uri.Provider)
	if err != nil {
		return fs.Item{}, err
	}
	return backend.Create(ctx, token, uri.DriveID, relPath, strings.NewReader(""))
}

func (v *VFS) Copy(ctx context.Context, src, dst *fs.URI, opts fs.CopyOptions) error {
	srcBackend, srcRel, err := v.selectBackend(src)
	if err != nil {
		return err
	}
	dstBackend, dstRel, err := v.selectBackend(dst)
	if err != nil {
		return err
	}

	token, err := v.getToken(ctx, src.Provider)
	if err != nil {
		return err
	}

	if srcBackend == dstBackend {
		if adv, ok := srcBackend.(fs.AdvancedBackend); ok && srcBackend.Capabilities().CanCopy {
			return adv.Copy(ctx, token, src.DriveID, srcRel, dstRel)
		}
	}
	r, err := srcBackend.Open(ctx, token, src.DriveID, srcRel)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = dstBackend.Create(ctx, token, dst.DriveID, dstRel, r)
	return err
}

func (v *VFS) Move(ctx context.Context, src, dst *fs.URI) error {
	srcBackend, srcRel, err := v.selectBackend(src)
	if err != nil {
		return err
	}
	dstBackend, dstRel, err := v.selectBackend(dst)
	if err != nil {
		return err
	}

	token, err := v.getToken(ctx, src.Provider)
	if err != nil {
		return err
	}

	if srcBackend == dstBackend {
		if adv, ok := srcBackend.(fs.AdvancedBackend); ok && srcBackend.Capabilities().CanMove {
			return adv.Move(ctx, token, src.DriveID, srcRel, dstRel)
		}
	}
	if err := v.Copy(ctx, src, dst, fs.CopyOptions{Overwrite: true}); err != nil {
		return err
	}
	return srcBackend.Remove(ctx, token, src.DriveID, srcRel)
}

func (v *VFS) Mounts() []string {
	prefixes := make([]string, 0, len(v.mounts))
	for p := range v.mounts {
		prefixes = append(prefixes, p)
	}
	return prefixes
}

func (v *VFS) getToken(ctx context.Context, provider string) (string, error) {
	// TODO: mind way to resolve identity provider for a backend
	provider = "microsoft"
	ids, err := v.identity.GetStore().List(ctx, provider)
	if err != nil || len(ids) == 0 {
		return "", fmt.Errorf("no identity found for provider %s", provider)
	}
	identityID := ids[0]

	accessToken, err := v.identity.GetStore().Get(ctx, provider, identityID)
	if err != nil {
		return "", err
	}

	rawToken, err := json.Marshal(accessToken)
	if err != nil {
		return "", err
	}

	return string(rawToken), nil
}
