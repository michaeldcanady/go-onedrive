package fs

import (
	"context"
	"io"
	"os"

	domainDrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

var _ domainfs.Service = (*Service2)(nil)

type Service2 struct {
	metadataRepo  file.MetadataRepository
	contentsRepo  file.FileContentsRepository
	logger        logging.Logger
	driveResolver domainDrive.DriveResolver
}

func NewService2(metadataRepo file.MetadataRepository, contentsRepo file.FileContentsRepository, driveResolver domainDrive.DriveResolver, logger logging.Logger) *Service2 {
	return &Service2{
		metadataRepo:  metadataRepo,
		contentsRepo:  contentsRepo,
		logger:        logger,
		driveResolver: driveResolver,
	}
}

func (s *Service2) buildLogger(ctx context.Context) logging.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)
}

func (s *Service2) Get(
	ctx context.Context,
	path string,
) (domainfs.Item, error) {
	logger := s.buildLogger(ctx).With(logging.String("path", path))

	logger.Debug("retrieving metadata", logging.String("event", eventFSGetStart))

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		return domainfs.Item{}, err
	}

	metadata, err := s.metadataRepo.GetByPath(ctx, driveID, path, file.MetadataGetOptions{})
	if err != nil {
		return domainfs.Item{}, err
	}

	return convertMetadataToItem(metadata), nil
}

// List implements [fs.Service].
func (s *Service2) List(
	ctx context.Context,
	path string,
	opts domainfs.ListOptions,
) ([]domainfs.Item, error) {
	logger := s.buildLogger(ctx).With(logging.String("path", path))

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		return nil, err
	}

	metadataOpts := file.MetadataListOptions{
		NoCache: opts.SkipCache,
		NoStore: opts.NoCache,
	}

	var results []domainfs.Item

	visited := map[string]bool{}

	var walk func(string) error
	walk = func(currentPath string) error {
		if visited[currentPath] {
			return nil
		}
		visited[currentPath] = true

		metadatas, err := s.metadataRepo.ListByPath(ctx, driveID, currentPath, metadataOpts)
		if err != nil {
			return err
		}

		for _, m := range metadatas {
			item := convertMetadataToItem(m)
			results = append(results, item)

			// Recurse only into folders
			if m.Type == file.ItemTypeFolder && opts.Recursive {
				childPath := currentPath + "/" + m.Name
				if err := walk(childPath); err != nil {
					return err
				}
			}
		}

		return nil
	}

	if err := walk(path); err != nil {
		logger.Error("listing failed", logging.Error(err))
		return nil, err
	}

	return results, nil
}

// Mkdir implements [fs.Service].
func (s *Service2) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	panic("unimplemented")
}

// Move implements [fs.Service].
func (s *Service2) Move(ctx context.Context, src string, dst string, opts domainfs.MoveOptions) error {
	panic("unimplemented")
}

// ReadFile implements [fs.Service].
func (s *Service2) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	logger := s.buildLogger(ctx)
	logger = logger.With(logging.String("path", path))

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		return nil, err
	}

	return s.contentsRepo.Download(ctx, driveID, path, file.DownloadOptions{})
}

// Remove implements [fs.Service].
func (s *Service2) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	panic("unimplemented")
}

// Stat implements [fs.Service].
func (s *Service2) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	panic("unimplemented")
}

// WriteFile implements [fs.Service].
func (s *Service2) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	logger := s.buildLogger(ctx)
	logger = logger.With(logging.String("path", path))

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		return domainfs.Item{}, err
	}

	metadata, err := s.contentsRepo.Upload(ctx, driveID, path, r, file.UploadOptions{
		Force: opts.Overwrite,
	})

	return convertMetadataToItem(metadata), err
}

func (s *Service2) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	_ = s.buildLogger(ctx)

	file, err := os.Open(src)
	if err != nil {
		return domainfs.Item{}, err
	}
	defer file.Close()

	return s.WriteFile(ctx, dst, file, domainfs.WriteOptions{
		Overwrite: opts.Overwrite,
	})
}
