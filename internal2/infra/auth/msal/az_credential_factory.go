package msal

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
)

func NewMSALCredentialFactory() *MSALCredentialFactory {
	return &MSALCredentialFactory{}
}

type MSALCredentialFactory struct{}

func (_ *MSALCredentialFactory) Credential(record azidentity.AuthenticationRecord, opts *auth.Options) (auth.CredentialProvider, error) {
	switch opts.Method {
	case auth.MethodInteractiveBrowser:
		return newInteractiveBrowserCredential(record, opts)
	}
	return nil, errors.New("unsupported token provider")
}
