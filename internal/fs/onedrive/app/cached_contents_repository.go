package app

import (
	"bytes"
	"context"
	"io"

	logger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/fs/onedrive/infra"
)

type CachedFileContentsRepository struct {
	gateway       infra.GraphContentsGateway
	contentCache  infra.ContentsCache
	metadataCache infra.MetadataCache
	pathIDCache   infra.PathIDCache
	log           logger.Logger
}

func NewCachedFileContentsRepository(
	gateway infra.GraphContentsGateway,
	contentCache infra.ContentsCache,
	metadataCache infra.MetadataCache,
	pathIDCache infra.PathIDCache,
	l logger.Logger,
) *CachedFileContentsRepository {
	return &CachedFileContentsRepository{
		gateway:       gateway,
		contentCache:  contentCache,
		metadataCache: metadataCache,
		pathIDCache:   pathIDCache,
		log:           l,
	}
}

func (r *CachedFileContentsRepository) Download(ctx context.Context, driveID, path string, opts domain.DownloadOptions) (io.ReadCloser, error) {
	var cachedData []byte
	var etag string

	if !opts.NoCache {
		cacheKey := path
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			cacheKey = id
		}

		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			cachedData = entry.Data
			if entry.CTag != "" {
				etag = entry.CTag
			}
		}
	}

	fresh, ctag, err := r.gateway.Download(ctx, driveID, path, etag)
	if err != nil {
		return nil, err
	}

	if fresh == nil {
		// 304
		return io.NopCloser(bytes.NewReader(cachedData)), nil
	}

	if !opts.NoStore && ctag != "" {
		cacheKey := path
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			cacheKey = id
		}
		_ = r.contentCache.Put(ctx, cacheKey, &domain.Contents{
			CTag: ctag,
			Data: fresh,
		})
	}

	return io.NopCloser(bytes.NewReader(fresh)), nil
}

func (r *CachedFileContentsRepository) Upload(ctx context.Context, driveID, path string, body io.Reader, opts domain.UploadOptions) (*domain.Metadata, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	ifMatch := opts.IfMatch
	if ifMatch == "" && !opts.Overwrite {
		cacheKey := path
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			cacheKey = id
		}
		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			ifMatch = entry.CTag
		}
	}

	metadata, ctag, err := r.gateway.Upload(ctx, driveID, path, data, ifMatch)
	if err != nil {
		return nil, err
	}

	if !opts.NoStore {
		_ = r.contentCache.Put(ctx, metadata.ID, &domain.Contents{
			CTag: ctag,
			Data: data,
		})
		_ = r.metadataCache.Put(ctx, metadata.ID, metadata)
		_ = r.pathIDCache.Put(ctx, path, metadata.ID)
	}

	return metadata, nil
}
