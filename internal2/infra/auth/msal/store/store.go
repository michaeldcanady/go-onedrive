package msal

import (
	"context"
	"encoding/json"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"

	infracache "github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

type AccountRecord struct {
	HomeAccountID     string `json:"home_account_id"`
	Environment       string `json:"environment"`
	Realm             string `json:"realm"`
	PreferredUsername string `json:"preferred_username"`
}

func FromAccount(a public.Account) AccountRecord {
	return AccountRecord{
		HomeAccountID:     a.HomeAccountID,
		Environment:       a.Environment,
		Realm:             a.Realm,
		PreferredUsername: a.PreferredUsername,
	}
}

func (r AccountRecord) ToAccount() public.Account {
	return public.Account{
		HomeAccountID:     r.HomeAccountID,
		Environment:       r.Environment,
		Realm:             r.Realm,
		PreferredUsername: r.PreferredUsername,
	}
}

type AccountStore interface {
	Load(ctx context.Context, profile string) (public.Account, error)
	Save(ctx context.Context, profile string, acct public.Account) error
	Delete(ctx context.Context, profile string) error
}

type cacheBackedStore struct {
	cache infracache.Cache2
}

func NewAccountStore(cache infracache.Cache2) AccountStore {
	return &cacheBackedStore{cache: cache}
}

func (s *cacheBackedStore) Load(ctx context.Context, profile string) (public.Account, error) {
	var rec AccountRecord
	if err := s.cache.Get(ctx, func() ([]byte, error) { return json.Marshal(profile) }, func(data []byte) error { return json.Unmarshal(data, &rec) }); err != nil {
		return public.Account{}, err
	}
	if rec == (AccountRecord{}) || rec.HomeAccountID == "" {
		return public.Account{}, ErrAccountNotFound
	}
	return rec.ToAccount(), nil
}

func (s *cacheBackedStore) Save(ctx context.Context, profile string, acct public.Account) error {
	rec := FromAccount(acct)
	return s.cache.Set(ctx, func() ([]byte, error) { return json.Marshal(profile) }, func() ([]byte, error) { return json.Marshal(rec) })
}

func (s *cacheBackedStore) Delete(ctx context.Context, profile string) error {
	return s.cache.Delete(ctx, func() ([]byte, error) { return json.Marshal(profile) })
}
