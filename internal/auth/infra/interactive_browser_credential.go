package infra

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
)

func newInteractiveBrowserCredential(record azidentity.AuthenticationRecord, opts *domainauth.Options) (*azidentity.InteractiveBrowserCredential, error) {

	azOpts := &azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord:           record,
		ClientID:                       opts.ClientID,
		TenantID:                       opts.TenantID,
		DisableAutomaticAuthentication: true,
	}

	return azidentity.NewInteractiveBrowserCredential(azOpts)
}
