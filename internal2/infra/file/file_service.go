package file

import (
	"context"
	"errors"
	"fmt"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

type Service2 struct {
	graph  clienter
	logger logging.Logger
	cache  domaincache.CacheService
}

func New2(graph clienter, logger logging.Logger, cache domaincache.CacheService) *Service2 {
	s := &Service2{
		graph:  graph,
		cache:  cache,
		logger: logger,
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
		if isOdataErr := errors.As(err, &odataErr); isOdataErr {
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
		s.logger.Error("unexpected error type returned from getDriveRoot", logging.Error(err))
		return nil, err
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

	var headers *abstractions.RequestHeaders
	if cached.ETag != "" {
		s.logger.Info("cached ETag found", logging.String("etag", cached.ETag))
		headers = abstractions.NewRequestHeaders()
		headers.Add("If-None-Match", fmt.Sprintf("\"%s\"", cached.ETag))
	} else {
		s.logger.Info("no cached ETag found for path", logging.String("path", normalizedPath))
	}

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

	s.logger.Info("retrieved drive children successfully", logging.String("path", normalized), logging.Int("count", len(items)))
	s.logger.Debug("children_count", logging.Any("count", len(items)))

	// Cache updated children
	if etag := driveItem.GetETag(); etag != nil && *etag != "" {
		s.cache.SetDrive(ctx, cacheKey, domaincache.CachedChildren{
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
	urlTemplate := rootURITemplate2
	subs := make(stduritemplate.Substitutions)
	subs["baseurl"] = baseURL
	subs["drive_id"] = driveID

	if normalizedPath != "" {
		urlTemplate = rootRelativeURITemplate2
		subs["path"] = normalizedPath
	}

	uri, _ := stduritemplate.Expand(urlTemplate, subs)

	return drives.NewItemRootRequestBuilder(uri, client.RequestAdapter)
}

func (s *Service2) childrenBuilder(client *msgraphsdkgo.GraphServiceClient, driveID, normalizedPath string) *drives.ItemItemsRequestBuilder {
	urlTemplate := rootChildrenURITemplate2
	subs := make(stduritemplate.Substitutions)
	subs["baseurl"] = baseURL
	subs["drive_id"] = driveID

	if normalizedPath != "" {
		urlTemplate = rootRelativeChildrenURITemplate2
		subs["path"] = normalizedPath
	}

	uri, _ := stduritemplate.Expand(urlTemplate, subs)

	return drives.NewItemItemsRequestBuilder(uri, client.RequestAdapter)
}
