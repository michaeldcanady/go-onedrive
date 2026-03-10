package file

import (
	"bytes"
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	nethttplibrary "github.com/microsoft/kiota-http-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
)

// ContentsRepository provides methods for downloading and uploading file
// content to OneDrive. It integrates with caching mechanisms for improved
// performance and handles Graph API interactions.
type ContentsRepository struct {
	client        abstractions.RequestAdapter
	contentCache  ContentsCache
	metadataCache MetadataCache
	pathIDCache   PathIDCache
	log           logger.Logger
}

// NewContentsRepository initializes a new ContentsRepository with the provided
// request adapter and cache implementations.
func NewContentsRepository(client abstractions.RequestAdapter, contentCache ContentsCache, metadataCache MetadataCache, pathIDCache PathIDCache, l logger.Logger) *ContentsRepository {
	return &ContentsRepository{
		client:        client,
		contentCache:  contentCache,
		metadataCache: metadataCache,
		pathIDCache:   pathIDCache,
		log:           l,
	}
}

// Download retrieves file content for a given path in a drive.
// It supports conditional GET requests using cached ETags (CTags) unless
// caching is explicitly disabled via DownloadOptions.
func (r *ContentsRepository) Download(
	ctx context.Context,
	driveID,
	path string,
	opts file.DownloadOptions,
) (io.ReadCloser, error) {

	var (
		cached io.ReadCloser
		config = drives.ItemRootContentRequestBuilderGetRequestConfiguration{
			Headers: abstractions.NewRequestHeaders(),
			Options: []abstractions.RequestOption{},
		}
		headerOpt *nethttplibrary.HeadersInspectionOptions
	)

	path = normalizePath(path)
	r.log.Debug("Download: starting retrieval", logger.String("path", path), logger.Bool("noCache", opts.NoCache))

	if !opts.NoStore {
		headerOpt = nethttplibrary.NewHeadersInspectionOptions()
		headerOpt.InspectResponseHeaders = true
		config.Options = append(config.Options, headerOpt)
	}

	// Try cache
	if !opts.NoCache {
		// Use ID if we have it
		cacheKey := path
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			r.log.Debug("Download: path-to-ID hit", logger.String("path", path), logger.String("id", id))
			cacheKey = id
		}

		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			r.log.Debug("Download: contents cache hit", logger.String("key", cacheKey))
			cached = io.NopCloser(bytes.NewReader(entry.Data))
			if entry.CTag != "" {
				config.Headers.Add("If-None-Match", entry.CTag)
			}
		} else {
			r.log.Debug("Download: contents cache miss", logger.String("key", cacheKey))
		}
	}

	r.log.Debug("Download: requesting from OneDrive", logger.String("path", path))
	uri := expandPathTemplate("", rootRelativeContentURITemplate2, driveID, path)
	builder := drives.NewItemRootContentRequestBuilder(uri, r.client)

	resp, err := builder.Get(ctx, &config)
	if err := mapGraphError2(err); err != nil {
		r.log.Error("Download: request failed", logger.String("path", path), logger.Error(err))
		return nil, err
	}

	// 304 Not Modified
	if resp == nil {
		r.log.Info("Download: 304 Not Modified", logger.String("path", path))
		return cached, nil
	}

	r.log.Info("Download: received fresh content", logger.String("path", path), logger.Int("size", len(resp)))

	// Cache new content
	if !opts.NoStore && headerOpt != nil {
		// Extract ETag
		headers := headerOpt.GetResponseHeaders()
		ctagValues := headers.Get("CTag")
		if len(ctagValues) == 0 {
			ctagValues = headers.Get("ctag")
		}
		var ctag string
		if len(ctagValues) > 0 {
			ctag = ctagValues[0]
		}

		if len(ctag) > 0 {
			// Prefer ID for cache key if we have it
			cacheKey := path
			if id, ok := r.pathIDCache.Get(ctx, path); ok {
				cacheKey = id
			}

			r.log.Debug("Download: updating contents cache", logger.String("key", cacheKey))
			// update contents cache
			if err := r.contentCache.Put(ctx, cacheKey, &file.Contents{
				CTag: ctag,
				Data: resp,
			}); err != nil {
				r.log.Warn("Download: failed to update cache", logger.Error(err))
				return nil, err
			}
		}
	}

	return io.NopCloser(bytes.NewReader(resp)), nil
}

func (r *ContentsRepository) Upload(
	ctx context.Context,
	driveID,
	path string,
	body io.Reader,
	opts file.UploadOptions,
) (*file.Metadata, error) {
	config := &drives.ItemRootContentRequestBuilderPutRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	path = normalizePath(path)
	r.log.Info("Upload: starting upload", logger.String("path", path))

	if opts.IfMatch != "" {
		r.log.Debug("Upload: adding If-Match header from options", logger.String("etag", opts.IfMatch))
		config.Headers.Add("If-Match", opts.IfMatch)
	} else if !opts.Force {
		cacheKey := path
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			cacheKey = id
		}

		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			if entry.CTag != "" && len(entry.Data) > 0 {
				r.log.Debug("Upload: adding If-Match header from cache", logger.String("ctag", entry.CTag))
				config.Headers.Add("If-Match", entry.CTag)
			}
		}
	}

	data, err := io.ReadAll(body)
	if err != nil {
		r.log.Error("Upload: failed to read upload body", logger.Error(err))
		return nil, err
	}

	// 3. Upload
	r.log.Debug("Upload: sending Put request to OneDrive", logger.String("path", path))
	uri := expandPathTemplate("", rootRelativeContentURITemplate2, driveID, path)
	builder := drives.NewItemRootContentRequestBuilder(uri, r.client)

	item, err := builder.Put(ctx, data, config)
	if err := mapGraphError2(err); err != nil {
		r.log.Error("Upload: request failed", logger.String("path", path), logger.Error(err))
		return nil, err
	}

	metadata := mapItemToMetadata(item)
	r.log.Info("Upload: upload successful", logger.String("path", path), logger.String("id", metadata.ID))

	if !opts.NoStore {
		r.log.Debug("Upload: updating caches", logger.String("id", metadata.ID))
		// update contents cache
		if err := r.contentCache.Put(ctx, metadata.ID, &file.Contents{
			CTag: *item.GetCTag(),
			Data: data,
		}); err != nil {
			r.log.Warn("Upload: failed to update contents cache", logger.Error(err))
			return nil, err
		}

		// update metadata cache
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			r.log.Warn("Upload: failed to update metadata cache", logger.Error(err))
			return nil, err
		}

		// update path-to-id mapping
		if err := r.pathIDCache.Put(ctx, path, metadata.ID); err != nil {
			r.log.Warn("Upload: failed to update path-to-ID cache", logger.String("path", path), logger.Error(err))
		}
	}

	return metadata, nil
}
