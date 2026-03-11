package infra

import (
	"context"

	graphinfra "github.com/michaeldcanady/go-onedrive/internal/core/graph/infra"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type GraphDriveGateway struct {
	client abstractions.RequestAdapter
	log    domainlogger.Logger
}

func NewGraphDriveGateway(client abstractions.RequestAdapter, log domainlogger.Logger) *GraphDriveGateway {
	return &GraphDriveGateway{
		client: client,
		log:    log,
	}
}

func (g *GraphDriveGateway) ListDrives(ctx context.Context) ([]*domaindrive.Drive, error) {
	resp, err := users.NewItemDrivesRequestBuilderInternal(map[string]string{"user%2Did": "me-token-to-replace"}, g.client).Get(ctx, nil)
	if err != nil {
		return nil, graphinfra.MapGraphError(err, true)
	}

	var out []*domaindrive.Drive

	pageIterator, err := msgraphsdkcore.NewPageIterator[models.Driveable](resp, g.client, models.CreateDriveCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(d models.Driveable) bool {
		out = append(out, toDomainDrive(d))
		return true
	})

	if err != nil {
		return nil, graphinfra.MapGraphError(err, true)
	}

	return out, nil
}

func (g *GraphDriveGateway) GetPersonalDrive(ctx context.Context) (*domaindrive.Drive, error) {
	resp, err := users.NewItemDriveRequestBuilderInternal(map[string]string{"user%2Did": "me-token-to-replace"}, g.client).Get(ctx, nil)
	if err != nil {
		return nil, graphinfra.MapGraphError(err, true)
	}

	return toDomainDrive(resp), nil
}
