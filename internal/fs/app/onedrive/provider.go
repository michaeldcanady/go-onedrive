package app

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	logger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domainDrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

var _ domain.Service = (*Provider)(nil)

type Provider struct {
	metadataRepo      domain.MetadataRepository
	contentsRepo      domain.FileContentsRepository
	log               logger.Logger
	driveResolver     domainDrive.DriveResolver
	driveAliasService domainDrive.DriveAliasService
}

func NewProvider(metadataRepo domain.MetadataRepository, contentsRepo domain.FileContentsRepository, driveResolver domainDrive.DriveResolver, driveAliasService domainDrive.DriveAliasService, l logger.Logger) *Provider {
	return &Provider{
		metadataRepo:      metadataRepo,
		contentsRepo:      contentsRepo,
		log:               l,
		driveResolver:     driveResolver,
		driveAliasService: driveAliasService,
	}
}

func (s *Provider) resolveDrive(ctx context.Context, path string) (string, string, error) {
	// If path is "shared-docs:/Folder/domain.txt"
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

func (s *Provider) Get(ctx context.Context, path string) (domain.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))

	log.Debug("retrieving metadata")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domain.Item{}, err
	}

	metadata, err := s.metadataRepo.GetByPath(ctx, driveID, cleanPath, domain.MetadataGetOptions{})
	if err != nil {
		return domain.Item{}, err
	}

	return convertMetadataToItem(metadata), nil
}

// List implements [domain.Service].
func (s *Provider) List(ctx context.Context, path string, opts domain.ListOptions) ([]domain.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return nil, err
	}

	metadataOpts := domain.MetadataListOptions{
		NoCache: opts.SkipCache,
		NoStore: opts.NoCache,
	}

	var results []domain.Item

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
			if m.Type == domain.ItemTypeFolder && opts.Recursive {
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

// Mkdir implements [domain.Service].
func (s *Provider) Mkdir(ctx context.Context, path string, opts domain.MKDirOptions) error {
	log := s.buildLogger(ctx)

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		log.Warn("failed to resolve driveID", logger.Error(err))
		return err
	}

	parentPath, name := filepath.Split(cleanPath)

	request := domain.MetadataCreateRequest{
		Name: name,
		Type: domain.ItemTypeFolder,
	}

	_, err = s.metadataRepo.CreateByPath(ctx, driveID, parentPath, request, domain.MetadataCreateOptions{
		CreateParents: opts.Parents,
	})

	return err
}

func (s *Provider) Copy(ctx context.Context, src, dst string, opts domain.CopyOptions) error {
	log := s.buildLogger(ctx).With(
		logger.String("src", src),
		logger.String("dst", dst),
	)
	log.Debug("Copy: starting")

	r, err := s.ReadFile(ctx, src, domain.ReadOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = s.WriteFile(ctx, dst, r, domain.WriteOptions{
		Overwrite: opts.Overwrite,
	})
	return err
}

// Move implements [domain.Service].
func (s *Provider) Move(ctx context.Context, src string, dst string, opts domain.MoveOptions) error {
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

	request := domain.MetadataUpdateRequest{
		Name:       name,
		ParentPath: parentPath,
	}

	_, err = s.metadataRepo.UpdateByPath(ctx, srcDriveID, cleanSrc, request, domain.MetadataUpdateOptions{})
	return err
}

// ReadFile implements [domain.Service].
func (s *Provider) ReadFile(ctx context.Context, path string, opts domain.ReadOptions) (io.ReadCloser, error) {
	s.buildLogger(ctx).With(logger.String("path", path))

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return nil, err
	}

	return s.contentsRepo.Download(ctx, driveID, cleanPath, domain.DownloadOptions{})
}

// Remove implements [domain.Service].
func (s *Provider) Remove(ctx context.Context, path string, opts domain.RemoveOptions) error {
	log := s.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("Remove: deleting file")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return err
	}

	return s.metadataRepo.DeleteByPath(ctx, driveID, cleanPath, domain.MetadataDeleteOptions{})
}

// Stat implements [domain.Service].
func (s *Provider) Stat(ctx context.Context, path string, opts domain.StatOptions) (domain.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("Stat: retrieving metadata")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domain.Item{}, err
	}

	metadata, err := s.metadataRepo.GetByPath(ctx, driveID, cleanPath, domain.MetadataGetOptions{})
	if err != nil {
		return domain.Item{}, err
	}

	return convertMetadataToItem(metadata), nil
}

// WriteFile implements [domain.Service].
func (s *Provider) WriteFile(ctx context.Context, path string, r io.Reader, opts domain.WriteOptions) (domain.Item, error) {
	s.buildLogger(ctx).With(logger.String("path", path))

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domain.Item{}, err
	}

	metadata, err := s.contentsRepo.Upload(ctx, driveID, cleanPath, r, domain.UploadOptions{
		Overwrite: opts.Overwrite,
		IfMatch:   opts.IfMatch,
	})

	return convertMetadataToItem(metadata), err
}

func (s *Provider) Upload(ctx context.Context, src, dst string, opts domain.UploadOptions) (domain.Item, error) {
	_ = s.buildLogger(ctx)

	file, err := os.Open(src)
	if err != nil {
		return domain.Item{}, err
	}
	defer file.Close()

	return s.WriteFile(ctx, dst, file, domain.WriteOptions{
		Overwrite: opts.Overwrite,
	})
}

// Touch implements [domain.Service].
func (s *Provider) Touch(ctx context.Context, path string, opts domain.TouchOptions) (domain.Item, error) {
	log := s.buildLogger(ctx).With(logger.String("path", path))
	log.Debug("Touch: starting")

	driveID, cleanPath, err := s.resolveDrive(ctx, path)
	if err != nil {
		return domain.Item{}, err
	}

	parentPath, name := filepath.Split(cleanPath)

	request := domain.MetadataCreateRequest{
		Name: name,
		Type: domain.ItemTypeFile,
	}

	metadata, err := s.metadataRepo.CreateByPath(ctx, driveID, parentPath, request, domain.MetadataCreateOptions{})
	if err != nil {
		return domain.Item{}, err
	}

	return convertMetadataToItem(metadata), nil
}
