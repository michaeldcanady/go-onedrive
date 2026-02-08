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
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Service struct {
	// due to golang's runtime generics there is no type safe way to manage these caches.
	// this is our work around until a better solution is possible.
	authCache          abstractions.Cache[string, azidentity.AuthenticationRecord]
	configurationCache abstractions.Cache[string, config.Configuration3]
	driveCache         abstractions.Cache[string, *domaincache.CachedChildren]
	fileCache          abstractions.Cache[string, *domaincache.CachedItem]
	logger             logging.Logger
	profileCache       abstractions.Cache[string, domainprofile.Profile]
}

func New(authCachePath, driveCachePath, fileCachePath string, logger logging.Logger) (*Service, error) {
	logger.Info("initializing cache service",
		logging.String("event", "cache_init"),
		logging.String("auth_cache_path", authCachePath),
		logging.String("drive_cache_path", driveCachePath),
		logging.String("file_cache_path", fileCachePath),
	)

	parent, _ := filepath.Split(authCachePath)
	if err := os.MkdirAll(parent, os.ModePerm); err != nil {
		logger.Error("failed to create cache directory",
			logging.String("event", "cache_init"),
			logging.Error(err),
		)
		return nil, err
	}

	authCache, err := disk.New(authCachePath, &JSONSerializerDeserializer[string]{}, &JSONSerializerDeserializer[azidentity.AuthenticationRecord]{})
	if err != nil {
		logger.Error("failed to initialize auth cache",
			logging.String("event", "cache_init"),
			logging.Error(err),
		)
		return nil, err
	}

	configurationCache := memory.New[*abstractions.Entry[string, config.Configuration3]]()

	driveCache, err := disk.New(driveCachePath, &JSONSerializerDeserializer[string]{}, NewKiotaJSONSerializerDeserializer[*domaincache.CachedChildren](domaincache.CreateCachedChildrenFromDiscriminatorValue))
	if err != nil {
		logger.Error("failed to initialize drive cache",
			logging.String("event", "cache_init"),
			logging.Error(err),
		)
		return nil, err
	}

	fileCache, err := disk.New(fileCachePath, &JSONSerializerDeserializer[string]{}, NewKiotaJSONSerializerDeserializer[*domaincache.CachedItem](domaincache.CreateCachedChildrenFromDiscriminatorValue))
	if err != nil {
		logger.Error("failed to initialize file cache",
			logging.String("event", "cache_init"),
			logging.Error(err),
		)
		return nil, err
	}

	profileCache := memory.New[*abstractions.Entry[string, domainprofile.Profile]]()

	logger.Info("cache service initialized successfully",
		logging.String("event", "cache_init"),
	)

	return &Service{
		authCache:          authCache,
		configurationCache: configurationCache,
		logger:             logger,
		driveCache:         driveCache,
		fileCache:          fileCache,
		profileCache:       profileCache,
	}, nil
}

