package fs

import (
	"context"
	"io"

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

func (s *Service2) buildLogger(ctx context.Context) logging.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)
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

	// -----------------------------
	// NON‑RECURSIVE LIST
	// -----------------------------
	if !opts.Recursive {
		metadatas, err := s.metadataRepo.ListByPath(ctx, driveID, path, metadataOpts)
		if err != nil {
			return nil, err
		}

		items := make([]domainfs.Item, len(metadatas))
		for i, m := range metadatas {
			items[i] = convertMetadataToItem(m)
		}
		return items, nil
	}

	// -----------------------------
	// RECURSIVE LIST
	// -----------------------------
	var results []domainfs.Item

	// Track visited paths to avoid cycles (defensive)
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
			if m.Type == file.ItemTypeFolder {
				childPath := currentPath + "/" + m.Name
				if err := walk(childPath); err != nil {
					return err
				}
			}
		}

		return nil
	}

	if err := walk(path); err != nil {
		logger.Error("recursive list failed", logging.Error(err))
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
