package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"

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
			if m == nil {
				continue
			}
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
	logger := s.buildLogger(ctx)

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		logger.Warn("failed to resolve driveID", logging.Error(err))
		return err
	}

	parentPath, name := filepath.Split(path)

	request := file.MetadataCreateRequest{
		Name: name,
		Type: file.ItemTypeFolder,
	}

	_, err = s.metadataRepo.CreateByPath(ctx, driveID, parentPath, request, file.MetadataCreateOptions{
		CreateParents: opts.Parents,
	})

	return err
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
	logger := s.buildLogger(ctx).With(logging.String("path", path))
	logger.Debug("Stat: retrieving metadata")

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

// WriteFile implements [fs.Service].
func (s *Service2) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	logger := s.buildLogger(ctx)
	logger = logger.With(logging.String("path", path))

	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	if err != nil {
		return domainfs.Item{}, err
	}

	metadata, err := s.contentsRepo.Upload(ctx, driveID, path, r, file.UploadOptions{
		Force:   opts.Overwrite,
		IfMatch: opts.IfMatch,
	})

	return convertMetadataToItem(metadata), err
}

func (s *Service2) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	logger := s.buildLogger(ctx).With(
		logging.String("src", src),
		logging.String("dst", dst),
	)

	info, err := os.Stat(src)
	if err != nil {
		logger.Error("failed to stat source", logging.Error(err))
		return domainfs.Item{}, err
	}

	if !info.IsDir() {
		return s.uploadFile(ctx, src, dst, opts)
	}

	// It's a directory
	return s.uploadDir(ctx, src, dst, opts)
}

func (s *Service2) uploadFile(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	f, err := os.Open(src)
	if err != nil {
		return domainfs.Item{}, err
	}
	defer f.Close()

	return s.WriteFile(ctx, dst, f, domainfs.WriteOptions{
		Overwrite: opts.Overwrite,
	})
}

func (s *Service2) uploadDir(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
	logger := s.buildLogger(ctx).With(
		logging.String("src", src),
		logging.String("dst", dst),
	)

	// 1. Ensure the directory exists in OneDrive
	if err := s.Mkdir(ctx, dst, domainfs.MKDirOptions{Parents: true}); err != nil {
		// If it's already a folder, Mkdir might fail depending on implementation.
		// For now we assume Mkdir handles "already exists" or we ignore it.
		logger.Debug("Mkdir possibly failed or directory already exists", logging.Error(err))
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return domainfs.Item{}, err
	}

	for _, entry := range entries {
		entrySrc := filepath.Join(src, entry.Name())
		entryDst := dst + "/" + entry.Name()

		if entry.IsDir() {
			if opts.Recursive {
				if _, err := s.uploadDir(ctx, entrySrc, entryDst, opts); err != nil {
					return domainfs.Item{}, err
				}
			}
			continue
		}

		// It's a file
		_, err := s.uploadFile(ctx, entrySrc, entryDst, opts)
		if err != nil {
			return domainfs.Item{}, err
		}
	}

	// For a directory upload, returning the item for the directory itself might be more appropriate.
	return s.Get(ctx, dst)
}
