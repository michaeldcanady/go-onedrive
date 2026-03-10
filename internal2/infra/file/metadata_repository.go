package file

import (
	"context"
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// MetadataRepository provides methods for fetching and listing drive item
// metadata from OneDrive. It supports multi-level caching (individual items
// and folder listings) and handles path normalization.
type MetadataRepository struct {
	client               abstractions.RequestAdapter
	metadataCache        MetadataCache
	metadataListingCache ListingCache
	pathIDCache          PathIDCache
	log                  logger.Logger
}

// NewMetadataRepository initializes a new MetadataRepository with the provided
// request adapter and cache implementations.
func NewMetadataRepository(client abstractions.RequestAdapter, metadataCache MetadataCache, metadataListingCache ListingCache, pathIDCache PathIDCache, l logger.Logger) *MetadataRepository {
	return &MetadataRepository{
		client:               client,
		metadataCache:        metadataCache,
		metadataListingCache: metadataListingCache,
		pathIDCache:          pathIDCache,
		log:                  l,
	}
}

// Listing represents a cached collection of drive item IDs at a specific path.
type Listing struct {
	// ETag is the entity tag of the parent folder when the listing was fetched.
	ETag string
	// ChildIDs is a slice of IDs for the items contained within the folder.
	ChildIDs []string
}

func (r *MetadataRepository) getByPath(ctx context.Context, driveID, path string, opts file.MetadataGetOptions) (*file.Metadata, bool, error) {
	log := r.log.WithContext(ctx)

	config := &drives.ItemRootRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	path = normalizePath(path)
	log.Debug("getByPath: starting retrieval", logger.String("path", path), logger.Bool("noCache", opts.NoCache))

	// 1. Cache lookup unless disabled
	var cached *file.Metadata
	if !opts.NoCache {
		// Try path-to-ID cache first if we have it
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			log.Debug("getByPath: path-to-ID hit", logger.String("path", path), logger.String("id", id))
			if m, ok := r.metadataCache.Get(ctx, id); ok && m != nil {
				log.Debug("getByPath: metadata cache hit (by id)", logger.String("id", id))
				cached = m
				if !opts.Force {
					config.Headers.Add("If-None-Match", m.ETag)
				}
			} else {
				log.Debug("getByPath: metadata cache miss (by id)", logger.String("id", id))
			}
		} else {
			log.Debug("getByPath: path-to-ID miss", logger.String("path", path))
			// Fallback to path lookup in metadata cache
			if m, ok := r.metadataCache.Get(ctx, path); ok && m != nil {
				log.Debug("getByPath: metadata cache hit (by path)", logger.String("path", path))
				cached = m
				if !opts.Force {
					config.Headers.Add("If-None-Match", m.ETag)
				}
			}
		}
	}

	// 2. Fetch from OneDrive
	log.Debug("getByPath: requesting from OneDrive", logger.String("path", path))
	uri := expandPathTemplate(rootURITemplate2, rootRelativeURITemplate2, driveID, path)
	builder := drives.NewItemRootRequestBuilder(uri, r.client)

	item, err := builder.Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		log.Error("getByPath: request failed", logger.String("path", path), logger.Error(err))
		return nil, false, err
	}

	// 3. 304 Not Modified → return cached
	if item == nil {
		log.Info("getByPath: 304 Not Modified", logger.String("path", path))
		return cached, false, nil
	}

	// 4. Fresh metadata
	metadata := mapItemToMetadata(item)
	log.Info("getByPath: fresh metadata received", logger.String("path", path), logger.String("id", metadata.ID))

	if !opts.NoStore {
		log.Debug("getByPath: updating cache", logger.String("id", metadata.ID))
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			log.Warn("getByPath: failed to update metadata cache", logger.String("id", metadata.ID), logger.Error(err))
		}
		if err := r.pathIDCache.Put(ctx, path, metadata.ID); err != nil {
			log.Warn("getByPath: failed to update path-to-ID cache", logger.String("path", path), logger.Error(err))
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
	log := r.log.WithContext(ctx)

	log.Debug("ListByPath: starting retrieval", logger.String("path", path))
	parent, updated, err := r.getByPath(ctx, driveID, path, file.MetadataGetOptions{
		NoCache: opts.NoCache,
		NoStore: opts.NoStore,
		Force:   opts.Force,
	})
	if err != nil {
		log.Error("ListByPath: failed to get parent metadata", logger.String("path", path), logger.Error(err))
		return nil, err
	}

	if parent == nil {
		log.Warn("ListByPath: parent metadata is nil", logger.String("path", path))
		return nil, ErrNotFound
	}

	// 2. If parent not updated and listing cache allowed → return cached listing
	if !updated && !opts.NoCache {
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			log.Debug("ListByPath: listing cache hit", logger.String("path", path))
			children := make([]*file.Metadata, 0, len(listing.ChildIDs))
			allFound := true
			for _, id := range listing.ChildIDs {
				m, ok := r.metadataCache.Get(ctx, id)
				if !ok || m == nil {
					log.Warn("ListByPath: inconsistent cache, child metadata missing", logger.String("id", id))
					allFound = false
					break
				}
				children = append(children, m)
			}
			if allFound {
				log.Info("ListByPath: returning cached listing", logger.String("path", path), logger.Int("count", len(children)))
				return children, nil
			}
			// Inconsistent cache: listing exists but some items are missing.
			// Fall through to fetch from Graph.
			log.Info("ListByPath: invalidating inconsistent listing cache", logger.String("path", path))
			if err := r.metadataListingCache.Invalidate(ctx, path); err != nil {
				log.Warn("ListByPath: failed to invalidate listing cache", logger.String("path", path), logger.Error(err))
			}
		}
	}

	// 3. Prepare conditional GET for children
	config := &drives.ItemItemsRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	if !opts.Force && parent.ETag != "" {
		config.Headers.Add("If-None-Match", parent.ETag)
	}

	path = normalizePath(path)

	// 4. Fetch children
	log.Debug("ListByPath: requesting children from OneDrive", logger.String("path", path))
	uri := expandPathTemplate(rootChildrenURITemplate2, rootRelativeChildrenURITemplate2, driveID, path)
	builder := drives.NewItemItemsRequestBuilder(uri, r.client)

	items, err := builder.Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		log.Error("ListByPath: failed to fetch children", logger.String("path", path), logger.Error(err))
		return nil, err
	}

	// 5. 304 Not Modified → return cached listing
	if items == nil && !opts.NoCache {
		log.Info("ListByPath: 304 Not Modified", logger.String("path", path))
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			children := make([]*file.Metadata, 0, len(listing.ChildIDs))
			allFound := true
			for _, id := range listing.ChildIDs {
				m, ok := r.metadataCache.Get(ctx, id)
				if !ok || m == nil {
					log.Warn("ListByPath: inconsistent cache after 304, child metadata missing", logger.String("id", id))
					allFound = false
					break
				}
				children = append(children, m)
			}
			if allFound {
				return children, nil
			}
			log.Info("ListByPath: invalidating inconsistent listing cache after 304", logger.String("path", path))
			if err := r.metadataListingCache.Invalidate(ctx, path); err != nil {
				log.Warn("ListByPath: failed to invalidate listing cache after 304", logger.String("path", path), logger.Error(err))
			}
		}
	}

	if items == nil {
		log.Warn("ListByPath: items is nil and no cache found", logger.String("path", path))
		return nil, nil
	}

	// 6. Fresh listing
	realItems := items.GetValue()
	log.Info("ListByPath: fresh children received", logger.String("path", path), logger.Int("count", len(realItems)))
	metadatas := make([]*file.Metadata, len(realItems))
	listing := &Listing{
		ETag:     parent.ETag,
		ChildIDs: make([]string, len(realItems)),
	}

	for i, item := range realItems {
		if item == nil {
			log.Warn("ListByPath: received nil item from OneDrive", logger.Int("index", i))
			continue
		}
		m := mapItemToMetadata(item)
		if m == nil {
			log.Warn("ListByPath: mapping item to metadata returned nil", logger.Int("index", i))
			continue
		}
		metadatas[i] = m
		listing.ChildIDs[i] = m.ID

		if !opts.NoStore {
			log.Debug("ListByPath: caching child metadata", logger.String("id", m.ID))
			if err := r.metadataCache.Put(ctx, m.ID, m); err != nil {
				log.Warn("ListByPath: failed to cache child metadata", logger.String("id", m.ID), logger.Error(err))
			}
		}
	}

	if !opts.NoStore {
		log.Debug("ListByPath: caching listing", logger.String("path", path))
		if err := r.metadataListingCache.Put(ctx, path, listing); err != nil {
			log.Warn("ListByPath: failed to cache listing", logger.String("path", path), logger.Error(err))
		}
	}

	return metadatas, nil
}