// GetProfile returns the currently cached profile by name.
func (s *Service) GetProfile(ctx context.Context, name string) (azidentity.AuthenticationRecord, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("retrieving cached auth profile",
		logging.String("event", "cache_get_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	var record azidentity.AuthenticationRecord
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while retrieving profile",
			logging.String("event", "cache_get_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return record, err
	}

	if s.authCache == nil {
		s.logger.Error("auth cache is nil",
			logging.String("event", "cache_get_profile"),
			logging.String("correlation_id", cid),
		)
		return record, errors.New("profile cache is nil")
	}

	entry, err := s.authCache.GetEntry(ctx, name)
	if err != nil {
		if errors.Is(err, core.ErrKeyNotFound) {
			s.logger.Debug("auth profile not found",
				logging.String("event", "cache_get_profile"),
				logging.String("profile", name),
				logging.String("correlation_id", cid),
			)
			return record, nil
		}

		s.logger.Error("failed to retrieve auth profile",
			logging.String("event", "cache_get_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return record, errors.Join(errors.New("unable to retrieve profile"), err)
	}

	record = entry.GetValue()

	s.logger.Debug("auth profile retrieved",
		logging.String("event", "cache_get_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return record, nil
}

// SetProfile caches the provided profile by name.
func (s *Service) SetProfile(ctx context.Context, name string, record azidentity.AuthenticationRecord) error {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("saving auth profile",
		logging.String("event", "cache_set_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while saving profile",
			logging.String("event", "cache_set_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if s.authCache == nil {
		s.logger.Error("auth cache is nil",
			logging.String("event", "cache_set_profile"),
			logging.String("correlation_id", cid),
		)
		return errors.New("profile cache is nil")
	}

	entry, err := s.authCache.GetEntry(ctx, name)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		s.logger.Error("failed to retrieve existing profile entry",
			logging.String("event", "cache_set_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if entry == nil {
		entry, err = s.authCache.NewEntry(ctx, name)
		if err != nil {
			s.logger.Error("failed to create new profile entry",
				logging.String("event", "cache_set_profile"),
				logging.Error(err),
				logging.String("correlation_id", cid),
			)
			return err
		}
	}

	entry.SetValue(record)

	if err := s.authCache.SetEntry(ctx, entry); err != nil {
		s.logger.Error("failed to persist auth profile",
			logging.String("event", "cache_set_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("auth profile saved",
		logging.String("event", "cache_set_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return nil
}

func (s *Service) DeleteProfile(ctx context.Context, name string) error {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Info("deleting auth profile",
		logging.String("event", "cache_delete_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while deleting profile",
			logging.String("event", "cache_delete_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if s.authCache == nil {
		s.logger.Error("auth cache is nil",
			logging.String("event", "cache_delete_profile"),
			logging.String("correlation_id", cid),
		)
		return errors.New("profile cache is nil")
	}

	if err := s.authCache.Remove(name); err != nil {
		s.logger.Error("failed to delete auth profile",
			logging.String("event", "cache_delete_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("auth profile deleted",
		logging.String("event", "cache_delete_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return nil
}

func (s *Service) GetConfiguration(ctx context.Context, name string) (config.Configuration3, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("retrieving configuration",
		logging.String("event", "cache_get_config"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	var record config.Configuration3

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while retrieving configuration",
			logging.String("event", "cache_get_config"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return record, err
	}

	// Cache missing
	if s.configurationCache == nil {
		s.logger.Error("configuration cache is nil",
			logging.String("event", "cache_get_config"),
			logging.String("correlation_id", cid),
		)
		return record, errors.New("configuration cache is nil")
	}

	// Retrieve entry
	entry, err := s.configurationCache.GetEntry(ctx, name)
	if err != nil {
		if errors.Is(err, core.ErrKeyNotFound) {
			s.logger.Debug("configuration not found",
				logging.String("event", "cache_get_config"),
				logging.String("profile", name),
				logging.String("correlation_id", cid),
			)
			return record, nil
		}

		s.logger.Error("failed to retrieve configuration",
			logging.String("event", "cache_get_config"),
			logging.Error(err),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return record, err
	}

	if entry != nil {
		record = entry.GetValue()
	}

	s.logger.Debug("configuration retrieved",
		logging.String("event", "cache_get_config"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return record, nil
}

func (s *Service) SetConfiguration(ctx context.Context, name string, record config.Configuration3) error {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("saving configuration",
		logging.String("event", "cache_set_config"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while saving configuration",
			logging.String("event", "cache_set_config"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	// Cache missing
	if s.configurationCache == nil {
		s.logger.Error("configuration cache is nil",
			logging.String("event", "cache_set_config"),
			logging.String("correlation_id", cid),
		)
		return errors.New("configuration cache is nil")
	}

	// Retrieve or create entry
	entry, err := s.configurationCache.GetEntry(ctx, name)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		s.logger.Error("failed to retrieve existing configuration entry",
			logging.String("event", "cache_set_config"),
			logging.Error(err),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if entry == nil {
		s.logger.Debug("creating new configuration entry",
			logging.String("event", "cache_set_config"),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)

		entry, err = s.configurationCache.NewEntry(ctx, name)
		if err != nil {
			s.logger.Error("failed to create configuration entry",
				logging.String("event", "cache_set_config"),
				logging.Error(err),
				logging.String("profile", name),
				logging.String("correlation_id", cid),
			)
			return err
		}
	}

	// Set value
	entry.SetValue(record)

	// Persist entry
	if err := s.configurationCache.SetEntry(ctx, entry); err != nil {
		s.logger.Error("failed to persist configuration",
			logging.String("event", "cache_set_config"),
			logging.Error(err),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("configuration saved",
		logging.String("event", "cache_set_config"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return nil
}

func (s *Service) GetCLIProfile(ctx context.Context, name string) (domainprofile.Profile, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("retrieving CLI profile",
		logging.String("event", "cache_get_cli_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	var profile domainprofile.Profile

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while retrieving CLI profile",
			logging.String("event", "cache_get_cli_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return profile, err
	}

	// Cache missing
	if s.profileCache == nil {
		s.logger.Error("CLI profile cache is nil",
			logging.String("event", "cache_get_cli_profile"),
			logging.String("correlation_id", cid),
		)
		return profile, errors.New("profile cache is nil")
	}

	// Retrieve entry
	entry, err := s.profileCache.GetEntry(ctx, name)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		s.logger.Error("failed to retrieve CLI profile",
			logging.String("event", "cache_get_cli_profile"),
			logging.Error(err),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return profile, errors.Join(errors.New("unable to retrieve profile"), err)
	}

	// Cache miss
	if entry == nil {
		s.logger.Debug("CLI profile not found",
			logging.String("event", "cache_get_cli_profile"),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return profile, nil
	}

	// Cache hit
	profile = entry.GetValue()
	if profile == (domainprofile.Profile{}) {
		s.logger.Warn("CLI profile entry was empty",
			logging.String("event", "cache_get_cli_profile"),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return domainprofile.Profile{}, errors.New("empty profile entry")
	}

	s.logger.Debug("CLI profile retrieved",
		logging.String("event", "cache_get_cli_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return profile, nil
}

func (s *Service) SetCLIProfile(ctx context.Context, name string, profile domainprofile.Profile) error {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("saving CLI profile",
		logging.String("event", "cache_set_cli_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while saving CLI profile",
			logging.String("event", "cache_set_cli_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	// Cache missing
	if s.profileCache == nil {
		s.logger.Error("CLI profile cache is nil",
			logging.String("event", "cache_set_cli_profile"),
			logging.String("correlation_id", cid),
		)
		return errors.New("profile cache is nil")
	}

	// Retrieve or create entry
	entry, err := s.profileCache.GetEntry(ctx, name)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		s.logger.Error("failed to retrieve existing CLI profile entry",
			logging.String("event", "cache_set_cli_profile"),
			logging.Error(err),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if entry == nil {
		s.logger.Debug("creating new CLI profile entry",
			logging.String("event", "cache_set_cli_profile"),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)

		entry, err = s.profileCache.NewEntry(ctx, name)
		if err != nil {
			s.logger.Error("failed to create CLI profile entry",
				logging.String("event", "cache_set_cli_profile"),
				logging.Error(err),
				logging.String("profile", name),
				logging.String("correlation_id", cid),
			)
			return err
		}
	}

	// Set value
	entry.SetValue(profile)

	// Persist entry
	if err := s.profileCache.SetEntry(ctx, entry); err != nil {
		s.logger.Error("failed to persist CLI profile",
			logging.String("event", "cache_set_cli_profile"),
			logging.Error(err),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("CLI profile saved",
		logging.String("event", "cache_set_cli_profile"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return nil
}

func (s *Service) GetDrive(ctx context.Context, name string) (domaincache.CachedChildren, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("retrieving cached drive",
		logging.String("event", "cache_get_drive"),
		logging.String("drive", name),
		logging.String("correlation_id", cid),
	)

	var record domaincache.CachedChildren

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while retrieving drive",
			logging.String("event", "cache_get_drive"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return record, err
	}

	// Cache missing
	if s.driveCache == nil {
		s.logger.Error("drive cache is nil",
			logging.String("event", "cache_get_drive"),
			logging.String("correlation_id", cid),
		)
		return record, errors.New("drive cache is nil")
	}

	// Retrieve entry
	entry, err := s.driveCache.GetEntry(ctx, name)
	if err != nil {
		if errors.Is(err, core.ErrKeyNotFound) {
			s.logger.Debug("drive not found in cache",
				logging.String("event", "cache_get_drive"),
				logging.String("drive", name),
				logging.String("correlation_id", cid),
			)
			return record, nil
		}

		s.logger.Error("failed to retrieve drive",
			logging.String("event", "cache_get_drive"),
			logging.Error(err),
			logging.String("drive", name),
			logging.String("correlation_id", cid),
		)
		return record, errors.Join(errors.New("unable to retrieve drive"), err)
	}

	// Cache hit
	if entry != nil {
		val := entry.GetValue()
		if val != nil {
			record = *val
		}
	}

	s.logger.Debug("drive retrieved",
		logging.String("event", "cache_get_drive"),
		logging.String("drive", name),
		logging.String("correlation_id", cid),
	)

	return record, nil
}

func (s *Service) SetDrive(ctx context.Context, name string, record domaincache.CachedChildren) error {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("saving cached drive",
		logging.String("event", "cache_set_drive"),
		logging.String("drive", name),
		logging.String("correlation_id", cid),
	)

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while saving drive",
			logging.String("event", "cache_set_drive"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	// Cache missing
	if s.driveCache == nil {
		s.logger.Error("drive cache is nil",
			logging.String("event", "cache_set_drive"),
			logging.String("correlation_id", cid),
		)
		return errors.New("drive cache is nil")
	}

	// Retrieve or create entry
	entry, err := s.driveCache.GetEntry(ctx, name)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		s.logger.Error("failed to retrieve existing drive entry",
			logging.String("event", "cache_set_drive"),
			logging.Error(err),
			logging.String("drive", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if entry == nil {
		s.logger.Debug("creating new drive entry",
			logging.String("event", "cache_set_drive"),
			logging.String("drive", name),
			logging.String("correlation_id", cid),
		)

		entry, err = s.driveCache.NewEntry(ctx, name)
		if err != nil {
			s.logger.Error("failed to create drive entry",
				logging.String("event", "cache_set_drive"),
				logging.Error(err),
				logging.String("drive", name),
				logging.String("correlation_id", cid),
			)
			return err
		}
	}

	// Set value
	entry.SetValue(&record)

	// Persist entry
	if err := s.driveCache.SetEntry(ctx, entry); err != nil {
		s.logger.Error("failed to persist drive",
			logging.String("event", "cache_set_drive"),
			logging.Error(err),
			logging.String("drive", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("drive cached",
		logging.String("event", "cache_set_drive"),
		logging.String("drive", name),
		logging.String("correlation_id", cid),
	)

	return nil
}

func (s *Service) GetItem(ctx context.Context, name string) (domaincache.CachedItem, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("retrieving cached item",
		logging.String("event", "cache_get_item"),
		logging.String("item", name),
		logging.String("correlation_id", cid),
	)

	var record domaincache.CachedItem

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while retrieving item",
			logging.String("event", "cache_get_item"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return record, err
	}

	// Cache missing
	if s.fileCache == nil {
		s.logger.Error("file cache is nil",
			logging.String("event", "cache_get_item"),
			logging.String("correlation_id", cid),
		)
		return record, errors.New("file cache is nil")
	}

	// Retrieve entry
	entry, err := s.fileCache.GetEntry(ctx, name)
	if err != nil {
		if errors.Is(err, core.ErrKeyNotFound) {
			s.logger.Debug("item not found in cache",
				logging.String("event", "cache_get_item"),
				logging.String("item", name),
				logging.String("correlation_id", cid),
			)
			return record, nil
		}

		s.logger.Error("failed to retrieve item",
			logging.String("event", "cache_get_item"),
			logging.Error(err),
			logging.String("item", name),
			logging.String("correlation_id", cid),
		)
		return record, errors.Join(errors.New("unable to retrieve item"), err)
	}

	// Cache hit
	if entry != nil {
		val := entry.GetValue()
		if val != nil {
			record = *val
		}
	}

	s.logger.Debug("item retrieved",
		logging.String("event", "cache_get_item"),
		logging.String("item", name),
		logging.String("correlation_id", cid),
	)

	return record, nil
}

func (s *Service) SetItem(ctx context.Context, name string, record domaincache.CachedItem) error {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("saving cached item",
		logging.String("event", "cache_set_item"),
		logging.String("item", name),
		logging.String("correlation_id", cid),
	)

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while saving item",
			logging.String("event", "cache_set_item"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	// Cache missing
	if s.fileCache == nil {
		s.logger.Error("file cache is nil",
			logging.String("event", "cache_set_item"),
			logging.String("correlation_id", cid),
		)
		return errors.New("file cache is nil")
	}

	// Retrieve or create entry
	entry, err := s.fileCache.GetEntry(ctx, name)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		s.logger.Error("failed to retrieve existing item entry",
			logging.String("event", "cache_set_item"),
			logging.Error(err),
			logging.String("item", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	if entry == nil {
		s.logger.Debug("creating new item entry",
			logging.String("event", "cache_set_item"),
			logging.String("item", name),
			logging.String("correlation_id", cid),
		)

		entry, err = s.fileCache.NewEntry(ctx, name)
		if err != nil {
			s.logger.Error("failed to create item entry",
				logging.String("event", "cache_set_item"),
				logging.Error(err),
				logging.String("item", name),
				logging.String("correlation_id", cid),
			)
			return err
		}
	}

	// Set value
	entry.SetValue(&record)

	// Persist entry
	if err := s.fileCache.SetEntry(ctx, entry); err != nil {
		s.logger.Error("failed to persist item",
			logging.String("event", "cache_set_item"),
			logging.Error(err),
			logging.String("item", name),
			logging.String("correlation_id", cid),
		)
		return err
	}

	s.logger.Info("item cached",
		logging.String("event", "cache_set_item"),
		logging.String("item", name),
		logging.String("correlation_id", cid),
	)

	return nil
}
