package file

import (
	"bytes"
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	nethttplibrary "github.com/microsoft/kiota-http-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

// ContentsRepository provides methods for downloading and uploading file
// content to OneDrive. It integrates with caching mechanisms for improved
// performance and handles Graph API interactions.
type ContentsRepository struct {
	client        abstractions.RequestAdapter
	contentCache  ContentsCache
	metadataCache MetadataCache
	pathIDCache   PathIDCache
}

// NewContentsRepository initializes a new ContentsRepository with the provided
// request adapter and cache implementations.
func NewContentsRepository(client abstractions.RequestAdapter, contentCache ContentsCache, metadataCache MetadataCache, pathIDCache PathIDCache) *ContentsRepository {
	return &ContentsRepository{
		client:        client,
		contentCache:  contentCache,
		metadataCache: metadataCache,
		pathIDCache:   pathIDCache,
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
			cacheKey = id
		}

		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			cached = io.NopCloser(bytes.NewReader(entry.Data))
			if entry.CTag != "" {
				config.Headers.Add("If-None-Match", entry.CTag)
			}
		}
	}

	resp, err := r.relativePathContentsBuilder(r.client, driveID, path).Get(ctx, &config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	// 304 Not Modified
	if resp == nil {
		return cached, nil
	}

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
			// update contents cache
			if err := r.contentCache.Put(ctx, path, &file.Contents{
				CTag: ctag,
				Data: resp,
			}); err != nil {
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

	if !opts.Force {
		cacheKey := path
		if id, ok := r.pathIDCache.Get(ctx, path); ok {
			cacheKey = id
		}

		if entry, ok := r.contentCache.Get(ctx, cacheKey); ok {
			if entry.CTag != "" && len(entry.Data) > 0 {
				config.Headers.Add("If-Match", entry.CTag)
			}
		}
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	// 3. Upload
	item, err := r.relativePathContentsBuilder(r.client, driveID, path).Put(ctx, data, config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	metadata := mapItemToMetadata(item)

	if !opts.NoStore {
		// update contents cache
		if err := r.contentCache.Put(ctx, metadata.ID, &file.Contents{
			CTag: *item.GetCTag(),
			Data: data,
		}); err != nil {
			return nil, err
		}

		// update metadata cache
		if err := r.metadataCache.Put(ctx, metadata.ID, metadata); err != nil {
			return nil, err
		}

		// update path-to-id mapping
		_ = r.pathIDCache.Put(ctx, path, metadata.ID)
	}

	return metadata, nil
}

func (s *ContentsRepository) relativePathContentsBuilder(client abstractions.RequestAdapter, driveID, normalizedPath string) *drives.ItemRootContentRequestBuilder {
	urlTemplate := rootRelativeContentURITemplate2
	subs := make(stduritemplate.Substitutions)
	subs["baseurl"] = baseURL
	subs["drive_id"] = driveID
	subs["path"] = normalizedPath

	uri, _ := stduritemplate.Expand(urlTemplate, subs)

	return drives.NewItemRootContentRequestBuilder(uri, client)
}
