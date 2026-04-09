package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"gopkg.in/yaml.v3"
)

// YAMLService is an implementation of the Service interface that loads configuration from YAML files.
type YAMLService struct {
	// resolver is used to look up configuration paths if not in the override map.
	resolver profile.PathResolver
	// state is used for looking up configuration path overrides.
	state state.Service
	// log is the logger used for reporting configuration events.
	log logger.Logger
}

// NewYAMLService creates a new instance of YAMLService.
func NewYAMLService(resolver profile.PathResolver, state state.Service, log logger.Logger) *YAMLService {
	return &YAMLService{
		resolver: resolver,
		state:    state,
		log:      log,
	}
}

// GetPath retrieves the configuration file path for the active profile.
func (s *YAMLService) GetPath(ctx context.Context) (string, bool) {
	l := s.log.WithContext(ctx)

	// Check for transient override in state service first
	if path, err := s.state.Get(state.KeyConfigOverride); err == nil && path != "" {
		l.Debug("configuration override detected", logger.String("path", path))
		return path, true
	}

	// Fetch current profile from state
	profile, err := s.state.Get(state.KeyProfile)
	if err != nil {
		l.Warn("failed to fetch active profile from state", logger.Error(err))
		return "", false
	}

	// Then check resolver if available
	if s.resolver != nil {
		path, err := s.resolver.ResolvePath(ctx, profile)
		if err != nil {
			l.Debug("resolver could not find path for profile", logger.String("profile", profile), logger.Error(err))
			return "", false
		}
		if path != "" {
			l.Debug("configuration path resolved", logger.String("profile", profile), logger.String("path", path))
			return path, true
		}
	}

	l.Warn("no configuration path could be resolved", logger.String("profile", profile))
	return "", false
}

// GetConfig reads and unmarshals the YAML configuration for the active profile.
func (s *YAMLService) GetConfig(ctx context.Context) (Config, error) {
	l := s.log.WithContext(ctx)
	path, ok := s.GetPath(ctx)

	if !ok || path == "" {
		l.Info("using default configuration (no path resolved)")
		return s.defaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			l.Info("configuration file not found, using defaults", logger.String("path", path))
			return s.defaultConfig(), nil
		}
		return Config{}, errors.NewAppError(
			errors.CodeReadError,
			err,
			"could not read the configuration file",
			fmt.Sprintf("Check if the file exists at '%s' and you have read permissions.", path),
		)
	}

	cfg := s.defaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, errors.NewAppError(
			errors.CodeInvalidConfig,
			err,
			"the configuration file format is invalid",
			"Ensure the YAML format is correct and matches the supported schema.",
		)
	}

	l.Debug("configuration loaded successfully", logger.String("path", path))
	return cfg, nil
}

// SaveConfig writes the provided configuration to the YAML file associated with the active profile.
func (s *YAMLService) SaveConfig(ctx context.Context, cfg Config) error {
	l := s.log.WithContext(ctx)
	path, ok := s.GetPath(ctx)

	if !ok || path == "" {
		return errors.NewAppError(
			errors.CodeNotFound,
			nil,
			"no configuration path could be resolved",
			"Check if an active profile is set. Use 'odc profile list' to see available profiles.",
		)
	}

	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return errors.NewAppError(
			errors.CodeWriteError,
			err,
			"could not create configuration directory",
			fmt.Sprintf("Ensure you have write permissions to the directory '%s'.", filepath.Dir(path)),
		)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.NewAppError(
			errors.CodeInternal,
			err,
			"failed to serialize configuration",
			"This is likely an internal error. Please report it if it persists.",
		)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return errors.NewAppError(
			errors.CodeWriteError,
			err,
			"could not write the configuration file",
			fmt.Sprintf("Ensure you have write permissions to '%s'.", path),
		)
	}

	l.Info("configuration saved successfully", logger.String("path", path))
	return nil
}

// defaultConfig returns the fallback configuration used when no file is found.
func (s *YAMLService) defaultConfig() Config {
	return Config{
		Auth: AuthenticationConfig{
			Provider:    AuthProviderMicrosoft,
			ClientID:    "6b1e6ec0-ad93-4175-a0e0-84c02e13f206",
			TenantID:    "common",
			RedirectURI: "http://localhost:8400",
		},
		Logging: LoggingConfig{
			Level: logger.LevelInfo,
		},
	}
}
