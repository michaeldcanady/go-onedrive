package onedrive

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainDrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

var _ domainfs.Service = (*Provider)(nil)

type Provider struct {
	metadataRepo      file.MetadataRepository
	contentsRepo      file.FileContentsRepository
	log               logger.Logger
	driveResolver     domainDrive.DriveResolver
	driveAliasService domainDrive.DriveAliasService
}

func NewProvider(metadataRepo file.MetadataRepository, contentsRepo file.FileContentsRepository, driveResolver domainDrive.DriveResolver, driveAliasService domainDrive.DriveAliasService, l logger.Logger) *Provider {
	return &Provider{
		metadataRepo:      metadataRepo,
		contentsRepo:      contentsRepo,
		log:               l,
		driveResolver:     driveResolver,
		driveAliasService: driveAliasService,
	}
}

func (s *Provider) resolveDrive(ctx context.Context, path string) (string, string, error) {
	// If path is "shared-docs:/Folder/File.txt"
	if !strings.HasPrefix(path, "/") && strings.Contains(path, ":") {
		driveAlias, cleanPath, _ := strings.Cut(path, ":")
		driveID, err := s.driveAliasService.Resolve(ctx, driveAlias)
		return driveID, cleanPath, err
	}

	// Default to the current selected drive
	driveID, err := s.driveResolver.CurrentDriveID(ctx)
	return driveID, path, err
}

func (s *Provider) buildLogger(ctx context.Context) logger.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)
}

func (s *Provider) Get(ctx context.Context, path string) (domainfs.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))

	log.Debug("retrieving metadata")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}

	metadata, err := s.metadataRepo.GetByPath(ctx, driveID, cleanPath, file.MetadataGetOptions{})
	if err != nil {
		return domainfs.Item{}, err
	}

	return convertMetadataToItem(metadata), nil
}

// List implements [fs.Service].
func (s *Provider) List(ctx context.Context, path string, opts domainfs.ListOptions) ([]domainfs.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
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

	if err := walk(cleanPath); err != nil {
		log.Error("listing failed", logger.Error(err))
		return nil, err
	}

	return results, nil
}

// Mkdir implements [fs.Service].
func (s *Provider) Mkdir(ctx context.Context, path string, opts domainfs.MKDirOptions) error {
	log := s.buildLogger(ctx)

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		log.Warn("failed to resolve driveID", logger.Error(err))
		return err
	}

	parentPath, name := filepath.Split(cleanPath)

	request := file.MetadataCreateRequest{
		Name: name,
		Type: file.ItemTypeFolder,
	}

	_, err = s.metadataRepo.CreateByPath(ctx, driveID, parentPath, request, file.MetadataCreateOptions{
		CreateParents: opts.Parents,
	})

	return err
}

func (s *Provider) Copy(ctx context.Context, src, dst string, opts domainfs.CopyOptions) error {
	log := s.buildLogger(ctx).With(
		logger.String("src", src),
		logger.String("dst", dst),
	)
	log.Debug("Copy: starting")

	srcItem, err := s.Stat(ctx, src, domainfs.StatOptions{})
	if err != nil {
		return err
	}

	if srcItem.Type == domainfs.ItemTypeFolder {
		if !opts.Recursive {
			return errors.New("source is a directory, use recursive flag")
		}

		if err := s.Mkdir(ctx, dst, domainfs.MKDirOptions{Parents: true}); err != nil {
			return err
		}

		children, err := s.List(ctx, src, domainfs.ListOptions{Recursive: false})
		if err != nil {
			return err
		}

		for _, child := range children {
			childSrc := src + "/" + child.Name
			childDst := dst + "/" + child.Name
			if err := s.Copy(ctx, childSrc, childDst, opts); err != nil {
				return err
			}
		}
		return nil
	}

	r, err := s.ReadFile(ctx, src, domainfs.ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = s.WriteFile(ctx, dst, r, domainfs.WriteOptions{
		Overwrite: opts.Overwrite,
	})
	return err
}

// Move implements [fs.Service].
func (s *Provider) Move(ctx context.Context, src string, dst string, opts domainfs.MoveOptions) error {
	log := s.buildLogger(ctx).With(
		logger.String("src", src),
		logger.String("dst", dst),
	)
	log.Debug("Move: starting")

	srcDriveID, cleanSrc, err := s.resolveDrive(ctx, src)
	if err != nil {
		return err
	}

	dstDriveID, cleanDst, err := s.resolveDrive(ctx, dst)
	if err != nil {
		return err
	}

	if srcDriveID != dstDriveID {
		// Cross-drive move: Copy + Delete (not implemented here yet)
		return errors.New("cross-drive move not supported yet")
	}

	parentPath, name := filepath.Split(cleanDst)

	request := file.MetadataUpdateRequest{
		Name:       name,
		ParentPath: parentPath,
	}

	_, err = s.metadataRepo.UpdateByPath(ctx, srcDriveID, cleanSrc, request, file.MetadataUpdateOptions{})
	return err
}

// ReadFile implements [fs.Service].
func (s *Provider) ReadFile(ctx context.Context, path string, opts domainfs.ReadOptions) (io.ReadCloser, error) {
	s.buildLogger(ctx).With(logger.String("path", path))

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return nil, err
	}

	return s.contentsRepo.Download(ctx, driveID, cleanPath, file.DownloadOptions{})
}

// Remove implements [fs.Service].
func (s *Provider) Remove(ctx context.Context, path string, opts domainfs.RemoveOptions) error {
	log := s.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("Remove: deleting file")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return err
	}

	return s.metadataRepo.DeleteByPath(ctx, driveID, cleanPath, file.MetadataDeleteOptions{
		Permanent: opts.Permanent,
	})
}

// Stat implements [fs.Service].
func (s *Provider) Stat(ctx context.Context, path string, opts domainfs.StatOptions) (domainfs.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("Stat: retrieving metadata")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}

	metadata, err := s.metadataRepo.GetByPath(ctx, driveID, cleanPath, file.MetadataGetOptions{})
	if err != nil {
		return domainfs.Item{}, err
	}

	return convertMetadataToItem(metadata), nil
}

// WriteFile implements [fs.Service].
func (s *Provider) WriteFile(ctx context.Context, path string, r io.Reader, opts domainfs.WriteOptions) (domainfs.Item, error) {
	s.buildLogger(ctx).With(logger.String("path", path))

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}

	metadata, err := s.contentsRepo.Upload(ctx, driveID, cleanPath, r, file.UploadOptions{
		Force:   opts.Overwrite,
		IfMatch: opts.IfMatch,
	})

	return convertMetadataToItem(metadata), err
}

func (s *Provider) Upload(ctx context.Context, src, dst string, opts domainfs.UploadOptions) (domainfs.Item, error) {
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

// Touch implements [fs.Service].
func (s *Provider) Touch(ctx context.Context, path string, opts domainfs.TouchOptions) (domainfs.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("Touch: starting")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domainfs.Item{}, err
	}

	parentPath, name := filepath.Split(cleanPath)

	request := file.MetadataCreateRequest{
		Name: name,
		Type: file.ItemTypeFile,
	}

	metadata, err := s.metadataRepo.CreateByPath(ctx, driveID, parentPath, request, file.MetadataCreateOptions{})
	if err != nil {
		return domainfs.Item{}, err
	}

	return convertMetadataToItem(metadata), nil
}
