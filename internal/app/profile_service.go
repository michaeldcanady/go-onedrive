package app

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/cache/abstractions"
)

const profileKey = "profile.json"

type ProfileService interface {
	// Load loads the profile from storage, or returns (nil, nil) if not found.
	Load(context.Context) (*azidentity.AuthenticationRecord, error)

	// Save persists the given profile. A nil profile could mean "delete/clear".
	Save(context.Context, *azidentity.AuthenticationRecord) error

	Clear(context.Context) error
}

type ProfileServiceImpl struct {
	store abstractions.Store
	codec abstractions.Codec
}

func NewProfileService(store abstractions.Store, codec abstractions.Codec) *ProfileServiceImpl {
	return &ProfileServiceImpl{
		store: store,
		codec: codec,
	}
}

func (c *ProfileServiceImpl) Load(ctx context.Context) (*azidentity.AuthenticationRecord, error) {
	data, err := c.store.LoadBytes(ctx, profileKey)
	if err != nil {
		return nil, fmt.Errorf("load profile bytes: %w", err)
	}
	if data == nil {
		// cache miss
		return nil, nil
	}

	var p azidentity.AuthenticationRecord
	if err := c.codec.Decode(data, &p); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}
	// if file existed but effectively empty/invalid -> treat as miss?
	if isZeroProfile(&p) {
		return nil, nil
	}

	return &p, nil
}

func (c *ProfileServiceImpl) Save(ctx context.Context, p *azidentity.AuthenticationRecord) error {
	if p == nil || isZeroProfile(p) {
		// You can choose to delete here in the future; for now, write empty object
		p = &azidentity.AuthenticationRecord{}
	}

	data, err := c.codec.Encode(p)
	if err != nil {
		return fmt.Errorf("encode profile: %w", err)
	}
	if err := c.store.SaveBytes(ctx, profileKey, data); err != nil {
		return fmt.Errorf("save profile bytes: %w", err)
	}
	return nil
}

func (c *ProfileServiceImpl) Clear(ctx context.Context) error {
	if err := c.store.Delete(ctx, profileKey); err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	return nil
}

func isZeroProfile(p *azidentity.AuthenticationRecord) bool {
	return p == nil
}
