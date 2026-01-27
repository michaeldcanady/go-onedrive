package di2

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type ProfileService interface {
	// Load loads the profile from storage, or returns (nil, nil) if not found.
	Load(context.Context) (*azidentity.AuthenticationRecord, error)

	// Save persists the given profile. A nil profile could mean "delete/clear".
	Save(context.Context, *azidentity.AuthenticationRecord) error
}
