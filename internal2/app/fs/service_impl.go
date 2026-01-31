// internal/fs/service_impl.go
package fs

import (
	"context"
	"io"

	domainDrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

var _ domainfs.Service = (*Service)(nil)

type Service struct {
	files         domainfile.FileService
	driveResolver domainDrive.DriveResolver
}

func NewService(files domainfile.FileService, driveResolver domainDrive.DriveResolver) *Service {
	return &Service{files: files, driveResolver: driveResolver}
}

// Mkdir implements [Service].
func (s *Service) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	panic("unimplemented")
}

// Move implements [Service].
func (s *Service) Move(ctx context.Context, src string, dst string, opts domainfs.MoveOptions) error {
	panic("unimplemented")
}

// ReadFile implements [Service].
func (s *Service) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	panic("unimplemented")
}

// Remove implements [Service].
func (s *Service) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	panic("unimplemented")
}

// WriteFile implements [Service].
func (s *Service) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) error {
	panic("unimplemented")
}

func (s *Service) resolvePath(ctx context.Context, p string) (*infrafile.DriveItem, error) {
	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		return nil, err
	}

	return s.files.ResolveItem(ctx, driveID, p)
}

func (s *Service) listChildren(ctx context.Context, p string, recursive bool) ([]*infrafile.DriveItem, error) {
	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		return nil, err
	}

	// List direct children
	items, err := s.files.ListChildren(ctx, driveID, p)
	if err != nil {
		return nil, err
	}

	// If not recursive, we're done
	if !recursive {
		return items, nil
	}

	// Accumulate all items (direct + recursive)
	all := make([]*infrafile.DriveItem, 0, len(items))
	all = append(all, items...)

	// Recursively fetch children of folders
	for _, item := range items {
		if item.IsFolder {
			childPath := item.PathWithoutDrive + "/" + item.Name

			children, err := s.listChildren(ctx, childPath, true)
			if err != nil {
				return nil, err
			}

			all = append(all, children...)
		}
	}

	return all, nil
}

func (s *Service) List(ctx context.Context, p string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	item, err := s.resolvePath(ctx, p)
	if err != nil {
		return nil, err
	}
	if !item.IsFolder {
		return []domainfs.Item{mapToFSItem(item, p)}, nil
	}

	children, err := s.listChildren(ctx, p, opts.Recursive)
	if err != nil {
		return nil, err
	}
	return mapToFSItems(children, p), nil
}

func (s *Service) Stat(ctx context.Context, p string, opts domainfs.StatOptions) (domainfs.Item, error) {
	item, err := s.resolvePath(ctx, p)
	if err != nil {
		return domainfs.Item{}, err
	}
	return mapToFSItem(item, p), nil
}
