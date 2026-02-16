package fs

import (
	"bytes"
	"context"
	"errors"
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

const (
	eventFSResolveStart   = "fs.resolve.start"
	eventFSResolveSuccess = "fs.resolve.success"
	eventFSResolveFailure = "fs.resolve.failure"
	eventFSReadFailure    = "fs.read.failure"

	eventFSListStart    = "fs.list.start"
	eventFSListChildren = "fs.list.children"
	eventFSListSuccess  = "fs.list.success"
	eventFSListFailure  = "fs.list.failure"

	eventFSStatStart   = "fs.stat.start"
	eventFSStatSuccess = "fs.stat.success"
	eventFSStatFailure = "fs.stat.failure"

	eventFSRecursiveStart = "fs.list.recursive.start"
	eventFSRecursiveStep  = "fs.list.recursive.step"
	eventFSRecursiveError = "fs.list.recursive.error"

	eventFSNotImplemented = "fs.not_implemented"
)

func (s *Service) buildLogger(ctx context.Context) logging.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)
}

func (s *Service) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	logger := s.buildLogger(ctx)
	logger = logger.With(logging.String("path", path))

	logger.Error("Mkdir is not implemented",
		logging.String("event", eventFSNotImplemented),
	)
	panic("unimplemented")
}

func (s *Service) Move(ctx context.Context, src string, dst string, opts domainfs.MoveOptions) error {
	correlationID := util.CorrelationIDFromContext(ctx)
	s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("src", src),
		logging.String("dst", dst),
		logging.String("event", eventFSNotImplemented),
	).Error("Move is not implemented")
	panic("unimplemented")
}

func (s *Service) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	logger := s.buildLogger(ctx)
	logger = logger.With(logging.String("path", path))

	item, err := s.resolvePath(ctx, path)
	if err != nil {
		logger.Warn("failed to resolve item path",
			logging.Error(err),
		)
		return nil, err
	}

	if item.IsFolder {
		return nil, errors.New("can't read directories")
	}

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		logger.Error("failed to resolve current drive ID",
			logging.String("event", eventFSResolveFailure),
			logging.Error(err),
		)
		return nil, err
	}

	content, err := s.files.GetFileContents(ctx, driveID, path)
	if err != nil {
		logger.Error("failed to retrieve file contents",
			logging.String("event", eventFSReadFailure),
			logging.Error(err),
		)
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(content)), nil
}

func (s *Service) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	correlationID := util.CorrelationIDFromContext(ctx)
	s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("path", path),
		logging.String("event", eventFSNotImplemented),
	).Error("Remove is not implemented")
	panic("unimplemented")
}

func (s *Service) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	correlationID := util.CorrelationIDFromContext(ctx)
	s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("path", path),
	)

	logger := s.buildLogger(ctx)
	logger = logger.With(logging.String("path", path))

	item, err := s.resolvePath(ctx, path)
	if err != nil {
		logger.Warn("failed to resolve item path",
			logging.Error(err),
		)
		return domainfs.Item{}, err
	}

	if item.IsFolder {
		return domainfs.Item{}, errors.New("can't write directories")
	}

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		logger.Error("failed to resolve current drive ID",
			logging.String("event", eventFSResolveFailure),
			logging.Error(err),
		)
		return domainfs.Item{}, err
	}

	result, err := s.files.WriteFile(ctx, driveID, path, r)
	if err != nil {

	}

	return mapToFSItem(result), nil
}

func (s *Service) resolvePath(ctx context.Context, p string) (*infrafile.DriveItem, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("path", p),
	)

	logger.Debug("resolving filesystem path",
		logging.String("event", eventFSResolveStart),
	)

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		logger.Error("failed to resolve current drive ID",
			logging.String("event", eventFSResolveFailure),
			logging.Error(err),
		)
		return nil, err
	}

	item, err := s.files.ResolveItem(ctx, driveID, p)
	if err != nil {
		logger.Error("failed to resolve drive item",
			logging.String("event", eventFSResolveFailure),
			logging.Error(err),
		)
		return nil, err
	}

	logger.Debug("filesystem path resolved",
		logging.String("event", eventFSResolveSuccess),
		logging.String("drive_id", driveID),
		logging.String("item_name", item.Name),
	)

	return item, nil
}

func (s *Service) listChildren(ctx context.Context, p string, recursive bool) ([]*infrafile.DriveItem, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("path", p),
		logging.Bool("recursive", recursive),
	)

	logger.Debug("listing children",
		logging.String("event", eventFSListChildren),
	)

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		logger.Error("failed to resolve current drive ID",
			logging.String("event", eventFSListFailure),
			logging.Error(err),
		)
		return nil, err
	}

	items, err := s.files.ListChildren(ctx, driveID, p)
	if err != nil {
		logger.Error("failed to list children",
			logging.String("event", eventFSListFailure),
			logging.Error(err),
		)
		return nil, err
	}

	if !recursive {
		logger.Debug("non-recursive list complete",
			logging.Int("count", len(items)),
		)
		return items, nil
	}

	logger.Debug("starting recursive listing",
		logging.String("event", eventFSRecursiveStart),
	)

	all := make([]*infrafile.DriveItem, 0, len(items))
	all = append(all, items...)

	for _, item := range items {
		if item.IsFolder {
			childPath := item.PathWithoutDrive + "/" + item.Name

			logger.Debug("recursing into folder",
				logging.String("event", eventFSRecursiveStep),
				logging.String("child_path", childPath),
			)

			children, err := s.listChildren(ctx, childPath, true)
			if err != nil {
				logger.Error("recursive listing failed",
					logging.String("event", eventFSRecursiveError),
					logging.String("child_path", childPath),
					logging.Error(err),
				)
				return nil, err
			}

			all = append(all, children...)
		}
	}

	logger.Debug("recursive listing complete",
		logging.Int("total_count", len(all)),
	)

	return all, nil
}

func (s *Service) List(ctx context.Context, p string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("path", p),
	)

	logger.Info("listing filesystem items",
		logging.String("event", eventFSListStart),
		logging.Bool("recursive", opts.Recursive),
	)

	item, err := s.resolvePath(ctx, p)
	if err != nil {
		logger.Error("failed to resolve path",
			logging.String("event", eventFSListFailure),
			logging.Error(err),
		)
		return nil, err
	}

	if !item.IsFolder {
		logger.Info("path is a file; returning single item",
			logging.String("event", eventFSListSuccess),
		)
		return []domainfs.Item{mapToFSItem(item)}, nil
	}

	children, err := s.listChildren(ctx, p, opts.Recursive)
	if err != nil {
		logger.Error("failed to list children",
			logging.String("event", eventFSListFailure),
			logging.Error(err),
		)
		return nil, err
	}

	logger.Info("filesystem list complete",
		logging.String("event", eventFSListSuccess),
		logging.Int("count", len(children)),
	)

	return mapToFSItems(children), nil
}

func (s *Service) Stat(ctx context.Context, p string, opts domainfs.StatOptions) (domainfs.Item, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("path", p),
	)

	logger.Info("stat on filesystem item",
		logging.String("event", eventFSStatStart),
	)

	item, err := s.resolvePath(ctx, p)
	if err != nil {
		logger.Error("failed to resolve path",
			logging.String("event", eventFSStatFailure),
			logging.Error(err),
		)
		return domainfs.Item{}, err
	}

	logger.Info("stat successful",
		logging.String("event", eventFSStatSuccess),
		logging.String("item_name", item.Name),
		logging.Bool("is_folder", item.IsFolder),
	)

	return mapToFSItem(item), nil
}
