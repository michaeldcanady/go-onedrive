package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

type MetadataRepository struct {
	client               abstractions.RequestAdapter
	metadataCache        MetadataCache
	metadataListingCache ListingCache
}

func (r *MetadataRepository) getByPath(
	ctx context.Context,
	driveID, path string,
	opts file.MetadataGetOptions,
) (*file.Metadata, bool, error) {

	config := &drives.ItemRootRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	// 1. Cache lookup unless disabled
	var cached *file.Metadata
	if !opts.NoCache {
		if m, ok := r.metadataCache.Get(ctx, path); ok && m != nil {
			cached = m
			if !opts.Force {
				config.Headers.Add("If-None-Match", m.ETag)
			}
		}
	}

	// 2. Fetch from OneDrive
	item, err := r.driveItemBuilder(r.client, driveID, normalizePath(path)).Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		return nil, false, err
	}

	// 3. 304 Not Modified → return cached
	if item == nil {
		return cached, false, nil
	}

	// 4. Fresh metadata
	metadata := mapItemToMetadata(item)

	if !opts.NoStore {
		_ = r.metadataCache.Put(ctx, metadata)
	}

	return metadata, true, nil
}

func (r *MetadataRepository) GetByPath(
	ctx context.Context,
	driveID, path string,
	opts file.MetadataGetOptions,
) (*file.Metadata, error) {

	metadata, _, err := r.getByPath(ctx, driveID, path, opts)
	return metadata, err
}

func (r *MetadataRepository) ListByPath(
	ctx context.Context,
	driveID, path string,
	opts file.MetadataListOptions,
) ([]*file.Metadata, error) {

	// 1. Fetch parent metadata (CTag)
	parent, updated, err := r.getByPath(ctx, driveID, path, file.MetadataGetOptions{
		NoCache: opts.NoCache,
		NoStore: opts.NoStore,
		Force:   opts.Force,
	})
	if err != nil {
		return nil, err
	}

	// 2. If parent not updated and listing cache allowed → return cached listing
	if !updated && !opts.NoCache {
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			children := make([]*file.Metadata, len(listing.ChildIDs))
			for i, id := range listing.ChildIDs {
				m, _ := r.metadataCache.Get(ctx, id)
				children[i] = m
			}
			return children, nil
		}
	}

	// 3. Prepare conditional GET for children
	config := &drives.ItemItemsRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	if !opts.Force && parent != nil && parent.CTag != "" {
		config.Headers.Add("If-None-Match", parent.CTag)
	}

	// 4. Fetch children
	items, err := r.childrenBuilder(r.client, driveID, normalizePath(path)).Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	// 5. 304 Not Modified → return cached listing
	if items == nil && !opts.NoCache {
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			children := make([]*file.Metadata, len(listing.ChildIDs))
			for i, id := range listing.ChildIDs {
				children[i], _ = r.metadataCache.Get(ctx, id)
			}
			return children, nil
		}
	}

	// 6. Fresh listing
	realItems := items.GetValue()
	metadatas := make([]*file.Metadata, len(realItems))
	listing := &Listing{
		CTag:     parent.CTag,
		ChildIDs: make([]string, len(realItems)),
	}

	for i, item := range realItems {
		m := mapItemToMetadata(item)
		metadatas[i] = m
		listing.ChildIDs[i] = m.ID

		if !opts.NoStore {
			_ = r.metadataCache.Put(ctx, m)
		}
	}

	if !opts.NoStore {
		_ = r.metadataListingCache.Put(ctx, path, listing)
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
