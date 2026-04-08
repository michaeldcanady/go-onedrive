package fs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Registry manages a collection of filesystem providers.
type Registry struct {
	mu sync.RWMutex
	// providers maps provider names to their implementations.
	providers map[string]Service
	// state is the service for tracking the active provider and aliases.
	state state.Service
	// alias is the drive alias management service.
	alias alias.Service
	// logger is the logger instance.
	logger logger.Logger
}

// NewRegistry initializes a new instance of the Registry.
func NewRegistry(state state.Service, alias alias.Service, log logger.Logger) *Registry {
	return &Registry{
		providers: make(map[string]Service),
		state:     state,
		alias:     alias,
		logger:    log,
	}
}

var (
	ErrProviderNotRegistered  = errors.New("provider not registered")
	ErrUnknownProviderOrAlias = errors.New("unknown provider or alias")
)

// Register associates a service with its provider name.
func (r *Registry) Register(provider string, svc Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Wrap the service with the validation decorator
	r.providers[provider] = NewValidationDecorator(svc, r.logger)
}

// Get retrieves the service for a given provider name.
func (r *Registry) Get(provider string) (Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.providers[provider]
	if !ok {
		return nil, ErrProviderNotRegistered
	}
	return svc, nil
}

// Resolve identifies the provider for a path based on its prefix (e.g., "onedrive:path").
// If no prefix is present, it defaults to "onedrive".
func (r *Registry) Resolve(ctx context.Context, path string) (Service, string, error) {
	log := r.logger.WithContext(ctx).With(logger.String("path", path))
	prefix, rest, found := strings.Cut(path, ":")
	if !found {
		// No prefix, default to onedrive
		log.Debug("no provider prefix found, using default", logger.String("default", DefaultProviderPrefix))
		p, err := r.Get(DefaultProviderPrefix)
		if err != nil {
			if !errors.Is(err, ErrProviderNotRegistered) {
				log.Error("error retrieving default provider", logger.String("provider", DefaultProviderPrefix), logger.Error(err))
				return nil, "", err
			}
			// If the default provider is not registered, this shouldn't happen
			log.Warn("default provider not registered, checking if path is an alias", logger.String("path", path))
			return nil, "", ErrUnknownProviderOrAlias
		}
		return p, path, nil
	}

	// Check if prefix is a registered provider
	p, err := r.Get(prefix)
	if err != nil {
		if !errors.Is(err, ErrProviderNotRegistered) {
			log.Error("error retrieving provider", logger.String("prefix", prefix), logger.Error(err))
			return nil, "", err
		}
		log.Debug("prefix is not a registered provider, checking aliases", logger.String("prefix", prefix))
		// Check if prefix is an alias
		driveID, err := r.alias.GetDriveIDByAlias(prefix)
		if err != nil {
			if !errors.Is(err, alias.ErrDriveIDNotFound) {
				log.Error("error retrieving drive ID for alias", logger.String("alias", prefix), logger.Error(err))
				return nil, "", err
			}
			log.Error("unknown provider or alias", logger.String("prefix", prefix))
			return nil, "", ErrUnknownProviderOrAlias
		}
		// If it's an alias, use the default provider (onedrive) and prepend the drive ID
		defaultProvider, err := r.Get(DefaultProviderPrefix)
		if err != nil {
			log.Error("error retrieving default provider", logger.String("provider", DefaultProviderPrefix), logger.Error(err))
			return nil, "", err
		}
		// Construct path with drive ID for the default provider
		rest = fmt.Sprintf("%s:%s", driveID, rest)
		log.Debug("resolved alias to drive ID", logger.String("alias", prefix), logger.String("drive_id", driveID))
		return defaultProvider, rest, nil
	}
	// If it's a registered provider, use it directly
	log.Debug("resolved path to registered provider", logger.String("provider", prefix))
	return p, rest, nil
}

func (r *Registry) RegisteredNames() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names, nil
}

// DefaultProviderPrefix is the default provider to use when no prefix is specified.
const DefaultProviderPrefix = "onedrive"

// Cleanup removes all providers from the registry.
func (r *Registry) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = make(map[string]Service)
}
