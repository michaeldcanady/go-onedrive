package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
)

type identityServiceAdapter struct {
	identity identity.Service
}

// NewIdentityServiceAdapter returns a new [IdentityService] that wraps an [identity.Service].
func NewIdentityServiceAdapter(is identity.Service) IdentityService {
	return &identityServiceAdapter{identity: is}
}

func (a *identityServiceAdapter) GetIdentity(ctx context.Context, id string) (*Identity, error) {
	iden, err := a.identity.GetIdentity(ctx, id)
	if err != nil {
		return nil, err
	}
	if iden == nil {
		return nil, nil
	}
	return &Identity{
		ID:       iden.ID,
		Provider: iden.Provider,
	}, nil
}

func (a *identityServiceAdapter) List(ctx context.Context) ([]*Identity, error) {
	identities, err := a.identity.List(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]*Identity, len(identities))
	for i, iden := range identities {
		results[i] = &Identity{
			ID:       iden.ID,
			Provider: iden.Provider,
		}
	}
	return results, nil
}

type tokenServiceAdapter struct {
	token identity.TokenService
}

// NewTokenServiceAdapter returns a new [TokenService] that wraps an [identity.TokenService].
func NewTokenServiceAdapter(ts identity.TokenService) TokenService {
	return &tokenServiceAdapter{token: ts}
}

func (a *tokenServiceAdapter) GetToken(ctx context.Context, provider, identityID string) (*Token, error) {
	t, err := a.token.GetToken(ctx, provider, identityID)
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken: t.AccessToken,
	}, nil
}
