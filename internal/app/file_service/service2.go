package driveservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

type Service2 struct {
	graph     clienter
	publisher event.Publisher
	logger    logging.Logger
	cache     CacheService
}

func New2(graph clienter, publisher event.Publisher, logger logging.Logger, cache CacheService) *Service2 {
	s := &Service2{
		graph:     graph,
		cache:     cache,
		logger:    logger,
		publisher: publisher,
	}

	return s
}

func (s *Service2) cacheKey(driveID, normalizedPath string) string {
	return driveID + ":" + normalizedPath
}

func (s *Service2) ResolveItem(ctx context.Context, driveID, path string) (*DriveItem, error) {
	// Internally calls your existing getDriveRoot logic
	item, err := s.getDriveRoot(ctx, driveID, normalizePath(path))
	if err != nil {
		var odataErr odataerrors.ODataErrorable
		errors.As(err, &odataErr)
		mainErr := odataErr.GetErrorEscaped()
		errDetails := mainErr.GetDetails()
		details := make([]logging.Field, len(errDetails)+1)
		details[0] = logging.String("msg", *mainErr.GetMessage())
		for i, errDetail := range errDetails {
			i = i + 1
			detail := logging.String(fmt.Sprintf("detail[%d]", i), *errDetail.GetMessage())
			details[i] = detail
		}

		s.logger.Error("failed to get drive root",
			details...,
		)
		return nil, mapGraphError(err)
	}

	if item == nil {
		return nil, &DomainError{
			Kind:    ErrNotFound,
			DriveID: driveID,
			Path:    path,
			Err:     errors.New("item not found"),
		}
	}

	return toDomainItem(driveID, item), nil
}

// getDriveRoot fetches the DriveItem for the given path, using ETag caching.
func (s *Service2) getDriveRoot(ctx context.Context, driveID, normalizedPath string) (models.DriveItemable, error) {
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

	config := &drives.ItemRootRequestBuilderGetRequestConfiguration{
		Headers: headers,
	}

	return s.driveItemBuilder(client, driveID, normalizedPath).Get(ctx, config)
}

// getChildren retrieves folder children, using ETag caching and event publishing.
func (s *Service2) getChildren(ctx context.Context, driveID, folderPath string) (models.DriveItemCollectionResponseable, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return nil, err
	}

	normalized := normalizePath(folderPath)
	s.logger.Debug("normalized folder path", logging.String("path", normalized))

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

func (s *Service2) ListChildren(ctx context.Context, driveID, path string) ([]*DriveItem, error) {
	resp, err := s.getChildren(ctx, driveID, path)
	if err != nil {
		return nil, mapGraphError(err)
	}

	values := resp.GetValue()
	out := make([]*DriveItem, 0, len(values))

	for _, it := range values {
		out = append(out, toDomainItem(driveID, it))
	}

	return out, nil
}

func (s *Service2) driveItemBuilder(client *msgraphsdkgo.GraphServiceClient, driveID, normalizedPath string) *drives.ItemRootRequestBuilder {
	if normalizedPath == "" {
		return drives.NewItemRootRequestBuilder(fmt.Sprintf(rootURITemplate, driveID), client.RequestAdapter)
	}
	return drives.NewItemRootRequestBuilder(fmt.Sprintf(rootRelativeURITemplate, driveID, normalizedPath), client.RequestAdapter)
}

func (s *Service2) childrenBuilder(client *msgraphsdkgo.GraphServiceClient, driveID, normalizedPath string) *drives.ItemItemsRequestBuilder {
	if normalizedPath != "" {
		return drives.NewItemItemsRequestBuilder(fmt.Sprintf(rootRelativeChildrenURITemplate, driveID, normalizedPath), client.RequestAdapter)
	}

	return drives.NewItemItemsRequestBuilder(fmt.Sprintf(rootChildrenURITemplate, driveID), client.RequestAdapter)
}
