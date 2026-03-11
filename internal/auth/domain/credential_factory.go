package domain

import domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"

type CredentialFactory interface {
	Credential(record domainaccount.Account, opts *Options) (CredentialProvider, error)
}
