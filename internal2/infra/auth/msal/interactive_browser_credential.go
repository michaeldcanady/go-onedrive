package msal

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

func newInteractiveBrowserCredential(record azidentity.AuthenticationRecord, opts *Options) (*azidentity.InteractiveBrowserCredential, error) {

	azOpts := &azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord: record,
		ClientID:             opts.ClientID,
		TenantID:             opts.TenantID,
	}

	return azidentity.NewInteractiveBrowserCredential(azOpts)
}
