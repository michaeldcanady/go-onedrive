package app

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	logger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/fs/infra"
)

// CachedMetadataRepository implements domain.MetadataRepository with caching and path-to-ID resolution.
type CachedMetadataRepository struct {
	gateway              infra.GraphGateway
	metadataCache        infra.MetadataCache
	metadataListingCache infra.ListingCache
	pathIDCache          infra.PathIDCache
	log                  logger.Logger
}

func NewCachedMetadataRepository(
	gateway infra.GraphGateway,
	metadataCache infra.MetadataCache,
	metadataListingCache infra.ListingCache,
	pathIDCache infra.PathIDCache,
	l logger.Logger,
) *CachedMetadataRepository {
	return &CachedMetadataRepository{
		gateway:              gateway,
		metadataCache:        metadataCache,
		metadataListingCache: metadataListingCache,
		pathIDCache:          pathIDCache,
		log:                  l,
	}
}

func (r *CachedMetadataRepository) GetByPath(ctx context.Context, driveID string, path string, opts domain.MetadataGetOptions) (*domain.Metadata, error) {
	m, _, err := r.getByPath(ctx, driveID, path, opts)
	return m, err
}

func (r *CachedMetadataRepository) getByPath(ctx context.Context, driveID, path string, opts domain.MetadataGetOptions) (*domain.Metadata, bool, error) {
	log := r.log.WithContext(ctx).With(logger.String("path", path))

	// 1. Cache lookup
	var cached *domain.Metadata
	var etag string
	if !opts.NoCache {
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			if m, ok := r.metadataCache.Get(ctx, id); ok {
				cached = m
				if !opts.Force {
					etag = m.ETag
				}
			}
		}
	}

	// 2. Gateway call
	metadata, err := r.gateway.GetByPath(ctx, driveID, path, etag)
	if err != nil {
		return nil, false, err
	}

	// 3. 304 Not Modified
	if metadata == nil {
		return cached, false, nil
	}

	// 4. Update cache
	if !opts.NoStore {
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			log.Warn("failed to update metadata cache", logger.Error(err))
		}
		if err := r.pathIDCache.Put(ctx, path, metadata.ID); err != nil {
			log.Warn("failed to update path-to-ID cache", logger.Error(err))
		}
	}

	return metadata, true, nil
}

func (r *CachedMetadataRepository) ListByPath(ctx context.Context, driveID string, path string, opts domain.MetadataListOptions) ([]*domain.Metadata, error) {
	log := r.log.WithContext(ctx).With(logger.String("path", path))

	// 1. Get parent metadata
	parent, updated, err := r.getByPath(ctx, driveID, path, domain.MetadataGetOptions{
		NoCache: opts.NoCache,
		NoStore: opts.NoStore,
		Force:   opts.Force,
	})
	if err != nil {
		return nil, err
	}

	// 2. Cache lookup for children
	if !updated && !opts.NoCache {
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			children := make([]*domain.Metadata, 0, len(listing.ChildIDs))
			allFound := true
			for _, id := range listing.ChildIDs {
				if m, ok := r.metadataCache.Get(ctx, id); ok {
					children = append(children, m)
				} else {
					allFound = false
					break
				}
			}
			if allFound {
				return children, nil
			}
		}
	}

	// 3. Gateway call for children
	var parentEtag string
	if !opts.Force && parent != nil {
		parentEtag = parent.ETag
	}

	metadatas, err := r.gateway.ListByPath(ctx, driveID, path, parentEtag)
	if err != nil {
		return nil, err
	}

	// 4. 304 Not Modified
	if metadatas == nil {
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			children := make([]*domain.Metadata, 0, len(listing.ChildIDs))
			allFound := true
			for _, id := range listing.ChildIDs {
				if m, ok := r.metadataCache.Get(ctx, id); ok {
					children = append(children, m)
				} else {
					allFound = false
					break
				}
			}
			if allFound {
				return children, nil
			}
		}
		// Inconsistent cache or truly empty? Re-fetch without ETag might be safer but gateway handled it.
		return nil, nil
	}

	// 5. Update cache
	if !opts.NoStore {
		childIDs := make([]string, len(metadatas))
		for i, m := range metadatas {
			childIDs[i] = m.ID
			if err := r.metadataCache.Put(ctx, m.ID, m); err != nil {
				log.Warn("failed to cache child metadata", logger.Error(err))
			}
		}
		listing := &domain.Listing{
			ETag:     parent.ETag,
			ChildIDs: childIDs,
		}
		if err := r.metadataListingCache.Put(ctx, path, listing); err != nil {
			log.Warn("failed to cache listing", logger.Error(err))
		}
	}

	return metadatas, nil
}

func (r *CachedMetadataRepository) CreateByPath(ctx context.Context, driveID, parentPath string, body domain.MetadataCreateRequest, opts domain.MetadataCreateOptions) (*domain.Metadata, error) {
	metadata, err := r.gateway.CreateByPath(ctx, driveID, parentPath, body)
	if err != nil {
		return nil, err
	}

	if !opts.NoStore {
		_ = r.metadataCache.Put(ctx, metadata.ID, metadata)
		_ = r.pathIDCache.Put(ctx, parentPath, metadata.ID) // This might be wrong, should be full path
	}

	return metadata, nil
}

func (r *CachedMetadataRepository) UpdateByPath(ctx context.Context, driveID, path string, body domain.MetadataUpdateRequest, opts domain.MetadataUpdateOptions) (*domain.Metadata, error) {
	metadata, err := r.gateway.UpdateByPath(ctx, driveID, path, body)
	if err != nil {
		return nil, err
	}

	if !opts.NoStore {
		_ = r.metadataCache.Put(ctx, metadata.ID, metadata)
		_ = r.pathIDCache.Put(ctx, metadata.Path, metadata.ID)
	}

	return metadata, nil
}

func (p *CachedMetadataRepository) buildLogger(ctx context.Context) logger.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return p.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)
}

func (r *CachedMetadataRepository) DeleteByPath(ctx context.Context, driveID, path string, opts domain.MetadataDeleteOptions) error {
	log := r.buildLogger(ctx).With(logger.String("path", path))

	// TODO: add support for permanent deletion: https://learn.microsoft.com/en-us/graph/api/driveitem-permanentdelete?view=graph-rest-1.0
	if err := r.gateway.DeleteByPath(ctx, driveID, path); err != nil {
		return err
	}

	if id, ok := r.pathIDCache.Get(ctx, path); ok {
		if err := r.metadataCache.Invalidate(ctx, id); err != nil {
			log.Warn("failed to invalidate cache", logger.Error(err))
		}
	}

	return nil
}
