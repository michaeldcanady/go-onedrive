package auth

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

type CredentialFactory interface {
	Credential(record azidentity.AuthenticationRecord, opts *Options) (CredentialProvider, error)
}
