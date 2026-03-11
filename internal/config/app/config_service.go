package app

import (
	"context"
	"strings"

	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal/config/domain"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
)

type ConfigService struct {
	paths  map[string]string
	loader domainconfig.Loader
	log    domainlogger.Logger
}

func New(loader domainconfig.Loader, log domainlogger.Logger) *ConfigService {
	return &ConfigService{
		paths:  make(map[string]string),
		loader: loader,
		log:    log,
	}
}

// ───────────────────────────────────────────────────────────────────────────────
// Event Taxonomy (configdomain.service)
// ───────────────────────────────────────────────────────────────────────────────

const (
	eventConfigGetStart         = "configdomain.get.start"
	eventConfigGetLoadStart     = "configdomain.get.load.start"
	eventConfigGetLoadSuccess   = "configdomain.get.load.success"
	eventConfigGetLoadFailure   = "configdomain.get.load.failure"
	eventConfigGetNotRegistered = "configdomain.get.not_registered"
	eventConfigGetPathMissing   = "configdomain.get.path_missing"
)

func (s *ConfigService) AddPath(name, path string) error {
	if _, exists := s.paths[name]; exists {
		return domainconfig.ErrAlreadyRegistered
	}
	s.paths[name] = path
	return nil
}

func (s *ConfigService) getDefaultConfig() domainconfig.Configuration {
	return domainconfig.Configuration{
		Auth: domainconfig.AuthenticationConfig{
			Type:        domainauth.MethodInteractiveBrowser,
			ClientID:    "6b1e6ec0-ad93-4175-a0e0-84c02e13f206",
			TenantID:    "common",
			RedirectURI: "http://localhost:8400",
		},
	}
}

func (s *ConfigService) GetConfiguration(ctx context.Context, name string) (domainconfig.Configuration, error) {
	if err := ctx.Err(); err != nil {
		return domainconfig.Configuration{}, err
	}

	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		domainlogger.String("correlation_id", correlationID),
		domainlogger.String("config_name", name),
	)

	log.Info("starting configuration retrieval",
		domainlogger.String("event", eventConfigGetStart),
	)

	path, ok := s.paths[name]
	if !ok {
		log.Error("configuration name not registered",
			domainlogger.String("event", eventConfigGetNotRegistered),
		)
		return domainconfig.Configuration{}, domainconfig.ErrNotRegistered
	}

	if strings.TrimSpace(path) == "" {
		log.Error("registered configuration path is empty",
			domainlogger.String("event", eventConfigGetPathMissing),
		)
		return domainconfig.Configuration{}, domainconfig.ErrPathMissing
	}

	log.Info("loading configuration from disk",
		domainlogger.String("event", eventConfigGetLoadStart),
		domainlogger.String("path", path),
	)

	loadedCfg, err := s.loader.Load(path)
	if err != nil {
		log.Error("failed to load configuration from disk",
			domainlogger.String("event", eventConfigGetLoadFailure),
			domainlogger.Error(err),
		)
		return domainconfig.Configuration{}, err
	}

	log.Info("configuration loaded successfully",
		domainlogger.String("event", eventConfigGetLoadSuccess),
	)

	return loadedCfg, nil
}
