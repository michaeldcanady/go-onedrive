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

type ContentsRepository struct {
	client        abstractions.RequestAdapter
	contentCache  ContentsCache
	metadataCache MetadataCache
}

func NewContentsRepository(client abstractions.RequestAdapter, contentCache ContentsCache, metadataCache MetadataCache) *ContentsRepository {
	return &ContentsRepository{
		client:        client,
		contentCache:  contentCache,
		metadataCache: metadataCache,
	}
}

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

	if !opts.NoStore {
		headerOpt = nethttplibrary.NewHeadersInspectionOptions()
		headerOpt.InspectResponseHeaders = true
		config.Options = append(config.Options, headerOpt)
	}

	// Try cache
	if !opts.NoCache {
		if entry, ok := r.contentCache.Get(ctx, path); ok {
			cached = io.NopCloser(bytes.NewReader(entry.Data))
			if entry.CTag != "" {
				config.Headers.Add("If-None-Match", entry.CTag)
			}
		}
	}

	resp, err := r.relativePathContentsBuilder(r.client, driveID, normalizePath(path)).Get(ctx, &config)
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
) error {
	config := &drives.ItemRootContentRequestBuilderPutRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	if !opts.Force {
		if entry, ok := r.contentCache.Get(ctx, path); ok {
			if entry.CTag != "" && len(entry.Data) > 0 {
				config.Headers.Add("If-Match", entry.CTag)
			}
		}
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	// 3. Upload
	item, err := r.relativePathContentsBuilder(r.client, driveID, normalizePath(path)).Put(ctx, data, config)
	if err := mapGraphError2(err); err != nil {
		return err
	}

	if !opts.NoStore {
		// update contents cache
		if err := r.contentCache.Put(ctx, path, &file.Contents{
			CTag: *item.GetCTag(),
			Data: data,
		}); err != nil {
			return err
		}

		// update metadata cache
		if err := r.metadataCache.Put(ctx, mapItemToMetadata(item)); err != nil {
			return err
		}
	}

	return nil
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
