package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// GraphMetadataGateway implements GraphGateway for interacting with OneDrive metadata via Microsoft Graph.
type GraphMetadataGateway struct {
	client abstractions.RequestAdapter
	log    logger.Logger
}

// NewGraphMetadataGateway creates a new instance of GraphMetadataGateway.
func NewGraphMetadataGateway(client abstractions.RequestAdapter, l logger.Logger) *GraphMetadataGateway {
	return &GraphMetadataGateway{
		client: client,
		log:    l,
	}
}

func (g *GraphMetadataGateway) GetByPath(ctx context.Context, driveID, path string, etag string) (*file.Metadata, error) {
	config := &drives.ItemRootRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	if etag != "" {
		config.Headers.Add("If-None-Match", etag)
	}

	uri := expandPathTemplate(rootURITemplate2, rootRelativeURITemplate2, driveID, path)
	builder := drives.NewItemRootRequestBuilder(uri, g.client)

	item, err := builder.Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil // 304 Not Modified
	}

	return mapItemToMetadata(item), nil
}

func (g *GraphMetadataGateway) ListByPath(ctx context.Context, driveID, path string, parentEtag string) ([]*file.Metadata, error) {
	config := &drives.ItemItemsRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	if parentEtag != "" {
		config.Headers.Add("If-None-Match", parentEtag)
	}

	uri := expandPathTemplate(rootChildrenURITemplate2, rootRelativeChildrenURITemplate2, driveID, path)
	builder := drives.NewItemItemsRequestBuilder(uri, g.client)

	items, err := builder.Get(ctx, config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	if items == nil {
		return nil, nil // 304 Not Modified
	}

	realItems := items.GetValue()
	metadatas := make([]*file.Metadata, 0, len(realItems))
	for _, item := range realItems {
		if item == nil {
			continue
		}
		metadatas = append(metadatas, mapItemToMetadata(item))
	}

	return metadatas, nil
}

func (g *GraphMetadataGateway) CreateByPath(ctx context.Context, driveID, parentPath string, request file.MetadataCreateRequest) (*file.Metadata, error) {
	requestBody := models.NewDriveItem()
	name := request.Name
	requestBody.SetName(&name)

	switch request.Type {
	case file.ItemTypeFolder:
		folder := models.NewFolder()
		requestBody.SetFolder(folder)
	case file.ItemTypeFile:
		file := models.NewFile()
		requestBody.SetFile(file)
	}

	config := &drives.ItemItemsRequestBuilderPostRequestConfiguration{}

	uri := expandPathTemplate(rootChildrenURITemplate2, rootRelativeChildrenURITemplate2, driveID, parentPath)
	builder := drives.NewItemItemsRequestBuilder(uri, g.client)

	item, err := builder.Post(ctx, requestBody, config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	return mapItemToMetadata(item), nil
}

func (g *GraphMetadataGateway) UpdateByPath(ctx context.Context, driveID, path string, request file.MetadataUpdateRequest) (*file.Metadata, error) {
	requestBody := models.NewDriveItem()
	if request.Name != "" {
		requestBody.SetName(&request.Name)
	}

	if request.ParentPath != "" {
		parentRef := models.NewItemReference()
		p := normalizePath(request.ParentPath)
		if p == "" {
			p = "/"
		}
		parentRef.SetPath(&p)
		requestBody.SetParentReference(parentRef)
	}

	config := &drives.ItemItemsDriveItemItemRequestBuilderPatchRequestConfiguration{}

	uri := expandPathTemplate(rootURITemplate2, rootRelativeURITemplate2, driveID, path)
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(uri, g.client)

	item, err := builder.Patch(ctx, requestBody, config)
	if err := mapGraphError2(err); err != nil {
		return nil, err
	}

	return mapItemToMetadata(item), nil
}
