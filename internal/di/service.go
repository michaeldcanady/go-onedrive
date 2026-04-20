package di

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	registry "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/mount"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/local"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/onedrive"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger/zap"
	bolt "go.etcd.io/bbolt"
)

// DefaultContainer provides a concrete implementation of the Container interface.
// It orchestrates the lifecycle and wiring of core application services.
type DefaultContainer struct {
	// logger is the centralized logging service.
	logger logger.Service
	// config is the service for managing application settings.
	config config.Service
	// mounts is the service for managing VFS mount points.
	mounts mount.Service
	// identityService is the registry for managing multiple identity providers, providing both authenticators and authorizers.
	identityService identity.Service
	// profile is the service for managing user configuration profiles.
	profile profile.Service
	// manager is the orchestrated filesystem manager.
	manager registry.Service
	// environment is the environment-related service.
	environment environment.Service
	// editor is the external editor service.
	editor editor.Service
	// drive is the OneDrive drive management service.
	drive drive.Service
	// uriFactory is the service for creating and resolving URIs.
	uriFactory *registry.URIFactory
}

// NewDefaultContainer initializes a new instance of the DefaultContainer with all core services wired.
func NewDefaultContainer() (*DefaultContainer, error) {
	ctx := context.Background()
	c := &DefaultContainer{}

	if err := c.initBaseServices(); err != nil {
		return nil, err
	}

	if err := c.initIdentityServices(ctx); err != nil {
		return nil, err
	}

	if err := c.initVFSServices(ctx); err != nil {
		return nil, err
	}

	if err := c.initDriveServices(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *DefaultContainer) initBaseServices() error {
	c.environment = environment.NewDefaultService("odc")
	if err := c.environment.EnsureAll(); err != nil {
		return fmt.Errorf("failed to ensure environment directories: %w", err)
	}

	c.logger = zap.NewZapService(c.environment)

	profileSvc, err := profile.NewDefaultService(c.environment)
	if err != nil {
		return fmt.Errorf("failed to initialize profile service: %w", err)
	}
	c.profile = profileSvc

	return nil
}

func (c *DefaultContainer) initIdentityServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")

	stateDir, err := c.environment.StateDir()
	if err != nil {
		return fmt.Errorf("failed to get state directory: %w", err)
	}
	dbPath := filepath.Join(stateDir, "identity.db")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to open identity database: %w", err)
	}

	tokenRepo := identity.NewBoltRepository(db)
	// Create the MicrosoftAuthenticator.
	msAuth := microsoft.NewMicrosoftAuthenticator()

	// Create a registry for Identity providers.
	authRegistry := identity.NewRegistry(tokenRepo, cliLog)
	authRegistry.RegisterAuthenticator("microsoft", msAuth) // Register Authenticator

	// Create and register the MicrosoftAuthorizer.
	msAuthorizer := microsoft.NewMicrosoftAuthorizer(tokenRepo)
	authRegistry.RegisterAuthorizer("microsoft", msAuthorizer)

	c.identityService = authRegistry // Assigning the registry to identityService

	return nil
}

func (c *DefaultContainer) initDriveServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")

	// c.manager is the VFS.
	c.drive = drive.NewDefaultService(c.manager.(*registry.VFS), cliLog)

	return nil
}

func (c *DefaultContainer) initVFSServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")
	yamlSvc := config.NewConfigService(c.profile, cliLog)
	c.config = yamlSvc
	c.mounts = mount.NewMountService(c.config)

	vfs := registry.NewVFS(c.identityService)

	// Load mount configurations and register them with the VFS.
	mounts, err := c.mounts.ListMounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list mounts: %w", err)
	}

	for _, m := range mounts {
		var backend fs.Backend
		switch m.Type {
		case "onedrive":
			backend = onedrive.NewBackend(m.Options)
		case "local":
			backend = local.NewBackend(m.Path)
		default:
			cliLog.Warn("unsupported backend type", logger.String("type", m.Type))
			continue
		}
		vfs.Mount(m.Path, backend)
	}

	c.manager = vfs
	c.uriFactory = registry.NewURIFactory(vfs)
	return nil
}

func (c *DefaultContainer) Logger() logger.Service           { return c.logger }
func (c *DefaultContainer) Config() config.Service           { return c.config }
func (c *DefaultContainer) Mounts() mount.Service            { return c.mounts }
func (c *DefaultContainer) Identity() identity.Service       { return c.identityService }
func (c *DefaultContainer) Profile() profile.Service         { return c.profile }
func (c *DefaultContainer) FS() registry.Service             { return c.manager }
func (c *DefaultContainer) Environment() environment.Service { return c.environment }
func (c *DefaultContainer) Editor() editor.Service           { return c.editor }
func (c *DefaultContainer) Drive() drive.Service             { return c.drive }
func (c *DefaultContainer) URIFactory() *registry.URIFactory {
	return c.uriFactory
}