func (r *MetadataRepository) CreateByPath(ctx context.Context, driveID, parentPath string, body file.MetadataCreateRequest, opts file.MetadataCreateOptions) (*file.Metadata, error) {
	log := r.log.WithContext(ctx)

	requestBody := models.NewDriveItem()
	name := body.Name
	requestBody.SetName(&name)

	switch body.Type {
	case file.ItemTypeFolder:
		folder := models.NewFolder()
		requestBody.SetFolder(folder)
	case file.ItemTypeFile:
		file := models.NewFile()
		requestBody.SetFile(file)
	default:
		log.Warn("unsupported file type", logger.String("file_type", body.Type.String()))
		return nil, errors.New("unsupported file type")
	}

	config := &drives.ItemItemsRequestBuilderPostRequestConfiguration{}

	uri := expandPathTemplate(rootChildrenURITemplate2, rootRelativeChildrenURITemplate2, driveID, parentPath)
	builder := drives.NewItemItemsRequestBuilder(uri, r.client)

	item, err := builder.Post(ctx, requestBody, config)
	if err := mapGraphError2(err); err != nil {
		log.Error("CreateByPath: request failed", logger.String("path", parentPath), logger.Error(err))
		return nil, err
	}

	metadata := mapItemToMetadata(item)
	log.Info("CreateByPath: fresh metadata received", logger.String("path", parentPath), logger.String("id", metadata.ID))

	if !opts.NoStore {
		log.Debug("CreateByPath: updating cache", logger.String("id", metadata.ID))
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			log.Warn("CreateByPath: failed to update metadata cache", logger.String("id", metadata.ID), logger.Error(err))
		}
		if err := r.pathIDCache.Put(ctx, parentPath, metadata.ID); err != nil {
			log.Warn("CreateByPath: failed to update path-to-ID cache", logger.String("path", parentPath), logger.Error(err))
		}
	}

	return metadata, nil
}

