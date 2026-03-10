package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type GraphDriveGateway struct {
	client abstractions.RequestAdapter
	log    logger.Logger
}

func NewGraphDriveGateway(client abstractions.RequestAdapter, log logger.Logger) *GraphDriveGateway {
	return &GraphDriveGateway{
		client: client,
		log:    log,
	}
}

func (g *GraphDriveGateway) ListDrives(ctx context.Context) ([]*drive.Drive, error) {
	resp, err := users.NewItemDrivesRequestBuilderInternal(map[string]string{"user%2Did": "me-token-to-replace"}, g.client).Get(ctx, nil)
	if err != nil {
		return nil, mapGraphError(err)
	}

	// TODO: use page iterator, in case someone has pages of drives
	out := make([]*drive.Drive, 0, len(resp.GetValue()))
	for _, d := range resp.GetValue() {
		out = append(out, toDomainDrive(d))
	}

	return out, nil
}

func (g *GraphDriveGateway) GetPersonalDrive(ctx context.Context) (*drive.Drive, error) {
	resp, err := users.NewItemDriveRequestBuilderInternal(map[string]string{"user%2Did": "me-token-to-replace"}, g.client).Get(ctx, nil)
	if err != nil {
		return nil, mapGraphError(err)
	}

	return toDomainDrive(resp), nil
}
