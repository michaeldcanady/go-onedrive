package cache

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/disk"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/memory"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type Service struct {
	// due to golang's runtime generics there is no type safe way to manage these caches.
	// this is our work around until a better solution is possible.
	profileCache       abstractions.Cache[string, azidentity.AuthenticationRecord]
	configurationCache abstractions.Cache[string, config.Configuration3]
	driveCache         abstractions.Cache[string, *domaincache.CachedChildren]
	fileCache          abstractions.Cache[string, *domaincache.CachedItem]
	logger             logging.Logger
}

func New(profileCachePath, driveCachePath, fileCachePath string, logger logging.Logger) (*Service, error) {
	parent, _ := filepath.Split(profileCachePath)
	if err := os.MkdirAll(parent, os.ModePerm); err != nil {
		return nil, err
	}

	profileCache, err := disk.New(profileCachePath, &JSONSerializerDeserializer[string]{}, &JSONSerializerDeserializer[azidentity.AuthenticationRecord]{})
	if err != nil {
		return nil, err
	}

	configurationCache := memory.New[*abstractions.Entry[string, config.Configuration3], string, config.Configuration3]()

	driveCache, err := disk.New(driveCachePath, &JSONSerializerDeserializer[string]{}, NewKiotaJSONSerializerDeserializer[*domaincache.CachedChildren](domaincache.CreateCachedChildrenFromDiscriminatorValue))
	if err != nil {
		return nil, err
	}

	fileCache, err := disk.New(fileCachePath, &JSONSerializerDeserializer[string]{}, NewKiotaJSONSerializerDeserializer[*domaincache.CachedItem](domaincache.CreateCachedChildrenFromDiscriminatorValue))
	if err != nil {
		return nil, err
	}

	return &Service{
		profileCache:       profileCache,
		configurationCache: configurationCache,
		logger:             logger,
		driveCache:         driveCache,
		fileCache:          fileCache,
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

func (s *Service) GetConfiguration(ctx context.Context, name string) (config.Configuration3, error) {
	var record config.Configuration3

	if err := ctx.Err(); err != nil {
		return record, err
	}

	if s.configurationCache == nil {
		return record, errors.New("configuration cache is nil")
	}

	entry, err := s.configurationCache.GetEntry(ctx, name)
	if err != nil {
		return record, err
	}

	if entry != nil {
		record = entry.GetValue()
	}

	return record, nil
}

func (s *Service) SetConfiguration(ctx context.Context, name string, record config.Configuration3) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.configurationCache == nil {
		return errors.New("configuration cache is nil")
	}

	entry, err := s.configurationCache.GetEntry(ctx, name)
	if err != nil {
		// ok if key isn't found as that means no configuration cached previously
		if !errors.Is(err, core.ErrKeyNotFound) {
			return err
		}
		if entry, err = s.configurationCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	if entry == nil {
		if entry, err = s.configurationCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	entry.SetValue(record)

	return s.configurationCache.SetEntry(ctx, entry)
}

func (s *Service) GetCLIProfile(ctx context.Context, name string) (domainprofile.Profile, error) {
	var profile domainprofile.Profile
	if err := ctx.Err(); err != nil {
		return profile, err
	}

	return profile, nil
}

func (s *Service) SetCLIProfile(ctx context.Context, name string, profile domainprofile.Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetDrive(ctx context.Context, name string) (domaincache.CachedChildren, error) {
	var record domaincache.CachedChildren
	if err := ctx.Err(); err != nil {
		return record, err
	}

	if s.driveCache == nil {
		return record, errors.New("profile cache is nil")
	}

	entry, err := s.driveCache.GetEntry(ctx, name)
	if err != nil {
		// ok if key isn't found as that means no profile cached
		if !errors.Is(err, core.ErrKeyNotFound) {
			return record, errors.Join(errors.New("unable to retrieve profile"), err)
		}
	}
	if entry != nil {
		val := entry.GetValue()
		if val != nil {
			record = *val
		}
	}

	return record, nil
}

func (s *Service) SetDrive(ctx context.Context, name string, record domaincache.CachedChildren) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.driveCache == nil {
		return errors.New("profile cache is nil")
	}

	entry, err := s.driveCache.GetEntry(ctx, name)
	if err != nil {
		// ok if key isn't found as that means no profile cached previously
		if !errors.Is(err, core.ErrKeyNotFound) {
			return err
		}
		if entry, err = s.driveCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	if entry == nil {
		if entry, err = s.driveCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	entry.SetValue(&record)

	return s.driveCache.SetEntry(ctx, entry)
}

func (s *Service) GetItem(ctx context.Context, name string) (domaincache.CachedItem, error) {
	var record domaincache.CachedItem
	if err := ctx.Err(); err != nil {
		return record, err
	}

	if s.fileCache == nil {
		return record, errors.New("profile cache is nil")
	}

	entry, err := s.fileCache.GetEntry(ctx, name)
	if err != nil {
		// ok if key isn't found as that means no profile cached
		if !errors.Is(err, core.ErrKeyNotFound) {
			return record, errors.Join(errors.New("unable to retrieve profile"), err)
		}
	}
	if entry != nil {
		val := entry.GetValue()
		if val != nil {
			record = *val
		}
	}

	return record, nil
}

func (s *Service) SetItem(ctx context.Context, name string, record domaincache.CachedItem) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.fileCache == nil {
		return errors.New("profile cache is nil")
	}

	entry, err := s.fileCache.GetEntry(ctx, name)
	if err != nil {
		// ok if key isn't found as that means no profile cached previously
		if !errors.Is(err, core.ErrKeyNotFound) {
			return err
		}
		if entry, err = s.fileCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	if entry == nil {
		if entry, err = s.fileCache.NewEntry(ctx, name); err != nil {
			return err
		}
	}
	entry.SetValue(&record)

	return s.fileCache.SetEntry(ctx, entry)
}