func (r *MetadataRepository) UpdateByPath(ctx context.Context, driveID, path string, body file.MetadataUpdateRequest, opts file.MetadataUpdateOptions) (*file.Metadata, error) {
	log := r.log.WithContext(ctx)

	requestBody := models.NewDriveItem()
	if body.Name != "" {
		requestBody.SetName(&body.Name)
	}

	if body.ParentPath != "" {
		parentRef := models.NewItemReference()
		p := normalizePath(body.ParentPath)
		if p == "" {
			p = "/"
		}
		parentRef.SetPath(&p)
		requestBody.SetParentReference(parentRef)
	}

	config := &drives.ItemItemsDriveItemItemRequestBuilderPatchRequestConfiguration{}

	uri := expandPathTemplate(rootURITemplate2, rootRelativeURITemplate2, driveID, path)
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(uri, r.client)

	item, err := builder.Patch(ctx, requestBody, config)
	if err := mapGraphError2(err); err != nil {
		log.Error("UpdateByPath: request failed", logger.String("path", path), logger.Error(err))
		return nil, err
	}

	metadata := mapItemToMetadata(item)
	log.Info("UpdateByPath: fresh metadata received", logger.String("path", path), logger.String("id", metadata.ID))

	if !opts.NoStore {
		log.Debug("UpdateByPath: updating cache", logger.String("id", metadata.ID))
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			log.Warn("UpdateByPath: failed to update metadata cache", logger.String("id", metadata.ID), logger.Error(err))
		}
		if err := r.pathIDCache.Put(ctx, metadata.FullPath, metadata.ID); err != nil {
			log.Warn("UpdateByPath: failed to update path-to-ID cache", logger.String("path", metadata.FullPath), logger.Error(err))
		}
	}

	return metadata, nil
}
