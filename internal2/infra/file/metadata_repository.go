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

func (r *MetadataRepository) getByPath(ctx context.Context, path string) (*file.Metadata, bool, error) {
	config := &drives.ItemRootRequestBuilderGetRequestConfiguration{}

	metadata, ok := r.metadataCache.Get(ctx, path)
	if ok && metadata != nil {
		config.Headers.Add("If-None-Match", metadata.ETag)
	}

	item, err := r.driveItemBuilder(r.client, "drive-id", normalizePath(path)).Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		return nil, false, err
	}

	if item == nil {
		return metadata, false, nil
	}

	metadata = mapItemToMetadata(item)

	if err := r.metadataCache.Put(ctx, metadata); err != nil {
		return nil, false, err
	}
	return metadata, true, nil
}

func (r *MetadataRepository) GetByPath(ctx context.Context, path string) (*file.Metadata, error) {
	metadata, _, err := r.getByPath(ctx, path)

	return metadata, err
}

func (r *MetadataRepository) ListByPath(ctx context.Context, path string) ([]*file.Metadata, error) {
	// 1. Always fetch parent first (to get CTag)
	parent, updated, err := r.getByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	// 2. If parent was not updated, try listing cache
	if !updated {
		if listing, ok := r.metadataListingCache.Get(ctx, path); ok {
			children := make([]*file.Metadata, len(listing.ChildIDs))
			for i, id := range listing.ChildIDs {
				m, _ := r.metadataCache.Get(ctx, id)
				children[i] = m
			}
			return children, nil
		}
	}

	// 3. Prepare conditional GET for children using parent's CTag
	config := &drives.ItemItemsRequestBuilderGetRequestConfiguration{}
	if parent != nil && parent.CTag != "" {
		config.Headers.Add("If-None-Match", parent.CTag)
	}

	// 4. Fetch children
	items, err := r.childrenBuilder(r.client, "drive-id", normalizePath(path)).Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	// 5. 304 Not Modified → return cached listing
	if items == nil {
		listing, ok := r.metadataListingCache.Get(ctx, path)
		if ok {
			children := make([]*file.Metadata, len(listing.ChildIDs))
			for i, id := range listing.ChildIDs {
				children[i], _ = r.metadataCache.Get(ctx, id)
			}
			return children, nil
		}
	}

	// 6. Fresh listing → update caches
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

		_ = r.metadataCache.Put(ctx, m)
	}

	if err := r.metadataListingCache.Put(ctx, path, listing); err != nil {
		return nil, err
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
