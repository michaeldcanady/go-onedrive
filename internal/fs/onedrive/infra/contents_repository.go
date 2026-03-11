package infra

import (
	"context"

	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	nethttplibrary "github.com/microsoft/kiota-http-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
)

// GraphContentsGateway implements GraphContentsGateway for interacting with OneDrive content via Microsoft Graph.
type GraphFileContentsGateway struct {
	client abstractions.RequestAdapter
	log    domainlogger.Logger
}

// NewGraphFileContentsGateway initializes a new GraphFileContentsGateway.
func NewGraphFileContentsGateway(client abstractions.RequestAdapter, l domainlogger.Logger) *GraphFileContentsGateway {
	return &GraphFileContentsGateway{
		client: client,
		log:    l,
	}
}

func (g *GraphFileContentsGateway) Download(ctx context.Context, driveID, path string, etag string) ([]byte, string, error) {
	config := drives.ItemRootContentRequestBuilderGetRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
		Options: []abstractions.RequestOption{},
	}

	if etag != "" {
		config.Headers.Add("If-None-Match", etag)
	}

	headerOpt := nethttplibrary.NewHeadersInspectionOptions()
	headerOpt.InspectResponseHeaders = true
	config.Options = append(config.Options, headerOpt)

	uri := expandPathTemplate("", rootRelativeContentURITemplate2, driveID, path)
	builder := drives.NewItemRootContentRequestBuilder(uri, g.client)

	resp, err := builder.Get(ctx, &config)
	if err := mapGraphError2(err); err != nil {
		return nil, "", err
	}

	if resp == nil {
		return nil, "", nil // 304 Not Modified
	}

	// Extract CTag
	headers := headerOpt.GetResponseHeaders()
	ctagValues := headers.Get("CTag")
	if len(ctagValues) == 0 {
		ctagValues = headers.Get("ctag")
	}
	var ctag string
	if len(ctagValues) > 0 {
		ctag = ctagValues[0]
	}

	return resp, ctag, nil
}

func (g *GraphFileContentsGateway) Upload(ctx context.Context, driveID, path string, data []byte, ifMatch string) (*domainfs.Metadata, string, error) {
	config := &drives.ItemRootContentRequestBuilderPutRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}

	if ifMatch != "" {
		config.Headers.Add("If-Match", ifMatch)
	}

	uri := expandPathTemplate("", rootRelativeContentURITemplate2, driveID, path)
	builder := drives.NewItemRootContentRequestBuilder(uri, g.client)

	item, err := builder.Put(ctx, data, config)
	if err := mapGraphError2(err); err != nil {
		return nil, "", err
	}

	metadata := mapItemToMetadata(item)
	ctag := ""
	if item.GetCTag() != nil {
		ctag = *item.GetCTag()
	}

	return metadata, ctag, nil
}
