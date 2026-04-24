package formatting

import (
	"reflect"
	"sync"
)

// Registry manages the mapping between types and their presentation definitions.
type Registry struct {
	mu     sync.RWMutex
	tables map[reflect.Type][]Column
}

// GlobalRegistry is the default registry for formatting definitions.
var GlobalRegistry = &Registry{
	tables: make(map[reflect.Type][]Column),
}

// RegisterTable adds a table definition for the given type.
func (r *Registry) RegisterTable(t reflect.Type, cols []Column) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tables[t] = cols
}

// GetTable retrieves the table definition for the given type.
func (r *Registry) GetTable(t reflect.Type) ([]Column, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cols, ok := r.tables[t]
	return cols, ok
}
