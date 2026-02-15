package msal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
)

func newInteractiveBrowserCredential(record azidentity.AuthenticationRecord, opts *auth.Options) (*azidentity.InteractiveBrowserCredential, error) {

	azOpts := &azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord: record,
		ClientID:             opts.ClientID,
		TenantID:             opts.TenantID,
	}

	return azidentity.NewInteractiveBrowserCredential(azOpts)
}
