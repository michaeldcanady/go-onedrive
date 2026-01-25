package driveservice

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"path"

	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	rootChildrenURITemplate         = "https://graph.microsoft.com/v1.0/drives/%s/root/children"
	rootRelativeChildrenURITemplate = "https://graph.microsoft.com/v1.0/drives/%s/root:%s:/children"
	rootRelativeURITemplate         = "https://graph.microsoft.com/v1.0/drives/%s/root:%s:"
	rootURITemplate                 = "https://graph.microsoft.com/v1.0/drives/%s/root"
)

type Service struct {
	graph     clienter
	publisher event.Publisher
	logger    logging.Logger
	cache     CacheService
}

func New(graph clienter, publisher event.Publisher, logger logging.Logger, cache CacheService) *Service {
	return &Service{
		graph:     graph,
		publisher: publisher,
		logger:    logger,
		cache:     cache,
	}
}

// normalizePath ensures paths like "Documents", "/Documents", "Documents/" all become "/Documents"
func normalizePath(p string) string {
	if p == "" || p == "/" || p == "." {
		return ""
	}
	return path.Clean("/" + p)
}

func (s *Service) cacheKey(driveID, normalizedPath string) string {
	return driveID + ":" + normalizedPath
}

// getUserDriveID retrieves the user's default drive ID
func (s *Service) getUserDriveID(ctx context.Context) (string, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return "", errors.Join(errors.New("unable to instantiate client"), err)
	}

	drive, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		s.logger.Error("failed to retrieve user drive", logging.Any("error", err))
		return "", fmt.Errorf("failed to retrieve user drive: %w", err)
	}

	id := drive.GetId()
	if id == nil || *id == "" {
		s.logger.Error("user drive ID is empty")
		return "", fmt.Errorf("user drive ID is empty")
	}

	s.logger.Info("retrieved user drive ID")
	s.logger.Debug("drive_id", logging.String("id", *id))

	return *id, nil
}

// getDriveRoot fetches the DriveItem for the given path, using ETag caching.
func (s *Service) getDriveRoot(ctx context.Context, driveID, normalizedPath string) (models.Driveable, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return nil, err
	}

	// Load cached ETag
	cacheKey := s.cacheKey(driveID, normalizedPath)
	cached, err := s.cache.GetDrive(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	headers := abstractions.NewRequestHeaders()
	headers.Add("If-None-Match", fmt.Sprintf("\"%s\"", cached.ETag))

	config := &drives.DriveItemRequestBuilderGetRequestConfiguration{
		Headers: headers,
	}

	return s.driveItemBuilder(client, driveID, normalizedPath).Get(ctx, config)
}

// getChildren retrieves folder children, using ETag caching and event publishing.
func (s *Service) getChildren(ctx context.Context, folderPath string) (models.DriveItemCollectionResponseable, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return nil, err
	}

	normalized := normalizePath(folderPath)
	s.logger.Debug("normalized folder path", logging.String("path", normalized))

	driveID, err := s.getUserDriveID(ctx)
	if err != nil {
		return nil, err
	}

	// Check root/folder metadata with ETag
	driveItem, err := s.getDriveRoot(ctx, driveID, normalized)
	if err != nil {
		return nil, err
	}

	cacheKey := s.cacheKey(driveID, normalized)

	// 304 Not Modified → use cached children
	if driveItem == nil {
		cached, err := s.cache.GetDrive(ctx, cacheKey)
		if err != nil {
			return nil, err
		}
		return cached.Items, nil
	}

	// Build children URL
	resp, err := s.childrenBuilder(client, driveID, normalized).Get(ctx, nil)
	if err != nil {
		s.logger.Error("failed to retrieve children", logging.String("path", normalized), logging.Any("error", err))
		return nil, fmt.Errorf("failed to retrieve children: %w", err)
	}

	items := resp.GetValue()

	// Publish event
	if s.publisher != nil {
		s.logger.Debug("publishing drive.children.loaded event")
		if err := s.publisher.Publish(ctx, newDriveChildrenLoadedEvent(normalized, items)); err != nil {
			s.logger.Error("failed to publish drive.children.loaded event", logging.Any("error", err))
		}
	} else {
		s.logger.Warn("no event publisher configured; skipping drive.children.loaded event")
	}

	s.logger.Info("retrieved drive children successfully", logging.String("path", normalized))
	s.logger.Debug("children_count", logging.Any("count", len(items)))

	// Cache updated children
	if etag := driveItem.GetETag(); etag != nil && *etag != "" {
		s.cache.SetDrive(ctx, cacheKey, CachedChildren{
			ETag:  *etag,
			Items: resp,
		})
	} else {
		s.logger.Warn("drive etag unavailable; response not cached")
	}

	return resp, nil
}

// ChildrenIterator returns an iterator over DriveItem children.
func (s *Service) ChildrenIterator(ctx context.Context, folderPath string) iter.Seq2[models.DriveItemable, error] {
	resp, err := s.getChildren(ctx, folderPath)
	if err != nil {
		s.logger.Error("unable to retrieve children", logging.String("path", folderPath), logging.Any("error", err))
		return func(yield func(models.DriveItemable, error) bool) {
			yield(nil, fmt.Errorf("unable to retrieve children: %w", err))
		}
	}

	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return func(yield func(models.DriveItemable, error) bool) {
			yield(nil, err)
		}
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.DriveItemable](
		resp,
		client.GetAdapter(),
		models.CreateDriveItemFromDiscriminatorValue,
	)
	if err != nil {
		s.logger.Error("unable to create page iterator", logging.Any("error", err))
		return func(yield func(models.DriveItemable, error) bool) {
			yield(nil, fmt.Errorf("unable to create page iterator: %w", err))
		}
	}

	s.logger.Info("iterating drive children", logging.String("path", folderPath))

	return func(yield func(models.DriveItemable, error) bool) {
		if iterErr := pageIterator.Iterate(ctx, func(item models.DriveItemable) bool {
			return yield(item, nil)
		}); iterErr != nil {
			s.logger.Error("error during children iteration", logging.Any("error", iterErr))
			yield(nil, iterErr)
		}
	}
}

func (s *Service) driveItemBuilder(client *msgraphsdkgo.GraphServiceClient, driveID, normalizedPath string) *drives.DriveItemRequestBuilder {
	if normalizedPath == "" {
		return drives.NewDriveItemRequestBuilder(fmt.Sprintf(rootURITemplate, driveID), client.RequestAdapter)
	}
	return drives.NewDriveItemRequestBuilder(fmt.Sprintf(rootRelativeURITemplate, driveID, normalizedPath), client.RequestAdapter)
}

func (s *Service) childrenBuilder(client *msgraphsdkgo.GraphServiceClient, driveID, normalizedPath string) *drives.ItemItemsRequestBuilder {
	if normalizedPath != "" {
		return drives.NewItemItemsRequestBuilder(fmt.Sprintf(rootRelativeChildrenURITemplate, driveID, normalizedPath), client.RequestAdapter)
	}

	return drives.NewItemItemsRequestBuilder(fmt.Sprintf(rootChildrenURITemplate, driveID), client.RequestAdapter)
}
