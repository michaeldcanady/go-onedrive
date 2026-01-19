package cacheservice

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/cachev2/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal/cachev2/core"
	"github.com/michaeldcanady/go-onedrive/internal/cachev2/disk"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

type Service struct {
	profileCache abstractions.Cache[string, azidentity.AuthenticationRecord]
	logger       logging.Logger
}

func New(cachePath string, logger logging.Logger) (*Service, error) {
	parent, _ := filepath.Split(cachePath)
	if err := os.MkdirAll(parent, os.ModePerm); err != nil {
		return nil, err
	}

	profileCache, err := disk.New(cachePath, &JSONSerializerDeserializer[string]{}, &JSONSerializerDeserializer[azidentity.AuthenticationRecord]{})
	if err != nil {
		return nil, err
	}

	return &Service{
		profileCache: profileCache,
		logger:       logger,
	}, nil
}

// GetProfile returns the currently cached profile by name.
func (s *Service) GetProfile(ctx context.Context, name string) (azidentity.AuthenticationRecord, error) {
	var record azidentity.AuthenticationRecord
	if err := ctx.Err(); err != nil {
		return record, err
	}

	if s.profileCache == nil {
		return record, errors.New("profile cache is nil")
	}

	entry, err := s.profileCache.GetEntry(ctx, name)
	if err != nil {
		// ok if key isn't found as that means no profile cached
		if !errors.Is(err, core.ErrKeyNotFound) {
			return record, errors.Join(errors.New("unable to retrieve profile"), err)
		}
	}
	if entry != nil {
		record = entry.GetValue()
	}

	return record, nil
}

// SetProfile caches the provided profile by name.
func (s *Service) SetProfile(ctx context.Context, name string, record azidentity.AuthenticationRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.profileCache == nil {
		return errors.New("profile cache is nil")
	}

	entry, err := s.profileCache.GetEntry(ctx, name)
	if err != nil {
		// ok if key isn't found as that means no profile cached previously
		if !errors.Is(err, core.ErrKeyNotFound) {
			return err
		}
		if entry, err = s.profileCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	if entry == nil {
		if entry, err = s.profileCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	entry.SetValue(record)

	return s.profileCache.SetEntry(ctx, entry)
}
