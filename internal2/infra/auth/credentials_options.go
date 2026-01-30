package auth

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

type CredentialOptions struct {
	Type                 string
	ClientID             string
	TenantID             string
	RedirectURI          string
	AuthenticationRecord azidentity.AuthenticationRecord
}
