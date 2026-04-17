package graph

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	idshared "github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	msgraphsdkcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// GraphDriveGateway implements the drive.Gateway interface using Microsoft Graph.
type GraphDriveGateway struct {
	auth idshared.Authenticator
	log  logger.Logger
}

// NewGraphDriveGateway initializes a new instance of the GraphDriveGateway.
func NewGraphDriveGateway(auth idshared.Authenticator, log logger.Logger) *GraphDriveGateway {
	return &GraphDriveGateway{
		auth: auth,
		log:  log,
	}
}

func (g *GraphDriveGateway) ListDrives(ctx context.Context, identityID string) ([]drive.Drive, error) {
	cred, err := g.auth.GetCredential(ctx, identityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential for identity %s: %w", identityID, err)
	}

	p := microsoft.NewGraphProvider(cred.(azcore.TokenCredential), g.log)
	adapter, err := p.Adapter(ctx)
	if err != nil {
		return nil, err
	}

	// Using "me" as a shortcut for the authenticated user.
	resp, err := users.NewItemDrivesRequestBuilderInternal(map[string]string{"user%2Did": "me"}, adapter).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var out []drive.Drive
	pageIterator, err := msgraphsdkcore.NewPageIterator[models.Driveable](resp, adapter, models.CreateDriveCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, err
	}

	err = pageIterator.Iterate(ctx, func(d models.Driveable) bool {
		out = append(out, toDrive(d))
		return true
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

// GetPersonalDrive retrieves the user's default personal drive.
func (g *GraphDriveGateway) GetPersonalDrive(ctx context.Context, identityID string) (drive.Drive, error) {
	cred, err := g.auth.GetCredential(ctx, identityID)
	if err != nil {
		return drive.Drive{}, fmt.Errorf("failed to get credential for identity %s: %w", identityID, err)
	}

	p := microsoft.NewGraphProvider(cred.(azcore.TokenCredential), g.log)
	adapter, err := p.Adapter(ctx)
	if err != nil {
		return drive.Drive{}, err
	}

	resp, err := users.NewItemDriveRequestBuilderInternal(map[string]string{"user%2Did": "me"}, adapter).Get(ctx, nil)
	if err != nil {
		return drive.Drive{}, err
	}

	return toDrive(resp), nil
}

func toDrive(d models.Driveable) drive.Drive {
	if d == nil {
		return drive.Drive{}
	}

	id := ""
	if d.GetId() != nil {
		id = *d.GetId()
	}

	name := ""
	if d.GetName() != nil {
		name = *d.GetName()
	}

	driveType := drive.DriveTypeUnknown
	if d.GetDriveType() != nil {
		driveType = drive.NewDriveType(*d.GetDriveType())
	}

	owner := ""
	if d.GetOwner() != nil && d.GetOwner().GetUser() != nil && d.GetOwner().GetUser().GetDisplayName() != nil {
		owner = *d.GetOwner().GetUser().GetDisplayName()
	}

	readOnly := false
	// Logic for read-only could be more complex depending on permissions, but keeping it simple for now.

	return drive.Drive{
		ID:       id,
		Name:     name,
		Type:     driveType,
		Owner:    owner,
		ReadOnly: readOnly,
	}
}
