package providers

import (
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// Factory is a function that creates a new instance of a filesystem provider.
type Factory func(deps Dependencies) (fs.Service, error)

// Descriptor describes a filesystem provider and how to instantiate it.
type Descriptor struct {
	Name    string
	Factory Factory
}

// Dependencies defines the interface for services that providers might need during initialization.
// This allows providers to be decoupled from the concrete DI container.
type Dependencies interface {
	Logger() logger.Logger
	Get(key string) (any, bool)
}

var (
	registry = make(map[string]Descriptor)
)

// Register adds a provider descriptor to the global registry.
func Register(desc Descriptor) {
	registry[desc.Name] = desc
}

// Get returns the descriptor for a given provider name.
func Get(name string) (Descriptor, bool) {
	desc, ok := registry[name]
	return desc, ok
}

// RegisteredNames returns a list of all registered provider names.
func RegisteredNames() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
