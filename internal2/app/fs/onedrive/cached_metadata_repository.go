package onedrive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

// CachedMetadataRepository implements file.MetadataRepository with caching and path-to-ID resolution.
type CachedMetadataRepository struct {
	gateway              infrafile.GraphGateway
	metadataCache        infrafile.MetadataCache
	metadataListingCache infrafile.ListingCache
	pathIDCache          infrafile.PathIDCache
	log                  logger.Logger
}

func NewCachedMetadataRepository(
	gateway infrafile.GraphGateway,
	metadataCache infrafile.MetadataCache,
	metadataListingCache infrafile.ListingCache,
	pathIDCache infrafile.PathIDCache,
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

func (r *CachedMetadataRepository) GetByPath(ctx context.Context, driveID string, path string, opts file.MetadataGetOptions) (*file.Metadata, error) {
	m, _, err := r.getByPath(ctx, driveID, path, opts)
	return m, err
}

func (r *CachedMetadataRepository) getByPath(ctx context.Context, driveID, path string, opts file.MetadataGetOptions) (*file.Metadata, bool, error) {
	log := r.log.WithContext(ctx).With(logger.String("path", path))

	// 1. Cache lookup
	var cached *file.Metadata
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

func (r *CachedMetadataRepository) ListByPath(ctx context.Context, driveID string, path string, opts file.MetadataListOptions) ([]*file.Metadata, error) {
	log := r.log.WithContext(ctx).With(logger.String("path", path))

	// 1. Get parent metadata
	parent, updated, err := r.getByPath(ctx, driveID, path, file.MetadataGetOptions{
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
			children := make([]*file.Metadata, 0, len(listing.ChildIDs))
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
			children := make([]*file.Metadata, 0, len(listing.ChildIDs))
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
		listing := &file.Listing{
			ETag:     parent.ETag,
			ChildIDs: childIDs,
		}
		if err := r.metadataListingCache.Put(ctx, path, listing); err != nil {
			log.Warn("failed to cache listing", logger.Error(err))
		}
	}

	return metadatas, nil
}

func (r *CachedMetadataRepository) CreateByPath(ctx context.Context, driveID, parentPath string, body file.MetadataCreateRequest, opts file.MetadataCreateOptions) (*file.Metadata, error) {
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

func (r *CachedMetadataRepository) UpdateByPath(ctx context.Context, driveID, path string, body file.MetadataUpdateRequest, opts file.MetadataUpdateOptions) (*file.Metadata, error) {
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
