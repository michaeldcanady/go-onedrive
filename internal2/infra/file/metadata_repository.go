package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

// MetadataRepository provides methods for fetching and listing drive item
// metadata from OneDrive. It supports multi-level caching (individual items
// and folder listings) and handles path normalization.
type MetadataRepository struct {
	client               abstractions.RequestAdapter
	metadataCache        MetadataCache
	metadataListingCache ListingCache
	pathIDCache          PathIDCache
	logger               logging.Logger
}

// NewMetadataRepository initializes a new MetadataRepository with the provided
// request adapter and cache implementations.
func NewMetadataRepository(client abstractions.RequestAdapter, metadataCache MetadataCache, metadataListingCache ListingCache, pathIDCache PathIDCache, logger logging.Logger) *MetadataRepository {
	return &MetadataRepository{
		client:               client,
		metadataCache:        metadataCache,
		metadataListingCache: metadataListingCache,
		pathIDCache:          pathIDCache,
		logger:               logger,
	}
}

func (r *MetadataRepository) getByPath(ctx context.Context, driveID, path string, opts file.MetadataGetOptions) (*file.Metadata, bool, error) {
	config := &drives.ItemRootRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	path = normalizePath(path)
	r.logger.Debug("getByPath: starting retrieval", logging.String("path", path), logging.Bool("noCache", opts.NoCache))

	// 1. Cache lookup unless disabled
	var cached *file.Metadata
	if !opts.NoCache {
		// Try path-to-ID cache first if we have it
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			r.logger.Debug("getByPath: path-to-ID hit", logging.String("path", path), logging.String("id", id))
			if m, ok := r.metadataCache.Get(ctx, id); ok && m != nil {
				r.logger.Debug("getByPath: metadata cache hit (by id)", logging.String("id", id))
				cached = m
				if !opts.Force {
					config.Headers.Add("If-None-Match", m.ETag)
				}
			} else {
				r.logger.Debug("getByPath: metadata cache miss (by id)", logging.String("id", id))
			}
		} else {
			r.logger.Debug("getByPath: path-to-ID miss", logging.String("path", path))
			// Fallback to path lookup in metadata cache
			if m, ok := r.metadataCache.Get(ctx, path); ok && m != nil {
				r.logger.Debug("getByPath: metadata cache hit (by path)", logging.String("path", path))
				cached = m
				if !opts.Force {
					config.Headers.Add("If-None-Match", m.ETag)
				}
			}
		}
	}

	// 2. Fetch from OneDrive
	r.logger.Debug("getByPath: requesting from OneDrive", logging.String("path", path))
	item, err := r.driveItemBuilder(r.client, driveID, path).Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		r.logger.Error("getByPath: request failed", logging.String("path", path), logging.Error(err))
		return nil, false, err
	}

	// 3. 304 Not Modified → return cached
	if item == nil {
		r.logger.Info("getByPath: 304 Not Modified", logging.String("path", path))
		return cached, false, nil
	}

	// 4. Fresh metadata
	metadata := mapItemToMetadata(item)
	r.logger.Info("getByPath: fresh metadata received", logging.String("path", path), logging.String("id", metadata.ID))

	if !opts.NoStore {
		r.logger.Debug("getByPath: updating cache", logging.String("id", metadata.ID))
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			r.logger.Warn("getByPath: failed to update metadata cache", logging.String("id", metadata.ID), logging.Error(err))
		}
		if err := r.pathIDCache.Put(ctx, path, metadata.ID); err != nil {
			r.logger.Warn("getByPath: failed to update path-to-ID cache", logging.String("path", path), logging.Error(err))
		}
	}

	return metadata, true, nil
}

// GetByPath retrieves metadata for a single drive item at the specified path.
func (r *MetadataRepository) GetByPath(ctx context.Context, driveID, path string, opts file.MetadataGetOptions) (*file.Metadata, error) {
	metadata, _, err := r.getByPath(ctx, driveID, path, opts)
	return metadata, err
}

