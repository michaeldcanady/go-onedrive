package graph

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	proto "github.com/michaeldcanady/go-onedrive/internal/identity/proto"
	msgraphsdkcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// GraphDriveGateway implements the drive.Gateway interface using Microsoft Graph.
type GraphDriveGateway struct {
	auth identity.Service // Changed from identity.Authorizer to identity.Service
	log  logger.Logger
}

// NewGraphDriveGateway initializes a new instance of the GraphDriveGateway.
func NewGraphDriveGateway(auth identity.Service, log logger.Logger) *GraphDriveGateway {
	return &GraphDriveGateway{
		auth: auth,
		log:  log,
	}
}

func (g *GraphDriveGateway) ListDrives(ctx context.Context, identityID string) ([]drive.Drive, error) {
	req := &proto.GetTokenRequest{
		IdentityId: identityID,
		Scopes:     []string{"Files.Read"},
	}
	resp, err := g.auth.Token(ctx, "microsoft", req) // Pass provider
	if err != nil {
		return nil, fmt.Errorf("failed to get token for identity %s: %w", identityID, err)
	}


	accessToken := identity.FromProtoAccessToken(resp.GetToken(), identityID)
	cred := microsoft.NewStaticTokenCredential(accessToken)

	p := microsoft.NewGraphProvider(cred, g.log)
	adapter, err := p.Adapter(ctx)
	if err != nil {
		return nil, err
	}

	// Using "me" as a shortcut for the authenticated user.
	respDrives, err := users.NewItemDrivesRequestBuilderInternal(map[string]string{"user%2Did": "me"}, adapter).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var out []drive.Drive
	pageIterator, err := msgraphsdkcore.NewPageIterator[models.Driveable](respDrives, adapter, models.CreateDriveCollectionResponseFromDiscriminatorValue)
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
	req := &proto.GetTokenRequest{
		IdentityId: identityID,
		Scopes:     []string{"Files.Read"},
	}
	resp, err := g.auth.Token(ctx, "microsoft", req) // Pass provider
	if err != nil {
		return drive.Drive{}, fmt.Errorf("failed to get token for identity %s: %w", identityID, err)
	}

	accessToken := identity.FromProtoAccessToken(resp.GetToken(), identityID)
	cred := microsoft.NewStaticTokenCredential(accessToken)

	p := microsoft.NewGraphProvider(cred, g.log)
	adapter, err := p.Adapter(ctx)
	if err != nil {
		return drive.Drive{}, err
	}

	respDrive, err := users.NewItemDriveRequestBuilderInternal(map[string]string{"user%2Did": "me"}, adapter).Get(ctx, nil)
	if err != nil {
		return drive.Drive{}, err
	}

	return toDrive(respDrive), nil
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
	// Logic for read-only could be inferred from permissions
	
	return drive.Drive{
		ID:       id,
		Name:     name,
		Type:     driveType,
		Owner:    owner,
		ReadOnly: readOnly,
	}
}
