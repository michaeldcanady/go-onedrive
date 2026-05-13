package drive

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

type DriveService struct {
	repo       Repository
	plugins    plugins.Manager
	identities IdentityService
	tokens     TokenService
	logger     logger.Service
}

// NewDriveService returns a new [*DriveService] initialized with the required domain dependencies.
// It uses the provided [plugins.Manager] to communicate with external storage providers.
func NewDriveService(repo Repository, pm plugins.Manager, is IdentityService, ts TokenService, l logger.Service) *DriveService {
	return &DriveService{
		repo:       repo,
		plugins:    pm,
		identities: is,
		tokens:     ts,
		logger:     l,
	}
}

func (s *DriveService) List(ctx context.Context, identityID string) ([]*Drive, error) {

	l := logger.WithContext(s.logger, ctx)

	// 1. Discover storage plugins
	allPlugins, err := s.plugins.ListPlugins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list plugins: %w", err)
	}

	var storagePlugins []*plugins.Metadata
	for _, p := range allPlugins {
		if p.Type == "storage" {
			storagePlugins = append(storagePlugins, p)
		}
	}

	// 2. Get target identities
	targetIdentities, err := s.getTargetIdentities(ctx, identityID)
	if err != nil {
		return nil, err
	}

	// 3. Aggregate results
	var allDrives []*Drive
	seenDrives := make(map[string]bool)

	addDrives := func(drives []*Drive) {
		for _, d := range drives {
			if !seenDrives[d.ID] {
				allDrives = append(allDrives, d)
				seenDrives[d.ID] = true
				if err := s.repo.Save(d); err != nil {
					l.Warn("failed to cache drive in repository", "drive_id", d.ID, "error", err)
				}
			}
		}
	}

	for _, p := range storagePlugins {
		if len(p.SupportedProviders) == 0 {
			// Case B: Plugin is provider-agnostic (e.g., local)
			drives, err := s.listFromPlugin(ctx, p.PluginPath, nil)
			if err != nil {
				l.Warn("failed to list drives from provider-agnostic plugin", "plugin", p.Name, "error", err)
				continue
			}
			addDrives(drives)
			continue
		}

		// Case A: Plugin supports specific providers
		for _, iden := range targetIdentities {
			if !s.isProviderSupported(p, iden.Provider) {
				continue
			}

			drives, err := s.listFromPlugin(ctx, p.PluginPath, iden)
			if err != nil {
				l.Warn("failed to list drives from plugin", "plugin", p.Name, "identity", iden.ID, "error", err)
				continue
			}
			addDrives(drives)
		}
	}

	if len(allDrives) == 0 && identityID != "" {
		// Try fallback to cache
		return s.repo.ListByIdentity(identityID)
	}

	return allDrives, nil
}

func (s *DriveService) getTargetIdentities(ctx context.Context, identityID string) ([]*Identity, error) {
	if identityID == "" {
		return s.identities.List(ctx)
	}

	iden, err := s.identities.GetIdentity(ctx, identityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity %s: %w", identityID, err)
	}
	if iden == nil {
		return nil, fmt.Errorf("identity %s not found", identityID)
	}
	return []*Identity{iden}, nil
}

func (s *DriveService) isProviderSupported(p *plugins.Metadata, provider string) bool {
	for _, sp := range p.SupportedProviders {
		if sp == provider {
			return true
		}
	}
	return false
}

func (s *DriveService) listFromPlugin(ctx context.Context, pluginPath string, iden *Identity) ([]*Drive, error) {
	client, err := s.plugins.GetStoragePlugin(pluginPath)
	if err != nil {
		return nil, err
	}

	options := make(map[string]string)
	if iden != nil {
		token, err := s.tokens.GetToken(ctx, iden.Provider, iden.ID)
		if err != nil {
			return nil, err
		}
		options["token"] = token.AccessToken
	}

	resp, err := client.ListDrives(ctx, &storage_proto.ListDrivesRequest{
		Options: options,
	})
	if err != nil {
		return nil, err
	}

	drives := make([]*Drive, len(resp.Drives))
	for i, d := range resp.Drives {
		idenID := ""
		if iden != nil {
			idenID = iden.ID
		}
		drives[i] = &Drive{
			ID:         d.Id,
			Name:       d.Name,
			IdentityID: idenID,
			Type:       d.Type,
		}
	}
	return drives, nil
}

func (s *DriveService) Get(ctx context.Context, driveID string) (*Drive, error) {
	return s.repo.ByID(driveID)
}

func (s *DriveService) FindDrive(ctx context.Context, query string) (*Drive, error) {
	// 1. Try exact ID match in cache
	if d, err := s.repo.ByID(query); err == nil && d != nil {
		return d, nil
	}

	// 2. Search by name in cache
	drives, err := s.repo.ListByIdentity("")
	if err != nil {
		return nil, err
	}

	for _, d := range drives {
		if d.Name == query {
			return d, nil
		}
	}

	// 3. If not found in cache, try listing from plugins to see if it's new
	// This might be expensive, so we only do it if not found in cache.
	// For now, let's just return nil if not in cache to keep it fast.
	return nil, nil
}
