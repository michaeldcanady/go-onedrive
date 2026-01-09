package di

import (
	"context"
	"fmt"

	clientservice "github.com/michaeldcanady/go-onedrive/internal/app/client_service"
	credentialservice "github.com/michaeldcanady/go-onedrive/internal/app/credential_service"
	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/drive_service"
	profileservice "github.com/michaeldcanady/go-onedrive/internal/app/profile_service"
	profileservice2 "github.com/michaeldcanady/go-onedrive/internal/app/profile_service2"
	"github.com/michaeldcanady/go-onedrive/internal/cache/fsstore"
	jsoncodec "github.com/michaeldcanady/go-onedrive/internal/cache/json_codex"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"go.uber.org/zap"
)

type Container struct {
	Ctx    context.Context
	Config config.Config
	Logger logging.Logger
	// Deprecated: use ProfileService2 instead.
	// ProfileService is the profile service.
	ProfileService     ProfileService
	ProfileService2    ProfileService2
	CredentialService  CredentialService
	GraphClientService Clienter
	DriveService       ChildrenIterator
	EventBus           *event.InMemoryBus
}

func initializeLogger(logCfg config.LoggingConfig) (logging.Logger, error) {
	if logCfg == nil {
		// TODO: apply default logger config
		return nil, ErrMissingLoggingConfig
	}

	cfg := zap.NewProductionConfig()

	switch logCfg.GetLevel() {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("unknown logging level: %s", logCfg.GetLevel())
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return logging.NewZapLoggerAdapter(zapLogger), nil
}

func NewContainer(ctx context.Context, cfg config.Config) (*Container, error) {
	c := &Container{Ctx: ctx, Config: cfg}

	// logger
	logger, _ := initializeLogger(cfg.GetLoggingConfig())
	c.Logger = logger

	// event bus
	bus := event.NewInMemoryBus(logger)
	c.EventBus = bus

	// services
	store := fsstore.New(cfg.GetAuthenticationConfig().GetProfileCache())
	codec := jsoncodec.New()

	c.ProfileService = profileservice.New(store, codec, bus, logger)
	c.ProfileService2 = profileservice2.New(store, bus, logger)
	c.CredentialService = credentialservice.New(c.ProfileService, bus, logger)
	c.GraphClientService = clientservice.New(c.CredentialService, bus, logger)
	c.DriveService = driveservice.New(c.GraphClientService, bus, logger)

	// wiring listeners
	bus.Subscribe(profileservice.ProfileClearedTopic,
		event.ListenerFunc(func(ctx context.Context, evt event.Topicer) error {
			_, err := c.CredentialService.LoadCredential(ctx, nil)
			return err
		}),
	)

	bus.Subscribe(profileservice2.ProfileDeletedEventTopic,
		event.ListenerFunc(func(ctx context.Context, evt event.Topicer) error {
			_, err := c.CredentialService.LoadCredential(ctx, nil)
			return err
		}),
	)

	bus.Subscribe(profileservice2.ProfileUpdatedEventTopic,
		event.ListenerFunc(func(ctx context.Context, evt event.Topicer) error {
			evt2, ok := evt.(*profileservice2.ProfileEvent)
			if !ok {
				return fmt.Errorf("unexpected event type: %T", evt)
			}

			if evt2.Old() != evt2.Profile() {
				_, err := c.CredentialService.LoadCredential(ctx, evt2.Profile())
				return err
			}
			return nil
		}),
	)

	bus.Subscribe(credentialservice.CredentialLoadedTopic,
		event.ListenerFunc(func(ctx context.Context, evt event.Topicer) error {
			_, err := c.GraphClientService.Client(ctx)
			return err
		}),
	)

	return c, nil
}
