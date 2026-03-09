package file

import (
	"bytes"
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
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
	logger        logging.Logger
}

// NewContentsRepository initializes a new ContentsRepository with the provided
// request adapter and cache implementations.
func NewContentsRepository(client abstractions.RequestAdapter, contentCache ContentsCache, metadataCache MetadataCache, pathIDCache PathIDCache, logger logging.Logger) *ContentsRepository {
	return &ContentsRepository{
		client:        client,
		contentCache:  contentCache,
		metadataCache: metadataCache,
		pathIDCache:   pathIDCache,
		logger:        logger,
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
	r.logger.Debug("Download: starting retrieval", logging.String("path", path), logging.Bool("noCache", opts.NoCache))

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
			r.logger.Debug("Download: path-to-ID hit", logging.String("path", path), logging.String("id", id))
			cacheKey = id
		}

		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			r.logger.Debug("Download: contents cache hit", logging.String("key", cacheKey))
			cached = io.NopCloser(bytes.NewReader(entry.Data))
			if entry.CTag != "" {
				config.Headers.Add("If-None-Match", entry.CTag)
			}
		} else {
			r.logger.Debug("Download: contents cache miss", logging.String("key", cacheKey))
		}
	}

	r.logger.Debug("Download: requesting from OneDrive", logging.String("path", path))
	uri := expandPathTemplate("", rootRelativeContentURITemplate2, driveID, path)
	builder := drives.NewItemRootContentRequestBuilder(uri, r.client)

	resp, err := builder.Get(ctx, &config)
	if err := mapGraphError2(err); err != nil {
		r.logger.Error("Download: request failed", logging.String("path", path), logging.Error(err))
		return nil, err
	}

	// 304 Not Modified
	if resp == nil {
		r.logger.Info("Download: 304 Not Modified", logging.String("path", path))
		return cached, nil
	}

	r.logger.Info("Download: received fresh content", logging.String("path", path), logging.Int("size", len(resp)))

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

			r.logger.Debug("Download: updating contents cache", logging.String("key", cacheKey))
			// update contents cache
			if err := r.contentCache.Put(ctx, cacheKey, &file.Contents{
				CTag: ctag,
				Data: resp,
			}); err != nil {
				r.logger.Warn("Download: failed to update cache", logging.Error(err))
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
	r.logger.Info("Upload: starting upload", logging.String("path", path))

	if opts.IfMatch != "" {
		r.logger.Debug("Upload: adding If-Match header from options", logging.String("etag", opts.IfMatch))
		config.Headers.Add("If-Match", opts.IfMatch)
	} else if !opts.Force {
		cacheKey := path
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			cacheKey = id
		}

		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			if entry.CTag != "" && len(entry.Data) > 0 {
				r.logger.Debug("Upload: adding If-Match header from cache", logging.String("ctag", entry.CTag))
				config.Headers.Add("If-Match", entry.CTag)
			}
		}
	}

	data, err := io.ReadAll(body)
	if err != nil {
		r.logger.Error("Upload: failed to read upload body", logging.Error(err))
		return nil, err
	}

	// 3. Upload
	r.logger.Debug("Upload: sending Put request to OneDrive", logging.String("path", path))
	uri := expandPathTemplate("", rootRelativeContentURITemplate2, driveID, path)
	builder := drives.NewItemRootContentRequestBuilder(uri, r.client)

	item, err := builder.Put(ctx, data, config)
	if err := mapGraphError2(err); err != nil {
		r.logger.Error("Upload: request failed", logging.String("path", path), logging.Error(err))
		return nil, err
	}

	metadata := mapItemToMetadata(item)
	r.logger.Info("Upload: upload successful", logging.String("path", path), logging.String("id", metadata.ID))

	if !opts.NoStore {
		r.logger.Debug("Upload: updating caches", logging.String("id", metadata.ID))
		// update contents cache
		if err := r.contentCache.Put(ctx, metadata.ID, &file.Contents{
			CTag: *item.GetCTag(),
			Data: data,
		}); err != nil {
			r.logger.Warn("Upload: failed to update contents cache", logging.Error(err))
			return nil, err
		}

		// update metadata cache
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			r.logger.Warn("Upload: failed to update metadata cache", logging.Error(err))
			return nil, err
		}

		// update path-to-id mapping
		if err := r.pathIDCache.Put(ctx, path, metadata.ID); err != nil {
			r.logger.Warn("Upload: failed to update path-to-ID cache", logging.String("path", path), logging.Error(err))
		}
	}

	return metadata, nil
}
