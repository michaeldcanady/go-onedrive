package infra

import (
	"errors"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
)

func NewMSALCredentialFactory() *MSALCredentialFactory {
	return &MSALCredentialFactory{}
}

type MSALCredentialFactory struct{}

func (_ *MSALCredentialFactory) Credential(account domainaccount.Account, opts *domainauth.Options) (domainauth.CredentialProvider, error) {
	record := domainaccount.AccountToMSAuthRecord(account)

	switch opts.Method {
	case domainauth.MethodInteractiveBrowser:
		return newInteractiveBrowserCredential(record, opts)
	}
	return nil, errors.New("unsupported token provider")
}