// ListByPath returns metadata for all children of the folder at the specified path.
// It utilizes both the listing cache and the individual item metadata cache.
func (r *MetadataRepository) ListByPath(ctx context.Context, driveID, path string, opts file.MetadataListOptions) ([]*file.Metadata, error) {
	r.logger.Debug("ListByPath: starting retrieval", logging.String("path", path))
	parent, updated, err := r.getByPath(ctx, driveID, path, file.MetadataGetOptions{
		NoCache: opts.NoCache,
		NoStore: opts.NoStore,
		Force:   opts.Force,
	})
	if err != nil {
		r.logger.Error("ListByPath: failed to get parent metadata", logging.String("path", path), logging.Error(err))
		return nil, err
	}

	if parent == nil {
		r.logger.Warn("ListByPath: parent metadata is nil", logging.String("path", path))
		return nil, ErrNotFound
	}

	// 2. If parent not updated and listing cache allowed → return cached listing
	if !updated && !opts.NoCache {
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			r.logger.Debug("ListByPath: listing cache hit", logging.String("path", path))
			children := make([]*file.Metadata, 0, len(listing.ChildIDs))
			allFound := true
			for _, id := range listing.ChildIDs {
				m, ok := r.metadataCache.Get(ctx, id)
				if !ok || m == nil {
					r.logger.Warn("ListByPath: inconsistent cache, child metadata missing", logging.String("id", id))
					allFound = false
					break
				}
				children = append(children, m)
			}
			if allFound {
				r.logger.Info("ListByPath: returning cached listing", logging.String("path", path), logging.Int("count", len(children)))
				return children, nil
			}
			// Inconsistent cache: listing exists but some items are missing.
			// Fall through to fetch from Graph.
			r.logger.Info("ListByPath: invalidating inconsistent listing cache", logging.String("path", path))
			if err := r.metadataListingCache.Invalidate(ctx, path); err != nil {
				r.logger.Warn("ListByPath: failed to invalidate listing cache", logging.String("path", path), logging.Error(err))
			}
		}
	}

	// 3. Prepare conditional GET for children
	config := &drives.ItemItemsRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	if !opts.Force && parent != nil && parent.ETag != "" {
		config.Headers.Add("If-None-Match", parent.ETag)
	}

	path = normalizePath(path)

	// 4. Fetch children
	r.logger.Debug("ListByPath: requesting children from OneDrive", logging.String("path", path))
	items, err := r.childrenBuilder(r.client, driveID, path).Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		r.logger.Error("ListByPath: failed to fetch children", logging.String("path", path), logging.Error(err))
		return nil, err
	}

	// 5. 304 Not Modified → return cached listing
	if items == nil && !opts.NoCache {
		r.logger.Info("ListByPath: 304 Not Modified", logging.String("path", path))
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			children := make([]*file.Metadata, 0, len(listing.ChildIDs))
			allFound := true
			for _, id := range listing.ChildIDs {
				m, ok := r.metadataCache.Get(ctx, id)
				if !ok || m == nil {
					r.logger.Warn("ListByPath: inconsistent cache after 304, child metadata missing", logging.String("id", id))
					allFound = false
					break
				}
				children = append(children, m)
			}
			if allFound {
				return children, nil
			}
			r.logger.Info("ListByPath: invalidating inconsistent listing cache after 304", logging.String("path", path))
			if err := r.metadataListingCache.Invalidate(ctx, path); err != nil {
				r.logger.Warn("ListByPath: failed to invalidate listing cache after 304", logging.String("path", path), logging.Error(err))
			}
		}
	}

	if items == nil {
		r.logger.Warn("ListByPath: items is nil and no cache found", logging.String("path", path))
		return nil, nil
	}

	// 6. Fresh listing
	realItems := items.GetValue()
	r.logger.Info("ListByPath: fresh children received", logging.String("path", path), logging.Int("count", len(realItems)))
	metadatas := make([]*file.Metadata, len(realItems))
	listing := &Listing{
		ETag:     parent.ETag,
		ChildIDs: make([]string, len(realItems)),
	}

	for i, item := range realItems {
		if item == nil {
			r.logger.Warn("ListByPath: received nil item from OneDrive", logging.Int("index", i))
			continue
		}
		m := mapItemToMetadata(item)
		if m == nil {
			r.logger.Warn("ListByPath: mapping item to metadata returned nil", logging.Int("index", i))
			continue
		}
		metadatas[i] = m
		listing.ChildIDs[i] = m.ID

		if !opts.NoStore {
			r.logger.Debug("ListByPath: caching child metadata", logging.String("id", m.ID))
			if err := r.metadataCache.Put(ctx, m.ID, m); err != nil {
				r.logger.Warn("ListByPath: failed to cache child metadata", logging.String("id", m.ID), logging.Error(err))
			}
		}
	}

	if !opts.NoStore {
		r.logger.Debug("ListByPath: caching listing", logging.String("path", path))
		if err := r.metadataListingCache.Put(ctx, path, listing); err != nil {
			r.logger.Warn("ListByPath: failed to cache listing", logging.String("path", path), logging.Error(err))
		}
	}

	return metadatas, nil
}

func (s *MetadataRepository) driveItemBuilder(client abstractions.RequestAdapter, driveID, normalizedPath string) *drives.ItemRootRequestBuilder {
	urlTemplate := rootURITemplate2
	subs := make(stduritemplate.Substitutions)
	subs["baseurl"] = baseURL
	subs["drive_id"] = driveID

	if normalizedPath != "" {
		urlTemplate = rootRelativeURITemplate2
		subs["path"] = normalizedPath
	}

	uri, _ := stduritemplate.Expand(urlTemplate, subs)

	return drives.NewItemRootRequestBuilder(uri, client)
}

func (s *MetadataRepository) childrenBuilder(client abstractions.RequestAdapter, driveID, normalizedPath string) *drives.ItemItemsRequestBuilder {
	urlTemplate := rootChildrenURITemplate2
	subs := make(stduritemplate.Substitutions)
	subs["baseurl"] = baseURL
	subs["drive_id"] = driveID

	if normalizedPath != "" {
		urlTemplate = rootRelativeChildrenURITemplate2
		subs["path"] = normalizedPath
	}

	uri, _ := stduritemplate.Expand(urlTemplate, subs)

	return drives.NewItemItemsRequestBuilder(uri, client)
}
