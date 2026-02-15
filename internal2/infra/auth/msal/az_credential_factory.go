package msal

import (
	"errors"

	accountdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
)

func NewMSALCredentialFactory() *MSALCredentialFactory {
	return &MSALCredentialFactory{}
}

type MSALCredentialFactory struct{}

func (_ *MSALCredentialFactory) Credential(account accountdomain.Account, opts *auth.Options) (auth.CredentialProvider, error) {
	record := accountdomain.AccountToMSAuthRecord(account)

	switch opts.Method {
	case auth.MethodInteractiveBrowser:
		return newInteractiveBrowserCredential(record, opts)
	}
	return nil, errors.New("unsupported token provider")
}
