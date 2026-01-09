package di

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/app"
	clientservice "github.com/michaeldcanady/go-onedrive/internal/app/client_service"
	credentialservice "github.com/michaeldcanady/go-onedrive/internal/app/credential_service"
	profileservice "github.com/michaeldcanady/go-onedrive/internal/app/profile_service"
	"github.com/michaeldcanady/go-onedrive/internal/cache/fsstore"
	jsoncodec "github.com/michaeldcanady/go-onedrive/internal/cache/json_codex"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"go.uber.org/zap"
)

type Container struct {
	Ctx                context.Context
	Config             config.Config
	Logger             logging.Logger
	ProfileService     ProfileService
	CredentialService  CredentialService
	GraphClientService Clienter
	DriveService       *app.DriveService
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
	logger, err := initializeLogger(cfg.GetLoggingConfig())
	if err != nil {
		return nil, err
	}
	c.Logger = logger

	// profile + credentials
	store := fsstore.New(cfg.GetAuthenticationConfig().GetProfileCache())
	codec := jsoncodec.New()
	c.ProfileService = profileservice.New(store, codec, nil, c.Logger)
	c.CredentialService = credentialservice.New(c.ProfileService, nil, c.Logger)

	// graph client
	c.GraphClientService = clientservice.New(c.CredentialService, nil, c.Logger)

	// drive
	c.DriveService = app.NewDriveService(c.GraphClientService)

	return c, nil
}
