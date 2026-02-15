package auth

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	accountdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/account"
)

type CredentialOptions struct {
	Type        string
	ClientID    string
	TenantID    string
	RedirectURI string
	// DEPRECATED: use [Account] instead.
	AuthenticationRecord azidentity.AuthenticationRecord
	Account              accountdomain.Account
}
