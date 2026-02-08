// internal/fs/service_impl.go
package fs

import (
	"context"
	"io"

	domainDrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

var _ domainfs.Service = (*Service)(nil)

type Service struct {
	files         domainfile.FileService
	driveResolver domainDrive.DriveResolver
	logger        logging.Logger
}

func NewService(
	files domainfile.FileService,
	driveResolver domainDrive.DriveResolver,
	logger logging.Logger,
) *Service {
	return &Service{files: files, driveResolver: driveResolver, logger: logger}
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
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("resolving filesystem path",
		logging.String("event", "fs_resolve_path_start"),
		logging.String("path", p),
		logging.String("correlation_id", cid),
	)

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		s.logger.Error("failed to resolve current drive ID",
			logging.String("event", "fs_resolve_path_drive_error"),
			logging.Error(err),
			logging.String("path", p),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	s.logger.Debug("resolving item via file service",
		logging.String("event", "fs_resolve_path_fileservice"),
		logging.String("drive_id", driveID),
		logging.String("path", p),
		logging.String("correlation_id", cid),
	)

	item, err := s.files.ResolveItem(ctx, driveID, p)
	if err != nil {
		s.logger.Error("failed to resolve item",
			logging.String("event", "fs_resolve_path_error"),
			logging.String("drive_id", driveID),
			logging.String("path", p),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	s.logger.Info("path resolved successfully",
		logging.String("event", "fs_resolve_path_success"),
		logging.String("drive_id", driveID),
		logging.String("path", p),
		logging.Bool("is_folder", item.IsFolder),
		logging.String("correlation_id", cid),
	)

	return item, nil
}

func (s *Service) listChildren(ctx context.Context, p string, recursive bool) ([]*infrafile.DriveItem, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("listing children",
		logging.String("event", "fs_list_children_start"),
		logging.String("path", p),
		logging.Bool("recursive", recursive),
		logging.String("correlation_id", cid),
	)

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		s.logger.Error("failed to resolve current drive ID",
			logging.String("event", "fs_list_children_drive_error"),
			logging.Error(err),
			logging.String("path", p),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	items, err := s.files.ListChildren(ctx, driveID, p)
	if err != nil {
		s.logger.Error("failed to list children",
			logging.String("event", "fs_list_children_error"),
			logging.String("drive_id", driveID),
			logging.String("path", p),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	s.logger.Debug("direct children retrieved",
		logging.String("event", "fs_list_children_direct"),
		logging.Int("count", len(items)),
		logging.String("path", p),
		logging.String("correlation_id", cid),
	)

	if !recursive {
		return items, nil
	}

	all := make([]*infrafile.DriveItem, 0, len(items))
	all = append(all, items...)

	for _, item := range items {
		if item.IsFolder {
			childPath := item.PathWithoutDrive + "/" + item.Name

			s.logger.Debug("recursively listing folder",
				logging.String("event", "fs_list_children_recurse"),
				logging.String("folder", childPath),
				logging.String("correlation_id", cid),
			)

			children, err := s.listChildren(ctx, childPath, true)
			if err != nil {
				return nil, err
			}

			all = append(all, children...)
		}
	}

	s.logger.Info("children listed successfully",
		logging.String("event", "fs_list_children_success"),
		logging.Int("total_count", len(all)),
		logging.String("path", p),
		logging.Bool("recursive", recursive),
		logging.String("correlation_id", cid),
	)

	return all, nil
}

func (s *Service) List(ctx context.Context, p string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("listing filesystem items",
		logging.String("event", "fs_list_start"),
		logging.String("path", p),
		logging.Bool("recursive", opts.Recursive),
		logging.String("correlation_id", cid),
	)

	item, err := s.resolvePath(ctx, p)
	if err != nil {
		s.logger.Error("failed to resolve path for listing",
			logging.String("event", "fs_list_resolve_error"),
			logging.String("path", p),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	if !item.IsFolder {
		s.logger.Info("path is a file; returning single item",
			logging.String("event", "fs_list_file"),
			logging.String("path", p),
			logging.String("correlation_id", cid),
		)
		return []domainfs.Item{mapToFSItem(item)}, nil
	}

	children, err := s.listChildren(ctx, p, opts.Recursive)
	if err != nil {
		s.logger.Error("failed to list children",
			logging.String("event", "fs_list_children_error"),
			logging.String("path", p),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	items := mapToFSItems(children)

	s.logger.Info("filesystem items listed successfully",
		logging.String("event", "fs_list_success"),
		logging.String("path", p),
		logging.Int("count", len(items)),
		logging.Bool("recursive", opts.Recursive),
		logging.String("correlation_id", cid),
	)

	return items, nil
}

func (s *Service) Stat(ctx context.Context, p string, opts domainfs.StatOptions) (domainfs.Item, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("statting filesystem item",
		logging.String("event", "fs_stat_start"),
		logging.String("path", p),
		logging.String("correlation_id", cid),
	)

	item, err := s.resolvePath(ctx, p)
	if err != nil {
		s.logger.Error("failed to stat path",
			logging.String("event", "fs_stat_error"),
			logging.String("path", p),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return domainfs.Item{}, err
	}

	s.logger.Info("stat completed successfully",
		logging.String("event", "fs_stat_success"),
		logging.String("path", p),
		logging.Bool("is_folder", item.IsFolder),
		logging.String("correlation_id", cid),
	)

	return mapToFSItem(item), nil
}
