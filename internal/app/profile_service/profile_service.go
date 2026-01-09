package profileservice

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const profileKey = "profile.json"

type Service struct {
	store     abstractions.Store
	codec     abstractions.Codec
	publisher event.Publisher
	logger    logging.Logger
}

// New creates a new instance of ProfileService.
func New(store abstractions.Store, codec abstractions.Codec, publisher event.Publisher, logger logging.Logger) *Service {
	return &Service{
		store:     store,
		codec:     codec,
		publisher: publisher,
		logger:    logger,
	}
}

// Load loads the profile from storage, or returns (nil, nil) if not found.
func (c *Service) Load(ctx context.Context) (*azidentity.AuthenticationRecord, error) {
	data, err := c.store.LoadBytes(ctx, profileKey)
	if err != nil {
		c.logger.Error("failed to load profile bytes", logging.Any("error", err))
		return nil, fmt.Errorf("load profile bytes: %w", err)
	}

	if data == nil {
		c.logger.Debug("profile not found in store")
		return nil, nil
	}

	var p azidentity.AuthenticationRecord
	if err := c.codec.Decode(data, &p); err != nil {
		c.logger.Error("failed to decode profile", logging.Any("error", err))
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	if isZeroProfile(&p) {
		c.logger.Debug("profile exists but is zero/invalid")
		return nil, nil
	}

	if c.publisher == nil {
		c.logger.Warn("no event publisher configured; skipping profile.loaded event")
	} else {
		c.logger.Debug("publishing profile.loaded event")
		evt := newProfileLoadedEvent(p)
		if err := c.publisher.Publish(ctx, evt); err != nil {
			c.logger.Warn("failed to publish profile.loaded event", logging.Any("error", err))
		}
	}

	c.logger.Info("profile loaded successfully")

	return &p, nil
}

// Save persists the given profile.
func (c *Service) Save(ctx context.Context, p *azidentity.AuthenticationRecord) error {
	if p == nil || isZeroProfile(p) {
		c.logger.Debug("saving zero profile; writing empty record")
		p = &azidentity.AuthenticationRecord{}
	}

	data, err := c.codec.Encode(p)
	if err != nil {
		c.logger.Error("failed to encode profile", logging.Any("error", err))
		return fmt.Errorf("encode profile: %w", err)
	}

	if err := c.store.SaveBytes(ctx, profileKey, data); err != nil {
		c.logger.Error("failed to save profile bytes", logging.Any("error", err))
		return fmt.Errorf("save profile bytes: %w", err)
	}

	if c.publisher == nil {
		c.logger.Warn("no event publisher configured; skipping profile.saved event")
	} else {
		c.logger.Debug("publishing profile.saved event")
		evt := newProfileSavedEvent(*p)
		if err := c.publisher.Publish(ctx, evt); err != nil {
			c.logger.Warn("failed to publish profile.saved event", logging.Any("error", err))
		}
	}

	c.logger.Info("profile saved successfully")

	return nil
}

func isZeroProfile(p *azidentity.AuthenticationRecord) bool {
	return p == nil
}
