package auth

import accountdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/account"

type CredentialFactory interface {
	Credential(record accountdomain.Account, opts *Options) (CredentialProvider, error)
}
