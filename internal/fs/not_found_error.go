package fs

import "fmt"

type NotFoundError struct {
	_ype    string
	Subject string
}

func NewNotFoundError(_ype, subject string) *NotFoundError {
	return &NotFoundError{_ype: _ype, Subject: subject}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s '%s' not found", e._ype, e.Subject)
}

// UnregisteredProvider indicates that a provider is not registered in the registry.
type UnregisteredProvider struct {
	Provider string
}

// NewUnregisteredProvider creates a new UnregisteredProvider error for the given provider name.
func NewUnregisteredProvider(provider string) *UnregisteredProvider {
	return &UnregisteredProvider{Provider: provider}
}

// Error returns a formatted error message indicating that the provider is not registered.
func (e *UnregisteredProvider) Error() string {
	return fmt.Sprintf("provider '%s' is not registered", e.Provider)
}
