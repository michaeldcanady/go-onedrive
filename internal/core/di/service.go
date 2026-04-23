package di
import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	registry "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"

	"github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	"github.com/michaeldcanady/go-onedrive/internal/features/editor/domain"
	"github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
	"github.com/michaeldcanady/go-onedrive/internal/features/storage"
	"github.com/michaeldcanady/go-onedrive/internal/features/storage/backend/local"
	"github.com/michaeldcanady/go-onedrive/internal/features/storage/backend/onedrive"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger/zap"
)

// DefaultContainer provides a concrete implementation of the Container interface.
// It orchestrates the lifecycle and wiring of core application services.
type DefaultContainer struct {
	logger          logger.Service
	config          config.Service
	mounts          mount.Service
	identityService identity.Service
	profile         profile.Service
	manager         registry.Service
	environment     environment.Service
	editor          editor.Service
	drive           drive.Service
	uriFactory      *registry.URIFactory
	storageService  storage.Service
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
    c.storageService = storage.NewDefaultService()

	return nil
}

func (c *DefaultContainer) initIdentityServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")

	stateDir, err := c.environment.StateDir()
	if err != nil {
		return fmt.Errorf("failed to get state directory: %w", err)
	}
	dbPath := filepath.Join(stateDir, "identity.db")
	db, err := c.storageService.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open identity database: %w", err)
	}

	tokenRepo := identity.NewBoltRepository(db)
	msAuth := microsoft.NewMicrosoftAuthenticator()

	authRegistry := identity.NewRegistry(tokenRepo, cliLog)
	authRegistry.RegisterAuthenticator("microsoft", msAuth)

	msAuthorizer := microsoft.NewMicrosoftAuthorizer(tokenRepo)
	authRegistry.RegisterAuthorizer("microsoft", msAuthorizer)

	c.identityService = authRegistry

	return nil
}

type driveLogger struct {
	l logger.Logger
}

func (l *driveLogger) Debug(msg string, fields ...logger.Field) {
	l.l.Debug(msg, fields...)
}

func (l *driveLogger) Error(msg string, fields ...logger.Field) {
	l.l.Error(msg, fields...)
}

func (c *DefaultContainer) initDriveServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")

	c.drive = drive.NewDefaultService(c.manager.(*registry.VFS), &driveLogger{l: cliLog})

	return nil
}

type mountConfigAdapter struct {
	svc config.Service
}

func (a *mountConfigAdapter) GetMounts(ctx context.Context) ([]mount.MountConfig, error) {
	cfg, err := a.svc.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	var mounts []mount.MountConfig
	for _, m := range cfg.Mounts {
		mounts = append(mounts, mount.MountConfig{
			Path:       m.Path,
			Type:       m.Type,
			IdentityID: m.IdentityID,
			Options:    m.Options,
		})
	}
	return mounts, nil
}

func (a *mountConfigAdapter) SaveMounts(ctx context.Context, mounts []mount.MountConfig) error {
	cfg, err := a.svc.GetConfig(ctx)
	if err != nil {
		return err
	}
	var newMounts []config.MountConfig
	for _, m := range mounts {
		newMounts = append(newMounts, config.MountConfig{
			Path:       m.Path,
			Type:       m.Type,
			IdentityID: m.IdentityID,
			Options:    m.Options,
		})
	}
	cfg.Mounts = newMounts
	return a.svc.SaveConfig(ctx, cfg)
}

func (c *DefaultContainer) initVFSServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")
	yamlSvc := config.NewConfigService(c.profile, cliLog)
	c.config = yamlSvc
	c.mounts = mount.NewMountService(&mountConfigAdapter{svc: yamlSvc})
	c.mounts.RegisterValidator("onedrive", onedrive.NewBackend(nil))

	vfs := registry.NewVFS(c.identityService)

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
