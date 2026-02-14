package msal

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func NewMSALCredentialFactory() *MSALCredentialFactory {
	return &MSALCredentialFactory{}
}

type MSALCredentialFactory struct{}

func (_ *MSALCredentialFactory) Credential(record azidentity.AuthenticationRecord, opts *Options) (CredentialProvider, error) {
	switch opts.Method {
	case MethodInteractiveBrowser:
		return newInteractiveBrowserCredential(record, opts)
	}
	return nil, errors.New("unsupported token provider")
}
